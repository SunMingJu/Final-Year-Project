import { DeleteChatItemReq, GetChatHistoryMsgReq, GetChatHistoryMsgRes, GetChatListRes, PersonalLetterReq } from "@/types/personal/chat/chat";
import { GetFavoriteVideoListReq, GetFavoriteVideoListRes } from "@/types/personal/collect/collectList";
import { createFavoritesReq, SaveCreateFavoritesDataReq } from "@/types/personal/collect/createFavorites";
import { DeleteFavoritesReq, FavoriteVideoReq, GetFavoritesListByFavoriteVideoReq, getFavoritesListByFavoriteVideoRes, GetFavoritesListRes } from "@/types/personal/collect/favorites";
import { GetLiveDataRes, SaveLiveDataReq } from "@/types/personal/live/setUp";
import { GetNoticeListReq, GetNoticeListRes } from "@/types/personal/notice/notice";
import { DeleteRecordByIDReq, GetRecordListReq, GetRecordListRes } from "@/types/personal/record/record";
import { changePasswordReq } from "@/types/personal/userInfo/security";
import { AttentionReq, DetermineNameExistsReq, DetermineNameExistsRes, SetUserInfoRes, UpdateAvatarReq, UserInfoRes } from "@/types/personal/userInfo/userInfo";
import httpRequest from "@/utils/requst";
//Get user information
export const getUserInfoRequist = () => {
    return httpRequest.post<UserInfoRes>('/user/getUserInfo');
}
//Determine if a username exists
export const determineNameExistsRequist = (params: DetermineNameExistsReq) => {
    return httpRequest.post<DetermineNameExistsRes>('/user/determineNameExists', params);
}
//Setting up user information
export const setUserInfoRequist = (params: UserInfoRes) => {
    return httpRequest.post<SetUserInfoRes>('/user/setUserInfo', params);
}
//Update user avatar
export const updateAvatarRequist = (params: UpdateAvatarReq) => {
    return httpRequest.post('/user/updateAvatar', params);
}
//Get user live information
export const getLiveDataRequist = () => {
    return httpRequest.post<GetLiveDataRes>('/user/getLiveData');
}
//Setting up live messages
export const saveLiveDataRequist = (params: SaveLiveDataReq) => {
    return httpRequest.post('/user/saveLiveData', params);
}

//change your password
export const changePassword = (params: changePasswordReq) => {
    return httpRequest.post('/user/changePassword', params);
}
//Change PasswordSend Verification Code
export const sendEmailVerificationCodeByChangePassword = () => {
    return httpRequest.post('/user/sendEmailVerificationCodeByChangePassword');
}
//focus
export const attention = (params: AttentionReq) => {
    return httpRequest.post('/user/attention', params);
}
//Create a favourite
export const createFavorites = (params: SaveCreateFavoritesDataReq | createFavoritesReq) => {
    return httpRequest.post('/user/createFavorites', params);
}
//Get favourites
export const getFavoritesList = () => {
    return httpRequest.post<GetFavoritesListRes>('/user/getFavoritesList');
}
//Delete favourites
export const deleteFavorites = (params: DeleteFavoritesReq) => {
    return httpRequest.post('/user/deleteFavorites', params);
}
//Collection Video
export const favoriteVideo = (params: FavoriteVideoReq) => {
    return httpRequest.post('/user/favoriteVideo', params);
}

//Get favourites in video page page need video ids
export const getFavoritesListByFavoriteVideo = (params: GetFavoritesListByFavoriteVideoReq) => {
    return httpRequest.post<getFavoritesListByFavoriteVideoRes>('/user/getFavoritesListByFavoriteVideo', params);
}
//Get favourite videos
export const getFavoriteVideoList = (params: GetFavoriteVideoListReq) => {
    return httpRequest.post<GetFavoriteVideoListRes>('/user/getFavoriteVideoList', params);
}

//Get History
export const getRecordList = (params: GetRecordListReq) => {
    return httpRequest.post<GetRecordListRes>('/user/getRecordList', params);
}

//Clear History
export const clearRecord = () => {
    return httpRequest.post('/user/clearRecord');
}
//Delete History
export const deleteRecordByID = (params: DeleteRecordByIDReq) => {
    return httpRequest.post('/user/deleteRecordByID', params);
}
//Getting Message Notifications
export const getNoticeList = (params: GetNoticeListReq) => {
    return httpRequest.post<GetNoticeListRes>('/user/getNoticeList', params);
}
//Get a list of recent chats
export const getChatList = () => {
    return httpRequest.post<GetChatListRes>('/user/getChatList');
}
//Actions triggered when clicking on a private message
export const personalLetter = (params: PersonalLetterReq) => {
    return httpRequest.post('/user/personalLetter', params);
}
//Delete Recent Chats
export const deleteChatItem = (params: DeleteChatItemReq) => {
    return httpRequest.post('/user/deleteChatItem', params);
}
//Load History Chat
export const getChatHistoryMsg = (params: GetChatHistoryMsgReq) => {
    return httpRequest.post<GetChatHistoryMsgRes>('/user/getChatHistoryMsg', params);
}