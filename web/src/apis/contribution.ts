import { AreateArticleContributionReq, GetArticleContributionListByUserReq, GetArticleClassificationListRes, getArticleTotalInfoRes, UpdateArticleContributionReq } from "@/types/creation/contribute/contributePage/articleContribution";
import { CreateVideoContributionReq, UpdateVideoContributionReq } from "@/types/creation/contribute/contributePage/vdeoContribution";
import { DeleteArticleByIDReq, GetArticleManagementListReq, GetArticleManagementListRes } from "@/types/creation/manuscript/article";
import { DeleteVideoByIDReq, GetVideoManagementListReq, GetVideoManagementListRes } from "@/types/creation/manuscript/video";
import { GetDiscussArticleListReq, GetDiscussArticleListRes, GetDiscussVideoListReq, GetDiscussVideoListRes } from "@/types/creation/discuss/comment";
import { ArticlePostCommentReq, GetArticleCommentReq, GetArticleCommentRes, GetArticleContributionByIDReq, GetArticleContributionByIDRes } from "@/types/show/article/article";
import { GetVideoContributionByIDReq, GetVideoContributionByIDRes, VideoPostCommentReq, GetVideoBarrageListReq, GetVideoBarrageListRes, SendVideoBarrageReq, GetVideoCommentReq, GetVideoCommentRes, LikeVideoReq } from "@/types/show/video/video";
import { GetArticleContributionListReq, GetArticleContributionListRes } from "@/types/home/column";
import { GetArticleContributionListByUserRes } from "@/types/live/liveRoom";
import httpRequest from "@/utils/requst"
import { GetDiscussBarrageListReq, GetDiscussBarrageListRes } from "@/types/creation/discuss/barrage";


//Post Video
export const createVideoContribution = (params: CreateVideoContributionReq) => {
    return httpRequest.post('/contribution/createVideoContribution', params);
}

//Update Video
export const updateVideoContribution = (params: UpdateVideoContributionReq) => {
    return httpRequest.post('/contribution/updateVideoContribution', params);
}

//Publishing Column
export const createArticleContribution = (params: AreateArticleContributionReq) => {
    return httpRequest.post('/contribution/createArticleContribution', params);
}

//Updated columns
export const updateArticleContribution = (params: UpdateArticleContributionReq) => {
    return httpRequest.post('/contribution/updateArticleContribution', params);
}

//Enquiry Column List
export const getArticleContributionList = (params: GetArticleContributionListReq) => {
    return httpRequest.post<GetArticleContributionListRes>('/contribution/getArticleContributionList', params);
}
//Access to column information based on users
export const getArticleContributionListByUser = (params: GetArticleContributionListByUserReq) => {
    return httpRequest.post<GetArticleContributionListByUserRes>('/contribution/getArticleContributionListByUser', params);
}

//Get article information according to article ID
export const getArticleContributionByID = (params: GetArticleContributionByIDReq) => {
    return httpRequest.post<GetArticleContributionByIDRes>('/contribution/getArticleContributionByID', params);
}

//Article Posting Comments
export const articlePostComment = (params: ArticlePostCommentReq) => {
    return httpRequest.post('/contribution/articlePostComment', params);
}

//Get article comments separately
export const getArticleComment = (params: GetArticleCommentReq) => {
    return httpRequest.post<GetArticleCommentRes>('/contribution/getArticleComment', params);
}

//Get video information based on id
export const getVideoContributionByID = (params: GetVideoContributionByIDReq) => {
    return httpRequest.post<GetVideoContributionByIDRes>('/contribution/getVideoContributionByID', params);
}

export const danmakuApi = import.meta.env.VITE_BASE_URL + '/contribution/video/barrage/'

//Get the list of visual pop-ups
export const getVideoBarrageList = (params: GetVideoBarrageListReq) => {
    return httpRequest.get<GetVideoBarrageListRes>('/contribution/getVideoBarrageList', params);
}

//Get the list of visual pop-ups
export const sendVideoBarrage = (params: SendVideoBarrageReq) => {
    return httpRequest.post('/contribution/video/barrage/v3/', params);
}

//Article Posting Comments
export const videoPostComment = (params: VideoPostCommentReq) => {
    return httpRequest.post('/contribution/videoPostComment', params);
}

//Get video reviews separately
export const getVideoComment = (params: GetVideoCommentReq) => {
    return httpRequest.post<GetVideoCommentRes>('/contribution/getVideoComment', params);
}


//Get Article Categories
export const getArticleClassificationList = () => {
    return httpRequest.post<GetArticleClassificationListRes>('/contribution/getArticleClassificationList');
}

//Get article category information
export const getArticleTotalInfo = () => {
    return httpRequest.post<getArticleTotalInfoRes>('/contribution/getArticleTotalInfo');
}

//Creative Centre to get personal release videos
export const getVideoManagementList = (params: GetVideoManagementListReq) => {
    return httpRequest.post<GetVideoManagementListRes>('/contribution/getVideoManagementList',params);
}

//Delete video by id
export const deleteVideoByID = (params: DeleteVideoByIDReq) => {
    return httpRequest.post('/contribution/deleteVideoByID',params);
}

//Creative Writing Centre access to personal publishing columns
export const getArticleManagementList = (params: GetArticleManagementListReq) => {
    return httpRequest.post<GetArticleManagementListRes>('/contribution/getArticleManagementList',params);
}


//Delete column by id
export const deleteArticleByID = (params: DeleteArticleByIDReq) => {
    return httpRequest.post('/contribution/deleteArticleByID',params);
}

//Get ReviewsManage Video Reviews
export const getDiscussVideoList = (params: GetDiscussVideoListReq) => {
    return httpRequest.post<GetDiscussVideoListRes>('/contribution/getDiscussVideoList',params);
}

//Get CommentsManage Article Comments
export const getDiscussArticleList = (params: GetDiscussArticleListReq) => {
    return httpRequest.post<GetDiscussArticleListRes>('/contribution/getDiscussArticleList',params);
}

//Get pop-up management
export const getDiscussBarrageList = (params: GetDiscussBarrageListReq) => {
    return httpRequest.post<GetDiscussBarrageListRes>('/contribution/getDiscussBarrageList',params);
}

//Get pop-up management
export const likeVideo = (params: LikeVideoReq) => {
    return httpRequest.post('/contribution/likeVideo',params);
}
