package controllers

import (
	"encoding/base64"
	"fmt"
	"log"
	"models"
	"net/url"

	"github.com/golang/protobuf/proto"
	"github.com/kataras/iris"
)

const (
	KPathCk = "/ck"
)

// /ck/:info
// idx从1开始计数，0是非法
func CkHandler(ctx *iris.Context) {
	info := ctx.Param("info")
	if info == "" {
		ctx.NotFound()
		return
	}

	b64Data, err := url.QueryUnescape(info)
	if err != nil {
		ctx.NotFound()
		return
	}

	data, err := base64.StdEncoding.DecodeString(b64Data)
	if err != nil {
		ctx.NotFound()
		return
	}
	msg := &models.Msg{}
	proto.Unmarshal(data, msg)
	if msg.ArticleId == 0 || msg.PkgIndex == 0 {
		ctx.NotFound()
		return
	}

	var urls map[int32][]string
	if msg.NoAd {
		urls = GetArticleDownloadUrls(msg.ArticleId)
		articleKey := genArticleKey(msg.ArticleId)
		IncCkCount(articleKey)
	} else {
		// 通过广告链接跳转
		urls = GetArticleAdUrls(msg.ArticleId)
	}

	if len(urls) == 0 {
		ctx.NotFound()
		return
	}

	ctx.Redirect(GetCkLeastUrl(ctx.RemoteAddr(), urls[msg.PkgIndex]))
}

// 编码之后的点击广告跳转链接
func GenEncodedCkAdUrl(articleId int32, index int32) string {
	return encodedCkUrl(articleId, index, true)
}

// 编码之后的点击实际跳转链接
func GenEncodedCkDownloadUrl(articleId int32, index int32) string {
	return encodedCkUrl(articleId, index, false)
}

func encodedCkUrl(articleId int32, index int32, isAd bool) string {
	msg := &models.Msg{
		ArticleId: articleId,
		PkgIndex:  index,
		NoAd:      !isAd,
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Printf("info proto marshal error:%v", err)
		return ""
	}

	info := url.QueryEscape(base64.StdEncoding.EncodeToString(data))
	return fmt.Sprintf("%s/%s", KPathCk, info)
}
