package manager

import (
	"basic/models"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/zfeng/cache"
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
	Roles       map[string]bool     `json:"roles" validate:"required"`
	UserName    string              `json:"username" validate:"required"`
	Domain      string              `json:"domain" validate:"required"`
	Password    string              `json:"password" validate:"required"`
	SecurityKey string              `json:"security_key" validate:"required"`
	ExpiredAt   string              `json:"expired_at" validate:"required"`
	Data        *router.TokenParams `json:"data" validate:"required"`
}

// Create 新增管理员
func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params := new(createParams)
	_ = body.ReadJSON(r, params)

	//  验证参数
	err := validator.Instantiate.Struct(params)
	if err != nil {
		body.ErrorJSON(w, err.Error(), -1)
		return
	}

	adminId := router.TokenManager.GetContextClaims(r).AdminId
	settingAdminId := models.NewAdminUser(nil).GetSettingAdminId(adminId)
	location, _ := time.LoadLocation(models.NewAdminSetting(nil).GetAdminFieldString(settingAdminId, "site_timezone"))
	// 开始插入数据
	tx := database.DbPool.GetTx()
	defer tx.Rollback()
	nowTime := time.Now()
	expiredAt := nowTime.Unix()
	if params.ExpiredAt != "" {
		expireTime, err := time.ParseInLocation("2006/01/02", params.ExpiredAt, location)
		if err != nil {
			panic(err)
		}
		expiredAt = expireTime.Unix()
	}
	// 判断管理员是否存在
	oldAdminModel := models.NewAdminUser(nil)
	oldAdminModel.AndWhere("username=?", params.UserName)
	oldAdminInfo := oldAdminModel.FindOne()
	if oldAdminInfo != nil {
		body.ErrorJSON(w, "用户名已存在", -1)
		return
	}

	//	判断域名是否存在
	domainModel := models.NewAdminUser(nil)
	domainModel.AndWhere("domain like ?", "%"+params.Domain+"%")
	domainInfo := domainModel.FindOne()
	if domainInfo != nil {
		body.ErrorJSON(w, "域名已存在", -1)
		return
	}

	rds := cache.RedisPool.Get()
	defer rds.Close()

	if params.Data.Key == "" {
		params.Data.Key = "8888"
	}
	if params.Data.Expire == 0 {
		params.Data.Expire = 3600
	}
	adminDataBytes, _ := json.Marshal(params.Data)

	//  模型插入数据
	newAdminId, err := models.NewAdminUser(tx).
		Field("parent_id", "username", "password", "security_key", "domain", "data", "expired_at", "created_at").
		Args(adminId, params.UserName, utils.PasswordEncrypt(params.Password), utils.PasswordEncrypt(params.SecurityKey), params.Domain, string(adminDataBytes), expiredAt, nowTime.Unix()).
		Insert()
	if err != nil {
		panic(err)
	}

	//	如果是代理初始化一些数据
	if adminId == models.AdminUserSupermanId {
		// 插入管理员语言
		langModel := models.NewLang(nil)
		langModel.AndWhere("admin_id=?", adminId)
		initLangList := langModel.FindMany()
		for _, lang := range initLangList {
			_, err = models.NewLang(tx).
				Field("admin_id", "name", "alias", "icon", "sort", "status", "data", "created_at").
				Args(newAdminId, lang.Name, lang.Alias, lang.Icon, lang.Sort, lang.Status, lang.Data, nowTime.Unix()).Insert()
			if err != nil {
				panic(err)
			}
		}

		// 插入管理员语言包
		dictionaryModel := models.NewLangDictionary(nil)
		dictionaryModel.AndWhere("admin_id=?", adminId)
		initDictionaryList := dictionaryModel.FindMany()
		for _, dictionary := range initDictionaryList {
			_, err = models.NewLangDictionary(nil).
				Field("admin_id", "type", "alias", "name", "field", "value", "data", "created_at").
				Args(newAdminId, dictionary.Type, dictionary.Alias, dictionary.Name, dictionary.Field, dictionary.Value, dictionary.Data, dictionary.CreatedAt).
				Insert()
			if err != nil {
				panic(err)
			}

			//	缓存本地语言
			//	TODO...
		}

		// 插入管理员国家
		countryModel := models.NewCountry(nil)
		countryModel.AndWhere("admin_id=?", adminId)
		initCountryList := countryModel.FindMany()
		for _, country := range initCountryList {
			_, err = models.NewCountry(tx).
				Field("admin_id", "lang_id", "name", "alias", "iso1", "icon", "code", "sort", "status", "data", "created_at").
				Args(newAdminId, country.LangId, country.Name, country.Alias, country.Iso1, country.Icon, country.Code, country.Sort, country.Status, country.Data, nowTime.Unix()).
				Insert()
			if err != nil {
				panic(err)
			}
		}

		// 插入用户等级
		levelModel := models.NewUserLevel(nil)
		levelModel.AndWhere("admin_id=?", adminId)
		initLevelList := levelModel.FindMany()
		for _, level := range initLevelList {
			_, err = models.NewUserLevel(tx).
				Field("admin_id", "name", "icon", "level", "money", "days", "status", "data", "created_at", "updated_at").
				Args(newAdminId, level.Name, level.Icon, level.Level, level.Money, level.Days, level.Status, level.Data, level.CreatedAt, level.UpdatedAt).
				Insert()
			if err != nil {
				panic(err)
			}
		}

		// 插入支付配置
		paymentModel := models.NewWalletPayment(nil)
		paymentModel.AndWhere("admin_id=?", adminId)
		initPaymentList := paymentModel.FindMany()
		for _, payment := range initPaymentList {
			_, err = models.NewWalletPayment(tx).
				Field("admin_id", "icon", "mode", "type", "name", "account_name", "account_code", "sort", "status", "description", "data", "expand", "created_at", "updated_at").
				Args(newAdminId, payment.Icon, payment.Mode, payment.Type, payment.Name, payment.AccountName, payment.AccountCode, payment.Sort, payment.Status, payment.Description, payment.Data, payment.Expand, payment.CreatedAt, payment.UpdatedAt).
				Insert()
			if err != nil {
				panic(err)
			}
		}

		//	插入资产管理
		assetsModel := models.NewAssets(nil)
		assetsModel.AndWhere("admin_id=?", adminId)
		initAssetsList := assetsModel.FindMany()
		for _, assets := range initAssetsList {
			_, err = models.NewAssets(tx).
				Field("admin_id", "name", "icon", "type", "data", "status", "created_at", "updated_at").
				Args(newAdminId, assets.Name, assets.Icon, assets.Type, assets.Data, assets.Status, assets.CreatedAt, assets.UpdatedAt).
				Insert()
			if err != nil {
				panic(err)
			}
		}

		// 插入管理员配置
		adminSettingModel := models.NewAdminSetting(nil)
		adminSettingModel.AndWhere("admin_id=?", adminId)
		initAdminSettingList := adminSettingModel.FindMany()
		for _, setting := range initAdminSettingList {
			_, err = models.NewAdminSetting(tx).
				Field("admin_id", "group_id", "name", "type", "field", "value", "data").
				Args(newAdminId, setting.GroupId, setting.Name, setting.Type, setting.Field, setting.Value, setting.Data).
				Insert()

			// 新增前台Token参数
			if setting.Type == models.SettingTypeJson && setting.Field == models.UpdateAdminTokenParamsField {
				tokenParams := new(router.TokenParams)
				_ = json.Unmarshal([]byte(setting.Value), &tokenParams)
				//	需要设置ID
				newSettingAdminId := models.NewAdminUser(tx).GetSettingAdminId(newAdminId)
				router.TokenManager.SetTokenParams(rds, models.TokenParamsPrefix(models.HomePrefixTokenKey, newSettingAdminId), tokenParams)
			}

			if err != nil {
				panic(err)
			}
		}

		// 继承管理产品分类｜产品列表
		_ = models.NewProductCategory(tx).InheritAdminProduct(adminId, newAdminId, 0)
	}

	// 	验证角色是否存在，并且添加角色
	adminRoles := models.NewAdminAuthAssignment(nil).GetAdminRoleList(adminId)
	adminRolesList := models.NewAdminAuthItem(nil).GetAdminRoleCheckedList(adminId, adminRoles)
	for role, checked := range params.Roles {
		if _, ok := adminRolesList[role]; ok && checked {
			_, err = models.NewAdminAuthAssignment(tx).Field("item_name", "user_id", "created_at").
				Args(role, newAdminId, nowTime.Unix()).
				Insert()
			if err != nil {
				panic(err)
			}
		}
	}

	// 	新增后台Token参数 并且 设置当前用户权限路由
	newAdminIdStr := strconv.FormatInt(newAdminId, 10)
	router.TokenManager.SetTokenParams(rds, utils.PasswordEncrypt(models.AdminPrefixTokenKey+"_"+newAdminIdStr), params.Data)
	newAdminRoles := models.NewAdminAuthAssignment(tx).GetAdminRoleList(newAdminId)
	router.TokenManager.SetTokenAdminRolesRouter(rds, newAdminId, utils.GetMapValues(models.NewAdminAuthChild(tx).GetRolesRouteList(newAdminRoles)))

	_ = tx.Commit()
	body.SuccessJSON(w, "ok")
}
