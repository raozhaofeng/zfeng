package category

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
	ParentName string                 `json:"parent_name"`
	AdminName  string                 `json:"admin_name"`
	Name       string                 `json:"name"`
	Status     int64                  `json:"status"`
	Type       int64                  `json:"type"`
	Recommend  int64                  `json:"recommend"`
	DateTime   *define.RangeTimeParam `json:"updated_at"`
	Pagination *define.Pagination     `json:"pagination"`
}

type indexData struct {
	models.ProductCategoryAttrs
	Data       *models.ProductCategoryData `json:"data"`        //数据
	AdminName  string                      `json:"admin_name"`  //管理员名称
	ParentName string                      `json:"parent_name"` //父级名称
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
	model := models.NewProductCategory(nil)
	model.Db.AndWhere("admin_id in ("+strings.Join(adminIds, ",")+")").AndWhere("status>?", models.ProductCategoryStatusDelete)

	define.NewFilterEmpty(model.Db).
		String("name like ?", "%"+params.Name+"%").
		Int64("status=?", params.Status).
		Int64("type=?", params.Type).
		Int64("recommend=?", params.Recommend).
		RangeTime("updated_at between ? and ?", params.DateTime, location).
		Pagination(params.Pagination)

	if params.AdminName != "" {
		model.Db.AndWhere("admin_id in (" + strings.Join(models.NewAdminUser(nil).FindAdminLikeNameIds(params.AdminName), ",") + ")")
	}
	if params.ParentName != "" {
		parentIds := models.NewProductCategory(nil).Field("id").AndWhere("name like ?", "%"+params.ParentName+"%").ColumnString()
		if len(parentIds) > 0 {
			model.Db.AndWhere("parent_id in (" + strings.Join(parentIds, ",") + ")")
		}
	}

	data := make([]*indexData, 0)
	model.Query(func(rows *sql.Rows) {
		tmp := new(indexData)
		var categoryData string
		_ = rows.Scan(&tmp.Id, &tmp.ParentId, &tmp.AdminId, &tmp.Type, &tmp.Name, &tmp.Image, &tmp.Sort, &tmp.Status, &tmp.Recommend, &categoryData, &tmp.UpdatedAt, &tmp.CreatedAt)
		adminModel := models.NewAdminUser(nil)
		adminModel.AndWhere("id=?", tmp.AdminId)
		adminInfo := adminModel.FindOne()
		if adminInfo != nil {
			tmp.AdminName = adminInfo.UserName
		}
		if tmp.ParentId > 0 {
			parentModel := models.NewProductCategory(nil)
			parentModel.AndWhere("id=?", tmp.ParentId)
			parentInfo := parentModel.FindOne()
			if parentInfo != nil {
				tmp.ParentName = parentInfo.Name
			}
		}

		tmp.Data = new(models.ProductCategoryData)
		_ = json.Unmarshal([]byte(categoryData), &tmp.Data)
		data = append(data, tmp)
	})

	body.SuccessJSON(w, &body.IndexData{
		Items: data,
		Count: model.Count(),
	})
}
