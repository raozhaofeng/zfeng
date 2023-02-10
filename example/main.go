package main

import (
	"basic/home"
	"basic/models"
	"github.com/raozhaofeng/zfeng"
)

func main() {
	// 启动后台接口
	adminApp := zfeng.NewApp(models.AccessLogsTypeAdmin, "./")
	go adminApp.ListenAndServe("0.0.0.0:8001")

	// 启动前台接口
	homeApp := zfeng.NewApp(models.AccessLogsTypeHome, "./") //	初始化
	homeApp.SetRouteHandle(home.Router())                    //	载入前台路由
	homeApp.SetCallbackAccessFunc(models.RouteAccessFunc)    //	设置访问日志
	homeApp.ListenAndServe("0.0.0.0:8002")                   //	启动监听
}
