package pools

import (
	"context"

	"github.com/net-auto/resourceManager/ent/schema"
	"log"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/resource"
	_ "github.com/net-auto/resourceManager/ent/runtime"
)

func getContext() context.Context {
	schema.InitializeAdminRoles("OWNER")
	ctx := context.Background()
	ctx = schema.WithIdentity(ctx, "fb", "fb-user", "ROLE1, OWNER, ABCD", "network-admin")
	return ctx
}

func getContextWithFailingRbac() context.Context {
	schema.InitializeAdminRoles("ROLE NOT MATCHING")
	ctx := context.Background()
	ctx = schema.WithIdentity(ctx, "fb", "fb-user", "ROLE1, OWNER, ABCD", "network-admin")
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

func getResourceType(ctx context.Context, client *ent.Client) (*ent.ResourceType, error) {
	propType, err := client.PropertyType.Create().
		SetName("vlan").
		SetType("int").
		SetIntVal(0).
		SetMandatory(true).
		Save(ctx)

	if err != nil {
		log.Printf("Failed to create property type: %v", err)
		return nil, err
	}
	resType, err := client.ResourceType.Create().
		SetName("vlan").
		AddPropertyTypes(propType).
		Save(ctx)
	if err != nil {
		log.Printf("Failed to create resource type: %v", err)
		return nil, err
	}

	return resType, nil
}

func assertDb(ctx context.Context, client *ent.Client, t *testing.T, count ...int) {
	assertInstancesInDb(client.PropertyType.Query().AllX(ctx), count[0], t)
	assertInstancesInDb(client.ResourceType.Query().AllX(ctx), count[1], t)
	assertInstancesInDb(client.ResourcePool.Query().AllX(ctx), count[2], t)
	assertInstancesInDb(client.Property.Query().AllX(ctx), count[3], t)
	assertInstancesInDb(client.Resource.Query().AllX(ctx), count[4], t)
}

func assertDbResourceStates(ctx context.Context, client *ent.Client, t *testing.T, count ...int) {
	assertDbResourceState(ctx, client, t, count[0], resource.StatusFree)
	assertDbResourceState(ctx, client, t, count[1], resource.StatusClaimed)
	assertDbResourceState(ctx, client, t, count[2], resource.StatusBench)
	assertDbResourceState(ctx, client, t, count[3], resource.StatusRetired)
}

func assertDbResourceState(ctx context.Context, client *ent.Client, t *testing.T, expected int, expectedStatus resource.Status) {
	freeResourceCount, _ := client.Resource.Query().Where(resource.StatusEQ(expectedStatus)).Count(ctx)
	if freeResourceCount != expected {
		t.Fatalf("%d different instances of resources in state: %s expected, got: %d",
			expected, expectedStatus.String(), freeResourceCount)
	}
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

func OpenTestDb(ctx context.Context) *ent.Client {
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
