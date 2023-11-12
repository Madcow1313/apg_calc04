package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"webCalc/controller"

	"github.com/gin-gonic/gin"
)

func GET(router *gin.Engine, Controller *controller.Controller) {
	router.GET("/", func(ctx *gin.Context) {
		Controller.Expression = Controller.LastResult
		fmt.Println(Controller.LastResult)
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"result": Controller.LastResult,
		})
	})
	router.GET("/help.html", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "help.html", nil)
	})
	router.GET("/graph_window.html", func(ctx *gin.Context) {
		_, err := Controller.InvokeGraphic(Controller.XMin, Controller.XMax, Controller.XMin, Controller.XMax)
		if err == nil {
			ctx.HTML(http.StatusOK, "graph_window.html", nil)
		}
	})
}

func POST(router *gin.Engine, Controller *controller.Controller) {
	router.POST("/", func(ctx *gin.Context) {
		q, _ := ctx.GetQuery("body")
		q = strings.Trim(q, "'")
		if q == " plus " {
			q = " + "
		} else if q == " divide " {
			q = " / "
		}
		if strings.HasPrefix(q, "x= ") {
			q = strings.TrimPrefix(q, "x= ")
			Controller.Expression = strings.ReplaceAll(Controller.Expression, "X", q)
		} else if strings.HasPrefix(q, "xy_min= ") {
			q = strings.TrimPrefix(q, "xy_min= ")
			Controller.XMin, _ = strconv.ParseFloat(q, 64)
		} else if strings.HasPrefix(q, "xy_max= ") {
			q = strings.TrimPrefix(q, "xy_max= ")
			Controller.XMax, _ = strconv.ParseFloat(q, 64)
		} else if q == "clear" {
			Controller.Expression = ""
		} else if q == "=" {
			Controller.InvokeModel()
		} else {
			Controller.HandleMessage(q)
		}
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"result": Controller.LastResult,
		})
	})
}

func main() {
	var Controller controller.Controller
	Controller.Init()
	router := gin.Default()

	router.LoadHTMLFiles("index.html", "help.html", "graph_window.html")
	router.StaticFile("./jsScript/index.js", "jsScript/index.js")
	router.StaticFile("./test.png", "test.png")

	GET(router, &Controller)
	POST(router, &Controller)
	router.Run()
}
