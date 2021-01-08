// @Author: abbeymart | Abi Akindele | @Created: 2020-12-28 | @Updated: 2020-12-28
// @Company: mConnect.biz | @License: MIT
// @Description: go: mConnect

package mccrud

import (
	"github.com/abbeymart/mcauditlog"
	"github.com/abbeymart/mccrud/helper"
	"github.com/abbeymart/mctypes"
	"time"
)

const TestTable = "audits_test1"
const DeleteAllTable = "audits_test2"
const TestAuditTable = "audits"

var CreateTableFields = []string{
	"table_name",
	"log_records",
	"log_type",
	"log_by",
	"log_at",
}

var UpdateTableFields = []string{
	"new_log_records",
	"log_type",
	"log_at",
}

var GetTableFields = []string{
	"id",
	"table_name",
	"log_records",
	"new_log_records",
	"log_type",
	"log_by",
	"log_at",
}

type TestParam struct {
	Name     string  `json:"name"`
	Desc     string  `json:"desc"`
	Url      string  `json:"url"`
	Priority int     `json:"priority"`
	Cost     float64 `json:"cost"`
}

const UserId = "085f48c5-8763-4e22-a1c6-ac1a68ba07de"
var TestUserInfo = mctypes.UserInfoType{
	UserId:    "085f48c5-8763-4e22-a1c6-ac1a68ba07de",
	LoginName: "abbeymart",
	Email:     "abbeya1@yahoo.com",
	Language:  "en-US",
	FirstName: "Abi",
	LastName:  "Akindele",
	Token:     "",
	Expire:    0,
	Group:     "TBD",
}

var Recs = TestParam{Name: "Abi", Desc: "Testing only", Url: "localhost:9000", Priority: 1, Cost: 1000.00}
var TableRecords, _ = helper.DataToValueParam(Recs)

//fmt.Println("table-records-json", string(tableRecords))
var NewRecs = TestParam{Name: "Abi Akindele", Desc: "Testing only - updated", Url: "localhost:9900", Priority: 1, Cost: 2000.00}
var NewTableRecords, _ = helper.DataToValueParam(NewRecs)

//fmt.Println("new-table-records-json", string(newTableRecords))
//var ReadP = map[string][]string{"keywords": {"lagos", "nigeria", "ghana", "accra"}}
//var ReadParams, _ = json.Marshal(ReadP)

var TestCrudParamOptions = mctypes.CrudOptionsType{
	AuditTable:    "audits",
	UserTable:     "users",
	ServiceTable:  "services",
	AccessTable:   "access_keys",
	VerifyTable:   "verify_users",
	RoleTable:     "roles",
	LogCreate:     true,
	LogUpdate:     true,
	LogDelete:     true,
	LogRead:       true,
	LogLogin:      true,
	LogLogout:     true,
	MaxQueryLimit: 100000,
	MsgFrom:       "support@mconnect.biz",
}

// create record(s)
var CreateRecordA = mcauditlog.AuditRecord{
	TableName:  "services",
	LogRecords: TableRecords,
	LogBy:      UserId,
	LogType:    mcauditlog.CreateLog,
	LogAt:      time.Now(),
}
var CreateRecordB = mcauditlog.AuditRecord{
	TableName:  "services",
	LogRecords: TableRecords,
	LogBy:      UserId,
	LogType:    mcauditlog.CreateLog,
	LogAt:      time.Now(),
}
var valParam1, _ = helper.DataToValueParam(CreateRecordA)
var valParam2, _ = helper.DataToValueParam(CreateRecordB)
var CreateActionParams = mctypes.ActionParamsType{
	valParam1,
	valParam2,
}

// update record(s)
type UpdateRecordType struct {
	Id            string
	TableName     string
	LogRecords    interface{}
	NewLogRecords interface{}
	LogBy         string
	LogType       string
	LogAt         time.Time
}

