// @Author: abbeymart | Abi Akindele | @Created: 2020-12-01 | @Updated: 2020-12-01
// @Company: mConnect.biz | @License: MIT
// @Description: get / query record(s)

package mccrud

import (
	"context"
	"fmt"
	"github.com/abbeymart/mccrud/helper"
	"github.com/abbeymart/mcresponse"
)

func (crud Crud) GetById(tableFieldPointers ...*interface{}) mcresponse.ResponseMessage {
	var tableFields []string
	// TODO: compose tableFields | ignore error, if tablesFields-len==0,
	// TODO: SELECT/scan fields in order specified by the table-model/type definition

	if tFields, err := helper.ComputeGetFields(crud.ProjectParams); err != nil {
		return mcresponse.GetResMessage("getError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing get/read-query-fields: %v", err.Error()),
			Value:   nil,
		})
	} else {
		tableFields = tFields
	}
	if queryRes, err := helper.ComputeSelectQueryById(crud.TableName, tableFields, crud.RecordIds); err != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing select/read-query: %v", err.Error()),
			Value:   queryRes,
		})
	} else {
		// TODO: perform crud-task action, include options (skip, limit, sort etc.):

		// by tableFields
		if getQuery, err := helper.ComputeSelectQueryById(crud.TableName, []string{}, crud.RecordIds); err != nil {
			return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Error computing select/read-query: %v", err.Error()),
				Value:   getQuery,
			})
		} else {
			// exit of currentRec-length is less than recordIds-length
			rows, err := crud.AppDb.Query(context.Background(), getQuery)
			if err != nil {
				errMsg := fmt.Sprintf("Db query Error: %v", err.Error())
				return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
					Message: errMsg,
					Value:   nil,
				})
			}
			defer rows.Close()
			// check rows count
			var rowCount = 0
			for rows.Next() {
				//var id string
				if err := rows.Scan(tableFieldPointers); err == nil {
					rowCount += 1
					// crud.CurrentRecords = append(crud.CurrentRecords, id)
					// parse the current-records for audit-log
					if parseVal, err := helper.ParseRawValues(rows.RawValues()); err != nil {
						return mcresponse.GetResMessage("parseError", mcresponse.ResponseMessageOptions{
							Message: fmt.Sprintf("Error parsing raw-record-values: %v", err.Error()),
							Value:   nil,
						})
					} else {
						// update instance CurrentRecords
						crud.CurrentRecords = append(crud.CurrentRecords, parseVal)
					}
				}
			}
			if rowCount < len(crud.RecordIds) {
				return mcresponse.GetResMessage("fewRecords", mcresponse.ResponseMessageOptions{
					Message: fmt.Sprintf("Fewer records (%v) less than expected (%v)", rowCount, len(crud.RecordIds)),
					Value:   nil,
				})
			}
			if err := rows.Err(); err != nil {
				return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
					Message: fmt.Sprintf("Error reading/getting records: %v", err.Error()),
					Value:   nil,
				})
			}
		}


		// by len(tableField) == 0 (i.e. all fields / *) => tableFields == TableModel-params

	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "success",
		Value:   nil,
	})
}

func (crud Crud) GetByParam() mcresponse.ResponseMessage {
	var tableFields []string
	// compose tableFields
	if tFields, err := helper.ComputeGetFields(crud.ProjectParams); err != nil {
		return mcresponse.GetResMessage("getError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing get/read-query-fields: %v", err.Error()),
			Value:   nil,
		})
	} else {
		tableFields = tFields
	}
	if queryRes, err := helper.ComputeSelectQueryByParam(crud.TableName, crud.QueryParams, tableFields); err != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing select/read-query: %v", err.Error()),
			Value:   queryRes,
		})
	} else {
		// TODO: perform crud-task action, include options (skip, limit, sort etc.):
		// by tableFields


		// by len(tableField) == 0 (i.e. all fields / *) => tableFields == TableModel-params


	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "success",
		Value:   nil,
	})
}

func (crud Crud) GetAll() mcresponse.ResponseMessage {
	var tableFields []string
	// compose tableFields
	if tFields, err := helper.ComputeGetFields(crud.ProjectParams); err != nil {
		return mcresponse.GetResMessage("getError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing get/read-query-fields: %v", err.Error()),
			Value:   nil,
		})
	} else {
		tableFields = tFields
	}
	if queryRes, err := helper.ComputeSelectQueryAll(crud.TableName, tableFields); err != nil {
		return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing select/read-query: %v", err.Error()),
			Value:   queryRes,
		})
	} else {
		// TODO: perform crud-task action, include options (skip, limit, sort etc.):

	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "success",
		Value:   nil,
	})
}
