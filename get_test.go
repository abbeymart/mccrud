// @Author: abbeymart | Abi Akindele | @Created: 2020-12-24 | @Updated: 2020-12-24
// @Company: mConnect.biz | @License: MIT
// @Description: get/read records test cases

package mccrud

import (
	"encoding/json"
	"fmt"
	"github.com/abbeymart/mccrud/test/config/secure"
	"github.com/abbeymart/mcdb"
	"github.com/abbeymart/mctest"
	"testing"
)

func TestGet(t *testing.T) {
	myDb := secure.MyDb
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
	modelRef := Audit{}
	crudParams := CrudParamsType{
		AppDb:       dbc.DbConn,
		ModelRef:    modelRef,
		TableName:   AuditTable,
		UserInfo:    TestUserInfo,
		RecordIds:   []string{},
		QueryParams: QueryParamType{},
	}
	var crud = NewCrud(crudParams, CrudParamOptions)

	mctest.McTest(mctest.OptionValue{
		Name: "should get records by Id and return success:",
		TestFunc: func() {
			crud.RecordIds = []string{GetAuditById}
			res := crud.GetRecord()
			fmt.Printf("get-by-id-response: %#v\n\n", res)
			value, _ := res.Value.(GetResultType)
			fmt.Printf("get-by-id-value: %#v\n", value.Records)
			fmt.Printf("get-by-param-count: %v\n", value.Stats.RecordsCount)
			jsonRecs, _ := json.Marshal(value.Records)
			fmt.Printf("json-records: %v\n\n", string(jsonRecs))
			mctest.AssertEquals(t, res.Code, "success", "get-task should return code: success")
			mctest.AssertEquals(t, value.Stats.RecordsCount, 1, "get-task-count should be: 1")
			mctest.AssertEquals(t, len(value.Records), 1, "get-result-count should be: 1")
		},
	})

	mctest.McTest(mctest.OptionValue{
		Name: "should get records by Ids and return success:",
		TestFunc: func() {
			crud.RecordIds = GetAuditByIds
			crud.QueryParams = QueryParamType{}
			res := crud.GetByIds()
			fmt.Printf("get-by-id-response: %#v\n\n", res)
			value, _ := res.Value.(GetResultType)
			fmt.Printf("get-by-id-value: %#v\n", value.Records)
			fmt.Printf("get-by-param-count: %v\n", value.Stats.RecordsCount)
			jsonRecs, _ := json.Marshal(value.Records)
			fmt.Printf("json-records: %v\n\n", string(jsonRecs))
			mctest.AssertEquals(t, res.Code, "success", "get-task should return code: success")
			mctest.AssertEquals(t, value.Stats.RecordsCount, 2, "get-task-count should be: 2")
			mctest.AssertEquals(t, len(value.Records), 2, "get-result-count should be: 2")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should get records by query-params and return success:",
		TestFunc: func() {
			crud.RecordIds = []string{}
			crud.QueryParams = GetAuditByParams
			res := crud.GetByParam()
			//fmt.Printf("get-by-param-response: %#v\n", res)
			value, _ := res.Value.(GetResultType)
			fmt.Printf("get-by-param-value: %#v\n", value.Records)
			fmt.Printf("get-by-param-count: %v\n", value.Stats.RecordsCount)
			mctest.AssertEquals(t, res.Code, "success", "get-task should return code: success")
			mctest.AssertEquals(t, value.Stats.RecordsCount > 0, true, "get-task-count should be >= 0")
			mctest.AssertEquals(t, len(value.Records) > 0, true, "get-result-count should be >= 0")
		},
	})

	mctest.McTest(mctest.OptionValue{
		Name: "should get all records and return success:",
		TestFunc: func() {
			crud.RecordIds = []string{}
			crud.QueryParams = QueryParamType{}
			res := crud.GetAll()
			value, _ := res.Value.(CrudResultType)
			fmt.Printf("get-by-all-value[0]: %#v\n", value.Records[0])
			fmt.Printf("get-by-all-value[1]: %#v\n", value.Records[1])
			fmt.Printf("get-by-all-count: %v\n", value.RecordsCount)
			mctest.AssertEquals(t, res.Code, "success", "get-task should return code: success")
			mctest.AssertEquals(t, value.RecordsCount > 20, true, "get-task-count should be >= 10")
			mctest.AssertEquals(t, len(value.Records) > 20, true, "get-result-count should be >= 10")
		},
	})
	mctest.McTest(mctest.OptionValue{
		Name: "should get all records by limit/skip(offset) and return success:",
		TestFunc: func() {
			crud.RecordIds = []string{}
			crud.QueryParams = QueryParamType{}
			crud.Skip = 0
			crud.Limit = 20
			res := crud.GetAll()
			value, _ := res.Value.(GetResultType)
			fmt.Printf("get-by-all-value[0]: %#v\n", value.Records[0])
			fmt.Printf("get-by-all-value[1]: %#v\n", value.Records[1])
			fmt.Printf("get-by-all-limit-count: %v\n", value.Stats.RecordsCount)
			mctest.AssertEquals(t, res.Code, "success", "get-task should return code: success")
			mctest.AssertEquals(t, value.Stats.RecordsCount == 20, true, "get-task-count should be = 20")
			mctest.AssertEquals(t, len(value.Records) == 20, true, "get-result-count should be = 20")
		},
	})
	mctest.PostTestResult()

}
