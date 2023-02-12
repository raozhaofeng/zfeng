package dictionary

import (
	"basic/models"
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/cache"
	"github.com/raozhaofeng/zfeng/router"
	"github.com/raozhaofeng/zfeng/utils/body"
	"github.com/raozhaofeng/zfeng/validator"
	"net/http"
	"time"
)

type createParams struct {
	Alias string `json:"alias" validate:"required"`
	Type  int64  `json:"type" validate:"required,oneof=1 2 10"`
	Name  string `json:"name" validate:"required,max=50"`
	Field string `json:"field" validate:"required,max=50"`
	Value string `json:"value" validate:"required"`
	Data  string `json:"data" validate:"max=255"`
}

func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(createParams)
	_ = body.ReadJSON(r, params)

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	//  模型插入数据
	nowTime := time.Now()
	adminId := router.TokenManager.GetContextClaims(r).AdminId

	oldLangModel := models.NewLang(nil)
	oldLangModel.AndWhere("admin_id=?", adminId).AndWhere("alias=?", params.Alias).AndWhere("status>?", models.LangStatusDisabled)
	oldLangInfo := oldLangModel.FindOne()
	if oldLangInfo == nil {
		body.ErrorJSON(w, "语言别名不存在", -1)
		return
	}

	_, err = models.NewLangDictionary(nil).
		Field("admin_id", "type", "alias", "name", "field", "value", "data", "created_at").
		Args(adminId, params.Type, params.Alias, params.Name, params.Field, params.Value, params.Data, nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}

	// 如果更新对类型是接口，或者数据翻译， 那么重载语言配置
	if params.Value != "" && params.Field != "" {
		rds := cache.RedisPool.Get()
		defer rds.Close()
		// TODO...
	}
	body.SuccessJSON(w, "ok")
}
