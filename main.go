package main

import (
	"github.com/gin-gonic/gin"
	"go-jwt-study/router"
	"log"
)

func main() {
	gin.SetMode("debug")
	r := gin.Default()
	frontMode := false
	if frontMode {
		log.Println("已开启前后端整合模式！")
		r.LoadHTMLGlob("static/index.html")
		r.Static("/static", "static")
		r.GET("/", func(context *gin.Context) {
			context.HTML(200, "index.html", "")
		})
	}
	router.RegRouter(r)
	log.Println("路由注册完成！当前端口为:8000")
	r.Run(":8000")
}
