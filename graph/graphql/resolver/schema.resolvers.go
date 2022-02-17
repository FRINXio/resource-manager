package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/facebook/ent/dialect/sql"
	"github.com/facebook/ent/dialect/sql/sqljson"
	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/allocationstrategy"
	"github.com/net-auto/resourceManager/ent/predicate"
	"github.com/net-auto/resourceManager/ent/propertytype"
	"github.com/net-auto/resourceManager/ent/resource"
	resourcePool "github.com/net-auto/resourceManager/ent/resourcepool"
	"github.com/net-auto/resourceManager/ent/resourcetype"
	"github.com/net-auto/resourceManager/graph/graphql/generated"
	"github.com/net-auto/resourceManager/graph/graphql/model"
	log "github.com/net-auto/resourceManager/logging"
	p "github.com/net-auto/resourceManager/pools"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func (r *mutationResolver) CreateTag(ctx context.Context, input model.CreateTagInput) (*model.CreateTagPayload, error) {
	var client = r.ClientFrom(ctx)
	tagEnt, err := createTag(ctx, client, input.TagText)

	if err != nil {
		log.Error(ctx, err, "Unable to create new tag")
		return &model.CreateTagPayload{Tag: nil}, gqlerror.Errorf("Unable to create tag: %v", err)
	}

	return &model.CreateTagPayload{Tag: tagEnt}, nil
}

func (r *mutationResolver) UpdateTag(ctx context.Context, input model.UpdateTagInput) (*model.UpdateTagPayload, error) {
	var client = r.ClientFrom(ctx)
	tagEnt, err := client.Tag.UpdateOneID(input.TagID).SetTag(input.TagText).Save(ctx)
	if err != nil {
		log.Error(ctx, err, "Unable to update tag ID %d", input.TagID)
		return &model.UpdateTagPayload{Tag: nil}, gqlerror.Errorf("Unable to update tag: %v", err)
	}
	return &model.UpdateTagPayload{Tag: tagEnt}, nil
}

func (r *mutationResolver) DeleteTag(ctx context.Context, input model.DeleteTagInput) (*model.DeleteTagPayload, error) {
	var client = r.ClientFrom(ctx)
	err := client.Tag.DeleteOneID(input.TagID).Exec(ctx)
	if err != nil {
		log.Error(ctx, err, "Unable to delete tag ID %d", input.TagID)
		return &model.DeleteTagPayload{TagID: input.TagID}, gqlerror.Errorf("Unable to delete tag: %v", err)
	}

	return &model.DeleteTagPayload{TagID: input.TagID}, nil
}

func (r *mutationResolver) TagPool(ctx context.Context, input model.TagPoolInput) (*model.TagPoolPayload, error) {
	var client = r.ClientFrom(ctx)
	tag, err := client.Tag.UpdateOneID(input.TagID).AddPoolIDs(input.PoolID).Save(ctx)
	if err != nil {
		log.Error(ctx, err, "Unable to tag pool ID %d", input.PoolID)
		return &model.TagPoolPayload{Tag: nil}, gqlerror.Errorf("Unable to tag pool: %v", err)
	}
	return &model.TagPoolPayload{Tag: tag}, nil
}

