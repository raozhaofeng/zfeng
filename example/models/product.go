package models

import (
	"database/sql"
	"github.com/raozhaofeng/zfeng/database"
	"github.com/raozhaofeng/zfeng/database/define"
)

const (
	ProductNumsUnlimited       = -1                                                    //	没有限制购买
	ProductModeReturnPrincipal = 1                                                     //	返回本金
	ProductModeReturnInterest  = 2                                                     //	返回利息
	ProductDataJSON            = "{\"expire\": -1, \"interval\": 11, \"returns\": 10}" //	初始值
	ProductStatusActivate      = 10
	ProductStatusDisabled      = -1
	ProductStatusDelete        = -2
	ProductRecommend           = 10
)

// ProductAttrs 数据库模型属性
type ProductAttrs struct {
	Id         int64   `json:"id"`          //主键
	AdminId    int64   `json:"admin_id"`    //管理员ID
	CategoryId int64   `json:"category_id"` //类目ID
	AssetsId   int64   `json:"assets_id"`   //资产ID
	Name       string  `json:"name"`        //标题
	Images     string  `json:"images"`      //图片列表
	Money      float64 `json:"money"`       //金额
	Sort       int64   `json:"sort"`        //排序
	Status     int64   `json:"status"`      //状态 -2删除 -1禁用 10启用
	Recommend  int64   `json:"recommend"`   //推荐 -1关闭 10推荐
	Sales      int64   `json:"sales"`       //销售量
	Nums       int64   `json:"nums"`        //限购 -1无限
	Mode       int64   `json:"mode"`        //模式 1返本 2返息
	Used       int64   `json:"used"`        //已使用
	Total      int64   `json:"total"`       //总数
	Data       string  `json:"data"`        //数据
	Describes  string  `json:"describes"`   //数据
	UpdatedAt  int64   `json:"updated_at"`  //更新时间
	CreatedAt  int64   `json:"created_at"`  //创建时间
}

type ProductDataAttrs struct {
	Expire   int64   `json:"expire"`   //	过期时间（天）
	Interval int64   `json:"interval"` // 	间隔时间(小时)
	Returns  float64 `json:"returns"`  // 	回报率(%)
}

// Product 数据库模型
type Product struct {
	define.Db
}

// NewProduct 创建数据库模型
func NewProduct(tx *sql.Tx) *Product {
	return &Product{
		database.DbPool.NewDb(tx).Table("product"),
	}
}

// FindOne 查询单挑
func (c *Product) FindOne() *ProductAttrs {
	attrs := new(ProductAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.AdminId, &attrs.CategoryId, &attrs.AssetsId, &attrs.Name, &attrs.Images, &attrs.Money, &attrs.Sort, &attrs.Status, &attrs.Recommend, &attrs.Sales, &attrs.Nums, &attrs.Mode, &attrs.Used, &attrs.Total, &attrs.Data, &attrs.Describes, &attrs.UpdatedAt, &attrs.CreatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *Product) FindMany() []*ProductAttrs {
	data := make([]*ProductAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(ProductAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.CategoryId, &tmp.AssetsId, &tmp.Name, &tmp.Images, &tmp.Money, &tmp.Sort, &tmp.Status, &tmp.Recommend, &tmp.Sales, &tmp.Nums, &tmp.Mode, &tmp.Used, &tmp.Total, &tmp.Data, &tmp.Describes, &tmp.UpdatedAt, &tmp.CreatedAt)
		data = append(data, tmp)
	})
	return data
}
