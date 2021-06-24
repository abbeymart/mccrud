// @Author: abbeymart | Abi Akindele | @Created: 2020-12-05 | @Updated: 2020-12-05
// @Company: mConnect.biz | @License: MIT
// @Description: mongoDB CRUD base type / behaviours

package mongo

import (
	"encoding/json"
	"fmt"
	"github.com/abbeymart/mccrud/types"
)

type CrudMongo struct {
	types.MongoCrudTaskType
	types.MongoCrudOptionsType
	HashKey string // Unique for exactly the same query
}

// NewCrudMongo constructor
func NewCrudMongo(params types.MongoCrudTaskType, options types.MongoCrudOptionsType) CrudMongo {
	result := CrudMongo{}
	// compute crud params
	result.AppDb = params.AppDb
	result.TableName = params.TableName
	result.UserInfo = params.UserInfo
	result.ActionParams = params.ActionParams
	result.RecordIds = params.RecordIds
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
	result.LogCrud = options.LogCrud
	result.LogRead = options.LogRead
	result.LogCreate = options.LogCreate
	result.LogUpdate = options.LogUpdate
	result.LogDelete = options.LogDelete
	result.CheckAccess = options.CheckAccess
	// Compute HashKey from TableName, QueryParams, SortParams, ProjectParams and RecordIds
	qParam, _ := json.Marshal(params.QueryParams)
	sParam, _ := json.Marshal(params.SortParams)
	pParam, _ := json.Marshal(params.ProjectParams)
	dIds, _ := json.Marshal(params.RecordIds)
	result.HashKey = params.TableName + string(qParam) + string(sParam) + string(pParam) + string(dIds)

	// TODO: TransLog instance
	//result.TransLog = newLog(result.AuditDb, result.AuditTable)

	return result
}

// methods => separate go-files

// String() function implementation
func (crud CrudMongo) String() string {
	//appDb := fmt.Sprintf("Application DB: %v", crud.AppDb)
	return fmt.Sprintf(`
	Application DB: %v \n Table Name: %v \n
	`,
		crud.AppDb,
		crud.TableName)
}
