package contribution

import (
	"simple-video-net/logic/contribution/videoSocket"
	"simple-video-net/utils/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// VideoSocket
func (c Controllers) VideoSocket(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	conn, _ := ctx.Get("conn")
	ws := conn.(*websocket.Conn)
	//Determine if a video socket room is created
	id, _ := strconv.Atoi(ctx.Query("videoID"))
	videoID := uint(id)
	if videoSocket.Severe.VideoRoom[videoID] == nil {
		//Unwatched Active Creation
		videoSocket.Severe.VideoRoom[videoID] = make(videoSocket.UserMapChannel, 10)
	}
	err := videoSocket.CreateVideoSocket(uid, videoID, ws)
	if err != nil {
		response.ErrorWs(ws, "Failed to create socket")
	}
}
