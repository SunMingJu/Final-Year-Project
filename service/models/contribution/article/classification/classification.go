package classification

import (
	"simple-video-net/global"
	"simple-video-net/models/common"
)

type Classification struct {
	common.PublicModel
	AID   uint   `json:"a_id" gorm:"a_id"`
	Label string `json:"label" gorm:"label"`
}

type ClassificationsList []Classification

func (Classification) TableName() string {
	return "lv_article_classification"
}

func (cl *ClassificationsList) FindAll() error {
	return global.Db.Find(cl).Error
}
