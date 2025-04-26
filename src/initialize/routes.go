package initialize

import (
	"app/config"
	"app/middlewares"
	"app/routes"
	"embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itchyny/timefmt-go"
)

func InitRoutes(e *gin.Engine, embededFiles embed.FS) {
	// 静态文件
	if config.DEBUG {
		e.Static("/statics", "./assets/statics")
		e.LoadHTMLFiles("./assets/index.html")
	} else {
		e.GET("/statics/*filepath", func(ctx *gin.Context) {
			ctx.FileFromFS("assets/statics/" + ctx.Param("filepath"), http.FS(embededFiles))
		})

		rootTemplate := template.New("").Funcs(e.FuncMap)
		templ := template.Must(rootTemplate, func () error {
			_, err := rootTemplate.ParseFS(embededFiles, "assets/index.html")
			if err != nil {
				return err
			}
			return nil
		}())
		e.SetHTMLTemplate(templ)
	}


	e.Use(gin.LoggerWithFormatter(func (param gin.LogFormatterParams) string {
		headers := "{"
		for k, v := range param.Request.Header {
			line := ""
			for _, item := range v {
				line = line + item + ";"
			}
			headers = headers + k + ":" + line + " "
		}
		headers = headers + "}"

		return fmt.Sprintf("%s - [%s] %s %s %s %d %s \"%s\" \"%s\"\n",
				param.ClientIP,
				timefmt.Format(param.TimeStamp, "%Y-%m-%d %H:%M:%S %z"),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				headers,
				param.ErrorMessage,
		)
	}))
	e.Use(middlewares.ErrorHandleMiddlewareFactory())

	v1 := e.Group(config.Config.Server.BaseUrl + "/api/v1", middlewares.BasicAuthMiddleware, middlewares.RawBodyMiddleware)
	routes.Init(v1, e, embededFiles)
}