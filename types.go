// @Author: abbeymart | Abi Akindele | @Created: 2021-01-05 | @Updated: 2021-01-05
// @Company: mConnect.biz | @License: MIT
// @Description: crud local/shared types

package mccrud

import "github.com/abbeymart/mctypes"

type InsertedResultType struct {
	TableName   string   `json:"table_name"`
	RecordIds   []string `json:"record_ids"`
	RecordCount int      `json:"record_count"`
}

type UpdatedResultType struct {
	TableName   string                 `json:"table_name"`
	QueryParam  mctypes.WhereParamType `json:"query_param"`
	RecordIds   []string               `json:"record_ids"`
	RecordCount int                    `json:"record_count"`
}

type DeletedResultType struct {
	TableName   string                 `json:"table_name"`
	QueryParam  mctypes.WhereParamType `json:"query_param"`
	RecordIds   []string               `json:"record_ids"`
	RecordCount int                    `json:"record_count"`
}

type GetResultType struct {
	TableName    string                 `json:"table_name"`
	QueryParam   mctypes.WhereParamType `json:"query_param"`
	RecordIds    []string               `json:"record_ids"`
	RecordCount  int                    `json:"record_count"`
	RecordValues []interface{}          `json:"record_values"`
}

type LogRecordsType struct {
	TableFields  []string
	TableRecords []interface{}
}
