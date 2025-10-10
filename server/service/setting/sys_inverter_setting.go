
package setting

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/setting"
    settingReq "github.com/flipped-aurora/gin-vue-admin/server/model/setting/request"
)

type SysInverterSettingService struct {}
// CreateSysInverterSetting 创建逆变器设置记录
// Author [yourname](https://github.com/yourname)
func (sysInverterSettingService *SysInverterSettingService) CreateSysInverterSetting(ctx context.Context, sysInverterSetting *setting.SysInverterSetting) (err error) {
	err = global.GVA_DB.Create(sysInverterSetting).Error
	return err
}

// DeleteSysInverterSetting 删除逆变器设置记录
// Author [yourname](https://github.com/yourname)
func (sysInverterSettingService *SysInverterSettingService)DeleteSysInverterSetting(ctx context.Context, id string) (err error) {
	err = global.GVA_DB.Delete(&setting.SysInverterSetting{},"id = ?",id).Error
	return err
}

// DeleteSysInverterSettingByIds 批量删除逆变器设置记录
// Author [yourname](https://github.com/yourname)
func (sysInverterSettingService *SysInverterSettingService)DeleteSysInverterSettingByIds(ctx context.Context, ids []string) (err error) {
	err = global.GVA_DB.Delete(&[]setting.SysInverterSetting{},"id in ?",ids).Error
	return err
}

// UpdateSysInverterSetting 更新逆变器设置记录
// Author [yourname](https://github.com/yourname)
func (sysInverterSettingService *SysInverterSettingService)UpdateSysInverterSetting(ctx context.Context, sysInverterSetting setting.SysInverterSetting) (err error) {
	err = global.GVA_DB.Model(&setting.SysInverterSetting{}).Where("id = ?",sysInverterSetting.Id).Updates(&sysInverterSetting).Error
	return err
}

// GetSysInverterSetting 根据id获取逆变器设置记录
// Author [yourname](https://github.com/yourname)
func (sysInverterSettingService *SysInverterSettingService)GetSysInverterSetting(ctx context.Context, id string) (sysInverterSetting setting.SysInverterSetting, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&sysInverterSetting).Error
	return
}
// GetSysInverterSettingInfoList 分页获取逆变器设置记录
// Author [yourname](https://github.com/yourname)
func (sysInverterSettingService *SysInverterSettingService)GetSysInverterSettingInfoList(ctx context.Context, info settingReq.SysInverterSettingSearch) (list []setting.SysInverterSetting, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&setting.SysInverterSetting{})
    var sysInverterSettings []setting.SysInverterSetting
    // 如果有条件搜索 下方会自动创建搜索语句
    
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }
        var OrderStr string
        orderMap := make(map[string]bool)
         	orderMap["inverter_no"] = true
       if orderMap[info.Sort] {
          OrderStr = info.Sort
          if info.Order == "descending" {
             OrderStr = OrderStr + " desc"
          }
          db = db.Order(OrderStr)
       }

	if limit != 0 {
       db = db.Limit(limit).Offset(offset)
    }

	err = db.Find(&sysInverterSettings).Error
	return  sysInverterSettings, total, err
}
func (sysInverterSettingService *SysInverterSettingService)GetSysInverterSettingPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
