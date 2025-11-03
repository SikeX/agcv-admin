package agvc_main

import (
	agcv_main "github.com/flipped-aurora/gin-vue-admin/server/api/v1/agvc_main"
)

var (
	Router     = new(router)
	apiDevice  = agcv_main.Api.Device
	apiHistory = agcv_main.Api.History
	apiAGC     = agcv_main.Api.AGC
	apiAVC     = agcv_main.Api.AVC
)

type router struct {
	Device  Device
	History History
	AGC     agc
	AVC     avc
}
