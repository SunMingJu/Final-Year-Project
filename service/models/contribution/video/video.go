package video

import (
	"simple-video-net/global"
	"simple-video-net/models/common"
	"simple-video-net/models/contribution/video/barrage"
	"simple-video-net/models/contribution/video/comments"
	"simple-video-net/models/contribution/video/like"
	"simple-video-net/models/users"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type VideosContribution struct {
	common.PublicModel
	Uid           uint           `json:"uid" gorm:"column:uid"`
	Title         string         `json:"title" gorm:"column:title"`
	Video         datatypes.JSON `json:"video" gorm:"column:video"` //默认1080p
	Video720p     datatypes.JSON `json:"video_720p" gorm:"column:video_720p"`
	Video480p     datatypes.JSON `json:"video_480p" gorm:"column:video_480p"`
	Video360p     datatypes.JSON `json:"video_360p" gorm:"column:video_360p"`
	MediaID string `json:"media_id" gorm:"column:media_id"`
	Cover         datatypes.JSON `json:"cover" gorm:"column:cover"`
	VideoDuration int64          `json:"video_duration" gorm:"column:video_duration"`
	Reprinted     int8           `json:"reprinted" gorm:"column:reprinted"`
	Timing        int8           `json:"timing" gorm:"column:timing"`
	TimingTime    time.Time      `json:"timingTime"  gorm:"column:timing_Time"`
	Label         string         `json:"label" gorm:"column:label"`
	Introduce     string         `json:"introduce" gorm:"column:introduce"`
	Heat          int            `json:"heat" gorm:"column:heat"

	UserInfo users.User           `json:"user_info" gorm:"foreignKey:Uid"`
	Likes    like.LikesList       `json:"likes" gorm:"foreignKey:VideoID" `
	Comments comments.CommentList `json:"comments" gorm:"foreignKey:VideoID"`
	Barrage  barrage.BarragesList `json:"barrage" gorm:"foreignKey:VideoID"`
}

type VideosContributionList []VideosContribution

func (VideosContribution) TableName() string {
	return "lv_video_contribution"
}

// Create
func (vc *VideosContribution) Create() bool {
	err := global.Db.Create(&vc).Error
	if err != nil {
		return false
	}
	return true
}

// Delete
func (vc *VideosContribution) Delete(id uint, uid uint) bool {
	if global.Db.Where("id", id).Find(&vc).Error != nil {
		return false
	}
	if vc.Uid != uid {
		return false
	}
	if global.Db.Delete(&vc).Error != nil {
		return false
	}
	return true
}

// Update
func (vc *VideosContribution) Update(info map[string]interface{}) bool {
	err := global.Db.Model(vc).Updates(info).Error
	if err != nil {
		return false
	}
	return true
}

func (vc *VideosContribution) Save() bool {
	err := global.Db.Save(vc).Error
	if err != nil {
		return false
	}
	return true
}

// FindByID
func (vc *VideosContribution) FindByID(id uint) error {
	return global.Db.Where("id", id).Preload("Likes").Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Preload("UserInfo").Order("created_at desc")
	}).Preload("Barrage").Preload("UserInfo").Order("created_at desc").Find(&vc).Error
}

// GetVideoComments
func (vc *VideosContribution) GetVideoComments(ID uint, info common.PageInfo) bool {
	err := global.Db.Where("id", ID).Preload("Likes").Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Preload("UserInfo").Order("created_at desc").Limit(info.Size).Offset((info.Page - 1) * info.Size)
	}).Find(vc).Error
	if err != nil {
		return false
	}
	return true
}

// Watch
func (vc *VideosContribution) Watch(id uint) error {
	return global.Db.Model(vc).Where("id", id).Updates(map[string]interface{}{"heat": gorm.Expr("Heat  + ?", 1)}).Error
}

// GetVideoListBySpace
func (vl *VideosContributionList) GetVideoListBySpace(id uint) error {
	return global.Db.Where("uid", id).Preload("Likes").Preload("Comments").Preload("Barrage").Order("created_at desc").Find(&vl).Error
}

// GetDiscussVideoCommentList
func (vl *VideosContributionList) GetDiscussVideoCommentList(id uint) error {
	return global.Db.Where("uid", id).Preload("Comments").Find(&vl).Error
}

func (vl *VideosContributionList) GetVideoManagementList(info common.PageInfo, uid uint) error {
	return global.Db.Where("uid", uid).Preload("Likes").Preload("Comments").Preload("Barrage").Limit(info.Size).Offset((info.Page - 1) * info.Size).Order("created_at desc").Find(&vl).Error
}
func (vl *VideosContributionList) GetHoneVideoList(info common.PageInfo) error {
	var offset int
	if info.Page == 1 {
		info.Size = 11
		offset = (info.Page - 1) * info.Size
	}
	offset = (info.Page-2)*info.Size + 11

	return global.Db.Preload("Likes").Preload("Comments").Preload("Barrage").Preload("UserInfo").Limit(info.Size).Offset(offset).Order("created_at desc").Find(&vl).Error
}

// GetRecommendList
func (vl *VideosContributionList) GetRecommendList() error {
	return global.Db.Preload("Likes").Preload("Comments").Preload("Barrage").Preload("UserInfo").Order("created_at desc").Limit(7).Find(&vl).Error
}

func (vl *VideosContributionList) Search(info common.PageInfo) error {
	return global.Db.Where("`title` LIKE ?", "%"+info.Keyword+"%").Preload("Likes").Preload("Comments").Preload("Barrage").Preload("UserInfo").Limit(info.Size).Offset((info.Page - 1) * info.Size).Order("created_at desc").Find(&vl).Error

}
