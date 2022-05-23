package pools

import (
	"context"
	"github.com/facebook/ent/dialect/sql"
	"github.com/facebook/ent/dialect/sql/sqljson"
	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/predicate"
	resource "github.com/net-auto/resourceManager/ent/resource"
	resourcePool "github.com/net-auto/resourceManager/ent/resourcepool"
	"github.com/net-auto/resourceManager/ent/schema"
	log "github.com/net-auto/resourceManager/logging"
	"github.com/pkg/errors"
	"strconv"
	"time"
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
		log.Error(ctx, err, "Unable to create set-pool")
		return nil, nil, err
	}

	return &SetPool{poolBase{pool, ctx, client}}, pool, nil
}

// Destroy removes the pool from DB if there are no more claims
func (pool SetPool) Destroy() error {
	// Check if there are no more claims
	claims, err := pool.QueryResources()
	if err != nil {
		log.Error(pool.ctx, err, "Unable to retrieve allocated resources")
		return err
	}

	if len(claims) > 0 {
		log.Warn(pool.ctx, "Unable to destroy pool \"%s\", there are claimed resources", pool.Name)
		return errors.Errorf("Unable to destroy pool \"%s\", there are claimed resources",
			pool.Name)
	}

	// Delete props
	resources, err := pool.findResources().All(pool.ctx)
	if err != nil {
		log.Error(pool.ctx, err, "Cannot find resources for pool \"%s\"", pool.Name)
		return errors.Wrapf(err, "Cannot destroy pool \"%s\". Unable to cleanup resoruces", pool.Name)
	}
	for _, res := range resources {
		props, err := res.QueryProperties().All(pool.ctx)
		if err != nil {
			log.Error(pool.ctx, err, "Cannot find properties for pool \"%s\"", pool.Name)
			return errors.Wrapf(err, "Cannot destroy pool \"%s\". Unable to cleanup resources", pool.Name)
		}

		for _, prop := range props {
			pool.client.Property.DeleteOne(prop).Exec(pool.ctx) //TODO missing error handling
		}
	}

	// Delete resources
	_, err = pool.client.Resource.Delete().Where(resource.HasPoolWith(resourcePool.ID(pool.ID))).Exec(pool.ctx)
	if err != nil {
		log.Error(pool.ctx, err, "Cannot resources of pool with ID %d", pool.ID)
		return errors.Wrapf(err, "Cannot destroy pool \"%s\". Unable to cleanup resoruces", pool.Name)
	}

	// Delete pool itself
	err = pool.client.ResourcePool.DeleteOne(pool.ResourcePool).Exec(pool.ctx)
	if err != nil {
		log.Error(pool.ctx, err, "Cannot delete pool with ID %d", pool.ID)
		return errors.Wrapf(err, "Cannot destroy pool \"%s\"", pool.Name)
	}

	return nil
}

// ClaimResource allocates the next available resource
func (pool SetPool) ClaimResource(userInput map[string]interface{}, description *string, alternativeId map[string]interface{}) (*ent.Resource, error) {

	// Allocate new resource for this tag
	unclaimedRes, err := pool.queryUnclaimedResourceEager()
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to find unclaimed resource in pool \"%s\"",
			pool.Name)
	}

	err = pool.client.Resource.
		UpdateOne(unclaimedRes).
		SetStatus(resource.StatusClaimed).
		SetNillableDescription(description).
		SetAlternateID(alternativeId).
		Exec(pool.ctx)

	if err != nil {
		err := errors.Wrapf(err, "Unable to claim a resource in pool \"%s\"", pool.Name)
		log.Error(pool.ctx, err, "Unable to claim a resource")
		return nil, err
	}
	return unclaimedRes, err
}

