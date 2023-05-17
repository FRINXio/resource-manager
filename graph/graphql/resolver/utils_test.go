package resolver_test

// mockup test for functions inside utils.go file

import (
	"context"
	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/schema"
	"github.com/net-auto/resourceManager/graph/graphql/resolver"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getContext() context.Context {
	schema.InitializeAdminRoles("OWNER")
	ctx := context.Background()
	ctx = schema.WithIdentity(ctx, "fb", "fb-user", "ROLE1, OWNER, ABCD", "network-admin")
	return ctx
}

type testSetup struct {
	ctx    context.Context
	client *ent.Client
}

func setup(t *testing.T) testSetup {
	ctx := getContext()
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		t.Fatal(err)
	}
	if err := client.Schema.Create(ctx); err != nil {
		t.Fatal(err)
	}
	return testSetup{ctx: ctx, client: client}
}

func createResourceType(ctx context.Context, t *testing.T, client *ent.Client, name string) *ent.ResourceType {
	resourceType, err := client.ResourceType.Create().
		SetName(name).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	return resourceType
}

func createAllocationStrategy(ctx context.Context, t *testing.T, client *ent.Client, name string) *ent.AllocationStrategy {
	allocationStrategy, err := client.AllocationStrategy.Create().
		SetName(name).
		SetScript("Hello World!").
		SetLang("js").
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	return allocationStrategy
}

func createPropertyType(ctx context.Context, t *testing.T, client *ent.Client, name string, rtId int) *ent.PropertyType {
	propertyType, err := client.PropertyType.Create().
		SetName(name).
		SetType("string").
		SetStringVal("test").
		SetResourceTypeID(rtId).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	return propertyType
}

func createProperty(ctx context.Context, t *testing.T, client *ent.Client, name string, propertyId int) *ent.Property {
	property, err := client.Property.Create().
		SetTypeID(propertyId).
		SetStringVal("test").
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	return property
}

func createResource(ctx context.Context, t *testing.T, client *ent.Client, name string, poolId int, propertyId int) *ent.Resource {
	resource, err := client.Resource.Create().
		SetDescription(name).
		SetAlternateID(map[string]interface{}{"test": "test"}).
		SetPoolID(poolId).
		AddPropertyIDs(propertyId).
		SetStatus("claimed").
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	return resource
}

func createPoolProperties(ctx context.Context, t *testing.T, client *ent.Client, propertyId int, rtId int) *ent.PoolProperties {
	poolProperties, err := client.PoolProperties.Create().
		AddPropertyIDs(propertyId).
		AddResourceTypeIDs(rtId).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	return poolProperties
}

func createResourcePool(ctx context.Context, t *testing.T, client *ent.Client, name string, poolPropId int, rtId int, allocStratId int) *ent.ResourcePool {
	resourcePool, err := client.ResourcePool.Create().
		SetName(name).
		SetPoolType("allocating").
		SetPoolPropertiesID(poolPropId).
		SetResourceTypeID(rtId).
		SetAllocationStrategyID(allocStratId).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	return resourcePool
}

func TestPool_Filtering_ByResources(t *testing.T) {
	s := setup(t)
	defer s.client.Close()

	allocStrategy := createAllocationStrategy(s.ctx, t, s.client, "testAllocStrat")
	resourceType := createResourceType(s.ctx, t, s.client, "testType")
	propertyType := createPropertyType(s.ctx, t, s.client, "testPropType", resourceType.ID)
	property := createProperty(s.ctx, t, s.client, "testProp", propertyType.ID)
	poolProperties := createPoolProperties(s.ctx, t, s.client, property.ID, resourceType.ID)
	resourcePool := createResourcePool(s.ctx, t, s.client, "testPool", poolProperties.ID, resourceType.ID, allocStrategy.ID)
	claimedResource := createResource(s.ctx, t, s.client, "testResource", resourcePool.ID, property.ID)

	ids, err := resolver.FilterResourcePoolByAllocatedResources(s.ctx, s.client.ResourcePool.Query(), map[string]interface{}{"testPropType": "test"})

	assert.Nil(t, err)
	assert.Equal(t, []int{property.ID}, ids)

	// should return empty array
	ids2, err := resolver.FilterResourcePoolByAllocatedResources(s.ctx, s.client.ResourcePool.Query(), map[string]interface{}{"testPropType": "not-existing-value"})
	assert.Nil(t, err)
	assert.Equal(t, []int{}, ids2)

	// should return empty array when not existing property type name is passed
	ids3, err := resolver.FilterResourcePoolByAllocatedResources(s.ctx, s.client.ResourcePool.Query(), map[string]interface{}{"not-existing-prop-type": "test"})
	assert.Nil(t, err)
	assert.Equal(t, []int{}, ids3)

	assert.NotNil(t, claimedResource)
}
