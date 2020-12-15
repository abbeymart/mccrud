// @Author: abbeymart | Abi Akindele | @Created: 2020-12-15 | @Updated: 2020-12-15
// @Company: mConnect.biz | @License: MIT
// @Description: go: mConnect

package helper

import (
	"errors"
	"github.com/abbeymart/mccrud"
)

func ComputeTableFields(actionParams mccrud.ActionParamsType, projectParams mccrud.ProjectParamType) ([]string, error) {
	if len(actionParams) < 1 {
		return nil, errors.New("actionParams is required")
	}
	// obtain tableFields from api consumer (ProjectParams)
	var tableFields []string
	if len(projectParams) > 0 {
		for fieldName, ok := range projectParams {
			if ok {
				tableFields = append(tableFields, fieldName)
			}
		}
		// include default fields (id) for select-query only
		//if !ArrayStringContains(tableFields, "id") {
		//	tableFields = append(tableFields, "id")
		//}
	}
	if len(tableFields) < 1 {
		// obtain tableFields from actionParams[0]
		for fieldName := range actionParams[0] {
			tableFields = append(tableFields, fieldName)
		}
	}

	return tableFields, nil
}
