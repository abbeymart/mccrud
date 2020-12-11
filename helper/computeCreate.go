// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: compute create-SQL script

package helper

import (
	"errors"
	"fmt"
	"github.com/abbeymart/mccrudgo"
)

// ComputeCreateScript computes insert SQL script. It returns createScripts []string, fieldNames []string and err error
func ComputeCreateScript(tableName string, actionParams mccrud.ActionParamsType) (createScript []string, fieldNames []string, err error) {
	var insertScripts []string
	var fNames []string	// fieldNames array of strings in order of SQL statement

	if tableName == "" || len(actionParams) < 1 {
		return nil, nil, errors.New("table/collection name and action-params are required for the create operation")
	}
	// compute fieldNames, from one of the actionParams items/records
	for fName := range actionParams[0] {
		fNames = append(fNames, fName)
	}

	// computed create script from actionParams
	for _, val := range actionParams {
		var itemScript = "INSERT INTO " + tableName + " ("
		var itemValues = " VALUES("
		var (
			fieldCount = 0
		)
		fieldsLength := len(val)
		for fieldName, fieldValue := range val {
			if fieldsLength > 1 && fieldCount < fieldsLength - 1 && fieldCount > 0 {
				itemScript += ", " + fieldName
				itemValues += ", " + fmt.Sprintf("%v", fieldValue)
			} else {
				itemScript += " " + fieldName
				itemValues += " " + fmt.Sprintf("%v", fieldValue)
			}
			fieldCount += 1
		}
		// close item-script
		itemScript += " )"
		itemValues += " )"

		// add item-script to the createScripts
		insertScripts = append(insertScripts, itemScript + itemValues)
	}

	return insertScripts, fNames, nil
}
