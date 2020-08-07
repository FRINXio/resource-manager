package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"

	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/allocationstrategy"
	"github.com/net-auto/resourceManager/graph/graphql/generated"
	p "github.com/net-auto/resourceManager/pools"
)

func (r *mutationResolver) CreateAllocationStrategy(ctx context.Context, name string, script string) (*ent.AllocationStrategy, error) {
	var client = r.ClientFrom(ctx)
	strat, err := client.AllocationStrategy.Create().SetName(name).SetScript(script).Save(ctx)
	return strat, err
}

func (r *mutationResolver) DeleteAllocationStrategy(ctx context.Context, allocationStrategyID int) (*ent.AllocationStrategy, error) {
	var client = r.ClientFrom(ctx)
	if strat, err := client.AllocationStrategy.Query().
		Where(allocationstrategy.ID(allocationStrategyID)).
		Only(ctx); err != nil {
		return nil, err
	} else {

		if dependentPools, err := strat.QueryPools().All(ctx); len(dependentPools) > 0 && err != nil {
			return nil, fmt.Errorf("Unable to delete, Allocation strategy is still in use")
		}

		if err := client.AllocationStrategy.DeleteOneID(allocationStrategyID).Exec(ctx); err != nil {
			return nil, err
		}
		return strat, nil
	}
}

func (r *mutationResolver) ClaimResource(ctx context.Context, poolName string) (*ent.Resource, error) {
	pool, err := p.ExistingPool(ctx, r.ClientFrom(ctx), poolName)
	if err != nil {
		return nil, err
	}

	return pool.ClaimResource()
}

func (r *mutationResolver) FreeResource(ctx context.Context, input map[string]interface{}, poolName string) (string, error) {
	pool, err := p.ExistingPool(ctx, r.ClientFrom(ctx), poolName)
	if err != nil {
		return err.Error(), err
	}
	err = pool.FreeResource(input)
	if err == nil {
		return "Resource removed successfully", nil
	}

	return err.Error(), err
}

func (r *mutationResolver) CreateSetPool(ctx context.Context, resourceTypeID int, poolName string, poolDealocationSafetyPeriod int, poolValues []map[string]interface{}) (*ent.ResourcePool, error) {
	var client = r.ClientFrom(ctx)

	resType, _ := client.ResourceType.Get(ctx, resourceTypeID)
	_, rp, err := p.NewSetPoolWithMeta(ctx, client, resType, p.ToRawTypes(poolValues), poolName, poolDealocationSafetyPeriod)
	return rp, err
}

func (r *mutationResolver) CreateSingletonPool(ctx context.Context, resourceTypeID int, poolName string, poolValues []map[string]interface{}) (*ent.ResourcePool, error) {
	var client = r.ClientFrom(ctx)

	resType, _ := client.ResourceType.Get(ctx, resourceTypeID)
	if len(poolValues) == 1 {
		_, rp, err := p.NewSingletonPoolWithMeta(ctx, client, resType, p.ToRawTypes(poolValues)[0], poolName)
		return rp, err
	} else {
		return nil, fmt.Errorf("Cannot create singleton pool, no resource provided")
	}
}

func (r *mutationResolver) CreateAllocatingPool(ctx context.Context, resourceTypeID int, poolName string, allocationStrategyID int, poolDealocationSafetyPeriod int) (*ent.ResourcePool, error) {
	var client = r.ClientFrom(ctx)

	resType, errRes := client.ResourceType.Get(ctx, resourceTypeID)
	if errRes != nil {
		return nil, errRes
	}
	allocationStrat, errAlloc := client.AllocationStrategy.Get(ctx, allocationStrategyID)
	if errAlloc != nil {
		return nil, errAlloc
	}

	_, rp, err := p.NewAllocatingPoolWithMeta(ctx, client, resType, allocationStrat, poolName, poolDealocationSafetyPeriod)
	return rp, err
}

