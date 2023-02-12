package models

import (
	"basic/chat"
	"database/sql"
	"github.com/raozhaofeng/zfeng/database"
	"github.com/raozhaofeng/zfeng/database/define"
	"github.com/raozhaofeng/zfeng/utils"
	"net/http"
	"time"
)

const (
	ChatSessionTypeAdminToUser    int64 = 10 //	管理对用户
	ChatSessionTypeAdminToTourist int64 = 11 //	管理对临时用户
	ChatSessionTypeUserToUser     int64 = 20 //	用户对用户
)

// ChatSessionAttrs 数据库模型属性
type ChatSessionAttrs struct {
	Id               int64  `json:"id"`                 //主键
	SessionKey       string `json:"session_key"`        //会话Key
	AdminId          int64  `json:"admin_id"`           //管理ID
	MainUser         int64  `json:"main_user"`          //主用户
	MainUserUnread   int64  `json:"main_user_unread"`   //主用户未读
	ClientUser       int64  `json:"client_user"`        //客户用户
	ClientUserUnread int64  `json:"client_user_unread"` //客户用户未读
	Type             int64  `json:"type"`               //10管理会话[管理对用户] 11临时会话[管理对临时] 20用户会话[用户对用户]
	Data             string `json:"data"`               //最后消息内容
	Ip4              string `json:"ip4"`                //IP4地址
	UpdatedAt        int64  `json:"updated_at"`         //更新时间
}

// ChatSession 数据库模型
type ChatSession struct {
	define.Db
}

// NewChatSession 创建数据库模型
func NewChatSession(tx *sql.Tx) *ChatSession {
	return &ChatSession{
		database.DbPool.NewDb(tx).Table("chat_session"),
	}
}

// FindOne 查询单挑
func (c *ChatSession) FindOne() *ChatSessionAttrs {
	attrs := new(ChatSessionAttrs)
	c.Field("id", "session_key", "admin_id", "main_user", "main_user_unread", "client_user", "client_user_unread", "type", "data", "INET_NTOA(ip4)", "updated_at")
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.SessionKey, &attrs.AdminId, &attrs.MainUser, &attrs.MainUserUnread, &attrs.ClientUser, &attrs.ClientUserUnread, &attrs.Type, &attrs.Data, &attrs.Ip4, &attrs.UpdatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *ChatSession) FindMany() []*ChatSessionAttrs {
	data := make([]*ChatSessionAttrs, 0)
	c.Field("id", "session_key", "admin_id", "main_user", "main_user_unread", "client_user", "client_user_unread", "type", "data", "INET_NTOA(ip4)", "updated_at")
	c.Query(func(rows *sql.Rows) {
		tmp := new(ChatSessionAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.SessionKey, &tmp.AdminId, &tmp.MainUser, &tmp.MainUserUnread, &tmp.ClientUser, &tmp.ClientUserUnread, &tmp.Type, &tmp.Data, &tmp.Ip4, &tmp.UpdatedAt)
		data = append(data, tmp)
	})
	return data
}

// InitSession 初始化会话
func (c *ChatSession) InitSession(adminId, mainUser, clientUser int64, sessionType int64, r *http.Request) *ChatSessionAttrs {
	c.AndWhere("admin_id=?", adminId).AndWhere("main_user=?", mainUser).AndWhere("client_user=?", clientUser)
	sessionInfo := c.FindOne()

	if sessionInfo == nil {
		nowTime := time.Now()
		sessionKey := utils.NewRandom().String(32)
		ip4 := utils.GetUserRealIP(r)
		sessionId, err := NewChatSession(nil).Field("session_key", "admin_id", "main_user", "client_user", "type", "data", "ip4", "updated_at").
			Value("?", "?", "?", "?", "?", "?", "INET_ATON(?)", "?").
			Args(sessionKey, adminId, mainUser, clientUser, sessionType, "", utils.GetUserRealIP(r), nowTime.Unix()).Insert()
		if err != nil {
			panic(err)
		}
		sessionInfo = &ChatSessionAttrs{
			Id: sessionId, SessionKey: sessionKey, AdminId: adminId, MainUser: mainUser, MainUserUnread: 0,
			ClientUser: clientUser, ClientUserUnread: 0, Type: sessionType, Data: "", Ip4: ip4, UpdatedAt: nowTime.Unix(),
		}
	}
	return sessionInfo
}

// OnlineAndOfflineNotify 上线｜下线	通知
func (c *ChatSession) OnlineAndOfflineNotify(isOnline bool, clientUser int64) {
	c.AndWhere("client_user=?", clientUser).AndWhere("type<>?", ChatSessionTypeUserToUser)
	sessionInfo := c.FindOne()

	if sessionInfo != nil {
		msgType := chat.MessageTypeOnline
		if !isOnline {
			msgType = chat.MessageTypeOffline
		}

		chat.Manager.SendAdminNotify(sessionInfo.Id, sessionInfo.MainUser, msgType)
	}
}

// GetChatUserMessageRole 获取用户消息角色
func (c *ChatSession) GetChatUserMessageRole(sessionInfo *ChatSessionAttrs) int64 {
	switch sessionInfo.Type {
	case ChatSessionTypeAdminToTourist:
		return chat.MessageRoleTouristToAdmin
	case ChatSessionTypeUserToUser:

	case ChatSessionTypeAdminToUser:
		return chat.MessageRoleUserToAdmin
	}
	panic("不存在的会话类型")
}
