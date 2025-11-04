package agcv_main

import (
    "fmt"
    "math"
    "time"

    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/agvc_main"
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/agvc_main/request"
    "github.com/flipped-aurora/gin-vue-admin/server/service/agvc/cons"
    "go.uber.org/zap"
)

type agc struct {
    running    bool
    stopChan   chan struct{}
    bwdNoChans map[int]chan struct{} // 每个并网点一个停止通道
}

var AGC = new(agc)

// Initialize 初始化AGC服务
func (s *agc) Initialize() {
    s.bwdNoChans = make(map[int]chan struct{})
    global.GVA_LOG.Info("AGC控制服务初始化成功")
}

// StartAGC 启动AGC控制循环
func (s *agc) StartAGC(bwdNo int) error {
    // 检查是否已经在运行
    if _, exists := s.bwdNoChans[bwdNo]; exists {
        return fmt.Errorf("电站%s的AGC控制已在运行", bwdNo)
    }

    // 获取并网点配置
    config, err := s.GetAGCConfig(bwdNo)
    if err != nil {
        return fmt.Errorf("获取并网点配置失败: %v", err)
    }

    // 创建停止通道
    stopChan := make(chan struct{})
    s.bwdNoChans[bwdNo] = stopChan

    // 启动控制循环
    go s.agcControlLoop(bwdNo, config, stopChan)

    global.GVA_LOG.Info("AGC控制循环已启动", zap.String("bwdNo", fmt.Sprintf("%d", bwdNo)))
    return nil
}

// StopAGC 停止AGC控制循环
func (s *agc) StopAGC(bwdNo int) error {
    stopChan, exists := s.bwdNoChans[bwdNo]
    if !exists {
        return fmt.Errorf("电站%s的AGC控制未运行", bwdNo)
    }

    close(stopChan)
    delete(s.bwdNoChans, bwdNo)

    global.GVA_LOG.Info("AGC控制循环已停止", zap.String("bwdNo", fmt.Sprintf("%d", bwdNo)))
    return nil
}

// agcControlLoop AGC控制主循环
func (s *agc) agcControlLoop(bwdNo int, config agvc.AgvcBwdSetting, stopChan chan struct{}) {
    ticker := time.NewTicker(time.Duration(*config.AgcControlPeriod) * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            if err := s.executeAGCCycle(bwdNo); err != nil {
                global.GVA_LOG.Error("AGC控制周期执行失败",
                    zap.String("bwdNo", fmt.Sprintf("%d", bwdNo)),
                    zap.Error(err))
            }
        case <-stopChan:
            global.GVA_LOG.Info("AGC控制循环退出", zap.String("bwdNo", fmt.Sprintf("%d", bwdNo)))
            return
        }
    }
}

