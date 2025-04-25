package routes

import (
	"app/controllers/route"

	"github.com/gin-gonic/gin"
)

func InitRoute(r *gin.RouterGroup) {
	r.GET("/:serverId/routes", route.Index)
	r.POST("/:serverId/routes", route.Create)
	r.GET("/:serverId/routes/:routeId", route.Show)
	r.PUT("/:serverId/routes/:routeId", route.Update)
	r.DELETE("/:serverId/routes/:routeId", route.Delete)
}
