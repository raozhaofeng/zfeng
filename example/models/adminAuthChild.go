package models

import (
	"database/sql"
	"github.com/raozhaofeng/zfeng/database"
	"github.com/raozhaofeng/zfeng/database/define"
	"github.com/raozhaofeng/zfeng/utils"
	"strings"
)

// AdminAuthChild 模型
type AdminAuthChild struct {
	define.Db
}

// AdminAuthChildAttrs 属性
type AdminAuthChildAttrs struct {
	Parent string `json:"parent"` //	父级
	Child  string `json:"child"`  //	子级
	Type   int    `json:"type"`   //	类型
}

// NewAdminAuthChild 创建模型
func NewAdminAuthChild(tx *sql.Tx) *AdminAuthChild {
	return &AdminAuthChild{
		database.DbPool.NewDb(tx).Table("admin_auth_child"),
	}
}

func (c *AdminAuthChild) FindOne() *AdminAuthChildAttrs {
	attrs := new(AdminAuthChildAttrs)
	NewAdminAuthItem(nil).QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Parent, &attrs.Child, &attrs.Type)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// GetRolesRouteList 获取角色组路由键值对
func (c *AdminAuthChild) GetRolesRouteList(roleList []string) map[string]string {
	data := map[string]string{}
	c.Table("admin_auth_child as a").Field("b.parent, b.child").
		AndWhere("a.parent in ("+strings.Join(define.InString(roleList), ",")+")").
		Join("LEFT JOIN", "admin_auth_child as b", "a.child=b.parent").
		Query(func(rows *sql.Rows) {
			var prent, child string
			_ = rows.Scan(&prent, &child)
			data[prent] = child
		})
	return data
}

// GetRouteRoleCheckedList 所有路由角色选中对列表
func (c *AdminAuthChild) GetRouteRoleCheckedList(rolesRouteList []string) map[string]bool {
	data := map[string]bool{}

	NewAdminAuthItem(nil).Field("name").
		AndWhere("type=?", AdminAuthItemTypeRouteName).
		Query(func(rows *sql.Rows) {
			var name string
			_ = rows.Scan(&name)
			data[name] = utils.SliceStringIndexOf(name, rolesRouteList) > -1
		})
	return data
}
