package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/allocationstrategy"
	"github.com/net-auto/resourceManager/graph/graphql/generated"
	p "github.com/net-auto/resourceManager/pools"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func (r *mutationResolver) CreateAllocationStrategy(ctx context.Context, name string, script string) (*ent.AllocationStrategy, error) {
	var client = r.ClientFrom(ctx)
	strat, err := client.AllocationStrategy.Create().SetName(name).SetScript(script).Save(ctx)
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
	if strats, err := client.ResourceType.Query().All(ctx); err != nil {
		return nil, gqlerror.Errorf("Unable to query resource types: %v", err)
	} else {
		return strats, nil
	}
}

func (r *queryResolver) QueryResourcePools(ctx context.Context) ([]*ent.ResourcePool, error) {
	client := r.ClientFrom(ctx)
	if strats, err := client.ResourcePool.Query().All(ctx); err != nil {
		return nil, gqlerror.Errorf("Unable to query resource pools: %v", err)
	} else {
		return strats, nil
	}
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

// Mutation returns generated.MutationResolver implementation.
//  Mutation() function removed in favour of resolver.go.Mutation()

// PropertyType returns generated.PropertyTypeResolver implementation.
func (r *Resolver) PropertyType() generated.PropertyTypeResolver { return &propertyTypeResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Resource returns generated.ResourceResolver implementation.
func (r *Resolver) Resource() generated.ResourceResolver { return &resourceResolver{r} }

type mutationResolver struct{ *Resolver }
type propertyTypeResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type resourceResolver struct{ *Resolver }
