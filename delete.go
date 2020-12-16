// @Author: abbeymart | Abi Akindele | @Created: 2020-12-01 | @Updated: 2020-12-01
// @Company: mConnect.biz | @License: MIT
// @Description: delete or remove record(s)

package mccrud

import (
	"context"
	"fmt"
	"github.com/abbeymart/mccrud/helper"
	"github.com/abbeymart/mcresponse"
)

func (crud Crud) DeleteById() mcresponse.ResponseMessage {
	// TODO: get current records, for audit-log | for delete tableField = []string{}
	if crud.LogDelete {
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
				var id string
				if err := rows.Scan(&id); err == nil {
					rowCount += 1
					// crud.CurrentRecords = append(crud.CurrentRecords, id)
					// parse the current-records for audit-log
					if parseVal, err := helper.ParseRawValues(rows.RawValues()); err != nil{
						return mcresponse.GetResMessage("parseError", mcresponse.ResponseMessageOptions{
							Message: fmt.Sprintf("Error parsing raw-record-values: %v", err.Error()),
							Value:   nil,
						})
					} else {
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
	}

	// TODO: perform crud-task action, include where-
	// compute delete script from where

	// perform audit-log


	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "success",
		Value:   nil,
	})
}

func (crud Crud) DeleteByParam() mcresponse.ResponseMessage {
	// TODO: get current records, for audit-log | for delete tableField = []string{}
	if crud.LogDelete {
		if getQuery, err := helper.ComputeSelectQueryByParam(crud.TableName, []string{}, crud.QueryParams); err != nil {
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
				var id string
				if err := rows.Scan(&id); err == nil {
					rowCount += 1
					// crud.CurrentRecords = append(crud.CurrentRecords, id)
					// parse the current-records for audit-log
					if parseVal, err := helper.ParseRawValues(rows.RawValues()); err != nil{
						return mcresponse.GetResMessage("parseError", mcresponse.ResponseMessageOptions{
							Message: fmt.Sprintf("Error parsing raw-record-values: %v", err.Error()),
							Value:   nil,
						})
					} else {
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
	}

	// TODO: perform crud-task action, include where-query(params):
	// compute delete script from where

	// perform audit-log

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "success",
		Value:   nil,
	})
}

// DeleteAll function removes all records in the tables. Recommended for admin-users only
// Use if and only if you know what you are doing
func (crud Crud) DeleteAll() mcresponse.ResponseMessage {
	// TODO: perform crud-task action:
	// compute delete script from where

	// perform audit-log | record only the action performed, excluding table-records

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "success",
		Value:   nil,
	})
}