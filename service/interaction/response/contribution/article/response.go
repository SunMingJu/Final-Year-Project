package response

import (
	"easy-video-net/consts"
	"easy-video-net/models/common"
	"easy-video-net/models/contribution/article"
	"easy-video-net/models/contribution/article/classification"
	"easy-video-net/models/users"
	"easy-video-net/utils/conversion"
	"encoding/json"
	"github.com/dlclark/regexp2"
	"time"
)

//Comments Info
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

type GetArticleContributionListByUserResponseStruct struct {
	Id             uint      `json:"id"`
	Uid            uint      `json:"uid" `
	Title          string    `json:"title" `
	Cover          string    `json:"cover" `
	Label          []string  `json:"label" `
	Content        string    `json:"content"`
	IsComments     bool      `json:"is_comments"`
	Heat           int       `json:"heat"`
	LikesNumber    int       `json:"likes_number"`
	CommentsNumber int       `json:"comments_number"`
	Classification string    `json:"classification"`
	CreatedAt      time.Time `json:"created_at"`
}

type GetArticleContributionByIDResponseStruct struct {
	Id             uint             `json:"id"`
	Uid            uint             `json:"uid" `
	Title          string           `json:"title" `
	Cover          string           `json:"cover" `
	Label          []string         `json:"label" `
	Content        string           `json:"content"`
	IsComments     bool             `json:"is_comments"`
	Heat           int              `json:"heat"`
	LikesNumber    int              `json:"likes_number"`
	Comments       commentsInfoList `json:"comments"`
	CommentsNumber int              `json:"comments_number"`
	CreatedAt      time.Time        `json:"created_at"`
}

type GetArticleContributionListByUserResponseList []GetArticleContributionListByUserResponseStruct

type GetArticleContributionCommentsResponseStruct struct {
	Id             uint             `json:"id"`
	Comments       commentsInfoList `json:"comments"`
	CommentsNumber int              `json:"comments_number"`
}

//Getting the hierarchical structure
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

//Spanning Tree Structure
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

func GetArticleContributionListByUserResponse(l *article.ArticlesContributionList) GetArticleContributionListByUserResponseList {
	response := make(GetArticleContributionListByUserResponseList, 0)
	for _, v := range *l {
		coverSrc, _ := conversion.FormattingJsonSrc(v.Cover)

		//Regular replacement of first text
		reg := regexp2.MustCompile(`<(\S*?)[^>]*>.*?|<.*? />`, 0)
		match, _ := reg.Replace(v.Content, "", -1, -1)
		matchRune := []rune(match)
		if len(matchRune) > 100 {
			v.Content = string(matchRune[:100]) + "..."
		} else {
			v.Content = match
		}

		//Show only one label
		label := conversion.StringConversionMap(v.Label)
		if len(label) >= 3 {
			label = label[:1]
		}

		response = append(response, GetArticleContributionListByUserResponseStruct{
			Id:             v.ID,
			Uid:            v.Uid,
			Title:          v.Title,
			Cover:          coverSrc,
			Label:          label,
			Content:        v.Content,
			Classification: v.Classification.Label,
			IsComments:     conversion.Int8TurnBool(v.IsComments),
			Heat:           v.Heat,
			LikesNumber:    len(v.Likes),
			CommentsNumber: len(v.Comments),
			CreatedAt:      v.CreatedAt,
		})
	}
	return response
}

func GetArticleContributionListResponse(l *article.ArticlesContributionList) GetArticleContributionListByUserResponseList {
	response := make(GetArticleContributionListByUserResponseList, 0)
	for _, v := range *l {
		coverSrc, _ := conversion.FormattingJsonSrc(v.Cover)

		//Regular replacement of first text
		reg := regexp2.MustCompile(`<(\S*?)[^>]*>.*?|<.*? />`, 0)
		match, _ := reg.Replace(v.Content, "", -1, -1)
		matchRune := []rune(match)
		if len(matchRune) > 100 {
			v.Content = string(matchRune[:100]) + "..."
		} else {
			v.Content = match
		}

		//Show only one label
		label := conversion.StringConversionMap(v.Label)
		if len(label) >= 2 {
			label = label[:1]
		}

		response = append(response, GetArticleContributionListByUserResponseStruct{
			Id:             v.ID,
			Uid:            v.Uid,
			Title:          v.Title,
			Cover:          coverSrc,
			Label:          label,
			Content:        v.Content,
			Classification: v.Classification.Label,
			IsComments:     conversion.Int8TurnBool(v.IsComments),
			Heat:           v.Heat,
			LikesNumber:    len(v.Likes),
			CommentsNumber: len(v.Comments),
			CreatedAt:      v.CreatedAt,
		})
	}
	return response
}

func GetArticleContributionByIDResponse(vc *article.ArticlesContribution) GetArticleContributionByIDResponseStruct {
	coverSrc, _ := conversion.FormattingJsonSrc(vc.Cover)

	prefix, _ := conversion.SwitchTypeAsUrlPrefix(vc.ContentStorageType)
	//Regular Replacement src
	reg := regexp2.MustCompile(`(?<=(img[^>]*src="))[^"]*?`+consts.UrlPrefixSubstitutionEscape, 0)
	match, _ := reg.Replace(vc.Content, prefix, -1, -1)
	vc.Content = match

	label := conversion.StringConversionMap(vc.Label)
	if len(label) >= 2 {
		label = label[:1]
	}
	//comments
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
	response := GetArticleContributionByIDResponseStruct{
		Id:             vc.ID,
		Uid:            vc.Uid,
		Title:          vc.Title,
		Cover:          coverSrc,
		Label:          label,
		Content:        vc.Content,
		IsComments:     conversion.Int8TurnBool(vc.IsComments),
		Heat:           vc.Heat,
		LikesNumber:    len(vc.Likes),
		Comments:       commentsList,
		CommentsNumber: len(vc.Comments),
		CreatedAt:      vc.CreatedAt,
	}
	return response
}

