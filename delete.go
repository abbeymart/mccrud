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
	"github.com/abbeymart/mctypes/tasks"
	"time"
)

// DeleteById method deletes or removes record(s) by record-id(s)
func (crud *Crud) DeleteById() mcresponse.ResponseMessage {
	// get current records, for audit-log | for delete tableFields = []string{}
	tabFields := []string{"id"}
	if crud.LogDelete {
		if getQuery, err := helper.ComputeSelectQueryById(crud.TableName, crud.RecordIds, tabFields); err != nil {
			return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error computing select/read-query: %v", err.Error()),
				Value:   getQuery,
			})
		} else {
			fmt.Printf("delete-by-id-get-query: %v\n", getQuery)
			rows, rqErr := crud.AppDb.Query(context.Background(), getQuery)
			if rqErr != nil {
				errMsg := fmt.Sprintf("Db query Error: %v", rqErr.Error())
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
				if getErr := rows.Scan(&id); getErr != nil {
					return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
						Message: fmt.Sprintf("Error getting records: %v", getErr.Error()),
						Value:   nil,
					})
				}
				rowCount += 1
				// parse the current-records for audit-log
				if parseVal, parseErr := helper.ParseRawValues(rows.RawValues()); parseErr != nil {
					return mcresponse.GetResMessage("parseError", mcresponse.ResponseMessageOptions{
						Message: fmt.Sprintf("Error parsing raw-record-values: %v", parseErr.Error()),
						Value:   nil,
					})
				} else {
					// update instance CurrentRecords
					crud.CurrentRecords = append(crud.CurrentRecords, parseVal)
				}
			}
			// exit if currentRec-length is less than recordIds-length
			if rowCount < len(crud.RecordIds) {
				return mcresponse.GetResMessage("fewRecords", mcresponse.ResponseMessageOptions{
					Message: fmt.Sprintf("Fewer records (%v) less than expected (%v)", rowCount, len(crud.RecordIds)),
					Value:   nil,
				})
			}
			if rowErr := rows.Err(); rowErr != nil {
				return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
					Message: fmt.Sprintf("Error reading/getting records: %v", rowErr.Error()),
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
			TableName: crud.TableName,
			LogRecords: LogRecordsType{
				TableFields:  []string{},
				TableRecords: crud.CurrentRecords,
			},
		}
		if logRes, logErr := crud.TransLog.AuditLog(tasks.Delete, crud.UserInfo.UserId, auditInfo); logErr != nil {
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

// DeleteByParam method deletes or removes record(s) by query-parameters or where conditions
func (crud *Crud) DeleteByParam() mcresponse.ResponseMessage {
	// get current records, for audit-log | for delete tableFields = []string{}
	var tabFields []string
	if crud.LogDelete {
		if getQuery, err := helper.ComputeSelectQueryByParam(crud.TableName, crud.QueryParams, tabFields); err != nil {
			return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error computing select/read-query: %v", err.Error()),
				Value:   getQuery,
			})
		} else {
			fmt.Printf("delete-by-params-get-query: %v\n", getQuery)
			rows, rqErr := crud.AppDb.Query(context.Background(), getQuery)
			if rqErr != nil {
				errMsg := fmt.Sprintf("Db query Error: %v", rqErr.Error())
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
				if getErr := rows.Scan(&id); getErr != nil {
					return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
						Message: fmt.Sprintf("Error getting records: %v", getErr.Error()),
						Value:   nil,
					})
				}
				rowCount += 1
				// parse the current-records for audit-log
				if parseVal, parseErr := helper.ParseRawValues(rows.RawValues()); parseErr != nil {
					return mcresponse.GetResMessage("parseError", mcresponse.ResponseMessageOptions{
						Message: fmt.Sprintf("Error parsing raw-record-values: %v", parseErr.Error()),
						Value:   nil,
					})
				} else {
					// update instance CurrentRecords
					crud.CurrentRecords = append(crud.CurrentRecords, parseVal)
				}
			}
			if rowErr := rows.Err(); rowErr != nil {
				return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
					Message: fmt.Sprintf("Error reading/getting records: %v", rowErr.Error()),
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
			TableName: crud.TableName,
			LogRecords: LogRecordsType{
				TableFields:  []string{},
				TableRecords: crud.CurrentRecords,
			},
		}
		if logRes, logErr := crud.TransLog.AuditLog(tasks.Delete, crud.UserInfo.UserId, auditInfo); logErr != nil {
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

// DeleteAll method deletes or removes all records in the tables. Recommended for admin-users only
// Use if and only if you know what you are doing
func (crud *Crud) DeleteAll() mcresponse.ResponseMessage {
	// ***** perform DELETE-ALL-RECORDS FROM A TABLE, IF RELATIONS/CONSTRAINTS PERMIT *****
	// ***** && IF-AND-ONLY-IF-YOU-KNOW-WHAT-YOU-ARE-DOING *****
	// compute delete query
	delQuery := fmt.Sprintf("DELETE FROM %v", crud.TableName)
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
			TableName: crud.TableName,
			LogRecords: LogRecordsType{
				TableRecords: []interface{}{fmt.Sprintf("All records deleted from %v table, by user: %v [id: %v, email: %v], at %v", crud.TableName, crud.UserInfo.LoginName, crud.UserInfo.UserId, crud.UserInfo.Email, time.Now())},
			},
		}
		if logRes, logErr := crud.TransLog.AuditLog(tasks.Delete, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: logMessage,
		Value:   commandTag.Delete(),
	})
}
