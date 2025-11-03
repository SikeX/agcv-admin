package agcv_main

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service/agvc/agcv_main"
)

var (
	Api            = new(api)
	serviceDevice  = agcv_main.Service.Device
	serviceHistory = agcv_main.Service.History
	serviceAGC     = agcv_main.Service.AGC
	serviceAVC     = agcv_main.Service.AVC
)

type api struct {
	Device  device
	History history
	AGC     agc
	AVC     avc
}

// GetById 通用ID查询结构
type GetById = request.GetById
