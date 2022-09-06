package pools

import (
	"context"
	"fmt"
	"github.com/net-auto/resourceManager/graph/graphql/model"
	"reflect"
	"testing"
	"time"

	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/allocationstrategy"
	"github.com/net-auto/resourceManager/ent/resource"
	"github.com/net-auto/resourceManager/ent/schema"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/net-auto/resourceManager/ent/runtime"
)

type mockInvoker struct {
	toBeReturned      map[string]interface{}
	toBeReturnedError error
}

func (m mockInvoker) invokeJs(
	strategyScript string,
	userInput map[string]interface{},
	resourcePool model.ResourcePoolInput,
	currentResources []*model.ResourceInput,
	poolPropertiesMaps map[string]interface{},
	functionName string,
) (map[string]interface{}, string, error) {
	if m.toBeReturnedError != nil {
		return nil, "", m.toBeReturnedError
	} else {
		return m.toBeReturned, "", nil
	}
}

func (m mockInvoker) invokePy(
	strategyScript string,
	userInput map[string]interface{},
	resourcePool model.ResourcePoolInput,
	currentResources []*model.ResourceInput,
	poolPropertiesMaps map[string]interface{},
	functionName string,
) (map[string]interface{}, string, error) {
	if m.toBeReturnedError != nil {
		return nil, "", m.toBeReturnedError
	} else {
		return m.toBeReturned, "", nil
	}
}

type testSetup struct {
	ctx    context.Context
	client *ent.Client
	pool   Pool
}

func mockStrategy(client *ent.Client, ctx context.Context) (*ent.AllocationStrategy, error) {
	return client.AllocationStrategy.Create().
		SetName("testStrat").
		SetLang(allocationstrategy.LangJs).
		SetScript("Hello World!").
		Save(ctx)
}

func CreateTestSetup(t *testing.T, mockInvoker mockInvoker, poolDealocationSafetyPeriod int) testSetup {
	ctx := getContext()
	client := openDb(ctx)

	resType, err := getResourceType(ctx, client)
	if err != nil {
		t.Fatalf("Unable to create resource type: %s", err)
	}

	strat, err := mockStrategy(client, ctx)
	if err != nil {
		t.Fatalf("Unable to create mock strategy: %s", err)
	}

	pool, _, err := newAllocatingPoolWithMetaInternal(
		ctx, client, resType, strat, "testAllocatingPool", nil,
		mockInvoker, poolDealocationSafetyPeriod, nil)

	if err != nil {
		t.Fatalf("Unable to create pool %s", err)
	}
	assertInstancesInDb(client.PropertyType.Query().AllX(ctx), 1, t)
	assertInstancesInDb(client.ResourceType.Query().AllX(ctx), 1, t)
	assertInstancesInDb(client.ResourcePool.Query().AllX(ctx), 1, t)

	return testSetup{ctx: ctx, client: client, pool: pool}
}

func (ts testSetup) Close() {
	ts.client.Close()
}

func TestAllocatingPool_ReclaimImmediately(t *testing.T) {
	propsAsMap := RawResourceProps{"vlan": 1}
	mockInvoker := mockInvoker{propsAsMap, nil}
	ts := CreateTestSetup(t, mockInvoker, schema.ResourcePoolDealocationImmediately)
	defer ts.Close()

	userInput := make(map[string]interface{})
	resource, err := ts.pool.ClaimResource(userInput, nil, nil)
	if err != nil {
		t.Fatalf("Unable to claim resource: %s", err)
	}
	assertInstancesInDb(ts.client.Resource.Query().AllX(ts.ctx), 1, t)
	assertInstancesInDb(ts.client.Property.Query().AllX(ts.ctx), 1, t)

	props, _ := resource.QueryProperties().WithType().All(ts.ctx)
	toMap, err := PropertiesToMap(props)

	if !reflect.DeepEqual(toMap, propsAsMap) {
		t.Fatalf("Unexpected props in claimed resource: %v, should be %v", toMap, propsAsMap)
	}

	err = ts.pool.FreeResource(propsAsMap)
	if err != nil {
		t.Fatalf("Unable to free resource: %s", err)
	}

	assertInstancesInDb(ts.client.Resource.Query().AllX(ts.ctx), 0, t)
	assertInstancesInDb(ts.client.Property.Query().AllX(ts.ctx), 0, t)

	// reclaiming is possible
	_, err = ts.pool.ClaimResource(userInput, nil, nil)
	if err != nil {
		t.Fatalf("Unable to claim resource: %s", err)
	}
	assertInstancesInDb(ts.client.Resource.Query().AllX(ts.ctx), 1, t)
	assertInstancesInDb(ts.client.Property.Query().AllX(ts.ctx), 1, t)
}

