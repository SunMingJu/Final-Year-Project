package home

import (
	"easy-video-net/controllers"
	receive "easy-video-net/interaction/receive/home"
	"easy-video-net/logic/home"
	"github.com/gin-gonic/gin"
)

type Controllers struct {
	controllers.BaseControllers
}

//GetHomeInfo 
func (c Controllers) GetHomeInfo(ctx *gin.Context) {
	if rec, err := controllers.ShouldBind(ctx, new(receive.GetHomeInfoReceiveStruct)); err == nil {
		results, err := home.GetHomeInfo(rec)
		c.Response(ctx, results, err)
	}
}
