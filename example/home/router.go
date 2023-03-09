package home

import (
	"github.com/raozhaofeng/zfeng/example/home/controllers/index"
	"github.com/raozhaofeng/zfeng/router"
)

func Router() []*router.Handle {
	return []*router.Handle{
		router.NewTokenHandle("Token", "POST", "/token", index.Index),
	}
}
