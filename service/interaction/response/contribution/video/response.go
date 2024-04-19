package response

import (
	"encoding/json"
	"simple-video-net/models/common"
	"simple-video-net/models/contribution/video"
	"simple-video-net/models/contribution/video/barrage"
	"simple-video-net/models/users"
	"simple-video-net/utils/conversion"
	"time"
)

// video info
type Info struct {
	ID             uint             `json:"id"`
	Uid            uint             `json:"uid" `
	Title          string           `json:"title" `
	Video          string           `json:"video"`
	Video720p      string           `json:"video_720p"`
	Video480p      string           `json:"video_480p"`
	Video360p      string           `json:"video_360p"`
	Cover          string           `json:"cover" `
	VideoDuration  int64            `json:"video_duration"`
	Label          []string         `json:"label"`
	Introduce      string           `json:"introduce"`
	Heat           int              `json:"heat"`
	BarrageNumber  int              `json:"barrageNumber"`
	Comments       commentsInfoList `json:"comments"`
	IsLike         bool             `json:"is_like"`
	IsCollect      bool             `json:"is_collect"`
	CommentsNumber int              `json:"comments_number"`
	CreatorInfo    creatorInfo      `json:"creatorInfo"`
	CreatedAt      time.Time        `json:"created_at"`
}

// creatorInfo
type creatorInfo struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	Avatar      string `json:"avatar"`
	Signature   string `json:"signature"`
	IsAttention bool   `json:"is_attention"`
}

// recommendVideo info
type recommendVideo struct {
	ID            uint      `json:"id"`
	Uid           uint      `json:"uid" `
	Title         string    `json:"title" `
	Video         string    `json:"video"`
	Cover         string    `json:"cover" `
	VideoDuration int64     `json:"video_duration"`
	Label         []string  `json:"label"`
	Introduce     string    `json:"introduce"`
	Heat          int       `json:"heat"`
	BarrageNumber int       `json:"barrageNumber"`
	Username      string    `json:"username"`
	CreatedAt     time.Time `json:"created_at"`
}
type RecommendList []recommendVideo

type Response struct {
	VideoInfo     Info          `json:"videoInfo"`
	RecommendList RecommendList `json:"recommendList"`
}

func GetVideoContributionByIDResponse(vc *video.VideosContribution, recommendVideoList *video.VideosContributionList, isAttention bool, isLike bool, isCollect bool) Response {
	//Processing video key information
	creatorAvatar, _ := conversion.FormattingJsonSrc(vc.UserInfo.Photo)
	cover, _ := conversion.FormattingJsonSrc(vc.Cover)
	videoSrc, _ := conversion.FormattingJsonSrc(vc.Video)
	video720pSrc, _ := conversion.FormattingJsonSrc(vc.Video720p)
	video480pSrc, _ := conversion.FormattingJsonSrc(vc.Video480p)
	video360pSrc, _ := conversion.FormattingJsonSrc(vc.Video360p)

	//comment
	comments := commentsInfoList{}
	for _, v := range vc.Comments {
		photo, _ := conversion.FormattingJsonSrc(v.UserInfo.Photo)
		commentUser := users.User{}
		commentUser.Find(v.CommentUserID)
		comments = append(comments, &commentsInfo{
			ID:              v.ID,
			CommentID:       v.CommentID,
			CommentFirstID:  v.CommentFirstID,
			CommentUserID:   v.CommentUserID,
			CommentUserName: commentUser.Username,
			CreatedAt:       v.CreatedAt,
			Context:         v.Context,
			Uid:             v.UserInfo.ID,
			Username:        v.UserInfo.Username,
			Photo:           photo,
		})
	}

	commentsList := comments.getChildComment()

	response := Response{
		VideoInfo: Info{
			ID:             vc.ID,
			Uid:            vc.Uid,
			Title:          vc.Title,
			Video:          videoSrc,
			Cover:          cover,
			VideoDuration:  vc.VideoDuration,
			Label:          conversion.StringConversionMap(vc.Label),
			Introduce:      vc.Introduce,
			Heat:           vc.Heat,
			BarrageNumber:  len(vc.Barrage),
			Comments:       commentsList,
			CommentsNumber: len(commentsList),
			IsLike:         isLike,
			IsCollect:      isCollect,
			CreatorInfo: creatorInfo{
				ID:          vc.UserInfo.ID,
				Username:    vc.UserInfo.Username,
				Avatar:      creatorAvatar,
				Signature:   vc.UserInfo.Signature,
				IsAttention: isAttention,
			},
			CreatedAt: vc.CreatedAt,
		},
	}
	//Handling of testimonial videos
	rl := make(RecommendList, 0)
	for _, lk := range *recommendVideoList {
		cover, _ := conversion.FormattingJsonSrc(lk.Cover)
		videoSrc, _ := conversion.FormattingJsonSrc(lk.Video)
		info := recommendVideo{
			ID:            lk.ID,
			Uid:           lk.Uid,
			Title:         lk.Title,
			Video:         videoSrc,
			Video720p:      video720pSrc,
			Video480p:      video480pSrc,
			Video360p:      video360pSrc,
			Cover:         cover,
			VideoDuration: lk.VideoDuration,
			Label:         conversion.StringConversionMap(lk.Label),
			Introduce:     lk.Introduce,
			Heat:          lk.Heat,
			BarrageNumber: len(lk.Barrage),
			Username:      lk.UserInfo.Username,
			CreatedAt:     lk.CreatedAt,
		}
		rl = append(rl, info)
	}
	response.RecommendList = rl
	return response
}