func TestAllocatingPool_ScriptFailure(t *testing.T) {
	mockInvoker := mockInvoker{nil, fmt.Errorf("Fail")}
	ts := CreateTestSetup(t, mockInvoker, schema.ResourcePoolDealocationImmediately)
	defer ts.Close()

	userInput := make(map[string]interface{})
	_, err := ts.pool.ClaimResource(userInput, nil, nil)
	if err == nil {
		t.Fatalf("Resource claim should have failed")
	}

	assertInstancesInDb(ts.client.Resource.Query().AllX(ts.ctx), 0, t)
	assertInstancesInDb(ts.client.Property.Query().AllX(ts.ctx), 0, t)
}

func TestAllocatingPool_RetiredResource(t *testing.T) {
	propsAsMap := RawResourceProps{"vlan": 1}
	mockInvoker := mockInvoker{propsAsMap, nil}
	ts := CreateTestSetup(t, mockInvoker, schema.ResourcePoolDealocationRetire)
	defer ts.Close()

	userInput := make(map[string]interface{})

	// claim resource vlan:1
	_, err := ts.pool.ClaimResource(userInput, nil, nil)
	if err != nil {
		t.Fatalf("Unable to claim resource: %s", err)
	}
	assertInstancesInDb(ts.client.Resource.Query().AllX(ts.ctx), 1, t)

	// free it (should be retired)
	err = ts.pool.FreeResource(propsAsMap)
	if err != nil {
		t.Fatalf("Unable to free resource: %s", err)
	}

	assertInstancesInDb(ts.client.Resource.Query().Where(resource.StatusEQ(resource.StatusRetired)).AllX(ts.ctx), 1, t)

	// it should not be allowed to be reclaimed
	_, err = ts.pool.ClaimResource(userInput, nil, nil)
	if err == nil {
		t.Fatalf("Second resource claim should have failed")
	}
	assertInstancesInDb(ts.client.Resource.Query().Where(resource.StatusEQ(resource.StatusRetired)).AllX(ts.ctx), 1, t)
}

func TestAllocatingPool_WithSafetyPeriod(t *testing.T) {
	propsAsMap := RawResourceProps{"vlan": 1}
	mockInvoker := mockInvoker{propsAsMap, nil}
	safetySeconds := 2
	ts := CreateTestSetup(t, mockInvoker, safetySeconds)
	defer ts.Close()

	userInput := make(map[string]interface{})

	// claim resource vlan:1
	_, err := ts.pool.ClaimResource(userInput, nil, nil)
	if err != nil {
		t.Fatalf("Unable to claim resource: %s", err)
	}
	assertInstancesInDb(ts.client.Resource.Query().AllX(ts.ctx), 1, t)

	// free it (should sit on bench for safetySeconds)
	err = ts.pool.FreeResource(propsAsMap)
	if err != nil {
		t.Fatalf("Unable to free resource: %s", err)
	}

	assertInstancesInDb(ts.client.Resource.Query().Where(resource.StatusEQ(resource.StatusBench)).AllX(ts.ctx), 1, t)

	// it should not be allowed to be reclaimed immediately
	_, err = ts.pool.ClaimResource(userInput, nil, nil)
	if err == nil {
		t.Fatalf("Second resource claim should have failed")
	}
	assertInstancesInDb(ts.client.Resource.Query().Where(resource.StatusEQ(resource.StatusBench)).AllX(ts.ctx), 1, t)
	// sleep
	time.Sleep(time.Duration(safetySeconds) * time.Second)
	// reclaim
	_, err = ts.pool.ClaimResource(userInput, nil, nil)
	if err != nil {
		t.Fatalf("Unable to claim resource: %s", err)
	}
	assertInstancesInDb(ts.client.Resource.Query().AllX(ts.ctx), 1, t)
	assertInstancesInDb(ts.client.Resource.Query().Where(resource.StatusEQ(resource.StatusClaimed)).AllX(ts.ctx), 1, t)
}
