package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/allocationstrategy"
	"github.com/net-auto/resourceManager/ent/predicate"
	resourcePool "github.com/net-auto/resourceManager/ent/resourcepool"
	tagWhere "github.com/net-auto/resourceManager/ent/tag"
	"github.com/net-auto/resourceManager/graph/graphql/generated"
	"github.com/net-auto/resourceManager/graph/graphql/model"
	p "github.com/net-auto/resourceManager/pools"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func (r *mutationResolver) CreateTag(ctx context.Context, tag string) (*ent.Tag, error) {
	var client = r.ClientFrom(ctx)
	tagEnt, err := client.Tag.Create().SetTag(tag).Save(ctx)
	if err != nil {
		return nil, gqlerror.Errorf("Unable to create tag: %v", err)
	}
	return tagEnt, err
}

func (r *mutationResolver) UpdateTag(ctx context.Context, tagID int, tag string) (*ent.Tag, error) {
	var client = r.ClientFrom(ctx)
	tagEnt, err := client.Tag.UpdateOneID(tagID).SetTag(tag).Save(ctx)
	if err != nil {
		return nil, gqlerror.Errorf("Unable to update tag: %v", err)
	}
	return tagEnt, nil
}

func (r *mutationResolver) DeleteTag(ctx context.Context, tagID int) (*ent.Tag, error) {
	var client = r.ClientFrom(ctx)
	err := client.Tag.DeleteOneID(tagID).Exec(ctx)
	if err != nil {
		return nil, gqlerror.Errorf("Unable to delete tag: %v", err)
	}
	return nil, nil
}

func (r *mutationResolver) TagPool(ctx context.Context, tagID int, poolID int) (*ent.Tag, error) {
	var client = r.ClientFrom(ctx)
	tag, err := client.Tag.UpdateOneID(tagID).AddPoolIDs(poolID).Save(ctx)
	if err != nil {
		return nil, gqlerror.Errorf("Unable to tag pool: %v", err)
	}
	return tag, nil
}

func (r *mutationResolver) CreateAllocationStrategy(ctx context.Context, name string, script string, lang allocationstrategy.Lang) (*ent.AllocationStrategy, error) {
	var client = r.ClientFrom(ctx)
	strat, err := client.AllocationStrategy.Create().SetName(name).SetScript(script).SetLang(lang).Save(ctx)
	if err != nil {
		return nil, gqlerror.Errorf("Unable to create strategy: %v", err)
	}
	return strat, err
}

func (r *mutationResolver) DeleteAllocationStrategy(ctx context.Context, allocationStrategyID int) (*ent.AllocationStrategy, error) {
	var client = r.ClientFrom(ctx)
	if strat, err := client.AllocationStrategy.Query().
		Where(allocationstrategy.ID(allocationStrategyID)).
		Only(ctx); err != nil {
		return nil, gqlerror.Errorf("Unable to delete strategy: %v", err)
	} else {

		if dependentPools, err := strat.QueryPools().All(ctx); len(dependentPools) > 0 && err != nil {
			return nil, gqlerror.Errorf("Unable to delete, Allocation strategy is still in use")
		}

		if err := client.AllocationStrategy.DeleteOneID(allocationStrategyID).Exec(ctx); err != nil {
			return nil, err
		}
		return strat, nil
	}
}

func (r *mutationResolver) TestAllocationStrategy(ctx context.Context, allocationStrategyID int, resourcePool model.ResourcePoolInput, currentResources []*model.ResourceInput, userInput map[string]interface{}) (map[string]interface{}, error) {
	var client = r.ClientFrom(ctx)
	strat, err := client.AllocationStrategy.Query().
		Where(allocationstrategy.ID(allocationStrategyID)).
		Only(ctx)
	if err != nil {
		return nil, gqlerror.Errorf("Unable to get strategy: %v", err)
	}
	// TODO keep just single instance
	wasmer, err := p.NewWasmerUsingEnvVars()
	if err != nil {
		return nil, gqlerror.Errorf("Unable to create scripting engine: %v", err)
	}

	parsedOutputFromStrat, stdErr, err := p.InvokeAllocationStrategy(wasmer, strat, userInput, resourcePool, currentResources)
	if err != nil {
		return nil, gqlerror.Errorf("Error while running the script: %v", err)
	}
	result := make(map[string]interface{})
	result["stdout"] = parsedOutputFromStrat
	result["stderr"] = stdErr
	return result, nil
}

func (r *mutationResolver) ClaimResource(ctx context.Context, poolID int, userInput map[string]interface{}) (*ent.Resource, error) {
	pool, err := p.ExistingPoolFromId(ctx, r.ClientFrom(ctx), poolID)
	if err != nil {
		return nil, gqlerror.Errorf("Unable to claim resource: %v", err)
	}

	if res, err := pool.ClaimResource(userInput); err != nil {
		return nil, gqlerror.Errorf("Unable to claim resource: %v", err)
	} else {
		return res, nil
	}
}

