package pools

import (
	"context"
	"fmt"
	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/allocationstrategy"
	"github.com/net-auto/resourceManager/ent/schema"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/net-auto/resourceManager/ent/runtime"
)

type mockInvoker struct {
	toBeReturned      map[string]interface{}
	toBeReturnedError error
}

func (m mockInvoker) invokeJs(strategyScript string) (map[string]interface{}, error) {
	if m.toBeReturnedError != nil {
		return nil, m.toBeReturnedError
	} else {
		return m.toBeReturned, nil
	}
}

func TestAllocatingPool(t *testing.T) {
	ctx := getContext()
	client := openDb(ctx)
	defer client.Close()
	resType := getResourceType(ctx, client)

	strat, _ := mockStrategy(client, ctx)
	propsAsMap := RawResourceProps{"vlan": 1}
	mockInvoker := mockInvoker{propsAsMap, nil}

	pool, _, err := newAllocatingPoolWithMetaInternal(
		ctx, client, resType, strat, "testAllocatingPool",
		mockInvoker, schema.ResourcePoolDealocationImmediatelly)

	if err != nil {
		t.Fatalf("Unable to create pool %s", err)
	}

	assertDb(ctx, client, t, 1, 1, 1, 0, 0)

	resource, err := pool.ClaimResource()
	if err != nil {
		t.Fatalf("Unable to claim resource: %s", err)
	}
	assertDb(ctx, client, t, 1, 1, 1, 1, 1)

	props, _ := resource.QueryProperties().WithType().All(ctx)
	toMap, err := PropertiesToMap(props)

	if !reflect.DeepEqual(toMap, propsAsMap) {
		t.Fatalf("Unexpected props in claimed resource: %v, should be %v", toMap, propsAsMap)
	}

	err = pool.FreeResource(propsAsMap)
	if err != nil {
		t.Fatalf("Unable to free resource: %s", err)
	}

	assertDb(ctx, client, t, 1, 1, 1, 0, 0)
}

func mockStrategy(client *ent.Client, ctx context.Context) (*ent.AllocationStrategy, error) {
	return client.AllocationStrategy.Create().
		SetName("testStrat").
		SetLang(allocationstrategy.LangJs).
		SetScript("Hello World!").
		Save(ctx)
}

func TestAllocatingPoolFailure(t *testing.T) {
	ctx := getContext()
	client := openDb(ctx)
	defer client.Close()
	resType := getResourceType(ctx, client)

	strat, _ := mockStrategy(client, ctx)
	mockInvoker := mockInvoker{nil, fmt.Errorf("Fail")}

	pool, _, _ := newAllocatingPoolWithMetaInternal(
		ctx, client, resType, strat, "testAllocatingPool",
		mockInvoker, schema.ResourcePoolDealocationRetire)

	_, err := pool.ClaimResource()
	if err == nil {
		t.Fatalf("Resource claim should have failed")
	}

	assertDb(ctx, client, t, 1, 1, 1, 0, 0)
}
