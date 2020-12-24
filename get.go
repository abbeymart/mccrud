// @Author: abbeymart | Abi Akindele | @Created: 2020-12-01 | @Updated: 2020-12-01
// @Company: mConnect.biz | @License: MIT
// @Description: get / query record(s)

package mccrud

import (
	"context"
	"fmt"
	"github.com/abbeymart/mcauditlog"
	"github.com/abbeymart/mccrud/helper"
	"github.com/abbeymart/mcresponse"
)

func (crud Crud) GetById(tableFields []string, tableFieldPointers ...interface{}) mcresponse.ResponseMessage {
	// TODO: SELECT/scan to tableFieldPointers, in order specified by the tableFields
	// i.e. tableFields and tableFieldPointers order must match
	logMessage := ""
	if getQuery, err := helper.ComputeSelectQueryById(crud.TableName, crud.RecordIds, tableFields); err != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing select/read-query: %v", err.Error()),
			Value:   getQuery,
		})
	} else {
		// perform crud-task action
		// include options: limit... TODO: sort?
		if crud.Limit > 0 {
			getQuery += getQuery + fmt.Sprintf(" LIMIT %v", crud.Limit)
		}
		// exit of currentRec-length is less than recordIds-length
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
			//var id string
			if err := rows.Scan(tableFieldPointers...); err == nil {
				rowCount += 1
			}
		}

		if err := rows.Err(); err != nil {
			return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error reading/getting records: %v", err.Error()),
				Value:   nil,	// TODO: handle result on requester-side (i.e. struct-spec)
			})
		}

		// perform audit-log
		if crud.LogRead {
			auditInfo := mcauditlog.PgxAuditLogOptionsType{
				TableName:  crud.TableName,
				LogRecords: crud.RecordIds,
			}
			if logRes, logErr := crud.TransLog.AuditLog(CrudTasks().Read, crud.UserInfo.UserId, auditInfo); logErr != nil {
				logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
			} else {
				logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
			}
		}
	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: logMessage,
		Value:   nil,
	})
}

func (crud Crud) GetByParam(tableFields []string, tableFieldPointers ...interface{}) mcresponse.ResponseMessage {
	// SELECT/scan to tableFieldPointers, in order specified by the tableFields
	// i.e. tableFields and tableFieldPointers order must match
	logMessage := ""
	if getQuery, err := helper.ComputeSelectQueryByParam(crud.TableName, crud.QueryParams, tableFields); err != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing select/read-query: %v", err.Error()),
			Value:   getQuery,
		})
	} else {
		// perform crud-task action
		// include options: limit TODO: sort?
		if crud.Skip > 0 {
			getQuery += getQuery + fmt.Sprintf(" SKIP %v", crud.Skip)
		}
		if crud.Limit > 0 {
			getQuery += getQuery + fmt.Sprintf(" LIMIT %v", crud.Limit)
		}
		// exit of currentRec-length is less than recordIds-length
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
			//var id string
			if err := rows.Scan(tableFieldPointers...); err == nil {
				rowCount += 1
			}
		}

		if err := rows.Err(); err != nil {
			return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error reading/getting records: %v", err.Error()),
				Value:   nil,	// TODO: handle result on requester-side (i.e. struct-spec)
			})
		}

		// perform audit-log
		if crud.LogRead {
			auditInfo := mcauditlog.PgxAuditLogOptionsType{
				TableName:  crud.TableName,
				LogRecords: crud.RecordIds,
			}
			if logRes, logErr := crud.TransLog.AuditLog(CrudTasks().Read, crud.UserInfo.UserId, auditInfo); logErr != nil {
				logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
			} else {
				logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
			}
		}
	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: logMessage,
		Value:   nil,
	})
}

func (crud Crud) GetAll(tableFields []string, tableFieldPointers ...interface{}) mcresponse.ResponseMessage {
	// SELECT/scan to tableFieldPointers, in order specified by the tableFields
	// i.e. tableFields and tableFieldPointers order must match
	logMessage := ""
	if getQuery, err := helper.ComputeSelectQueryAll(crud.TableName, tableFields); err != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing select/read-query: %v", err.Error()),
			Value:   getQuery,
		})
	} else {
		// perform crud-task action
		// include options: skip && limit TODO: sort?
		if crud.Skip > 0 {
			getQuery += getQuery + fmt.Sprintf(" SKIP %v", crud.Skip)
		}
		if crud.Limit > 0 {
			getQuery += getQuery + fmt.Sprintf(" LIMIT %v", crud.Limit)
		}
		// exit of currentRec-length is less than recordIds-length
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
			//var id string
			if err := rows.Scan(tableFieldPointers...); err == nil {
				rowCount += 1
			}
		}

		if err := rows.Err(); err != nil {
			return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error reading/getting records: %v", err.Error()),
				Value:   nil,	// TODO: handle result on requester-side (i.e. struct-spec)
			})
		}

		// perform audit-log
		if crud.LogRead {
			auditInfo := mcauditlog.PgxAuditLogOptionsType{
				TableName:  crud.TableName,
				LogRecords: crud.RecordIds,
			}
			if logRes, logErr := crud.TransLog.AuditLog(CrudTasks().Read, crud.UserInfo.UserId, auditInfo); logErr != nil {
				logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
			} else {
				logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
			}
		}
	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: logMessage,
		Value:   nil,
	})
}
