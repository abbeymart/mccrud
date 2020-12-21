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

// ComputeSelectQueryAll compose select SQL script to retrieve all table-records
// The query may be limit/response may be controlled, by the user, by appending skip and limit options
func ComputeSelectQueryAll(tableName string, tableFields []string) (string, error) {
	if tableName == "" {
		return "", errors.New("table-name is required to perform the select operation")
	}
	selectQuery := ""
	fieldLen := len(tableFields)
	switch fieldLen {
	case 0:
		selectQuery = "SELECT * FROM " + tableName
	default:
		selectQuery = "SELECT " + strings.Join(tableFields, ", ")
		selectQuery += fmt.Sprintf(" FROM %v ", tableName)
	}
	return selectQuery, nil
}

// ComputeSelectQueryById compose select SQL script by id(s)
func ComputeSelectQueryById(tableName string, recordIds []string, tableFields []string) (string, error) {
	if tableName == "" || len(recordIds) < 1 {
		return "", errors.New("table-name and record-ids are required to perform the select operation")
	}
	selectQuery := ""
	fieldLen := len(tableFields)
	switch fieldLen {
	case 0:
		selectQuery = "SELECT * FROM " + tableName + " WHERE id IN (" + strings.Join(recordIds, ", ") + ")"
	default:
		// get record(s) based on projected/provided field names ([]string)
		// select field/column-names
		selectQuery = "SELECT " + strings.Join(tableFields, ", ")
		// from / where condition
		selectQuery += fmt.Sprintf(" FROM %v WHERE id IN(", tableName)
		// where-in-values
		selectQuery += strings.Join(recordIds, ", ") + ")"
	}
	return selectQuery, nil
}

// ComputeSelectQueryByParam compose SELECT query from the where-parameters
func ComputeSelectQueryByParam(tableName string, where mccrud.WhereParamType, tableFields []string ) (string, error) {
	if tableName == "" || len(where) < 1 {
		return "", errors.New("table-name and where-params are required to perform the select operation")
	}
	selectQuery := ""
	fieldLen := len(tableFields)
	switch fieldLen {
	case 0:
		selectQuery = "SELECT * FROM " + tableName
	default:
		// get record(s) based on projected/provided field names ([]string)
		// select field/column-names
		selectQuery = "SELECT " + strings.Join(tableFields, ", ")
		// from
		selectQuery += fmt.Sprintf(" FROM %v", tableName)
	}
	// add where-params condition
	if whereScript, err := ComputeWhereQuery(where, tableFields); err == nil {
		selectQuery += selectQuery + whereScript
		return selectQuery, nil
	} else {
		return "", errors.New(fmt.Sprintf("error computing where-query condition(s): %v", err.Error()))
	}
}

// TODO: select-query functions for relational tables (eager & lazy queries) and data aggregation

