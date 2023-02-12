package product

import (
	"basic/models"
	"database/sql"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/database/define"
	"github.com/raozhaofeng/zfeng/router"
	"github.com/raozhaofeng/zfeng/utils/body"
	"net/http"
	"strings"
	"time"
)

type indexParams struct {
	CategoryName string                 `json:"category_name"`
	AdminName    string                 `json:"admin_name"`
	AssetsName   string                 `json:"assets_name"`
	Name         string                 `json:"name"`
	Status       int64                  `json:"status"`
	Recommend    int64                  `json:"recommend"`
	DateTime     *define.RangeTimeParam `json:"updated_at"`
	Pagination   *define.Pagination     `json:"pagination"`
}

type indexData struct {
	Data *models.ProductDataAttrs `json:"data"` //	数据
	models.ProductAttrs
	AdminName    string              `json:"admin_name"`    //管理员名称
	CategoryName string              `json:"category_name"` //分类ID
	ImagesList   []map[string]string `json:"images_list"`   // 产品图片
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(indexParams)
	_ = body.ReadJSON(r, params)

	//  获取子级包括自己ID
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
	location, _ := time.LoadLocation(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "site_timezone"))
	//  实例化模型
	model := models.NewProduct(nil)
	model.Db.AndWhere("admin_id in ("+strings.Join(adminIds, ",")+")").AndWhere("status>?", models.ProductStatusDelete)

	define.NewFilterEmpty(model.Db).
		String("name like ?", "%"+params.Name+"%").
		Int64("status=?", params.Status).
		Int64("recommend=?", params.Recommend).
		RangeTime("updated_at between ? and ?", params.DateTime, location).
		Pagination(params.Pagination)

	if params.AdminName != "" {
		model.Db.AndWhere("admin_id in (" + strings.Join(models.NewAdminUser(nil).FindAdminLikeNameIds(params.AdminName), ",") + ")")
	}

	if params.AssetsName != "" {
		assetsIds := models.NewAssets(nil).Field("id").AndWhere("name like ?", "%"+params.AssetsName+"%").ColumnString()
		if len(assetsIds) > 0 {
			model.Db.AndWhere("assets_id in (" + strings.Join(assetsIds, ",") + ")")
		} else {
			model.Db.AndWhere("assets_id in (-1)")
		}
	}

	if params.CategoryName != "" {
		categoryIds := models.NewProductCategory(nil).Field("id").AndWhere("name like ?", "%"+params.CategoryName+"%").ColumnString()
		if len(categoryIds) > 0 {
			model.Db.AndWhere("category_id in (" + strings.Join(categoryIds, ",") + ")")
		} else {
			model.Db.AndWhere("category_id in (-1)")
		}
	}

	data := make([]*indexData, 0)
	model.Query(func(rows *sql.Rows) {
		tmp := new(indexData)
		var productDataStr string
		imagesList := make([]map[string]string, 0)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.CategoryId, &tmp.AssetsId, &tmp.Name, &tmp.Images, &tmp.Money, &tmp.Sort, &tmp.Status, &tmp.Recommend, &tmp.Sales, &tmp.Nums, &tmp.Mode, &tmp.Used, &tmp.Total, &productDataStr, &tmp.Describes, &tmp.UpdatedAt, &tmp.CreatedAt)
		adminModel := models.NewAdminUser(nil)
		adminModel.AndWhere("id=?", tmp.AdminId)
		adminInfo := adminModel.FindOne()
		if adminInfo != nil {
			tmp.AdminName = adminInfo.UserName
		}
		if tmp.CategoryId > 0 {
			categoryModel := models.NewProductCategory(nil)
			categoryModel.AndWhere("id=?", tmp.CategoryId)
			categoryInfo := categoryModel.FindOne()
			if categoryInfo != nil {
				tmp.CategoryName = categoryInfo.Name
			}
		}
		if tmp.Images != "" {
			_ = json.Unmarshal([]byte(tmp.Images), &imagesList)
		}
		tmp.ImagesList = imagesList

		if tmp.AssetsId == 0 {
			tmp.AssetsId = -1
		}

		tmp.Data = new(models.ProductDataAttrs)
		_ = json.Unmarshal([]byte(productDataStr), &tmp.Data)

		data = append(data, tmp)
	})

	body.SuccessJSON(w, &body.IndexData{
		Items: data,
		Count: model.Count(),
	})
}
