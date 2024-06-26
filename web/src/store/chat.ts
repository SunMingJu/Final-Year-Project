import { ChatList, MessageInfo } from "@/types/store/chat"
import { timetoRFC3339 } from "@/utils/conversion/timeConversion"
import { defineStore } from "pinia"
import { ref } from "vue"

export const useChatListStore = defineStore("caht", () => {
    const isShow = ref(false)
    const chatListData = ref(<ChatList>[])
    const tid = ref(0)
    const tUsername = ref("")

    const setChatListData = (info: ChatList) => {
        chatListData.value = info
    }

    //Delete Record
    const deleteItemByid = (tid: number) => {
        chatListData.value = chatListData.value.filter((item) => {
            return item.to_id != tid
        })
    }

    //Clear unread messages
    const clearUnread = (tid: number) => {
        chatListData.value = chatListData.value.filter((item) => {
            if (item.to_id == tid) {
                item.unread = 0
            }
            return item
        })
    }
    //Add message
    const addMessage = (tid: number, info: MessageInfo) => {
        chatListData.value = chatListData.value.filter((item) => {
            if (item.to_id == tid) {
                //Replace the last record 
                item.last_message = info.message
                item.last_at = timetoRFC3339(new Date())
                item.message_list = [...item.message_list, info]

            }
            return item
        })
    }


    return {
        isShow,
        chatListData,
        setChatListData,
        tid,
        tUsername,
        deleteItemByid,
        clearUnread,
        addMessage
    }
}, {
    //Only persist chat related data
    persist: [
        {
            paths: ['chatListData'],
            storage: localStorage,
        }
    ]
})