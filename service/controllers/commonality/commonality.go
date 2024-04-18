package commonality

import (
	"simple-video-net/controllers"
	receive "simple-video-net/interaction/receive/commonality"
	"simple-video-net/logic/commonality"

	"github.com/gin-gonic/gin"
)

type Controllers struct {
	controllers.BaseControllers
}

// OssSTS
func (c *Controllers) OssSTS(ctx *gin.Context) {
	results, err := commonality.OssSTS()
	c.Response(ctx, results, err)
}

// Upload
func (c *Controllers) Upload(ctx *gin.Context) {
	file, _ := ctx.FormFile("file")
	results, err := commonality.Upload(file, ctx)
	c.Response(ctx, results, err)
}

// UploadSlice
func (c *Controllers) UploadSlice(ctx *gin.Context) {
	file, _ := ctx.FormFile("file")
	results, err := commonality.UploadSlice(file, ctx)
	c.Response(ctx, results, err)
}

// UploadCheck
func (c *Controllers) UploadCheck(ctx *gin.Context) {
	if rec, err := controllers.ShouldBind(ctx, new(receive.UploadCheckStruct)); err == nil {
		results, err := commonality.UploadCheck(rec)
		c.Response(ctx, results, err)
	}
}

// UploadMerge
func (c *Controllers) UploadMerge(ctx *gin.Context) {
	if rec, err := controllers.ShouldBind(ctx, new(receive.UploadMergeStruct)); err == nil {
		results, err := commonality.UploadMerge(rec)
		c.Response(ctx, results, err)
	}
}

// UploadingMethod
func (c *Controllers) UploadingMethod(ctx *gin.Context) {
	if rec, err := controllers.ShouldBind(ctx, new(receive.UploadingMethodStruct)); err == nil {
		results, err := commonality.UploadingMethod(rec)
		c.Response(ctx, results, err)
	}
}

// UploadingDir
func (c *Controllers) UploadingDir(ctx *gin.Context) {
	if rec, err := controllers.ShouldBind(ctx, new(receive.UploadingDirStruct)); err == nil {
		results, err := commonality.UploadingDir(rec)
		c.Response(ctx, results, err)
	}
}

// GetFullPathOfImage
func (c *Controllers) GetFullPathOfImage(ctx *gin.Context) {
	if rec, err := controllers.ShouldBind(ctx, new(receive.GetFullPathOfImageMethodStruct)); err == nil {
		results, err := commonality.GetFullPathOfImage(rec)
		c.Response(ctx, results, err)
	}
}

// Search
func (c *Controllers) Search(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.SearchStruct)); err == nil {
		results, err := commonality.Search(rec, uid)
		c.Response(ctx, results, err)
	}
}
