package resolver_test

// mockup test for functions inside utils.go file

import (
	"context"
	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/allocationstrategy"
	"github.com/net-auto/resourceManager/ent/property"
	"github.com/net-auto/resourceManager/ent/resource"
	resourcePool "github.com/net-auto/resourceManager/ent/resourcepool"
	"github.com/net-auto/resourceManager/ent/resourcetype"
	"github.com/net-auto/resourceManager/ent/schema"
	"github.com/net-auto/resourceManager/graph/graphql/resolver"
	pools2 "github.com/net-auto/resourceManager/pools"
	pools "github.com/net-auto/resourceManager/pools/allocating_strategies"
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

func TestPool_Filtering_ByResources(t *testing.T) {
	s := setup(t)
	defer s.client.Close()

	err := pools.LoadBuiltinTypes(s.ctx, s.client)

	if err != nil {
		t.Fatal(err)
	}

	ipv4PropertyType, err := s.client.ResourceType.Query().Where(resourcetype.Name("ipv4_prefix")).Only(s.ctx)
	ipv6PropertyType, err1 := s.client.ResourceType.Query().Where(resourcetype.Name("ipv6_prefix")).Only(s.ctx)

	if err != nil || err1 != nil {
		t.Fatal(err, err1)
	}

	ipv4PoolProperties, err := pools2.CreatePoolProperties(s.ctx, s.client, []map[string]interface{}{
		{
			"address": "10.0.0.0",
			"prefix":  24,
			"subnet":  true,
		},
	}, ipv4PropertyType)
	ipv6PoolProperties, err1 := pools2.CreatePoolProperties(s.ctx, s.client, []map[string]interface{}{
		{
			"address": "2001:db8::",
			"prefix":  64,
			"subnet":  true,
		},
	}, ipv6PropertyType)

	if err != nil || err1 != nil {
		t.Fatal(err, err1)
	}

	var description *string = nil
	ipv4ResourceType, err := s.client.ResourceType.Query().Where(resourcetype.Name("ipv4_prefix")).Only(s.ctx)
	ipv6ResourceType, err1 := s.client.ResourceType.Query().Where(resourcetype.Name("ipv6_prefix")).Only(s.ctx)

	if err != nil || err1 != nil {
		t.Fatal(err, err1)
	}

	ipv4AllocationStrategy, err := s.client.AllocationStrategy.Query().Where(allocationstrategy.Name("ipv4_prefix")).Only(s.ctx)
	ipv6AllocationStrategy, err1 := s.client.AllocationStrategy.Query().Where(allocationstrategy.Name("ipv6_prefix")).Only(s.ctx)

	if err != nil || err1 != nil {
		t.Fatal(err, err1)
	}

	ipv4ResPool, _, err := pools2.NewAllocatingPoolWithMeta(s.ctx, s.client, ipv4ResourceType, ipv4AllocationStrategy, "test_ipv4", description, 0, ipv4PoolProperties)
	ipv6ResPool, _, err1 := pools2.NewAllocatingPoolWithMeta(s.ctx, s.client, ipv6ResourceType, ipv6AllocationStrategy, "test_ipv6", description, 0, ipv6PoolProperties)

	if err != nil || err1 != nil {
		t.Fatal(err, err1)
	}

	_, err = ipv4ResPool.ClaimResource(map[string]interface{}{
		"desiredSize": 5,
	}, description, map[string]interface{}{"status": "blacklisted"})

	_, err1 = ipv6ResPool.ClaimResource(map[string]interface{}{
		"desiredSize": 5,
	}, description, map[string]interface{}{"status": "blacklisted"})

	if err != nil || err1 != nil {
		t.Fatal(err, err1)
	}

	expectedIpv4PoolIds, err := s.client.ResourcePool.Query().Where(resourcePool.HasClaimsWith(resource.HasPropertiesWith(property.StringValContains("10.0.0.0")))).IDs(s.ctx)
	expectedIpv6PoolIds, err1 := s.client.ResourcePool.Query().Where(resourcePool.HasClaimsWith(resource.HasPropertiesWith(property.StringValContains("2001:db8::")))).IDs(s.ctx)

	if err != nil || err1 != nil {
		t.Fatal(err, err1)
	}

	ipv4IDs, err := resolver.FilterResourcePoolByAllocatedResources(s.ctx, s.client.ResourcePool.Query(), map[string]interface{}{"address": "10.0.0.0"})
	ipv6IDs, err1 := resolver.FilterResourcePoolByAllocatedResources(s.ctx, s.client.ResourcePool.Query(), map[string]interface{}{"address": "2001:db8::"})

	if err != nil {
		t.Fatal(err)
	}

	actualFilteredPoolIds, err := s.client.ResourcePool.Query().Where(resourcePool.HasClaimsWith(resource.HasPropertiesWith(property.IDIn(ipv4IDs...)))).IDs(s.ctx)
	assert.Nil(t, err)
	assert.Equal(t, actualFilteredPoolIds, expectedIpv4PoolIds)

	actualFilteredIpv6PoolIds, err1 := s.client.ResourcePool.Query().Where(resourcePool.HasClaimsWith(resource.HasPropertiesWith(property.IDIn(ipv6IDs...)))).IDs(s.ctx)
	assert.Nil(t, err1)
	assert.Equal(t, actualFilteredIpv6PoolIds, expectedIpv6PoolIds)

	// should return empty array
	ids2, err := resolver.FilterResourcePoolByAllocatedResources(s.ctx, s.client.ResourcePool.Query(), map[string]interface{}{})
	assert.Nil(t, err)
	assert.Equal(t, 0, len(ids2))

	// should return empty array when not existing property type name is passed
	ids3, err := resolver.FilterResourcePoolByAllocatedResources(s.ctx, s.client.ResourcePool.Query(), map[string]interface{}{"not-existing-prop-type": "test"})
	assert.Nil(t, err)
	assert.Equal(t, []int{}, ids3)
}
