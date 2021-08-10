// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: compute where-SQL script

package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/abbeymart/mccrud"
	"github.com/asaskevich/govalidator"
	"time"
)

func whereErrMessage(errMsg string) (mccrud.WhereQueryObject, error) {
	return mccrud.WhereQueryObject{
		WhereQuery:  "",
		FieldValues: nil,
	}, errors.New(errMsg)
}

// ComputeWhereQuery function computes the multi-cases where-conditions for crud-operations
func ComputeWhereQuery(queryParams mccrud.QueryParamType, fieldLength int) (mccrud.WhereQueryObject, error) {
	if len(queryParams) < 1 || fieldLength < 1 {
		return whereErrMessage("queryParams condition and fieldLength(the start of the where-query-value-placeholders are required")
	}
	// compute queryParams script from queryParams
	whereQuery := "WHERE "
	var fieldValues []interface{}
	fieldCount := 0
	whereFieldLength := len(queryParams)
	for fieldName, fieldValue := range queryParams {
		// update fieldValues by fieldValue-type, for correct postgres-SQL-parsing
		var currentFieldValue interface{}
		switch fieldValue.(type) {
		case time.Time:
			if fVal, ok := fieldValue.(time.Time); !ok {
				return whereErrMessage(fmt.Sprintf("field_name: %v [date-type] | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = "'" + fVal.Format("2006-01-02 15:04:05.000000") + "'"
			}
		case string:
			if fVal, ok := fieldValue.(string); !ok {
				return whereErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				if govalidator.IsUUID(fVal) {
					currentFieldValue = fVal
				} else if govalidator.IsJSON(fVal) {
					if fValue, err := govalidator.ToJSON(fieldValue); err != nil {
						return whereErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fmt.Printf("string-toJson-value: %v\n\n", fValue)
						currentFieldValue = fValue
					}
				} else {
					currentFieldValue = "'" + fVal + "'"
				}
			}
		case int, uint, float32, float64, bool:
			currentFieldValue = fieldValue
		default:
			// json-stringify fieldValue
			if fVal, err := json.Marshal(fieldValue); err != nil {
				return whereErrMessage(fmt.Sprintf("Unknown or Unsupported field-value type: %v", err.Error()))
			} else {
				currentFieldValue = string(fVal)
			}
		}

		fieldValues = append(fieldValues, currentFieldValue)
		whereQuery += fmt.Sprintf("%v=$%v", fieldName, fieldLength)
		fieldCount += 1
		fieldLength += 1
		if whereFieldLength > 1 && fieldCount < whereFieldLength {
			whereQuery += ", "
		}
	}

	// if all went well, return valid queryParams script
	return mccrud.WhereQueryObject{
		WhereQuery:  whereQuery,
		FieldValues: fieldValues,
	}, nil
}
