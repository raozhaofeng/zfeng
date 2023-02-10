package database

import (
	"database/sql"
	"github.com/raozhaofeng/zfeng/database/define"
	"github.com/raozhaofeng/zfeng/database/mysql"
)

// DbPool 数据管理
var DbPool *Pool

// Pool 数据池
type Pool struct {
	db *sql.DB
}

// NewDb 创建新的接口Db
func (c *Pool) NewDb(tx *sql.Tx) define.Db {
	if tx == nil {
		return mysql.NewDb(tx, c.db)
	}
	return mysql.NewDb(tx, nil)
}

// GetDb 获取Db
func (c *Pool) GetDb() *sql.DB {
	return c.db
}

// GetTx 获取Tx
func (c *Pool) GetTx() *sql.Tx {
	tx, _ := c.db.Begin()
	return tx
}

// InitializationDb 初始化数据库
func InitializationDb(conf *define.DatabaseConf) {
	DbPool = &Pool{mysql.InitializationDb(conf)}
}