func (r *mutationResolver) DeleteResourcePool(ctx context.Context, resourcePoolID int) (string, error) {
	client := r.ClientFrom(ctx)
	pool, err := p.ExistingPoolFromId(ctx, client, resourcePoolID)

	if err != nil {
		return "error", err
	}

	// Do not allow removing pools with allocated resources
	allocatedResources, err2 := pool.QueryResources()
	if len(allocatedResources) > 0 || err2 != nil {
		return "error", errors.New("resource pool has allocated resources, deallocate those first")
	}

	if err := client.ResourcePool.DeleteOneID(resourcePoolID).Exec(ctx); err != nil {
		return "error", err
	} else {
		return "ok", nil
	}
}

func (r *mutationResolver) CreateResourceType(ctx context.Context, resourceName string, resourceProperties map[string]interface{}) (*ent.ResourceType, error) {
	var client = r.ClientFrom(ctx)
	//TODO check error

	var propertyTypes []*ent.PropertyType
	for propName, rawPropType := range resourceProperties {
		var propertyType, err = p.CreatePropertyType(ctx, client, propName, rawPropType)
		if err != nil {
			return nil, err
		}
		propertyTypes = append(propertyTypes, propertyType)
	}

	resType, err2 := client.ResourceType.Create().
		SetName(resourceName).
		AddPropertyTypes(propertyTypes...).
		Save(ctx)
	if err2 != nil {
		return nil, err2
	}

	return resType, nil
}

func (r *mutationResolver) DeleteResourceType(ctx context.Context, resourceTypeID int) (string, error) {
	client := r.ClientFrom(ctx)
	resourceType, err := client.ResourceType.Get(ctx, resourceTypeID)

	if err != nil {
		return "error", err
	}

	pools, err2 := client.ResourceType.QueryPools(resourceType).All(ctx)

	if err2 != nil {
		return "error", err2
	}

	if len(pools) > 0 {
		return "not ok", errors.New("resourceType has pools, can't delete (delete resource pools first)")
	}

	if err := client.ResourcePool.DeleteOneID(resourceTypeID).Exec(ctx); err == nil {
		return "ok", nil
	} else {
		return "error", err
	}
}

func (r *mutationResolver) UpdateResourceTypeName(ctx context.Context, resourceTypeID int, resourceName string) (*ent.ResourceType, error) {
	var client = r.ClientFrom(ctx)
	return client.ResourceType.UpdateOneID(resourceTypeID).SetName(resourceName).Save(ctx)
}

func (r *propertyTypeResolver) Type(ctx context.Context, obj *ent.PropertyType) (string, error) {
	// Just converts enum to string
	return obj.Type.String(), nil
}

func (r *queryResolver) QueryResource(ctx context.Context, input map[string]interface{}, poolName string) (*ent.Resource, error) {
	pool, err := p.ExistingPool(ctx, r.ClientFrom(ctx), poolName)
	if err != nil {
		return nil, err
	}
	return pool.QueryResource(input)
}

func (r *queryResolver) QueryResources(ctx context.Context, poolName string) ([]*ent.Resource, error) {
	pool, err := p.ExistingPool(ctx, r.ClientFrom(ctx), poolName)
	if err != nil {
		return nil, err
	}
	return pool.QueryResources()
}

func (r *queryResolver) QueryAllocationStrategy(ctx context.Context, allocationStrategyName string) (*ent.AllocationStrategy, error) {
	client := r.ClientFrom(ctx)
	return client.AllocationStrategy.Query().Where(allocationstrategy.Name(allocationStrategyName)).Only(ctx)
}

func (r *queryResolver) QueryAllocationStrategies(ctx context.Context) ([]*ent.AllocationStrategy, error) {
	client := r.ClientFrom(ctx)
	return client.AllocationStrategy.Query().All(ctx)
}

func (r *queryResolver) QueryResourceTypes(ctx context.Context) ([]*ent.ResourceType, error) {
	client := r.ClientFrom(ctx)
	return client.ResourceType.Query().All(ctx)
}

func (r *queryResolver) QueryResourcePools(ctx context.Context) ([]*ent.ResourcePool, error) {
	client := r.ClientFrom(ctx)
	return client.ResourcePool.Query().All(ctx)
}

func (r *resourceResolver) Properties(ctx context.Context, obj *ent.Resource) (map[string]interface{}, error) {
	props, err := obj.QueryProperties().WithType().All(ctx)
	if err != nil {
		return nil, err
	}

	return p.PropertiesToMap(props)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

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
