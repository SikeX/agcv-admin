package monitor

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	InverterMonitorRouter
	SysGridConnectionPointHistoryRouter
}

var (
	inverterMonitorApi               = api.ApiGroupApp.MonitorApiGroup.InverterMonitorApi
	sysGridConnectionPointHistoryApi = api.ApiGroupApp.MonitorApiGroup.SysGridConnectionPointHistoryApi
)
