// @Author: abbeymart | Abi Akindele | @Created: 2020-12-01 | @Updated: 2020-12-01
// @Company: mConnect.biz | @License: MIT
// @Description: save (create / update) record(s)

package mccrud

import (
	"context"
	"fmt"
	"github.com/abbeymart/mcauditlog"
	"github.com/abbeymart/mccrud/helper"
	"github.com/abbeymart/mcresponse"
	"github.com/abbeymart/mctypes"
	"github.com/abbeymart/mctypes/tasks"
	"github.com/jackc/pgx/v4"
)

// Save method creates new record(s) or updates existing record(s)
func (crud *Crud) Save(tableFields []string) mcresponse.ResponseMessage {
	//fmt.Printf("save-action-params: %#v \n\n", crud.ActionParams)
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
		return crud.Update(updateRecs, tableFields, recIds)
	}

	// update record(s) by recordIds | CONTROL ACCESS (by api-user)
	if len(updateRecs) == 1 && len(crud.RecordIds) > 0 {
		return crud.UpdateById(updateRecs, tableFields)
	}

	// update record(s) by queryParams | CONTROL ACCESS (by api-user)
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
	// perform audit-log
	logMessage := ""
	if crud.LogCreate {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:  crud.TableName,
			LogRecords: crud.ActionParams,
		}
		if logRes, logErr := crud.TransLog.AuditLog(tasks.Create, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: logMessage,
		Value: InsertedResultType{
			TableName:   crud.TableName,
			RecordIds:   insertIds,
			RecordCount: insertCount,
		},
	})
}

// CreateBatch method creates new record(s) by placeholder values from copy-create-query
// TODO: resolve sql-values parsing error, for create_batch & create_copy
// Error updating information/record(s) | Error updating record(s): ERROR: invalid input syntax for type json (SQLSTATE 22P02) 304
func (crud Crud) CreateBatch(createRecs mctypes.ActionParamsType, tableFields []string) mcresponse.ResponseMessage {
	// create from createRecs (actionParams)
	fmt.Printf("action-params: %#v \n\n", createRecs)
	// compute query
	createQuery, qErr := helper.ComputeCreateCopyQuery(crud.TableName, tableFields, createRecs)
	if qErr != nil {
		return mcresponse.GetResMessage("insertError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing create-query: %v", qErr.Error()),
			Value:   nil,
		})
	}
	fmt.Printf("create-batch-query: %v \n", createQuery)
	// perform create/insert action, via transaction/copy-protocol:
	tx, txErr := crud.AppDb.Begin(context.Background())
	if txErr != nil {
		return mcresponse.GetResMessage("insertError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error creating new record(s): %v", txErr.Error()),
			Value:   nil,
		})
	}
	fmt.Printf("transaction-start\n\n")
	defer tx.Rollback(context.Background())

	// perform records' creation
	insertCount := 0
	var insertIds []string
	var insertId string
	for _, insertValues := range createQuery.FieldValues {
		fmt.Printf("query: %v\n\n", createQuery.CreateQuery)
		fmt.Printf("query-value: %v \n\n", insertValues)
		insertErr := tx.QueryRow(context.Background(), createQuery.CreateQuery, insertValues...).Scan(&insertId)
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
	fmt.Printf("before-commit\n\n")
	// commit
	txcErr := tx.Commit(context.Background())
	if txcErr != nil {
		_ = tx.Rollback(context.Background())
		return mcresponse.GetResMessage("insertError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error creating new record(s): %v", txcErr.Error()),
			Value:   nil,
		})
	}
	//fmt.Println("before-log")
	//fmt.Println("")
	// perform audit-log
	logMessage := ""
	if crud.LogCreate {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:  crud.TableName,
			LogRecords: crud.ActionParams,
		}
		//fmt.Printf("\n***Audit-Table***: %v\n\n", crud.AuditTable)
		if logRes, logErr := crud.TransLog.AuditLog(tasks.Create, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: logMessage,
		Value: InsertedResultType{
			TableName:   crud.TableName,
			RecordIds:   insertIds,
			RecordCount: insertCount,
		},
	})
}

