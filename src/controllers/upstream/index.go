package upstream

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
	var data []models.UpStream

	di.Container.DB.Find(&data)

	schemas.MakeResponse(ctx, data, nil)
}
func Create(ctx *gin.Context) {
	var parser helper.JSONParser
	parser.GetJSONBody(ctx)

	name, _ := parser.GetJSONString("Name")
	schema, _ := parser.GetJSONString("Scheme")
	address, _ := parser.GetJSONString("Address")

	obj := models.UpStream{
		Name: name,
		Scheme: schema,
		Address: address,
	}
	results := di.Container.DB.Create(&obj)
	if results.Error != nil {
		schemas.MakeErrorResponse(ctx, results.Error, 500)
		return
	}

	schemas.MakeResponse(ctx, obj, nil)
}

func Delete(ctx *gin.Context) {
	var obj models.UpStream
	results := di.Container.DB.First(&obj, ctx.Param("id"))
	if results.Error != nil {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			schemas.MakeErrorResponse(ctx, "Authentication not found", 404)
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
	var obj models.UpStream
	results := di.Container.DB.First(&obj, ctx.Param("id"))
	if results.Error != nil {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			schemas.MakeErrorResponse(ctx, "Authentication not found", 404)
		} else {
			schemas.MakeErrorResponse(ctx, results.Error, 500)
		}
		return
	}

	var parser helper.JSONParser
	parser.GetJSONBody(ctx)

	name, _ := parser.GetJSONString("Name")
	schema, _ := parser.GetJSONString("Scheme")
	address, _ := parser.GetJSONString("Address")

	obj.Name = name
	obj.Scheme = schema
	obj.Address = address

	results = di.Container.DB.Save(&obj)

	if results.Error != nil {
		schemas.MakeErrorResponse(ctx, results.Error, 400)
		return
	}

	schemas.MakeResponse(ctx, obj, nil)
}
