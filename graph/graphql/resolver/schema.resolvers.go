package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/allocationstrategy"
	"github.com/net-auto/resourceManager/ent/predicate"
	"github.com/net-auto/resourceManager/ent/propertytype"
	"github.com/net-auto/resourceManager/ent/resource"
	resourcePool "github.com/net-auto/resourceManager/ent/resourcepool"
	"github.com/net-auto/resourceManager/ent/resourcetype"
	tagWhere "github.com/net-auto/resourceManager/ent/tag"
	"github.com/net-auto/resourceManager/graph/graphql/generated"
	"github.com/net-auto/resourceManager/graph/graphql/model"
	p "github.com/net-auto/resourceManager/pools"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func (r *mutationResolver) CreateTag(ctx context.Context, input model.CreateTagInput) (*model.CreateTagPayload, error) {
	var client = r.ClientFrom(ctx)
	tagEnt, err := client.Tag.Create().SetTag(input.TagText).Save(ctx)

	if err != nil {
		return &model.CreateTagPayload{Tag: nil}, gqlerror.Errorf("Unable to create tag: %v", err)
	}

	return &model.CreateTagPayload{Tag: tagEnt}, nil
}

func (r *mutationResolver) UpdateTag(ctx context.Context, input model.UpdateTagInput) (*model.UpdateTagPayload, error) {
	var client = r.ClientFrom(ctx)
	tagEnt, err := client.Tag.UpdateOneID(input.TagID).SetTag(input.TagText).Save(ctx)
	if err != nil {
		return &model.UpdateTagPayload{Tag: nil}, gqlerror.Errorf("Unable to update tag: %v", err)
	}
	return &model.UpdateTagPayload{Tag: tagEnt}, nil
}

func (r *mutationResolver) DeleteTag(ctx context.Context, input model.DeleteTagInput) (*model.DeleteTagPayload, error) {
	var client = r.ClientFrom(ctx)
	err := client.Tag.DeleteOneID(input.TagID).Exec(ctx)
	if err != nil {
		return &model.DeleteTagPayload{TagID: input.TagID}, gqlerror.Errorf("Unable to delete tag: %v", err)
	}

	return &model.DeleteTagPayload{TagID: input.TagID}, nil
}

func (r *mutationResolver) TagPool(ctx context.Context, input model.TagPoolInput) (*model.TagPoolPayload, error) {
	var client = r.ClientFrom(ctx)
	tag, err := client.Tag.UpdateOneID(input.TagID).AddPoolIDs(input.PoolID).Save(ctx)
	if err != nil {
		return &model.TagPoolPayload{Tag: nil}, gqlerror.Errorf("Unable to tag pool: %v", err)
	}
	return &model.TagPoolPayload{Tag: tag}, nil
}

func (r *mutationResolver) UntagPool(ctx context.Context, input model.UntagPoolInput) (*model.UntagPoolPayload, error) {
	var client = r.ClientFrom(ctx)
	tag, err := client.Tag.UpdateOneID(input.TagID).RemovePoolIDs(input.PoolID).Save(ctx)
	if err != nil {
		return &model.UntagPoolPayload{Tag: nil}, gqlerror.Errorf("Unable to un-tag pool: %v", err)
	}
	return &model.UntagPoolPayload{Tag: tag}, nil
}

func (r *mutationResolver) CreateAllocationStrategy(ctx context.Context, input *model.CreateAllocationStrategyInput) (*model.CreateAllocationStrategyPayload, error) {
	var client = r.ClientFrom(ctx)
	strat, err := client.AllocationStrategy.Create().
		SetName(input.Name).
		SetNillableDescription(input.Description).
		SetScript(input.Script).
		SetLang(input.Lang).
		Save(ctx)
	if err != nil {
		return &model.CreateAllocationStrategyPayload{Strategy: nil}, gqlerror.Errorf("Unable to create strategy: %v", err)
	}

	return &model.CreateAllocationStrategyPayload{Strategy: strat}, nil
}

