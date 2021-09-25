// @Author: abbeymart | Abi Akindele | @Created: 2020-12-22 | @Updated: 2020-12-22
// @Company: mConnect.biz | @License: MIT
// @Description: crud operations' types - updated

package mccrud

import (
	"fmt"
	"github.com/abbeymart/mcresponse"
	//"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jmoiron/sqlx"
	"time"
)

type DbConnectionType *sqlx.DB

type DbSecureType struct {
	SecureAccess bool   `json:"secureAccess"`
	SecureCert   string `json:"secureCert"`
	SecureKey    string `json:"secureKey"`
	SslMode      string `json:"sslMode"`
}

type DbConfigType struct {
	Host         string       `json:"host"`
	Username     string       `json:"username"`
	Password     string       `json:"password"`
	DbName       string       `json:"dbName"`
	Filename     string       `json:"filename"`
	Location     string       `json:"location"`
	Port         uint32       `json:"port"`
	DbType       string       `json:"dbType"`
	PoolSize     uint         `json:"poolSize"`
	Url          string       `json:"url"`
	SecureOption DbSecureType `json:"secureOption"`
}

type DbConnectOptions map[string]interface{}

type DbConfig struct {
	DbType       string           `json:"dbType"`
	Host         string           `json:"host"`
	Username     string           `json:"username"`
	Password     string           `json:"password"`
	DbName       string           `json:"dbName"`
	Filename     string           `json:"filename"`
	Location     string           `json:"location"`
	Port         uint32           `json:"port"`
	PoolSize     uint             `json:"poolSize"`
	Url          string           `json:"url"`
	Timezone     string           `json:"timezone"`
	SecureOption DbSecureType     `json:"secureOption"`
	Options      DbConnectOptions `json:"options"`
}

type CrudTasksType struct {
	Create string
	Insert string
	Update string
	Read   string
	Delete string
	Remove string
	Login  string
	Logout string
	Other  string
}

func CrudTasks() CrudTasksType {
	return CrudTasksType{
		Create: "create",
		Insert: "insert",
		Update: "update",
		Read:   "read",
		Delete: "delete",
		Remove: "remove",
		Login:  "login",
		Logout: "logout",
		Other:  "other",
	}
}

const (
	CreateTask  = "create"
	InsertTask  = "insert"
	UpdateTask  = "update"
	ReadTask    = "read"
	DeleteTask  = "delete"
	RemoveTask  = "remove"
	LoginTask   = "login"
	LogoutTask  = "logout"
	SystemTask  = "system"
	AppTask     = "app"
	UnknownTask = "unknown"
)

type UserInfoType struct {
	UserId    string `json:"userId" form:"userId" mcorm:"user_id"`
	Firstname string `json:"firstname" mcorm:"firstname"`
	Lastname  string `json:"lastname" mcorm:"lastname"`
	Language  string `json:"language" mcorm:"language"`
	LoginName string `json:"loginName" form:"loginName" mcorm:"login_name"`
	Token     string `json:"token" mcorm:"token"`
	Expire    int64  `json:"expire" mcorm:"expire"`
	Email     string `json:"email" form:"email" mcorm:"email"`
	Role      string `json:"role" mcorm:"role"`
}

type BaseModelType struct {
	Id        string    `json:"id" mcorm:"id"`
	Language  string    `json:"language" mcorm:"language"`
	Desc      string    `json:"desc" mcorm:"desc"`
	AppId     string    `json:"appId" mcorm:"app_id"`       // application-id in a multi-hosted apps environment (e.g. cloud-env)
	IsActive  bool      `json:"isActive" mcorm:"is_active"` // => activate by modelOptionsType settings...
	CreatedBy string    `json:"createdBy" mcorm:"created_by"`
	CreatedAt time.Time `json:"createdAt" mcorm:"created_at"`
	UpdatedBy string    `json:"updatedBy" mcorm:"updated_by"`
	UpdatedAt time.Time `json:"updatedAt" mcorm:"updated_at"`
	DeletedAt time.Time `json:"deletedAt" mcorm:"deleted_at"`
}

type EmailAddressType = map[string]string

type Profile struct {
	BaseModelType
	UserID        string
	Firstname     string             `json:"firstname"`
	Lastname      string             `json:"lastname"`
	Middlename    string             `json:"middlename"`
	Fullname      string             `json:"fullname"`
	Phone         string             `json:"phone"`
	RecEmail      string             `json:"recEmail"`
	RoleId        string             `json:"roleId"`
	Emails        []EmailAddressType `json:"emails"`
	DateOfBirth   time.Time          `json:"dateOfBirth"`
	TwoFactorAuth bool               `json:"twoFactorAuth"`
	AuthAgent     string             `json:"authAgent"`
	AuthPhone     string             `json:"authPhone"`
	PostalCode    string             `json:"postalCode"`
	UserInfo      UserInfoType       `json:"userInfo"`
	ReqUrl        string             `json:"reqUrl"`
}

