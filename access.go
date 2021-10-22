// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: optional access methods, to be used as middleware, prior to CRUD operation

package mccrud

import (
	"errors"
	"fmt"
	"github.com/abbeymart/mcresponse"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

// AccessInfoType for CheckUserAccess method value (interface{}) response,
// and to assert returned value
type AccessInfoType struct {
	UserId   string
	RoleId   string
	RoleIds  []string
	IsAdmin  bool
	IsActive bool
}

// TaskPermissionType for TaskPermission method value (interface{}) response,
// and to assert returned value
type TaskPermissionType struct {
	Ok       bool
	IsAdmin  bool
	IsActive bool
	UserId   string
	Role     string
	Roles    []string
}

func (crud *Crud) CheckTaskType() string {
	taskType := ""
	if len(crud.ActionParams) > 0 {
		actParam := crud.ActionParams[0]
		_, ok := actParam["id"]
		if !ok {
			if len(crud.RecordIds) > 0 || len(crud.QueryParams) > 0 {
				taskType = UpdateTask
			} else {
				taskType = CreateTask
			}
		} else {
			taskType = UpdateTask
		}
	}
	return taskType
}

// CheckTaskAccess method determines the access by role-assignment
func (crud *Crud) CheckTaskAccess() mcresponse.ResponseMessage {
	// validate current user active status: by token (API) and user/loggedIn-status
	accessRes := crud.CheckUserAccess()
	if accessRes.Code != "success" {
		return accessRes
	}
	// set current-user info for next steps
	var (
		uId      string
		roleId   string
		roleIds  []string
		isAdmin  bool
		isActive bool
	)
	val, ok := accessRes.Value.(AccessInfoType)
	if !ok {
		return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
			Message: "Error parsing user access information/value",
			Value:   nil,
		})
	}
	uId = val.UserId
	roleId = val.RoleId
	roleIds = val.RoleIds
	isAdmin = val.IsAdmin
	isActive = val.IsActive
	// determine records/documents ownership, for all records (atomic)
	ownerPermitted := false
	idLen := len(crud.RecordIds)
	if idLen > 0 && uId != "" && isActive {
		// SQL script
		inValues := ""
		for idCount, id := range crud.RecordIds {
			inValues += "'" + id + "'"
			if idLen > 1 && idCount < idLen-1 {
				inValues += ", "
			}
		}
		sqlScript := fmt.Sprintf("SELECT id FROM %v WHERE id IN (%v) AND created_by = $1", crud.TableName, inValues)
		rows, err := crud.AccessDb.Queryx(sqlScript, uId)
		if err != nil {
			errMsg := fmt.Sprintf("Db query Error: %v", err.Error())
			return mcresponse.GetResMessage("readError", mcresponse.ResponseMessageOptions{
				Message: errMsg,
				Value:   nil,
			})
		}
		defer func(rows *sqlx.Rows) {
			err := rows.Close()
			if err != nil {

			}
		}(rows)
		// check rows count
		var rowCount = 0
		for rows.Next() {
			var id string
			if recErr := rows.Scan(&id); recErr == nil {
				rowCount += 1
			}
		}
		// ensure complete records count, as requested
		if rowCount == len(crud.RecordIds) {
			ownerPermitted = true
		}
	}
	// if all the above checks passed, check for role-services access by taskType
	// obtain table/collName id(_id) from serviceTable/Coll (repo for all resources)
	var (
		serviceId string
		category  string
	)
	serviceScript := fmt.Sprintf("SELECT id, category from %v WHERE name=$1", crud.ServiceTable)
	serviceRow := crud.AccessDb.QueryRow(serviceScript, crud.TableName)
	// check error
	if err := serviceRow.Scan(&serviceId, &category); err != nil {
		return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Unauthorized: user information not found or inactive | %v", err.Error()),
			Value:   nil,
		})
	}
	// if permitted, include table/collId and recordIds in serviceIds
	tableId := ""
	serviceIds := crud.RecordIds
	catLowercase := strings.ToLower(category)
	if catLowercase == "table" || catLowercase == "collection" {
		tableId = serviceId
		serviceIds = append(serviceIds, serviceId)
	}
	// compute service-items/records
	var roleServices []RoleServiceType
	var rsErr error
	if len(serviceIds) > 0 {
		roleServices, rsErr = crud.GetRoleServices(crud.AccessDb, crud.RoleTable, roleId, serviceIds)
		if rsErr != nil {
			return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Action un-authorised / not-permitted | %v", rsErr.Error()),
				Value:   nil,
			})
		}
	}

	permittedRes := CheckAccessType{
		UserId:         uId,
		RoleId:         roleId,
		RoleIds:        roleIds,
		IsActive:       isActive,
		IsAdmin:        isAdmin,
		RoleServices:   roleServices,
		TableId:        tableId,
		OwnerPermitted: ownerPermitted,
	}

	if permittedRes.IsActive && permittedRes.IsAdmin {
		return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
			Message: "Action authorised / permitted.",
			Value:   permittedRes,
		})
	}
	recLen := len(permittedRes.RoleServices)
	if permittedRes.IsActive && recLen > 0 && recLen >= len(crud.RecordIds) {
		return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Access permitted for %v of %v service-items/records", recLen, len(crud.RecordIds)),
			Value:   permittedRes,
		})
	}
	return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
		Message: "Action authorised / permitted.",
		Value:   permittedRes,
	})
}

