package models

import (
	"database/sql"
	"github.com/raozhaofeng/zfeng/database"
	"github.com/raozhaofeng/zfeng/database/define"
	"time"
)

const (
	UserBillTypeSystemDeposit         int64 = 1  //	系统充值
	UserBillTypeSystemDeduction       int64 = 2  //	系统扣除
	UserBillTypeDeposit               int64 = 3  //	用户充值
	UserBillTypeWithdraw              int64 = 4  //	用户提现
	UserBillTypeWithdrawRefuse        int64 = 5  //	提现拒绝
	UserBillTypeBuyLevel              int64 = 10 //	购买等级
	UserBillTypeBuyUpgradeLevel       int64 = 11 //	升级等级
	UserBillTypeRegisterRewards       int64 = 15 //	注册奖励
	UserBillTypeTaskRewards           int64 = 16 //	任务奖励
	UserBillTypeExtraRewards          int64 = 17 //	额外奖励
	UserBillTypeInviteRewards         int64 = 18 //	邀请奖励
	UserBillTypeBuyProduct            int64 = 20 //	购买产品
	UserBillTypeReturnProductAmount   int64 = 21 //	退回本金
	UserBillTypeProductProfit         int64 = 22 //	产品利润
	UserBillTypeBuyProductEarnings    int64 = 30 //	分销购买产品收益
	UserBillTypeProductProfitEarnings int64 = 31 //	分销产品利润收益
	UserBillTypeTaskEarnings          int64 = 32 //	分销任务收益

	UserBillTypeSystemAssetsDeposit         int64 = 101 //	系统资产充值
	UserBillTypeSystemAssetsDeduction       int64 = 102 //	系统资产扣除
	UserBillTypeAssetsDeposit               int64 = 103 //	资产充值
	UserBillTypeAssetsWithdraw              int64 = 104 //	资产提现
	UserBillTypeAssetsWithdrawRefuse        int64 = 105 //	资产提现拒绝
	UserBillTypeAssetsTaskRewards           int64 = 116 //	资产任务奖励
	UserBillTypeAssetsExtraRewards          int64 = 117 //	资产额外奖励
	UserBillTypeAssetsBuyProduct            int64 = 120 //	资产购买产品
	UserBillTypeAssetsReturnProductAmount   int64 = 121 //	资产退回本金
	UserBillTypeAssetsProductProfit         int64 = 122 //	资产产品利润
	UserBillTypeAssetsBuyProductEarnings    int64 = 130 //	分销资产购买产品收益
	UserBillTypeAssetsProductProfitEarnings int64 = 131 //	分销资产产品利润收益
	UserBillTypeAssetsTaskEarnings          int64 = 132 //	分销资产任务收益
)

// UserBillTypeNameMap 语言字典名称
var UserBillTypeNameMap = map[int64]string{
	UserBillTypeSystemDeposit: "systemDeposit", UserBillTypeSystemDeduction: "systemDeduction", UserBillTypeDeposit: "deposit",
	UserBillTypeWithdraw: "withdraw", UserBillTypeWithdrawRefuse: "withdrawRefuse", UserBillTypeBuyLevel: "buyLevel",
	UserBillTypeBuyUpgradeLevel: "buyUpgradeLevel", UserBillTypeRegisterRewards: "registerRewards", UserBillTypeTaskRewards: "taskRewards",
	UserBillTypeExtraRewards: "extraRewards", UserBillTypeInviteRewards: "inviteRewards", UserBillTypeBuyProduct: "buyProduct", UserBillTypeReturnProductAmount: "returnProductAmount",
	UserBillTypeProductProfit: "productProfit", UserBillTypeBuyProductEarnings: "buyProductEarnings", UserBillTypeProductProfitEarnings: "productProfitEarnings",
	UserBillTypeTaskEarnings: "taskEarnings",

	UserBillTypeSystemAssetsDeposit: "systemAssetsDeposit", UserBillTypeSystemAssetsDeduction: "systemAssetsDeduction",
	UserBillTypeAssetsDeposit: "assetsDeposit", UserBillTypeAssetsWithdraw: "assetsWithdraw", UserBillTypeAssetsWithdrawRefuse: "assetsWithdrawRefuse",
	UserBillTypeAssetsTaskRewards: "assetsTaskRewards", UserBillTypeAssetsExtraRewards: "assetsExtraRewards",
	UserBillTypeAssetsBuyProduct: "assetsBuyProduct", UserBillTypeAssetsReturnProductAmount: "assetsReturnProductAmount",
	UserBillTypeAssetsProductProfit: "assetsProductProfit", UserBillTypeAssetsBuyProductEarnings: "assetsBuyProductEarnings",
	UserBillTypeAssetsProductProfitEarnings: "assetsProductProfitEarnings", UserBillTypeAssetsTaskEarnings: "assetsTaskEarnings",
}

