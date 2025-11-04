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

type avc struct {
    psidChans map[int]chan struct{} // 每个电站一个停止通道
}

var AVC = new(avc)

// Initialize 初始化AVC服务
func (s *avc) Initialize() {
    s.psidChans = make(map[int]chan struct{})
    global.GVA_LOG.Info("AVC控制服务初始化成功")
}

// AutoStartAllGridPoints 自动启动所有并网点的AVC功能
func (s *avc) AutoStartAllGridPoints() {
    // 查询所有并网点配置
    var settings []agvc.AgvcBwdSetting
    if err := global.GVA_DB.Find(&settings).Error; err != nil {
        global.GVA_LOG.Error("查询并网点配置失败", zap.Error(err))
        return
    }

    if len(settings) == 0 {
        global.GVA_LOG.Info("未找到并网点配置，跳过AVC自动启动")
        return
    }

    global.GVA_LOG.Info("开始自动启动所有并网点的AVC功能", zap.Int("并网点数量", len(settings)))

    successCount := 0
    for _, setting := range settings {
        if setting.Number == nil {
            continue
        }

        // 将并网点编号转换为整数
        var bwdNo int
        if _, err := fmt.Sscanf(*setting.Number, "%d", &bwdNo); err != nil {
            global.GVA_LOG.Warn("并网点编号格式错误，跳过",
                zap.String("number", *setting.Number),
                zap.Error(err))
            continue
        }

        // 检查AVC是否启用
        if setting.AvcIsEnabled == nil || *setting.AvcIsEnabled == 0 {
            global.GVA_LOG.Debug("并网点AVC未启用，跳过",
                zap.Int("bwdNo", bwdNo),
                zap.String("name", func() string {
                    if setting.Name != nil {
                        return *setting.Name
                    }
                    return ""
                }()))
            continue
        }

        // 启动AVC
        if err := s.StartAVC(bwdNo); err != nil {
            global.GVA_LOG.Warn("自动启动AVC失败",
                zap.Int("bwdNo", bwdNo),
                zap.String("name", func() string {
                    if setting.Name != nil {
                        return *setting.Name
                    }
                    return ""
                }()),
                zap.Error(err))
        } else {
            successCount++
            global.GVA_LOG.Info("自动启动AVC成功",
                zap.Int("bwdNo", bwdNo),
                zap.String("name", func() string {
                    if setting.Name != nil {
                        return *setting.Name
                    }
                    return ""
                }()))
        }
    }

    global.GVA_LOG.Info("AVC自动启动完成",
        zap.Int("成功数量", successCount),
        zap.Int("总数量", len(settings)))
}

// StartAVC 启动AVC控制循环
func (s *avc) StartAVC(bwdNo int) error {
    // 检查是否已经在运行
    if _, exists := s.psidChans[bwdNo]; exists {
        return fmt.Errorf("电站%s的AVC控制已在运行", bwdNo)
    }

    // 获取AVC配置
    config, err := s.GetAVCConfig(bwdNo)
    if err != nil {
        return fmt.Errorf("获取AVC配置失败: %v", err)
    }

    // 创建停止通道
    stopChan := make(chan struct{})
    s.psidChans[bwdNo] = stopChan

    // 启动控制循环
    go s.avcControlLoop(bwdNo, config, stopChan)

    global.GVA_LOG.Info("AVC控制循环已启动", zap.Int("bwdNo", bwdNo))
    return nil
}

// StopAVC 停止AVC控制循环
func (s *avc) StopAVC(bwdNo int) error {
    stopChan, exists := s.psidChans[bwdNo]
    if !exists {
        return fmt.Errorf("电站%s的AVC控制未运行", bwdNo)
    }

    close(stopChan)
    delete(s.psidChans, bwdNo)

    global.GVA_LOG.Info("AVC控制循环已停止", zap.Int("bwdNo", bwdNo))
    return nil
}

