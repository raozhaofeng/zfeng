package account

import (
	"basic/models"
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/router"
	"github.com/raozhaofeng/zfeng/utils/body"
	"github.com/raozhaofeng/zfeng/validator"
	"net/http"
	"time"
)

type createParams struct {
	UserName   string `json:"username" validate:"required"`
	PaymentId  int64  `json:"payment_id" validate:"required"`
	Name       string `json:"name" validate:"required"`
	RealName   string `json:"real_name" validate:"required"`
	CardNumber string `json:"card_number" validate:"required"`
	Address    string `json:"address"`
	Data       string `json:"data"`
}

func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(createParams)
	_ = body.ReadJSON(r, params)

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	paymentModel := models.NewWalletPayment(nil)
	paymentModel.AndWhere("id=?", params.PaymentId)
	paymentInfo := paymentModel.FindOne()
	if paymentInfo == nil {
		body.ErrorJSON(w, "提现方式不存在", -1)
		return
	}

	userModel := models.NewUser(nil)
	userModel.AndWhere("username=?", params.UserName)
	userInfo := userModel.FindOne()
	if userInfo == nil {
		body.ErrorJSON(w, "用户不存在", -1)
		return
	}
	//  获取管理员ID
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	if adminId != models.AdminUserSupermanId && adminId != userInfo.AdminId {
		body.ErrorJSON(w, "权限不足", -1)
		return
	}

	nowTime := time.Now()
	_, err = models.NewUserWalletAccount(nil).
		Field("admin_id", "user_id", "payment_id", "name", "real_name", "card_number", "address", "data", "updated_at", "created_at").
		Args(userInfo.AdminId, userInfo.Id, params.PaymentId, params.Name, params.RealName, params.CardNumber, params.Address, params.Data, nowTime.Unix(), nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}

	body.SuccessJSON(w, "ok")
}
