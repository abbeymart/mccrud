// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: compute create-SQL script, for bulk/copy insert operation

package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/abbeymart/mccrud"
	"github.com/asaskevich/govalidator"
	"time"
)

func errMessage(errMsg string) (mccrud.CreateQueryObject, error) {
	return mccrud.CreateQueryObject{
		CreateQuery: "",
		FieldNames:  nil,
		FieldValues: nil,
	}, errors.New(errMsg)
}

// ComputeCreateQuery function computes insert SQL scripts. It returns createScripts []string and err error
func ComputeCreateQuery(tableName string, actionParams mccrud.ActionParamsType) (mccrud.CreateQueryObject, error) {
	if tableName == "" || len(actionParams) < 1 {
		return errMessage("table-name is required for the create operation")
	}

	// declare slice variable for create/insert queries
	var createQuery string
	var fieldNames []string
	var fieldValues [][]interface{}

	// compute create script and associated values () for all the records in actionParams
	// compute create-query from the first actionParams
	itemQuery := fmt.Sprintf("INSERT INTO %v(", tableName)
	itemValuePlaceholder := " VALUES("
	fieldsLength := len(actionParams[0])
	fieldCount := 0
	for fieldName := range actionParams[0] {
		fieldCount += 1
		fieldNames = append(fieldNames, fieldName)
		itemQuery += fmt.Sprintf(" %v", fieldName)
		itemValuePlaceholder += fmt.Sprintf(" $%v", fieldCount)
		if fieldsLength > 1 && fieldCount < fieldsLength {
			itemQuery += ", "
			itemValuePlaceholder += ", "
		}
	}
	// close item-script/value-placeholder
	itemQuery += " )"
	itemValuePlaceholder += " )"
	// add/append item-script & value-placeholder to the createScript
	createQuery = itemQuery + itemValuePlaceholder + " RETURNING id"

	// compute create-record-values from actionParams/records, in order of the fields-sequence
	// value-computation for each of the actionParams / records must match the record-fields
	for recIndex, rec := range actionParams {
		// item-values-computation variable
		var recFieldValues []interface{}
		for _, fieldName := range fieldNames {
			fieldValue, ok := rec[fieldName]
			// check for required field in each record
			if !ok {
				return errMessage(fmt.Sprintf("Record #%v [%#v]: required field_name[%v] has field_value of %v ", recIndex, rec, fieldName, fieldValue))
			}
			// update recFieldValues by fieldValue-type, for correct postgres-SQL-parsing
			var currentFieldValue interface{}
			switch fieldValue.(type) {
			case time.Time:
				if fVal, ok := fieldValue.(time.Time); !ok {
					return errMessage(fmt.Sprintf("field_name: %v [date-type] | field_value: %v error: ", fieldName, fieldValue))
				} else {
					currentFieldValue = "'" + fVal.Format("2006-01-02 15:04:05.000000") + "'"
				}
			case string:
				if fVal, ok := fieldValue.(string); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					if govalidator.IsUUID(fVal) {
						currentFieldValue = fVal
					} else if govalidator.IsJSON(fVal) {
						if fValue, err := govalidator.ToJSON(fieldValue); err != nil {
							return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
						} else {
							fmt.Printf("string-toJson-value: %v\n\n", fValue)
							currentFieldValue = fValue
						}
					} else {
						currentFieldValue = "'" + fVal + "'"
					}
				}
			case bool:
				if fVal, ok := fieldValue.(bool); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					currentFieldValue = fVal
				}
			case int8:
				if fVal, ok := fieldValue.(int8); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					currentFieldValue = fVal
				}
			case int16:
				if fVal, ok := fieldValue.(int16); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					currentFieldValue = fVal
				}
			case int32:
				if fVal, ok := fieldValue.(int32); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					currentFieldValue = fVal
				}
			case int64:
				if fVal, ok := fieldValue.(int64); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					currentFieldValue = fVal
				}
			case int:
				if fVal, ok := fieldValue.(int); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					currentFieldValue = fVal
				}
			case uint8:
				if fVal, ok := fieldValue.(uint8); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					currentFieldValue = fVal
				}
			case uint16:
				if fVal, ok := fieldValue.(uint16); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					currentFieldValue = fVal
				}
			case uint32:
				if fVal, ok := fieldValue.(uint32); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					currentFieldValue = fVal
				}
			case uint64:
				if fVal, ok := fieldValue.(uint64); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					currentFieldValue = fVal
				}
			case uint:
				if fVal, ok := fieldValue.(uint); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					currentFieldValue = fVal
				}
			case float32:
				if fVal, ok := fieldValue.(float32); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					currentFieldValue = fVal
				}
			case float64:
				if fVal, ok := fieldValue.(float64); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					currentFieldValue = fVal
				}
			case []string:
				if fVal, ok := fieldValue.([]string); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					currentFieldValue = fVal
				}
			case []int:
				if fVal, ok := fieldValue.([]int); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					currentFieldValue = fVal
				}
			case []float32:
				if fVal, ok := fieldValue.([]float32); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					currentFieldValue = fVal
				}
			case []float64:
				if fVal, ok := fieldValue.([]float64); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					currentFieldValue = fVal
				}
			case []struct{}:
				if fVal, ok := fieldValue.([]struct{}); !ok {
					return errMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					currentFieldValue = fVal
				}
			default:
				// json-stringify fieldValue
				if fVal, err := json.Marshal(fieldValue); err != nil {
					return errMessage(fmt.Sprintf("Unknown or Unsupported field-value type: %v", err.Error()))
				} else {
					currentFieldValue = string(fVal)
				}
			}
			// add itemValue
			recFieldValues = append(recFieldValues, currentFieldValue)
		}
		// update fieldValues
		fieldValues = append(fieldValues, recFieldValues)
		// re-initialise recFieldValues, for next update
		recFieldValues = []interface{}{}
	}

	// result
	return mccrud.CreateQueryObject{
		CreateQuery: createQuery,
		FieldNames: fieldNames,
		FieldValues: fieldValues,
	}, nil
}
