// @Author: abbeymart | Abi Akindele | @Created: 2020-12-01 | @Updated: 2020-12-01
// @Company: mConnect.biz | @License: MIT
// @Description: Base type/function CRUD operations for PgDB

package mccrud

import (
	"encoding/json"
	"fmt"
	"github.com/abbeymart/mcauditlog"
)

type Crud struct {
	CrudTaskType
	CrudOptionsType
	TransLog mcauditlog.LogParam
	HashKey string // Unique for exactly the same query
}

// constructor
func NewCrud(params CrudTaskType, options CrudOptionsType) Crud {
	result := Crud{}
	// compute crud params
	result.AppDb = params.AppDb
	result.TableName = params.TableName
	result.UserInfo = params.UserInfo
	result.ActionParams = params.ActionParams
	result.DocIds = params.DocIds
	result.QueryParams = params.QueryParams
	result.SortParams = params.SortParams
	result.ProjectParams = params.ProjectParams
	result.ExistParams = params.ExistParams
	result.Token = params.Token
	result.TaskName = params.TaskName

	// Options
	result.Skip = options.Skip
	result.Limit = options.Limit
	//result.DefaultLimit = defaultLimit

	result.AuditTable = options.AuditTable
	result.AccessTable = options.AccessTable
	result.RoleTable = options.RoleTable
	result.UserTable = options.UserTable
	result.AuditDb = options.AuditDb
	result.AccessDb = options.AccessDb
	result.LogAll = options.LogAll
	result.LogRead = options.LogRead
	result.LogCreate = options.LogCreate
	result.LogUpdate = options.LogUpdate
	result.LogDelete = options.LogDelete
	result.CheckAccess = options.CheckAccess
	// Compute HashKey from TableName, QueryParams, SortParams, ProjectParams and DocIds
	qParam, _ := json.Marshal(params.QueryParams)
	sParam, _ := json.Marshal(params.SortParams)
	pParam, _ := json.Marshal(params.ProjectParams)
	dIds, _ := json.Marshal(params.DocIds)
	result.HashKey = params.TableName + string(qParam) + string(sParam) + string(pParam) + string(dIds)

	// TODO: TransLog instance
	result.TransLog = mcauditlog.NewAuditLog(result.AuditDb, result.AuditTable)

	return result
}

// methods => separate go-files

// String() function implementation
func (crud Crud) String() string {
	//appDb := fmt.Sprintf("Application DB: %v", crud.AppDb)
	return fmt.Sprintf(`
	Application DB: %v \n Table Name: %v \n
	`,
		crud.AppDb,
		crud.TableName)
}
