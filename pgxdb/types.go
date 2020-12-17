// @Author: abbeymart | Abi Akindele | @Created: 2020-12-11 | @Updated: 2020-12-11
// @Company: mConnect.biz | @License: MIT
// @Description: go: mConnect

package pgxdb

import (
	"github.com/abbeymart/mccrud"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PgxCheckAccessParamsType struct {
	accessDb     *pgxpool.Pool
	userInfo     mccrud.UserInfoType
	tableName    string
	docIds       []string // for update, delete and read tasks
	accessTable  string
	userTable    string
	roleTable    string
	serviceTable string
}

type PgxCrudTaskType struct {
	AppDb         *pgxpool.Pool			// use *pgxpool.Pool, preferred || *pgx.Conn
	TableName     string
	UserInfo      mccrud.UserInfoType
	ActionParams  mccrud.ActionParamsType
	ExistParams   mccrud.ExistParamsType
	QueryParams   mccrud.QueryParamType
	DocIds        []string
	ProjectParams mccrud.ProjectParamType
	SortParams    mccrud.SortParamType
	Token         string
	Options       mccrud.CrudOptionsType
	TaskName      string
}

type PgxCrudOptionsType struct {
	Skip                  uint
	Limit                 uint
	ParentTables          []string
	ChildTables           []string
	RecursiveDelete       bool
	CheckAccess           bool
	AccessDb              *pgxpool.Pool
	AuditDb               *pgxpool.Pool
	ServiceDb             *pgxpool.Pool
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
	ModelOptions          mccrud.ModelOptionsType
	LoginTimeout          uint
	UsernameExistsMessage string
	EmailExistsMessage    string
	MsgFrom               string
}

type PgxCrudParamType struct {
	appDb           *pgxpool.Pool
	tableName       string
	token           string
	userInfo        mccrud.UserInfoType
	userId          string
	group           string
	groups          []string
	docIds          []string
	actionParams    mccrud.ActionParamsType
	queryParams     mccrud.QueryParamType
	existParams     mccrud.ExistParamsType
	projectParams   mccrud.ProjectParamType
	sortParams      mccrud.SortParamType
	skip            uint
	limit           uint
	parentTables    []string
	childTables     []string
	recursiveDelete bool
	checkAccess     bool
	accessDb        *pgxpool.Pool
	auditDb         *pgxpool.Pool
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
	createItems         mccrud.ActionParamsType
	updateItems         mccrud.ActionParamsType
	currentRecs         mccrud.ActionParamsType
	roleServices        []mccrud.RoleServiceType
	subItems            []bool
	cacheExpire         uint
	params              mccrud.CrudParamsType
}
