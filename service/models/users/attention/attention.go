package attention

import (
	"easy-video-net/global"
	"easy-video-net/models/common"
	"easy-video-net/models/users"
)

type Attention struct {
	common.PublicModel
	Uid         uint `json:"uid" gorm:"uid"`
	AttentionID uint `json:"attention_id" gorm:"attention_id"`

	UserInfo          users.User `json:"user_info" gorm:"foreignKey:Uid"`
	AttentionUserInfo users.User `json:"attention_user_info" gorm:"foreignKey:AttentionID"`
}

type AttentionsList []Attention

func (Attention) TableName() string {
	return "lv_users_attention"
}

//Create 
func (at *Attention) Create() bool {
	err := global.Db.Create(&at).Error
	if err != nil {
		return false
	}
	return true
}

//Delete 
func (at *Attention) Delete() bool {
	err := global.Db.Where("uid", at.Uid).Updates(&at).Error
	if err != nil {
		return false
	}
	return true
}

//Attention 
func (at *Attention) Attention(uid uint, aid uint) bool {
	err := global.Db.Where(Attention{Uid: uid, AttentionID: aid}).Find(&at).Error
	if at.ID > 0 {
		//Followed
		err = global.Db.Where("id", at.ID).Delete(&at).Error
	} else {
		//not following
		err = global.Db.Create(&Attention{Uid: uid, AttentionID: aid}).Error
	}

	if err != nil {
		return false
	}
	return true
}

//IsAttention  
func (at *Attention) IsAttention(uid uint, aid uint) bool {
	_ = global.Db.Where(Attention{Uid: uid, AttentionID: aid}).Find(&at).Error
	if at.ID > 0 {
		return true
	} else {
		return false
	}
}

//GetAttentionNum 
func (at *Attention) GetAttentionNum(uid uint) (*int64, error) {
	num := new(int64)
	err := global.Db.Model(at).Where(Attention{Uid: uid}).Count(num).Error
	if err != nil {
		return num, err
	}
	return num, nil
}

//GetVermicelliNum 
func (at *Attention) GetVermicelliNum(uid uint) (*int64, error) {
	num := new(int64)
	err := global.Db.Model(at).Where(Attention{AttentionID: uid}).Count(num).Error
	if err != nil {
		return num, err
	}
	return num, nil
}

//GetAttentionList 
func (al *AttentionsList) GetAttentionList(uid uint) error {
	err := global.Db.Where("uid", uid).Preload("AttentionUserInfo").Find(al).Error
	if err != nil {
		return err
	}
	return nil
}

//GetVermicelliList 
func (al *AttentionsList) GetVermicelliList(uid uint) error {
	err := global.Db.Where("attention_id", uid).Preload("UserInfo").Find(al).Error
	if err != nil {
		return err
	}
	return nil
}

//GetAttentionListByIdArr 
func (al *AttentionsList) GetAttentionListByIdArr(uid uint) (arr []uint, err error) {
	arr = make([]uint, 0)
	err = global.Db.Where("uid", uid).Find(al).Error
	if err != nil {
		return arr, err
	}
	//self-requested array
	for _, v := range *al {
		arr = append(arr, v.AttentionID)
	}
	return arr, nil
}