// GetRoleServices method process and returns the permission to user / user-group/roleId for the specified service items
func (crud *Crud) GetRoleServices(accessDb *sqlx.DB, roleTable string, userRoleId string, serviceIds []string) ([]RoleServiceType, error) {
	var roleServices []RoleServiceType
	// where-in-values
	inValues := ""
	idLen := len(serviceIds)
	for idCount, id := range serviceIds {
		inValues += "'" + id + "'"
		if idLen > 1 && idCount < idLen-1 {
			inValues += ", "
		}
	}
	roleScript := fmt.Sprintf("SELECT role_id, service_id, service_category, can_read, can_create, can_delete, can_update, can_crud from %v WHERE service_id IN (%v) AND role_id=$1 AND is_active=$2", roleTable, inValues)
	rows, err := accessDb.Queryx(roleScript, userRoleId, true)
	if err != nil {
		//errMsg := fmt.Sprintf("Db query Error: %v", err.Error())
		return roleServices, errors.New(fmt.Sprintf("%v", err.Error()))
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)
	var (
		roleId, serviceId, serviceCategory                string
		canRead, canCreate, canDelete, canUpdate, canCrud bool
	)
	for rows.Next() {
		if err := rows.Scan(&roleId, &serviceId, &serviceCategory, &canRead, &canCreate, &canDelete, &canUpdate, &canCrud); err == nil {
			roleServices = append(roleServices, RoleServiceType{
				ServiceId:       serviceId,
				RoleId:          roleId,
				ServiceCategory: serviceCategory,
				CanRead:         canRead,
				CanCreate:       canCreate,
				CanUpdate:       canUpdate,
				CanDelete:       canDelete,
				CanCrud:         canCrud,
			})
		}
	}

	return roleServices, nil
}

