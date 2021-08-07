// @Author: abbeymart | Abi Akindele | @Created: 2020-12-01 | @Updated: 2020-12-01
// @Company: mConnect.biz | @License: MIT
// @Description: save (create / update) record(s)

package mccrud

import (
	"context"
	"fmt"
	"github.com/abbeymart/mcauditlog"
	"github.com/abbeymart/mccache"
	"github.com/abbeymart/mccrud/helper"
	"github.com/abbeymart/mcresponse"
	"github.com/jackc/pgx/v4"
)

// Save method creates new record(s) or updates existing record(s)
func (crud *Crud) Save(modelRef interface{}, recs interface{}, batch int) mcresponse.ResponseMessage {
	//  determine taskType from actionParams: create or update
	//  iterate through actionParams: update createRecs, updateRecs & crud.recordIds
	var (
		createRecs ActionParamsType // records without id field-value
		updateRecs ActionParamsType // records with id field-value
		recIds     []string         // capture recordIds for separate/multiple updates
	)
	for _, rec := range crud.ActionParams {
		// determine if record existed (update) or is new (create)
		if fieldValue, ok := rec["id"]; ok && fieldValue != "" {
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
		return crud.CreateBatch(createRecs, tableFields)
	}

	// update each record by its recordId
	if len(updateRecs) >= 1 && (len(recIds) == len(updateRecs)) {
		return crud.Update(updateRecs, tableFields)
	}

	// update record(s) by recordIds | CONTROL ACCESS (by api-user)
	if len(updateRecs) == 1 && len(crud.RecordIds) > 0 {
		return crud.UpdateById(updateRecs, tableFields)
	}

	// update record(s) by queryParams | CONTROL ACCESS (by api-user)
	if len(updateRecs) == 1 && len(crud.QueryParams) > 0 {
		return crud.UpdateByParam(updateRecs, tableFields)
	}

	// otherwise, return saveError
	return mcresponse.GetResMessage("saveError", mcresponse.ResponseMessageOptions{
		Message: "Save error: incomplete or invalid action/query-params provided",
		Value:   nil,
	})
}

// Create method creates new record(s)
func (crud *Crud) Create(modelRef interface{}, recs interface{}, batch int) mcresponse.ResponseMessage {
	// compute query
	createQuery, qErr := helper.ComputeCreateQuery(crud.TableName, createRecs, tableFields)
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
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {

		}
	}(tx, context.Background())

	// perform records' creation
	insertCount := 0
	var insertIds []string
	var insertId string
	for _, insertQuery := range createQuery {
		insertErr := tx.QueryRow(context.Background(), insertQuery).Scan(&insertId)
		if insertErr != nil {
			_ = tx.Rollback(context.Background())
			return mcresponse.GetResMessage("updateError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error updating record(s): %v", insertErr.Error()),
				Value:   nil,
			})
		}
		insertCount += 1
		insertIds = append(insertIds, insertId)
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
	// delete cache
	_ = mccache.DeleteHashCache(crud.TableName, crud.HashKey, "hash")

	// perform audit-log
	logMessage := ""
	if crud.LogCreate {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:  crud.TableName,
			LogRecords: crud.ActionParams,
		}
		if logRes, logErr := crud.TransLog.AuditLog(CreateTask, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: logMessage,
		Value: CrudResultType{
			RecordIds:   insertIds,
			RecordCount: insertCount,
		},
	})
}

// CreateBatch method creates new record(s) by placeholder values from copy-create-query
// resolve sql-values parsing error: only time.Time and String value requires '' wrapping
// uuid, json and others (int/bool/float) should not be wrapped as placeholder values
func (crud *Crud) CreateBatch(recs interface{}, batch int) mcresponse.ResponseMessage {
	// create from createRecs (actionParams)
	// compute query
	createQuery, qErr := helper.ComputeCreateCopyQuery(crud.TableName, createRecs, tableFields)
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
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {

		}
	}(tx, context.Background())

	// perform records' creation
	insertCount := 0
	var insertIds []string
	var insertId string
	for _, iValues := range createQuery.FieldValues {
		//fmt.Printf("query: %v\n\n", createQuery.CreateQuery)
		//fmt.Printf("query-value: %v \n\n", iValues)
		insertErr := tx.QueryRow(context.Background(), createQuery.CreateQuery, iValues...).Scan(&insertId)
		if insertErr != nil {
			_ = tx.Rollback(context.Background())
			return mcresponse.GetResMessage("updateError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error updating record(s): %v", insertErr.Error()),
				Value:   nil,
			})
		}
		insertCount += 1
		insertIds = append(insertIds, insertId)
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
	// delete cache
	_ = mccache.DeleteHashCache(crud.TableName, crud.HashKey, "hash")

	// perform audit-log
	logMessage := ""
	if crud.LogCreate {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:  crud.TableName,
			LogRecords: crud.ActionParams,
		}
		if logRes, logErr := crud.TransLog.AuditLog(CreateTask, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: logMessage,
		Value: CrudResultType{
			RecordIds:   insertIds,
			RecordCount: insertCount,
		},
	})
}

