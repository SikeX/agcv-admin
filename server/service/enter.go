package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service/agvc"
	"github.com/flipped-aurora/gin-vue-admin/server/service/example"
	"github.com/flipped-aurora/gin-vue-admin/server/service/monitor"
	"github.com/flipped-aurora/gin-vue-admin/server/service/setting"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	SystemServiceGroup  system.ServiceGroup
	ExampleServiceGroup example.ServiceGroup
	SettingServiceGroup setting.ServiceGroup
	MonitorServiceGroup monitor.ServiceGroup
	AgvcServiceGroup    agvc.ServiceGroup
}
