package routes

import (
	"app/controllers/caddyfile"

	"github.com/gin-gonic/gin"
)

func InitCaddyfile(r *gin.RouterGroup) {
	r.GET("/config", caddyfile.Index)
	r.POST("/reload", caddyfile.Reload)
	r.GET("/certificates", caddyfile.ListCertificate)
}
