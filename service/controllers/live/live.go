package live

import (
	"simple-video-net/controllers"
	receive "simple-video-net/interaction/receive/live"
	"simple-video-net/logic/live"

	"github.com/gin-gonic/gin"
)

type LivesControllers struct {
	controllers.BaseControllers
}

// GetLiveRoom
func (lv LivesControllers) GetLiveRoom(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	results, err := live.GetLiveRoom(uid)
	lv.Response(ctx, results, err)
}

// GetLiveRoomInfo
func (lv LivesControllers) GetLiveRoomInfo(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.GetLiveRoomInfoReceiveStruct)); err == nil {
		results, err := live.GetLiveRoomInfo(rec, uid)
		lv.Response(ctx, results, err)
	}
}

// GetBeLiveList
func (lv LivesControllers) GetBeLiveList(ctx *gin.Context) {
	results, err := live.GetBeLiveList()
	lv.Response(ctx, results, err)
}
