package chat

import (
	"basic/models"
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/router"
	"github.com/raozhaofeng/zfeng/utils/body"
	"github.com/raozhaofeng/zfeng/validator"
	"net/http"
)

type messageParams struct {
	Id int64 `json:"id" validate:"required"`
}

// Message 会话消息记录
func Message(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(messageParams)
	_ = body.ReadJSON(r, params)

	//  验证参数
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	adminId := router.TokenManager.GetContextClaims(r).AdminId
	sessionModel := models.NewChatSession(nil)
	sessionModel.AndWhere("id=?", params.Id).AndWhere("main_user=?", adminId).AndWhere("type < ?", models.ChatSessionTypeUserToUser)
	sessionInfo := sessionModel.FindOne()
	if sessionInfo == nil {
		body.ErrorJSON(w, "当前用户会话不存在", -1)
		return
	}

	messageModel := models.NewChatMessage(nil)
	messageModel.AndWhere("session_id=?", sessionInfo.Id).OffsetLimit(0, 100)

	body.SuccessJSON(w, messageModel.FindMany())
}
