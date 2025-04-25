package routes

import (
	"app/controllers/upstream"

	"github.com/gin-gonic/gin"
)

func InitUpstream(r *gin.RouterGroup) {
	r.GET("/", upstream.Index)
	r.POST("/", upstream.Create)
	r.GET("/:upstreamId", upstream.Show)
	r.PUT("/:upstreamId", upstream.Update)
	r.DELETE("/:upstreamId", upstream.Delete)
}
