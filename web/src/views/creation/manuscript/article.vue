<template>
    <div class="overall">
        <pageHeader title="Column manuscript" icon-nmae="column" :animate="false" :whiteWhale="false"></pageHeader>
        <div class="content" v-loading="isLoading" :infinite-scroll-disabled="isTheEnd">
            <div class="article-list" v-show="articleList.length > 0" v-infinite-scroll="scrollBottom"
                infinite-scroll-delay="1000" :infinite-scroll-disabled="isTheEnd">
                <div :class="{ 'animate__animated': true, 'animate__fadeOutLeftBig': item.is_delete }"
                    v-for="(item, index) in articleList" :key="item.id" placement="top">
                    <div class="article-item">
                        <div class="item-left" @click="jump(item)">
                            <el-image class="item-img" :src="item.cover" fit="cover" />
                        </div>
                        <div class="item-right">
                            <div class="data">
                                <div class="item-title" @click="jump(item)">
                                    <VueEllipsis3 :visibleLine="1" :text="item.title">
                                    </VueEllipsis3>
                                </div>
                                <div class="item-info">
                                    <div class="icon-item">
                                        Number of reads : {{ item.heat }}
                                    </div>
                                </div>
                            </div>
                            <div class="function">
                                <el-button type="primary" v-removeFocus :icon="Edit" circle @click="editRecord(item)" />
                                <el-button type="info" v-removeFocus :icon="Delete" @click="delRecord(item.id)" circle />
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            <div class="record-empty" v-show="articleList.length == 0 && isLoading == false">
                <el-empty description="The column has not been published yet, please publish it soon~" />
            </div>
            <div class="load-more" v-show="articleList.length > 0 && isLoadMore" v-loading="true">
            </div>
            <!--Open the bottom -->
            <div class="spread-bottom">
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { useArticleProp, useDelRecord, useEditRecord, useJump, useLoadData } from '@/logic/creation/manuscript/article';
import { GetArticleManagementListItem } from '@/types/creation/manuscript/Article';
import { vRemoveFocus } from "@/utils/customInstruction/focus";
import { Delete, Edit } from '@element-plus/icons-vue';
import { watch } from 'vue';
import { VueEllipsis3 } from 'vue-ellipsis-3';


components: {
    VueEllipsis3
}

const { route, router, articleList, isLoading, pageInfo, loading, isLoadMore, isTheEnd, editArticleStore } = useArticleProp()

const delRecord = (id: number) => {
    useDelRecord(articleList, id)
}

const jump = (item: GetArticleManagementListItem) => {
    useJump(item, router)
}

const editRecord = (item: GetArticleManagementListItem) => {
    useEditRecord(item, loading, editArticleStore, router)
}

//load bottom
const scrollBottom = async () => {
    //Cancel loading more when there is no data
    isLoadMore.value = true
    if (articleList.value.length <= 0) return false
    useLoadData(articleList, isLoading, pageInfo, isTheEnd)
    isLoadMore.value = false
}

watch(() => route.path, async () => {
    articleList.value = []
    isLoading.value = true
    useLoadData(articleList, isLoading, pageInfo, isTheEnd)
}, { immediate: true, deep: true })

</script>

<style lang="scss" scoped>
@import "@/assets/styles/views/creation/manuscript/article.scss";
</style>
