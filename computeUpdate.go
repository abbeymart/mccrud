// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: compute update-SQL scripts

package mccrud

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"time"
)

func updateErrMessage(errMsg string) UpdateQueryResult {
	return UpdateQueryResult{
		UpdateQueryObject: UpdateQueryObject{
			UpdateQuery: "",
			FieldNames:  nil,
			FieldValues: nil,
		},
		Ok:      false,
		Message: errMsg,
	}
}

func updatesErrMessage(errMsg string) MultiUpdateQueryResult {
	return MultiUpdateQueryResult{
		UpdateQueryObjects: []UpdateQueryObject{},
		Ok:                 false,
		Message:            errMsg,
	}
}

// TODO: review/refactor

// ComputeUpdateQuery function computes update SQL script. It returns updateScript, updateValues []interface{} and/or err error
func ComputeUpdateQuery(tableName string, actionParams ActionParamsType) MultiUpdateQueryResult {
	if tableName == "" || len(actionParams) < 1 {
		return updatesErrMessage("tableName and actionParam are required for the update operation")
	}
	var updateQueryObjects []UpdateQueryObject
	for _, actParam := range actionParams {
		// compute update script and associated place-holder values for the actionParam/record
		updateQuery := fmt.Sprintf("UPDATE %v SET ", tableName)
		var fieldValues []interface{}
		var fieldNames []string
		var fieldNamesUnderscore []string
		fieldsLength := len(actParam)
		fieldCount := 0
		recordId := ""
		for fieldName, fieldValue := range actParam {
			// next placeholder-value-position
			fieldCount += 1
			// skip fieldName=="id"
			if fieldName == "id" {
				recordId = fmt.Sprintf("%v", actParam["id"])
				fieldCount -= 1
				fieldsLength -= 1
				continue
			}
			fieldNameUnderScore := govalidator.CamelCaseToUnderscore(fieldName)
			fieldNames = append(fieldNames, fieldName)
			fieldNamesUnderscore = append(fieldNamesUnderscore, fieldNameUnderScore)
			// TODO: update fieldValues by fieldValue-type, for correct postgres-SQL-parsing
			var currentFieldValue interface{}
			switch fieldValue.(type) {
			case time.Time:
				if fVal, ok := fieldValue.(time.Time); !ok {
					return updatesErrMessage(fmt.Sprintf("field_name: %v [date-type] | field_value: %v error: ", fieldName, fieldValue))
				} else {
					currentFieldValue = "'" + fVal.Format("2006-01-02 15:04:05.000000") + "'"
				}
			case string:
				if fVal, ok := fieldValue.(string); !ok {
					return updatesErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
				} else {
					if govalidator.IsUUID(fVal) {
						currentFieldValue = fVal
					} else if govalidator.IsJSON(fVal) {
						if fValue, err := govalidator.ToJSON(fieldValue); err != nil {
							return updatesErrMessage(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
						} else {
							fmt.Printf("string-toJson-value: %v\n\n", fValue)
							currentFieldValue = fValue
						}
					} else {
						currentFieldValue = "'" + fVal + "'"
					}
				}
			default:
				currentFieldValue = fieldValue
			}

			fieldValues = append(fieldValues, currentFieldValue)
			updateQuery += fmt.Sprintf("%v=$%v", fieldNameUnderScore, fieldCount)
			if fieldsLength > 1 && fieldCount < fieldsLength {
				updateQuery += ", "
			}
		}
		// add where condition by id and the placeholder-value position
		fieldCount += 1
		updateQuery += fmt.Sprintf(" WHERE id=$%v", fieldCount)
		updateQuery += " RETURNING id"
		// add id-placeholder-value
		fieldValues = append(fieldValues, recordId)
		// update result
		updateQueryObjects = append(updateQueryObjects, UpdateQueryObject{
			UpdateQuery: updateQuery,
			FieldNames:  fieldNames,
			FieldValues: fieldValues,
		})
	}

	// result
	return MultiUpdateQueryResult{
		UpdateQueryObjects: updateQueryObjects,
		Ok:                 true,
		Message:            "success",
	}
}

// ComputeUpdateQueryById function computes update SQL scripts by recordId. It returns updateScript, updateValues []interface{} and/or err error
func ComputeUpdateQueryById(tableName string, actionParam ActionParamType, recordId string) UpdateQueryResult {
	if tableName == "" || len(actionParam) < 1 || actionParam == nil || recordId == "" {
		return updateErrMessage("table-name, recordId and actionParam are required for the update operation")
	}
	// compute update script and associated place-holder values for the actionParam/record
	updateQuery := fmt.Sprintf("UPDATE %v SET ", tableName)
	var fieldValues []interface{}
	var fieldNames []string
	var fieldNamesUnderscore []string
	fieldsLength := len(actionParam)
	fieldCount := 0
	for fieldName, fieldValue := range actionParam {
		// next placeholder-value-position
		fieldCount += 1
		// skip fieldName=="id"
		if fieldName == "id" {
			fieldCount -= 1
			fieldsLength -= 1
			continue
		}
		fieldNameUnderScore := govalidator.CamelCaseToUnderscore(fieldName)
		fieldNames = append(fieldNames, fieldName)
		fieldNamesUnderscore = append(fieldNamesUnderscore, fieldNameUnderScore)
		// TODO: update fieldValues by fieldValue-type, for correct postgres-SQL-parsing
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
		default:
			currentFieldValue = fieldValue
		}

		fieldValues = append(fieldValues, currentFieldValue)
		updateQuery += fmt.Sprintf("%v=$%v", fieldNameUnderScore, fieldCount)
		if fieldsLength > 1 && fieldCount < fieldsLength {
			updateQuery += ", "
		}
	}
	// add where condition by id and the placeholder-value position
	fieldCount += 1
	updateQuery += fmt.Sprintf(" WHERE id=$%v", fieldCount)
	updateQuery += " RETURNING id"
	// add id-placeholder-value
	fieldValues = append(fieldValues, recordId)

	// result
	return UpdateQueryResult{
		UpdateQueryObject: UpdateQueryObject{
			UpdateQuery: updateQuery,
			FieldNames:  fieldNames,
			FieldValues: fieldValues,
		},
		Ok:      true,
		Message: "success",
	}
}

