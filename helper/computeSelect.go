// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: compute select-SQL script

package helper

import (
	"errors"
	"fmt"
	"github.com/abbeymart/mccrud"
	"strings"
)

func selectErrMessage(errMsg string) (mccrud.SelectQueryObject, error) {
	return mccrud.SelectQueryObject{
		SelectQuery: "",
		FieldValues: nil,
		WhereQuery:  mccrud.WhereQueryObject{},
	}, errors.New(errMsg)
}

// ComputeSelectQueryAll compose select SQL script to retrieve all table-records.
// The query may be constraint by skip(offset) and limit options
func ComputeSelectQueryAll(modelRef interface{}, tableName string) (string, error) {
	if tableName == "" || modelRef == nil {
		return "", errors.New("model(struct) and table-name are required to perform the select operation")
	}
	// compute map[string]interface (underscore_fields) from the modelRef (struct)
	mapMod, mapErr := mccrud.StructToMapUnderscore(modelRef)
	if mapErr != nil {
		return "", mapErr
	}
	// compute table-fields
	var tableFields []string
	for fieldName := range mapMod {
		tableFields = append(tableFields, fieldName)
	}
	// get records for the model-defined fields/columns
	selectQuery := fmt.Sprintf("SELECT %v FROM %v ", strings.Join(tableFields, ", "), tableName)
	return selectQuery, nil
}

// ComputeSelectQueryById compose select SQL script by id
func ComputeSelectQueryById(modelRef interface{}, tableName string, recordId string) (string, error) {
	if tableName == "" || recordId == "" || modelRef == nil {
		return "", errors.New("model (struct), table-name and record-id are required to perform the select operation")
	}
	// TODO: compute map[string]interface (underscore_fields) from the modelRef (struct)
	mapMod, mapErr := mccrud.StructToMapUnderscore(modelRef)
	if mapErr != nil {
		return "", mapErr
	}
	var tableFields []string
	for fieldName := range mapMod {
		tableFields = append(tableFields, fieldName)
	}
	// get record(s) based on projected/provided field names ([]string)
	selectQuery := fmt.Sprintf("SELECT %v FROM %v ", strings.Join(tableFields, ", "), tableName)
	// from / where condition (where-in-values)
	selectQuery += fmt.Sprintf("WHERE id = %v", recordId)
	return selectQuery, nil
}

// ComputeSelectQueryByIds compose select SQL script by ids
func ComputeSelectQueryByIds(modelRef interface{}, tableName string, recordIds []string) (string, error) {
	if tableName == "" || len(recordIds) < 1 || modelRef == nil {
		return "", errors.New("model (struct), table-name and record-ids are required to perform the select operation")
	}
	// TODO: compute map[string]interface (underscore_fields) from the modelRef (struct)
	mapMod, mapErr := mccrud.StructToMapUnderscore(modelRef)
	if mapErr != nil {
		return "", mapErr
	}
	// compute select-fields
	var tableFields []string
	for fieldName := range mapMod {
		tableFields = append(tableFields, fieldName)
	}
	// get record(s) based on projected/provided field names ([]string)
	selectQuery := fmt.Sprintf("SELECT %v FROM %v ", strings.Join(tableFields, ", "), tableName)
	// from / where condition (where-in-values)
	whereIds := ""
	idLen := len(recordIds)
	for idCount, id := range recordIds {
		whereIds += "'" + id + "'"
		if idLen > 1 && idCount < idLen-1 {
			whereIds += ", "
		}
	}
	selectQuery += fmt.Sprintf("WHERE id IN( %v )", whereIds)
	return selectQuery, nil
}

// ComputeSelectQueryByParam compose SELECT query from the where-parameters
func ComputeSelectQueryByParam(modelRef interface{}, tableName string, queryParam mccrud.QueryParamType) (mccrud.SelectQueryObject, error) {
	if tableName == "" || len(queryParam) < 1 || modelRef == nil {
		return selectErrMessage("model (struct), table-name, and queryParam are required to perform the select operation")
	}
	// TODO: compute map[string]interface (underscore_fields) from the modelRef (struct)
	mapMod, mapErr := mccrud.StructToMapUnderscore(modelRef)
	if mapErr != nil {
		return selectErrMessage(fmt.Sprintf("%v", mapErr.Error()))
	}
	// compute select-fields
	var tableFields []string
	for fieldName := range mapMod {
		tableFields = append(tableFields, fieldName)
	}
	// get record(s) based on projected/provided field names ([]string)
	selectQuery := fmt.Sprintf("SELECT %v FROM %v ", strings.Join(tableFields, ", "), tableName)
	// add queryParam-params condition
	whereQuery, err := ComputeWhereQuery(queryParam, 1)
	if err == nil {
		return mccrud.SelectQueryObject{
			SelectQuery: selectQuery,
			FieldValues: nil,
			WhereQuery:  whereQuery,
		}, nil
	} else {
		return selectErrMessage(fmt.Sprintf("error computing queryParam-query condition(s): %v", err.Error()))
	}
}

// TODO: select-query functions for relational tables (eager & lazy queries) and data aggregation
