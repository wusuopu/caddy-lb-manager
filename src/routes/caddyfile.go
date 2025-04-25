package routes

import (
	"app/controllers/caddyfile"

	"github.com/gin-gonic/gin"
)

func InitCaddyfile(r *gin.RouterGroup) {
	r.GET("/", caddyfile.Index)
	r.POST("/reload", caddyfile.Reload)
}
