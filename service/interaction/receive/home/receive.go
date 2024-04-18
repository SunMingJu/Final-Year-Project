package home

import "simple-video-net/models/common"

type GetHomeInfoReceiveStruct struct {
	PageInfo common.PageInfo `json:"page_info" binding:"required"`
}