// avcControlLoop AVC控制主循环
func (s *avc) avcControlLoop(bwdNo int, config agvc_main.AVCConfig, stopChan chan struct{}) {
    ticker := time.NewTicker(time.Duration(config.RegPeriod) * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            if err := s.executeAVCCycle(bwdNo); err != nil {
                global.GVA_LOG.Error("AVC控制周期执行失败",
                    zap.Int("bwdNo", bwdNo),
                    zap.Error(err))
            }
        case <-stopChan:
            global.GVA_LOG.Info("AVC控制循环退出", zap.Int("bwdNo", bwdNo))
            return
        }
    }
}

// executeAVCCycle 执行一个AVC控制周期
func (s *avc) executeAVCCycle(bwdNo int) error {
    // 步骤1：获取最新配置
    config, err := s.GetAVCConfig(bwdNo)
    if err != nil {
        return fmt.Errorf("获取AVC配置失败: %v", err)
    }

    // 步骤2：判断控制权限模式
    var isRemoteControl bool
    if config.ControlAuth != nil && *config.ControlAuth == 1 {
        isRemoteControl = true // 远程调度控制
    } else {
        isRemoteControl = false // 本地控制
    }

    // 步骤3：检查AVC投退信号
    var avcEnabled bool
    if isRemoteControl {
        // 远程模式：从内存读取AVC投退信号
        pointID, err := PointMapper.GetPointID(cons.TYPE_BWG, cons.AVC_YX_SIGNAL)
        if err != nil {
            global.GVA_LOG.Warn("获取AVC投退信号点位失败，使用数据库配置", zap.Error(err))
            avcEnabled = config.IsActive != nil && *config.IsActive == 1
        } else {
            signalVal, err := DispatchStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_BWG, cons.YX, pointID)
            if err != nil {
                global.GVA_LOG.Debug("从调度存储读取AVC投退信号失败，使用数据库配置", zap.Error(err))
                avcEnabled = config.IsActive != nil && *config.IsActive == 1
            } else {
                avcEnabled = signalVal == 1
                global.GVA_LOG.Debug("从调度存储读取AVC投退信号",
                    zap.Int("bwdNo", bwdNo),
                    zap.Float64("signalVal", signalVal),
                    zap.Bool("avcEnabled", avcEnabled))
            }
        }
    } else {
        // 本地模式：从数据库读取
        avcEnabled = config.IsActive != nil && *config.IsActive == 1
    }

    if !avcEnabled {
        global.GVA_LOG.Debug("AVC系统未投入", zap.Int("bwdNo", bwdNo), zap.Bool("isRemoteControl", isRemoteControl))
        return nil
    }

    // 步骤4：检查AVC就地远方控制模式
    if isRemoteControl {
        pointID, err := PointMapper.GetPointID(cons.TYPE_BWG, cons.AVC_YX_CONTROL_MODE)
        if err == nil {
            controlMode, err := DispatchStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_BWG, cons.YX, pointID)
            if err == nil && controlMode != 1 {
                global.GVA_LOG.Debug("AVC控制模式不是远程模式，跳过调控",
                    zap.Int("bwdNo", bwdNo),
                    zap.Float64("controlMode", controlMode))
                return nil
            }
        }
    }

    // 步骤5：检查AVC开/闭环状态
    var isOpenLoop bool
    if isRemoteControl {
        pointID, err := PointMapper.GetPointID(cons.TYPE_BWG, cons.AVC_YX_LOOP_STATUS)
        if err == nil {
            loopStatus, err := DispatchStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_BWG, cons.YX, pointID)
            if err == nil {
                isOpenLoop = loopStatus == 1
            } else {
                isOpenLoop = config.RunMode != nil && *config.RunMode == 1
            }
        } else {
            isOpenLoop = config.RunMode != nil && *config.RunMode == 1
        }
    } else {
        isOpenLoop = config.RunMode != nil && *config.RunMode == 1
    }

    // 步骤6：获取目标电压范围
    targetLow := config.TargetVoltageLow
    targetHigh := config.TargetVoltageHigh

    // 步骤7：采集并网点数据
    collectData, err := s.collectAVCData(bwdNo)
    if err != nil {
        return fmt.Errorf("采集数据失败: %v", err)
    }

    pointVoltage := collectData["pointVoltage"].(float64)
    totalReactive := collectData["totalReactive"].(float64)
    systemFreq := collectData["systemFreq"].(float64)

    global.GVA_LOG.Debug("AVC数据采集",
        zap.Int("bwdNo", bwdNo),
        zap.Bool("isOpenLoop", isOpenLoop),
        zap.Float64("并网点电压", pointVoltage),
        zap.Float64("总无功", totalReactive),
        zap.Float64("系统频率", systemFreq))

    // 步骤5：判断电压是否在合格范围
    if pointVoltage >= targetLow && pointVoltage <= targetHigh {
        global.GVA_LOG.Debug("电压在合格范围，无需调节",
            zap.Int("bwdNo", bwdNo),
            zap.Float64("电压", pointVoltage),
            zap.Float64("范围", targetLow),
            zap.Float64("到", targetHigh))
        return nil
    }

    // 步骤6：计算电压偏差
    var deltaV float64
    var targetVoltage float64
    if pointVoltage < targetLow {
        deltaV = targetLow - pointVoltage
        targetVoltage = targetLow
    } else {
        deltaV = targetHigh - pointVoltage
        targetVoltage = targetHigh
    }

    // 判断是否在死区内
    if math.Abs(deltaV) <= config.VoltageDeadZone {
        global.GVA_LOG.Debug("电压偏差在死区内，无需调节",
            zap.Int("bwdNo", bwdNo),
            zap.Float64("偏差", deltaV),
            zap.Float64("死区", config.VoltageDeadZone))
        return nil
    }

    // 步骤7：检查闭锁信号
    if deltaV > 0 && config.UpRegLock != nil && *config.UpRegLock == 1 {
        global.GVA_LOG.Warn("上调节闭锁，无法增加电压", zap.Int("bwdNo", bwdNo))
        return nil
    }
    if deltaV < 0 && config.DownRegLock != nil && *config.DownRegLock == 1 {
        global.GVA_LOG.Warn("下调节闭锁，无法降低电压", zap.Int("bwdNo", bwdNo))
        return nil
    }

    // 步骤8：计算所需无功调节量
    requiredDeltaQ := s.calcRequiredReactive(deltaV, pointVoltage, config.ReactiveSensitivity)

    global.GVA_LOG.Debug("AVC计算结果",
        zap.Int("bwdNo", bwdNo),
        zap.Float64("电压偏差", deltaV),
        zap.Float64("需求无功", requiredDeltaQ))

    // 步骤9：筛选可用设备（逆变器）
    availableDevices, err := s.filterAvailableDevices(bwdNo, requiredDeltaQ)
    if err != nil || len(availableDevices) == 0 {
        return fmt.Errorf("没有可用的调节设备: %v", err)
    }

    // 步骤10：平均分配无功调节量
    regulationDetails := s.assignReactiveToDevices(bwdNo, requiredDeltaQ, availableDevices)

    // 步骤11：记录调节开始
    record := agvc_main.AVCRegulationRecord{
        BwdNo:              bwdNo,
        TargetVoltage:      targetVoltage,
        ActualVoltage:      pointVoltage,
        VoltageDeviation:   deltaV,
        RequiredReactive:   requiredDeltaQ,
        RegulationReactive: requiredDeltaQ,
        TotalReactive:      totalReactive,
        SystemFreq:         systemFreq,
        Status:             "executing",
        Message:            fmt.Sprintf("开始调节，分配%d个设备", len(availableDevices)),
    }

    if err := global.GVA_DB.Create(&record).Error; err != nil {
        return fmt.Errorf("记录调节失败: %v", err)
    }

    // 步骤12：下发无功调节指令
    successCount := s.sendReactiveCommands(bwdNo, record.ID, regulationDetails)

    // 步骤13：更新调节记录
    if successCount == len(regulationDetails) {
        record.Status = "success"
        record.Message = fmt.Sprintf("调节成功，%d个设备全部响应", successCount)
    } else if successCount > 0 {
        record.Status = "partial"
        record.Message = fmt.Sprintf("部分调节成功，%d/%d个设备响应", successCount, len(regulationDetails))
    } else {
        record.Status = "failed"
        record.Message = "调节失败，所有设备无响应"
    }

    global.GVA_DB.Model(&record).Updates(map[string]interface{}{
        "status":  record.Status,
        "message": record.Message,
    })

    global.GVA_LOG.Info("AVC调节周期完成",
        zap.Int("bwdNo", bwdNo),
        zap.String("status", record.Status),
        zap.Int("成功数", successCount),
        zap.Int("总数", len(regulationDetails)))

    return nil
}