func (pool SetPool) Capacity() (string, string, error) {
	claimedResources, err := pool.QueryResources()

	if err != nil {
		log.Error(pool.ctx, err, "Unable to retrieve resources for pool ID %d", pool.ID)
		return "0", "0", err
	}

	resources, err := pool.client.Resource.Query().Where(resource.HasPoolWith(resourcePool.ID(pool.ID))).All(pool.ctx)

	if err != nil {
		log.Error(pool.ctx, err, "Unable to retrieve resources for pool ID %d", pool.ID)
		return "0", "0", err
	}

	return strconv.Itoa(len(resources) - len(claimedResources)), strconv.Itoa(len(claimedResources)), nil
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
		err := errors.Wrapf(err, "Unable to find resource in pool: \"%s\"", pool.Name)
		log.Error(pool.ctx, err, "Unable to find resource")
		return err
	}
	res, err := query.
		WithProperties().
		Only(pool.ctx)

	if err != nil {
		err := errors.Wrapf(err, "Unable to free a resource in pool \"%s\". Unable to find resource", pool.Name)
		log.Error(pool.ctx, err, "Unable to free a resource in pool")
		return err
	}

	if res.Status != resource.StatusClaimed {
		log.Warn(pool.ctx, "Unable to free a resource in pool \"%s\". It has not been claimed", pool.Name)
		return errors.Wrapf(err, "Unable to free a resource in pool \"%s\". It has not been claimed", pool.Name)
	}

	// Make sure there are no nested pools attached
	if nestedPool, err := res.QueryNestedPool().First(pool.ctx); err != nil && !ent.IsNotFound(err) {
		log.Error(pool.ctx, err, "Unable to free a resource in pool ID %d", pool.ID)
		return errors.Wrapf(err, "Unable to free a resource in pool \"%s\". "+
			"Unable to check nested pools", pool.Name)
	} else if nestedPool != nil {
		err := errors.Errorf("Unable to free a resource in pool \"%s\". "+
			"There is a nested pool attached to it \"%v\"", pool.Name, nestedPool.ID)
		log.Error(pool.ctx, err, "Unable to free a resource in pool ID %d", pool.ID)
		return err
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
		err := errors.Wrapf(err, "Unable to free a resource in pool \"%s\". Unable to unclaim", pool.Name)
		log.Error(pool.ctx, err, "Unable to free a resource in pool")
		return err
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
		err := errors.Wrapf(err, "Unable to find resource in pool: \"%s\"", pool.Name)
		log.Error(pool.ctx, err, "Unable to find resource in pool")
		return nil, err
	}

	resources := pool.findResources()
	var resourceComparator []predicate.Resource
	for _, propPred := range propComparator {
		resourceComparator = append(resourceComparator, resource.HasPropertiesWith(propPred))

	}
	return resources.Where(resourceComparator...), nil
}

// QueryResource returns a resource identified by its properties
func (pool SetPool) QueryResource(raw RawResourceProps) (*ent.Resource, error) {
	query, err := pool.findResource(raw)
	if err != nil {
		log.Error(pool.ctx, err, "Unable to find resource %+v in pool %d", raw, pool.ID)
		return nil, err
	}
	only, err2 := query.
		Where(resource.StatusEQ(resource.StatusClaimed)).
		Only(pool.ctx)

	if err2 != nil {
		log.Error(pool.ctx, err, "Unable retrieve resources for pool ID %d", pool.ID)
	}

	return only, err2
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
		err := errors.Wrapf(err, "No more free resources in the pool: \"%s\"", pool.Name)
		log.Error(pool.ctx, err, "No free resources in pool ID %d", pool.ID)
		return nil, err
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

	if err != nil {
		log.Error(pool.ctx, err, "Unable retrieve resources for pool ID %d", pool.ID)
	}

	return res, err
}

// QueryResourcesByAltId returns resources if alt Id matches
func (pool SetPool) QueryResourcesByAltId(alternativeId map[string]interface{}) ([]*ent.Resource, error) {
	res, err := pool.client.Resource.Query().
		Where(func(selector *sql.Selector) {
			for k, v := range alternativeId {
				selector.Where(sqljson.ValueEQ("alternate_id", v, sqljson.Path(k)))
			}
		}).All(pool.ctx)

	if err != nil {
		log.Error(pool.ctx, err, "Unable to retrieve resources in pool ID %d", pool.ID)
		return nil, err
	}

	if res != nil {
		return res, nil
	}

	log.Warn(pool.ctx, "There is not such resource with alternative ID %v for pool ID %d", alternativeId, pool.ID)

	return nil, errors.New("No such resource with given alternative ID")
}
