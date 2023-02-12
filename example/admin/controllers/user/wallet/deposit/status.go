package deposit

import (
	"basic/models"
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/database"
	"github.com/raozhaofeng/zfeng/database/define"
	"github.com/raozhaofeng/zfeng/router"
	"github.com/raozhaofeng/zfeng/utils/body"
	"github.com/raozhaofeng/zfeng/validator"
	"net/http"
	"strings"
	"time"
)

type statusParams struct {
	Id     int64  `json:"id" validate:"required"`
	Status int64  `json:"status" validate:"required,oneof=-1 20"`
	Data   string `json:"data"`
}

func Status(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(statusParams)
	_ = body.ReadJSON(r, params)

	//  参数验证
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	// 是否存在订单
	walletOrderModel := models.NewUserWalletOrder(nil)
	walletOrderModel.AndWhere("id=?", params.Id).AndWhere("status=?", models.WalletOrderStatusPending)
	walletOrderInfo := walletOrderModel.FindOne()
	if walletOrderInfo == nil {
		body.ErrorJSON(w, "订单不存在", -1)
		return
	}

	userModel := models.NewUser(nil)
	userModel.AndWhere("id=?", walletOrderInfo.UserId)
	userInfo := userModel.FindOne()
	if userInfo == nil {
		body.ErrorJSON(w, "用户不存在", -1)
		return
	}

	//  实例化模型
	tx := database.DbPool.GetTx()
	defer tx.Rollback()

	model := models.NewUserWalletOrder(tx)
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	nowTime := time.Now()
	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		String("data=?", params.Data).
		Int64("status=?", params.Status).
		Float64("balance=?", userInfo.Money).
		Int64("updated_at=?", nowTime.Unix())

	if adminId != models.AdminUserSupermanId {
		adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
		model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	}
	_, err = model.AndWhere("id=?", params.Id).AndWhere("type=?", models.WalletOrderTypeDeposit).Update()
	if err != nil {
		panic(err)
	}

	// 如果是成功，那么新增余额
	if params.Status == models.WalletOrderStatusComplete {
		err = models.UserFundingChanges(tx, userInfo.AdminId, userInfo.Id, userInfo.ParentId, nil, 0, models.UserBillTypeDeposit, userInfo.Money, walletOrderInfo.Money)
	}

	_ = tx.Commit()
	body.SuccessJSON(w, "ok")
}
