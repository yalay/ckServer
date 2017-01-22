package controllers

import (
	"common"
	//"conf"
	"fmt"

	"github.com/kataras/iris"
)

const (
	KPathCk = "/ck"
)

// /ck?aid=123&idx=1
// idx从1开始计数，0是非法
func CkHandler(ctx *iris.Context) {
	articleId := common.Atoi(ctx.URLParam("aid"))
	if articleId == 0 {
		ctx.NotFound()
		return
	}

	urls := GetArticleAdUrls(articleId)
	//urls := conf.GetUrlsByArticleKey(articleKey)
	if len(urls) == 0 {
		ctx.NotFound()
		return
	}

	articleKey := genArticleKey(articleId)
	IncCkCount(articleKey)
	pkgIdx := common.Atoi(ctx.URLParam("idx"))
	if pkgIdx == 0 {
		ctx.NotFound()
		return
	}

	ctx.Redirect(GetCkLeastUrl(ctx.RemoteAddr(), urls[pkgIdx-1]))
}

// 动态点击链接
func GenDynamicCkUrl(articleId int, index int) string {
	return fmt.Sprintf("%s?aid=%d&idx=%d", KPathCk, articleId, index)
}