// CreateCopy method creates new record(s) using Pg CopyFrom
func (crud Crud) CreateCopy(createRecs mctypes.ActionParamsType, tableFields []string) mcresponse.ResponseMessage {
	// create from createRecs (actionParams)
	//fmt.Printf("action-params: %#v \n\n", createRecs)
	// compute query
	createQuery, qErr := helper.ComputeCreateCopyQuery(crud.TableName, tableFields, createRecs)
	if qErr != nil {
		return mcresponse.GetResMessage("insertError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing create-query: %v", qErr.Error()),
			Value:   nil,
		})
	}
	fmt.Printf("create-query: %v \n", createQuery)
	fmt.Printf("create-query-fields: %v \n", createQuery.FieldNames)
	fmt.Printf("create-query-values: %v \n\n", createQuery.FieldValues)
	// perform create/insert action, via transaction/copy-protocol:
	tx, txErr := crud.AppDb.Begin(context.Background())
	if txErr != nil {
		return mcresponse.GetResMessage("insertError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error creating new record(s): %v", txErr.Error()),
			Value:   nil,
		})
	}
	fmt.Printf("transaction-start\n\n")
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
	fmt.Printf("before-commit\n\n")
	// commit
	txcErr := tx.Commit(context.Background())
	if txcErr != nil {
		_ = tx.Rollback(context.Background())
		return mcresponse.GetResMessage("insertError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error creating new record(s): %v", txcErr.Error()),
			Value:   nil,
		})
	}
	fmt.Println("before-log")
	fmt.Println("")
	// perform audit-log
	logMessage := ""
	if crud.LogCreate {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:  crud.TableName,
			LogRecords: crud.ActionParams,
		}
		if logRes, logErr := crud.TransLog.AuditLog(tasks.Create, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: logMessage,
		Value:   copyCount,
	})
}

// Update method updates existing record(s)
func (crud Crud) Update(updateRecs mctypes.ActionParamsType, tableFields []string, recordIds []string) mcresponse.ResponseMessage {
	// get current records, for audit-log
	if crud.LogUpdate {
		if getQuery, err := helper.ComputeSelectQueryById(crud.TableName, recordIds, tableFields); err != nil {
			return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error computing select/read-query: %v", err.Error()),
				Value:   getQuery,
			})
		} else {
			rows, err := crud.AppDb.Query(context.Background(), getQuery)
			if err != nil {
				errMsg := fmt.Sprintf("Db query Error: %v", err.Error())
				return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
					Message: errMsg,
					Value:   nil,
				})
			}
			defer rows.Close()
			// check rows count
			var rowCount = 0
			for rows.Next() {
				var id string
				//var getRecords []interface{}
				if err := rows.Scan(&id); err == nil {
					rowCount += 1
					// parse the current-records for audit-log
					if parseVal, err := helper.ParseRawValues(rows.RawValues()); err != nil {
						// capture recordId
						crud.CurrentRecords = append(crud.CurrentRecords, id)
						//return mcresponse.GetResMessage("parseError", mcresponse.ResponseMessageOptions{
						//	Message: fmt.Sprintf("Error parsing raw-record-values: %v", err.Error()),
						//	Value:   nil,
						//})
					} else {
						// update instance CurrentRecords
						crud.CurrentRecords = append(crud.CurrentRecords, parseVal)
					}
				}
			}
			// exit if currentRec-length is less than recordIds-length
			if rowCount < len(crud.RecordIds) {
				return mcresponse.GetResMessage("fewRecords", mcresponse.ResponseMessageOptions{
					Message: fmt.Sprintf("Fewer records (%v) less than expected (%v)", rowCount, len(crud.RecordIds)),
					Value:   nil,
				})
			}
			if err := rows.Err(); err != nil {
				return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
					Message: fmt.Sprintf("Error reading/getting records: %v", err.Error()),
					Value:   nil,
				})
			}
		}
	}

	// create from updatedRecs (actionParams)
	updateQuery, err := helper.ComputeUpdateQuery(crud.TableName, tableFields, updateRecs)
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
	// perform audit-log
	logMessage := ""
	if crud.LogUpdate {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName: crud.TableName,
			LogRecords: LogRecordsType{
				TableFields:  tableFields,
				TableRecords: crud.CurrentRecords,
			},
			NewLogRecords: crud.ActionParams,
		}
		if logRes, logErr := crud.TransLog.AuditLog(tasks.Update, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: logMessage,
		Value:   updateCount,
	})
}

