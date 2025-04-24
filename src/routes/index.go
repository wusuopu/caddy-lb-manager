package routes

import (
	"app/config"
	"app/middlewares"
	"time"

	"github.com/gin-gonic/gin"
)

func Init(router *gin.RouterGroup, engine *gin.Engine) {
	// 静态文件
	engine.Static("/statics", "./assets/statics")
	engine.LoadHTMLFiles("./assets/index.html")

	// engine.StaticFile("/", "./assets/index.html")
	engine.GET("/", middlewares.BasicAuthMiddleware, func(ctx *gin.Context) {
		debug := gin.Mode() != gin.ReleaseMode
		ctx.HTML(200, "index.html", gin.H{
			"debug": debug,
			"version": config.Config["VERSION"],
			"t": time.Now().Unix(),
		})
	})

	engine.GET("_health", func(ctx *gin.Context) {
		ctx.String(200, "ok")
	})
}