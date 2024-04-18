<template>
    <div class="column">
        <div class="column-list" v-show="!isLoading || columnList.length > 0" v-infinite-scroll="scrollBottom"
            infinite-scroll-delay="1000">
            <!--skeleton screen-->
            <el-skeleton style="width: 100%; height: 18rem; margin-bottom: 8rem; " class="video-card"
                v-for="(item, index) in columnList.length ? columnList : quickCreationArr(6) " :key="item.id"
                :loading="!columnList.length" animated>
                <template #template>
                    <el-skeleton-item variant="text" style="  width: 100%;height: 100%;" />
                    <div style="text-align: start; margin-top: 0.2rem;">
                        <el-skeleton-item variant="h3" style="width: 100%" />
                        <div>
                            <el-skeleton-item variant="h3" style="width: 100%" />
                        </div>
                    </div>
                </template>
                <template #default>
                    <!--single card-->
                    <div :class="{ mouseover: item.is_stay, mouseleave: !item.is_stay }" class="column-item shadow-box "
                        @mouseover="mouseOver(item)" @mouseleave="mouseleave(item)" @click="jumpArticle(item.id)">
                        <div :class="{ 'item-image': true, 'right': index % 2 == 1 }">
                            <img :src="item.cover" class="el-image__inner image" style="object-fit: cover;">
                        </div>
                        <div class="item-text">
                            <div class="post-meta">
                                <SvgIcon name="camera" class="icon-small"></SvgIcon> " posted on
                                {{ dayjs(item.created_at).format('YYYY.MM.DD.hh.mm') }} "
                            </div>
                            <h3>{{ item.title }}</h3>
                            <div class="post-meta" style="margin-bottom: 15px;">
                                <SvgIcon name="hot" class="icon-small"></SvgIcon>
                                <span>
                                    {{ item.heat }} heat
                                </span>
                                <SvgIcon name="comments" class="icon-small"></SvgIcon>
                                <span>
                                    {{ item.comments_number }} comments
                                </span>
                                <SvgIcon name="like" class="icon-small"></SvgIcon>
                                <span>
                                    {{ item.likes_number }}like
                                </span>
                            </div>
                            <div class="recent-post-desc">
                                <VueEllipsis3 :text="item.content" :visibleLine="4">
                                    <template v-slot:ellipsisNode>
                                        <span>...</span>
                                    </template>
                                </VueEllipsis3>
                            </div>
                            <div class="sort-label">
                                <div class="lable-item" style="margin-right: 12px;">
                                    <SvgIcon name="classification" class="icon-small" /> {{ item.classification }}
                                </div>
                                <div class="lable-item" style="margin-right: 12px;" v-for="label in item.label" :key="label"
                                    v-show="label">
                                    <SvgIcon name="label" class="icon-small" /> {{ label }}
                                </div>
                            </div>
                        </div>
                    </div>
                </template>
            </el-skeleton>
        </div>
        <div class="load-more" v-show="isLoadMore" v-loading="isLoadMore">
        </div>
        <!--Open the bottom-->
        <div class="no-more" v-show="isTheEnd">
            no more~
        </div>
        <div class="spread-bottom">
        </div>
    </div>
    <div class="column-empty" v-show="columnList.length == 0 && isLoading == true">
        <el-empty description="No one has posted a column yet, please post it soon~" />
    </div>
</template>

<script setup lang="ts">

import { getArticleContributionList } from "@/apis/contribution";
import { GetArticleContributionListReq } from "@/types/home/column";
import { PageInfo } from '@/types/idnex';
import { GetArticleContributionListByUserRes, GetArticleContributionListByUserResItem } from "@/types/live/liveRoom";
import dayjs from "dayjs";
import { onMounted, ref } from 'vue';
import { VueEllipsis3 } from 'vue-ellipsis-3';
import { useRouter } from 'vue-router';

components: {
    VueEllipsis3
}

//Column list
const columnList = ref<GetArticleContributionListByUserRes>([])
const router = useRouter()
const pageInfo = ref(<PageInfo>{
    page: 1,
    size: 6
})
//Whether to load for the first time
const isLoading = ref(false)
//Is loading more
const isLoadMore = ref(false)
//Whether all loading is completed
const isTheEnd = ref(false)


//Load bottom
const scrollBottom = async () => {
    if (!isLoading.value) return false
    if (isTheEnd.value) return false
    //Cancel loading when no dataMore
    if (columnList.value.length <= 0) return false
    isLoadMore.value = true
    await loadData()

}

const mouseOver = (item: GetArticleContributionListByUserResItem) => {
    item.is_stay = true
}
const mouseleave = (item: GetArticleContributionListByUserResItem) => {
    item.is_stay = false
}

const loadData = async () => {
    setTimeout(async () => {
        try {
            const response = await getArticleContributionList(<GetArticleContributionListReq>{
                page_info: pageInfo.value
            })
            if (!response.data) return false
            if (response.data.length == 0) isTheEnd.value = true
            //Add whether to stay the mouse
            response.data?.filter((item) => {
                item.is_stay = false
                return true
            })
            columnList.value = [...columnList.value, ...response.data]
            console.log(columnList)
            //Next paging +1 after successful request
            pageInfo.value.page++
            isLoading.value = true
            isLoadMore.value = false
        } catch (err) {
            console.log(err)
        }
    }, 500);

}

const quickCreationArr = (num: number): Array<GetArticleContributionListByUserResItem> => {
    let arr = []
    for (let i = 0; i < num; i++) {
        arr.push(<GetArticleContributionListByUserResItem>{

        })
    }
    return arr
}

const jumpArticle = (articleID: number) => {
    router.push({ name: "ArticleShow", query: { articleID } })
}
onMounted(() => {
    loadData()
})
</script>
<style scoped lang="scss">
@import "./static/style/column.scss"
</style>