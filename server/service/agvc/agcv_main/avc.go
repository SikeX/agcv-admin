package agcv_main

import (
	"fmt"
	"math"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc/agvc_main"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc/agvc_main/request"
	"go.uber.org/zap"
)

type avc struct {
	psidChans map[string]chan struct{} // 每个电站一个停止通道
}

var AVC = new(avc)

// Initialize 初始化AVC服务
func (s *avc) Initialize() {
	s.psidChans = make(map[string]chan struct{})
	global.GVA_LOG.Info("AVC控制服务初始化成功")
}

// StartAVC 启动AVC控制循环
func (s *avc) StartAVC(psid string) error {
	// 检查是否已经在运行
	if _, exists := s.psidChans[psid]; exists {
		return fmt.Errorf("电站%s的AVC控制已在运行", psid)
	}

	// 获取AVC配置
	config, err := s.GetAVCConfig(psid)
	if err != nil {
		return fmt.Errorf("获取AVC配置失败: %v", err)
	}

	// 创建停止通道
	stopChan := make(chan struct{})
	s.psidChans[psid] = stopChan

	// 启动控制循环
	go s.avcControlLoop(psid, config, stopChan)

	global.GVA_LOG.Info("AVC控制循环已启动", zap.String("psid", psid))
	return nil
}

// StopAVC 停止AVC控制循环
func (s *avc) StopAVC(psid string) error {
	stopChan, exists := s.psidChans[psid]
	if !exists {
		return fmt.Errorf("电站%s的AVC控制未运行", psid)
	}

	close(stopChan)
	delete(s.psidChans, psid)

	global.GVA_LOG.Info("AVC控制循环已停止", zap.String("psid", psid))
	return nil
}

// avcControlLoop AVC控制主循环
func (s *avc) avcControlLoop(psid string, config agvc_main.AVCConfig, stopChan chan struct{}) {
	ticker := time.NewTicker(time.Duration(config.RegPeriod) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.executeAVCCycle(psid); err != nil {
				global.GVA_LOG.Error("AVC控制周期执行失败",
					zap.String("psid", psid),
					zap.Error(err))
			}
		case <-stopChan:
			global.GVA_LOG.Info("AVC控制循环退出", zap.String("psid", psid))
			return
		}
	}
}

// executeAVCCycle 执行一个AVC控制周期
func (s *avc) executeAVCCycle(psid string) error {
	// 步骤1：获取最新配置
	config, err := s.GetAVCConfig(psid)
	if err != nil {
		return fmt.Errorf("获取AVC配置失败: %v", err)
	}

	// 步骤2：检查AVC是否投入
	if config.IsActive == nil || *config.IsActive == 0 {
		global.GVA_LOG.Debug("AVC系统未投入", zap.String("psid", psid))
		return nil
	}

	// 步骤3：获取目标电压范围
	targetLow := config.TargetVoltageLow
	targetHigh := config.TargetVoltageHigh

	// 步骤4：采集并网点数据
	collectData, err := s.collectAVCData(psid)
	if err != nil {
		return fmt.Errorf("采集数据失败: %v", err)
	}

	pointVoltage := collectData["pointVoltage"].(float64)
	totalReactive := collectData["totalReactive"].(float64)
	systemFreq := collectData["systemFreq"].(float64)

	global.GVA_LOG.Debug("AVC数据采集",
		zap.String("psid", psid),
		zap.Float64("并网点电压", pointVoltage),
		zap.Float64("总无功", totalReactive),
		zap.Float64("系统频率", systemFreq))

	// 步骤5：判断电压是否在合格范围
	if pointVoltage >= targetLow && pointVoltage <= targetHigh {
		global.GVA_LOG.Debug("电压在合格范围，无需调节",
			zap.String("psid", psid),
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
			zap.String("psid", psid),
			zap.Float64("偏差", deltaV),
			zap.Float64("死区", config.VoltageDeadZone))
		return nil
	}

	// 步骤7：检查闭锁信号
	if deltaV > 0 && config.UpRegLock != nil && *config.UpRegLock == 1 {
		global.GVA_LOG.Warn("上调节闭锁，无法增加电压", zap.String("psid", psid))
		return nil
	}
	if deltaV < 0 && config.DownRegLock != nil && *config.DownRegLock == 1 {
		global.GVA_LOG.Warn("下调节闭锁，无法降低电压", zap.String("psid", psid))
		return nil
	}

	// 步骤8：计算所需无功调节量
	requiredDeltaQ := s.calcRequiredReactive(deltaV, pointVoltage, config.ReactiveSensitivity)

	global.GVA_LOG.Debug("AVC计算结果",
		zap.String("psid", psid),
		zap.Float64("电压偏差", deltaV),
		zap.Float64("需求无功", requiredDeltaQ))

	// 步骤9：筛选可用设备（逆变器）
	availableDevices, err := s.filterAvailableDevices(psid, requiredDeltaQ)
	if err != nil || len(availableDevices) == 0 {
		return fmt.Errorf("没有可用的调节设备: %v", err)
	}

	// 步骤10：平均分配无功调节量
	regulationDetails := s.assignReactiveToDevices(psid, requiredDeltaQ, availableDevices)

	// 步骤11：记录调节开始
	record := agvc_main.AVCRegulationRecord{
		PSID:               psid,
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
	successCount := s.sendReactiveCommands(psid, record.ID, regulationDetails)

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
		zap.String("psid", psid),
		zap.String("status", record.Status),
		zap.Int("成功数", successCount),
		zap.Int("总数", len(regulationDetails)))

	return nil
}

