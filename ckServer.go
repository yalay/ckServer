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
	iris.StaticWeb("img", common.TEMPLATE_PATH+"/img")
	iris.StaticWeb("css", common.TEMPLATE_PATH+"/css")
	iris.StaticWeb("js", common.TEMPLATE_PATH+"/js")
	iris.StaticWeb("fonts", common.TEMPLATE_PATH+"/fonts")

	iris.Use(logger.New())
	iris.UseTemplate(html.New(html.Config{Layout: iris.NoLayout})).Directory(common.TEMPLATE_PATH, ".html")

	iris.Get("/ck/:info", controllers.CkHandler)
	iris.Get("/im/:id", controllers.ImHandler)
	iris.Get("/encodes/:id/:type/:index", controllers.EncodesGetHandler)
	iris.Listen(":" + strconv.Itoa(listenPort))
}
