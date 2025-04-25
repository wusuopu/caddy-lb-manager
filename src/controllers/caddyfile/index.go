package caddyfile

import (
	"app/di"
	"app/schemas"

	"github.com/gin-gonic/gin"
)

func Index(ctx *gin.Context) {
	content, err := di.Service.CaddyfileService.GenCaddyfile()
	if err != nil {
		schemas.MakeErrorResponse(ctx, err, 400)
		return
	}

	schemas.MakeResponse(ctx, content, nil)
}

func Reload(ctx *gin.Context) {
	content, err := di.Service.CaddyfileService.GenCaddyfile()
	if err != nil {
		schemas.MakeErrorResponse(ctx, err, 400)
		return
	}

	ret, err := di.Service.CaddyfileService.Reload(content)
	if ret != true {
		schemas.MakeErrorResponse(ctx, err, 400)
		return
	}

	schemas.MakeResponse(ctx, nil, nil)
}