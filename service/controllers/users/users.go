package users

import (
	"simple-video-net/controllers"
	receive "simple-video-net/interaction/receive/users"
	"simple-video-net/logic/users"

	"github.com/gin-gonic/gin"
)

type UserControllers struct {
	controllers.BaseControllers
}

// GetUserInfo
func (us UserControllers) GetUserInfo(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	results, err := users.GetUserInfo(uid)
	us.Response(ctx, results, err)
}

// SetUserInfo
func (us UserControllers) SetUserInfo(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.SetUserInfoReceiveStruct)); err == nil {
		results, err := users.SetUserInfo(rec, uid)
		us.Response(ctx, results, err)
	}
}

// DetermineNameExists
func (us UserControllers) DetermineNameExists(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.DetermineNameExistsStruct)); err == nil {
		results, err := users.DetermineNameExists(rec, uid)
		us.Response(ctx, results, err)
	}
}

// UpdateAvatar
func (us UserControllers) UpdateAvatar(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.UpdateAvatarStruct)); err == nil {
		results, err := users.UpdateAvatar(rec, uid)
		us.Response(ctx, results, err)
	}
}

// GetLiveData
func (us UserControllers) GetLiveData(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	results, err := users.GetLiveData(uid)
	us.Response(ctx, results, err)
}

// SaveLiveData
func (us UserControllers) SaveLiveData(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.SaveLiveDataReceiveStruct)); err == nil {
		results, err := users.SaveLiveData(rec, uid)
		us.Response(ctx, results, err)
	}
}

// SendEmailVerificationCodeByChangePassword
func (us UserControllers) SendEmailVerificationCodeByChangePassword(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	results, err := users.SendEmailVerificationCodeByChangePassword(uid)
	us.Response(ctx, results, err)
}

// ChangePassword
func (us UserControllers) ChangePassword(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.ChangePasswordReceiveStruct)); err == nil {
		results, err := users.ChangePassword(rec, uid)
		us.Response(ctx, results, err)
	}
}

// Attention
func (us UserControllers) Attention(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.AttentionReceiveStruct)); err == nil {
		results, err := users.Attention(rec, uid)
		us.Response(ctx, results, err)
	}
}

// CreateFavorites
func (us UserControllers) CreateFavorites(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.CreateFavoritesReceiveStruct)); err == nil {
		results, err := users.CreateFavorites(rec, uid)
		us.Response(ctx, results, err)
	}
}

// DeleteFavorites
func (us UserControllers) DeleteFavorites(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.DeleteFavoritesReceiveStruct)); err == nil {
		results, err := users.DeleteFavorites(rec, uid)
		us.Response(ctx, results, err)
	}
}

// GetFavoritesList
func (us UserControllers) GetFavoritesList(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	results, err := users.GetFavoritesList(uid)
	us.Response(ctx, results, err)
}

// FavoriteVideo
func (us UserControllers) FavoriteVideo(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.FavoriteVideoReceiveStruct)); err == nil {
		results, err := users.FavoriteVideo(rec, uid)
		us.Response(ctx, results, err)
	}
}

// GetFavoritesListByFavoriteVideo
func (us UserControllers) GetFavoritesListByFavoriteVideo(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.GetFavoritesListByFavoriteVideoReceiveStruct)); err == nil {
		results, err := users.GetFavoritesListByFavoriteVideo(rec, uid)
		us.Response(ctx, results, err)
	}
}

// GetFavoriteVideoList
func (us UserControllers) GetFavoriteVideoList(ctx *gin.Context) {
	if rec, err := controllers.ShouldBind(ctx, new(receive.GetFavoriteVideoListReceiveStruct)); err == nil {
		results, err := users.GetFavoriteVideoList(rec)
		us.Response(ctx, results, err)
	}
}

// GetRecordList
func (us UserControllers) GetRecordList(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.GetRecordListReceiveStruct)); err == nil {
		results, err := users.GetRecordList(rec, uid)
		us.Response(ctx, results, err)
	}
}

// ClearRecord
func (us UserControllers) ClearRecord(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	results, err := users.ClearRecord(uid)
	us.Response(ctx, results, err)
}

// DeleteRecordByID
func (us UserControllers) DeleteRecordByID(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.DeleteRecordByIDReceiveStruct)); err == nil {
		results, err := users.DeleteRecordByID(rec, uid)
		us.Response(ctx, results, err)
	}
}

// GetNoticeList
func (us UserControllers) GetNoticeList(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.GetNoticeListReceiveStruct)); err == nil {
		results, err := users.GetNoticeList(rec, uid)
		us.Response(ctx, results, err)
	}
}

// GetChatList
func (us UserControllers) GetChatList(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	results, err := users.GetChatList(uid)
	us.Response(ctx, results, err)
}

// GetChatHistoryMsg
func (us UserControllers) GetChatHistoryMsg(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.GetChatHistoryMsgStruct)); err == nil {
		results, err := users.GetChatHistoryMsg(rec, uid)
		us.Response(ctx, results, err)
	}
}

// PersonalLetter
func (us UserControllers) PersonalLetter(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.PersonalLetterReceiveStruct)); err == nil {
		results, err := users.PersonalLetter(rec, uid)
		us.Response(ctx, results, err)
	}
}

// DeleteChatItem
func (us UserControllers) DeleteChatItem(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if rec, err := controllers.ShouldBind(ctx, new(receive.DeleteChatItemReceiveStruct)); err == nil {
		results, err := users.DeleteChatItem(rec, uid)
		us.Response(ctx, results, err)
	}
}
