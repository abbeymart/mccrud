// @Author: abbeymart | Abi Akindele | @Created: 2020-12-01 | @Updated: 2020-12-01
// @Company: mConnect.biz | @License: MIT
// @Description: get / query record(s)

package mccrud

import (
	"fmt"
	"github.com/abbeymart/mccrud/helper"
	"github.com/abbeymart/mcresponse"
)

func (crud Crud) GetAll() mcresponse.ResponseMessage {
	var tableFields []string
	// compose tableFields
	if tFields, err := helper.ComputeGetFields(crud.ProjectParams); err == nil {
		tableFields = tFields
	}
	if queryRes, err := helper.ComputeSelectQueryAll(crud.TableName, tableFields); err != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing select/read-query: %v", err.Error()),
			Value:   queryRes,
		})
	} else {
		// TODO: perform crud-task action, include options (skip, limit, sort etc.):

	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "success",
		Value:   nil,
	})
}

func (crud Crud) GetById() mcresponse.ResponseMessage {
	var tableFields []string
	// compose tableFields
	if tFields, err := helper.ComputeGetFields(crud.ProjectParams); err == nil {
		tableFields = tFields
	}
	if queryRes, err := helper.ComputeSelectQueryById(crud.TableName, tableFields, crud.RecordIds); err != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing select/read-query: %v", err.Error()),
			Value:   queryRes,
		})
	} else {
		// TODO: perform crud-task action, include options (skip, limit, sort etc.):
	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "success",
		Value:   nil,
	})
}

func (crud Crud) GetByParam() mcresponse.ResponseMessage {
	var tableFields []string
	// compose tableFields
	if tFields, err := helper.ComputeGetFields(crud.ProjectParams); err == nil {
		tableFields = tFields
	}
	if queryRes, err := helper.ComputeSelectQueryByParam(crud.TableName, tableFields, crud.QueryParams); err != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing select/read-query: %v", err.Error()),
			Value:   queryRes,
		})
	} else {
		// TODO: perform crud-task action, include options (skip, limit, sort etc.):

	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "success",
		Value:   nil,
	})
}
