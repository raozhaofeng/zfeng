package models

import (
	"database/sql"
	"github.com/raozhaofeng/beego/db"
	"github.com/raozhaofeng/beego/db/define"
)

// AdminAuthAssignment 模型
type AdminAuthAssignment struct {
	define.Db
}

// AdminAuthAssignmentsAttrs 属性
type AdminAuthAssignmentsAttrs struct {
	ItemName  string `json:"item_name"`  //	名称
	UserId    int64  `json:"user_id"`    //	用户ID
	CreatedAt int64  `json:"created_at"` //	创建时间
}

// NewAdminAuthAssignment 创建模型
func NewAdminAuthAssignment(tx *sql.Tx) *AdminAuthAssignment {
	return &AdminAuthAssignment{
		db.Manager.NewInterfaceDb(tx).Table("admin_auth_assignment"),
	}
}

// GetAdminRoleList 获取管理员角色列表
func (c *AdminAuthAssignment) GetAdminRoleList(adminId int64) []string {
	var roleList []string
	c.Field("item_name").
		AndWhere("user_id = ?", adminId).
		Query(func(rows *sql.Rows) {
			var role string
			_ = rows.Scan(&role)
			roleList = append(roleList, role)
		})
	return roleList
}
