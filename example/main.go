package main

import (
	"github.com/raozhaofeng/zfeng"
	"github.com/raozhaofeng/zfeng/example/home"
)

func main() {
	homeApp := zfeng.NewApp(1, "./")
	homeApp.SetRouteHandle(home.Router())
	homeApp.ListenAndServe("0.0.0.0:8010")
}
