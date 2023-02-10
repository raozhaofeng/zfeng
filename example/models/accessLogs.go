package models

import (
	"database/sql"
	"fmt"
	"github.com/raozhaofeng/zfeng/database"
	"github.com/raozhaofeng/zfeng/database/define"
	"github.com/raozhaofeng/zfeng/logs"
	"github.com/raozhaofeng/zfeng/router"
	"github.com/raozhaofeng/zfeng/utils"
	"github.com/raozhaofeng/zfeng/utils/body"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

const (
	AccessLogsTypeAdmin   = 1       //	后端日志
	AccessLogsTypeHome    = 2       //	前端日志
	AccessLogsModuleAdmin = "admin" //	后端
	AccessLogsModuleHome  = "home"  //	后端
)

type AccessLogsAttrs struct {
	Id        int64  `json:"id"`         //主键
	AdminId   int64  `json:"admin_id"`   //管理员ID
	UserId    int64  `json:"user_id"`    //用户ID
	Type      int64  `json:"type"`       //日志类型
	Name      string `json:"name"`       //标题
	Ip4       string `json:"ip4"`        //IP4地址
	UserAgent string `json:"user_agent"` //ua信息
	Lang      string `json:"lang"`       //语言信息
	Route     string `json:"route"`      //操作路由
	Data      string `json:"data"`       //数据
	CreatedAt int64  `json:"created_at"` //时间
}

type AccessLogs struct {
	define.Db
}

func NewAccessLogs(tx *sql.Tx) *AccessLogs {
	return &AccessLogs{
		database.DbPool.NewDb(tx).Table("access_logs"),
	}
}

// FindOne 查询单挑
func (c *AccessLogs) FindOne() *AccessLogsAttrs {
	attrs := new(AccessLogsAttrs)
	c.QueryRow(func(row *sql.Row) {
		err := row.Scan(&attrs.Id, &attrs.AdminId, &attrs.UserId, &attrs.Type, &attrs.Name, &attrs.Ip4, &attrs.UserAgent, &attrs.Lang, &attrs.Route, &attrs.Data, &attrs.CreatedAt)
		if err != nil {
			attrs = nil
		}
	})
	return attrs
}

// FindMany 查询多条
func (c *AccessLogs) FindMany() []*AccessLogsAttrs {
	data := make([]*AccessLogsAttrs, 0)
	c.Query(func(rows *sql.Rows) {
		tmp := new(AccessLogsAttrs)
		_ = rows.Scan(&tmp.Id, &tmp.AdminId, &tmp.UserId, &tmp.Type, &tmp.Name, &tmp.Ip4, &tmp.UserAgent, &tmp.Lang, &tmp.Route, &tmp.Data, &tmp.CreatedAt)
		data = append(data, tmp)
	})
	return data
}

// UniqueVisitor 网站独立访客
func (c *AccessLogs) UniqueVisitor(adminIds []string, betweenTime []int64) int64 {
	c.Field("count(DISTINCT ip4)").Where("AND", "admin_id in ("+strings.Join(adminIds, ",")+")", nil)
	if len(betweenTime) == 2 {
		c.AndWhere("created_at between ? and ?", []any{betweenTime[0], betweenTime[1]})
	}
	return c.Count()
}

// RouteAccessFunc 路由日志方法
func RouteAccessFunc(accessType int64, handleParams *router.Handle, r *http.Request, claims *router.Claims) {
	fmt.Println(" -> ", handleParams.Name, handleParams.Route, handleParams.Method, time.Now().Format("2006-01-02 15:04:05"))
	//	验证的路由， 没有验证的路由， 后端跟前端
	var adminId, userId int64
	if claims != nil {
		adminId = claims.AdminId
		userId = claims.UserId
	}

	//	所有没有验证的方法,  没有管理ID,
	if adminId == 0 {
		adminId = NewAdminUser(nil).GetDomainAdminId(r)
	}

	data := `{"GET": ` + r.URL.Query().Encode() + `, "POST": ` + body.GetBody(r) + `}`
	if strings.Contains(handleParams.Route, "login") {
		data = ""
	}
	nowTime := time.Now().Unix()

	logs.Logger.Info(" -> ", zap.String("name", handleParams.Name), zap.String("method", handleParams.Method), zap.String("router", handleParams.Route))
	_, _ = NewAccessLogs(nil).
		Field("admin_id", "user_id", "type", "name", "ip4", "user_agent", "lang", "route", "data", "created_at").
		Value("?", "?", "?", "?", "INET_ATON(?)", "?", "?", "?", "?", "?").
		Args(adminId, userId, accessType, handleParams.Name, utils.GetUserRealIP(r), r.Header.Get("User-Agent"), r.Header.Get("Accept-Language"), handleParams.Route, data, nowTime).
		Insert()
}
