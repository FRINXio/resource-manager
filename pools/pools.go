package pools

import (
	"context"
	log "github.com/net-auto/resourceManager/logging"

	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/resource"
	resourcePool "github.com/net-auto/resourceManager/ent/resourcepool"
	"github.com/pkg/errors"
)

// Pool is a resource provider
type Pool interface {
	ClaimResource(userInput map[string]interface{}) (*ent.Resource, error)
	FreeResource(RawResourceProps) error
	QueryResource(RawResourceProps) (*ent.Resource, error)
	QueryResources() (ent.Resources, error)
	Destroy() error
	ResourceType() (*ent.ResourceType, error)
	Capacity() (float64, float64, error)
}

type poolBase struct {
	*ent.ResourcePool

	ctx    context.Context
	client *ent.Client
}

func (pool poolBase) ResourceType() (*ent.ResourceType, error) {
	return pool.ResourcePool.QueryResourceType().Only(pool.ctx)
}

// SetPool is a pool providing resources from a finite/predefined set of resources
type SetPool struct {
	poolBase
}

// SingletonPool always provides the same resource and never deallocates it
type SingletonPool struct {
	SetPool
}

// AllocatingPool provides resources based on allocation strategy
type AllocatingPool struct {
	SetPool
	invoker ScriptInvoker
}

// Raw representation of resource property values such as ["a": 2, "b": "value"]
type RawResourceProps map[string]interface{}

func newFixedPoolInner(ctx context.Context,
	client *ent.Client,
	resourceType *ent.ResourceType,
	propertyValues []RawResourceProps,
	poolName string,
	description *string,
	poolType resourcePool.PoolType,
	poolDealocationSafetyPeriod int) (*ent.ResourcePool, error) {
	pool, err := client.ResourcePool.Create().
		SetName(poolName).
		SetPoolType(poolType).
		SetNillableDescription(description).
		SetResourceType(resourceType).
		SetDealocationSafetyPeriod(poolDealocationSafetyPeriod).
		Save(ctx)

	if err != nil {
		log.Error(ctx, err, "Unable to create new pool \"%s\". Error creating pool", poolName)
		return nil, errors.Wrapf(err, "Unable to create new pool \"%s\". Error creating pool", poolName)
	}

	// Pre-create all resources
	_, resourceErr := PreCreateResources(ctx, client, propertyValues, pool, resourceType, resource.StatusFree)

	if resourceErr != nil {
		log.Error(ctx, err, "Unable to create new pool \"%s\"", poolName)
		return nil, errors.Wrapf(resourceErr, "Unable to create pool")
	}

	return pool, nil
}

// ExistingPool wraps existing pool entity by ID
func ExistingPoolFromId(
	ctx context.Context,
	client *ent.Client,
	poolId int) (Pool, error) {

	pool, err := client.ResourcePool.Query().
		Where(resourcePool.ID(poolId)).
		Only(ctx)

	if err != nil {
		log.Error(ctx, err, "Unable to find pool ID %d", poolId)
		return nil, errors.Wrapf(err, "Cannot create pool from existing entity")
	}

	return existingPool(ctx, client, pool)
}

func existingPool(
	ctx context.Context,
	client *ent.Client,
	pool *ent.ResourcePool) (Pool, error) {

	switch pool.PoolType {
	case resourcePool.PoolTypeSingleton:
		return &SingletonPool{SetPool{poolBase{pool, ctx, client}}}, nil
	case resourcePool.PoolTypeSet:
		return &SetPool{poolBase{pool, ctx, client}}, nil
	case resourcePool.PoolTypeAllocating:
		wasmer, err := NewWasmerUsingEnvVars()
		if err != nil {
			log.Error(ctx, err, "Unable to create wasmer for %d", pool.ID)
			return nil, err
		}
		return &AllocatingPool{SetPool{poolBase{pool, ctx, client}}, wasmer}, nil
	default:
		err := errors.Errorf("Unknown pool type \"%s\"", pool.PoolType)
		log.Error(ctx, err, "cannot create pool")
		return nil, err
	}
}
