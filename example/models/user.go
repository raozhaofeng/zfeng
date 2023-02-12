package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/raozhaofeng/zfeng/database"
	"github.com/raozhaofeng/zfeng/database/define"
	"strings"
	"time"
)

const (
	UserStatusDelete   = -2 //	用户删除状态
	UserStatusDisabled = -1 //	用户禁用状态
	UserStatusFreeze   = 1  //	用户冻结状态
	UserStatusActivate = 10 //	激活用户状态
	UserTypeVirtual    = -1 //	虚拟用户
	UserTypeReality    = 10 //	真实用户
)

// UserAttrs 数据库模型属性
type UserAttrs struct {
	Id          int64   `json:"id"`           //主键
	AdminId     int64   `json:"admin_id"`     //管理员ID
	ParentId    int64   `json:"parent_id"`    //上级ID
	CountryId   int64   `json:"country_id"`   //国家ID
	UserName    string  `json:"username"`     //用户名
	Nickname    string  `json:"nickname"`     //昵称
	Email       string  `json:"email"`        //邮箱
	Telephone   string  `json:"telephone"`    //手机号码
	Avatar      string  `json:"avatar"`       //头像
	Sex         int64   `json:"sex"`          //类型 -1未知 1男 2女
	Birthday    int64   `json:"birthday"`     //生日
	Password    string  `json:"password"`     //密码
	SecurityKey string  `json:"security_key"` //安全密钥
	Money       float64 `json:"money"`        //金额
	FreezeMoney float64 `json:"freeze_money"` //冻结金额
	Type        int64   `json:"type"`         //类型 -1虚拟 10普通
	Status      int64   `json:"status"`       //状态 -2删除｜-1禁用｜10启用
	Data        string  `json:"data"`         //数据
	Ip4         string  `json:"ip4"`          //IP4地址
	CreatedAt   int64   `json:"created_at"`   //创建时间
	UpdatedAt   int64   `json:"updated_at"`   //更新时间
}

type UserVerifyInfo struct {
	Status int64  `json:"status"` //	验证状态
	Data   string `json:"data"`   //	信息
}

type UserLevelInfo struct {
	Id        int64   `json:"id"`        //	Id
	Level     int64   `json:"level"`     //	等级
	Name      string  `json:"name"`      //	名称
	Icon      string  `json:"icon"`      //	图标
	Days      int64   `json:"days"`      //	天数
	Money     float64 `json:"money"`     //	金额
	CreatedAt int64   `json:"createdAt"` //	创建时间戳
	UpdatedAt int64   `json:"updatedAt"` //	过期时间戳
}

type UserInfo struct {
	Id          int64           `json:"id"`           //	用户ID
	CountryId   int64           `json:"countryId"`    //	国家ID
	UserName    string          `json:"username"`     //	用户名
	Nickname    string          `json:"nickname"`     //	昵称
	Email       string          `json:"email"`        //	邮箱
	Telephone   string          `json:"telephone"`    //	手机号
	Avatar      string          `json:"avatar"`       //	头像
	Sex         int64           `json:"sex"`          //	类型 -1未知 1男 2女
	Birthday    int64           `json:"birthday"`     //	生日
	Money       float64         `json:"money"`        //	金额
	FreezeMoney float64         `json:"freeze_money"` //	冻结金额
	Data        string          `json:"data"`         //	数据
	VerifyInfo  *UserVerifyInfo `json:"verifyInfo"`   //	是否验证 -1验证失败 0未验证 10 正在验证 20已验证
	LevelInfo   *UserLevelInfo  `json:"levelInfo"`    //	等级信息
	InviteCode  string          `json:"inviteCode"`   //	邀请码
	CreatedAt   int64           `json:"createdAt"`    //	创建时间
	UpdatedAt   int64           `json:"updatedAt"`    //	更新时间
}

