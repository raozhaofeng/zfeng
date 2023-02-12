package sql

import (
	"basic/tools/utils"
)

const HomeProductTableName = "product"
const HomeProductTableComment = "产品商品"
const CreateHomeProduct = `CREATE TABLE ` + HomeProductTableName + ` (
	id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
	admin_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '管理员ID',
	category_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '类目ID',
	assets_id INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '资产ID',
	name VARCHAR(255) NOT NULL DEFAULT '' COMMENT '标题',
	images VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '图片列表',
	money DECIMAL(12, 2) NOT NULL DEFAULT 0 COMMENT '金额',
	sort SMALLINT UNSIGNED NOT NULL DEFAULT 99 COMMENT '排序',
	status TINYINT NOT NULL DEFAULT 10 COMMENT '状态 -2删除 -1禁用 10启用',
	recommend TINYINT NOT NULL DEFAULT -1 COMMENT '推荐 -1关闭 10推荐',
	sales INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '销售量',
	nums TINYINT NOT NULL DEFAULT -1 COMMENT '限购 -1 无限制',
	mode TINYINT NOT NULL DEFAULT 1 COMMENT '模式 1返本 2返息',
	used INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '已用',
	total INT UNSIGNED NOT NULL DEFAULT 1000 COMMENT '总数',
	data VARCHAR(2048) NOT NULL DEFAULT '' COMMENT '数据',
	describes VARCHAR(255) NOT NULL DEFAULT '' COMMENT '描述',
	updated_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间',
	created_at INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间'
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='` + HomeProductTableComment + `';`

const InsertHomeProduct = `INSERT INTO ` + HomeProductTableName + `(admin_id, category_id, name, images, money, recommend, data, describes) VALUES
(1, 1, 'productTeddyDog', '[{"label": "", "value": "/assets/images/product/teddy1.jpeg"}, {"label": "", "value": "/assets/images/product/teddy2.jpeg"}, {"label": "", "value": "/assets/images/product/teddy3.jpeg"}, {"label": "", "value": "/assets/images/product/teddy4.jpeg"}]', 100, 10, '{"expire": 3, "interval": 1, "returns": 1.53}', 'productTeddyDogDescribe'),
(1, 1, 'productCorgiDog', '[{"label": "", "value": "/assets/images/product/corgi1.jpeg"}, {"label": "", "value": "/assets/images/product/corgi2.jpeg"}, {"label": "", "value": "/assets/images/product/corgi3.jpeg"}, {"label": "", "value": "/assets/images/product/corgi4.jpeg"}]', 300, 10, '{"expire": 30, "interval": 1, "returns": 8.68}', 'productCorgiDogDescribe'),
(1, 2, 'productBlueCat', '[{"label": "", "value": "/assets/images/product/blue1.jpeg"}, {"label": "", "value": "/assets/images/product/blue2.jpeg"}, {"label": "", "value": "/assets/images/product/blue3.jpeg"}, {"label": "", "value": "/assets/images/product/blue4.jpeg"}]', 100, 10, '{"expire": 3, "interval": 1, "returns": 1.53}', 'productBlueCatDescribe'),
(1, 2, 'productPuppetCat', '[{"label": "", "value": "/assets/images/product/puppet1.jpeg"}, {"label": "", "value": "/assets/images/product/puppet2.jpeg"}, {"label": "", "value": "/assets/images/product/puppet3.jpeg"}, {"label": "", "value": "/assets/images/product/puppet4.jpeg"}]', 500, 10, '{"expire": 30, "interval": 1, "returns": 10.88}', 'productPuppetCatDescribe');`

var BasicHomeProduct = &utils.InitTable{
	Name:        HomeProductTableName,
	Comment:     HomeProductTableComment,
	CreateTable: CreateHomeProduct,
	InsertTable: InsertHomeProduct,
}
