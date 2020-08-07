package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"

	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/resourcepool"
	"github.com/net-auto/resourceManager/graph/graphql/generated"
	p "github.com/net-auto/resourceManager/pools"
)

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

func (r *mutationResolver) CreatePool(ctx context.Context, poolType *resourcepool.PoolType, resourceTypeID int, poolName string, poolValues []map[string]interface{}, allocationScript string) (*ent.ResourcePool, error) {
	var client = r.ClientFrom(ctx)

	resType, _ := client.ResourceType.Get(ctx, resourceTypeID)

	var rawProps = p.ToRawTypes(poolValues)

	if resourcepool.PoolTypeSet == *poolType {
		_, rp, err := p.NewSetPoolWithMeta(ctx, client, resType, rawProps, poolName)
		return rp, err
	} else if resourcepool.PoolTypeSingleton == *poolType {
		if len(rawProps) > 0 {
			_, rp, err := p.NewSingletonPoolWithMeta(ctx, client, resType, rawProps[0], poolName)
			return rp, err
		} else {
			return nil, fmt.Errorf("Cannot create singleton pool, no resource provided")
		}
	}

	return nil, fmt.Errorf("Cannot create singleton pool, something went wrong")
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

func (r *mutationResolver) UpdateResourcePool(ctx context.Context, resourcePoolID int, poolName string, poolValues []map[string]interface{}, allocationScript string) (string, error) {
	client := r.ClientFrom(ctx)

	resourcePool, err := client.ResourcePool.UpdateOneID(resourcePoolID).SetName(poolName).Save(ctx) //TODO also set allocationScript
	if err != nil {
		return "error", err
	}

	resourceType, err2 := resourcePool.QueryResourceType().Only(ctx)
	if err2 != nil {
		return "error", err2
	}

	var rawProps = p.ToRawTypes(poolValues)

	if err := p.PreCreateResources(ctx, client, rawProps, resourcePool, resourceType); err == nil {
		return "ok", nil
	} else {
		return "error", err
	}
}

func (r *mutationResolver) CreateResourceType(ctx context.Context, resourceName string, resourceProperties map[string]interface{}) (*ent.ResourceType, error) {
	var client = r.ClientFrom(ctx)
	//TODO property and resource name the same?
	//TODO check error
	var propType, err = p.CreatePropertyType(ctx, client, resourceName, resourceProperties["type"], resourceProperties["init"])
	if err != nil {
		return nil, err
	}

	resType, err2 := client.ResourceType.Create().
		SetName(resourceName).
		AddPropertyTypes(propType).
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

func (r *mutationResolver) AddResourceTypeProperty(ctx context.Context, resourceTypeID int, resourceProperties map[string]interface{}) (*ent.ResourceType, error) {
	var client = r.ClientFrom(ctx)

	exist, resourceType := p.CheckIfPoolsExist(ctx, client, resourceTypeID)

	if exist {
		return nil, fmt.Errorf("Cannot modify resource type, there are pools attached to it")
	}

	propertyType, err := p.CreatePropertyType(ctx, client, resourceType.Name, resourceProperties["type"], resourceProperties["init"])
	if err != nil {
		return nil, err
	}

	return client.ResourceType.UpdateOneID(resourceTypeID).AddPropertyTypeIDs(propertyType.ID).Save(ctx)
}

func (r *mutationResolver) AddExistingPropertyToResourceType(ctx context.Context, resourceTypeID int, propertyTypeID int) (int, error) {
	var client = r.ClientFrom(ctx)
	if err := client.ResourceType.UpdateOneID(resourceTypeID).AddPropertyTypeIDs(propertyTypeID).Exec(ctx); err == nil {
		return propertyTypeID, nil
	} else {
		return -1, err
	}
}

func (r *mutationResolver) RemoveResourceTypeProperty(ctx context.Context, resourceTypeID int, propertyTypeID int) (*ent.ResourceType, error) {
	var client = r.ClientFrom(ctx)
	exist, _ := p.CheckIfPoolsExist(ctx, client, resourceTypeID)

	if exist {
		return nil, fmt.Errorf("Cannot modify resource type, there are pools attached to it")
	}

	if resourceType, err := client.ResourceType.UpdateOneID(resourceTypeID).RemovePropertyTypeIDs(propertyTypeID).Save(ctx); err == nil {
		if err := client.PropertyType.DeleteOneID(propertyTypeID).Exec(ctx); err != nil {
			return nil, err
		} else {
			return resourceType, nil
		}
	} else {
		return nil, err
	}
}

func (r *propertyTypeResolver) Type(ctx context.Context, obj *ent.PropertyType) (string, error) {
	panic(fmt.Errorf("not implemented"))
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

func (r *queryResolver) QueryResourceTypes(ctx context.Context) ([]*ent.ResourceType, error) {
	client := r.ClientFrom(ctx)
	return client.ResourceType.Query().All(ctx)
}

func (r *queryResolver) QueryResourcePools(ctx context.Context) ([]*ent.ResourcePool, error) {
	client := r.ClientFrom(ctx)
	return client.ResourcePool.Query().All(ctx)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// PropertyType returns generated.PropertyTypeResolver implementation.
func (r *Resolver) PropertyType() generated.PropertyTypeResolver { return &propertyTypeResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type propertyTypeResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