// UserTree 用户树
type UserTree struct {
	Id          int64       `json:"id"`           //	用户ID
	Header      string      `json:"header"`       //	层级
	Avatar      string      `json:"avatar"`       //	头像
	UserName    string      `json:"username"`     //	用户名
	SumPeople   int64       `json:"sum_people"`   //	总人数
	SumAmount   float64     `json:"sum_amount"`   //	总购买
	SumEarnings float64     `json:"sum_earnings"` //	总收益
	Children    []*UserTree `json:"children"`     //	子集
}

// User 数据库模型
type User struct {
	define.Db
}

// NewUser 创建数据库模型
func NewUser(tx *sql.Tx) *User {
	return &User{
		database.DbPool.NewDb(tx).Table("user"),
	}
}

// FindOne 查询单挑
func (c *User) FindOne() *UserAttrs {
	attrs := new(UserAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.AdminId, &attrs.ParentId, &attrs.CountryId, &attrs.UserName, &attrs.Nickname, &attrs.Email, &attrs.Telephone, &attrs.Avatar, &attrs.Sex, &attrs.Birthday, &attrs.Password, &attrs.SecurityKey, &attrs.Money, &attrs.FreezeMoney, &attrs.Type, &attrs.Status, &attrs.Data, &attrs.Ip4, &attrs.CreatedAt, &attrs.UpdatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *User) FindMany() []*UserAttrs {
	data := make([]*UserAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(UserAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.ParentId, &tmp.CountryId, &tmp.UserName, &tmp.Nickname, &tmp.Email, &tmp.Telephone, &tmp.Avatar, &tmp.Sex, &tmp.Birthday, &tmp.Password, &tmp.SecurityKey, &tmp.Money, &tmp.FreezeMoney, &tmp.Type, &tmp.Status, &tmp.Data, &tmp.Ip4, &tmp.CreatedAt, &tmp.UpdatedAt)
		data = append(data, tmp)
	})
	return data
}

// FindUserLikeNameIds 获取用户名称IDS
func (c *User) FindUserLikeNameIds(username string) []string {
	data := NewUser(nil).
		Field("id").
		AndWhere("username like ?", "%"+username+"%").ColumnString()
	if len(data) == 0 {
		return []string{"-1"}
	}
	return data
}

// GetUserTree 获取用户树
func (c *User) GetUserTree(adminId int64, parentId int64) []*UserTree {
	adminIds := NewAdminUser(nil).GetAdminChildrenParentIds(adminId)

	data := make([]*UserTree, 0)
	NewUser(nil).Field("id", "avatar", "username").
		AndWhere("admin_id in ("+strings.Join(adminIds, ",")+")").
		AndWhere("status>?", UserStatusDelete).AndWhere("parent_id=?", parentId).Query(func(rows *sql.Rows) {
		temp := new(UserTree)
		_ = rows.Scan(&temp.Id, &temp.Avatar, &temp.UserName)
		if parentId == 0 {
			temp.Header = "root"
		} else {
			temp.Header = "generic"
		}

		// 总人数
		NewUser(nil).Field("count(*)").
			AndWhere("parent_id=?", temp.Id).
			AndWhere("status>?", UserStatusDelete).QueryRow(func(row *sql.Row) {
			_ = row.Scan(&temp.SumPeople)
		})

		// 总购买
		NewUserBill(nil).Field("sum(money)").AndWhere("user_id=?", temp.Id).AndWhere("type=?", UserBillTypeBuyProduct).QueryRow(func(row *sql.Row) {
			_ = row.Scan(&temp.SumAmount)
		})
		// 总收益
		NewUserBill(nil).Field("sum(money)").AndWhere("user_id=?", temp.Id).AndWhere("type=?", UserBillTypeProductProfit).QueryRow(func(row *sql.Row) {
			_ = row.Scan(&temp.SumEarnings)
		})

		// 子集用户
		temp.Children = NewUser(nil).GetUserTree(adminId, temp.Id)
		data = append(data, temp)
	})

	return data
}

