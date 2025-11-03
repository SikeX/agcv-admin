package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/agvc/service"
)

var (
	Api            = new(api)
	serviceDevice  = service.Service.Device
	serviceHistory = service.Service.History
	serviceAGC     = service.Service.AGC
	serviceAVC     = service.Service.AVC
)

type api struct {
	Device  device
	History history
	AGC     agc
	AVC     avc
}

// GetById 通用ID查询结构
type GetById = request.GetById