func (r *mutationResolver) UntagPool(ctx context.Context, input model.UntagPoolInput) (*model.UntagPoolPayload, error) {
	var client = r.ClientFrom(ctx)
	tag, err := client.Tag.UpdateOneID(input.TagID).RemovePoolIDs(input.PoolID).Save(ctx)
	if err != nil {
		log.Error(ctx, err, "Unable to un-tag pool ID %d", input.PoolID)
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
		log.Error(ctx, err, "Unable create a new allocation strategy")
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
		log.Error(ctx, err, "Unable to find allocation strategy ID %d", input.AllocationStrategyID)
		return &emptyRetVal, gqlerror.Errorf("Unable to delete strategy: %v", err)
	} else {

		if len(strat.Edges.Pools) > 0 {
			log.Error(ctx, err, "Unable to delete allocation strategy ID %d because it is used by %d pool(s)", input.AllocationStrategyID, len(strat.Edges.Pools))
			return &emptyRetVal, gqlerror.Errorf("Unable to delete, Allocation strategy is still in use")
		}

		if err := client.AllocationStrategy.DeleteOneID(input.AllocationStrategyID).Exec(ctx); err != nil {
			log.Error(ctx, err, "Unable to delete allocation strategy ID %d", input.AllocationStrategyID)
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
		log.Error(ctx, err, "Unable to find allocation strategy %d", allocationStrategyID)
		return nil, gqlerror.Errorf("Unable to get strategy: %v", err)
	}
	// TODO keep just single instance
	wasmer, err := p.NewWasmerUsingEnvVars()
	if err != nil {
		log.Error(ctx, err, "Unable to create a scripting engine (wasmer)")
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
		log.Error(ctx, err, "Error while running script on pool \"%s\" strategy ID %d", resourcePool.ResourcePoolName, allocationStrategyID)
		return nil, gqlerror.Errorf("Error while running the script: %v", err)
	}
	result := make(map[string]interface{})
	result["stdout"] = parsedOutputFromStrat
	result["stderr"] = stdErr
	return result, nil
}

func (r *mutationResolver) ClaimResource(ctx context.Context, poolID int, description *string, userInput map[string]interface{}) (*ent.Resource, error) {
	return r.ClaimResourceWithAltID(ctx, poolID, description, userInput, nil)
}

func (r *mutationResolver) ClaimResourceWithAltID(ctx context.Context, poolID int, description *string, userInput map[string]interface{}, alternativeID map[string]interface{}) (*ent.Resource, error) {
	pool, err := p.ExistingPoolFromId(ctx, r.ClientFrom(ctx), poolID)
	if err != nil {
		return nil, gqlerror.Errorf("Unable to claim resource: %v", err)
	}

	if res, err := pool.ClaimResource(userInput, description, alternativeID); err != nil {
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

	log.Error(ctx, err, "Unable to free resource on pool ID %d with properties %+v", poolID, input)
	return "", gqlerror.Errorf("Unable to free resource: %v", err)
}

func (r *mutationResolver) CreateSetPool(ctx context.Context, input model.CreateSetPoolInput) (*model.CreateSetPoolPayload, error) {
	var client = r.ClientFrom(ctx)

	resType, err := client.ResourceType.Get(ctx, input.ResourceTypeID)
	if err != nil {
		log.Error(ctx, err, "Unable to retrieve resource type for the set-pool (resource type ID: %d)", input.ResourceTypeID)
		return &model.CreateSetPoolPayload{Pool: nil}, gqlerror.Errorf("Unable to create pool: %v", err)
	}
	_, rp, err := p.NewSetPoolWithMeta(ctx, client, resType, p.ToRawTypes(input.PoolValues),
		input.PoolName, input.Description, input.PoolDealocationSafetyPeriod)

	if err := createTagsAndTagPool(ctx, client, rp, input.Tags); err != nil {
		log.Error(ctx, err, "Unable to tag the pool with tags: %v", input.Tags)
		return nil, err
	}

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
			Tags:                        input.Tags,
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

	resType, err := client.ResourceType.Get(ctx, input.ResourceTypeID)

	if err != nil {
		log.Error(ctx, err, "Unable to retrieve resource type for the set-pool (resource type ID: %d)", input.ResourceTypeID)
	}

	if len(input.PoolValues) == 1 {
		_, rp, err := p.NewSingletonPoolWithMeta(ctx, client, resType, p.ToRawTypes(input.PoolValues)[0],
			input.PoolName, input.Description)

		if err := createTagsAndTagPool(ctx, client, rp, input.Tags); err != nil {
			log.Error(ctx, err, "Unable to tag the pool with tags: %v", input.Tags)
			return nil, err
		}

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
			Tags:           input.Tags,
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
			log.Error(ctx, err2, "Unable to create a resource-type for a root pool \"%s\"", input.PoolName)
			return &emptyRetVal, gqlerror.Errorf("Unable to create pool: %v", err2)
		}

		resPropertyType = rp.ResourceType
	}

	var poolProperties *ent.PoolProperties = nil

	// only root pool
	if resPropertyType != nil {
		pp, err := p.CreatePoolProperties(ctx, client, []map[string]interface{}{input.PoolProperties}, resPropertyType)
		if err != nil {
			log.Error(ctx, err, "Unable to create pool properties for a root pool \"%s\"", input.PoolName)
			return &emptyRetVal, gqlerror.Errorf("Unable to create pool properties: %v", err)
		}
		poolProperties = pp
	}

	resType, errRes := client.ResourceType.Get(ctx, input.ResourceTypeID)
	if errRes != nil {
		log.Error(ctx, errRes, "Unable to retrieve resource type for pool (resource type ID: %d)", input.ResourceTypeID)
		return &emptyRetVal, gqlerror.Errorf("Unable to create pool: %v", errRes)
	}
	allocationStrat, errAlloc := client.AllocationStrategy.Get(ctx, input.AllocationStrategyID)
	if errAlloc != nil {
		log.Error(ctx, errAlloc, "Unable to retrieve allocation strategy for pool (strategy ID: %d)", input.AllocationStrategyID)
		return &emptyRetVal, gqlerror.Errorf("Unable to create pool: %v", errAlloc)
	}

	_, rp, err := p.NewAllocatingPoolWithMeta(ctx, client, resType, allocationStrat,
		input.PoolName, input.Description, input.PoolDealocationSafetyPeriod, poolProperties)

	if err := createTagsAndTagPool(ctx, client, rp, input.Tags); err != nil {
		log.Error(ctx, err, "Unable to tag the pool with tags: %v", input.Tags)
		return nil, err
	}

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
			Tags:                        input.Tags,
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
		log.Error(ctx, err2, "Unable to create a new resource type")
		return &model.CreateResourceTypePayload{ResourceType: nil}, gqlerror.Errorf("Unable to create resource type: %v", err2)
	}

	return &model.CreateResourceTypePayload{ResourceType: resType}, nil
}