// TaskPermissionById method determines the access permission by owner, role/group (on coll/table or doc/record(s)) or admin
// for various : create/insert, update, delete/remove, read
func (crud *Crud) TaskPermissionById(taskType string) mcresponse.ResponseMessage {
	// permit crud : by owner, role (on table or record(s)) or admin
	// task permission access variables
	var (
		taskPermitted   = false
		ownerPermitted  = false
		recordPermitted = false
		tablePermitted  = false
		isAdmin         = false
		isActive        = false
		userId          = ""
		tableId         = ""
		group           = ""
		groups          []string
		roleServices    []RoleServiceType
	)
	// check role-based access
	accessRes := crud.CheckTaskAccess()
	// capture roleServices value
	if accessRes.Code != "success" {
		return accessRes
	}
	// get access-record
	accessRec, ok := accessRes.Value.(CheckAccessType)
	if !ok {
		return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
			Message: "Error parsing task access information/value",
			Value:   nil,
		})
	}
	// set access status variables
	ownerPermitted = accessRec.OwnerPermitted
	isAdmin = accessRec.IsAdmin
	isActive = accessRec.IsActive
	roleServices = accessRec.RoleServices
	userId = accessRec.UserId
	group = accessRec.RoleId
	groups = accessRec.RoleIds
	tableId = accessRec.TableId
	// validate active status
	if !isActive {
		return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
			Message: "Account is not active. Validate active status",
			Value:   nil,
		})
	}
	// validate task (roleServices) permission, for non-admin users
	if !isAdmin && len(roleServices) < 1 {
		return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
			Message: "You are not authorized to perform the requested action/task",
			Value:   nil,
		})
	}
	// filter the roleServices by categories ("collection | table" or "record | document")
	recordIds := crud.RecordIds
	collTabFunc := func(item RoleServiceType) bool {
		return item.ServiceCategory == tableId
	}
	recordFunc := func(item RoleServiceType) bool {
		return ArrayStringContains(recordIds, item.ServiceCategory)
	}

	var (
		roleTables, roleRecords []RoleServiceType
	)
	if len(roleServices) > 0 {
		for _, v := range roleServices {
			if collTabFunc(v) {
				roleTables = append(roleTables, v)
			}
		}
		for _, v := range roleServices {
			if recordFunc(v) {
				roleRecords = append(roleRecords, v)
			}
		}
	}
	// helper functions
	canCreateFunc := func(item RoleServiceType) bool {
		return item.CanCreate
	}
	canUpdateFunc := func(item RoleServiceType) bool {
		return item.CanUpdate
	}
	canDeleteFunc := func(item RoleServiceType) bool {
		return item.CanDelete
	}
	canReadFunc := func(item RoleServiceType) bool {
		return item.CanRead
	}

	roleUpdateFunc := func(it1 string, it2 RoleServiceType) bool {
		return it2.ServiceId == it1 && it2.CanUpdate
	}
	roleDeleteFunc := func(it1 string, it2 RoleServiceType) bool {
		return it2.ServiceId == it1 && it2.CanDelete
	}
	roleReadFunc := func(it1 string, it2 RoleServiceType) bool {
		return it2.ServiceId == it1 && it2.CanRead
	}

	roleRecFunc := func(it1 string, roleRecs []RoleServiceType, roleFunc RoleFuncType) bool {
		// test if any or some of the roleRecords it1/it2 met the access condition
		for _, it2 := range roleRecs {
			if roleFunc(it1, it2) {
				return true
			}
		}
		return false
	}
	// taskType specific permission(s)
	if !isAdmin && len(roleServices) > 0 {
		switch taskType {
		case CreateTask, InsertTask:
			// collection/table level access | only tableId was included in serviceIds
			// must be able to perform create on the specified tableId(s)
			if len(roleTables) > 0 {
				tablePermitted = func() bool {
					for _, v := range roleTables {
						if !canCreateFunc(v) {
							return false
						}
					}
					return true
				}()
			}
		case UpdateTask:
			// collection/table level access
			if len(roleTables) > 0 {
				tablePermitted = func() bool {
					for _, v := range roleTables {
						if !canUpdateFunc(v) {
							return false
						}
					}
					return true
				}()
			}
			// document/record level access: all recordIds must have at least a match in the roleRecords
			if len(recordIds) > 0 {
				recordPermitted = func() bool {
					for _, v := range recordIds {
						if !roleRecFunc(v, roleRecords, roleUpdateFunc) {
							return false
						}
					}
					return true
				}()
			}
		case DeleteTask, RemoveTask:
			// collection/table level access
			if len(roleTables) > 0 {
				tablePermitted = func() bool {
					for _, v := range roleTables {
						if !canDeleteFunc(v) {
							return false
						}
					}
					return true
				}()
			}
			// document/record level access: all recordIds must have at least a match in the roleRecords
			if len(recordIds) > 0 {
				recordPermitted = func() bool {
					for _, v := range recordIds {
						if !roleRecFunc(v, roleRecords, roleDeleteFunc) {
							return false
						}
					}
					return true
				}()
			}
		case ReadTask:
			// collection/table level access
			if len(roleTables) > 0 {
				tablePermitted = func() bool {
					for _, v := range roleTables {
						if !canReadFunc(v) {
							return false
						}
					}
					return true
				}()
			}
			// document/record level access: all recordIds must have at least a match in the roleRecords
			if len(recordIds) > 0 {
				recordPermitted = func() bool {
					for _, v := range recordIds {
						if !roleRecFunc(v, roleRecords, roleReadFunc) {
							return false
						}
					}
					return true
				}()
			}
		default:
			return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
				Message: "Unknown access type or access type not specified.",
				Value:   nil,
			})
		}
	}

	// overall access permitted
	taskPermitted = recordPermitted || tablePermitted || ownerPermitted || isAdmin

	if !taskPermitted {
		return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
			Message: "You are not authorized to perform the requested action/task.",
			Value: TaskPermissionType{
				Ok: taskPermitted,
			},
		})
	}
	// if all went well
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "Action authorised / permitted.",
		Value: TaskPermissionType{
			Ok:       taskPermitted,
			IsAdmin:  isAdmin,
			IsActive: isActive,
			UserId:   userId,
			Role:     group,
			Roles:    groups,
		},
	})
}