// collectAVCData 采集AVC控制所需数据
func (s *avc) collectAVCData(bwdNo int) (map[string]interface{}, error) {
    data := make(map[string]interface{})

    // 采集并网点电压（AB线电压）
    voltagePointID, err := PointMapper.GetPointID(cons.TYPE_BWG, cons.BWG_YC_VOLTAGE_AB)
    if err != nil {
        voltagePointID = "1" // 使用默认值
    }
    pointVoltage, err := DataStorage.GetDataAsFloat64(bwdNo, cons.TYPE_BWG, cons.YC, voltagePointID)
    if err != nil {
        return nil, fmt.Errorf("采集并网点电压失败: %v", err)
    }
    data["pointVoltage"] = pointVoltage

    // 采集总无功：优先从DataStorage获取，如无则从逆变器聚合
    totalReactive, err := PowerAggregator.GetBwgReactivePower(bwdNo)
    if err != nil {
        // 无功不是必须的，使用默认值
        totalReactive = 0.0
    }
    data["totalReactive"] = totalReactive

    // 采集系统频率
    freqPointID, err := PointMapper.GetPointID(cons.TYPE_BWG, cons.BWG_YC_FREQUENCY)
    if err != nil {
        freqPointID = "10" // 使用默认值
    }
    systemFreq, err := DataStorage.GetDataAsFloat64(bwdNo, cons.TYPE_BWG, cons.YC, freqPointID)
    if err != nil {
        systemFreq = 50.0
    }
    data["systemFreq"] = systemFreq

    return data, nil
}

