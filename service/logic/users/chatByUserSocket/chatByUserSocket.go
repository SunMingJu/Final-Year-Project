package chatByUserSocket

import (
	"easy-video-net/consts"
	"easy-video-net/global"
	receive "easy-video-net/interaction/receive/socket"
	socketResponse "easy-video-net/interaction/response/socket"
	"easy-video-net/logic/users/chatSocket"
	userModel "easy-video-net/models/users"
	"easy-video-net/models/users/chat/chatList"
	"easy-video-net/models/users/notice"
	"easy-video-net/utils/response"
	"encoding/json"
	"github.com/gorilla/websocket"
)

type Engine struct {
	//video room
	UserMapChannel map[uint]*UserChannel

	Register     chan *UserChannel
	Cancellation chan *UserChannel
}

type ChanInfo struct {
	Type string
	Data interface{}
}

//UserChannel 
type UserChannel struct {
	UserInfo *userModel.User
	Tid      uint
	Socket   *websocket.Conn
	MsgList  chan ChanInfo
}

var Severe = &Engine{
	UserMapChannel: make(map[uint]*UserChannel, 10),
	Register:       make(chan *UserChannel, 10),
	Cancellation:   make(chan *UserChannel, 10),
}

// Start 
func (e *Engine) Start() {
	for {
		select {
		//registered event
		case registerMsg := <-e.Register:
			//Add Member
			e.UserMapChannel[registerMsg.UserInfo.ID] = registerMsg
			//Clear unread messages
			cl := new(chatList.ChatsListInfo)
			err := cl.UnreadEmpty(registerMsg.UserInfo.ID, registerMsg.Tid)
			//Add Online Record
			if _, ok := chatSocket.Severe.UserMapChannel[registerMsg.UserInfo.ID]; ok {
				//Chatting with someone online
				chatSocket.Severe.UserMapChannel[registerMsg.UserInfo.ID].ChatList[registerMsg.Tid] = registerMsg.Socket
			}
			if err != nil {
				global.Logger.Error("uid %d tid %d Failed to clear the number of unread messages", registerMsg.UserInfo.ID, registerMsg.Tid)
			}

		case cancellationMsg := <-e.Cancellation:
			//Delete member
			delete(e.UserMapChannel, cancellationMsg.UserInfo.ID)
			//Delete Online Record
			if _, ok := chatSocket.Severe.UserMapChannel[cancellationMsg.UserInfo.ID]; ok {
				//Chatting with someone online
				delete(chatSocket.Severe.UserMapChannel[cancellationMsg.UserInfo.ID].ChatList, cancellationMsg.Tid)
			}
		}
	}
}

func CreateChatByUserSocket(uid uint, tid uint, conn *websocket.Conn) (err error) {
	//Creating a UserChannel
	userChannel := new(UserChannel)
	//Binding ws
	userChannel.Socket = conn
	user := &userModel.User{}
	user.Find(uid)
	userChannel.UserInfo = user
	userChannel.Tid = tid
	userChannel.MsgList = make(chan ChanInfo, 10)

	Severe.Register <- userChannel

	go userChannel.Read()
	go userChannel.Writer()
	return nil
}

//Writer Listening for write data
func (lre *UserChannel) Writer() {
	for {
		select {
		case msg := <-lre.MsgList:
			response.SuccessWs(lre.Socket, msg.Type, msg.Data)
		}
	}
}

//Read retrieve data
func (lre *UserChannel) Read() {
	//Link broken for offline
	defer func() {
		Severe.Cancellation <- lre
		err := lre.Socket.Close()
		if err != nil {
			return
		}
	}()
	//Listening to business channels
	for {
		//Checking for Tonda ping passes
		lre.Socket.PongHandler()
		_, text, err := lre.Socket.ReadMessage()
		if err != nil {
			return
		}
		info := new(receive.Receive)
		if err = json.Unmarshal(text, info); err != nil {
			response.ErrorWs(lre.Socket, "Message formatting error")
		}
		switch info.Type {
		case "sendChatMsgText":
			sendChatMsgText(lre, lre.UserInfo.ID, lre.Tid, info)
		}
	}
}

func (lre *UserChannel) NoticeMessage(tp string) {
	//Get unread messages
	nl := new(notice.Notice)
	num := nl.GetUnreadNum(lre.UserInfo.ID)
	if num == nil {
		global.Logger.Error("The notification id is%d User unread message failure", lre.UserInfo.ID)
	}
	lre.MsgList <- ChanInfo{
		Type: consts.NoticeSocketTypeMessage,
		Data: socketResponse.NoticeMessageStruct{
			NoticeType: tp,
			Unread:     num,
		},
	}
}
