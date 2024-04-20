<template>
  <div class="mian">
    <div class="head">
      <topNavigation color="#fff" scroll :displaySearch="true"></topNavigation>
    </div>
    <!--Cover image -->
    <div class="cover-picture">
    </div>
    <!--Top channel popular categories -->
    <homeHeaderChannel @click="notOpen()"></homeHeaderChannel>
    <!--Subject -->
    <div class="middle" :infinite-scroll-distance="770" v-infinite-scroll="scrollBottom" :record-empty="40"
      :infinite-scroll-delay="1000">
      <!--Carousel frame screen -->
      <div class="rotograph-skeleton">
        <el-skeleton style="width: 100%; height: 70%;" class="video-card" :loading="!homeInfo.rotograph.length" animated>
          <template #template>
            <el-skeleton-item variant="text" style="  width: 100%;height: 100%;" />
            <div style="text-align: start; margin-top: 0.3rem;">
              <el-skeleton-item variant="text" style="width: 100% ;height: 5rem;" />
            </div>
          </template>
          <template #default>
            <!--Carousel image -->
            <homeRotograph :rotograph="homeInfo?.rotograph"></homeRotograph>
          </template>
        </el-skeleton>
      </div>

      <!--Video -->
      <!--Video skeleton screen -->
      <div class="video-card" v-for="(videoInfo, index) in videoList.length ? videoList : quickCreationArr(11)"
        :key="videoInfo.id">
        <el-skeleton style="width: 100%; height: 13rem;" class="video-card" :loading="!videoInfo.id" animated>
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
            <Card :id="videoInfo.id" :title="videoInfo.title" :username="videoInfo.username"
              :video_duration="videoInfo.video_duration" v-bind:barrage-number="videoInfo.barrageNumber"
              :heat="videoInfo.heat" :cover="videoInfo.cover" :created_at="videoInfo.created_at"></Card>
          </template>
        </el-skeleton>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { getHomeInfo } from "@/apis/home";
import globalScss from "@/assets/styles/global/export.module.scss";
import homeHeaderChannel from "@/components/homeHeaderChannel/homeHeaderChannel.vue";
import homeRotograph from "@/components/homeRotograph/homeRotograph.vue";
import Card from "@/components/homeVideoList/card.vue";
import topNavigation from "@/components/topNavigation/topNavigation.vue";
import { GetHomeInfoReq, GetHomeInfoRes, VideoInfo } from "@/types/home/home";
import Swal from "sweetalert2";
import { Ref, UnwrapNestedRefs, computed, onMounted, reactive, ref } from "vue";
components: {
  homeRotograph
  Card
  homeHeaderChannel
  topNavigation
}

const pageInfo = reactive(<GetHomeInfoReq>{
  page_info: {
    page: 1,
    size: 15,
  }
})


let homeInfo = ref<GetHomeInfoRes>({
  rotograph: [],
  videoList: []
})

const videoList = computed(() => {
  let list = [] as Array<VideoInfo>
  //Judge the current page number to be the second page
  if (pageInfo.page_info.page == 2) {
    if (homeInfo.value.videoList.length % 11 == 0) {
      list = [...list, ...homeInfo.value.videoList, ...quickCreationArr(15)]
      //and load data
      loadData(homeInfo, pageInfo)
    } else {
      list = [...homeInfo.value.videoList]
    }
  } else if (pageInfo.page_info.page > 2) {
    if ((homeInfo.value.videoList.length - 11) % pageInfo.page_info.size == 0) {
      list = [...list, ...homeInfo.value.videoList, ...quickCreationArr(15)]
    } else {
      list = [...homeInfo.value.videoList]
    }
  } else {
    list = [...homeInfo.value.videoList]
  }

  return list
})


//Generate placeholder skeleton screen
const quickCreationArr = (num: number): Array<VideoInfo> => {
  let arr = []
  for (let i = 0; i < num; i++) {
    arr.push(<VideoInfo>{
    })
  }
  return arr
}

const loadData = async (homeInfo: Ref<GetHomeInfoRes>, pageInfo: UnwrapNestedRefs<GetHomeInfoReq>) => {
  const response = await getHomeInfo(pageInfo)
  if (!response.data) return
  homeInfo.value.rotograph = response.data.rotograph
  homeInfo.value.videoList = [...homeInfo.value.videoList, ...response.data.videoList]
//Next paging +1 after successful request
  pageInfo.page_info.page = pageInfo.page_info.page + 1
}

let timer: NodeJS.Timeout | null = null

//load bottom
const scrollBottom = () => {
  console.log("bottom")
  //Cancel loading more when there is no data
  if (homeInfo.value.videoList.length <= 0) return false
  loadData(homeInfo, pageInfo)

}

const notOpen = () => {
  Swal.fire({
    title: "Stay tuned",
    heightAuto: false,
    confirmButtonColor: globalScss.colorButtonTheme,
    icon: "info",
  })
}

onMounted(() => {
  loadData(homeInfo, pageInfo)
})

</script>

<style scoped lang="scss">
@import "@/assets/styles/views/home/home.scss";
</style>
