package main

import (
	_ "simple-video-net/global/database/mysql"
	_ "simple-video-net/global/database/redis"
	"simple-video-net/logic/contribution/videoSocket"
	liveSocket "simple-video-net/logic/live/socket"
	"simple-video-net/logic/users/chatByUserSocket"
	"simple-video-net/logic/users/chatSocket"
	"simple-video-net/logic/users/noticeSocket"
	"simple-video-net/router"
	"simple-video-net/utils/testing"
)

func main() {
	//Inspection of live streaming services
	testing.LiveSeverTesting()
	//Enable live streaming and video sockets
	go liveSocket.Severe.Start()
	go videoSocket.Severe.Start()
	go noticeSocket.Severe.Start()
	go chatSocket.Severe.Start()
	go chatByUserSocket.Severe.Start()
	//ces
	router.InitRouter()

}
