package liveInfo

import (
	"simple-video-net/global"
	"simple-video-net/models/common"

	"gorm.io/datatypes"
)

type LiveInfo struct {
	common.PublicModel
	Uid   uint           `json:"uid" gorm:"column:uid"`
	Title string         `json:"title" gorm:"column:title"`
	Img   datatypes.JSON `json:"img" gorm:"type:json;comment:img"`
}

func (LiveInfo) TableName() string {
	return "lv_live_info"
}

// Update
func (li *LiveInfo) Update() bool {
	err := global.Db.Where("uid", li.Uid).Updates(&li).Error
	if err != nil {
		return false
	}
	return true
}

// Create
func (li *LiveInfo) Create() bool {
	err := global.Db.Create(&li).Error
	if err != nil {
		return false
	}
	return true
}

func (li *LiveInfo) UpdateInfo() bool {
	l := new(LiveInfo)
	if l.IsExistByField("uid", li.Uid) {
		return li.Update()
	} else {
		return li.Create()
	}
}

// IsExistByField
func (li *LiveInfo) IsExistByField(field string, value any) bool {
	err := global.Db.Where(field, value).Find(&li).Error
	if err != nil {
		return false
	}
	if li.ID <= 0 {
		return false
	}
	return true
}
