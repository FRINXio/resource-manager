package pools

import (
	"context"
	"github.com/net-auto/resourceManager/ent/resource"

	"github.com/net-auto/resourceManager/ent"
	resourcePool "github.com/net-auto/resourceManager/ent/resourcepool"
	"github.com/net-auto/resourceManager/ent/schema"
)

// NewSingletonPool creates a brand new pool allocating DB entities in the process
func NewSingletonPool(
	ctx context.Context,
	client *ent.Client,
	resourceType *ent.ResourceType,
	propertyValues RawResourceProps,
	poolName string,
	description *string) (Pool, error) {
	pool, _, err := NewSingletonPoolWithMeta(ctx, client, resourceType, propertyValues, poolName, description)
	return pool, err
}

// NewSingletonPoolWithMeta creates a brand new pool + returns the pools underlying meta information
func NewSingletonPoolWithMeta(
	ctx context.Context,
	client *ent.Client,
	resourceType *ent.ResourceType,
	propertyValues RawResourceProps,
	poolName string,
	description *string) (Pool, *ent.ResourcePool, error) {

	pool, err := newFixedPoolInner(ctx, client, resourceType, []RawResourceProps{propertyValues},
		poolName, description, resourcePool.PoolTypeSingleton, schema.ResourcePoolDealocationImmediately)

	if err != nil {
		return nil, nil, err
	}

	return &SingletonPool{SetPool{poolBase{pool, ctx, client}}}, pool, nil
}

func (pool SingletonPool) ClaimResource(userInput map[string]interface{}) (*ent.Resource, error) {
	_, err := pool.client.Resource.Update().
		SetStatus(resource.StatusClaimed).
		Where(resource.HasPoolWith(resourcePool.ID(pool.ID))).
		Save(pool.ctx)

	if err != nil {
		return nil, err
	}

	return pool.client.Resource.Query().Where(resource.HasPoolWith(resourcePool.ID(pool.ID))).Only(pool.ctx)
}

func (pool SingletonPool) FreeResource(raw RawResourceProps) error {
	pool.client.Resource.Update().
		SetStatus(resource.StatusFree).
		Where(resource.HasPoolWith(resourcePool.ID(pool.ID))).
		Save(pool.ctx)
	return nil
}

// TODO add capacity implementation
func (pool SingletonPool) Capacity() (int, error) {
	return 1, nil
}

// QueryResource returns always the same resource
func (pool SingletonPool) QueryResource(raw RawResourceProps) (*ent.Resource, error) {
	resources, err := pool.QueryResources()

	if err != nil {
		return nil, err
	}

	return resources[0], nil
}

func (pool SingletonPool) QueryResources() (ent.Resources, error) {
	return pool.client.Resource.Query().Where(
		resource.And(
			resource.HasPoolWith(resourcePool.ID(pool.ID)),
			resource.StatusIn(resource.StatusBench, resource.StatusClaimed))).All(pool.ctx)
}

func (pool SingletonPool) Destroy() error {
	_, err := pool.client.Resource.Delete().Where(resource.HasPoolWith(resourcePool.ID(pool.ID))).Exec(pool.ctx)

	if err != nil {
		return err
	}

	err = pool.client.ResourcePool.DeleteOneID(pool.ID).Exec(pool.ctx)

	if err != nil {
		return err
	}

	return nil
}