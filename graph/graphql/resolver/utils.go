package resolver

import (
	"context"
	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/predicate"
	"github.com/net-auto/resourceManager/ent/resourcepool"
	resourcePool "github.com/net-auto/resourceManager/ent/resourcepool"
	"github.com/net-auto/resourceManager/ent/tag"
	"github.com/net-auto/resourceManager/graph/graphql/model"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func decodeCursor(cursorAsString * string) (*ent.Cursor, error) {
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
	if pool.PoolType == resourcepool.PoolTypeAllocating {

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
