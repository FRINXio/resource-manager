package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"regexp"
	"strconv"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/allocationstrategy"
	"github.com/net-auto/resourceManager/ent/predicate"
	"github.com/net-auto/resourceManager/ent/property"
	"github.com/net-auto/resourceManager/ent/propertytype"
	"github.com/net-auto/resourceManager/ent/resource"
	resourcePool "github.com/net-auto/resourceManager/ent/resourcepool"
	"github.com/net-auto/resourceManager/ent/resourcetype"
	"github.com/net-auto/resourceManager/graph/graphql/generated"
	"github.com/net-auto/resourceManager/graph/graphql/model"
	log "github.com/net-auto/resourceManager/logging"
	p "github.com/net-auto/resourceManager/pools"
	src "github.com/net-auto/resourceManager/pools/allocating_strategies/strategies/generated"
	src2 "github.com/net-auto/resourceManager/pools/allocating_strategies/strategies/src"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// CreateTag is the resolver for the CreateTag field.
func (r *mutationResolver) CreateTag(ctx context.Context, input model.CreateTagInput) (*model.CreateTagPayload, error) {
	var client = r.ClientFrom(ctx)
	tagEnt, err := createTag(ctx, client, input.TagText)

	if err != nil {
		log.Error(ctx, err, "Unable to create new tag")
		return &model.CreateTagPayload{Tag: nil}, gqlerror.Errorf("Unable to create tag: %v", err)
	}

	return &model.CreateTagPayload{Tag: tagEnt}, nil
}

// UpdateTag is the resolver for the UpdateTag field.
func (r *mutationResolver) UpdateTag(ctx context.Context, input model.UpdateTagInput) (*model.UpdateTagPayload, error) {
	var client = r.ClientFrom(ctx)
	tagEnt, err := client.Tag.UpdateOneID(input.TagID).SetTag(input.TagText).Save(ctx)
	if err != nil {
		log.Error(ctx, err, "Unable to update tag ID %d", input.TagID)
		return &model.UpdateTagPayload{Tag: nil}, gqlerror.Errorf("Unable to update tag: %v", err)
	}
	return &model.UpdateTagPayload{Tag: tagEnt}, nil
}

// DeleteTag is the resolver for the DeleteTag field.
func (r *mutationResolver) DeleteTag(ctx context.Context, input model.DeleteTagInput) (*model.DeleteTagPayload, error) {
	var client = r.ClientFrom(ctx)
	err := client.Tag.DeleteOneID(input.TagID).Exec(ctx)
	if err != nil {
		log.Error(ctx, err, "Unable to delete tag ID %d", input.TagID)
		return &model.DeleteTagPayload{TagID: input.TagID}, gqlerror.Errorf("Unable to delete tag: %v", err)
	}

	return &model.DeleteTagPayload{TagID: input.TagID}, nil
}

// TagPool is the resolver for the TagPool field.
func (r *mutationResolver) TagPool(ctx context.Context, input model.TagPoolInput) (*model.TagPoolPayload, error) {
	var client = r.ClientFrom(ctx)
	tag, err := client.Tag.UpdateOneID(input.TagID).AddPoolIDs(input.PoolID).Save(ctx)
	if err != nil {
		log.Error(ctx, err, "Unable to tag pool ID %d", input.PoolID)
		return &model.TagPoolPayload{Tag: nil}, gqlerror.Errorf("Unable to tag pool: %v", err)
	}
	return &model.TagPoolPayload{Tag: tag}, nil
}

// UntagPool is the resolver for the UntagPool field.
func (r *mutationResolver) UntagPool(ctx context.Context, input model.UntagPoolInput) (*model.UntagPoolPayload, error) {
	var client = r.ClientFrom(ctx)
	tag, err := client.Tag.UpdateOneID(input.TagID).RemovePoolIDs(input.PoolID).Save(ctx)
	if err != nil {
		log.Error(ctx, err, "Unable to un-tag pool ID %d", input.PoolID)
		return &model.UntagPoolPayload{Tag: nil}, gqlerror.Errorf("Unable to un-tag pool: %v", err)
	}
	return &model.UntagPoolPayload{Tag: tag}, nil
}

