package pools

import (
	"context"
	"log"
	"reflect"
	"testing"

	"github.com/facebookincubator/symphony/pkg/authz"
	"github.com/facebookincubator/symphony/pkg/authz/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/net-auto/resourceManager/ent"
	_ "github.com/net-auto/resourceManager/ent/runtime"
)

func getContext() context.Context {
	ctx := context.Background()
	ctx = authz.NewContext(ctx, &models.PermissionSettings{
		CanWrite:        true,
		WorkforcePolicy: authz.NewWorkforcePolicy(true, true)})
	return ctx
}

func openDb(ctx context.Context) *ent.Client {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	// run the auto migration tool.
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	return client
}

func getResourceType(ctx context.Context, client *ent.Client) *ent.ResourceType {
	propType, _ := client.PropertyType.Create().
		SetName("vlan").
		SetType("int").
		SetIntVal(0).
		SetMandatory(true).
		Save(ctx)

	resType, _ := client.ResourceType.Create().
		SetName("vlan").
		AddPropertyTypes(propType).
		Save(ctx)

	return resType
}

func TestNewSingletonPool(t *testing.T) {
	ctx := getContext()
	client := openDb(ctx)
	defer client.Close()
	resType := getResourceType(ctx, client)

	pool, err := NewSingletonPool(ctx, client, resType, map[string]interface{}{
		"vlan": 44,
	}, "singleton")

	if err != nil {
		t.Fatal(err)
	}

	err = pool.Destroy()
	if err != nil {
		t.Fatal(err)
	}
}

func TestClaimResoourceSetPool(t *testing.T) {
	ctx := getContext()
	client := openDb(ctx)
	defer client.Close()
	resType := getResourceType(ctx, client)

	pool, _ := NewSetPool(ctx, client, resType, []RawResourceProps{
		RawResourceProps{"vlan": 44},
		RawResourceProps{"vlan": 45},
	}, "set")

	claims, err := pool.QueryResources()
	if len(claims) != 0 {
		t.Fatalf("Expected 0 claims, got: %d", len(claims))
	}

	claim1, err := pool.ClaimResource()
	t.Log(claim1)
	claims, err = pool.QueryResources()
	if len(claims) != 1 {
		t.Fatalf("Expected 1 claims, got: %d", len(claims))
	}
	claim2, err := pool.ClaimResource()
	t.Log(claim2)

	if _, err := pool.ClaimResource(); err == nil {
		t.Fatalf("Claiming resource from exhausted pool should return error")
	}

	entityPool := pool.(*SetPool).ResourcePool
	if claim1.QueryPool().OnlyX(ctx).ID != entityPool.ID {
		t.Fatalf("Wrong resource pool set expected: %s but was: %s",
			entityPool, claim1.QueryPool().OnlyX(ctx))
	}
	if claim2.QueryPool().OnlyX(ctx).ID != entityPool.ID {
		t.Fatalf("Wrong resource pool set expected %s but was: %s",
			entityPool, claim2.QueryPool().OnlyX(ctx))
	}

	claimProps1 := claim1.QueryProperties().AllX(ctx)
	claimProps2 := claim2.QueryProperties().AllX(ctx)
	if len(claimProps1) != 1 {
		t.Fatalf("Missing properties in resource claim: %s", claim1)
	}
	if *claimProps1[0].IntVal > 45 || *claimProps1[0].IntVal < 44 {
		t.Fatalf("Wrong property in resource claim: %s", claim1)
	}
	if *claimProps2[0].IntVal > 45 || *claimProps2[0].IntVal < 44 {
		t.Fatalf("Wrong property in resource claim: %s", claim1)
	}

	claims, err = pool.QueryResources()
	if len(claims) != 2 {
		t.Fatalf("Expected 2 claims, got: %d", len(claims))
	}
	assertDb(ctx, client, t, 1, 1, 1, 2, 2)

	if err := pool.Destroy(); err == nil {
		t.Fatalf("Destroying pool with active claims should return error")
	}

	if err = pool.FreeResource(RawResourceProps{"vlan": *claimProps1[0].IntVal}); err != nil {
		t.Fatal(err)
	}
	pool.FreeResource(RawResourceProps{"vlan": *claimProps2[0].IntVal})
	assertDb(ctx, client, t, 1, 1, 1, 2, 2)

	if err = pool.Destroy(); err != nil {
		t.Error(err)
	}
	assertDb(ctx, client, t, 1, 1, 0, 0, 0)
}

func assertDb(ctx context.Context, client *ent.Client, t *testing.T, count ...int) {
	assertInstancesInDb(client.PropertyType.Query().AllX(ctx), count[0], t)
	assertInstancesInDb(client.ResourceType.Query().AllX(ctx), count[1], t)
	assertInstancesInDb(client.ResourcePool.Query().AllX(ctx), count[2], t)
	assertInstancesInDb(client.Property.Query().AllX(ctx), count[3], t)
	assertInstancesInDb(client.Resource.Query().AllX(ctx), count[4], t)
}

func assertInstancesInDb(instances interface{}, expected int, t *testing.T) {
	slice := reflect.ValueOf(instances)
	if slice.Kind() != reflect.Slice {
		t.Fatalf("%s is not a slice, cannot assert length", instances)
	}

	if slice.Len() != expected {
		t.Fatalf("%d different instances of %s expected, got: %s", expected, slice.Type(), slice)
	}
}
