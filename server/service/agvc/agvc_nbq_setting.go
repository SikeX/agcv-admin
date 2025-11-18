package agvc

import (
    "context"

    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
    agvcReq "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/request"
)

type AgvcNbqSettingService struct{}

// CreateAgvcNbqSetting 创建逆变器配置记录
// Author [yourname](https://github.com/yourname)
func (agvcNbqSettingService *AgvcNbqSettingService) CreateAgvcNbqSetting(ctx context.Context, agvcNbqSetting *agvc.AgvcNbqSetting) (err error) {
    err = global.GVA_DB.Create(agvcNbqSetting).Error
    return err
}

// DeleteAgvcNbqSetting 删除逆变器配置记录
// Author [yourname](https://github.com/yourname)
func (agvcNbqSettingService *AgvcNbqSettingService) DeleteAgvcNbqSetting(ctx context.Context, id string) (err error) {
    err = global.GVA_DB.Delete(&agvc.AgvcNbqSetting{}, "id = ?", id).Error
    return err
}

// DeleteAgvcNbqSettingByIds 批量删除逆变器配置记录
// Author [yourname](https://github.com/yourname)
func (agvcNbqSettingService *AgvcNbqSettingService) DeleteAgvcNbqSettingByIds(ctx context.Context, ids []string) (err error) {
    err = global.GVA_DB.Delete(&[]agvc.AgvcNbqSetting{}, "id in ?", ids).Error
    return err
}

// UpdateAgvcNbqSetting 更新逆变器配置记录
// Author [yourname](https://github.com/yourname)
func (agvcNbqSettingService *AgvcNbqSettingService) UpdateAgvcNbqSetting(ctx context.Context, agvcNbqSetting agvc.AgvcNbqSetting) (err error) {
    err = global.GVA_DB.Model(&agvc.AgvcNbqSetting{}).Where("id = ?", agvcNbqSetting.ID).Updates(&agvcNbqSetting).Error
    return err
}

// GetAgvcNbqSetting 根据id获取逆变器配置记录
// Author [yourname](https://github.com/yourname)
func (agvcNbqSettingService *AgvcNbqSettingService) GetAgvcNbqSetting(ctx context.Context, id string) (agvcNbqSetting agvc.AgvcNbqSetting, err error) {
    err = global.GVA_DB.Where("id = ?", id).First(&agvcNbqSetting).Error
    return
}

// GetAgvcNbqSettingInfoList 分页获取逆变器配置记录
// Author [yourname](https://github.com/yourname)
func (agvcNbqSettingService *AgvcNbqSettingService) GetAgvcNbqSettingInfoList(ctx context.Context, info agvcReq.AgvcNbqSettingSearch) (list []agvc.AgvcNbqSetting, total int64, err error) {
    limit := info.PageSize
    offset := info.PageSize * (info.Page - 1)
    // 创建db
    db := global.GVA_DB.Model(&agvc.AgvcNbqSetting{})
    var agvcNbqSettings []agvc.AgvcNbqSetting
    // 如果有条件搜索 下方会自动创建搜索语句

    if info.Psid != nil {
        db = db.Where("psid = ?", *info.Psid)
    }
    if info.InverterNo != nil {
        db = db.Where("inverter_no Like ?", "%"+*info.InverterNo+"%")
    }
    if info.Name != nil && *info.Name != "" {
        db = db.Where("name LIKE ?", "%"+*info.Name+"%")
    }
    if info.IsParticipateAdjust != nil {
        db = db.Where("is_participate_adjust = ?", *info.IsParticipateAdjust)
    }
    if info.IsBenchmarkInverter != nil {
        db = db.Where("is_benchmark_inverter = ?", *info.IsBenchmarkInverter)
    }
    err = db.Count(&total).Error
    if err != nil {
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

    err = db.Find(&agvcNbqSettings).Error
    return agvcNbqSettings, total, err
}
func (agvcNbqSettingService *AgvcNbqSettingService) GetAgvcNbqSettingPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}

// GetAgvcNbqSettingByInverterNo 根据逆变器编号获取逆变器配置记录
// Author [yourname](https://github.com/yourname)
func (agvcNbqSettingService *AgvcNbqSettingService) GetAgvcNbqSettingByInverterNo(ctx context.Context, inverterNo string) (agvcNbqSetting agvc.AgvcNbqSetting, err error) {
    err = global.GVA_DB.Where("inverter_no = ?", inverterNo).First(&agvcNbqSetting).Error
    return
}
