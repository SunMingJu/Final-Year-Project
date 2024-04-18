package contribution

import (
	"fmt"
	receive "simple-video-net/interaction/receive/contribution/discuss"
	response "simple-video-net/interaction/response/contribution/discuss"
	"simple-video-net/models/contribution/article"
	articleComments "simple-video-net/models/contribution/article/comments"
	"simple-video-net/models/contribution/video"
	"simple-video-net/models/contribution/video/barrage"
	videoComments "simple-video-net/models/contribution/video/comments"
)

func GetDiscussVideoList(data *receive.GetDiscussVideoListReceiveStruct, uid uint) (results interface{}, err error) {
	//Get user-posted videos
	videoList := new(video.VideosContributionList)
	err = videoList.GetDiscussVideoCommentList(uid)
	if err != nil {
		return nil, fmt.Errorf("Failed to query video related information")
	}
	videoIDs := make([]uint, 0)
	for _, v := range *videoList {
		videoIDs = append(videoIDs, v.ID)
	}
	//Getting video information
	cml := new(videoComments.CommentList)
	err = cml.GetCommentListByIDs(videoIDs, data.PageInfo)
	if err != nil {
		return nil, fmt.Errorf("Failed to query video comment information")
	}

	return response.GetDiscussVideoListResponse(cml), nil
}

func GetDiscussArticleList(data *receive.GetDiscussArticleListReceiveStruct, uid uint) (results interface{}, err error) {
	//Get user-posted columns
	articleList := new(article.ArticlesContributionList)
	err = articleList.GetDiscussArticleCommentList(uid)
	if err != nil {
		return nil, fmt.Errorf("Failed to search for column-related information")
	}
	articleIDs := make([]uint, 0)
	for _, v := range *articleList {
		articleIDs = append(articleIDs, v.ID)
	}
	//Get article information
	cml := new(articleComments.CommentList)
	err = cml.GetCommentListByIDs(articleIDs, data.PageInfo)
	if err != nil {
		return nil, fmt.Errorf("Failed to query article comment information")
	}
	return response.GetDiscussArticleListResponse(cml), nil
}

func GetDiscussBarrageList(data *receive.GetDiscussBarrageListReceiveStruct, uid uint) (results interface{}, err error) {
	//Get user-posted videos
	videoList := new(video.VideosContributionList)
	err = videoList.GetDiscussVideoCommentList(uid)
	if err != nil {
		return nil, fmt.Errorf("Failed to query video related information")
	}
	videoIDs := make([]uint, 0)
	for _, v := range *videoList {
		videoIDs = append(videoIDs, v.ID)
	}
	//Get video pop-up information
	cml := new(barrage.BarragesList)
	err = cml.GetBarrageListByIDs(videoIDs, data.PageInfo)
	if err != nil {
		return nil, fmt.Errorf("Failed to query video pop-up information")
	}
	return response.GetDiscussBarrageListResponse(cml), nil
}
