// @Author: abbeymart | Abi Akindele | @Created: 2020-12-01 | @Updated: 2020-12-01
// @Company: mConnect.biz | @License: MIT
// @Description: save (create / update) record(s)

package mccrud

import (
	"context"
	"fmt"
	"github.com/abbeymart/mccrud/helper"
	"github.com/abbeymart/mcresponse"
	"github.com/abbeymart/mctypes"
	"github.com/jackc/pgx/v4"
)

// Save method creates new record(s) or updates existing record(s)
func (crud *Crud) Save(tableFields []string) mcresponse.ResponseMessage {
	//  determine taskType from actionParams: create or update
	//  iterate through actionParams, update createRecs, updateRecs & crud.recordIds
	var (
		createRecs mctypes.ActionParamsType // records without id or _id field-value
		updateRecs mctypes.ActionParamsType // records with id or _id field-value
		recIds     []string                 // capture recordIds for separate/multiple updates
	)
	for _, rec := range crud.ActionParams {
		// determine if record existed (update) or is new (create)
		if fieldValue, ok := rec["id"]; ok && fieldValue != nil {
			// validate fieldValue as string
			switch fieldValue.(type) {
			case string:
				updateRecs = append(updateRecs, rec)
				recIds = append(recIds, fieldValue.(string))
			default:
				// invalid fieldValue type (string)
				return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
					Message: fmt.Sprintf("Invalid fieldValue type for fieldName: id, in record: %v", rec),
					Value:   nil,
				})
			}
		} else if len(crud.ActionParams) == 1 && (len(crud.RecordIds) > 0 || len(crud.QueryParams) > 0) {
			updateRecs = append(updateRecs, rec)
		} else {
			createRecs = append(createRecs, rec)
		}
	}

	// permit only create or update, not both at the same time
	if len(createRecs) > 0 && len(updateRecs) > 0 {
		return mcresponse.GetResMessage("saveError", mcresponse.ResponseMessageOptions{
			Message: "You may only create or update record(s), not both at the same time",
			Value:   nil,
		})
	}

	if len(createRecs) > 0 {
		// save-record(s): create/insert new record(s), recordIds = @[], if len(createRecs) > 0
		return crud.Create(createRecs, tableFields)
	}

	// update each record by it's recordId
	if len(updateRecs) >= 1 && len(recIds) == len(updateRecs) {
		return crud.Update(updateRecs, tableFields)
	}

	// update record(s) by recordIds | CONTROL ACCESS
	if len(updateRecs) == 1 && len(crud.RecordIds) > 0 {
		return crud.UpdateById(updateRecs, tableFields)
	}
	// update record(s) by queryParams | CONTROL ACCESS
	if len(updateRecs) == 1 && len(crud.QueryParams) > 0 {
		return crud.UpdateByParam(updateRecs, tableFields)
	}

	// otherwise return saveError
	return mcresponse.GetResMessage("saveError", mcresponse.ResponseMessageOptions{
		Message: "Save error: incomplete or invalid action/query-params provided",
		Value:   nil,
	})
}

// Create method creates new record(s)
func (crud Crud) Create(createRecs mctypes.ActionParamsType, tableFields []string) mcresponse.ResponseMessage {
	// create from createRecs (actionParams)
	// compute query
	createQuery, qErr := helper.ComputeCreateQuery(crud.TableName, tableFields, createRecs)
	if qErr != nil {
		return mcresponse.GetResMessage("insertError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing create-query: %v", qErr.Error()),
			Value:   nil,
		})
	}
	// perform create/insert action, via transaction/copy-protocol:
	tx, txErr := crud.AppDb.Begin(context.Background())
	if txErr != nil {
		return mcresponse.GetResMessage("insertError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error creating new record(s): %v", txErr.Error()),
			Value:   nil,
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
		_ = tx.Rollback(context.Background())
		return mcresponse.GetResMessage("insertError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error creating new record(s): %v", cErr.Error()),
			Value:   nil,
		})
	}
	// commit
	txcErr := tx.Commit(context.Background())
	if txcErr != nil {
		_ = tx.Rollback(context.Background())
		return mcresponse.GetResMessage("insertError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error creating new record(s): %v", txcErr.Error()),
			Value:   nil,
		})
	}
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "success",
		Value:   copyCount,
	})
}

