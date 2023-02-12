package dictionary

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
	LangName   string                 `json:"lang_name"`
	Type       int64                  `json:"type"`
	Name       string                 `json:"name"`
	Field      string                 `json:"field"`
	Value      string                 `json:"value"`
	Data       string                 `json:"data"`
	DateTime   *define.RangeTimeParam `json:"created_at"`
	Pagination *define.Pagination     `json:"pagination"`
}

type indexData struct {
	models.LangDictionaryAttrs
	AdminName string `json:"admin_name"`
	LangName  string `json:"lang_name"`
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(indexParams)
	_ = body.ReadJSON(r, params)

	//  实例化模型
	model := models.NewLangDictionary(nil)
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	adminIds := models.NewAdminUser(nil).GetAdminChildrenParentIds(adminId)
	model.Db.AndWhere("admin_id in (" + strings.Join(adminIds, ",") + ")")
	location, _ := time.LoadLocation(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "site_timezone"))

	define.NewFilterEmpty(model.Db).
		Int64("type=?", params.Type).
		String("name like ?", "%"+params.Name+"%").
		String("field like ?", "%"+params.Field+"%").
		String("value like ?", "%"+params.Value+"%").
		String("data like ?", "%"+params.Data+"%").
		RangeTime("created_at between ? and ?", params.DateTime, location).
		Pagination(params.Pagination)

	if params.LangName != "" {
		langNameAlias := models.NewLang(nil).Field("alias").AndWhere("name like ?", "%"+params.LangName+"%").ColumnString()
		if len(langNameAlias) > 0 {
			var langNameAliasList []string
			for _, alias := range langNameAlias {
				langNameAliasList = append(langNameAliasList, "'"+alias+"'")
			}
			model.Db.AndWhere("alias in (" + strings.Join(langNameAliasList, ",") + ")")
		} else {
			model.Db.AndWhere("alias in ('none')")
		}
	}
	if params.AdminName != "" {
		model.Db.AndWhere("admin_id in (" + strings.Join(models.NewAdminUser(nil).FindAdminLikeNameIds(params.AdminName), ",") + ")")
	}

	data := make([]*indexData, 0)
	model.Query(func(rows *sql.Rows) {
		tmp := new(indexData)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.Type, &tmp.Alias, &tmp.Name, &tmp.Field, &tmp.Value, &tmp.Data, &tmp.CreatedAt)

		// 当前管理员信息
		adminModel := models.NewAdminUser(nil)
		adminModel.AndWhere("id=?", tmp.AdminId)
		adminInfo := adminModel.FindOne()
		if adminInfo != nil {
			tmp.AdminName = adminInfo.UserName
		}

		// 当前语言信息
		langModel := models.NewLang(nil)
		langModel.AndWhere("alias=?", tmp.Alias)
		langInfo := langModel.FindOne()
		if langInfo != nil {
			tmp.LangName = langInfo.Name
		}
		data = append(data, tmp)
	})

	body.SuccessJSON(w, &body.IndexData{
		Items: data,
		Count: model.Count(),
	})
}
