package pools

import (
	"context"
	"time"

	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/resource"
	resourcePool "github.com/net-auto/resourceManager/ent/resourcepool"
	"github.com/net-auto/resourceManager/graph/graphql/model"
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

	// TODO keep just single instance
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
		return errors.Errorf("Unable to destroy pool #%d, there are claimed resources",
			pool.ID)
	}

	// Delete pool itself
	err = pool.client.ResourcePool.DeleteOne(pool.ResourcePool).Exec(pool.ctx)
	if err != nil {
		return errors.Wrapf(err, "Cannot destroy pool #%d", pool.ID)
	}

	return nil
}

func (pool AllocatingPool) AllocationStrategy() (*ent.AllocationStrategy, error) {
	return pool.ResourcePool.QueryAllocationStrategy().Only(pool.ctx)
}

// ClaimResource allocates the next available resource
func (pool AllocatingPool) ClaimResource(userInput map[string]interface{}) (*ent.Resource, error) {

	strat, err := pool.AllocationStrategy()
	if err != nil {
		return nil, errors.Wrapf(err,
			"Unable to claim resource from pool #%d, allocation strategy loading error ", pool.ID)
	}
	resourceType, err := pool.ResourceType()
	if err != nil {
		return nil, errors.Wrapf(err,
			"Unable to claim resource from pool #%d, resource type loading error ", pool.ID)
	}

	var resourcePool model.ResourcePoolInput
	resourcePool.ResourcePoolName = pool.Name

	var currentResources []*model.ResourceInput
	claimedResources, err := pool.findResources().WithProperties(
		func(propertyQuery *ent.PropertyQuery) { propertyQuery.WithType() }).All(pool.ctx)
	if err != nil {
		return nil, errors.Wrapf(err,
			"Unable get claimed resources from pool #%d, resource loading error ", pool.ID)
	}
	for _, claimedResource := range claimedResources {
		var r model.ResourceInput
		r.UpdatedAt = claimedResource.UpdatedAt.String()
		r.Status = claimedResource.Status.String()
		for _, prop := range claimedResource.Edges.Properties {
			name := prop.Edges.Type.Name
			var pi model.PropertyInput
			pi.Name = name
			pi.Type = prop.Edges.Type.Type.String()
			pi.Mandatory = prop.Edges.Type.Mandatory
			pi.IntVal = prop.IntVal
			pi.FloatVal = prop.FloatVal
			pi.StringVal = prop.StringVal
			r.Properties = append(r.Properties, &pi)
		}
		currentResources = append(currentResources, &r)
	}

	resourceProperties, _ /*TODO do something with logs */, err := InvokeAllocationStrategy(
		pool.invoker, strat, userInput, resourcePool, currentResources)
	if err != nil {
		return nil, errors.Wrapf(err,
			"Unable to claim resource from pool #%d, allocation strategy \"%s\" failed", pool.ID, strat.Name)
	}

	// Query to check whether this resource already exists.
	// 1. construct query
	query, err := pool.findResource(RawResourceProps(resourceProperties))
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot query for resource based on pool #%d and properties \"%s\"", pool.ID, resourceProperties)
	}

	// 2. Try to find the resource in DB
	foundResources, err := query.WithProperties().All(pool.ctx)

	if len(foundResources) == 0 {
		// 3a. Nothing found - create new resource
		created, err := PreCreateResources(pool.ctx, pool.client, []RawResourceProps{resourceProperties},
			pool.ResourcePool, resourceType, resource.StatusClaimed)
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to create resource in pool #%d", pool.ID)
		}
		if len(created) > 1 {
			return nil, errors.Errorf(
				"Unexpected error creating resource in pool #%d, properties \"%s\" . "+
					"Created %d resources instead of one.", pool.ID, resourceProperties, len(created))
		}
		return created[0], nil
	} else if len(foundResources) > 1 {
		return nil, errors.Errorf(
			"Unable to claim resource with properties \"%s\" from pool #%d, database contains more than one result", resourceProperties, pool.ID)
	}
	res := foundResources[0]
	// 3b. Claim found resource if possible
	if res.Status == resource.StatusClaimed || res.Status == resource.StatusRetired {
		return nil, errors.Errorf("Resource #%d is in incorrect state \"%s\"", res.ID, res.Status)
	} else if res.Status == resource.StatusBench {
		cutoff := res.UpdatedAt.Add(time.Duration(pool.DealocationSafetyPeriod) * time.Second)
		if time.Now().Before(cutoff) {
			return nil, errors.Errorf(
				"Unable to claim resource #%d from pool #%d, resource cannot be claimed before %s", res.ID, pool.ID, cutoff)
		}
	}
	res.Status = resource.StatusClaimed
	err = pool.client.Resource.UpdateOne(res).SetStatus(res.Status).Exec(pool.ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot update resource #%d", res.ID)
	}
	return res, nil
}

// FreeResource deallocates the resource identified by its properties
func (pool AllocatingPool) FreeResource(raw RawResourceProps) error {
	return pool.freeResourceInner(raw, pool.retireResource, pool.freeResourceImmediately, pool.benchResource)
}

func (pool AllocatingPool) freeResourceImmediately(res *ent.Resource) error {
	// Delete props
	for _, prop := range res.Edges.Properties {
		if err := pool.client.Property.DeleteOne(prop).Exec(pool.ctx); err != nil {
			return errors.Wrapf(err, "Cannot free resource from #%d. Unable to cleanup properties", pool.ID)
		}
	}

	// Delete resource
	err := pool.client.Resource.DeleteOne(res).Exec(pool.ctx)
	if err != nil {
		return errors.Wrapf(err, "Cannot free resource from #%d. Unable to cleanup resource", pool.ID)
	}

	return nil
}
