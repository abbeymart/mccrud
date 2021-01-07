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
	"github.com/abbeymart/mctypes/tasks"
)

// GetById method fetches/gets/reads record(s) that met the specified record-id(s),
// constrained by optional skip and limit parameters
func (crud Crud) GetById(tableFields []string, getChan chan int, tableFieldPointers ...interface{}) mcresponse.ResponseMessage {
	// TODO: SELECT/scan to tableFieldPointers, in order specified by the tableFields
	// tableFields and tableFieldPointers length and order must match
	if len(tableFields) != len(tableFieldPointers) {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("tableFields Count [%v] and tableFieldPointer Count [%v] must be the same", len(tableFields), len(tableFieldPointers)),
			Value:   nil,
		})
	}
	logMessage := ""
	getQuery, err := helper.ComputeSelectQueryById(crud.TableName, crud.RecordIds, tableFields)
	if err != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing select/read-query: %v", err.Error()),
			Value:   getQuery,
		})
	}
	// perform crud-task action
	// include options: limit... TODO: sort?
	if crud.Limit > 0 {
		getQuery += getQuery + fmt.Sprintf(" LIMIT %v", crud.Limit)
	}
	// exit of currentRec-length is less than recordIds-length
	rows, err := crud.AppDb.Query(context.Background(), getQuery)
	if err != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Db query Error: %v", err.Error()),
			Value:   nil,
		})
	}
	defer rows.Close()
	// check rows count
	//var getRecords []interface{}
	var rowCount = 0
	for rows.Next() {
		//var id string
		if err := rows.Scan(tableFieldPointers...); err == nil {
			getChan <- rowCount		// pass the scanned result alert to getChan | will block until read
			rowCount += 1
			//getRecords = append(getRecords, tableFieldPointers)
		}
	}

	if err := rows.Err(); err != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error reading/getting records: %v", err.Error()),
			Value:   nil,
		})
	}

	// perform audit-log
	if crud.LogRead {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:  crud.TableName,
			LogRecords: LogRecordsType{
				TableFields:  tableFields,
				TableRecords: []interface{}{crud.RecordIds},
			},
		}
		if logRes, logErr := crud.TransLog.AuditLog(tasks.Read, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: logMessage,
		Value:   nil, // handle result on requester-side (i.e. struct-spec)... tableFieldPointers
	})
}

// GetByParam method fetches/gets/reads record(s) that met the specified query-params or where conditions,
// constrained by optional skip and limit parameters
func (crud Crud) GetByParam(tableFields []string, getChan chan int, tableFieldPointers ...interface{}) mcresponse.ResponseMessage {
	// SELECT/scan to tableFieldPointers, in order specified by the tableFields
	if len(tableFields) != len(tableFieldPointers) {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("tableFields Count [%v] and tableFieldPointer Count [%v] must be the same", len(tableFields), len(tableFieldPointers)),
			Value:   nil,
		})
	}
	logMessage := ""
	getQuery, err := helper.ComputeSelectQueryByParam(crud.TableName, crud.QueryParams, tableFields)
	if err != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing select/read-query: %v", err.Error()),
			Value:   getQuery,
		})
	}
	// perform crud-task action
	// include options: limit TODO: sort?
	if crud.Limit > 0 {
		getQuery += getQuery + fmt.Sprintf(" LIMIT %v", crud.Limit)
	}
	// exit of currentRec-length is less than recordIds-length
	rows, err := crud.AppDb.Query(context.Background(), getQuery)
	if err != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Db query Error: %v", err.Error()),
			Value:   nil,
		})
	}
	defer rows.Close()
	// check rows count
	var rowCount = 0
	for rows.Next() {
		//var id string
		if err := rows.Scan(tableFieldPointers...); err == nil {
			getChan <- rowCount		// pass the scanned result alert to getChan | will block until read
			rowCount += 1
		}
	}

	if err := rows.Err(); err != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error reading/getting records: %v", err.Error()),
			Value:   nil,
		})
	}

	// perform audit-log
	if crud.LogRead {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:  crud.TableName,
			LogRecords: LogRecordsType{
				TableFields:  tableFields,
				TableRecords: []interface{}{crud.QueryParams},
			},
		}
		if logRes, logErr := crud.TransLog.AuditLog(tasks.Read, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: logMessage,
		Value:   nil, // handle result on requester-side (i.e. struct-spec)... tableFieldPointers
	})
}

// GetAll method fetches/gets/reads all record(s), constrained by optional skip and limit parameters
func (crud Crud) GetAll(tableFields []string, getChan chan int, tableFieldPointers ...interface{}) mcresponse.ResponseMessage {
	// SELECT/scan to tableFieldPointers, in order specified by the tableFields
	if len(tableFields) != len(tableFieldPointers) {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("tableFields Count [%v] and tableFieldPointer Count [%v] must be the same", len(tableFields), len(tableFieldPointers)),
			Value:   nil,
		})
	}
	logMessage := ""
	getQuery, err := helper.ComputeSelectQueryAll(crud.TableName, tableFields)
	if err != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing select/read-query: %v", err.Error()),
			Value:   getQuery,
		})
	}
	// perform crud-task action
	// include options: skip && limit TODO: sort?
	if crud.Limit > 0 {
		getQuery += getQuery + fmt.Sprintf(" LIMIT %v", crud.Limit)
	}
	if crud.Skip > 0 {
		getQuery += getQuery + fmt.Sprintf(" OFFSET %v", crud.Skip)
	}
	// exit if currentRec-length is less than recordIds-length
	rows, err := crud.AppDb.Query(context.Background(), getQuery)
	if err != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Db query Error: %v", err.Error()),
			Value:   nil,
		})
	}
	defer rows.Close()
	// check rows count
	var rowCount = 0
	for rows.Next() {
		//var id string
		if err := rows.Scan(tableFieldPointers...); err == nil {
			getChan <- rowCount		// pass the scanned result alert to getChan | will block until read
			rowCount += 1
		}
	}

	if err := rows.Err(); err != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error reading/getting records: %v", err.Error()),
			Value:   nil,
		})
	}

	// perform audit-log
	if crud.LogRead {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:  crud.TableName,
			LogRecords: LogRecordsType{
				TableFields:  tableFields,
				TableRecords: []interface{}{"all-records"},
			},
		}
		if logRes, logErr := crud.TransLog.AuditLog(tasks.Read, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: logMessage,
		Value:   nil, // handle result on requester-side (i.e. struct-spec)... tableFieldPointers
	})
}
