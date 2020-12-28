// @Author: abbeymart | Abi Akindele | @Created: 2020-12-28 | @Updated: 2020-12-28
// @Company: mConnect.biz | @License: MIT
// @Description: go: mConnect

package tests

import (
	"encoding/json"
	"github.com/abbeymart/mcauditlog"
	"github.com/abbeymart/mctypes"
	"time"
)

const TestTable = "audits_test1"
const DeleteAllTable = "audits_test2"
const TestAuditTable = "audits"

var CreateTableFields = []string{
	"table_name",
	"log_records",
	"new_log_records",
	"log_type",
	"log_by",
	"log_at",
}

var UpdateTableFields = []string{
	"id",
	"table_name",
	"log_records",
	"new_log_records",
	"log_type",
	"log_by",
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
var TableRecords, _ = json.Marshal(Recs)

//fmt.Println("table-records-json", string(tableRecords))
var NewRecs = TestParam{Name: "Abi Akindele", Desc: "Testing only - updated", Url: "localhost:9900", Priority: 1, Cost: 2000.00}
var NewTableRecords, _ = json.Marshal(NewRecs)

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
	LogRecords: string(TableRecords),
	LogBy:      UserId,
	LogType:    mcauditlog.CreateLog,
	LogAt:      time.Now(),
}

var CreateRecordB = mcauditlog.AuditRecord{
	TableName:  "services",
	LogRecords: string(TableRecords),
	LogBy:      UserId,
	LogType:    mcauditlog.CreateLog,
	LogAt:      time.Now(),
}

var valParam1, _ = DataToValueParam(CreateRecordA)
var valParam2, _ = DataToValueParam(CreateRecordB)

var CreateActionParams = mctypes.ActionParamsType{
	valParam1,
	valParam2,
}

// update record(s)
var UpdateRecordA = mcauditlog.AuditRecord{
	TableName:  "services",
	LogRecords: string(TableRecords),
	LogBy:      UserId,
	LogType:    mcauditlog.CreateLog,
	LogAt:      time.Now(),
}

var UpdateRecordB = mcauditlog.AuditRecord{
	TableName:  "services",
	LogRecords: string(TableRecords),
	LogBy:      UserId,
	LogType:    mcauditlog.CreateLog,
	LogAt:      time.Now(),
}

var UpdateRecordById = mcauditlog.AuditRecord{
	TableName:  "services",
	LogRecords: string(TableRecords),
	LogBy:      UserId,
	LogType:    mcauditlog.CreateLog,
	LogAt:      time.Now(),
}

var UpdateRecordByParam = mcauditlog.AuditRecord{
	TableName:  "services",
	LogRecords: string(TableRecords),
	LogBy:      UserId,
	LogType:    mcauditlog.CreateLog,
	LogAt:      time.Now(),
}

var updateRec1, _ = DataToValueParam(UpdateRecordA)
var updateRec2, _ = DataToValueParam(UpdateRecordB)

var TestUpdateRecords = mctypes.ActionParamsType{
	updateRec1,
	updateRec2,
}

var TestUpdateRecordIds = []string{
	"abc",
	"xyz",
}

// get record(s)
// by id
var TestGetRecordIds = []string{
	"abc",
	"xyz",
}

// by param

// delete record(s)
// by id
var TestDeleteRecordIds = []string{
	"abc",
	"xyz",
}

// by param
