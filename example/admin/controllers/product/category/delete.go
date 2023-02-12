package category

import (
	"basic/models"
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/router"
	"github.com/raozhaofeng/zfeng/utils/body"
	"github.com/raozhaofeng/zfeng/validator"
	"net/http"
	"strings"
	"time"
)

type deleteParams struct {
	Ids []int64 `json:"id" validate:"required"`
}

// Delete 删除分类
func Delete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(deleteParams)
	_ = body.ReadJSON(r, params)

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	nowTime := time.Now()
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	for _, v := range params.Ids {
		model := models.NewProductCategory(nil)
		if adminId != models.AdminUserSupermanId {
			adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
			model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
		}
		_, err = model.Value("status=?", "updated_at=?").Args(models.ProductCategoryStatusDelete, nowTime.Unix()).
			AndWhere("id=?", v).Update()
		if err != nil {
			panic(err)
		}
	}

	body.SuccessJSON(w, "ok")
}
