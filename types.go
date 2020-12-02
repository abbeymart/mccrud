// @Author: abbeymart | Abi Akindele | @Created: 2020-12-01 | @Updated: 2020-12-01
// @Company: mConnect.biz | @License: MIT
// @Description: types definition

package mccrud

import (
	"database/sql"
)

// constants | enum
const (
	// CRUD Tasks
	CREATE = "create"
	INSERT = "insert"
	UPDATE = "update"
	READ   = "read"
	DELETE = "delete"
	REMOVE = "remove"
	// ORM Relations
	OneToOne   = "1To1"
	OneToMany  = "1ToM"
	ManyToMany = "MToM"
	ManyToOne  = "MTo1"
	// Relation Actions
	RESTRICT = "restrict" // must remove target-record(s), prior to removing source-record
	CASCADE  = "cascade"  // default for ON UPDATE | update foreignKey value or delete foreignKey record/value
	NoAction = "noAction" // leave the foreignKey value, as-is
	DEFAULT  = "default"  // set foreignKey to specified default value
	NULL     = "null"     // set foreignKey value to null or ""
	// DataTypes
	STRING = "string"
	// STRINGALPHA = "stringalpha"
	// STRINGALPHANUMERIC = "stringalphanumeric"
	POSTALCODE     = "postalcode"
	MONGODBID      = "objectid"
	UUID           = "uuid"
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
	EMAIL          = "email"
	URL            = "url"
	PORT           = "port"
	IP             = "ipaddress"
	JWT            = "jwt"
	LATLONG        = "latlong"
	ISO2           = "iso2"
	ISO3           = "iso3"
	MACADDRESS     = "macaddress"
	MIME           = "mime"
	CREDITCARD     = "creditcard"
	CURRENCY       = "currency"
	IMEI           = "imei"
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
	userId    string
	firstName string
	lastName  string
	language  string
	loginName string
	token     string
	expire    uint
	group     string
	email     string
}

type RoleServiceType struct {
	serviceId           string
	roleId              string
	serviceCategory     string
	canRead             bool
	canCreate           bool
	canUpdate           bool
	canDelete           bool
	tableAccessPermitted bool
}

type CheckAccessType struct {
	userId       string
	group        string
	groups       []string
	isActive     bool
	isAdmin      bool
	roleServices []RoleServiceType
	tableId       string
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
	appDb         *sql.DB
	Table         string
	userInfo      UserInfoType
	actionParams  ActionParamsType
	existParams   ExistParamsType
	queryParams   QueryParamType
	docIds        []string
	projectParams ProjectParamType
	sortParams    SortParamType
	token         string
	options       CrudOptionsType
	taskName      string
}

type ModelOptionsType struct {
	timeStamp    bool // auto-add: createdAt and updatedAt | default: true
	actorStamp   bool // auto-add: createdBy and updatedBy | default: true
	activeStamp  bool // auto-add isActive, if not already set | default: true
	docValueDesc DocDescType
	docValue     ValueParamType
}

type CrudOptionsType struct {
	skip                  uint
	limit                 uint
	parentColls           []string
	childColls            []string
	recursiveDelete       bool
	checkAccess           bool
	accessDb              *sql.DB
	auditDb               *sql.DB
	serviceDb             *sql.DB
	auditTable            string
	serviceTable          string
	userTable             string
	roleTable             string
	accessTable           string
	verifyTable           string
	maxQueryLimit         uint
	logAll                bool
	logCreate             bool
	logUpdate             bool
	logRead               bool
	logDelete             bool
	logLogin              bool
	logLogout             bool
	unAuthorizedMessage   string
	recExistMessage       string
	cacheExpire           uint
	modelOptions          ModelOptionsType
	loginTimeout          uint
	usernameExistsMessage string
	emailExistsMessage    string
	msgFrom               string
}

type ErrorType map[string]string
type ValidateResponseType struct {
	ok     bool
	errors ErrorType
}
type OkResponse struct {
	ok bool
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
	fieldType       string
	fieldLength     uint   // default: 255 for DataTypes.STRING
	fieldPattern    string // "/^[0-9]{10}$/" => includes 10 digits, 0 to 9 | "/^[0-9]{6}.[0-9]{2}$/ => max 16 digits and 2 decimal places
	allowNull       bool   // default: true
	unique          bool
	indexable       bool
	primaryKey      bool
	minValue        uint
	maxValue        uint
	setValue        SetValueType       // set/transform fieldValue prior to save(create/insert), T=>fieldType
	defaultValue    interface{}        // result/T must be of fieldType
	validate        ValidateMethodType // T=>fieldType, returns a bool (valid=true/invalid=false)
	validateMessage string
}

type ModelRelationType struct {
	sourceTable   string
	targetTable   string
	sourceField   string
	targetField   string
	relationType  string
	sourceModel   ModelType
	targetModel   ModelType
	foreignField  string // source-to-targetField map
	relationField string // relation-targetField, for many-to-many
	relationTable string // optional tableName for many-to-many | default: source_target TableTableor sourceTarget
	onDelete      string
	onUpdate      string
}

type ModelType struct {
	tableName        string
	docDesc         DocDescType
	timeStamp       bool // auto-add: createdAt and updatedAt | default: true
	actorStamp      bool // auto-add: createdBy and updatedBy | default: true
	activeStamp     bool // record active status, isActive (true | false) | default: true
	relations       []ModelRelationType
	computedMethods ComputedMethodsType // model-level functions, e.g fullName(a, b: T): T
	validateMethods ValidateMethodsType
	alterSyncTable  bool // create / alter table/collection and sync existing data, if there was a change to the table structure | default: true
	// if alterSyncTable: false it will create/re-create the table, with no data sync
}
