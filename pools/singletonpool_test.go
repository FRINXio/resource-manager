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
	resType, err := getResourceType(ctx, client)
	if err != nil {
		t.Fatalf("Unable to create resource type: %s", err)
	}

	pool, err := NewSingletonPool(ctx, client, resType, map[string]interface{}{
		"vlan": 44,
	}, "singleton", nil)

	if err != nil {
		t.Fatal(err)
	}

	err = pool.Destroy()
	if err != nil {
		t.Fatal(err)
	}
}
