package videoSocket

import (
	"encoding/json"
	"simple-video-net/consts"
	receive "simple-video-net/interaction/receive/socket"
	userModel "simple-video-net/models/users"
	"simple-video-net/utils/response"

	"github.com/gorilla/websocket"
)

type Engine struct {
	//VideoRoom
	VideoRoom map[uint]UserMapChannel

	Register     chan VideoRoomEvent
	Cancellation chan VideoRoomEvent
}

type UserMapChannel map[uint]*UserChannel

type ChanInfo struct {
	Type string
	Data interface{}
}

// UserChannel
type UserChannel struct {
	UserInfo userModel.User
	Socket   *websocket.Conn
	MsgList  chan ChanInfo
}

// VideoRoomEvent
type VideoRoomEvent struct {
	VideoID uint
	Channel *UserChannel
}

var Severe = &Engine{
	VideoRoom:    make(map[uint]UserMapChannel, 10),
	Register:     make(chan VideoRoomEvent, 10),
	Cancellation: make(chan VideoRoomEvent, 10),
}

// Start
func (e *Engine) Start() {
	for {
		select {
		//registered event
		case registerMsg := <-e.Register:
			//Add Member
			e.VideoRoom[registerMsg.VideoID][registerMsg.Channel.UserInfo.ID] = registerMsg.Channel
			//Number of online viewers of broadcasts
			num := len(e.VideoRoom[registerMsg.VideoID])
			r := struct {
				People int `json:"people"`
			}{
				People: num,
			}
			res := ChanInfo{
				Type: consts.VideoSocketTypeNumberOfViewers,
				Data: r,
			}
			for _, v := range e.VideoRoom[registerMsg.VideoID] {
				v.MsgList <- res
			}
		case cancellationMsg := <-e.Cancellation:
			//Number of online viewers of broadcasts
			num := len(e.VideoRoom[cancellationMsg.VideoID]) - 1
			r := struct {
				People int `json:"people"`
			}{
				People: num,
			}
			res := ChanInfo{
				Type: consts.VideoSocketTypeNumberOfViewers,
				Data: r,
			}
			for _, v := range e.VideoRoom[cancellationMsg.VideoID] {
				v.MsgList <- res
			}
			delete(e.VideoRoom[cancellationMsg.VideoID], cancellationMsg.Channel.UserInfo.ID)
		}
	}
}

func CreateVideoSocket(userID uint, videoID uint, conn *websocket.Conn) (err error) {
	//Creating a UserChannel
	userChannel := new(UserChannel)
	//Binding ws
	userChannel.Socket = conn
	//Binding user information
	user := userModel.User{}
	user.Find(userID)
	userChannel.UserInfo = user
	//Preventing blockage
	userChannel.MsgList = make(chan ChanInfo, 10)

	//Create User
	userLiveEvent := VideoRoomEvent{
		VideoID: videoID,
		Channel: userChannel,
	}
	Severe.Register <- userLiveEvent

	go userLiveEvent.Read()
	go userLiveEvent.Writer()
	return nil

}

// Writer Listening for write data
func (lre VideoRoomEvent) Writer() {
	for {
		select {
		case msg := <-lre.Channel.MsgList:
			response.SuccessWs(lre.Channel.Socket, msg.Type, msg.Data)
		}
	}
}

// Read
func (lre VideoRoomEvent) Read() {
	//Link broken for offline
	defer func() {
		Severe.Cancellation <- lre
		err := lre.Channel.Socket.Close()
		if err != nil {
			return
		}
	}()
	//Listening to business channels
	for {
		//Checking for Tonda ping passes
		lre.Channel.Socket.PongHandler()
		_, text, err := lre.Channel.Socket.ReadMessage()
		if err != nil {
			return
		}
		info := new(receive.Receive)
		if err = json.Unmarshal(text, info); err != nil {
			response.ErrorWs(lre.Channel.Socket, "Message formatting error")
		}
		switch info.Type {

		}
	}
}
