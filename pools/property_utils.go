package pools

import (
	"context"
	"fmt"
	"strconv"
	"encoding/json"

	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/predicate"
	"github.com/net-auto/resourceManager/ent/property"
	"github.com/net-auto/resourceManager/ent/propertytype"
	"github.com/pkg/errors"
)

func GetValue(prop *ent.Property) (interface{}, error) {
	// TODO is there a better way of parsing individual types ? Reuse something from inv ?
	// TODO add additional types
	// TODO we have this switch in 2 places
	switch prop.Edges.Type.Type {
	case "int":
		if prop.IntVal == nil {
			return nil, nil
		}
		return *prop.IntVal, nil
	case "string":
		if prop.StringVal == nil {
			return nil, nil
		}
		return *prop.StringVal, nil
	case "float":
		if prop.FloatVal == nil {
			return nil, nil
		}
		return *prop.FloatVal, nil
	case "bool":
		if prop.BoolVal == nil {
			return nil, nil
		}
		return *prop.BoolVal, nil
	default:
		return nil, fmt.Errorf("Unsupported property type \"%s\"", prop.Edges.Type.Type)
	}
}

func PropertiesToMap(props []*ent.Property) (RawResourceProps, error) {
	var asMap = make(map[string]interface{})

	for _, prop := range props {
		value, err := GetValue(prop)
		if err != nil {
			return nil, err
		}
		if value != nil {
			asMap[prop.Edges.Type.Name] = value
		}
	}
	return asMap, nil
}

func CompareProps(
	ctx context.Context,
	resourceType *ent.ResourceType,
	propertyValues RawResourceProps) ([]predicate.Property, error) {

	var predicates []predicate.Property
	for pN, pV := range propertyValues {
		// FIXME: N+1 selects problem
		pT, err := resourceType.QueryPropertyTypes().Where(propertytype.NameEQ(pN)).Only(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "Unknown property: \"%s\" for resource type: \"%s\"", pN, resourceType)
		}
		// TODO nil handling

		propPredict := property.HasTypeWith(propertytype.ID(pT.ID))

		// TODO is there a better way of parsing individual types ? Reuse something from inv ?
		// TODO add additional types
		// TODO we have this switch in 2 places
		switch pT.Type {
		case "int":
			var intVal int
			switch t := pV.(type) {
			case int:
				intVal = t
			case int32:
				intVal = int(t)
			case int64:
				intVal = int(t)
			case float32:
				intVal = int(t)
			case float64:
				intVal = int(t)
			case json.Number:
				intVal64, err := pV.(json.Number).Int64()
				if err != nil {
					return nil, errors.Errorf("Unable to convert a json number, error: %v", err)
				}
				intVal = int(intVal64)
			default:
				return nil, errors.Errorf("Unsupported int conversion from %T", t)
			}
			propPredict = property.And(propPredict, property.IntValEQ(intVal))
		case "string":
			propPredict = property.And(propPredict, property.StringValEQ(pV.(string)))
		case "float":
			propPredict = property.And(propPredict, property.FloatValEQ(pV.(float64)))
		case "bool":
			propPredict = property.And(propPredict, property.BoolValEQ(pV.(bool)))
		default:
			return nil, errors.Errorf("Unsupported property type \"%s\"", pT.Type)
		}

		predicates = append(predicates, propPredict)
	}

	return predicates, nil
}

// ParseProps turns a map such as ["a": 3, "b": "value"] into a list of properties and stores them in DB
//  uses resource type to find out what are the predefined types for each value
func ParseProps(
	ctx context.Context,
	tx *ent.Client,
	resourceType *ent.ResourceType,
	propertyValues RawResourceProps) (ent.Properties, error) {

	var props ent.Properties
	propTypes, err := resourceType.QueryPropertyTypes().All(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to determine property types for \"%s\"", resourceType)
	}

	for _, pt := range propTypes {
		pv := propertyValues[pt.Name]

		if pt.Mandatory {
			if pv == nil {
				return nil, errors.Errorf("Missing mandatory property \"%s\"", pt.Name)
			}
		} else {
			if pv == nil {
				continue
			}
		}

		ppBuilder := tx.Property.Create().SetType(pt)

		// TODO is there a better way of parsing individual types ? Reuse something from inv ?
		// TODO add additional types
		switch pt.Type {
		case "int":
			// Parse the int from string to be sure
			atoi, err := strconv.Atoi(fmt.Sprintf("%v", pv))
			if err != nil {
				return nil, errors.Wrapf(err, "Unable to parse int value from \"%s\"", pv)
			}
			ppBuilder.SetIntVal(atoi)
		case "string":
			ppBuilder.SetStringVal(pv.(string))
		case "float":
			// Parse the float from string to be sure
			parsedFloat, err := strconv.ParseFloat(fmt.Sprintf("%v", pv), 64)
			if err != nil {
				return nil, errors.Wrapf(err, "Unable to parse float value from \"%s\"", pv)
			}
			ppBuilder.SetFloatVal(parsedFloat)
		case "bool":
			parsedBool, err := strconv.ParseBool(fmt.Sprintf("%v", pv))
			if err != nil {
				return nil, errors.Wrapf(err, "Unable to parse bool value from \"%s\"", pv)
			}
			ppBuilder.SetBoolVal(parsedBool)
		default:
			return nil, errors.Errorf("Unsupported property type \"%s\"", pt.Type)
		}

		pp, err := ppBuilder.Save(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to instantiate property of type \"%s\"", pt.Type)
		}
		props = append(props, pp)
	}

	return props, nil
}

// ToRawTypes converts between []map[string]interface{} and []RawResourceProps
//  which is the same thing ... but not to the compiler
func ToRawTypes(poolValues []map[string]interface{}) []RawResourceProps {
	var rawProps []RawResourceProps
	for _, v := range poolValues {
		rawProps = append(rawProps, v)
	}
	return rawProps
}
