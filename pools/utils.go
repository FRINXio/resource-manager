package pools

import (
	"context"
	"fmt"

	log "github.com/net-auto/resourceManager/logging"

	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/propertytype"
	"github.com/net-auto/resourceManager/ent/resource"
	"github.com/pkg/errors"
)

func GetResourceFromPool (
	ctx context.Context, obj *ent.ResourcePool) ([]*ent.Resource, error) {
	if es, err := obj.Edges.ClaimsOrErr(); !ent.IsNotLoaded(err) {
		log.Error(ctx, err, "Loading resources for resource pool %d failed", obj.ID)
		return es, err
	}
	resources, err := obj.QueryClaims().All(ctx)

	if err != nil {
		log.Error(ctx, err, "Loading resources for resource pool %d failed", obj.ID)
	}

	return resources, err
}

func CreatePropertyType(
	ctx context.Context,
	client *ent.Client,
	name string,
	typeName interface{}) (*ent.PropertyType, error) {

	propertyTypeNameString := fmt.Sprintf("%v", typeName)
	propertyTypeName := propertytype.Type(propertyTypeNameString)
	if err := propertytype.TypeValidator(propertyTypeName); err != nil {
		err := errors.Wrapf(err, "Unknown property type: %s", typeName)
		log.Error(ctx, err, "Unknown property type")
		return nil, err
	}

	return client.PropertyType.Create().
		SetName(name).
		SetType(propertyTypeName).
		SetMandatory(true).
		Save(ctx)
}

func PreCreateResources(ctx context.Context, client *ent.Client, propertyValues []RawResourceProps, pool *ent.ResourcePool, resourceType *ent.ResourceType, claimed resource.Status, description *string) ([]*ent.Resource, error) {

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
			SetNillableDescription(description).
			SetStatus(claimed).
			AddProperties(props...).
			Save(ctx)
		created = append(created, resource)

		if err != nil {
			log.Error(ctx, err, "Error creating resource")
			return nil, errors.Wrapf(err, "Error creating resource")
		}
	}

	return created, nil
}
