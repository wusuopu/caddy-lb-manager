package routes

import (
	"app/config"
	"app/controllers/dashboard"
	"app/middlewares"
	"app/utils"
	"embed"
	"fmt"
	"path"
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
	// 静态文件
	fmt.Println("Is debug:", config.DEBUG)
	if config.DEBUG {
		engine.Static("/statics", "./assets/statics")
		engine.LoadHTMLFiles("./assets/index.html")
	} else {
		tmpDir, err := utils.ExpandEmbed(embededFiles)
		if err != nil {
			panic(err)
		}
		fmt.Printf("embed tmpDir: %s\n", tmpDir)
		engine.Static("/statics", path.Join(tmpDir, "assets/statics"))
		engine.LoadHTMLFiles(path.Join(tmpDir, "assets/index.html"))
	}

	// engine.StaticFile("/", "./assets/index.html")
	engine.GET(config.Config.Server.BaseUrl + "/", middlewares.BasicAuthMiddleware, homePage)

	engine.GET("_health", func(ctx *gin.Context) {
		ctx.String(200, "ok")
	})

	router.GET("/dashboard", dashboard.Index)

	serverGroup := router.Group("/servers")
	InitServer(serverGroup)
	InitRoute(serverGroup)
	InitUpstream(router.Group("/upstreams"))
}