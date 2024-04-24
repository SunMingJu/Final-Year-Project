package chatList

import (
	"fmt"
	"simple-video-net/global"
	"simple-video-net/models/common"
	"simple-video-net/models/users"
	"time"
)

type ChatsListInfo struct {
	common.PublicModel
	Uid         uint      `json:"uid" gorm:"column:uid"`
	Tid         uint      `json:"tid"  gorm:"column:tid"`
	Unread      int       `json:"unread" gorm:"column:unread"`
	LastMessage string    `json:"last_message" gorm:"column:last_message"`
	LastAt      time.Time `json:"last_at" gorm:"column:last_at"`

	ToUserInfo users.User `json:"toUserInfo"  gorm:"foreignKey:tid"`
}

type ChatList []ChatsListInfo

func (ChatsListInfo) TableName() string {
	return "lv_users_chat_list"
}

// AddChat
func (i *ChatsListInfo) AddChat() error {
	//Determine the presence or absence of
	is := &ChatsListInfo{}
	err := global.Db.Where("uid = ? And tid = ?", i.Uid, i.Tid).Find(is).Error
	if err != nil {
		return err
	}
	if is.ID != 0 {
		//Presence is renewal
		global.Db.Model(is).Updates(map[string]interface{}{"last_at": i.LastAt, "last_message": i.LastMessage})
		return nil
	} else {
		return global.Db.Create(i).Error
	}
}

// DeleteChat
func (i *ChatsListInfo) DeleteChat(tid uint, uid uint) error {
	return global.Db.Where("uid = ? and tid = ?", uid, tid).Delete(i).Error
}

// GetListByIO
func (cl *ChatList) GetListByIO(uid uint) error {
	return global.Db.Where("uid", uid).Preload("ToUserInfo").Order("updated_at desc").Find(cl).Error
}

func (i *ChatsListInfo) UnreadAutocorrection(uid uint, tid uint) error {
	err := global.Db.Where(ChatsListInfo{Uid: uid, Tid: tid}).Find(i).Error
	if err != nil {
		return err
	}
	if i.ID > 0 {
		i.Unread++
		return global.Db.Save(i).Error
	}
	return fmt.Errorf("No record")
}

func (i *ChatsListInfo) UnreadEmpty(uid uint, tid uint) error {
	err := global.Db.Where(ChatsListInfo{Uid: uid, Tid: tid}).Find(i).Error
	if err != nil {
		return err
	}
	if i.ID > 0 {
		i.Unread = 0
		return global.Db.Save(i).Error
	}
	return fmt.Errorf("case failure")
}

func (i *ChatsListInfo) FindByID(uid uint, tid uint) error {
	return global.Db.Where(ChatsListInfo{Uid: uid, Tid: tid}).Find(i).Error
}

// GetUnreadNumber
func (i ChatsListInfo) GetUnreadNumber(uid uint) *int64 {
	num := new(int64)
	err := global.Db.Model(i).Select("IFNULL(unread,0) as total_unread").Where(ChatsListInfo{Uid: uid}).Scan(num).Error
	if err != nil {
		fmt.Println(err)
	}

	return num
}
