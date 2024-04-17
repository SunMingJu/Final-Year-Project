package users

import (
	"crypto/md5"
	"easy-video-net/consts"
	"easy-video-net/global"
	receive "easy-video-net/interaction/receive/users"
	response "easy-video-net/interaction/response/users"
	"easy-video-net/models/common"
	"easy-video-net/models/users"
	"easy-video-net/models/users/attention"
	"easy-video-net/models/users/chat/chatList"
	"easy-video-net/models/users/chat/chatMsg"
	"easy-video-net/models/users/collect"
	"easy-video-net/models/users/favorites"
	"easy-video-net/models/users/liveInfo"
	"easy-video-net/models/users/notice"
	"easy-video-net/models/users/record"
	"easy-video-net/utils/conversion"
	"easy-video-net/utils/email"
	"easy-video-net/utils/jwt"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/datatypes"
)

func GetUserInfo(uid uint) (results interface{}, err error) {
	user := new(users.User)
	user.IsExistByField("id", uid)
	res := response.UserSetInfoResponse(user)
	return res, nil
}

func SetUserInfo(data *receive.SetUserInfoReceiveStruct, uid uint) (results interface{}, err error) {
	user := &users.User{
		PublicModel: common.PublicModel{ID: uid},
	}
	update := map[string]interface{}{
		"Username":  data.Username,
		"Gender":    0,
		"BirthDate": data.BirthDate,
		"IsVisible": conversion.BoolTurnInt8(*data.IsVisible),
		"Signature": data.Signature,
	}

	return user.UpdatePureZero(update), nil
}

func DetermineNameExists(data *receive.DetermineNameExistsStruct, uid uint) (results interface{}, err error) {
	user := new(users.User)
	is := user.IsExistByField("username", data.Username)
	//Determine if no changes have been made
	if user.ID == uid {
		return false, nil
	} else if is {
		return true, nil
	} else {
		return false, nil
	}
}

func UpdateAvatar(data *receive.UpdateAvatarStruct, uid uint) (results interface{}, err error) {
	photo, _ := json.Marshal(common.Img{
		Src: data.ImgUrl,
		Tp:  data.Tp,
	})
	user := &users.User{PublicModel: common.PublicModel{ID: uid}, Photo: photo}
	if user.Update() {
		return conversion.FormattingSrc(data.ImgUrl), nil
	} else {
		return nil, fmt.Errorf("update failure")
	}
}

func GetLiveData(uid uint) (results interface{}, err error) {
	info := new(liveInfo.LiveInfo)
	if info.IsExistByField("uid", uid) {
		results, err = response.GetLiveDataResponse(info)
		if err != nil {
			return nil, fmt.Errorf("Failed to get")
		}
		return results, nil
	}
	return common.Img{}, nil
}

func SaveLiveData(data *receive.SaveLiveDataReceiveStruct, uid uint) (results interface{}, err error) {
	img, _ := json.Marshal(common.Img{
		Src: data.ImgUrl,
		Tp:  data.Tp,
	})
	info := &liveInfo.LiveInfo{
		Uid:   uid,
		Title: data.Title,
		Img:   datatypes.JSON(img),
	}
	if info.UpdateInfo() {
		return "Modified successfully", nil
	} else {
		return nil, fmt.Errorf("Modification Failure")
	}

}

func SendEmailVerificationCodeByChangePassword(uid uint) (results interface{}, err error) {
	user := new(users.User)
	user.Find(uid)
	//sender
	mailTo := []string{user.Email}
	// Email Subject
	code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(1000000))
	subject := "CAPTCHA"
	// Body of the email
	body := fmt.Sprintf("You are changing your password, your verification code is.%s,5 minutes. Please do not forward to others.", code)
	err = email.SendMail(mailTo, subject, body)
	if err != nil {
		return nil, err
	}
	err = global.RedisDb.Set(fmt.Sprintf("%s%s", consts.EmailVerificationCodeByChangePassword, user.Email), code, 5*time.Minute).Err()
	if err != nil {
		return nil, err
	}
	return "Sent successfully", nil

}

func ChangePassword(data *receive.ChangePasswordReceiveStruct, uid uint) (results interface{}, err error) {
	user := new(users.User)
	user.Find(uid)

	if data.Password != data.ConfirmPassword {
		return nil, fmt.Errorf("The two passwords do not match!")
	}

	//Determine if the CAPTCHA is correct
	verCode, err := global.RedisDb.Get(fmt.Sprintf("%s%s", consts.EmailVerificationCodeByChangePassword, user.Email)).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("Captcha expired!")
	}

	if verCode != data.VerificationCode {
		return nil, fmt.Errorf("CAPTCHA error")
	}
	//Generate password salt 8 bits
	salt := make([]byte, 6)
	for i := range salt {
		salt[i] = jwt.SaltStr[rand.Int63()%int64(len(jwt.SaltStr))]
	}
	password := []byte(fmt.Sprintf("%s%s%s", salt, data.Password, salt))
	passwordMd5 := fmt.Sprintf("%x", md5.Sum(password))

	user.Salt = string(salt)
	user.Password = passwordMd5

	registerRes := user.Update()
	if !registerRes {
		return nil, fmt.Errorf("Modification Failure")
	}
	return "Modified successfully", nil
}

