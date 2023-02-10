package models

import (
	"database/sql"
	"encoding/json"
	"github.com/raozhaofeng/beego/db"
	"github.com/raozhaofeng/beego/db/define"
	"github.com/raozhaofeng/beego/router"
	"github.com/raozhaofeng/beego/utils"
	"net/http"
	"net/url"
	"strconv"
)

const (
	AdminUserSupermanId     = 1       //	超级管理员ID
	AdminPrefixTokenKey     = "admin" //	后端TokenKey前缀
	HomePrefixTokenKey      = "home"  //	前端TokenKey前缀
	AdminUserStatusActivate = 10      //	启用状态
	AdminUserStatusDelete   = -2      //	删除状态
)

// AdminUserAttrs 允许只能3级 超管->代理->管理
type AdminUserAttrs struct {
	Id          int64   `json:"id"`           //主键
	ParentId    int64   `json:"parent_id"`    //上级ID
	UserName    string  `json:"username"`     //用户名
	Email       string  `json:"email"`        //邮件
	Nickname    string  `json:"nickname"`     //昵称
	Avatar      string  `json:"avatar"`       //头像
	Password    string  `json:"password"`     //密码
	SecurityKey string  `json:"security_key"` //安全密钥
	Money       float64 `json:"money"`        //金额
	Status      int64   `json:"status"`       //状态 -2删除 -1禁用 10启用
	Data        string  `json:"data"`         //数据
	Extra       string  `json:"extra"`        //额外
	Domain      string  `json:"domain"`       //域名
	ExpiredAt   int64   `json:"expired_at"`   //过期时间
	CreatedAt   int64   `json:"created_at"`   //创建时间
	UpdatedAt   int64   `json:"updated_at"`   //更新时间
}

type AdminUser struct {
	define.Db
}

func NewAdminUser(tx *sql.Tx) *AdminUser {
	return &AdminUser{
		db.Manager.NewInterfaceDb(tx).Table("admin_user"),
	}
}

func (c *AdminUser) FindOne() *AdminUserAttrs {
	attrs := new(AdminUserAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.ParentId, &attrs.UserName, &attrs.Email, &attrs.Nickname, &attrs.Avatar, &attrs.Password,
			&attrs.SecurityKey, &attrs.Money, &attrs.Status, &attrs.Data, &attrs.Extra, &attrs.Domain, &attrs.ExpiredAt, &attrs.CreatedAt, &attrs.UpdatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// GetAdminChildrenParentIds 获取所有下级父级IDS
func (c *AdminUser) GetAdminChildrenParentIds(adminId int64) []string {
	adminIds := []string{strconv.FormatInt(adminId, 10)}
	if adminId == AdminUserSupermanId {
		adminIds = append(adminIds, "0")
	}
	NewAdminUser(nil).
		Field("id").
		AndWhere("parent_id = ?", adminId).
		Query(func(rows *sql.Rows) {
			var tmpAdminId int64
			_ = rows.Scan(&tmpAdminId)
			adminIds = append(adminIds, c.GetAdminChildrenParentIds(tmpAdminId)...)
		})
	return adminIds
}

// FindAdminLikeNameIds 获取管理名称IDS
func (c *AdminUser) FindAdminLikeNameIds(username string) []string {
	data := NewAdminUser(nil).
		Field("id").
		AndWhere("username like ?", "%"+username+"%").ColumnString()
	if len(data) == 0 {
		return []string{"0"}
	}
	return data
}

// GetSettingAdminId 获取配置管理ID
func (c *AdminUser) GetSettingAdminId(adminId int64) int64 {
	if adminId == AdminUserSupermanId || adminId == 0 {
		return adminId
	}

	adminModel := NewAdminUser(nil)
	adminModel.AndWhere("id=?", adminId)
	adminInfo := adminModel.FindOne()
	if adminInfo == nil {
		panic("管理设置ID不存在")
	}
	if adminInfo.ParentId != AdminUserSupermanId {
		return c.GetSettingAdminId(adminInfo.ParentId)
	}
	return adminInfo.Id
}

// GetDomainAdminId 获取域名管理ID
func (c *AdminUser) GetDomainAdminId(r *http.Request) int64 {
	u, _ := url.Parse(r.Header.Get("Origin"))
	var adminId int64

	if u.Host != "" {
		c.Field("id").AndWhere("domain like ?", "%"+u.Host+"%").QueryRow(func(row *sql.Row) {
			_ = row.Scan(&adminId)
		})
	}

	if adminId == 0 {
		return AdminUserSupermanId
	}
	return adminId
}

// GetAdminTokenParams 获取管理【后端｜前端】Token参数
func GetAdminTokenParams() map[string]*router.TokenParams {
	//	添加后台管理员Token参数
	data := map[string]*router.TokenParams{}

	// 后台管理
	NewAdminUser(nil).Field("id", "data").AndWhere("status > ?", AdminUserStatusDelete).Query(func(rows *sql.Rows) {
		var adminId int64
		var tmpData string
		_ = rows.Scan(&adminId, &tmpData)
		dataKey := TokenParamsPrefix(AdminPrefixTokenKey, adminId)

		data[dataKey] = new(router.TokenParams)
		_ = json.Unmarshal([]byte(tmpData), data[dataKey])
	})

	//	前台管理
	NewAdminSetting(nil).Field("admin_id", "value").AndWhere("field=?", "site_token").Query(func(rows *sql.Rows) {
		var adminId int64
		var tmpData string
		_ = rows.Scan(&adminId, &tmpData)
		dataKey := TokenParamsPrefix(HomePrefixTokenKey, adminId)
		data[dataKey] = new(router.TokenParams)
		_ = json.Unmarshal([]byte(tmpData), data[dataKey])
	})
	return data
}

// GetAdminRolesRouter 获取管理角色路由
func GetAdminRolesRouter() map[int64][]string {
	data := map[int64][]string{}

	NewAdminUser(nil).Field("id").Query(func(rows *sql.Rows) {
		var adminId int64
		_ = rows.Scan(&adminId)

		adminRoles := NewAdminAuthAssignment(nil).GetAdminRoleList(adminId)
		data[adminId] = utils.GetMapValues(NewAdminAuthChild(nil).GetRolesRouteList(adminRoles))
	})
	return data
}

// TokenParamsPrefix Token缓存参数名称
func TokenParamsPrefix(prefix string, adminId int64) string {
	adminIdStr := strconv.FormatInt(adminId, 10)
	return utils.PasswordEncrypt(prefix + "_" + adminIdStr)
}
