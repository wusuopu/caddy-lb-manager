package routes

import (
	"app/config"
	"app/controllers/dashboard"
	"app/middlewares"
	"embed"
	"time"

	"github.com/gin-gonic/gin"
)

func homePage (ctx *gin.Context) {
	debug := gin.Mode() != gin.ReleaseMode
	ctx.HTML(200, "index.html", gin.H{
		"debug": debug,
		"version": config.Config.Server.Version,
		"base_url": config.Config.Server.BaseUrl,
		"t": time.Now().Unix(),
	})
}

func Init(router *gin.RouterGroup, engine *gin.Engine, embededFiles embed.FS) {
	engine.GET(config.Config.Server.BaseUrl + "/", middlewares.BasicAuthMiddleware, homePage)

	engine.GET("_health", func(ctx *gin.Context) {
		ctx.String(200, "ok")
	})

	router.GET("/dashboard", dashboard.Index)

	serverGroup := router.Group("/servers")
	InitServer(serverGroup)
	InitRoute(serverGroup)
	InitUpstream(router.Group("/upstreams"))
	InitCaddyfile(router.Group("/caddy"))
	InitAuthentication(router.Group("/authentications"))
}