// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: compute select-SQL script

package helper

import (
	"errors"
	"fmt"
	"github.com/abbeymart/mctypes"
	"strings"
)

// ComputeSelectQueryAll compose select SQL script to retrieve all table-records
// The query may be limit/response may be controlled, by the user, by appending skip and limit options
func ComputeSelectQueryAll(tableName string, tableFields []string) (string, error) {
	if tableName == "" {
		return "", errors.New("table-name is required to perform the select operation")
	}
	selectQuery := ""
	if len(tableFields) > 0 {
		selectQuery = fmt.Sprintf("SELECT %v FROM %v", strings.Join(tableFields, ", "), tableName)
	} else {
		selectQuery = fmt.Sprintf("SELECT * FROM %v", tableName)
	}

	return selectQuery, nil
}

// ComputeSelectQueryById compose select SQL script by id(s)
func ComputeSelectQueryById(tableName string, recordIds []string, tableFields []string) (string, error) {
	if tableName == "" || len(recordIds) < 1 {
		return "", errors.New("table-name and record-ids are required to perform the select operation")
	}
	selectQuery := ""
	if len(tableFields) > 0 {
		// get record(s) based on projected/provided field names ([]string)
		selectQuery = fmt.Sprintf("SELECT %v", strings.Join(tableFields, ", "))
	} else {
		selectQuery = fmt.Sprintf("SELECT * FROM %v", tableName)
	}
	// from / where condition (where-in-values)
	selectQuery += fmt.Sprintf(" FROM %v WHERE id IN(%v)", tableName, strings.Join(recordIds, ", "))
	return selectQuery, nil
}

// ComputeSelectQueryByParam compose SELECT query from the where-parameters
func ComputeSelectQueryByParam(tableName string, where mctypes.WhereParamType, tableFields []string) (string, error) {
	if tableName == "" || len(where) < 1 {
		return "", errors.New("table-name and where-params are required to perform the select operation")
	}
	selectQuery := ""
	if len(tableFields) > 0 {
		// get record(s) based on projected/provided field names ([]string)
		selectQuery = fmt.Sprintf("SELECT %v", strings.Join(tableFields, ", "))
		// from
		selectQuery += fmt.Sprintf(" FROM %v", tableName)
	} else {
		selectQuery = fmt.Sprintf("SELECT * FROM %v", tableName)
	}
	// add where-params condition
	if whereScript, err := ComputeWhereQuery(where, tableFields); err == nil {
		selectQuery += whereScript
		return selectQuery, nil
	} else {
		return "", errors.New(fmt.Sprintf("error computing where-query condition(s): %v", err.Error()))
	}
}

// TODO: select-query functions for relational tables (eager & lazy queries) and data aggregation
