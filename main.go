package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	router.LoadHTMLFiles("index.html")
	router.StaticFile("./jsScript/index.js", "jsScript/index.js")
	router.GET("/index.html", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})
	router.POST("/index.html", func(ctx *gin.Context) {
		q, _ := ctx.GetQuery("body")
		fmt.Println("q is ", q)
	})
	router.Run()
}
