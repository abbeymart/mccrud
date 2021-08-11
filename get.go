// @Author: abbeymart | Abi Akindele | @Created: 2020-12-01 | @Updated: 2020-12-01
// @Company: mConnect.biz | @License: MIT
// @Description: get / query record(s)

package mccrud

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/abbeymart/mcauditlog"
	"github.com/abbeymart/mccache"
	"github.com/abbeymart/mccrud/helper"
	"github.com/abbeymart/mcresponse"
)

// GetById method fetches/gets/reads record that met the specified record-id,
// constrained by optional skip and limit parameters
func (crud *Crud) GetById(modelRef interface{}, id string) mcresponse.ResponseMessage {
	// check cache
	getCacheRes := mccache.GetHashCache(crud.CacheKey, crud.TableName)
	val, ok := getCacheRes.Value.([]interface{})
	if getCacheRes.Ok && ok && len(val) > 0 {
		return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
			Message: "records successfully retrieved from the cache",
			Value:   val,
		})
	}
	logMessage := ""
	getQuery, err := helper.ComputeSelectQueryById(modelRef, crud.TableName, id)
	if err != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing select/read-query: %v", err.Error()),
			Value:   nil,
		})
	}
	// include options: limit 
	if crud.Limit > 0 {
		getQuery += fmt.Sprintf(" LIMIT %v", crud.Limit)
	}
	if crud.Skip > 0 {
		getQuery += fmt.Sprintf(" OFFSET %v", crud.Skip)
	}

	// totalRecordsCount from the table
	var totalRows int
	countQuery := fmt.Sprintf("SELECT COUNT(*) AS totalRows FROM %v", crud.TableName)
	countRows, tRowErr := crud.AppDb.Query(context.Background(), countQuery)
	if tRowErr != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Db query Error: %v", tRowErr.Error()),
			Value:   nil,
		})
	}
	for countRows.Next() {
		cErr := countRows.Scan(&totalRows)
		if cErr != nil {
			return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Db query Error: %v", cErr.Error()),
				Value:   nil,
			})
		}
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
	for rows.Next() {
		if rowScanErr := rows.Scan(&modelRef); rowScanErr != nil {
			return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error reading/getting records[row-scan]: %v", rowScanErr.Error()),
				Value:   nil,
			})
		} else {
			// get snapshot value from the model-struct-record(pointer) | transform value to json-value-format
			jByte, jErr := json.Marshal(modelRef)
			if jErr != nil {
				return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
					Message: fmt.Sprintf("Error transforming json-value to map-record-format: %v", jErr.Error()),
					Value:   nil,
				})
			}
			var gValue map[string]interface{}
			jErr = json.Unmarshal(jByte, &gValue)
			if jErr != nil {
				return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
					Message: fmt.Sprintf("Error transforming result-value into json-value-format: %v", jErr.Error()),
					Value:   nil,
				})
			}
			getResults = append(getResults, gValue)
			rowCount += 1
		}
	}

	if rowErr := rows.Err(); rowErr != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error reading/getting records: %v", rowErr.Error()),
			Value: GetResultType{
				Records:  nil,
				Stats:    GetStatType{},
				TaskType: crud.TaskType,
				LogRes:   mcresponse.ResponseMessage{},
			},
		})
	}

	// perform audit-log
	logRes := mcresponse.ResponseMessage{}
	var logErr error
	if crud.LogRead {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:  crud.TableName,
			LogRecords: crud.QueryParams,
		}
		if logRes, logErr = crud.TransLog.AuditLog(ReadTask, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}

	getValue := GetResultType{
		Records: getResults,
		Stats: GetStatType{
			Skip:              crud.Skip,
			Limit:             crud.Limit,
			RecordsCount:      rowCount,
			TotalRecordsCount: totalRows,
			QueryParam:        crud.QueryParams,
			RecordIds:         crud.RecordIds,
		},
		TaskType: crud.TaskType,
		LogRes:   logRes,
	}

	// update cache
	_ = mccache.SetHashCache(crud.CacheKey, crud.TableName, getValue, uint(crud.CacheExpire))

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: logMessage,
		Value:   getValue,
	})
}

