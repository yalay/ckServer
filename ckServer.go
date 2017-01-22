package main

import (
	"common"
	"controllers"

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

	controllers.AddArticle(169, "title-169", "desc-169", "https://ss0.bdstatic.com/5aV1bjqh_Q23odCf/static/superman/img/logo/bd_logo1_31bdc765.png")
	controllers.AddArticleAdUrl(169, 0, "http://ad1.com")
	controllers.AddArticleAdUrl(169, 0, "http://ad1-1.com")
	controllers.AddArticleAdUrl(169, 1, "http://ad2.com")
	controllers.AddArticleAdUrl(169, 1, "http://ad2-1.com")
	controllers.DeleteAdUrl("http://ad1-1.com")

	controllers.AddArticleDownloadUrl(169, 0, "http://localhost/1")
	controllers.AddArticleDownloadUrl(169, 1, "http://localhost/2")
	controllers.DeleteDownloadUrl("http://localhost/2")

	iris.Get("/article-169.html", controllers.ImHandler)
	iris.Get("/ck", controllers.CkHandler)
	iris.Get("/im", controllers.ImHandler)
	iris.Get("/articles/:id", controllers.ArticleGetHandler)
	iris.Post("/articles/:id", controllers.ArticlePostHandler)
	iris.Get("/links/:id/:type", controllers.LinksGetHandler)
	iris.Post("/links/:id/:type", controllers.LinksPostHandler)
	iris.Get("/encodes/:id/:type/:index", controllers.EncodesGetHandler)
	iris.Listen(":8080")
}
