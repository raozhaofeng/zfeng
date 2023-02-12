package models

import (
	"database/sql"
	"github.com/raozhaofeng/zfeng/database"
	"github.com/raozhaofeng/zfeng/database/define"
)

const (
	AssetsStatusDelete   = -2
	AssetsStatusDisabled = -1
	AssetsStatusActivate = 10
	AssetsTypeETH        = 1
	AssetsTypeBSC        = 2
	AssetsTypeTRX        = 3
)

// AssetsTypeList 资产列表
var AssetsTypeList = map[int64]string{
	AssetsTypeETH: "ETH", AssetsTypeBSC: "BSC", AssetsTypeTRX: "TRX",
}

// AssetsAttrs 数据库模型属性
type AssetsAttrs struct {
	Id        int64  `json:"id"`         //主键
	AdminId   int64  `json:"admin_id"`   //管理员ID
	Name      string `json:"name"`       //名称
	Icon      string `json:"icon"`       //图标
	Type      int64  `json:"type"`       //类型 1ETH 2BSC 3TRX
	Data      string `json:"data"`       //数据
	Status    int64  `json:"status"`     //状态 -2删除｜-1禁用｜10启用
	CreatedAt int64  `json:"created_at"` //创建时间
	UpdatedAt int64  `json:"updated_at"` //更新时间
}

// Assets 数据库模型
type Assets struct {
	define.Db
}

// NewAssets 创建数据库模型
func NewAssets(tx *sql.Tx) *Assets {
	return &Assets{
		database.DbPool.NewDb(tx).Table("assets"),
	}
}

// FindOne 查询单挑
func (c *Assets) FindOne() *AssetsAttrs {
	attrs := new(AssetsAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.AdminId, &attrs.Name, &attrs.Icon, &attrs.Type, &attrs.Data, &attrs.Status, &attrs.CreatedAt, &attrs.UpdatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *Assets) FindMany() []*AssetsAttrs {
	data := make([]*AssetsAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(AssetsAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.Name, &tmp.Icon, &tmp.Type, &tmp.Data, &tmp.Status, &tmp.CreatedAt, &tmp.UpdatedAt)
		data = append(data, tmp)
	})
	return data
}
