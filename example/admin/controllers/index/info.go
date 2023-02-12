package index

import (
	"basic/models"
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/router"
	"github.com/raozhaofeng/zfeng/utils/body"
	"net/http"
)

// Info 管理员信息
func Info(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	adminId := router.TokenManager.GetContextClaims(r).AdminId

	adminModel := models.NewAdminUser(nil)
	adminModel.AndWhere("id=?", adminId)
	adminInfo := adminModel.FindOne()

	//	获取在线人数
	var onlineNums int64

	//	获取未读消息数
	var unreadNums int64

	data := &userInfo{
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
	}

	body.SuccessJSON(w, data)
}
