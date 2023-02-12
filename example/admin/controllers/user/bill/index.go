package bill

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
	UserName   string                 `json:"username"`
	Name       string                 `json:"name"`
	Type       int64                  `json:"type"`
	DateTime   *define.RangeTimeParam `json:"created_at"`
	Pagination *define.Pagination     `json:"pagination"`
}

type indexData struct {
	models.UserBillAttrs
	AdminName string `json:"admin_name"` //管理员名称
	UserName  string `json:"username"`   //用户名
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(indexParams)
	_ = body.ReadJSON(r, params)

	//  获取子级包括自己ID
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
	location, _ := time.LoadLocation(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "site_timezone"))

	model := models.NewUserBill(nil)
	model.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")

	define.NewFilterEmpty(model.Db).
		String("name like ?", "%"+params.Name+"%").
		Int64("type=?", params.Type).
		RangeTime("created_at between ? and ?", params.DateTime, location).
		Pagination(params.Pagination)
	// 管理员名称
	if params.AdminName != "" {
		model.Db.AndWhere("admin_id in (" + strings.Join(models.NewAdminUser(nil).FindAdminLikeNameIds(params.AdminName), ",") + ")")
	}

	// 用户名称
	if params.UserName != "" {
		model.Db.AndWhere("user_id in (" + strings.Join(models.NewUser(nil).FindUserLikeNameIds(params.UserName), ",") + ")")
	}

	data := make([]*indexData, 0)
	model.Query(func(rows *sql.Rows) {
		tmp := new(indexData)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.UserId, &tmp.SourceId, &tmp.Name, &tmp.Type, &tmp.Balance, &tmp.Money, &tmp.Data, &tmp.CreatedAt)

		// 当前用户信息
		userModel := models.NewUser(nil)
		userModel.AndWhere("id=?", tmp.UserId)
		userInfo := userModel.FindOne()
		if userInfo != nil {
			tmp.UserName = userInfo.UserName
		}

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