type UserBillAttrs struct {
	Id        int64   `json:"id"`         //主键
	AdminId   int64   `json:"admin_id"`   //管理员ID
	UserId    int64   `json:"user_id"`    //用户ID
	SourceId  int64   `json:"source_id"`  //来源ID
	Name      string  `json:"name"`       //标题
	Type      int64   `json:"type"`       //类型
	Balance   float64 `json:"balance"`    //余额
	Money     float64 `json:"money"`      //金额
	Data      string  `json:"data"`       //数据
	CreatedAt int64   `json:"created_at"` //创建时间
}

type UserBill struct {
	define.Db
}

func NewUserBill(tx *sql.Tx) *UserBill {
	return &UserBill{
		database.DbPool.NewDb(tx).Table("user_bill"),
	}
}

func (c *UserBill) FindOne() *UserBillAttrs {
	attrs := new(UserBillAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.AdminId, &attrs.UserId, &attrs.SourceId, &attrs.Name, &attrs.Type, &attrs.Balance, &attrs.Money, &attrs.Data, &attrs.CreatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// WriteUserBill 写入账单
func (c *UserBill) WriteUserBill(adminId, userId, sourceId, billType int64, beforeMoney, money float64) {
	nowTime := time.Now()
	_, err := c.Field("admin_id", "user_id", "source_id", "name", "type", "balance", "money", "created_at").
		Args(adminId, userId, sourceId, UserBillTypeNameMap[billType], billType, beforeMoney, money, nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}
}

// GetBillFundingChangesMoney 获取账单变动金额
func GetBillFundingChangesMoney(billType int64, beforeMoney, money float64) float64 {
	var currentMoney float64
	switch billType {
	case UserBillTypeSystemDeduction, UserBillTypeWithdraw, UserBillTypeBuyLevel, UserBillTypeBuyUpgradeLevel,
		UserBillTypeBuyProduct:
		currentMoney = beforeMoney - money
	default:
		currentMoney = beforeMoney + money
	}
	return currentMoney
}

// GetBillFundingChangesEarningsType 获取资金变动分销类型
func GetBillFundingChangesEarningsType(billType int64) (string, int64) {
	switch billType {
	case UserBillTypeTaskRewards:
		return "taskProfit", UserBillTypeTaskEarnings
	case UserBillTypeBuyProduct:
		return "buyProduct", UserBillTypeBuyProductEarnings
	case UserBillTypeProductProfit:
		return "productProfit", UserBillTypeProductProfitEarnings
	case UserBillTypeAssetsTaskRewards:
		return "taskProfit", UserBillTypeAssetsTaskEarnings
	case UserBillTypeAssetsBuyProduct:
		return "buyProduct", UserBillTypeAssetsBuyProductEarnings
	case UserBillTypeAssetsProductProfit:
		return "productProfit", UserBillTypeAssetsProductProfitEarnings
	}
	return "", 0
}
