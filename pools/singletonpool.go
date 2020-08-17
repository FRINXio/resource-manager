package pools

import (
	"context"

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
	poolName string) (Pool, error) {
	pool, _, err := NewSingletonPoolWithMeta(ctx, client, resourceType, propertyValues, poolName)
	return pool, err
}

// NewSingletonPoolWithMeta creates a brand new pool + returns the pools underlying meta information
func NewSingletonPoolWithMeta(
	ctx context.Context,
	client *ent.Client,
	resourceType *ent.ResourceType,
	propertyValues RawResourceProps,
	poolName string) (Pool, *ent.ResourcePool, error) {

	pool, err := newFixedPoolInner(ctx, client, resourceType, []RawResourceProps{propertyValues},
		poolName, resourcePool.PoolTypeSingleton, schema.ResourcePoolDealocationImmediately)

	if err != nil {
		return nil, nil, err
	}

	return &SingletonPool{SetPool{poolBase{pool, ctx, client}}}, pool, nil
}

// ClaimResource returns always the same resource
func (pool SingletonPool) ClaimResource(userInput map[string]interface{}) (*ent.Resource, error) {
	return pool.queryUnclaimedResourceEager()
}

// FreeResource does nothing
func (pool SingletonPool) FreeResource(raw RawResourceProps) error {
	return nil
}

// QueryResource returns always the same resource
func (pool SingletonPool) QueryResource(raw RawResourceProps) (*ent.Resource, error) {
	return pool.QueryResource(raw)
}

// QueryResource returns always the same resource
func (pool SingletonPool) QueryResources() (ent.Resources, error) {
	return pool.QueryResources()
}
