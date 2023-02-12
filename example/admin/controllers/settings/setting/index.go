package setting

import (
	"basic/models"
	"database/sql"
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/router"
	"github.com/raozhaofeng/zfeng/utils/body"
	"net/http"
)

type settingAttrs struct {
	models.AdminSettingAttrs
	AdminName string `json:"admin_name"`
}

type groupSetting struct {
	Id       int64           `json:"id"`
	Name     string          `json:"name"`
	Children []*settingAttrs `json:"children"`
}

type indexData struct {
	Groups []*groupSetting `json:"groups"`
}

var settingGroupList = []*groupSetting{
	{Id: models.SettingGroupBasic, Name: "基本设置"},
	{Id: models.SettingGroupHome, Name: "首页设置"},
	{Id: models.SettingGroupFinance, Name: "财务设置"},
	{Id: models.SettingGroupTemplate, Name: "模版设置"},
	{Id: models.SettingGroupHelpers, Name: "帮助中心"},
}

// Index 设置
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	adminId := router.TokenManager.GetContextClaims(r).AdminId

	data := make([]*groupSetting, 0)
	for _, group := range settingGroupList {
		model := models.NewAdminSetting(nil)
		model.AndWhere("admin_id=?", adminId)

		items := &groupSetting{
			Id: group.Id, Name: group.Name, Children: make([]*settingAttrs, 0),
		}

		model.AndWhere("group_id=?", group.Id)
		model.Query(func(rows *sql.Rows) {
			tmp := new(settingAttrs)
			_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.GroupId, &tmp.Name, &tmp.Type, &tmp.Field, &tmp.Value, &tmp.Data)
			adminModel := models.NewAdminUser(nil)
			adminModel.AndWhere("id=?", tmp.AdminId)
			adminInfo := adminModel.FindOne()
			if adminInfo != nil {
				tmp.AdminName = adminInfo.UserName
			}
			items.Children = append(items.Children, tmp)
		})
		data = append(data, items)
	}

	body.SuccessJSON(w, &indexData{
		Groups: data,
	})
}