type RoleServiceType struct {
	ServiceId            string   `json:"serviceId"`
	RoleId               string   `json:"roleId"`
	RoleIds              []string `json:"roleIds"`
	ServiceCategory      string   `json:"serviceCategory"`
	CanRead              bool     `json:"canRead"`
	CanCreate            bool     `json:"canCreate"`
	CanUpdate            bool     `json:"canUpdate"`
	CanDelete            bool     `json:"canDelete"`
	CanCrud              bool     `json:"canCrud"`
	TableAccessPermitted bool     `json:"tableAccessPermitted"`
}

type CheckAccessType struct {
	UserId       string            `json:"userId" mcorm:"userId"`
	RoleId       string            `json:"roleId" mcorm:"roleId"`
	RoleIds      []string          `json:"roleIds" mcorm:"roleIds"`
	IsActive     bool              `json:"isActive" mcorm:"isActive"`
	IsAdmin      bool              `json:"isAdmin" mcorm:"isAdmin"`
	RoleServices []RoleServiceType `json:"roleServices" mcorm:"roleServices"`
	TableId      string            `json:"tableId" mcorm:"tableId"`
}

type CheckAccessParamsType struct {
	AccessDb     *sqlx.DB     `json:"accessDb"`
	UserInfo     UserInfoType `json:"userInfo"`
	TableName    string       `json:"tableName"`
	RecordIds    []string     `json:"recordIds"` // for update, delete and read tasks
	AccessTable  string       `json:"accessTable"`
	UserTable    string       `json:"userTable"`
	RoleTable    string       `json:"roleTable"`
	ServiceTable string       `json:"serviceTable"`
	ProfileTable string       `json:"profileTable"`
}

type RoleFuncType func(it1 string, it2 RoleServiceType) bool
type FieldValueType interface{}
type ActionParamType map[string]interface{}
type ValueToDataType map[string]interface{}
type ActionParamsType []ActionParamType
type SortParamType map[string]int     // 1 for "asc", -1 for "desc"
type ProjectParamType map[string]bool // 1 or true for inclusion, 0 or false for exclusion
type QueryParamType map[string]interface{}

type QueryParamItemType struct {
	Query    QueryParamType `json:"query"`
	Order    int            `json:"order"`    // order
	Operator string         `json:"operator"` // relationship to the next group (AND, OR), the last group-operator is "" or ignored
}
type QueryParamsType []QueryParamItemType

type ModelOptionsType struct {
	TimeStamp   bool
	ActiveStamp bool
	ActorStamp  bool
}

// CrudParamsType is the struct type for receiving, composing and passing CRUD inputs
type CrudParamsType struct {
	ModelRef      interface{}      `json:"-"`
	ModelPointer  interface{}      `json:"-"`
	AppDb         *sqlx.DB         `json:"-"`
	TableName     string           `json:"-"`
	UserInfo      UserInfoType     `json:"userInfo"`
	ActionParams  ActionParamsType `json:"actionParams"`
	QueryParams   QueryParamType   `json:"queryParams"`
	RecordIds     []string         `json:"recordIds"`
	ProjectParams ProjectParamType `json:"projectParams"`
	SortParams    SortParamType    `json:"sortParams"`
	Token         string           `json:"token"`
	Skip          int              `json:"skip"`
	Limit         int              `json:"limit"`
	TaskName      string           `json:"-"`
	TaskType      string           `json:"-"`
}

type CrudOptionsType struct {
	CheckAccess           bool
	BulkCreate            bool
	AccessDb              *sqlx.DB
	AuditDb               *sqlx.DB
	ServiceDb             *sqlx.DB
	AuditTable            string
	ServiceTable          string
	UserTable             string
	RoleTable             string
	AccessTable           string
	VerifyTable           string
	ProfileTable          string
	MaxQueryLimit         int
	LogCrud               bool
	LogCreate             bool
	LogUpdate             bool
	LogRead               bool
	LogDelete             bool
	LogLogin              bool
	LogLogout             bool
	UnAuthorizedMessage   string
	RecExistMessage       string
	CacheExpire           int
	LoginTimeout          int
	UsernameExistsMessage string
	EmailExistsMessage    string
	MsgFrom               string
	ModelOptions          ModelOptionsType
	FieldSeparator        string
}

type SelectQueryOptions struct {
	Skip  int
	Limit int
}

type MessageObject map[string]string

type ValidateResponseType struct {
	Ok     bool          `json:"ok"`
	Errors MessageObject `json:"errors"`
}
type OkResponse struct {
	Ok bool `json:"ok"`
}

// CRUD operations

type CreateQueryObject struct {
	CreateQuery string
	FieldNames  []string
	FieldValues [][]interface{}
}

type WhereQueryObject struct {
	WhereQuery  string
	FieldValues []interface{}
}

type UpdateQueryObject struct {
	UpdateQuery string
	FieldNames  []string
	FieldValues []interface{}
	WhereQuery  WhereQueryObject
}

