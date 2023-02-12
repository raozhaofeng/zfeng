package role

import (
	"basic/models"
	"database/sql"
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/router"
	"github.com/raozhaofeng/zfeng/utils/body"
	"net/http"
)

// PermissionsList 权限列表
func PermissionsList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	adminId := router.TokenManager.GetContextClaims(r).AdminId

	data := make([]map[string]string, 0)
	models.NewAdminAuthItem(nil).Field("name").
		AndWhere("type=?", models.AdminAuthItemTypeRouteName).
		Query(func(rows *sql.Rows) {
			var name string
			_ = rows.Scan(&name)
			if adminId == models.AdminUserSupermanId || name != "所有权限" {
				data = append(data, map[string]string{"label": name, "value": name})
			}
		})

	body.SuccessJSON(w, data)
}
