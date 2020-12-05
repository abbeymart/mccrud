// @Author: abbeymart | Abi Akindele | @Created: 2020-12-01 | @Updated: 2020-12-01
// @Company: mConnect.biz | @License: MIT
// @Description: types definition

package mccrud

import (
	"database/sql"
)

type TaskType struct {
	Create string
	Insert string
	Update string
	Read   string
	Delete string
	Remove string
}

type ORMRelationType struct {
	OneToOne   string
	OneToMany  string
	ManyToMany string
	ManyToOne  string
}

type RelationActionType struct {
	Restrict string // must remove target-record(s), prior to removing source-record
	Cascade  string // default for ON UPDATE | update foreignKey value or delete foreignKey record/value
	NoAction string // leave the foreignKey value, as-is
	Default  string // set foreignKey to specified default value
	Null     string // set foreignKey value to null or ""
}

// constants | enum
func CrudTasks() TaskType {
	return TaskType{
		// CRUD Tasks "create"
		Create: "Create",
		Insert: "Insert",
		Update: "Update",
		Read:   "Read",
		Delete: "Delete",
		Remove: "Remove",
	}
}

func ORMRelations() ORMRelationType {
	return ORMRelationType{
		OneToOne:   "OneToOne",
		OneToMany:  "OneToMany",
		ManyToMany: "ManyToMany",
		ManyToOne:  "ManyToOne",
	}
}

func RelationActions() RelationActionType {
	return RelationActionType{
		Restrict: "Restrict", // must remove target-record(s), prior to removing source-record
		Cascade:  "Cascade",  // default for ON UPDATE | update foreignKey value or delete foreignKey record/value
		NoAction: "NoAction", // leave the foreignKey value, as-is
		Default:  "Default",  // set foreignKey to specified default value
		Null:     "Null",     // set foreignKey value to null or ""
	}
}

const (
	// CRUD Tasks "create"
	CREATE = "create"
	INSERT = "insert"
	UPDATE = "update"
	READ   = "read"
	DELETE = "delete"
	REMOVE = "remove"
	// DataTypes
	STRING = "string"
	// STRINGALPHA = "stringalpha"
	// STRINGALPHANUMERIC = "stringalphanumeric"
	POSTALCODE = "postalcode"
	MONGODBID   = "objectid"
	UUID        = "uuid"
	UUID3       = "uuid3"
	UUID4       = "uuid4"
	UUID5       = "uuid5"
	MD4         = "md4"
	MD5         = "md5"
	SHA1        = "sha1"
	SHA256      = "sha256"
	SHA384      = "sha384"
	SHA512      = "sha512"
	NUMBER         = "number"
	INTEGER        = "integer"
	DECIMAL        = "decimal"
	FLOAT          = "float"
	BIGINT         = "bigint"
	BIGFLOAT       = "bigfloat"
	OBJECT         = "object" // struct, map...
	ARRAY          = "array"
	ARRAYOFSTRING  = "arrayofstring"
	ARRAYOFNUMBER  = "arrayofnumber"
	ARRAYOFBOOLEAN = "arrayofboolean"
	ARRAYOFOBJECT  = "arrayofobject"
	BOOLEAN        = "bool"
	JSON           = "json"
	DATETIME       = "datetime"
	DATE           = "date"
	TIME           = "time"
	TIMESTAMP      = "timestamp"
	TIMESTAMPZ     = "timestampz"
	POSITIVE       = "positive"
	NATURAL        = "natural"
	NEGATIVE       = "negative"
	EMAIL          = "email"
	URL            = "url"
	DOMAINNAME     = "domainname"
	PORT           = "port"
	IP             = "ipaddress"
	IP4            = "ipaddress4"
	IP6            = "ipaddress6"
	JWT            = "jwt"
	LATITUDE       = "latitude"
	LONGITUDE      = "longitude"
	//LATLONG        = "latlong"
	ISO2       = "iso2"
	ISO3       = "iso3"
	MACADDRESS = "macaddress"
	MIME       = "mime"
	CREDITCARD = "creditcard"
	CURRENCY   = "currency"
	IMEI       = "imei"
	// ENUM = "enum"       // Enumerations
	SET = "set" // Unique values set
	// WEAKSET = "weakset"
	MAP = "map" // Table/Map/Dictionary
	// WEAKMAP = "weakmap"
	MCDB = "mcdb" // Database connection handle
	// MODEL = "model"   // Model record definition
	// MODELVALUE = "modelvalue"
	UNDEFINED = "undefined"
)

type UserInfoType struct {
	UserId    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Language  string `json:"language"`
	LoginName string `json:"loginName"`
	Token     string `json:"token"`
	Expire    uint   `json:"expire"`
	Group     string `json:"group"`
	Email     string `json:"email"`
}

type RoleServiceType struct {
	ServiceId            string `json:"serviceId"`
	RoleId               string `json:"roleId"`
	ServiceCategory      string `json:"serviceCategory"`
	CanRead              bool   `json:"canRead"`
	CanCreate            bool   `json:"canCreate"`
	CanUpdate            bool   `json:"canUpdate"`
	CanDelete            bool   `json:"canDelete"`
	TableAccessPermitted bool   `json:"tableAccessPermitted"`
}

type CheckAccessType struct {
	UserId       string            `json:"userId"`
	Group        string            `json:"group"`
	Groups       []string          `json:"groups"`
	IsActive     bool              `json:"isActive"`
	IsAdmin      bool              `json:"isAdmin"`
	RoleServices []RoleServiceType `json:"roleServices"`
	TableId      string            `json:"tableId"`
}

