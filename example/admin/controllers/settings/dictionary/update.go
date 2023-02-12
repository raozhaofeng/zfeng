package dictionary

import (
	"basic/models"
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/cache"
	"github.com/raozhaofeng/zfeng/database/define"
	"github.com/raozhaofeng/zfeng/router"
	"github.com/raozhaofeng/zfeng/utils/body"
	"github.com/raozhaofeng/zfeng/validator"
	"net/http"
	"strings"
)

type updateParams struct {
	Id    int64  `json:"id" validate:"required"`
	Alias string `json:"alias"`
	Type  int64  `json:"type" validate:"omitempty,oneof=1 2 10"`
	Name  string `json:"name" validate:"max=50"`
	Field string `json:"field" validate:"max=50"`
	Value string `json:"value"`
	Data  string `json:"data" validate:"max=255"`
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
	model := models.NewLangDictionary(nil)
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		String("alias=?", params.Alias).
		Int64("type=?", params.Type).
		String("name=?", params.Name).
		String("field=?", params.Field).
		String("value=?", params.Value).
		String("data=?", params.Data)

	//  模型增加where条件并更新
	if adminId != models.AdminUserSupermanId {
		adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
		model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	}
	_, err = model.AndWhere("id=?", params.Id).Update()
	if err != nil {
		panic(err)
	}

	// 如果更新对类型是接口，或者数据翻译， 那么重载语言配置
	if params.Value != "" && params.Field != "" {
		rds := cache.RedisPool.Get()
		defer rds.Close()
		// TODO..

	}

	body.SuccessJSON(w, "ok")
}
