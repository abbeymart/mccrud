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
	"time"
)

// GetById method fetches/gets/reads record(s) that met the specified record-id(s),
// constrained by optional skip and limit parameters
func (crud *Crud) GetById(modelRef interface{}, id string) mcresponse.ResponseMessage {
	// check cache
	getCacheRes := mccache.GetHashCache(crud.TableName, crud.HashKey)
	val, ok := getCacheRes.Value.([]interface{})
	if getCacheRes.Ok && ok && len(val) > 0 {
		return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
			Message: "records successfully retrieved from the cache",
			Value: CrudResultType{
				QueryParam:   crud.QueryParams,
				RecordIds:    crud.RecordIds,
				RecordCount:  len(val),
				TableRecords: val,
			},
		})
	}
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
			// get snapshot value from the pointer | transform value to json-value-format
			jByte, jErr := json.Marshal(getResult)
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
					Message: fmt.Sprintf("Error transforming result-value into json-value-format: %v", jErr.Error()),
					Value:   nil,
				})
			}
			getResults = append(getResults, gValue)
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
	// update cache
	_ = mccache.SetHashCache(crud.TableName, crud.HashKey, getResults, uint(crud.CacheExpire))

	// perform audit-log
	logMessage := ""
	if crud.LogRead {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:  crud.TableName,
			LogRecords: crud.RecordIds,
		}
		if logRes, logErr := crud.TransLog.AuditLog(ReadTask, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: logMessage,
		Value: CrudResultType{
			QueryParam:   crud.QueryParams,
			RecordIds:    crud.RecordIds,
			RecordCount:  rowCount,
			TableRecords: getResults,
		},
	})
}

func (crud Crud) GetByIds(modelRef interface{}) mcresponse.ResponseMessage {
	// check cache
	getCacheRes := mccache.GetHashCache(crud.TableName, crud.HashKey)
	val, ok := getCacheRes.Value.([]interface{})
	if getCacheRes.Ok && ok && len(val) > 0 {
		return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
			Message: "records successfully retrieved from the cache",
			Value: CrudResultType{
				QueryParam:   crud.QueryParams,
				RecordIds:    crud.RecordIds,
				RecordCount:  len(val),
				TableRecords: val,
			},
		})
	}
	if len(crud.RecordIds) < 1 {
		return mcresponse.GetResMessage("paramsError",
			mcresponse.ResponseMessageOptions{
				Message: "recordIds param is required to get-record-by-id",
				Value:   nil,
			})
	}
	// perform get-query
	result := crud.GormDb.Limit(crud.Limit).Offset(crud.Skip).Where("id in ?", crud.RecordIds).Find(&modelRef)
	if result.Error != nil {
		return mcresponse.GetResMessage("readError",
			mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("%v", result.Error.Error()),
				Value:   nil,
			})
	}
	// rows
	var records []interface{}
	rows, err := result.Rows()
	if err != nil {
		errMsg := fmt.Sprintf("%v", result.Error.Error())
		return mcresponse.GetResMessage("readError",
			mcresponse.ResponseMessageOptions{
				Message: errMsg,
				Value:   nil,
			})
	}
	for rows.Next() {
		err = crud.GormDb.ScanRows(rows, &modelRef)
		if err != nil {
			return mcresponse.GetResMessage("readError",
				mcresponse.ResponseMessageOptions{
					Message: fmt.Sprintf("%v", result.Error.Error()),
					Value:   nil,
				})
		}
		// get snapshot value from the pointer | transform value to json-value([]byte)
		jByte, jErr := json.Marshal(modelRef)
		if jErr != nil {
			return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error transforming record(row-value) into json-value([]byte): %v", jErr.Error()),
				Value:   nil,
			})
		}
		var gValue map[string]interface{}
		jErr = json.Unmarshal(jByte, &gValue)
		if jErr != nil {
			return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error transforming json-value to result-value: %v", jErr.Error()),
				Value:   nil,
			})
		}
		records = append(records, gValue)
	}
	var totalRecordsCount int64
	var _ = crud.GormDb.Find(modelRef).Count(&totalRecordsCount)
	// logRead
	var logRes mcresponse.ResponseMessage
	if crud.LogRead {
		logRes, err = crud.TransLog.AuditLog(CrudTasks().Read, crud.UserInfo.UserId, AuditLogOptionsType{
			LogRecords: crud.RecordIds,
			TableName:  crud.TableName,
		})
		if err != nil {
			logRes = mcresponse.ResponseMessage{
				Code:    "logError",
				Message: fmt.Sprintf("Audit-log error: %v", err.Error()),
				Value:   nil,
			}
		}
	}
	return mcresponse.GetResMessage("success",
		mcresponse.ResponseMessageOptions{
			Message: "Task completed successfully",
			Value: GetResultType{
				Records: records,
				Stats: GetStatType{
					Skip:              crud.Skip,
					Limit:             crud.Limit,
					RecordsCount:      int(result.RowsAffected),
					TotalRecordsCount: int(totalRecordsCount),
				},
				LogRes: logRes,
			},
		})
}

// GetByParam method fetches/gets/reads record(s) that met the specified query-params or where conditions,
// constrained by optional skip and limit parameters
func (crud *Crud) GetByParam(modelRef interface{}) mcresponse.ResponseMessage {
	// check cache
	getCacheRes := mccache.GetHashCache(crud.TableName, crud.HashKey)
	val, ok := getCacheRes.Value.([]interface{})
	if getCacheRes.Ok && ok && len(val) > 0 {
		return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
			Message: "records successfully retrieved from the cache",
			Value: CrudResultType{
				QueryParam:   crud.QueryParams,
				RecordIds:    crud.RecordIds,
				RecordCount:  len(val),
				TableRecords: val,
			},
		})
	}

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
	//fmt.Printf("getQuery-param: %v\n", getQuery)
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
			// get snapshot value from the pointer | transform value to json-value-format
			jByte, jErr := json.Marshal(getResult)
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
			Value: CrudResultType{
				QueryParam:   crud.QueryParams,
				RecordIds:    crud.RecordIds,
				RecordCount:  rowCount,
				TableRecords: getResults,
			},
		})
	}

	// update cache
	_ = mccache.SetHashCache(crud.TableName, crud.HashKey, getResults, uint(crud.CacheExpire))

	// perform audit-log
	if crud.LogRead {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:  crud.TableName,
			LogRecords: crud.QueryParams,
		}
		if logRes, logErr := crud.TransLog.AuditLog(ReadTask, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: logMessage,
		Value: CrudResultType{
			QueryParam:   crud.QueryParams,
			RecordIds:    crud.RecordIds,
			RecordCount:  rowCount,
			TableRecords: getResults,
		},
	})

}

// GetAll method fetches/gets/reads all record(s), constrained by optional skip and limit parameters
func (crud *Crud) GetAll(modelRef interface{}) mcresponse.ResponseMessage {
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
	getResult := map[string]interface{}{}
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
			// get snapshot value from the pointer | transform value to json-value-format
			jByte, jErr := json.Marshal(getResult)
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
			Value:   nil,
		})
	}

	// perform audit-log
	if crud.LogRead {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:  crud.TableName,
			LogRecords: map[string]string{"query_desc": "all-records"},
		}
		if logRes, logErr := crud.TransLog.AuditLog(ReadTask, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: logMessage,
		Value: CrudResultType{
			QueryParam:   crud.QueryParams,
			RecordIds:    crud.RecordIds,
			RecordCount:  rowCount,
			TableRecords: getResults,
		},
	})
}
