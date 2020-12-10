// @Author: abbeymart | Abi Akindele | @Created: 2020-12-08 | @Updated: 2020-12-08
// @Company: mConnect.biz | @License: MIT
// @Description: optional access methods, to be used as middleware, prior to CRUD operation

package mccrud

import (
	"fmt"
	"github.com/abbeymart/mcresponsego"
	mcutils "github.com/abbeymart/mcutilsgo"
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
			if err := rows.Scan(&id); err != nil {
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
			Message: "AYou are not authorized to perform the requested action/task",
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
		Ok bool
		IsAdmin bool
		IsActive bool
		UserId string
		Group string
		Groups []string
	}{
		Ok: taskPermitted,
		IsAdmin: isAdmin,
		IsActive: isActive,
		UserId: userId,
		Group: group,
		Groups: groups,
	}
	// if all went well
	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "Action authorised / permitted.",
		Value:   value,
	})
}

func (crud Crud) CheckTaskAccess(params CheckAccessParamsType) mcresponse.ResponseMessage {
	// validate current user active status: by token (API) and user/loggedIn-status

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
