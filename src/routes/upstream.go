package routes

import (
	"app/controllers/upstream"

	"github.com/gin-gonic/gin"
)

func InitUpstream(r *gin.RouterGroup) {
	r.GET("/", upstream.Index)
	r.POST("/", upstream.Create)
	r.PUT("/:id", upstream.Update)
	r.DELETE("/:id", upstream.Delete)
}
