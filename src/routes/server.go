package routes

import (
	"app/controllers/server"

	"github.com/gin-gonic/gin"
)

func InitServer(r *gin.RouterGroup) {
	r.GET("/", server.Index)
	r.POST("/", server.Create)
	r.GET("/:serverId", server.Show)
	r.PUT("/:serverId", server.Update)
	r.DELETE("/:serverId", server.Delete)
}
