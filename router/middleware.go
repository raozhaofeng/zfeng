package router

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// SetRouteHandle 设置路由函数
func (c *Router) SetRouteHandle(routeHandle []*Handle) *Router {
	for _, handle := range routeHandle {
		handleController := c.middlewareToken(&Handle{
			Name:      handle.Name,
			Method:    handle.Method,
			Route:     handle.Route,
			RouteAuth: handle.RouteAuth,
			TokenAuth: handle.TokenAuth,
		}, httprouter.Handle(handle.Controller))

		//	添加路由函数
		c.httpRouter.Handle(handle.Method, handle.Route, handleController)
	}
	return c
}

// middlewareToken Token中间件
func (c *Router) middlewareToken(handle *Handle, next httprouter.Handle) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		crossDomainRequest(writer, request)
		var claims *Claims
		if handle.TokenAuth {
			rds := c.redisPool.Get()
			defer rds.Close()

			//	验证Token
			claims = TokenManager.Verify(rds, request)
			if claims == nil {
				c.StatusUnauthorized(writer)
				return
			}

			//	如果需要验证路由
			if handle.RouteAuth && !TokenManager.AuthRouter(rds, claims.AdminId, handle.Route) {
				c.StatusUnauthorized(writer)
				return
			}

			//	复制Token信息传递
			ctx := context.WithValue(request.Context(), ClaimsKey, claims)
			request = request.WithContext(ctx)
		}

		//	正常返回
		if c.CallbackAccessFunc != nil {
			c.CallbackAccessFunc(handle, request, claims)
		}
		next(writer, request, params)
	}
}

// StatusUnauthorized 返回没有权限
func (c *Router) StatusUnauthorized(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusUnauthorized)
	_, _ = writer.Write([]byte("401 Unauthorized"))
}
