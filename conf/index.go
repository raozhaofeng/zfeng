package conf

import (
	"github.com/raozhaofeng/zfeng/database/define"
	"github.com/spf13/viper"
	"time"
)

// RedisConf Redis配置文件
type RedisConf struct {
	Network         string        //	网络
	Server          string        //	地址
	Port            int           //	端口
	Pass            string        //	密码
	Dbname          int           // 	库名
	ConnectTimeout  time.Duration //	连接超时时间
	ReadTimeout     time.Duration //	读取超时时间
	WriteTimeout    time.Duration //	写入超时时间
	MaxOpenConn     int           // 	设置最大连接数
	ConnMaxIdleTime time.Duration // 	空闲连接超时
	MaxIdleConn     int           // 	最大空闲连接数
	Wait            bool          // 	如果超过最大连接数是否等待
}

// LogsConf 日志配置文件
type LogsConf struct {
	OutputPaths []string //	日志输出路径
}

// Config 配置文件
type Config struct {
	Debug    bool                 // 是否调试
	Database *define.DatabaseConf // 数据库配置
	Redis    *RedisConf           // 缓存配置文件
	Logs     *LogsConf            // 日志配置
}

// ReadConfigFile 读取配置文件
func ReadConfigFile(confPath string) *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(confPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	return &Config{
		Debug: viper.GetBool("Debug"),
		Database: &define.DatabaseConf{
			Drive:           viper.GetString("Database.Drive"),
			User:            viper.GetString("Database.User"),
			Pass:            viper.GetString("Database.Pass"),
			Dbname:          viper.GetString("Database.Dbname"),
			Network:         viper.GetString("Database.Network"),
			Server:          viper.GetString("Database.Server"),
			Port:            viper.GetInt("Database.Port"),
			ConnMaxLifeTime: 0,
			MaxOpenConn:     0,
			ConnMaxIdleTime: 0,
			MaxIdleConn:     0,
			Params:          map[string]any{},
		},
		Redis: &RedisConf{
			Network:         viper.GetString("Redis.Network"),
			Server:          viper.GetString("Redis.Server"),
			Port:            viper.GetInt("Redis.Port"),
			Pass:            viper.GetString("Redis.Pass"),
			Dbname:          viper.GetInt("Redis.Dbname"),
			ConnectTimeout:  30,
			ReadTimeout:     30,
			WriteTimeout:    30,
			MaxOpenConn:     0,
			ConnMaxIdleTime: 30,
			MaxIdleConn:     0,
			Wait:            false,
		},
		Logs: &LogsConf{
			OutputPaths: viper.GetStringSlice("Logs.OutputPaths"),
		},
	}
}
