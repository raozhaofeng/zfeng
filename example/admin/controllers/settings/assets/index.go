package assets

import (
	"basic/models"
	"database/sql"
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/database/define"
	"github.com/raozhaofeng/zfeng/router"
	"github.com/raozhaofeng/zfeng/utils/body"
	"net/http"
	"strings"
	"time"
)

type indexParams struct {
	AdminName  string                 `json:"admin_name"`
	Name       string                 `json:"name"`
	Status     int64                  `json:"status"`
	Type       int64                  `json:"type"`
	DateTime   *define.RangeTimeParam `json:"updated_at"`
	Pagination *define.Pagination     `json:"pagination"`
}

type indexData struct {
	models.AssetsAttrs
	AdminName string `json:"admin_name"` //管理名称
}

// Index 资产列表
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(indexParams)
	_ = body.ReadJSON(r, params)

	//  实例化模型
	model := models.NewAssets(nil)
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
	model.Db.AndWhere("admin_id in ("+strings.Join(adminIds, ",")+")").AndWhere("status>?", models.AssetsStatusDelete)
	location, _ := time.LoadLocation(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "site_timezone"))

	define.NewFilterEmpty(model.Db).
		String("name like ?", "%"+params.Name+"%").
		Int64("status=?", params.Status).
		Int64("type=?", params.Type).
		RangeTime("created_at between ? and ?", params.DateTime, location).
		Pagination(params.Pagination)

	if params.AdminName != "" {
		model.Db.AndWhere("admin_id in (" + strings.Join(models.NewAdminUser(nil).FindAdminLikeNameIds(params.AdminName), ",") + ")")
	}

	data := make([]*indexData, 0)
	model.Query(func(rows *sql.Rows) {
		tmp := new(indexData)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.Name, &tmp.Icon, &tmp.Type, &tmp.Data, &tmp.Status, &tmp.CreatedAt, &tmp.UpdatedAt)
		// 当前管理员信息
		adminModel := models.NewAdminUser(nil)
		adminModel.AndWhere("id=?", tmp.AdminId)
		adminInfo := adminModel.FindOne()
		if adminInfo != nil {
			tmp.AdminName = adminInfo.UserName
		}
		data = append(data, tmp)
	})

	body.SuccessJSON(w, &body.IndexData{
		Items: data,
		Count: model.Count(),
	})
}
