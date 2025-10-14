
package system

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
    systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
)

type SysSvgSvcSettingService struct {}
// CreateSysSvgSvcSetting 创建SVG/SVC设置记录
// Author [yourname](https://github.com/yourname)
func (sysSvgSvcSettingService *SysSvgSvcSettingService) CreateSysSvgSvcSetting(ctx context.Context, sysSvgSvcSetting *system.SysSvgSvcSetting) (err error) {
	err = global.GVA_DB.Create(sysSvgSvcSetting).Error
	return err
}

// DeleteSysSvgSvcSetting 删除SVG/SVC设置记录
// Author [yourname](https://github.com/yourname)
func (sysSvgSvcSettingService *SysSvgSvcSettingService)DeleteSysSvgSvcSetting(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&system.SysSvgSvcSetting{},"id = ?",ID).Error
	return err
}

// DeleteSysSvgSvcSettingByIds 批量删除SVG/SVC设置记录
// Author [yourname](https://github.com/yourname)
func (sysSvgSvcSettingService *SysSvgSvcSettingService)DeleteSysSvgSvcSettingByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]system.SysSvgSvcSetting{},"id in ?",IDs).Error
	return err
}

// UpdateSysSvgSvcSetting 更新SVG/SVC设置记录
// Author [yourname](https://github.com/yourname)
func (sysSvgSvcSettingService *SysSvgSvcSettingService)UpdateSysSvgSvcSetting(ctx context.Context, sysSvgSvcSetting system.SysSvgSvcSetting) (err error) {
	err = global.GVA_DB.Model(&system.SysSvgSvcSetting{}).Where("id = ?",sysSvgSvcSetting.ID).Updates(&sysSvgSvcSetting).Error
	return err
}

// GetSysSvgSvcSetting 根据ID获取SVG/SVC设置记录
// Author [yourname](https://github.com/yourname)
func (sysSvgSvcSettingService *SysSvgSvcSettingService)GetSysSvgSvcSetting(ctx context.Context, ID string) (sysSvgSvcSetting system.SysSvgSvcSetting, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&sysSvgSvcSetting).Error
	return
}
// GetSysSvgSvcSettingInfoList 分页获取SVG/SVC设置记录
// Author [yourname](https://github.com/yourname)
func (sysSvgSvcSettingService *SysSvgSvcSettingService)GetSysSvgSvcSettingInfoList(ctx context.Context, info systemReq.SysSvgSvcSettingSearch) (list []system.SysSvgSvcSetting, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&system.SysSvgSvcSetting{})
    var sysSvgSvcSettings []system.SysSvgSvcSetting
    // 如果有条件搜索 下方会自动创建搜索语句
    if info.ID != nil {
        db = db.Where("id = ?", *info.ID)
    }
    if info.WugongName != nil && *info.WugongName != "" {
        db = db.Where("wugong_name LIKE ?", "%"+*info.WugongName+"%")
    }
    if info.IsAdjustment != nil && *info.IsAdjustment != "" {
        db = db.Where("is_adjustment = ?", *info.IsAdjustment)
    }
    
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	if limit != 0 {
       db = db.Limit(limit).Offset(offset)
    }

	err = db.Find(&sysSvgSvcSettings).Error
	return  sysSvgSvcSettings, total, err
}
func (sysSvgSvcSettingService *SysSvgSvcSettingService)GetSysSvgSvcSettingPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
