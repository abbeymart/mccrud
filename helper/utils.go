// @Author: abbeymart | Abi Akindele | @Created: 2020-12-09 | @Updated: 2020-12-09
// @Company: mConnect.biz | @License: MIT
// @Description: crud utility functions

package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/abbeymart/mctypes"
	"github.com/asaskevich/govalidator"
	"reflect"
)

func ArrayStringContains(arr []string, val string) bool {
	for _, a := range arr {
		if a == val {
			return true
		}
	}
	return false
}

func ArrayIntContains(arr []int, val int) bool {
	for _, a := range arr {
		if a == val {
			return true
		}
	}
	return false
}

func ArraySQLInStringValues(arr []string) string {
	result := ""
	for ind, val := range arr {
		result += "'" + val + "'"
		if ind < len(arr) - 1 {
			result += ", "
		}
	}
	return result
}

// JsonDataETL method converts json inputs to equivalent struct data type specification
// rec must be a pointer to a type matching the jsonRec
func JsonDataETL(jsonRec []byte, rec interface{}) error {
	if err := json.Unmarshal(jsonRec, &rec); err == nil {
		return nil
	} else {
		return errors.New(fmt.Sprintf("Error converting json-to-record-format: %v", err.Error()))
	}
}

// DataToValueParam method accepts only a struct record/param (type/model) and returns the ValueParamType
// data camel/Pascal-case keys are converted to underscore-keys to match table-field/columns specs
func DataToValueParam(rec interface{}) (mctypes.ValueParamType, error) {
	switch rec.(type) {
	case struct{}:
		dataValue := mctypes.ValueParamType{}
		v := reflect.ValueOf(rec)
		typeOfS := v.Type()

		for i := 0; i < v.NumField(); i++ {
			dataValue[govalidator.CamelCaseToUnderscore(typeOfS.Field(i).Name)] = v.Field(i).Interface()
			//fmt.Printf("Field: %s\tValue: %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
		}
		return dataValue, nil
	default:
		return nil, errors.New("invalid type - requires parameter of type struct only")
	}
}
