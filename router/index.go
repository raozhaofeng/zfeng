package router

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Router struct {
	redisPool          *redis.Pool
	httpRouter         *httprouter.Router                                    // 路由实例
	CallbackAccessFunc func(handle *Handle, r *http.Request, claims *Claims) //	访问日志
}

// NewRoute 创建路由
func NewRoute(redisPool *redis.Pool) *Router {
	httpRouter := httprouter.New()
	// 开启跨域请求
	httpRouter.GlobalOPTIONS = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		crossDomainRequest(writer, request)
		writer.WriteHeader(http.StatusNoContent)
	})

	// 全局异常拦截
	httpRouter.PanicHandler = func(writer http.ResponseWriter, request *http.Request, i interface{}) {
		fmt.Println(i)
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte("500 Internal Server Error"))
	}

	return &Router{
		redisPool:  redisPool,
		httpRouter: httpRouter,
	}
}

// ServeFiles 开启静态资源
func (c *Router) ServeFiles(filePath string) *Router {
	c.httpRouter.ServeFiles("/"+filePath+"/*filepath", http.Dir("./"+filePath))
	return c
}

// SetCallbackAccessFunc 设置访问日志函数
func (c *Router) SetCallbackAccessFunc(fun func(handle *Handle, r *http.Request, claims *Claims)) *Router {
	c.CallbackAccessFunc = fun
	return c
}

// ListenAndServe 监听服务
func (c *Router) ListenAndServe(addr string) {
	fmt.Println("Listen", addr, "Successful")
	err := http.ListenAndServe(addr, c.httpRouter)
	if err != nil {
		panic(err)
	}
}

// crossDomainRequest 设置跨域请求
func crossDomainRequest(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Credentials", "true")
	writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers", "Token, Token-Key, Content-Type, If-Match, If-Modified-Since, If-None-Match, If-Unmodified-Since, X-Requested-With")
	origin := request.Header.Get("origin")
	if origin == "" {
		origin = "*"
	}
	writer.Header().Set("Access-Control-Allow-Origin", origin)
}
