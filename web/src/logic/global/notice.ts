import globalScss from "@/assets/styles/global/export.module.scss";
import { useUserStore } from "@/store/main";
import Swal from "sweetalert2";
import { useRouter } from "vue-router";



export const useInitNoticeSocket = () => {
    const router = useRouter()
    const userStore = useUserStore()
    let socket: WebSocket
    const open = () => {
        console.log("Notify websocket connection is successful ")
    }
    const error = () => {
        console.error("Notify websocket connection failed")
    }
    const getMessage = async (msg: any) => {
        let data = JSON.parse(msg.data)
        switch (data.type) {
            case "error":
                console.error("Notification socket returns error")
                break;
            case "messageNotice":
                messageNotice(data.data.unread, userStore)
                break;
        }
    }

    if (typeof (WebSocket) === "undefined") {
        Swal.fire({
            title: "Your browser does not support socket",
            heightAuto: false,
            confirmButtonColor: globalScss.colorButtonTheme,
            icon: "error",
        })
        router.back()
        return
    } else {
        //Instantiate socket
        socket = new WebSocket(import.meta.env.VITE_SOCKET_URL + "/ws/noticeSocket?token=" + userStore.userInfoData.token)
        //Listen to socket connection
        socket.onopen = open
        //Listen for socket error messages
        socket.onerror = error
        //Listen to socket messages
        socket.onmessage = getMessage
    }
}


//Message processing function

const messageNotice = (num: number, userStore: any) => {
    userStore.setUnreadNotice(num)
}