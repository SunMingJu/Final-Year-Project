import { danmakuApi, getVideoBarrageList, getVideoContributionByID, likeVideo, sendVideoBarrage } from "@/apis/contribution";
import globalScss from "@/assets/styles/global/export.module.scss";
import { useGlobalStore, useUserStore } from "@/store/main";
import { GetVideoBarrageListReq, GetVideoContributionByIDReq, LikeVideoReq, SendVideoBarrageReq, VideoInfo } from "@/types/show/video/video";
import DPlayer, { DPlayerDanmakuItem, DPlayerVideoQuality } from "dplayer";
import Swal from 'sweetalert2';
import { Ref, UnwrapNestedRefs, reactive, ref } from "vue";
import { RouteLocationNormalizedLoaded, Router, useRoute, useRouter } from "vue-router";
import { numberOfViewers, responseBarrageNum } from './socketFun';


export const useVideoProp = () => {
  const route = useRoute()
  const router = useRouter()
  const global = useGlobalStore()
  const userStore = useUserStore()
  const videoRef = ref()
  const videoID = ref<number>(0)
  const videoInfo = reactive(<VideoInfo>{})
  const barrageInput = ref("")
  const barrageListShow = ref(false)
  const videoBarrage = ref(true)
  const liveNumber = ref(0)
  //Reply to secondary comments
  const replyCommentsDialog = reactive({
    show: false,
    commentsID: 0,
  })

  return {
    route,
    router,
    userStore,
    videoRef,
    videoID,
    videoInfo,
    barrageInput,
    barrageListShow,
    liveNumber,
    replyCommentsDialog,
    videoBarrage,
    global
  }
}

export const useSendBarrage = async (text: Ref<string>, dpaler: DPlayer, userStore: any, videoInfo: UnwrapNestedRefs<VideoInfo>, socket: WebSocket) => {
  const res = await sendVideoBarrage(<SendVideoBarrageReq>{
    author: userStore.userInfoData.username,
    color: 16777215,
    id: videoInfo.videoInfo.id.toString(),
    text: text.value,
    time: dpaler.video.currentTime,
    type: 0,
    token: userStore.userInfoData.token,
  })

  console.log(userStore.userInfoData)
  if (res.code != 0) {
    Swal.fire({
      title: "Barrage service exception",
      heightAuto: false,
      confirmButtonColor: globalScss.colorButtonTheme,
      icon: "error",
    })
    return
  }
  const danmaku = <DPlayerDanmakuItem>{
    text: text.value,
    color: '#fff',
    type: 'right',
  };

  dpaler.danmaku.draw(danmaku);

  text.value = ""

  let data = JSON.stringify({
    type: "sendBarrage",
    data: ""
  })
  socket.send(data)

}

export const useLikeVideo = async (videoInfo: UnwrapNestedRefs<VideoInfo>) => {
  try {
    await likeVideo(<LikeVideoReq>{
      video_id: videoInfo.videoInfo.id
    })
    videoInfo.videoInfo.is_like = !videoInfo.videoInfo.is_like
  } catch (err) {
    Swal.fire({
      title: "Failed to like",
      heightAuto: false,
      confirmButtonColor: globalScss.colorButtonTheme,
      icon: "error",
    })
  }
}
export const useInit = async (videoRef: Ref, route: RouteLocationNormalizedLoaded, Router: Router, videoID: Ref<Number>, videoInfo: UnwrapNestedRefs<VideoInfo>, global: any) => {
  try {
    //Bind video id
    if (!route.params.id) {
      Router.back()
      Swal.fire({
        title: "Failed to get video",
        heightAuto: false,
        confirmButtonColor: globalScss.colorButtonTheme,
        icon: "error",
      })
      Router.back()
      return
    }
    global.globalData.loading.loading = true
    videoID.value = Number(route.params.id)
   //Get video information
    const vinfo = await getVideoContributionByID(<GetVideoContributionByIDReq>{
      video_id: videoID.value
    })
    if (!vinfo.data) return false
    videoInfo.videoInfo = vinfo.data.videoInfo
    videoInfo.recommendList = vinfo.data.recommendList

    //get clarity list
    let quality: DPlayerVideoQuality[] = []
    if (videoInfo.videoInfo.video) {
      quality = [...quality, {
        name: "1080P",
        url: videoInfo.videoInfo.video
      }]
    }
    if (videoInfo.videoInfo.video_720p) {
      quality = [...quality, {
        name: "720P",
        url: videoInfo.videoInfo.video_720p
      }]
    }
    if (videoInfo.videoInfo.video_480p) {
      quality = [...quality, {
        name: "408P",
        url: videoInfo.videoInfo.video_480p
      }]
    }
    if (videoInfo.videoInfo.video_360p) {
      quality = [...quality, {
        name: "360P",
        url: videoInfo.videoInfo.video_360p
      }]
    }

    //Get video barrage information
    const barrageList = await getVideoBarrageList(<GetVideoBarrageListReq>{
      id: videoID.value.toString()
    })
    if (!barrageList.data) return false
    videoInfo.barrageList = barrageList.data
    //Get current user information
const userStore = useUserStore()
    //Initialize the player
    const dp = new DPlayer({
      container: videoRef.value,
      loop: true, //loop playback
      lang: "zh-cn", //language
      logo: "",
      autoplay: true,
      danmaku: {
        id: videoID.value.toString(),
        api: danmakuApi,
        token: userStore.userInfoData.token
      },
      mutex: false, //Mutually exclusive, prevent multiple players from playing at the same time
      video: {
        quality: quality,
        defaultQuality: 0,
        url: "Leave blank", // Video link
        pic: videoInfo.videoInfo.cover
      },
    });
    global.globalData.loading.loading = false
    return dp
  } catch (err) {
    global.globalData.loading.loading = false
    console.log(err)
  }
}

export const useWebSocket = (userStore: any, videoInfo: UnwrapNestedRefs<VideoInfo>, Router: Router, liveNumber: Ref<number>) => {
  let socket: WebSocket
  const open = () => {
    console.log("websocket connection succeeded ")
  }
  const error = () => {
    console.error("websocket Connection failed")
  }
  const getMessage = async (msg: any) => {
    let data = JSON.parse(msg.data)
    console.log(data)
    switch (data.type) {
      case "numberOfViewers":
        numberOfViewers(liveNumber, data.data.people)
        break;
      case "responseBarrageNum":
        responseBarrageNum(videoInfo)
        break;
    }
  }

  if (typeof (WebSocket) === "undefined") {
    Swal.fire({
      title: "Your browser does not support sockets",
      heightAuto: false,
      confirmButtonColor: globalScss.colorButtonTheme,
      icon: "error",
    })
    Router.back()
    return
  } else {
    //Instantiate socket
    socket = new WebSocket(import.meta.env.VITE_SOCKET_URL + "/ws/videoSocket?token=" + userStore.userInfoData.token + "&videoID=" + videoInfo.videoInfo.id)
    //Listen to socket connection
    socket.onopen = open
    //Listen for socket error messages
    socket.onerror = error
    //Listen to socket messages
    socket.onmessage = getMessage
  }

  return socket
}
