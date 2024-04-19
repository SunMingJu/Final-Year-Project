package contribution

import (
	"encoding/json"
	"fmt"
	"simple-video-net/consts"
	"simple-video-net/global"
	receive "simple-video-net/interaction/receive/contribution/video"
	response "simple-video-net/interaction/response/contribution/video"
	"simple-video-net/logic/contribution/sokcet"
	"simple-video-net/logic/users/notice"
	"simple-video-net/models/common"
	"simple-video-net/models/contribution/video"
	"simple-video-net/models/contribution/video/barrage"
	"simple-video-net/models/contribution/video/comments"
	"simple-video-net/models/contribution/video/like"
	"simple-video-net/models/users/attention"
	"simple-video-net/models/users/collect"
	"simple-video-net/models/users/favorites"
	noticeModel "easy-video-net/models/users/notice"
	"simple-video-net/models/users/record"
	"simple-video-net/utils/conversion"
	"strconv"
)

func CreateVideoContribution(data *receive.CreateVideoContributionReceiveStruct, uid uint) (results interface{}, err error) {
	//publish video
	videoSrc, _ := json.Marshal(common.Img{
		Src: data.Video,
		Tp:  data.VideoUploadType,
	})
	coverImg, _ := json.Marshal(common.Img{
		Src: data.Cover,
		Tp:  data.CoverUploadType,
	})
	videoContribution := video.VideosContribution{
		Uid:           uid,
		Title:         data.Title,
		Video:         videoSrc,
		Cover:         coverImg,
		VideoDuration: data.VideoDuration,
		Reprinted:     conversion.BoolTurnInt8(*data.Reprinted),
		Timing:        conversion.BoolTurnInt8(*data.Timing),
		TimingTime:    data.TimingTime,
		Label:         conversion.MapConversionString(data.Label),
		Introduce:     data.Introduce,
		Heat:          0,
	}
	if *data.Timing {
		//Push related after posting a video (to be developed)
	}
	if !videoContribution.Create() {
		return nil, fmt.Errorf("fail to save")
	}
	return "Save Successful", nil
}

func UpdateVideoContribution(data *receive.UpdateVideoContributionReceiveStruct, uid uint) (results interface{}, err error) {
	//Update Video
	videoInfo := new(video.VideosContribution)
	err = videoInfo.FindByID(data.ID)
	if err != nil {
		return nil, fmt.Errorf("Operation video does not exist")
	}
	if videoInfo.Uid != uid {
		return nil, fmt.Errorf("unauthorised operation")
	}
	coverImg, _ := json.Marshal(common.Img{
		Src: data.Cover,
		Tp:  data.CoverUploadType,
	})
	updateList := map[string]interface{}{
		"cover":     coverImg,
		"title":     data.Title,
		"label":     conversion.MapConversionString(data.Label),
		"reprinted": conversion.BoolTurnInt8(*data.Reprinted),
		"introduce": data.Introduce,
	}
	//进行视频资料更新
	if !videoInfo.Update(updateList) {
		return nil, fmt.Errorf("Failed to update data")
	}
	return "Successful update", nil
}

func DeleteVideoByID(data *receive.DeleteVideoByIDReceiveStruct, uid uint) (results interface{}, err error) {
	vo := new(video.VideosContribution)
	if !vo.Delete(data.ID, uid) {
		return nil, fmt.Errorf("Failed to delete")
	}
	return "Deleted successfully", nil
}

func GetVideoContributionByID(data *receive.GetVideoContributionByIDReceiveStruct, uid uint) (results interface{}, err error) {
	videoInfo := new(video.VideosContribution)
	//Get video information
	err = videoInfo.FindByID(data.VideoID)
	if err != nil {
		return nil, fmt.Errorf("Failed to query information")
	}
	isAttention := false
	isLike := false
	isCollect := false
	if uid != 0 {
		//Perform video playback additions
		if !global.RedisDb.SIsMember(consts.VideoWatchByID+strconv.Itoa(int(data.VideoID)), uid).Val() {
			//No recent broadcasts
			global.RedisDb.SAdd(consts.VideoWatchByID+strconv.Itoa(int(data.VideoID)), uid)
			if videoInfo.Watch(data.VideoID) != nil {
				global.Logger.Error("Add playback error", videoInfo.Watch(data.VideoID))
			}
		}
		//Get attention or not
		at := new(attention.Attention)
		isAttention = at.IsAttention(uid, videoInfo.UserInfo.ID)

		//Get attention or not
		lk := new(like.Likes)
		isLike = lk.IsLike(uid, videoInfo.ID)

		//Determine whether a collection has been made
		fl := new(favorites.FavoriteList)
		err = fl.GetFavoritesList(uid)
		if err != nil {
			return nil, fmt.Errorf("Enquiry Failure")
		}
		flIDs := make([]uint, 0)
		for _, v := range *fl {
			flIDs = append(flIDs, v.ID)
		}
		//Determine if you are in your favourites
		cl := new(collect.CollectsList)
		isCollect = cl.FindIsCollectByFavorites(data.VideoID, flIDs)

		//Add History
		rd := new(record.Record)
		err = rd.AddVideoRecord(uid, data.VideoID)
		if err != nil {
			return nil, fmt.Errorf("Failed to add history")
		}

	}
	//Get Recommended List
	recommendList := new(video.VideosContributionList)
	err = recommendList.GetRecommendList()
	if err != nil {
		return nil, err
	}
	res := response.GetVideoContributionByIDResponse(videoInfo, recommendList, isAttention, isLike, isCollect)
	return res, nil
}

