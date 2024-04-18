import {  VideoInfo } from "@/types/show/video/video";

import { Ref, UnwrapNestedRefs } from "vue";

//Number of online viewers
export const numberOfViewers = (liveNumber: Ref<number>, people: number) => {
    liveNumber.value = people
}

//Someone sent a barrage
export const responseBarrageNum = (info:  UnwrapNestedRefs<VideoInfo>)  => {
    info.videoInfo.barrageNumber++
}
