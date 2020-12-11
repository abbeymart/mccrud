// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: optional access methods, to be used as middleware, prior to CRUD operation

package mccrud

import (
	"database/sql"
	"fmt"
	"github.com/abbeymart/mcresponsego"
	"github.com/abbeymart/mcutilsgo"
	"strings"
	"time"
)

func (crud Crud) TaskPermission(taskType string) mcresponse.ResponseMessage {
	// TaskType: "create", "update", "delete"/"remove", "read"
	// permit task(crud): by owner, role/group (on coll/table or doc/record(s)) or admin
	// # validation access variables
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
	accessRes := crud.CheckTaskAccess(
		CheckAccessParamsType{
			accessDb:     crud.AccessDb,
			userInfo:     crud.UserInfo,
			tableName:    crud.TableName,
			docIds:       crud.DocIds,
			accessTable:  crud.AccessTable,
			userTable:    crud.UserTable,
			roleTable:    crud.RoleTable,
			serviceTable: crud.ServiceTable,
		})
	// capture roleServices value
	if accessRes.Code != "success" {
		return accessRes
	}

	// get access info value
	accessResValue := accessRes.Value
	accessInfo := accessResValue.(CheckAccessType)
	accessUserId := accessInfo.UserId
	var docIds []string

	// determine records/documents ownership
	if len(crud.DocIds) > 0 && accessUserId != "" && accessInfo.IsActive {
		// set docIds
		docIds = crud.DocIds
		// SQL script
		sqlScript := fmt.Sprintf("SELECT id FROM %v WHERE id IN $1 AND created_by = $2", crud.TableName)
		rows, err := crud.AppDb.Query(sqlScript, crud.DocIds, accessUserId)
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
			}
		}
		if rowCount == len(crud.DocIds) {
			ownerPermitted = true
		}
	}
	isAdmin = accessInfo.IsAdmin
	isActive = accessInfo.IsActive
	roleServices = accessInfo.RoleServices
	userId = accessInfo.UserId
	group = accessInfo.Group
	groups = accessInfo.Groups
	tableId = accessInfo.TableId

	// validate active status
	if !isActive {
		return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
			Message: "Account is not active. Validate active status",
			Value:   nil,
		})
	}
	// validate roleServices permission, for non-admin users
	if !isActive && len(roleServices) < 1 {
		return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
			Message: "You are not authorized to perform the requested action/task",
			Value:   nil,
		})
	}
	// filter the roleServices by categories ("collection | table" and "record or document")
	collTabFunc := func(item RoleServiceType) bool {
		return item.ServiceCategory == tableId
	}
	recordFunc := func(item RoleServiceType) bool {
		return mcutils.ArrayContains(docIds, item.ServiceCategory)
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

	//roleCreateFunc := func(it1 string, it2 RoleServiceType) bool {
	//	return it2.ServiceId == it1 && it2.CanCreate
	//}
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
		// test if any/some of the roleRecords it1/it2 met the
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
		case CrudTasks().Create, CrudTasks().Insert:
			// collection/table level access | only tableName Id was included in serviceIds
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
		case CrudTasks().Update:
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
			// document/record level access: all docIds must have at least a match in the roleRec
			if len(docIds) > 0 {
				recordPermitted = func() bool {
					for _, v := range docIds {
						if !roleRecFunc(v, roleRecords, roleUpdateFunc) {
							return false
						}
					}
					return true
				}()
			}
		case CrudTasks().Delete, CrudTasks().Remove:
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
			// document/record level access: all docIds must have at least a match in the roleRec
			if len(docIds) > 0 {
				recordPermitted = func() bool {
					for _, v := range docIds {
						if !roleRecFunc(v, roleRecords, roleDeleteFunc) {
							return false
						}
					}
					return true
				}()
			}
		case CrudTasks().Read:
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
			// document/record level access: all docIds must have at least a match in the roleRec
			if len(docIds) > 0 {
				recordPermitted = func() bool {
					for _, v := range docIds {
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
			Value:   taskPermitted,
		})
	}
	// const value = {...ok, ...{isAdmin, isActive, userId, group, groups}};
	value := struct {
		Ok       bool
		IsAdmin  bool
		IsActive bool
		UserId   string
		Group    string
		Groups   []string
	}{
		Ok:       taskPermitted,
		IsAdmin:  isAdmin,
		IsActive: isActive,
		UserId:   userId,
		Group:    group,
		Groups:   groups,
	}
	// if all went well
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "Action authorised / permitted.",
		Value:   value,
	})
}