var upRecs = TestParam{Name: "Abi100", Desc: "Testing only100", Url: "localhost:9000", Priority: 1, Cost: 1000.00}
var upTableRecords, _ = helper.DataToValueParam(upRecs)
var upRecs2 = TestParam{Name: "Abi200", Desc: "Testing only200", Url: "localhost:9000", Priority: 1, Cost: 1000.00}
var upTableRecords2, _ = helper.DataToValueParam(upRecs2)
var UpdateRecordA = UpdateRecordType{
	Id:            "record-id",
	TableName:     "services",
	LogRecords:    upTableRecords,
	NewLogRecords: NewTableRecords,
	LogBy:         UserId,
	LogType:       mcauditlog.UpdateLog,
	LogAt:         time.Now(),
}

var UpdateRecordB = UpdateRecordType{
	Id:            "record-id",
	TableName:     "services2",
	LogRecords:    upTableRecords2,
	NewLogRecords: NewTableRecords,
	LogBy:         UserId,
	LogType:       mcauditlog.UpdateLog,
	LogAt:         time.Now(),
}

var UpdateRecordById = UpdateRecordType{
	TableName:     "services2",
	LogRecords:    upTableRecords,
	NewLogRecords: NewTableRecords,
	LogBy:         UserId,
	LogType:       mcauditlog.UpdateLog,
	LogAt:         time.Now(),
}

var UpdateRecordByParam = mcauditlog.AuditRecord{
	TableName:     "services3",
	LogRecords:    upTableRecords2,
	NewLogRecords: NewTableRecords,
	LogBy:         UserId,
	LogType:       mcauditlog.UpdateLog,
	LogAt:         time.Now(),
}

var UpdateIds = []string{"id1", "id2"}
var UpdateParams = mctypes.QueryParamType{
	mctypes.QueryGroupType{

	},
	mctypes.QueryGroupType{

	},
}

var updateRec1, _ = helper.DataToValueParam(UpdateRecordA)
var updateRec2, _ = helper.DataToValueParam(UpdateRecordB)
var updateRecId, _ = helper.DataToValueParam(UpdateRecordById)
var updateRecParam, _ = helper.DataToValueParam(UpdateRecordByParam)

var TestUpdateRecords = mctypes.ActionParamsType{
	updateRec1,
	updateRec2,
}

var TestUpdateRecordIds = mctypes.ActionParamsType{
	updateRecId,
}
var TestUpdateRecordParam = mctypes.ActionParamsType{
	updateRecParam,
}

// get record(s)
type GetRecordType struct {
	Id            string
	TableName     string
	LogRecords    interface{}
	NewLogRecords interface{}
	LogBy         string
	LogType       string
	LogAt         time.Time
}

var (
	id            string
	tableName     string
	logRecords    interface{}
	newLogRecords interface{}
	logBy         string
	logType       string
	logAt         time.Time
)

// tablePointers for scanning
var tableFieldPointers = []interface{}{
	&id,
	&tableName,
	&logRecords,
	&newLogRecords,
	&logBy,
	&logType,
	&logAt,
}

var getResults []GetRecordType
// synchronize get records, via a buffered channel
var getChan = make(chan int, 1)
// TODO: check getChan to update getResults one at a time
var getResult = GetRecordType{
	Id:            id,
	TableName:     tableName,
	LogRecords:    logRecords,
	NewLogRecords: newLogRecords,
	LogBy:         logBy,
	LogType:       logType,
	LogAt:         logAt,
}

// get by ids & params
var GetIds = []string{"id1", "id2"}
var GetParams = mctypes.QueryParamType{
	mctypes.QueryGroupType{

	},
	mctypes.QueryGroupType{

	},
}

// delete record(s) by ids & params
var DeleteIds = []string{"id1", "id2"}
var DeleteParams = mctypes.QueryParamType{
	mctypes.QueryGroupType{

	},
	mctypes.QueryGroupType{

	},
}