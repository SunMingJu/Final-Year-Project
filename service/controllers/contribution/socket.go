package contribution

import (
	"simple-video-net/logic/contribution/socket"
	"simple-video-net/utils/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// socket
func (c Controllers) socket(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	conn, _ := ctx.Get("conn")
	ws := conn.(*websocket.Conn)
	//Determine if a video socket room is created
	id, _ := strconv.Atoi(ctx.Query("videoID"))
	videoID := uint(id)
	if socket.Severe.VideoRoom[videoID] == nil {
		//Unwatched Active Creation
		socket.Severe.VideoRoom[videoID] = make(socket.UserMapChannel, 10)
	}
	err := socket.Createsocket(uid, videoID, ws)
	if err != nil {
		response.ErrorWs(ws, "Failed to create socket")
	}
}
