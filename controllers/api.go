package controllers

import (
	"common"
	"conf"

	"github.com/kataras/iris"
)

const (
	kLinkTypeAd       = "ad"
	kLinkTypeDownload = "dl"
)

type AjaxRsp struct {
	Success int32  `json:"success"`
	Msg     string `json:"msg"`
	DMsg    string `json:"d_msg"`
	DTxt    string `json:"d_txt"`
	DLink   string `json:"d_link"`
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

func AjaxGetHandler(ctx *iris.Context) {
	ajaxType := ctx.Param("type")
	postId := common.Atoi32(ctx.Param("id"))
	if postId == 0 {
		ctx.JSON(iris.StatusOK, AjaxRsp{
			Msg: "非法请求",
		})
		return
	}

	ctx.SetHeader("Access-Control-Allow-Origin", "*")
	switch ajaxType {
	case kLinkTypeDownload:
		adLinks := conf.GetArticleAdLinks(postId)
		sTitle := conf.GetArticleStitle(postId)
		if len(adLinks) == 0 {
			ctx.JSON(iris.StatusOK, AjaxRsp{
				Msg: "资源还未上传，请稍等片刻",
			})
			return
		}
		ckUrl := GetCkLeastUrl(ctx.RemoteAddr(), adLinks)
		ctx.JSON(iris.StatusOK, AjaxRsp{
			Success: 1,
			Msg:     "vip资源免费下载",
			DMsg:    sTitle + ".zip",
			DTxt:    "非会员有广告页面，点击跳过广告等按钮才会进入网盘下载",
			DLink:   ckUrl,
		})
		return
	default:
		ctx.JSON(iris.StatusOK, AjaxRsp{
			Msg: "非法请求",
		})
		return
	}
}
