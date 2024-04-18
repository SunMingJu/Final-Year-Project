<template>
    <div class="comment-box" v-loading="isLoading" :infinite-scroll-distance="40" v-infinite-scroll="scrollBottom"
        :infinite-scroll-disabled="isTheEnd" ref="scrollRef" :style="{ height: scrollHeight + 'px' }">
        <div class="comment-item" v-for="item in commentList" :key="item.id">
            <div class="item-left">
                <div class="avatar"><el-avatar :size="52" :src="item.photo" />
                </div>
                <div class="info">
                    <div class="top">
                        <div class="username"><span>{{ item.username }}</span></div>
                        <div class="time"><span>{{ dayjs(item.created_at).format('YYYY.MM.DD.hh.mm') }}</span></div>
                    </div>
                    <div class="comment-text">
                        <div class="comment-content">comments : </div>
                        <VueEllipsis3 :visibleLine="1" :text="item.comment">
                        </VueEllipsis3>
                    </div>
                </div>
            </div>
            <div class="item-right">
                <div class="video-info">
                    <el-image class="item-img" :src="item.cover" fit="cover" />
                    <div class="title">
                        <VueEllipsis3 :visibleLine="1" :text="item.title">
                        </VueEllipsis3>
                    </div>
                </div>
            </div>
        </div>

        <div class="load-more" v-show="commentList.length > 0 && isLoadMore" v-loading="true">
        </div>
        <!--Open the bottom -->
        <div class="spread-bottom">
        </div>
    </div>
</template>

<script lang="ts" setup>
import { GetDiscussVideoListReq, GetDiscussVideoListRes } from '@/types/creation/discuss/comment';
import Swal from 'sweetalert2';
import { nextTick, onMounted, ref } from 'vue';
import globalScss from "@/assets/styles/global/export.module.scss"
import { VueEllipsis3 } from 'vue-ellipsis-3';
import { PageInfo } from '@/types/idnex';
import { getDiscussVideoList } from '@/apis/contribution';
import dayjs from "dayjs"
import { Console } from 'console';

components: {
    VueEllipsis3
}

const commentList = ref(<GetDiscussVideoListRes>[])
const pageInfo = ref(<PageInfo>{
    page: 1,
    size: 9,
})
const scrollHeight = ref(0)
const scrollRef = ref()
//Whether it is loaded for the first time
const isLoading = ref(true)
//Is loading more
const isLoadMore = ref(false)
//Whether all loading is completed
const isTheEnd = ref(false)


const loadData = async () => {
    try {
        const response = await getDiscussVideoList(<GetDiscussVideoListReq>{
            page_info: pageInfo.value
        })
        if (!response.data) return false
        if (response.data.length == 0) isTheEnd.value = true
        commentList.value = [...commentList.value, ...response.data]
        pageInfo.value.page++
        console.log(response)

    } catch (err) {
        console.log(err)
        Swal.fire({
            title: "Failed to load data",
            heightAuto: false,
            confirmButtonColor: globalScss.colorButtonTheme,
            icon: "error",
        })
    }
}

//Load bottom
const scrollBottom = async () => {
    // //Cancel loading when no dataMore
    if (isLoading.value == true) return false
    isLoadMore.value = true
    if (commentList.value.length <= 0) return false
    await loadData()
    isLoadMore.value = false

}


onMounted(async () => {
    await loadData()
    isLoading.value = false
    nextTick(()=>{
        scrollHeight.value = document.documentElement.clientHeight - scrollRef.value.offsetTop - 20
    })
})


</script>

<style lang="scss" scoped>
@import "@/assets/styles/views/creation/discuss/comment/video.scss"
</style>