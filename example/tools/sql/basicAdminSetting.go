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
(1, 1, '等级购买方式', 'select', 'buy_level_mode', 'premium', '[{"label": "补差模式", "value": "premium"}, {"label": "等价模式", "value": "equivalence"}]'),
(1, 1, '产品收益模式', 'select', 'product_earnings_mode', 'manual', '[{"label": "手动领取", "value": "manual"}, {"label": "自动领取", "value": "automatic"}]'),
(1, 1, '时区设置[参考网站: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones]', 'text', 'site_timezone', 'Asia/Shanghai', ''),
(1, 1, '下载设置', 'json', 'site_down', '{"android": "", "ios": ""}', '[[{"field": "android", "name": "上传安卓包", "type": "file"}, {"field": "ios", "name": "上传苹果包", "type": "file"}]]'),
(1, 1, 'Token设置', 'json', 'site_token', '{"key": "gvfsor2sj51tknuirzjhsose30oqgw0n", "only": true, "expire": 86400, "whitelist": "", "blacklist": ""}', '[[{"field": "key", "name": "密钥Key", "type": "text"}, {"field": "only", "name": "登陆唯一", "type": "select", "data": [{"label": "唯一登陆", "value": true}, {"label": "无限登陆", "value": false}]}, {"field": "expire", "name": "过期时间", "type": "number"}], [{"field": "whitelist", "name": "白名单", "type": "textarea"}], [{"field": "blacklist", "name": "黑名单", "type": "textarea"}]]'),

(1, 2, '首页轮播图(没有数据则不显示)', 'images', 'home_banner', '[{"label": "", "value": "/assets/images/banner/banner1.jpg"},{"label": "", "value": "/assets/images/banner/banner2.jpg"},{"label": "", "value": "/assets/images/banner/banner3.jpg"},{"label": "", "value": "/assets/images/banner/banner4.jpg"},{"label": "", "value": "/assets/images/banner/banner5.jpg"}]', ''),
(1, 2, '滚屏金额范围(最大值为0不显示)', 'json', 'home_money', '{"min": 100, "max": 10000}', '[[{"field": "min", "name": "金额最小值", "type": "number"}, {"field": "max", "name": "金额最大值", "type": "number"}]]'),
(1, 2, '站点介绍(语言字典 - 可显示富文本内容)', 'text', 'home_introduce', 'homeIntroduceText', ''),
(1, 2, '站点公告(语言字典 - 可显示富文本内容)', 'text', 'home_notice', 'homeNoticeText', ''),
(1, 2, '隐私设置(语言字典 - 可显示富文本内容)', 'text', 'home_privacy', 'privacyPolicyText', ''),
(1, 2, '服务协议(语言字典 - 可显示富文本内容)', 'text', 'home_protocol', 'serviceAgreementText', ''),
(1, 2, '客服链接(内置或三方客服 - 不填则不显示)', 'text', 'home_online', '/online', ''),
(1, 2, '底部导航设置(不填则不显示)', 'children', 'admin_tabs', '[{"label": "home", "router": "/", "icon": "/assets/images/tabs/home.png", "active_icon": "/assets/images/tabs/active_home.png"}, {"label": "product", "router": "/product", "icon": "/assets/images/tabs/shop.png", "active_icon": "/assets/images/tabs/active_shop.png"}, {"label": "order", "router": "/product/order",  "icon": "/assets/images/tabs/order.png", "active_icon": "/assets/images/tabs/active_order.png"}, {"label": "helpers", "router": "/service", "icon": "/assets/images/tabs/service.png", "active_icon": "/assets/images/tabs/active_service.png"}, {"label": "user", "router": "/user", "icon": "/assets/images/tabs/user.png", "active_icon": "/assets/images/tabs/active_user.png"}]', '[{"field": "label", "name": "名称(语言字典)", "type": "text"}, {"field": "router", "name": "路由", "type": "text"}, {"field": "icon", "name": "图标", "type": "icon"}, {"field": "active_icon", "name": "激活图标", "type": "icon"}]'),

(1, 3, '充值描述(语言字典 - 可显示富文本内容)', 'text', 'finance_deposit_tip', 'depositTip', ''),
(1, 3, '提现描述(语言字典 - 可显示富文本内容)', 'text', 'finance_withdraw_tip', 'withdrawTip', ''),
(1, 3, '充值范围设置', 'json', 'finance_deposit_range', '{"min": 100, "max": 10000}', '[[{"field": "min", "name": "充值最小值", "type": "number"}, {"field": "max", "name": "充值最大值", "type": "number"}]]'),
(1, 3, '提现范围设置', 'json', 'finance_withdraw_range', '{"min": 100, "max": 10000}', '[[{"field": "min", "name": "提现最小值", "type": "number"}, {"field": "max", "name": "提现最大值", "type": "number"}]]'),
(1, 3, '提现时间设置', 'children', 'finance_withdraw_times', '[{"sta_time": "00:00:00", "end_time": "23:59:59"}]', '[{"field": "sta_time", "name": "开始时间", "type": "timePicker"}, {"field": "end_time", "name": "结束时间", "type": "timePicker"}]'),
(1, 3, '提现次数设置', 'json', 'finance_withdraw_nums', '{"days": 7, "nums": 2}', '[[{"field": "days", "name": "间隔天数", "type": "number"}, {"field": "nums", "name": "提现次数", "type": "number"}]]'),
(1, 3, '提现手续费(%)', 'number', 'finance_withdraw_fee', '1', ''),
(1, 3, '钱包绑定账户(单类型个数)', 'number', 'financial_wallet_num', '1', ''),
(1, 3, '注册奖励', 'number', 'register_rewards', '5', ''),
(1, 3, '邀请奖励', 'number', 'invite_rewards', '100', ''),
(1, 3, '分销类型设置', 'checkbox', 'pyramid_items', '{"buyProduct": true, "productProfit": true, "taskProfit": true}', '[{"label": "购买产品分销", "value": "buyProduct"}, {"label": "产品利润分销", "value": "productProfit"}, {"label": "任务奖励分销", "value": "taskProfit"}]'),
(1, 3, '分销等级设置', 'children', 'pyramid_level', '[{"label": "pyramidLevel1", "value": 10}, {"label": "pyramidLevel2", "value": 5}, {"label": "pyramidLevel3", "value": 3}]', '[{"field": "label", "name": "名称(语言字典)", "type": "text"}, {"field": "value", "name": "收益比例(%)", "type": "number"}]'),

