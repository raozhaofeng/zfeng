package bill

import (
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/cache"
	"github.com/raozhaofeng/zfeng/utils/body"
	"net/http"
)

// TypesList 账单类型
func TypesList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make([]map[string]interface{}, 0)
	//adminId := router.TokenManager.GetContextClaims(r).AdminId

	rds := cache.RedisPool.Get()
	defer rds.Close()

	//for k, v := range models.UserBillTypeNameMap {
	// TODO...
	//data = append(data, map[string]interface{}{"label": beego.LocalesManager.GetAdminLocales(rds, adminId, "zh-CN", v), "value": k})
	//}

	body.SuccessJSON(w, data)
}