// Update method updates existing record(s)
func (crud Crud) Update(updateRecs mctypes.ActionParamsType, tableFields []string) mcresponse.ResponseMessage {
	// create from updatedRecs (actionParams)
	if updateQuery, err := helper.ComputeUpdateQuery(crud.TableName, tableFields, updateRecs); err != nil {
		return mcresponse.GetResMessage("updateError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing update-query: %v", err.Error()),
			Value:   nil,
		})
	} else {
		// perform update action, via transaction/copy-protocol:
		tx, txErr := crud.AppDb.Begin(context.Background())
		if txErr != nil {
			return mcresponse.GetResMessage("updateError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error updating record(s): %v", txErr.Error()),
				Value:   nil,
			})
		}
		defer tx.Rollback(context.Background())
		// perform records' updates
		updateCount := 0
		for _, upQuery := range updateQuery {
			commandTag, updateErr := tx.Exec(context.Background(), upQuery)
			if updateErr != nil {
				_ = tx.Rollback(context.Background())
				return mcresponse.GetResMessage("updateError", mcresponse.ResponseMessageOptions{
					Message: fmt.Sprintf("Error updating record(s): %v", updateErr.Error()),
					Value:   nil,
				})
			}
			updateCount += int(commandTag.RowsAffected())
		}
		// commit
		txcErr := tx.Commit(context.Background())
		if txcErr != nil {
			_ = tx.Rollback(context.Background())
			return mcresponse.GetResMessage("updateError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error updating record(s): %v", txcErr.Error()),
				Value:   nil,
			})
		}
		return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
			Message: "success",
			Value:   updateCount,
		})
	}
}

// UpdateById method updates existing records (in batch) that met the specified record-id(s)
func (crud Crud) UpdateById(updateRecs mctypes.ActionParamsType, tableFields []string) mcresponse.ResponseMessage {
	// create from updatedRecs (actionParams)
	if updateQuery, err := helper.ComputeUpdateQueryById(crud.TableName, tableFields, updateRecs, crud.RecordIds); err != nil {
		return mcresponse.GetResMessage("updateError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing update-query: %v", err.Error()),
			Value:   nil,
		})
	} else {
		// perform update action, via transaction/copy-protocol:
		tx, txErr := crud.AppDb.Begin(context.Background())
		if txErr != nil {
			return mcresponse.GetResMessage("updateError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error updating record(s): %v", txErr.Error()),
				Value:   nil,
			})
		}
		defer tx.Rollback(context.Background())
		commandTag, updateErr := tx.Exec(context.Background(), updateQuery)
		if updateErr != nil {
			_ = tx.Rollback(context.Background())
			return mcresponse.GetResMessage("updateError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error updating record(s): %v", updateErr.Error()),
				Value:   nil,
			})
		}
		// commit
		txcErr := tx.Commit(context.Background())
		if txcErr != nil {
			_ = tx.Rollback(context.Background())
			return mcresponse.GetResMessage("updateError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error updating record(s): %v", txcErr.Error()),
				Value:   nil,
			})
		}
		return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
			Message: "success",
			Value:   commandTag.RowsAffected(),
		})
	}
}

// UpdateByParam method updates existing records (in batch) that met the specified query-params or where conditions
func (crud Crud) UpdateByParam(updateRecs mctypes.ActionParamsType, tableFields []string) mcresponse.ResponseMessage {
	// create from updatedRecs (actionParams)
	if updateQuery, err := helper.ComputeUpdateQueryByParam(crud.TableName, tableFields, updateRecs, crud.QueryParams); err != nil {
		return mcresponse.GetResMessage("updateError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing update-query: %v", err.Error()),
			Value:   nil,
		})
	} else {
		// perform update action, via transaction/copy-protocol:
		tx, txErr := crud.AppDb.Begin(context.Background())
		if txErr != nil {
			return mcresponse.GetResMessage("updateError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error updating record(s): %v", txErr.Error()),
				Value:   nil,
			})
		}
		defer tx.Rollback(context.Background())
		commandTag, updateErr := tx.Exec(context.Background(), updateQuery)
		if updateErr != nil {
			_ = tx.Rollback(context.Background())
			return mcresponse.GetResMessage("updateError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error updating record(s): %v", updateErr.Error()),
				Value:   nil,
			})
		}
		// commit
		txcErr := tx.Commit(context.Background())
		if txcErr != nil {
			_ = tx.Rollback(context.Background())
			return mcresponse.GetResMessage("updateError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error updating record(s): %v", txcErr.Error()),
				Value:   nil,
			})
		}
		return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
			Message: "success",
			Value:   commandTag.RowsAffected(),
		})
	}
}
