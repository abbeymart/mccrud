// @Author: abbeymart | Abi Akindele | @Created: 2020-12-14 | @Updated: 2020-12-14
// @Company: mConnect.biz | @License: MIT
// @Description: go: mConnect

package mccrud

import (
	"fmt"
	"github.com/abbeymart/mcauditlog"
	"github.com/abbeymart/mcdb"
	"github.com/abbeymart/mctest"
	"github.com/abbeymart/mctypes"
	"testing"
)

func TestSave(t *testing.T) {
	myDb := mcdb.DbConfig{
		DbType:   "postgres",
		Host:     "localhost",
		Username: "postgres",
		Password: "ab12testing",
		Port:     5432,
		DbName:   "mcdev",
		Filename: "testdb.db",
		PoolSize: 20,
		Url:      "localhost:5432",
	}
	myDb.Options = mcdb.DbConnectOptions{}

	// db-connection
	dbc, err := myDb.OpenPgxDbPool()
	//fmt.Printf("*****dbc-info: %v\n", dbc)
	// defer dbClose
	defer myDb.ClosePgxDbPool()
	// check db-connection-error
	if err != nil {
		fmt.Printf("*****db-connection-error: %v\n", err.Error())
		return
	}
	// expected db-connection result
	mcLogResult := mcauditlog.PgxLogParam{AuditDb: dbc.DbConn, AuditTable: TestAuditTable}
	// audit-log instance
	mcLog := mcauditlog.NewAuditLogPgx(dbc.DbConn, TestAuditTable)

	// json-records
	//fmt.Println("table-records-json", string(TableRecords))
	//fmt.Println("new-table-records-json", string(NewTableRecords))

	createCrudParams := mctypes.CrudParamsType{
		AppDb:        dbc.DbConn,
		TableName:    TestTable,
		UserInfo:     TestUserInfo,
		ActionParams: CreateActionParams,
	}

	var crudCreate interface{} = NewCrud(createCrudParams, TestCrudParamOptions)

	mctest.McTest(mctest.OptionValue{
		Name: "should connect to the Audit-DB and return an instance object:",
		TestFunc: func() {
			mctest.AssertEquals(t, err, nil, "error-response should be: nil")
			mctest.AssertEquals(t, mcLog, mcLogResult, "db-connection instance should be: "+mcLogResult.String())
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should connect to the CRUD-object and return an instance object:",
		TestFunc: func() {
			_, ok := crudCreate.(Crud)
			mctest.AssertEquals(t, ok, true, "crudCreate should be instance of mccrud.Crud")
		},
	})

	//mctest.McTest(mctest.OptionValue{
	//	Name: "should store create-transaction log and return success:",
	//	TestFunc: func() {
	//		res, err := mcLog.AuditLog(mcauditlog.CreateLog, UserId, mcauditlog.PgxAuditLogOptionsType{
	//			TableName:  TestTable,
	//			LogRecords: string(TableRecords),
	//		})
	//		//fmt.Printf("create-log: %v", res)
	//		mctest.AssertEquals(t, err, nil, "error-response should be: nil")
	//		mctest.AssertEquals(t, res.Code, "success", "log-action response-code should be: success")
	//	},
	//})

	mctest.PostTestResult()
}