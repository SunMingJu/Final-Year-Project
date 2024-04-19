package users

import (
	"simple-video-net/logic/users/chat"
	"simple-video-net/logic/users/chatUser"
	"simple-video-net/logic/users/notice"
	"simple-video-net/utils/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// notice
func (us UserControllers) notice(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	conn, _ := ctx.Get("conn")
	ws := conn.(*websocket.Conn)
	err := notice.Createnotice(uid, ws)
	if err != nil {
		response.ErrorWs(ws, "Failed to create notification socket")
	}
}

// chatUser
func (us UserControllers) chatUser(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	conn, _ := ctx.Get("conn")
	ws := conn.(*websocket.Conn)
	err := chatUser.CreatechatUser(uid, ws)
	if err != nil {
		response.ErrorWs(ws, "Failed to create chat socket")
	}
}

func (us UserControllers) chat(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	conn, _ := ctx.Get("conn")
	//Determine if a video socket room is created
	tidQuery, _ := strconv.Atoi(ctx.Query("tid"))
	tid := uint(tidQuery)
	ws := conn.(*websocket.Conn)
	err := chat.Createchat(uid, tid, ws)
	if err != nil {
		response.ErrorWs(ws, "Failed to create user chat socket")
	}
}
