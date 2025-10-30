package agvc

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
	agvcReq "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/request"
)

type AgvcQxySettingService struct{}

// CreateAgvcQxySetting 创建气象仪配置记录
// Author [yourname](https://github.com/yourname)
func (agvcQxySettingService *AgvcQxySettingService) CreateAgvcQxySetting(ctx context.Context, agvcQxySetting *agvc.AgvcQxySetting) (err error) {
	err = global.GVA_DB.Create(agvcQxySetting).Error
	return err
}

// DeleteAgvcQxySetting 删除气象仪配置记录
// Author [yourname](https://github.com/yourname)
func (agvcQxySettingService *AgvcQxySettingService) DeleteAgvcQxySetting(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&agvc.AgvcQxySetting{}, "id = ?", ID).Error
	return err
}

// DeleteAgvcQxySettingByIds 批量删除气象仪配置记录
// Author [yourname](https://github.com/yourname)
func (agvcQxySettingService *AgvcQxySettingService) DeleteAgvcQxySettingByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]agvc.AgvcQxySetting{}, "id in ?", IDs).Error
	return err
}

// UpdateAgvcQxySetting 更新气象仪配置记录
// Author [yourname](https://github.com/yourname)
func (agvcQxySettingService *AgvcQxySettingService) UpdateAgvcQxySetting(ctx context.Context, agvcQxySetting agvc.AgvcQxySetting) (err error) {
	err = global.GVA_DB.Model(&agvc.AgvcQxySetting{}).Where("id = ?", agvcQxySetting.ID).Updates(&agvcQxySetting).Error
	return err
}

// GetAgvcQxySetting 根据ID获取气象仪配置记录
// Author [yourname](https://github.com/yourname)
func (agvcQxySettingService *AgvcQxySettingService) GetAgvcQxySetting(ctx context.Context, ID string) (agvcQxySetting agvc.AgvcQxySetting, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&agvcQxySetting).Error
	return
}

// GetAgvcQxySettingInfoList 分页获取气象仪配置记录
// Author [yourname](https://github.com/yourname)
func (agvcQxySettingService *AgvcQxySettingService) GetAgvcQxySettingInfoList(ctx context.Context, info agvcReq.AgvcQxySettingSearch) (list []agvc.AgvcQxySetting, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&agvc.AgvcQxySetting{})
	var agvcQxySettings []agvc.AgvcQxySetting
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.Number != "" {
		db = db.Where("number LIKE ?", "%"+info.Number+"%")
	}
	if info.DeviceName != "" {
		db = db.Where("device_name LIKE ?", "%"+info.DeviceName+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&agvcQxySettings).Error
	return agvcQxySettings, total, err
}
func (agvcQxySettingService *AgvcQxySettingService) GetAgvcQxySettingPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
