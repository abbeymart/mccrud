// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: compute update-SQL scripts

package helper

import (
	"errors"
	"fmt"
	"github.com/abbeymart/mctypes"
	"strings"
)

func ComputeUpdateQuery(tableName string, tableFields []string, actionParams mctypes.ActionParamsType) ([]string, error) {
	if tableName == "" || len(actionParams) < 1 || len(tableFields) < 1 {
		return nil, errors.New("table-name, table-fields and action-params are required for the update operation")
	}
	// compute update script from queryParams
	var updateQuery []string
	validUpdateItemCount := 0
	invalidUpdateItemCount := 0

	for recNum, rec := range actionParams {
		itemScript := fmt.Sprintf("UPDATE %v SET", tableName)
		fieldCount := 0
		fieldLen := len(rec)
		for _, fieldName := range tableFields {
			fieldValue, ok := rec[fieldName]
			// check for the required fields in each record
			if !ok {
				return nil, errors.New(fmt.Sprintf("Record #%v [%#v]: required field_name[%v] is missing", recNum, rec, fieldName))
			}
			fieldCount += 1
			itemScript += fmt.Sprintf(" %v=%v", fieldName, fieldValue)

			if fieldLen > 1 && fieldCount < fieldLen {
				itemScript += ", "
			}
		}

		// add where condition by id
		itemScript += fmt.Sprintf(" WHERE id=%v", rec["id"])
		//validate/update script content based on valid field specifications
		if fieldCount > 0 && fieldCount == fieldLen {
			validUpdateItemCount += 1
			updateQuery = append(updateQuery, itemScript)
		} else {
			invalidUpdateItemCount += 1
		}
	}
	// check is there was invalid update items
	if invalidUpdateItemCount > 0 {
		return nil, errors.New(fmt.Sprintf("Invalid action-params [%v]", invalidUpdateItemCount))
	}
	return updateQuery, nil
}

func ComputeUpdateQueryById(tableName string, tableFields []string, actionParams mctypes.ActionParamsType, recordIds []string) (string, error) {
	if tableName == "" || len(actionParams) < 1 || len(tableFields) < 1 || len(recordIds) < 1 {
		return "", errors.New("table-name, table-fields, action-params and record/doc-Ids are required for the update-by-id operation")
	}
	// compute update script from query-ids
	var updateQuery string
	itemScript := fmt.Sprintf("UPDATE %v SET", tableName)
	whereQuery := fmt.Sprintf(" WHERE id IN(%v)", strings.Join(recordIds, ", "))
	invalidUpdateItemCount := 0
	validUpdateItemCount := 0

	// only one actionParams record is required for update by docIds
	rec := actionParams[0]
	fieldCount := 0
	fieldLen := len(rec)
	for _, fieldName := range tableFields {
		fieldValue, ok := rec[fieldName]
		// check for the required fields in each record
		if !ok {
			return "", errors.New(fmt.Sprintf("Record [%#v]: required field_name[%v] is missing", rec, fieldName))
		}
		fieldCount += 1
		itemScript += fmt.Sprintf(" %v=%v", fieldName, fieldValue)

		if fieldLen > 1 && fieldCount < fieldLen {
			itemScript += ", "
		}
	}
	//validate/update script content based on valid field specifications
	if fieldCount > 0 && fieldCount == fieldLen {
		validUpdateItemCount += 1
		updateQuery = itemScript + whereQuery
	} else {
		invalidUpdateItemCount += 1
	}

	// check is there was invalid update items
	if invalidUpdateItemCount > 0 {
		return "", errors.New(fmt.Sprintf("Invalid action-params [%v]", invalidUpdateItemCount))
	}
	return updateQuery, nil
}

func ComputeUpdateQueryByParam(tableName string, tableFields []string, actionParams mctypes.ActionParamsType, where mctypes.WhereParamType) (string, error) {
	if tableName == "" || len(actionParams) < 1 || len(tableFields) < 1 || len(where) < 1 {
		return "", errors.New("table-name, table-fields, action-params and where-params are required for the update-by-params operation")
	}
	// compute update script from queryParams
	var updateQuery string
	invalidUpdateItemCount := 0
	validUpdateItemCount := 0

	// only one actionParams record is required for update by where-params
	rec := actionParams[0]
	itemScript := fmt.Sprintf("UPDATE %v SET", tableName)
	fieldCount := 0
	fieldLen := len(rec)
	for _, fieldName := range tableFields {
		fieldValue, ok := rec[fieldName]
		// check for the required fields in each record
		if !ok {
			return "", errors.New(fmt.Sprintf("Record [%#v]: required field_name[%v] is missing", rec, fieldName))
		}
		fieldCount += 1
		itemScript += fmt.Sprintf(" %v=%v", fieldName, fieldValue)

		if fieldLen > 1 && fieldCount < fieldLen {
			itemScript += ", "
		}
	}
	//validate/update script content based on valid field specifications
	if fieldCount > 0 && fieldCount == fieldLen {
		validUpdateItemCount += 1
		updateQuery = itemScript
	} else {
		invalidUpdateItemCount += 1
	}

	// check is there was invalid update items
	if invalidUpdateItemCount > 0 {
		return "", errors.New(fmt.Sprintf("Invalid action-params [%v]", invalidUpdateItemCount))
	}

	if whereScript, err := ComputeWhereQuery(where, tableFields); err == nil {
		updateQuery += whereScript
		return updateQuery, nil
	} else {
		return "", errors.New(fmt.Sprintf("error computing where-query condition(s): %v", err.Error()))
	}
}
