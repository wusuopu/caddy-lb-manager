package initialize

import (
	"app/config"
	"embed"

	"github.com/gin-gonic/gin"
)

func commonInit(e *gin.Engine, embededFiles embed.FS) *gin.Engine {
	InitServices()

	var engine *gin.Engine
	if e == nil {
		engine = gin.New()
	} else {
		engine = e
	}
	InitDB()
	InitLogger()
	InitRoutes(engine, embededFiles)
	return engine
}

func Init(e *gin.Engine, embededFiles embed.FS) *gin.Engine {
	// 先加载 .env 文件
	InitEnv()
	config.Load()

	if config.Config.Server.GO_ENV == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	return commonInit(e, embededFiles)
}

func InitTest(e *gin.Engine, embededFiles embed.FS) *gin.Engine {
	// 先加载 .env.test 文件
	InitEnv(".env.test")
	config.Load()
	config.Config.Server.GO_ENV = "test"

	gin.SetMode(gin.TestMode)
	return commonInit(e, embededFiles)
}