// CreateCopy method creates new record(s) using Pg CopyFrom
// TODO: resolve sql-values parsing error (incorrect binary data format (SQLSTATE 22P03) - ?uuid primary key?)
func (crud *Crud) CreateCopy(createRecs ActionParamsType, tableFields []string) mcresponse.ResponseMessage {
	// create from createRecs (actionParams)
	// compute query
	createQuery, qErr := helper.ComputeCreateCopyQuery(crud.TableName, createRecs, tableFields)
	if qErr != nil {
		return mcresponse.GetResMessage("insertError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing create-query: %v", qErr.Error()),
			Value:   nil,
		})
	}
	//fmt.Printf("create-query: %v \n", createQuery)
	//fmt.Printf("create-query-fields: %v \n", createQuery.FieldNames)
	//fmt.Printf("create-query-values: %v \n\n", createQuery.FieldValues)
	// perform create/insert action, via transaction/copy-protocol:
	tx, txErr := crud.AppDb.Begin(context.Background())
	if txErr != nil {
		return mcresponse.GetResMessage("insertError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error creating new record(s): %v", txErr.Error()),
			Value:   nil,
		})
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {

		}
	}(tx, context.Background())

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

	// delete cache
	_ = mccache.DeleteHashCache(crud.TableName, crud.HashKey, "hash")

	// perform audit-log
	logMessage := ""
	if crud.LogCreate {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:  crud.TableName,
			LogRecords: crud.ActionParams,
		}
		if logRes, logErr := crud.TransLog.AuditLog(CreateTask, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: logMessage,
		Value: CrudResultType{
			RecordIds:   crud.RecordIds,
			RecordCount: int(copyCount),
		},
	})
}

// Update method updates existing record(s)
func (crud *Crud) Update(modelRef interface{}, recs interface{}) mcresponse.ResponseMessage {
	// create from updatedRecs (actionParams)
	updateQuery, err := helper.ComputeUpdateQuery(crud.TableName, updateRecs, tableFields)
	if err != nil {
		return mcresponse.GetResMessage("updateError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing update-query: %v", err.Error()),
			Value:   nil,
		})
	}
	// perform update action, via transaction:
	tx, txErr := crud.AppDb.Begin(context.Background())
	if txErr != nil {
		return mcresponse.GetResMessage("updateError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error updating record(s): %v", txErr.Error()),
			Value:   nil,
		})
	}
	//defer tx.Rollback(context.Background())
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {

		}
	}(tx, context.Background())
	// perform records' updates
	updateCount := 0
	//fmt.Printf("update-queries: %v\n", updateQuery)
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

	// delete cache
	_ = mccache.DeleteHashCache(crud.TableName, crud.HashKey, "hash")

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "Record(s) update completed successfully",
		Value: CrudResultType{
			QueryParam:  crud.QueryParams,
			RecordIds:   crud.RecordIds,
			RecordCount: updateCount,
		},
	})
}

// UpdateById method updates existing records (in batch) that met the specified record-id(s)
func (crud *Crud) UpdateById(modelRef interface{}, recs interface{}, id string) mcresponse.ResponseMessage {
	// create from updatedRecs (actionParams)
	updateQuery, err := helper.ComputeUpdateQueryById(crud.TableName, updateRecs, crud.RecordIds, tableFields)
	if err != nil {
		return mcresponse.GetResMessage("updateError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing update-query: %v", err.Error()),
			Value:   nil,
		})
	}
	// perform update action, via transaction:
	tx, txErr := crud.AppDb.Begin(context.Background())
	if txErr != nil {
		return mcresponse.GetResMessage("updateError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error updating record(s): %v", txErr.Error()),
			Value:   nil,
		})
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {

		}
	}(tx, context.Background())
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

	// delete cache
	_ = mccache.DeleteHashCache(crud.TableName, crud.HashKey, "hash")

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "Record(s) update completed successfully",
		Value: CrudResultType{
			QueryParam:  crud.QueryParams,
			RecordIds:   crud.RecordIds,
			RecordCount: int(commandTag.RowsAffected()),
		},
	})
}

