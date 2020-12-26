// @Author: abbeymart | Abi Akindele | @Created: 2020-12-14 | @Updated: 2020-12-14
// @Company: mConnect.biz | @License: MIT
// @Description: go: mConnect

package tests

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/abbeymart/mctypes"
	"github.com/abbeymart/mctypes/groupOperators"
)

// data
var jsonData = `[
	{"name": "Abi", "age": 10, "location_id": "CA", "phone_number": "123-456-9999"},
	{"name": "Abi", "age": 10, "location_id": "CA", "phone_number": "123-456-9999"},
	{"name": "Abi", "age": 10, "location_id": "CA", "phone_number": "123-456-9999"}
]`

var jsonQueryParams = `
	{
		"group_name": "abc",
		"group_order": 1,
		"group_link_op": "and",
		"group_items": [
			group_item: {"name": { "eq": Paul"}, "age":{"eq": 10}, "location": {"eq": Toronto"}},
			
		]
	},
	{},
	{},
`
// convert/decode jsonQueryParams to queryParams
var queryParams mctypes.QueryParamType = mctypes.QueryParamType{
	mctypes.QueryGroupType{
		GroupName: "abc",
		GroupOrder: 1,
		GroupLinkOp: groupOperators.AND,
		GroupItems: []mctypes.QueryItemType{
			{},
			{},
		},
	},
	mctypes.QueryGroupType{},
	mctypes.QueryGroupType{},
}

type Person struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Age         int    `json:"age"`
	LocationId  string `json:"location_id"`
	PhoneNumber string `json:"phone_number"`
}

// convert/decode jsonData to []model-type => action-params
var actionParams = mctypes.ActionParamsType{
	mctypes.ValueParamType{},
	mctypes.ValueParamType{},
	mctypes.ValueParamType{},
	mctypes.ValueParamType{},
}

func jsonDataETL(personJson []byte) (Person, error) {
	var person Person
	if err := json.Unmarshal(personJson, &person); err == nil {
		return person, nil
	} else {
		return Person{}, errors.New(fmt.Sprintf("Error converting json-to-struct: %v", err.Error()))
	}
}

func main() {
	if person, err := jsonDataETL([]byte(jsonData)); err == nil {
		fmt.Printf("Person's record: %+v", person)
	} else {
		fmt.Printf("Error coverting json-data: %v", err.Error())
	}

	// TODO: table-fields | tableFieldPointers | queryParams

}