func (crud *Crud) TaskPermissionByParam(taskType string) mcresponse.ResponseMessage {
	// ids of records, from queryParams
	var recordIds []string
	if len(crud.CurrentRecords) < 1 {
		currentRecRes := crud.GetByParam()
		if currentRecRes.Code != "success" {
			return currentRecRes
		}
		result, ok := currentRecRes.Value.(GetResultType)
		if !ok {
			return mcresponse.GetResMessage("notFound", mcresponse.ResponseMessageOptions{
				Message: "Missing or Invalid record(s) for task-permission-by-queryParams",
				Value:   result,
			})
		}
		crud.CurrentRecords = result.Records
	}
	for _, rec := range crud.CurrentRecords {
		//val, _ := rec.(ActionParamType)
		id, ok := rec["id"].(string)
		if !ok {
			return mcresponse.GetResMessage("notFound", mcresponse.ResponseMessageOptions{
				Message: "Missing record(s) for task-permission-by-queryParams",
				Value:   rec,
			})
		}
		recordIds = append(recordIds, id)
	}
	crud.RecordIds = recordIds
	return crud.TaskPermissionById(taskType)
}

// CheckUserAccess method determines the user access status: active, valid login and admin
func (crud *Crud) CheckUserAccess() mcresponse.ResponseMessage {
	// validate current user active status: by token (API) and user/loggedIn-status
	// get the accessKey information for the user
	accessScript := fmt.Sprintf("SELECT expire from %v WHERE user_id=$1 AND token=$2 AND login_name=$3", crud.AccessTable)
	rowAccess := crud.AccessDb.QueryRow(accessScript, crud.UserInfo.UserId, crud.UserInfo.Token, crud.UserInfo.LoginName)
	// check login-status/expiration
	var accessExpire int64
	if err := rowAccess.Scan(&accessExpire); err != nil {
		return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Unauthorized: please ensure that you are logged-in: %v", err.Error()),
			Value:   nil,
		})
	} else {
		if (time.Now().Unix() * 1000) > accessExpire {
			return mcresponse.GetResMessage("tokenExpired", mcresponse.ResponseMessageOptions{
				Message: "Access expired: please login to continue",
				Value:   nil,
			})
		}
	}
	// check the current-user status/info
	var (
		uId      string
		roleIds  interface{} // IDs type
		isAdmin  bool
		isActive bool
		profile  interface{} // Profile type
	)
	userScript := fmt.Sprintf("SELECT id, role_ids, is_admin, profile, is_active from %v WHERE id=$1 AND is_active=$2", crud.UserTable)
	rowUser := crud.AccessDb.QueryRow(userScript, crud.UserInfo.UserId, true)
	if err := rowUser.Scan(&uId, &roleIds, &isAdmin, &profile, &isActive); err != nil {
		return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Unauthorized: user information not found or is inactive: %v", err.Error()),
			Value:   nil,
		})
	}
	// transform, roleIds and profile (base64String-values) to roleIdsVal and profileVal
	roleIdsModel := IDs{}
	profileModel := Profile{}

	rIdsVal, rErr := ConvertJsonBase64StringToTypeValue(roleIds, &roleIdsModel)
	if rErr != nil {
		return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error parsing roleIds-json-value: %v", rErr.Error()),
			Value:   nil,
		})
	}
	pVal, pErr := ConvertJsonBase64StringToTypeValue(profile, &profileModel)
	if pErr != nil {
		return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error parsing user-profile-json-value: %v", pErr.Error()),
			Value:   nil,
		})
	}

	roleIdsVal, rOk := rIdsVal.(IDs)
	if !rOk {
		return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error asserting-type of the parsed roleIds-json-value"),
			Value:   nil,
		})
	}
	profileVal, pOk := pVal.(Profile)
	if !pOk {
		return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Error asserting-type of the parsed user-profile-json-value"),
			Value:   nil,
		})
	}

	// if all went well
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "Action authorised / permitted.",
		Value: AccessInfoType{
			UserId:   uId,
			RoleId:   profileVal.RoleId,
			RoleIds:  roleIdsVal,
			IsAdmin:  isAdmin,
			IsActive: isActive,
		},
	})
}

