package callback

import (
	"simple-video-net/global"
	receive "simple-video-net/interaction/receive/callback"
	"simple-video-net/models/common"
	"simple-video-net/models/contribution/video"
	transcodingTask "simple-video-net/models/sundry/transcoding"
	"encoding/json"
)

func AliyunTranscodingMedia(data *receive.AliyunMediaCallback[receive.AliyunTranscodingMediaStruct]) (results interface{}	//Query tasks
	taskInfo := new(transcodingTask.TranscodingTask)
	err = taskInfo.GetInfoByTaskID(data.MessageBody.ParentJobId)
	if err != nil {
		global.Logger.Errorf("Alibaba Cloud Media Service video transcoding callback failed，task id %s", data.MessageBody.Jobs)
		return nil, nil

	}
	if taskInfo.ID <= 0 {
		global.Logger.Errorf("Alibaba Cloud Media Service video transcoding callback failed: task does not exist，task id %s", data.MessageBody.Jobs)
		return nil, nil
	}
	if taskInfo.Status == 1 {

		global.Logger.Infof("Alibaba Cloud Media Service video transcoding callback has been processed and does not need to be processed, task id %s", data.MessageBody.Jobs)
		return nil, nil
	}
	if data.MessageBody.Status != "Success" {
		taskInfo.Status = 2
		if !taskInfo.Save() {
			global.Logger.Errorf("Alibaba Cloud Media Service fails to update task information after successful callback, task id %s", data.MessageBody.Jobs)
		}
		global.Logger.Errorf("Alibaba Cloud Media Service video transcoding failed, task id %s", data.MessageBody.Jobs)
		return nil, nil
	}
	videoInfo := new(video.VideosContribution)
	err = videoInfo.FindByID(taskInfo.VideoID)
	if err != nil {
		global.Logger.Errorf("Alibaba Cloud Media Service callback failed to obtain video information, task id %s", data.MessageBody.Jobs)
		return nil, nil
	}
	if videoInfo.ID <= 0 {
		global.Logger.Errorf("Alibaba Cloud Media Service callback video has been deleted, task id %s", data.MessageBody.Jobs)
		return nil, nil
	}
	src, _ := json.Marshal(common.Img{
		Src: taskInfo.Dst,
		Tp:  "aliyunOss",
	})
	switch taskInfo.Resolution {
	case 1080:
		videoInfo.Video = src
	case 720:
		videoInfo.Video720p = src
	case 480:
		videoInfo.Video480p = src
	case 360:
		videoInfo.Video360p = src
	}
	if !videoInfo.Save() {
		global.Logger.Errorf("Alibaba Cloud Media Service callback successful, failed to save video information, task id %s", data.MessageBody.Jobs)
	}
	taskInfo.Status = 1
	if !taskInfo.Save() {
		global.Logger.Errorf("Alibaba Cloud Media Service fails to update task information after successful callback, task id %s", data.MessageBody.Jobs)
	}

	return nil, nil
}