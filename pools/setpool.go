package pools

import (
	"context"
	"time"

	"github.com/net-auto/resourceManager/ent"
	resource "github.com/net-auto/resourceManager/ent/resource"
	resourcePool "github.com/net-auto/resourceManager/ent/resourcepool"
	"github.com/net-auto/resourceManager/ent/schema"
	"github.com/pkg/errors"
)

// NewSetPool creates a brand new pool allocating DB entities in the process
func NewSetPool(
	ctx context.Context,
	client *ent.Client,
	resourceType *ent.ResourceType,
	propertyValues []RawResourceProps,
	poolName string,
	description *string,
	poolDealocationSafetyPeriod int) (Pool, error) {
	pool, _, err := NewSetPoolWithMeta(ctx, client, resourceType, propertyValues, poolName, description, poolDealocationSafetyPeriod)
	return pool, err
}

// NewSetPoolWithMeta creates a brand new pool + returns the pools underlying meta information
func NewSetPoolWithMeta(
	ctx context.Context,
	client *ent.Client,
	resourceType *ent.ResourceType,
	propertyValues []RawResourceProps,
	poolName string,
	description *string,
	poolDealocationSafetyPeriod int) (Pool, *ent.ResourcePool, error) {

	// TODO check that propertyValues are unique

	pool, err := newFixedPoolInner(ctx, client, resourceType, propertyValues,
		poolName, description, resourcePool.PoolTypeSet, poolDealocationSafetyPeriod)

	if err != nil {
		return nil, nil, err
	}

	return &SetPool{poolBase{pool, ctx, client}}, pool, nil
}

// Destroy removes the pool from DB if there are no more claims
func (pool SetPool) Destroy() error {
	// Check if there are no more claims
	claims, err := pool.QueryResources()
	if err != nil {
		return err
	}

	if len(claims) > 0 {
		return errors.Errorf("Unable to destroy pool \"%s\", there are claimed resources",
			pool.Name)
	}

	// Delete props
	resources, err := pool.findResources().All(pool.ctx)
	if err != nil {
		return errors.Wrapf(err, "Cannot destroy pool \"%s\". Unable to cleanup resoruces", pool.Name)
	}
	for _, res := range resources {
		props, err := res.QueryProperties().All(pool.ctx)
		if err != nil {
			return errors.Wrapf(err, "Cannot destroy pool \"%s\". Unable to cleanup resoruces", pool.Name)
		}

		for _, prop := range props {
			pool.client.Property.DeleteOne(prop).Exec(pool.ctx)
		}
		if err != nil {
			return errors.Wrapf(err, "Cannot destroy pool \"%s\". Unable to cleanup resoruces", pool.Name)
		}
	}

	// Delete resources
	_, err = pool.client.Resource.Delete().Where(resource.HasPoolWith(resourcePool.ID(pool.ID))).Exec(pool.ctx)
	if err != nil {
		return errors.Wrapf(err, "Cannot destroy pool \"%s\". Unable to cleanup resoruces", pool.Name)
	}

	// Delete pool itself
	err = pool.client.ResourcePool.DeleteOne(pool.ResourcePool).Exec(pool.ctx)
	if err != nil {
		return errors.Wrapf(err, "Cannot destroy pool \"%s\"", pool.Name)
	}

	return nil
}

// ClaimResource allocates the next available resource
func (pool SetPool) ClaimResource(userInput map[string]interface{}) (*ent.Resource, error) {
	// Allocate new resource for this tag
	unclaimedRes, err := pool.queryUnclaimedResourceEager()
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to find unclaimed resource in pool \"%s\"",
			pool.Name)
	}

	err = pool.client.Resource.UpdateOne(unclaimedRes).SetStatus(resource.StatusClaimed).Exec(pool.ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to claim a resource in pool \"%s\"", pool.Name)
	}
	return unclaimedRes, err
}

// FreeResource deallocates the resource identified by its properties
func (pool SetPool) FreeResource(raw RawResourceProps) error {
	return pool.freeResourceInner(raw, pool.retireResource, pool.freeResourceImmediately, pool.benchResource)
}

