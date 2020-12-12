// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: compute update-SQL script

package helper

import (
	"errors"
	mccrud "github.com/abbeymart/mccrudgo"
)

func ComputeUpdateQueryById(tableName string, actionParams mccrud.ActionParamsType , docIds []string) (mccrud.UpdateScriptResponseType, error) {

	return mccrud.UpdateScriptResponseType{}, errors.New("development")
}
