// @Author: abbeymart | Abi Akindele | @Created: 2020-12-14 | @Updated: 2020-12-14
// @Company: mConnect.biz | @License: MIT
// @Description: go: mConnect

package mccrud

import (
	"fmt"
	"github.com/abbeymart/mcauditlog"
	"github.com/abbeymart/mcdb"
	"github.com/abbeymart/mctest"
	"testing"
)

func TestSaveGroup(t *testing.T) {
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
	mcLogResult := mcauditlog.PgxLogParam{AuditDb: dbc.DbConn, AuditTable: AuditTable}
	// audit-log instance
	mcLog := mcauditlog.NewAuditLogPgx(dbc.DbConn, AuditTable)

	// group-table-records
	groupModelRef := Group{}
	createCrudParams := CrudParamsType{
		AppDb:        dbc.DbConn,
		TableName:    GroupTable,
		UserInfo:     TestUserInfo,
		ActionParams: GroupCreateActionParams,
	}
	updateCrudParams := CrudParamsType{
		AppDb:        dbc.DbConn,
		TableName:    GroupTable,
		UserInfo:     TestUserInfo,
		ActionParams: GroupUpdateActionParams,
	}
	updateCrudParamsById := CrudParamsType{
		AppDb:        dbc.DbConn,
		TableName:    GroupTable,
		UserInfo:     TestUserInfo,
		ActionParams: ActionParamsType{GroupUpdateRecordById},
		RecordIds:    GroupByIds,
	}
	updateCrudParamsByParam := CrudParamsType{
		AppDb:        dbc.DbConn,
		TableName:    GroupTable,
		UserInfo:     TestUserInfo,
		ActionParams: ActionParamsType{GroupUpdateRecordByParam},
		QueryParams:  GroupByParams,
	}

	//fmt.Printf("test-action-params: %#v \n", createCrudParams.ActionParams)

	var crud interface{} = NewCrud(createCrudParams, CrudParamOptions)
	var updateCrud = NewCrud(updateCrudParams, CrudParamOptions)
	var updateIdCrud = NewCrud(updateCrudParamsById, CrudParamOptions)
	var updateParamCrud = NewCrud(updateCrudParamsByParam, CrudParamOptions)

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
			_, ok := crud.(*Crud)
			mctest.AssertEquals(t, ok, true, "crud should be instance of mccrud.Crud")
		},
	})

	// group-table test-cases
	mctest.McTest(mctest.OptionValue{
		Name: "should create two new group-records and return success:",
		TestFunc: func() {
			crud, ok := crud.(*Crud)
			if !ok {
				mctest.AssertEquals(t, ok, true, "crud should be instance of mccrud.Crud")
			}
			res := crud.SaveRecord(groupModelRef)
			fmt.Println(res.Message, res.ResCode)
			value, _ := res.Value.(CrudResultType)
			mctest.AssertEquals(t, res.Code, "success", "save-create should return code: success")
			mctest.AssertEquals(t, value.RecordCount, 2, "save-create-count should be: 2")
			mctest.AssertEquals(t, len(value.RecordIds), 2, "save-create-recordIds-length should be: 2")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should update two group-records and return success:",
		TestFunc: func() {
			res := updateCrud.SaveRecord(groupModelRef)
			fmt.Printf("updates: %v : %v \n", res.Message, res.ResCode)
			mctest.AssertEquals(t, res.Code, "success", "update should return code: success")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should update two group-records by Ids and return success:",
		TestFunc: func() {
			res := updateIdCrud.SaveRecord(groupModelRef)
			fmt.Printf("update-by-ids: %v : %v \n", res.Message, res.ResCode)
			mctest.AssertEquals(t, res.Code, "success", "update-by-id should return code: success")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should update two group-records by query-params and return success:",
		TestFunc: func() {
			res := updateParamCrud.SaveRecord(groupModelRef)
			fmt.Printf("update-by-params: %v : %v \n", res.Message, res.ResCode)
			mctest.AssertEquals(t, res.Code, "success", "update-by-params should return code: success")
		},
	})

	mctest.PostTestResult()

}
