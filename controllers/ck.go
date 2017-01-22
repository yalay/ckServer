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

// /ck?aid=123&idx=1
// idx从1开始计数，0是非法
func CkHandler(ctx *iris.Context) {
	info := ctx.URLParam("info")
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

	urls := GetArticleAdUrls(msg.ArticleId)
	//urls := conf.GetUrlsByArticleKey(articleKey)
	if len(urls) == 0 {
		ctx.NotFound()
		return
	}

	articleKey := genArticleKey(msg.ArticleId)
	IncCkCount(articleKey)
	ctx.Redirect(GetCkLeastUrl(ctx.RemoteAddr(), urls[msg.PkgIndex-1]))
}

// 编码之后的点击链接
func GenEncodedCkUrl(articleId int32, index int32) string {
	msg := &models.Msg{
		ArticleId: articleId,
		PkgIndex:  index,
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Printf("info proto marshal error:%v", err)
		return ""
	}

	info := url.QueryEscape(base64.StdEncoding.EncodeToString(data))
	return fmt.Sprintf("%s?info=%s", KPathCk, info)
}
