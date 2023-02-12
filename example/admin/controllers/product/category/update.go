package category

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
)

type updateParams struct {
	Id        int64                       `json:"id" validate:"required"`
	ParentId  int64                       `json:"parent_id" validate:"omitempty,gt=0"`
	Name      string                      `json:"name"`
	Image     string                      `json:"image"`
	Sort      int64                       `json:"sort" validate:"omitempty,gt=0"`
	Status    int64                       `json:"status" validate:"omitempty,oneof=-1 10"`
	Type      int64                       `json:"type" validate:"omitempty,oneof=1 10"`
	Recommend int64                       `json:"recommend" validate:"omitempty,oneof=-1 10"`
	Data      *models.ProductCategoryData `json:"data"`
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
	model := models.NewProductCategory(nil)
	adminId := router.TokenManager.GetContextClaims(r).AdminId

	// 判断父级ID是否存在
	if params.ParentId > 0 {
		parentModel := models.NewProductCategory(nil)
		parentModel.AndWhere("id=?", params.ParentId)
		parentInfo := parentModel.FindOne()
		if parentInfo == nil {
			body.ErrorJSON(w, "父级ID不存在", -1)
			return
		}
	}

	productCategoryData := ""
	if params.Data != nil {
		productCategoryDataBytes, _ := json.Marshal(params.Data)
		productCategoryData = string(productCategoryDataBytes)
	}

	define.NewFilterEmpty(model.Db).SetUpdateOpt().
		Int64("parent_id=?", params.ParentId).
		String("name=?", params.Name).
		String("image=?", params.Image).
		Int64("sort=?", params.Sort).
		Int64("status=?", params.Status).
		Int64("type=?", params.Type).
		Int64("recommend=?", params.Recommend).
		String("data=?", productCategoryData)

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
