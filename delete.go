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

// DeleteById method deletes or removes record(s) by record-id(s)
func (crud *Crud) DeleteById() mcresponse.ResponseMessage {
	// perform crud-task action, include where-
	// compute delete query by record-ids
	deleteQuery, dQErr := helper.ComputeDeleteQueryById(crud.TableName, crud.RecordIds)
	if dQErr != nil {
		return mcresponse.GetResMessage("deleteError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing delete-query: %v", dQErr.Error()),
			Value:   nil,
		})
	}
	commandTag, delErr := crud.AppDb.Exec(context.Background(), deleteQuery)
	if delErr != nil {
		return mcresponse.GetResMessage("deleteError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error deleting record(s): %v", delErr.Error()),
			Value:   nil,
		})
	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "Record(s) deleted successfully",
		Value:   commandTag.Delete(),
	})
}

// DeleteByParam method deletes or removes record(s) by query-parameters or where conditions
func (crud *Crud) DeleteByParam() mcresponse.ResponseMessage {
	// perform crud-task action, include where-query(params):
	// compute delete query by query-params
	deleteQuery, dQErr := helper.ComputeDeleteQueryByParam(crud.TableName, crud.QueryParams, []string{})
	if dQErr != nil {
		return mcresponse.GetResMessage("deleteError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error computing delete-query: %v", dQErr.Error()),
			Value:   nil,
		})
	}
	commandTag, delErr := crud.AppDb.Exec(context.Background(), deleteQuery)
	if delErr != nil {
		return mcresponse.GetResMessage("deleteError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error deleting record(s): %v", delErr.Error()),
			Value:   nil,
		})
	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "Record(s) deleted successfully",
		Value:   commandTag.Delete(),
	})
}

// DeleteAll method deletes or removes all records in the tables. Recommended for admin-users only
// Use if and only if you know what you are doing
func (crud *Crud) DeleteAll() mcresponse.ResponseMessage {
	// ***** perform DELETE-ALL-RECORDS FROM A TABLE, IF RELATIONS/CONSTRAINTS PERMIT *****
	// ***** && IF-AND-ONLY-IF-YOU-KNOW-WHAT-YOU-ARE-DOING *****
	// compute delete query
	delQuery := fmt.Sprintf("DELETE FROM %v", crud.TableName)
	commandTag, delErr := crud.AppDb.Exec(context.Background(), delQuery)
	if delErr != nil {
		return mcresponse.GetResMessage("deleteError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error deleting record(s): %v", delErr.Error()),
			Value:   nil,
		})
	}

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "Record(s) deleted successfully",
		Value:   commandTag.Delete(),
	})
}
