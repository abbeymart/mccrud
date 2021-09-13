// @Author: abbeymart | Abi Akindele | @Created: 2020-12-01 | @Updated: 2020-12-01
// @Company: mConnect.biz | @License: MIT
// @Description: save (create / update) record(s)

package mccrud

import (
	"context"
	"fmt"
	"github.com/abbeymart/mcauditlog"
	"github.com/abbeymart/mccache"
	"github.com/abbeymart/mcresponse"
	"github.com/jackc/pgx/v4"
)

// Create method creates new record(s)
func (crud *Crud) Create(recs ActionParamsType) mcresponse.ResponseMessage {
	// compute query
	createQueryRes := ComputeCreateQuery(crud.TableName, recs)
	if !createQueryRes.Ok {
		return mcresponse.GetResMessage("insertError", mcresponse.ResponseMessageOptions{
			Message: createQueryRes.Message,
			Value:   nil,
		})
	}
	fmt.Printf("create-query: %v", createQueryRes.CreateQueryObject.CreateQuery)
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
	// create new records by fieldValues
	for _, fValues := range createQueryRes.CreateQueryObject.FieldValues {
		insertErr := tx.QueryRow(context.Background(), createQueryRes.CreateQueryObject.CreateQuery, fValues...).Scan(&insertId)
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
	_ = mccache.DeleteHashCache(crud.TableName, crud.CacheKey, "key")
	// perform audit-log
	logMessage := ""
	logRes := mcresponse.ResponseMessage{}
	var logErr error
	if crud.LogCreate {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:  crud.TableName,
			LogRecords: crud.ActionParams,
		}
		if logRes, logErr = crud.TransLog.AuditLog(CreateTask, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}
	// response
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: logMessage,
		Value: CrudResultType{
			RecordIds:    insertIds,
			RecordsCount: insertCount,
			TaskType:     crud.TaskType,
			LogRes:       logRes,
		},
	})
}

// CreateCopy method creates new record(s) using Pg CopyFrom
// resolve sql-values parsing error (incorrect binary data format (SQLSTATE 22P03) - ?uuid primary key?)
func (crud *Crud) CreateCopy(recs ActionParamsType) mcresponse.ResponseMessage {
	// create from createRecs (actionParams)
	// compute query
	createQueryRes := ComputeCreateQuery(crud.TableName, recs)
	if !createQueryRes.Ok {
		return mcresponse.GetResMessage("insertError", mcresponse.ResponseMessageOptions{
			Message: createQueryRes.Message,
			Value:   nil,
		})
	}
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
		createQueryRes.CreateQueryObject.FieldNames,
		pgx.CopyFromRows(createQueryRes.CreateQueryObject.FieldValues),
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
	_ = mccache.DeleteHashCache(crud.TableName, crud.CacheKey, "hash")
	// perform audit-log
	logMessage := ""
	if crud.LogCreate || crud.LogCrud {
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
			RecordIds:    crud.RecordIds,
			RecordsCount: int(copyCount),
		},
	})
}

// Update method updates existing record(s)
func (crud *Crud) Update(recs ActionParamsType) mcresponse.ResponseMessage {
	// include audit-log feature
	if crud.LogUpdate || crud.LogCrud {
		getRes := crud.GetByIds()
		value, _ := getRes.Value.(CrudResultType)
		crud.CurrentRecords = value.Records
	}
	// create from updatedRecs (actionParams)
	updateQueryRes := ComputeUpdateQuery(crud.TableName, recs)
	if !updateQueryRes.Ok {
		return mcresponse.GetResMessage("updateError", mcresponse.ResponseMessageOptions{
			Message: updateQueryRes.Message,
			Value:   nil,
		})
	}
	fmt.Printf("update-query: %v", updateQueryRes.UpdateQueryObjects[0].UpdateQuery)
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
	for _, upQuery := range updateQueryRes.UpdateQueryObjects {
		commandTag, updateErr := tx.Exec(context.Background(), upQuery.UpdateQuery, upQuery.FieldValues...)
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
	_ = mccache.DeleteHashCache(crud.TableName, crud.CacheKey, "key")
	// perform audit-log
	logMessage := ""
	logRes := mcresponse.ResponseMessage{}
	var logErr error
	if crud.LogUpdate || crud.LogCrud {
		currentRecs := map[string]interface{}{"currentRecords": crud.CurrentRecords, "recordIds": crud.RecordIds}
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:     crud.TableName,
			LogRecords:    currentRecs,
			NewLogRecords: crud.ActionParams,
		}
		if logRes, logErr = crud.TransLog.AuditLog(UpdateTask, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}
	// response
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: fmt.Sprintf("Record(s) update completed successfully [log-message: %v]", logMessage),
		Value: CrudResultType{
			QueryParam:   crud.QueryParams,
			RecordIds:    crud.RecordIds,
			RecordsCount: updateCount,
			TaskType:     crud.TaskType,
			LogRes:       logRes,
		},
	})
}