func (r *mutationResolver) DeleteResourceType(ctx context.Context, input model.DeleteResourceTypeInput) (*model.DeleteResourceTypePayload, error) {
	client := r.ClientFrom(ctx)
	resourceType, err := client.ResourceType.Get(ctx, input.ResourceTypeID)
	retValue := &model.DeleteResourceTypePayload{ResourceTypeID: input.ResourceTypeID}
	if err != nil {
		log.Error(ctx, err, "Unable to retrieve resource type ID %d", input.ResourceTypeID)
		return retValue, gqlerror.Errorf("Unable to delete resource type - cannot find by ID %d: %v", input.ResourceTypeID, err)
	}

	pools, err := client.ResourceType.QueryPools(resourceType).All(ctx)

	if err != nil {
		log.Error(ctx, err, "Unable to retrieve pools associated with resource type ID %d", input.ResourceTypeID)
		return retValue, gqlerror.Errorf("Unable to delete resource type - error obtaining pools: %v", err)
	}

	if len(pools) > 0 {
		log.Warn(ctx, "Unable delete resource type ID %d - there are %d pool(s) associated with it", input.ResourceTypeID, len(pools))
		return retValue, gqlerror.Errorf("Unable to delete resource type, there are pools attached to it")
	}

	// delete property types
	_, err = client.PropertyType.Delete().Where(propertytype.HasResourceTypeWith(resourcetype.ID(resourceType.ID))).Exec(ctx)
	if err != nil {
		log.Error(ctx, err, "Unable delete resource type ID %d", input.ResourceTypeID)
		return retValue, gqlerror.Errorf("Unable to delete resource type - error deleting property types: %v", err)
	}

	// delete resource type
	if err := client.ResourceType.DeleteOneID(input.ResourceTypeID).Exec(ctx); err == nil {
		return &model.DeleteResourceTypePayload{ResourceTypeID: input.ResourceTypeID}, nil
	} else {
		log.Error(ctx, err, "Unable delete resource type ID %d", input.ResourceTypeID)
		return retValue, gqlerror.Errorf("Unable to delete resource type: %v", err)
	}
}

