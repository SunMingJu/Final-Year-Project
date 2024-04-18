import { FileUpload } from "@/types/idnex"
export interface CreateCollectRmation extends FileUpload {
    id : number
    title: string, //title
    content: string //content
}

export interface SaveCreateFavoritesDataReq {
    id : number 
    type: string
    title: string
    content: string
    cover: string
}

//for single-title requests
export  interface createFavoritesReq {
    title : string
}