// UpdateByIds method updates existing records (in batch) that met the specified record-id(s)
func (crud *Crud) UpdateByIds(modelRef interface{}, recs interface{}) mcresponse.ResponseMessage {
	// create from updatedRecs (actionParams)
	updateQuery, err := helper.ComputeUpdateQueryById(crud.TableName, updateRecs, crud.RecordIds, tableFields)
	if err != nil {
		return mcresponse.GetResMessage("updateError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing update-query: %v", err.Error()),
			Value:   nil,
		})
	}
	// perform update action, via transaction:
	tx, txErr := crud.AppDb.Begin(context.Background())
	if txErr != nil {
		return mcresponse.GetResMessage("updateError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error updating record(s): %v", txErr.Error()),
			Value:   nil,
		})
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {

		}
	}(tx, context.Background())
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

	// delete cache
	_ = mccache.DeleteHashCache(crud.TableName, crud.HashKey, "hash")

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "Record(s) update completed successfully",
		Value: CrudResultType{
			QueryParam:  crud.QueryParams,
			RecordIds:   crud.RecordIds,
			RecordCount: int(commandTag.RowsAffected()),
		},
	})
}

// UpdateByParam method updates existing records (in batch) that met the specified query-params or where conditions
func (crud *Crud) UpdateByParam(modelRef interface{}, recs interface{}) mcresponse.ResponseMessage {
	// create from updatedRecs (actionParams)
	updateQuery, err := helper.ComputeUpdateQueryByParam(crud.TableName, updateRecs, crud.QueryParams, tableFields)
	if err != nil {
		return mcresponse.GetResMessage("updateError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing update-query: %v", err.Error()),
			Value:   nil,
		})
	}
	// perform update action, via transaction:
	tx, txErr := crud.AppDb.Begin(context.Background())
	if txErr != nil {
		return mcresponse.GetResMessage("updateError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error updating record(s): %v", txErr.Error()),
			Value:   nil,
		})
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {

		}
	}(tx, context.Background())
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

	// delete cache
	_ = mccache.DeleteHashCache(crud.TableName, crud.HashKey, "hash")

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "Record(s) update completed successfully",
		Value: CrudResultType{
			QueryParam:  crud.QueryParams,
			RecordIds:   crud.RecordIds,
			RecordCount: int(commandTag.RowsAffected()),
		},
	})
}

func (crud *Crud) UpdateLog(updateRecs ActionParamsType, tableFields []string, upTableFields []string, tableFieldPointers []interface{}) mcresponse.ResponseMessage {
	// get records to update, for audit-log
	if crud.LogUpdate && len(tableFields) == len(tableFieldPointers) {
		getRes := crud.GetById(tableFields, tableFieldPointers)
		value, _ := getRes.Value.(CrudResultType)
		crud.CurrentRecords = value.TableRecords
	}

	// perform update
	updateRes := crud.Update(updateRecs, upTableFields)

	// perform audit-log
	logMessage := ""
	if crud.LogUpdate {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:     crud.TableName,
			LogRecords:    crud.CurrentRecords,
			NewLogRecords: crud.ActionParams,
		}
		if logRes, logErr := crud.TransLog.AuditLog(UpdateTask, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}

	// overall response
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: updateRes.Message + " | " + logMessage,
		Value:   updateRes.Value,
	})
}

func (crud *Crud) UpdateByIdLog(updateRecs ActionParamsType, tableFields []string, upTableFields []string, tableFieldPointers []interface{}) mcresponse.ResponseMessage {
	// get records to update, for audit-log
	if crud.LogUpdate && len(tableFields) == len(tableFieldPointers) {
		getRes := crud.GetById(tableFields, tableFieldPointers)
		value, _ := getRes.Value.(CrudResultType)
		crud.CurrentRecords = value.TableRecords
	}

	// perform update-by-id
	updateRes := crud.UpdateById(updateRecs, upTableFields)

	// perform audit-log
	logMessage := ""
	if crud.LogUpdate {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:     crud.TableName,
			LogRecords:    crud.CurrentRecords,
			NewLogRecords: crud.ActionParams,
		}
		if logRes, logErr := crud.TransLog.AuditLog(UpdateTask, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}

	// overall response
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: updateRes.Message + " | " + logMessage,
		Value:   updateRes.Value,
	})
}

func (crud *Crud) UpdateByParamLog(updateRecs ActionParamsType, tableFields []string, upTableFields []string, tableFieldPointers []interface{}) mcresponse.ResponseMessage {
	// get records to update, for audit-log
	if crud.LogUpdate && len(tableFields) == len(tableFieldPointers) {
		getRes := crud.GetByParam(tableFields, tableFieldPointers)
		value, _ := getRes.Value.(CrudResultType)
		crud.CurrentRecords = value.TableRecords
	}

	// perform update-by-id
	updateRes := crud.UpdateByParam(updateRecs, upTableFields)

	// perform audit-log
	logMessage := ""
	if crud.LogUpdate {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:     crud.TableName,
			LogRecords:    crud.CurrentRecords,
			NewLogRecords: crud.ActionParams,
		}
		if logRes, logErr := crud.TransLog.AuditLog(UpdateTask, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}

	// overall response
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: updateRes.Message + " | " + logMessage,
		Value:   updateRes.Value,
	})
}
