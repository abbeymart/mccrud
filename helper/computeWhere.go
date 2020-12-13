// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: compute where-SQL script

package helper

import (
	"errors"
	"fmt"
	"github.com/abbeymart/mccrud"
	"sort"
	"strings"
	"time"
)

func whereScriptErr(errMsg string) (mccrud.WhereScriptResponseType, error) {
	return mccrud.WhereScriptResponseType{
		WhereScript: nil,
		FieldValues: nil,
	}, errors.New(errMsg)
}

func ComputeWhereQuery(where mccrud.WhereParamType, tableFields []string) (string, error) {
	// groups length/size
	groupsLen := len(where)
	if len(where) < 1 || len(tableFields) < 1 {
		return "", errors.New("where condition and table-field names are required")
	}
	// compute where script from where
	// initialize group validation variable | to determine group with empty/no fieldItems
	unspecifiedGroupCount := 0
	// sort where by groupOrder (ASC)
	sort.SliceStable(where, func(i, j int) bool {
		return where[i].GroupOrder < where[j].GroupOrder
	})
	// variable for valid group count, i.e. group with groupItems
	groupCount := 0
	// iterate through where (groups)
	whereQuery := " WHERE "
	for _, group := range where {
		var (
			unspecifiedGroupItemCount = 0 // variable to determine unspecified/invalid fieldName or fieldValue
			groupItemCount            = 0 // valid groupItem count, i.e. group item with valid name and value
		)
		groupItemsLen := len(group.GroupItems) // total items in a group
		// check groupItems length, if 0 continue to the next group
		if groupItemsLen < 1 {
			unspecifiedGroupCount += 1
			continue
		}
		// count valid group, i.e. group with groupItems
		groupCount += 1
		// sort group items by groupItem/fieldOrder (ASC)
		gItems := group.GroupItems
		sort.SliceStable(gItems, func(i, j int) bool {
			return gItems[i].GroupItemOrder < gItems[j].GroupItemOrder
		})
		// compute the field-where-script
		fieldQuery := " ("
		for _, groupItem := range gItems {
			// check groupItem's fieldName, fieldOperator and fieldValue
			fieldName := ""
			fieldOperator := ""
			var fieldValue interface{}

			for fName, opVal := range groupItem.GroupItem {
				fieldName = fName
				for fOp, val := range opVal {
					fieldOperator = fOp
					fieldValue = val
				}
			}
			// validate fieldName
			if !ArrayStringContains(tableFields, fieldName) {
				return "", errors.New(fmt.Sprintf("invalid field name [%v] specified in where condition", fieldName))
			}
			if fieldName == "" || fieldOperator == "" || fieldValue == nil {
				// skip missing field or return error
				unspecifiedGroupItemCount += 1
				continue
				//return "", errors.New("field-name, operator and/or value are required")
			}
			// count valid groupItem | TODO: include date-type
			groupItemCount += 1
			switch fieldOperator {
			case mccrud.FieldOperators().Equals, strings.ToLower(mccrud.FieldOperators().Equals):
				switch fieldValue.(type) {
				case string:
					if fVal, ok := fieldValue.(string); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v=%v", fieldName, fVal)
					}
				case time.Time:
					if fVal, ok := fieldValue.(time.Time); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v=%v", fieldName, fVal)
					}
				case bool:
					if fVal, ok := fieldValue.(bool); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v=%v", fieldName, fVal)
					}
				case int8:
					if fVal, ok := fieldValue.(int8); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v=%v", fieldName, fVal)
					}
				case int16:
					if fVal, ok := fieldValue.(int16); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v=%v", fieldName, fVal)
					}
				case int32:
					if fVal, ok := fieldValue.(int32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v=%v", fieldName, fVal)
					}
				case int64:
					if fVal, ok := fieldValue.(int64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v=%v", fieldName, fVal)
					}
				case int:
					if fVal, ok := fieldValue.(int); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v=%v", fieldName, fVal)
					}
				case uint8:
					if fVal, ok := fieldValue.(uint8); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v=%v", fieldName, fVal)
					}
				case uint16:
					if fVal, ok := fieldValue.(uint16); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v=%v", fieldName, fVal)
					}
				case uint32:
					if fVal, ok := fieldValue.(uint32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v=%v", fieldName, fVal)
					}
				case uint64:
					if fVal, ok := fieldValue.(uint64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v=%v", fieldName, fVal)
					}
				case uint:
					if fVal, ok := fieldValue.(uint); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v=%v", fieldName, fVal)
					}
				case float32:
					if fVal, ok := fieldValue.(float32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v=%v", fieldName, fVal)
					}
				case float64:
					if fVal, ok := fieldValue.(float64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v=%v", fieldName, fVal)
					}
				case []string:
					if fVal, ok := fieldValue.([]string); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v=%v", fieldName, fVal)
					}
				case []int:
					if fVal, ok := fieldValue.([]int); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v=%v", fieldName, fVal)
					}
				case []struct{}:
					if fVal, ok := fieldValue.([]struct{}); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v=%v", fieldName, fVal)
					}
				default:
					return "", errors.New(fmt.Sprintf("Unsupported field[%v], format for field-value %v", fieldName, fieldValue))
				}
				groupItemOps := []string{"and", "or"}
				groupItemOp := groupItem.GroupItemOp
				if groupItemOp != "" && !ArrayStringContains(groupItemOps, strings.ToLower(groupItemOp)) {
					groupItemOp = mccrud.GroupOperators().AND // use GroupOpTypes.AND as default operator
				}
				if groupItemsLen > 1 && groupItemCount < (groupItemsLen - unspecifiedGroupItemCount){
					if groupItem.GroupItemOp == mccrud.GroupOperators().AND {
						fieldQuery = fieldQuery + " " + groupItem.GroupItemOp + " "
					}
				}
			case mccrud.FieldOperators().NotEquals, strings.ToLower(mccrud.FieldOperators().NotEquals):
				switch fieldValue.(type) {
				case string:
					if fVal, ok := fieldValue.(string); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<>%v", fieldName, fVal)
					}
				case time.Time:
					if fVal, ok := fieldValue.(time.Time); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<>%v", fieldName, fVal)
					}
				case bool:
					if fVal, ok := fieldValue.(bool); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<>%v", fieldName, fVal)
					}
				case int8:
					if fVal, ok := fieldValue.(int8); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<>%v", fieldName, fVal)
					}
				case int16:
					if fVal, ok := fieldValue.(int16); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<>%v", fieldName, fVal)
					}
				case int32:
					if fVal, ok := fieldValue.(int32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<>%v", fieldName, fVal)
					}
				case int64:
					if fVal, ok := fieldValue.(int64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<>%v", fieldName, fVal)
					}
				case int:
					if fVal, ok := fieldValue.(int); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<>%v", fieldName, fVal)
					}
				case uint8:
					if fVal, ok := fieldValue.(uint8); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<>%v", fieldName, fVal)
					}
				case uint16:
					if fVal, ok := fieldValue.(uint16); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<>%v", fieldName, fVal)
					}
				case uint32:
					if fVal, ok := fieldValue.(uint32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<>%v", fieldName, fVal)
					}
				case uint64:
					if fVal, ok := fieldValue.(uint64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<>%v", fieldName, fVal)
					}
				case uint:
					if fVal, ok := fieldValue.(uint); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<>%v", fieldName, fVal)
					}
				case float32:
					if fVal, ok := fieldValue.(float32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<>%v", fieldName, fVal)
					}
				case float64:
					if fVal, ok := fieldValue.(float64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<>%v", fieldName, fVal)
					}
				case []string:
					if fVal, ok := fieldValue.([]string); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<>%v", fieldName, fVal)
					}
				case []int:
					if fVal, ok := fieldValue.([]int); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<>%v", fieldName, fVal)
					}
				case []struct{}:
					if fVal, ok := fieldValue.([]struct{}); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<>%v", fieldName, fVal)
					}
				default:
					return "", errors.New(fmt.Sprintf("Unsupported field[%v], format for field-value %v", fieldName, fieldValue))
				}
				groupItemOps := []string{"and", "or"}
				groupItemOp := groupItem.GroupItemOp
				if groupItemOp != "" && !ArrayStringContains(groupItemOps, strings.ToLower(groupItemOp)) {
					groupItemOp = mccrud.GroupOperators().AND // use GroupOpTypes.AND as default operator
				}
				if groupItemsLen > 1 && groupItemCount < (groupItemsLen - unspecifiedGroupItemCount){
					if groupItem.GroupItemOp == mccrud.GroupOperators().AND {
						fieldQuery = fieldQuery + " " + groupItem.GroupItemOp + " "
					}
				}
			case mccrud.FieldOperators().LessThan, strings.ToLower(mccrud.FieldOperators().LessThan):
				switch fieldValue.(type) {
				case string:
					if fVal, ok := fieldValue.(string); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<%v", fieldName, fVal)
					}
				case time.Time:
					if fVal, ok := fieldValue.(time.Time); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<%v", fieldName, fVal)
					}
				case bool:
					if fVal, ok := fieldValue.(bool); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<%v", fieldName, fVal)
					}
				case int8:
					if fVal, ok := fieldValue.(int8); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<%v", fieldName, fVal)
					}
				case int16:
					if fVal, ok := fieldValue.(int16); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<%v", fieldName, fVal)
					}
				case int32:
					if fVal, ok := fieldValue.(int32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<%v", fieldName, fVal)
					}
				case int64:
					if fVal, ok := fieldValue.(int64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<%v", fieldName, fVal)
					}
				case int:
					if fVal, ok := fieldValue.(int); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<%v", fieldName, fVal)
					}
				case uint8:
					if fVal, ok := fieldValue.(uint8); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<%v", fieldName, fVal)
					}
				case uint16:
					if fVal, ok := fieldValue.(uint16); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<%v", fieldName, fVal)
					}
				case uint32:
					if fVal, ok := fieldValue.(uint32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<%v", fieldName, fVal)
					}
				case uint64:
					if fVal, ok := fieldValue.(uint64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<%v", fieldName, fVal)
					}
				case uint:
					if fVal, ok := fieldValue.(uint); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<%v", fieldName, fVal)
					}
				case float32:
					if fVal, ok := fieldValue.(float32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<%v", fieldName, fVal)
					}
				case float64:
					if fVal, ok := fieldValue.(float64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<%v", fieldName, fVal)
					}
				case []string:
					if fVal, ok := fieldValue.([]string); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<%v", fieldName, fVal)
					}
				case []int:
					if fVal, ok := fieldValue.([]int); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<%v", fieldName, fVal)
					}
				case []struct{}:
					if fVal, ok := fieldValue.([]struct{}); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<%v", fieldName, fVal)
					}
				default:
					return "", errors.New(fmt.Sprintf("Unsupported field: %v, format for field-value: %v", fieldName, fieldValue))
				}
				groupItemOps := []string{"and", "or"}
				groupItemOp := groupItem.GroupItemOp
				if groupItemOp != "" && !ArrayStringContains(groupItemOps, strings.ToLower(groupItemOp)) {
					groupItemOp = mccrud.GroupOperators().AND // use GroupOpTypes.AND as default operator
				}
				if groupItemsLen > 1 && groupItemCount < (groupItemsLen - unspecifiedGroupItemCount){
					if groupItem.GroupItemOp == mccrud.GroupOperators().AND {
						fieldQuery = fieldQuery + " " + groupItem.GroupItemOp + " "
					}
				}
			case mccrud.FieldOperators().LessThanOrEquals, strings.ToLower(mccrud.FieldOperators().LessThanOrEquals):
				switch fieldValue.(type) {
				case string:
					if fVal, ok := fieldValue.(string); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<=%v", fieldName, fVal)
					}
				case time.Time:
					if fVal, ok := fieldValue.(time.Time); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<=%v", fieldName, fVal)
					}
				case bool:
					if fVal, ok := fieldValue.(bool); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<=%v", fieldName, fVal)
					}
				case int8:
					if fVal, ok := fieldValue.(int8); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<=%v", fieldName, fVal)
					}
				case int16:
					if fVal, ok := fieldValue.(int16); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<=%v", fieldName, fVal)
					}
				case int32:
					if fVal, ok := fieldValue.(int32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<=%v", fieldName, fVal)
					}
				case int64:
					if fVal, ok := fieldValue.(int64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<=%v", fieldName, fVal)
					}
				case int:
					if fVal, ok := fieldValue.(int); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<=%v", fieldName, fVal)
					}
				case uint8:
					if fVal, ok := fieldValue.(uint8); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<=%v", fieldName, fVal)
					}
				case uint16:
					if fVal, ok := fieldValue.(uint16); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<=%v", fieldName, fVal)
					}
				case uint32:
					if fVal, ok := fieldValue.(uint32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<=%v", fieldName, fVal)
					}
				case uint64:
					if fVal, ok := fieldValue.(uint64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<=%v", fieldName, fVal)
					}
				case uint:
					if fVal, ok := fieldValue.(uint); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<=%v", fieldName, fVal)
					}
				case float32:
					if fVal, ok := fieldValue.(float32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<=%v", fieldName, fVal)
					}
				case float64:
					if fVal, ok := fieldValue.(float64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<=%v", fieldName, fVal)
					}
				case []string:
					if fVal, ok := fieldValue.([]string); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<=%v", fieldName, fVal)
					}
				case []int:
					if fVal, ok := fieldValue.([]int); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<=%v", fieldName, fVal)
					}
				case []struct{}:
					if fVal, ok := fieldValue.([]struct{}); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v<=%v", fieldName, fVal)
					}
				default:
					return "", errors.New(fmt.Sprintf("Unsupported field[%v], format for field-value %v", fieldName, fieldValue))
				}
				groupItemOps := []string{"and", "or"}
				groupItemOp := groupItem.GroupItemOp
				if groupItemOp != "" && !ArrayStringContains(groupItemOps, strings.ToLower(groupItemOp)) {
					groupItemOp = mccrud.GroupOperators().AND // use GroupOpTypes.AND as default operator
				}
				if groupItemsLen > 1 && groupItemCount < (groupItemsLen - unspecifiedGroupItemCount){
					if groupItem.GroupItemOp == mccrud.GroupOperators().AND {
						fieldQuery = fieldQuery + " " + groupItem.GroupItemOp + " "
					}
				}
			case mccrud.FieldOperators().GreaterThan, strings.ToLower(mccrud.FieldOperators().GreaterThan):
				switch fieldValue.(type) {
				case string:
					if fVal, ok := fieldValue.(string); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>%v", fieldName, fVal)
					}
				case time.Time:
					if fVal, ok := fieldValue.(time.Time); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>%v", fieldName, fVal)
					}
				case bool:
					if fVal, ok := fieldValue.(bool); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>%v", fieldName, fVal)
					}
				case int8:
					if fVal, ok := fieldValue.(int8); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>%v", fieldName, fVal)
					}
				case int16:
					if fVal, ok := fieldValue.(int16); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>%v", fieldName, fVal)
					}
				case int32:
					if fVal, ok := fieldValue.(int32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>%v", fieldName, fVal)
					}
				case int64:
					if fVal, ok := fieldValue.(int64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>%v", fieldName, fVal)
					}
				case int:
					if fVal, ok := fieldValue.(int); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>%v", fieldName, fVal)
					}
				case uint8:
					if fVal, ok := fieldValue.(uint8); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>%v", fieldName, fVal)
					}
				case uint16:
					if fVal, ok := fieldValue.(uint16); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>%v", fieldName, fVal)
					}
				case uint32:
					if fVal, ok := fieldValue.(uint32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>%v", fieldName, fVal)
					}
				case uint64:
					if fVal, ok := fieldValue.(uint64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>%v", fieldName, fVal)
					}
				case uint:
					if fVal, ok := fieldValue.(uint); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>%v", fieldName, fVal)
					}
				case float32:
					if fVal, ok := fieldValue.(float32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>%v", fieldName, fVal)
					}
				case float64:
					if fVal, ok := fieldValue.(float64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>%v", fieldName, fVal)
					}
				case []string:
					if fVal, ok := fieldValue.([]string); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>%v", fieldName, fVal)
					}
				case []int:
					if fVal, ok := fieldValue.([]int); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>%v", fieldName, fVal)
					}
				case []struct{}:
					if fVal, ok := fieldValue.([]struct{}); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>%v", fieldName, fVal)
					}
				default:
					return "", errors.New(fmt.Sprintf("Unsupported field[%v], format for field-value %v", fieldName, fieldValue))
				}
				groupItemOps := []string{"and", "or"}
				groupItemOp := groupItem.GroupItemOp
				if groupItemOp != "" && !ArrayStringContains(groupItemOps, strings.ToLower(groupItemOp)) {
					groupItemOp = mccrud.GroupOperators().AND // use GroupOpTypes.AND as default operator
				}
				if groupItemsLen > 1 && groupItemCount < (groupItemsLen - unspecifiedGroupItemCount){
					if groupItem.GroupItemOp == mccrud.GroupOperators().AND {
						fieldQuery = fieldQuery + " " + groupItem.GroupItemOp + " "
					}
				}
			case mccrud.FieldOperators().GreaterThanOrEquals, strings.ToLower(mccrud.FieldOperators().GreaterThanOrEquals):
				switch fieldValue.(type) {
				case string:
					if fVal, ok := fieldValue.(string); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>=%v", fieldName, fVal)
					}
				case time.Time:
					if fVal, ok := fieldValue.(time.Time); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>=%v", fieldName, fVal)
					}
				case bool:
					if fVal, ok := fieldValue.(bool); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>=%v", fieldName, fVal)
					}
				case int8:
					if fVal, ok := fieldValue.(int8); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>=%v", fieldName, fVal)
					}
				case int16:
					if fVal, ok := fieldValue.(int16); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>=%v", fieldName, fVal)
					}
				case int32:
					if fVal, ok := fieldValue.(int32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>=%v", fieldName, fVal)
					}
				case int64:
					if fVal, ok := fieldValue.(int64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>=%v", fieldName, fVal)
					}
				case int:
					if fVal, ok := fieldValue.(int); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>=%v", fieldName, fVal)
					}
				case uint8:
					if fVal, ok := fieldValue.(uint8); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>=%v", fieldName, fVal)
					}
				case uint16:
					if fVal, ok := fieldValue.(uint16); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>=%v", fieldName, fVal)
					}
				case uint32:
					if fVal, ok := fieldValue.(uint32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>=%v", fieldName, fVal)
					}
				case uint64:
					if fVal, ok := fieldValue.(uint64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>=%v", fieldName, fVal)
					}
				case uint:
					if fVal, ok := fieldValue.(uint); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>=%v", fieldName, fVal)
					}
				case float32:
					if fVal, ok := fieldValue.(float32); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>=%v", fieldName, fVal)
					}
				case float64:
					if fVal, ok := fieldValue.(float64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>=%v", fieldName, fVal)
					}
				case []string:
					if fVal, ok := fieldValue.([]string); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>=%v", fieldName, fVal)
					}
				case []int:
					if fVal, ok := fieldValue.([]int); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>=%v", fieldName, fVal)
					}
				case []struct{}:
					if fVal, ok := fieldValue.([]struct{}); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v>=%v", fieldName, fVal)
					}
				default:
					return "", errors.New(fmt.Sprintf("Unsupported field[%v], format for field-value %v", fieldName, fieldValue))
				}
				groupItemOps := []string{"and", "or"}
				groupItemOp := groupItem.GroupItemOp
				if groupItemOp != "" && !ArrayStringContains(groupItemOps, strings.ToLower(groupItemOp)) {
					groupItemOp = mccrud.GroupOperators().AND // use GroupOpTypes.AND as default operator
				}
				if groupItemsLen > 1 && groupItemCount < (groupItemsLen - unspecifiedGroupItemCount){
					if groupItem.GroupItemOp == mccrud.GroupOperators().AND {
						fieldQuery = fieldQuery + " " + groupItem.GroupItemOp + " "
					}
				}
			case mccrud.FieldOperators().In, strings.ToLower(mccrud.FieldOperators().In):
				switch fieldValue.(type) {
				case []string:
					if fVal, ok := fieldValue.([]string); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v=%v", fieldName, fVal)
					}
				case []int:
					if fVal, ok := fieldValue.([]int); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v=%v", fieldName, fVal)
					}
				case []float32, []float64:
					if fVal, ok := fieldValue.([]float64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v=%v", fieldName, fVal)
					}
				default:
					return "", errors.New(fmt.Sprintf("Unsupported field[%v], format for field-value %v", fieldName, fieldValue))
				}
				groupItemOps := []string{"and", "or"}
				groupItemOp := groupItem.GroupItemOp
				if groupItemOp != "" && !ArrayStringContains(groupItemOps, strings.ToLower(groupItemOp)) {
					groupItemOp = mccrud.GroupOperators().AND // use GroupOpTypes.AND as default operator
				}
				if groupItemsLen > 1 && groupItemCount < (groupItemsLen - unspecifiedGroupItemCount){
					if groupItem.GroupItemOp == mccrud.GroupOperators().AND {
						fieldQuery = fieldQuery + " " + groupItem.GroupItemOp + " "
					}
				}
			case mccrud.FieldOperators().NotIn, strings.ToLower(mccrud.FieldOperators().NotIn):
				switch fieldValue.(type) {
				case []string:
					if fVal, ok := fieldValue.([]string); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v=%v", fieldName, fVal)
					}
				case []int:
					if fVal, ok := fieldValue.([]int); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v=%v", fieldName, fVal)
					}
				case []float32, []float64:
					if fVal, ok := fieldValue.([]float64); !ok {
						return "", errors.New(fmt.Sprintf("field_name: %v | field_value: %v error: ", fieldName, fieldValue))
					} else {
						fieldQuery += fmt.Sprintf(" %v=%v", fieldName, fVal)
					}
				default:
					return "", errors.New(fmt.Sprintf("Unsupported field[%v], format for field-value %v", fieldName, fieldValue))
				}
				groupItemOps := []string{"and", "or"}
				groupItemOp := groupItem.GroupItemOp
				if groupItemOp != "" && !ArrayStringContains(groupItemOps, strings.ToLower(groupItemOp)) {
					groupItemOp = mccrud.GroupOperators().AND // use GroupOpTypes.AND as default operator
				}
				if groupItemsLen > 1 && groupItemCount < (groupItemsLen - unspecifiedGroupItemCount){
					if groupItem.GroupItemOp == mccrud.GroupOperators().AND {
						fieldQuery = fieldQuery + " " + groupItem.GroupItemOp + " "
					}
				}
			default:
				return "", errors.New(fmt.Sprintf("Unknown or unsupported field(%v) operator: %v", fieldName, fieldOperator))
			}
			// continue to the next group iteration, if fieldItems is empty for the current group
			if unspecifiedGroupItemCount == groupItemsLen {
				continue
			}
			// add closing bracket to complete the group-items query/script
			fieldQuery = fieldQuery + " ) "
			//validate acceptable groupLinkOperators (and || or)
			grpLinkOp := group.GroupLinkOp
			groupLnOps := []string{"and", "or"}
			if grpLinkOp != "" && !ArrayStringContains(groupLnOps, strings.ToLower(grpLinkOp)) {
				grpLinkOp = mccrud.GroupOperators().AND // use GroupOpTypes.AND as default operator
			}
			// add groupLinkOp, if groupsLen > 1
			if groupsLen > 1 && groupCount < (groupsLen - unspecifiedGroupCount)  {
				fieldQuery = fieldQuery + " " + strings.ToUpper(grpLinkOp) + " "
			}
			// compute where-script from the group-script, append in sequence by groupOrder
			whereQuery = whereQuery + " " + fieldQuery
		}
		// check WHERE script contains at least one condition, otherwise raise an exception
		if unspecifiedGroupCount == groupsLen {
			return "", errors.New("no valid where condition specified")
		}
	}
	// if all went well, return valid where script
	return whereQuery, nil
}
