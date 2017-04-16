package main

import (
	"common"
	"controllers"
	"flag"
	"strconv"

	"github.com/iris-contrib/middleware/logger"
	"github.com/kataras/go-template/html"
	"github.com/kataras/iris"
)

var (
	listenPort int
)

func init() {
	flag.IntVar(&listenPort, "p", 1320, "p=1320")
	flag.Parse()
}

func main() {
	iris.StaticWeb("static/img", common.TEMPLATE_PATH+"/img")
	iris.StaticWeb("static/css", common.TEMPLATE_PATH+"/css")
	iris.StaticWeb("static/js", common.TEMPLATE_PATH+"/js")
	iris.StaticWeb("static/fonts", common.TEMPLATE_PATH+"/fonts")

	iris.Use(logger.New())
	iris.UseTemplate(html.New(html.Config{Layout: iris.NoLayout})).Directory(common.TEMPLATE_PATH, ".html")

	iris.Get("/ck/:info", controllers.CkHandler)
	iris.Get("/im/:type/:id", controllers.ImHandler)
	iris.Get("/encodes/:id/:type/:index", controllers.EncodesGetHandler)
	iris.Get("/ajax/:type/:id", controllers.AjaxGetHandler)
	iris.Listen(":" + strconv.Itoa(listenPort))
}
