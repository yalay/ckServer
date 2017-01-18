package main

import (
	"common"
	"controllers"

	"github.com/iris-contrib/middleware/logger"
	"github.com/kataras/go-template/django"
	"github.com/kataras/iris"
)

func main() {
	iris.Config.IsDevelopment = true // this will reload the templates on each request
	iris.StaticWeb("/img", common.TEMPLATE_PATH+"/img")
	iris.StaticWeb("/css", common.TEMPLATE_PATH+"/css")
	iris.StaticWeb("/js", common.TEMPLATE_PATH+"/js")
	iris.StaticWeb("/fonts", common.TEMPLATE_PATH+"/fonts")

	iris.Use(logger.New())
	iris.UseTemplate(django.New()).Directory(common.TEMPLATE_PATH, ".html")

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
