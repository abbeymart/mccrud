// @Author: abbeymart | Abi Akindele | @Created: 2020-12-14 | @Updated: 2020-12-14
// @Company: mConnect.biz | @License: MIT
// @Description: go: mConnect

package tests

import (
	"encoding/json"
	"errors"
	"fmt"
)

// data
var jsonData = `{"name": "Abi", "age": 10, "location_id": "CA", "phone_number": "123-456-9999"}`

type Person struct {
	Name        string `json:"name"`
	Age         int    `json:"age"`
	LocationId  string `json:"location_id"`
	PhoneNumber string `json:"phone_number"`
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
}
