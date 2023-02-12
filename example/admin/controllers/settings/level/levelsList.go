package level

import (
	"basic/models"
	"database/sql"
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/router"
	"github.com/raozhaofeng/zfeng/utils/body"
	"net/http"
)

func LevelsList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	adminId := router.TokenManager.GetContextClaims(r).AdminId

	data := make([]map[string]any, 0)
	models.NewUserLevel(nil).Field("id", "name").
		AndWhere("admin_id=?", adminId).
		Query(func(rows *sql.Rows) {
			var id int64
			var name string
			_ = rows.Scan(&id, &name)
			data = append(data, map[string]any{
				"label": name, "value": id,
			})
		})
	body.SuccessJSON(w, data)
}
