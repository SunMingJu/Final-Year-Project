<template>
    <div class="overall">
        <pageHeader title="my collection" icon-nmae="collection"></pageHeader>
        <div class="content">
            <div class="video-list">
                <div class="title">
                    <SvgIcon name="video" class="icon"></SvgIcon>video works
                </div>
                <div class="vidoe-box" v-show="releaseInformation?.videoList?.length > 0 || isLoading == false">
                    <div class="video-item"
                        v-for="videoInfo in releaseInformation?.videoList?.length ? releaseInformation.videoList : quickCreationArr(16)">
                        <el-skeleton style="width: 100%; height: 154px;" class="video-card"
                            :loading="!releaseInformation?.videoList?.length" animated>
                            <template #template>
                                <el-skeleton-item variant="text" style="  width: 100%;height: 100%;" />
                                <div style="text-align: start; margin-top: 0.6rem;">
                                    <el-skeleton-item variant="h3" style="width: 100%" />
                                    <div>
                                        <el-skeleton-item variant="h3" style="width: 80%" />
                                        <el-skeleton-item variant="h3" style="width: 60%" />
                                    </div>
                                </div>
                            </template>
                            <template #default>
                                <VideoCard :id="videoInfo.id" :title="videoInfo.title"
                                    :cover="videoInfo.cover" :video_duration="videoInfo.video_duration" :created_at="videoInfo.created_at">
                                </VideoCard>
                            </template>
                        </el-skeleton>
                    </div>
                </div>
                <div class="vidoe-empty" v-show="releaseInformation?.videoList?.length == 0 && isLoading == true">
                    <el-empty description="No favorite videos yet~" />
                </div>
            </div>
            <!--Open the bottom -->
            <div class="spread-bottom">
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>

import VideoCard from "@/components/collectListVideoCard/videoCard.vue"
import Column from "@/components/spaceCard/columnCard.vue";
import { onMounted } from "vue";
import { VideoInfo } from "@/types/space/space";
import { useCollectListProp ,useInit} from "@/logic/personal/create/collectList";

const {  route , releaseInformation , isLoading } = useCollectListProp()

//Generate placeholder skeleton screen
const quickCreationArr = (num: number): Array<VideoInfo> => {
    let arr = []
    for (let i = 0; i < num; i++) {
        arr.push(<VideoInfo>{
        })
    }
    return arr
}

components: {
    Column
    VideoCard
}


onMounted(() => {
    setTimeout(() => {
        useInit(route, releaseInformation, isLoading)
    }, 1500);
})

</script>

<style lang="scss" scoped>
@import  "@/assets/styles/views/personal/collect/collectList.scss";
</style>
