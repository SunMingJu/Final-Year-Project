package barrage

import (
	"simple-video-net/global"
	"simple-video-net/models/common"
	"simple-video-net/models/users"

	"gorm.io/datatypes"
)

type Barrage struct {
	common.PublicModel
	Uid     uint    `json:"uid" gorm:"uid"`
	VideoID uint    `json:"video_id" gorm:"video_id"`
	Time    float64 `json:"time" gorm:"time"`
	Author  string  `json:"author" gorm:"author"`
	Type    uint    `json:"type" gorm:"type"`
	Text    string  `json:"text" gorm:"text"`
	Color   uint    `json:"color" gorm:"color"`

	UserInfo  users.User `json:"user_info" gorm:"foreignKey:Uid"`
	VideoInfo VideoInfo  `json:"video_info" gorm:"foreignKey:VideoID"`
}

type BarragesList []Barrage

func (Barrage) TableName() string {
	return "lv_video_contribution_barrage"
}

// VideoInfo Temporarily add a video model to solve the dependency loop
type VideoInfo struct {
	common.PublicModel
	Uid   uint           `json:"uid" gorm:"uid"`
	Title string         `json:"title" gorm:"title"`
	Video datatypes.JSON `json:"video" gorm:"video"`
	Cover datatypes.JSON `json:"cover" gorm:"cover"`
}

func (VideoInfo) TableName() string {
	return "lv_video_contribution"
}

func (b *Barrage) Create() bool {
	err := global.Db.Create(&b).Error
	if err != nil {
		return false
	}
	return true
}

// GetVideoBarrageByID Enquire about video pop-ups
func (bl *BarragesList) GetVideoBarrageByID(id uint) bool {
	err := global.Db.Where("video_id", id).Find(&bl).Error
	if err != nil {
		return false
	}
	return true
}

func (bl *BarragesList) GetBarrageListByIDs(ids []uint, info common.PageInfo) error {
	return global.Db.Where("video_id", ids).Preload("UserInfo").Preload("VideoInfo").Limit(info.Size).Offset((info.Page - 1) * info.Size).Order("created_at desc").Find(&bl).Error
}