func GetVideoBarrageResponse(list *barrage.BarragesList) interface{} {
	barrageInfoList := make([][]interface{}, 0)
	for _, v := range *list {
		info := make([]interface{}, 0)
		info = append(info, v.Time)
		info = append(info, v.Type)
		info = append(info, v.Color)
		info = append(info, v.Author)
		info = append(info, v.Text)
		barrageInfoList = append(barrageInfoList, info)
	}
	return barrageInfoList
}

// Get video pop-up response
type barrageInfo struct {
	Time     int       `json:"time"`
	Text     string    `json:"text"`
	SendTime time.Time `json:"sendTime"`
}

type barrageInfoList []barrageInfo

func GetVideoBarrageListResponse(list *barrage.BarragesList) interface{} {
	barrageList := make(barrageInfoList, 0)
	for _, v := range *list {
		info := barrageInfo{
			Time:     int(v.Time),
			Text:     v.Text,
			SendTime: v.PublicModel.CreatedAt,
		}
		barrageList = append(barrageList, info)
	}
	return barrageList
}

// Comments
type commentsInfo struct {
	ID              uint             `json:"id"`
	CommentID       uint             `json:"comment_id"`
	CommentFirstID  uint             `json:"comment_first_id"`
	CreatedAt       time.Time        `json:"created_at"`
	Context         string           `json:"context"`
	Uid             uint             `json:"uid"`
	Username        string           `json:"username"`
	Photo           string           `json:"photo"`
	CommentUserID   uint             `json:"comment_user_id"`
	CommentUserName string           `json:"comment_user_name"`
	LowerComments   commentsInfoList `json:"lowerComments"`
}

type commentsInfoList []*commentsInfo

type GetArticleContributionCommentsResponseStruct struct {
	Id             uint             `json:"id"`
	Comments       commentsInfoList `json:"comments"`
	CommentsNumber int              `json:"comments_number"`
}

// Getting the hierarchical structure
func (l commentsInfoList) getChildComment() commentsInfoList {
	topList := commentsInfoList{}
	for _, v := range l {
		if v.CommentID == 0 {
			//the top of a building
			topList = append(topList, v)
		}
	}
	return commentsInfoListSecondTree(topList, l)
}

