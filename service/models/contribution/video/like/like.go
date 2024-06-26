package like

import (
	"simple-video-net/global"
	"simple-video-net/models/common"
	"simple-video-net/models/users/notice"

	"gorm.io/gorm"
)

type Likes struct {
	common.PublicModel
	Uid     uint `json:"uid" gorm:"column:uid"`
	VideoID uint `json:"video_id"  gorm:"column:video_id"`
}

type LikesList []Likes

func (Likes) TableName() string {
	return "lv_video_contribution_like"
}

func (l *Likes) IsLike(uid uint, videoID uint) bool {
	err := global.Db.Where(Likes{Uid: uid, VideoID: videoID}).Find(l).Error
	if err != nil {
		return false
	}
	if l.ID <= 0 {
		return false
	}
	return true
}
func (l *Likes) Like(uid uint, videoID uint, videoUid uint) error {
	err := global.Db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("uid", uid).Where("video_id", videoID).Find(l).Error
		if err != nil {
			return err
		}
		if l.ID > 0 {
			err = tx.Where("uid", uid).Where("video_id", videoID).Delete(l).Error
			if err != nil {
				return err
			}
			//No notification for liking your own work
			if videoUid == uid {
				return nil
			}
			//Delete Message Notification
			ne := new(notice.Notice)
			err = ne.Delete(videoUid, uid, videoID, notice.VideoLike)
			if err != nil {
				return err
			}
		} else {
			l.Uid = uid
			l.VideoID = videoID
			err = global.Db.Create(l).Error
			if err != nil {
				return err
			}
			//No notification for liking your own work
			if videoUid == uid {
				return nil
			}
			//Add Message Notification
			ne := new(notice.Notice)
			err = ne.AddNotice(videoUid, uid, videoID, notice.VideoLike, "赞了您的作品")
			if err != nil {
				return err
			}
		}
		// Return nil Commit transaction
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
