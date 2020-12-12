// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: compute create-SQL script, for bulk/copy insert operation

package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/abbeymart/mccrudgo"
)

func errMessage(errMsg string) (mccrud.CreateScriptResponseType, error) {
	return mccrud.CreateScriptResponseType{
		CreateScript: "",
		FieldNames:   nil,
		FieldValues:  nil,
	}, errors.New(errMsg)
}

// ComputeCreateScript computes insert SQL script. It returns createScripts []string, fieldNames []string and err error
func ComputeCreateQuery(tableName string, actionParams mccrud.ActionParamsType, tableFields []string) (mccrud.CreateScriptResponseType, error) {
	var insertScript string
	var fValues [][]interface{} // fieldValues array of ValueParamType

	if tableName == "" || len(actionParams) < 1 || len(tableFields) < 1 {
		return errMessage("table-name, action-params and table-fields are required for the create operation")
	}
	// value-computation for each of the actionParams' records must match the fieldNames
	// compute create script for all the create-task, with value-placeholders
	var itemScript = fmt.Sprintf("INSERT INTO %v(", tableName)
	var itemValuePlaceholder = " VALUES("
	fieldsLength := len(tableFields)
	for fieldIndex, fieldName := range tableFields {
		if fieldsLength > 1 && fieldIndex < fieldsLength-1 {
			itemScript += fmt.Sprintf(", %v", fieldName)
			itemValuePlaceholder += fmt.Sprintf(", $%v", fieldIndex+1)
		} else {
			itemScript += fmt.Sprintf(" %v", fieldName)
			itemValuePlaceholder += fmt.Sprintf(" $%v", fieldIndex+1)
		}
	}
	// close item-script/value-placeholder
	itemScript += " )"
	itemValuePlaceholder += " )"
	// add/append item-script & value-placeholder to the createScripts
	insertScript = itemScript + itemValuePlaceholder

	// computed create values from actionParams
	for _, rec := range actionParams {
		// initial item-values-computation variables
		var recFieldValues []interface{}
		for fieldName, fieldValue := range rec {
			// check for missing field in each record
			if !ArrayStringContains(tableFields, fieldName) {
				return errMessage(fmt.Sprintf("Missing field: %v from record %v", fieldName, rec))
			}
			// update recFieldValues by fieldValue-type
			switch fieldValue.(type) {
			case string:
				if fVal, ok := fieldValue.(string); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					recFieldValues = append(recFieldValues, fVal)
				}
			case bool:
				if fVal, ok := fieldValue.(bool); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					recFieldValues = append(recFieldValues, fVal)
				}
			case int8:
				if fVal, ok := fieldValue.(int8); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					recFieldValues = append(recFieldValues, fVal)
				}
			case int16:
				if fVal, ok := fieldValue.(int16); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					recFieldValues = append(recFieldValues, fVal)
				}
			case int32:
				if fVal, ok := fieldValue.(int32); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					recFieldValues = append(recFieldValues, fVal)
				}
			case int64:
				if fVal, ok := fieldValue.(int64); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					recFieldValues = append(recFieldValues, fVal)
				}
			case int:
				if fVal, ok := fieldValue.(int); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					recFieldValues = append(recFieldValues, fVal)
				}
			case uint8:
				if fVal, ok := fieldValue.(uint8); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					recFieldValues = append(recFieldValues, fVal)
				}
			case uint16:
				if fVal, ok := fieldValue.(uint16); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					recFieldValues = append(recFieldValues, fVal)
				}
			case uint32:
				if fVal, ok := fieldValue.(uint32); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					recFieldValues = append(recFieldValues, fVal)
				}
			case uint64:
				if fVal, ok := fieldValue.(uint64); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					recFieldValues = append(recFieldValues, fVal)
				}
			case uint:
				if fVal, ok := fieldValue.(uint); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					recFieldValues = append(recFieldValues, fVal)
				}
			case float32:
				if fVal, ok := fieldValue.(float32); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					recFieldValues = append(recFieldValues, fVal)
				}
			case float64:
				if fVal, ok := fieldValue.(float64); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					recFieldValues = append(recFieldValues, fVal)
				}
			case []string:
				if fVal, ok := fieldValue.([]string); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					recFieldValues = append(recFieldValues, fVal)
				}
			case []int:
				if fVal, ok := fieldValue.([]int); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					recFieldValues = append(recFieldValues, fVal)
				}
			case []struct{}:
				if fVal, ok := fieldValue.([]struct{}); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					recFieldValues = append(recFieldValues, fVal)
				}
			default:
				// json-stringify fieldValue
				if fVal, err := json.Marshal(fieldValue); err != nil {
					return errMessage(fmt.Sprintf("Unknown or Unsupported field-value type: %v", err.Error()))
				} else {
					recFieldValues = append(recFieldValues, fVal)
				}
			}
		}
		// update fieldValues
		fValues = append(fValues, recFieldValues)
		// re-initialise recFieldValues, for next update
		recFieldValues = []interface{}{}
	}
	// result
	return mccrud.CreateScriptResponseType{
		CreateScript: insertScript,
		FieldNames:   tableFields,
		FieldValues:  fValues,
	}, nil
}
