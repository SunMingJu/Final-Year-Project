package uploadMethod

import (
	"simple-video-net/global"
	"simple-video-net/models/common"
)

type UploadMethod struct {
	common.PublicModel
	Interfaces string `json:"interface"  gorm:"interface"`
	Method     string `json:"method"  gorm:"method"`
	Path       string `json:"path" gorm:"path"`
}

func (UploadMethod) TableName() string {
	return "lv_upload_method"
}

// IsExistByField Determine if a user exists based on a field
func (um *UploadMethod) IsExistByField(field string, value any) bool {
	err := global.Db.Where(field, value).Find(&um).Error
	if err != nil {
		return false
	}
	if um.ID <= 0 {
		return false
	}
	return true
}
