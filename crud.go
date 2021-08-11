// @Author: abbeymart | Abi Akindele | @Created: 2020-12-01 | @Updated: 2020-12-01
// @Company: mConnect.biz | @License: MIT
// @Description: Base type/function CRUD operations for PgDB

package mccrud

import (
	"encoding/json"
	"fmt"
	"github.com/abbeymart/mcauditlog"
	"github.com/abbeymart/mcresponse"
	"time"
)

// Crud object / struct
type Crud struct {
	CrudParamsType
	CrudOptionsType
	CurrentRecords []interface{}
	TransLog       mcauditlog.PgxLogParam
	CacheKey       string // Unique for exactly the same query
}

// NewCrud constructor returns a new crud-instance
func NewCrud(params CrudParamsType, options CrudOptionsType) (crudInstance *Crud) {
	crudInstance = &Crud{}
	// compute crud params
	crudInstance.AppDb = params.AppDb
	crudInstance.TableName = params.TableName
	crudInstance.UserInfo = params.UserInfo
	crudInstance.ActionParams = params.ActionParams
	crudInstance.RecordIds = params.RecordIds
	crudInstance.QueryParams = params.QueryParams
	crudInstance.SortParams = params.SortParams
	crudInstance.ProjectParams = params.ProjectParams
	crudInstance.Token = params.Token
	crudInstance.TaskName = params.TaskName
	crudInstance.Skip = params.Skip
	crudInstance.Limit = params.Limit

	// crud options
	crudInstance.MaxQueryLimit = options.MaxQueryLimit
	crudInstance.AuditTable = options.AuditTable
	crudInstance.AccessTable = options.AccessTable
	crudInstance.RoleTable = options.RoleTable
	crudInstance.UserTable = options.UserTable
	crudInstance.ProfileTable = options.ProfileTable
	crudInstance.ServiceTable = options.ServiceTable
	crudInstance.AuditDb = options.AuditDb
	crudInstance.AccessDb = options.AccessDb
	crudInstance.LogCrud = options.LogCrud
	crudInstance.LogRead = options.LogRead
	crudInstance.LogCreate = options.LogCreate
	crudInstance.LogUpdate = options.LogUpdate
	crudInstance.LogDelete = options.LogDelete
	crudInstance.CheckAccess = options.CheckAccess // Dec 09/2020: user to implement auth as a middleware
	crudInstance.CacheExpire = options.CacheExpire // cache expire in secs
	// Default values
	if crudInstance.AuditTable == "" {
		crudInstance.AuditTable = "audits"
	}
	if crudInstance.AccessTable == "" {
		crudInstance.AccessTable = "access_keys"
	}
	if crudInstance.RoleTable == "" {
		crudInstance.RoleTable = "roles"
	}
	if crudInstance.UserTable == "" {
		crudInstance.UserTable = "users"
	}
	if crudInstance.ProfileTable == "" {
		crudInstance.ProfileTable = "profiles"
	}
	if crudInstance.ServiceTable == "" {
		crudInstance.ServiceTable = "services"
	}
	if crudInstance.AuditDb == nil {
		crudInstance.AuditDb = crudInstance.AppDb
	}
	if crudInstance.AccessDb == nil {
		crudInstance.AccessDb = crudInstance.AppDb
	}
	if crudInstance.Skip < 0 {
		crudInstance.Skip = 0
	}

	if crudInstance.MaxQueryLimit <= 0 {
		crudInstance.MaxQueryLimit = 10000
	}

	if crudInstance.Limit > crudInstance.MaxQueryLimit {
		crudInstance.Limit = crudInstance.MaxQueryLimit
	}

	if crudInstance.CacheExpire <= 0 {
		crudInstance.CacheExpire = 300 // 300 secs, 5 minutes
	}
	// Compute CacheKey from TableName, QueryParams, SortParams, ProjectParams and RecordIds
	qParam, _ := json.Marshal(params.QueryParams)
	sParam, _ := json.Marshal(params.SortParams)
	pParam, _ := json.Marshal(params.ProjectParams)
	dIds, _ := json.Marshal(params.RecordIds)
	//crudInstance.CacheKey = params.TableName + string(qParam) + string(sParam) + string(pParam) + string(dIds)
	crudInstance.CacheKey = fmt.Sprintf("%v-%v-%v-%v-%v-%v-%v", params.TableName, string(qParam), string(sParam), string(pParam), string(dIds), crudInstance.Skip, crudInstance.Limit)

	// Audit/TransLog instance
	crudInstance.TransLog = mcauditlog.NewAuditLogPgx(crudInstance.AuditDb, crudInstance.AuditTable)

	return crudInstance
}

// String() function implementation for crud instance/object
func (crud Crud) String() string {
	return fmt.Sprintf("CRUD Instance Information: %#v \n\n", crud)
}

// Methods

