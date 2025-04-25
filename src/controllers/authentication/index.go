package authentication

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
	// don't return password field
	var data []struct {
		gorm.Model
		Name string
		Username string
		HashedPw string
	}

	di.Container.DB.Model(&models.Authentication{}).Find(&data)

	schemas.MakeResponse(ctx, data, nil)
}

func Create(ctx *gin.Context) {
	var parser helper.JSONParser
	parser.GetJSONBody(ctx)

	name, _ := parser.GetJSONString("Name")
	username, _ := parser.GetJSONString("Username")
	password, _ := parser.GetJSONString("Password")

	obj := models.Authentication{
		Name: name,
		Username: username,
		Password: password,
	}
	hashedPassword, err := di.Service.CaddyfileService.HashPassword(obj.Password)
	if err != nil {
		schemas.MakeErrorResponse(ctx, err, 500)
		return
	}
	obj.HashedPw = hashedPassword

	results := di.Container.DB.Create(&obj)
	if results.Error != nil {
		schemas.MakeErrorResponse(ctx, results.Error, 500)
		return
	}

	schemas.MakeResponse(ctx, obj, nil)
}

func Update(ctx *gin.Context) {
	var obj models.Authentication
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
	username, _ := parser.GetJSONString("Username")
	password, _ := parser.GetJSONString("Password")

	obj.Name = name
	obj.Username = username
	obj.Password = password

	hashedPassword, err := di.Service.CaddyfileService.HashPassword(obj.Password)
	if err != nil {
		schemas.MakeErrorResponse(ctx, err, 500)
		return
	}
	obj.HashedPw = hashedPassword

	results = di.Container.DB.Save(&obj)

	if results.Error != nil {
		schemas.MakeErrorResponse(ctx, results.Error, 400)
		return
	}

	schemas.MakeResponse(ctx, obj, nil)
}

func Delete(ctx *gin.Context) {
	var obj models.Authentication
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