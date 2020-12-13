// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: compute delete-SQL script

package helper

import (
	"errors"
	"fmt"
	mccrud "github.com/abbeymart/mccrud"
	"strings"
)

func deleteScriptErr(errMsg string) (mccrud.DeleteScriptResponseType, error) {
	return mccrud.DeleteScriptResponseType{
		DeleteScript: nil,
		WhereScript:  nil,
		FieldValues:  nil,
	}, errors.New(errMsg)
}

// ComputeDeleteQueryById function computes delete SQL script by id(s)
func ComputeDeleteQueryById(tableName string, docIds []string) (string, error) {
	if tableName == "" || len(docIds) < 1 {
		return "", errors.New("table/collection name and doc-Ids are required for the delete-by-id operation")
	}
	deleteScripts := "DELETE FROM " + tableName + " WHERE id IN("
	// validated docIds, strictly contains string/UUID values, to avoid SQL-injection
	deleteIdValues := strings.Join(docIds, ", ")
	deleteScripts += deleteIdValues + " )"
	return deleteScripts, nil
}

func ComputeDeleteQueryByParam(tableName string, where mccrud.WhereParamType, tableFields []string) (string, error) {
	if tableName == "" || len(where) < 1 {
		return "", errors.New("table/collection name and where/query-condition are required for the delete-by-param operation")
	}
	if whereParam, err := ComputeWhereQuery(where, tableFields); err == nil {
		deleteScripts := fmt.Sprintf("DELETE FROM %v %v", tableName, whereParam)
		return deleteScripts, nil
	} else {
		return "", errors.New(fmt.Sprintf("error computing where-query condition(s): %v", err.Error()))
	}
}
