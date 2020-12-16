// @Author: abbeymart | Abi Akindele | @Created: 2020-12-01 | @Updated: 2020-12-01
// @Company: mConnect.biz | @License: MIT
// @Description: delete or remove record(s)

package mccrud

import (
	"github.com/abbeymart/mccrud/helper"
	"github.com/abbeymart/mcresponse"
)

func (crud Crud) DeleteById() mcresponse.ResponseMessage {
	// TODO: get current records, for audit-log | for delete tableField = []string{}
	if getQuery, err := helper.ComputeSelectQueryById(crud.TableName, []string{}, crud.RecordIds); err != nil {
		// exit on error

	} else {
		// exit of currentRec-length is less than recordIds-length


	}
	// TODO: perform crud-task action, include where-
	// compute delete script from where

	// perform audit-log


	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "success",
		Value:   nil,
	})
}

func (crud Crud) DeleteByParam() mcresponse.ResponseMessage {
	// TODO: get current records, for audit-log | for delete tableField = []string{}
	if getQuery, err := helper.ComputeSelectQueryByParam(crud.TableName, []string{}, crud.QueryParams); err != nil {
		// exit on error

	} else {
		// exit of currentRec-length is less than 1

	}
	// TODO: perform crud-task action, include where-query(params):
	// compute delete script from where

	// perform audit-log

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "success",
		Value:   nil,
	})
}

// DeleteAll function removes all records in the tables. Recommended for admin-users only
// Use if and only if you know what you are doing
func (crud Crud) DeleteAll() mcresponse.ResponseMessage {
	// TODO: perform crud-task action:
	// compute delete script from where

	// perform audit-log | record only the action performed, excluding table-records

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "success",
		Value:   nil,
	})
}