// ComputeUpdateQueryByIds function computes update SQL scripts by recordIds. It returns updateScript, updateValues []interface{} and/or err error
func ComputeUpdateQueryByIds(tableName string, actionParam ActionParamType, recordIds []string) UpdateQueryResult {
	if tableName == "" || len(actionParam) < 1 || actionParam == nil || len(recordIds) < 1 {
		return updateErrMessage("tableName, recordIds and actionParam are required for the update operation")
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
	// compute update script and associated place-holder values for the actionParam/record
	updateQuery := fmt.Sprintf("UPDATE %v SET ", tableName)
	var fieldValues []interface{}
	var fieldNames []string
	var fieldNamesUnderscore []string
	fieldsLength := len(actionParam)
	fieldCount := 0
	for fieldName, fieldValue := range actionParam {
		// next placeholder-value-position
		fieldCount += 1
		// skip fieldName=="id"
		if fieldName == "id" {
			fieldCount -= 1
			fieldsLength -= 1
			continue
		}
		fieldNameUnderScore := govalidator.CamelCaseToUnderscore(fieldName)
		fieldNames = append(fieldNames, fieldName)
		fieldNamesUnderscore = append(fieldNamesUnderscore, fieldNameUnderScore)
		// TODO: update fieldValues by fieldValue-type, for correct postgres-SQL-parsing
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
		default:
			currentFieldValue = fieldValue
		}
		fieldValues = append(fieldValues, currentFieldValue)
		updateQuery += fmt.Sprintf("%v=$%v", fieldNameUnderScore, fieldCount)
		if fieldsLength > 1 && fieldCount < fieldsLength {
			updateQuery += ", "
		}
	}
	// add where condition by id and the placeholder-value position
	updateQuery += whereQuery

	// result
	return UpdateQueryResult{
		UpdateQueryObject: UpdateQueryObject{
			UpdateQuery: updateQuery,
			FieldNames:  fieldNames,
			FieldValues: fieldValues,
		},
		Ok:      true,
		Message: "success",
	}
}

// ComputeUpdateQueryByParam function computes update SQL scripts by queryParams. It returns updateScript, updateValues []interface{} and/or err error
func ComputeUpdateQueryByParam(tableName string, actionParam ActionParamType, queryParam QueryParamType) UpdateQueryResult {
	if tableName == "" || len(actionParam) < 1 || actionParam == nil || len(queryParam) < 1 {
		return updateErrMessage("table-name, queryParam and actionParam are required for the update operation")
	}
	// compute update script and associated place-holder values for the actionParam/record
	updateQuery := fmt.Sprintf("UPDATE %v SET ", tableName)
	var fieldValues []interface{}
	var fieldNames []string
	var fieldNamesUnderscore []string
	fieldsLength := len(actionParam)
	fieldCount := 0
	for fieldName, fieldValue := range actionParam {
		// skip fieldName=="id"
		if fieldName == "id" {
			fieldsLength -= 1
			continue
		}
		// next placeholder-value-position
		fieldCount += 1
		fieldNameUnderScore := govalidator.CamelCaseToUnderscore(fieldName)
		fieldNames = append(fieldNames, fieldName)
		fieldNamesUnderscore = append(fieldNamesUnderscore, fieldNameUnderScore)
		// TODO: update fieldValues by fieldValue-type, for correct postgres-SQL-parsing
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
		default:
			currentFieldValue = fieldValue
		}

		fieldValues = append(fieldValues, currentFieldValue)
		updateQuery += fmt.Sprintf("%v=$%v", fieldNameUnderScore, fieldCount)
		if fieldsLength > 1 && fieldCount < fieldsLength {
			updateQuery += ", "
		}
	}
	// where-query
	whereRes := ComputeWhereQuery(queryParam, fieldCount+1)
	if !whereRes.Ok {
		return updateErrMessage(fmt.Sprintf("error computing where-query condition(s): %v", whereRes.Message))
	}

	updateQuery += " " + whereRes.WhereQueryObject.WhereQuery

	// result
	return UpdateQueryResult{
		UpdateQueryObject: UpdateQueryObject{
			UpdateQuery: updateQuery,
			FieldNames:  fieldNames,
			FieldValues: append(fieldValues, whereRes.WhereQueryObject.FieldValues...),
		},
		Ok:      true,
		Message: "success",
	}
}