// executeAGCCycle 执行一个AGC控制周期
func (s *agc) executeAGCCycle(bwdNo int) error {
    // 步骤1：获取最新配置
    config, err := s.GetAGCConfig(bwdNo)
    if err != nil {
        return fmt.Errorf("获取AGC配置失败: %v", err)
    }

    // 步骤2：检查AGC是否投入
    if config.AgcIsEnabled == nil {
        global.GVA_LOG.Debug("AGC系统未投入", zap.String("bwdNo", fmt.Sprintf("%d", bwdNo)))
        return nil
    }

    // 步骤3：判断控制权限，获取执行值
    var execVal float64
    if config.ControlAuth != nil && *config.ControlAuth == 1 {
        // 调度控制：优先从内存获取调度下发的值
        pointID, err := PointMapper.GetPointID(cons.TYPE_BWG, "有功功率P(kW)")
        if err != nil {
            pointID = "7" // 使用默认值
        }
        
        // 从调度存储获取值
        dispatchVal, err := DispatchStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_BWG, cons.YT, pointID)
        if err == nil {
            // 成功从调度存储获取
            execVal = dispatchVal
            global.GVA_LOG.Debug("从调度存储获取AGC目标值",
                zap.Int("bwdNo", bwdNo),
                zap.Float64("execVal", execVal))
        } else {
            // 调度存储没有值，使用配置中的值
            execVal = *config.DispatchExecValue
            global.GVA_LOG.Debug("调度存储无数据，使用配置中的调度执行值",
                zap.Int("bwdNo", bwdNo),
                zap.Float64("execVal", execVal))
        }
    } else {
        // 站内控制
        execVal = *config.StationExecValue
    }

    // 步骤4：判断是否开环运行
    if config.RunMode != nil && *config.RunMode == 1 {
        // 开环运行，直接执行目标值
        global.GVA_LOG.Debug("AGC开环运行，直接执行目标值",
            zap.String("bwdNo", fmt.Sprintf("%d", bwdNo)),
            zap.Float64("execVal", execVal))
        return s.executeOpenLoopControl(bwdNo, execVal)
    }

    // 步骤5：闭环运行，采集数据
    collectData, err := s.collectAGCData(bwdNo)
    if err != nil {
        return fmt.Errorf("采集数据失败: %v", err)
    }

    // 步骤6：计算出力偏差
    targetOutput := execVal
    actualOutput := collectData["actualOutput"].(float64)
    outputDeviation := targetOutput - actualOutput

    global.GVA_LOG.Debug("AGC数据采集",
        zap.String("并网点编号", fmt.Sprintf("%d", bwdNo)),
        zap.Float64("目标出力", targetOutput),
        zap.Float64("实际出力", actualOutput),
        zap.Float64("出力偏差", outputDeviation))

    // 步骤7：判断偏差是否在抖动区间内
    if math.Abs(outputDeviation) <= *config.AgcVibrationRange {
        global.GVA_LOG.Debug("偏差在抖动区间内，无需调节",
            zap.String("bwdNo", fmt.Sprintf("%d", bwdNo)),
            zap.Float64("偏差", outputDeviation),
            zap.Float64("抖动区间", *config.AgcVibrationRange))
        return nil
    }

    // 步骤8：检查闭锁信号
    // if outputDeviation > 0 && config.UpRegLock != nil && *config.UpRegLock == 1 {
    //     global.GVA_LOG.Warn("上调节闭锁，无法增加出力", zap.String("bwdNo", bwdNo))
    //     return nil
    // }
    // if outputDeviation < 0 && config.DownRegLock != nil && *config.DownRegLock == 1 {
    //     global.GVA_LOG.Warn("下调节闭锁，无法减少出力", zap.String("bwdNo", bwdNo))
    //     return nil
    // }

    // 步骤9：获取可用逆变器
    inverters, err := Device.GetOnlineInvertersByBwdNo(bwdNo)
    if err != nil || len(inverters) == 0 {
        return fmt.Errorf("没有可用的逆变器: %v", err)
    }

    // 步骤10：平均分配调节量
    perInvDeviation := outputDeviation / float64(len(inverters))
    regulationDetails := make([]agvc_main.InverterRegulation, 0, len(inverters))

    for _, inv := range inverters {
        // 限制调节量不超过逆变器最大调节能力
        actualReg := perInvDeviation
        if math.Abs(actualReg) > *inv.RatedActivePower {
            if actualReg > 0 {
                actualReg = *inv.RatedActivePower
            } else {
                actualReg = -*inv.RatedActivePower
            }
        }

        // 获取逆变器当前功率
        currentPower, err := DataStorage.GetDataAsFloat64(*inv.InverterNo, 2, cons.YC, "28")
        if err != nil {
            global.GVA_LOG.Warn("获取逆变器当前功率失败",
                zap.String("逆变器编号", fmt.Sprintf("%d", *inv.InverterNo)),
                zap.Error(err))
            currentPower = 0
        }

        regulationDetails = append(regulationDetails, agvc_main.InverterRegulation{
            BwdNo:           bwdNo,
            EQID:            *inv.InverterNo,
            RegulationPower: actualReg,
            BeforePower:     currentPower,
            Status:          "pending",
        })
    }

    // 步骤11：记录调节开始
    record := agvc_main.AGCRegulationRecord{
        BwdNo:           bwdNo,
        TargetPower:     targetOutput,
        ActualPower:     actualOutput,
        PowerDeviation:  outputDeviation,
        RegulationPower: outputDeviation,
        SystemFreq:      collectData["systemFreq"].(float64),
        Status:          "executing",
        Message:         fmt.Sprintf("开始调节，分配%d个逆变器", len(inverters)),
    }

    // if err := global.GVA_DB.Create(&record).Error; err != nil {
    //     return fmt.Errorf("记录调节失败: %v", err)
    // }

    // 步骤12：下发调节指令
    successCount := 0
    for i := range regulationDetails {
        detail := &regulationDetails[i]
        detail.RecordID = record.ID

        // 计算目标功率
        targetPower := detail.BeforePower + detail.RegulationPower

        // 构建指令
        commands := map[int]interface{}{
            401: targetPower, // 有功功率降额执行值（遥调）
        }

        // 发送CoAP指令
        host := CoapSender.GetDefaultCoapHost()
        port := CoapSender.GetDefaultCoapPort()
        psid := 1
        if err := CoapSender.SendInverterCommand(host, port, psid, detail.EQID, commands); err != nil {
            global.GVA_LOG.Error("发送逆变器调节指令失败",
                zap.Int("eqid", detail.EQID),
                zap.Error(err))
            detail.Status = "failed"
            detail.AfterPower = detail.BeforePower
        } else {
            detail.Status = "success"
            detail.AfterPower = targetPower
            successCount++
        }

        // 保存调节详情
        // global.GVA_DB.Create(detail)
    }

    // 步骤13：更新调节记录
    if successCount == len(regulationDetails) {
        record.Status = "success"
        record.Message = fmt.Sprintf("调节成功，%d个逆变器全部响应", successCount)
    } else if successCount > 0 {
        record.Status = "partial"
        record.Message = fmt.Sprintf("部分调节成功，%d/%d个逆变器响应", successCount, len(regulationDetails))
    } else {
        record.Status = "failed"
        record.Message = "调节失败，所有逆变器无响应"
    }

    global.GVA_DB.Model(&record).Updates(map[string]interface{}{
        "status":  record.Status,
        "message": record.Message,
    })

    global.GVA_LOG.Info("AGC调节周期完成",
        zap.Int("bwdNo", bwdNo),
        zap.String("status", record.Status),
        zap.Int("成功数", successCount),
        zap.Int("总数", len(regulationDetails)))

    return nil
}

