import { GetHomeInfoReq, GetHomeInfoRes } from "@/types/home/home";
import httpRequest from "@/utils/requst"

//Get homepage information
export const getHomeInfo = (params: GetHomeInfoReq) => {
    return httpRequest.post<GetHomeInfoRes>('/home/getHomeInfo',params);
}