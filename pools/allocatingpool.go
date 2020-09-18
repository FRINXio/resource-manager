package pools

import (
	"context"
	"github.com/net-auto/resourceManager/ent/poolproperties"
	"github.com/net-auto/resourceManager/ent/predicate"
	"github.com/net-auto/resourceManager/ent/property"
	"time"

	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/resource"
	resourcePool "github.com/net-auto/resourceManager/ent/resourcepool"
	"github.com/net-auto/resourceManager/graph/graphql/model"
	"github.com/pkg/errors"
)

func DeletePoolProperties(ctx context.Context, client *ent.Client, poolId int) error {
	poolProperties, err1 := client.PoolProperties.Query().Where(poolproperties.HasPoolWith(resourcePool.ID(poolId))).WithProperties().Only(ctx)

	if err1 != nil && !ent.IsNotFound(err1) {
		return err1
	}

	//pool properties does not exist for this pool (it is nested)
	if ent.IsNotFound(err1) {
		return nil
	}

	rp, err := client.ResourcePool.Query().Where(resourcePool.ID(poolId)).WithParentResource().Only(ctx)
	if err != nil {
		return err
	}

	//if we are deleting root pool we need to delete the individual properties as well
	if rp.Edges.ParentResource == nil {
		propertyIdsToDelete := make([]predicate.Property, len(poolProperties.Edges.Properties))
		for i, prop := range poolProperties.Edges.Properties {
			propertyIdsToDelete[i] = property.ID(prop.ID)
		}

		_, err := client.Property.Delete().Where(property.Or(propertyIdsToDelete...)).Exec(ctx)
		if err != nil {
			return err
		}
	}

	if err := client.PoolProperties.DeleteOne(poolProperties).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func CreatePoolProperties(ctx context.Context, client *ent.Client, pp []map[string]interface{},resPropertyType *ent.ResourceType ) (*ent.PoolProperties, error) {
	var propTypes = ToRawTypes(pp)

	//this loops only once
	for _, propType := range propTypes {
		properties, err := ParseProps(ctx, client, resPropertyType, propType)
		if err != nil {
			return nil, err
		}
		return client.PoolProperties.Create().AddProperties(properties...).AddResourceType(resPropertyType).Save(ctx)
	}

	return nil, errors.New("Unable to create pool properties")
}

// NewAllocatingPool creates a brand new pool allocating DB entities in the process
func NewAllocatingPool(
	ctx context.Context,
	client *ent.Client,
	resourceType *ent.ResourceType,
	allocationStrategy *ent.AllocationStrategy,
	poolName string,
	poolDealocationSafetyPeriod int) (Pool, error) {
	pool, _, err := NewAllocatingPoolWithMeta(ctx, client, resourceType, allocationStrategy, poolName, nil, poolDealocationSafetyPeriod, nil)
	return pool, err
}

// NewAllocatingPoolWithMeta creates a brand new pool + returns the pools underlying meta information
func NewAllocatingPoolWithMeta(ctx context.Context,
	client *ent.Client,
	resourceType *ent.ResourceType,
	allocationStrategy *ent.AllocationStrategy,
	poolName string,
	description *string,
	poolDealocationSafetyPeriod int,
	poolProperties *ent.PoolProperties) (Pool, *ent.ResourcePool, error) {

	// TODO keep just single instance
	wasmer, err := NewWasmerUsingEnvVars()
	if err != nil {
		return nil, nil, errors.Wrap(err, "Cannot create resource pool")
	}
	return newAllocatingPoolWithMetaInternal(
		ctx, client, resourceType, allocationStrategy, poolName, description, wasmer, poolDealocationSafetyPeriod, poolProperties)
}

func newAllocatingPoolWithMetaInternal(
	ctx context.Context,
	client *ent.Client,
	resourceType *ent.ResourceType,
	allocationStrategy *ent.AllocationStrategy,
	poolName string,
	description *string,
	invoker ScriptInvoker,
	poolDealocationSafetyPeriod int,
	poolProperties *ent.PoolProperties) (Pool, *ent.ResourcePool, error) {

	pool, err := client.ResourcePool.Create().
		SetName(poolName).
		SetNillableDescription(description).
		SetPoolType(resourcePool.PoolTypeAllocating).
		SetResourceType(resourceType).
		SetAllocationStrategy(allocationStrategy).
		SetDealocationSafetyPeriod(poolDealocationSafetyPeriod).
		Save(ctx)

	if err != nil {
		return nil, nil, errors.Wrap(err, "Cannot create resource pool")
	}

	if poolProperties != nil {
		_, err = pool.Update().SetPoolProperties(poolProperties).Save(ctx)
	}

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

func (pool AllocatingPool) PoolProperties() ([]*ent.Property, error) {
	return pool.QueryPoolProperties().QueryProperties().WithType().All(pool.ctx)
}

// TODO add capacity implementation
func (pool AllocatingPool) Capacity() (int, error) {
	return 1, nil
}

// ClaimResource allocates the next available resource
func (pool AllocatingPool) ClaimResource(userInput map[string]interface{}) (*ent.Resource, error) {

	strat, err := pool.AllocationStrategy()
	if err != nil {
		return nil, errors.Wrapf(err,
			"Unable to claim resource from pool #%d, allocation strategy loading error ", pool.ID)
	}

	ps, err := pool.PoolProperties()

	if err != nil {
		return nil, errors.Wrapf(err,
			"Unable to claim resource from pool #%d, resource type loading error ", pool.ID)
	}

	propMap, propErr := convertProperties(ps)

	if propErr != nil {
		return nil, errors.Wrapf(propErr, "Unable to convert value from property")
	}

	resourceType, err := pool.ResourceType()
	if err != nil {
		return nil, errors.Wrapf(err,
			"Unable to claim resource from pool #%d, resource type loading error ", pool.ID)
	}

	var resourcePool model.ResourcePoolInput
	resourcePool.ResourcePoolName = pool.Name

	currentResources, err := pool.loadClaimedResources()
	if err != nil {
		return nil, errors.Wrapf(err,
			"Unable to claim resource from pool #%d, resource loading error ", pool.ID)
	}

	resourceProperties, _ /*TODO do something with logs */, err := InvokeAllocationStrategy(
		pool.invoker, strat, userInput, resourcePool, currentResources, propMap)
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

func convertProperties(ps []*ent.Property) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for _, p := range ps {
		value, err := GetValue(p)
		if err != nil {
			return nil, err
		}

		result[p.Edges.Type.Name] = value
	}

	return result, nil
}

func  (pool AllocatingPool) loadClaimedResources() ([]*model.ResourceInput, error) {
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
		if propsToMap, err := PropertiesToMap(claimedResource.Edges.Properties); err != nil {
			return nil, errors.Wrapf(err, "Unable to serialize resource properties")
		} else {
			r.Properties = propsToMap
		}

		currentResources = append(currentResources, &r)
	}
	return currentResources, nil
}

func clearBenchedResources(pool AllocatingPool) error {
	_, err := pool.client.Resource.Delete().
		Where(resource.StatusEQ(resource.StatusBench)).
		Where(resource.UpdatedAtLT(time.Now().Add(time.Duration(-pool.ResourcePool.DealocationSafetyPeriod) * time.Second))).
		Exec(pool.ctx)

	return err
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
