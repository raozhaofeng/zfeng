package sql

import (
	"basic/tools/utils"
)

const BasicChatSessionTableName = "chat_session"
const BasicChatSessionTableComment = "聊天会话"
const CreateBasicChatSession = `CREATE TABLE ` + BasicChatSessionTableName + ` (
	id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
	session_key CHAR(32) NOT NULL DEFAULT '' UNIQUE KEY COMMENT '会话key',
	admin_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '管理ID',
	main_user INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '主用户',
	main_user_unread SMALLINT NOT NULL DEFAULT 0 COMMENT '主用户未读',
	client_user INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '客户用户',
	client_user_unread SMALLINT NOT NULL DEFAULT 0 COMMENT '客户用户未读',
	type TINYINT NOT NULL DEFAULT 10 COMMENT '10管理会话[管理对用户] 11临时会话[管理对临时] 20用户会话[用户对用户]',
	data TEXT COMMENT '最后消息内容',
	ip4 INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'IP4地址',
    updated_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间'
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + BasicChatSessionTableComment + `';`

const InsertBasicChatSession = ``

var BasicChatSession = &utils.InitTable{
	Name:        BasicChatSessionTableName,
	Comment:     BasicChatSessionTableComment,
	CreateTable: CreateBasicChatSession,
	InsertTable: InsertBasicChatSession,
}