// collectAVCData 采集AVC控制所需数据
func (s *avc) collectAVCData(psid string) (map[string]interface{}, error) {
	data := make(map[string]interface{})

	// 从并网点采集数据
	eqid := "0000"
	eqType := "02" // 并网点

	// 采集并网点电压（遥测401：电压）
	pointVoltage, err := DataStorage.GetDataAsFloat64(psid, eqid, eqType, "02", "401")
	if err != nil {
		return nil, fmt.Errorf("采集并网点电压失败: %v", err)
	}
	data["pointVoltage"] = pointVoltage

	// 采集总无功（遥测403：无功功率）
	totalReactive, err := DataStorage.GetDataAsFloat64(psid, eqid, eqType, "02", "403")
	if err != nil {
		// 无功不是必须的，使用默认值
		totalReactive = 0.0
	}
	data["totalReactive"] = totalReactive

	// 采集系统频率
	systemFreq, err := DataStorage.GetDataAsFloat64(psid, eqid, eqType, "02", "404")
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

// filterAvailableDevices 筛选可用设备
func (s *avc) filterAvailableDevices(psid string, requiredQ float64) ([]agvc_main.Device, error) {
	// 获取所有在线逆变器
	inverters, err := Device.GetOnlineInvertersByPSID(psid)
	if err != nil {
		return nil, err
	}

	available := make([]agvc_main.Device, 0)
	for _, inv := range inverters {
		// 检查设备是否有无功能力
		if inv.MaxReact <= 0 {
			continue
		}

		// 检查设备当前无功（可选）
		// currentReactive, _ := DataStorage.GetDataAsFloat64(psid, inv.EQID, "01", "02", "402")

		// 简单策略：所有有无功能力的设备都可用
		available = append(available, inv)
	}

	return available, nil
}

// assignReactiveToDevices 分配无功调节量到设备
func (s *avc) assignReactiveToDevices(psid string, requiredQ float64, devices []agvc_main.Device) []agvc_main.DeviceReactiveRegulation {
	details := make([]agvc_main.DeviceReactiveRegulation, 0, len(devices))

	// 平均分配策略
	perDeviceQ := requiredQ / float64(len(devices))

	for _, dev := range devices {
		// 限制调节量不超过设备最大无功容量
		actualQ := perDeviceQ
		if math.Abs(actualQ) > dev.MaxReact {
			if actualQ > 0 {
				actualQ = dev.MaxReact
			} else {
				actualQ = -dev.MaxReact
			}
		}

		// 获取设备当前无功
		currentReactive, err := DataStorage.GetDataAsFloat64(psid, dev.EQID, "01", "02", "402")
		if err != nil {
			currentReactive = 0
		}

		details = append(details, agvc_main.DeviceReactiveRegulation{
			PSID:               psid,
			EQID:               dev.EQID,
			EQType:             dev.EQType,
			RegulationReactive: actualQ,
			BeforeReactive:     currentReactive,
			Status:             "pending",
		})
	}

	return details
}

// sendReactiveCommands 下发无功调节指令
func (s *avc) sendReactiveCommands(psid string, recordID uint, details []agvc_main.DeviceReactiveRegulation) int {
	successCount := 0
	host := CoapSender.GetDefaultCoapHost()
	port := CoapSender.GetDefaultCoapPort()

	for i := range details {
		detail := &details[i]
		detail.RecordID = recordID

		// 计算目标无功
		targetReactive := detail.BeforeReactive + detail.RegulationReactive

		// 构建指令
		commands := map[string]interface{}{
			"402_reactive": targetReactive, // 无功功率补偿执行值（遥调）
		}

		// 发送CoAP指令
		if err := CoapSender.SendInverterCommand(host, port, psid, detail.EQID, commands); err != nil {
			global.GVA_LOG.Error("发送设备无功调节指令失败",
				zap.String("eqid", detail.EQID),
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
func (s *avc) GetAVCConfig(psid string) (agvc_main.AVCConfig, error) {
	var config agvc_main.AVCConfig
	err := global.GVA_DB.Where("psid = ?", psid).First(&config).Error
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