func (r *mutationResolver) FreeResource(ctx context.Context, input map[string]interface{}, poolID int) (string, error) {
	pool, err := p.ExistingPoolFromId(ctx, r.ClientFrom(ctx), poolID)
	if err != nil {
		return "", gqlerror.Errorf("Unable to free resource: %v", err)
	}
	err = pool.FreeResource(input)
	if err == nil {
		return "Resource freed successfully", nil
	}

	return "", gqlerror.Errorf("Unable to free resource: %v", err)
}

func (r *mutationResolver) CreateSetPool(ctx context.Context, resourceTypeID int, poolName string, poolDealocationSafetyPeriod int, poolValues []map[string]interface{}) (*ent.ResourcePool, error) {
	var client = r.ClientFrom(ctx)

	resType, err := client.ResourceType.Get(ctx, resourceTypeID)
	if err != nil {
		return nil, gqlerror.Errorf("Unable to create pool: %v", err)
	}
	_, rp, err := p.NewSetPoolWithMeta(ctx, client, resType, p.ToRawTypes(poolValues), poolName, poolDealocationSafetyPeriod)
	if err != nil {
		return nil, gqlerror.Errorf("Unable to create pool: %v", err)
	}
	return rp, nil
}

func (r *mutationResolver) CreateSingletonPool(ctx context.Context, resourceTypeID int, poolName string, poolValues []map[string]interface{}) (*ent.ResourcePool, error) {
	var client = r.ClientFrom(ctx)

	resType, _ := client.ResourceType.Get(ctx, resourceTypeID)
	if len(poolValues) == 1 {
		_, rp, err := p.NewSingletonPoolWithMeta(ctx, client, resType, p.ToRawTypes(poolValues)[0], poolName)
		return rp, gqlerror.Errorf("Cannot create singleton pool: %v", err)
	} else {
		return nil, gqlerror.Errorf("Cannot create singleton pool, no resource provided")
	}
}

func (r *mutationResolver) CreateAllocatingPool(ctx context.Context, resourceTypeID int, poolName string, allocationStrategyID int, poolDealocationSafetyPeriod int) (*ent.ResourcePool, error) {
	var client = r.ClientFrom(ctx)

	resType, errRes := client.ResourceType.Get(ctx, resourceTypeID)
	if errRes != nil {
		return nil, gqlerror.Errorf("Unable to create pool: %v", errRes)
	}
	allocationStrat, errAlloc := client.AllocationStrategy.Get(ctx, allocationStrategyID)
	if errAlloc != nil {
		return nil, gqlerror.Errorf("Unable to create pool: %v", errAlloc)
	}

	_, rp, err := p.NewAllocatingPoolWithMeta(ctx, client, resType, allocationStrat, poolName, poolDealocationSafetyPeriod)
	if err != nil {
		return nil, gqlerror.Errorf("Unable to create pool: %v", err)
	}
	return rp, err
}

func (r *mutationResolver) DeleteResourcePool(ctx context.Context, resourcePoolID int) (string, error) {
	client := r.ClientFrom(ctx)
	pool, err := p.ExistingPoolFromId(ctx, client, resourcePoolID)

	if err != nil {
		return "", gqlerror.Errorf("Unable to delete pool: %v", err)
	}

	// Do not allow removing pools with allocated resources
	allocatedResources, err2 := pool.QueryResources()
	if len(allocatedResources) > 0 || err2 != nil {
		return "", gqlerror.Errorf("Unable to delete pool, pool has allocated resources, deallocate those first")
	}

	if err := client.ResourcePool.DeleteOneID(resourcePoolID).Exec(ctx); err != nil {
		return "", gqlerror.Errorf("Unable to delete pool: %v", err)
	} else {
		return "ok", nil
	}
}

func (r *mutationResolver) CreateResourceType(ctx context.Context, resourceName string, resourceProperties map[string]interface{}) (*ent.ResourceType, error) {
	var client = r.ClientFrom(ctx)

	var propertyTypes []*ent.PropertyType
	for propName, rawPropType := range resourceProperties {
		var propertyType, err = p.CreatePropertyType(ctx, client, propName, rawPropType)
		if err != nil {
			return nil, gqlerror.Errorf("Unable to create resource type: %v", err)
		}
		propertyTypes = append(propertyTypes, propertyType)
	}

	resType, err2 := client.ResourceType.Create().
		SetName(resourceName).
		AddPropertyTypes(propertyTypes...).
		Save(ctx)
	if err2 != nil {
		return nil, gqlerror.Errorf("Unable to create resource type: %v", err2)
	}

	return resType, nil
}

