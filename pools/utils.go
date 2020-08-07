package pools

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/facebookincubator/symphony/pkg/authz"
	"github.com/facebookincubator/symphony/pkg/authz/models"
	"github.com/net-auto/resourceManager/ent/propertytype"

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

func UpdatePropertyType(
	ctx context.Context,
	client *ent.Client,
	propertyTypeID int,
	name string,
	typeName interface{},
	initValue interface{}) error {

	propertyTypeNameString := fmt.Sprintf("%v", typeName)
	propertyTypeName := propertytype.Type(propertyTypeNameString)
	if err := propertytype.TypeValidator(propertyTypeName); err != nil {
		return errors.Wrapf(err, "Unknown property type: %s", typeName)
	}

	prop := client.PropertyType.UpdateOneID(propertyTypeID).
		SetName(name).
		SetType(propertytype.TypeInt).
		SetMandatory(true)

	in := ProcessInitValue(initValue)
	//we set the property type value (we don't know what we get from the user)
	// TODO same is attempted in pools.go using switch case -> unify
	reflect.ValueOf(prop).MethodByName("Set" + strings.Title(propertyTypeNameString) + "Val").Call(in)

	_, err := prop.Save(ctx)
	return err
}

func CreatePropertyType(
	ctx context.Context,
	client *ent.Client,
	name string,
	typeName interface{},
	initValue interface{}) (*ent.PropertyType, error) {

	propertyTypeNameString := fmt.Sprintf("%v", typeName)
	propertyTypeName := propertytype.Type(propertyTypeNameString)
	if err := propertytype.TypeValidator(propertyTypeName); err != nil {
		return nil, errors.Wrapf(err, "Unknown property type: %s", typeName)
	}

	prop := client.PropertyType.Create().
		SetName(name).
		SetType(propertyTypeName).
		SetMandatory(true)

	in := ProcessInitValue(initValue)
	//we set the property type value (we don't know what we get from the user)
	// TODO same is attempted in pools.go using switch case -> unify
	reflect.ValueOf(prop).MethodByName("Set" + strings.Title(propertyTypeNameString) + "Val").Call(in)

	return prop.Save(ctx)
}

func ProcessInitValue(initValue interface{}) []reflect.Value {
	//TODO we support int, but we always get int64 instead of int
	if reflect.TypeOf(initValue).String() == "int64" {
		initValue = int(initValue.(int64))
	}

	return []reflect.Value{reflect.ValueOf(initValue)}
}

func ToRawTypes(poolValues []map[string]interface{}) []RawResourceProps {
	//TODO we support int, but we always get int64 instead of int
	for i, v := range poolValues {
		for k, val := range v {
			fmt.Printf("key[%s] value[%s]\n", k, v)
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
	resourceType *ent.ResourceType) error {
	for _, rawResourceProps := range propertyValues {
		// Parse & create the props
		var err error = nil
		var props ent.Properties
		if props, err = ParseProps(ctx, client, resourceType, rawResourceProps); err != nil {
			//TODO logging
			return errors.Wrapf(err, "Unable to create new pool \"%s\". Error parsing properties", pool.Name)
		}

		// Create pre-allocated resource
		_, err = client.Resource.Create().
			SetPool(pool).
			SetClaimed(false).
			AddProperties(props...).
			Save(ctx)

		if err != nil {
			//TODO logging
			return errors.Wrapf(err, "Unable to create new pool \"%s\". Error creating resource", pool.Name)
		}
	}

	return nil
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
