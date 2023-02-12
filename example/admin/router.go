package admin

import (
	"basic/admin/controllers/chat"
	"basic/admin/controllers/index"
	"basic/admin/controllers/manager/logs"
	"basic/admin/controllers/manager/manager"
	"basic/admin/controllers/manager/menu"
	"basic/admin/controllers/manager/role"
	"basic/admin/controllers/product/category"
	"basic/admin/controllers/product/order"
	"basic/admin/controllers/product/product"
	settingAssets "basic/admin/controllers/settings/assets"
	"basic/admin/controllers/settings/country"
	"basic/admin/controllers/settings/dictionary"
	"basic/admin/controllers/settings/lang"
	"basic/admin/controllers/settings/level"
	"basic/admin/controllers/settings/payment"
	"basic/admin/controllers/settings/setting"
	"basic/admin/controllers/user/assets"
	"basic/admin/controllers/user/bill"
	"basic/admin/controllers/user/user"
	"basic/admin/controllers/user/verify"
	"basic/admin/controllers/user/wallet/account"
	"basic/admin/controllers/user/wallet/deposit"
	"basic/admin/controllers/user/wallet/withdraw"
	"github.com/raozhaofeng/zfeng/router"
)

func Router() []*router.Handle {
	return []*router.Handle{
		//	基础路由
		router.NewRouterTokenHandle("数据库表", "POST", "/tables/index", index.Tables),
		router.NewRouterTokenHandle("数据表信息", "POST", "/tables/columns", index.Columns),
		router.NewRouterTokenHandle("上传文件", "POST", "/upload", index.Upload),
		router.NewHandle("生成验证码", "GET", "/captcha/generate", index.GenerateCaptcha),
		router.NewHandle("显示验证码", "GET", "/captcha/image", index.ImageCaptcha),
		router.NewHandle("管理员登陆", "POST", "/login", index.Login),
		router.NewRouterTokenHandle("首页信息", "POST", "/index", index.Index),
		router.NewRouterTokenHandle("管理信息", "POST", "/info", index.Info),
		router.NewRouterTokenHandle("更新密码", "POST", "/update/password", index.UpdatePassword),
		router.NewRouterTokenHandle("更新信息", "POST", "/update", index.Update),

		//	聊天路由
		router.NewHandle("管理Websocket连接", "GET", "/ws", chat.Websocket),
		router.NewRouterTokenHandle("会话列表", "POST", "/chat/index", chat.Index),
		router.NewRouterTokenHandle("会话信息", "POST", "/chat/info", chat.Info),
		router.NewRouterTokenHandle("发送消息", "POST", "/chat/send", chat.Send),
		router.NewRouterTokenHandle("清除未读", "POST", "/chat/unread", chat.Unread),
		router.NewRouterTokenHandle("会话消息", "POST", "/chat/message", chat.Message),

		//	管理员列表
		router.NewRouterTokenHandle("管理员列表", "POST", "/manage/index", manager.Index),
		router.NewRouterTokenHandle("管理员新增", "POST", "/manage/create", manager.Create),
		router.NewRouterTokenHandle("管理员更新", "POST", "/manage/update", manager.Update),
		router.NewRouterTokenHandle("管理员删除", "POST", "/manage/delete", manager.Delete),

		// 角色
		router.NewRouterTokenHandle("角色数组", "POST", "/role/roles", role.RolesList),
		router.NewRouterTokenHandle("权限数组", "POST", "/role/permissions", role.PermissionsList),
		router.NewRouterTokenHandle("角色列表", "POST", "/role/index", role.Index),
		router.NewRouterTokenHandle("角色更新", "POST", "/role/update", role.Update),
		router.NewRouterTokenHandle("角色新增", "POST", "/role/create", role.Create),
		router.NewRouterTokenHandle("角色删除", "POST", "/role/delete", role.Delete),

		// 菜单管理
		router.NewRouterTokenHandle("菜单列表", "POST", "/menu/index", menu.Index),
		router.NewRouterTokenHandle("菜单更新", "POST", "/menu/update", menu.Update),
		// 操作日志
		router.NewRouterTokenHandle("操作日志", "POST", "/logs/index", logs.Index),

		// 用户等级
		router.NewRouterTokenHandle("等级管理", "POST", "/level/index", level.Index),
		router.NewRouterTokenHandle("等级更新", "POST", "/level/update", level.Update),
		router.NewRouterTokenHandle("等级删除", "POST", "/level/delete", level.Delete),
		router.NewRouterTokenHandle("等级新增", "POST", "/level/create", level.Create),
		router.NewRouterTokenHandle("等级数组", "POST", "/level/levels", level.LevelsList),

		// 钱包支付方式
		router.NewRouterTokenHandle("支付管理", "POST", "/payment/index", payment.Index),
		router.NewRouterTokenHandle("支付更新", "POST", "/payment/update", payment.Update),
		router.NewRouterTokenHandle("支付删除", "POST", "/payment/delete", payment.Delete),
		router.NewRouterTokenHandle("支付新增", "POST", "/payment/create", payment.Create),
		router.NewRouterTokenHandle("支付提现类型", "POST", "/payment/withdraw", payment.WithdrawsList),

		// 平台资产
		router.NewRouterTokenHandle("平台资产管理", "POST", "/assets/index", settingAssets.Index),
		router.NewRouterTokenHandle("平台资产数组", "POST", "/assets/list", settingAssets.List),
		router.NewRouterTokenHandle("平台资产新增", "POST", "/assets/create", settingAssets.Create),
		router.NewRouterTokenHandle("平台资产更新", "POST", "/assets/update", settingAssets.Update),
		router.NewRouterTokenHandle("平台资产删除", "POST", "/assets/delete", settingAssets.Delete),

		// 用户国家
		router.NewRouterTokenHandle("国家管理", "POST", "/country/index", country.Index),
		router.NewRouterTokenHandle("国家更新", "POST", "/country/update", country.Update),
		router.NewRouterTokenHandle("国家删除", "POST", "/country/delete", country.Delete),
		router.NewRouterTokenHandle("国家新增", "POST", "/country/create", country.Create),

		// 用户语言
		router.NewRouterTokenHandle("语言管理", "POST", "/lang/index", lang.Index),
		router.NewRouterTokenHandle("语言更新", "POST", "/lang/update", lang.Update),
		router.NewRouterTokenHandle("语言删除", "POST", "/lang/delete", lang.Delete),
		router.NewRouterTokenHandle("语言新增", "POST", "/lang/create", lang.Create),

		// 语言字典
		router.NewRouterTokenHandle("语言字典管理", "POST", "/dictionary/index", dictionary.Index),
		router.NewRouterTokenHandle("语言字典更新", "POST", "/dictionary/update", dictionary.Update),
		router.NewRouterTokenHandle("语言字典删除", "POST", "/dictionary/delete", dictionary.Delete),
		router.NewRouterTokenHandle("语言字典新增", "POST", "/dictionary/create", dictionary.Create),
		router.NewRouterTokenHandle("语言字典下载", "POST", "/dictionary/download", dictionary.Download),
		router.NewRouterTokenHandle("语言字典上传", "POST", "/dictionary/upload", dictionary.Upload),

		// 管理员配置
		router.NewRouterTokenHandle("管理配置列表", "POST", "/setting/index", setting.Index),
		router.NewRouterTokenHandle("管理配置更新", "POST", "/setting/update", setting.Update),

		// 前台用户
		router.NewRouterTokenHandle("用户管理", "POST", "/user/index", user.Index),
		router.NewRouterTokenHandle("用户更新", "POST", "/user/update", user.Update),
		router.NewRouterTokenHandle("用户删除", "POST", "/user/delete", user.Delete),
		router.NewRouterTokenHandle("用户新增", "POST", "/user/create", user.Create),
		router.NewRouterTokenHandle("用户关系", "POST", "/user/relation", user.Relation),
		router.NewRouterTokenHandle("用户加减款", "POST", "/user/amount", user.Amount),

		// 用户资产
		router.NewRouterTokenHandle("用户资产", "POST", "/user/assets/index", assets.Index),
		router.NewRouterTokenHandle("资产更新", "POST", "/user/assets/update", assets.Update),
		router.NewRouterTokenHandle("资产删除", "POST", "/user/assets/delete", assets.Delete),
		router.NewRouterTokenHandle("资产新增", "POST", "/user/assets/create", assets.Create),
		router.NewRouterTokenHandle("资产加减款", "POST", "/user/assets/amount", assets.Amount),

		// 用户账单
		router.NewRouterTokenHandle("用户账单列表", "POST", "/user/bill/index", bill.Index),
		router.NewRouterTokenHandle("更新用户账单", "POST", "/user/bill/update", bill.Update),
		router.NewRouterTokenHandle("用户账单类型", "POST", "/user/bill/types", bill.TypesList),

		// 用户认证管理
		router.NewRouterTokenHandle("用户认证列表", "POST", "/user/verify/index", verify.Index),
		router.NewRouterTokenHandle("用户认证删除", "POST", "/user/verify/delete", verify.Delete),
		router.NewRouterTokenHandle("用户认证新增", "POST", "/user/verify/create", verify.Create),
		router.NewRouterTokenHandle("用户认证审核", "POST", "/user/verify/status", verify.Status),

		// 用户提现账户
		router.NewRouterTokenHandle("用户提现账户列表", "POST", "/wallet/account/index", account.Index),
		router.NewRouterTokenHandle("用户提现账户更新", "POST", "/wallet/account/update", account.Update),
		router.NewRouterTokenHandle("用户提现账户删除", "POST", "/wallet/account/delete", account.Delete),
		router.NewRouterTokenHandle("用户提现账户新增", "POST", "/wallet/account/create", account.Create),

		// 用户充值管理
		router.NewRouterTokenHandle("用户充值列表", "POST", "/wallet/deposit/index", deposit.Index),
		router.NewRouterTokenHandle("用户充值更新", "POST", "/wallet/deposit/update", deposit.Update),
		router.NewRouterTokenHandle("用户充值删除", "POST", "/wallet/deposit/delete", deposit.Delete),
		router.NewRouterTokenHandle("用户充值新增", "POST", "/wallet/deposit/create", deposit.Create),
		router.NewRouterTokenHandle("用户充值审核", "POST", "/wallet/deposit/status", deposit.Status),

		// 用户提现管理
		router.NewRouterTokenHandle("用户提现列表", "POST", "/wallet/withdraw/index", withdraw.Index),
		router.NewRouterTokenHandle("用户提现更新", "POST", "/wallet/withdraw/update", withdraw.Update),
		router.NewRouterTokenHandle("用户提现删除", "POST", "/wallet/withdraw/delete", withdraw.Delete),
		router.NewRouterTokenHandle("用户提现新增", "POST", "/wallet/withdraw/create", withdraw.Create),
		router.NewRouterTokenHandle("用户提现审核", "POST", "/wallet/withdraw/status", withdraw.Status),

		// 产品分类管理
		router.NewRouterTokenHandle("产品分类列表", "POST", "/product/category/index", category.Index),
		router.NewRouterTokenHandle("产品分类新增", "POST", "/product/category/create", category.Create),
		router.NewRouterTokenHandle("产品分类更新", "POST", "/product/category/update", category.Update),
		router.NewRouterTokenHandle("产品分类删除", "POST", "/product/category/delete", category.Delete),

		// 产品商品管理
		router.NewRouterTokenHandle("产品商品列表", "POST", "/product/index/index", product.Index),
		router.NewRouterTokenHandle("产品商品新增", "POST", "/product/index/create", product.Create),
		router.NewRouterTokenHandle("产品商品更新", "POST", "/product/index/update", product.Update),
		router.NewRouterTokenHandle("产品商品删除", "POST", "/product/index/delete", product.Delete),

		// 产品订单管理
		router.NewRouterTokenHandle("产品订单列表", "POST", "/product/order/index", order.Index),
		router.NewRouterTokenHandle("产品订单新增", "POST", "/product/order/create", order.Create),
		router.NewRouterTokenHandle("产品订单更新", "POST", "/product/order/update", order.Update),
		router.NewRouterTokenHandle("产品订单删除", "POST", "/product/order/delete", order.Delete),
	}
}
