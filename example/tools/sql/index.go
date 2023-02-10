package sql

import "basic/tools/utils"

// InitializeTables 初始化数据表
func InitializeTables() []*utils.InitTable {
	tables := []*utils.InitTable{
		BasicAdminRBACItem,       //	权限表
		BasicAdminRBACChild,      //	权限表
		BasicAdminRBACAssignment, //	权限表
		BasicAccessLogs,          //	访问日志
		BasicAdminMenu,           //	管理菜单
		BasicAdminSetting,        //	管理设置
		BasicAdminUser,           //	后台管理
	}

	// 项目表
	tables = append(tables, []*utils.InitTable{}...)
	return tables
}

// InitializeAuth 初始化RBAC
func InitializeAuth() *utils.Permission {
	return &utils.Permission{
		Roles: []string{"代理管理员", "普通管理员"},
		//	通过的权限
		RoleOnlyRouter: map[string][]string{
			"普通管理员": {},
		},
		//	过滤的权限
		RoleFilterRouter: map[string][]string{
			"代理管理员": {"数据库表", "数据表信息", "权限数组", "角色列表", "角色更新", "角色新增", "角色删除", "菜单列表", "菜单更新"},
		},
	}
}
