// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: compute where-SQL script

package helper

import (
	"errors"
	mccrud "github.com/abbeymart/mccrudgo"
)

func whereScriptErr(errMsg string) (mccrud.WhereScriptResponseType, error) {
	return mccrud.WhereScriptResponseType{
		WhereScript: nil,
		FieldValues:  nil,
	}, errors.New(errMsg)
}

func ComputeWhereQuery(where mccrud.WhereParamType) (mccrud.WhereScriptResponseType, error) {


	return whereScriptErr("development")
}

