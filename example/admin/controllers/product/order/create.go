package order

import (
	"basic/models"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/database"
	"github.com/raozhaofeng/zfeng/router"
	"github.com/raozhaofeng/zfeng/utils"
	"github.com/raozhaofeng/zfeng/utils/body"
	"github.com/raozhaofeng/zfeng/validator"
	"net/http"
	"strconv"
	"time"
)

type createParams struct {
	UserName  string `json:"username" validate:"required"`
	ProductId int64  `json:"product_id" validate:"required,gt=0"`
	Nums      int64  `json:"nums" validate:"required,gt=0"`
}

func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(createParams)
	_ = body.ReadJSON(r, params)

	//  验证参数
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	// 查询用户说会否存在
	userModel := models.NewUser(nil)
	userModel.AndWhere("username=?", params.UserName)
	userInfo := userModel.FindOne()
	if userInfo == nil {
		body.ErrorJSON(w, "用户不存在", -1)
		return
	}

	// 查询产品Id是否存在
	productModel := models.NewProduct(nil)
	productModel.AndWhere("id=?", params.ProductId).AndWhere("status>?", models.ProductStatusDelete)
	productInfo := productModel.FindOne()
	if productInfo == nil {
		body.ErrorJSON(w, "产品不存在", -1)
		return
	}

	// 如果是超级管理员能修改所有用户， 不是超级管理员只能修改自己用户
	adminId := router.TokenManager.GetContextClaims(r).AdminId
	nowTime := time.Now()

	if adminId != models.AdminUserSupermanId && adminId != userInfo.AdminId {
		body.ErrorJSON(w, "权限不足", -1)
		return
	}

	tx := database.DbPool.GetTx()
	defer tx.Rollback()

	// 计算过期时间， 收益时间
	productDataInfo := new(models.ProductDataAttrs)
	err = json.Unmarshal([]byte(productInfo.Data), &productDataInfo)
	if err != nil {
		panic(err)
	}
	expiredTime := time.Now().AddDate(0, 0, int(productDataInfo.Expire))
	afterHour, _ := time.ParseDuration(strconv.FormatInt(productDataInfo.Interval, 10) + "h")
	updatedTime := time.Now().Add(afterHour)

	_, err = models.NewProductOrder(tx).
		Field("admin_id", "user_id", "product_id", "order_sn", "money", "nums", "data", "expired_at", "updated_at", "created_at").
		Args(userInfo.AdminId, userInfo.Id, params.ProductId, utils.NewRandom().OrderSn(), productInfo.Money, params.Nums, productInfo.Data, expiredTime.Unix(), updatedTime.Unix(), nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}

	// 记录用户消费账单
	spendAmount := productInfo.Money * float64(params.Nums)
	err = models.UserFundingChanges(tx, userInfo.AdminId, userInfo.Id, userInfo.ParentId, nil, 0, models.UserBillTypeBuyProduct, userInfo.Money, spendAmount)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	_ = tx.Commit()
	body.SuccessJSON(w, "ok")
}
