package home

import (
	"basic/home/controllers/index"
	"github.com/raozhaofeng/zfeng/router"
)

func Router() []*router.Handle {
	return []*router.Handle{
		router.NewHandle("上传图片", "POST", "/upload", index.Upload),
		router.NewHandle("测试", "GET", "/test", index.Test),
	}
}
