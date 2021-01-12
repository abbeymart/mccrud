// @Author: abbeymart | Abi Akindele | @Created: 2021-01-05 | @Updated: 2021-01-05
// @Company: mConnect.biz | @License: MIT
// @Description: crud local/shared types

package mccrud

import "github.com/abbeymart/mctypes"

type InsertedResultType struct {
	RecordIds   []string `json:"record_ids"`
	RecordCount int      `json:"record_count"`
}

type UpdatedResultType struct {
	QueryParam  mctypes.WhereParamType `json:"query_param"`
	RecordIds   []string               `json:"record_ids"`
	RecordCount int                    `json:"record_count"`
}

type DeletedResultType struct {
	QueryParam   mctypes.WhereParamType `json:"query_param"`
	RecordIds    []string               `json:"record_ids"`
	RecordCount  int                    `json:"record_count"`
	TableRecords []interface{}          `json:"table_records"`
}

type GetResultType struct {
	QueryParam   mctypes.WhereParamType `json:"query_param"`
	RecordIds    []string               `json:"record_ids"`
	RecordCount  int                    `json:"record_count"`
	TableRecords []interface{}          `json:"table_records""`
}

type LogRecordsType struct {
	TableFields  []string               `json:"table_fields"`
	TableRecords []interface{}          `json:"table_records"`
	QueryParam   mctypes.WhereParamType `json:"query_param"`
	RecordIds    []string               `json:"record_ids"`
}