// UpdateById method updates existing records (in batch) that met the specified record-id(s)
func (crud *Crud) UpdateById(rec ActionParamType, id string) mcresponse.ResponseMessage {
	// include audit-log feature
	if crud.LogUpdate || crud.LogCrud {
		getRes := crud.GetById(id)
		value, _ := getRes.Value.(CrudResultType)
		crud.CurrentRecords = value.Records
	}
	// create from updatedRecs (actionParams)
	updateQueryRes := ComputeUpdateQueryById(crud.TableName, rec, id)
	if !updateQueryRes.Ok {
		return mcresponse.GetResMessage("updateError", mcresponse.ResponseMessageOptions{
			Message: updateQueryRes.Message,
			Value:   nil,
		})
	}
	fmt.Printf("update-query: %v", updateQueryRes.UpdateQueryObject.UpdateQuery)
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
	commandTag, updateErr := tx.Exec(context.Background(), updateQueryRes.UpdateQueryObject.UpdateQuery, updateQueryRes.UpdateQueryObject.FieldValues...)
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
	_ = mccache.DeleteHashCache(crud.TableName, crud.CacheKey, "key")
	// perform audit-log
	logMessage := ""
	logRes := mcresponse.ResponseMessage{}
	var logErr error
	if crud.LogUpdate || crud.LogCrud {
		currentRecs := map[string]interface{}{"currentRecords": crud.CurrentRecords, "recordIds": []string{id}}
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:     crud.TableName,
			LogRecords:    currentRecs,
			NewLogRecords: crud.ActionParams,
		}
		if logRes, logErr = crud.TransLog.AuditLog(UpdateTask, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}
	// response
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: fmt.Sprintf("Record(s) update completed successfully [log-message: %v]", logMessage),
		Value: CrudResultType{
			QueryParam:   crud.QueryParams,
			RecordIds:    crud.RecordIds,
			RecordsCount: int(commandTag.RowsAffected()),
			TaskType:     crud.TaskType,
			LogRes:       logRes,
		},
	})
}

// UpdateByIds method updates existing records (in batch) that met the specified record-id(s)
func (crud *Crud) UpdateByIds(rec ActionParamType) mcresponse.ResponseMessage {
	// include audit-log feature
	if crud.LogUpdate || crud.LogCrud {
		getRes := crud.GetByIds()
		value, _ := getRes.Value.(CrudResultType)
		crud.CurrentRecords = value.Records
	}
	// create from updatedRecs (actionParams)
	updateQueryRes := ComputeUpdateQueryByIds(crud.TableName, rec, crud.RecordIds)
	if !updateQueryRes.Ok {
		return mcresponse.GetResMessage("updateError", mcresponse.ResponseMessageOptions{
			Message: updateQueryRes.Message,
			Value:   nil,
		})
	}
	fmt.Printf("update-query: %v", updateQueryRes.UpdateQueryObject.UpdateQuery)
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
	commandTag, updateErr := tx.Exec(context.Background(), updateQueryRes.UpdateQueryObject.UpdateQuery, updateQueryRes.UpdateQueryObject.FieldValues...)
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
	_ = mccache.DeleteHashCache(crud.TableName, crud.CacheKey, "key")
	// perform audit-log
	logMessage := ""
	logRes := mcresponse.ResponseMessage{}
	var logErr error
	if crud.LogUpdate || crud.LogCrud {
		currentRecs := map[string]interface{}{"currentRecords": crud.CurrentRecords, "recordIds": crud.RecordIds}
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:     crud.TableName,
			LogRecords:    currentRecs,
			NewLogRecords: crud.ActionParams,
		}
		if logRes, logErr = crud.TransLog.AuditLog(UpdateTask, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}
	// response
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: fmt.Sprintf("Record(s) update completed successfully [log-message: %v]", logMessage),
		Value: CrudResultType{
			QueryParam:   crud.QueryParams,
			RecordIds:    crud.RecordIds,
			RecordsCount: int(commandTag.RowsAffected()),
			TaskType:     crud.TaskType,
			LogRes:       logRes,
		},
	})
}

// UpdateByParam method updates existing records (in batch) that met the specified query-params or where conditions
func (crud *Crud) UpdateByParam(rec ActionParamType) mcresponse.ResponseMessage {
	// include audit-log feature
	if crud.LogUpdate || crud.LogCrud {
		getRes := crud.GetByParam()
		value, _ := getRes.Value.(CrudResultType)
		crud.CurrentRecords = value.Records
	}
	// create from updatedRecs (actionParams)
	updateQueryRes := ComputeUpdateQueryByParam(crud.TableName, rec, crud.QueryParams)
	if !updateQueryRes.Ok {
		return mcresponse.GetResMessage("updateError", mcresponse.ResponseMessageOptions{
			Message: updateQueryRes.Message,
			Value:   nil,
		})
	}
	fmt.Printf("update-query: %v", updateQueryRes.UpdateQueryObject.UpdateQuery)
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
	updateFieldValues := updateQueryRes.UpdateQueryObject.FieldValues
	commandTag, updateErr := tx.Exec(context.Background(), updateQueryRes.UpdateQueryObject.UpdateQuery, updateFieldValues...)
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
	_ = mccache.DeleteHashCache(crud.TableName, crud.CacheKey, "key")
	// perform audit-log
	logMessage := ""
	logRes := mcresponse.ResponseMessage{}
	var logErr error
	if crud.LogUpdate || crud.LogCrud {
		currentRecs := map[string]interface{}{"currentRecords": crud.CurrentRecords, "queryParams": crud.QueryParams}
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:     crud.TableName,
			LogRecords:    currentRecs,
			NewLogRecords: crud.ActionParams,
		}
		if logRes, logErr = crud.TransLog.AuditLog(UpdateTask, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}
	// response
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: fmt.Sprintf("Record(s) update completed successfully [log-message: %v]", logMessage),
		Value: CrudResultType{
			QueryParam:   crud.QueryParams,
			RecordIds:    crud.RecordIds,
			RecordsCount: int(commandTag.RowsAffected()),
			TaskType:     crud.TaskType,
			LogRes:       logRes,
		},
	})
}
