package sql

import (
	"basic/tools/utils"
)

const AdminSettingTableName = "admin_setting"
const AdminSettingTableComment = "管理设置"
const CreateAdminSetting = `CREATE TABLE ` + AdminSettingTableName + ` (
	id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
	admin_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '管理员ID',
	group_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '组ID',
	name VARCHAR(255) NOT NULL DEFAULT '' COMMENT '名称',
	type VARCHAR(64) NOT NULL DEFAULT '' COMMENT '类型',
	field VARCHAR(64) NOT NULL DEFAULT '' COMMENT '健名',
	value TEXT COMMENT '健值',
	data TEXT COMMENT '数据',
	KEY admin_setting_key (field)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + AdminSettingTableComment + `';`

const InsertAdminSetting = `INSERT INTO ` + AdminSettingTableName + `(admin_id, group_id, name, type, field, value, data) VALUES
(1, 1, '站点名称(语言字典)', 'text', 'site_name', 'siteName', ''),
(1, 1, '站点LOGO', 'image', 'site_logo', '/assets/images/logo.png', ''),
(1, 1, '模版缓存时间(s)', 'number', 'cache_expire', '60', ''),
(1, 1, '时区设置[参考网站: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones]', 'text', 'site_timezone', 'Asia/Shanghai', ''),
(1, 1, '下载设置', 'json', 'site_down', '{"android": "", "ios": ""}', '[[{"field": "android", "name": "上传安卓包", "type": "file"}, {"field": "ios", "name": "上传苹果包", "type": "file"}]]'),
(1, 1, 'Token设置', 'json', 'site_token', '{"key": "gvfsor2sj51tknuirzjhsose30oqgw0n", "only": true, "expire": 86400, "whitelist": "", "blacklist": ""}', '[[{"field": "key", "name": "密钥Key", "type": "text"}, {"field": "only", "name": "登陆唯一", "type": "select", "data": [{"label": "唯一登陆", "value": true}, {"label": "无限登陆", "value": false}]}, {"field": "expire", "name": "过期时间", "type": "number"}], [{"field": "whitelist", "name": "白名单", "type": "textarea"}], [{"field": "blacklist", "name": "黑名单", "type": "textarea"}]]');`

var BasicAdminSetting = &utils.InitTable{
	Name:        AdminSettingTableName,
	Comment:     AdminSettingTableComment,
	CreateTable: CreateAdminSetting,
	InsertTable: InsertAdminSetting,
}