func GetArticleContributionCommentsResponse(vc *article.ArticlesContribution) GetArticleContributionCommentsResponseStruct {
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

//ArticleClassificationInfo 
type ArticleClassificationInfo struct {
	ID       uint                          `json:"id"`
	AID      uint                          `json:"aid"`
	Label    string                        `json:"label"`
	Children ArticleClassificationInfoList `json:"children"`
}

type ArticleClassificationInfoList []*ArticleClassificationInfo

//Getting the hierarchical structure
func (l ArticleClassificationInfoList) getChildComment() ArticleClassificationInfoList {
	topList := ArticleClassificationInfoList{}
	for _, v := range l {
		if v.AID == 0 {
			//the top of a building
			topList = append(topList, &ArticleClassificationInfo{
				ID:       v.ID,
				AID:      v.AID,
				Label:    v.Label,
				Children: nil,
			})
		}
	}
	return classificationInfoListTree(topList, l)
}

//Spanning Tree Structure
func classificationInfoListTree(menus ArticleClassificationInfoList, allData ArticleClassificationInfoList) ArticleClassificationInfoList {
	//Cycle through all first level menus
	for k, v := range menus {
		//Query all submenus under this menu
		var nodes ArticleClassificationInfoList //Defining the child node directory
		for _, av := range allData {
			if av.AID == v.ID {
				nodes = append(nodes, &ArticleClassificationInfo{
					ID:       av.ID,
					AID:      av.AID,
					Label:    av.Label,
					Children: nil,
				})
			}
		}
		for kk, _ := range nodes {
			menus[k].Children = append(menus[k].Children, nodes[kk])
		}
		//Just query the sub-menu for recursion, query the three-level menu and four-level menu
		classificationListSecondTree(nodes, allData)
	}
	return menus
}

func classificationListSecondTree(menus ArticleClassificationInfoList, allData ArticleClassificationInfoList) ArticleClassificationInfoList {
	//Cycle through all first level menus
	for k, v := range menus {
		//Query all submenus under this menu
		var nodes ArticleClassificationInfoList //Defining the child node directory
		for _, av := range allData {
			if av.AID == v.ID {
				nodes = append(nodes, av)
			}
		}
		for kk, _ := range nodes {
			menus[k].Children = append(menus[k].Children, nodes[kk])
		}
		//Just query the sub-menu for recursion, query the three-level menu and four-level menu
		classificationListSecondTree(nodes, allData)
	}
	return menus
}

func GetArticleClassificationListResponse(cl *classification.ClassificationsList) ArticleClassificationInfoList {
	response := make(ArticleClassificationInfoList, 0)
	for _, v := range *cl {
		response = append(response, &ArticleClassificationInfo{
			ID:       v.ID,
			AID:      v.AID,
			Label:    v.Label,
			Children: make(ArticleClassificationInfoList, 0),
		})
	}
	return response.getChildComment()
}

type GetArticleTotalInfoResponseStruct struct {
	Classification    ArticleClassificationInfoList `json:"classification"`
	ArticleNum        int64                         `json:"article_num"`
	ClassificationNum int64                         `json:"classification_num"`
}

func GetArticleTotalInfoResponse(cl *classification.ClassificationsList, articleNum *int64, clNum int64) interface{} {
	classificationInfo := make(ArticleClassificationInfoList, 0)
	for _, v := range *cl {
		classificationInfo = append(classificationInfo, &ArticleClassificationInfo{
			ID:       v.ID,
			AID:      v.AID,
			Label:    v.Label,
			Children: make(ArticleClassificationInfoList, 0),
		})
	}
	classificationInfo = classificationInfo.getChildComment()

	return GetArticleTotalInfoResponseStruct{
		Classification:    classificationInfo,
		ArticleNum:        *articleNum,
		ClassificationNum: clNum,
	}
}

type GetArticleManagementListItem struct {
	ID               uint     `json:"id"`
	ClassificationID uint     `json:"classification_id"`
	Title            string   `json:"title"`
	Cover            string   `json:"cover"`
	CoverUrl         string   `json:"cover_url"`
	CoverUploadType  string   `json:"cover_upload_type"`
	Label            []string `json:"label"`
	Content          string   `json:"content"`
	IsComments       bool     `json:"is_comments" `
	Heat             int      `json:"heat"`
}

type GetArticleManagementListResponseStruct []GetArticleManagementListItem

func GetArticleManagementListResponse(al *article.ArticlesContributionList) (interface{}, error) {
	list := make(GetArticleManagementListResponseStruct, 0)
	for _, v := range *al {
		coverJson := new(common.Img)
		_ = json.Unmarshal(v.Cover, coverJson)
		cover, _ := conversion.FormattingJsonSrc(v.Cover)
		prefix, _ := conversion.SwitchTypeAsUrlPrefix(v.ContentStorageType)
		//Regular Replacement src
		reg := regexp2.MustCompile(`(?<=(img[^>]*src="))[^"]*?`+consts.UrlPrefixSubstitutionEscape, 0)
		match, _ := reg.Replace(v.Content, prefix, -1, -1)
		v.Content = match

		list = append(list, GetArticleManagementListItem{
			ID:               v.ID,
			ClassificationID: v.ClassificationID,
			Title:            v.Title,
			Cover:            cover,
			CoverUrl:         coverJson.Src,
			CoverUploadType:  coverJson.Tp,
			Label:            conversion.StringConversionMap(v.Label),
			Content:          v.Content,
			IsComments:       conversion.Int8TurnBool(v.IsComments),
			Heat:             v.Heat,
		})
	}

	return list, nil
}