// UserFundingChanges 用户资金变动
// productId 如果需要增加资产金额
func UserFundingChanges(tx *sql.Tx, adminId, userId, userParentId int64, assetsInfo *AssetsAttrs, sourceId, billType int64, beforeMoney, money float64) error {
	if money <= 0 {
		return errors.New("IncorrectAmount")
	}

	if assetsInfo != nil {
		userAssetsModel := NewUserAssets(tx)
		userAssetsModel.AndWhere("user_id=?", userId).AndWhere("assets_id=?", assetsInfo.Id)
		userAssetsInfo := userAssetsModel.FindOne()

		// 更新用户资产金额
		if userAssetsInfo == nil {
			nowTime := time.Now()
			_, err := NewUserAssets(tx).Field("admin_id", "user_id", "assets_id", "name", "money", "created_at", "updated_at").
				Args(adminId, userId, assetsInfo.Id, assetsInfo.Name, money, nowTime.Unix(), nowTime.Unix()).
				Insert()
			if err != nil {
				panic(err)
			}
		} else {
			// 获取资金变动金额
			beforeMoney = userAssetsInfo.Money
			currentUserMoney := GetBillFundingChangesMoney(billType, userAssetsInfo.Money, money)
			if currentUserMoney <= 0 {
				return errors.New("insufficientBalance")
			}
			_, err := NewUserAssets(tx).Value("money=?").Args(currentUserMoney).
				AndWhere("user_id=?", userId).AndWhere("assets_id=?", assetsInfo.Id).
				Update()
			if err != nil {
				return err
			}
		}
	} else {
		// 获取资金变动金额
		currentUserMoney := GetBillFundingChangesMoney(billType, beforeMoney, money)
		if currentUserMoney <= 0 {
			return errors.New("insufficientBalance")
		}
		_, err := NewUser(tx).Value("money=?").Args(currentUserMoney).AndWhere("id=?", userId).Update()
		if err != nil {
			return err
		}
	}

	// 是否需要分销操作
	UserFundingEarnings(tx, adminId, userId, userParentId, billType, assetsInfo, money)

	// 写入账单
	NewUserBill(tx).WriteUserBill(adminId, userId, sourceId, billType, beforeMoney, money)
	return nil
}

// UserFundingEarnings 用户资金变动上级收益
func UserFundingEarnings(tx *sql.Tx, adminId, userId, userParentId, billType int64, assetsInfo *AssetsAttrs, money float64) {
	if userParentId <= 0 {
		return
	}

	isPyramidKey, newBillType := GetBillFundingChangesEarningsType(billType)
	settingAdminId := NewAdminUser(nil).GetSettingAdminId(adminId)
	adminSettingList := NewAdminSetting(nil).GetAdminFieldAllString(settingAdminId)
	if adminSettingList["pyramid_level"] == "" || adminSettingList["pyramid_items"] == "" {
		return
	}
	// 收益是否开启
	var pyramidItems map[string]bool
	_ = json.Unmarshal([]byte(adminSettingList["pyramid_items"]), &pyramidItems)
	if !pyramidItems[isPyramidKey] {
		return
	}

	// 等级收益比例
	var pyramidLevel []map[string]float64
	_ = json.Unmarshal([]byte(adminSettingList["pyramid_level"]), &pyramidLevel)

	// 分销者收益
	currentParentId := userParentId
	for i := 0; i < len(pyramidLevel); i++ {
		if pyramidLevel[i]["value"] <= 0 {
			continue
		}
		// 上级收益
		currentUserModel := NewUser(nil)
		currentUserModel.AndWhere("id=?", currentParentId)
		currentUserInfo := currentUserModel.FindOne()
		pyramidLevelProfitMoney := money * pyramidLevel[i]["value"] / 100
		_ = UserFundingChanges(tx, currentUserInfo.AdminId, currentUserInfo.Id, currentUserInfo.ParentId, assetsInfo, userId, newBillType, currentUserInfo.Money, pyramidLevelProfitMoney)
	}
}
