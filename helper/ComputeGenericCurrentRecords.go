// @Author: abbeymart | Abi Akindele | @Created: 2020-12-16 | @Updated: 2020-12-16
// @Company: mConnect.biz | @License: MIT
// @Description: go: mConnect

package helper

import (
	"encoding/json"
	"errors"
	"fmt"
)

func ParseRawValues(rawValues [][]byte) ([]interface{}, error) {
	// variables
	var v interface{}
	var va []interface{}
	// parse the current-raw-values
	for _, val := range rawValues {
		if err := json.Unmarshal(val, &v); err != nil {
			return nil, errors.New(fmt.Sprintf("Error parsing raw-row-value: %v", err.Error()))
		} else {
			va = append(va, v)
		}
	}
	return va, nil
}
