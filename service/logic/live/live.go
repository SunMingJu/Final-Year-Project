package live

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"simple-video-net/global"
	receive "simple-video-net/interaction/receive/live"
	response "simple-video-net/interaction/response/live"
	"simple-video-net/models/users"
	"simple-video-net/models/users/record"
	"strconv"
	"strings"
)

func GetLiveRoom(uid uint) (results interface{}, err error) {
	//Request Live Server
	url := global.Config.LiveConfig.Agreement + "://" + global.Config.LiveConfig.IP + ":" + global.Config.LiveConfig.Api + "/control/get?room="
	url = url + strconv.Itoa(int(uid))
	// Create http get request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	// Parses the body data from the http request into the structure we defined.
	ReqGetRoom := new(receive.ReqGetRoom)
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(ReqGetRoom); err != nil {
		return nil, err
	}
	if ReqGetRoom.Status != 200 {
		return nil, fmt.Errorf("Failed to get live address")
	}
	return response.GetLiveRoomResponse("rtmp://"+global.Config.LiveConfig.IP+":"+global.Config.LiveConfig.RTMP+"/live", ReqGetRoom.Data), nil
}

func GetLiveRoomInfo(data *receive.GetLiveRoomInfoReceiveStruct, uid uint) (results interface{}, err error) {
	userInfo := new(users.User)
	userInfo.FindLiveInfo(data.RoomID)
	flv := global.Config.LiveConfig.Agreement + "://" + global.Config.LiveConfig.IP + ":" + global.Config.LiveConfig.FLV + "/live/" + strconv.Itoa(int(data.RoomID)) + ".flv"

	if uid > 0 {
		//Add History
		rd := new(record.Record)
		err = rd.AddLiveRecord(uid, data.RoomID)
		if err != nil {
			return nil, fmt.Errorf("Failed to add history")
		}
	}
	return response.GetLiveRoomInfoResponse(userInfo, flv), nil
}

func GetBeLiveList() (results interface{}, err error) {
	//Takes the id of the user who opens the playlist.
	url := global.Config.LiveConfig.Agreement + "://" + global.Config.LiveConfig.IP + ":" + global.Config.LiveConfig.Api + "/stat/livestat"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	// Parses the body data from the http request into the structure we defined.
	livestat := new(receive.LivestatRes)
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(livestat); err != nil {
		return nil, fmt.Errorf("parsing failure")
	}
	if livestat.Status != 200 {
		return nil, fmt.Errorf("Failed to get live list")
	}
	//Get the list of active values in live
	keys := make([]uint, 0)
	for _, kv := range livestat.Data.Publishers {
		ka := strings.Split(kv.Key, "live/")
		uintKey, _ := strconv.ParseUint(ka[1], 10, 19)
		keys = append(keys, uint(uintKey))
	}
	userList := new(users.UserList)
	if len(keys) > 0 {
		err = userList.GetBeLiveList(keys)
		if err != nil {
			return nil, fmt.Errorf("Enquiry Failure")
		}
	}
	return response.GetBeLiveListResponse(userList), nil
}
