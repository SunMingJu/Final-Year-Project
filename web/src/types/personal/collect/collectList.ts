export interface GetFavoriteVideoListReq {
    favorite_id : number
}


//Video information
export interface VideoInfo {
    id : number
	uid :number 
	title : string
	video : string
	cover : string
    video_duration : number
    created_at : string
}

export type VideoInfoList = Array<VideoInfo>

export interface GetFavoriteVideoListRes {
    videoList :VideoInfoList
}