package assets

import (
	"basic/models"
	"database/sql"
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/router"
	"github.com/raozhaofeng/zfeng/utils/body"
	"net/http"
)

// List 资产列表
func List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)

	data := make([]map[string]any, 0)
	models.NewAssets(nil).Field("id", "name", "type").
		AndWhere("admin_id=?", settingAdminId).AndWhere("status>?", models.AssetsStatusDelete).
		Query(func(rows *sql.Rows) {
			var id, assetsType int64
			var name string
			_ = rows.Scan(&id, &name, &assetsType)
			data = append(data, map[string]any{
				"label": models.AssetsTypeList[assetsType] + "(" + name + ")", "value": id,
			})
		})
	body.SuccessJSON(w, data)
}
