// @Author: abbeymart | Abi Akindele | @Created: 2020-12-28 | @Updated: 2020-12-28
// @Company: mConnect.biz | @License: MIT
// @Description: test-cases data

package mccrud

import (
	"encoding/json"
	"time"
)

// Models

type Audit struct {
	Id            string      `json:"id" db:"id"`
	TableName     string      `json:"tableName" db:"table_name"`
	LogRecords    interface{} `json:"logRecords" db:"log_records"`
	NewLogRecords interface{} `json:"newLogRecords" db:"new_log_records"`
	LogType       string      `json:"logType" db:"log_type"`
	LogBy         string      `json:"logBy" db:"log_by"`
	LogAt         time.Time   `json:"logAt" db:"log_at"`
}

type Group struct {
	BaseModelType
	Name string `json:"name" gorm:"unique" db:"name"`
}

type Category struct {
	BaseModelType
	Name      string  `json:"name"  db:"name"`
	Path      string  `json:"path" db:"path"`
	Priority  uint    `json:"priority" db:"priority"`
	ParentId  *string `json:"parentId" db:"parent_id"`
	GroupId   string  `json:"groupId" db:"group_id"`
	GroupName string  `json:"groupName" db:"group_name"`
	IconStyle string  `json:"iconStyle" db:"icon_style"`
}

const GroupTable = "groups"
const CategoryTable = "categories"
const AuditTable = "audits"
const GetTable = "audits_get"
const DeleteTable = "audits_delete"
const DeleteAllTable = "audits_delete_all"
const UpdateTable = "audits_update"

const UserId = "085f48c5-8763-4e22-a1c6-ac1a68ba07de" // TODO: review/update
const UId = "faea411c-e82d-454f-8ee7-574c4e753d06"

var UserInfo = map[string]interface{}{
	"userId":    "085f48c5-8763-4e22-a1c6-ac1a68ba07de",
	"loginName": "abbeymart",
	"email":     "abbeya1@yahoo.com",
	"language":  "en-US",
	"firstname": "Abi",
	"lastname":  "Akindele",
	"token":     "",
	"expire":    0,
	"role":      "win-20_000_000",
}

var TestUserInfo = UserInfoType{
	UserId:    "085f48c5-8763-4e22-a1c6-ac1a68ba07de",
	LoginName: "abbeymart",
	Email:     "abbeya1@yahoo.com",
	Language:  "en-US",
	Firstname: "Abi",
	Lastname:  "Akindele",
	Token:     "",
	Expire:    0,
	Role:      "TBD: win-20_000_000",
}

var CrudParamOptions = CrudOptionsType{
	CheckAccess:   false,
	AuditTable:    "audits",
	UserTable:     "users",
	ProfileTable:  "profiles",
	ServiceTable:  "services",
	AccessTable:   "accesses",
	VerifyTable:   "verify_users",
	RoleTable:     "roles",
	LogCrud:       false,
	LogCreate:     false,
	LogUpdate:     false,
	LogDelete:     false,
	LogRead:       false,
	LogLogin:      false,
	LogLogout:     false,
	MaxQueryLimit: 100000,
	MsgFrom:       "support@mconnect.biz",
}

// TODO: create/update, get & delete records for groups & categories tables

var LogRecords = ActionParamType{
	"name":     "Abi",
	"desc":     "Testing only",
	"url":      "localhost:9000",
	"priority": 100,
	"cost":     1000.00,
}

var NewLogRecords = ActionParamType{
	"name":     "Abi Akindele",
	"desc":     "Testing only - updated",
	"url":      "localhost:9900",
	"priority": 1,
	"cost":     2000.00,
}

var LogRecords2 = ActionParamType{
	"name":     "Ola",
	"desc":     "Testing only - 2",
	"url":      "localhost:9000",
	"priority": 1,
	"cost":     10000.00,
}
var NewLogRecords2 = ActionParamType{
	"name":     "Ola",
	"desc":     "Testing only - 2 - updated",
	"url":      "localhost:9000",
	"priority": 1,
	"cost":     20000.00,
}

// create record(s)

var GroupCreateRec1 = ActionParamType{
	"name": "Location",
	"desc": "Location group",
}
var GroupCreateRec2 = ActionParamType{
	"name": "Address",
	"desc": "Address group",
}
var GroupUpdateRec1 = ActionParamType{
	"id":   "tbd",
	"name": "Location",
	"desc": "location group - updated",
}
var GroupUpdateRec2 = ActionParamType{
	"id":   "tbd",
	"name": "Address",
	"desc": "address group - updated",
}

var LogRecs, _ = json.Marshal(LogRecords)
var NewLogRecs, _ = json.Marshal(NewLogRecords)
var LogRecs2, _ = json.Marshal(LogRecords2)
var NewLogRecs2, _ = json.Marshal(NewLogRecords2)

