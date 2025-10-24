package monitor

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	InverterMonitorApi
	SysGridConnectionPointHistoryApi
}

var (
	inverterMonitorService               = service.ServiceGroupApp.MonitorServiceGroup.InverterMonitorService
	sysGridConnectionPointHistoryService = service.ServiceGroupApp.MonitorServiceGroup.SysGridConnectionPointHistoryService
)
