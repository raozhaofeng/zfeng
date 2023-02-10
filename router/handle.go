package router

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// ControllerFunc 控制器方法名
type ControllerFunc func(w http.ResponseWriter, r *http.Request, p httprouter.Params)

type Handle struct {
	Name       string         //	名称
	Method     string         // 	类型
	Route      string         // 	路由
	Controller ControllerFunc //	控制器
	RouteAuth  bool           //	路由验证
	TokenAuth  bool           //	Token验证
}

// NewHandle 创建路由
func NewHandle(name, method, route string, controllerFunc ControllerFunc) *Handle {
	return &Handle{
		Name: name, Method: method, Route: route, Controller: controllerFunc, RouteAuth: false, TokenAuth: false,
	}
}

// NewTokenHandle 创建Token验证
func NewTokenHandle(name, method, route string, controllerFunc ControllerFunc) *Handle {
	return &Handle{
		Name: name, Method: method, Route: route, Controller: controllerFunc, RouteAuth: false, TokenAuth: true,
	}
}

// NewRouterTokenHandle 创建路由Token验证
func NewRouterTokenHandle(name, method, route string, controllerFunc ControllerFunc) *Handle {
	return &Handle{
		Name: name, Method: method, Route: route, Controller: controllerFunc, RouteAuth: true, TokenAuth: true,
	}
}