type DeleteQueryObject struct {
	DeleteQuery string
	FieldValues []interface{}
	WhereQuery  WhereQueryObject
}

type SelectQueryObject struct {
	SelectQuery string
	FieldValues []interface{}
	WhereQuery  WhereQueryObject
}

type CreateQueryResult struct {
	CreateQueryObject CreateQueryObject
	Ok                bool
	Message           string
}

type UpdateQueryResult struct {
	UpdateQueryObject UpdateQueryObject
	Ok                bool
	Message           string
}

type MultiUpdateQueryResult struct {
	UpdateQueryObjects []UpdateQueryObject
	Ok                 bool
	Message            string
}

type DeleteQueryResult struct {
	DeleteQueryObject DeleteQueryObject
	Ok                bool
	Message           string
}

type SelectQueryResult struct {
	SelectQueryObject SelectQueryObject
	Ok                bool
	Message           string
}

type WhereQueryResult struct {
	WhereQueryObject WhereQueryObject
	Ok               bool
	Message          string
}

// ErrorType provides the structure for error reporting
type ErrorType struct {
	Code    string
	Message string
}

type SaveError ErrorType
type CreateError ErrorType
type UpdateError ErrorType
type DeleteError ErrorType
type ReadError ErrorType
type AuthError ErrorType
type ConnectError ErrorType
type SelectQueryError ErrorType
type WhereQueryError ErrorType
type CreateQueryError ErrorType
type UpdateQueryError ErrorType
type DeleteQueryError ErrorType

// sample Error() implementation
func (err ErrorType) Error() string {
	return fmt.Sprintf("Error-code: %v | Error-message: %v", err.Code, err.Message)
}

type LogRecordsType struct {
	TableFields  []string       `json:"tableFields"`
	TableRecords []interface{}  `json:"tableRecords"`
	QueryParam   QueryParamType `json:"queryParam"`
	RecordIds    []string       `json:"recordIds"`
}

type CrudResultType struct {
	QueryParam   QueryParamType             `json:"queryParam"`
	RecordIds    []string                   `json:"recordIds"`
	RecordsCount int                        `json:"recordsCount"`
	Records      []map[string]interface{}   `json:"records"`
	TaskType     string                     `json:"taskType"`
	LogRes       mcresponse.ResponseMessage `json:"logRes"`
}

type GetStatType struct {
	Skip              int            `json:"skip"`
	Limit             int            `json:"limit"`
	RecordsCount      int            `json:"recordsCount"`
	TotalRecordsCount int            `json:"totalRecordsCount"`
	QueryParam        QueryParamType `json:"queryParam"`
	RecordIds         []string       `json:"recordIds"`
	Expire            int            `json:"expire"`
}

type GetResultType struct {
	Records  []map[string]interface{}   `json:"value"`
	Stats    GetStatType                `json:"stats"`
	LogRes   mcresponse.ResponseMessage `json:"logRes"`
	TaskType string                     `json:"taskType"`
}

type SaveResultType struct {
	QueryParam   QueryParamType             `json:"queryParam"`
	RecordIds    []string                   `json:"recordIds"`
	RecordsCount int                        `json:"recordsCount"`
	TaskType     string                     `json:"taskType"`
	LogRes       mcresponse.ResponseMessage `json:"logRes"`
}

// TODO: remove, not required

type SaveParamsType struct {
	UserInfo    UserInfoType   `json:"userInfo"`
	QueryParams QueryParamType `json:"queryParams"`
	RecordIds   []string       `json:"recordIds"`
	//ActionParams ActionParamsType `json:"actionParams"`
}

type DeleteParamsType struct {
	UserInfo    UserInfoType   `json:"userInfo"`
	RecordIds   []string       `json:"recordIds"`
	QueryParams QueryParamType `json:"queryParams"`
}

type GetParamsType struct {
	UserInfo     UserInfoType     `json:"userInfo"`
	Skip         int              `json:"skip"`
	Limit        int              `json:"limit"`
	RecordIds    []string         `json:"recordIds"`
	QueryParams  QueryParamType   `json:"queryParams"`
	SortParam    SortParamType    `json:"sortParams"`
	ProjectParam ProjectParamType `json:"projectParam"`
}

type SaveCrudParamsType struct {
	CrudParams         CrudParamsType
	CrudOptions        CrudOptionsType
	CreateTableFields  []string
	UpdateTableFields  []string
	GetTableFields     []string
	TableFieldPointers []interface{}
	AuditLog           bool
}

type DeleteCrudParamsType struct {
	CrudParams         CrudParamsType
	CrudOptions        CrudOptionsType
	GetTableFields     []string
	TableFieldPointers []interface{}
	AuditLog           bool
}

type GetCrudParamsType struct {
	CrudParams         CrudParamsType
	CrudOptions        CrudOptionsType
	GetTableFields     []string
	TableFieldPointers []interface{}
	AuditLog           bool
}
