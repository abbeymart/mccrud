// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: compute select-SQL script

package helper

import (
	"errors"
	mccrud "github.com/abbeymart/mccrudgo"
)

func ComputeSelectQueryById(tableName string, docIds []string, fields []string) (mccrud.SelectScriptResponseType, error) {

	return mccrud.SelectScriptResponseType{}, errors.New("development")
}

func ComputeSelectQuery(tableName string, where mccrud.WhereParamType, fields []string, queryType mccrud.TaskType) (mccrud.SelectScriptResponseType, error) {

	return mccrud.SelectScriptResponseType{}, errors.New("development")
}
