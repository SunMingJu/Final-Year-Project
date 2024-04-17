package socket

import (
	"easy-video-net/consts"
	"easy-video-net/global"
	userModel "easy-video-net/models/users"
	"easy-video-net/proto/pb"
	"easy-video-net/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type Engine struct {
	//liveroom
	LiveRoom map[uint]UserMapChannel

	Register     chan LiveRoomEvent
	Cancellation chan LiveRoomEvent
}

type UserMapChannel map[uint]*UserChannel

//UserChannel 
type UserChannel struct {
	UserInfo userModel.User
	Socket   *websocket.Conn
	MsgList  chan []byte
}

//LiveRoomEvent 
type LiveRoomEvent struct {
	RoomID  uint
	Channel *UserChannel
}

var Severe = &Engine{
	LiveRoom:     make(map[uint]UserMapChannel, 10),
	Register:     make(chan LiveRoomEvent, 10),
	Cancellation: make(chan LiveRoomEvent, 10),
}

// Start 
func (e *Engine) Start() {
	//Create live sockets for each user
	type userList []userModel.User
	users := new(userList)
	global.Db.Select("id").Find(users)
	for _, v := range *users {
		e.LiveRoom[v.ID] = make(UserMapChannel, 10)
	}
	//Listening to business channel information
	//Register offline event listener
	for {
		select {
		//registered event
		case registerMsg := <-e.Register:
			logrus.Infof("registered event %v", registerMsg)
			//No room to launch straight away
			if _, ok := e.LiveRoom[registerMsg.RoomID]; !ok {
				//Formatting the response
				message := &pb.Message{
					MsgType: consts.Error,
					Data:    []byte("Message formatting error"),
				}
				res, _ := proto.Marshal(message)
				_ = registerMsg.Channel.Socket.WriteMessage(websocket.BinaryMessage, res)
				return
			}
			//add member
			e.LiveRoom[registerMsg.RoomID][registerMsg.Channel.UserInfo.ID] = registerMsg.Channel
			//Broadcast users on-line
			err := serviceOnlineAndOfflineRemind(registerMsg, true)
			if err != nil {
				response.ErrorWs(registerMsg.Channel.Socket, err.Error())
			}
			//Send a history message to the user
			err = serviceResponseLiveRoomHistoricalBarrage(registerMsg)
			if err != nil {
				response.ErrorWs(registerMsg.Channel.Socket, err.Error())
			}

		case cancellationMsg := <-e.Cancellation:
			logrus.Infof("occurrence of an offline event %v", cancellationMsg)
			delete(e.LiveRoom[cancellationMsg.RoomID], cancellationMsg.Channel.UserInfo.ID)
			//Broadcast users offline
			err := serviceOnlineAndOfflineRemind(cancellationMsg, false)
			if err != nil {
				response.ErrorWs(cancellationMsg.Channel.Socket, err.Error())
			}
		}
	}
}

func CreateSocket(ctx *gin.Context, userId uint, roomID uint, conn *websocket.Conn) (err error) {
	//Creating a UserChannel
	userChannel := new(UserChannel)
	//Binding ws
	userChannel.Socket = conn
	//Binding user information
	user := userModel.User{}
	user.Find(userId)
	userChannel.UserInfo = user
	//Preventing blockage
	userChannel.MsgList = make(chan []byte, 10)

	//Create User
	userLiveEvent := LiveRoomEvent{
		RoomID:  roomID,
		Channel: userChannel,
	}
	Severe.Register <- userLiveEvent

	go userLiveEvent.Read()
	go userLiveEvent.Writer()
	return nil
}
