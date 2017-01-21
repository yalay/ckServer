package main

import (
	"common"
	"controllers"
	"models"

	"github.com/iris-contrib/middleware/logger"
	"github.com/kataras/go-template/django"
	"github.com/kataras/iris"
)

/*
  article-169:
    -
      - http://test1.com
      - http://test1-1.com
    -
      - http://test2.com
      - http://test2-1.com
*/

func main() {
	iris.Config.IsDevelopment = true // this will reload the templates on each request
	iris.StaticWeb("/img", common.TEMPLATE_PATH+"/img")
	iris.StaticWeb("/css", common.TEMPLATE_PATH+"/css")
	iris.StaticWeb("/js", common.TEMPLATE_PATH+"/js")
	iris.StaticWeb("/fonts", common.TEMPLATE_PATH+"/fonts")

	iris.Use(logger.New())
	iris.UseTemplate(django.New()).Directory(common.TEMPLATE_PATH, ".html")

	models.UpdateArticle(169, "title-169", "desc-169")
	models.AddArticleAdUrl(169, 0, "http://ad1.com")
	models.AddArticleAdUrl(169, 0, "http://ad1-1.com")
	models.AddArticleAdUrl(169, 1, "http://ad2.com")
	models.AddArticleAdUrl(169, 1, "http://ad2-1.com")

	models.AddArticleDownloadUrl(169, 0, "http://localhost/1")
	models.AddArticleDownloadUrl(169, 1, "http://localhost/2")

	iris.Get("/", hi)
	iris.Get("/article-169.html", controllers.ImHandler)
	iris.Get("/ck", controllers.CkHandler)
	iris.Get("/im", controllers.ImHandler)
	iris.Listen(":8080")
}

func hi(ctx *iris.Context) {
	ctx.Log("%s", ctx.Request.Referer())
	ctx.MustRender("center.html", struct{ Name string }{Name: "iris"})
}
