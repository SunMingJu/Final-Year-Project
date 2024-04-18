import httpRequest from "@/utils/requst"

import { GetLiveRoomRes } from "@/types/home/home"
import { GetLiveRoomInfoReq, GetLiveRoomInfoRes } from "@/types/live/liveRoom";
import { GetBeLiveListRes } from "@/types/home/live";


export const getLiveRoom = () => {
    return httpRequest.post<GetLiveRoomRes>('/live/getLiveRoom');
}

//Get Live Streaming Information
export const getLiveRoomInfo = (params: GetLiveRoomInfoReq) => {
    return httpRequest.post<GetLiveRoomInfoRes>('/live/getLiveRoomInfo', params);
}

//Get a list of what's being streamed
export const getBeLiveList = () => {
    return httpRequest.post<GetBeLiveListRes>('/live/getBeLiveList');
}
