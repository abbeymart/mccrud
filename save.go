// @Author: abbeymart | Abi Akindele | @Created: 2020-12-01 | @Updated: 2020-12-01
// @Company: mConnect.biz | @License: MIT
// @Description: save (create / update) record(s)

package mccrud

import (
	"context"
	"fmt"
	"github.com/abbeymart/mccrud/helper"
	"github.com/abbeymart/mcresponse"
	"github.com/jackc/pgx/v4"
)

func (crud *Crud) Save() mcresponse.ResponseMessage {
	//  determine taskType from actionParams: create or update
	//  iterate through actionParams, update createRecs, updateRecs & crud.recordIds
	var (
		createRecs ActionParamsType // records without id or _id field-value
		updateRecs ActionParamsType // records with id or _id field-value
	)
	for _, rec := range crud.ActionParams {
		// determine if record existed (update) or is new (create)
		if fieldValue, ok := rec["id"]; ok && fieldValue != nil {
			updateRecs = append(updateRecs, rec)
			// validate fieldValue as string
			switch fieldValue.(type) {
			case string:
				crud.RecordIds = append(crud.RecordIds, fieldValue.(string))
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

	// permit only create or update, not both at the same time
	if len(createRecs) > 0 && len(updateRecs) > 0 {
		return mcresponse.GetResMessage("saveError", mcresponse.ResponseMessageOptions{
			Message: "you may only create or update record(s), not both at the same time",
			Value:   nil,
		})
	}

	if len(createRecs) > 0 {
		// save-record(s): create/insert new record(s), recordIds = @[], if len(createRecs) > 0
		return crud.Create(createRecs)
	}

	// update record(s) by recordIds or by queryParams
	if len(updateRecs) > 0 && len(crud.RecordIds) > 0 {
		// update-record(s): update existing record(s), recordIds != @[], if len(updateRecs) > 0
		return crud.UpdateById(updateRecs)
	}
	if len(updateRecs) > 0 && len(crud.QueryParams) > 0 {
		// update-record(s): update existing record(s), recordIds != @[], if len(updateRecs) > 0
		return crud.UpdateByParam(updateRecs)
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
	// compose tableFields
	if tFields, err := helper.ComputeSaveFields(createRecs, crud.ProjectParams); err != nil {
		return mcresponse.GetResMessage("insertError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing create-query: %v", err.Error()),
			Value:   tFields,
		})
	} else {
		tableFields = tFields
	}
	// compute query
	createQuery, qErr := helper.ComputeCreateQuery(crud.TableName, tableFields, createRecs)
	if qErr != nil {
		return mcresponse.GetResMessage("insertError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing create-query: %v", qErr.Error()),
			Value:   createQuery,
		})
	}

	// perform create/insert action, via transaction/copy-protocol:
	tx, txErr := crud.AppDb.Begin(context.Background())
	if txErr != nil {
		return mcresponse.GetResMessage("insertError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error creating new records: %v", txErr.Error()),
			Value:   createQuery,
		})
	}
	defer tx.Rollback(context.Background())

	// bulk create
	copyCount, cErr := tx.CopyFrom(
		context.Background(),
		pgx.Identifier{crud.TableName},
		createQuery.FieldNames,
		pgx.CopyFromRows(createQuery.FieldValues),
	)
	if cErr != nil {
		return mcresponse.GetResMessage("insertError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error creating new records: %v", cErr.Error()),
			Value:   createQuery,
		})
	}
	// commit
	txcErr := tx.Commit(context.Background())
	if txcErr != nil {
		return mcresponse.GetResMessage("insertError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error creating new records: %v", txcErr.Error()),
			Value:   createQuery,
		})
	}
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "success",
		Value:   copyCount,
	})
}

func (crud Crud) UpdateById(updateRecs ActionParamsType) mcresponse.ResponseMessage {
	// create from updatedRecs (actionParams)
	var tableFields []string
	// compose tableFields
	if tFields, err := helper.ComputeSaveFields(updateRecs, crud.ProjectParams); err != nil {
		return mcresponse.GetResMessage("insertError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing create-query: %v", err.Error()),
			Value:   tFields,
		})
	} else {
		tableFields = tFields
	}

	if updateQuery, err := helper.ComputeUpdateQueryById(crud.TableName, tableFields, updateRecs, crud.RecordIds); err != nil {
		return mcresponse.GetResMessage("insertError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing update-query: %v", err.Error()),
			Value:   updateQuery,
		})
	} else {
		// TODO: perform create/insert action:
	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "success",
		Value:   nil,
	})
}

func (crud Crud) UpdateByParam(updateRecs ActionParamsType) mcresponse.ResponseMessage {
	// create from updatedRecs (actionParams)
	var tableFields []string
	// compose tableFields
	if tFields, err := helper.ComputeSaveFields(updateRecs, crud.ProjectParams); err != nil {
		return mcresponse.GetResMessage("insertError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing create-query: %v", err.Error()),
			Value:   tFields,
		})
	} else {
		tableFields = tFields
	}

	if updateQuery, err := helper.ComputeUpdateQueryByParam(crud.TableName, tableFields, updateRecs, crud.QueryParams); err != nil {
		return mcresponse.GetResMessage("insertError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing update-query: %v", err.Error()),
			Value:   updateQuery,
		})
	} else {
		// TODO: perform create/insert action:
	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "success",
		Value:   nil,
	})
}
