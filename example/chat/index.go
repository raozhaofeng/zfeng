package chat

import "github.com/gorilla/websocket"

// Manager 聊天管理
var Manager = &Chat{
	AdminConn:   map[int64]*websocket.Conn{},
	UserConn:    map[int64]*websocket.Conn{},
	TouristConn: map[int64]*websocket.Conn{},
}

const (
	MessageRoleSystemToAdmin  int64 = 0  //	系统发送管理
	MessageRoleUserToAdmin    int64 = 10 //	用户发送管理
	MessageRoleTouristToAdmin int64 = 11 //	临时用户发送管理
	MessageRoleAdminToUser    int64 = 12 //	管理发送用户
	MessageRoleAdminToTourist int64 = 13 //	管理发送临时用户
	MessageRoleMainUserToUser int64 = 14 //	主用户发送用户
	MessageRoleUserToMainUser int64 = 15 //	用户发送主用户
)

// Chat 聊天
type Chat struct {
	AdminConn   map[int64]*websocket.Conn //	管理连接
	UserConn    map[int64]*websocket.Conn //	用户连接
	TouristConn map[int64]*websocket.Conn //	临时用户
}