(1, 4, '主题模版', 'select', 'template', 'default', '[{"label": "默认模版", "value": "default"}]'),
(1, 4, '主题主色', 'color', 'color_primary', '#1976D2', ''),
(1, 4, '主题辅色', 'color', 'color_secondary', '#26A69A', ''),
(1, 4, '主题强色', 'color', 'color_accent', '#9C27B0', ''),
(1, 4, '基础配置', 'checkbox', 'template_basic', '{"update_password": true, "update_security": true, "update_avatar": true, "update_nickname": true, "update_email": true, "update_telephone": true, "update_sex": true, "update_birthday": true}', '[{"label": "开启更新密码", "value": "update_password"}, {"label": "开启更新安全密钥", "value": "update_security"}, {"label": "开启更新头像", "value": "update_avatar"}, {"label": "开启更新昵称", "value": "update_nickname"}, {"label": "开启更新邮箱", "value": "update_email"}, {"label": "开启更新手机号码", "value": "update_telephone"}, {"label": "开启更新性别", "value": "update_sex"}, {"label": "开启更新生日", "value": "update_birthday"}]'),
(1, 4, '登陆配置', 'checkbox', 'template_login', '{"show_register": true, "show_logo": true, "show_name": true}', '[{"label": "显示注册", "value": "show_register"}, {"label": "显示LOGO", "value": "show_logo"}, {"label": "显示名称", "value": "show_name"}]'),
(1, 4, '注册配置', 'checkbox', 'template_register', '{"telephone": true, "nickname": false, "email": false, "confirm_password": false, "security_key": true, "invite_code": true, "show_logo": true, "show_name": true}', '[{"label": "开启确认密码", "value": "confirm_password"}, {"label": "开启安全密钥", "value": "security_key"}, {"label": "开启邀请码", "value": "invite_code"}, {"label": "开启手机号码", "value": "telephone"}, {"label": "开启用户昵称", "value": "nickname"}, {"label": "开启用户邮箱", "value": "email"}, {"label": "显示LOGO", "value": "show_logo"}, {"label": "显示名称", "value": "show_name"}]'),
(1, 4, '钱包配置', 'checkbox', 'template_wallet', '{"update": true, "delete": true, "security_key": true, "withdraw_security_key": true, "withdraw_freeze": true, "withdraw_verify": false}', '[{"label": "开启修改", "value": "update"}, {"label": "开启删除", "value": "delete"}, {"label": "修改钱包账户密钥验证", "value": "security_key"}, {"label": "提现密钥验证", "value": "withdraw_security_key"}, {"label": "冻结不可提现", "value": "withdraw_freeze"}, {"label": "实名提现", "value": "withdraw_verify"}]'),
(1, 4, '实名配置', 'checkbox', 'template_verify', '{"real_name": true, "id_number": true, "email": true, "telephone": false, "photo_front": true, "photo_back": true, "photo_hold": true, "address": true}', '[{"label": "开启真实姓名", "value": "real_name"}, {"label": "开启证件号码", "value": "id_number"}, {"label": "开启手机号码", "value": "telephone"}, {"label": "开启邮箱", "value": "email"}, {"label": "开启证件照正面", "value": "photo_front"}, {"label": "开启证件照反面", "value": "photo_back"}, {"label": "开启手持证件照", "value": "photo_hold"}, {"label": "开启证件照地址", "value": "address"}]'),

(1, 5, '帮助背景图', 'image', 'service_image', '/assets/images/bg_service.png', ''),
(1, 5, '邀请背景图', 'image', 'invite_image', '/assets/images/bg_invite.jpg', ''),
(1, 5, '帮助中心', 'children', 'helpers', '[{"title": "helpersTitle1", "content": "helpersContent1"}]', '[{"field": "title", "name": "标题(语言字典)", "type": "text"}, {"field": "content", "name": "内容(语言字典)", "type": "text"}]'),
(1, 5, '联系方式', 'children', 'contacts', '[{"avatar": "/assets/images/contact/line.png", "link": "", "name": "line", "desc": "contactLine"}, {"avatar": "/assets/images/contact/whatsapp.png", "link": "", "name": "whatsapp", "desc": "contactWhatsapp"}, {"avatar": "/assets/images/contact/telegram.png", "link": "", "name": "telegram", "desc": "contactTelegram"}]', '[{"field": "avatar", "name": "图标", "type": "icon"}, {"field": "name", "name": "名称", "type": "text"}, {"field": "link", "name": "链接", "type": "text"}, {"field": "desc", "name": "描述(语言字典)", "type": "text"}]');`

var BasicAdminSetting = &utils.InitTable{
	Name:        AdminSettingTableName,
	Comment:     AdminSettingTableComment,
	CreateTable: CreateAdminSetting,
	InsertTable: InsertAdminSetting,
}
