<template>
    <div class="overall">
        <pageHeader title="我的收藏" icon-nmae="collection"></pageHeader>
        <div class="collect-list" ref="scrollRef" :style="{ height: scrollHeight + 'px', overflow: isLoading ?  'scroll' : 'hidden'  }"
            v-show="isLoading == false || favoritesList.length > 0">
            <div class="collect-item" v-for="item in favoritesList.length > 0 ? favoritesList : quickCreationArr(7)">
                <!--Skeleton screen -->
                <el-skeleton style="width: 100%; height: 7.4rem;  margin:  6px 20px 20px 0px;" :loading="!item.id" animated>
                    <template #template>
                        <el-skeleton-item variant="text" style="width: 100%; height: 100%;" />
                        <div style="text-align: start; margin-top: 0.6rem;">
                        </div>
                    </template>
                    <template #default>
                        <div class="cover">
                            <div class="shadow-first"></div>
                            <div class="shadow-second"></div>
                            <el-image style="width: 100%; height: 100%" :src="item.cover" fit="cover" />
                        </div>
                        <div class="info">
                            <div class="title">{{ item.title }}</div>
                            <div class="creator">creator : {{ item.userInfo.username }} Views : 0</div>
                            <div class="btn">
                                <el-button type="primary" round size="small" @click="viewContent(item.id)"> View content </el-button>
                            </div>
                        </div>
                        <div class="function">
                            <el-popover placement="bottom" trigger="hover">
                                <template #reference>
                                    <el-icon class="more-filled" :size="20" color="#99a2aa">
                                        <MoreFilled />
                                    </el-icon>
                                </template>
                                <div style="text-align: center;">
                                    <el-button type="primary" round size="small" @click="updateFavorites(item)"> Revise
                                    </el-button>
                                    <el-button type="danger" round size="small" @click="delFavorites(item.id)"> delete
                                    </el-button>
                                </div>
                            </el-popover>
                        </div>
                    </template>
                </el-skeleton>
            </div>
            <div class="placeholder"></div>
        </div>
        <div class="collect-empty" v-show="favoritesList.length == 0 && isLoading == true">
            <el-empty description="No favorites created yet~" />
        </div>
        <div class="dialog">
            <el-dialog v-model="createCollectDialogShow" title="Create favorites" width="50%" center :close-on-click-modal="true"
                :append-to-body="true" :before-close="createCollectDialogShowClose" align-center>
                <CreateFavorites :info="updataInfo" :type="false" @shutDown="shutDown"></CreateFavorites>
            </el-dialog>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { MoreFilled } from '@element-plus/icons-vue'
import { useFavoritesProp, useInit, useDelFavorites, useUpdateFavorites } from "@/logic/personal/create/favorites"
import { onMounted, ref, watch } from 'vue';
import CreateFavorites from "@/components/createFavorites/createFavorites.vue"
import { GetFavoritesListItem } from '@/types/personal/collect/favorites';

comments: {
    CreateFavorites
}

const { route ,router, favoritesList, createCollectDialogShow, updataInfo } = useFavoritesProp()
const scrollHeight = ref(0)
const scrollRef = ref()
const isLoading = ref(false)


onMounted(() => {
    scrollHeight.value = document.documentElement.clientHeight - scrollRef.value.offsetTop
})

const viewContent = (id :number) => {
    router.push({name:"CollectList" , params : {
        "id" : id
    }})
}
const delFavorites = (id: number) => {
    useDelFavorites(id, favoritesList)
}

const updateFavorites = (item: GetFavoritesListItem) => {
    useUpdateFavorites(item, createCollectDialogShow, updataInfo)
}

const createCollectDialogShowClose = (done: () => void) => {
    done()
}

//Generate placeholder skeleton screen
const quickCreationArr = (num: number): Array<GetFavoritesListItem> => {
    let arr = []
    for (let i = 0; i < num; i++) {
        arr.push(<GetFavoritesListItem>{
        })
    }
    return arr
}

//Close create collect dialog show
const shutDown = () => {
    createCollectDialogShow.value = false
    //Close is for data update
    favoritesList.value = []
    isLoading.value = false
    setTimeout(() => {
        useInit(favoritesList, isLoading)
    }, 3000)
}


watch(() => route.path, async () => {
    favoritesList.value = []
    isLoading.value = false
    setTimeout(() => {
        useInit(favoritesList, isLoading)
    }, 3000)
}, { immediate: true, deep: true })

</script>

<style lang="scss" scoped>
@import "@/assets/styles/views/personal/collect/collect.scss";
</style>