type RoleFuncType func(it1 string, it2 RoleServiceType) bool
type FieldValueType interface{}
type ValueParamType map[string]interface{}
type ValueToDataType map[string]interface{}
type ActionParamsType []ValueParamType
type QueryParamType map[string]interface{}
type ExistParamType map[string]interface{}
type ExistParamsType []ExistParamType
type SortParamType map[string]int     // 1 for "asc", -1 for "desc"
type ProjectParamType map[string]bool // 1 or true for inclusion, 0 or false for exclusion

type CrudTaskType struct {
	AppDb         *sql.DB
	TableName     string
	UserInfo      UserInfoType
	ActionParams  ActionParamsType
	ExistParams   ExistParamsType
	QueryParams   QueryParamType
	DocIds        []string
	ProjectParams ProjectParamType
	SortParams    SortParamType
	Token         string
	Options       CrudOptionsType
	TaskName      string
}

type ModelOptionsType struct {
	TimeStamp    bool // auto-add: createdAt and updatedAt | default: true
	ActorStamp   bool // auto-add: createdBy and updatedBy | default: true
	ActiveStamp  bool // auto-add isActive, if not already set | default: true
	DocValueDesc DocDescType
	DocValue     ValueParamType
}

type CrudOptionsType struct {
	Skip                  uint
	Limit                 uint
	ParentTables          []string
	ChildTables           []string
	RecursiveDelete       bool
	CheckAccess           bool
	AccessDb              *sql.DB
	AuditDb               *sql.DB
	ServiceDb             *sql.DB
	AuditTable            string
	ServiceTable          string
	UserTable             string
	RoleTable             string
	AccessTable           string
	VerifyTable           string
	MaxQueryLimit         uint
	LogAll                bool
	LogCreate             bool
	LogUpdate             bool
	LogRead               bool
	LogDelete             bool
	LogLogin              bool
	LogLogout             bool
	UnAuthorizedMessage   string
	RecExistMessage       string
	CacheExpire           uint
	ModelOptions          ModelOptionsType
	LoginTimeout          uint
	UsernameExistsMessage string
	EmailExistsMessage    string
	MsgFrom               string
}

type CrudParamTYpe struct {
	appDb           *sql.DB
	tableName       string
	token           string
	userInfo        UserInfoType
	userId          string
	group           string
	groups          []string
	docIds          []string
	actionParams    ActionParamsType
	queryParams     QueryParamType
	existParams     ExistParamsType
	projectParams   ProjectParamType
	sortParams      SortParamType
	skip            uint
	limit           uint
	parentTables    []string
	childTables     []string
	recursiveDelete bool
	checkAccess     bool
	accessDb        *sql.DB
	auditDb         *sql.DB
	auditTable      string
	serviceTable    string
	userTable       string
	roleTable       string
	accessTable     string
	maxQueryLimit   uint
	logAll          bool
	logCreate       bool
	logUpdate       bool
	logRead         bool
	logDelete       bool
	//transLog AuditLog
	hashKey             string
	isRecExist          bool
	actionAuthorized    bool
	unAuthorizedMessage string
	recExistMessage     string
	isAdmin             bool
	createItems         ActionParamsType
	updateItems         ActionParamsType
	currentRecs         ActionParamsType
	roleServices        []RoleServiceType
	subItems            []bool
	cacheExpire         uint
	params              CrudTaskType
}

type ErrorType map[string]string
type ValidateResponseType struct {
	Ok     bool      `json:"ok"`
	Errors ErrorType `json:"errors"`
}
type OkResponse struct {
	Ok bool `json:"ok"`
}

// ORM types
type DocValueType map[string]ValueParamType
type DocDescType map[string]interface{}

type GetValueType func() interface{}
type SetValueType func(val interface{}) interface{}
type DefaultValueType func(val interface{}) interface{}
type ValidateMethodType func(val interface{}) bool
type ValidateMethodResponseType func(val interface{}) ValidateResponseType
type ComputedValueType func(val interface{}) interface{}

type ValidateMethodsType map[string]ValidateMethodResponseType
type ComputedMethodsType map[string]ComputedValueType

type FieldDescType struct {
	FieldType       string
	FieldLength     uint   // default: 255 for DataTypes.STRING
	FieldPattern    string // "/^[0-9]{10}$/" => includes 10 digits, 0 to 9 | "/^[0-9]{6}.[0-9]{2}$/ => max 16 digits and 2 decimal places
	AllowNull       bool   // default: true
	Unique          bool
	Indexable       bool
	PrimaryKey      bool
	MinValue        uint
	MaxValue        uint
	SetValue        SetValueType       // set/transform fieldValue prior to save(create/insert), T=>fieldType
	DefaultValue    interface{}        // result/T must be of fieldType
	Validate        ValidateMethodType // T=>fieldType, returns a bool (valid=true/invalid=false)
	ValidateMessage string
}

type ModelRelationType struct {
	SourceTable   string
	TargetTable   string
	SourceField   string
	TargetField   string
	RelationType  string
	SourceModel   ModelType
	TargetModel   ModelType
	ForeignField  string // source-to-targetField map
	RelationField string // relation-targetField, for many-to-many
	RelationTable string // optional tableName for many-to-many | default: source_target TableTableor sourceTarget
	OnDelete      string
	OnUpdate      string
}

type ModelType struct {
	TableName       string
	DocDesc         DocDescType
	TimeStamp       bool // auto-add: createdAt and updatedAt | default: true
	ActorStamp      bool // auto-add: createdBy and updatedBy | default: true
	ActiveStamp     bool // record active status, isActive (true | false) | default: true
	Relations       []ModelRelationType
	ComputedMethods ComputedMethodsType // model-level functions, e.g fullName(a, b: T): T
	ValidateMethods ValidateMethodsType
	AlterSyncTable  bool // create / alter table/collection and sync existing data, if there was a change to the table structure | default: true
	// if alterSyncTable: false it will create/re-create the table, with no data sync
}
