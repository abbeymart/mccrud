// @Author: abbeymart | Abi Akindele | @Created: 2020-12-28 | @Updated: 2020-12-28
// @Company: mConnect.biz | @License: MIT
// @Description: go: mConnect

package tests

import "encoding/json"

const TestTable = "audits_test1"
const DeleteAllTable = "audits_test2"

type TestParam struct {
	Name     string
	Desc     string
	Url      string
	Priority int
	Cost     float64
}

const UserId = "085f48c5-8763-4e22-a1c6-ac1a68ba07de"
var Recs = TestParam{Name: "Abi", Desc: "Testing only", Url: "localhost:9000", Priority: 1, Cost: 1000.00}
var TableRecords, _ = json.Marshal(Recs)
//fmt.Println("table-records-json", string(tableRecords))
var NewRecs = TestParam{Name: "Abi Akindele", Desc: "Testing only - updated", Url: "localhost:9900", Priority: 1, Cost: 2000.00}
var NewTableRecords, _ = json.Marshal(NewRecs)
//fmt.Println("new-table-records-json", string(newTableRecords))
var ReadP = map[string][]string{"keywords": {"lagos", "nigeria", "ghana", "accra"}}
var ReadParams, _ = json.Marshal(ReadP)

// create record(s)

// update record(s)

// get record(s)

// delete record(s)