func (r *mutationResolver) UpdateResourceTypeName(ctx context.Context, input model.UpdateResourceTypeNameInput) (*model.UpdateResourceTypeNamePayload, error) {
	var client = r.ClientFrom(ctx)
	retValue := &model.UpdateResourceTypeNamePayload{ResourceTypeID: input.ResourceTypeID}
	if _, err := client.ResourceType.UpdateOneID(input.ResourceTypeID).SetName(input.ResourceName).Save(ctx); err != nil {
		log.Error(ctx, err, "Unable to update resource type ID %d", input.ResourceTypeID)
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

func (r *queryResolver) QueryResourcesByAltID(ctx context.Context, input map[string]interface{}, poolID *int) ([]*ent.Resource, error) {
	if poolID != nil {
		pool, err := p.ExistingPoolFromId(ctx, r.ClientFrom(ctx), *poolID)
		if err != nil {
			return nil, gqlerror.Errorf("Unable to query resources: %v", err)
		}
		typeFixedAlternativeId, errFix := p.ConvertValuesToFloat64(ctx, input)
		if errFix != nil {
			return nil, gqlerror.Errorf("Unable to process input data", err)
		}
		return pool.QueryResourcesByAltId(typeFixedAlternativeId)
	}

	res, err := r.ClientFrom(ctx).Resource.Query().
		Where(func(selector *sql.Selector) {
			for k, v := range input {
				selector.Where(sqljson.ValueEQ("alternate_id", v, sqljson.Path(k)))
			}
		}).All(ctx)
	if err != nil {
		log.Error(ctx, err, "Unable to retrieve resources with alternative ID %v", input)
		return nil, gqlerror.Errorf("Unable to query resources: %v", err)
	}

	if res != nil {
		return res, nil
	}

	log.Warn(ctx, "There is not such resources with alternative ID %v", input)
	return nil, gqlerror.Errorf("Unable to query resources: %v", err)
}

func (r *queryResolver) QueryAllocationStrategy(ctx context.Context, allocationStrategyID int) (*ent.AllocationStrategy, error) {
	client := r.ClientFrom(ctx)
	if strats, err := client.AllocationStrategy.Query().Where(allocationstrategy.ID(allocationStrategyID)).Only(ctx); err != nil {
		log.Error(ctx, err, "Unable to retrieve allocation strategy ID %d", allocationStrategyID)
		return nil, gqlerror.Errorf("Unable to query strategy: %v", err)
	} else {
		return strats, nil
	}
}

func (r *queryResolver) QueryAllocationStrategies(ctx context.Context, byName *string) ([]*ent.AllocationStrategy, error) {
	client := r.ClientFrom(ctx)
	query := client.AllocationStrategy.Query()

	if byName != nil {
		query = query.Where(allocationstrategy.Name(*byName))
	}

	if strats, err := query.All(ctx); err != nil {
		log.Error(ctx, err, "Unable to retrieve allocation strategies")
		return nil, gqlerror.Errorf("Unable to query strategies: %v", err)
	} else {
		return strats, nil
	}
}

func (r *queryResolver) QueryResourceTypes(ctx context.Context, byName *string) ([]*ent.ResourceType, error) {
	client := r.ClientFrom(ctx)
	query := client.ResourceType.Query()

	// Filter out pool properties that are stored in resource type table
	query = query.Where(resourcetype.Not(resourcetype.HasPoolProperties()))

	if byName != nil {
		query = query.Where(resourcetype.Name(*byName))
	}

	if resourceTypes, err := query.All(ctx); err != nil {
		log.Error(ctx, err, "Unable to retrieve resource types")
		return nil, gqlerror.Errorf("Unable to query resource types: %v", err)
	} else {
		return resourceTypes, nil
	}
}

func (r *queryResolver) QueryResourcePool(ctx context.Context, poolID int) (*ent.ResourcePool, error) {
	rp, err := r.ClientFrom(ctx).ResourcePool.Get(ctx, poolID)

	if err != nil {
		log.Error(ctx, err, "Unable to retrieve resource pool")
	}

	return rp, err
}

func (r *queryResolver) QueryResourcePools(ctx context.Context, resourceTypeID *int, tags *model.TagOr) ([]*ent.ResourcePool, error) {
	client := r.ClientFrom(ctx)
	query := client.ResourcePool.Query()

	if resourceTypeID != nil {
		query.Where(resourcePool.HasResourceTypeWith(resourcetype.ID(*resourceTypeID)))
	}

	if tags != nil {
		// TODO make sure all tags exist
		query.Where(resourcePoolTagPredicate(tags))
	}

	if resourcePools, err := query.All(ctx); err != nil {
		log.Error(ctx, err, "Unable to retrieve resource pools")
		return nil, gqlerror.Errorf("Unable to query resource pools: %v", err)
	} else {
		return resourcePools, nil
	}
}

func (r *queryResolver) QueryRecentlyActiveResourcePools(ctx context.Context, fromDatetime string, toDatetime *string) ([]*ent.ResourcePool, error) {
	client := r.ClientFrom(ctx)
	query := client.ResourcePool.Query()
	dateFrom, err := time.Parse("2006-01-02-15", fromDatetime)
	if err != nil {
		log.Error(ctx, err, "Unable to parse date from: "+fromDatetime+". Must be in format: YYYY-MM-DD-hh.")
		return nil, gqlerror.Errorf("Unable to parse date from: "+fromDatetime+
			". Must be in format: YYYY-MM-DD-hh. Error: %v", err)
	}

	if toDatetime != nil && len(*toDatetime) != 0 {
		dateTo, err := time.Parse("2006-01-02-15", *toDatetime)
		if err != nil {
			log.Error(ctx, err, "Unable to parse date to: "+*toDatetime+". Must be in format: YYYY-MM-DD-hh.")
			return nil, gqlerror.Errorf("Unable to parse date to: "+*toDatetime+
				". Must be in format: YYYY-MM-DD-hh. Error: %v", err)
		}
		query.Where(resourcePool.HasClaimsWith(resource.And(
			resource.UpdatedAtGTE(dateFrom), resource.UpdatedAtLTE(dateTo))))
	} else {
		currentDate := time.Now()
		query.Where(resourcePool.HasClaimsWith(resource.And(
			resource.UpdatedAtGTE(dateFrom), resource.UpdatedAtLTE(currentDate))))
	}

	if resourcePools, err := query.All(ctx); err != nil {
		log.Error(ctx, err, "Unable to retrieve resource pools")
		return nil, gqlerror.Errorf("Unable to query resource pools: %v", err)
	} else {
		return resourcePools, nil
	}
}

func (r *queryResolver) QueryResourcePoolHierarchyPath(ctx context.Context, poolID int) ([]*ent.ResourcePool, error) {
	client := r.ClientFrom(ctx)
	currentPool, err := queryPoolWithParent(ctx, poolID, client)
	if err != nil {
		log.Error(ctx, err, "Unable to find pool")
		return nil, gqlerror.Errorf("Unable to find pool: %v", err)
	}

	var hierarchy []*ent.ResourcePool

	for hasParent(currentPool) {
		parentPool := currentPool.Edges.ParentResource.Edges.Pool
		hierarchy = append([]*ent.ResourcePool{parentPool}, hierarchy...)
		if currentPool, err = queryPoolWithParent(ctx, parentPool.ID, client); err != nil {
			log.Error(ctx, err, "Unable to find pool")
			return nil, gqlerror.Errorf("Unable to find pool: %v", err)
		}
	}

	return hierarchy, nil
}

func (r *queryResolver) QueryRootResourcePools(ctx context.Context, resourceTypeID *int, tags *model.TagOr) ([]*ent.ResourcePool, error) {
	client := r.ClientFrom(ctx)
	query := client.ResourcePool.
		Query().
		Where(resourcePool.Not(resourcePool.HasParentResource()))

	if resourceTypeID != nil {
		query.Where(resourcePool.HasResourceTypeWith(resourcetype.ID(*resourceTypeID)))
	}

	if tags != nil {
		// TODO make sure all tags exist
		query.Where(resourcePoolTagPredicate(tags))
	}

	if resourcePools, err := query.All(ctx); err != nil {
		log.Error(ctx, err, "Unable to retrieve root resource pools")
		return nil, gqlerror.Errorf("Unable to query resource pools: %v", err)
	} else {
		return resourcePools, nil
	}
}

func (r *queryResolver) QueryLeafResourcePools(ctx context.Context, resourceTypeID *int, tags *model.TagOr) ([]*ent.ResourcePool, error) {
	client := r.ClientFrom(ctx)
	query := client.ResourcePool.
		Query().
		Where(resourcePool.HasParentResource()).
		Where(resourcePool.Not(resourcePool.HasClaimsWith(resource.HasNestedPool())))

	if resourceTypeID != nil {
		query.Where(resourcePool.HasResourceTypeWith(resourcetype.ID(*resourceTypeID)))
	}

	if tags != nil {
		// TODO make sure all tags exist
		query.Where(resourcePoolTagPredicate(tags))
	}

	if resourcePools, err := query.All(ctx); err != nil {
		log.Error(ctx, err, "Unable to retrieve leaf resource pools")
		return nil, gqlerror.Errorf("Unable to query resource pools: %v", err)
	} else {
		return resourcePools, nil
	}
}

func (r *queryResolver) SearchPoolsByTags(ctx context.Context, tags *model.TagOr) ([]*ent.ResourcePool, error) {
	var client = r.ClientFrom(ctx)

	var condition predicate.ResourcePool

	if tags != nil {
		// TODO make sure all tags exist
		condition = resourcePoolTagPredicate(tags)
	}

	var (
		matchedPools []*ent.ResourcePool
		err          error
	)
	if condition == nil {
		matchedPools, err = client.ResourcePool.Query().All(ctx)
	} else {
		matchedPools, err = client.ResourcePool.Query().Where(condition).All(ctx)
	}

	if err != nil {
		log.Error(ctx, err, "Unable to retrieve pools by tags")
		return nil, gqlerror.Errorf("Unable to query pools: %v", err)
	}
	return matchedPools, nil
}

func (r *queryResolver) QueryTags(ctx context.Context) ([]*ent.Tag, error) {
	var client = r.ClientFrom(ctx)
	tags, err := client.Tag.Query().All(ctx)
	if err != nil {
		log.Error(ctx, err, "Unable to retrieve tags")
		return nil, gqlerror.Errorf("Unable to query tags: %v", err)
	}
	return tags, nil
}

func (r *queryResolver) Node(ctx context.Context, id int) (ent.Noder, error) {
	var client = r.ClientFrom(ctx)
	node, err := client.Noder(ctx, id)

	if err != nil {
		log.Error(ctx, err, "Unable to retrieve node with ID %d", id)
	}

	return node, err
}

func (r *resourceResolver) NestedPool(ctx context.Context, obj *ent.Resource) (*ent.ResourcePool, error) {
	if es, err := obj.Edges.NestedPoolOrErr(); !ent.IsNotLoaded(err) {
		log.Error(ctx, err, "Unable to retrieve nested pool for resource with ID %d", obj.ID)
		return es, err
	}
	if pool, err := obj.QueryNestedPool().First(ctx); ent.IsNotFound(err) {
		log.Warn(ctx, "No nested resource pool found for resource with ID %d", obj.ID)
		return nil, nil
	} else if err == nil {
		return pool, nil
	} else {
		log.Error(ctx, err, "Unable to retrieve nested pool for resource with ID %d", obj.ID)
		return nil, gqlerror.Errorf("Unable to query nested pool: %v", err)
	}
}

func (r *resourceResolver) ParentPool(ctx context.Context, obj *ent.Resource) (*ent.ResourcePool, error) {
	if es, err := obj.Edges.PoolOrErr(); !ent.IsNotLoaded(err) {
		log.Error(ctx, err, "Unable to retrieve pool for resource with ID %d", obj.ID)
		return es, err
	}
	if pool, err := obj.QueryPool().Only(ctx); err == nil {
		return pool, nil
	} else {
		log.Error(ctx, err, "Unable to retrieve parent pool for resource with ID %d", obj.ID)
		return nil, gqlerror.Errorf("Unable to query parent pool: %v", err)
	}
}

func (r *resourceResolver) Properties(ctx context.Context, obj *ent.Resource) (map[string]interface{}, error) {
	props, err := obj.QueryProperties().WithType().All(ctx)
	if err != nil {
		log.Error(ctx, err, "Unable to retrieve properties for resource with ID %d", obj.ID)
		return nil, gqlerror.Errorf("Unable to query properties: %v", err)
	}

	if props, err := p.PropertiesToMap(props); err != nil {
		log.Error(ctx, err, "Unable to process properties for resource with ID %d", obj.ID)
		return nil, gqlerror.Errorf("Unable to query properties: %v", err)
	} else {
		return props, nil
	}
}

func (r *resourceResolver) AlternativeID(ctx context.Context, obj *ent.Resource) (map[string]interface{}, error) {
	return obj.AlternateID, nil
}

func (r *resourcePoolResolver) AllocationStrategy(ctx context.Context, obj *ent.ResourcePool) (*ent.AllocationStrategy, error) {
	if obj.PoolType != resourcePool.PoolTypeAllocating {
		log.Warn(ctx, "Pool with ID %d does not have an allocation strategy", obj.ID)
		return nil, nil
	}
	if es, err := obj.Edges.AllocationStrategyOrErr(); !ent.IsNotLoaded(err) {
		log.Error(ctx, err, "Loading allocation strategy for pool ID %d failed", obj.ID)
		return es, err
	}
	strategy, err := obj.QueryAllocationStrategy().Only(ctx)

	if err != nil {
		log.Error(ctx, err, "Loading allocation strategy for pool ID %d failed", obj.ID)
	}

	return strategy, err
}

func (r *resourcePoolResolver) Capacity(ctx context.Context, obj *ent.ResourcePool) (*model.PoolCapacityPayload, error) {
	return r.Query().QueryPoolCapacity(ctx, obj.ID)
}

func (r *resourcePoolResolver) ParentResource(ctx context.Context, obj *ent.ResourcePool) (*ent.Resource, error) {
	if es, err := obj.Edges.ParentResourceOrErr(); !ent.IsNotLoaded(err) {
		log.Error(ctx, err, "Loading parent resource for pool ID %d failed", obj.ID)
		return es, err
	}

	if pr, err := obj.QueryParentResource().Only(ctx); err != nil {
		log.Error(ctx, err, "Loading parent resource for pool ID %d failed", obj.ID)
		return nil, err
	} else {
		return pr, err
	}
}

func (r *resourcePoolResolver) PoolProperties(ctx context.Context, obj *ent.ResourcePool) (map[string]interface{}, error) {
	var (
		props *ent.PoolProperties
		err   error
	)

	props, err = obj.Edges.PoolPropertiesOrErr()

	if err != nil && ent.IsNotLoaded(err) {

		props, err = obj.
			QueryPoolProperties().
			WithProperties(func(query *ent.PropertyQuery) {
				query.WithType()
			}).First(ctx)

		if ent.IsNotFound(err) {
			return make(map[string]interface{}), nil
		}
		if err != nil {
			log.Error(ctx, err, "Loading pool properties for pool ID %d failed", obj.ID)
			return nil, err
		}

	} else if err != nil {
		log.Error(ctx, err, "Loading pool properties for pool ID %d failed", obj.ID)
		return nil, err
	}

	if props, err := p.PropertiesToMap(props.Edges.Properties); err != nil {
		log.Error(ctx, err, "Unable to process properties for pool with ID %d", obj.ID)
		return nil, gqlerror.Errorf("Unable to query properties: %v", err)
	} else {
		return props, nil
	}
}

func (r *resourcePoolResolver) ResourceType(ctx context.Context, obj *ent.ResourcePool) (*ent.ResourceType, error) {
	if es, err := obj.Edges.ResourceTypeOrErr(); !ent.IsNotLoaded(err) {
		log.Error(ctx, err, "Unable to retrieve resource type for pool with ID %d", obj.ID)
		return es, err
	}
	rt, err := obj.QueryResourceType().Only(ctx)

	if err != nil {
		log.Error(ctx, err, "Unable to retrieve resource type for pool with ID %d", obj.ID)
	}

	return rt, err
}

func (r *resourcePoolResolver) Resources(ctx context.Context, obj *ent.ResourcePool) ([]*ent.Resource, error) {
	resources, err := p.GetResourceFromPool(ctx, obj)

	if err != nil {
		log.Error(ctx, err, "Unable to retrieve resources for pool with ID %d", obj.ID)
	}

	return resources, err
}

func (r *resourcePoolResolver) Tags(ctx context.Context, obj *ent.ResourcePool) ([]*ent.Tag, error) {
	if es, err := obj.Edges.TagsOrErr(); !ent.IsNotLoaded(err) {
		log.Error(ctx, err, "Loading tags for pool ID %d failed", obj.ID)
		return es, err
	}

	tags, err := obj.QueryTags().All(ctx)

	if err != nil {
		log.Error(ctx, err, "Loading tags for pool ID %d failed", obj.ID)
	}

	return tags, err
}

func (r *resourcePoolResolver) AllocatedResources(ctx context.Context, obj *ent.ResourcePool, first *int, last *int, before *string, after *string) (*ent.ResourceConnection, error) {
	//pagination https://relay.dev/graphql/connections.htm

	//we query resources only for a specific pool
	onlyForPool := func(rq *ent.ResourceQuery) (*ent.ResourceQuery, error) {
		return rq.Where(resource.HasPoolWith(resourcePool.ID(obj.ID))), nil
	}

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

	resourceConnection, err := r.ClientFrom(ctx).Resource.Query().Paginate(ctx, afterCursor, first, beforeCursor, last, ent.WithResourceFilter(onlyForPool))

	if err != nil {
		log.Error(ctx, errB, "Loading resources for a pagination query for pool ID %d failed", obj.ID)
	}

	return resourceConnection, err
}

func (r *resourceTypeResolver) Pools(ctx context.Context, obj *ent.ResourceType) ([]*ent.ResourcePool, error) {
	if es, err := obj.Edges.PoolsOrErr(); !ent.IsNotLoaded(err) {
		log.Error(ctx, err, "Loading resource pools for resource type %d failed", obj.ID)
		return es, err
	}
	pools, err := obj.QueryPools().All(ctx)

	if err != nil {
		log.Error(ctx, err, "Loading resource pools for resource type %d failed", obj.ID)
	}

	return pools, err
}

func (r *resourceTypeResolver) PropertyTypes(ctx context.Context, obj *ent.ResourceType) ([]*ent.PropertyType, error) {
	if es, err := obj.Edges.PropertyTypesOrErr(); !ent.IsNotLoaded(err) {
		return es, err
	}
	return obj.QueryPropertyTypes().All(ctx)
}

func (r *tagResolver) Pools(ctx context.Context, obj *ent.Tag) ([]*ent.ResourcePool, error) {
	if es, err := obj.Edges.PoolsOrErr(); !ent.IsNotLoaded(err) {
		log.Error(ctx, err, "Loading resource pools for tag ID %d failed", obj.ID)
		return es, err
	}

	pools, err := obj.QueryPools().All(ctx)

	if err != nil {
		log.Error(ctx, err, "Loading resource pools for tag ID %d failed", obj.ID)
	}

	return pools, err
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

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
