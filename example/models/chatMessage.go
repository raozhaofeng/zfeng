package models

import (
	"basic/chat"
	"database/sql"
	"encoding/json"
	"github.com/raozhaofeng/zfeng/database"
	"github.com/raozhaofeng/zfeng/database/define"
	"github.com/raozhaofeng/zfeng/utils"
	"net/http"
	"time"
)

// ChatMessageAttrs 数据库模型属性
type ChatMessageAttrs struct {
	Id         int64  `json:"id"`          //主键
	SessionId  int64  `json:"session_id"`  //会话ID
	RoleId     int64  `json:"role_id"`     //10用户发送管理 | 11临时用户发送管理 | 12管理发送用户 | 13管理发送临时用户 | 14主用户发送用户 | 15用户发送主用户
	SenderId   int64  `json:"sender_id"`   //发送者
	ReceiverId int64  `json:"receiver_id"` //接收者
	Type       int64  `json:"type"`        //消息类型
	Data       string `json:"data"`        //数据
	Extra      string `json:"extra"`       //额外数据
	CreatedAt  int64  `json:"created_at"`  //创建时间
}

// ChatMessage 数据库模型
type ChatMessage struct {
	define.Db
}

// NewChatMessage 创建数据库模型
func NewChatMessage(tx *sql.Tx) *ChatMessage {
	return &ChatMessage{
		database.DbPool.NewDb(tx).Table("chat_message"),
	}
}

// FindOne 查询单挑
func (c *ChatMessage) FindOne() *ChatMessageAttrs {
	attrs := new(ChatMessageAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.SessionId, &attrs.RoleId, &attrs.SenderId, &attrs.ReceiverId, &attrs.Type, &attrs.Data, &attrs.Extra, &attrs.CreatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *ChatMessage) FindMany() []*ChatMessageAttrs {
	data := make([]*ChatMessageAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(ChatMessageAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.SessionId, &tmp.RoleId, &tmp.SenderId, &tmp.ReceiverId, &tmp.Type, &tmp.Data, &tmp.Extra, &tmp.CreatedAt)
		data = append(data, tmp)
	})
	return data
}

// NewSendMessage 创建新的发送消息
func (c *ChatMessage) NewSendMessage(msg *chat.Message, r *http.Request) error {
	if msg.Time == 0 {
		msg.Time = time.Now().Unix()
	}
	messageId, err := NewChatMessage(c.GetTx()).
		Field("session_id", "role_id", "sender_id", "receiver_id", "type", "data", "extra", "created_at").
		Args(msg.SessionId, msg.RoleId, msg.SenderId, msg.ReceiverId, msg.Type, msg.Data, msg.Extra, msg.Time).
		Insert()
	if err != nil {
		return err
	}

	//	发送消息
	msg.Id = messageId
	chat.Manager.Send(msg)

	// 更新会话
	msgBytes, _ := json.Marshal(msg)
	model := NewChatSession(c.GetTx()).Value("data=?", "updated_at=?").
		Args(string(msgBytes), time.Now().Unix())

	//	替换更新内容
	switch msg.RoleId {
	case chat.MessageRoleAdminToUser, chat.MessageRoleAdminToTourist, chat.MessageRoleMainUserToUser:
		model.Value("client_user_unread=(client_user_unread+1)")
	default:
		model.Value("main_user_unread=(main_user_unread+1)", "ip4=INET_ATON(?)").Args(utils.GetUserRealIP(r))
	}

	//	更新会话
	_, err = model.AndWhere("id=?", msg.SessionId).Update()
	if err != nil {
		return err
	}
	return nil
}
