// @Author: abbeymart | Abi Akindele | @Created: 2020-12-01 | @Updated: 2020-12-01
// @Company: mConnect.biz | @License: MIT
// @Description: crud model specs

package mccrud

import "fmt"
import "github.com/abbeymart/mcresponsego"

type Model struct {
	model ModelType
	options ModelOptionsType
}

// Methods
func(model Model) ComputeDocValueType(docValue ValueParamType) ValueToDataType {
	var computedType ValueToDataType
	fmt.Println(docValue)
	// perform computation of doc-value-types

	return computedType
}

func(model Model) UpdateDefaultValue(docValue ValueParamType) ValueParamType {
	setDocValue := docValue
	// perform update of default/set-values for the doc-values => modelDocValue

	return setDocValue
}

func(model Model) ValidateDocValue(modelDocValue ValueParamType, taskName string) ValidateResponseType {
	setDocValue := modelDocValue
	fmt.Println(setDocValue)
	fmt.Println(taskName)
	// perform validation of model-doc-value

	// return success
	var errMsg = ErrorType{}
	return ValidateResponseType{ok: true, errors: errMsg}
}

func(model Model) Save(params CrudTaskType, options CrudOptionsType) mcresponse.ResponseMessage {

	return mcresponse.ResponseMessage{}
}

func(model Model) Get(params CrudTaskType, options CrudOptionsType) mcresponse.ResponseMessage {

	return mcresponse.ResponseMessage{}
}

func(model Model) GetStream(params CrudTaskType, options CrudOptionsType) mcresponse.ResponseMessage {

	return mcresponse.ResponseMessage{}
}

func(model Model) Delete(params CrudTaskType, options CrudOptionsType) mcresponse.ResponseMessage {

	return mcresponse.ResponseMessage{}
}
