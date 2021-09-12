// @Author: abbeymart | Abi Akindele | @Created: 2020-12-01 | @Updated: 2020-12-01
// @Company: mConnect.biz | @License: MIT
// @Description: Base type/method CRUD operations for PgDB

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
	crudInstance.ModelRef = params.ModelRef
	crudInstance.ModelFieldsRef = params.ModelFieldsRef
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
	crudInstance.BulkCreate = options.BulkCreate

	// Default values
	if crudInstance.AuditTable == "" {
		crudInstance.AuditTable = "audits"
	}
	if crudInstance.AccessTable == "" {
		crudInstance.AccessTable = "accesses"
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
		crudInstance.MaxQueryLimit = 100000
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

// String() method implementation for crud instance/object
func (crud Crud) String() string {
	return fmt.Sprintf("CRUD Instance Information: %#v \n\n", crud)
}

// Methods

// SaveRecord method creates new record(s) or updates existing record(s)
func (crud *Crud) SaveRecord() mcresponse.ResponseMessage {
	//  compute taskType-records from actionParams: create or update
	var (
		createRecs = ActionParamsType{} // records without id field-value
		updateRecs = ActionParamsType{} // records with id field-value
		recIds     []string             // capture recordIds for separate/multiple updates
	)
	for _, rec := range crud.ActionParams {
		// determine if record exists (update), cast id into string or new (create)
		recId, ok := rec["id"]
		recIdStr, idOk := recId.(string)
		if ok && recId != nil && idOk && recIdStr != "" {
			rec["updatedBy"] = crud.UserInfo.UserId
			rec["updatedAt"] = time.Now()
			recIds = append(recIds, recIdStr)
			updateRecs = append(updateRecs, rec)
		} else if len(crud.RecordIds) > 0 || len(crud.QueryParams) > 0 {
			rec["updatedBy"] = crud.UserInfo.UserId
			rec["updatedAt"] = time.Now()
			updateRecs = append(updateRecs, rec)
		} else {
			rec["createdBy"] = crud.UserInfo.UserId
			rec["createdAt"] = time.Now()
			createRecs = append(createRecs, rec)
		}
	}
	// validate and set task-type, create or update
	if len(createRecs) > 0 && len(updateRecs) > 0 {
		return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
			Message: "You may only create or update record(s), not both at the same time",
			Value:   nil,
		})
	}
	// set task-type
	if len(createRecs) > 0 {
		crud.TaskType = CreateTask
	} else if len(updateRecs) > 0 {
		crud.TaskType = UpdateTask
	} else {
		return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
			Message: "Inputs errors: actionParams required to complete create or update task.",
			Value:   nil,
		})
	}

	// create/insert new record(s)
	if crud.TaskType == CreateTask && len(createRecs) > 0 {
		// check task-permission
		if crud.CheckAccess {
			accessRes := crud.CheckTaskAccess()
			if accessRes.Code != "success" {
				return accessRes
			}
		}
		if crud.BulkCreate {
			return crud.CreateCopy(createRecs)
		}
		return crud.Create(createRecs)
	}

	// update existing record(s), by record-id(s) or queryParams | or perform multiple updates
	if crud.TaskType == UpdateTask {
		if len(updateRecs) == 1 {
			if len(crud.RecordIds) == 1 {
				// check task-permission
				if crud.CheckAccess {
					accessRes := crud.TaskPermissionById(crud.TaskType)
					if accessRes.Code != "success" {
						return accessRes
					}
				}
				return crud.UpdateById(updateRecs[0], crud.RecordIds[0])
			}
			if len(crud.RecordIds) > 1 {
				// check task-permission
				if crud.CheckAccess {
					accessRes := crud.TaskPermissionById(crud.TaskType)
					if accessRes.Code != "success" {
						return accessRes
					}
				}
				return crud.UpdateByIds(updateRecs[0])
			}
			if len(crud.QueryParams) > 0 {
				// check task-permission
				if crud.CheckAccess {
					accessRes := crud.TaskPermissionByParam(crud.TaskType)
					if accessRes.Code != "success" {
						return accessRes
					}
				}
				return crud.UpdateByParam(updateRecs[0])
			}
		}
		// update multiple records
		crud.RecordIds = recIds
		// check task-permission
		if crud.CheckAccess {
			accessRes := crud.TaskPermissionById(crud.TaskType)
			if accessRes.Code != "success" {
				return accessRes
			}
		}
		return crud.Update(updateRecs)
	}
	// otherwise, return saveError
	return mcresponse.GetResMessage("saveError", mcresponse.ResponseMessageOptions{
		Message: "Save error: incomplete or invalid parameters (action/query-params/record-ids) provided",
		Value:   nil,
	})
}

// DeleteRecord method deletes/removes record(s) by recordIds or queryParams
func (crud *Crud) DeleteRecord() mcresponse.ResponseMessage {
	if len(crud.RecordIds) == 1 {
		if crud.CheckAccess {
			accessRes := crud.TaskPermissionById(DeleteTask)
			if accessRes.Code != "success" {
				return accessRes
			}
		}
		return crud.DeleteById(crud.RecordIds[0])
	}
	if len(crud.RecordIds) > 1 {
		if crud.CheckAccess {
			accessRes := crud.TaskPermissionById(DeleteTask)
			if accessRes.Code != "success" {
				return accessRes
			}
		}
		return crud.DeleteByIds()
	}
	if crud.QueryParams != nil && len(crud.QueryParams) > 0 {
		if crud.CheckAccess {
			accessRes := crud.TaskPermissionByParam(DeleteTask)
			if accessRes.Code != "success" {
				return accessRes
			}
		}
		return crud.DeleteByParam()
	}
	// delete-all ***RESTRICTED***
	// otherwise return error
	return mcresponse.GetResMessage("removeError", mcresponse.ResponseMessageOptions{
		Message: "You may delete records by recordIds or queryParams only.",
		Value:   nil,
	})
}

// GetRecord method fetches records by recordIds, queryParams or all
func (crud *Crud) GetRecord() mcresponse.ResponseMessage {
	if len(crud.RecordIds) == 1 {
		if crud.CheckAccess {
			accessRes := crud.TaskPermissionById(ReadTask)
			if accessRes.Code != "success" {
				return accessRes
			}
		}
		return crud.GetById(crud.RecordIds[0])
	}
	if len(crud.RecordIds) > 1 {
		if crud.CheckAccess {
			accessRes := crud.TaskPermissionById(ReadTask)
			if accessRes.Code != "success" {
				return accessRes
			}
		}
		return crud.GetByIds()
	}
	if crud.QueryParams != nil && len(crud.QueryParams) > 0 {
		if crud.CheckAccess {
			accessRes := crud.TaskPermissionByParam(ReadTask)
			if accessRes.Code != "success" {
				return accessRes
			}
		}
		return crud.GetByParam()
	}
	if crud.CheckAccess {
		accessRes := crud.CheckTaskAccess()
		if accessRes.Code != "success" {
			return accessRes
		}
	}
	return crud.GetAll()
}

// GetRecords method fetches records by recordIds, queryParams or all - lookup-items (no-access-constraint)
func (crud *Crud) GetRecords() mcresponse.ResponseMessage {
	if len(crud.RecordIds) == 1 {
		return crud.GetById(crud.RecordIds[0])
	}
	if len(crud.RecordIds) > 1 {
		return crud.GetByIds()
	}
	if crud.QueryParams != nil && len(crud.QueryParams) > 0 {
		return crud.GetByParam()
	}
	return crud.GetAll()
}
