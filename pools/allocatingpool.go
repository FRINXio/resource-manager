package pools

import (
	"context"
	"github.com/net-auto/resourceManager/ent/allocationstrategy"
	"github.com/net-auto/resourceManager/ent/poolproperties"
	"github.com/net-auto/resourceManager/ent/predicate"
	"github.com/net-auto/resourceManager/ent/property"
	"github.com/net-auto/resourceManager/graph/graphql/model"
	log "github.com/net-auto/resourceManager/logging"
	"time"

	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/resource"
	resourcePool "github.com/net-auto/resourceManager/ent/resourcepool"
	"github.com/pkg/errors"
)

var manualSqlExecutionStrategies = map[string]bool{
	"unique_id": true,
}

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
		log.Error(ctx, err, "Unable to retrieve resource pool %d", poolId)
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
			log.Error(ctx, err, "Unable to delete properties")
			return err
		}
	}

	if err := client.PoolProperties.DeleteOne(poolProperties).Exec(ctx); err != nil {
		log.Error(ctx, err, "Unable to delete pool-properties with ID %d", poolProperties.ID)
		return err
	}

	return nil
}

func CreatePoolProperties(ctx context.Context, client *ent.Client, pp []map[string]interface{}, resPropertyType *ent.ResourceType) (*ent.PoolProperties, error) {
	var propTypes = ToRawTypes(pp)

	//this loops only once
	for _, propType := range propTypes {
		properties, err := ParseProps(ctx, client, resPropertyType, propType)
		if err != nil {
			log.Error(ctx, err, "Unable to parse property %+v", propType)
			return nil, err
		}
		return client.PoolProperties.Create().AddProperties(properties...).AddResourceType(resPropertyType).Save(ctx)
	}

	return nil, errors.New("Unable to create pool properties")
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
		log.Error(ctx, err, "Creating wasmer instance failed")
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

	query, _ := client.ResourcePool.Query().Where(resourcePool.NameEQ(poolName)).Only(ctx)
	if query != nil {
		log.Error(ctx, nil, "Unable to create a resource pool: resource pool with name "+query.Name+" already exists.")
		return nil, nil, errors.New("Unable to create a resource pool: resource pool with name " + query.Name + " already exists. " +
			"Resource pool name must be unique, use another name.")
	}

	pool, err := client.ResourcePool.Create().
		SetName(poolName).
		SetNillableDescription(description).
		SetPoolType(resourcePool.PoolTypeAllocating).
		SetResourceType(resourceType).
		SetAllocationStrategy(allocationStrategy).
		SetDealocationSafetyPeriod(poolDealocationSafetyPeriod).
		Save(ctx)

	if err != nil {
		log.Error(ctx, err, "Unable to create a resource pool")
		return nil, nil, errors.Wrap(err, "Cannot create resource pool")
	}

	if poolProperties != nil {
		_, err = pool.Update().SetPoolProperties(poolProperties).Save(ctx)
	}

	if err != nil {
		log.Error(ctx, err, "Unable to create resource properties")
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
		log.Warn(pool.ctx, "Unable to delete pool with ID %d because it has allocated resources", pool.ID)
		return errors.Errorf("Unable to destroy pool #%d, there are claimed resources",
			pool.ID)
	}

	err = DeletePoolProperties(pool.ctx, pool.client, pool.ID)

	if err != nil {
		return err
	}

	// Delete pool itself
	err = pool.client.ResourcePool.DeleteOne(pool.ResourcePool).Exec(pool.ctx)
	if err != nil {
		log.Error(pool.ctx, err, "Unable delete pool with ID %d", pool.ID)
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

func (pool AllocatingPool) Capacity() (string, string, error) {

	strat, err := pool.AllocationStrategy()
	if err != nil {
		log.Error(pool.ctx, err, "Unable to retrieve allocation-strategy for pool %d", pool.ID)
		return "0", "0", errors.Wrapf(err,
			"Unable to retrieve allocation-strategy for pool %d, allocation strategy loading error", pool.ID)
	}
	var currentResources []*model.ResourceInput

	if !manualSqlExecutionStrategies[strat.Name] {
		currentResources, err = getFullListOfResources(pool)
		if err != nil {
			log.Error(pool.ctx, err, "Unable to load resources for pool %d", pool.ID)
			return "0", "0", errors.Wrapf(err,
				"Unable to load resources for pool %d, resource loading error", pool.ID)
		}
	} else {
		// The query call doesn't create a transaction, so we have to create it manually
		tx, err := pool.client.Tx(pool.ctx)
		if err != nil {
			log.Error(pool.ctx, err, "Unable to open new read transaction for pool %d", pool.ID)
			return "0", "0", errors.Wrapf(err, "Unable to open new read transaction for pool %d", pool.ID)
		}
		pool.ctx = context.WithValue(pool.ctx, ent.TxCtxKey{}, tx)
		defer func(tx *ent.Tx) {
			err := tx.Commit()
			if err != nil {
				log.Error(pool.ctx, err, "Unable to commit read transaction for pool %d", pool.ID)
			}
		}(tx)
	}

	ps, err := pool.PoolProperties()

	if err != nil {
		log.Error(pool.ctx, err, "Unable to load resources for pool %d", pool.ID)
		return "0", "0", errors.Wrapf(err,
			"Unable to get properties from pool #%d, resource type loading error ", pool.ID)
	}

	propMap, propErr := convertProperties(ps)

	if propErr != nil {
		log.Error(pool.ctx, propErr, "Unable to convert value from property")
		return "0", "0", errors.Wrapf(propErr, "Unable to convert value from property")
	}

	var emptyMap = map[string]interface{}{}
	result, _, err := InvokeAllocationStrategy(pool.ctx, pool.invoker, strat, emptyMap, model.ResourcePoolInput{
		ResourcePoolID:   pool.ID,
		PoolProperties:   emptyMap,
		ResourcePoolName: pool.Name,
	}, currentResources, propMap, "capacity()")
	if err != nil || result == nil {
		log.Error(pool.ctx, err, "Invoking allocation strategy failed")
		return "0", "0", errors.Wrapf(err,
			"Unable to compute capacity pool #%d, allocation strategy \"%s\" failed", pool.ID, strat.Name)
	}

	var resultFreeCapacity = "0"
	var resultUtilizedCapacity = "0"

	if result["freeCapacity"] != nil {
		resultFreeCapacity = result["freeCapacity"].(string)
	}

	if result["utilizedCapacity"] != nil {
		resultUtilizedCapacity = result["utilizedCapacity"].(string)
	}

	return resultFreeCapacity, resultUtilizedCapacity, nil
}

// ClaimResource allocates the next available resource
func (pool AllocatingPool) ClaimResource(userInput map[string]interface{}, description *string, alternativeId map[string]interface{}) (*ent.Resource, error) {

	strat, err := pool.AllocationStrategy()
	if err != nil {
		log.Error(pool.ctx, err, "Unable to retrieve allocation-strategy for pool %d", pool.ID)
		return nil, errors.Wrapf(err,
			"Unable to claim resource from pool #%d, allocation strategy loading error ", pool.ID)
	}

	ps, err := pool.PoolProperties()

	if err != nil {
		log.Error(pool.ctx, err, "Unable to retrieve pool-properties for pool %d", pool.ID)
		return nil, errors.Wrapf(err,
			"Unable to claim resource from pool #%d, resource type loading error ", pool.ID)
	}

	propMap, propErr := convertProperties(ps)

	if propErr != nil {
		log.Error(pool.ctx, propErr, "Unable to convert value from property")
		return nil, errors.Wrapf(propErr, "Unable to convert value from property")
	}

	resourceType, err := pool.ResourceType()
	if err != nil {
		log.Error(pool.ctx, err, "Unable retrieve resource type for pool with ID: %d", pool.ID)
		return nil, errors.Wrapf(err,
			"Unable to claim resource from pool #%d, resource type loading error ", pool.ID)
	}

	var resourcePool model.ResourcePoolInput
	resourcePool.ResourcePoolName = pool.Name
	resourcePool.ResourcePoolID = pool.ID
	var currentResources []*model.ResourceInput

	if !manualSqlExecutionStrategies[strat.Name] {
		currentResources, err = getFullListOfResources(pool)
		if err != nil {
			log.Error(pool.ctx, err, "Unable retrieve already claimed resources for pool with ID: %d", pool.ID)
			return nil, errors.Wrapf(err,
				"Unable to claim resource from pool #%d, resource loading error ", pool.ID)
		}
	}
	var functionName string

	if strat.Lang == allocationstrategy.LangPy {
		functionName = "script_fun()"
	} else {
		functionName = "invoke()"
	}
	resourceProperties, _ /*TODO do something with logs */, err := InvokeAllocationStrategy(
		pool.ctx, pool.invoker, strat, userInput, resourcePool, currentResources, propMap, functionName)
	if err != nil {
		log.Error(pool.ctx, err, "Unable to claim resource with pool with ID: %d, invoking strategy failed", pool.ID)
		return nil, errors.Wrapf(err,
			"Unable to claim resource from pool #%d, allocation strategy \"%s\" failed", pool.ID, strat.Name)
	}

	// Query to check whether this resource already exists.
	// 1. construct query
	query, err := pool.findResource(RawResourceProps(resourceProperties))
	if err != nil {
		log.Error(pool.ctx, err, "Cannot query for resource based on pool with ID %d", pool.ID)
		return nil, errors.Wrapf(err, "Cannot query for resource based on pool #%d and properties \"%s\"", pool.ID, resourceProperties)
	}

	// 2. Try to find the resource in DB
	foundResources, err := query.WithProperties().All(pool.ctx)

	//TODO - what if foundResources is nil ?? do we continue?
	if err != nil {
		log.Error(pool.ctx, err, "Unable to retrieve allocated resources for pool %d", pool.ID)
	}

	if len(foundResources) == 0 {
		// 3a. Nothing found - create new resource
		created, err := PreCreateResources(pool.ctx, pool.client, []RawResourceProps{resourceProperties},
			pool.ResourcePool, resourceType, resource.StatusClaimed, description, alternativeId)
		if err != nil {
			log.Error(pool.ctx, err, "Unable to create resource in pool %d", pool.ID)
			return nil, errors.Wrapf(err, "Unable to create resource in pool #%d", pool.ID)
		}
		if len(created) > 1 {
			// TODO this seems serious, shouldn't we delete those resources or something more than log it?
			log.Error(pool.ctx, err, "Unexpected error creating resource in pool %d"+
				" multiple resources created (count: %d)", pool.ID, len(created))
			return nil, errors.Errorf(
				"Unexpected error creating resource in pool #%d, properties \"%s\" . "+
					"Created %d resources instead of one.", pool.ID, resourceProperties, len(created))
		}
		return created[0], nil
	} else if len(foundResources) > 1 {
		log.Error(pool.ctx, err, "Unable to claim resource for pool ID %d, database contains more than one result", pool.ID)
		return nil, errors.Errorf(
			"Unable to claim resource with properties \"%s\" from pool #%d, database contains more than one result", resourceProperties, pool.ID)
	}
	res := foundResources[0]
	// 3b. Claim found resource if possible
	if res.Status == resource.StatusClaimed || res.Status == resource.StatusRetired {
		log.Error(pool.ctx, err, "Resource with ID %d is in an incorrect state %+v", res.ID, res.Status)
		return nil, errors.Errorf("Resource #%d is in incorrect state \"%s\"", res.ID, res.Status)
	} else if res.Status == resource.StatusBench {
		cutoff := res.UpdatedAt.Add(time.Duration(pool.DealocationSafetyPeriod) * time.Second)
		if time.Now().Before(cutoff) {
			log.Error(pool.ctx, err, "Unable to claim resource %d from pool %d, resource cannot be claimed before %s", res.ID, pool.ID, cutoff)
			return nil, errors.Errorf(
				"Unable to claim resource #%d from pool #%d, resource cannot be claimed before %s", res.ID, pool.ID, cutoff)
		}
	}
	res.Status = resource.StatusClaimed
	err = pool.client.Resource.
		UpdateOne(res).
		SetStatus(res.Status).
		SetNillableDescription(description).
		SetAlternateID(alternativeId).
		Exec(pool.ctx)

	//TODO what does this mean? should we somehow rollback everything that transpired until this point??
	if err != nil {
		log.Error(pool.ctx, err, "Cannot update resource %d", res.ID)
		return nil, errors.Wrapf(err, "Cannot update resource #%d", res.ID)
	}
	return res, nil
}

func getFullListOfResources(pool AllocatingPool) ([]*model.ResourceInput, error) {
	return pool.loadClaimedResources()
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

func (pool AllocatingPool) loadClaimedResources() ([]*model.ResourceInput, error) {
	var currentResources []*model.ResourceInput
	claimedResources, err := pool.findResources().WithProperties(
		func(propertyQuery *ent.PropertyQuery) { propertyQuery.WithType() }).All(pool.ctx)
	if err != nil {
		log.Error(pool.ctx, err, "Unable to get claimed resources from pool %d", pool.ID)
		return nil, errors.Wrapf(err,
			"Unable to get claimed resources from pool #%d, resource loading error", pool.ID)
	}
	for _, claimedResource := range claimedResources {
		var r model.ResourceInput
		r.UpdatedAt = claimedResource.UpdatedAt.String()
		r.Status = claimedResource.Status.String()
		if propsToMap, err := PropertiesToMap(claimedResource.Edges.Properties); err != nil {
			log.Error(pool.ctx, err, "Unable to serialize resource properties")
			return nil, errors.Wrapf(err, "Unable to serialize resource properties")
		} else {
			r.Properties = propsToMap
		}

		currentResources = append(currentResources, &r)
	}
	return currentResources, nil
}

// FreeResource deallocates the resource identified by its properties
func (pool AllocatingPool) FreeResource(raw RawResourceProps) error {
	return pool.freeResourceInner(raw, pool.retireResource, pool.freeResourceImmediately, pool.benchResource)
}

func (pool AllocatingPool) freeResourceImmediately(res *ent.Resource) error {
	// Delete props
	for _, prop := range res.Edges.Properties {
		if err := pool.client.Property.DeleteOne(prop).Exec(pool.ctx); err != nil {
			log.Error(pool.ctx, err, "Cannot delete properties from on resource ID %d", res.ID)
			return errors.Wrapf(err, "Cannot free resource from #%d. Unable to cleanup properties", pool.ID)
		}
	}

	// Delete resource
	err := pool.client.Resource.DeleteOne(res).Exec(pool.ctx)
	if err != nil {
		log.Error(pool.ctx, err, "Cannot delete resource ID %d", res.ID)
		return errors.Wrapf(err, "Cannot free resource from #%d. Unable to cleanup resource", pool.ID)
	}

	return nil
}
