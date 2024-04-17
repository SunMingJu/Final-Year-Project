package socket

import (
	"github.com/gorilla/websocket"
)

//Writer Listening for write data
func (lre LiveRoomEvent) Writer() {
	for {
		select {
		case msg := <-lre.Channel.MsgList:
			err := lre.Channel.Socket.WriteMessage(websocket.BinaryMessage, msg)
			if err != nil {
				return
			}
		}
	}
}
