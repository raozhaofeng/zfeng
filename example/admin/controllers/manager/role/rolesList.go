package role

import (
	"basic/models"
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/router"
	"github.com/raozhaofeng/zfeng/utils/body"
	"net/http"
)

// RolesList 角色列表
func RolesList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	// 获取管理员的角色
	adminRoles := models.NewAdminAuthItem(nil).GetAdminRoleCheckedList(adminId, []string{})
	data := make([]map[string]string, 0)
	for role, _ := range adminRoles {
		data = append(data, map[string]string{"label": role, "value": role})
	}

	body.SuccessJSON(w, data)
}
