package pools

import (
	"context"
	"encoding/json"
	"github.com/facebook/ent/dialect/sql"
	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/propertytype"
	"github.com/net-auto/resourceManager/ent/resource"
	"strconv"

	log "github.com/net-auto/resourceManager/logging"
	"time"
)

func ResourcePropertyPropertyTypeJoins(ctx context.Context, poolId int, tx *ent.Tx) ([]*ent.Resource, error) {
	var solvedResources []*ent.Resource
	query := "SELECT DISTINCT resources.id, resources.status, resources.alternate_id, resources.updated_at, " +
		"properties.id, properties.int_val, properties.bool_val, properties.float_val, properties.string_val, " +
		"property_types.id, property_types.type, property_types.name FROM resources " +
		"INNER JOIN properties ON properties.resource_properties = resources.id " +
		"INNER JOIN property_types ON properties.property_type = property_types.id " +
		"WHERE resources.resource_pool_claims = " + strconv.Itoa(poolId) + ";"

	var resDTO resourceDTO
	var propDTO propertyDTO
	var propTypeDTO propertyTypeDTO

	rows := &sql.Rows{}
	var args []interface{}
	err := tx.UnderlyingTx().Query(ctx, query, args, rows)
	if err != nil {
		log.Error(ctx, err, "Error while executing query: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&resDTO.id, &resDTO.status, &resDTO.alternativeID, &resDTO.updatedAt, &propDTO.id,
			&propDTO.intVal, &propDTO.boolVal, &propDTO.floatVal, &propDTO.stringVal, &propTypeDTO.id,
			&propTypeDTO.types, &propTypeDTO.name)
		if err != nil {
			log.Error(ctx, err, "Error while scanning results: %v", err)
		}

		currentResource := assignResourceValues(resDTO)
		currentProperty := assignPropertyValues(propDTO)
		currentPropertyType := assignPropertyTypeValues(propTypeDTO)

		var newResource = true
		currentProperty.Edges.Type = currentPropertyType
		for index, res := range solvedResources {
			if res.ID == currentResource.ID {
				newResource = false
				res.Edges.Properties = append(res.Edges.Properties, currentProperty)
				solvedResources[index] = res
				break
			}
		}
		if newResource {
			currentResource.Edges.Properties = append(currentResource.Edges.Properties, currentProperty)
			solvedResources = append(solvedResources, currentResource)
		}
	}

	if err != nil {
		log.Error(ctx, err, "Failed to execute query: %v", err)
	}
	return solvedResources, err
}

// assignResourceValues assigns the values that were returned from sql.Rows (after scanning)
// to the Resource fields.
func assignResourceValues(dto resourceDTO) *ent.Resource {
	var r ent.Resource

	r.ID = dto.id
	if dto.status != nil {
		r.Status = resource.Status(*dto.status)
	}
	if dto.alternativeID != nil {
		json.Unmarshal(dto.alternativeID, &r.AlternateID)
	}
	if dto.updatedAt != nil {
		r.UpdatedAt = *dto.updatedAt
	}
	r.Description = new(string)
	return &r
}

// assignPropertyValues assigns the values that were returned from sql.Rows (after scanning)
// to the Property fields.
func assignPropertyValues(dto propertyDTO) *ent.Property {
	var pr ent.Property

	pr.ID = dto.id
	if dto.intVal.Valid {
		val := int(dto.intVal.Int64)
		pr.IntVal = &val
	}
	if dto.boolVal.Valid {
		pr.BoolVal = &dto.boolVal.Bool
	}
	if dto.floatVal.Valid {
		pr.FloatVal = &dto.floatVal.Float64
	}
	if dto.floatVal.Valid {
		pr.FloatVal = &dto.floatVal.Float64
	}
	if dto.stringVal.Valid {
		pr.StringVal = &dto.stringVal.String
	}
	return &pr
}

// assignPropertyTypeValues assigns the values that were returned from sql.Rows (after scanning)
// to the PropertyType fields.
func assignPropertyTypeValues(dto propertyTypeDTO) *ent.PropertyType {
	var pt ent.PropertyType

	pt.ID = dto.id
	if dto.types.Valid {
		pt.Type = propertytype.Type(dto.types.String)
	}
	if dto.name.Valid {
		pt.Name = dto.name.String
	}
	return &pt
}

type resourceDTO struct {
	id            int
	status        *string
	alternativeID []byte
	updatedAt     *time.Time
}

type propertyDTO struct {
	id        int
	intVal    sql.NullInt64
	boolVal   sql.NullBool
	floatVal  sql.NullFloat64
	stringVal sql.NullString
}

type propertyTypeDTO struct {
	id    int
	types sql.NullString
	name  sql.NullString
}
