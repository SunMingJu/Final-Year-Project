package callback

import (
	"simple-video-net/controllers"
	receive "simple-video-net/interaction/receive/callback"
	"simple-video-net/logic/callback"
	"github.com/gin-gonic/gin"
)

type Controllers struct {
	controllers.BaseControllers
}

//AliyunTranscodingMedia Alibaba Cloud Media transcoding successful callback
func (c *Controllers) AliyunTranscodingMedia(ctx *gin.Context) {
	if rec, err := controllers.ShouldBind(ctx, new(receive.AliyunMediaCallback[receive.AliyunTranscodingMediaStruct])); err == nil {
		results, err := callback.AliyunTranscodingMedia(rec)
		c.Response(ctx, results, err)
	}
}