// @Author: abbeymart | Abi Akindele | @Created: 2020-12-14 | @Updated: 2020-12-14
// @Company: mConnect.biz | @License: MIT
// @Description: mccrud create & update records test-cases

package mccrud

import (
	"fmt"
	"github.com/abbeymart/mcauditlog"
	"github.com/abbeymart/mccrud/test/config/secure"
	"github.com/abbeymart/mcdb"
	"github.com/abbeymart/mctest"
	"testing"
)

func TestSaveGroup(t *testing.T) {
	myDb := secure.MyDb
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
	mcLogResult := mcauditlog.PgxLogParam{AuditDb: dbc.DbConn, AuditTable: AuditTable}
	// audit-log instance
	mcLog := mcauditlog.NewAuditLogPgx(dbc.DbConn, AuditTable)
	// group-table-records
	modelRef := Audit{}
	crudParams := CrudParamsType{
		AppDb:        dbc.DbConn,
		ModelRef:     modelRef,
		TableName:    "",
		UserInfo:     TestUserInfo,
		ActionParams: nil,
		RecordIds:    []string{},
		QueryParams:  QueryParamType{},
	}
	var crud = NewCrud(crudParams, CrudParamOptions)

	mctest.McTest(mctest.OptionValue{
		Name: "should connect to the Audit-DB and return an instance object:",
		TestFunc: func() {
			mctest.AssertEquals(t, err, nil, "error-response should be: nil")
			mctest.AssertEquals(t, mcLog, mcLogResult, "db-connection instance should be: "+mcLogResult.String())
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should connect to the CRUD-DB and return an instance object:",
		TestFunc: func() {
			mctest.AssertEquals(t, crud != nil, true, "crud should be instance of mccrud.Crud")
		},
	})

	// group-table test-cases
	mctest.McTest(mctest.OptionValue{
		Name: "should create two new group-records and return success:",
		TestFunc: func() {
			res := crud.SaveRecord()
			fmt.Println(res.Message, res.ResCode)
			value, _ := res.Value.(CrudResultType)
			mctest.AssertEquals(t, res.Code, "success", "save-create should return code: success")
			mctest.AssertEquals(t, value.RecordsCount, 2, "save-create-count should be: 2")
			mctest.AssertEquals(t, len(value.RecordIds), 2, "save-create-recordIds-length should be: 2")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should update two group-records and return success:",
		TestFunc: func() {
			res := crud.SaveRecord()
			fmt.Printf("updates: %v : %v \n", res.Message, res.ResCode)
			mctest.AssertEquals(t, res.Code, "success", "update should return code: success")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should update two group-records by Ids and return success:",
		TestFunc: func() {
			res := crud.SaveRecord()
			fmt.Printf("update-by-ids: %v : %v \n", res.Message, res.ResCode)
			mctest.AssertEquals(t, res.Code, "success", "update-by-id should return code: success")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should update two group-records by query-params and return success:",
		TestFunc: func() {
			res := crud.SaveRecord()
			fmt.Printf("update-by-params: %v : %v \n", res.Message, res.ResCode)
			mctest.AssertEquals(t, res.Code, "success", "update-by-params should return code: success")
		},
	})

	mctest.PostTestResult()

}
