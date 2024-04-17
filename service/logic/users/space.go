package users

import (
	receive "easy-video-net/interaction/receive/users"
	response "easy-video-net/interaction/response/users"
	"easy-video-net/models/contribution/article"
	"easy-video-net/models/contribution/video"
	"easy-video-net/models/users"
	"easy-video-net/models/users/attention"
	"fmt"
)

func GetSpaceIndividual(data *receive.GetSpaceIndividualReceiveStruct, uid uint) (results interface{}, err error) {
	//Get user information
	user := new(users.User)
	user.Find(data.ID)
	isAttention := false
	at := new(attention.Attention)
	if uid != 0 {
		//Get attention or not
		isAttention = at.IsAttention(uid, data.ID)
	}
	//Getting Followers and Fans
	attentionNum, err := at.GetAttentionNum(data.ID)
	if err != nil {
		return nil, fmt.Errorf("Enquiry Failure")
	}
	vermicelliNum, err := at.GetVermicelliNum(data.ID)
	if err != nil {
		return nil, fmt.Errorf("Enquiry Failure")
	}
	return response.GetSpaceIndividualResponse(user, isAttention, attentionNum, vermicelliNum)
}

func GetReleaseInformation(data *receive.GetReleaseInformationReceiveStruct) (results interface{}, err error) {
	//Get Video List
	videoList := new(video.VideosContributionList)
	err = videoList.GetVideoListBySpace(data.ID)
	if err != nil {
		return nil, fmt.Errorf("Failed to query space video")
	}
	articleList := new(article.ArticlesContributionList)
	err = articleList.GetArticleBySpace(data.ID)
	if err != nil {
		return nil, fmt.Errorf("Failed to query space column")
	}
	res, err := response.GetReleaseInformationResponse(videoList, articleList)
	if err != nil {
		return nil, fmt.Errorf("Response Failure")
	}
	return res, nil
}

func GetAttentionList(data *receive.GetAttentionListReceiveStruct, id uint) (results interface{}, err error) {
	//Get user's followers list
	al := new(attention.AttentionsList)
	err = al.GetAttentionList(data.ID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get")
	}
	//Get the users you follow
	ual := new(attention.AttentionsList)
	arr, err := ual.GetAttentionListByIdArr(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to get")
	}
	res, err := response.GetAttentionListResponse(al, arr)
	if err != nil {
		return nil, fmt.Errorf("Response Failure")
	}
	return res, nil
}

func GetVermicelliList(data *receive.GetVermicelliListReceiveStruct, id uint) (results interface{}, err error) {
	//Get user's fan list
	al := new(attention.AttentionsList)
	err = al.GetVermicelliList(data.ID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get")
	}
	//Get the users you follow
	ual := new(attention.AttentionsList)
	arr, err := ual.GetAttentionListByIdArr(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to get")
	}
	res, err := response.GetVermicelliListResponse(al, arr)
	if err != nil {
		return nil, fmt.Errorf("Response Failure")
	}
	return res, nil
}
