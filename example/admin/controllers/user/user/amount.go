package user

import (
	"basic/models"
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/database"
	"github.com/raozhaofeng/zfeng/router"
	"github.com/raozhaofeng/zfeng/utils"
	"github.com/raozhaofeng/zfeng/utils/body"
	"github.com/raozhaofeng/zfeng/validator"
	"net/http"
	"strings"
	"time"
)

type amountParams struct {
	UserName string  `json:"username" validate:"required"`
	Type     int64   `json:"type" validate:"required,oneof=1 2"`
	Money    float64 `json:"money" validate:"required,number"`
}

// Amount 修改用户余额
func Amount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(amountParams)
	_ = body.ReadJSON(r, params)
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	adminId := router.TokenManager.GetContextClaims(r).AdminId
	userMode := models.NewUser(nil)
	userMode.AndWhere("username=?", params.UserName)
	if adminId != models.AdminUserSupermanId {
		adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
		userMode.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	}
	userInfo := userMode.FindOne()
	if userInfo == nil {
		body.ErrorJSON(w, "用户不存在", -1)
		return
	}

	tx := database.DbPool.GetTx()
	defer tx.Rollback()

	err = models.UserFundingChanges(tx, userInfo.AdminId, userInfo.Id, userInfo.ParentId, nil, 0, params.Type, userInfo.Money, params.Money)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	// 写入钱包订单
	nowTime := time.Now()
	walletOrderType := models.WalletOrderTypeSystemDeposit
	if params.Type == models.UserBillTypeSystemDeduction {
		walletOrderType = models.WalletOrderTypeSystemWithdraw
	}
	_, err = models.NewUserWalletOrder(nil).Field("order_sn", "admin_id", "user_id", "user_type", "type", "money", "balance", "status", "updated_at", "created_at").
		Args(utils.NewRandom().OrderSn(), userInfo.AdminId, userInfo.Id, userInfo.Type, walletOrderType, params.Money, userInfo.Money, models.WalletOrderStatusComplete, nowTime.Unix(), nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}

	_ = tx.Commit()
	body.SuccessJSON(w, "ok")
}
