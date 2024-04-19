package contribution

import (
	"encoding/json"
	"fmt"
	"simple-video-net/consts"
	receive "simple-video-net/interaction/receive/contribution/article"
	response "simple-video-net/interaction/response/contribution/article"
	"simple-video-net/logic/users/notice"
	"simple-video-net/models/common"
	"simple-video-net/models/contribution/article"
	"simple-video-net/models/contribution/article/classification"
	"simple-video-net/models/contribution/article/comments"
	noticeModel "easy-video-net/models/users/notice"
	"simple-video-net/models/users/record"
	"simple-video-net/utils/conversion"

	"github.com/dlclark/regexp2"
)

func CreateArticleContribution(data *receive.CreateArticleContributionReceiveStruct, uid uint) (results interface{}, err error) {
	//Make content judgements
	for _, v := range data.Label {
		vRune := []rune(v)
		if len(vRune) > 7 {
			return nil, fmt.Errorf("Label length cannot be greater than 7 bits")
		}
	}

	coverImg, _ := json.Marshal(common.Img{
		Src: data.Cover,
		Tp:  data.CoverUploadType,
	})
	//Regular match replace url
	//Fetch url prefix
	prefix, err := conversion.SwitchTypeAsUrlPrefix(data.ArticleContributionUploadType)
	if err != nil {
		return nil, fmt.Errorf("Resource preservation method does not exist")
	}
	//regular matching substitution
	reg := regexp2.MustCompile(`(?<=(img[^>]*src="))[^"]*?`+prefix, 0)
	match, err := reg.Replace(data.Content, consts.UrlPrefixSubstitution, -1, -1)
	data.Content = match
	//insert data
	articlesContribution := article.ArticlesContribution{
		Uid:                uid,
		ClassificationID:   data.ClassificationID,
		Title:              data.Title,
		Cover:              coverImg,
		Timing:             conversion.BoolTurnInt8(*data.Timing),
		TimingTime:         data.TimingTime,
		Label:              conversion.MapConversionString(data.Label),
		Content:            data.Content,
		ContentStorageType: data.ArticleContributionUploadType,
		IsComments:         conversion.BoolTurnInt8(*data.Comments),
		Heat:               0,
	}

	if *data.Timing {
		//Push related after posting a video (to be developed)
	}
	if !articlesContribution.Create() {
		return nil, fmt.Errorf("fail to save")
	}
	return "Save Successful", nil
}

func UpdateArticleContribution(data *receive.UpdateArticleContributionReceiveStruct, uid uint) (results interface{}, err error) {
	//Updated columns
	articleInfo := new(article.ArticlesContribution)
	if !articleInfo.GetInfoByID(data.ID) {
		return nil, fmt.Errorf("Operation video does not exist")
	}
	if articleInfo.Uid != uid {
		return nil, fmt.Errorf("unauthorised operation")
	}
	coverImg, _ := json.Marshal(common.Img{
		Src: data.Cover,
		Tp:  data.CoverUploadType,
	})
	updateList := map[string]interface{}{
		"cover":             coverImg,
		"title":             data.Title,
		"label":             conversion.MapConversionString(data.Label),
		"content":           data.Content,
		"is_comments":       data.Comments,
		"classification_id": data.ClassificationID,
	}
	//Update video data
	if !articleInfo.Update(updateList) {
		return nil, fmt.Errorf("Failed to update data")
	}
	return "Successful update", nil
}

func DeleteArticleByID(data *receive.DeleteArticleByIDReceiveStruct, uid uint) (results interface{}, err error) {
	al := new(article.ArticlesContribution)
	if !al.Delete(data.ID, uid) {
		return nil, fmt.Errorf("Failed to delete")
	}
	return "Deleted successfully", nil
}

