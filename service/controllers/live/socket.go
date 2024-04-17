package live

import (
	"easy-video-net/consts"
	"easy-video-net/logic/live/socket"
	"easy-video-net/proto/pb"
	"easy-video-net/utils/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

func (lv LivesControllers) LiveSocket(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	conn, _ := ctx.Get("conn")
	ws := conn.(*websocket.Conn)

	//Determine whether to create a live room
	liveRoom, _ := strconv.Atoi(ctx.Query("liveRoom"))
	liveRoomID := uint(liveRoom)
	if socket.Severe.LiveRoom[liveRoomID] == nil {
		message := &pb.Message{
			MsgType: consts.Error,
			Data:    []byte("The live stream is not on"),
		}
		res, _ := proto.Marshal(message)
		_ = ws.WriteMessage(websocket.BinaryMessage, res)
		return
	}

	err := socket.CreateSocket(ctx, uid, liveRoomID, ws)
	if err != nil {
		response.ErrorWs(ws, err.Error())
		return
	}
}
