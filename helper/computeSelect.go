// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: compute select-SQL script

package helper

import (
	"errors"
	mccrud "github.com/abbeymart/mccrudgo"
)

func ComputeSelectQueryById(tableName string, docIds []string, fields []string) (string, error) {

	return "", errors.New("development")
}

func ComputeSelectQuery(tableName string, where mccrud.WhereParamType, fields []string, queryType mccrud.TaskType) (string, error) {

	return "", errors.New("development")
}
