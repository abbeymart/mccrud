// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: compute delete-SQL script

package helper

import (
	"errors"
	mccrud "github.com/abbeymart/mccrudgo"
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
	docIdsLen := len(docIds)
	if tableName == "" || docIdsLen < 1 {
		return "", errors.New("table/collection name and doc-Ids are required for the delete-by-id operation")
	}
	deleteScripts := "DELETE FROM " + tableName + " WHERE id IN("
	// TODO: sanitize docIds to ensure it's strictly contains UUID values
	deleteIdValues := strings.Join(docIds, ", ")  // perform at the DbQuery task

	// for ind := range docIds {
	//	if docIdsLen > 1 && ind < docIdsLen-1 {
	//		deleteScripts += ", $" + fmt.Sprintf("%v", ind+1)
	//	} else {
	//		deleteScripts += " $" + fmt.Sprintf("%v", ind+1)
	//	}
	//}

	deleteScripts += deleteIdValues + " )"

	return deleteScripts, nil
}

func ComputeDeleteQueryByParam(tableName string, where mccrud.WhereParamType) (mccrud.DeleteScriptResponseType, error) {
	if tableName == "" || len(where) < 1 {
		return deleteScriptErr("table/collection name and where/query-condition are required for the delete-by-param operation")
	}
	whereParam, _ := ComputeWhereQuery(where)
	deleteScripts := "DELETE FROM " + tableName + " " + whereParam.WhereScript

	return mccrud.DeleteScriptResponseType{
		DeleteScript: deleteScripts,
		WhereScript:  whereParam.WhereScript,
		FieldValues:  whereParam.FieldValues,
	}, nil
}