// SaveRecord function creates new record(s) or updates existing record(s)
func (crud *Crud) SaveRecord() mcresponse.ResponseMessage {
	// transform actionParams ([]map[string]interface) camelCase fields to underscore
	actParams, err := ArrayMapToMapUnderscore(crud.ActionParams)
	if err != nil {
		return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("actParams-records format error: %v", err.Error()),
			Value:   nil,
		})
	}

	//  compute taskType-records from actionParams: create or update
	var (
		createRecs = ActionParamsType{} // records without id field-value
		updateRecs = ActionParamsType{} // records with id field-value
		recIds     []string             // capture recordIds for separate/multiple updates
	)
	for _, rec := range actParams {
		// determine if record exists (update), cast id into string or is new (create)
		recId, ok := rec["id"]
		recIdStr, idOk := recId.(string)
		if ok && recId != nil && idOk && recIdStr != "" {
			rec["updated_by"] = crud.UserInfo.UserId
			rec["updated_at"] = time.Now()
			recIds = append(recIds, recIdStr)
			updateRecs = append(updateRecs, rec)
		} else {
			rec["created_by"] = crud.UserInfo.UserId
			rec["created_at"] = time.Now()
			createRecs = append(createRecs, rec)
		}
	}

	// set action-type (create or update)
	if len(createRecs) > 0 && len(updateRecs) > 0 {
		// return only create or update task permitted
		return mcresponse.GetResMessage("saveError", mcresponse.ResponseMessageOptions{
			Message: "You may only create or update record(s), not both at the same time",
			Value:   nil,
		})
	}
	if len(updateRecs) > 1 && len(recIds) > 1 && len(recIds) == len(updateRecs) {
		crud.TaskType = UpdateTask
		crud.RecordIds = recIds
	} else if len(updateRecs) == 1 && (len(crud.RecordIds) > 0 || crud.QueryParams != nil) {
		crud.TaskType = UpdateTask
	} else if len(recIds) == 0 && len(createRecs) > 0 {
		crud.TaskType = CreateTask
	} else {
		// return error, if above conditions could not be met
		return mcresponse.GetResMessage("saveError", mcresponse.ResponseMessageOptions{
			Message: "Incomplete params to perform the data-operation-task",
			Value:   nil,
		})
	}

	// create/insert new record(s)
	if crud.TaskType == CrudTasks().Create {
		// check task-permission
		if crud.CheckAccess {
			accessRes := crud.TaskPermission(crud.TaskType)
			if accessRes.Code != "success" {
				return accessRes
			}
		}
		// save-record(s): create/insert new record(s): len(recordIds) = 0 && len(createRecs) > 0
		return crud.Create(createRecs)
	}

	if crud.TaskType == CrudTasks().Update {
		// check task-permission
		if crud.CheckAccess {
			accessRes := crud.TaskPermission(crud.TaskType)
			if accessRes.Code != "success" {
				return accessRes
			}
		}
		// update 1 or more records by ids or queryParams
		if len(crud.ActionParams) == 1 || len(updateRecs) == 1 {
			upRec := updateRecs[0]
			// update the record by recordId
			if len(crud.RecordIds) == 1 {
				return crud.UpdateById(upRec, crud.RecordIds[0])
			}
			// update record(s) by recordIds
			if len(crud.RecordIds) > 1 {
				return crud.UpdateByIds(upRec)
			}
			// update record(s) by queryParams
			if len(crud.QueryParams) > 0 {
				return crud.UpdateByParam(upRec)
			}
		}
		// update multiple records
		if len(crud.ActionParams) > 1 {
			return crud.Update(updateRecs)
		}
	}

	// otherwise, return saveError
	return mcresponse.GetResMessage("saveError", mcresponse.ResponseMessageOptions{
		Message: "Save error: incomplete or invalid action/query-params provided",
		Value:   nil,
	})

}

// DeleteRecord function deletes/removes record(s) by id(s) or params
func (crud *Crud) DeleteRecord() mcresponse.ResponseMessage {
	// check task-permission - delete
	if crud.CheckAccess {
		accessRes := crud.TaskPermission(DeleteTask)
		if accessRes.Code != "success" {
			return accessRes
		}
	}

	if len(crud.RecordIds) == 1 {
		return crud.DeleteById(crud.RecordIds[0])
	}
	if len(crud.RecordIds) > 1 {
		return crud.DeleteByIds()
	}
	if crud.QueryParams != nil && len(crud.QueryParams) > 0 {
		return crud.DeleteByParam()
	}
	// delete-all ***RESTRICTED***
	// otherwise return error
	return mcresponse.GetResMessage("removeError", mcresponse.ResponseMessageOptions{
		Message: "Remove error: incomplete or invalid query-conditions provided",
		Value:   nil,
	})
}

// GetRecord function get records by id, params or all
func (crud *Crud) GetRecord(modelRef interface{}) mcresponse.ResponseMessage {
	// check task-permission - get/read
	if crud.CheckAccess {
		accessRes := crud.TaskPermission(ReadTask)
		if accessRes.Code != "success" {
			return accessRes
		}
	}

	if len(crud.RecordIds) == 1 {
		return crud.GetById(modelRef, crud.RecordIds[0])
	}
	if len(crud.RecordIds) > 1 {
		return crud.GetByIds(modelRef)
	}
	if crud.QueryParams != nil && len(crud.QueryParams) > 0 {
		return crud.GetByParam(modelRef)
	}
	return crud.GetAll(modelRef)
}

// GetRecords function get records by id, params or all - lookup-items
func (crud *Crud) GetRecords(modelRef interface{}) mcresponse.ResponseMessage {
	if len(crud.RecordIds) == 1 {
		return crud.GetById(modelRef, crud.RecordIds[0])
	}
	if len(crud.RecordIds) > 1 {
		return crud.GetByIds(modelRef)
	}
	if crud.QueryParams != nil && len(crud.QueryParams) > 0 {
		return crud.GetByParam(modelRef)
	}
	return crud.GetAll(modelRef)
}
