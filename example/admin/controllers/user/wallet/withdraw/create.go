package withdraw

import (
	"basic/models"
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/database"
	"github.com/raozhaofeng/zfeng/router"
	"github.com/raozhaofeng/zfeng/utils"
	"github.com/raozhaofeng/zfeng/utils/body"
	"github.com/raozhaofeng/zfeng/validator"
	"net/http"
	"time"
)

type createParams struct {
	UserName string  `json:"username"`
	Money    float64 `json:"money"`
}

func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(createParams)
	_ = body.ReadJSON(r, params)

	//  验证参数
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	// 查询用户是否存在
	userModel := models.NewUser(nil)
	userModel.AndWhere("username=?", params.UserName)
	userInfo := userModel.FindOne()
	if userInfo == nil {
		body.ErrorJSON(w, "用户不存在", -1)
		return
	}

	// 判断用户金额是否足够
	if userInfo.Money < params.Money {
		body.ErrorJSON(w, "用户金额不足", -1)
		return
	}

	// 是否有支付方式
	paymentModel := models.NewUserWalletAccount(nil)
	paymentModel.AndWhere("user_id=?", userInfo.Id)
	paymentInfo := paymentModel.FindOne()
	if paymentInfo == nil {
		body.ErrorJSON(w, "没有提现方式", -1)
		return
	}

	//  获取管理员ID
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	if adminId != models.AdminUserSupermanId && adminId != userInfo.AdminId {
		body.ErrorJSON(w, "权限不足", -1)
		return
	}

	tx := database.DbPool.GetTx()
	defer tx.Rollback()

	nowTime := time.Now()
	orderSn := utils.NewRandom().OrderSn()
	withdrawFee := models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "finance_withdraw_fee")
	_, err = models.NewUserWalletOrder(tx).
		Field("order_sn", "admin_id", "user_id", "user_type", "type", "payment_id", "money", "fee", "data", "updated_at", "created_at").
		Args(orderSn, userInfo.AdminId, userInfo.Id, userInfo.Type, models.WalletOrderTypeWithdraw, paymentInfo.Id, params.Money, withdrawFee, "", nowTime.Unix(), nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}

	// 用户资金变动
	err = models.UserFundingChanges(tx, userInfo.AdminId, userInfo.Id, userInfo.ParentId, nil, 0, models.UserBillTypeWithdraw, userInfo.Money, params.Money)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	_ = tx.Commit()
	body.SuccessJSON(w, "ok")
}
