// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: compute select-SQL script

package helper

import (
	"errors"
	"github.com/abbeymart/mccrud"
)

func ComputeSelectQueryAll(tableName string, tableFields []string, recordIds []string, fields []string) (string, error) {

	return "", errors.New("development")
}

func ComputeSelectQueryById(tableName string, tableFields []string, recordIds []string, fields []string) (string, error) {

	return "", errors.New("development")
}

func ComputeSelectQueryByParam(tableName string, where mccrud.WhereParamType, fields []string, queryType mccrud.TaskType) (string, error) {

	return "", errors.New("development")
}
