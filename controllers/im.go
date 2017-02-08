package controllers

import (
	"common"
	"conf"

	"github.com/kataras/iris"
)

const (
	KPathIm = "/im"
)

type imPageParams struct {
	Id           int32
	Im           int
	Ck           int
	Title        string
	STitle       string
	Desc         string
	Cover        string
	DownloadUrls []string
}

var emptyParams = imPageParams{
	Title: "请从文章页点击下载",
	Desc:  "没找到作品下载内容，请从文章页点击下载",
	Cover: "http://blog.vvniu.com/img/niu.jpg",
}

// im/article/169
func ImHandler(ctx *iris.Context) {
	referUrl := ctx.Request.Referer()
	if !conf.IsInWhiteList(referUrl) {
		ctx.MustRender("im.html", emptyParams)
		return
	}

	channel := ctx.Param("type")
	if channel != "article" {
		ctx.MustRender("im.html", emptyParams)
		return
	}

	articleId := common.Atoi32(ctx.Param("id"))
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

	title, sTitle, desc, cover := GetArticleAttrs(articleId)
	params := imPageParams{
		Id:           articleId,
		Im:           GetImCount(articleKey),
		Ck:           GetCkCount(articleKey),
		Title:        title,
		STitle:       sTitle,
		Desc:         desc,
		Cover:        cover,
		DownloadUrls: downloadUrls,
	}
	ctx.MustRender("im.html", params)
}
