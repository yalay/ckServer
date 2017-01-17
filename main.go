package main

import (
	"conf"
	"net/url"
	"path"
	"strings"

	"github.com/iris-contrib/middleware/logger"
	"github.com/kataras/iris"
)

const (
	VERSION       = "v1"
	TEMPLATE_PATH = "./templates/"
)

func main() {
	iris.Config.IsDevelopment = true // this will reload the templates on each request
	iris.StaticWeb("/img", TEMPLATE_PATH+VERSION+"/img")
	iris.StaticWeb("/css", TEMPLATE_PATH+VERSION+"/css")
	iris.StaticWeb("/js", TEMPLATE_PATH+VERSION+"/js")

	iris.Use(logger.New())
	iris.Get("/", hi)
	iris.Get("/article-169.html", hi)
	iris.Get("/index.php", hi)
	iris.Get("/ck", ckHandler)
	iris.Listen(":8080")
}

func hi(ctx *iris.Context) {
	ctx.Log("%s", ctx.Request.Referer())
	ctx.MustRender(VERSION+"/center.html", struct{ Name string }{Name: "iris"})
}

// article-169.html
// index.php?c=Article&id=169
func ckHandler(ctx *iris.Context) {
	rawRefer := ctx.Request.Referer()
	if rawRefer == "" {
		ctx.NotFound()
		return
	}

	var articleKey string
	// 伪静态
	if strings.HasSuffix(rawRefer, ".html") {
		fileName := path.Base(rawRefer)
		if !strings.HasPrefix(fileName, "article") {
			ctx.NotFound()
			return
		}
		articleKey = strings.TrimSuffix(fileName, ".html")
	} else {
		referUrl, err := url.Parse(rawRefer)
		if err != nil {
			ctx.NotFound()
			return
		}
		referValues := referUrl.Query()
		className := strings.ToLower(referValues.Get("c"))
		if className != "article" {
			ctx.NotFound()
			return
		}
		articleId := referValues.Get("id")
		articleKey = "article-" + articleId
	}

	urls := conf.GetUrlsByArticleKey(articleKey)
	if len(urls) == 0 {
		ctx.NotFound()
		return
	}

	ctx.Redirect(conf.GetCkLeastUrl(ctx.RemoteAddr(), urls[0]))
}
