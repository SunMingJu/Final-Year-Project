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
	transcodingTask "simple-video-net/models/sundry/transcoding"
	"simple-video-net/models/users/attention"
	"simple-video-net/models/users/collect"
	"simple-video-net/models/users/favorites"
	noticeModel "simple-video-net/models/users/notice"
	"simple-video-net/models/users/record"
	"simple-video-net/utils/calculate"
	"simple-video-net/utils/conversion"
	"simple-video-net/utils/oss"
	"math"
	"os/exec"
	"strconv"
	"strings"
	"time"
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
	var width, height int
	if data.VideoUploadType == "local" {
		//If uploading locally
		width, height, err = calculate.GetVideoResolution(data.Video)
		if err != nil {
			global.Logger.Error("Failed to get video resolution")
			return nil, fmt.Errorf("Failed to get video resolution")
		}
	} else {
		mediaInfo, err := oss.GetMediaInfo(data.Media)
		if err != nil {
			return nil, fmt.Errorf("Failed to obtain video information,try again later")
		}
		width, _ = strconv.Atoi(*mediaInfo.Body.MediaInfo.FileInfoList[0].FileBasicInfo.Width)
		height, _ = strconv.Atoi(*mediaInfo.Body.MediaInfo.FileInfoList[0].FileBasicInfo.Height)
	}
	videoContribution := &video.VideosContribution{
		Uid:       uid,
		Title:     data.Title,
		Cover:     coverImg,
		Reprinted: conversion.BoolTurnInt8(*data.Reprinted),
		Label:     conversion.MapConversionString(data.Label),
		Introduce: data.Introduce,
		MediaID:   *data.Media,
		Heat:      0,
	}
	// Define a list of transcoding resolutions
	resolutions := []int{1080, 720, 480, 360}
	if height >= 1080 {
		resolutions = resolutions[1:]
		videoContribution.Video = videoSrc
	} else if height >= 720 && height < 1080 {
		resolutions = resolutions[2:]
		videoContribution.Video720p = videoSrc
	} else if height >= 480 && height < 720 {
		resolutions = resolutions[3:]
		videoContribution.Video480p = videoSrc
	} else if height >= 360 && height < 480 {
		resolutions = resolutions[4:]
		videoContribution.Video360p = videoSrc
	} else {
		global.Logger.Error("The uploaded video resolution is too low")
		return nil, fmt.Errorf("The uploaded video resolution is too low")
	}

	if !videoContribution.Create() {
		return nil, fmt.Errorf("fail to save")
	}
	//Transcode video
	go func(width, height int, video *video.VideosContribution) {
		//If the uploaded video is local, transcoding will begin.
		if data.VideoUploadType == "local" {
			//Local ffmpeg processing
			inputFile := data.Video
			sr := strings.Split(inputFile, ".")
		
			for _, r := range resolutions {
				//Calculating the transcoded width and height requires rounding
				w := int(math.Ceil(float64(r) / float64(height) * float64(width)))
				h := int(math.Ceil(float64(r)))
				if h >= height {
					continue
				}
				dst := sr[0] + fmt.Sprintf("_output_%dp."+sr[1], r)
				// TODO: Call the transcoding interface and save the transcoded video to the specified directory
				cmd := exec.Command("ffmpeg",
					"-i",
					inputFile,
					"-vf",
					fmt.Sprintf("scale=%d:%d", w, h),
					"-c:a",
					"copy",
					"-c:v",
					"libx264",
					"-preset",
					"medium",
					"-crf",
					"23",
					"-y",
					dst)
				err = cmd.Run()
				if err != nil {
					global.Logger.Errorf("video :%s :Transcode%d*%D failed cmd : %s ,err :%s", inputFile, w, h, cmd, err)
					continue
				}
				src, _ := json.Marshal(common.Img{
					Src: dst,
					Tp:  "local",
				})
				switch r {
				case 1080:
					videoContribution.Video = src
				case 720:
					videoContribution.Video720p = src
				case 480:
					videoContribution.Video480p = src
				case 360:
					videoContribution.Video360p = src
				}
				if !videoContribution.Save() {
					global.Logger.Errorf("video :%s : Transcode%d*%dAfter saving the video to the database failed", inputFile, w, h)
				}
				global.Logger.Infof("video :%s : Transcode%d*%D success", inputFile, w, h)
			}
			} else {
			inputFile := data.Video
			sr := strings.Split(inputFile, ".")
			//Cloud transcoding processing
			for _, r := range resolutions {
				//Get transcoding template
				var template string
				dst := sr[0] + fmt.Sprintf("_output_%dp."+sr[1], r)
				src, _ := json.Marshal(common.Img{
					Src: dst,
					Tp:  data.VideoUploadType,
				})
				switch r {
				case 1080:
					template = global.Config.AliyunOss.TranscodingTemplate1080p
					videoContribution.Video = src
				case 720:
					template = global.Config.AliyunOss.TranscodingTemplate720p
					videoContribution.Video720p = src
				case 480:
					template = global.Config.AliyunOss.TranscodingTemplate480p
					videoContribution.Video480p = src
				case 360:
					template = global.Config.AliyunOss.TranscodingTemplate360p
					videoContribution.Video360p = src
				}
				outputUrl, _ := conversion.SwitchIngStorageFun(data.VideoUploadType, dst)
				taskName := "Transcode : " + *data.Media + "time :" + time.Now().Format("2006.01.02 15:04:05") + " template : " + template
				jobInfo, err := oss.SubmitTranscodeJob(taskName, video.MediaID, outputUrl, template)
				if err != nil {
					global.Logger.Errorf("Video cloud transcoding : %s fail err : %s", outputUrl, err.Error())
					continue
				}
				task := &transcodingTask.TranscodingTask{
					TaskID:     *jobInfo.TranscodeParentJob.ParentJobId,
					VideoID:    video.ID,
					Resolution: r,
					Dst:        dst,
					Status:     0,
					Type:       transcodingTask.Aliyun,
				}
				if !task.AddTask() {
					global.Logger.Errorf("Video cloud transcoding task name: %s Video tasks later Failed to save to database", taskName)
				}
			}
		}
	}(width, height, videoContribution)
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
	//Update video data
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
				global.Logger.Error("Add playback volume error video video_id:", videoInfo.Watch(data.VideoID))videoInfo.Watch(data.VideoID))
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
