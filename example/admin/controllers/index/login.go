package index

import (
	"basic/models"
	"github.com/dchest/captcha"
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/cache"
	"github.com/raozhaofeng/zfeng/router"
	"github.com/raozhaofeng/zfeng/utils"
	"github.com/raozhaofeng/zfeng/utils/body"
	"github.com/raozhaofeng/zfeng/validator"
	"net/http"
	"time"
)

type loginParams struct {
	UserName     string `json:"username" validate:"required"`
	Password     string `json:"password" validate:"required"`
	CaptchaId    string `json:"captcha_id" validate:"required"`
	CaptchaValue string `json:"captcha_value" validate:"required"`
}

type userInfo struct {
	Id         int64   `json:"id"`          //	管理员ID
	Username   string  `json:"username"`    //	用户名
	Nickname   string  `json:"nickname"`    //	昵称
	Email      string  `json:"email"`       //	邮箱
	Avatar     string  `json:"avatar"`      //	头像
	Money      float64 `json:"money"`       //	金额
	Data       string  `json:"data"`        //	数据
	InviteCode string  `json:"invite_code"` //	邀请码
	OnlineNums int64   `json:"online_nums"` //	在线人数
	UnreadNums int64   `json:"unread_nums"` //	未读消息
	ExpiredAt  int64   `json:"expired_at"`  //	过期时间
	UpdatedAt  int64   `json:"updatedAt"`   //	更新时间
}

type loginData struct {
	Menu       []*models.AdminMenuList `json:"menu"`        //	菜单
	RouterList []string                `json:"router_list"` //	路由列表
	UserInfo   *userInfo               `json:"info"`        //	用户信息
	Token      string                  `json:"token"`       //	Token
	TokenKey   string                  `json:"token_key"`   // TokenKey
}

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(loginParams)
	_ = body.ReadJSON(r, params)

	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	// 验证码是否正确
	if !captcha.VerifyString(params.CaptchaId, params.CaptchaValue) {
		body.ErrorJSON(w, "验证码错误", -1)
		return
	}

	//	获取用户信息 [是否存在，密码是否正确，是否过期]
	nowTime := time.Now().Unix()
	adminModel := models.NewAdminUser(nil)
	adminModel.AndWhere("username=?", params.UserName).AndWhere("status=?", models.AdminUserStatusActivate)
	adminInfo := adminModel.FindOne()
	if adminInfo == nil || adminInfo.Password != utils.PasswordEncrypt(params.Password) {
		body.ErrorJSON(w, "账号或密码错误", -1)
		return
	}
	if adminInfo.ExpiredAt > 0 && adminInfo.ExpiredAt < nowTime {
		body.ErrorJSON(w, "账号已过期, 请联系管理员", -1)
		return
	}

	//	获取在线人数
	var onlineNums int64

	//	获取未读消息数
	var unreadNums int64

	rds := cache.RedisPool.Get()
	defer rds.Close()
	tokenKey := models.TokenParamsPrefix(models.AdminPrefixTokenKey, adminInfo.Id)
	body.SuccessJSON(w, &loginData{
		Menu:       models.NewAdminMenu(nil).GetAdminMenuList(adminInfo.Id),
		RouterList: utils.GetMapValues(models.NewAdminAuthChild(nil).GetRolesRouteList(models.NewAdminAuthAssignment(nil).GetAdminRoleList(adminInfo.Id))),
		UserInfo: &userInfo{
			Id:         adminInfo.Id,
			Username:   adminInfo.UserName,
			Nickname:   adminInfo.Nickname,
			Email:      adminInfo.Email,
			Avatar:     adminInfo.Avatar,
			Money:      adminInfo.Money,
			Data:       adminInfo.Data,
			InviteCode: models.NewUserInvite(nil).GetInviteCode(adminInfo.Id, 0),
			OnlineNums: onlineNums,
			UnreadNums: unreadNums,
			ExpiredAt:  adminInfo.ExpiredAt,
			UpdatedAt:  adminInfo.UpdatedAt,
		},
		Token:    router.TokenManager.Generate(rds, tokenKey, adminInfo.Id, 0),
		TokenKey: tokenKey,
	})
}
