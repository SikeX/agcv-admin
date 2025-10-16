package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/setting"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(setting.SysInverterSetting{}, system.SysGridConnectionPoint{}, system.SysHisEvent{}, system.SysSvgSvcSetting{}, system.SysQixiangyiSetting{}, system.SysTestPoint{})
	if err != nil {
		return err
	}
	return nil
}
