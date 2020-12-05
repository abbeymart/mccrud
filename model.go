// @Author: abbeymart | Abi Akindele | @Created: 2020-12-01 | @Updated: 2020-12-01
// @Company: mConnect.biz | @License: MIT
// @Description: crud model specs

package mccrud

import (
	"encoding/json"
	"fmt"
	"strconv"
)
import "github.com/abbeymart/mcresponsego"
import "github.com/asaskevich/govalidator"

type Model struct {
	modelParams ModelType
	options     ModelOptionsType
}

// Methods
func (model Model) GetParentRelations() []ModelRelationType {
	// extract relations/collections where targetTable == model-TableName
	var parentRelations []ModelRelationType
	modelRelations := model.modelParams.Relations
	for _, item := range modelRelations {
		if item.TargetTable == model.modelParams.TableName {
			parentRelations = append(parentRelations, item)
		}
	}
	return parentRelations
}

func (model Model) GetChildRelations() []ModelRelationType {
	// extract relations/collections where sourceTable == model-TableName
	var childRelations []ModelRelationType
	modelRelations := model.modelParams.Relations
	for _, item := range modelRelations {
		if item.SourceTable == model.modelParams.TableName {
			childRelations = append(childRelations, item)
		}
	}
	return childRelations
}

func (model Model) GetParentTables() []string {
	var parentTables []string
	parentRelations := model.GetParentRelations()
	for _, rel := range parentRelations {
		parentTables = append(parentTables, rel.SourceTable)
	}
	return parentTables
}

func (model Model) GetChildTables() []string {
	var childTables []string
	childRelations := model.GetChildRelations()
	for _, rel := range childRelations {
		childTables = append(childTables, rel.TargetTable)
	}
	return childTables
}

func (model Model) ComputeDocValueType(docValue ValueParamType) ValueToDataType {
	computedType := ValueToDataType{}
	// perform computation of doc-value-types
	for key, val := range docValue {
		// array check
		//if govalidator.IsType(val, "string") {}
		switch fmt.Sprintf("%T", val) {
		case "[]string":
			computedType[key] = ARRAYOFSTRING
		case "[]int":
			computedType[key] = ARRAYOFNUMBER
		case "[]float64":
			computedType[key] = ARRAYOFNUMBER
		case "[]bool":
			computedType[key] = ARRAYOFBOOLEAN
		case "[]map":
			computedType[key] = ARRAYOFOBJECT
		case "[]struct":
			computedType[key] = ARRAYOFOBJECT
		case "[]":
			computedType[key] = ARRAY
		case "map":
			computedType[key] = MAP
		case "struct":
			computedType[key] = OBJECT
		case "string":
			// compute string value
			jsonStr, _ := json.Marshal(val)
			strVal := string(jsonStr)
			var strToNum float64
			if val, err := strconv.Atoi(strVal); err == nil {
				strToNum = float64(val)
			}
			// check all string-based formats
			// TODO: ISO2, ISO3, Currency, Mime, JWT, PostalCode
			if govalidator.IsEmail(strVal) {
				computedType[key] = EMAIL
			} else if govalidator.IsUnixTime(strVal) {
				computedType[key] = DATETIME
			} else if govalidator.IsTime(strVal, "HH:MM:SS") {
				computedType[key] = TIME
			} else if govalidator.IsMongoID(strVal) {
				computedType[key] = MONGODBID
			} else if govalidator.IsUUID(strVal) {
				computedType[key] = UUID
			} else if govalidator.IsUUIDv3(strVal) {
				computedType[key] = UUID3
			} else if govalidator.IsUUIDv4(strVal) {
				computedType[key] = UUID4
			} else if govalidator.IsUUIDv5(strVal) {
				computedType[key] = UUID5
			} else if govalidator.IsMD4(strVal) {
				computedType[key] = MD4
			} else if govalidator.IsMD5(strVal) {
				computedType[key] = MD5
			} else if govalidator.IsSHA1(strVal) {
				computedType[key] = SHA1
			} else if govalidator.IsSHA256(strVal) {
				computedType[key] = SHA256
			} else if govalidator.IsSHA384(strVal) {
				computedType[key] = SHA384
			} else if govalidator.IsSHA512(strVal) {
				computedType[key] = SHA512
			} else if govalidator.IsJSON(strVal) {
				computedType[key] = JSON
			} else if govalidator.IsCreditCard(strVal) {
				computedType[key] = CREDITCARD
			} else if govalidator.IsURL(strVal) {
				computedType[key] = URL
			} else if govalidator.IsDNSName(strVal) {
				computedType[key] = DOMAINNAME
			} else if govalidator.IsPort(strVal) {
				computedType[key] = PORT
			} else if govalidator.IsIP(strVal) {
				computedType[key] = IP
			} else if govalidator.IsIPv4(strVal) {
				computedType[key] = IP4
			} else if govalidator.IsIPv6(strVal) {
				computedType[key] = IP6
			} else if govalidator.IsIMEI(strVal) {
				computedType[key] = IMEI
			} else if govalidator.IsLatitude(strVal) {
				computedType[key] = LATITUDE
			} else if govalidator.IsLongitude(strVal) {
				computedType[key] = LONGITUDE
			} else if govalidator.IsMAC(strVal) {
				computedType[key] = MACADDRESS
			} else if govalidator.IsInt(strVal) {
				computedType[key] = INTEGER
			} else if govalidator.IsPositive(strToNum) {
				computedType[key] = POSITIVE
			} else if govalidator.IsNegative(strToNum) {
				computedType[key] = NEGATIVE
			} else if govalidator.IsNatural(strToNum) {
				computedType[key] = NATURAL
			} else {
				computedType[key] = STRING
			}
		case "int":
			computedType[key] = INTEGER
		case "float64":
			computedType[key] = FLOAT
		case "bool":
			computedType[key] = BOOLEAN
		default:
			computedType[key] = UNDEFINED
		}
	}
	return computedType
}

func (model Model) UpdateDefaultValue(docValue ValueParamType) ValueParamType {
	// set default values, for null fields | then setValue (transform), if specified
	// set base docValue
	setDocValue := docValue
	// perform update of default/set-values for the doc-values => modelDocValue

	return setDocValue
}

func (model Model) ValidateDocValue(modelDocValue ValueParamType, taskName string) ValidateResponseType {
	setDocValue := modelDocValue
	fmt.Println(setDocValue)
	fmt.Println(taskName)
	// perform validation of model-doc-value

	// return success
	var errMsg = ErrorType{}
	return ValidateResponseType{Ok: true, Errors: errMsg}
}

func (model Model) Save(params CrudTaskType, options CrudOptionsType) mcresponse.ResponseMessage {

	return mcresponse.ResponseMessage{}
}

func (model Model) Get(params CrudTaskType, options CrudOptionsType) mcresponse.ResponseMessage {

	return mcresponse.ResponseMessage{}
}

func (model Model) GetStream(params CrudTaskType, options CrudOptionsType) mcresponse.ResponseMessage {

	return mcresponse.ResponseMessage{}
}

func (model Model) Delete(params CrudTaskType, options CrudOptionsType) mcresponse.ResponseMessage {

	return mcresponse.ResponseMessage{}
}