// UpdateById method updates existing records (in batch) that met the specified record-id(s)
func (crud Crud) UpdateById(updateRecs mctypes.ActionParamsType, tableFields []string) mcresponse.ResponseMessage {
	// get current records, for audit-log
	if crud.LogUpdate {
		if getQuery, err := helper.ComputeSelectQueryById(crud.TableName, crud.RecordIds, tableFields); err != nil {
			return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error computing select/read-query: %v", err.Error()),
				Value:   getQuery,
			})
		} else {
			rows, err := crud.AppDb.Query(context.Background(), getQuery)
			if err != nil {
				errMsg := fmt.Sprintf("Db query Error: %v", err.Error())
				return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
					Message: errMsg,
					Value:   nil,
				})
			}
			defer rows.Close()
			// check rows count
			var rowCount = 0
			for rows.Next() {
				var id string
				if err := rows.Scan(&id); err == nil {
					rowCount += 1
					// parse the current-records for audit-log
					if parseVal, err := helper.ParseRawValues(rows.RawValues()); err != nil {
						return mcresponse.GetResMessage("parseError", mcresponse.ResponseMessageOptions{
							Message: fmt.Sprintf("Error parsing raw-record-values: %v", err.Error()),
							Value:   nil,
						})
					} else {
						// update instance CurrentRecords
						crud.CurrentRecords = append(crud.CurrentRecords, parseVal)
					}
				}
			}
			// exit if currentRec-length is less than recordIds-length
			if rowCount < len(crud.RecordIds) {
				return mcresponse.GetResMessage("fewRecords", mcresponse.ResponseMessageOptions{
					Message: fmt.Sprintf("Fewer records (%v) less than expected (%v)", rowCount, len(crud.RecordIds)),
					Value:   nil,
				})
			}
			if err := rows.Err(); err != nil {
				return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
					Message: fmt.Sprintf("Error reading/getting records: %v", err.Error()),
					Value:   nil,
				})
			}
		}
	}

	// create from updatedRecs (actionParams)
	updateQuery, err := helper.ComputeUpdateQueryById(crud.TableName, tableFields, updateRecs, crud.RecordIds)
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
	// perform audit-log
	logMessage := ""
	if crud.LogUpdate {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName: crud.TableName,
			LogRecords: LogRecordsType{
				TableFields:  tableFields,
				TableRecords: crud.CurrentRecords,
			},
			NewLogRecords: crud.ActionParams,
		}
		if logRes, logErr := crud.TransLog.AuditLog(tasks.Update, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: logMessage,
		Value:   commandTag.RowsAffected(),
	})

}

// UpdateByParam method updates existing records (in batch) that met the specified query-params or where conditions
func (crud Crud) UpdateByParam(updateRecs mctypes.ActionParamsType, tableFields []string) mcresponse.ResponseMessage {
	// get current records, for audit-log
	if crud.LogUpdate {
		if getQuery, err := helper.ComputeSelectQueryByParam(crud.TableName, crud.QueryParams, tableFields); err != nil {
			return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error computing select/read-query: %v", err.Error()),
				Value:   getQuery,
			})
		} else {
			// exit if currentRec-length is less than recordIds-length
			rows, err := crud.AppDb.Query(context.Background(), getQuery)
			if err != nil {
				errMsg := fmt.Sprintf("Db query Error: %v", err.Error())
				return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
					Message: errMsg,
					Value:   nil,
				})
			}
			defer rows.Close()
			// check rows count
			var rowCount = 0
			for rows.Next() {
				var id string
				if err := rows.Scan(&id); err == nil {
					rowCount += 1
					// parse the current-records for audit-log
					if parseVal, err := helper.ParseRawValues(rows.RawValues()); err != nil {
						return mcresponse.GetResMessage("parseError", mcresponse.ResponseMessageOptions{
							Message: fmt.Sprintf("Error parsing raw-record-values: %v", err.Error()),
							Value:   nil,
						})
					} else {
						// update instance CurrentRecords
						crud.CurrentRecords = append(crud.CurrentRecords, parseVal)
					}
				}
			}
			if err := rows.Err(); err != nil {
				return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
					Message: fmt.Sprintf("Error reading/getting records: %v", err.Error()),
					Value:   nil,
				})
			}
		}
	}

	// create from updatedRecs (actionParams)
	updateQuery, err := helper.ComputeUpdateQueryByParam(crud.TableName, tableFields, updateRecs, crud.QueryParams)
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

	// perform audit-log
	logMessage := ""
	if crud.LogUpdate {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName: crud.TableName,
			LogRecords: LogRecordsType{
				TableFields:  tableFields,
				TableRecords: crud.CurrentRecords,
			},
			NewLogRecords: crud.ActionParams,
		}
		if logRes, logErr := crud.TransLog.AuditLog(tasks.Update, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: logMessage,
		Value:   commandTag.RowsAffected(),
	})
}
