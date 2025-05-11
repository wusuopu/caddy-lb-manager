package caddyfile

import (
	"app/config"
	"app/di"
	"app/schemas"
	"os"

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

	os.WriteFile(config.Config.Caddy.ConfigPath, []byte(content), 0644)

	ret, err = di.Service.CaddyfileService.TouchReloadTime()
	if ret != true {
		schemas.MakeErrorResponse(ctx, err, 400)
		return
	}

	schemas.MakeResponse(ctx, nil, nil)
}

func ListCertificate(ctx *gin.Context) {
	data := di.Service.CaddyfileService.ListCertificate()
	schemas.MakeResponse(ctx, data, nil)
}