// @Author: abbeymart | Abi Akindele | @Created: 2020-12-24 | @Updated: 2020-12-24
// @Company: mConnect.biz | @License: MIT
// @Description: get/read records test cases

package mccrud

import (
	"encoding/json"
	"fmt"
	"github.com/abbeymart/mcdb"
	"github.com/abbeymart/mctest"
	"testing"
)

func TestGet(t *testing.T) {
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
	// defer dbClose
	defer myDb.ClosePgxDbPool()

	// check db-connection-error
	if err != nil {
		fmt.Printf("*****db-connection-error: %v\n", err.Error())
		return
	}

	getCrudParams := CrudParamsType{
		AppDb:       dbc.DbConn,
		TableName:   AuditTable,
		UserInfo:    TestUserInfo,
		RecordIds:   GetAuditByIds,
		QueryParams: GetAuditByParams,
	}
	//modelRefGroup := Group{}
	modelRefAudit := Audit{}
	recId := ""		// TODO: set id
	var getCrud = NewCrud(getCrudParams, CrudParamOptions)

	mctest.McTest(mctest.OptionValue{
		Name: "should get records by Id and return success:",
		TestFunc: func() {
			res := getCrud.GetById(modelRefAudit, recId)
			fmt.Printf("get-by-id-response: %#v\n\n", res)

			value, _ := res.Value.(CrudResultType)
			fmt.Printf("get-by-id-value: %#v\n", value.TableRecords)
			fmt.Printf("get-by-param-count: %v\n", value.RecordCount)
			jsonRecs, _ := json.Marshal(value.TableRecords)
			fmt.Printf("json-records: %v\n\n", string(jsonRecs))
			mctest.AssertEquals(t, res.Code, "success", "get-task should return code: success")
			mctest.AssertEquals(t, value.RecordCount, 2, "get-task-count should be: 2")
			mctest.AssertEquals(t, len(value.TableRecords), 2, "get-result-count should be: 2")
		},
	})

	mctest.McTest(mctest.OptionValue{
		Name: "should get records by Ids and return success:",
		TestFunc: func() {
			getCrud.RecordIds = GetAuditByIds
			getCrud.QueryParams = QueryParamType{}
			res := getCrud.GetByIds(modelRefAudit)
			fmt.Printf("get-by-id-response: %#v\n\n", res)

			value, _ := res.Value.(CrudResultType)
			fmt.Printf("get-by-id-value: %#v\n", value.TableRecords)
			fmt.Printf("get-by-param-count: %v\n", value.RecordCount)
			jsonRecs, _ := json.Marshal(value.TableRecords)
			fmt.Printf("json-records: %v\n\n", string(jsonRecs))
			mctest.AssertEquals(t, res.Code, "success", "get-task should return code: success")
			mctest.AssertEquals(t, value.RecordCount, 2, "get-task-count should be: 2")
			mctest.AssertEquals(t, len(value.TableRecords), 2, "get-result-count should be: 2")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should get records by query-params and return success:",
		TestFunc: func() {
			getCrud.RecordIds = []string{}
			getCrud.QueryParams = GetAuditByParams
			res := getCrud.GetByParam(modelRefAudit)
			//fmt.Printf("get-by-param-response: %#v\n", res)
			value, _ := res.Value.(CrudResultType)
			fmt.Printf("get-by-param-value: %#v\n", value.TableRecords)
			fmt.Printf("get-by-param-count: %v\n", value.RecordCount)
			mctest.AssertEquals(t, res.Code, "success", "get-task should return code: success")
			mctest.AssertEquals(t, value.RecordCount >= 0, true, "get-task-count should be >= 0")
			mctest.AssertEquals(t, len(value.TableRecords) >= 0, true, "get-result-count should be >= 0")
		},
	})

	mctest.McTest(mctest.OptionValue{
		Name: "should get all records and return success:",
		TestFunc: func() {
			getCrud.RecordIds = []string{}
			getCrud.QueryParams = QueryParamType{}
			res := getCrud.GetAll(modelRefAudit)
			value, _ := res.Value.(CrudResultType)
			fmt.Printf("get-by-all-value[0]: %#v\n", value.TableRecords[0])
			fmt.Printf("get-by-all-value[1]: %#v\n", value.TableRecords[1])
			fmt.Printf("get-by-all-count: %v\n", value.RecordCount)
			mctest.AssertEquals(t, res.Code, "success", "get-task should return code: success")
			mctest.AssertEquals(t, value.RecordCount >= 100, true, "get-task-count should be >= 10")
			mctest.AssertEquals(t, len(value.TableRecords) >= 100, true, "get-result-count should be >= 10")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should get all records by limit/skip(offset) and return success:",
		TestFunc: func() {
			getCrud.TableName = AuditTable
			modelRef := Audit{}
			getCrud.Skip = 0
			getCrud.Limit = 20
			res := getCrud.GetAll(modelRef)
			value, _ := res.Value.(CrudResultType)
			fmt.Printf("get-by-all-value[0]: %#v\n", value.TableRecords[0])
			fmt.Printf("get-by-all-value[1]: %#v\n", value.TableRecords[1])
			fmt.Printf("get-by-all-limit-count: %v\n", value.RecordCount)
			mctest.AssertEquals(t, res.Code, "success", "get-task should return code: success")
			mctest.AssertEquals(t, value.RecordCount == 20, true, "get-task-count should be = 20")
			mctest.AssertEquals(t, len(value.TableRecords) == 20, true, "get-result-count should be = 20")
		},
	})
	mctest.PostTestResult()

}
