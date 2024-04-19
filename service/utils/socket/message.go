package socket

import (
	"simple-video-net/logic/contribution/sokcet"
	liveSocket "simple-video-net/logic/live/socket"
	"simple-video-net/logic/users/chat"
	"simple-video-net/logic/users/chatUser"
	"simple-video-net/logic/users/notice"
)

func init() {
	//Initialize all sockets
	go liveSocket.Severe.Start()
	go sokcet.Severe.Start()
	go notice.Severe.Start()
	go chat.Severe.Start()
	go chatUser.Severe.Start()
}
