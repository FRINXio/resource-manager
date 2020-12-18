package pools

import (
	"context"

	"github.com/net-auto/resourceManager/ent/resource"
	"github.com/pkg/errors"

	log "github.com/net-auto/resourceManager/logging"

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
		log.Error(ctx, err, "Unable to create pool")
		return nil, nil, err
	}

	return &SingletonPool{SetPool{poolBase{pool, ctx, client}}}, pool, nil
}

func (pool SingletonPool) ClaimResources(userInput map[string]interface{}, description *string) ([]ent.Resource, error) {
	_, err := pool.client.Resource.Update().
		SetStatus(resource.StatusClaimed).
		Where(resource.HasPoolWith(resourcePool.ID(pool.ID))).
		Save(pool.ctx)

	if err != nil {
		log.Error(pool.ctx, err, "Unable to claim resource in pool ID %d", pool.ID)
		return nil, err
	}

	if description != nil {
		log.Warn(pool.ctx, "Description for a resource from singleton pool will be ignored")
	}
	resource, err := pool.client.Resource.Query().Where(resource.HasPoolWith(resourcePool.ID(pool.ID))).Only(pool.ctx)
	if err != nil {
		return nil, err
	}
	return []ent.Resource{*resource}, nil
}

func (pool SingletonPool) FreeResource(raw RawResourceProps) error {
	pool.client.Resource.Update().
		SetStatus(resource.StatusFree).
		Where(resource.HasPoolWith(resourcePool.ID(pool.ID))).
		Save(pool.ctx)
	return nil
}

func (pool SingletonPool) Capacity() (float64, float64, error) {
	claimedResources, err := pool.QueryResources()

	if err != nil {
		log.Error(pool.ctx, err, "Unable to retrieve resources in pool ID %d", pool.ID)
		return 0, 0, err
	}

	return float64(1 - len(claimedResources)), float64(len(claimedResources)), nil
}

// QueryResource returns always the same resource
func (pool SingletonPool) QueryResource(raw RawResourceProps) (*ent.Resource, error) {
	resources, err := pool.QueryResources()

	if err != nil {
		log.Error(pool.ctx, err, "Unable to retrieve resources in pool ID %d", pool.ID)
		return nil, err
	}

	return resources[0], nil
}

func (pool SingletonPool) QueryResources() (ent.Resources, error) {
	all, err := pool.client.Resource.Query().Where(
		resource.And(
			resource.HasPoolWith(resourcePool.ID(pool.ID)),
			resource.StatusIn(resource.StatusBench, resource.StatusClaimed))).All(pool.ctx)

	if err != nil {
		log.Error(pool.ctx, err, "Unable retrieve resources for pool ID %d", pool.ID)
	}

	return all, err
}

func (pool SingletonPool) Destroy() error {
	claims, errQr := pool.QueryResources()

	if errQr != nil {
		log.Error(pool.ctx, errQr, "Unable to retrieve resources in pool ID %d", pool.ID)
		return errQr
	}

	if len(claims) > 0 {
		log.Warn(pool.ctx, "Unable to delete pool ID %d there are claimed resources", pool.ID)
		return errors.Errorf("Unable to destroy pool \"%s\", there are claimed resources",
			pool.Name)
	}

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
