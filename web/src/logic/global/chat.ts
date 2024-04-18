import globalScss from "@/assets/styles/global/export.module.scss";
import { useChatListStore } from "@/store/chat";
import { useUserStore } from "@/store/main";
import { ResultDataWs } from "@/types/idnex";
import { ChatOnlineUnreadNotice, ChatUnreadNotic } from "@/types/socket/chat";
import { ChatInfo, MessageInfo } from "@/types/store/chat";
import Swal from "sweetalert2";
import { useRouter } from "vue-router";



export const useInitChatSocket = () => {
    const router = useRouter()
    const userStore = useUserStore()
    console.log(useChatListStore())
    let socket: WebSocket
   const open = () => {
        console.log("Chat websocket connection successful")
    }
    const error = () => {
        console.error("Chat websocket connection failed")
    }
    const getMessage = async (msg: any) => {
        let data = <ResultDataWs>JSON.parse(msg.data)
        switch (data.type) {
            case "error":
                console.error("Chat socket returns error")
                break;
            case "chatUnreadNotice":
                data as ResultDataWs<ChatUnreadNotic>
                chatUnreadNotice(data.data)
                break;
            case "chatOnlineUnreadNotice":
                data as ResultDataWs<ChatUnreadNotic>
                chatOnlineUnreadNotice(data.data)
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
        router.back()
        return
    } else {
        //Instantiate socket
        socket = new WebSocket(import.meta.env.VITE_SOCKET_URL + "/ws/chatSocket?token=" + userStore.userInfoData.token)
        //Listen to socket connection
        socket.onopen = open
        //Listen for socket error messages
        socket.onerror = error
        //Listen to socket messages
        socket.onmessage = getMessage
    }
}


//Message processing function
const chatUnreadNotice = (data: ChatUnreadNotic) => {
    const chatListStore = useChatListStore()
    chatListStore.chatListData = chatListStore.chatListData.filter((item: ChatInfo) => {
        if (item.to_id == data.uid) {
            item.updated_at = data.last_message_info.created_at
            item.last_message = data.last_message
            item.unread = data.unread
            //Add to message list
            chatListStore.addMessage(data.last_message_info.uid, <MessageInfo>{
                uid: data.last_message_info.uid,
                username: data.last_message_info.username,
                photo: data.last_message_info.photo,
                tid: data.last_message_info.tid,
                message: data.last_message_info.message,
                type: data.last_message_info.type,
                created_at: data.last_message_info.created_at
            })
        }
        return item
    })
}

const chatOnlineUnreadNotice = (data: ChatOnlineUnreadNotice) => {
    const chatListStore = useChatListStore()
    try {
        const chatList = data.filter((item) => {
            item.message_list = item.message_list.reverse()
            return item
        })
        chatListStore.chatListData = chatList
    } catch (err: any) {
        console.error(err)
    }
}