// executeOpenLoopControl 执行开环控制
func (s *agc) executeOpenLoopControl(bwdNo int, execVal float64) error {
    // 获取可用逆变器
    inverters, err := Device.GetOnlineInvertersByBwdNo(bwdNo)
    if err != nil || len(inverters) == 0 {
        return fmt.Errorf("没有可用的逆变器: %v", err)
    }

    // 平均分配目标功率
    perInvPower := execVal / float64(len(inverters))

    for _, inv := range inverters {
        commands := map[int]interface{}{
            401: perInvPower,
        }

        host := CoapSender.GetDefaultCoapHost()
        port := CoapSender.GetDefaultCoapPort()
        // bwdNo := *inv.BwdNo
        psid := 1
        CoapSender.SendInverterCommand(host, port, psid, *inv.InverterNo, commands)
    }

    return nil
}

// collectAGCData 采集AGC控制所需数据
func (s *agc) collectAGCData(bwdNo int) (map[string]interface{}, error) {
    data := make(map[string]interface{})

    // 采集实际出力（有功功率）：优先从DataStorage获取，如无则从逆变器聚合
    actualOutput, err := PowerAggregator.GetBwgActivePower(bwdNo)
    if err != nil {
        return nil, fmt.Errorf("采集实际出力失败: %v", err)
    }
    data["actualOutput"] = actualOutput

    // 采集系统频率
    pointID, err := PointMapper.GetPointID(cons.TYPE_BWG, "F(频率Hz)")
    if err != nil {
        pointID = "10" // 使用默认值
    }
    
    systemFreq, err := DataStorage.GetDataAsFloat64(bwdNo, cons.TYPE_BWG, cons.YC, pointID)
    if err != nil {
        // 频率不是必须的，使用默认值
        systemFreq = 50.0
    }
    data["systemFreq"] = systemFreq

    return data, nil
}

