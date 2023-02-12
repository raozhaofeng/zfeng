package account

import (
	"basic/models"
	"database/sql"
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/database/define"
	"github.com/raozhaofeng/zfeng/router"
	"github.com/raozhaofeng/zfeng/utils/body"
	"net/http"
	"strings"
	"time"
)

type indexParams struct {
	AdminName   string                 `json:"admin_name"`
	UserName    string                 `json:"username"`
	PaymentName string                 `json:"payment_name"`
	Name        string                 `json:"name"`
	RealName    string                 `json:"real_name"`
	CardNumber  string                 `json:"card_number"`
	Address     string                 `json:"address"`
	Status      int64                  `json:"status"`
	DateTime    *define.RangeTimeParam `json:"updated_at"`
	Pagination  *define.Pagination     `json:"pagination"`
}

type indexData struct {
	models.UserWalletAccountAttrs
	AdminName   string `json:"admin_name"`
	UserName    string `json:"username"`
	PaymentName string `json:"payment_name"`
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(indexParams)
	_ = body.ReadJSON(r, params)

	//  获取子级包括自己ID
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
	location, _ := time.LoadLocation(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "site_timezone"))

	model := models.NewUserWalletAccount(nil)
	model.AndWhere("admin_id in ("+strings.Join(adminIds, ",")+")").AndWhere("status<>?", models.WalletPaymentStatusDelete)

	define.NewFilterEmpty(model.Db).
		String("name like ?", "%"+params.Name+"%").
		String("real_name like ?", "%"+params.RealName+"%").
		String("card_number like ?", "%"+params.CardNumber+"%").
		String("address like ?", "%"+params.Address+"%").
		Int64("status=?", params.Status).
		RangeTime("updated_at between ? and ?", params.DateTime, location).
		Pagination(params.Pagination)

	// 管理员名称
	if params.AdminName != "" {
		model.Db.AndWhere("admin_id in (" + strings.Join(models.NewAdminUser(nil).FindAdminLikeNameIds(params.AdminName), ",") + ")")
	}

	// 用户名称
	if params.UserName != "" {
		model.Db.AndWhere("user_id in (" + strings.Join(models.NewUser(nil).FindUserLikeNameIds(params.UserName), ",") + ")")
	}

	// 支付名称
	if params.PaymentName != "" {
		paymentNameIds := models.NewWalletPayment(nil).Field("id").AndWhere("name like ?", "%"+params.PaymentName+"%").ColumnString()
		if len(paymentNameIds) == 0 {
			paymentNameIds = append(paymentNameIds, "-1")
		}
		model.Db.AndWhere("payment_id in (" + strings.Join(paymentNameIds, ",") + ")")
	}

	data := make([]*indexData, 0)
	model.Query(func(rows *sql.Rows) {
		tmp := new(indexData)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.UserId, &tmp.PaymentId, &tmp.Name, &tmp.RealName, &tmp.CardNumber, &tmp.Address, &tmp.Status, &tmp.Sort, &tmp.Data, &tmp.UpdatedAt, &tmp.CreatedAt)
		// 当前用户信息
		userModel := models.NewUser(nil)
		userModel.AndWhere("id=?", tmp.UserId)
		userInfo := userModel.FindOne()
		if userInfo != nil {
			tmp.UserName = userInfo.UserName
		}

		// 当前管理员信息
		adminModel := models.NewAdminUser(nil)
		adminModel.AndWhere("id=?", tmp.AdminId)
		adminInfo := adminModel.FindOne()
		if adminInfo != nil {
			tmp.AdminName = adminInfo.UserName
		}

		// 支付名称
		paymentModel := models.NewWalletPayment(nil)
		paymentModel.AndWhere("id=?", tmp.PaymentId)
		paymentInfo := paymentModel.FindOne()
		if paymentInfo != nil {
			tmp.PaymentName = paymentInfo.Name
		}

		data = append(data, tmp)
	})

	body.SuccessJSON(w, &body.IndexData{
		Items: data,
		Count: model.Count(),
	})
}
