// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: go: mConnect

package mccrud

import (
	"fmt"
	"github.com/abbeymart/mcresponsego"
)

func (crud Crud) TaskPermission() mcresponse.ResponseMessage {
	fmt.Println(crud)

	return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
		Message: "action-params is required to perform save operation.",
		Value:   nil,
	})
}

func (crud Crud) CheckTaskAccess() mcresponse.ResponseMessage {
	fmt.Println(crud)

	return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
		Message: "action-params is required to perform save operation.",
		Value:   nil,
	})
}

func (crud Crud) GetRoleServices() []RoleServiceType {
	fmt.Println(crud)

	return []RoleServiceType{}
}

func (crud Crud) GetCurrentRecord() mcresponse.ResponseMessage {
	fmt.Println(crud)

	return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
		Message: "action-params is required to perform save operation.",
		Value:   nil,
	})
}
