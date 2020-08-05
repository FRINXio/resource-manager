package pools

import (
	"context"
	"github.com/net-auto/resourceManager/ent"
	allocationStrategy "github.com/net-auto/resourceManager/ent/allocationstrategy"
	resource "github.com/net-auto/resourceManager/ent/resource"
	resourcePool "github.com/net-auto/resourceManager/ent/resourcepool"
	"github.com/pkg/errors"
)

// NewAllocatingPool creates a brand new pool allocating DB entities in the process
func NewAllocatingPool(
	ctx context.Context,
	client *ent.Client,
	resourceType *ent.ResourceType,
	allocationStrategy *ent.AllocationStrategy,
	poolName string) (Pool, error) {
	pool, _, err := NewAllocatingPoolWithMeta(ctx, client, resourceType, allocationStrategy, poolName)
	return pool, err
}

// NewAllocatingPoolWithMeta creates a brand new pool + returns the pools underlying meta information
func NewAllocatingPoolWithMeta(
	ctx context.Context,
	client *ent.Client,
	resourceType *ent.ResourceType,
	allocationStrategy *ent.AllocationStrategy,
	poolName string) (Pool, *ent.ResourcePool, error) {

	return newAllocatingPoolWithMetaInternal(
		ctx, client, resourceType, allocationStrategy, poolName, NewWasmerDefault())
}

func newAllocatingPoolWithMetaInternal(
	ctx context.Context,
	client *ent.Client,
	resourceType *ent.ResourceType,
	allocationStrategy *ent.AllocationStrategy,
	poolName string,
	invoker ScriptInvoker) (Pool, *ent.ResourcePool, error) {

	pool, err := client.ResourcePool.Create().
		SetName(poolName).
		SetPoolType(resourcePool.PoolTypeAllocating).
		SetResourceType(resourceType).
		SetAllocationStrategy(allocationStrategy).
		Save(ctx)

	if err != nil {
		return nil, nil, err
	}

	return &AllocatingPool{
			poolBase{pool, ctx, client},
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

func (pool AllocatingPool) AddLabel(label PoolLabel) error {
	// TODO implement labeling
	return errors.Errorf("NOT IMPLEMENTED")
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
		pool.ResourcePool, resourceType, true)
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
	// TODO Python
	default:
		return nil, errors.Errorf("Unknown language \"%s\" for strategy \"%s\"", strat.Lang, strat.Name)
	}
}

// TODO extract common code between this and set pool

// FreeResource deallocates the resource identified by its properties
func (pool AllocatingPool) FreeResource(raw RawResourceProps) error {
	return pool.freeResourceInner(raw)
}

func (pool AllocatingPool) freeResourceInner(raw RawResourceProps) error {
	query, err := pool.findResource(raw)
	if err != nil {
		return errors.Wrapf(err, "Unable to find resource in pool: \"%s\"", pool.Name)
	}
	res, err := query.
		WithProperties().
		Only(pool.ctx)

	if err != nil {
		return errors.Wrapf(err, "Unable to free a resource in pool \"%s\". Unable to find resource", pool.Name)
	}

	// Delete props
	for _, prop := range res.Edges.Properties {
		if err := pool.client.Property.DeleteOne(prop).Exec(pool.ctx); err != nil {
			return errors.Wrapf(err, "Cannot free resource from \"%s\". Unable to cleanup properties", pool.Name)
		}
	}

	// Delete resource
	err = pool.client.Resource.DeleteOne(res).Exec(pool.ctx)
	if err != nil {
		return errors.Wrapf(err, "Cannot free resource from \"%s\". Unable to cleanup resource", pool.Name)
	}

	return nil
}

func (pool AllocatingPool) findResource(raw RawResourceProps) (*ent.ResourceQuery, error) {
	propComparator, err := compareProps(pool.ctx, pool.QueryResourceType().OnlyX(pool.ctx), raw)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to find resource in pool: \"%s\"", pool.Name)
	}

	return pool.findResources().
		Where(resource.HasPropertiesWith(propComparator...)), nil
}

// QueryResource returns a resource identified by its properties
func (pool AllocatingPool) QueryResource(raw RawResourceProps) (*ent.Resource, error) {
	query, err := pool.findResource(raw)
	if err != nil {
		return nil, err
	}
	return query.
		Where(resource.Claimed(true)).
		Only(pool.ctx)
}

func (pool AllocatingPool) findResources() *ent.ResourceQuery {
	return pool.client.Resource.Query().
		Where(resource.HasPoolWith(resourcePool.ID(pool.ID)))
}

// QueryResources returns all allocated resources
func (pool AllocatingPool) QueryResources() (ent.Resources, error) {
	res, err := pool.findResources().
		Where(resource.Claimed(true)).
		All(pool.ctx)

	return res, err
}