func SendVideoBarrage(data *receive.SendVideoBarrageReceiveStruct, uid uint) (results interface{}, err error) {
	//Save pop-ups
	videoID, _ := strconv.ParseUint(data.ID, 0, 19)
	bg := barrage.Barrage{
		Uid:     uid,
		VideoID: uint(videoID),
		Time:    data.Time,
		Author:  data.Author,
		Type:    data.Type,
		Text:    data.Text,
		Color:   data.Color,
	}
	if !bg.Create() {
		return data, fmt.Errorf("Failed to send pop-up")
	}
	//socket message notification
	res := sokcet.ChanInfo{
		Type: consts.sokcetTypeResponseBarrageNum,
		Data: nil,
	}
	for _, v := range sokcet.Severe.VideoRoom[uint(videoID)] {
		v.MsgList <- res
	}

	return data, nil
}

func GetVideoBarrage(data *receive.GetVideoBarrageReceiveStruct) (results interface{}, err error) {
	//Get video pop-ups
	list := new(barrage.BarragesList)
	videoID, _ := strconv.ParseUint(data.ID, 0, 19)
	if !list.GetVideoBarrageByID(uint(videoID)) {
		return nil, fmt.Errorf("Enquiry Failure")
	}
	res := response.GetVideoBarrageResponse(list)
	return res, nil
}

func GetVideoBarrageList(data *receive.GetVideoBarrageListReceiveStruct) (results interface{}, err error) {
	//Get video pop-ups
	list := new(barrage.BarragesList)
	videoID, _ := strconv.ParseUint(data.ID, 0, 19)
	if !list.GetVideoBarrageByID(uint(videoID)) {
		return nil, fmt.Errorf("Enquiry Failure")
	}
	res := response.GetVideoBarrageListResponse(list)
	return res, nil
}

func VideoPostComment(data *receive.VideosPostCommentReceiveStruct, uid uint) (results interface{}, err error) {
	videoInfo := new(video.VideosContribution)
	err = videoInfo.FindByID(data.VideoID)
	if err != nil {
		return nil, fmt.Errorf("Video does not exist")
	}

	ct := comments.Comment{
		PublicModel: common.PublicModel{ID: data.ContentID},
	}
	CommentFirstID := ct.GetCommentFirstID()

	ctu := comments.Comment{
		PublicModel: common.PublicModel{ID: data.ContentID},
	}
	CommentUserID := ctu.GetCommentUserID()
	comment := comments.Comment{
		Uid:            uid,
		VideoID:        data.VideoID,
		Context:        data.Content,
		CommentID:      data.ContentID,
		CommentUserID:  CommentUserID,
		CommentFirstID: CommentFirstID,
	}
	if !comment.Create() {
		return nil, fmt.Errorf("Failure to publish")
	}

	//Socket push (when online)
	if _, ok := notice.Severe.UserMapChannel[videoInfo.UserInfo.ID]; ok {
		userChannel := notice.Severe.UserMapChannel[videoInfo.UserInfo.ID]
		userChannel.NoticeMessage(noticeModel.VideoComment)
	}

	return "Publish Successfully", nil
}

func GetVideoComment(data *receive.GetVideoCommentReceiveStruct) (results interface{}, err error) {
	videosContribution := new(video.VideosContribution)
	if !videosContribution.GetVideoComments(data.VideoID, data.PageInfo) {
		return nil, fmt.Errorf("Enquiry Failure")
	}
	return response.GetVideoContributionCommentsResponse(videosContribution), nil
}

func GetVideoManagementList(data *receive.GetVideoManagementListReceiveStruct, uid uint) (results interface{}, err error) {
	//Getting information on individual video releases
	list := new(video.VideosContributionList)
	err = list.GetVideoManagementList(data.PageInfo, uid)
	if err != nil {
		return nil, fmt.Errorf("Enquiry Failure")
	}
	res, err := response.GetVideoManagementListResponse(list)
	if err != nil {
		return nil, fmt.Errorf("Response Failure")
	}
	return res, nil
}

func LikeVideo(data *receive.LikeVideoReceiveStruct, uid uint) (results interface{}, err error) {
	//LIKE VIDEO
	videoInfo := new(video.VideosContribution)
	err = videoInfo.FindByID(data.VideoID)
	if err != nil {
		return nil, fmt.Errorf("Video does not exist")
	}
	lk := new(like.Likes)
	err = lk.Like(uid, data.VideoID, videoInfo.UserInfo.ID)
	if err != nil {
		return nil, fmt.Errorf("failure of an operation")
	}

	//Socket push (when online)
	if _, ok := notice.Severe.UserMapChannel[videoInfo.UserInfo.ID]; ok {
		userChannel := notice.Severe.UserMapChannel[videoInfo.UserInfo.ID]
		userChannel.NoticeMessage(noticeModel.VideoLike)
	}

	return "The operation was successful.", nil
}
