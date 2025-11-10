package agvc

import (
    "context"
    "fmt"

    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
    agvcReq "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/request"
    "github.com/flipped-aurora/gin-vue-admin/server/service/agvc/agcv_main"
    "go.uber.org/zap"
)

type AgvcNbqSettingService struct{}

// CreateAgvcNbqSetting 创建逆变器配置记录
// Author [yourname](https://github.com/yourname)
func (agvcNbqSettingService *AgvcNbqSettingService) CreateAgvcNbqSetting(ctx context.Context, agvcNbqSetting *agvc.AgvcNbqSetting) (err error) {
    err = global.GVA_DB.Create(agvcNbqSetting).Error
    if err != nil {
        return err
    }

    // 同步逆变器配置到调度系统
    if err := agvcNbqSettingService.syncNbqSettingToDispatch(ctx, *agvcNbqSetting); err != nil {
        global.GVA_LOG.Warn("同步逆变器配置到调度失败",
            zap.Uint("id", agvcNbqSetting.ID),
            zap.Error(err))
        // 不影响创建操作，仅记录警告日志
    }

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
    if err != nil {
        return err
    }

    // 同步逆变器配置到调度系统
    if err := agvcNbqSettingService.syncNbqSettingToDispatch(ctx, agvcNbqSetting); err != nil {
        global.GVA_LOG.Warn("同步逆变器配置到调度失败",
            zap.Uint("id", agvcNbqSetting.ID),
            zap.Error(err))
        // 不影响更新操作，仅记录警告日志
    }

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

// syncNbqSettingToDispatch 同步逆变器配置到调度系统
func (agvcNbqSettingService *AgvcNbqSettingService) syncNbqSettingToDispatch(ctx context.Context, setting agvc.AgvcNbqSetting) error {
    // 检查必要字段
    if setting.InverterNo == nil {
        return fmt.Errorf("逆变器编号不能为空")
    }

    // 获取并网点编号
    if setting.BwdNo == nil {
        global.GVA_LOG.Warn("逆变器配置缺少并网点编号，跳过同步",
            zap.Int("inverterNo", *setting.InverterNo))
        return nil
    }

    var bwdNo int
    if _, err := fmt.Sscanf(*setting.BwdNo, "%d", &bwdNo); err != nil {
        return fmt.Errorf("并网点编号格式错误: %v", err)
    }

    // 构建同步数据
    results := make(map[string]interface{})

    // 逆变器配置相关信息
    if setting.RatedActivePower != nil {
        results["ratedActivePower"] = *setting.RatedActivePower
    }
    if setting.RatedReactivePower != nil {
        results["ratedReactivePower"] = *setting.RatedReactivePower
    }
    if setting.IsParticipateAdjust != nil {
        results["isParticipateAdjust"] = *setting.IsParticipateAdjust
    }
    if setting.IsBenchmarkInverter != nil {
        results["isBenchmarkInverter"] = *setting.IsBenchmarkInverter
    }
    if setting.UpgradePriority != nil {
        results["upgradePriority"] = *setting.UpgradePriority
    }
    if setting.DowngradePriority != nil {
        results["downgradePriority"] = *setting.DowngradePriority
    }
    if setting.JitterRange != nil {
        results["jitterRange"] = *setting.JitterRange
    }
    if setting.DeadbandRange != nil {
        results["deadbandRange"] = *setting.DeadbandRange
    }

    // 添加逆变器编号信息
    results["inverterNo"] = *setting.InverterNo

    // 如果没有需要同步的数据，直接返回
    if len(results) <= 1 { // 只有inverterNo
        global.GVA_LOG.Debug("逆变器配置无需同步数据",
            zap.Int("inverterNo", *setting.InverterNo))
        return nil
    }

    // 发送到调度系统
    if err := agcv_main.CoapSender.SendAGCResultToDispatch(bwdNo, results); err != nil {
        return fmt.Errorf("发送逆变器配置到调度失败: %v", err)
    }

    global.GVA_LOG.Info("逆变器配置已同步到调度",
        zap.Int("inverterNo", *setting.InverterNo),
        zap.Int("bwdNo", bwdNo),
        zap.Int("dataCount", len(results)))

    return nil
}
