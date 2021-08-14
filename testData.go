// @Author: abbeymart | Abi Akindele | @Created: 2020-12-28 | @Updated: 2020-12-28
// @Company: mConnect.biz | @License: MIT
// @Description: go: mConnect

package mccrud

import (
	"time"
)

// Models

type Group struct {
	BaseModelType
	Name string `json:"name" gorm:"unique" mcorm:"name"`
}

type Category struct {
	BaseModelType
	Name      string    `json:"name"  mcorm:"name"`
	OwnerId   string    `json:"ownerId" mcorm:"owner_id"`
	Path      string    `json:"path" mcorm:"path"`
	Priority  uint      `json:"priority" mcorm:"priority"`
	ParentId  *string   `json:"parentId" mcorm:"parent_id"`
	GroupId   string    `json:"groupId" mcorm:"group_id"`
	Group     Group     `json:"group" mcorm:"group"`
	Parent    *Category `json:"parent" mcorm:"parent"`
	IconStyle string    `json:"iconStyle" mcorm:"icon_style"`
}

const GroupTable = "groups"
const CategoryTable = "categories"
const AuditTable = "audits"
const DeleteAllTable = "audits_test2"
const TestAuditTable = "audits"

const UserId = "085f48c5-8763-4e22-a1c6-ac1a68ba07de"

var TestUserInfo = UserInfoType{
	UserId:    "085f48c5-8763-4e22-a1c6-ac1a68ba07de",
	LoginName: "abbeymart",
	Email:     "abbeya1@yahoo.com",
	Language:  "en-US",
	Firstname: "Abi",
	Lastname:  "Akindele",
	Token:     "",
	Expire:    0,
	Role:      "TBD",
}

var userInfo = map[string]interface{}{
	"userId":    "085f48c5-8763-4e22-a1c6-ac1a68ba07de",
	"loginName": "abbeymart",
	"email":     "abbeya1@yahoo.com",
	"language":  "en-US",
	"firstname": "Abi",
	"lastname":  "Akindele",
	"token":     "",
	"expire":    0,
	"role":      "guest",
}

var TestCrudParamOptions = CrudOptionsType{
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
	"name": "services",
}
var GroupCreateRec2 = ActionParamType{
	"name": "services",
}

var GroupUpdateRec1 = ActionParamType{
	"name": "services",
}
var GroupUpdateRec2 = ActionParamType{
	"name": "services",
}

var CategoryCreateRec1 = ActionParamType{
	"name": "services",
}

var CategoryCreateRec2 = ActionParamType{
	"name": "services",
}

var CategoryUpdateRec1 = ActionParamType{
	"name": "services",
}

var CategoryUpdateRec2 = ActionParamType{
	"name": "services",
}

var GroupCreateActionParams = ActionParamsType{
	GroupCreateRec1,
	GroupCreateRec2,
}

var CategoryCreateActionParams = ActionParamsType{
	CategoryCreateRec1,
	CategoryCreateRec2,
}

// TODO: update and delete params (ids, queryParams)

var GroupUpdateRecordById = ActionParamType{
	"name": "services2",
}

var CategoryUpdateRecordById = ActionParamType{
	"name": "services2",
}

var GroupUpdateRecordByParam = ActionParamType{
	"name": "services2",
}

var CategoryUpdateRecordByParam = ActionParamType{
	"name": "services2",
}

var GroupUpdateIds = []string{"6900d9f9-2ceb-450f-9a9e-527eb66c962f", "122d0f0e-3111-41a5-9103-24fa81004550"}
var GroupUpdateParams = QueryParamType{}

var CategoryUpdateIds = []string{"6900d9f9-2ceb-450f-9a9e-527eb66c962f", "122d0f0e-3111-41a5-9103-24fa81004550"}
var CategoryUpdateParams = QueryParamType{}

var UpdateIds = []string{"6900d9f9-2ceb-450f-9a9e-527eb66c962f", "122d0f0e-3111-41a5-9103-24fa81004550"}
var UpdateParams = QueryParamType{}

var GroupUpdateActionParams = ActionParamsType{
	GroupUpdateRec1,
	GroupUpdateRec2,
}

var GroupUpdateActionParamsById = ActionParamsType{
	GroupUpdateRecordById,
}
var GroupUpdateActionParamsByParam = ActionParamsType{
	GroupUpdateRecordByParam,
	GroupUpdateRecordByParam,
}

// GetRecordType get record(s)
type GetRecordType struct {
	Id            string
	TableName     string
	LogRecords    interface{}
	NewLogRecords interface{}
	LogBy         string
	LogType       string
	LogAt         time.Time
}

// GetIds get by ids & params
var GetIds = []string{"6900d9f9-2ceb-450f-9a9e-527eb66c962f", "122d0f0e-3111-41a5-9103-24fa81004550"}
var GetParams = QueryParamType{}

// DeleteIds delete record(s) by ids & params
var DeleteIds = []string{"dba4adbb-4482-4f3d-bb05-0db80c30876b", "02f83bc1-8fa3-432a-8432-709f0df3f3b0"}
var DeleteParams = QueryParamType{}