// CheckLoginStatus method checks if the user exists and has active login status/token
func (crud *Crud) CheckLoginStatus() mcresponse.ResponseMessage {
	params := crud.UserInfo
	// check if user exists, from users table
	emailUsername := EmailUsername(params.LoginName)
	email := emailUsername.Email
	username := emailUsername.Username
	var uId string
	if email != "" {
		query := fmt.Sprintf("SELECT id from %v WHERE id=$1 AND email=$2", crud.UserTable)
		row := crud.AccessDb.QueryRow(query, params.UserId, email)
		err := row.Scan(&uId)
		if err != nil {
			return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Record not found for %v. Register a new account: %v", params.LoginName, err.Error()),
				Value:   nil,
			})
		}
	} else if username != "" {
		query := fmt.Sprintf("SELECT id from %v WHERE id=$1 AND username=$2", crud.UserTable)
		row := crud.AccessDb.QueryRow(query, params.UserId, username)
		err := row.Scan(&uId)
		if err != nil {
			return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
				Message: fmt.Sprintf("Record not found for %v. Register a new account: %v", params.LoginName, err.Error()),
				Value:   nil,
			})
		}
	} else {
		// invalid user-information provided
		return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
			Message: "Invalid user-information provided.",
			Value:   nil,
		})
	}
	// check loginName, userId and token validity... from access_keys table
	var expire int64
	query := fmt.Sprintf("SELECT expire from %v WHERE id=$1 AND login_name=$2 AND token=$3", crud.AccessTable)
	row := crud.AccessDb.QueryRow(query, params.UserId, params.LoginName, params.Token)
	err := row.Scan(&expire)
	if err != nil {
		return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
			Message: fmt.Sprintf("Access information for %v not found. Login first, or contact system administrator: %v", params.LoginName, err.Error()),
			Value:   nil,
		})
	}
	if (time.Now().Unix() * 1000) > expire {
		// Delete the expired access_keys | remove access-info from access_keys table
		delQuery := fmt.Sprintf("DELETE FROM %v WHERE id=$1 AND token=$2", crud.AccessTable)
		_, _ = crud.AppDb.Exec(delQuery, params.UserId, params.Token)
		return mcresponse.GetResMessage("tokenExpired", mcresponse.ResponseMessageOptions{
			Message: "Access expired: please login to continue",
			Value:   nil,
		})
	}
	// if all went well
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "Action authorised / Access permitted.",
		Value:   uId,
	})
}
