package pools

import (
	"context"
	"github.com/facebookincubator/symphony/pkg/authz"
	"github.com/facebookincubator/symphony/pkg/authz/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/net-auto/resourceManager/ent"
	_ "github.com/net-auto/resourceManager/ent/runtime"
	"log"
	"reflect"
	"testing"
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
