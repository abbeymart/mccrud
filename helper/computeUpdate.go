// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: compute update-SQL scripts

package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/abbeymart/mccrud"
	"github.com/asaskevich/govalidator"
	"time"
)

func updateErrMessage(errMsg string) (mccrud.UpdateQueryObject, error) {
	return mccrud.UpdateQueryObject{
		UpdateQuery: "",
		FieldNames:  nil,
		WhereQuery:  mccrud.WhereQueryObject{},
		FieldValues: nil,
	}, errors.New(errMsg)
}

// ComputeUpdateQuery function computes update SQL script. It returns updateScript, updateValues []interface{} and/or err error
func ComputeUpdateQuery(tableName string, actionParam mccrud.ActionParamType) (mccrud.UpdateQueryObject, error) {
	if tableName == "" || len(actionParam) < 1 || actionParam == nil {
		return updateErrMessage("table-name and actionParam are required for the update operation")
	}
	// validate actionParam/record-id
	recId, ok := actionParam["id"]
	if !ok || recId == "" {
		return updateErrMessage(fmt.Sprintf("actionParam/record-is is required for the update operation: %v", actionParam))
	}

	// declare slice variable for create/insert queries
	var whereQuery string
	var fieldNames []string
	var fieldValues []interface{}

	// compute update script and associated values () for all the actionParam/record
	updateQuery := fmt.Sprintf("UPDATE %v SET", tableName)
	fieldsLength := len(actionParam)
	fieldCount := 0
	for fieldName := range actionParam {
		// skip fieldName=="id"
		if fieldName == "id" {
			continue
		}
		fieldCount += 1
		fieldNames = append(fieldNames, fieldName)
		updateQuery += fmt.Sprintf(" %v=%v", fieldName, fieldCount)
		if fieldsLength > 1 && fieldCount < fieldsLength {
			updateQuery += ", "
		}
	}
	// close item-script/value-placeholder
	updateQuery += " )"
	// add where condition by id
	updateQuery += fmt.Sprintf(" WHERE id='%v'", actionParam["id"])

	// compute update-field-values from actionParams/records, in order of the fields-sequence
	// value-computation for each of the actionParam / record must match the record-fields

	// item-values-computation variable
	for _, fieldName := range fieldNames {
		fieldValue, ok := actionParam[fieldName]
		// check for required field in each record
		if !ok {
			return updateErrMessage(fmt.Sprintf("Record #%v: required field_name[%v] has field_value of %v ", actionParam, fieldName, fieldValue))
		}
		// update fieldValues by fieldValue-type, for correct postgres-SQL-parsing
		var currentFieldValue interface{}
		switch fieldValue.(type) {
		case time.Time:
			if fVal, ok := fieldValue.(time.Time); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v [date-type] | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = "'" + fVal.Format("2006-01-02 15:04:05.000000") + "'"
			}
		case string:
			if fVal, ok := fieldValue.(string); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				if govalidator.IsUUID(fVal) {
					currentFieldValue = fVal
				} else if govalidator.IsJSON(fVal) {
					if fValue, err := govalidator.ToJSON(fieldValue); err != nil {
						return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
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
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case int8:
			if fVal, ok := fieldValue.(int8); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case int16:
			if fVal, ok := fieldValue.(int16); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case int32:
			if fVal, ok := fieldValue.(int32); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case int64:
			if fVal, ok := fieldValue.(int64); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case int:
			if fVal, ok := fieldValue.(int); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case uint8:
			if fVal, ok := fieldValue.(uint8); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case uint16:
			if fVal, ok := fieldValue.(uint16); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case uint32:
			if fVal, ok := fieldValue.(uint32); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case uint64:
			if fVal, ok := fieldValue.(uint64); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case uint:
			if fVal, ok := fieldValue.(uint); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case float32:
			if fVal, ok := fieldValue.(float32); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case float64:
			if fVal, ok := fieldValue.(float64); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case []string:
			if fVal, ok := fieldValue.([]string); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case []int:
			if fVal, ok := fieldValue.([]int); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case []float32:
			if fVal, ok := fieldValue.([]float32); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case []float64:
			if fVal, ok := fieldValue.([]float64); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case []struct{}:
			if fVal, ok := fieldValue.([]struct{}); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		default:
			// json-stringify fieldValue
			if fVal, err := json.Marshal(fieldValue); err != nil {
				return updateErrMessage(fmt.Sprintf("Unknown or Unsupported field-value type: %v", err.Error()))
			} else {
				currentFieldValue = string(fVal)
			}
		}
		// add itemValue
		fieldValues = append(fieldValues, currentFieldValue)
	}

	// result
	return mccrud.UpdateQueryObject{
		UpdateQuery: updateQuery,
		FieldNames:  fieldNames,
		WhereQuery: mccrud.WhereQueryObject{
			WhereQuery:  whereQuery,
			FieldValues: nil,
		},
		FieldValues: fieldValues,
	}, nil
}

// ComputeUpdateQueryById function computes update SQL script by recordIds. It returns updateScript, updateValues []interface{} and/or err error
func ComputeUpdateQueryById(tableName string, actionParam mccrud.ActionParamType, recordIds []string) (mccrud.UpdateQueryObject, error) {
	if tableName == "" || len(actionParam) < 1 || actionParam == nil {
		return updateErrMessage("table-name and actionParam are required for the update operation")
	}
	// validate actionParam/record-ids
	if len(recordIds) < 1 {
		return updateErrMessage(fmt.Sprintf("actionParam/recordIds are required for the update operation: %v", actionParam))
	}
	// from / where condition (where-in-values)
	whereIds := ""
	idLen := len(recordIds)
	for idCount, id := range recordIds {
		whereIds += "'" + id + "'"
		if idLen > 1 && idCount < idLen-1 {
			whereIds += ", "
		}
	}
	whereQuery := fmt.Sprintf(" WHERE id IN(%v)", whereIds)

	// declare slice variable for create/insert queries
	var fieldNames []string
	var fieldValues []interface{}

	// compute update script and associated values () for all the actionParam/record
	updateQuery := fmt.Sprintf("UPDATE %v SET", tableName)
	fieldsLength := len(actionParam)
	fieldCount := 0
	for fieldName := range actionParam {
		// skip fieldName="id"
		if fieldName == "id" {
			continue
		}
		fieldCount += 1
		fieldNames = append(fieldNames, fieldName)
		updateQuery += fmt.Sprintf(" %v=%v", fieldName, fieldCount)
		if fieldsLength > 1 && fieldCount < fieldsLength {
			updateQuery += ", "
		}
	}
	// close item-script/value-placeholder
	updateQuery += " )"
	// add where condition by id
	updateQuery += whereQuery
	// compute update-field-values from actionParams/records, in order of the fields-sequence
	// value-computation for each of the actionParam / record must match the record-fields
	// item-values-computation variable
	for _, fieldName := range fieldNames {
		fieldValue, ok := actionParam[fieldName]
		// check for required field in each record
		if !ok {
			return updateErrMessage(fmt.Sprintf("Record #%v: required field_name[%v] has field_value of %v ", actionParam, fieldName, fieldValue))
		}
		// update fieldValues by fieldValue-type, for correct postgres-SQL-parsing
		var currentFieldValue interface{}
		switch fieldValue.(type) {
		case time.Time:
			if fVal, ok := fieldValue.(time.Time); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v [date-type] | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = "'" + fVal.Format("2006-01-02 15:04:05.000000") + "'"
			}
		case string:
			if fVal, ok := fieldValue.(string); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				if govalidator.IsUUID(fVal) {
					currentFieldValue = fVal
				} else if govalidator.IsJSON(fVal) {
					if fValue, err := govalidator.ToJSON(fieldValue); err != nil {
						return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
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
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case int8:
			if fVal, ok := fieldValue.(int8); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case int16:
			if fVal, ok := fieldValue.(int16); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case int32:
			if fVal, ok := fieldValue.(int32); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case int64:
			if fVal, ok := fieldValue.(int64); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case int:
			if fVal, ok := fieldValue.(int); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case uint8:
			if fVal, ok := fieldValue.(uint8); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case uint16:
			if fVal, ok := fieldValue.(uint16); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case uint32:
			if fVal, ok := fieldValue.(uint32); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case uint64:
			if fVal, ok := fieldValue.(uint64); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case uint:
			if fVal, ok := fieldValue.(uint); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case float32:
			if fVal, ok := fieldValue.(float32); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case float64:
			if fVal, ok := fieldValue.(float64); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case []string:
			if fVal, ok := fieldValue.([]string); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case []int:
			if fVal, ok := fieldValue.([]int); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case []float32:
			if fVal, ok := fieldValue.([]float32); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case []float64:
			if fVal, ok := fieldValue.([]float64); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case []struct{}:
			if fVal, ok := fieldValue.([]struct{}); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		default:
			// json-stringify fieldValue
			if fVal, err := json.Marshal(fieldValue); err != nil {
				return updateErrMessage(fmt.Sprintf("Unknown or Unsupported field-value type: %v", err.Error()))
			} else {
				currentFieldValue = string(fVal)
			}
		}
		// add itemValue
		fieldValues = append(fieldValues, currentFieldValue)
	}

	// result
	return mccrud.UpdateQueryObject{
		UpdateQuery: updateQuery,
		FieldNames:  fieldNames,
		WhereQuery: mccrud.WhereQueryObject{
			WhereQuery:  whereQuery,
			FieldValues: nil,
		},
		FieldValues: fieldValues,
	}, nil
}

// ComputeUpdateQueryByParam function computes update SQL scripts by queryParams. It returns updateScript, updateValues []interface{} and/or err error
func ComputeUpdateQueryByParam(tableName string, actionParam mccrud.ActionParamType, queryParams mccrud.QueryParamType) (mccrud.UpdateQueryObject, error) {
	if tableName == "" || len(actionParam) < 1 || actionParam == nil {
		return updateErrMessage("table-name and actionParam are required for the update operation")
	}
	// validate actionParam/record-ids
	if len(queryParams) < 1 {
		return updateErrMessage(fmt.Sprintf("queryParams is required for the update operation: %v", actionParam))
	}
	// declare slice variable for create/insert queries
	var fieldNames []string
	var fieldValues []interface{}

	// compute update script and associated values () for all the actionParam/record
	updateQuery := fmt.Sprintf("UPDATE %v SET", tableName)
	fieldsLength := len(actionParam)
	fieldCount := 0
	for fieldName := range actionParam {
		// skip fieldName="id"
		if fieldName == "id" {
			continue
		}
		fieldCount += 1
		fieldNames = append(fieldNames, fieldName)
		updateQuery += fmt.Sprintf(" %v=$%v", fieldName, fieldCount)
		if fieldsLength > 1 && fieldCount < fieldsLength {
			updateQuery += ", "
		}
	}
	// close item-script/value-placeholder
	updateQuery += " )"
	// add where condition by queryParams | TODO: refactor ComputeWhereQuery to support query-value-placeholders
	whereQuery, err := ComputeWhereQuery(queryParams, fieldsLength)
	if err != nil {
		return updateErrMessage(fmt.Sprintf("error computing where-query condition(s): %v", err.Error()))
	}
	// compute update-field-values from actionParams/records, in order of the fields-sequence
	// value-computation for each of the actionParam / record must match the record-fields
	// item-values-computation variable
	for _, fieldName := range fieldNames {
		fieldValue, ok := actionParam[fieldName]
		// check for required field in each record
		if !ok {
			return updateErrMessage(fmt.Sprintf("Record #%v: required field_name[%v] has field_value of %v ", actionParam, fieldName, fieldValue))
		}
		// update fieldValues by fieldValue-type, for correct postgres-SQL-parsing
		var currentFieldValue interface{}
		switch fieldValue.(type) {
		case time.Time:
			if fVal, ok := fieldValue.(time.Time); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v [date-type] | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = "'" + fVal.Format("2006-01-02 15:04:05.000000") + "'"
			}
		case string:
			if fVal, ok := fieldValue.(string); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				if govalidator.IsUUID(fVal) {
					currentFieldValue = fVal
				} else if govalidator.IsJSON(fVal) {
					if fValue, err := govalidator.ToJSON(fieldValue); err != nil {
						return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
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
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case int8:
			if fVal, ok := fieldValue.(int8); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case int16:
			if fVal, ok := fieldValue.(int16); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case int32:
			if fVal, ok := fieldValue.(int32); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case int64:
			if fVal, ok := fieldValue.(int64); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case int:
			if fVal, ok := fieldValue.(int); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case uint8:
			if fVal, ok := fieldValue.(uint8); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case uint16:
			if fVal, ok := fieldValue.(uint16); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case uint32:
			if fVal, ok := fieldValue.(uint32); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case uint64:
			if fVal, ok := fieldValue.(uint64); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case uint:
			if fVal, ok := fieldValue.(uint); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case float32:
			if fVal, ok := fieldValue.(float32); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case float64:
			if fVal, ok := fieldValue.(float64); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case []string:
			if fVal, ok := fieldValue.([]string); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case []int:
			if fVal, ok := fieldValue.([]int); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case []float32:
			if fVal, ok := fieldValue.([]float32); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case []float64:
			if fVal, ok := fieldValue.([]float64); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		case []struct{}:
			if fVal, ok := fieldValue.([]struct{}); !ok {
				return updateErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
			} else {
				currentFieldValue = fVal
			}
		default:
			// json-stringify fieldValue
			if fVal, err := json.Marshal(fieldValue); err != nil {
				return updateErrMessage(fmt.Sprintf("Unknown or Unsupported field-value type: %v", err.Error()))
			} else {
				currentFieldValue = string(fVal)
			}
		}
		// add itemValue
		fieldValues = append(fieldValues, currentFieldValue)
	}

	// result
	return mccrud.UpdateQueryObject{
		UpdateQuery: updateQuery,
		FieldNames:  fieldNames,
		WhereQuery:  whereQuery,
		FieldValues: fieldValues,
	}, nil
}
