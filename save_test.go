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

	mctest.McTest(mctest.OptionValue{
		Name: "should create two new records and return success:",
		TestFunc: func() {
			crud.TableName = AuditTable
			crud.ActionParams = AuditCreateActionParams
			crud.RecordIds = []string{}
			crud.QueryParams = QueryParamType{}
			recLen := len(crud.ActionParams)
			res := crud.SaveRecord()
			fmt.Println(res.Message, res.ResCode)
			value, _ := res.Value.(CrudResultType)
			mctest.AssertEquals(t, res.Code, "success", "create should return code: success")
			mctest.AssertEquals(t, value.RecordsCount, recLen, fmt.Sprintf("save-create-count should be: %v", recLen))
			mctest.AssertEquals(t, len(value.RecordIds), recLen, fmt.Sprintf("save-create-recordIds-length should be: %v", recLen))
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should update two existing records and return success:",
		TestFunc: func() {
			crud.TableName = UpdateTable
			crud.ActionParams = AuditUpdateActionParams
			crud.RecordIds = []string{}
			crud.QueryParams = QueryParamType{}
			recLen := len(crud.ActionParams)
			res := crud.SaveRecord()
			fmt.Printf("updates: %v : %v \n", res.Message, res.ResCode)
			value, _ := res.Value.(CrudResultType)
			mctest.AssertEquals(t, res.Code, "success", "update should return code: success")
			mctest.AssertEquals(t, value.RecordsCount, recLen, fmt.Sprintf("save-create-count should be: %v", recLen))
			mctest.AssertEquals(t, len(value.RecordIds), recLen, fmt.Sprintf("save-create-recordIds-length should be: %v", recLen))
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should update a record by Id and return success:",
		TestFunc: func() {
			crud.TableName = UpdateTable
			crud.ActionParams = ActionParamsType{AuditUpdateRecordById}
			crud.RecordIds = []string{UpdateAuditById}
			crud.QueryParams = QueryParamType{}
			idLen := len(crud.RecordIds)
			res := crud.SaveRecord()
			fmt.Printf("update-by-ids: %v : %v \n", res.Message, res.ResCode)
			value, _ := res.Value.(CrudResultType)
			mctest.AssertEquals(t, res.Code, "success", "update-by-id should return code: success")
			mctest.AssertEquals(t, value.RecordsCount, idLen, fmt.Sprintf("update-by-id-count should be: %v", idLen))
			mctest.AssertEquals(t, len(value.RecordIds), idLen, fmt.Sprintf("update-by-id-recordIds-length should be: %v", idLen))
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should update records by Ids and return success:",
		TestFunc: func() {
			crud.TableName = UpdateTable
			crud.ActionParams = ActionParamsType{AuditUpdateRecordByIds}
			crud.RecordIds = UpdateAuditByIds
			crud.QueryParams = QueryParamType{}
			idLen := len(crud.RecordIds)
			res := crud.SaveRecord()
			fmt.Printf("update-by-ids: %v : %v \n", res.Message, res.ResCode)
			mctest.AssertEquals(t, res.Code, "success", "update-by-id should return code: success")
			value, _ := res.Value.(CrudResultType)
			mctest.AssertEquals(t, res.Code, "success", "update-by-ids should return code: success")
			mctest.AssertEquals(t, value.RecordsCount, idLen, fmt.Sprintf("update-by-ids-count should be: %v", idLen))
			mctest.AssertEquals(t, len(value.RecordIds), idLen, fmt.Sprintf("update-by-ids-recordIds-length should be: %v", idLen))
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should update records by query-params and return success:",
		TestFunc: func() {
			crud.TableName = UpdateTable
			crud.ActionParams = ActionParamsType{AuditUpdateRecordByParam}
			crud.RecordIds = []string{}
			crud.QueryParams = UpdateAuditByParams
			recLen := len(crud.ActionParams)
			res := crud.SaveRecord()
			fmt.Printf("update-by-params: %v : %v \n", res.Message, res.ResCode)
			value, _ := res.Value.(CrudResultType)
			mctest.AssertEquals(t, res.Code, "success", "update-by-params should return code: success")
			mctest.AssertEquals(t, value.RecordsCount > recLen, true, fmt.Sprintf("update-by-params-count should be >: %v", recLen))
			mctest.AssertEquals(t, len(value.RecordIds) > recLen, true, fmt.Sprintf("update-by-params-recordIds-length should be > : %v", recLen))
		},
	})

	mctest.PostTestResult()

}
