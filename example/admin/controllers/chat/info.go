package chat

import (
	"basic/chat"
	"basic/models"
	"database/sql"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/router"
	"github.com/raozhaofeng/zfeng/utils/body"
	"github.com/raozhaofeng/zfeng/validator"
	"net/http"
	"strconv"
)

type infoParams struct {
	Id int64 `json:"id" validate:"required"`
}

func Info(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(infoParams)
	_ = body.ReadJSON(r, params)

	// 验证参数
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

	//	返回数据
	data := &sessionItem{
		Id:         sessionInfo.Id,
		SessionKey: sessionInfo.SessionKey,
		Type:       sessionInfo.Type,
		Unread:     sessionInfo.MainUserUnread,
		UpdatedAt:  sessionInfo.UpdatedAt,
	}

	//	最后读消息内容
	_ = json.Unmarshal([]byte(sessionInfo.Data), &data.Data)

	// 用户信息
	data.UserInfo = new(chatUserInfo)
	data.UserInfo.Ip4 = sessionInfo.Ip4
	if sessionInfo.Type == models.ChatSessionTypeAdminToTourist {
		//	临时用户
		data.UserInfo.Online = chat.Manager.TouristOnlineStatus(sessionInfo.ClientUser)
		data.UserInfo.Id = sessionInfo.ClientUser
		data.UserInfo.UserName = "临时用户" + strconv.FormatInt(sessionInfo.ClientUser, 10)
	} else {
		//	登陆用户信息
		data.UserInfo.Online = chat.Manager.UserOnlineStatus(sessionInfo.ClientUser)
		models.NewUser(nil).Field("id", "avatar", "username").
			AndWhere("id=?", sessionInfo.ClientUser).QueryRow(func(row *sql.Row) {
			_ = row.Scan(&data.UserInfo.Id, &data.UserInfo.Avatar, &data.UserInfo.UserName)
		})
	}

	//	获取IP4地址
	ip2location, _ := models.GetIp2Location(data.UserInfo.Ip4)
	if ip2location != nil {
		data.UserInfo.Address = ip2location.Country_long + "." + ip2location.Region + "." + ip2location.City
	}

	body.SuccessJSON(w, data)
}
