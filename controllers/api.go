package controllers

import (
	"common"

	"github.com/kataras/iris"
)

const (
	kLinkTypeAd       = "ad"
	kLinkTypeDownload = "dl"
)

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
