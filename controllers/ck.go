package controllers

import (
	"common"
	"conf"
	"fmt"

	"github.com/kataras/iris"
)

const (
	KPathCk = "/ck"
)

// /ck?aid=123&idx=1
// idx从1开始计数，0是非法
func CkHandler(ctx *iris.Context) {
	articleId := ctx.URLParam("aid")
	if articleId == "" {
		ctx.NotFound()
		return
	}

	articleKey := genArticleKey(articleId)
	urls := conf.GetUrlsByArticleKey(articleKey)
	if len(urls) == 0 {
		ctx.NotFound()
		return
	}

	pkgIdx := common.Atoi(ctx.URLParam("idx"))
	if pkgIdx == 0 || pkgIdx > len(urls) {
		ctx.NotFound()
		return
	}

	ctx.Redirect(GetCkLeastUrl(ctx.RemoteAddr(), urls[pkgIdx-1]))
}

// 动态点击链接
func GenDynamicCkUrl(articleId string, index int) string {
	return fmt.Sprintf("%s?aid=%s&idx=%d", KPathCk, articleId, index)
}
