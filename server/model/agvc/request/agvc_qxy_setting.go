package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type AgvcQxySettingSearch struct {
	Number     string `json:"number" form:"number"`
	DeviceName string `json:"deviceName" form:"deviceName"`
	request.PageInfo
}
