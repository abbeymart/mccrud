// @Author: abbeymart | Abi Akindele | @Created: 2020-12-15 | @Updated: 2020-12-15
// @Company: mConnect.biz | @License: MIT
// @Description: go: mConnect

package tests

import "github.com/abbeymart/mccrud"

var PersonModel = mccrud.ModelType{
	TableName: "persons",
	RecordDesc: map[string]mccrud.FieldDescType{
		"id": {
			FieldType: mccrud.DataTypes().String,
			FieldLength: 100,
			FieldPattern: "",
			AllowNull: false,
			Unique: false,
			Indexable: false,
			PrimaryKey: false,
			ValidateMessage: "Length must not be longer than 100",
		},
		"name": {
			FieldType: mccrud.DataTypes().String,
			FieldLength: 100,
			FieldPattern: "",
			AllowNull: false,
			Unique: false,
			Indexable: false,
			PrimaryKey: false,
			ValidateMessage: "Length must not be longer than 100",
		},
	},
	Relations: nil,
	TimeStamp: true,
	ActorStamp: true,
	ActiveStamp: true,
}

var person = mccrud.NewModel(PersonModel)