// CreateAllocationStrategy is the resolver for the CreateAllocationStrategy field.
func (r *mutationResolver) CreateAllocationStrategy(ctx context.Context, input *model.CreateAllocationStrategyInput) (*model.CreateAllocationStrategyPayload, error) {
	var client = r.ClientFrom(ctx)

	var propertyTypes []*ent.PropertyType
	var strat *ent.AllocationStrategy
	var err error
	if input.ExpectedPoolPropertyTypes != nil {
		for propName, rawPropType := range input.ExpectedPoolPropertyTypes {
			var propertyType, err = p.CreatePropertyType(ctx, client, propName, rawPropType)
			if err != nil {
				return &model.CreateAllocationStrategyPayload{Strategy: nil}, gqlerror.Errorf("Unable to create expected resource type: %v", err)
			}
			propertyTypes = append(propertyTypes, propertyType)
		}

		strat, err = client.AllocationStrategy.Create().
			SetName(input.Name).
			SetNillableDescription(input.Description).
			SetScript(input.Script).
			SetLang(input.Lang).
			AddPoolPropertyTypes(propertyTypes...).
			Save(ctx)
		if err != nil {
			log.Error(ctx, err, "Unable create a new allocation strategy")
			return &model.CreateAllocationStrategyPayload{Strategy: nil}, gqlerror.Errorf("Unable to create strategy: %v", err)
		}
	} else {
		strat, err = client.AllocationStrategy.Create().
			SetName(input.Name).
			SetNillableDescription(input.Description).
			SetScript(input.Script).
			SetLang(input.Lang).
			Save(ctx)
		if err != nil {
			log.Error(ctx, err, "Unable create a new allocation strategy")
			return &model.CreateAllocationStrategyPayload{Strategy: nil}, gqlerror.Errorf("Unable to create strategy: %v", err)
		}
	}

	return &model.CreateAllocationStrategyPayload{Strategy: strat}, nil
}

// DeleteAllocationStrategy is the resolver for the DeleteAllocationStrategy field.
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

// TestAllocationStrategy is the resolver for the TestAllocationStrategy field.
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

	parsedOutputFromStrat, stdErr, err := p.InvokeAllocationStrategy(ctx, wasmer, strat, userInput, resourcePool, currentResources, poolPropertiesMaps, functionName)
	if err != nil {
		log.Error(ctx, err, "Error while running script on pool \"%s\" strategy ID %d", resourcePool.ResourcePoolName, allocationStrategyID)
		return nil, gqlerror.Errorf("Error while running the script: %v", err)
	}
	result := make(map[string]interface{})
	result["stdout"] = parsedOutputFromStrat
	result["stderr"] = stdErr
	return result, nil
}

// ClaimResource is the resolver for the ClaimResource field.
func (r *mutationResolver) ClaimResource(ctx context.Context, poolID int, description *string, userInput map[string]interface{}) (*ent.Resource, error) {
	pool, err := p.ExistingPoolFromId(ctx, r.ClientFrom(ctx), poolID)
	if err != nil {
		return nil, gqlerror.Errorf("Resource pool is not existing, for you to be able to claim resource: %v", err)
	}

	return ClaimResource(pool, userInput, description, nil)
}

// ClaimResourceWithAltID is the resolver for the ClaimResourceWithAltId field.
func (r *mutationResolver) ClaimResourceWithAltID(ctx context.Context, poolID int, description *string, userInput map[string]interface{}, alternativeID map[string]interface{}) (*ent.Resource, error) {
	pool, err := p.ExistingPoolFromId(ctx, r.ClientFrom(ctx), poolID)
	if err != nil {
		return nil, gqlerror.Errorf("Resource pool is not existing, for you to be able to claim resource: %v", err)
	}

	return ClaimResource(pool, userInput, description, alternativeID)
}

// FreeResource is the resolver for the FreeResource field.
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

// CreateSetPool is the resolver for the CreateSetPool field.
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

// CreateNestedSetPool is the resolver for the CreateNestedSetPool field.
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

// CreateSingletonPool is the resolver for the CreateSingletonPool field.
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

// CreateNestedSingletonPool is the resolver for the CreateNestedSingletonPool field.
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

