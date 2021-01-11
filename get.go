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
	"time"
)

// GetById method fetches/gets/reads record(s) that met the specified record-id(s),
// constrained by optional skip and limit parameters
func (crud *Crud) GetById(tableFields []string, tableFieldPointers []interface{}) mcresponse.ResponseMessage {
	// SELECT/scan to tableFieldPointers, in order specified by the tableFields
	// tableFields and tableFieldPointers length and order must match
	if len(tableFields) != len(tableFieldPointers) {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("tableFields Count [%v] and tableFieldPointer Count [%v] must be the same", len(tableFields), len(tableFieldPointers)),
			Value:   nil,
		})
	}
	getQuery, err := helper.ComputeSelectQueryById(crud.TableName, crud.RecordIds, tableFields)
	if err != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing select/read-query: %v", err.Error()),
			Value:   getQuery,
		})
	}
	// include options: limit... TODO: sort?
	if crud.Limit > 0 {
		getQuery += fmt.Sprintf(" LIMIT %v", crud.Limit)
	}
	// perform crud-task action
	rows, qRowErr := crud.AppDb.Query(context.Background(), getQuery)
	if qRowErr != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Db query Error: %v", qRowErr.Error()),
			Value:   nil,
		})
	}
	defer rows.Close()
	// check rows count
	var rowCount = 0
	var getResults []interface{}
	var getResult = map[string]interface{}{}
	for rows.Next() {
		if rowScanErr := rows.Scan(tableFieldPointers...); rowScanErr != nil {
			return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error reading/getting records[row-scan]: %v", rowScanErr.Error()),
				Value:   nil,
			})
		} else {
			// extract values from tableFieldPointers
			for i, fieldPointer := range tableFieldPointers {
				switch fieldPointer.(type) {
				case *time.Time:
					val := fieldPointer.(*time.Time)
					getResult[tableFields[i]] = *val
				case *string:
					val := fieldPointer.(*string)
					getResult[tableFields[i]] = *val
				case *int:
					val := fieldPointer.(*int)
					getResult[tableFields[i]] = *val
				case *float64:
					val := fieldPointer.(*float64)
					getResult[tableFields[i]] = *val
				case *interface{}:
					val := fieldPointer.(*interface{})
					getResult[tableFields[i]] = *val
				default:
					// avoid panic, return unsupported type
					return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
						Message: fmt.Sprintf("Unsupportted fieldName [%v] type %v", tableFields[i], fieldPointer),
						Value:   nil,
					})
				}
			}
			// getChan <- rowCount // pass the scanned result alert to getChan | will block until read
			getResults = append(getResults, getResult)
			rowCount += 1
		}
	}
	// close channel
	//close(getChan)

	if err := rows.Err(); err != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error reading/getting records: %v", err.Error()),
			Value:   nil,
		})
	}

	// perform audit-log
	logMessage := ""
	if crud.LogRead {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName: crud.TableName,
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
		Value: GetResultType{
			TableName:    crud.TableName,
			QueryParam:   crud.QueryParams,
			RecordIds:    crud.RecordIds,
			RecordCount:  rowCount,
			RecordValues: getResults,
		},
	})
}

