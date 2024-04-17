package controllers

import (
	"easy-video-net/utils/response"
	"easy-video-net/utils/validator"
	"fmt"
	"github.com/gin-gonic/gin"
)

type BaseControllers struct {
}

//Response 
func (c BaseControllers) Response(ctx *gin.Context, results interface{}, err error) {
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, results)
}

//ShouldBind 
func ShouldBind[T interface{}](ctx *gin.Context, data T) (t T, err error) {
	if err := ctx.ShouldBind(data); err != nil {
		fmt.Println(err)
		validator.CheckParams(ctx, err)
		return t, err
	}
	return data, nil
}
