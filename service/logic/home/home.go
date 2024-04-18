package home

import (
	receive "simple-video-net/interaction/receive/home"
	response "simple-video-net/interaction/response/home"
	"simple-video-net/models/contribution/video"
	"simple-video-net/models/home/rotograph"
)

func GetHomeInfo(data *receive.GetHomeInfoReceiveStruct) (results interface{}, err error) {
	//Get homepage rotator
	rotographList := new(rotograph.List)
	err = rotographList.GetAll()
	if err != nil {
		return nil, err
	}

	//Get homepage testimonials
	videoList := new(video.VideosContributionList)
	err = videoList.GetHoneVideoList(data.PageInfo)

	if err != nil {
		return nil, err
	}
	res := &response.GetHomeInfoResponse{}
	res.Response(rotographList, videoList)

	return res, nil
}
