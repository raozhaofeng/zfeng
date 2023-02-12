package lang

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
	Name  string `json:"name" validate:"required,max=50"`
	Alias string `json:"alias" validate:"required,max=50"`
	Icon  string `json:"icon" validate:"required,max=255"`
}

// Create 新增国家语言
func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(createParams)
	_ = body.ReadJSON(r, params)

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	nowTime := time.Now()
	adminId := router.TokenManager.GetContextClaims(r).AdminId

	//	判断语言别名是否存在
	oldLangModel := models.NewLang(nil)
	oldLangModel.AndWhere("admin_id=?", adminId).AndWhere("alias=?", params.Alias)
	oldLangInfo := oldLangModel.FindOne()
	if oldLangInfo != nil {
		body.ErrorJSON(w, "当前语言别名已存在", -1)
		return
	}

	//  模型插入数据
	_, err = models.NewLang(nil).
		Field("admin_id", "name", "alias", "icon", "created_at").
		Args(adminId, params.Name, params.Alias, params.Icon, nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}
	body.SuccessJSON(w, "ok")
}
