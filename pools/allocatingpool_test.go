package pools

import (
	"context"
	"fmt"
	"github.com/facebook/ent/dialect/sql"
	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/allocationstrategy"
	"github.com/net-auto/resourceManager/ent/property"
	"github.com/net-auto/resourceManager/ent/resource"
	resourcePool "github.com/net-auto/resourceManager/ent/resourcepool"
	"github.com/net-auto/resourceManager/graph/graphql/model"
	logger "github.com/net-auto/resourceManager/logging"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/net-auto/resourceManager/ent/runtime"
)

type mockInvoker struct {
	toBeReturned      []map[string]interface{}
	toBeReturnedError error
	counter           int
}

func (m *mockInvoker) invokeJs(
	strategyScript string,
	userInput map[string]interface{},
	resourcePool model.ResourcePoolInput,
	currentResources []*model.ResourceInput,
	poolPropertiesMaps map[string]interface{},
	functionName string,
) (map[string]interface{}, string, error) {
	m.counter++
	if m.toBeReturnedError != nil {
		return nil, "", m.toBeReturnedError
	} else {
		return m.toBeReturned[m.counter], "", nil
	}
}

func (m *mockInvoker) invokePy(
	strategyScript string,
	userInput map[string]interface{},
	resourcePool model.ResourcePoolInput,
	currentResources []*model.ResourceInput,
	poolPropertiesMaps map[string]interface{},
	functionName string,
) (map[string]interface{}, string, error) {
	m.counter++

	if m.toBeReturnedError != nil {
		return nil, "", m.toBeReturnedError
	} else {
		return m.toBeReturned[m.counter], "", nil
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
	logger.Init("a.log", "debug", false)

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
		&mockInvoker, poolDealocationSafetyPeriod, nil)

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

//func TestAllocatingPool_ReclaimImmediately(t *testing.T) {
//	propsAsMap := RawResourceProps{"vlan": 1}
//	mockInvoker := mockInvoker{propsAsMap, nil}
//	ts := CreateTestSetup(t, mockInvoker, schema.ResourcePoolDealocationImmediately)
//	defer ts.Close()
//
//	userInput := make(map[string]interface{})
//	resource, err := ts.pool.ClaimResource(userInput, nil, nil)
//	if err != nil {
//		t.Fatalf("Unable to claim resource: %s", err)
//	}
//	assertInstancesInDb(ts.client.Resource.Query().AllX(ts.ctx), 1, t)
//	assertInstancesInDb(ts.client.Property.Query().AllX(ts.ctx), 1, t)
//
//	props, _ := resource.QueryProperties().WithType().All(ts.ctx)
//	toMap, err := PropertiesToMap(props)
//
//	if !reflect.DeepEqual(toMap, propsAsMap) {
//		t.Fatalf("Unexpected props in claimed resource: %v, should be %v", toMap, propsAsMap)
//	}
//
//	err = ts.pool.FreeResource(propsAsMap)
//	if err != nil {
//		t.Fatalf("Unable to free resource: %s", err)
//	}
//
//	assertInstancesInDb(ts.client.Resource.Query().AllX(ts.ctx), 0, t)
//	assertInstancesInDb(ts.client.Property.Query().AllX(ts.ctx), 0, t)
//
//	// reclaiming is possible
//	_, err = ts.pool.ClaimResource(userInput, nil, nil)
//	if err != nil {
//		t.Fatalf("Unable to claim resource: %s", err)
//	}
//	assertInstancesInDb(ts.client.Resource.Query().AllX(ts.ctx), 1, t)
//	assertInstancesInDb(ts.client.Property.Query().AllX(ts.ctx), 1, t)
//}

//func TestAllocatingPool_ScriptFailure(t *testing.T) {
//	mockInvoker := mockInvoker{nil, fmt.Errorf("Fail")}
//	ts := CreateTestSetup(t, mockInvoker, schema.ResourcePoolDealocationImmediately)
//	defer ts.Close()
//
//	userInput := make(map[string]interface{})
//	_, err := ts.pool.ClaimResource(userInput, nil, nil)
//	if err == nil {
//		t.Fatalf("Resource claim should have failed")
//	}
//
//	assertInstancesInDb(ts.client.Resource.Query().AllX(ts.ctx), 0, t)
//	assertInstancesInDb(ts.client.Property.Query().AllX(ts.ctx), 0, t)
//}

//func TestAllocatingPool_RetiredResource(t *testing.T) {
//	propsAsMap := RawResourceProps{"vlan": 1}
//	mockInvoker := mockInvoker{propsAsMap, nil}
//	ts := CreateTestSetup(t, mockInvoker, schema.ResourcePoolDealocationRetire)
//	defer ts.Close()
//
//	userInput := make(map[string]interface{})
//
//	// claim resource vlan:1
//	_, err := ts.pool.ClaimResource(userInput, nil, nil)
//	if err != nil {
//		t.Fatalf("Unable to claim resource: %s", err)
//	}
//	assertInstancesInDb(ts.client.Resource.Query().AllX(ts.ctx), 1, t)
//
//	// free it (should be retired)
//	err = ts.pool.FreeResource(propsAsMap)
//	if err != nil {
//		t.Fatalf("Unable to free resource: %s", err)
//	}
//
//	assertInstancesInDb(ts.client.Resource.Query().Where(resource.StatusEQ(resource.StatusRetired)).AllX(ts.ctx), 1, t)
//
//	// it should not be allowed to be reclaimed
//	_, err = ts.pool.ClaimResource(userInput, nil, nil)
//	if err == nil {
//		t.Fatalf("Second resource claim should have failed")
//	}
//	assertInstancesInDb(ts.client.Resource.Query().Where(resource.StatusEQ(resource.StatusRetired)).AllX(ts.ctx), 1, t)
//}

func TestTTTTTTTTTTTTTTTT(t *testing.T) {
	allProps := make([]map[string]interface{}, 6)
	allProps[0] = RawResourceProps{"vlan": 1}
	allProps[1] = RawResourceProps{"vlan": 2}
	allProps[2] = RawResourceProps{"vlan": 3}
	allProps[3] = RawResourceProps{"vlan": 4}
	allProps[4] = RawResourceProps{"vlan": 5}
	allProps[5] = RawResourceProps{"vlan": 6}

	mockInvoker := mockInvoker{allProps, nil, 0}

	ts := CreateTestSetup(t, mockInvoker, 2)
	pool := ts.pool.(*AllocatingPool)
	userInput := make(map[string]interface{})

	pool.ClaimResource(userInput, nil, nil)
	pool.ClaimResource(RawResourceProps{"vlan": 2}, nil, nil)
	pool.ClaimResource(RawResourceProps{"vlan": 3}, nil, nil)
	pool.ClaimResource(RawResourceProps{"vlan": 4}, nil, nil)

	var x []struct {
		Id     int `json:"id,omitempty"`
		IntVal int `json:"int_val,omitempty"`
	}

	err := ts.client.Property.Query().Select("id", "int_val").Scan(pool.ctx, &x)
	fmt.Println(err)
	fmt.Println(x)

	var v []struct {
		//Id     int `json:"id,omitempty"`
		IntVal int    `json:"int_val,omitempty"`
		Status string `json:"status,omitempty"`
	}

	ts.client.Resource.Query().
		Where(
			resource.HasPoolWith(resourcePool.ID(pool.ID)),
			func(resourceSelector *sql.Selector) {
				t := sql.Table(property.Table)
				resourceSelector.Join(t).On(resourceSelector.C(resource.FieldID), t.C(property.ResourcesColumn))
			},
		).Select("`resources`.`status`", "`t0`.`int_val`").Scan(pool.ctx, &v)

	fmt.Println(v)

	all, _ := ts.client.Resource.Query().
		Where(
			resource.HasPoolWith(resourcePool.ID(pool.ID)),
			func(resourceSelector *sql.Selector) {
				t := sql.Table(property.Table)
				resourceSelector.Join(t).On(resourceSelector.C(resource.FieldID), t.C(property.ResourcesColumn)).As("properties")
			},
		).CollectFields(pool.ctx, "properties").
		//.WithProperties(func(propertyQuery *ent.PropertyQuery) { propertyQuery.WithType() }).
		All(pool.ctx)

	fmt.Println(len(all))

	for _, r := range all {
		fmt.Println(r)
		fmt.Println(r.Edges.Properties)
		for _, p := range r.Edges.Properties {
			fmt.Println(p.Edges.Type)
		}
	}
}

//func TestAllocatingPool_WithSafetyPeriod(t *testing.T) {
//	propsAsMap := RawResourceProps{"vlan": 1}
//
//	mockInvoker := mockInvoker{propsAsMap, nil}
//	safetySeconds := 2
//	ts := CreateTestSetup(t, mockInvoker, safetySeconds)
//	defer ts.Close()
//
//	userInput := make(map[string]interface{})
//
//	// claim resource vlan:1
//	_, err := ts.pool.ClaimResource(userInput, nil, nil)
//	if err != nil {
//		t.Fatalf("Unable to claim resource: %s", err)
//	}
//	assertInstancesInDb(ts.client.Resource.Query().AllX(ts.ctx), 1, t)
//
//	// free it (should sit on bench for safetySeconds)
//	err = ts.pool.FreeResource(propsAsMap)
//	if err != nil {
//		t.Fatalf("Unable to free resource: %s", err)
//	}
//
//	assertInstancesInDb(ts.client.Resource.Query().Where(resource.StatusEQ(resource.StatusBench)).AllX(ts.ctx), 1, t)
//
//	// it should not be allowed to be reclaimed immediately
//	_, err = ts.pool.ClaimResource(userInput, nil, nil)
//	if err == nil {
//		t.Fatalf("Second resource claim should have failed")
//	}
//	assertInstancesInDb(ts.client.Resource.Query().Where(resource.StatusEQ(resource.StatusBench)).AllX(ts.ctx), 1, t)
//	// sleep
//	time.Sleep(time.Duration(safetySeconds) * time.Second)
//	// reclaim
//	_, err = ts.pool.ClaimResource(userInput, nil, nil)
//	if err != nil {
//		t.Fatalf("Unable to claim resource: %s", err)
//	}
//	assertInstancesInDb(ts.client.Resource.Query().AllX(ts.ctx), 1, t)
//	assertInstancesInDb(ts.client.Resource.Query().Where(resource.StatusEQ(resource.StatusClaimed)).AllX(ts.ctx), 1, t)
//}
