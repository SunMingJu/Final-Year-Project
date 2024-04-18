
import DPlayer, { DPlayerDanmaku } from "dplayer";
import Swal from 'sweetalert2'
import globalScss from "@/assets/styles/global/export.module.scss"
import { nextTick, reactive, Ref, ref, UnwrapNestedRefs } from "vue";
import flvJs from "flv.js";
import { decodeMessage } from "@/proto/pb/live"
import { webClientBarrageDeal, webClientEnterLiveRoomDeal, webClientHistoricalBarrageRes, webError } from "@/logic/live/socketFun"
import { useRoute, RouteLocationNormalizedLoaded, useRouter, Router } from "vue-router"
import { useUserStore } from "@/store/main";
import { getLiveRoomInfo } from "@/apis/live";
import { GetLiveRoomInfoReq, GetLiveRoomInfoRes } from "@/types/live/liveRoom";
export const useLiveRoomProp = () => {
  const route = useRoute()
  const router = useRouter()
  const videoRef = ref()
  const userStore = useUserStore()
  const roomID = ref<number>(0)
  const liveInfo = reactive(<GetLiveRoomInfoRes>{})

  return {
    route,
    router,
    videoRef,
    userStore,
    roomID,
    liveInfo
  }
}

export const useWebSocket = (dp: DPlayer, userStore: any, sideRef: Ref<any>, roomID: Ref<Number>, Router: Router) => {
  let socket: WebSocket
  const initWebSocket = (() => {
    const open = () => {
      console.log("websocket connection successful")
    }
    const error = () => {
      console.error("websocket connection failed")
    }
    const getMessage = async (msg: any) => {
      console.log(msg.data)
      //Escape uint8ARRAY
      const reader = new FileReader();
      reader.readAsArrayBuffer(msg.data);
      reader.onload = function (e) {
        let buf = new Uint8Array(reader.result as ArrayBuffer);
        const response = decodeMessage(buf)
        switch (response.msgType) {
          case "error" :
            webError(response)
            break;
          case "webClientBarrageRes":
            webClientBarrageDeal(response, dp, sideRef)
            break;
          case "webClientEnterLiveRoomRes":
            webClientEnterLiveRoomDeal(response, dp, sideRef)
            break;
          case "webClientHistoricalBarrageRes":
            webClientHistoricalBarrageRes(response, dp, sideRef)
            break;
          default: console.error("Unsupported message type")
            break;
        }
        console.log(response)
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
      socket = new WebSocket( import.meta.env.VITE_SOCKET_URL + "/ws/liveSocket?token=" + userStore.userInfoData.token + "&liveRoom=" + roomID.value)
      //Listen to socket connection
      socket.onopen = open
      //Listen for socket error messages
      socket.onerror = error
      //Listen to socket messages
      socket.onmessage = getMessage
    }
  })()

  const sendMessage = (msg: string | ArrayBufferLike | Blob | ArrayBufferView) => {
    console.log("Send a message", msg)
    socket.send(msg)
  }

  return {
    sendMessage
  }

}


export const useInit = async (videoRef: Ref, route: RouteLocationNormalizedLoaded, Router: Router, roomID: Ref<Number>,liveInfo : UnwrapNestedRefs<GetLiveRoomInfoRes>) => {
  try {
    //Bind room
    if (!route.query.roomID) {
      Router.back()
      Swal.fire({
        title: "Failed to access room",
        heightAuto: false,
        confirmButtonColor: globalScss.colorButtonTheme,
        icon: "error",
      })
      Router.back()
      return
    }
    roomID.value = Number(route.query.roomID)

    //Get live broadcast information

    const li = await getLiveRoomInfo(<GetLiveRoomInfoReq>{
      room_id :  roomID.value
    })

    if(!li.data){
      Swal.fire({
        title: "This user has not configured live broadcast",
        heightAuto: false,
        confirmButtonColor: globalScss.colorButtonTheme,
        icon: "error",
      })
      return
    }
    liveInfo.live_title = li.data?.live_title
    liveInfo.photo = li.data?.photo
    liveInfo.username = li.data?.username
    liveInfo.flv = li.data?.flv



    //Initialize player
    console.log(videoRef)
    const dp = new DPlayer({
     container: videoRef.value, //container
      loop: true, //loop playback
      lang: "zh-cn", //Language, optional 'en', 'zh-cn', 'zh-tw',
      logo: "", //Display a logo in the upper left corner
      autoplay: true,
      danmaku: true as unknown as DPlayerDanmaku, //The official document gives true, but the types specified in ts are inconsistent.
      apiBackend: {
        read: function (options) {
          console.log('Pretend to connect WebSocket');
          options.success([]);
        },
        send: function (options) {
          console.log('Pretend to send danmaku via WebSocket', options.data);
          options.success();
        },
      },
      mutex: false, //Mutually exclusive, prevent multiple players from playing at the same time
      video: { //Video information
        type: "customFlv", //Video type optional "auto", "hls", "flv", "dash"..
        url: liveInfo.flv, //video link
        customType: {
          customFlv: (video: any, player: any) => {
            const flvPlayer = flvJs.createPlayer({
              type: "flv",
              isLive: true,
              hasAudio: true,
              url: liveInfo.flv,
            });
            nextTick(() =>{
              flvPlayer.attachMediaElement(video);
              flvPlayer.load();
            })
          }
        },
      },
    });
    return dp
  } catch (err) {
    console.log(err)
  }
}
