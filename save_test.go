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
	updateCrudParams := mctypes.CrudParamsType{
		AppDb:        dbc.DbConn,
		TableName:    TestTable,
		UserInfo:     TestUserInfo,
		ActionParams: UpdateActionParams,
	}
	deleteCrudParams := mctypes.CrudParamsType{
		AppDb:        dbc.DbConn,
		TableName:    TestTable,
		UserInfo:     TestUserInfo,
	}
	deleteAllCrudParams := mctypes.CrudParamsType{
		AppDb:        dbc.DbConn,
		TableName:    DeleteAllTable,
		UserInfo:     TestUserInfo,
	}
	getCrudParams := mctypes.CrudParamsType{
		AppDb:        dbc.DbConn,
		TableName:    TestTable,
		UserInfo:     TestUserInfo,
	}

	//fmt.Printf("test-action-params: %#v \n", createCrudParams.ActionParams)

	var crud interface{} = NewCrud(createCrudParams, TestCrudParamOptions)
	var updateCrud = NewCrud(updateCrudParams, TestCrudParamOptions)
	var getCrud = NewCrud(getCrudParams, TestCrudParamOptions)
	var deleteCrud = NewCrud(deleteCrudParams, TestCrudParamOptions)
	var deleteAllCrud = NewCrud(deleteAllCrudParams, TestCrudParamOptions)

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

	mctest.McTest(mctest.OptionValue{
		Name: "should create two new records and return success:",
		TestFunc: func() {
			crud, ok := crud.(*Crud)
			if !ok {
				mctest.AssertEquals(t, ok, true, "crud should be instance of mccrud.Crud")
			}
			res := crud.Save(CreateTableFields)
			fmt.Println(res.Message, res.ResCode)
			value, _ := res.Value.(InsertedResultType)
			mctest.AssertEquals(t, res.Code, "success", "save-create should return code: success")
			mctest.AssertEquals(t, value.TableName, TestTable, "save-create-table should be: " + TestTable)
			mctest.AssertEquals(t, value.RecordCount, 2, "save-create-count should be: 2")
			mctest.AssertEquals(t, len(value.RecordIds), 2, "save-create-recordIds-length should be: 2")
		},
	})

	mctest.McTest(mctest.OptionValue{
		Name: "should update two records and return success:",
		TestFunc: func() {
			crud, ok := crud.(*Crud)
			if !ok {
				mctest.AssertEquals(t, ok, true, "crud should be instance of mccrud.Crud")
			}
			res := crud.Save(UpdateTableFields)
			//fmt.Println(res.Message, res.ResCode)
			value, _ := res.Value.(InsertedResultType)
			mctest.AssertEquals(t, res.Code, "success", "save-create should return code: success")
			mctest.AssertEquals(t, value.TableName, TestTable, "save-create-table should be: " + TestTable)
			mctest.AssertEquals(t, value.RecordCount, 2, "save-create-count should be: 2")
			mctest.AssertEquals(t, len(value.RecordIds), 2, "save-create-recordIds-length should be: 2")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should update two records by Ids and return success:",
		TestFunc: func() {
			crud, ok := crud.(*Crud)
			if !ok {
				mctest.AssertEquals(t, ok, true, "crud should be instance of mccrud.Crud")
			}
			res := crud.Save(CreateTableFields)
			fmt.Println(res.Message, res.ResCode)
			value, _ := res.Value.(InsertedResultType)
			mctest.AssertEquals(t, res.Code, "success", "save-create should return code: success")
			mctest.AssertEquals(t, value.TableName, TestTable, "save-create-table should be: " + TestTable)
			mctest.AssertEquals(t, value.RecordCount, 2, "save-create-count should be: 2")
			mctest.AssertEquals(t, len(value.RecordIds), 2, "save-create-recordIds-length should be: 2")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should update two records by query-params and return success:",
		TestFunc: func() {
			crud, ok := crud.(*Crud)
			if !ok {
				mctest.AssertEquals(t, ok, true, "crud should be instance of mccrud.Crud")
			}
			res := crud.Save(CreateTableFields)
			fmt.Println(res.Message, res.ResCode)
			value, _ := res.Value.(InsertedResultType)
			mctest.AssertEquals(t, res.Code, "success", "save-create should return code: success")
			mctest.AssertEquals(t, value.TableName, TestTable, "save-create-table should be: " + TestTable)
			mctest.AssertEquals(t, value.RecordCount, 2, "save-create-count should be: 2")
			mctest.AssertEquals(t, len(value.RecordIds), 2, "save-create-recordIds-length should be: 2")
		},
	})

	mctest.McTest(mctest.OptionValue{
		Name: "should get records by Ids and return success:",
		TestFunc: func() {
			crud, ok := crud.(*Crud)
			if !ok {
				mctest.AssertEquals(t, ok, true, "crud should be instance of mccrud.Crud")
			}
			res := crud.Save(CreateTableFields)
			fmt.Println(res.Message, res.ResCode)
			value, _ := res.Value.(InsertedResultType)
			mctest.AssertEquals(t, res.Code, "success", "save-create should return code: success")
			mctest.AssertEquals(t, value.TableName, TestTable, "save-create-table should be: " + TestTable)
			mctest.AssertEquals(t, value.RecordCount, 2, "save-create-count should be: 2")
			mctest.AssertEquals(t, len(value.RecordIds), 2, "save-create-recordIds-length should be: 2")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should get records by query-params and return success:",
		TestFunc: func() {
			crud, ok := crud.(*Crud)
			if !ok {
				mctest.AssertEquals(t, ok, true, "crud should be instance of mccrud.Crud")
			}
			res := crud.Save(CreateTableFields)
			fmt.Println(res.Message, res.ResCode)
			value, _ := res.Value.(InsertedResultType)
			mctest.AssertEquals(t, res.Code, "success", "save-create should return code: success")
			mctest.AssertEquals(t, value.TableName, TestTable, "save-create-table should be: " + TestTable)
			mctest.AssertEquals(t, value.RecordCount, 2, "save-create-count should be: 2")
			mctest.AssertEquals(t, len(value.RecordIds), 2, "save-create-recordIds-length should be: 2")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should get all records and return success:",
		TestFunc: func() {
			crud, ok := crud.(*Crud)
			if !ok {
				mctest.AssertEquals(t, ok, true, "crud should be instance of mccrud.Crud")
			}
			res := crud.Save(CreateTableFields)
			fmt.Println(res.Message, res.ResCode)
			value, _ := res.Value.(InsertedResultType)
			mctest.AssertEquals(t, res.Code, "success", "save-create should return code: success")
			mctest.AssertEquals(t, value.TableName, TestTable, "save-create-table should be: " + TestTable)
			mctest.AssertEquals(t, value.RecordCount, 2, "save-create-count should be: 2")
			mctest.AssertEquals(t, len(value.RecordIds), 2, "save-create-recordIds-length should be: 2")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should get all records by limit/skip(offset) and return success:",
		TestFunc: func() {
			crud, ok := crud.(*Crud)
			if !ok {
				mctest.AssertEquals(t, ok, true, "crud should be instance of mccrud.Crud")
			}
			res := crud.Save(CreateTableFields)
			fmt.Println(res.Message, res.ResCode)
			value, _ := res.Value.(InsertedResultType)
			mctest.AssertEquals(t, res.Code, "success", "save-create should return code: success")
			mctest.AssertEquals(t, value.TableName, TestTable, "save-create-table should be: " + TestTable)
			mctest.AssertEquals(t, value.RecordCount, 2, "save-create-count should be: 2")
			mctest.AssertEquals(t, len(value.RecordIds), 2, "save-create-recordIds-length should be: 2")
		},
	})

	mctest.McTest(mctest.OptionValue{
		Name: "should delete two records by Ids and return success:",
		TestFunc: func() {
			crud, ok := crud.(*Crud)
			if !ok {
				mctest.AssertEquals(t, ok, true, "crud should be instance of mccrud.Crud")
			}
			res := crud.Save(CreateTableFields)
			fmt.Println(res.Message, res.ResCode)
			value, _ := res.Value.(InsertedResultType)
			mctest.AssertEquals(t, res.Code, "success", "save-create should return code: success")
			mctest.AssertEquals(t, value.TableName, TestTable, "save-create-table should be: " + TestTable)
			mctest.AssertEquals(t, value.RecordCount, 2, "save-create-count should be: 2")
			mctest.AssertEquals(t, len(value.RecordIds), 2, "save-create-recordIds-length should be: 2")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should delete two records by query-params and return success:",
		TestFunc: func() {
			crud, ok := crud.(*Crud)
			if !ok {
				mctest.AssertEquals(t, ok, true, "crud should be instance of mccrud.Crud")
			}
			res := crud.Save(CreateTableFields)
			fmt.Println(res.Message, res.ResCode)
			value, _ := res.Value.(InsertedResultType)
			mctest.AssertEquals(t, res.Code, "success", "save-create should return code: success")
			mctest.AssertEquals(t, value.TableName, TestTable, "save-create-table should be: " + TestTable)
			mctest.AssertEquals(t, value.RecordCount, 2, "save-create-count should be: 2")
			mctest.AssertEquals(t, len(value.RecordIds), 2, "save-create-recordIds-length should be: 2")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should delete all table records and return success:",
		TestFunc: func() {
			crud, ok := crud.(*Crud)
			// change crud-table-name
			if !ok {
				mctest.AssertEquals(t, ok, true, "crud should be instance of mccrud.Crud")
			}
			res := crud.Save(CreateTableFields)
			fmt.Println(res.Message, res.ResCode)
			value, _ := res.Value.(InsertedResultType)
			mctest.AssertEquals(t, res.Code, "success", "save-create should return code: success")
			mctest.AssertEquals(t, value.TableName, TestTable, "save-create-table should be: " + TestTable)
			mctest.AssertEquals(t, value.RecordCount, 2, "save-create-count should be: 2")
			mctest.AssertEquals(t, len(value.RecordIds), 2, "save-create-recordIds-length should be: 2")
		},
	})

	mctest.PostTestResult()

}
