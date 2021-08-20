// @Author: abbeymart | Abi Akindele | @Created: 2020-12-01 | @Updated: 2020-12-01
// @Company: mConnect.biz | @License: MIT
// @Description: delete or remove record(s)

package mccrud

import (
	"context"
	"fmt"
	"github.com/abbeymart/mcauditlog"
	"github.com/abbeymart/mccache"
	"github.com/abbeymart/mccrud/helper"
	"github.com/abbeymart/mcresponse"
)

// DeleteById method deletes or removes record(s) by record-id(s)
func (crud *Crud) DeleteById(modelRef interface{}, id string) mcresponse.ResponseMessage {
	// audit-log
	// get records to delete, for audit-log
	if crud.LogDelete || crud.LogCrud {
		getRes := crud.GetById(modelRef, id)
		value, _ := getRes.Value.(CrudResultType)
		crud.CurrentRecords = value.TableRecords
	}
	// compute delete query by record-ids
	deleteQuery, dQErr := helper.ComputeDeleteQueryById(crud.TableName, id)
	if dQErr != nil {
		return mcresponse.GetResMessage("deleteError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing delete-query: %v", dQErr.Error()),
			Value:   nil,
		})
	}
	_, delErr := crud.AppDb.Exec(context.Background(), deleteQuery)
	if delErr != nil {
		return mcresponse.GetResMessage("deleteError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error deleting record(s): %v", delErr.Error()),
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
			LogRecords: crud.CurrentRecords,
		}
		if logRes, logErr = crud.TransLog.AuditLog(CreateTask, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: fmt.Sprintf("Record(s) deleted successfully: [log-message: %v]",logMessage),
		Value: CrudResultType{
			QueryParam:  crud.QueryParams,
			//RecordCount: int(commandTag.RowsAffected()),
			TaskType:    crud.TaskType,
			LogRes:      logRes,
		},
	})
}

// DeleteAuditByIds method deletes or removes record(s) by record-id(s)
func (crud *Crud) DeleteByIds(modelRef interface{}) mcresponse.ResponseMessage {
	// audit-log
	if crud.LogDelete || crud.LogCrud {
		getRes := crud.GetByIds(modelRef)
		value, _ := getRes.Value.(CrudResultType)
		crud.CurrentRecords = value.TableRecords
	}
	// compute delete query by record-ids
	deleteQuery, dQErr := helper.ComputeDeleteQueryByIds(crud.TableName, crud.RecordIds)
	if dQErr != nil {
		return mcresponse.GetResMessage("deleteError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing delete-query: %v", dQErr.Error()),
			Value:   nil,
		})
	}
	_, delErr := crud.AppDb.Exec(context.Background(), deleteQuery)
	if delErr != nil {
		return mcresponse.GetResMessage("deleteError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error deleting record(s): %v", delErr.Error()),
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
			LogRecords: crud.CurrentRecords,
		}
		if logRes, logErr = crud.TransLog.AuditLog(CreateTask, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: fmt.Sprintf("Record(s) deleted successfully: [log-message: %v]",logMessage),
		Value: CrudResultType{
			QueryParam:  crud.QueryParams,
			//RecordCount: int(commandTag.RowsAffected()),
			TaskType:    crud.TaskType,
			LogRes:      logRes,
		},
	})
}

// DeleteByParam method deletes or removes record(s) by query-parameters or where conditions
func (crud *Crud) DeleteByParam(modelRef interface{}) mcresponse.ResponseMessage {
	// audit-log
	if crud.LogDelete || crud.LogCrud {
		getRes := crud.GetByParam(modelRef)
		value, _ := getRes.Value.(CrudResultType)
		crud.CurrentRecords = value.TableRecords
	}
	// compute delete query by query-params
	delQueryObj, dQErr := helper.ComputeDeleteQueryByParam(crud.TableName, crud.QueryParams)
	if dQErr != nil {
		return mcresponse.GetResMessage("deleteError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing delete-query: %v", dQErr.Error()),
			Value:   nil,
		})
	}
	deleteQuery := delQueryObj.DeleteQuery + delQueryObj.WhereQuery.WhereQuery
	_, delErr := crud.AppDb.Exec(context.Background(), deleteQuery, delQueryObj.WhereQuery.FieldValues...)
	if delErr != nil {
		return mcresponse.GetResMessage("deleteError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error deleting record(s): %v", delErr.Error()),
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
			LogRecords: crud.CurrentRecords,
		}
		if logRes, logErr = crud.TransLog.AuditLog(CreateTask, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: fmt.Sprintf("Record(s) deleted successfully: [log-message: %v]",logMessage),
		Value: CrudResultType{
			QueryParam:  crud.QueryParams,
			//RecordCount: int(commandTag.RowsAffected()),
			TaskType:    crud.TaskType,
			LogRes:      logRes,
		},
	})
}

// DeleteAll method deletes or removes all records in the tables. Recommended for admin-users only
// Use if and only if you know what you are doing
func (crud *Crud) DeleteAll() mcresponse.ResponseMessage {
	// ***** perform DELETE-ALL-RECORDS FROM A TABLE, IF RELATIONS/CONSTRAINTS PERMIT *****
	// ***** && IF-AND-ONLY-IF-YOU-KNOW-WHAT-YOU-ARE-DOING && AT-YOUR-OWN-RISK *****
	// compute delete query
	delQuery := fmt.Sprintf("DELETE FROM %v", crud.TableName)
	commandTag, delErr := crud.AppDb.Exec(context.Background(), delQuery)
	if delErr != nil {
		return mcresponse.GetResMessage("deleteError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error deleting record(s): %v", delErr.Error()),
			Value:   nil,
		})
	}

	// delete cache, by key (TableName)
	_ = mccache.DeleteHashCache(crud.TableName, crud.CacheKey, "key")

	// perform audit-log
	logMessage := ""
	if crud.LogDelete || crud.LogCrud {
		auditInfo := mcauditlog.PgxAuditLogOptionsType{
			TableName:  crud.TableName,
			LogRecords: map[string]string{"query_desc": "all-records"},
		}
		if logRes, logErr := crud.TransLog.AuditLog(DeleteTask, crud.UserInfo.UserId, auditInfo); logErr != nil {
			logMessage = fmt.Sprintf("Audit-log-error: %v", logErr.Error())
		} else {
			logMessage = fmt.Sprintf("Audit-log-code: %v | Message: %v", logRes.Code, logRes.Message)
		}
	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "Record(s) deleted successfully | " + logMessage,
		Value:   commandTag.Delete(),
	})

}
