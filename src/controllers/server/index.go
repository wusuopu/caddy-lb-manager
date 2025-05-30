package server

import (
	"app/di"
	"app/models"
	"app/schemas"
	"app/utils/helper"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Index(ctx *gin.Context) {
	var data []models.Server

	di.Container.DB.Find(&data)

	schemas.MakeResponse(ctx, data, nil)
}

func Create(ctx *gin.Context) {
	var parser helper.JSONParser
	parser.GetJSONBody(ctx)

	port, _ := parser.GetJSONInt64("Port")
	host, _ := parser.GetJSONString("Host")
	results := di.Container.DB.
		Where(&models.Server{Host: host, Port: uint(port)}).
		First(&models.Server{})

	if results.RowsAffected > 0 {
		schemas.MakeErrorResponse(ctx, "Server is already exists with this Host and Port", 400)
		return
	}

	name, _ := parser.GetJSONString("Name")
	enableSSL, _ := parser.GetJSONBool("EnableSSL")
	enable, _ := parser.GetJSONBool("Enable")

	obj := models.Server{
		Name: name,
		Host: host,
		Port: uint(port),
		EnableSSL: enableSSL,
		Enable: enable,
	}
	results = di.Container.DB.Create(&obj)
	if results.Error != nil {
		schemas.MakeErrorResponse(ctx, results.Error, 500)
		return
	}

	schemas.MakeResponse(ctx, obj, nil)
}
func Show(ctx *gin.Context) {
	var obj models.Server
	results := di.Container.DB.First(&obj, ctx.Param("serverId"))
	if results.Error != nil {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			schemas.MakeErrorResponse(ctx, "Server not found", 404)
		} else {
			schemas.MakeErrorResponse(ctx, results.Error, 500)
		}
		return
	}

	schemas.MakeResponse(ctx, obj, nil)
}

func Delete(ctx *gin.Context) {
	var obj models.Server
	results := di.Container.DB.First(&obj, ctx.Param("serverId"))
	if results.Error != nil {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			schemas.MakeErrorResponse(ctx, "Server not found", 404)
		} else {
			schemas.MakeErrorResponse(ctx, results.Error, 500)
		}
		return
	}

	results = di.Container.DB.Unscoped().Delete(&obj)
	if results.Error != nil {
		schemas.MakeErrorResponse(ctx, results.Error, 500)
		return
	}

	di.Container.DB.Unscoped().Delete(&models.Route{}, "server_id = ?", ctx.Param("serverId"))

	schemas.MakeResponse(ctx, obj, nil)
}
func Update(ctx *gin.Context) {
	var obj models.Server
	results := di.Container.DB.First(&obj, ctx.Param("serverId"))
	if results.Error != nil {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			schemas.MakeErrorResponse(ctx, "Server not found", 404)
		} else {
			schemas.MakeErrorResponse(ctx, results.Error, 500)
		}
		return
	}

	var parser helper.JSONParser
	parser.GetJSONBody(ctx)

	port, _ := parser.GetJSONInt64("Port")
	host, _ := parser.GetJSONString("Host")
	name, _ := parser.GetJSONString("Name")
	enableSSL, _ := parser.GetJSONBool("EnableSSL")
	enable, _ := parser.GetJSONBool("Enable")

	obj.Name = name
	obj.Host = host
	obj.Port = uint(port)
	obj.EnableSSL = enableSSL
	obj.Enable = enable

	results = di.Container.DB.Save(&obj)

	if results.Error != nil {
		schemas.MakeErrorResponse(ctx, results.Error, 400)
		return
	}

	schemas.MakeResponse(ctx, obj, nil)
}