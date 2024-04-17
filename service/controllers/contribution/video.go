package contribution

import (
	"easy-video-net/controllers"
	receive "easy-video-net/interaction/receive/contribution/video"
	"easy-video-net/logic/contribution"
	"easy-video-net/utils/response"
	"easy-video-net/utils/validator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Controllers struct {
	controllers.BaseControllers
}

//CreateVideoContribution 
func (c Controllers) CreateVideoContribution(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.CreateVideoContributionReceiveStruct)); err == nil {
		results, err := contribution.CreateVideoContribution(rec, uid)
		c.Response(ctx, results, err)
	}
}

//UpdateVideoContribution 
func (c Controllers) UpdateVideoContribution(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.UpdateVideoContributionReceiveStruct)); err == nil {
		results, err := contribution.UpdateVideoContribution(rec, uid)
		c.Response(ctx, results, err)
	}
}

// GetVideoContributionByID  
func (c Controllers) GetVideoContributionByID(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.GetVideoContributionByIDReceiveStruct)); err == nil {
		results, err := contribution.GetVideoContributionByID(rec, uid)
		c.Response(ctx, results, err)
	}

}

// SendVideoBarrage  
func (c Controllers) SendVideoBarrage(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	SendVideoBarrageReceive := new(receive.SendVideoBarrageReceiveStruct)
	//Solving duplicate binds with ShouldBindBodyWith
	if err := ctx.ShouldBindBodyWith(SendVideoBarrageReceive, binding.JSON); err == nil {
		results, err := contribution.SendVideoBarrage(SendVideoBarrageReceive, uid)
		if err != nil {
			response.Error(ctx, err.Error())
			return
		}
		response.BarrageSuccess(ctx, results)
	} else {
		validator.CheckParams(ctx, err)
	}
}

// GetVideoBarrage  
func (c Controllers) GetVideoBarrage(ctx *gin.Context) {
	GetVideoBarrageReceive := new(receive.GetVideoBarrageReceiveStruct)
	GetVideoBarrageReceive.ID = ctx.Query("id")
	results, err := contribution.GetVideoBarrage(GetVideoBarrageReceive)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.BarrageSuccess(ctx, results)

}

// GetVideoBarrageList  
func (c Controllers) GetVideoBarrageList(ctx *gin.Context) {
	GetVideoBarrageReceive := new(receive.GetVideoBarrageListReceiveStruct)
	GetVideoBarrageReceive.ID = ctx.Query("id")
	results, err := contribution.GetVideoBarrageList(GetVideoBarrageReceive)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.BarrageSuccess(ctx, results)
}

//VideoPostComment 
func (c Controllers) VideoPostComment(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.VideosPostCommentReceiveStruct)); err == nil {
		results, err := contribution.VideoPostComment(rec, uid)
		c.Response(ctx, results, err)
	}
}

//GetVideoComment 
func (c Controllers) GetVideoComment(ctx *gin.Context) {
	if rec, err := controllers.ShouldBind(ctx, new(receive.GetVideoCommentReceiveStruct)); err == nil {
		results, err := contribution.GetVideoComment(rec)
		c.Response(ctx, results, err)
	}
}

//GetVideoManagementList 
func (c Controllers) GetVideoManagementList(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.GetVideoManagementListReceiveStruct)); err == nil {
		results, err := contribution.GetVideoManagementList(rec, uid)
		c.Response(ctx, results, err)
	}
}

//DeleteVideoByID 
func (c Controllers) DeleteVideoByID(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.DeleteVideoByIDReceiveStruct)); err == nil {
		results, err := contribution.DeleteVideoByID(rec, uid)
		c.Response(ctx, results, err)
	}
}

//LikeVideo 
func (c Controllers) LikeVideo(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.LikeVideoReceiveStruct)); err == nil {
		results, err := contribution.LikeVideo(rec, uid)
		c.Response(ctx, results, err)
	}
}