// calcRequiredReactive 计算所需无功调节量
func (s *avc) calcRequiredReactive(deltaV, voltage, sensitivity float64) float64 {
    // 使用无功灵敏度计算
    // deltaQ = deltaV * sensitivity
    // 正值表示需要发出无功（增加电压），负值表示需要吸收无功（降低电压）
    return deltaV * sensitivity
}

// filterAvailableDevices 筛选可用设备（逆变器）
func (s *avc) filterAvailableDevices(bwdNo int, requiredQ float64) ([]agvc.AgvcNbqSetting, error) {
    // 获取所有在线逆变器
    inverters, err := Device.GetOnlineInvertersByBwdNo(bwdNo)
    if err != nil {
        return nil, err
    }

    available := make([]agvc.AgvcNbqSetting, 0)
    for _, inv := range inverters {
        // 检查设备是否有无功能力
        // if inv.MaxReact <= 0 {
        //     continue
        // }

        // 检查设备当前无功（可选）
        // currentReactive, _ := DataStorage.GetDataAsFloat64(psid, inv.EQID, "01", "02", "402")

        // 简单策略：所有有无功能力的设备都可用
        available = append(available, inv)
    }

    return available, nil
}

// assignReactiveToDevices 分配无功调节量到设备
func (s *avc) assignReactiveToDevices(nbqNo int, requiredQ float64, devices []agvc.AgvcNbqSetting) []agvc_main.DeviceReactiveRegulation {
    details := make([]agvc_main.DeviceReactiveRegulation, 0, len(devices))

    // 平均分配策略
    perDeviceQ := requiredQ / float64(len(devices))

    for _, dev := range devices {
        // 限制调节量不超过设备最大无功容量
        actualQ := perDeviceQ
        if math.Abs(actualQ) > *dev.RatedReactivePower {
            if actualQ > 0 {
                actualQ = *dev.RatedReactivePower
            } else {
                actualQ = -*dev.RatedReactivePower
            }
        }

        // 获取设备当前无功
        currentReactive, err := DataStorage.GetDataAsFloat64(nbqNo, cons.TYPE_BWG, cons.YC, "27")
        if err != nil {
            currentReactive = 0
        }

        psid := 1

        details = append(details, agvc_main.DeviceReactiveRegulation{
            PSID:               psid,
            EQID:               *dev.InverterNo,
            EQType:             cons.TYPE_NBQ,
            RegulationReactive: actualQ,
            BeforeReactive:     currentReactive,
            Status:             "pending",
        })
    }

    return details
}

