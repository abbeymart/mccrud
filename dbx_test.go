// @Author: abbeymart | Abi Akindele | @Created: 2020-12-04 | @Updated: 2020-12-04
// @Company: mConnect.biz | @License: MIT
// @Description: db testing

package mccrud

import (
	"fmt"
	"testing"
)
import "github.com/abbeymart/mctest"


func TestDbx(t *testing.T) {
	// test-data: db-configuration settings
	myDb := MyDb
	myDb.Options = DbConnectOptions{}

	sqliteDb := DbConfig{
		DbType: "sqlite3",
		Filename: "testdb.db",
	}

	mctest.McTest(mctest.OptionValue{
		Name: "should successfully connect to the PostgresDB",
		TestFunc: func() {
			dbc, err := myDb.OpenDbx()
			fmt.Printf("pg-dbc: %v\n", dbc)
			defer myDb.CloseDbx()
			fmt.Println(dbc)
			fmt.Println("*****************************************")
			mctest.AssertEquals(t, err, nil, "response-code should be: nil")
		},
	})

	mctest.McTest(mctest.OptionValue{
		Name: "should successfully connect to SQLite3 database",
		TestFunc: func() {
			dbc2, err := sqliteDb.OpenDbx()
			defer sqliteDb.CloseDbx()
			fmt.Println(dbc2)
			mctest.AssertEquals(t, err, nil, "response-code should be: nil")
		},
	})

	mctest.PostTestResult()
}
