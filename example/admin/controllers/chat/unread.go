package chat

import (
	"basic/models"
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/router"
	"github.com/raozhaofeng/zfeng/utils/body"
	"github.com/raozhaofeng/zfeng/validator"
	"net/http"
)

type unreadParams struct {
	Id int64 `json:"id" validate:"required"`
}

// Unread 消除未读数量
func Unread(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(unreadParams)
	_ = body.ReadJSON(r, params)

	//  验证参数
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	adminId := router.TokenManager.GetContextClaims(r).AdminId
	_, _ = models.NewChatSession(nil).Value("main_user_unread=0").
		AndWhere("main_user=?", adminId).AndWhere("id=?", params.Id).
		Update()

	body.SuccessJSON(w, "ok")
}
