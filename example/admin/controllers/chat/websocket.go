package chat

import (
	"basic/chat"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/cache"
	"github.com/raozhaofeng/zfeng/router"
	"net/http"
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Resolve cross-domain problems
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Websocket 管理员websocket 连接
func Websocket(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	rds := cache.RedisPool.Get()
	defer rds.Close()

	claims := router.TokenManager.Verify(rds, r)
	if claims == nil {
		_ = conn.Close()
		return
	}

	//	设置当前连接到管理员
	chat.Manager.SetAdminConn(claims.AdminId, conn)
	for {
		_, _, err = conn.ReadMessage()
		if err != nil {
			//	读取消息失败, 断开连接
			chat.Manager.DelAdminConn(claims.AdminId)
			return
		}
	}
}
