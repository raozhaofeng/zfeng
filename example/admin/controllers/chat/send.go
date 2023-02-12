package chat

import (
	"basic/chat"
	"basic/models"
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/database"
	"github.com/raozhaofeng/zfeng/router"
	"github.com/raozhaofeng/zfeng/utils/body"
	"github.com/raozhaofeng/zfeng/validator"
	"net/http"
	"time"
)

type sendParams struct {
	Id      int64  `json:"id" validate:"required"`
	Type    int64  `json:"type" validate:"required,oneof=1 2"`
	Message string `json:"message" validate:"required"`
	Extra   string `json:"extra"`
}

// Send 发送聊天消息
func Send(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(sendParams)
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

	//	创建消息, 更新会话未读数 并且 发送聊天消息
	tx := database.DbPool.GetTx()
	defer tx.Rollback()

	//	发送消息
	roleId := chat.MessageRoleAdminToUser
	if sessionInfo.Type == models.ChatSessionTypeAdminToTourist {
		roleId = chat.MessageRoleAdminToTourist
	}
	msg := &chat.Message{
		SessionId:  sessionInfo.Id,
		RoleId:     roleId,
		SenderId:   sessionInfo.MainUser,
		ReceiverId: sessionInfo.ClientUser,
		Type:       params.Type,
		Data:       params.Message,
		Extra:      params.Extra,
		Time:       time.Now().Unix(),
	}
	err = models.NewChatMessage(tx).NewSendMessage(msg, r)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	_ = tx.Commit()
	body.SuccessJSON(w, "ok")
}
