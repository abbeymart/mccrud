// @Author: abbeymart | Abi Akindele | @Created: 2020-12-01 | @Updated: 2020-12-01
// @Company: mConnect.biz | @License: MIT
// @Description: types definition

package mccrud

import (
	"database/sql"
	"github.com/abbeymart/mcutilsgo"
	"go.mongodb.org/mongo-driver/mongo"
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

type FieldOperatorType struct {
	Equals              string
	GreaterThan         string
	LessThan            string
	GreaterThanOrEquals string
	LessThanOrEquals    string
	NotEquals           string
	True                string
	False               string
	Includes            string
	NotIncludes         string
	StartsWith          string
	EndsWith            string
	NotStartsWith       string
	NotEndsWith         string
}

type GroupOperatorType struct {
	AND string
	OR  string
}

type DataType struct {
	String             string
	Text               string
	StringAlpha        string
	StringAlphaNumeric string
	PostCode           string
	MongoDBId          string
	UUID               string
	UUID3              string
	UUID4              string
	UUID5              string
	MD4                string
	MD5                string
	SHA1               string
	SHA256             string
	SHA384             string
	SHA512             string
	Number             string
	Integer            string
	Decimal            string
	Float              string
	Float32            string
	Float64            string
	BigInt             string
	BigFloat           string
	Object             string
	Array              string
	ArrayOfString      string
	ArrayOfNumber      string
	ArrayOfBoolean     string
	ArrayOfObject      string
	Boolean            string
	JSON               string
	DateTime           string
	Date               string
	Time               string
	TimeStamp          string
	TimeStampZ         string
	Positive           string
	Natural            string
	Negative           string
	Email              string
	URL                string
	DomainName         string
	Port               string
	IP                 string
	IP4                string
	IP6                string
	JWT                string
	Latitude           string
	Longitude          string
	ISO2               string
	ISO3               string
	MACAddress         string
	Mime               string
	CreditCard         string
	Currency           string
	IMEI               string
	Set                string
	Map                string
	Undefined          string
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

func FieldOperators() FieldOperatorType {
	return FieldOperatorType{
		Equals:              "Equals",
		GreaterThan:         "GreaterThan",
		LessThan:            "LessThan",
		GreaterThanOrEquals: "GreaterThanOrEquals",
		LessThanOrEquals:    "LessThanOrEquals",
		NotEquals:           "NotEquals",
		True:                "True",
		False:               "False",
		Includes:            "Includes",
		NotIncludes:         "NotIncludes",
		StartsWith:          "StartsWith",
		EndsWith:            "EndsWith",
		NotStartsWith:       "NotStartsWith",
		NotEndsWith:         "NotEndsWith",
	}
}

func GroupOperators() GroupOperatorType {
	return GroupOperatorType{
		AND: "And",
		OR:  "Or",
	}
}

func DataTypes() DataType {
	return DataType{
		String:             "string",
		Text:               "text",
		StringAlpha:        "stringalpha",
		StringAlphaNumeric: "stringalphanumeric",
		PostCode:           "postalcode",
		MongoDBId:          "mongodbid",
		UUID:               "uuid",
		UUID3:              "uuid3",
		UUID4:              "uuid4",
		UUID5:              "uuid5",
		MD4:                "md4",
		MD5:                "md5",
		SHA1:               "sha1",
		SHA256:             "sha2",
		SHA384:             "sha3",
		SHA512:             "sha5",
		Number:             "number",
		Integer:            "integer",
		Decimal:            "decimal",
		Float:              "float",
		Float32:            "float32",
		Float64:            "float64",
		BigInt:             "bigint",
		BigFloat:           "bigfloat",
		Object:             "object",
		Array:              "array",
		ArrayOfString:      "arrayofstring",
		ArrayOfNumber:      "arrayofnumber",
		ArrayOfBoolean:     "arrayofboolean",
		Boolean:            "boolean",
		JSON:               "json",
		DateTime:           "datetime",
		Date:               "date",
		Time:               "time",
		TimeStamp:          "timestamp",
		TimeStampZ:         "timestampz",
		Positive:           "positive",
		Natural:            "natural",
		Negative:           "negative",
		Email:              "email",
		URL:                "url",
		DomainName:         "domainname",
		Port:               "port",
		IP:                 "ip",
		IP4:                "ip4",
		IP6:                "ip6",
		JWT:                "jwt",
		Latitude:           "latitude",
		Longitude:          "longitude",
		ISO2:               "iso2",
		ISO3:               "iso3",
		MACAddress:         "macaddress",
		Mime:               "mime",
		CreditCard:         "creditcard",
		Currency:           "currency",
		IMEI:               "imei",
		Set:                "set",
		Map:                "map",
		Undefined:          "undefined",
	}
}

const (
	// TaskType CRUD Tasks
	CreateTask = "create"
	InsertTask = "insert"
	UpdateTask = "update"
	ReadTask   = "read"
	DeleteTask = "delete"
	RemoveTask = "remove"
	// Model Relations
	OneToOneRelation   = "onetoone"
	OneToManyRelation  = "onetomany"
	ManyToManyRelation = "manytomany"
	ManyToOneRelation  = "manytoone"
	// Model Relation Actions
	RestrictAction = "restrict" // must remove target-record(s), prior to removing source-record
	CascadeAction  = "cascade"  // default for ON UPDATE | update foreignKey value or delete foreignKey record/value
	NoAction       = "noaction" // leave the foreignKey value, as-is
	DefaultAction  = "default"  // set foreignKey to specified default value
	NullAction     = "null"     // set foreignKey value to null or ""
	// DataType
	STRING = "string"
	// STRINGALPHA = "stringalpha"
	// STRINGALPHANUMERIC = "stringalphanumeric"
	POSTALCODE     = "postalcode"
	MONGODBID      = "objectid"
	UUID           = "uuid"
	UUID3          = "uuid3"
	UUID4          = "uuid4"
	UUID5          = "uuid5"
	MD4            = "md4"
	MD5            = "md5"
	SHA1           = "sha1"
	SHA256         = "sha256"
	SHA384         = "sha384"
	SHA512         = "sha512"
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
	UserId    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Language  string `json:"language"`
	LoginName string `json:"login_name"`
	Token     string `json:"token"`
	Expire    uint   `json:"expire"`
	Group     string `json:"group"`
	Email     string `json:"email"`
}

type RoleServiceType struct {
	ServiceId            string `json:"service_id"`
	RoleId               string `json:"role_id"`
	ServiceCategory      string `json:"service_category"`
	CanRead              bool   `json:"can_read"`
	CanCreate            bool   `json:"can_create"`
	CanUpdate            bool   `json:"can_update"`
	CanDelete            bool   `json:"can_delete"`
	TableAccessPermitted bool   `json:"table_access_permitted"`
}

type CheckAccessType struct {
	UserId       string            `json:"user_id"`
	Group        string            `json:"group"`
	Groups       []string          `json:"groups"`
	IsActive     bool              `json:"is_active"`
	IsAdmin      bool              `json:"is_admin"`
	RoleServices []RoleServiceType `json:"role_services"`
	TableId      string            `json:"table_id"`
}

type CheckAccessParamsType struct {
	accessDb     *sql.DB
	userInfo     UserInfoType
	tableName    string
	docIds       []string // for update, delete and read tasks
	accessTable  string
	userTable    string
	roleTable    string
	serviceTable string
}

type RoleFuncType func(it1 string, it2 RoleServiceType) bool
type FieldValueType interface{}
type ValueParamType map[string]interface{}
type ValueToDataType map[string]interface{}
type ActionParamsType []ValueParamType
type ExistParamType map[string]interface{}
type ExistParamsType []ExistParamType
type SortParamType map[string]int     // 1 for "asc", -1 for "desc"
type ProjectParamType map[string]bool // 1 or true for inclusion, 0 or false for exclusion

type GroupItemType struct {
	GroupItem      map[string]map[string]interface{} // key1 => fieldName, key2 => fieldOperator, interface{}=> value(s)
	GroupItemOrder uint		// item/field order within the group
}

type GroupParamType struct {
	GroupName   string          // for group-items(fields) categorization
	GroupItems  []GroupItemType // group items to be composed by category
	GroupOrder  uint             // group order
	GroupLinkOp string          // group relationship to the next group (AND, OR), the last group groupLinkOp should be "" or will be ignored
}

type QueryParamType []GroupParamType
type WhereParamType []GroupParamType

type ModelOptionsType struct {
	TimeStamp    bool // auto-add: createdAt and updatedAt | default: true
	ActorStamp   bool // auto-add: createdBy and updatedBy | default: true
	ActiveStamp  bool // auto-add isActive, if not already set | default: true
	DocValueDesc RecordDescType
	DocValue     ValueParamType
}

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

type CrudParamType struct {
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

type MongoCrudTaskType struct {
	AppDb         *mongo.Client
	TableName     string
	UserInfo      UserInfoType
	ActionParams  ActionParamsType
	ExistParams   ExistParamsType
	QueryParams   QueryParamType
	DocIds        []string
	ProjectParams ProjectParamType
	SortParams    SortParamType
	Token         string
	Options       MongoCrudOptionsType
	TaskName      string
}

type MongoCrudOptionsType struct {
	Skip                  uint
	Limit                 uint
	ParentTables          []string
	ChildTables           []string
	RecursiveDelete       bool
	CheckAccess           bool
	AccessDb              *mongo.Client
	AuditDb               *mongo.Client
	ServiceDb             *mongo.Client
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

type MongoCrudParamType struct {
	appDb           *mongo.Client
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
	accessDb        *mongo.Client
	auditDb         *mongo.Client
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
	params              MongoCrudTaskType
}

type ErrorType map[string]string
type ValidateResponseType struct {
	Ok     bool                  `json:"ok"`
	Errors mcutils.MessageObject `json:"errors"`
}
type OkResponse struct {
	Ok bool `json:"ok"`
}

// ORM types
type RecordValueType map[string]ValueParamType
type RecordDescType map[string]FieldDescType

type GetValueType func() interface{}
type SetValueType func(val interface{}) interface{}
type DefaultValueType func() interface{}
type ValidateMethodType func(val interface{}) bool
type ValidateMethodResponseType func(val interface{}) ValidateResponseType
type ComputedValueType func(val interface{}) interface{}

type ValidateMethodsType map[string]ValidateMethodResponseType
type ComputedMethodsType map[string]ComputedValueType

type FieldDescType struct {
	FieldType       string
	FieldLength     uint   // default: 255 for DataType.STRING
	FieldPattern    string // "/^[0-9]{10}$/" => includes 10 digits, 0 to 9 | "/^[0-9]{6}.[0-9]{2}$/ => max 16 digits and 2 decimal places
	AllowNull       bool   // default: true
	Unique          bool
	Indexable       bool
	PrimaryKey      bool
	MinValue        uint
	MaxValue        uint
	SetValue        SetValueType       // set/transform fieldValue prior to save(create/insert), T=>fieldType
	DefaultValue    DefaultValueType   // result/T must be of fieldType
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
	RecordDesc      RecordDescType
	TimeStamp       bool // auto-add: createdAt and updatedAt | default: true
	ActorStamp      bool // auto-add: createdBy and updatedBy | default: true
	ActiveStamp     bool // record active status, isActive (true | false) | default: true
	Relations       []ModelRelationType
	ComputedMethods ComputedMethodsType // model-level functions, e.g fullName(a, b: T): T
	ValidateMethods ValidateMethodsType
	AlterSyncTable  bool // create / alter table/collection and sync existing data, if there was a change to the table structure | default: true
	// if alterSyncTable: false it will create/re-create the table, with no data sync
}