// GetAGCConfig 获取AGC配置
func (s *agc) GetAGCConfig(bwdNo int) (agvc.AgvcBwdSetting, error) {
    var config agvc.AgvcBwdSetting
    err := global.GVA_DB.Where("number = ?", bwdNo).First(&config).Error
    return config, err
}

// UpdateAGCConfig 更新AGC配置
func (s *agc) UpdateAGCConfig(req request.AGCConfigUpdate) error {
    updates := make(map[string]interface{})

    if req.IsActive != nil {
        updates["is_active"] = *req.IsActive
    }
    if req.ControlAuth != nil {
        updates["control_auth"] = *req.ControlAuth
    }
    if req.RunMode != nil {
        updates["run_mode"] = *req.RunMode
    }
    if req.JitterRange != 0 {
        updates["jitter_range"] = req.JitterRange
    }
    if req.RegPeriod != 0 {
        updates["reg_period"] = req.RegPeriod
    }
    if req.RegStep != 0 {
        updates["reg_step"] = req.RegStep
    }
    updates["power_upper_limit"] = req.PowerUpperLimit
    updates["power_lower_limit"] = req.PowerLowerLimit
    updates["power_exec_value"] = req.PowerExecValue
    if req.UpRegLock != nil {
        updates["up_reg_lock"] = *req.UpRegLock
    }
    if req.DownRegLock != nil {
        updates["down_reg_lock"] = *req.DownRegLock
    }
    updates["dispatch_exec_value"] = req.DispatchExecValue
    updates["station_exec_value"] = req.StationExecValue

    return global.GVA_DB.Model(&agvc_main.AGCConfig{}).
        Where("psid = ?", req.PSID).
        Updates(updates).Error
}

// GetAGCRecords 获取AGC调节记录
func (s *agc) GetAGCRecords(req request.AGCRegulationRecordSearch) ([]agvc_main.AGCRegulationRecord, int64, error) {
    var records []agvc_main.AGCRegulationRecord
    var total int64

    db := global.GVA_DB.Model(&agvc_main.AGCRegulationRecord{})

    if req.PSID != "" {
        db = db.Where("psid = ?", req.PSID)
    }
    if req.Status != "" {
        db = db.Where("status = ?", req.Status)
    }
    if req.StartTime != nil && req.EndTime != nil {
        db = db.Where("created_at BETWEEN ? AND ?",
            time.Unix(*req.StartTime, 0),
            time.Unix(*req.EndTime, 0))
    }

    err := db.Count(&total).Error
    if err != nil {
        return nil, 0, err
    }

    if req.PageSize > 0 {
        offset := (req.Page - 1) * req.PageSize
        db = db.Limit(req.PageSize).Offset(offset)
    }

    err = db.Order("created_at DESC").Find(&records).Error
    return records, total, err
}

// CreateOrUpdateConfig 创建或更新配置
func (s *agc) CreateOrUpdateConfig(config *agvc_main.AGCConfig) error {
    var existing agvc_main.AGCConfig
    err := global.GVA_DB.Where("psid = ?", config.PSID).First(&existing).Error

    if err != nil {
        // 不存在，创建新配置
        return global.GVA_DB.Create(config).Error
    }

    // 存在，更新配置
    return global.GVA_DB.Model(&existing).Updates(config).Error
}
