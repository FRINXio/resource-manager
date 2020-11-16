package pools

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/net-auto/resourceManager/ent/runtime"
)

func TestRbacFail(t *testing.T) {
	ctx := getContextWithFailingRbac()
	client := openDb(ctx)
	defer client.Close()
	_, err := getResourceType(ctx, client)
	if err == nil {
		t.Fatalf("Rbac should have failed")
	}
}

func TestRbacSuccess(t *testing.T) {
	ctx := getContext()
	client := openDb(ctx)
	defer client.Close()
	_, err := getResourceType(ctx, client)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}
