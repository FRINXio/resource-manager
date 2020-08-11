package pools

import (
	"context"

	"github.com/net-auto/resourceManager/ent"
	allocationStrategy "github.com/net-auto/resourceManager/ent/allocationstrategy"
	"github.com/net-auto/resourceManager/ent/resource"
	resourcePool "github.com/net-auto/resourceManager/ent/resourcepool"
	"github.com/pkg/errors"
)

// NewAllocatingPool creates a brand new pool allocating DB entities in the process
func NewAllocatingPool(
	ctx context.Context,
	client *ent.Client,
	resourceType *ent.ResourceType,
	allocationStrategy *ent.AllocationStrategy,
	poolName string,
	poolDealocationSafetyPeriod int) (Pool, error) {
	pool, _, err := NewAllocatingPoolWithMeta(ctx, client, resourceType, allocationStrategy, poolName, poolDealocationSafetyPeriod)
	return pool, err
}

// NewAllocatingPoolWithMeta creates a brand new pool + returns the pools underlying meta information
func NewAllocatingPoolWithMeta(
	ctx context.Context,
	client *ent.Client,
	resourceType *ent.ResourceType,
	allocationStrategy *ent.AllocationStrategy,
	poolName string,
	poolDealocationSafetyPeriod int) (Pool, *ent.ResourcePool, error) {

	wasmer, err := NewWasmerUsingEnvVars()
	if err != nil {
		return nil, nil, errors.Wrap(err, "Cannot create resource pool")
	}
	return newAllocatingPoolWithMetaInternal(
		ctx, client, resourceType, allocationStrategy, poolName, wasmer, poolDealocationSafetyPeriod)
}

func newAllocatingPoolWithMetaInternal(
	ctx context.Context,
	client *ent.Client,
	resourceType *ent.ResourceType,
	allocationStrategy *ent.AllocationStrategy,
	poolName string,
	invoker ScriptInvoker,
	poolDealocationSafetyPeriod int) (Pool, *ent.ResourcePool, error) {

	pool, err := client.ResourcePool.Create().
		SetName(poolName).
		SetPoolType(resourcePool.PoolTypeAllocating).
		SetResourceType(resourceType).
		SetAllocationStrategy(allocationStrategy).
		SetDealocationSafetyPeriod(poolDealocationSafetyPeriod).
		Save(ctx)

	if err != nil {
		return nil, nil, errors.Wrap(err, "Cannot create resource pool")
	}

	return &AllocatingPool{
			SetPool{poolBase{pool, ctx, client}},
			invoker},
		pool, nil
}

// Destroy removes the pool from DB if there are no more claims
func (pool AllocatingPool) Destroy() error {
	// Check if there are no more claims
	claims, err := pool.QueryResources()
	if err != nil {
		return err
	}

	if len(claims) > 0 {
		return errors.Errorf("Unable to destroy pool \"%s\", there are claimed resources",
			pool.Name)
	}

	// Delete pool itself
	err = pool.client.ResourcePool.DeleteOne(pool.ResourcePool).Exec(pool.ctx)
	if err != nil {
		return errors.Wrapf(err, "Cannot destroy pool \"%s\"", pool.Name)
	}

	return nil
}

func (pool AllocatingPool) AllocationStrategy() (*ent.AllocationStrategy, error) {
	return pool.ResourcePool.QueryAllocationStrategy().Only(pool.ctx)
}

// ClaimResource allocates the next available resource
func (pool AllocatingPool) ClaimResource() (*ent.Resource, error) {

	strat, err := pool.AllocationStrategy()
	if err != nil {
		return nil, errors.Wrapf(err,
			"Unable to claim resource from pool \"%s\", allocation strategy loading error ", pool.Name)
	}
	resourceType, err := pool.ResourceType()
	if err != nil {
		return nil, errors.Wrapf(err,
			"Unable to claim resource from pool \"%s\", resource type loading error ", pool.Name)
	}

	parsedOutputFromStrat, err := pool.invokeAllocationStrategy(strat)
	if err != nil {
		return nil, errors.Wrapf(err,
			"Unable to claim resource from pool \"%s\", allocation strategy \"%s\" failed", pool.Name, strat.Name)
	}
	created, err := PreCreateResources(pool.ctx, pool.client, []RawResourceProps{parsedOutputFromStrat},
		pool.ResourcePool, resourceType, resource.StatusClaimed)
	if len(created) > 1 {
		return nil, errors.Errorf(
			"Unable to claim resource from pool \"%s\", allocation strategy \"%s\" "+
				"returned more than 1 result \"%s\"", pool.Name, strat.Name, created)
	}
	return created[0], nil
}

func (pool AllocatingPool) invokeAllocationStrategy(strat *ent.AllocationStrategy) (map[string]interface{}, error) {
	switch strat.Lang {
	case allocationStrategy.LangJs:
		return pool.invoker.invokeJs(strat.Script)
	case allocationStrategy.LangPy:
		return pool.invoker.invokePy(strat.Script)
	default:
		return nil, errors.Errorf("Unknown language \"%s\" for strategy \"%s\"", strat.Lang, strat.Name)
	}
}

// FreeResource deallocates the resource identified by its properties
func (pool AllocatingPool) FreeResource(raw RawResourceProps) error {
	return pool.freeResourceInner(raw, pool.retireResource, pool.freeResourceImmediately, pool.benchResource)
}

func (pool AllocatingPool) freeResourceImmediately(res *ent.Resource) error {
	// Delete props
	for _, prop := range res.Edges.Properties {
		if err := pool.client.Property.DeleteOne(prop).Exec(pool.ctx); err != nil {
			return errors.Wrapf(err, "Cannot free resource from \"%s\". Unable to cleanup properties", pool.Name)
		}
	}

	// Delete resource
	err := pool.client.Resource.DeleteOne(res).Exec(pool.ctx)
	if err != nil {
		return errors.Wrapf(err, "Cannot free resource from \"%s\". Unable to cleanup resource", pool.Name)
	}

	return nil
}