func Attention(data *receive.AttentionReceiveStruct, uid uint) (results interface{}, err error) {
	at := new(attention.Attention)
	if at.Attention(uid, data.Uid) {
		if data.Uid == uid {
			return nil, fmt.Errorf("failure of an operation")
		}
		return "The operation was successful.", nil
	}
	return nil, fmt.Errorf("failure of an operation")
}

func CreateFavorites(data *receive.CreateFavoritesReceiveStruct, uid uint) (results interface{}, err error) {
	if data.ID == 0 {
		//Insertion modes
		if len(data.Title) == 0 {
			return nil, fmt.Errorf("Title is empty")
		}
		//Determine if there is only a title
		if data.ID <= 0 && len(data.Tp) == 0 && len(data.Content) == 0 && len(data.Cover) == 0 {
			//Single Title Creation
			fs := &favorites.Favorites{Uid: uid, Title: data.Title, Max: 1000}
			if !fs.Create() {
				return nil, fmt.Errorf("Creation Failure")
			}
			return fmt.Errorf("Created Successfully"), nil
		} else {
			//Well-documented creation
			cover, _ := json.Marshal(common.Img{
				Src: data.Cover,
				Tp:  data.Tp,
			})
			fs := &favorites.Favorites{
				Uid:     uid,
				Title:   data.Title,
				Content: data.Content,
				Cover:   cover,
				Max:     1000,
			}
			if !fs.Create() {
				return nil, fmt.Errorf("Creation Failure")
			}
			return fmt.Errorf("Created Successfully"), nil
		}
	} else {
		//carry out an update
		fs := new(favorites.Favorites)
		if !fs.Find(data.ID) {
			return nil, fmt.Errorf("Enquiry Failure")
		}
		if fs.Uid != uid {
			return nil, fmt.Errorf("Querying illegal operations")
		}
		cover, _ := json.Marshal(common.Img{
			Src: data.Cover,
			Tp:  data.Tp,
		})
		fs.Title = data.Title
		fs.Content = data.Content
		fs.Cover = cover
		if !fs.Update() {
			return nil, fmt.Errorf("update failure")
		}
		return "Successful update", nil
	}
}

