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
	createQueryObject, qErr := ComputeCreateQuery(crud.TableName, recs)
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
	// create new records by fieldValues
	for _, fValues := range createQueryObject.FieldValues {
		insertErr := tx.QueryRow(context.Background(), createQueryObject.CreateQuery, fValues...).Scan(&insertId)
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
	_ = mccache.DeleteHashCache(crud.TableName, crud.CacheKey, "hash")

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
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: logMessage,
		Value: CrudResultType{
			RecordIds:   insertIds,
			RecordCount: insertCount,
			TaskType:    crud.TaskType,
			LogRes:      logRes,
		},
	})
}

// CreateCopy method creates new record(s) using Pg CopyFrom
// resolve sql-values parsing error (incorrect binary data format (SQLSTATE 22P03) - ?uuid primary key?)
func (crud *Crud) CreateCopy(createRecs ActionParamsType) mcresponse.ResponseMessage {
	// create from createRecs (actionParams)
	// compute query
	createQuery, qErr := ComputeCreateQuery(crud.TableName, createRecs)
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
	_ = mccache.DeleteHashCache(crud.TableName, crud.CacheKey, "hash")

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
func (crud *Crud) Update(modelRef interface{}, updateRecs ActionParamsType) mcresponse.ResponseMessage {
	// include audit-log feature
	if crud.LogUpdate || crud.LogCrud {
		getRes := crud.GetByIds(modelRef)
		value, _ := getRes.Value.(CrudResultType)
		crud.CurrentRecords = value.TableRecords
	}
	// create from updatedRecs (actionParams)
	var updateQueryObjects []UpdateQueryObject
	for _, rec := range updateRecs {
		updateQueryObject, err := ComputeUpdateQuery(crud.TableName, rec)
		if err != nil {
			return mcresponse.GetResMessage("updateError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error computing update-query: %v", err.Error()),
				Value:   nil,
			})
		}
		updateQueryObjects = append(updateQueryObjects, updateQueryObject)
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
	for _, upQuery := range updateQueryObjects {
		updateQuery := upQuery.UpdateQuery + upQuery.WhereQuery.WhereQuery
		var updateFieldValues []interface{}
		updateFieldValues = append(upQuery.FieldValues, upQuery.WhereQuery.FieldValues...)
		commandTag, updateErr := tx.Exec(context.Background(), updateQuery, updateFieldValues...)
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
	_ = mccache.DeleteHashCache(crud.TableName, crud.CacheKey, "hash")

	// perform audit-log
	logMessage := ""
	logRes := mcresponse.ResponseMessage{}
	var logErr error
	if crud.LogUpdate || crud.LogCrud {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:     crud.TableName,
			LogRecords:    crud.CurrentRecords,
			NewLogRecords: crud.ActionParams,
		}
		if logRes, logErr = crud.TransLog.AuditLog(UpdateTask, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: fmt.Sprintf("Record(s) update completed successfully [log-message: %v]", logMessage),
		Value: CrudResultType{
			QueryParam:  crud.QueryParams,
			RecordIds:   crud.RecordIds,
			RecordCount: updateCount,
			TaskType:    crud.TaskType,
			LogRes:      logRes,
		},
	})

}

// UpdateById method updates existing records (in batch) that met the specified record-id(s)
func (crud *Crud) UpdateById(modelRef interface{}, updateRec ActionParamType, id string) mcresponse.ResponseMessage {
	// include audit-log feature
	if crud.LogUpdate || crud.LogCrud {
		getRes := crud.GetById(modelRef, id)
		value, _ := getRes.Value.(CrudResultType)
		crud.CurrentRecords = value.TableRecords
	}
	// create from updatedRecs (actionParams)
	upQueryObj, err := ComputeUpdateQueryById(crud.TableName, updateRec, id)
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
	commandTag, updateErr := tx.Exec(context.Background(), upQueryObj.UpdateQuery, upQueryObj.FieldValues...)
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
	_ = mccache.DeleteHashCache(crud.TableName, crud.CacheKey, "hash")

	// perform audit-log
	logMessage := ""
	logRes := mcresponse.ResponseMessage{}
	var logErr error
	if crud.LogUpdate || crud.LogCrud {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:     crud.TableName,
			LogRecords:    crud.CurrentRecords,
			NewLogRecords: crud.ActionParams,
		}
		if logRes, logErr = crud.TransLog.AuditLog(UpdateTask, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: fmt.Sprintf("Record(s) update completed successfully [log-message: %v]", logMessage),
		Value: CrudResultType{
			QueryParam:  crud.QueryParams,
			RecordIds:   crud.RecordIds,
			RecordCount: int(commandTag.RowsAffected()),
			TaskType:    crud.TaskType,
			LogRes:      logRes,
		},
	})

}

