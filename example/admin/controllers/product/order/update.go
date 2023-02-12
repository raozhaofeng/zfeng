package order

import (
	"basic/models"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/database/define"
	"github.com/raozhaofeng/zfeng/router"
	"github.com/raozhaofeng/zfeng/utils/body"
	"github.com/raozhaofeng/zfeng/validator"
	"net/http"
	"strings"
	"time"
)

type updateParams struct {
	Id        int64                    `json:"id" validate:"required"`
	Money     float64                  `json:"money" validate:"omitempty,gte=0"`
	Nums      int64                    `json:"nums" validate:"omitempty,gte=0"`
	Type      int64                    `json:"type"`
	Data      *models.ProductDataAttrs `json:"data"`
	ExpiredAt string                   `json:"expired_at"`
	CreatedAt string                   `json:"created_at"`
}

func Update(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(updateParams)
	_ = body.ReadJSON(r, params)

	//  参数验证
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}
	//  实例化模型
	model := models.NewProductOrder(nil)
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	location, _ := time.LoadLocation(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "site_timezone"))

	productOrderData := ""
	if params.Data != nil {
		productOrderDataBytes, _ := json.Marshal(params.Data)
		productOrderData = string(productOrderDataBytes)
	}

	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		Float64("money=?", params.Money).
		Int64("nums=?", params.Nums).
		Int64("type=?", params.Type).
		String("data=?", productOrderData).
		DateTime("created_at=?", params.CreatedAt, location).
		DateTime("expired_at=?", params.ExpiredAt, location)

	if adminId != models.AdminUserSupermanId {
		adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
		model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	}
	_, err = model.AndWhere("id = ?", params.Id).Update()
	if err != nil {
		panic(err)
	}
	body.SuccessJSON(w, "ok")
}