func GetArticleContributionList(data *receive.GetArticleContributionListReceiveStruct) (results interface{}, err error) {
	articlesContribution := new(article.ArticlesContributionList)
	if !articlesContribution.GetList(data.PageInfo) {
		return nil, fmt.Errorf("Enquiry Failure")
	}
	return response.GetArticleContributionListResponse(articlesContribution), nil
}

func GetArticleContributionListByUser(data *receive.GetArticleContributionListByUserReceiveStruct) (results interface{}, err error) {
	articlesContribution := new(article.ArticlesContributionList)
	if !articlesContribution.GetListByUid(data.UserID) {
		return nil, fmt.Errorf("Enquiry Failure")
	}
	return response.GetArticleContributionListByUserResponse(articlesContribution), nil
}

func GetArticleContributionByID(data *receive.GetArticleContributionByIDReceiveStruct, uid uint) (results interface{}, err error) {
	articlesContribution := new(article.ArticlesContribution)
	if !articlesContribution.GetInfoByID(data.ArticleID) {
		return nil, fmt.Errorf("Enquiry Failure")
	}
	if uid > 0 {
		//add history
		rd := new(record.Record)
		err = rd.AddArticleRecord(uid, data.ArticleID)
		if err != nil {
			return nil, fmt.Errorf("Failed to add history")
		}
	}
	return response.GetArticleContributionByIDResponse(articlesContribution), nil
}

func ArticlePostComment(data *receive.ArticlesPostCommentReceiveStruct, uid uint) (results interface{}, err error) {
	articleInfo := new(article.ArticlesContribution)
	if !articleInfo.GetInfoByID(data.ArticleID) {
		return nil, fmt.Errorf("Comment article does not exist")
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
		ArticleID:      data.ArticleID,
		Context:        data.Content,
		CommentID:      data.ContentID,
		CommentUserID:  CommentUserID,
		CommentFirstID: CommentFirstID,
	}
	if !comment.Create() {
		return nil, fmt.Errorf("Failure to publish")
	}

	//Socket push (when online)
	if _, ok := notice.Severe.UserMapChannel[articleInfo.UserInfo.ID]; ok {
		userChannel := notice.Severe.UserMapChannel[articleInfo.UserInfo.ID]
		userChannel.NoticeMessage(noticeModel.ArticleComment)
	}

	return "Publish Successfully", nil
}

func GetArticleComment(data *receive.GetArticleCommentReceiveStruct) (results interface{}, err error) {
	articlesContribution := new(article.ArticlesContribution)
	if !articlesContribution.GetArticleComments(data.ArticleID, data.PageInfo) {
		return nil, fmt.Errorf("Enquiry Failure")
	}
	return response.GetArticleContributionCommentsResponse(articlesContribution), nil
}

func GetArticleClassificationList() (results interface{}, err error) {
	cn := new(classification.ClassificationsList)
	err = cn.FindAll()
	if err != nil {
		return nil, fmt.Errorf("Enquiry Failure")
	}
	return response.GetArticleClassificationListResponse(cn), nil
}

func GetArticleTotalInfo() (results interface{}, err error) {
	//Query the number of articles
	articleNm := new(int64)
	al := new(article.ArticlesContributionList)
	al.GetAllCount(articleNm)
	//Query article classification information
	cn := make(classification.ClassificationsList, 0)
	err = cn.FindAll()
	if err != nil {
		return nil, fmt.Errorf("Enquiry Failure")
	}
	cnNum := int64(len(cn))
	return response.GetArticleTotalInfoResponse(&cn, articleNm, cnNum), nil
}

func GetArticleManagementList(data *receive.GetArticleManagementListReceiveStruct, uid uint) (results interface{}, err error) {
	//Access to personal publishing columns
	list := new(article.ArticlesContributionList)
	err = list.GetArticleManagementList(data.PageInfo, uid)
	if err != nil {
		return nil, fmt.Errorf("Enquiry Failure")
	}
	res, err := response.GetArticleManagementListResponse(list)
	if err != nil {
		return nil, fmt.Errorf("Response Failure")
	}
	return res, nil
}