// UpdateByIds method updates existing records (in batch) that met the specified record-id(s)
func (crud *Crud) UpdateByIds(modelRef interface{}, updateRec ActionParamType) mcresponse.ResponseMessage {
	// include audit-log feature
	if crud.LogUpdate || crud.LogCrud {
		getRes := crud.GetByIds(modelRef)
		value, _ := getRes.Value.(CrudResultType)
		crud.CurrentRecords = value.TableRecords
	}
	// create from updatedRecs (actionParams)
	upQueryObj, err := ComputeUpdateQueryByIds(crud.TableName, updateRec, crud.RecordIds)
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
	commandTag, updateErr := tx.Exec(context.Background(), upQueryObj.UpdateQuery, upQueryObj.FieldValues...)
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
	_ = mccache.DeleteHashCache(crud.TableName, crud.CacheKey, "hash")

	// perform audit-log
	logMessage := ""
	logRes := mcresponse.ResponseMessage{}
	var logErr error
	if crud.LogUpdate || crud.LogCrud {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:     crud.TableName,
			LogRecords:    crud.CurrentRecords,
			NewLogRecords: crud.ActionParams,
		}
		if logRes, logErr = crud.TransLog.AuditLog(UpdateTask, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: fmt.Sprintf("Record(s) update completed successfully [log-message: %v]", logMessage),
		Value: CrudResultType{
			QueryParam:  crud.QueryParams,
			RecordIds:   crud.RecordIds,
			RecordCount: int(commandTag.RowsAffected()),
			TaskType:    crud.TaskType,
			LogRes:      logRes,
		},
	})

}

// UpdateByParam method updates existing records (in batch) that met the specified query-params or where conditions
func (crud *Crud) UpdateByParam(modelRef interface{}, updateRec ActionParamType) mcresponse.ResponseMessage {
	// include audit-log feature
	if crud.LogUpdate || crud.LogCrud {
		getRes := crud.GetByParam(modelRef)
		value, _ := getRes.Value.(CrudResultType)
		crud.CurrentRecords = value.TableRecords
	}
	// create from updatedRecs (actionParams)
	updateQueryObject, err := ComputeUpdateQueryByParam(crud.TableName, updateRec, crud.QueryParams)
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
	updateQuery := updateQueryObject.UpdateQuery + updateQueryObject.WhereQuery.WhereQuery
	var updateFieldValues []interface{}
	updateFieldValues = append(updateQueryObject.FieldValues, updateQueryObject.WhereQuery.FieldValues...)
	commandTag, updateErr := tx.Exec(context.Background(), updateQuery, updateFieldValues...)
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
	_ = mccache.DeleteHashCache(crud.TableName, crud.CacheKey, "hash")

	// perform audit-log
	logMessage := ""
	logRes := mcresponse.ResponseMessage{}
	var logErr error
	if crud.LogUpdate || crud.LogCrud {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:     crud.TableName,
			LogRecords:    crud.CurrentRecords,
			NewLogRecords: crud.ActionParams,
		}
		if logRes, logErr = crud.TransLog.AuditLog(UpdateTask, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: fmt.Sprintf("Record(s) update completed successfully [log-message: %v]", logMessage),
		Value: CrudResultType{
			QueryParam:  crud.QueryParams,
			RecordIds:   crud.RecordIds,
			RecordCount: int(commandTag.RowsAffected()),
			TaskType:    crud.TaskType,
			LogRes:      logRes,
		},
	})
}
