package agvc

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
    AgvcBwdSettingApi
    AgvcBwdSettingSyncApi
    AgvcQxySettingApi
    AgvcNbqSettingApi
    AgvcEventHisApi
    AgvcNbqHisApi
    AgvcBwdHisApi
    AgvcQxyHisApi
}

var (
    agvcBwdSettingService = service.ServiceGroupApp.AgvcServiceGroup.AgvcBwdSettingService
    agvcQxySettingService = service.ServiceGroupApp.AgvcServiceGroup.AgvcQxySettingService
    agvcNbqSettingService = service.ServiceGroupApp.AgvcServiceGroup.AgvcNbqSettingService
    agvcEventHisService   = service.ServiceGroupApp.AgvcServiceGroup.AgvcEventHisService
    agvcNbqHisService     = service.ServiceGroupApp.AgvcServiceGroup.AgvcNbqHisService
    agvcBwdHisService     = service.ServiceGroupApp.AgvcServiceGroup.AgvcBwdHisService
    agvcQxyHisService     = service.ServiceGroupApp.AgvcServiceGroup.AgvcQxyHisService
    agvcAgcHisService     = service.ServiceGroupApp.AgvcServiceGroup.AgvcAgcHisService
    agvcAvcHisService     = service.ServiceGroupApp.AgvcServiceGroup.AgvcAvcHisService
)
