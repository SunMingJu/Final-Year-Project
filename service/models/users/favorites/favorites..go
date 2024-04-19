package favorites

import (
	"fmt"
	"simple-video-net/global"
	"simple-video-net/models/common"
	"simple-video-net/models/users"
	"simple-video-net/models/users/collect"

	"gorm.io/datatypes"
)

type Favorites struct {
	common.PublicModel
	Uid     uint           `json:"uid" gorm:"column:uid"`
	Title   string         `json:"title" gorm:"column:title"`
	Content string         `json:"content" gorm:"column:content"`
	Cover   datatypes.JSON `json:"cover" gorm:"type:json;comment:cover"`
	Max     int            `json:"max" gorm:"column:max"`

	UserInfo    users.User           `json:"userInfo" gorm:"foreignKey:Uid"`
	CollectList collect.CollectsList `json:"collectList"  gorm:"foreignKey:FavoritesID"`
}

type FavoriteList []Favorites

func (Favorites) TableName() string {
	return "lv_users_favorites"
}

// Find
func (f *Favorites) Find(id uint) bool {
	err := global.Db.Where("id", id).Preload("CollectList").Order("created_at desc").Find(&f).Error
	if err != nil {
		return false
	}
	return true
}

// Create
func (f *Favorites) Create() bool {
	err := global.Db.Create(&f).Error
	if err != nil {
		return false
	}
	return true
}

// AloneTitleCreate
func (f *Favorites) AloneTitleCreate() bool {
	err := global.Db.Create(&f).Error
	if err != nil {
		return false
	}
	return true
}

// Update
func (f *Favorites) Update() bool {
	err := global.Db.Updates(&f).Error
	if err != nil {
		return false
	}
	return true
}

// Delete
func (f *Favorites) Delete(id uint, uid uint) error {
	err := global.Db.Where("id", id).Find(&f).Error
	if err != nil {
		return fmt.Errorf("Enquiry Failure")
	}
	if f.ID <= 0 {
		return fmt.Errorf("Favourites do not exist")
	}
	err = global.Db.Delete(&f).Error
	if f.Uid != uid {
		return fmt.Errorf("Cannot be deleted by non-creators")
	}
	//Delete Collection Record
	cl := new(collect.Collect)
	if !cl.DetectByFavoritesID(id) {
		return fmt.Errorf("Failed to delete favourite records")
	}
	if err != nil {
		return fmt.Errorf("Failed to delete")
	}
	return nil
}

func (fl *FavoriteList) GetFavoritesList(id uint) error {
	err := global.Db.Where("uid", id).Preload("UserInfo").Preload("CollectList").Order("created_at desc").Find(fl).Error
	if err != nil {
		return err
	}
	return nil
}
