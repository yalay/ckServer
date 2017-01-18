package controllers

import (
	"conf"
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
	"title":        "empty-title",
	"desc":         "empty-desc",
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
	if articleId == "" {
		ctx.MustRender("im.html", emptyParams)
		return
	}

	urls := conf.GetUrlsByArticleKey(articleKey)
	pkgTotalNum := len(urls)
	downloadUrls := make([]string, pkgTotalNum)
	for i, _ := range downloadUrls {
		downloadUrls[i] = GenDynamicCkUrl(articleId, i+1)
	}
	params := map[string]interface{}{
		"id":           articleId,
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
		articleKey = referValues.Get("id")
	}
	return articleKey
}
