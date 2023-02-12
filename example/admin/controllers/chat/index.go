package chat

import (
	"basic/chat"
	"basic/models"
	"database/sql"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/router"
	"github.com/raozhaofeng/zfeng/utils/body"
	"net/http"
	"strconv"
)

type chatUserInfo struct {
	Id       int64  `json:"id"`       //	用户ID
	Avatar   string `json:"avatar"`   //	用户头像
	UserName string `json:"username"` //	用户昵称
	Online   bool   `json:"online"`   //	是否在线
	Ip4      string `json:"ip4"`      //	IP4
	Address  string `json:"address"`  //	地址
}

type sessionItem struct {
	Id         int64         `json:"id"`          //	会话ID
	SessionKey string        `json:"session_key"` //	会话Key
	Type       int64         `json:"type"`        //	会话类型
	UserInfo   *chatUserInfo `json:"user_info"`   //	用户信息
	Unread     int64         `json:"unread"`      //	未读消息
	Data       *chat.Message `json:"data"`        //	最后消息
	Ip4        string        `json:"ip4"`         //	IP地址
	UpdatedAt  int64         `json:"updated_at"`  //	更新时间
}

// Index 会话列表
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	adminId := router.TokenManager.GetContextClaims(r).AdminId

	data := make([]*sessionItem, 0)
	models.NewChatSession(nil).Field("id", "session_key", "type", "client_user", "data", "main_user_unread", "INET_NTOA(ip4)", "updated_at").
		AndWhere("main_user=?", adminId).AndWhere("type<?", models.ChatSessionTypeUserToUser).
		OrderBy("updated_at desc").OffsetLimit(0, 50).
		Query(func(rows *sql.Rows) {
			sessionTmp := new(sessionItem)
			var userId int64
			var messageStr string
			_ = rows.Scan(&sessionTmp.Id, &sessionTmp.SessionKey, &sessionTmp.Type, &userId, &messageStr, &sessionTmp.Unread, &sessionTmp.Ip4, &sessionTmp.UpdatedAt)

			//	最后读消息内容
			_ = json.Unmarshal([]byte(messageStr), &sessionTmp.Data)

			//	用户信息
			sessionTmp.UserInfo = new(chatUserInfo)
			sessionTmp.UserInfo.Ip4 = sessionTmp.Ip4
			if sessionTmp.Type == models.ChatSessionTypeAdminToTourist {
				//	临时用户
				sessionTmp.UserInfo.Online = chat.Manager.TouristOnlineStatus(userId)
				sessionTmp.UserInfo.Id = userId
				sessionTmp.UserInfo.UserName = "临时用户" + strconv.FormatInt(userId, 10)
			} else {
				//	登陆用户信息
				sessionTmp.UserInfo.Online = chat.Manager.UserOnlineStatus(userId)
				models.NewUser(nil).Field("id", "avatar", "username").
					AndWhere("id=?", userId).QueryRow(func(row *sql.Row) {
					_ = row.Scan(&sessionTmp.UserInfo.Id, &sessionTmp.UserInfo.Avatar, &sessionTmp.UserInfo.UserName)
				})
			}

			//	获取IP4地址
			ip2location, _ := models.GetIp2Location(sessionTmp.UserInfo.Ip4)
			if ip2location != nil {
				sessionTmp.UserInfo.Address = ip2location.Country_long + "." + ip2location.Region + "." + ip2location.City
			}

			data = append(data, sessionTmp)
		})

	body.SuccessJSON(w, data)
}
