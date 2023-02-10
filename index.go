package zfeng

import (
	"github.com/raozhaofeng/zfeng/cache"
	"github.com/raozhaofeng/zfeng/conf"
	"github.com/raozhaofeng/zfeng/database"
	"github.com/raozhaofeng/zfeng/logs"
	"github.com/raozhaofeng/zfeng/router"
	"github.com/raozhaofeng/zfeng/validator"
)

var (
	// ConfInfo 配置文件
	ConfInfo *conf.Config
)

// NewApp 创建框架对象
func NewApp(confPath string) *router.Router {
	// 读取配置文件
	ConfInfo = conf.ReadConfigFile(confPath)

	// 初始化日志对象
	logs.InitializationLogger(ConfInfo.Logs.OutputPaths, ConfInfo.Debug)

	// 初始化缓存对象
	cache.InitializationConnPool(ConfInfo.Redis)

	// 初始化数据库
	database.InitializationDb(ConfInfo.Database)

	// 实例化验证器
	validator.InitializeValidator()

	// 启动线程
	return router.NewRoute(cache.RedisPool)
}
