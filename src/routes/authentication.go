package routes

import (
	"app/controllers/authentication"

	"github.com/gin-gonic/gin"
)

func InitAuthentication(r *gin.RouterGroup) {
	r.GET("/", authentication.Index)
	r.POST("/", authentication.Create)
	r.PUT("/:id", authentication.Update)
	r.DELETE("/:id", authentication.Delete)
}