// GetByIds method fetches/gets/reads records that met the specified record-ids,
// constrained by optional skip and limit parameters
func (crud Crud) GetByIds(modelRef interface{}) mcresponse.ResponseMessage {
	// check cache
	getCacheRes := mccache.GetHashCache(crud.TableName, crud.CacheKey)
	val, ok := getCacheRes.Value.([]interface{})
	if getCacheRes.Ok && ok && len(val) > 0 {
		return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
			Message: "records successfully retrieved from the cache",
			Value:   val,
		})
	}
	if len(crud.RecordIds) < 1 {
		return mcresponse.GetResMessage("paramsError",
			mcresponse.ResponseMessageOptions{
				Message: "recordIds param is required to get-record-by-id",
				Value:   nil,
			})
	}
	logMessage := ""
	getQuery, err := helper.ComputeSelectQueryByIds(modelRef, crud.TableName, crud.RecordIds)
	if err != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing select/read-query: %v", err.Error()),
			Value:   nil,
		})
	}
	// include options: limit 
	if crud.Limit > 0 {
		getQuery += fmt.Sprintf(" LIMIT %v", crud.Limit)
	}
	if crud.Skip > 0 {
		getQuery += fmt.Sprintf(" OFFSET %v", crud.Skip)
	}

	// totalRecordsCount from the table
	var totalRows int
	countQuery := fmt.Sprintf("SELECT COUNT(*) AS totalRows FROM %v", crud.TableName)
	countRows, tRowErr := crud.AppDb.Query(context.Background(), countQuery)
	if tRowErr != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Db query Error: %v", tRowErr.Error()),
			Value:   nil,
		})
	}
	for countRows.Next() {
		cErr := countRows.Scan(&totalRows)
		if cErr != nil {
			return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Db query Error: %v", cErr.Error()),
				Value:   nil,
			})
		}
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
	for rows.Next() {
		if rowScanErr := rows.Scan(&modelRef); rowScanErr != nil {
			return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error reading/getting records[row-scan]: %v", rowScanErr.Error()),
				Value:   nil,
			})
		} else {
			// get snapshot value from the model-struct-record(pointer) | transform value to json-value-format
			jByte, jErr := json.Marshal(modelRef)
			if jErr != nil {
				return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
					Message: fmt.Sprintf("Error transforming result-value into json-value-format: %v", jErr.Error()),
					Value:   nil,
				})
			}
			var gValue map[string]interface{}
			jErr = json.Unmarshal(jByte, &gValue)
			if jErr != nil {
				return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
					Message: fmt.Sprintf("Error transforming json-value to map-record-format: %v", jErr.Error()),
					Value:   nil,
				})
			}
			getResults = append(getResults, gValue)
			rowCount += 1
		}
	}

	if rowErr := rows.Err(); rowErr != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error reading/getting records: %v", rowErr.Error()),
			Value: GetResultType{
				Records:  nil,
				Stats:    GetStatType{},
				TaskType: crud.TaskType,
				LogRes:   mcresponse.ResponseMessage{},
			},
		})
	}

	// perform audit-log
	logRes := mcresponse.ResponseMessage{}
	var logErr error
	if crud.LogRead {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:  crud.TableName,
			LogRecords: crud.QueryParams,
		}
		if logRes, logErr = crud.TransLog.AuditLog(ReadTask, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}

	getValue := GetResultType{
		Records: getResults,
		Stats: GetStatType{
			Skip:              crud.Skip,
			Limit:             crud.Limit,
			RecordsCount:      rowCount,
			TotalRecordsCount: totalRows,
			QueryParam:        crud.QueryParams,
			RecordIds:         crud.RecordIds,
		},
		TaskType: crud.TaskType,
		LogRes:   logRes,
	}

	// update cache
	_ = mccache.SetHashCache(crud.CacheKey, crud.TableName, getValue, uint(crud.CacheExpire))

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: logMessage,
		Value:   getValue,
	})
}

// GetByParam method fetches/gets/reads records that met the specified query-params or where conditions,
// constrained by optional skip and limit parameters
func (crud *Crud) GetByParam(modelRef interface{}) mcresponse.ResponseMessage {
	// check cache
	getCacheRes := mccache.GetHashCache(crud.TableName, crud.CacheKey)
	val, ok := getCacheRes.Value.([]interface{})
	if getCacheRes.Ok && ok && len(val) > 0 {
		return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
			Message: "records successfully retrieved from the cache",
			Value:   val,
		})
	}

	logMessage := ""
	getQueryObj, err := helper.ComputeSelectQueryByParam(modelRef, crud.TableName, crud.QueryParams)
	if err != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing select/read-query: %v", err.Error()),
			Value:   nil,
		})
	}
	getQuery := getQueryObj.SelectQuery
	// include options: limit 
	if crud.Limit > 0 {
		getQuery += fmt.Sprintf(" LIMIT %v", crud.Limit)
	}
	if crud.Skip > 0 {
		getQuery += fmt.Sprintf(" OFFSET %v", crud.Skip)
	}
	// include where-condition / placeholders (add to query-exec)
	getQuery += getQueryObj.WhereQuery.WhereQuery

	// totalRecordsCount from the table
	var totalRows int
	countQuery := fmt.Sprintf("SELECT COUNT(*) AS totalRows FROM %v", crud.TableName)
	countRows, tRowErr := crud.AppDb.Query(context.Background(), countQuery)
	if tRowErr != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Db query Error: %v", tRowErr.Error()),
			Value:   nil,
		})
	}
	for countRows.Next() {
		cErr := countRows.Scan(&totalRows)
		if cErr != nil {
			return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Db query Error: %v", cErr.Error()),
				Value:   nil,
			})
		}
	}

	// perform crud-task action
	rows, qRowErr := crud.AppDb.Query(context.Background(), getQuery, getQueryObj.WhereQuery.FieldValues...)
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
	for rows.Next() {
		if rowScanErr := rows.Scan(&modelRef); rowScanErr != nil {
			return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error reading/getting records[row-scan]: %v", rowScanErr.Error()),
				Value:   nil,
			})
		} else {
			// get snapshot value from the model-struct-record(pointer) | transform value to json-value-format
			jByte, jErr := json.Marshal(modelRef)
			if jErr != nil {
				return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
					Message: fmt.Sprintf("Error transforming result-value into json-value-format: %v", jErr.Error()),
					Value:   nil,
				})
			}
			var gValue map[string]interface{}
			jErr = json.Unmarshal(jByte, &gValue)
			if jErr != nil {
				return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
					Message: fmt.Sprintf("Error transforming json-value to map-record-format: %v", jErr.Error()),
					Value:   nil,
				})
			}
			getResults = append(getResults, gValue)
			rowCount += 1
		}
	}

	if rowErr := rows.Err(); rowErr != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error reading/getting records: %v", rowErr.Error()),
			Value: GetResultType{
				Records:  nil,
				Stats:    GetStatType{},
				TaskType: crud.TaskType,
				LogRes:   mcresponse.ResponseMessage{},
			},
		})
	}

	// perform audit-log
	logRes := mcresponse.ResponseMessage{}
	var logErr error
	if crud.LogRead {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:  crud.TableName,
			LogRecords: crud.QueryParams,
		}
		if logRes, logErr = crud.TransLog.AuditLog(ReadTask, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}

	getValue := GetResultType{
		Records: getResults,
		Stats: GetStatType{
			Skip:              crud.Skip,
			Limit:             crud.Limit,
			RecordsCount:      rowCount,
			TotalRecordsCount: totalRows,
			QueryParam:        crud.QueryParams,
			RecordIds:         crud.RecordIds,
		},
		TaskType: crud.TaskType,
		LogRes:   logRes,
	}

	// update cache
	_ = mccache.SetHashCache(crud.CacheKey, crud.TableName, getValue, uint(crud.CacheExpire))

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: logMessage,
		Value:   getValue,
	})
}

