package setting

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct{ SysInverterSettingApi }

var sysInverterSettingService = service.ServiceGroupApp.SettingServiceGroup.SysInverterSettingService