// Spanning Tree Structure
func commentsInfoListTree(menus commentsInfoList, allData commentsInfoList) commentsInfoList {
	//Cycle through all first level menus
	for k, v := range menus {
		//Query all submenus under this menu
		var nodes commentsInfoList //Defining the child node directory
		for _, av := range allData {
			if av.CommentID == v.ID {
				nodes = append(nodes, av)
			}
		}
		for kk, _ := range nodes {
			menus[k].LowerComments = append(menus[k].LowerComments, nodes[kk])
		}
		//Just query the sub-menu for recursion, query the three-level menu and four-level menu
		commentsInfoListTree(nodes, allData)
	}
	return menus
}

func commentsInfoListSecondTree(menus commentsInfoList, allData commentsInfoList) commentsInfoList {
	//Cycle through all first level menus
	for k, v := range menus {
		//Query all submenus under this menu
		var nodes commentsInfoList //Defining the child node directory
		for _, av := range allData {
			if av.CommentFirstID == v.ID {
				nodes = append(nodes, av)
			}
		}
		for kk, _ := range nodes {
			menus[k].LowerComments = append(menus[k].LowerComments, nodes[kk])
		}
		//Just query the sub-menu for recursion, query the three-level menu and four-level menu
		commentsInfoListTree(nodes, allData)
	}
	return menus
}

func GetVideoContributionCommentsResponse(vc *video.VideosContribution) GetArticleContributionCommentsResponseStruct {
	//comment
	comments := commentsInfoList{}
	for _, v := range vc.Comments {
		photo, _ := conversion.FormattingJsonSrc(v.UserInfo.Photo)
		commentUser := users.User{}
		commentUser.Find(v.CommentUserID)
		comments = append(comments, &commentsInfo{
			ID:              v.ID,
			CommentID:       v.CommentID,
			CommentFirstID:  v.CommentFirstID,
			CommentUserID:   v.CommentUserID,
			CommentUserName: commentUser.Username,
			CreatedAt:       v.CreatedAt,
			Context:         v.Context,
			Uid:             v.UserInfo.ID,
			Username:        v.UserInfo.Username,
			Photo:           photo,
		})
	}
	commentsList := comments.getChildComment()
	//output
	response := GetArticleContributionCommentsResponseStruct{
		Id:             vc.ID,
		Comments:       commentsList,
		CommentsNumber: len(vc.Comments),
	}
	return response
}

type GetVideoManagementListItem struct {
	ID              uint      `json:"id"`
	Uid             uint      `json:"uid" `
	Title           string    `json:"title" `
	Video           string    `json:"video"`
	Cover           string    `json:"cover" `
	Reprinted       bool      `json:"reprinted"`
	CoverUrl        string    `json:"cover_url"`
	CoverUploadType string    `json:"cover_upload_type"`
	VideoDuration   int64     `json:"video_duration"`
	Label           []string  `json:"label"`
	Introduce       string    `json:"introduce"`
	Heat            int       `json:"heat"`
	BarrageNumber   int       `json:"barrageNumber"`
	CommentsNumber  int       `json:"comments_number"`
	CreatedAt       time.Time `json:"created_at"`
}

type GetVideoManagementList []GetVideoManagementListItem

func GetVideoManagementListResponse(vc *video.VideosContributionList) (interface{}, error) {
	list := make(GetVideoManagementList, 0)
	for _, v := range *vc {
		coverJson := new(common.Img)
		_ = json.Unmarshal(v.Cover, coverJson)
		cover, _ := conversion.FormattingJsonSrc(v.Cover)
		videoSrc, _ := conversion.FormattingJsonSrc(v.Video)
		info := GetVideoManagementListItem{
			ID:              v.ID,
			Uid:             v.Uid,
			Title:           v.Title,
			Video:           videoSrc,
			Cover:           cover,
			Reprinted:       conversion.Int8TurnBool(v.Reprinted),
			CoverUploadType: coverJson.Tp,
			CoverUrl:        coverJson.Src,
			VideoDuration:   v.VideoDuration,
			Label:           conversion.StringConversionMap(v.Label),
			Introduce:       v.Introduce,
			Heat:            v.Heat,
			BarrageNumber:   len(v.Barrage),
			CommentsNumber:  len(v.Comments),
			CreatedAt:       v.CreatedAt,
		}
		list = append(list, info)
	}
	return list, nil
}