// sendReactiveCommands 下发无功调节指令
func (s *avc) sendReactiveCommands(psid int, recordID uint, details []agvc_main.DeviceReactiveRegulation) int {
    successCount := 0
    host := CoapSender.GetDefaultCoapHost()
    port := CoapSender.GetDefaultCoapPort()

    for i := range details {
        detail := &details[i]
        detail.RecordID = recordID

        // 计算目标无功
        targetReactive := detail.BeforeReactive + detail.RegulationReactive

        // 构建指令
        commands := map[int]interface{}{
            402: targetReactive, // 无功功率补偿执行值（遥调）
        }

        // 发送CoAP指令
        if err := CoapSender.SendInverterCommand(host, port, psid, detail.EQID, commands); err != nil {
            global.GVA_LOG.Error("发送设备无功调节指令失败",
                zap.Int("eqid", detail.EQID),
                zap.Error(err))
            detail.Status = "failed"
            detail.AfterReactive = detail.BeforeReactive
        } else {
            detail.Status = "success"
            detail.AfterReactive = targetReactive
            successCount++
        }

        // 保存调节详情
        global.GVA_DB.Create(detail)
    }

    return successCount
}

// GetAVCConfig 获取AVC配置
func (s *avc) GetAVCConfig(bwdNo int) (agvc_main.AVCConfig, error) {
    var config agvc_main.AVCConfig
    err := global.GVA_DB.Where("psid = ?", bwdNo).First(&config).Error
    return config, err
}

// UpdateAVCConfig 更新AVC配置
func (s *avc) UpdateAVCConfig(req request.AVCConfigUpdate) error {
    updates := make(map[string]interface{})

    if req.IsActive != nil {
        updates["is_active"] = *req.IsActive
    }
    if req.ControlAuth != nil {
        updates["control_auth"] = *req.ControlAuth
    }
    if req.CommandStatus != nil {
        updates["command_status"] = *req.CommandStatus
    }
    if req.RunMode != nil {
        updates["run_mode"] = *req.RunMode
    }
    if req.TargetVoltageLow != 0 {
        updates["target_voltage_low"] = req.TargetVoltageLow
    }
    if req.TargetVoltageHigh != 0 {
        updates["target_voltage_high"] = req.TargetVoltageHigh
    }
    updates["voltage_exec_value"] = req.VoltageExecValue
    updates["reactive_exec_value"] = req.ReactiveExecValue
    updates["reactive_inc_cap"] = req.ReactiveIncCap
    updates["reactive_dec_cap"] = req.ReactiveDecCap
    if req.UpRegLock != nil {
        updates["up_reg_lock"] = *req.UpRegLock
    }
    if req.DownRegLock != nil {
        updates["down_reg_lock"] = *req.DownRegLock
    }
    if req.RegPeriod != 0 {
        updates["reg_period"] = req.RegPeriod
    }
    if req.VoltageDeadZone != 0 {
        updates["voltage_dead_zone"] = req.VoltageDeadZone
    }
    if req.ReactiveSensitivity != 0 {
        updates["reactive_sensitivity"] = req.ReactiveSensitivity
    }

    return global.GVA_DB.Model(&agvc_main.AVCConfig{}).
        Where("psid = ?", req.PSID).
        Updates(updates).Error
}

// GetAVCRecords 获取AVC调节记录
func (s *avc) GetAVCRecords(req request.AVCRegulationRecordSearch) ([]agvc_main.AVCRegulationRecord, int64, error) {
    var records []agvc_main.AVCRegulationRecord
    var total int64

    db := global.GVA_DB.Model(&agvc_main.AVCRegulationRecord{})

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
func (s *avc) CreateOrUpdateConfig(config *agvc_main.AVCConfig) error {
    var existing agvc_main.AVCConfig
    err := global.GVA_DB.Where("psid = ?", config.PSID).First(&existing).Error

    if err != nil {
        // 不存在，创建新配置
        return global.GVA_DB.Create(config).Error
    }

    // 存在，更新配置
    return global.GVA_DB.Model(&existing).Updates(config).Error
}
