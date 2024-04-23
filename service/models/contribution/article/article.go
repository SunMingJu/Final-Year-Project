package article

import (
	"simple-video-net/global"
	"simple-video-net/models/common"
	"simple-video-net/models/contribution/article/classification"
	"simple-video-net/models/contribution/article/comments"
	"simple-video-net/models/contribution/article/like"
	"simple-video-net/models/users"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ArticlesContribution struct {
	common.PublicModel
	Uid                uint           `json:"uid" gorm:"column:uid"`
	ClassificationID   uint           `json:"classification_id"  gorm:"column:classification_id"`
	Title              string         `json:"title" gorm:"column:title"`
	Cover              datatypes.JSON `json:"cover" gorm:"column:cover"`
	Timing             int8           `json:"timing" gorm:"column:timing"`
	Content            string         `json:"content" gorm:"column:content"`
	ContentStorageType string         `json:"content_Storage_Type" gorm:"column:content_storage_type"`
	IsComments         int8           `json:"is_comments" gorm:"column:is_comments"`
	Heat               int            `json:"heat" gorm:"column:heat"`

	//optical union table

	UserInfo       users.User                    `json:"user_info" gorm:"foreignKey:Uid"`
	Likes          like.LikesList                `json:"likes" gorm:"foreignKey:ArticleID" `
	Comments       comments.CommentList          `json:"comments" gorm:"foreignKey:ArticleID"`
	Classification classification.Classification `json:"classification"  gorm:"foreignKey:ClassificationID"`
}

type ArticlesContributionList []ArticlesContribution

func (ArticlesContribution) TableName() string {
	return "lv_article_contribution"
}

// Create Add Data
func (vc *ArticlesContribution) Create() bool {
	err := global.Db.Create(&vc).Error
	if err != nil {
		return false
	}
	return true
}

//Watch add play
func (vc *ArticlesContribution) Watch(id uint) error {
	return global.Db.Model(vc).Where("id", id).Updates(map[string]interface{}{"heat": gorm.Expr("Heat  + ?", 1)}).Error
}

// GetList Query Data Type
func (l *ArticlesContributionList) GetList(info common.PageInfo) bool {
	err := global.Db.Preload("Likes").Preload("Classification").Preload("UserInfo").Preload("Comments").Limit(info.Size).Offset((info.Page - 1) * info.Size).Order("created_at desc").Find(l).Error
	if err != nil {
		return false
	}
	return true
}

func (vc *ArticlesContribution) Update(info map[string]interface{}) bool {
	err := global.Db.Model(vc).Updates(info).Error
	if err != nil {
		return false
	}
	return true
}

// GetListByUid 查single user
func (l *ArticlesContributionList) GetListByUid(uid uint) bool {
	err := global.Db.Where("uid", uid).Preload("Likes").Preload("Classification").Preload("Comments").Order("created_at desc").Find(l).Error
	if err != nil {
		return false
	}
	return true
}

// GetAllCount Query the number of all articles
func (l *ArticlesContributionList) GetAllCount(cu *int64) bool {
	err := global.Db.Find(l).Count(cu).Error
	if err != nil {
		return false
	}
	return true
}

// GetInfoByID Search for individual articles
func (vc *ArticlesContribution) GetInfoByID(ID uint) bool {
	err := global.Db.Where("id", ID).Preload("Likes").Preload("UserInfo").Preload("Classification").Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Preload("UserInfo").Order("created_at desc")
	}).Find(vc).Error
	if err != nil {
		return false
	}
	return true
}

// GetArticleComments Get Comments
func (vc *ArticlesContribution) GetArticleComments(ID uint, info common.PageInfo) bool {
	err := global.Db.Where("id", ID).Preload("Likes").Preload("Classification").Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Preload("UserInfo").Order("created_at desc").Limit(info.Size).Offset((info.Page - 1) * info.Size)
	}).Find(vc).Error
	if err != nil {
		return false
	}
	return true
}

// GetArticleBySpace
func (l *ArticlesContributionList) GetArticleBySpace(id uint) error {
	err := global.Db.Where("uid", id).Preload("Likes").Preload("Classification").Preload("Comments").Order("created_at desc").Find(l).Error
	if err != nil {
		return err
	}
	return nil
}

// GetArticleManagementList
func (l *ArticlesContributionList) GetArticleManagementList(info common.PageInfo, id uint) error {
	err := global.Db.Where("uid", id).Preload("Likes").Preload("Classification").Preload("Comments").Limit(info.Size).Offset((info.Page - 1) * info.Size).Order("created_at desc").Find(l).Error
	if err != nil {
		return err
	}
	return nil
}

func (vc *ArticlesContribution) Delete(id uint, uid uint) bool {
	err := global.Db.Where("id", id).Find(vc).Error
	if err != nil {
		return false
	}
	if vc.Uid != uid {
		return false
	}
	err = global.Db.Delete(vc).Error
	if err != nil {
		return false
	}
	return true
}

// GetDiscussArticleCommentList
func (l *ArticlesContributionList) GetDiscussArticleCommentList(id uint) error {
	err := global.Db.Where("uid", id).Preload("Comments").Find(&l).Error
	return err
}
