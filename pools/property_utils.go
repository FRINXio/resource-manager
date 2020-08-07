package pools

import (
	"context"
	"fmt"
	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/predicate"
	"github.com/net-auto/resourceManager/ent/property"
	"github.com/net-auto/resourceManager/ent/propertytype"
	"github.com/pkg/errors"
	"strconv"
)

func PropertiesToMap(props []*ent.Property) (RawResourceProps, error) {
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
		case "float":
			asMap[prop.Edges.Type.Name] = *prop.FloatVal
		case "bool":
			asMap[prop.Edges.Type.Name] = *prop.BoolVal
		default:
			return nil, fmt.Errorf("Unsupported property type \"%s\"", prop.Edges.Type.Type)
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
		pT, err := resourceType.QueryPropertyTypes().Where(propertytype.NameEQ(pN)).Only(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "Unknown property: \"%s\" for resource type: \"%s\"", pN, resourceType)
		}

		propPredict := property.HasTypeWith(propertytype.ID(pT.ID))

		// TODO is there a better way of parsing individual types ? Reuse something from inv ?
		// TODO add additional types
		// TODO we have this switch in 2 places
		switch pT.Type {
		case "int":
			propPredict = property.And(propPredict, property.IntValEQ(pV.(int)))
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
