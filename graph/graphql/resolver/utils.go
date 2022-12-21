package resolver

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"fmt"
	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/predicate"
	"github.com/net-auto/resourceManager/ent/resource"
	resourcePool "github.com/net-auto/resourceManager/ent/resourcepool"
	"github.com/net-auto/resourceManager/ent/tag"
	"github.com/net-auto/resourceManager/graph/graphql/model"
	"github.com/net-auto/resourceManager/pools"
	"strconv"

	//"github.com/net-auto/resourceManager/graph/graphql/model"
	log "github.com/net-auto/resourceManager/logging"
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func decodeCursor(cursorAsString *string) (*ent.Cursor, error) {
	if cursorAsString == nil {
		return nil, nil
	}

	result := &ent.Cursor{
		ID:    1,
		Value: nil,
	}
	err := result.UnmarshalGQL(*cursorAsString)

	return result, err
}

func hasParent(currentPool *ent.ResourcePool) bool {
	return currentPool != nil &&
		currentPool.Edges.ParentResource != nil &&
		currentPool.Edges.ParentResource.Edges.Pool != nil
}
func queryPoolWithParent(ctx context.Context, poolID int, client *ent.Client) (*ent.ResourcePool, error) {
	return client.ResourcePool.
		Query().
		Where(resourcePool.ID(poolID)).
		WithParentResource(func(query *ent.ResourceQuery) {
			query.WithPool()
		}).Only(ctx)
}

func createNestedPool(ctx context.Context,
	parentResourceID int,
	client *ent.Client,
	nestedPoolFactory func() (*ent.ResourcePool, error)) (*ent.ResourcePool, error) {

	parentResource, errFindRes := client.Resource.Get(ctx, parentResourceID)
	if errFindRes != nil {
		return nil, gqlerror.Errorf("Cannot create nested pool, Unable to find parent resource: %v", errFindRes)
	}

	pool, errCreatePool := nestedPoolFactory()
	if errCreatePool != nil {
		return nil, errCreatePool
	}

	//create pool properties from parent resource type
	if pool.PoolType == resourcePool.PoolTypeAllocating {

		properties, err := parentResource.QueryProperties().All(ctx)

		if err != nil {
			return nil, gqlerror.Errorf("Cannot retrieve properties, error: %v", err)
		}

		poolProperties, err := client.PoolProperties.Create().AddProperties(properties...).Save(ctx)

		if err != nil {
			return nil, gqlerror.Errorf("Cannot create pool properties, error: %v", err)
		}

		_, err = pool.Update().SetPoolProperties(poolProperties).Save(ctx)

		if err != nil {
			return nil, gqlerror.Errorf("Unable to set pool properties on the given pool, error: %v", err)
		}
	}

	pool, errSetParentResource := pool.Update().SetParentResource(parentResource).Save(ctx)
	if errSetParentResource != nil {
		return pool, nil
	}

	return pool, nil
}

type ById []*ent.Resource

func (a ById) Len() int           { return len(a) }
func (a ById) Less(i, j int) bool { return a[i].ID < a[j].ID }
func (a ById) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func resourcePoolTagPredicate(tags *model.TagOr) predicate.ResourcePool {
	var predicateOr predicate.ResourcePool

	for _, tagOr := range tags.MatchesAny {
		// Join queries where tag equals to input by AND operation
		predicateAnd := resourcePool.HasTags()
		for _, tagAnd := range tagOr.MatchesAll {
			predicateAnd = resourcePool.And(predicateAnd, resourcePool.HasTagsWith(tag.Tag(tagAnd)))
		}

		// Join multiple AND tag queries with OR
		if predicateOr == nil {
			// If this is the first AND query, use the AND query as a starting point
			predicateOr = predicateAnd
		} else {
			predicateOr = resourcePool.Or(predicateOr, predicateAnd)
		}
	}
	return predicateOr
}

func createTagsAndTagPool(ctx context.Context, client *ent.Client, rp *ent.ResourcePool, tags []string) error {
	var tagsInDb []*ent.Tag
	for _, newTag := range tags {
		if tagInDb, err := createOrLoadTag(ctx, client, newTag); err != nil {
			return err
		} else {
			tagsInDb = append(tagsInDb, tagInDb)
		}

	}

	if tagPool(ctx, client, rp, tagsInDb) != nil {
		return tagPool(ctx, client, rp, tagsInDb)
	}

	return nil
}

func tagPool(ctx context.Context, client *ent.Client, rp *ent.ResourcePool, tagsInDb []*ent.Tag) error {
	return client.ResourcePool.UpdateOne(rp).AddTags(tagsInDb...).Exec(ctx)
}

func createOrLoadTag(ctx context.Context, client *ent.Client, newTag string) (*ent.Tag, error) {
	tagInDb, err := tagFromDb(ctx, client, newTag)
	if err == nil {
		return tagInDb, nil
	} else if !ent.IsNotFound(err) {
		return nil, err
	}

	return createTag(ctx, client, newTag)
}

func createTag(ctx context.Context, client *ent.Client, newTag string) (*ent.Tag, error) {
	return client.Tag.Create().SetTag(newTag).Save(ctx)
}

func tagFromDb(ctx context.Context, client *ent.Client, tagValue string) (*ent.Tag, error) {
	return client.Tag.Query().Where(tag.Tag(tagValue)).Only(ctx)
}

// QueryResourcesByAltId returns paginate resources if alt Id matches
func QueryResourcesByAltId(ctx context.Context, client *ent.Client, alternativeId map[string]interface{}, poolId *int, first *int,
	last *int, before *string, after *string) (*ent.ResourceConnection, error) {

	afterCursor, errA := decodeCursor(after)
	if errA != nil {
		log.Error(ctx, errA, "Unable to decode after value (\"%s\") for pagination", *after)
		return nil, errA
	}

	beforeCursor, errB := decodeCursor(before)
	if errB != nil {
		log.Error(ctx, errB, "Unable to decode before value (\"%s\") for pagination", *before)
		return nil, errB
	}

	if poolId != nil {
		res, err := client.Resource.Query().
			Where(resource.HasPoolWith(resourcePool.ID(*poolId))).
			Where(func(selector *sql.Selector) {
				for k, v := range alternativeId {
					selector.Where(sqljson.ValueContains("alternate_id", v, sqljson.Path(k)))
				}
			}).
			Paginate(ctx, afterCursor, first, beforeCursor, last)

		if err != nil {
			log.Error(ctx, err, "Unable to retrieve resources with alternative ID %v", alternativeId)
			return nil, gqlerror.Errorf("Unable to query resources: %v", err)
		}

		if res != nil {
			return res, nil
		}
	}

	log.Error(ctx, nil, "There is not such resource with alternative ID %v", alternativeId)
	return nil, errors.New("No such resource with given alternative ID")
}

func ClaimResource(pool pools.Pool, userInput map[string]interface{}, description *string, alternativeId map[string]interface{}) (*ent.Resource, error) {
	input := make(map[string]interface{})

	for key, value := range userInput {
		switch value.(type) {
		case string:
			intVal, intErr := strconv.Atoi(fmt.Sprintf("%v", value))
			if intErr == nil && key == "desiredSize" {
				input[key] = intVal
			} else {
				input[key] = value
			}
		default:
			input[key] = value
			break
		}
	}

	if res, err := pool.ClaimResource(input, description, alternativeId); err != nil {
		return nil, gqlerror.Errorf("Unable to claim resource: %v", err)
	} else {
		return res, nil
	}
}
