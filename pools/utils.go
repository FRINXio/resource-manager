package pools

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/net-auto/resourceManager/graph/graphql/model"
	log "github.com/net-auto/resourceManager/logging"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/propertytype"
	"github.com/net-auto/resourceManager/ent/resource"
	"github.com/pkg/errors"
)

func GetResourceFromPool(
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

func PreCreateResources(ctx context.Context, client *ent.Client, propertyValues []RawResourceProps, pool *ent.ResourcePool,
	resourceType *ent.ResourceType, claimed resource.Status, description *string, alternativeId map[string]interface{}) ([]*ent.Resource, error) {

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
			SetAlternateID(alternativeId).
			Save(ctx)
		created = append(created, resource)

		if err != nil {
			log.Error(ctx, err, "Error creating resource")
			return nil, errors.Wrapf(err, "Error creating resource")
		}
	}

	return created, nil
}

func ConvertValuesToFloat64(ctx context.Context, datamap map[string]interface{}) (map[string]interface{}, error) {
	for k, v := range datamap {
		switch t := v.(type) {
		case int:
			datamap[k] = float64(t)
		case int32:
			datamap[k] = float64(t)
		case int64:
			datamap[k] = float64(t)
		case float32:
			datamap[k] = float64(t)
		case json.Number:
			floatVal, err := v.(json.Number).Float64()
			if err != nil {
				log.Error(ctx, err, "Unable to convert a json number")
				return nil, errors.Errorf("Unable to convert a json number, error: %v", err)
			}
			datamap[k] = floatVal
		}
	}

	return datamap, nil
}

func CreateResourceType(ctx context.Context, client *ent.Client, input model.CreateResourceTypeInput) (*ent.ResourceType, error) {
	var propertyTypes []*ent.PropertyType
	for propName, rawPropType := range input.ResourceProperties {
		var propertyType, err = CreatePropertyType(ctx, client, propName, rawPropType)
		if err != nil {
			return nil, gqlerror.Errorf("Unable to create resource type: %v", err)
		}
		propertyTypes = append(propertyTypes, propertyType)
	}

	resType, err2 := client.ResourceType.Create().
		SetName(input.ResourceName).
		AddPropertyTypes(propertyTypes...).
		Save(ctx)

	if err2 != nil {
		return resType, gqlerror.Errorf("Unable to create resource type", err2)
	}

	return resType, nil
}
