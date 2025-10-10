package setting

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct{ SysInverterSettingRouter }

var sysInverterSettingApi = api.ApiGroupApp.SettingApiGroup.SysInverterSettingApi
