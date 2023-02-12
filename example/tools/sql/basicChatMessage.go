package sql

import (
	"basic/tools/utils"
)

const BasicChatMessageTableName = "chat_message"
const BasicChatMessageTableComment = "聊天消息"
const CreateBasicChatMessage = `CREATE TABLE ` + BasicChatMessageTableName + ` (
	id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
	session_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '会话ID',
	role_id TINYINT NOT NULL DEFAULT 10 COMMENT '10用户发送管理 | 11临时用户发送管理 | 12管理发送用户 | 13管理发送临时用户 | 14用户发送用户',
	sender_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '发送者',
	receiver_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '接收者',
	type TINYINT NOT NULL DEFAULT 1 COMMENT '消息类型',
	data TEXT COMMENT '数据',
	extra VARCHAR(2048) NOT NULL DEFAULT '' COMMENT '额外数据',
	created_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间'
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + BasicChatMessageTableComment + `';`

const InsertBasicChatMessage = ``

var BasicChatMessage = &utils.InitTable{
	Name:        BasicChatMessageTableName,
	Comment:     BasicChatMessageTableComment,
	CreateTable: CreateBasicChatMessage,
	InsertTable: InsertBasicChatMessage,
}
