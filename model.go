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

type CrudMethods interface {
	Save()
	Get()
	Delete()
}
type CrudSave interface {
	Save()
}
type CrudGet interface {
	Get()
}
type CrudDelete interface {
	Delete()
}

type Model struct {
	ModelType
	ModelOptionsType
}

// Methods
func (model Model) GetParentRelations() []ModelRelationType {
	// extract relations/collections where targetTable == model-TableName
	var parentRelations []ModelRelationType
	modelRelations := model.Relations
	for _, item := range modelRelations {
		if item.TargetTable == model.TableName {
			parentRelations = append(parentRelations, item)
		}
	}
	return parentRelations
}

func (model Model) GetChildRelations() []ModelRelationType {
	// extract relations/collections where sourceTable == model-TableName
	var childRelations []ModelRelationType
	modelRelations := model.Relations
	for _, item := range modelRelations {
		if item.SourceTable == model.TableName {
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

func (model Model) ComputeRecordValueType(recordValue ValueParamType) ValueToDataType {
	computedType := ValueToDataType{}
	// perform computation of doc-value-types
	for key, val := range recordValue {
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

func (model Model) UpdateDefaultValue(recordValue ValueParamType) ValueParamType {
	// set default values, for null fields | then setValue (pre-set/transform), if specified
	// set base recordValue
	setRecordValue := recordValue
	// perform update of default/set-values for the doc-values => modelRecordValue
	for key, fieldValue := range recordValue {
		// defaultValue setting applies to FieldDescType only | otherwise, the value is required (not null)
		// transform fieldDesc to interface{} for type checking
		var fieldDescType interface{} = model.RecordDesc[key]
		//fieldValue := recordValue[key]
		// set default values
		if fieldValue != nil {
			switch fieldDescType.(type) {
			case FieldDescType:
				// type of defaultValue and fieldValue must be equivalent (re: validateMethod)
				fieldDesc := model.RecordDesc[key]
				if fieldDesc.DefaultValue != nil {
					defaultValue := fieldDesc.DefaultValue()
					// defaultValue and fieldValue types must match => validation-check
					// update setRecordValue for the key/field-column
					setRecordValue[key] = defaultValue
				}
			}
		}
		// setValue / transform field-value prior-to/before save-task (create / update)
		switch fieldDescType.(type) {
		case FieldDescType:
			setFieldValue := setRecordValue[key]
			if setFieldValue != nil && model.RecordDesc[key].SetValue != nil {
				// set/pre-set setRecordValue for the key/field-column
				setRecordValue[key] = model.RecordDesc[key].SetValue(recordValue)
			}
		}
	}
	return setRecordValue
}

func (model Model) ValidateRecordValue(modelRecordValue ValueParamType, taskName string) ValidateResponseType {
	// perform validation of model-record-value
	// recommendation: use updated recordValue, defaultValues and setValues, prior to validation
	// get recordValue transformed types
	recordValueTypes := model.ComputeRecordValueType(modelRecordValue)
	// model-description/definition
	recordDesc := model.RecordDesc
	// combine errors/messages
	validateErrorMessage := map[string]string{}
	// perform model-recordValue validation
	for key, recordFieldValue := range modelRecordValue {
		// check field description / definition exists
		if recordFieldDesc, ok := recordDesc[key]; ok {
			// transform recordFieldDesc to interface{} for type checking
			var recordFieldDescType interface{} = recordFieldDesc
			switch recordFieldDescType.(type) {
			case FieldDescType:
				// validate fieldValue and fieldDesc (model) types
				if recordValueTypes[key] != recordFieldDesc.FieldType {
					errMsg := fmt.Sprintf("Invalid Type for:  %v. Expected %v, Got %v", key, recordFieldDesc.FieldType, recordValueTypes[key])
					if recordFieldDesc.ValidateMessage != "" {
						validateErrorMessage[key] = recordFieldDesc.ValidateMessage + " :: " + errMsg
					} else {
						validateErrorMessage[key] = errMsg
					}
				}
				// validate allowNull, fieldLength, min/maxValues...| user-defined-validation-methods
				// use values from transform docValue, including default/set-values
				// nullCheck, if recordField value is not specified
				if recordFieldValue == nil && !recordFieldDesc.AllowNull {
					errMsg := fmt.Sprintf("Value is required for: %v. Can't be Null", key)
					if recordFieldDesc.ValidateMessage != "" {
						validateErrorMessage[key+"-nullValidation"] = recordFieldDesc.ValidateMessage + " :: " + errMsg
					} else {
						validateErrorMessage[key+"-nullValidation"] = errMsg
					}
				}
				// validate field-value-type constraints: fieldLength, min/maxValues..
				switch recordFieldValue.(type) {
				case string:
					if fieldValue, ok := recordFieldValue.(string); ok {
						if recordFieldDesc.FieldLength > 0 {
							fieldLength := len(fieldValue)
							if uint(fieldLength) > recordFieldDesc.FieldLength {
								errMsg := fmt.Sprintf("Size of %v cannot be longer than %v", key, recordFieldDesc.FieldLength)
								if recordFieldDesc.ValidateMessage != "" {
									validateErrorMessage[key+"-lengthValidation"] = recordFieldDesc.ValidateMessage + " :: " + errMsg
								} else {
									validateErrorMessage[key+"-lengthValidation"] = errMsg
								}
							}
						}
						// Perform field level validation-methods
						if recordFieldDesc.Validate != nil {
							valRes := recordFieldDesc.Validate(fieldValue)
							if !valRes {
								validateErrorMessage[key+"-validationError"] = fmt.Sprintf("Error validating the field-value: %v", key)
							}
						}
					} else {
						validateErrorMessage[key+"-transformError"] = fmt.Sprintf("Error processing the field-value type / format for: %v", key)
					}
				case int:
					if fieldValue, ok := recordFieldValue.(int); ok {
						if uint(fieldValue) < recordFieldDesc.MinValue && uint(fieldValue) > recordFieldDesc.MaxValue {
							errMsg := fmt.Sprintf("Value of: %v must be greater than %v, and less than %v", key, recordFieldDesc.MinValue, recordFieldDesc.MaxValue)
							if recordFieldDesc.ValidateMessage != "" {
								validateErrorMessage[key+"-minMaxValidation"] = recordFieldDesc.ValidateMessage + " :: " + errMsg
							} else {
								validateErrorMessage[key+"-minMaxValidation"] = errMsg
							}
						} else if uint(fieldValue) < recordFieldDesc.MinValue {
							errMsg := fmt.Sprintf("Value of: %v must be greater than %v", key, recordFieldDesc.MinValue)
							if recordFieldDesc.ValidateMessage != "" {
								validateErrorMessage[key+"-minValidation"] = recordFieldDesc.ValidateMessage + " :: " + errMsg
							} else {
								validateErrorMessage[key+"-minValidation"] = errMsg
							}
						} else if uint(fieldValue) > recordFieldDesc.MaxValue {
							errMsg := fmt.Sprintf("Value of: %v must be less than %v", key, recordFieldDesc.MaxValue)
							if recordFieldDesc.ValidateMessage != "" {
								validateErrorMessage[key+"-maxValidation"] = recordFieldDesc.ValidateMessage + " :: " + errMsg
							} else {
								validateErrorMessage[key+"-maxValidation"] = errMsg
							}
						}
						// Perform field level validation-methods
						if recordFieldDesc.Validate != nil {
							valRes := recordFieldDesc.Validate(fieldValue)
							if !valRes {
								validateErrorMessage[key+"-validationError"] = fmt.Sprintf("Error validating the field-value: %v", key)
							}
						}
					} else {
						validateErrorMessage[key+"-transformError"] = fmt.Sprintf("Error processing the field-value type / format for: %v", key)
					}
				case float32, float64:
					if fieldValue, ok := recordFieldValue.(float64); ok {
						if uint(fieldValue) < recordFieldDesc.MinValue && uint(fieldValue) > recordFieldDesc.MaxValue {
							errMsg := fmt.Sprintf("Value of: %v must be greater than %v, and less than %v", key, recordFieldDesc.MinValue, recordFieldDesc.MaxValue)
							if recordFieldDesc.ValidateMessage != "" {
								validateErrorMessage[key+"-minMaxValidation"] = recordFieldDesc.ValidateMessage + " :: " + errMsg
							} else {
								validateErrorMessage[key+"-minMaxValidation"] = errMsg
							}
						} else if uint(fieldValue) < recordFieldDesc.MinValue {
							errMsg := fmt.Sprintf("Value of: %v must be greater than %v", key, recordFieldDesc.MinValue)
							if recordFieldDesc.ValidateMessage != "" {
								validateErrorMessage[key+"-minValidation"] = recordFieldDesc.ValidateMessage + " :: " + errMsg
							} else {
								validateErrorMessage[key+"-minValidation"] = errMsg
							}
						} else if uint(fieldValue) > recordFieldDesc.MaxValue {
							errMsg := fmt.Sprintf("Value of: %v must be less than %v", key, recordFieldDesc.MaxValue)
							if recordFieldDesc.ValidateMessage != "" {
								validateErrorMessage[key+"-maxValidation"] = recordFieldDesc.ValidateMessage + " :: " + errMsg
							} else {
								validateErrorMessage[key+"-maxValidation"] = errMsg
							}
						}
						// Perform field level validation-methods
						if recordFieldDesc.Validate != nil {
							valRes := recordFieldDesc.Validate(fieldValue)
							if !valRes {
								validateErrorMessage[key+"-validationError"] = fmt.Sprintf("Error validating the field-value: %v", key)
							}
						}
					} else {
						validateErrorMessage[key+"-transformError"] = fmt.Sprintf("Error processing the field-value type / format for: %v", key)
					}
				case []string:
					if fieldValue, ok := recordFieldValue.([]string); ok {
						// Perform field level validation-methods
						if recordFieldDesc.Validate != nil {
							valRes := recordFieldDesc.Validate(fieldValue)
							if !valRes {
								validateErrorMessage[key+"-validationError"] = fmt.Sprintf("Error validating the field-value: %v", key)
							}
						}
					} else {
						validateErrorMessage[key+"-transformError"] = fmt.Sprintf("Error processing the field-value type / format for: %v", key)
					}
				case []int:
					if fieldValue, ok := recordFieldValue.([]int); ok {
						// Perform field level validation-methods
						if recordFieldDesc.Validate != nil {
							valRes := recordFieldDesc.Validate(fieldValue)
							if !valRes {
								validateErrorMessage[key+"-validationError"] = fmt.Sprintf("Error validating the field-value: %v", key)
							}
						}
					} else {
						validateErrorMessage[key+"-transformError"] = fmt.Sprintf("Error processing the field-value type / format for: %v", key)
					}
				case []float64, []float32:
					if fieldValue, ok := recordFieldValue.([]float64); ok {
						// Perform field level validation-methods
						if recordFieldDesc.Validate != nil {
							valRes := recordFieldDesc.Validate(fieldValue)
							if !valRes {
								validateErrorMessage[key+"-validationError"] = fmt.Sprintf("Error validating the field-value: %v", key)
							}
						}
					} else {
						validateErrorMessage[key+"-transformError"] = fmt.Sprintf("Error processing the field-value type / format for: %v", key)
					}
				case []struct{}:
					if fieldValue, ok := recordFieldValue.([]struct{}); ok {
						// Perform field level validation-methods
						if recordFieldDesc.Validate != nil {
							valRes := recordFieldDesc.Validate(fieldValue)
							if !valRes {
								validateErrorMessage[key+"-validationError"] = fmt.Sprintf("Error validating the field-value: %v", key)
							}
						}
					} else {
						validateErrorMessage[key+"-transformError"] = fmt.Sprintf("Error processing the field-value type / format for: %v", key)
					}
				}
			default:
				// validate field-value/type
				// use values from transform docValue, including default/set-values
				//if fieldValue, ok := modelRecordValue[key]; ok {
				//	fmt.Println(fieldValue)
				//}
			}
		} else {
			validateErrorMessage[key] = fmt.Sprintf("Invalid key: %v is not defined in the model", key)
		}
	}

	// perform user-defined recordValue validation
	// get validate method for the recordValue task by taskName (e.g. registerUser, login, saveProfile etc.)
	if modelValidateMethod, ok := model.ValidateMethods[taskName]; ok {
		valRes := modelValidateMethod(modelRecordValue)
		if !valRes.Ok {
			var modelErrorMsg = ""
			for _, msg := range valRes.Errors {
				if modelErrorMsg != "" {
					modelErrorMsg += " | " + msg
				} else {
					modelErrorMsg = msg
				}
			}
			validateErrorMessage[model.TableName + "-validationError"] = modelErrorMsg
		}
	}

	// check validateErrors
	if len(validateErrorMessage) != 0 {
		return ValidateResponseType{
			Ok:     false,
			Errors: validateErrorMessage,
		}
	}

	// return success, if validation process has been completed without errors
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
