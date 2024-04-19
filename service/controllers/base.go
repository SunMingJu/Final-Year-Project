package controllers

import (
	
	"simple-video-net/utils/response"
	"simple-video-net/utils/validator"
	"simple-ideo-net/global"
	"github.com/gin-gonic/gin"
)

type BaseControllers struct {
}

// Response
func (c BaseControllers) Response(ctx *gin.Context, results interface{}, err error) {
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, results)
}

// ShouldBind
func ShouldBind[T interface{}](ctx *gin.Context, data T) (t T, err error) {
	if err := ctx.ShouldBind(data); err != nil {
		global.Logger.Errorf("Request incoming parameter binding failed type：%T ,wrong reason : %s ", t, err.Error())
		validator.CheckParams(ctx, err)
		return t, err
	}
	return data, nil
}
