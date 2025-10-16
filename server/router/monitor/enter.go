package monitor

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct{ InverterMonitorRouter }

var inverterMonitorApi = api.ApiGroupApp.MonitorApiGroup.InverterMonitorApi
