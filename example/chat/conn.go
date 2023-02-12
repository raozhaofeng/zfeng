package chat

import "github.com/gorilla/websocket"

// SetAdminConn 设置管理员连接
func (c *Chat) SetAdminConn(adminId int64, conn *websocket.Conn) {
	c.AdminConn[adminId] = conn
}

// DelAdminConn 删除管理员连接
func (c *Chat) DelAdminConn(adminId int64) {
	delete(c.AdminConn, adminId)
}

// SetUserConn 设置用户连接
func (c *Chat) SetUserConn(userId int64, conn *websocket.Conn) {
	c.UserConn[userId] = conn
}

// DelUserConn 删除用户连接
func (c *Chat) DelUserConn(userId int64) {
	delete(c.UserConn, userId)
}

// SetTouristConn 设置游客连接
func (c *Chat) SetTouristConn(userId int64, conn *websocket.Conn) {
	c.TouristConn[userId] = conn
}

// DelTouristConn 关闭游客连接
func (c *Chat) DelTouristConn(userId int64) {
	delete(c.TouristConn, userId)
}

// UserOnlineStatus 用户在线状态
func (c *Chat) UserOnlineStatus(userId int64) bool {
	if _, ok := c.UserConn[userId]; ok {
		return true
	}
	return false
}

// TouristOnlineStatus 临时用户在线状态
func (c *Chat) TouristOnlineStatus(userId int64) bool {
	if _, ok := c.TouristConn[userId]; ok {
		return true
	}
	return false
}
