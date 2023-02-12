package product

import (
	"basic/models"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/router"
	"github.com/raozhaofeng/zfeng/utils/body"
	"github.com/raozhaofeng/zfeng/validator"
	"net/http"
	"time"
)

type createParams struct {
	AssetsId   int64                    `json:"assets_id"`
	CategoryId int64                    `json:"category_id" validate:"required"`
	Name       string                   `json:"name" validate:"required"`
	Images     []map[string]string      `json:"images" validate:"required"`
	Money      float64                  `json:"money" validate:"gt=0"`
	Sales      int64                    `json:"sales" validate:"omitempty,gte=0"`
	Nums       int64                    `json:"nums" validate:"required,gte=-1"`
	Mode       int64                    `json:"mode" validate:"required,oneof=1 2"`
	Used       int64                    `json:"used" validate:"omitempty,gt=0"`
	Total      int64                    `json:"total" validate:"omitempty,gt=0"`
	Data       *models.ProductDataAttrs `json:"data"`
	Describes  string                   `json:"describes"`
}

func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(createParams)
	_ = body.ReadJSON(r, params)

	//  验证参数
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	// 查询分类ID是否存在
	categoryModel := models.NewProductCategory(nil)
	categoryModel.AndWhere("id=?", params.CategoryId).AndWhere("status=?", models.ProductCategoryStatusActivate)
	categoryInfo := categoryModel.FindOne()
	if categoryInfo == nil {
		body.ErrorJSON(w, "分类ID不存在", -1)
		return
	}

	productDataBytes, _ := json.Marshal(params.Data)
	imagesByte, _ := json.Marshal(params.Images)
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	nowTime := time.Now()
	//  模型插入数据
	_, err = models.NewProduct(nil).
		Field("admin_id", "category_id", "assets_id", "name", "images", "money", "sales", "nums", "mode", "used", "total", "data", "describes", "updated_at", "created_at").
		Args(adminId, params.CategoryId, params.AssetsId, params.Name, string(imagesByte), params.Money, params.Sales, params.Nums, params.Mode, params.Used, params.Total, string(productDataBytes), params.Describes, nowTime.Unix(), nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}

	body.SuccessJSON(w, "ok")
}
