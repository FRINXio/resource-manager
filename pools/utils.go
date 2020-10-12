package pools

import (
	"context"
	"fmt"

	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/propertytype"
	"github.com/net-auto/resourceManager/ent/resource"
	"github.com/pkg/errors"
)

func HasAllocatedResources(ctx context.Context, client *ent.Client, poolId int) (bool, error) {
	pool, errRp := client.ResourcePool.Get(ctx, poolId)

	if errRp != nil {
		return false, errRp
	}

	if pool == nil {
		return false, errors.New("Unable to find pool")
	}

	resources, err := GetResourceFromPool(ctx, pool)

	if err != nil {
		return false, err
	}

	for _, r := range resources {
		if r.Status == resource.StatusClaimed || r.Status == resource.StatusBench {
			return true, nil
		}
	}

	return false, nil
}

func GetResourceFromPool (
	ctx context.Context, obj *ent.ResourcePool) ([]*ent.Resource, error) {
	if es, err := obj.Edges.ClaimsOrErr(); !ent.IsNotLoaded(err) {
		return es, err
	}
	return obj.QueryClaims().All(ctx)
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
			return nil, errors.Wrapf(err, "Error parsing properties")
		}

		// FIXME: fail when this resource is already in DB (same logic as in FreeResource)
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
