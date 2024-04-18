package noticeSocket

import (
	"encoding/json"
	"simple-video-net/consts"
	"simple-video-net/global"
	receive "simple-video-net/interaction/receive/socket"
	socketResponse "simple-video-net/interaction/response/socket"
	userModel "simple-video-net/models/users"
	"simple-video-net/models/users/notice"
	"simple-video-net/utils/response"

	"github.com/gorilla/websocket"
)

type Engine struct {
	//Video Room
	UserMapChannel map[uint]*UserChannel

	Register     chan *UserChannel
	Cancellation chan *UserChannel
}

type ChanInfo struct {
	Type string
	Data interface{}
}

// UserChannel
type UserChannel struct {
	UserInfo *userModel.User
	Socket   *websocket.Conn
	MsgList  chan ChanInfo
}

var Severe = &Engine{
	UserMapChannel: make(map[uint]*UserChannel, 10),
	Register:       make(chan *UserChannel, 10),
	Cancellation:   make(chan *UserChannel, 10),
}

// Start Starting services
func (e *Engine) Start() {
	for {
		select {
		//registered event
		case registerMsg := <-e.Register:
			//Add Member
			e.UserMapChannel[registerMsg.UserInfo.ID] = registerMsg

		case cancellationMsg := <-e.Cancellation:
			//Delete member
			delete(e.UserMapChannel, cancellationMsg.UserInfo.ID)
		}
	}
}

func CreateNoticeSocket(uid uint, conn *websocket.Conn) (err error) {
	//Creating a UserChannel
	userChannel := new(UserChannel)
	//Binding ws
	userChannel.Socket = conn
	user := &userModel.User{}
	user.Find(uid)
	userChannel.UserInfo = user
	userChannel.MsgList = make(chan ChanInfo, 10)

	Severe.Register <- userChannel

	go userChannel.Read()
	go userChannel.Writer()
	return nil

}

// Writer Listening for write data
func (lre *UserChannel) Writer() {
	for {
		select {
		case msg := <-lre.MsgList:
			response.SuccessWs(lre.Socket, msg.Type, msg.Data)
		}
	}
}

// Read retrieve data
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

		}
	}
}

func (lre *UserChannel) NoticeMessage(tp string) {
	//Get unread messages
	nl := new(notice.Notice)
	num := nl.GetUnreadNum(lre.UserInfo.ID)
	if num == nil {
		global.Logger.Error("The notification id is%dUser unread message failure", lre.UserInfo.ID)
	}
	lre.MsgList <- ChanInfo{
		Type: consts.NoticeSocketTypeMessage,
		Data: socketResponse.NoticeMessageStruct{
			NoticeType: tp,
			Unread:     num,
		},
	}
}
