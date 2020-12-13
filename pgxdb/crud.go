// @Author: abbeymart | Abi Akindele | @Created: 2020-12-01 | @Updated: 2020-12-01
// @Company: mConnect.biz | @License: MIT
// @Description: Base type/function CRUD operations for PgDB

package pgxdb

import (
	"encoding/json"
	"fmt"
	"github.com/abbeymart/mcauditlog"
)

type Crud struct {
	PgxCrudTaskType
	PgxCrudOptionsType
	TransLog mcauditlog.PgxLogParam
	HashKey  string // Unique for exactly the same query
}

// constructor
func NewCrud(params PgxCrudTaskType, options PgxCrudOptionsType) Crud {
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

	// crud options
	result.Skip = options.Skip
	result.Limit = options.Limit
	result.MaxQueryLimit = options.MaxQueryLimit

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
	result.CheckAccess = options.CheckAccess // Dec 09/2020: user to implement auth as a middleware
	// Compute HashKey from TableName, QueryParams, SortParams, ProjectParams and RecordIds
	qParam, _ := json.Marshal(params.QueryParams)
	sParam, _ := json.Marshal(params.SortParams)
	pParam, _ := json.Marshal(params.ProjectParams)
	dIds, _ := json.Marshal(params.DocIds)
	result.HashKey = params.TableName + string(qParam) + string(sParam) + string(pParam) + string(dIds)

	// Audit/TransLog instance
	result.TransLog = mcauditlog.NewAuditLogPgx(result.AuditDb, result.AuditTable)

	// Default values
	if result.AuditTable == "" {
		result.AuditTable = "audits"
	}
	if result.AccessTable == "" {
		result.AccessTable = "accesskeys"
	}
	if result.RoleTable == "" {
		result.RoleTable = "roles"
	}
	if result.UserTable == "" {
		result.UserTable = "users"
	}
	if result.ServiceTable == "" {
		result.AuditTable = "services"
	}
	if result.AuditDb == nil {
		result.AuditDb = result.AppDb
	}
	if result.AccessDb == nil {
		result.AccessDb = result.AppDb
	}
	if result.Skip < 0 {
		result.Skip = 0
	}
	if result.Limit > result.MaxQueryLimit && result.MaxQueryLimit != 0 {
		result.Limit = result.MaxQueryLimit
	} else if result.Limit > 10000 {
		result.Limit = 10000
	}

	return result
}

// String() function implementation
func (crud Crud) String() string {
	//appDb := fmt.Sprintf("Application DB: %v", crud.AppDb)
	return fmt.Sprintf(`
	Application DB: %v \n Table Name: %v \n
	`,
		crud.AppDb,
		crud.TableName)
}

// Methods => separate go-files: auth.go...
