package router

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gomodule/redigo/redis"
	"github.com/raozhaofeng/zfeng/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type contextKey int

const (
	ClaimsKey                 contextKey = 1                   //	claims类型
	TokenParamsRedisName                 = "_tokenParams"      //	Token参数缓存名称
	TokenValuesRedisName                 = "_tokenValues"      //	Token值缓存名称
	TokenAdminRolesRouterName            = "_adminRolesRouter" //	管理角色路由缓存名称
)

// TokenManager Token管理
var TokenManager *Token

// Claims Token对象
type Claims struct {
	AdminId            int64 //	管理ID
	UserId             int64 //	用户ID
	jwt.StandardClaims       //	jwt基础参数
}

// TokenParams Token 参数
type TokenParams struct {
	Key       string        `json:"key"`       //	密钥
	Only      bool          `json:"only"`      //	是否唯一
	Expire    time.Duration `json:"expire"`    //	过期时间
	Whitelist string        `json:"whitelist"` //	白名单
	Blacklist string        `json:"blacklist"` //	黑名单
}

type Token struct {
}

// InitializationToken 初始化Token
func InitializationToken(rds redis.Conn, tokenParamsList map[string]*TokenParams, adminRolesRouter map[int64][]string) {
	TokenManager = &Token{}
	//	初始化设置Token参数
	for tokenKey, tokenParams := range tokenParamsList {
		TokenManager.SetTokenParams(rds, tokenKey, tokenParams)
	}

	// 初始化设置管理路由列表
	for adminId, rolesRouter := range adminRolesRouter {
		TokenManager.SetTokenAdminRolesRouter(rds, adminId, rolesRouter)
	}
}

// Generate 生成Token
func (c *Token) Generate(rds redis.Conn, tokenKey string, adminId, userId int64) string {
	tokenParams := c.GetTokenParams(rds, tokenKey)
	if tokenParams == nil {
		return ""
	}

	nowTime := time.Now()
	claims := &Claims{
		UserId:  userId,
		AdminId: adminId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: nowTime.Add(tokenParams.Expire * time.Second).Unix(), //	过期时间
			IssuedAt:  nowTime.Unix(),                                       //	签发时间
		},
	}

	//	生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(tokenParams.Key))
	if err != nil {
		panic(err)
	}

	c.SetTokenValue(rds, adminId, userId, tokenStr)
	return tokenStr
}

// Verify 验证Token
func (c *Token) Verify(rds redis.Conn, r *http.Request) *Claims {
	claims := &Claims{}

	tokenStr, tokenKey := c.GetHeaderTokenAndTokenKey(r)
	if tokenStr == "" || tokenKey == "" {
		return nil
	}

	// 验证JWT TokenStr
	tokenParams := c.GetTokenParams(rds, tokenKey)
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		token.Method = jwt.SigningMethodHS256
		return []byte(tokenParams.Key), nil
	})

	if err != nil || token == nil || !token.Valid {
		return nil
	}

	//	判断是否唯一Token
	if tokenParams.Only && c.GetTokenValue(rds, claims.AdminId, claims.UserId) != tokenStr {
		return nil
	}

	//	判断是否白名单
	userRealIP := utils.GetUserRealIP(r)
	if tokenParams.Whitelist != "" && utils.SliceStringIndexOf(userRealIP, strings.Split(tokenParams.Whitelist, ",")) == -1 {
		return nil
	}

	//	判断是否黑名单
	if tokenParams.Blacklist != "" && utils.SliceStringIndexOf(userRealIP, strings.Split(tokenParams.Blacklist, ",")) > -1 {
		return nil
	}

	return claims
}

// SetTokenParams 设置Token参数
func (c *Token) SetTokenParams(rds redis.Conn, tokenKey string, tokenParams *TokenParams) {
	tokenParamsBytes, err := json.Marshal(tokenParams)
	if err != nil {
		panic(err)
	}
	_, err = rds.Do("HSET", TokenParamsRedisName, tokenKey, tokenParamsBytes)
	if err != nil {
		panic(err)
	}
}

// GetTokenParams 获取Token参数
func (c *Token) GetTokenParams(rds redis.Conn, tokenKey string) *TokenParams {
	tokenParamsBytes, err := redis.Bytes(rds.Do("HGET", TokenParamsRedisName, tokenKey))
	if err != nil {
		return nil
	}
	tokenParams := new(TokenParams)
	err = json.Unmarshal(tokenParamsBytes, &tokenParams)
	if err != nil {
		panic(err)
	}
	return tokenParams
}

// GetTokenValue 获取Token值
func (c *Token) GetTokenValue(rds redis.Conn, adminId, userId int64) string {
	tokenStr, _ := redis.String(rds.Do("HGET", TokenValuesRedisName, c.GetTokenValueKey(adminId, userId)))
	return tokenStr
}

// SetTokenValue 设置Token值
func (c *Token) SetTokenValue(rds redis.Conn, adminId, userId int64, tokenStr string) {
	_, err := rds.Do("HSET", TokenValuesRedisName, c.GetTokenValueKey(adminId, userId), tokenStr)
	if err != nil {
		panic(err)
	}
}

// DelTokenValue 删除Token值
func (c *Token) DelTokenValue(rds redis.Conn, adminId, userId int64) {
	_, err := rds.Do("HDEL", TokenValuesRedisName, c.GetTokenValueKey(adminId, userId))
	if err != nil {
		panic(err)
	}
}

// GetTokenAdminRolesRouter 获取管理角色路由列表
func (c *Token) GetTokenAdminRolesRouter(rds redis.Conn, adminId int64) []string {
	adminRolesRouter, err := redis.String(rds.Do("HGET", TokenAdminRolesRouterName, adminId))
	if err != nil {
		return []string{}
	}
	return strings.Split(adminRolesRouter, ",")
}

// SetTokenAdminRolesRouter 设置管理角色路由列表
func (c *Token) SetTokenAdminRolesRouter(rds redis.Conn, adminId int64, rolesRouter []string) {
	_, err := rds.Do("HSET", TokenAdminRolesRouterName, adminId, strings.Join(rolesRouter, ","))
	if err != nil {
		panic(err)
	}
}

// AuthRouter 验证路由
func (c *Token) AuthRouter(rds redis.Conn, adminId int64, router string) bool {
	adminRolesRouter := c.GetTokenAdminRolesRouter(rds, adminId)
	for _, adminRouter := range adminRolesRouter {
		if router == adminRouter || adminRouter == "*" {
			return true
		}
	}
	return false
}

// GetTokenValueKey 获取Token值key
func (c *Token) GetTokenValueKey(adminId, userId int64) string {
	adminIdStr := strconv.FormatInt(adminId, 10)
	userIdStr := strconv.FormatInt(userId, 10)
	return adminIdStr + "_" + userIdStr
}

// GetContextClaims 获取当前Claims
func (c *Token) GetContextClaims(r *http.Request) *Claims {
	return r.Context().Value(ClaimsKey).(*Claims)
}

// GetHeaderTokenAndTokenKey 获取头信息Token参数
func (c *Token) GetHeaderTokenAndTokenKey(r *http.Request) (string, string) {
	tokenStr := r.Header.Get("Token")
	tokenKey := r.Header.Get("Token-Key")

	urlTokenStr := r.URL.Query().Get("token")
	urlTokenKey := r.URL.Query().Get("key")
	if urlTokenStr != "" {
		tokenStr = urlTokenStr
	}
	if urlTokenKey != "" {
		tokenKey = urlTokenKey
	}
	return tokenStr, tokenKey
}