var AuditCreateRec1 = ActionParamType{
	"tableName":  "audits",
	"logAt":      time.Now(),
	"logBy":      UserId,
	"logRecords": string(LogRecs),
	"logType":    CreateTask,
}
var AuditCreateRec2 = ActionParamType{
	"tableName":  "audits",
	"logAt":      time.Now(),
	"logBy":      UserId,
	"logRecords": string(LogRecs2),
	"logType":    CreateTask,
}
var AuditUpdateRec1 = ActionParamType{
	"id":            "f517ef7b-5457-4f51-a905-e427465defd0",
	"tableName":     "todos",
	"logAt":         time.Now(),
	"logBy":         UserId,
	"logRecords":    string(LogRecs),
	"newLogRecords": string(NewLogRecs),
	"logType":       UpdateTask,
}
var AuditUpdateRec2 = ActionParamType{
	"id":            "a66a3057-028d-4f64-aa18-05ea26b1d2dc",
	"tableName":     "todos",
	"logAt":         time.Now(),
	"logBy":         UserId,
	"logRecords":    string(LogRecs2),
	"newLogRecords": string(NewLogRecs2),
	"logType":       UpdateTask,
}
var GroupCreateActionParams = ActionParamsType{
	GroupCreateRec1,
	GroupCreateRec2,
}
var GroupUpdateActionParams = ActionParamsType{
	GroupUpdateRec1,
	GroupUpdateRec2,
}
var AuditCreateActionParams = ActionParamsType{
	AuditCreateRec1,
	AuditCreateRec2,
}
var AuditUpdateActionParams = ActionParamsType{
	AuditUpdateRec1,
	AuditUpdateRec2,
}

// TODO: update and delete params, by ids / queryParams

var GroupUpdateRecordById = ActionParamType{
	"name": "Location",
	"desc": "updated-by-id",
}
var GroupUpdateRecordByParam = ActionParamType{
	"name": "Address",
	"desc": "updated-by-param",
}

var AuditUpdateRecordById = ActionParamType{
	"id":            "03012156-19a4-43f9-b8ee-c1e9dd5d19b8",
	"tableName":     "groups",
	"logAt":         time.Now(),
	"logBy":         UserId,
	"logRecords":    string(LogRecs),
	"newLogRecords": string(NewLogRecs),
	"logType":       DeleteTask,
}
var AuditUpdateRecordByIds = ActionParamType{
	"id":            "be6b3f55-1724-4811-a23c-8ea481558f25",
	"tableName":     "users",
	"logAt":         time.Now(),
	"logBy":         UserId,
	"logRecords":    string(LogRecs),
	"newLogRecords": string(NewLogRecs),
	"logType":       CreateTask,
}
var AuditUpdateRecordByParam = ActionParamType{
	"id":            "03012156-19a4-43f9-b8ee-c1e9dd5d19b8",
	"tableName":     "contacts",
	"logAt":         time.Now(),
	"logBy":         UserId,
	"logRecords":    string(LogRecs),
	"newLogRecords": string(NewLogRecs),
	"logType":       UpdateTask,
}

// GetIds: for get-records by ids & params | TODO: update ids after create

var UpdateGroupById = "6900d9f9-2ceb-450f-9a9e-527eb66c962f"
var UpdateGroupByIds = []string{"6900d9f9-2ceb-450f-9a9e-527eb66c962f", "122d0f0e-3111-41a5-9103-24fa81004550"} // TODO: updated ids, after create-task
var UpdateGroupByParams = QueryParamType{
	"name": "Location",
}

var GetAuditById = "f517ef7b-5457-4f51-a905-e427465defd0"
var GetAuditByIds = []string{"f517ef7b-5457-4f51-a905-e427465defd0", "a66a3057-028d-4f64-aa18-05ea26b1d2dc"}
var GetAuditByParams = QueryParamType{
	"logType": "create",
}
var DeleteAuditById = "a6dbc263-86ee-436a-bb41-0d3b77274f79"
var DeleteAuditByIds = []string{
	"f517ef7b-5457-4f51-a905-e427465defd0",
	"a66a3057-028d-4f64-aa18-05ea26b1d2dc",
	"03012156-19a4-43f9-b8ee-c1e9dd5d19b8",
	"be6b3f55-1724-4811-a23c-8ea481558f25",
}
var DeleteAuditByParams = QueryParamType{
	"logType": "read",
}
var UpdateAuditById = "8883ef51-d730-4645-9bc8-f986cebb7881"
var UpdateAuditByIds = []string{
	"b8ca59b8-46f4-4a0c-82aa-f689d8ce295c",
	"a6dbc263-86ee-436a-bb41-0d3b77274f79",
	"19604025-cb69-46cd-883a-8f51e754817b",
}
var UpdateAuditByParams = QueryParamType{
	"logType": "read",
}
