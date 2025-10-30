
package agvc

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
    agvcReq "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/request"
)

type AgvcBwdSettingService struct {}
// CreateAgvcBwdSetting 创建并网点配置记录
// Author [yourname](https://github.com/yourname)
func (agvcBwdSettingService *AgvcBwdSettingService) CreateAgvcBwdSetting(ctx context.Context, agvcBwdSetting *agvc.AgvcBwdSetting) (err error) {
	err = global.GVA_DB.Create(agvcBwdSetting).Error
	return err
}

// DeleteAgvcBwdSetting 删除并网点配置记录
// Author [yourname](https://github.com/yourname)
func (agvcBwdSettingService *AgvcBwdSettingService)DeleteAgvcBwdSetting(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&agvc.AgvcBwdSetting{},"id = ?",ID).Error
	return err
}

// DeleteAgvcBwdSettingByIds 批量删除并网点配置记录
// Author [yourname](https://github.com/yourname)
func (agvcBwdSettingService *AgvcBwdSettingService)DeleteAgvcBwdSettingByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]agvc.AgvcBwdSetting{},"id in ?",IDs).Error
	return err
}

// UpdateAgvcBwdSetting 更新并网点配置记录
// Author [yourname](https://github.com/yourname)
func (agvcBwdSettingService *AgvcBwdSettingService)UpdateAgvcBwdSetting(ctx context.Context, agvcBwdSetting agvc.AgvcBwdSetting) (err error) {
	err = global.GVA_DB.Model(&agvc.AgvcBwdSetting{}).Where("id = ?",agvcBwdSetting.ID).Updates(&agvcBwdSetting).Error
	return err
}

// GetAgvcBwdSetting 根据ID获取并网点配置记录
// Author [yourname](https://github.com/yourname)
func (agvcBwdSettingService *AgvcBwdSettingService)GetAgvcBwdSetting(ctx context.Context, ID string) (agvcBwdSetting agvc.AgvcBwdSetting, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&agvcBwdSetting).Error
	return
}
// GetAgvcBwdSettingInfoList 分页获取并网点配置记录
// Author [yourname](https://github.com/yourname)
func (agvcBwdSettingService *AgvcBwdSettingService)GetAgvcBwdSettingInfoList(ctx context.Context, info agvcReq.AgvcBwdSettingSearch) (list []agvc.AgvcBwdSetting, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&agvc.AgvcBwdSetting{})
    var agvcBwdSettings []agvc.AgvcBwdSetting
    // 如果有条件搜索 下方会自动创建搜索语句
    if info.Name != "" {
     db = db.Where("name LIKE ?", "%"+info.Name+"%")
    }
    if info.Number != "" {
     db = db.Where("number LIKE ?", "%"+info.Number+"%")
    }
    
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	if limit != 0 {
       db = db.Limit(limit).Offset(offset)
    }

	err = db.Find(&agvcBwdSettings).Error
	return  agvcBwdSettings, total, err
}
func (agvcBwdSettingService *AgvcBwdSettingService)GetAgvcBwdSettingPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
