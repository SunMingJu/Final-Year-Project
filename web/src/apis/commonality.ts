import { GetFullPathOfImageRrq, GetuploadingDirReq, GetuploadingDirRes, GetUploadingMethodReq, GetUploadingMethodRes, GteossStsres, SearchReq, SearchRes, UploadCheckReq, UploadCheckRes, UploadMergeReq, UploadMergeRes } from "@/types/commonality/commonality";
import { FileSliceUpload, FileUpload } from "@/types/idnex";
import httpRequest from "@/utils/requst";

//get oss sts information
export const gteossSTS = () => {
    return httpRequest.post<GteossStsres>('/commonality/ossSTS');
}

//Get upload method
export function getuploadingMethod(params: GetUploadingMethodReq) {
    return httpRequest.post<GetUploadingMethodRes>('/commonality/uploadingMethod', params);
}

//Get save address
export function getuploadingDir(params: GetuploadingDirReq) {
    return httpRequest.post<GetuploadingDirRes>('/commonality/uploadingDir', params);
}

//Get the full path of the image
export function getFullPathOfImage(params: GetFullPathOfImageRrq) {
    return httpRequest.post<string>('/commonality/getFullPathOfImage', params);
}

//search function
export function search(params: SearchReq) {
    return httpRequest.post<SearchRes>('/commonality/search', params);
}

//Uploading files
export const uploadFile = (params: any, uploadConfig: FileUpload) => {
    return httpRequest.upload('/commonality/upload', params, uploadConfig);
}

//Uploading files in slices
export const UploadSliceFile = (params: any, uploadConfig: FileSliceUpload) => {
    return httpRequest.uploadSlice('/commonality/UploadSlice', params, uploadConfig);
}

//Upload file validation (validation operation)
export const uploadCheck = (params: UploadCheckReq) => {
    return httpRequest.post<UploadCheckRes>('/commonality/uploadCheck', params);
}

//Upload file validation (merge operation)
export const uploadMerge = (params: UploadMergeReq) => {
    return httpRequest.post<UploadMergeRes>('/commonality/uploadMerge', params);
}