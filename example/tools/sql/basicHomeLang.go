package sql

import (
	"basic/tools/utils"
)

const BasicHomeLangTableName = "lang"
const BasicHomeLangTableComment = "用户语言"
const CreateBasicHomeLang = `CREATE TABLE ` + BasicHomeLangTableName + ` (
	id         	 INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
    admin_id   	 INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '管理员ID',
    name       	 VARCHAR(50)  NOT NULL DEFAULT '' COMMENT '名称',
	alias		 VARCHAR(50)  NOT NULL DEFAULT '' COMMENT '别名',
    icon       	 VARCHAR(255) NOT NULL DEFAULT '' COMMENT '图标',
    sort       	 TINYINT      NOT NULL DEFAULT 99 COMMENT '排序',
    status     	 TINYINT      NOT NULL DEFAULT 10 COMMENT '状态 -1禁用｜10启用',
    data       	 VARCHAR(255) NOT NULL DEFAULT '' COMMENT '数据',
    created_at 	 INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间'
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + BasicHomeLangTableComment + `';`

const InsertBasicHomeLang = `INSERT INTO ` + BasicHomeLangTableName + `(admin_id, name, alias, icon, status) VALUES
(1, '简体中文', 'zh-CN', '/assets/images/country/china.png', 10),
(1, '繁体中文','zh-TW', '/assets/images/country/taiwan.png', 10),
(1, 'English','en-US', '/assets/images/country/usa.png', 10),
(1, '阿拉伯文', 'ar-AE', '/assets/images/country/united_arab_emirates.png', -1),
(1, '白俄罗斯文', 'be-BY', '/assets/images/country/belarus.png', -1),
(1, '保加利亚文', 'bg-BG', '/assets/images/country/bulgaria.png', -1),
(1, '捷克文', 'cs-CZ', '/assets/images/country/czech.png', -1),
(1, '丹麦文', 'da-DK', '/assets/images/country/denmark.png', -1),
(1, '德文', 'de-DE', '/assets/images/country/germany.png', -1),
(1, '希腊文', 'el-GR', '/assets/images/country/greece.png', -1),
(1, '西班牙文', 'es-ES', '/assets/images/country/spain.png', -1),
(1, '爱沙尼亚文', 'et-EE', '/assets/images/country/estonia.png', -1),
(1, '芬兰文', 'fi-FI', '/assets/images/country/finland.png', -1),
(1, '法文', 'fr-FR', '/assets/images/country/france.png', -1),
(1, '克罗地亚文', 'hr-HR', '/assets/images/country/croatia.png', -1),
(1, '匈牙利文', 'hu-HU', '/assets/images/country/hungary.png', -1),
(1, '冰岛文', 'is-IS', '/assets/images/country/iceland.png', -1),
(1, '意大利文', 'it-IT', '/assets/images/country/italy.png', -1),
(1, '日文', 'ja-JP', '/assets/images/country/japan.png', -1),
(1, '朝鲜文', 'ko-KR', '/assets/images/country/north_korea.png', -1),
(1, '立陶宛文', 'lt-LT', '/assets/images/country/lithuania.png', -1),
(1, '马其顿文', 'mk-MK', '/assets/images/country/macedonia.png', -1),
(1, '荷兰文', 'nl-NL', '/assets/images/country/netherlands.png', -1),
(1, '挪威文', 'no-NO', '/assets/images/country/norway.png', -1),
(1, '波兰文', 'pl-PL', '/assets/images/country/poland.png', -1),
(1, '葡萄牙文', 'pt-PT', '/assets/images/country/portugal.png', -1),
(1, '罗马尼亚文', 'ro-RO', '/assets/images/country/romania.png', -1),
(1, '俄文', 'ru-RU', '/assets/images/country/russia.png', -1),
(1, '克罗地亚文', 'sh-YU', '/assets/images/country/croatia.png', -1),
(1, '斯洛伐克文', 'sk-SK', '/assets/images/country/slovakia.png', -1),
(1, '斯洛文尼亚文', 'sl-SI', '/assets/images/country/slovenia.png', -1),
(1, '阿尔巴尼亚文', 'sq-AL', '/assets/images/country/albania.png', -1),
(1, '瑞典文', 'sv-SE', '/assets/images/country/sweden.png', -1),
(1, '泰文', 'th-TH', '/assets/images/country/thailand.png', -1),
(1, '土耳其文', 'tr-TR', '/assets/images/country/turkey.png', -1),
(1, '乌克兰文', 'uk-UA', '/assets/images/country/ukraine.png', -1),
(1, '拉托维亚文', 'lv-LV', '/assets/images/country/latvia.png', -1),
(1, '塞尔维亚文', 'sr-YU', '/assets/images/country/serbia.png', -1),
(1, '希伯来文', 'iw-IL', '/assets/images/country/israel.png', -1);`

var BasicHomeLang = &utils.InitTable{
	Name:        BasicHomeLangTableName,
	Comment:     BasicHomeLangTableComment,
	CreateTable: CreateBasicHomeLang,
	InsertTable: InsertBasicHomeLang,
}