func (r *mutationResolver) DeleteResourceType(ctx context.Context, resourceTypeID int) (string, error) {
	client := r.ClientFrom(ctx)
	resourceType, err := client.ResourceType.Get(ctx, resourceTypeID)

	if err != nil {
		return "nil", gqlerror.Errorf("Unable to delete resource type: %v", err)
	}

	pools, err2 := client.ResourceType.QueryPools(resourceType).All(ctx)

	if err2 != nil {
		return "", gqlerror.Errorf("Unable to create resource type: %v", err2)
	}

	if len(pools) > 0 {
		return "", gqlerror.Errorf("Unable to create resource type, there are pools attached to it")
	}

	if err := client.ResourcePool.DeleteOneID(resourceTypeID).Exec(ctx); err == nil {
		return "ok", nil
	} else {
		return "", gqlerror.Errorf("Unable to create resource type: %v", err2)
	}
}

func (r *mutationResolver) UpdateResourceTypeName(ctx context.Context, resourceTypeID int, resourceName string) (*ent.ResourceType, error) {
	var client = r.ClientFrom(ctx)
	if rT, err := client.ResourceType.UpdateOneID(resourceTypeID).SetName(resourceName).Save(ctx); err != nil {
		return nil, gqlerror.Errorf("Unable to update resource type: %v", err)
	} else {
		return rT, nil
	}
}

func (r *propertyTypeResolver) Type(ctx context.Context, obj *ent.PropertyType) (string, error) {
	// Just converts enum to string
	return obj.Type.String(), nil
}

func (r *queryResolver) QueryResource(ctx context.Context, input map[string]interface{}, poolID int) (*ent.Resource, error) {
	pool, err := p.ExistingPoolFromId(ctx, r.ClientFrom(ctx), poolID)
	if err != nil {
		return nil, gqlerror.Errorf("Unable to query resource: %v", err)
	}
	return pool.QueryResource(input)
}

func (r *queryResolver) QueryResources(ctx context.Context, poolID int) ([]*ent.Resource, error) {
	pool, err := p.ExistingPoolFromId(ctx, r.ClientFrom(ctx), poolID)
	if err != nil {
		return nil, gqlerror.Errorf("Unable to query resources: %v", err)
	}
	return pool.QueryResources()
}

func (r *queryResolver) QueryAllocationStrategy(ctx context.Context, allocationStrategyID int) (*ent.AllocationStrategy, error) {
	client := r.ClientFrom(ctx)
	if strats, err := client.AllocationStrategy.Query().Where(allocationstrategy.ID(allocationStrategyID)).Only(ctx); err != nil {
		return nil, gqlerror.Errorf("Unable to query strategy: %v", err)
	} else {
		return strats, nil
	}
}

func (r *queryResolver) QueryAllocationStrategies(ctx context.Context) ([]*ent.AllocationStrategy, error) {
	client := r.ClientFrom(ctx)
	if strats, err := client.AllocationStrategy.Query().All(ctx); err != nil {
		return nil, gqlerror.Errorf("Unable to query strategies: %v", err)
	} else {
		return strats, nil
	}
}

func (r *queryResolver) QueryResourceTypes(ctx context.Context) ([]*ent.ResourceType, error) {
	client := r.ClientFrom(ctx)
	if resourceTypes, err := client.ResourceType.Query().All(ctx); err != nil {
		return nil, gqlerror.Errorf("Unable to query resource types: %v", err)
	} else {
		return resourceTypes, nil
	}
}

func (r *queryResolver) QueryResourcePools(ctx context.Context) ([]*ent.ResourcePool, error) {
	client := r.ClientFrom(ctx)
	if resourcePools, err := client.ResourcePool.Query().All(ctx); err != nil {
		return nil, gqlerror.Errorf("Unable to query resource pools: %v", err)
	} else {
		return resourcePools, nil
	}
}

func (r *queryResolver) SearchPoolsByTags(ctx context.Context, tags *model.TagOr) ([]*ent.ResourcePool, error) {
	var client = r.ClientFrom(ctx)
	var predicateOr predicate.ResourcePool

	// TODO make sure all tags exist

	for _, tagOr := range tags.MatchesAny {
		// Join queries where tag equals to input by AND operation
		predicateAnd := resourcePool.HasTags()
		for _, tagAnd := range tagOr.MatchesAll {
			predicateAnd = resourcePool.And(predicateAnd, resourcePool.HasTagsWith(tagWhere.Tag(tagAnd)))
		}

		// Join multiple AND tag queries with OR
		if predicateOr == nil {
			// If this is the first AND query, use the AND query as a starting point
			predicateOr = predicateAnd
		} else {
			predicateOr = resourcePool.Or(predicateOr, predicateAnd)
		}
	}

	matchedPools, err := client.ResourcePool.Query().Where(predicateOr).All(ctx)
	if err != nil {
		return nil, gqlerror.Errorf("Unable to query pools: %v", err)
	}
	return matchedPools, nil
}

