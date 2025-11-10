package initialize

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/agvc_main"
    agvcMainService "github.com/flipped-aurora/gin-vue-admin/server/service/agvc/agcv_main"
    "go.uber.org/zap"
)

// CoapHandler CoAP处理器函数类型
type CoapHandler func(ctx context.Context, msg coapMessage) (code byte, payload []byte)

// coapRoutes CoAP路由表
var coapRoutes = map[string]map[byte]CoapHandler{
    // path -> method -> handler
    "/health": {
        coapCodeGet: handleHealthCheck,
    },
    "/agvc/data": {
        coapCodePost: handleAgvcData,
    },
    "/agvc/nbq_setting": {
        coapCodePost: handleNbqSettingUpdate,
    },
}

// handleHealthCheck 健康检查处理器
func handleHealthCheck(ctx context.Context, msg coapMessage) (code byte, payload []byte) {
    return coapCodeContent, []byte("ok")
}

// handleAgvcData AGVC数据接收处理器
func handleAgvcData(ctx context.Context, msg coapMessage) (code byte, payload []byte) {
    // 解析JSON数据
    var dataBatch agvc.AgvcDataBatch
    // 查看原始JSON数据
    fmt.Println("cccccccc", string(msg.payload))
    if err := json.Unmarshal(msg.payload, &dataBatch); err != nil {
        global.GVA_LOG.Error("Failed to parse AGVC data", zap.Error(err))
        return coapCodeBadRequest, []byte(`{"error":"invalid json"}`)
    }

    // 验证数据
    if len(dataBatch) == 0 {
        global.GVA_LOG.Warn("Received empty AGVC data batch")
        return coapCodeBadRequest, []byte(`{"error":"empty data"}`)
    }

    // 转换数据格式并存储到DataStorage
    convertedData := make([]agvc_main.AgvcDataItem, len(dataBatch))
    for i, item := range dataBatch {
        convertedData[i] = agvc_main.AgvcDataItem{
            Psid:     item.Psid,
            Eqid:     item.Eqid,
            EqType:   item.EqType,
            DataType: item.DataType,
            Point:    item.Point,
            Value:    item.Value,
        }
    }
    //查看转换后的数据
    global.GVA_LOG.Debug("Converted AGVC data", zap.Any("data", convertedData))

    // 存储到内存，由DataStorage每5分钟定时保存到InfluxDB
    agvcMainService.DataStorage.StoreAgvcDataBatch(convertedData)

    global.GVA_LOG.Debug("AGVC data stored to memory", zap.Int("count", len(dataBatch)))
    return coapCodeCreated, []byte(`{"success":true}`)
}

// handleNbqSettingUpdate 逆变器配置更新处理器（从调度接收）
func handleNbqSettingUpdate(ctx context.Context, msg coapMessage) (code byte, payload []byte) {
    // 解析JSON数据
    var updateData struct {
        InverterNo          int      `json:"inverterNo"`          // 逆变器编号
        RatedActivePower    *float64 `json:"ratedActivePower"`    // 额定有功功率
        RatedReactivePower  *float64 `json:"ratedReactivePower"`  // 额定无功功率
        JitterRange         *float64 `json:"jitterRange"`         // 功率抖动区间
        DeadbandRange       *float64 `json:"deadbandRange"`       // 功率死区区间
        UpgradePriority     *float64 `json:"upgradePriority"`     // 升额优先级
        DowngradePriority   *float64 `json:"downgradePriority"`   // 降级优先级
        IsParticipateAdjust *bool    `json:"isParticipateAdjust"` // 参与调节
        IsBenchmarkInverter *bool    `json:"isBenchmarkInverter"` // 标杆逆变器
    }

    if err := json.Unmarshal(msg.payload, &updateData); err != nil {
        global.GVA_LOG.Error("Failed to parse nbq setting data", zap.Error(err))
        return coapCodeBadRequest, []byte(`{"error":"invalid json"}`)
    }

    // 验证逆变器编号
    if updateData.InverterNo == 0 {
        global.GVA_LOG.Warn("Received nbq setting update without inverterNo")
        return coapCodeBadRequest, []byte(`{"error":"inverterNo required"}`)
    }

    // 从数据库查询逆变器配置
    var nbqSetting agvc.AgvcNbqSetting
    if err := global.GVA_DB.Where("inverter_no = ?", updateData.InverterNo).First(&nbqSetting).Error; err != nil {
        global.GVA_LOG.Error("Failed to find nbq setting", zap.Int("inverterNo", updateData.InverterNo), zap.Error(err))
        return coapCodeNotFound, []byte(`{"error":"nbq setting not found"}`)
    }

    // 更新字段
    needUpdate := false
    if updateData.RatedActivePower != nil {
        nbqSetting.RatedActivePower = updateData.RatedActivePower
        needUpdate = true
    }
    if updateData.RatedReactivePower != nil {
        nbqSetting.RatedReactivePower = updateData.RatedReactivePower
        needUpdate = true
    }
    if updateData.JitterRange != nil {
        nbqSetting.JitterRange = updateData.JitterRange
        needUpdate = true
    }
    if updateData.DeadbandRange != nil {
        nbqSetting.DeadbandRange = updateData.DeadbandRange
        needUpdate = true
    }
    if updateData.UpgradePriority != nil {
        nbqSetting.UpgradePriority = updateData.UpgradePriority
        needUpdate = true
    }
    if updateData.DowngradePriority != nil {
        nbqSetting.DowngradePriority = updateData.DowngradePriority
        needUpdate = true
    }
    if updateData.IsParticipateAdjust != nil {
        nbqSetting.IsParticipateAdjust = updateData.IsParticipateAdjust
        needUpdate = true
    }
    if updateData.IsBenchmarkInverter != nil {
        nbqSetting.IsBenchmarkInverter = updateData.IsBenchmarkInverter
        needUpdate = true
    }

    // 保存到数据库
    if needUpdate {
        if err := global.GVA_DB.Save(&nbqSetting).Error; err != nil {
            global.GVA_LOG.Error("Failed to update nbq setting",
                zap.Int("inverterNo", updateData.InverterNo),
                zap.Error(err))
            return coapCodeInternalServerError, []byte(`{"error":"update failed"}`)
        }

        global.GVA_LOG.Info("Nbq setting updated from dispatch",
            zap.Int("inverterNo", updateData.InverterNo))
    }

    return coapCodeCreated, []byte(`{"success":true}`)
}

// getCoapHandler 获取CoAP处理器
func getCoapHandler(path string, method byte) CoapHandler {
    if methods, exists := coapRoutes[path]; exists {
        if handler, exists := methods[method]; exists {
            return handler
        }
    }
    return nil
}
