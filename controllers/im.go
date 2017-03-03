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
	Title: "作品还未上传，暂时不能下载",
	Desc:  "已通知相关客服，请等待更新",
	Cover: "/img/logo.png",
}

// im/articles/169
func ImHandler(ctx *iris.Context) {
	channel := ctx.Param("type")
	if channel != "articles" {
		ctx.EmitError(iris.StatusNotFound)
		return
	}

	articleId := common.Atoi32(ctx.Param("id"))
	if articleId == 0 {
		ctx.EmitError(iris.StatusNotFound)
		return
	}

	adLinks := conf.GetArticleAdLinks(articleId)
	if len(adLinks) == 0 {
		ctx.MustRender("im.html", emptyParams)
		return
	}

	articleKey := genArticleKey(articleId)
	IncImCount(articleKey)
	downloadUrls := make([]string, 1)
	downloadUrls[0] = GenEncodedCkAdUrl(articleId, 1)

	title, sTitle, cover, desc := conf.GetArticleAttrs(articleId)
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
