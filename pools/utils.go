package pools

import (
	"context"
	"fmt"
	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/propertytype"
	"github.com/net-auto/resourceManager/ent/resource"
	"github.com/pkg/errors"
)

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

func PreCreateResources(ctx context.Context,
	client *ent.Client,
	propertyValues []RawResourceProps,
	pool *ent.ResourcePool,
	resourceType *ent.ResourceType,
	claimed resource.Status) ([]*ent.Resource, error) {

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
			SetStatus(claimed).
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

