package src

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/facebook/ent/dialect/sql"
	"github.com/net-auto/resourceManager/ent"
	log "github.com/net-auto/resourceManager/logging"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type UniqueId struct {
	ctx                    context.Context
	resourcePoolID         int
	resourcePoolProperties map[string]interface{}
	userInput              map[string]interface{}
}

func NewUniqueId(ctx context.Context,
	resourcePoolID int,
	resourcePoolProperties map[string]interface{},
	userInput map[string]interface{}) UniqueId {
	return UniqueId{ctx, resourcePoolID, resourcePoolProperties, userInput}
}

func (uniqueId *UniqueId) getNextFreeCounter(poolId int, fromValue int, toValue int, desiredValue int,
	ctx context.Context) (int, error) {

	transaction := ctx.Value(ent.TxCtxKey{})
	if transaction == nil {
		log.Error(ctx, nil, "Unable retrieve already opened transaction for pool with ID: %d", poolId)
		return -1, errors.Wrapf(nil, "Unable retrieve already opened transaction for pool with ID: %d", poolId)
	}
	tx := transaction.(*ent.Tx)
	if desiredValue >= 0 {
		query := "WITH RECURSIVE t(n) AS (VALUES (" + strconv.Itoa(fromValue) + ") " +
			"UNION ALL SELECT n+1 FROM t WHERE n < " + strconv.Itoa(toValue) + ") " +
			"SELECT n FROM t LEFT OUTER JOIN ( " +
			"SELECT properties.int_val FROM properties JOIN resources ON properties.resource_properties = resources.id " +
			"WHERE resources.resource_pool_claims = " + strconv.Itoa(poolId) + ") " +
			"AS pr ON n = pr.int_val WHERE pr.int_val IS null AND n = " + strconv.Itoa(desiredValue) + ";"
		valueExist, value, err := selectValueFromDB(ctx, tx, query)
		if err != nil {
			return 0, err
		}
		if valueExist == true {
			return int(value), nil
		}
		return 0, errors.New("Unique-id " + strconv.Itoa(desiredValue) + " was already claimed.")
	} else {
		query := "WITH RECURSIVE t(n) AS (VALUES (" + strconv.Itoa(fromValue) + ") " +
			"UNION ALL SELECT n+1 FROM t WHERE n < " + strconv.Itoa(toValue) + ") " +
			"SELECT n FROM t LEFT OUTER JOIN ( " +
			"SELECT properties.int_val FROM properties JOIN resources ON properties.resource_properties = resources.id " +
			"WHERE resources.resource_pool_claims = " + strconv.Itoa(poolId) + ") AS pr " +
			"ON n = pr.int_val WHERE pr.int_val IS null ORDER BY n ASC LIMIT 1;"
		valueExist, value, err := selectValueFromDB(ctx, tx, query)
		if err != nil {
			return 0, err
		}
		if valueExist == true {
			return int(value), nil
		}
		return 0, errors.New("Unique-id pool " + strconv.Itoa(poolId) + " is full.")
	}
}

