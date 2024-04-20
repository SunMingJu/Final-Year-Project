import { FileUpload } from "@/types/idnex"
//form form structure
export interface vdeoContributionForm {
    id: number | undefined,
    isShow: boolean,
    title: string,
    type: boolean,
    timing: boolean,
    date1time: string,
    labelInputVisible: boolean,
    labelText: string,
    label: Array<string>,
    introduce: string,
    videoDuration: number
}
//upload
export interface uploadFileformation extends FileUpload {
    media?: string
}

//api createVideoContribution need structure
export interface CreateVideoContributionReq {
    id: number,
    video: string,
    videoUploadType: string,
    cover: string,
    coverUploadType: string,
    title: string,
    reprinted: boolean,
    label: Array<string>,
    introduce: string,
    videoDuration: number,
    media: string
}

export interface UpdateVideoContributionReq {
    id: number,
    cover: string,
    coverUploadType: string,
    title: string,
    reprinted: boolean,
    label: Array<string>,
    introduce: string,
}