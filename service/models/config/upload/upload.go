package upload

import (
	"simple-video-net/global"
	"simple-video-net/models/common"
)

type Upload struct {
	common.PublicModel
	Interfaces string  `json:"interface"  gorm:"column:interface"`
	Method     string  `json:"method"  gorm:"column:method"`
	Path       string  `json:"path" gorm:"column:path"`
	Quality    float64 `json:"quality"  gorm:"column:quality"`
}

func (Upload) TableName() string {
	return "lv_Upload_method"
}

// IsExistByField Determine if a user exists based on a field
func (um *Upload) IsExistByField(field string, value any) bool {
	err := global.Db.Where(field, value).Find(&um).Error
	if err != nil {
		return false
	}
	if um.ID <= 0 {
		return false
	}
	return true
}
