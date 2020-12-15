// @Author: abbeymart | Abi Akindele | @Created: 2020-12-01 | @Updated: 2020-12-01
// @Company: mConnect.biz | @License: MIT
// @Description: save (create / update) record(s)

package mccrud

import (
	"fmt"
	"github.com/abbeymart/mccrud/helper"
	"github.com/abbeymart/mcresponse"
)

func (crud Crud) Save() mcresponse.ResponseMessage {
	//  determine taskType from actionParams: create or update
	//  iterate through actionParams, update createRecs, updateRecs & crud.recordIds
	var (
		createRecs ActionParamsType	// records without id or _id field-value
		updateRecs ActionParamsType	// records with id or _id field-value
	)
	for _, rec := range crud.ActionParams {
		// determine if record existed (update) or is new (create)
		if fieldValue, ok := rec["id"]; ok && fieldValue != nil {
			updateRecs = append(updateRecs, rec)
			// validate fieldValue as string
			switch fieldValue.(type) {
			case string:
				crud.RecordIds = append(crud.RecordIds, fieldValue.(string) )
			default:
				// invalid fieldValue type (string)
				return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
					Message: fmt.Sprintf("Invalid fieldValue type for fieldName: id, and in record: %v", rec),
					Value:   nil,
				})
			}
		} else {
			createRecs = append(createRecs, rec)
		}
	}
	// save-record(s): create/insert new record(s), recordIds = @[], if len(createRecs) > 0
	if len(createRecs) > 0{
		return crud.Create(createRecs)
	}
	// update-record(s): update existing record(s), recordIds != @[], if len(updateRecs) > 0
	if len(updateRecs) > 0 {
		return crud.Update(updateRecs)
	}
	// otherwise return saveError
	return mcresponse.GetResMessage("saveError", mcresponse.ResponseMessageOptions{
		Message: "Save error: incomplete or invalid action-params provided",
		Value:   nil,
	})
}

func (crud Crud) Create(createRecs ActionParamsType) mcresponse.ResponseMessage {
	// create from createRecs (actionParams)
	var tableFields []string
	if createQuery, err := helper.ComputeCreateQuery(crud.TableName, tableFields, createRecs); err != nil {
		return mcresponse.GetResMessage("insertError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing create-query: %v", err.Error()),
			Value:   createQuery,
		})
		// TODO: perform create/insert action, wrap in transaction:


	}


	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "success",
		Value:   nil,
	})
}

func (crud Crud) Update(updateRecs ActionParamsType) mcresponse.ResponseMessage {
	// create from updatedRecs (actionParams)

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "success",
		Value:   nil,
	})
}

func (crud Crud) ComputeRecords() mcresponse.ResponseMessage {

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "success",
		Value:   nil,
	})
}
