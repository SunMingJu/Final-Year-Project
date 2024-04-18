<template>
    <div class="overall">
        <pageHeader title="history record" icon-nmae="playRecording"></pageHeader>
        <div class="content" v-loading="isLoading">
            <div class="timeline" v-show="recordList.length > 0" v-infinite-scroll="scrollBottom"
                infinite-scroll-delay="1000">
                <el-timeline>
                    <el-timeline-item
                        :class="{ 'animate__animated': true, 'animate__fadeOutLeftBig': item.is_delete }"
                        v-for="(item, index) in recordList" :key="item.id" :timestamp="recordTimeFormat(item.updated_at)"
                        placement="top">
                        <div class="timeline-item" >
                            <div class="item-left" @click="jump(item)">
                                <el-image class="item-img" :src="item.cover" fit="cover" />
                            </div>
                            <div class="item-right" >
                                <div class="data">
                                    <div class="item-title" @click="jump(item)">
                                        <VueEllipsis3 :visibleLine="1" :text="item.title">
                                        </VueEllipsis3>
                                    </div>
                                    <div class="item-info">
                                        <el-avatar :size="38" :src="item.photo" />
                                        <span class="username">{{ item.username }}</span>
                                        <span class="partition">|</span>
                                        <span class="type">{{ item.type }}</span>
                                    </div>
                                </div>
                                <div class="function">
                                    <el-button type="info" v-removeFocus :icon="Delete" @click="delRecord(item.id)"
                                        circle />
                                </div>
                            </div>
                        </div>
                    </el-timeline-item>
                </el-timeline>
            </div>
            <div class="record-empty" v-show="recordList.length == 0 && isLoading == false">
                <el-empty description="There is no history record yet. Go and have a look.~" />
            </div>
            <div class="load-more" v-show="recordList.length > 0 && isLoadMore" v-loading="true">
            </div>
            <!--Open the bottom -->
            <div class="spread-bottom">
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { VueEllipsis3 } from 'vue-ellipsis-3';
import { Delete } from '@element-plus/icons-vue'
import { useRecordProp, useLoadData, useDelRecord, useJump } from '@/logic/personal/record/record';
import {  watch } from 'vue';
import { recordTimeFormat } from "@/utils/conversion/timeConversion"
import { vRemoveFocus } from "@/utils/customInstruction/focus"
import { GetRecordListItem } from '@/types/personal/record/record';


components: {
    VueEllipsis3
}

const { route, router, recordList, isLoading, pageInfo, loading, isLoadMore, isTheEnd } = useRecordProp()

const delRecord = (id: number) => {
    useDelRecord(recordList, id, loading)
}

const jump = (item: GetRecordListItem) => {
    useJump(item, router)
}

//load bottom
const scrollBottom = async () => {
    //Cancel loading more when there is no data
    if (isTheEnd.value) return false
    isLoadMore.value = true
    if (recordList.value.length <= 0) return false
    useLoadData(recordList, isLoading, pageInfo, isTheEnd)
    isLoadMore.value = false
}

watch(() => route.path, async () => {
    recordList.value = []
    isLoading.value = true
    setTimeout(() => {
        useLoadData(recordList, isLoading, pageInfo, isTheEnd)
    }, 2000)
}, { immediate: true, deep: true })

</script>

<style lang="scss" scoped>
@import "@/assets/styles/views/personal/record/record.scss";
</style>
