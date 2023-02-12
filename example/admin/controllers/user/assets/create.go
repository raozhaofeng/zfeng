package assets

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
	UserName string `json:"username" validate:"required"`
	Type     int64  `json:"type" validate:"required,oneof=1 2"`
	Name     string `json:"name" validate:"required"`
}

func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(createParams)
	_ = body.ReadJSON(r, params)

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	adminId := router.TokenManager.GetContextClaims(r).AdminId
	userModel := models.NewUser(nil)
	userModel.AndWhere("username=?", params.UserName).AndWhere("admin_id=?", adminId)
	userInfo := userModel.FindOne()
	if userInfo == nil {
		body.ErrorJSON(w, "用户名不存在", -1)
		return
	}

	// 如果资产已存在, 那么不能添加
	userAssetsModel := models.NewUserAssets(nil)
	userAssetsModel.AndWhere("user_id=?", userInfo.Id).AndWhere("name=?", params.Name).AndWhere("status>?", models.ProductStatusDelete)
	userAssetsInfo := userAssetsModel.FindOne()
	if userAssetsInfo != nil {
		body.ErrorJSON(w, "当前用户资产已存在", -1)
		return
	}

	nowTime := time.Now()
	_, err = models.NewUserAssets(nil).Field("admin_id", "user_id", "name", "type", "created_at", "updated_at").
		Args(userInfo.AdminId, userInfo.Id, params.Name, params.Type, nowTime.Unix(), nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}

	body.SuccessJSON(w, "ok")
}
