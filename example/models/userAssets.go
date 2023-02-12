package models

import (
	"database/sql"
	"github.com/raozhaofeng/zfeng/database"
	"github.com/raozhaofeng/zfeng/database/define"
)

const (
	UserAssetsStatusActivate = 10
	UserAssetsStatusDelete   = -2
)

// UserAssetsAttrs 数据库模型属性
type UserAssetsAttrs struct {
	Id          int64   `json:"id"`           //主键
	AdminId     int64   `json:"admin_id"`     //管理员ID
	UserId      int64   `json:"user_id"`      //用户ID
	AssetsId    int64   `json:"assets_id"`    //资产ID
	Name        string  `json:"name"`         //名称
	Money       float64 `json:"money"`        //金额
	FreezeMoney float64 `json:"freeze_money"` //冻结金额
	Data        string  `json:"data"`         //数据
	Status      int64   `json:"status"`       //状态 -2删除｜-1禁用｜10启用
	CreatedAt   int64   `json:"created_at"`   //创建时间
	UpdatedAt   int64   `json:"updated_at"`   //更新时间
}

// UserAssets 数据库模型
type UserAssets struct {
	define.Db
}

// NewUserAssets 创建数据库模型
func NewUserAssets(tx *sql.Tx) *UserAssets {
	return &UserAssets{
		database.DbPool.NewDb(tx).Table("user_assets"),
	}
}

// FindOne 查询单挑
func (c *UserAssets) FindOne() *UserAssetsAttrs {
	attrs := new(UserAssetsAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.AdminId, &attrs.UserId, &attrs.AssetsId, &attrs.Name, &attrs.Money, &attrs.FreezeMoney, &attrs.Data, &attrs.Status, &attrs.CreatedAt, &attrs.UpdatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *UserAssets) FindMany() []*UserAssetsAttrs {
	data := make([]*UserAssetsAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(UserAssetsAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.UserId, &tmp.AssetsId, &tmp.Name, &tmp.Money, &tmp.FreezeMoney, &tmp.Data, &tmp.Status, &tmp.CreatedAt, &tmp.UpdatedAt)
		data = append(data, tmp)
	})
	return data
}
