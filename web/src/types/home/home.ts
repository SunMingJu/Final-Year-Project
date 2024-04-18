import { type } from "os"
import { PageInfo } from "../idnex"

//getLiveRoom obtains the live broadcast interface request requirements
export interface GetLiveRoomRes {
    address: string,
    key: string
}
  
//Carousel image
export interface Rotograph {
    title: string
    cover:string
    color: string
    type: string
    to_id: number
}
export type RotographList = Array<Rotograph>

//Video information
export interface VideoInfo {
    id: number
uid:number
title : string
video:string
cover: string
video_duration : number
label : Array<string>
introduce :  string
	heat : number 
	barrageNumber : number
	username : string
    created_at : string
}

export type VideoInfoList = Array<VideoInfo>



export interface GetHomeInfoReq {
    page_info: PageInfo
}

export interface GetHomeInfoRes {
    rotograph: RotographList
    videoList : VideoInfoList
}