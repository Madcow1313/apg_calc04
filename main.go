package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"webCalc/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	var Controller controller.Controller
	Controller.Init()
	router := gin.Default()

	router.LoadHTMLFiles("index.html", "help.html", "graph_window.html")
	router.StaticFile("./jsScript/index.js", "jsScript/index.js")
	router.StaticFile("./test.png", "test.png")
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"result": Controller.LastResult,
		})
	})
	router.GET("/help.html", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "help.html", nil)
	})
	router.GET("/graph_window.html", func(ctx *gin.Context) {
		fmt.Println("xmin xmax", Controller.XMin, Controller.XMax)
		_, err := Controller.InvokeGraphic(Controller.XMin, Controller.XMax, Controller.XMin, Controller.XMax)
		if err == nil {
			ctx.HTML(http.StatusOK, "graph_window.html", nil)
		}
	})
	router.POST("/", func(ctx *gin.Context) {
		q, _ := ctx.GetQuery("body")
		q = strings.Trim(q, "'")
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
			fmt.Println(Controller.LastResult)
		} else {
			Controller.HandleMessage(q)
		}
		fmt.Println(Controller.Expression)
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"result": Controller.LastResult,
		})
	})
	router.Run()
}
