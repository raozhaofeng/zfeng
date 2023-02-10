package main

import (
	"github.com/raozhaofeng/zfeng"
	"github.com/raozhaofeng/zfeng/example/home"
)

func main() {
	// 启动后台接口
	adminApp := zfeng.NewApp("./")
	go adminApp.ListenAndServe("0.0.0.0:8001")

	// 启动前台接口
	homeApp := zfeng.NewApp("./").SetRouteHandle(home.Router())
	homeApp.ListenAndServe("0.0.0.0:8002")
}