func (r *queryResolver) QueryTags(ctx context.Context) ([]*ent.Tag, error) {
	var client = r.ClientFrom(ctx)
	tags, err := client.Tag.Query().All(ctx)
	if err != nil {
		return nil, gqlerror.Errorf("Unable to query tags: %v", err)
	}
	return tags, nil
}

func (r *resourceResolver) Properties(ctx context.Context, obj *ent.Resource) (map[string]interface{}, error) {
	props, err := obj.QueryProperties().WithType().All(ctx)
	if err != nil {
		return nil, gqlerror.Errorf("Unable to query properties: %v", err)
	}

	if props, err := p.PropertiesToMap(props); err != nil {
		return nil, gqlerror.Errorf("Unable to query properties: %v", err)
	} else {
		return props, nil
	}
}

func (r *resourcePoolResolver) ResourceType(ctx context.Context, obj *ent.ResourcePool) (*ent.ResourceType, error) {
	if es, err := obj.Edges.ResourceTypeOrErr(); !ent.IsNotLoaded(err) {
		return es, err
	}
	return obj.QueryResourceType().Only(ctx)
}

func (r *resourcePoolResolver) Resources(ctx context.Context, obj *ent.ResourcePool) ([]*ent.Resource, error) {
	if es, err := obj.Edges.ClaimsOrErr(); !ent.IsNotLoaded(err) {
		return es, err
	}
	return obj.QueryClaims().All(ctx)
}

func (r *resourcePoolResolver) Tags(ctx context.Context, obj *ent.ResourcePool) ([]*ent.Tag, error) {
	if es, err := obj.Edges.TagsOrErr(); !ent.IsNotLoaded(err) {
		return es, err
	}
	return obj.QueryTags().All(ctx)
}

func (r *resourcePoolResolver) AllocationStrategy(ctx context.Context, obj *ent.ResourcePool) (*ent.AllocationStrategy, error) {
	if obj.PoolType != resourcePool.PoolTypeAllocating {
		return nil, nil
	}
	if es, err := obj.Edges.AllocationStrategyOrErr(); !ent.IsNotLoaded(err) {
		return es, err
	}
	return obj.QueryAllocationStrategy().Only(ctx)
}

func (r *resourceTypeResolver) PropertyTypes(ctx context.Context, obj *ent.ResourceType) ([]*ent.PropertyType, error) {
	if es, err := obj.Edges.PropertyTypesOrErr(); !ent.IsNotLoaded(err) {
		return es, err
	}
	return obj.QueryPropertyTypes().All(ctx)
}

func (r *resourceTypeResolver) Pools(ctx context.Context, obj *ent.ResourceType) ([]*ent.ResourcePool, error) {
	if es, err := obj.Edges.PoolsOrErr(); !ent.IsNotLoaded(err) {
		return es, err
	}
	return obj.QueryPools().All(ctx)
}

func (r *tagResolver) Pools(ctx context.Context, obj *ent.Tag) ([]*ent.ResourcePool, error) {
	if es, err := obj.Edges.PoolsOrErr(); !ent.IsNotLoaded(err) {
		return es, err
	}
	return obj.QueryPools().All(ctx)
}

// Mutation returns generated.MutationResolver implementation.
//  Mutation() function removed in favour of resolver.go.Mutation()

// PropertyType returns generated.PropertyTypeResolver implementation.
func (r *Resolver) PropertyType() generated.PropertyTypeResolver { return &propertyTypeResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Resource returns generated.ResourceResolver implementation.
func (r *Resolver) Resource() generated.ResourceResolver { return &resourceResolver{r} }

// ResourcePool returns generated.ResourcePoolResolver implementation.
func (r *Resolver) ResourcePool() generated.ResourcePoolResolver { return &resourcePoolResolver{r} }

// ResourceType returns generated.ResourceTypeResolver implementation.
func (r *Resolver) ResourceType() generated.ResourceTypeResolver { return &resourceTypeResolver{r} }

// Tag returns generated.TagResolver implementation.
func (r *Resolver) Tag() generated.TagResolver { return &tagResolver{r} }

type mutationResolver struct{ *Resolver }
type propertyTypeResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type resourceResolver struct{ *Resolver }
type resourcePoolResolver struct{ *Resolver }
type resourceTypeResolver struct{ *Resolver }
type tagResolver struct{ *Resolver }
