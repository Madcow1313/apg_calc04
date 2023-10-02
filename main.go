package main

import (
	"fmt"
	"net/http"
	"strings"
	"webCalc/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	var Controller controller.Controller
	Controller.Init()
	router := gin.New()

	router.LoadHTMLFiles("index.html", "help.html", "graph_window.html")
	router.StaticFile("./jsScript/index.js", "jsScript/index.js")
	router.GET("/index.html", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})
	router.POST("/index.html", func(ctx *gin.Context) {
		q, _ := ctx.GetQuery("body")
		q = strings.Trim(q, "'")
		if strings.HasPrefix(q, "x= ") {
			q = strings.TrimLeft(q, "x= ")
			Controller.Expression = strings.ReplaceAll(Controller.Expression, "X", q)
		} else if q == "=" {
			Controller.InvokeModel()
			fmt.Println(Controller.LastResult)
		} else {
			Controller.HandleMessage(q)
		}
	})
	router.Run()
}
