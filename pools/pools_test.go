package pools

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/net-auto/resourceManager/ent/runtime"
)

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