func (pool SetPool) freeResourceInner(raw RawResourceProps,
	retireResource func(res *ent.Resource) error,
	freeResource func(res *ent.Resource) error,
	benchResource func(res *ent.Resource) error,
) error {
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

	if res.Status != resource.StatusClaimed {
		return errors.Wrapf(err, "Unable to free a resource in pool \"%s\". It has not been claimed", pool.Name)
	}

	// Make sure there are no nested pools attached
	if nestedPool, err := res.QueryNestedPool().First(pool.ctx); err != nil && !ent.IsNotFound(err) {
		return errors.Wrapf(err, "Unable to free a resource in pool \"%s\". " +
			"Unable to check nested pools", pool.Name)
	} else if nestedPool != nil {
		return errors.Wrapf(err, "Unable to free a resource in pool \"%s\". " +
			"There is a nested pool attached to it \"%v\"", pool.Name, nestedPool.ID)
	}

	switch pool.ResourcePool.DealocationSafetyPeriod {
	case schema.ResourcePoolDealocationRetire:
		err = retireResource(res)
	case schema.ResourcePoolDealocationImmediately:
		err = freeResource(res)
	default:
		err = benchResource(res)
	}

	if err != nil {
		return errors.Wrapf(err, "Unable to free a resource in pool \"%s\". Unable to unclaim", pool.Name)
	}

	return nil
}

func (pool SetPool) benchResource(res *ent.Resource) error {
	return pool.client.Resource.UpdateOne(res).SetStatus(resource.StatusBench).Exec(pool.ctx)
}

func (pool SetPool) freeResourceImmediately(res *ent.Resource) error {
	return pool.client.Resource.UpdateOne(res).SetStatus(resource.StatusFree).Exec(pool.ctx)
}

func (pool SetPool) retireResource(res *ent.Resource) error {
	return pool.client.Resource.UpdateOne(res).SetStatus(resource.StatusRetired).Exec(pool.ctx)
}

func (pool SetPool) findResource(raw RawResourceProps) (*ent.ResourceQuery, error) {
	propComparator, err := CompareProps(pool.ctx, pool.QueryResourceType().OnlyX(pool.ctx), raw)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to find resource in pool: \"%s\"", pool.Name)
	}

	return pool.findResources().
		Where(resource.HasPropertiesWith(propComparator...)), nil
}

// QueryResource returns a resource identified by its properties
func (pool SetPool) QueryResource(raw RawResourceProps) (*ent.Resource, error) {
	query, err := pool.findResource(raw)
	if err != nil {
		return nil, err
	}
	return query.
		Where(resource.StatusEQ(resource.StatusClaimed)).
		Only(pool.ctx)
}

// load eagerly with some edges, ready to be copied
func (pool SetPool) queryUnclaimedResourceEager() (*ent.Resource, error) {
	// Find first unclaimed
	res, err := pool.findResources().
		Where(resource.StatusEQ(resource.StatusFree)).
		First(pool.ctx)

	// No more free, try benched that have been benched before NOW - dealocationSafetyPeriod
	if ent.IsNotFound(err) {
		res, err = pool.findResources().
			Where(resource.StatusEQ(resource.StatusBench)).
			Where(resource.UpdatedAtLT(time.Now().Add(time.Duration(-pool.ResourcePool.DealocationSafetyPeriod) * time.Second))).
			First(pool.ctx)
	}

	// No more benched, its over
	if ent.IsNotFound(err) {
		return nil, errors.Wrapf(err, "No more free resources in the pool: \"%s\"", pool.Name)
	}

	return res, err
}

func (pool SetPool) findResources() *ent.ResourceQuery {
	return pool.client.Resource.Query().
		Where(resource.HasPoolWith(resourcePool.ID(pool.ID)))
}

// QueryResources returns all allocated resources
func (pool SetPool) QueryResources() (ent.Resources, error) {
	res, err := pool.findResources().
		Where(resource.StatusEQ(resource.StatusClaimed)).
		All(pool.ctx)

	return res, err
}
