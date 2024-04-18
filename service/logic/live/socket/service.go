package socket

import (
	"fmt"
	"simple-video-net/consts"
	"simple-video-net/global"
	"simple-video-net/proto/pb"
	"simple-video-net/utils/conversion"
	"strconv"

	"google.golang.org/protobuf/proto"
)

// Send a pop-up message
func serviceSendBarrage(lre LiveRoomEvent, text []byte) error {
	barrageInfo := &pb.WebClientSendBarrageReq{}
	if err := proto.Unmarshal(text, barrageInfo); err != nil {
		return fmt.Errorf("Message formatting error")
	}
	src, _ := conversion.FormattingJsonSrc(lre.Channel.UserInfo.Photo)
	response := &pb.WebClientSendBarrageRes{
		UserId:   float32(lre.Channel.UserInfo.ID),
		Username: lre.Channel.UserInfo.Username,
		Avatar:   src,
		Text:     barrageInfo.Text,
		Color:    barrageInfo.Color,
		Type:     barrageInfo.Type,
	}
	data, err := proto.Marshal(response)
	if err != nil {
		return fmt.Errorf("Message formatting error")
	}
	//Save pop-ups to recent messages
	str := conversion.Bytes2String(data)
	if studentLen, _ := global.RedisDb.LLen(consts.LiveRoomHistoricalBarrage + strconv.Itoa(int(lre.RoomID))).Result(); studentLen >= 10 {
		err := global.RedisDb.RPop(consts.LiveRoomHistoricalBarrage + strconv.Itoa(int(lre.RoomID))).Err()
		if err != nil {
			return err
		}
	}
	//Less than 20 messages Insert directly
	err = global.RedisDb.LPush(consts.LiveRoomHistoricalBarrage+strconv.Itoa(int(lre.RoomID)), str).Err()
	if err != nil {
		global.Logger.Errorf("Room ID： %d Recent failures to deposit pop-ups into Redis messages： %s", lre.RoomID, data)
		return err
	}
	//Formatting the response
	message := &pb.Message{
		MsgType: consts.WebClientBarrageRes,
		Data:    data,
	}
	res, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("Message formatting error")
	}
	for _, v := range Severe.LiveRoom[lre.RoomID] {
		v.MsgList <- res
	}

	return nil
}

// User offline alerts
func serviceOnlineAndOfflineRemind(lre LiveRoomEvent, isOnlineOndOffline bool) error {
	//Get all current users
	type userListStruct []*pb.EnterLiveRoom
	userList := make(userListStruct, 0)
	src, _ := conversion.FormattingJsonSrc(lre.Channel.UserInfo.Photo)
	for _, v := range Severe.LiveRoom[lre.RoomID] {
		itemSrc, _ := conversion.FormattingJsonSrc(v.UserInfo.Photo)
		item := &pb.EnterLiveRoom{
			UserId:   float32(v.UserInfo.ID),
			Username: v.UserInfo.Username,
			Avatar:   itemSrc,
		}
		userList = append(userList, item)
	}

	for i := 0; i < len(Severe.LiveRoom[lre.RoomID]); i++ {
	}
	response := &pb.WebClientEnterLiveRoomRes{
		UserId:   float32(lre.Channel.UserInfo.ID),
		Username: lre.Channel.UserInfo.Username,
		Avatar:   src,
		Type:     isOnlineOndOffline,
		List:     userList,
	}

	//response output
	data, err := proto.Marshal(response)
	if err != nil {
		return fmt.Errorf("Message formatting error")
	}
	message := &pb.Message{
		MsgType: consts.WebClientEnterLiveRoomRes,
		Data:    data,
	}
	res, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("Message formatting error")
	}
	for _, v := range Severe.LiveRoom[lre.RoomID] {
		v.MsgList <- res
	}
	return nil
}

// Responding to history message pop-ups
func serviceResponseLiveRoomHistoricalBarrage(lre LiveRoomEvent) error {
	//Get the history.
	val, err := global.RedisDb.LRange(consts.LiveRoomHistoricalBarrage+strconv.Itoa(int(lre.RoomID)), 0, -1).Result()

	if err != nil {
		return fmt.Errorf("Failed to get history popups")
	}
	historicalBarrage := &pb.WebClientHistoricalBarrageRes{}
	list := make([]*pb.WebClientSendBarrageRes, 0)
	for _, v := range val {
		barrage := &pb.WebClientSendBarrageRes{}
		if err := proto.Unmarshal(conversion.String2Bytes(v), barrage); err != nil {
			return fmt.Errorf("Message formatting error")
		}
		list = append(list, barrage)
	}
	historicalBarrage.List = list
	data, err := proto.Marshal(historicalBarrage)
	if err != nil {
		return fmt.Errorf("Message formatting error")
	}
	//Formatting the response
	message := &pb.Message{
		MsgType: consts.WebClientHistoricalBarrageRes,
		Data:    data,
	}
	res, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("Message formatting error")
	}
	for _, v := range Severe.LiveRoom[lre.RoomID] {
		v.MsgList <- res
	}

	return nil
}
