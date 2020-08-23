package pools

import (
	"testing"
	"time"

	"github.com/net-auto/resourceManager/ent/schema"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/net-auto/resourceManager/ent/runtime"
)

func TestClaimResoourceSetPool(t *testing.T) {
	ctx := getContext()
	client := openDb(ctx)
	defer client.Close()
	resType, err := getResourceType(ctx, client)
	if err != nil {
		t.Fatalf("Unable to create resource type: %s", err)
	}

	pool, _ := NewSetPool(ctx, client, resType, []RawResourceProps{
		RawResourceProps{"vlan": 44},
		RawResourceProps{"vlan": 45},
	}, "set", nil, schema.ResourcePoolDealocationImmediately)

	claims, err := pool.QueryResources()
	if len(claims) != 0 {
		t.Fatalf("Expected 0 claims, got: %d", len(claims))
	}
	userInput := make(map[string]interface{})
	claim1, err := pool.ClaimResource(userInput)
	t.Log(claim1)
	claims, err = pool.QueryResources()
	if len(claims) != 1 {
		t.Fatalf("Expected 1 claims, got: %d", len(claims))
	}
	claim2, err := pool.ClaimResource(userInput)
	t.Log(claim2)

	if _, err := pool.ClaimResource(userInput); err == nil {
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

func TestResourceRetirement(t *testing.T) {
	ctx := getContext()
	client := openDb(ctx)
	defer client.Close()
	resType, err := getResourceType(ctx, client)
	if err != nil {
		t.Fatalf("Unable to create resource type: %s", err)
	}

	description := "testPool"
	pool, _ := NewSetPool(ctx, client, resType, []RawResourceProps{
		RawResourceProps{"vlan": 44},
		RawResourceProps{"vlan": 45},
	}, "set", &description, schema.ResourcePoolDealocationRetire)

	userInput := make(map[string]interface{})
	claim1, _ := pool.ClaimResource(userInput)
	claim2, _ := pool.ClaimResource(userInput)

	pool.FreeResource(RawResourceProps{"vlan": *claim1.QueryProperties().AllX(ctx)[0].IntVal})
	pool.FreeResource(RawResourceProps{"vlan": *claim2.QueryProperties().AllX(ctx)[0].IntVal})

	if ress, _ := pool.QueryResources(); len(ress) > 0 {
		t.Fatalf("There should be no resources returned since all have been freed")
	}

	if _, err := pool.ClaimResource(userInput); err == nil {
		t.Fatalf("Expecting error, since all resources should have been retired")
	}
}

func TestResourceDealocationSafetyWindow(t *testing.T) {
	ctx := getContext()
	client := openDb(ctx)
	defer client.Close()
	resType, err := getResourceType(ctx, client)
	if err != nil {
		t.Fatalf("Unable to create resource type: %s", err)
	}

	pool, _ := NewSetPool(ctx, client, resType, []RawResourceProps{
		RawResourceProps{"vlan": 44},
		RawResourceProps{"vlan": 45},
	}, "set", nil,3)

	// Claim and free resource 44
	userInput := make(map[string]interface{})
	claim1, _ := pool.ClaimResource(userInput)
	assertDbResourceStates(ctx, client, t, 1, 1, 0, 0)

	pool.FreeResource(RawResourceProps{"vlan": *claim1.QueryProperties().AllX(ctx)[0].IntVal})
	assertDbResourceStates(ctx, client, t, 1, 0, 1, 0)

	if ress, _ := pool.QueryResources(); len(ress) > 0 {
		t.Fatalf("There should be no resources returned since all have been freed")
	}

	// Claim and free resource ... should be 45 since 44 is free but benched
	claim2, _ := pool.ClaimResource(userInput)
	assertDbResourceStates(ctx, client, t, 0, 1, 1, 0)

	if *claim2.QueryProperties().AllX(ctx)[0].IntVal != 45 {
		t.Fatalf("Wrong property in resource claim: %s. "+
			"Expected 45 since free resources should have priority over benched", claim1)
	}
	pool.FreeResource(RawResourceProps{"vlan": *claim2.QueryProperties().AllX(ctx)[0].IntVal})
	assertDbResourceStates(ctx, client, t, 0, 0, 2, 0)

	// All are benched
	if _, err := pool.ClaimResource(userInput); err == nil {
		t.Fatalf("Expecting error, since all resources should have been retired")
	}

	// Waiting for resources to become available
	time.Sleep(time.Duration(4) * time.Second)
	// Resources are not marked free automatically, only during resource claim call
	assertDbResourceStates(ctx, client, t, 0, 0, 2, 0)

	claim1Again, _ := pool.ClaimResource(userInput)
	assertDbResourceStates(ctx, client, t, 0, 1, 1, 0)

	pool.FreeResource(RawResourceProps{"vlan": *claim1Again.QueryProperties().AllX(ctx)[0].IntVal})
	assertDbResourceStates(ctx, client, t, 0, 0, 2, 0)
}
