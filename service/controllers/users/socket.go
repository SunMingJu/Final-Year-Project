package users

import (
	"easy-video-net/logic/users/chatByUserSocket"
	"easy-video-net/logic/users/chatSocket"
	"easy-video-net/logic/users/noticeSocket"
	"easy-video-net/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"strconv"
)

// NoticeSocket  
func (us UserControllers) NoticeSocket(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	conn, _ := ctx.Get("conn")
	ws := conn.(*websocket.Conn)
	err := noticeSocket.CreateNoticeSocket(uid, ws)
	if err != nil {
		response.ErrorWs(ws, "Failed to create notification socket")
	}
}

// ChatSocket   
func (us UserControllers) ChatSocket(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	conn, _ := ctx.Get("conn")
	ws := conn.(*websocket.Conn)
	err := chatSocket.CreateChatSocket(uid, ws)
	if err != nil {
		response.ErrorWs(ws, "Failed to create chat socket")
	}
}

func (us UserControllers) ChatByUserSocket(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	conn, _ := ctx.Get("conn")
	//Determine if a video socket room is created
	tidQuery, _ := strconv.Atoi(ctx.Query("tid"))
	tid := uint(tidQuery)
	ws := conn.(*websocket.Conn)
	err := chatByUserSocket.CreateChatByUserSocket(uid, tid, ws)
	if err != nil {
		response.ErrorWs(ws, "Failed to create user chat socket")
	}
}
