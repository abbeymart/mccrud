// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: compute delete-SQL scripts

package mccrud

import (
	"errors"
	"fmt"
)

func deleteErrMessage(errMsg string) (DeleteQueryObject, error) {
	return DeleteQueryObject{
		DeleteQuery: "",
		FieldValues: nil,
		WhereQuery:  WhereQueryObject{},
	}, errors.New(errMsg)
}

// ComputeDeleteQueryById function computes delete SQL script by id(s)
func ComputeDeleteQueryById(tableName string, recordId string) (string, error) {
	if tableName == "" || recordId == "" {
		return "", errors.New("table/collection name and record-id are required for the delete-by-id operation")
	}
	// validated recordIds, strictly contains string/UUID values, to avoid SQL-injection
	deleteQuery := "DELETE FROM " + tableName + " WHERE id =" + recordId
	return deleteQuery, nil
}

// ComputeDeleteQueryByIds function computes delete SQL script by id(s)
func ComputeDeleteQueryByIds(tableName string, recordIds []string) (string, error) {
	if tableName == "" || len(recordIds) < 1 {
		return "", errors.New("table/collection name and record-Ids are required for the delete-by-id operation")
	}
	// validated recordIds, strictly contains string/UUID values, to avoid SQL-injection
	// from / where condition (where-in-values)
	whereIds := ""
	idLen := len(recordIds)
	for idCount, id := range recordIds {
		whereIds += "'" + id + "'"
		if idLen > 1 && idCount < idLen-1 {
			whereIds += ", "
		}
	}
	deleteQuery := "DELETE FROM " + tableName + " WHERE id IN(" + whereIds + ")"
	return deleteQuery, nil
}

// ComputeDeleteQueryByParam function computes delete SQL script by parameter specifications
func ComputeDeleteQueryByParam(tableName string, where QueryParamType) (DeleteQueryObject, error) {
	if tableName == "" || len(where) < 1 {
		return deleteErrMessage("table/collection name and where/query-condition are required for the delete-by-param operation")
	}
	whereParam, err := ComputeWhereQuery(where, 1)
	if err == nil {
		//deleteScript := fmt.Sprintf("DELETE FROM %v %v", tableName, whereParam)
		return DeleteQueryObject{
			DeleteQuery: fmt.Sprintf("DELETE FROM %v ", tableName),
			WhereQuery: whereParam,
			FieldValues: nil,
		}, nil
	} else {
		return deleteErrMessage(fmt.Sprintf("error computing where-query condition(s): %v", err.Error()))
	}
}