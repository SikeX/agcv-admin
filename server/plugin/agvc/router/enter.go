package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/agvc/api"
)

var (
	Router     = new(router)
	apiDevice  = api.Api.Device
	apiHistory = api.Api.History
	apiAGC     = api.Api.AGC
	apiAVC     = api.Api.AVC
)

type router struct {
	Device  device
	History history
	AGC     agc
	AVC     avc
}
