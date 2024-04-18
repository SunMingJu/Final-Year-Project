import DPlayer from "dplayer"


export interface LiveRoomProp {
    dp?: DPlayer   //Definition when not initialized
}

export interface GetArticleContributionListByUserResItem {
    id: number
    uid: number
    title: string
    cover: string
    classification: string
    label: Array<string>
    content: string
    is_comments: boolean
    heat: number
    likes_number: number
    comments_number: number
    created_at: string
    is_stay: boolean
}

export type GetArticleContributionListByUserRes = Array<GetArticleContributionListByUserResItem>

export interface GetLiveRoomInfoReq {
    room_id: number
}

export interface GetLiveRoomInfoRes {
    username: string
    photo: string
    live_title: string
    flv :string
}
