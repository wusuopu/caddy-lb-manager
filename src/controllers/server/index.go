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
	body := helper.GetJSONBody(ctx)

	port, _ := helper.GetJSONInt64(body, "Port")
	results := di.Container.DB.
		Where(&models.Server{Host: helper.GetJSONString(body, "Host"), Port: uint(port)}).
		First(&models.Server{})

	if results.RowsAffected > 0 {
		schemas.MakeErrorResponse(ctx, "Server is already exists with this Host and Port", 400)
		return
	}

	obj := models.Server{
		Name: helper.GetJSONString(body, "Name"),
		Host: helper.GetJSONString(body, "Host"),
		Port: uint(port),
		EnableSSL: helper.GetJSONBool(body, "EnableSSL"),
		Enable: helper.GetJSONBool(body, "Enable"),
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

	body := helper.GetJSONBody(ctx)
	port, _ := helper.GetJSONInt64(body, "Port")

	obj.Name = helper.GetJSONString(body, "Name")
	obj.Host = helper.GetJSONString(body, "Host")
	obj.Port = uint(port)
	obj.EnableSSL = helper.GetJSONBool(body, "EnableSSL")
	obj.Enable = helper.GetJSONBool(body, "Enable")

	results = di.Container.DB.Save(&obj)

	if results.Error != nil {
		schemas.MakeErrorResponse(ctx, results.Error, 400)
		return
	}

	schemas.MakeResponse(ctx, obj, nil)
}