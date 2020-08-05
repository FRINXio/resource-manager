package pools

import (
	"context"
	"fmt"
	"github.com/facebookincubator/symphony/pkg/authz"
	"github.com/facebookincubator/symphony/pkg/authz/models"
	"github.com/net-auto/resourceManager/ent/propertytype"
	"log"
	"reflect"

	"github.com/net-auto/resourceManager/ent"
	"github.com/pkg/errors"
)

func GetContext() context.Context {
	ctx := context.Background()
	ctx = authz.NewContext(ctx, &models.PermissionSettings{
		CanWrite:        true,
		WorkforcePolicy: authz.NewWorkforcePolicy(true, true)})
	return ctx
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

func HasPropertyTypeExistingProperties(
	ctx context.Context,
	client *ent.Client,
	propertyTypeID int) bool {

	propertyType, err := client.PropertyType.Get(ctx, propertyTypeID)

	if err != nil {
		//TODO error handling
		return true //not sure
	}
	exists, err2 := client.PropertyType.QueryProperties(propertyType).Exist(ctx)

	if err2 != nil {
		//TODO error handling
		return true //not sure
	}
	//TODO also check resource type association??

	return exists
}

func CreatePropertyType(
	ctx context.Context,
	client *ent.Client,
	name string,
	typeName interface{}) (*ent.PropertyType, error) {

	propertyTypeNameString := fmt.Sprintf("%v", typeName)
	propertyTypeName := propertytype.Type(propertyTypeNameString)
	if err := propertytype.TypeValidator(propertyTypeName); err != nil {
		return nil, errors.Wrapf(err, "Unknown property type: %s", typeName)
	}

	return client.PropertyType.Create().
		SetName(name).
		SetType(propertyTypeName).
		SetMandatory(true).
		Save(ctx)
}

func ToRawTypes(poolValues []map[string]interface{}) []RawResourceProps {
	//TODO we support int, but we always get int64 instead of int
	for i, v := range poolValues {
		for k, val := range v {
			if reflect.TypeOf(val).String() == "int64" {
				poolValues[i][k] = int(val.(int64))
			}
		}
	}

	var rawProps []RawResourceProps

	for _, v := range poolValues {
		rawProps = append(rawProps, v)
	}

	return rawProps
}

func PreCreateResources(ctx context.Context,
	client *ent.Client,
	propertyValues []RawResourceProps,
	pool *ent.ResourcePool,
	resourceType *ent.ResourceType,
	claimed bool) ([]*ent.Resource, error) {

	var created []*ent.Resource
	for _, rawResourceProps := range propertyValues {
		// Parse & create the props
		var err error = nil
		var props ent.Properties
		if props, err = ParseProps(ctx, client, resourceType, rawResourceProps); err != nil {
			//TODO logging
			return nil, errors.Wrapf(err, "Error parsing properties")
		}

		// Create pre-allocated resource
		var resource *ent.Resource
		resource, err = client.Resource.Create().
			SetPool(pool).
			SetClaimed(claimed).
			AddProperties(props...).
			Save(ctx)
		created = append(created, resource)

		if err != nil {
			//TODO logging
			return nil, errors.Wrapf(err, "Error creating resource")
		}
	}

	return created, nil
}

func CheckIfPoolsExist(
	ctx context.Context,
	client *ent.Client,
	resourceTypeID int) (bool, *ent.ResourceType) {
	resourceType, err := client.ResourceType.Get(ctx, resourceTypeID)
	if err != nil {
		//TODO add annoying GO error handling
		return true, resourceType //fix we don't know
	}

	//there can't be any existing pools
	count, err2 := resourceType.QueryPools().Count(ctx)

	if err2 != nil || count > 0 {
		//TODO add annoying GO error handling
		return true, resourceType //fix we don't know
	}

	return false, resourceType
}

func PropertiesToMap(props []*ent.Property) (map[string]interface{}, error) {
	var asMap = make(map[string]interface{})

	for _, prop := range props {
		// TODO is there a better way of parsing individual types ? Reuse something from inv ?
		// TODO add additional types
		// TODO we have this switch in 2 places
		switch prop.Edges.Type.Type {
		case "int":
			asMap[prop.Edges.Type.Name] = *prop.IntVal
		case "string":
			asMap[prop.Edges.Type.Name] = *prop.StringVal
		default:
			return nil, fmt.Errorf("Unsupported property type \"%s\"", prop.Edges.Type.Type)
		}
	}

	return asMap, nil
}