func (crud Crud) CheckTaskAccess(params CheckAccessParamsType) mcresponse.ResponseMessage {
	// validate current user active status: by token (API) and user/loggedIn-status

	// get the accessKey information for the user
	accessScript := fmt.Sprintf("SELECT expire from %v WHERE user_id=$1 AND token=$2 AND login_name=$3", params.accessTable)
	rowAccess := crud.AccessDb.QueryRow(accessScript, params.userInfo.UserId, params.userInfo.Token, params.userInfo.LoginName)
	// check login-status/expiration
	var accessExpire int64
	if err := rowAccess.Scan(&accessExpire); err != nil {
		return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
			Message: "Unauthorized: please ensure that you are logged-in",
			Value:   nil,
		})
	} else {
		if time.Now().Unix() > accessExpire {
			return mcresponse.GetResMessage("tokenExpired", mcresponse.ResponseMessageOptions{
				Message: "Access expired: please login to continue",
				Value:   nil,
			})
		}
	}
	// TODO: check current-user status/info
	var (
		uId      string
		group    string
		groups   []string
		isAdmin  bool
		isActive bool
	)
	userScript := fmt.Sprintf("SELECT id, group, groups, isAdmin, isActive from %v WHERE id=$1 AND is_active=$2", params.userTable)
	rowUser := crud.AccessDb.QueryRow(userScript, params.userInfo.UserId, true)
	// check login-status/expiration
	if err := rowUser.Scan(&uId, &group, &groups, &isAdmin, &isActive); err != nil {
		return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
			Message: "Unauthorized: user information not found or inactive",
			Value:   nil,
		})
	}

	// TODO: if all the above checks passed, check for role-services access by taskType
	// obtain table/collName/collId (_id) from serviceColl (repo for all resources)
	var (
		serviceId string
		category  string
	)
	serviceScript := fmt.Sprintf("SELECT id, category from %v WHERE name=$1", params.serviceTable)
	rowService := crud.AccessDb.QueryRow(serviceScript, params.tableName)
	// check login-status/expiration
	if err := rowService.Scan(&serviceId, &category); err != nil {
		return mcresponse.GetResMessage("unAuthorized", mcresponse.ResponseMessageOptions{
			Message: "Unauthorized: user information not found or inactive",
			Value:   nil,
		})
	}
	// # if permitted, include collId and docIds in serviceIds
	tableId := ""
	serviceIds := params.docIds
	catLowercase := strings.ToLower(category)
	if catLowercase == "table" || catLowercase == "collection" {
		tableId = serviceId
		serviceIds = append(serviceIds, serviceId)
	}

	var roleServices []RoleServiceType
	if len(serviceIds) > 0 {
		roleServices = crud.GetRoleServices(params.accessDb, params.roleTable, group, serviceIds)
	}

	permittedRes := CheckAccessType{
		UserId:       uId,
		Group:        group,
		Groups:       groups,
		IsActive:     isActive,
		IsAdmin:      isAdmin,
		RoleServices: roleServices,
		TableId:      tableId,
	}

	// if all went well
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "Action authorised / permitted.",
		Value:   permittedRes,
	})
}

func (crud Crud) GetRoleServices(accessDb *sql.DB, roleTable string, group string, serviceIds []string) []RoleServiceType {
	var roleServices []RoleServiceType
	roleScript := fmt.Sprintf("SELECT service_id, role_id, service_category, can_read, can_create, can_delete, can_update from %v WHERE service_id IN $1 AND group=$2 AND is_active=$3", roleTable)
	rows, err := accessDb.Query(roleScript, serviceIds, group, true)
	if err != nil {
		//errMsg := fmt.Sprintf("Db query Error: %v", err.Error())
		return roleServices
	}
	defer rows.Close()
	var (
		serviceId, roleId, serviceCategory       string
		canRead, canCreate, canDelete, canUpdate bool
	)
	if rows.Next() {
		if err := rows.Scan(&serviceId, &roleId, &serviceCategory, &canRead, &canCreate, &canDelete, &canUpdate); err == nil {
			roleServices = append(roleServices, RoleServiceType{
				ServiceId:       serviceId,
				RoleId:          roleId,
				ServiceCategory: serviceCategory,
				CanRead:         canRead,
				CanCreate:       canCreate,
				CanUpdate:       canUpdate,
				CanDelete:       canDelete,
			})
		}
	}

	return roleServices
}

func (crud Crud) GetCurrentRecord(recordType interface{}) mcresponse.ResponseMessage {
	//var currentRecords []interface{}
	roleScript := fmt.Sprintf("SELECT * from %v WHERE id IN $1", crud.TableName)
	rows, err := crud.AppDb.Query(roleScript, crud.DocIds)
	if err != nil {
		//errMsg := fmt.Sprintf("Db query Error: %v", err.Error())
		return mcresponse.ResponseMessage{}
	}
	defer rows.Close()


	return mcresponse.GetResMessage("paramsError", mcresponse.ResponseMessageOptions{
		Message: "action-params is required to perform save operation.",
		Value:   nil,
	})
}
