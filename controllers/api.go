package controllers

import (
	"common"

	"github.com/kataras/iris"
)

const (
	kLinkTypeAd       = "ad"
	kLinkTypeDownload = "dl"
)

func ArticleGetHandler(ctx *iris.Context) {
	articleId := common.Atoi32(ctx.Param("id"))
	if articleId == 0 {
		ctx.NotFound()
		return
	}
	title, desc, cover := sqliteDb.GetArticleAttrs(articleId)
	if title == "" {
		ctx.NotFound()
		return
	}
	ctx.Writef("title:%s desc:%s cover:%s", title, desc, cover)
}

func ArticlePostHandler(ctx *iris.Context) {
	articleId := common.Atoi32(ctx.Param("id"))
	if articleId == 0 {
		ctx.NotFound()
		return
	}
	title := ctx.FormValue("title")
	desc := ctx.FormValue("desc")
	cover := ctx.FormValue("cover")
	sqliteDb.AddArticle(articleId, title, desc, cover)
	ctx.Writef("add success")
}

func LinksGetHandler(ctx *iris.Context) {
	articleId := common.Atoi32(ctx.Param("id"))
	if articleId == 0 {
		ctx.NotFound()
		return
	}

	var urls map[int32][]string
	linkType := ctx.Param("type")
	switch linkType {
	case kLinkTypeAd:
		urls = GetArticleAdUrls(articleId)
	case kLinkTypeDownload:
		urls = GetArticleDownloadUrls(articleId)
	default:
		ctx.NotFound()
		return
	}

	if len(urls) == 0 {
		ctx.NotFound()
		return
	}
	ctx.JSON(iris.StatusOK, urls)
}

func LinksPostHandler(ctx *iris.Context) {
	articleId := common.Atoi32(ctx.Param("id"))
	if articleId == 0 {
		ctx.Writef("id is invalid")
		return
	}

	ckUrl := ctx.FormValue("url")
	if ckUrl == "" {
		ctx.Writef("url is empty")
		return
	}

	var err error
	pkgIndex := common.Atoi32(ctx.FormValue("index"))
	linkType := ctx.Param("type")
	switch linkType {
	case kLinkTypeAd:
		err = AddArticleAdUrl(articleId, pkgIndex, ckUrl)
	case kLinkTypeDownload:
		err = AddArticleDownloadUrl(articleId, pkgIndex, ckUrl)
	default:
		ctx.NotFound()
		return
	}

	if err != nil {
		ctx.Writef(err.Error())
		return
	}
	ctx.Writef("success")
}

func EncodesGetHandler(ctx *iris.Context) {
	articleId := common.Atoi32(ctx.Param("id"))
	if articleId == 0 {
		ctx.Writef("id is invalid")
		return
	}

	pkgIndex := common.Atoi32(ctx.Param("index"))
	linkType := ctx.Param("type")
	switch linkType {
	case kLinkTypeAd:
		ctx.WriteString(GenEncodedCkAdUrl(articleId, pkgIndex))
	case kLinkTypeDownload:
		ctx.WriteString(GenEncodedCkDownloadUrl(articleId, pkgIndex))
	default:
		ctx.NotFound()
		return
	}
	return
}
