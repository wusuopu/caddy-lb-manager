package dashboard

import (
	"app/di"
	"app/models"
	"app/schemas"

	"github.com/gin-gonic/gin"
)

func Index(ctx *gin.Context) {
	var serverCount int64
	var routeCount int64
	var upstreamCount int64

	di.Container.DB.Model(&models.Server{}).Count(&serverCount)
	di.Container.DB.Model(&models.Route{}).Count(&routeCount)
	di.Container.DB.Model(&models.UpStream{}).Count(&upstreamCount)

	data := map[string]interface{}{
		"serverCount": serverCount,
		"routeCount": routeCount,
		"upstreamCount": upstreamCount,
	}
	schemas.MakeResponse(ctx, data, nil)
}