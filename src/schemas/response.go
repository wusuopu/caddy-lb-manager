package schemas

import (
	"app/utils/helper"

	"github.com/gin-gonic/gin"
)

type ResponseMeta struct {
	Pagination *helper.Pagination
}
type ResponseBody struct {
	Data any
	Error any
	Meta *ResponseMeta
}

func MakeResponse(ctx *gin.Context, data any, pagination *helper.Pagination) {
	body := ResponseBody{
		Data: data,
		Meta: &ResponseMeta{
			Pagination: pagination,
		},
	}
	ctx.JSON(200, body)
}

func MakeErrorResponse(ctx *gin.Context, err any, status int) {
	body := ResponseBody{}
	switch err.(type) {
		default:
			body.Error = err
		case error:
			body.Error = err.(error).Error()
	}
	ctx.JSON(status, body)
}