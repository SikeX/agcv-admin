package agvc

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	AgvcBwdSettingApi
	AgvcQxySettingApi
}

var (
	agvcBwdSettingService = service.ServiceGroupApp.AgvcServiceGroup.AgvcBwdSettingService
	agvcQxySettingService = service.ServiceGroupApp.AgvcServiceGroup.AgvcQxySettingService
)
