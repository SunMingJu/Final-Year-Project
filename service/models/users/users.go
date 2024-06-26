package users

import (
	"crypto/md5"
	"fmt"
	"simple-video-net/global"
	"simple-video-net/models/common"
	"simple-video-net/models/users/liveInfo"
	"time"

	"gorm.io/datatypes"
)

// User 表结构体
type User struct {
	common.PublicModel
	Email     string         `json:"email" gorm:"column:email"`
	Username  string         `json:"username" gorm:"column:username"`
	Openid    string         `json:"openid" gorm:"column:openid"`
	Salt      string         `json:"salt" gorm:"column:salt"`
	Password  string         `json:"password" gorm:"column:password"`
	Photo     datatypes.JSON `json:"photo" gorm:"column:photo"`
	Gender    int8           `json:"gender" gorm:"column:gender"`
	BirthDate time.Time      `json:"birth_date" gorm:"column:birth_date"`
	IsVisible int8           `json:"is_visible" gorm:"column:is_visible"`
	Signature string         `json:"signature" gorm:"column:signature"`

	LiveInfo liveInfo.LiveInfo `json:"liveInfo" gorm:"foreignKey:Uid"`
}

type UserList []User

func (User) TableName() string {
	return "lv_users"
}

// Update
func (us *User) Update() bool {
	err := global.Db.Where("id", us.ID).Updates(&us).Error
	if err != nil {
		return false
	}
	return true
}

// UpdatePureZero
func (us *User) UpdatePureZero(user map[string]interface{}) bool {
	err := global.Db.Model(&us).Where("id", us.ID).Updates(user).Error
	if err != nil {
		return false
	}
	return true
}

// Create
func (us *User) Create() bool {
	err := global.Db.Create(&us).Error
	if err != nil {
		return false
	}
	return true
}

// IsExistByField
func (us *User) IsExistByField(field string, value any) bool {
	err := global.Db.Where(field, value).Find(&us).Error
	if err != nil {
		return false
	}
	if us.ID <= 0 {
		return false
	}
	return true
}

// IfPasswordCorrect
func (us *User) IfPasswordCorrect(password string) bool {
	passwordImport := fmt.Sprintf("%s%s%s", us.Salt, password, us.Salt)
	passwordImport = fmt.Sprintf("%x", md5.Sum([]byte(passwordImport)))
	if passwordImport != us.Password {
		return false
	}
	return true
}

// Find
func (us *User) Find(id uint) {
	_ = global.Db.Where("id", id).Find(&us).Error
}

// FindLiveInfo
func (us *User) FindLiveInfo(id uint) {
	_ = global.Db.Where("id", id).Preload("LiveInfo").Find(&us).Error
}

func (l *UserList) GetBeLiveList(ids []uint) error {
	return global.Db.Where("id", ids).Preload("LiveInfo").Find(&l).Error
}

func (l *UserList) Search(info common.PageInfo) error {
	return global.Db.Where("`username` LIKE ?", "%"+info.Keyword+"%").Find(&l).Error
}
