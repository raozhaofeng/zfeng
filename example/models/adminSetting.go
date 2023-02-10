package models

import (
	"database/sql"
	"encoding/json"
	"github.com/raozhaofeng/zfeng/database"
	"github.com/raozhaofeng/zfeng/database/define"
)

const (
	SettingTypeText                          = "text"
	SettingTypeNumber                        = "number"
	SettingTypeEditor                        = "editor"
	SettingTypeImage                         = "image"
	SettingTypeImages                        = "images"
	SettingTypeSelect                        = "select"
	SettingTypeCheckbox                      = "checkbox"
	SettingTypeChildren                      = "children"
	SettingTypeJson                          = "json"
	SettingGroupBasic                        = 1             //	基本设置
	SettingGroupHome                         = 2             //	首页设置
	SettingGroupFinance                      = 3             //	财务设置
	SettingGroupTemplate                     = 4             //	模版设置
	SettingGroupHelpers                      = 5             //	帮助中心
	UpdateAdminTokenParamsField              = "site_token"  //	前端Token健铭
	AdminSettingBuyLevelModePremium          = "premium"     //	补价模式
	AdminSettingBuyLevelModeEquivalence      = "equivalence" //	等价模式
	AdminSettingProductEarningsModeManual    = "manual"      //	产品收益模式【手动】
	AdminSettingProductEarningsModeAutomatic = "automatic"   //	产品收益模式【自动】
)

// AdminSettingAttrs 数据库模型属性
type AdminSettingAttrs struct {
	Id      int64  `json:"id"`       //主键
	AdminId int64  `json:"admin_id"` //管理员ID
	GroupId int64  `json:"group_id"` //组ID
	Name    string `json:"name"`     //名称
	Type    string `json:"type"`     //类型
	Field   string `json:"field"`    //健名
	Value   string `json:"value"`    //健值
	Data    string `json:"data"`     //数据
}

// AdminSetting 数据库模型
type AdminSetting struct {
	define.Db
}

// NewAdminSetting 创建数据库模型
func NewAdminSetting(tx *sql.Tx) *AdminSetting {
	return &AdminSetting{
		database.DbPool.NewDb(tx).Table("admin_setting"),
	}
}

// FindOne 查询单条
func (c *AdminSetting) FindOne() *AdminSettingAttrs {
	attrs := new(AdminSettingAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.AdminId, &attrs.GroupId, &attrs.Name, &attrs.Type, &attrs.Field, &attrs.Value, &attrs.Data)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *AdminSetting) FindMany() []*AdminSettingAttrs {
	var data []*AdminSettingAttrs
	c.Query(func(rows *sql.Rows) {
		tmp := new(AdminSettingAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.GroupId, &tmp.Name, &tmp.Type, &tmp.Field, &tmp.Value, &tmp.Data)
		data = append(data, tmp)
	})
	return data
}

// GetAdminFieldString 获取索引值
func (c *AdminSetting) GetAdminFieldString(adminId int64, field string) string {
	data := ""
	c.Field("value").AndWhere("admin_id=?", adminId).AndWhere("field=?", field).QueryRow(func(row *sql.Row) {
		_ = row.Scan(&data)
	})
	return data
}

// GetAdminFieldAllString 获取管理员所有配置信息
func (c *AdminSetting) GetAdminFieldAllString(adminId int64) map[string]string {
	data := map[string]string{}
	c.Field("field", "value").AndWhere("admin_id=?", adminId).Query(func(rows *sql.Rows) {
		var fieldTmp, valueTmp string
		_ = rows.Scan(&fieldTmp, &valueTmp)

		data[fieldTmp] = valueTmp
	})
	return data
}

// GetAdminEarningsMode 获取所有管理收益模式
func (c *AdminSetting) GetAdminEarningsMode() map[int64]string {
	data := map[int64]string{}
	c.Field("admin_id", "value").AndWhere("field=?", "product_earnings_mode").Query(func(rows *sql.Rows) {
		var adminId int64
		var value string
		_ = rows.Scan(&adminId, &value)
		data[adminId] = value
	})
	return data
}

// AdminSettingValueToMapInterface 管理配置内容转 Map
func AdminSettingValueToMapInterface(val string) map[string]any {
	data := map[string]any{}
	_ = json.Unmarshal([]byte(val), &data)
	return data
}

// AdminSettingValueToMapInterfaces 管理员配置内容转 Map数组
func AdminSettingValueToMapInterfaces(val string) []map[string]any {
	data := make([]map[string]any, 0)
	_ = json.Unmarshal([]byte(val), &data)
	return data
}