// GetByParam method fetches/gets/reads record(s) that met the specified query-params or where conditions,
// constrained by optional skip and limit parameters
func (crud *Crud) GetByParam(tableFields []string, tableFieldPointers []interface{}) mcresponse.ResponseMessage {
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
	// include options: limit TODO: sort?
	if crud.Limit > 0 {
		getQuery += fmt.Sprintf(" LIMIT %v", crud.Limit)
	}
	// perform crud-task action
	fmt.Printf("getQuery-param: %v\n", getQuery)
	rows, qRowErr := crud.AppDb.Query(context.Background(), getQuery)
	if qRowErr != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Db query Error: %v", qRowErr.Error()),
			Value:   nil,
		})
	}
	defer rows.Close()
	// check rows count
	var rowCount = 0
	var getResults []interface{}
	var getResult = map[string]interface{}{}
	for rows.Next() {
		if rowScanErr := rows.Scan(tableFieldPointers...); rowScanErr != nil {
			return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error reading/getting records[row-scan]: %v", rowScanErr.Error()),
				Value:   nil,
			})
		} else {
			// extract values from tableFieldPointers
			for i, fieldPointer := range tableFieldPointers {
				switch fieldPointer.(type) {
				case *time.Time:
					val := fieldPointer.(*time.Time)
					getResult[tableFields[i]] = *val
				case *string:
					val := fieldPointer.(*string)
					getResult[tableFields[i]] = *val
				case *int:
					val := fieldPointer.(*int)
					getResult[tableFields[i]] = *val
				case *float64:
					val := fieldPointer.(*float64)
					getResult[tableFields[i]] = *val
				case *interface{}:
					val := fieldPointer.(*interface{})
					getResult[tableFields[i]] = *val
				default:
					// avoid panic, return unsupported type
					return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
						Message: fmt.Sprintf("Unsupportted fieldName [%v] type %v", tableFields[i], fieldPointer),
						Value:   nil,
					})
				}
			}
			getResults = append(getResults, getResult)
			rowCount += 1
		}
	}

	if rowErr := rows.Err(); rowErr != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error reading/getting records: %v", rowErr.Error()),
			Value: GetResultType{
				TableName:    crud.TableName,
				QueryParam:   crud.QueryParams,
				RecordIds:    crud.RecordIds,
				RecordCount:  rowCount,
				RecordValues: getResults,
			},
		})
	}

	// perform audit-log
	if crud.LogRead {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName: crud.TableName,
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
		Value: GetResultType{
			TableName:    crud.TableName,
			QueryParam:   crud.QueryParams,
			RecordIds:    crud.RecordIds,
			RecordCount:  rowCount,
			RecordValues: getResults,
		},
	})
}

// GetAll method fetches/gets/reads all record(s), constrained by optional skip and limit parameters
func (crud *Crud) GetAll(tableFields []string, tableFieldPointers []interface{}) mcresponse.ResponseMessage {
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
	// include options: skip && limit TODO: sort?
	if crud.Limit > 0 {
		getQuery += fmt.Sprintf(" LIMIT %v", crud.Limit)
	}
	if crud.Skip > 0 {
		getQuery += fmt.Sprintf(" OFFSET %v", crud.Skip)
	}
	// perform crud-task action
	rows, qRowErr := crud.AppDb.Query(context.Background(), getQuery)
	if qRowErr != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Db query Error: %v", qRowErr.Error()),
			Value:   nil,
		})
	}
	defer rows.Close()
	// check rows count
	var rowCount = 0
	var getResults []interface{}
	var getResult = map[string]interface{}{}
	for rows.Next() {
		if rowScanErr := rows.Scan(tableFieldPointers...); rowScanErr != nil {
			return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error reading/getting records[row-scan]: %v", rowScanErr.Error()),
				Value:   nil,
			})
		} else {
			// extract values from tableFieldPointers
			for i, fieldPointer := range tableFieldPointers {
				switch fieldPointer.(type) {
				case *time.Time:
					val := fieldPointer.(*time.Time)
					getResult[tableFields[i]] = *val
				case *string:
					val := fieldPointer.(*string)
					getResult[tableFields[i]] = *val
				case *int:
					val := fieldPointer.(*int)
					getResult[tableFields[i]] = *val
				case *float64:
					val := fieldPointer.(*float64)
					getResult[tableFields[i]] = *val
				case *interface{}:
					val := fieldPointer.(*interface{})
					getResult[tableFields[i]] = *val
				default:
					// avoid panic, return unsupported type
					return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
						Message: fmt.Sprintf("Unsupportted fieldName [%v] type %v", tableFields[i], fieldPointer),
						Value:   nil,
					})
				}
			}
			getResults = append(getResults, getResult)
			rowCount += 1
		}
	}

	if rowErr := rows.Err(); rowErr != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error reading/getting records: %v", rowErr.Error()),
			Value:   nil,
		})
	}

	// perform audit-log
	if crud.LogRead {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName: crud.TableName,
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
		Value: GetResultType{
			TableName:    crud.TableName,
			QueryParam:   crud.QueryParams,
			RecordIds:    crud.RecordIds,
			RecordCount:  rowCount,
			RecordValues: getResults,
		},
	})
}
