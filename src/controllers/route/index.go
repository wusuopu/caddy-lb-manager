package route

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
	var data []models.Route

	di.Container.DB.Where("server_id = ?", ctx.Param("serverId")).Find(&data)

	schemas.MakeResponse(ctx, data, nil)
}
func Create(ctx *gin.Context) {
	var server models.Server
	results := di.Container.DB.First(&server, ctx.Param("serverId"))
	if results.Error != nil {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			schemas.MakeErrorResponse(ctx, "This Server is not exists", 404)
		} else {
			schemas.MakeErrorResponse(ctx, results.Error, 500)
		}
		return
	}

	var parser helper.JSONParser
	parser.GetJSONBody(ctx)

	name, _ := parser.GetJSONString("Name")
	methods, _ := parser.GetJSONString("Methods")
	path, _ := parser.GetJSONString("Path")
	headerUp := parser.GetJSONItem("HeaderUp").MarshalTo(nil)
	headerDown := parser.GetJSONItem("HeaderDown").MarshalTo(nil)
	stripPath, _ := parser.GetJSONBool("StripPath")
	enable, _ := parser.GetJSONBool("Enable")
	upstreamId, _ := parser.GetJSONInt64("UpStreamId")
	authenticationId, _ := parser.GetJSONInt64("AuthenticationId")

	obj := models.Route{
		Name: name,
		Methods: methods,
		Path: path,
		HeaderUp: headerUp,
		HeaderDown: headerDown,
		StripPath: stripPath,
		Enable: enable,
		UpStreamId: uint(upstreamId),
		ServerId: uint(server.ID),
		AuthenticationId: uint(authenticationId),
	}
	results = di.Container.DB.Create(&obj)
	if results.Error != nil {
		schemas.MakeErrorResponse(ctx, results.Error, 500)
		return
	}

	schemas.MakeResponse(ctx, obj, nil)
}
func Delete(ctx *gin.Context) {
	var obj models.Route
	results := di.Container.DB.First(&obj, ctx.Param("routeId"))
	if results.Error != nil {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			schemas.MakeErrorResponse(ctx, "Route not found", 404)
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
	var obj models.Route
	results := di.Container.DB.First(&obj, ctx.Param("routeId"))
	if results.Error != nil {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			schemas.MakeErrorResponse(ctx, "Route not found", 404)
		} else {
			schemas.MakeErrorResponse(ctx, results.Error, 500)
		}
		return
	}

	var parser helper.JSONParser
	parser.GetJSONBody(ctx)

	name, _ := parser.GetJSONString("Name")
	methods, _ := parser.GetJSONString("Methods")
	path, _ := parser.GetJSONString("Path")
	headerUp := parser.GetJSONItem("HeaderUp").MarshalTo(nil)
	headerDown := parser.GetJSONItem("HeaderDown").MarshalTo(nil)
	stripPath, _ := parser.GetJSONBool("StripPath")
	enable, _ := parser.GetJSONBool("Enable")
	upstreamId, _ := parser.GetJSONInt64("UpStreamId")
	authenticationId, _ := parser.GetJSONInt64("AuthenticationId")

	obj.Name = name
	obj.Methods = methods
	obj.Path = path
	obj.HeaderUp = headerUp
	obj.HeaderDown = headerDown
	obj.StripPath = stripPath
	obj.Enable = enable
	obj.UpStreamId = uint(upstreamId)
	obj.AuthenticationId = uint(authenticationId)

	results = di.Container.DB.Save(&obj)
	if results.Error != nil {
		schemas.MakeErrorResponse(ctx, results.Error, 500)
		return
	}

	schemas.MakeResponse(ctx, obj, nil)
}
