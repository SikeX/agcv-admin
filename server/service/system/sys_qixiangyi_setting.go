
package system

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
)

type SysQixiangyiSettingService struct {}
// CreateSysQixiangyiSetting 创建气象仪配置记录
// Author [yourname](https://github.com/yourname)
func (sysQixiangyiSettingService *SysQixiangyiSettingService) CreateSysQixiangyiSetting(ctx context.Context, sysQixiangyiSetting *system.SysQixiangyiSetting) (err error) {
	err = global.GVA_DB.Create(sysQixiangyiSetting).Error
	return err
}

// DeleteSysQixiangyiSetting 删除气象仪配置记录
// Author [yourname](https://github.com/yourname)
func (sysQixiangyiSettingService *SysQixiangyiSettingService)DeleteSysQixiangyiSetting(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&system.SysQixiangyiSetting{},"id = ?",ID).Error
	return err
}

// DeleteSysQixiangyiSettingByIds 批量删除气象仪配置记录
// Author [yourname](https://github.com/yourname)
func (sysQixiangyiSettingService *SysQixiangyiSettingService)DeleteSysQixiangyiSettingByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]system.SysQixiangyiSetting{},"id in ?",IDs).Error
	return err
}

// UpdateSysQixiangyiSetting 更新气象仪配置记录
// Author [yourname](https://github.com/yourname)
func (sysQixiangyiSettingService *SysQixiangyiSettingService)UpdateSysQixiangyiSetting(ctx context.Context, sysQixiangyiSetting system.SysQixiangyiSetting) (err error) {
	err = global.GVA_DB.Model(&system.SysQixiangyiSetting{}).Where("id = ?",sysQixiangyiSetting.ID).Updates(&sysQixiangyiSetting).Error
	return err
}

// GetSysQixiangyiSetting 根据ID获取气象仪配置记录
// Author [yourname](https://github.com/yourname)
func (sysQixiangyiSettingService *SysQixiangyiSettingService)GetSysQixiangyiSetting(ctx context.Context, ID string) (sysQixiangyiSetting system.SysQixiangyiSetting, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&sysQixiangyiSetting).Error
	return
}
// GetSysQixiangyiSettingInfoList 分页获取气象仪配置记录
// Author [yourname](https://github.com/yourname)
func (sysQixiangyiSettingService *SysQixiangyiSettingService)GetSysQixiangyiSettingInfoList(ctx context.Context, info systemReq.SysQixiangyiSettingSearch) (list []system.SysQixiangyiSetting, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&system.SysQixiangyiSetting{})
    var sysQixiangyiSettings []system.SysQixiangyiSetting
    // 如果有条件搜索 下方会自动创建搜索语句
    if info.ID != 0 {
     db = db.Where("id = ?",info.ID)
    }
    if info.DeviceName != "" {
     db = db.Where("device_name LIKE ?", "%"+info.DeviceName+"%")
    }
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	if limit != 0 {
       db = db.Limit(limit).Offset(offset)
    }

	err = db.Find(&sysQixiangyiSettings).Error
	return  sysQixiangyiSettings, total, err
}
func (sysQixiangyiSettingService *SysQixiangyiSettingService)GetSysQixiangyiSettingPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