// GetAll method fetches/gets/reads all record(s), constrained by optional skip and limit parameters
func (crud *Crud) GetAll(modelRef interface{}) mcresponse.ResponseMessage {
	// compute select-query
	getQuery, err := helper.ComputeSelectQueryAll(modelRef, crud.TableName)
	if err != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing select/read-query: %v", err.Error()),
			Value:   getQuery,
		})
	}

	// include options: limit 
	if crud.Limit > 0 {
		getQuery += fmt.Sprintf(" LIMIT %v", crud.Limit)
	}
	if crud.Skip > 0 {
		getQuery += fmt.Sprintf(" OFFSET %v", crud.Skip)
	}

	// totalRecordsCount from the table
	var totalRows int
	countQuery := fmt.Sprintf("SELECT COUNT(*) AS totalRows FROM %v", crud.TableName)
	countRows, tRowErr := crud.AppDb.Query(context.Background(), countQuery)
	if tRowErr != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Db query Error: %v", tRowErr.Error()),
			Value:   nil,
		})
	}
	for countRows.Next() {
		cErr := countRows.Scan(&totalRows)
		if cErr != nil {
			return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Db query Error: %v", cErr.Error()),
				Value:   nil,
			})
		}
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
	for rows.Next() {
		rowScanErr := rows.Scan(&modelRef)
		if rowScanErr != nil {
			return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error reading/getting records[row-scan]: %v", rowScanErr.Error()),
				Value:   nil,
			})
		}
		// get snapshot value from the model-struct-record(pointer) | transform value to json-value-format
		jByte, jErr := json.Marshal(modelRef)
		if jErr != nil {
			return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error transforming result-value into json-value-format: %v", jErr.Error()),
				Value:   nil,
			})
		}
		var gValue map[string]interface{}
		jErr = json.Unmarshal(jByte, &gValue)
		if jErr != nil {
			return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error transforming json-value to map-record-format: %v", jErr.Error()),
				Value:   nil,
			})
		}
		getResults = append(getResults, gValue)
		rowCount += 1

	}

	if rowErr := rows.Err(); rowErr != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error reading/getting records: %v", rowErr.Error()),
			Value: GetResultType{
				Records:  nil,
				Stats:    GetStatType{},
				TaskType: crud.TaskType,
				LogRes:   mcresponse.ResponseMessage{},
			},
		})
	}

	// perform audit-log | initialize log-variables
	logMessage := ""
	logRes := mcresponse.ResponseMessage{}
	var logErr error
	if crud.LogRead {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:  crud.TableName,
			LogRecords: crud.QueryParams,
		}
		if logRes, logErr = crud.TransLog.AuditLog(ReadTask, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}

	getValue := GetResultType{
		Records: getResults,
		Stats: GetStatType{
			Skip:              crud.Skip,
			Limit:             crud.Limit,
			RecordsCount:      rowCount,
			TotalRecordsCount: totalRows,
			QueryParam:        crud.QueryParams,
			RecordIds:         crud.RecordIds,
		},
		TaskType: crud.TaskType,
		LogRes:   logRes,
	}

	// update cache | *****don't cache all-table-records, due to large/unknown size*****
	//_ = mccache.SetHashCache(crud.CacheKey, crud.TableName, getResults, uint(crud.CacheExpire))

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: logMessage,
		Value:   getValue,
	})
}
