<template>
    <div id="chat-box" class="chat-box" ref="boxRef" @scroll="boxScroll">
        <div class="chat-list">
            <div class="msg-list">
                <div v-for="msgItem in msgList" :key="msgItem.id">
                    <div class="msg-left" v-if="userStore.userInfoData.id != msgItem.uid">
                        <el-avatar shape="circle" :size="34" :src="msgItem.photo" />
                        <div class="msg-text">
                            <span class="msg-span">{{ msgItem.message }}</span>
                        </div>
                    </div>
                    <div class="msg-right" v-if="userStore.userInfoData.id == msgItem.uid">
                        <div class="msg-text">
                            <span class="msg-span">{{ msgItem.message }}</span>
                        </div>
                        <el-avatar shape="circle" :size="34" :src="userStore.userInfoData.photo" />
                    </div>
                    <div></div>
                </div>
            </div>
        </div>
        <div class="chat-input">
            <el-input v-model="input" resize="none" :input-style="{ width: '400px' }" autosize type="textarea"
                placeholder="What are we talking about~" />
            <el-button v-removeFocus @click="sendText(input)" class="send" :type="input ? 'primary' : 'info'" :icon="Check"
                circle />
        </div>
    </div>
</template>

<script lang="ts" setup>
import { getChatHistoryMsg } from "@/apis/personal";
import globalScss from "@/assets/styles/global/export.module.scss";
import { useChatListStore } from "@/store/chat";
import { useUserStore } from "@/store/main";
import { ResultDataWs } from "@/types/idnex";
import { GetChatHistoryMsgReq } from "@/types/personal/chat/chat";
import { ChatSendTextMsg } from "@/types/socket/chat";
import { MessageInfo } from '@/types/store/chat';
import { vRemoveFocus } from "@/utils/customInstruction/focus";
import { Check } from '@element-plus/icons-vue';
import { ElMessage } from "element-plus";
import Swal from "sweetalert2";
import { nextTick, onMounted, onUnmounted, ref, watch } from 'vue';
import { useRouter } from 'vue-router';

const props = defineProps({
    tid: {
        type: Number,
        required: true,
    },
    msgList: {
        type: Array as () => Array<MessageInfo> | undefined,
        required: true,
    }
})

var socket: WebSocket | undefined
const input = ref("")
const userStore = useUserStore()
const chatListStore = useChatListStore()
const router = useRouter()
const tid = ref(0)
const boxRef = ref()

const loadSocket = () => {
    let socket: WebSocket
    const open = () => {
        console.log("User chat websocket connection successful ")
    }
    const error = (err: any) => {
        console.log(err)
        console.error("User chat chat websocket connection failed")
    }
    const getMessage = async (msg: any) => {
        let data = <ResultDataWs>JSON.parse(msg.data)
        switch (data.type) {
            case "error":
                console.error("User chat chat socket returns error")
                ElMessage({
                    message: data.message,
                    type: 'error',
                    appendTo: document.getElementById("chat-box") as HTMLElement,
                })
                break;
            case "chatSendTextMsg":
                data as ResultDataWs<ChatSendTextMsg>
                chatSendTextMsg(data.data)
                break
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
        // Instantiate socket
        socket = new WebSocket(import.meta.env.VITE_SOCKET_URL + "/ws/chatByUserSocket?token=" + userStore.userInfoData.token + "&tid=" + tid.value)
        // Listen for socket connections
        socket.onopen = open
        // Listen for socket error messages
        socket.onerror = error
        // Listen for socket messages
        socket.onmessage = getMessage
    }
    return socket
}

const send = (type: string, msg: string) => {
    let data = JSON.stringify({
        type,
        "data": msg,
    })
    socket?.send(data)
}

const sendText = (msg: string) => {
    if (msg == "") return false
    send("sendChatMsgText", msg)
    //Add record
    input.value = ""
    chatListStore.addMessage(tid.value, <MessageInfo>{
        uid: userStore.userInfoData.id,
        username: userStore.userInfoData.username,
        photo: userStore.userInfoData.photo,
        tid: tid.value,
        message: msg,
        type: "text",
        created_at: Date().toString()
    })
    //scroll to bottom
    nextTick(() => {
        rollingBottom()
    })
}

const chatSendTextMsg = (data: ChatSendTextMsg) => {
    chatListStore.addMessage(data.uid, data)
    nextTick(() => {
        rollingBottom()
    })

}


const rollingBottom = () => {
    //Executed
    if (props.msgList) {
        nextTick(() => {
            if (boxRef?.value?.scrollHeight) {
                boxRef.value.scrollTop = boxRef?.value?.scrollHeight
            }

        })
    }

}

const boxScroll = async () => {
    console.log("touch top")

    if (boxRef.value.scrollTop == 0) {
        //Touch top to load more’
        const h = boxRef.value.scrollHeight
        try {
            let mixTime: number | string = new Date().getTime() //Minimum value defaults to current time
            chatListStore.chatListData.filter((item) => {
                if (item.to_id == chatListStore.tid) {
                    item.message_list.filter((ml) => {
                        console.log(new Date(ml.created_at).getTime() - new Date(mixTime).getTime())
                        if (new Date(ml.created_at).getTime() - new Date(mixTime).getTime() < 0) {
                            mixTime = ml.created_at
                        }
                    })
                }
            })
            console.log("minimum value", mixTime)
            const response = await getChatHistoryMsg(<GetChatHistoryMsgReq>{
                tid: props.tid,
                last_time: mixTime
            })
            if (!response.data) return false
            const chatList = response.data.reverse()
            chatListStore.chatListData = chatListStore.chatListData.filter((item) => {
                if (item.to_id == chatListStore.tid) {
                    item.message_list = [...chatList, ...item.message_list]
                }
                return item
            })
            nextTick(() => {
                boxRef.value.scrollTop = boxRef.value.scrollHeight - h
            })

            console.log(response)
        } catch (err) {
            console.log(err)
        }
    }
}

const watchTid = watch(() => { return chatListStore.tid }, (newVal, oldVal) => {
    socket?.close()
    if (newVal != oldVal) {
        tid.value = newVal
        chatListStore.clearUnread(tid.value)
        socket = loadSocket()
        rollingBottom()
    }
}, { immediate: true })



onMounted(() => {

})

//End listening and socket when closing
onUnmounted(() => {
    watchTid()
    socket?.close()
})

</script>

<style lang="scss" scoped>
@import "./static/style/chatBox.scss";
</style>