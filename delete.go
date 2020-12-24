// @Author: abbeymart | Abi Akindele | @Created: 2020-12-01 | @Updated: 2020-12-01
// @Company: mConnect.biz | @License: MIT
// @Description: delete or remove record(s)

package mccrud

import (
	"context"
	"fmt"
	"github.com/abbeymart/mcauditlog"
	"github.com/abbeymart/mccrud/helper"
	"github.com/abbeymart/mcresponse"
	"time"
)

func (crud Crud) DeleteById() mcresponse.ResponseMessage {
	// get current records, for audit-log | for delete tableFields = []string{}
	if crud.LogDelete {
		if getQuery, err := helper.ComputeSelectQueryById(crud.TableName, []string{}, crud.RecordIds); err != nil {
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
					// crud.CurrentRecords = append(crud.CurrentRecords, id)
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

	// perform crud-task action, include where-
	// compute delete query by record-ids
	deleteQuery, dQErr := helper.ComputeDeleteQueryById(crud.TableName, crud.RecordIds)
	if dQErr != nil {
		return mcresponse.GetResMessage("deleteError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing delete-query: %v", dQErr.Error()),
			Value:   nil,
		})
	}
	commandTag, delErr := crud.AppDb.Exec(context.Background(), deleteQuery)
	if delErr != nil {
		return mcresponse.GetResMessage("deleteError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error deleting record(s): %v", delErr.Error()),
			Value:   nil,
		})
	}

	// perform audit-log
	logMessage := ""
	if crud.LogDelete {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:  crud.TableName,
			LogRecords: crud.CurrentRecords,
		}
		if logRes, logErr := crud.TransLog.AuditLog(CrudTasks().Delete, crud.UserInfo.UserId, auditInfo); logErr != nil {
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

func (crud Crud) DeleteByParam() mcresponse.ResponseMessage {
	// get current records, for audit-log | for delete tableFields = []string{}
	if crud.LogDelete {
		if getQuery, err := helper.ComputeSelectQueryByParam(crud.TableName, crud.QueryParams, []string{}); err != nil {
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
					// crud.CurrentRecords = append(crud.CurrentRecords, id)
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

	// perform crud-task action, include where-query(params):
	// compute delete query by query-params
	deleteQuery, dQErr := helper.ComputeDeleteQueryByParam(crud.TableName, crud.QueryParams, []string{})
	if dQErr != nil {
		return mcresponse.GetResMessage("deleteError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing delete-query: %v", dQErr.Error()),
			Value:   nil,
		})
	}
	commandTag, delErr := crud.AppDb.Exec(context.Background(), deleteQuery)
	if delErr != nil {
		return mcresponse.GetResMessage("deleteError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error deleting record(s): %v", delErr.Error()),
			Value:   nil,
		})
	}

	// perform audit-log
	logMessage := ""
	if crud.LogDelete {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:  crud.TableName,
			LogRecords: crud.CurrentRecords,
		}
		if logRes, logErr := crud.TransLog.AuditLog(CrudTasks().Delete, crud.UserInfo.UserId, auditInfo); logErr != nil {
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

// DeleteAll method removes all records in the tables. Recommended for admin-users only
// Use if and only if you know what you are doing
func (crud Crud) DeleteAll() mcresponse.ResponseMessage {
	// ***** perform DELETE-ALL-RECORDS FROM A TABLE, IF RELATIONS/CONSTRAINTS PERMIT
	// && IF-AND-ONLY-IF-YOU-KNOW-WHAT-YOU-ARE-DOING*****
	// compute delete query
	delQuery := "DELETE * FROM " + crud.TableName
	commandTag, delErr := crud.AppDb.Exec(context.Background(), delQuery)
	if delErr != nil {
		return mcresponse.GetResMessage("deleteError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error deleting record(s): %v", delErr.Error()),
			Value:   nil,
		})
	}

	// perform audit-log
	logMessage := ""
	if crud.LogDelete {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:  crud.TableName,
			LogRecords: fmt.Sprintf("All Records Delete From %v table, by User: %v [id: %v, email: %v], at %v", crud.TableName, crud.UserInfo.LoginName, crud.UserInfo.UserId, crud.UserInfo.Email, time.Now()),
		}
		if logRes, logErr := crud.TransLog.AuditLog(CrudTasks().Delete, crud.UserInfo.UserId, auditInfo); logErr != nil {
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
