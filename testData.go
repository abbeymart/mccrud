// @Author: abbeymart | Abi Akindele | @Created: 2020-12-28 | @Updated: 2020-12-28
// @Company: mConnect.biz | @License: MIT
// @Description: test-cases data

package mccrud

import (
	"time"
)

// Models

type Audit struct {
	Id            string      `json:"id" mcorm:"id"`
	TableName     string      `json:"tableName" mcorm:"table_name"`
	LogRecords    interface{} `json:"logRecords" mcorm:"log_records"`
	NewLogRecords interface{} `json:"newLogRecords" mcorm:"new_log_records"`
	LogType       string      `json:"logType" mcorm:"log_type"`
	LogBy         string      `json:"logBy" mcorm:"log_by"`
	LogAt         time.Time   `json:"logAt" mcorm:"log_at"`
}

type Group struct {
	BaseModelType
	Name string `json:"name" gorm:"unique" mcorm:"name"`
}

type Category struct {
	BaseModelType
	Name      string  `json:"name"  mcorm:"name"`
	Path      string  `json:"path" mcorm:"path"`
	Priority  uint    `json:"priority" mcorm:"priority"`
	ParentId  *string `json:"parentId" mcorm:"parent_id"`
	GroupId   string  `json:"groupId" mcorm:"group_id"`
	GroupName string  `json:"groupName" mcorm:"group_name"`
	IconStyle string  `json:"iconStyle" mcorm:"icon_style"`
}

const GroupTable = "groups"
const CategoryTable = "categories"
const AuditTable = "audits"
const DeleteTable = "audits_test1"
const DeleteAllTable = "audits_test2"

const UserId = "085f48c5-8763-4e22-a1c6-ac1a68ba07de" // TODO: review/update

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
	Role:     "TBD: win-20_000_000",
}


var CrudParamOptions = CrudOptionsType{
	AuditTable:    "audits",
	UserTable:     "users",
	ProfileTable:  "profiles",
	ServiceTable:  "services",
	AccessTable:   "access_keys",
	VerifyTable:   "verify_users",
	RoleTable:     "roles",
	LogCrud:       true,
	LogCreate:     true,
	LogUpdate:     true,
	LogDelete:     true,
	LogRead:       true,
	LogLogin:      true,
	LogLogout:     true,
	MaxQueryLimit: 100000,
	MsgFrom:       "support@mconnect.biz",
}

// TODO: create/update, get & delete records for groups & categories tables

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

var CategoryCreateRec1 = ActionParamType{
	"name":      "Region",
	"groupId":   "tbd",
	"groupName": "Location",
	"path":      "reg",
	"desc":      "regional location",
}

var CategoryCreateRec2 = ActionParamType{
	"name":      "business",
	"groupId":   "tbd",
	"groupName": "Address",
	"path":      "bus",
	"desc":      "business address",
}

var CategoryUpdateRec1 = ActionParamType{
	"id":        "tbd",
	"name":      "Region",
	"groupId":   "tbd",
	"groupName": "Location",
	"path":      "reg-up",
	"desc":      "regional location - update",
}

var CategoryUpdateRec2 = ActionParamType{
	"id":        "tbd",
	"name":      "business",
	"groupId":   "tbd",
	"groupName": "Address",
	"path":      "bus-up",
	"desc":      "business address - update",
}

var GroupCreateActionParams = ActionParamsType{
	GroupCreateRec1,
	GroupCreateRec2,
}

var CategoryCreateActionParams = ActionParamsType{
	CategoryCreateRec1,
	CategoryCreateRec2,
}

var GroupUpdateActionParams = ActionParamsType{
	GroupUpdateRec1,
	GroupUpdateRec2,
}

var CategoryUpdateActionParams = ActionParamsType{
	CategoryUpdateRec1,
	CategoryUpdateRec2,
}

// TODO: update and delete params, by ids / queryParams

var GroupUpdateRecordById = ActionParamType{
	"name": "location",
	"desc": "updated-by-id",
}

var CategoryUpdateRecordById = ActionParamType{
	"name": "country",
	"desc": "updated-by-id",
}

var GroupUpdateRecordByParam = ActionParamType{
	"name": "address",
	"desc": "updated-by-param",
}

var CategoryUpdateRecordByParam = ActionParamType{
	"name": "business",
	"desc": "updated-by-param",
}

// GetIds: for get-records by ids & params

var GetGroupByIds = []string{"6900d9f9-2ceb-450f-9a9e-527eb66c962f", "122d0f0e-3111-41a5-9103-24fa81004550"}
var GetGroupByParams = QueryParamType{
	"name": "Location",
}

var GetCategoryByIds = []string{"6900d9f9-2ceb-450f-9a9e-527eb66c962f", "122d0f0e-3111-41a5-9103-24fa81004550"}
var GetCategoryByParams = QueryParamType{
	"name": "Region",
	"groupName": "Location",
}

// DeleteIds delete record(s) by ids & queryParams - temporary tables

var DeleteByIds = []string{"dba4adbb-4482-4f3d-bb05-0db80c30876b", "02f83bc1-8fa3-432a-8432-709f0df3f3b0"}
var DeleteByParams = QueryParamType{
	"logType": "create",
}
