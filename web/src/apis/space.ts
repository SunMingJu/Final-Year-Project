import { GetAttentionListReq, GetAttentionListRes, GetReleaseInformationReq, GetReleaseInformationRes, GetSpaceIndividualReq, GetSpaceIndividualRes, GetVermicelliListReq, GetVermicelliListRes } from "@/types/space/space";
import httpRequest from "@/utils/requst"
//Getting personal space information
export const getSpaceIndividual = (params: GetSpaceIndividualReq) => {
    return httpRequest.post<GetSpaceIndividualRes>('/space/getSpaceIndividual', params);
}

//Get videos and columns posted in the space
export const getReleaseInformation = (params: GetReleaseInformationReq) => {
    return httpRequest.post<GetReleaseInformationRes>('/space/getReleaseInformation', params);
}

//Get Followed List
export const getAttentionList = (params: GetAttentionListReq) => {
    return httpRequest.post<GetAttentionListRes>('/space/getAttentionList', params);
}

//Get Fans List
export const getVermicelliList = (params: GetVermicelliListReq) => {
    return httpRequest.post<GetVermicelliListRes>('/space/getVermicelliList', params);
}