func GetFavoritesList(uid uint) (results interface{}, err error) {
	fl := new(favorites.FavoriteList)
	err = fl.GetFavoritesList(uid)
	if err != nil {
		return nil, fmt.Errorf("Enquiry Failure")
	}
	res, err := response.GetFavoritesListResponse(fl)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func DeleteFavorites(data *receive.DeleteFavoritesReceiveStruct, uid uint) (results interface{}, err error) {
	fs := new(favorites.Favorites)
	err = fs.Delete(data.ID, uid)
	if err != nil {
		return nil, err
	}
	return "Deleted successfully", nil
}

func FavoriteVideo(data *receive.FavoriteVideoReceiveStruct, uid uint) (results interface{}, err error) {
	for _, k := range data.IDs {
		fs := new(favorites.Favorites)
		fs.Find(k)
		if fs.Uid != uid {
			return nil, fmt.Errorf("unauthorised operation")
		}
		if len(fs.CollectList)+1 > fs.Max {
			return nil, fmt.Errorf("Favourites are full")
		}

		cl := &collect.Collect{
			Uid:         uid,
			FavoritesID: k,
			VideoID:     data.VideoID,
		}
		if !cl.Create() {
			return nil, fmt.Errorf("Collection Failure")
		}
	}
	return "The operation was successful.", nil
}

func GetFavoritesListByFavoriteVideo(data *receive.GetFavoritesListByFavoriteVideoReceiveStruct, uid uint) (results interface{}, err error) {
	//Get favourites list
	fl := new(favorites.FavoriteList)
	err = fl.GetFavoritesList(uid)
	if err != nil {
		return nil, fmt.Errorf("Enquiry Failure")
	}
	//Check which videos have been bookmarked in those favourites.
	cl := new(collect.CollectsList)
	err = cl.FindVideoExistWhere(data.VideoID)
	if err != nil {
		return nil, fmt.Errorf("Enquiry Failure")
	}
	ids := make([]uint, 0)
	for _, v := range *cl {
		ids = append(ids, v.FavoritesID)
	}

	res, err := response.GetFavoritesListByFavoriteVideoResponse(fl, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func GetFavoriteVideoList(data *receive.GetFavoriteVideoListReceiveStruct) (results interface{}, err error) {
	cl := new(collect.CollectsList)
	err = cl.GetVideoInfo(data.FavoriteID)
	if err != nil {
		return nil, fmt.Errorf("Enquiry Failure")
	}
	res, err := response.GetFavoriteVideoListResponse(cl)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func GetRecordList(data *receive.GetRecordListReceiveStruct, uid uint) (results interface{}, err error) {
	rl := new(record.RecordsList)
	err = rl.GetRecordListByUid(uid, data.PageInfo)
	if err != nil {
		return nil, fmt.Errorf("Enquiry Failure")
	}
	res, err := response.GetRecordListResponse(rl)
	if err != nil {
		return nil, fmt.Errorf("Response Failure")
	}
	return res, nil
}

func ClearRecord(uid uint) (results interface{}, err error) {
	rl := new(record.Record)
	err = rl.ClearRecord(uid)
	if err != nil {
		return nil, fmt.Errorf("Emptying failed")
	}
	return "Emptying complete", nil
}

func DeleteRecordByID(data *receive.DeleteRecordByIDReceiveStruct, uid uint) (results interface{}, err error) {
	rl := new(record.Record)
	err = rl.DeleteRecordByID(data.ID, uid)
	if err != nil {
		return nil, fmt.Errorf("Failed to delete")
	}
	return "Deleted successfully", nil
}

func GetNoticeList(data *receive.GetNoticeListReceiveStruct, uid uint) (results interface{}, err error) {
	//Get user notifications
	messageType := make([]string, 0)
	nl := new(notice.NoticesList)
	switch data.Type {
	case "comment":
		messageType = append(messageType, notice.VideoComment, notice.ArticleComment)
		break
	case "like":
		messageType = append(messageType, notice.VideoLike, notice.ArticleLike)
	}

	err = nl.GetNoticeList(data.PageInfo, messageType, uid)
	if err != nil {
		return nil, fmt.Errorf("Enquiry Failure")
	}
	//Record all read
	n := new(notice.Notice)
	err = n.ReadAll(uid)
	if err != nil {
		return nil, fmt.Errorf("Read Message Failure")
	}
	res, err := response.GetNoticeListResponse(nl)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func GetChatList(uid uint) (results interface{}, err error) {
	//Getting a list of messages
	cList := new(chatList.ChatList)
	err = cList.GetListByIO(uid)
	if err != nil {
		return nil, fmt.Errorf("Enquiry Failure")
	}
	ids := make([]uint, 0)
	for _, v := range *cList {
		ids = append(ids, v.Tid)
	}
	msgList := make(map[uint]*chatMsg.MsgList, 0)
	for _, v := range ids {
		ml := new(chatMsg.MsgList)
		err = ml.FindList(uid, v)
		if err != nil {
			break
		}
		msgList[v] = ml
	}
	res, err := response.GetChatListResponse(cList, msgList)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func GetChatHistoryMsg(data *receive.GetChatHistoryMsgStruct, uid uint) (results interface{}, err error) {
	//Query History Messages
	cm := new(chatMsg.MsgList)
	err = cm.FindHistoryMsg(uid, data.Tid, data.LastTime)
	if err != nil {
		return nil, fmt.Errorf("Enquiry Failure")
	}
	res, err := response.GetChatHistoryMsgResponse(cm)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func PersonalLetter(data *receive.PersonalLetterReceiveStruct, uid uint) (results interface{}, err error) {
	cm := new(chatMsg.Msg)
	err = cm.GetLastMessage(uid, data.ID)
	if err != nil {
		return nil, fmt.Errorf("failure of an operation")
	}
	var lastTime time.Time
	if cm.ID > 0 {
		lastTime = cm.CreatedAt
	} else {
		lastTime = time.Now()
	}
	ci := &chatList.ChatsListInfo{
		Uid:         uid,
		Tid:         data.ID,
		LastMessage: cm.Message,
		LastAt:      lastTime,
	}
	err = ci.AddChat()
	if err != nil {
		return nil, fmt.Errorf("failure of an operation")
	}
	return "The operation was successful.", nil
}

func DeleteChatItem(data *receive.DeleteChatItemReceiveStruct, uid uint) (results interface{}, err error) {
	ci := new(chatList.ChatsListInfo)
	err = ci.DeleteChat(data.ID, uid)
	if err != nil {
		return nil, fmt.Errorf("Failed to delete")
	}
	return "The operation was successful.", nil
}
