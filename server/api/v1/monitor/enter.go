package monitor

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct{ InverterMonitorApi }

var inverterMonitorService = service.ServiceGroupApp.MonitorServiceGroup.InverterMonitorService
