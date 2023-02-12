package chat

import "time"

const (
	MessageTypeAudio   int64 = 1000 //	音频消息
	MessageTypeOnline  int64 = 1010 //	上线类型
	MessageTypeOffline int64 = 1011 //	下线类型
	MessageTypeText    int64 = 1    //	文本消息
	MessageTypeImage   int64 = 2    //	图片消息
)

type Message struct {
	Id         int64 `json:"id"`          //	消息ID
	SessionId  int64 `json:"session_id"`  //	会话ID
	RoleId     int64 `json:"role_id"`     //	角色ID
	SenderId   int64 `json:"sender_id"`   //	发送者
	ReceiverId int64 `json:"receiver_id"` //	接收者
	Type       int64 `json:"type"`        //	消息类型
	Data       any   `json:"data"`        //	消息内容
	Extra      any   `json:"extra"`       //	额外数据
	Time       int64 `json:"time"`        //	发送时间
}

// SendAdminAudio 发送音频消息
func (c *Chat) SendAdminAudio(adminId int64, name, audioPath string) {
	c.Send(&Message{
		RoleId:     MessageRoleSystemToAdmin,
		ReceiverId: adminId,
		Type:       MessageTypeAudio,
		Data:       name,
		Extra:      audioPath,
		Time:       time.Now().Unix(),
	})
}

// SendAdminNotify 发送消息通知
func (c *Chat) SendAdminNotify(sessionId int64, adminId int64, msgType int64) {
	c.Send(&Message{
		SessionId:  sessionId,
		RoleId:     MessageRoleSystemToAdmin,
		ReceiverId: adminId,
		Type:       msgType,
	})
}

// Send 发送消息
func (c *Chat) Send(msg *Message) {
	switch msg.RoleId {
	// 系统｜用户｜游客  => 管理
	case MessageRoleSystemToAdmin, MessageRoleUserToAdmin, MessageRoleTouristToAdmin:
		if _, ok := c.AdminConn[msg.ReceiverId]; ok {
			_ = c.AdminConn[msg.ReceiverId].WriteJSON(msg)
		}
	// 管理|用户 => 用户
	case MessageRoleAdminToUser, MessageRoleMainUserToUser, MessageRoleUserToMainUser:
		if _, ok := c.UserConn[msg.ReceiverId]; ok {
			_ = c.UserConn[msg.ReceiverId].WriteJSON(msg)
		}
	// 管理 => 临时用户
	case MessageRoleAdminToTourist:
		if _, ok := c.TouristConn[msg.ReceiverId]; ok {
			_ = c.TouristConn[msg.ReceiverId].WriteJSON(msg)
		}
	}
}