func (r *mutationResolver) DeleteAllocationStrategy(ctx context.Context, input *model.DeleteAllocationStrategyInput) (*model.DeleteAllocationStrategyPayload, error) {
	var client = r.ClientFrom(ctx)
	emptyRetVal := model.DeleteAllocationStrategyPayload{Strategy: nil}
	if strat, err := client.AllocationStrategy.Query().
		Where(allocationstrategy.ID(input.AllocationStrategyID)).
		WithPools().
		Only(ctx); err != nil {
		return &emptyRetVal, gqlerror.Errorf("Unable to delete strategy: %v", err)
	} else {

		if len(strat.Edges.Pools) > 0 {
			return &emptyRetVal, gqlerror.Errorf("Unable to delete, Allocation strategy is still in use")
		}

		if err := client.AllocationStrategy.DeleteOneID(input.AllocationStrategyID).Exec(ctx); err != nil {
			return &emptyRetVal, err
		}

		return &model.DeleteAllocationStrategyPayload{Strategy: strat}, nil
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

	poolPropertiesMaps := make(map[string]interface{})

	for key, element := range resourcePool.PoolProperties {
		poolPropertiesMaps[key] = element
	}

	var functionName string
	if strat.Lang == allocationstrategy.LangPy {
		functionName = "script_fun()"
	} else {
		functionName = "invoke()"
	}

	parsedOutputFromStrat, stdErr, err := p.InvokeAllocationStrategy(wasmer, strat, userInput, resourcePool, currentResources, poolPropertiesMaps, functionName)
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

func (r *mutationResolver) CreateSetPool(ctx context.Context, input model.CreateSetPoolInput) (*model.CreateSetPoolPayload, error) {
	var client = r.ClientFrom(ctx)

	resType, err := client.ResourceType.Get(ctx, input.ResourceTypeID)
	if err != nil {
		return &model.CreateSetPoolPayload{Pool: nil}, gqlerror.Errorf("Unable to create pool: %v", err)
	}
	_, rp, err := p.NewSetPoolWithMeta(ctx, client, resType, p.ToRawTypes(input.PoolValues),
		input.PoolName, input.Description, input.PoolDealocationSafetyPeriod)
	if err != nil {
		return &model.CreateSetPoolPayload{Pool: nil}, gqlerror.Errorf("Unable to create pool: %v", err)
	}
	return &model.CreateSetPoolPayload{Pool: rp}, nil
}

func (r *mutationResolver) CreateNestedSetPool(ctx context.Context, input model.CreateNestedSetPoolInput) (*model.CreateNestedSetPoolPayload, error) {
	var client = r.ClientFrom(ctx)

	pool, err2 := createNestedPool(ctx, input.ParentResourceID, client, func() (*ent.ResourcePool, error) {
		poolInput := model.CreateSetPoolInput{
			ResourceTypeID:              input.ResourceTypeID,
			PoolName:                    input.PoolName,
			Description:                 input.Description,
			PoolDealocationSafetyPeriod: input.PoolDealocationSafetyPeriod,
			PoolValues:                  input.PoolValues,
		}
		createSetPoolPayload, err := r.CreateSetPool(ctx, poolInput)
		if createSetPoolPayload != nil {
			return createSetPoolPayload.Pool, err
		} else {
			return nil, err
		}
	})

	return &model.CreateNestedSetPoolPayload{Pool: pool}, err2
}

func (r *mutationResolver) CreateSingletonPool(ctx context.Context, input *model.CreateSingletonPoolInput) (*model.CreateSingletonPoolPayload, error) {
	var client = r.ClientFrom(ctx)

	resType, _ := client.ResourceType.Get(ctx, input.ResourceTypeID)
	if len(input.PoolValues) == 1 {
		_, rp, err := p.NewSingletonPoolWithMeta(ctx, client, resType, p.ToRawTypes(input.PoolValues)[0],
			input.PoolName, input.Description)
		retVal := model.CreateSingletonPoolPayload{Pool: rp}
		if err != nil {
			return &retVal, gqlerror.Errorf("Cannot create singleton pool: %v", err)
		} else {
			return &retVal, nil
		}
	} else {
		return &model.CreateSingletonPoolPayload{Pool: nil}, gqlerror.Errorf("Cannot create singleton pool, no resource provided")
	}
}

func (r *mutationResolver) CreateNestedSingletonPool(ctx context.Context, input model.CreateNestedSingletonPoolInput) (*model.CreateNestedSingletonPoolPayload, error) {
	var client = r.ClientFrom(ctx)

	nestedPool, err2 := createNestedPool(ctx, input.ParentResourceID, client, func() (*ent.ResourcePool, error) {
		poolInput := model.CreateSingletonPoolInput{
			ResourceTypeID: input.ResourceTypeID,
			PoolName:       input.PoolName,
			Description:    input.Description,
			PoolValues:     input.PoolValues,
		}
		payload, err := r.CreateSingletonPool(ctx, &poolInput)
		if payload != nil {
			return payload.Pool, err
		} else {
			return nil, err
		}

	})

	return &model.CreateNestedSingletonPoolPayload{Pool: nestedPool}, err2
}

func (r *mutationResolver) CreateAllocatingPool(ctx context.Context, input *model.CreateAllocatingPoolInput) (*model.CreateAllocatingPoolPayload, error) {
	var client = r.ClientFrom(ctx)
	emptyRetVal := model.CreateAllocatingPoolPayload{Pool: nil}

	var resPropertyType *ent.ResourceType = nil
	//create additional resource type IFF we are not a nested type
	//only root pool
	if input.PoolPropertyTypes != nil {
		rp, err2 := r.CreateResourceType(ctx, model.CreateResourceTypeInput{
			ResourceName:       input.PoolName + "-ResourceType",
			ResourceProperties: input.PoolPropertyTypes,
		})

		if err2 != nil || rp == nil || rp.ResourceType == nil {
			return &emptyRetVal, gqlerror.Errorf("Unable to create pool: %v", err2)
		}

		resPropertyType = rp.ResourceType
	}

	var poolProperties *ent.PoolProperties = nil

	// only root pool
	if resPropertyType != nil {
		pp, err := p.CreatePoolProperties(ctx, client, []map[string]interface{}{input.PoolProperties}, resPropertyType)
		if err != nil {
			return &emptyRetVal, gqlerror.Errorf("Unable to create pool properties: %v", err)
		}
		poolProperties = pp
	}

	resType, errRes := client.ResourceType.Get(ctx, input.ResourceTypeID)
	if errRes != nil {
		return &emptyRetVal, gqlerror.Errorf("Unable to create pool: %v", errRes)
	}
	allocationStrat, errAlloc := client.AllocationStrategy.Get(ctx, input.AllocationStrategyID)
	if errAlloc != nil {
		return &emptyRetVal, gqlerror.Errorf("Unable to create pool: %v", errAlloc)
	}

	_, rp, err := p.NewAllocatingPoolWithMeta(ctx, client, resType, allocationStrat,
		input.PoolName, input.Description, input.PoolDealocationSafetyPeriod, poolProperties)
	if err != nil {
		return &emptyRetVal, gqlerror.Errorf("Unable to create pool: %v", err)
	}

	return &model.CreateAllocatingPoolPayload{Pool: rp}, err
}

func (r *mutationResolver) CreateNestedAllocatingPool(ctx context.Context, input model.CreateNestedAllocatingPoolInput) (*model.CreateNestedAllocatingPoolPayload, error) {
	var client = r.ClientFrom(ctx)

	pool, err2 := createNestedPool(ctx, input.ParentResourceID, client, func() (*ent.ResourcePool, error) {
		poolInput := model.CreateAllocatingPoolInput{
			ResourceTypeID:              input.ResourceTypeID,
			PoolName:                    input.PoolName,
			Description:                 input.Description,
			AllocationStrategyID:        input.AllocationStrategyID,
			PoolDealocationSafetyPeriod: input.PoolDealocationSafetyPeriod,
		}
		poolPayload, err := r.CreateAllocatingPool(ctx, &poolInput)
		if poolPayload != nil {
			return poolPayload.Pool, err
		} else {
			return nil, err
		}
	})

	return &model.CreateNestedAllocatingPoolPayload{Pool: pool}, err2
}

func (r *mutationResolver) DeleteResourcePool(ctx context.Context, input model.DeleteResourcePoolInput) (*model.DeleteResourcePoolPayload, error) {
	client := r.ClientFrom(ctx)
	retVal := model.DeleteResourcePoolPayload{ResourcePoolID: input.ResourcePoolID}

	pool, errPool := p.ExistingPoolFromId(ctx, client, input.ResourcePoolID)

	if errPool != nil {
		return &retVal, gqlerror.Errorf("Unable to retrieve pool: %v", errPool)
	}

	errDp := pool.Destroy()

	if errDp != nil {
		return &retVal, gqlerror.Errorf("Unable to delete pool: %v", errDp)
	}

	return &retVal, nil
}

func (r *mutationResolver) CreateResourceType(ctx context.Context, input model.CreateResourceTypeInput) (*model.CreateResourceTypePayload, error) {
	var client = r.ClientFrom(ctx)

	var propertyTypes []*ent.PropertyType
	for propName, rawPropType := range input.ResourceProperties {
		var propertyType, err = p.CreatePropertyType(ctx, client, propName, rawPropType)
		if err != nil {
			return &model.CreateResourceTypePayload{ResourceType: nil}, gqlerror.Errorf("Unable to create resource type: %v", err)
		}
		propertyTypes = append(propertyTypes, propertyType)
	}

	resType, err2 := client.ResourceType.Create().
		SetName(input.ResourceName).
		AddPropertyTypes(propertyTypes...).
		Save(ctx)
	if err2 != nil {
		return &model.CreateResourceTypePayload{ResourceType: nil}, gqlerror.Errorf("Unable to create resource type: %v", err2)
	}

	return &model.CreateResourceTypePayload{ResourceType: resType}, nil
}

func (r *mutationResolver) DeleteResourceType(ctx context.Context, input model.DeleteResourceTypeInput) (*model.DeleteResourceTypePayload, error) {
	client := r.ClientFrom(ctx)
	resourceType, err := client.ResourceType.Get(ctx, input.ResourceTypeID)
	retValue := &model.DeleteResourceTypePayload{ResourceTypeID: input.ResourceTypeID}
	if err != nil {
		return retValue, gqlerror.Errorf("Unable to delete resource type - cannot find by ID %d: %v", input.ResourceTypeID, err)
	}

	pools, err := client.ResourceType.QueryPools(resourceType).All(ctx)

	if err != nil {
		return retValue, gqlerror.Errorf("Unable to delete resource type - error obtaining pools: %v", err)
	}

	if len(pools) > 0 {
		return retValue, gqlerror.Errorf("Unable to delete resource type, there are pools attached to it")
	}

	// delete property types
	_, err = client.PropertyType.Delete().Where(propertytype.HasResourceTypeWith(resourcetype.ID(resourceType.ID))).Exec(ctx)
	if err != nil {
		return retValue, gqlerror.Errorf("Unable to delete resource type - error deleting property types: %v", err)
	}

	// delete resource type
	if err := client.ResourceType.DeleteOneID(input.ResourceTypeID).Exec(ctx); err == nil {
		return &model.DeleteResourceTypePayload{ResourceTypeID: input.ResourceTypeID}, nil
	} else {
		return retValue, gqlerror.Errorf("Unable to delete resource type: %v", err)
	}
}

func (r *mutationResolver) UpdateResourceTypeName(ctx context.Context, input model.UpdateResourceTypeNameInput) (*model.UpdateResourceTypeNamePayload, error) {
	var client = r.ClientFrom(ctx)
	retValue := &model.UpdateResourceTypeNamePayload{ResourceTypeID: input.ResourceTypeID}
	if _, err := client.ResourceType.UpdateOneID(input.ResourceTypeID).SetName(input.ResourceName).Save(ctx); err != nil {
		return retValue, gqlerror.Errorf("Unable to update resource type: %v", err)
	} else {
		return retValue, nil
	}
}

func (r *outputCursorResolver) ID(ctx context.Context, obj *ent.Cursor) (string, error) {
	//this will never be called because ent.Cursor will use its msgpack annotation
	return "", nil
}

func (r *propertyTypeResolver) Type(ctx context.Context, obj *ent.PropertyType) (string, error) {
	// Just converts enum to string
	return obj.Type.String(), nil
}

func (r *queryResolver) QueryPoolCapacity(ctx context.Context, poolID int) (*model.PoolCapacityPayload, error) {
	pool, err := p.ExistingPoolFromId(ctx, r.ClientFrom(ctx), poolID)

	if err != nil {
		return nil, gqlerror.Errorf("Unable to find pool: %v", err)
	}

	freeCapacity, utilizedCapacity, err2 := pool.Capacity()

	if err2 != nil {
		return nil, gqlerror.Errorf("Unable to compute capacity: %v", err2)
	}

	return &model.PoolCapacityPayload{
		FreeCapacity:     freeCapacity,
		UtilizedCapacity: utilizedCapacity,
	}, nil
}

func (r *queryResolver) QueryPoolTypes(ctx context.Context) ([]resourcePool.PoolType, error) {
	poolTypes := []resourcePool.PoolType{
		resourcePool.PoolTypeSingleton,
		resourcePool.PoolTypeSet,
		resourcePool.PoolTypeAllocating}
	return poolTypes, nil
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

func (r *queryResolver) QueryResourcePool(ctx context.Context, poolID int) (*ent.ResourcePool, error) {
	return r.ClientFrom(ctx).ResourcePool.Get(ctx, poolID)
}

func (r *queryResolver) QueryResourcePools(ctx context.Context, resourceTypeID *int) ([]*ent.ResourcePool, error) {
	client := r.ClientFrom(ctx)
	query := client.ResourcePool.Query()

	if resourceTypeID != nil {
		query.Where(resourcePool.HasResourceTypeWith(resourcetype.ID(*resourceTypeID)))
	}

	if resourcePools, err := query.All(ctx); err != nil {
		return nil, gqlerror.Errorf("Unable to query resource pools: %v", err)
	} else {
		return resourcePools, nil
	}
}

func (r *queryResolver) QueryRootResourcePools(ctx context.Context, resourceTypeID *int) ([]*ent.ResourcePool, error) {
	client := r.ClientFrom(ctx)
	query := client.ResourcePool.
		Query().
		Where(resourcePool.Not(resourcePool.HasParentResource()))

	if resourceTypeID != nil {
		query.Where(resourcePool.HasResourceTypeWith(resourcetype.ID(*resourceTypeID)))
	}

	if resourcePools, err := query.All(ctx); err != nil {
		return nil, gqlerror.Errorf("Unable to query resource pools: %v", err)
	} else {
		return resourcePools, nil
	}
}

func (r *queryResolver) QueryLeafResourcePools(ctx context.Context, resourceTypeID *int) ([]*ent.ResourcePool, error) {
	client := r.ClientFrom(ctx)
	query := client.ResourcePool.
		Query().
		Where(resourcePool.HasParentResource()).
		Where(resourcePool.Not(resourcePool.HasClaimsWith(resource.HasNestedPool())))

	if resourceTypeID != nil {
		query.Where(resourcePool.HasResourceTypeWith(resourcetype.ID(*resourceTypeID)))
	}

	if resourcePools, err := query.All(ctx); err != nil {
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

func (r *queryResolver) Node(ctx context.Context, id int) (ent.Noder, error) {
	var client = r.ClientFrom(ctx)
	return client.Noder(ctx, id)
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

func (r *resourceResolver) NestedPool(ctx context.Context, obj *ent.Resource) (*ent.ResourcePool, error) {
	if es, err := obj.Edges.NestedPoolOrErr(); !ent.IsNotLoaded(err) {
		return es, err
	}
	if pool, err := obj.QueryNestedPool().First(ctx); ent.IsNotFound(err) {
		return nil, nil
	} else if err == nil {
		return pool, nil
	} else {
		return nil, gqlerror.Errorf("Unable to query nested pool: %v", err)
	}
}

func (r *resourcePoolResolver) ResourceType(ctx context.Context, obj *ent.ResourcePool) (*ent.ResourceType, error) {
	if es, err := obj.Edges.ResourceTypeOrErr(); !ent.IsNotLoaded(err) {
		return es, err
	}
	return obj.QueryResourceType().Only(ctx)
}

func (r *resourcePoolResolver) Resources(ctx context.Context, obj *ent.ResourcePool) ([]*ent.Resource, error) {
	return p.GetResourceFromPool(ctx, obj)
}

func (r *resourcePoolResolver) AllocatedResources(ctx context.Context, obj *ent.ResourcePool, first *int, last *int, before *string, after *string) (*ent.ResourceConnection, error) {
	//pagination https://relay.dev/graphql/connections.htm

	//we query resources only for a specific pool
	onlyForPool := func(rq *ent.ResourceQuery) (*ent.ResourceQuery, error) {
		return rq.Where(resource.HasPoolWith(resourcePool.ID(obj.ID))), nil
	}

	afterCursor, errA := decodeCursor(after)
	if errA != nil {
		return nil, errA
	}

	beforeCursor, errB := decodeCursor(before)
	if errB != nil {
		return nil, errB
	}

	return r.ClientFrom(ctx).Resource.Query().Paginate(ctx, afterCursor, first, beforeCursor, last, ent.WithResourceFilter(onlyForPool))
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

// OutputCursor returns generated.OutputCursorResolver implementation.
func (r *Resolver) OutputCursor() generated.OutputCursorResolver { return &outputCursorResolver{r} }

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
type outputCursorResolver struct{ *Resolver }
type propertyTypeResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type resourceResolver struct{ *Resolver }
type resourcePoolResolver struct{ *Resolver }
type resourceTypeResolver struct{ *Resolver }
type tagResolver struct{ *Resolver }
