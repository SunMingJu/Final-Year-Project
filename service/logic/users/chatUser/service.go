package chatUser

import (
	"simple-video-net/consts"
	"simple-video-net/global"
	receive "simple-video-net/interaction/receive/socket"
	"simple-video-net/interaction/response/socket"
	"simple-video-net/logic/users/chat"
	"simple-video-net/models/users/chat/chatList"
	"simple-video-net/models/users/chat/chatMsg"
	"simple-video-net/utils/conversion"
	"simple-video-net/utils/response"
)

func sendChatMsgText(ler *UserChannel, uid uint, tid uint, info *receive.Receive) {

	//Add Message Log
	cm := chatMsg.Msg{
		Uid:     uid,
		Tid:     tid,
		Type:    info.Type,
		Message: info.Data,
	}
	err := cm.AddMessage()
	if err != nil {
		response.ErrorWs(ler.Socket, "Send Failure")
		return
	}
	//Message Enquiry
	msgInfo := new(chatMsg.Msg)
	err = msgInfo.FindByID(cm.ID)
	if err != nil {
		response.ErrorWs(ler.Socket, "Failed to send message")
		return
	}
	photo, _ := conversion.FormattingJsonSrc(msgInfo.UInfo.Photo)

	//Messaging yourself without pushing
	if uid == tid {
		return
	}

	if _, ok := chat.Severe.UserMapChannel[tid]; ok {
		//Online situation
		if _, ok := chat.Severe.UserMapChannel[tid].ChatList[uid]; ok {
			//In the chat window with yourself (direct push)
			response.SuccessWs(chat.Severe.UserMapChannel[tid].ChatList[uid], consts.ChatSendTextMsg, socket.ChatSendTextMsgStruct{
				ID:        msgInfo.ID,
				Uid:       msgInfo.Uid,
				Username:  msgInfo.UInfo.Username,
				Photo:     photo,
				Tid:       msgInfo.Tid,
				Message:   msgInfo.Message,
				Type:      msgInfo.Type,
				CreatedAt: msgInfo.CreatedAt,
			})
			return
		} else {
			//Adding unread records
			cl := new(chatList.ChatsListInfo)
			err := cl.UnreadAutocorrection(tid, uid)
			if err != nil {
				global.Logger.Error("uid %d tid %d 消息记录自增未读消息数量失败", tid, uid)
			}
			ci := new(chatList.ChatsListInfo)
			_ = ci.FindByID(uid, tid)
			//Push the main socket
			response.SuccessWs(chat.Severe.UserMapChannel[tid].Socket, consts.ChatUnreadNotice, socket.ChatUnreadNoticeStruct{
				Uid:         uid,
				Tid:         tid,
				LastMessage: ci.LastMessage,
				LastMessageInfo: socket.ChatSendTextMsgStruct{
					ID:        msgInfo.ID,
					Uid:       msgInfo.Uid,
					Username:  msgInfo.UInfo.Username,
					Photo:     photo,
					Tid:       msgInfo.Tid,
					Message:   msgInfo.Message,
					Type:      msgInfo.Type,
					CreatedAt: msgInfo.CreatedAt,
				},
				Unread: cl.Unread,
			})
		}
	} else {
		//offline
		cl := new(chatList.ChatsListInfo)
		err := cl.UnreadAutocorrection(tid, uid)
		if err != nil {
			global.Logger.Error("uid %d tid %d Failed to increment the number of unread messages in the message log", tid, uid)
		}
	}
}
