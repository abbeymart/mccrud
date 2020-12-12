// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: compute create-SQL script

package helper

import (
	"errors"
	"fmt"
	"github.com/abbeymart/mccrudgo"
)

func errMessage(errMsg string) (mccrud.CreateScriptResponseType, error) {
	return mccrud.CreateScriptResponseType{
		CreateScript: nil,
		FieldNames:   nil,
		FieldValues:  nil,
	}, errors.New(errMsg)
}

// ComputeCreateScript computes insert SQL script. It returns createScripts []string, fieldNames []string and err error
func ComputeCreateQuery(tableName string, actionParams mccrud.ActionParamsType) (mccrud.CreateScriptResponseType, error) {
	var insertScripts []string
	var fNames []string         // fieldNames array of strings in order of SQL statement
	var fValues [][]interface{} // fieldValues array of ValueParamType

	if tableName == "" || len(actionParams) < 1 {
		return errMessage("table/collection name and action-params are required for the create operation")
	}
	// compute fieldNames, from one of the actionParams items/records
	// script-computation for each of the actionParams' records must match the fieldNames
	for fName := range actionParams[0] {
		fNames = append(fNames, fName)
	}

	// computed create script from actionParams
	for _, rec := range actionParams {
		// initial item-script/value-computation variables
		var itemScript = "INSERT INTO " + tableName + " ("
		var itemValuePlaceholder = " VALUES("
		var recFieldValues []interface{}
		var (
			fieldCount = 0
		)
		fieldsLength := len(rec)
		for fieldName, fieldValue := range rec {
			if !ArrayStringContains(fNames, fieldName) {
				return errMessage(fmt.Sprintf("Missing field: %v from record %v", fieldName, rec))
			}
			if fieldsLength > 1 && fieldCount < fieldsLength-1 {
				itemScript += ", " + fieldName
				itemValuePlaceholder += ", $" + fmt.Sprintf("%v", fieldCount+1)
				//itemValuePlaceholder += ", ?" + fmt.Sprintf("%v", fieldValue)
			} else {
				itemScript += " " + fieldName
				itemValuePlaceholder += " $" + fmt.Sprintf("%v", fieldCount+1)
				//itemValuePlaceholder += " " + fmt.Sprintf("%v", fieldValue)
			}
			// increment fieldCount for the current record
			fieldCount += 1
			// update recFieldValues by fieldValue-type
			switch fieldValue.(type) {
			case string:
				if fVal, ok := fieldValue.(string); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					transformedFieldValue := fVal
					recFieldValues = append(recFieldValues, transformedFieldValue)
				}
			case bool:
				if fVal, ok := fieldValue.(bool); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					transformedFieldValue := fVal
					recFieldValues = append(recFieldValues, transformedFieldValue)
				}
			case int8:
				if fVal, ok := fieldValue.(int8); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					transformedFieldValue := fVal
					recFieldValues = append(recFieldValues, transformedFieldValue)
				}
			case int16:
				if fVal, ok := fieldValue.(int16); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					transformedFieldValue := fVal
					recFieldValues = append(recFieldValues, transformedFieldValue)
				}
			case int32:
				if fVal, ok := fieldValue.(int32); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					transformedFieldValue := fVal
					recFieldValues = append(recFieldValues, transformedFieldValue)
				}
			case int64:
				if fVal, ok := fieldValue.(int64); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					transformedFieldValue := fVal
					recFieldValues = append(recFieldValues, transformedFieldValue)
				}
			case int:
				if fVal, ok := fieldValue.(int); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					transformedFieldValue := fVal
					recFieldValues = append(recFieldValues, transformedFieldValue)
				}
			case uint8:
				if fVal, ok := fieldValue.(uint8); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					transformedFieldValue := fVal
					recFieldValues = append(recFieldValues, transformedFieldValue)
				}
			case uint16:
				if fVal, ok := fieldValue.(uint16); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					transformedFieldValue := fVal
					recFieldValues = append(recFieldValues, transformedFieldValue)
				}
			case uint32:
				if fVal, ok := fieldValue.(uint32); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					transformedFieldValue := fVal
					recFieldValues = append(recFieldValues, transformedFieldValue)
				}
			case uint64:
				if fVal, ok := fieldValue.(uint64); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					transformedFieldValue := fVal
					recFieldValues = append(recFieldValues, transformedFieldValue)
				}
			case uint:
				if fVal, ok := fieldValue.(uint); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					transformedFieldValue := fVal
					recFieldValues = append(recFieldValues, transformedFieldValue)
				}
			case float32:
				if fVal, ok := fieldValue.(float32); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					transformedFieldValue := fVal
					recFieldValues = append(recFieldValues, transformedFieldValue)
				}
			case float64:
				if fVal, ok := fieldValue.(float64); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					transformedFieldValue := fVal
					recFieldValues = append(recFieldValues, transformedFieldValue)
				}
			case []string:
				if fVal, ok := fieldValue.([]string); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					transformedFieldValue := fVal
					recFieldValues = append(recFieldValues, transformedFieldValue)
				}
			case []int:
				if fVal, ok := fieldValue.([]int); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					transformedFieldValue := fVal
					recFieldValues = append(recFieldValues, transformedFieldValue)
				}
			case []struct{}:
				if fVal, ok := fieldValue.([]struct{}); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					transformedFieldValue := fVal
					recFieldValues = append(recFieldValues, transformedFieldValue)
				}
			default:
				recFieldValues = append(recFieldValues, fieldValue)
			}
		}
		// update fieldValues
		fValues = append(fValues, recFieldValues)
		// re-initialise recFieldValues
		recFieldValues = []interface{}{}
		// close item-script/value-placeholder
		itemScript += " )"
		itemValuePlaceholder += " )"
		// add/append item-script & value-placeholder to the createScripts
		insertScripts = append(insertScripts, itemScript+itemValuePlaceholder)
	}
	// result
	return mccrud.CreateScriptResponseType{
		CreateScript: insertScripts,
		FieldNames:   fNames,
		FieldValues:  fValues,
	}, nil
}
