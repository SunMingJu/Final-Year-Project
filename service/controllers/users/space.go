package users

import (
	"simple-video-net/controllers"
	receive "simple-video-net/interaction/receive/users"
	"simple-video-net/logic/users"

	"github.com/gin-gonic/gin"
)

type SpaceControllers struct {
	controllers.BaseControllers
}

// GetSpaceIndividual
func (sp SpaceControllers) GetSpaceIndividual(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.GetSpaceIndividualReceiveStruct)); err == nil {
		results, err := users.GetSpaceIndividual(rec, uid)
		sp.Response(ctx, results, err)
	}
}

// GetReleaseInformation
func (sp SpaceControllers) GetReleaseInformation(ctx *gin.Context) {
	if rec, err := controllers.ShouldBind(ctx, new(receive.GetReleaseInformationReceiveStruct)); err == nil {
		results, err := users.GetReleaseInformation(rec)
		sp.Response(ctx, results, err)
	}
}

// GetAttentionList
func (sp SpaceControllers) GetAttentionList(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.GetAttentionListReceiveStruct)); err == nil {
		results, err := users.GetAttentionList(rec, uid)
		sp.Response(ctx, results, err)
	}
}

// GetVermicelliList
func (sp SpaceControllers) GetVermicelliList(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.GetVermicelliListReceiveStruct)); err == nil {
		results, err := users.GetVermicelliList(rec, uid)
		sp.Response(ctx, results, err)
	}
}
