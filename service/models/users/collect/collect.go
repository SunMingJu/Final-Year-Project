package collect

import (
	"easy-video-net/global"
	"easy-video-net/models/common"
	"easy-video-net/models/contribution/video"
	"easy-video-net/models/users"
)

type Collect struct {
	common.PublicModel
	Uid         uint `json:"uid" gorm:"uid"`
	FavoritesID uint `json:"favorites_id" gorm:"favorites_id"`
	VideoID     uint `json:"video_id" gorm:"video_id "`

	UserInfo  users.User               `json:"userInfo" gorm:"foreignKey:Uid"`
	VideoInfo video.VideosContribution `json:"videoInfo" gorm:"foreignKey:VideoID"`
}

type CollectsList []Collect

func (Collect) TableName() string {
	return "lv_users_collect"
}

//Create 
func (c *Collect) Create() bool {
	err := global.Db.Create(c).Error
	if err != nil {
		return false
	}
	return true
}

//DetectByFavoritesID 
func (c *Collect) DetectByFavoritesID(id uint) bool {
	err := global.Db.Where("favorites_id", id).Delete(c).Error
	if err != nil {
		return false
	}
	return true
}

func (l *CollectsList) FindVideoExistWhere(videoID uint) error {
	err := global.Db.Where("video_id", videoID).Find(l).Error
	return err
}

func (l *CollectsList) GetVideoInfo(id uint) error {
	err := global.Db.Where("favorites_id", id).Preload("VideoInfo").Find(l).Error
	return err
}

//FindIsCollectByFavorites 
func (l *CollectsList) FindIsCollectByFavorites(videoID uint, ids []uint) bool {
	//False if no favourite is created
	if len(ids) == 0 {
		return false
	}
	err := global.Db.Where("video_id", videoID).Where("favorites_id", ids).Find(l).Error
	if err != nil {
		return false
	}
	if len(*l) == 0 {
		return false
	}
	return true
}
