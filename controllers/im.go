package controllers

import (
	"net/url"
	"path"
	"strings"

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

// article-169.html
// index.php?c=Article&id=169
func ImHandler(ctx *iris.Context) {
	rawRefer := ctx.Request.Referer()
	if rawRefer == "" {
		ctx.MustRender("im.html", emptyParams)
		return
	}

	articleKey := getArticleKeyFromRefer(rawRefer)
	if articleKey == "" {
		ctx.MustRender("im.html", emptyParams)
		return
	}
	articleId := getIdFromArticleKey(articleKey)
	if articleId == 0 {
		ctx.MustRender("im.html", emptyParams)
		return
	}

	pkgCount := GetArticlePkgCount(articleId)
	if pkgCount == 0 {
		ctx.MustRender("im.html", emptyParams)
		return
	}

	IncImCount(articleKey)
	downloadUrls := make([]string, pkgCount)
	for i, _ := range downloadUrls {
		downloadUrls[i] = GenEncodedCkUrl(articleId, int32(i+1))
	}
	params := map[string]interface{}{
		"id":           articleId,
		"im":           GetImCount(articleKey),
		"ck":           GetCkCount(articleKey),
		"title":        "test-titledddd",
		"desc":         "test-desc",
		"cover":        "http://blog.vvniu.com/img/niu.jpg",
		"downloadUrls": downloadUrls,
	}
	ctx.MustRender("im.html", params)
}

func getArticleKeyFromRefer(rawRefer string) string {
	var articleKey = ""
	// 伪静态
	if strings.HasSuffix(rawRefer, ".html") {
		fileName := path.Base(rawRefer)
		if !strings.HasPrefix(fileName, kArticleKeyPrefix) {
			return ""
		}
		articleKey = strings.TrimSuffix(fileName, ".html")
	} else {
		referUrl, err := url.Parse(rawRefer)
		if err != nil {
			return ""
		}
		referValues := referUrl.Query()
		className := strings.ToLower(referValues.Get("c"))
		if className != "article" {
			return ""
		}
		articleKey = kArticleKeyPrefix + referValues.Get("id")
	}
	return articleKey
}
