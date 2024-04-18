import {FileUpload} from "@/types/idnex"
export  interface  LiveInformation extends FileUpload{
    title : string, //title
}

export interface SaveLiveDataReq {
    type : string
    title : string
    imgUrl :string
} 

export interface GetLiveDataRes {
    title : string
    img :string
} 