package controllers

import (
	"common"

	"github.com/kataras/iris"
)

const (
	KPathIm = "/im"
)

var emptyParams = map[string]interface{}{
	"id":           0,
	"im":           0,
	"ck":           0,
	"title":        "未知标题",
	"desc":         "没找到对应下载内容，请重新从文章页点击下载",
	"cover":        "http://blog.vvniu.com/img/niu.jpg",
	"downloadUrls": nil,
}

// im?c=article&id=169
func ImHandler(ctx *iris.Context) {
	channel := ctx.URLParam("c")
	if channel != "article" {
		ctx.MustRender("im.html", emptyParams)
		return
	}

	articleId := common.Atoi32(ctx.URLParam("id"))
	if articleId == 0 {
		ctx.MustRender("im.html", emptyParams)
		return
	}

	pkgCount := GetArticlePkgCount(articleId)
	if pkgCount == 0 {
		ctx.MustRender("im.html", emptyParams)
		return
	}

	articleKey := genArticleKey(articleId)
	IncImCount(articleKey)
	downloadUrls := make([]string, pkgCount)
	for i, _ := range downloadUrls {
		downloadUrls[i] = GenEncodedCkAdUrl(articleId, int32(i+1))
	}

	title, desc, cover := GetArticleAttrs(articleId)
	params := map[string]interface{}{
		"id":           articleId,
		"im":           GetImCount(articleKey),
		"ck":           GetCkCount(articleKey),
		"title":        title,
		"desc":         desc,
		"cover":        cover,
		"downloadUrls": downloadUrls,
	}
	ctx.MustRender("im.html", params)
}