// CreateAllocatingPool is the resolver for the CreateAllocatingPool field.
func (r *mutationResolver) CreateAllocatingPool(ctx context.Context, input *model.CreateAllocatingPoolInput) (*model.CreateAllocatingPoolPayload, error) {
	var client = r.ClientFrom(ctx)
	emptyRetVal := model.CreateAllocatingPoolPayload{Pool: nil}

	allocationStrat, errAlloc := client.AllocationStrategy.Get(ctx, input.AllocationStrategyID)
	if errAlloc != nil {
		log.Error(ctx, errAlloc, "Unable to retrieve allocation strategy for pool (strategy ID: %d)", input.AllocationStrategyID)
		return &emptyRetVal, gqlerror.Errorf("Unable to create pool: %v", errAlloc)
	}

	requiredPoolProperties, err := r.Query().QueryRequiredPoolProperties(ctx, allocationStrat.Name)
	if err != nil {
		log.Error(ctx, errAlloc, "Unable to retrieve required pool properties for pool (strategy ID: %d)", input.AllocationStrategyID)
		return &emptyRetVal, gqlerror.Errorf("Unable to create pool: %v", errAlloc)
	}

	if requiredPoolProperties != nil {
		for _, requiredPoolProperty := range requiredPoolProperties {
			if input.PoolPropertyTypes != nil {
				propertyExists := false
				inputPoolPropertiesMap := []map[string]interface{}{input.PoolPropertyTypes}
				for _, inputPoolPropertyMap := range inputPoolPropertiesMap {
					if inputPoolPropertyMap[requiredPoolProperty.Name] != nil &&
						inputPoolPropertyMap[requiredPoolProperty.Name].(string) == requiredPoolProperty.Type.String() {
						propertyExists = true
					}
				}
				if propertyExists != true {
					log.Error(ctx, nil, "In pool properties input missed property: %s - %s", requiredPoolProperty.Name, requiredPoolProperty.Type)
					return &emptyRetVal, gqlerror.Errorf("In pool properties input missed property: %s - %s", requiredPoolProperty.Name, requiredPoolProperty.Type)
				}
			}
		}
	}

	if input.PoolProperties != nil && (allocationStrat.Name == "ipv4_prefix" || allocationStrat.Name == "ipv4" || allocationStrat.Name == "ipv6_prefix" || allocationStrat.Name == "ipv6") {
		prefix, ok := input.PoolProperties["prefix"]
		if !ok {
			log.Error(ctx, nil, "In pool properties input missed property: prefix of type int")
			return &emptyRetVal, gqlerror.Errorf("In pool properties input missed property: prefix of type int")
		}

		subnet, ok := input.PoolProperties["subnet"]
		if !ok {
			log.Error(ctx, nil, "In pool properties input missed property: subnet of type bool")
			return &emptyRetVal, gqlerror.Errorf("In pool properties input missed property: subnet of type bool")
		}

		address, ok := input.PoolProperties["address"]
		if !ok {
			log.Error(ctx, nil, "In pool properties input missed property: address of type string")
			return &emptyRetVal, gqlerror.Errorf("In pool properties input missed property: address of type string")
		}

		isIpv4Ok, _ := regexp.Match(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}$`, []byte(address.(string)))
		if allocationStrat.Name == "ipv4_prefix" || allocationStrat.Name == "ipv4" {
			if !isIpv4Ok {
				log.Error(ctx, nil, "When creating allocating pool the invalid ipv4 address was provided")
				return &emptyRetVal, gqlerror.Errorf("When creating allocating pool the invalid ipv4 address was provided")
			}

			if networkAddress, err := src2.InetAton(address.(string)); err != nil || networkAddress%2 != 0 {
				log.Error(ctx, nil, "When creating allocating pool user needs to provide network address (the last bit needs to be even-numbered)")
				return &emptyRetVal, gqlerror.Errorf("When creating allocating pool user needs to provide network address (the last bit needs to be even-numbered)")
			}
		}

		if _, isIpv6Ok := src.Ipv6InetAton(address.(string)); (allocationStrat.Name == "ipv6_prefix" || allocationStrat.Name == "ipv6") && isIpv6Ok != nil {
			log.Error(ctx, nil, "When creating allocating pool the invalid ipv6 address was provided")
			return &emptyRetVal, gqlerror.Errorf("When creating allocating pool the invalid ipv6 address was provided")
		}

		prefixValue, isPrefixOk := src.NumberToInt(prefix)
		isSubnet, isSubnetOk := subnet.(bool)

		if isPrefixOk == nil && isSubnetOk && isSubnet && (prefixValue.(int) == 32 || prefixValue.(int) == 31 || prefixValue.(int) == 127 || prefixValue.(int) == 128) {
			return &emptyRetVal, gqlerror.Errorf("Unable to create pool, because you cannot create resource pool of prefix /%s together with subnet true", strconv.Itoa(prefixValue.(int)))
		}
	}

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

	_, rp, err := p.NewAllocatingPoolWithMeta(ctx, client, resType, allocationStrat,
		input.PoolName, input.Description, input.PoolDealocationSafetyPeriod, poolProperties)

	if err != nil {
		return &emptyRetVal, err
	}

	if err := createTagsAndTagPool(ctx, client, rp, input.Tags); err != nil {
		log.Error(ctx, err, "Unable to tag the pool with tags: %v", input.Tags)
		return nil, err
	}

	if err != nil {
		return &emptyRetVal, gqlerror.Errorf("Unable to create pool: %v", err)
	}

	return &model.CreateAllocatingPoolPayload{Pool: rp}, err
}

// CreateNestedAllocatingPool is the resolver for the CreateNestedAllocatingPool field.
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

// DeleteResourcePool is the resolver for the DeleteResourcePool field.
func (r *mutationResolver) DeleteResourcePool(ctx context.Context, input model.DeleteResourcePoolInput) (*model.DeleteResourcePoolPayload, error) {
	client := r.ClientFrom(ctx)
	retVal := model.DeleteResourcePoolPayload{ResourcePoolID: input.ResourcePoolID}

	pool, err := p.ExistingPoolFromId(ctx, client, input.ResourcePoolID)
	if err != nil {
		return &retVal, gqlerror.Errorf("Unable to retrieve pool: %v", err)
	}

	poolName, err := p.PoolName(ctx, client, input.ResourcePoolID)
	if err != nil {
		return &retVal, gqlerror.Errorf("Unable to retrieve pool name: %v", err)
	}

	errDp := pool.Destroy()

	if errDp != nil {
		return &retVal, gqlerror.Errorf("Unable to delete pool: %v", errDp)
	}

	resourceTypes, err := client.ResourceType.Query().Where(resourcetype.Name(poolName + "-ResourceType")).All(ctx)
	if err != nil {
		log.Error(ctx, err, "Unable to retrieve resource type from pool id %d", input.ResourcePoolID)
		return &retVal, gqlerror.Errorf("Unable to delete resource type - cannot find by pool ID %d: %v", input.ResourcePoolID, err)
	}
	for _, resourceType := range resourceTypes {

		rp, err := r.DeleteResourceType(ctx, model.DeleteResourceTypeInput{
			ResourceTypeID: resourceType.ID,
		})
		if err != nil {
			log.Error(ctx, err, "Unable to delete resource type ID %d", rp.ResourceTypeID)
			return &retVal, gqlerror.Errorf("Unable to delete resource type ID %d: %v", rp.ResourceTypeID, err)
		}
	}

	return &retVal, nil
}

// CreateResourceType is the resolver for the CreateResourceType field.
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

// DeleteResourceType is the resolver for the DeleteResourceType field.
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

// UpdateResourceTypeName is the resolver for the UpdateResourceTypeName field.
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

// UpdateResourceAltID is the resolver for the UpdateResourceAltId field.
func (r *mutationResolver) UpdateResourceAltID(ctx context.Context, input map[string]interface{}, poolID int, alternativeID map[string]interface{}) (*ent.Resource, error) {
	pool, err := p.ExistingPoolFromId(ctx, r.ClientFrom(ctx), poolID)
	if err != nil {
		return nil, gqlerror.Errorf("Unable to retrieve pool: %v", err)
	}
	queryResource, err := pool.QueryResource(input)
	if err != nil {
		return nil, gqlerror.Errorf("Unable to query resource: %v", err)
	}
	var client = r.ClientFrom(ctx)
	for k, v := range alternativeID {
		queryResource.AlternateID[k] = v
	}
	if _, err := client.Resource.UpdateOne(queryResource).SetAlternateID(queryResource.AlternateID).Save(ctx); err != nil {
		log.Error(ctx, err, "Unable to update resource alternative ID %v", alternativeID)
		return queryResource, gqlerror.Errorf("Unable to update resource alternative ID: %v", err)
	}
	return queryResource, nil
}

// ID is the resolver for the ID field.
func (r *outputCursorResolver) ID(ctx context.Context, obj *ent.Cursor) (string, error) {
	//this will never be called because ent.Cursor will use its msgpack annotation
	return "", nil
}

// Type is the resolver for the Type field.
func (r *propertyTypeResolver) Type(ctx context.Context, obj *ent.PropertyType) (string, error) {
	// Just converts enum to string
	return obj.Type.String(), nil
}

// QueryPoolCapacity is the resolver for the QueryPoolCapacity field.
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

// QueryPoolTypes is the resolver for the QueryPoolTypes field.
func (r *queryResolver) QueryPoolTypes(ctx context.Context) ([]resourcePool.PoolType, error) {
	poolTypes := []resourcePool.PoolType{
		resourcePool.PoolTypeSingleton,
		resourcePool.PoolTypeSet,
		resourcePool.PoolTypeAllocating}
	return poolTypes, nil
}

// QueryResource is the resolver for the QueryResource field.
func (r *queryResolver) QueryResource(ctx context.Context, input map[string]interface{}, poolID int) (*ent.Resource, error) {
	pool, err := p.ExistingPoolFromId(ctx, r.ClientFrom(ctx), poolID)
	if err != nil {
		return nil, gqlerror.Errorf("Unable to query resource: %v", err)
	}
	return pool.QueryResource(input)
}

// QueryResources is the resolver for the QueryResources field.
func (r *queryResolver) QueryResources(ctx context.Context, poolID int, first *int, last *int, before *string, after *string) (*ent.ResourceConnection, error) {
	pool, err := p.ExistingPoolFromId(ctx, r.ClientFrom(ctx), poolID)
	if err != nil {
		return nil, gqlerror.Errorf("Unable to query resources: %v", err)
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
	return pool.QueryPaginatedResources(first, last, afterCursor, beforeCursor)
}

// QueryResourcesByAltID is the resolver for the QueryResourcesByAltId field.
func (r *queryResolver) QueryResourcesByAltID(ctx context.Context, input map[string]interface{}, poolID *int, first *int, last *int, before *string, after *string) (*ent.ResourceConnection, error) {
	typeFixedAlternativeId, err := p.ConvertValuesToFloat64(ctx, input)
	if err != nil {
		return nil, gqlerror.Errorf("Unable to process input data %v", err)
	}

	if poolID != nil {
		_, err := p.ExistingPoolFromId(ctx, r.ClientFrom(ctx), *poolID)
		if err != nil {
			return nil, gqlerror.Errorf("Unable to query resources: %v", err)
		}
	}

	res, err := QueryResourcesByAltId(ctx, r.ClientFrom(ctx), typeFixedAlternativeId, poolID, first, last, before, after)

	if res != nil {
		return res, nil
	}

	log.Warn(ctx, "There is not such resources with alternative ID %v", input)
	return nil, gqlerror.Errorf("Unable to query resources: %v", err)
}

// QueryAllocationStrategy is the resolver for the QueryAllocationStrategy field.
func (r *queryResolver) QueryAllocationStrategy(ctx context.Context, allocationStrategyID int) (*ent.AllocationStrategy, error) {
	client := r.ClientFrom(ctx)
	if strats, err := client.AllocationStrategy.Query().Where(allocationstrategy.ID(allocationStrategyID)).Only(ctx); err != nil {
		log.Error(ctx, err, "Unable to retrieve allocation strategy ID %d", allocationStrategyID)
		return nil, gqlerror.Errorf("Unable to query strategy: %v", err)
	} else {
		return strats, nil
	}
}

// QueryAllocationStrategies is the resolver for the QueryAllocationStrategies field.
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

// QueryResourceTypes is the resolver for the QueryResourceTypes field.
func (r *queryResolver) QueryResourceTypes(ctx context.Context, byName *string) ([]*ent.ResourceType, error) {
	client := r.ClientFrom(ctx)
	query := client.ResourceType.Query()

	// Filter out pool properties that are stored in resource type table and pool names
	query = query.Where(resourcetype.Not(resourcetype.HasPoolProperties())).
		Where(resourcetype.Not(resourcetype.NameContains("-ResourceType")))

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

// QueryRequiredPoolProperties is the resolver for the QueryRequiredPoolProperties field.
func (r *queryResolver) QueryRequiredPoolProperties(ctx context.Context, allocationStrategyName string) ([]*ent.PropertyType, error) {
	allocationStrategy, err := r.ClientFrom(ctx).AllocationStrategy.Query().Where(allocationstrategy.Name(allocationStrategyName)).Only(ctx)
	if err != nil {
		log.Error(ctx, err, "Unable to retrieve required allocation strategy by name: %s", allocationStrategyName)
		return nil, gqlerror.Errorf("Unable to retrieve required allocation strategy by name: %s", allocationStrategyName)
	}

	requiredPropertyTypes, err := r.ClientFrom(ctx).PropertyType.Query().Where(func(s *sql.Selector) {
		s.Where(sql.EQ(allocationstrategy.PoolPropertyTypesColumn, allocationStrategy.ID))
	}).All(ctx)
	if err != nil {
		log.Error(ctx, err, "Unable to retrieve required pool properties for allocation strategy: %s", allocationStrategyName)
		return nil, gqlerror.Errorf("Unable to retrieve required pool properties for allocation strategy: %s", allocationStrategyName)
	}
	return requiredPropertyTypes, nil
}

// QueryResourcePool is the resolver for the QueryResourcePool field.
func (r *queryResolver) QueryResourcePool(ctx context.Context, poolID int) (*ent.ResourcePool, error) {
	rp, err := r.ClientFrom(ctx).ResourcePool.Get(ctx, poolID)

	if err != nil {
		log.Error(ctx, err, "Unable to retrieve resource pool")
	}

	return rp, err
}

// QueryEmptyResourcePools is the resolver for the QueryEmptyResourcePools field.
func (r *queryResolver) QueryEmptyResourcePools(ctx context.Context, resourceTypeID *int, first *int, last *int, before *ent.Cursor, after *ent.Cursor) (*ent.ResourcePoolConnection, error) {
	client := r.ClientFrom(ctx)
	query := client.ResourcePool.Query()

	if resourceTypeID != nil {
		query.Where(resourcePool.HasResourceTypeWith(resourcetype.ID(*resourceTypeID))).Where(resourcePool.Not(resourcePool.HasClaims()))
	} else {
		query.Where(resourcePool.Not(resourcePool.HasClaims()))
	}

	if resourcePools, err := query.Paginate(ctx, after, first, before, last); err != nil {
		log.Error(ctx, err, "Unable to retrieve resource pools")
		return nil, gqlerror.Errorf("Unable to query resource pools: %v", err)
	} else {
		return resourcePools, nil
	}
}

// QueryResourcePools is the resolver for the QueryResourcePools field.
func (r *queryResolver) QueryResourcePools(ctx context.Context, resourceTypeID *int, tags *model.TagOr, first *int, last *int, before *ent.Cursor, after *ent.Cursor, filterByResources map[string]interface{}) (*ent.ResourcePoolConnection, error) {
	client := r.ClientFrom(ctx)
	query := client.ResourcePool.Query()

	if resourceTypeID != nil {
		query.Where(resourcePool.HasResourceTypeWith(resourcetype.ID(*resourceTypeID)))
	}

	if filterByResources != nil {
		filteredOutResourceIDs, err := FilterResourcePoolByAllocatedResources(ctx, query, filterByResources)

		if err != nil {
			log.Error(ctx, err, "Unable to retrieve resource pools")
			return nil, gqlerror.Errorf("Unable to query resource pools: %v", err)
		}

		query.Where(resourcePool.HasClaimsWith(resource.HasPropertiesWith(property.IDIn(filteredOutResourceIDs...))))
	}

	if tags != nil {
		// TODO make sure all tags exist
		query.Where(resourcePoolTagPredicate(tags))
	}

	if resourcePools, err := query.Paginate(ctx, after, first, before, last); err != nil {
		log.Error(ctx, err, "Unable to retrieve resource pools")
		return nil, gqlerror.Errorf("Unable to query resource pools: %v", err)
	} else {
		return resourcePools, nil
	}
}

// QueryRecentlyActiveResources is the resolver for the QueryRecentlyActiveResources field.
func (r *queryResolver) QueryRecentlyActiveResources(ctx context.Context, fromDatetime string, toDatetime *string, first *int, last *int, before *string, after *string) (*ent.ResourceConnection, error) {
	client := r.ClientFrom(ctx)
	query := client.Resource.Query()
	dateFrom, err := time.Parse("2006-01-02-15", fromDatetime)
	if err != nil {
		log.Error(ctx, err, "Unable to parse date from: "+fromDatetime+". Must be in format: YYYY-MM-DD-hh.")
		return nil, gqlerror.Errorf("Unable to parse date from: "+fromDatetime+
			". Must be in format: YYYY-MM-DD-hh. Error: %v", err)
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

	var resources *ent.ResourceConnection

	if toDatetime != nil && len(*toDatetime) != 0 {
		dateTo, err := time.Parse("2006-01-02-15", *toDatetime)
		if err != nil {
			log.Error(ctx, err, "Unable to parse date to: "+*toDatetime+". Must be in format: YYYY-MM-DD-hh.")
			return nil, gqlerror.Errorf("Unable to parse date to: "+*toDatetime+
				". Must be in format: YYYY-MM-DD-hh. Error: %v", err)
		}
		resources, err = query.Where(resource.And(resource.UpdatedAtGTE(dateFrom), resource.UpdatedAtLTE(dateTo))).
			Paginate(ctx, afterCursor, first, beforeCursor, last)
		if err != nil {
			return nil, err
		}
	} else {
		currentDate := time.Now()
		resources, err = query.Where(resource.And(resource.UpdatedAtGTE(dateFrom), resource.UpdatedAtLTE(currentDate))).
			Paginate(ctx, afterCursor, first, beforeCursor, last)
		if err != nil {
			return nil, err
		}
	}
	return resources, nil
}

// QueryResourcePoolHierarchyPath is the resolver for the QueryResourcePoolHierarchyPath field.
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

// QueryRootResourcePools is the resolver for the QueryRootResourcePools field.
func (r *queryResolver) QueryRootResourcePools(ctx context.Context, resourceTypeID *int, tags *model.TagOr, first *int, last *int, before *ent.Cursor, after *ent.Cursor, filterByResources map[string]interface{}) (*ent.ResourcePoolConnection, error) {
	client := r.ClientFrom(ctx)
	query := client.ResourcePool.
		Query().
		Where(resourcePool.Not(resourcePool.HasParentResource()))

	if resourceTypeID != nil {
		query.Where(resourcePool.HasResourceTypeWith(resourcetype.ID(*resourceTypeID)))
	}

	if filterByResources != nil {
		filteredOutResourceIDs, err := FilterResourcePoolByAllocatedResources(ctx, query, filterByResources)

		if err != nil {
			log.Error(ctx, err, "Unable to retrieve resource pools")
			return nil, gqlerror.Errorf("Unable to query resource pools: %v", err)
		}

		query.Where(resourcePool.HasClaimsWith(resource.HasPropertiesWith(property.IDIn(filteredOutResourceIDs...))))
	}

	if tags != nil {
		// TODO make sure all tags exist
		query.Where(resourcePoolTagPredicate(tags))
	}

	if resourcePools, err := query.Paginate(ctx, after, first, before, last); err != nil {
		log.Error(ctx, err, "Unable to retrieve root resource pools")
		return nil, gqlerror.Errorf("Unable to query resource pools: %v", err)
	} else {
		return resourcePools, nil
	}
}

// QueryLeafResourcePools is the resolver for the QueryLeafResourcePools field.
func (r *queryResolver) QueryLeafResourcePools(ctx context.Context, resourceTypeID *int, tags *model.TagOr, first *int, last *int, before *ent.Cursor, after *ent.Cursor, filterByResources map[string]interface{}) (*ent.ResourcePoolConnection, error) {
	client := r.ClientFrom(ctx)
	query := client.ResourcePool.
		Query().
		Where(resourcePool.HasParentResource()).
		Where(resourcePool.Not(resourcePool.HasClaimsWith(resource.HasNestedPool())))

	if resourceTypeID != nil {
		query.Where(resourcePool.HasResourceTypeWith(resourcetype.ID(*resourceTypeID)))
	}

	if filterByResources != nil {
		filteredOutResourceIDs, err := FilterResourcePoolByAllocatedResources(ctx, query, filterByResources)

		if err != nil {
			log.Error(ctx, err, "Unable to retrieve resource pools")
			return nil, gqlerror.Errorf("Unable to query resource pools: %v", err)
		}

		query.Where(resourcePool.HasClaimsWith(resource.HasPropertiesWith(property.IDIn(filteredOutResourceIDs...))))
	}

	if tags != nil {
		// TODO make sure all tags exist
		query.Where(resourcePoolTagPredicate(tags))
	}

	if resourcePools, err := query.Paginate(ctx, after, first, before, last); err != nil {
		log.Error(ctx, err, "Unable to retrieve leaf resource pools")
		return nil, gqlerror.Errorf("Unable to query resource pools: %v", err)
	} else {
		return resourcePools, nil
	}
}

// SearchPoolsByTags is the resolver for the SearchPoolsByTags field.
func (r *queryResolver) SearchPoolsByTags(ctx context.Context, tags *model.TagOr, first *int, last *int, before *ent.Cursor, after *ent.Cursor) (*ent.ResourcePoolConnection, error) {
	var client = r.ClientFrom(ctx)

	var condition predicate.ResourcePool

	if tags != nil {
		// TODO make sure all tags exist
		condition = resourcePoolTagPredicate(tags)
	}

	var (
		matchedPools *ent.ResourcePoolConnection
		err          error
	)
	if condition == nil {
		matchedPools, err = client.ResourcePool.Query().Paginate(ctx, after, first, before, last)
	} else {
		matchedPools, err = client.ResourcePool.Query().Where(condition).Paginate(ctx, after, first, before, last)
	}

	if err != nil {
		log.Error(ctx, err, "Unable to retrieve pools by tags")
		return nil, gqlerror.Errorf("Unable to query pools: %v", err)
	}
	return matchedPools, nil
}

// QueryTags is the resolver for the QueryTags field.
func (r *queryResolver) QueryTags(ctx context.Context) ([]*ent.Tag, error) {
	var client = r.ClientFrom(ctx)
	tags, err := client.Tag.Query().All(ctx)
	if err != nil {
		log.Error(ctx, err, "Unable to retrieve tags")
		return nil, gqlerror.Errorf("Unable to query tags: %v", err)
	}
	return tags, nil
}

// Node is the resolver for the node field.
func (r *queryResolver) Node(ctx context.Context, id int) (ent.Noder, error) {
	var client = r.ClientFrom(ctx)
	node, err := client.Noder(ctx, id)

	if err != nil {
		log.Error(ctx, err, "Unable to retrieve node with ID %d", id)
	}

	return node, err
}

// ParentPool is the resolver for the ParentPool field.
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

// Properties is the resolver for the Properties field.
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

// AlternativeID is the resolver for the AlternativeId field.
func (r *resourceResolver) AlternativeID(ctx context.Context, obj *ent.Resource) (map[string]interface{}, error) {
	return obj.AlternateID, nil
}

// Capacity is the resolver for the Capacity field.
func (r *resourcePoolResolver) Capacity(ctx context.Context, obj *ent.ResourcePool) (*model.PoolCapacityPayload, error) {
	return r.Query().QueryPoolCapacity(ctx, obj.ID)
}

// PoolProperties is the resolver for the PoolProperties field.
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

// Resources is the resolver for the Resources field.
func (r *resourcePoolResolver) Resources(ctx context.Context, obj *ent.ResourcePool) ([]*ent.Resource, error) {
	resources, err := p.GetResourceFromPool(ctx, obj)

	if err != nil {
		log.Error(ctx, err, "Unable to retrieve resources for pool with ID %d", obj.ID)
	}

	return resources, err
}

// AllocatedResources is the resolver for the allocatedResources field.
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

type mutationResolver struct{ *Resolver }
type outputCursorResolver struct{ *Resolver }
type propertyTypeResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type resourceResolver struct{ *Resolver }
type resourcePoolResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
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

type resourceTypeResolver struct{ *Resolver }
type tagResolver struct{ *Resolver }
