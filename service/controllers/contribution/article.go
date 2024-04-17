package contribution

import (
	"easy-video-net/controllers"
	receive "easy-video-net/interaction/receive/contribution/article"
	"easy-video-net/logic/contribution"
	"github.com/gin-gonic/gin"
)

//CreateArticleContribution 
func (c Controllers) CreateArticleContribution(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.CreateArticleContributionReceiveStruct)); err == nil {
		results, err := contribution.CreateArticleContribution(rec, uid)
		c.Response(ctx, results, err)
	}
}

//UpdateArticleContribution 
func (c Controllers) UpdateArticleContribution(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.UpdateArticleContributionReceiveStruct)); err == nil {
		results, err := contribution.UpdateArticleContribution(rec, uid)
		c.Response(ctx, results, err)
	}
}

//DeleteArticleByID 
func (c Controllers) DeleteArticleByID(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.DeleteArticleByIDReceiveStruct)); err == nil {
		results, err := contribution.DeleteArticleByID(rec, uid)
		c.Response(ctx, results, err)
	}
}

//GetArticleContributionList 
func (c Controllers) GetArticleContributionList(ctx *gin.Context) {
	if rec, err := controllers.ShouldBind(ctx, new(receive.GetArticleContributionListReceiveStruct)); err == nil {
		results, err := contribution.GetArticleContributionList(rec)
		c.Response(ctx, results, err)
	}
}

//GetArticleContributionListByUser 
func (c Controllers) GetArticleContributionListByUser(ctx *gin.Context) {
	if rec, err := controllers.ShouldBind(ctx, new(receive.GetArticleContributionListByUserReceiveStruct)); err == nil {
		results, err := contribution.GetArticleContributionListByUser(rec)
		c.Response(ctx, results, err)
	}
}

//GetArticleContributionByID 
func (c Controllers) GetArticleContributionByID(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.GetArticleContributionByIDReceiveStruct)); err == nil {
		results, err := contribution.GetArticleContributionByID(rec, uid)
		c.Response(ctx, results, err)
	}
}

//ArticlePostComment 
func (c Controllers) ArticlePostComment(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.ArticlesPostCommentReceiveStruct)); err == nil {
		results, err := contribution.ArticlePostComment(rec, uid)
		c.Response(ctx, results, err)
	}
}

//GetArticleComment 
func (c Controllers) GetArticleComment(ctx *gin.Context) {
	if rec, err := controllers.ShouldBind(ctx, new(receive.GetArticleCommentReceiveStruct)); err == nil {
		results, err := contribution.GetArticleComment(rec)
		c.Response(ctx, results, err)
	}
}

//GetArticleClassificationList 
func (c Controllers) GetArticleClassificationList(ctx *gin.Context) {
	results, err := contribution.GetArticleClassificationList()
	c.Response(ctx, results, err)
}

//GetArticleTotalInfo 
func (c Controllers) GetArticleTotalInfo(ctx *gin.Context) {
	results, err := contribution.GetArticleTotalInfo()
	c.Response(ctx, results, err)
}

//GetArticleManagementList 
func (c Controllers) GetArticleManagementList(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.GetArticleManagementListReceiveStruct)); err == nil {
		results, err := contribution.GetArticleManagementList(rec, uid)
		c.Response(ctx, results, err)
	}
}