func (uniqueId *UniqueId) Invoke() (map[string]interface{}, error) {
	if uniqueId.resourcePoolProperties == nil {
		return nil, errors.New("Unable to extract resources")
	}
	var fromValue = 0 // max int
	resourcePoolfromValue, ok := uniqueId.resourcePoolProperties["from"]
	if ok {
		resourcePoolfromValue, _ := NumberToInt(resourcePoolfromValue)
		fromValue = resourcePoolfromValue.(int)
	}

	idFormat, ok := uniqueId.resourcePoolProperties["idFormat"]
	if !ok {
		return nil, errors.New("Missing idFormat in resources")
	}
	if !strings.Contains(idFormat.(string), "{counter}") {
		return nil, errors.New("Missing {counter} in idFormat")
	}

	var toValue = ^uint(0) >> 1 // max int
	resourcePooltoValue, ok := uniqueId.resourcePoolProperties["to"]
	if ok {
		resourcePooltoValue, _ := NumberToInt(resourcePooltoValue)
		toValue = uint(resourcePooltoValue.(int))
	}

	replacePoolProperties := make(map[string]interface{})
	for k, v := range uniqueId.resourcePoolProperties {
		if k != "idFormat" && k != "counterFormatWidth" && k != "from" && k != "to" {
			replacePoolProperties[k] = v
		}
	}
	desiredValue := -1
	if value, ok := uniqueId.userInput["desiredValue"]; ok {
		value, _ = NumberToInt(value)
		desiredValue = value.(int)
		if desiredValue > int(int64(toValue)) {
			return nil, errors.New("Unable to allocate Unique-id desiredValue: " + strconv.FormatInt(int64(value.(int)), 10) + "." +
				" Value is out of scope: " + strconv.FormatInt(int64(toValue), 10))
		}
		if desiredValue < fromValue {
			return nil, errors.New("Unable to allocate Unique-id desiredValue: " + strconv.FormatInt(int64(value.(int)), 10) + "." +
				" Value is out of scope: " + strconv.FormatInt(int64(fromValue), 10))
		}
	}

	nextFreeCounter, err := uniqueId.getNextFreeCounter(uniqueId.resourcePoolID, fromValue, resourcePooltoValue.(int),
		desiredValue, uniqueId.ctx)
	if err != nil {
		return nil, err
	}
	if prefixNumber, ok := uniqueId.resourcePoolProperties["counterFormatWidth"]; ok {
		replacePoolProperties["counter"] = fmt.Sprintf(
			"%0"+strconv.Itoa(prefixNumber.(int))+"d", int(nextFreeCounter))
	} else {
		replacePoolProperties["counter"] = nextFreeCounter
	}

	for k, v := range replacePoolProperties {
		switch v.(type) {
		case float64:
			v = fmt.Sprint(v.(float64))
		case int:
			v = strconv.Itoa(v.(int))
		case json.Number:
			intVal64, err := v.(json.Number).Float64()
			if err != nil {
				return nil, errors.New("Unable to convert a json number")
			}
			v = fmt.Sprint(intVal64)
		}

		idFormat = strings.Replace(idFormat.(string), "{"+k+"}", v.(string), 1)
	}

	var result = make(map[string]interface{})
	result["text"] = idFormat
	result["counter"] = nextFreeCounter
	return result, nil
}

func (uniqueId *UniqueId) Capacity() (map[string]interface{}, error) {
	ctx := uniqueId.ctx

	transaction := ctx.Value(ent.TxCtxKey{})
	if transaction == nil {
		log.Error(ctx, nil, "Unable retrieve already opened transaction for pool with ID: %d", 1)
		return nil, errors.Wrapf(nil, "Unable retrieve already opened transaction for pool with ID: %d", 1)
	}
	tx := transaction.(*ent.Tx)

	var result = make(map[string]interface{})
	var fromValue float64
	var toValue float64
	to, ok := uniqueId.resourcePoolProperties["to"]
	if ok {
		toValue = float64(to.(int))
	} else {
		toValue = float64(^uint(0) >> 1)
	}
	from, ok := uniqueId.resourcePoolProperties["from"]
	if ok {
		fromValue = float64(from.(int))
	} else {
		fromValue = float64(0)
	}
	query := "WITH RECURSIVE t(n) AS (VALUES (" + strconv.Itoa(from.(int)) + ") " +
		"UNION ALL SELECT n+1 FROM t WHERE n < " + strconv.Itoa(to.(int)) + ") " +
		"SELECT COUNT(n) FROM t LEFT OUTER JOIN ( " +
		"SELECT properties.int_val FROM properties JOIN resources ON properties.resource_properties = resources.id " +
		"WHERE resources.resource_pool_claims = " + strconv.Itoa(uniqueId.resourcePoolID) + ") AS pr " +
		"ON n = pr.int_val WHERE pr.int_val IS null;"
	valueExist, value, err := selectValueFromDB(ctx, tx, query)
	if err != nil {
		return nil, err
	}
	if valueExist == true {
		result["freeCapacity"] = float64(value)
		result["utilizedCapacity"] = toValue - float64(value) - fromValue + 1
	}
	return result, nil
}

func selectValueFromDB(ctx context.Context, tx *ent.Tx, query string) (valueExist bool, resultValue int64, error error) {
	rows := &sql.Rows{}
	var args []interface{}
	err := tx.UnderlyingTx().Query(ctx, query, args, rows)
	if err != nil {
		log.Error(ctx, err, "Error while executing query: %v", err)
		return false, 0, err
	}
	defer rows.Close()

	type value struct {
		intValue sql.NullInt64
	}
	var result value
	for rows.Next() {
		err = rows.Scan(&result.intValue)
		if err != nil {
			log.Error(ctx, err, "Error while scanning results: %v", err)
			return false, 0, err
		}
	}

	if err != nil {
		log.Error(ctx, err, "Failed to execute query: %v", err)
		return false, 0, err
	}
	return result.intValue.Valid, result.intValue.Int64, err
}
