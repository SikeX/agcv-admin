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
		return fmt.Errorf("电站%d的AVC控制已在运行", bwdNo)
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
		return fmt.Errorf("电站%d的AVC控制未运行", bwdNo)
	}

	close(stopChan)
	delete(s.psidChans, bwdNo)

	global.GVA_LOG.Info("AVC控制循环已停止", zap.Int("bwdNo", bwdNo))
	return nil
}

// avcControlLoop AVC控制主循环
func (s *avc) avcControlLoop(bwdNo int, config agvc.AgvcBwdSetting, stopChan chan struct{}) {
	ticker := time.NewTicker(time.Duration(*config.AvcStepPeriod) * time.Second)
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

	// 步骤2：判断控制权限模式（从数据库配置读取）
	var isRemoteControl bool
	pointID, err := PointMapper.GetPointID(cons.TYPE_AVC, cons.AVC_YX_CONTROL_MODE)
	if err != nil {
		pointID = "402" // 使用默认点标识
	}
	signalVal, err := DispatchStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_AVC, cons.YK, pointID)
	if err != nil {
		global.GVA_LOG.Warn("获取AVC就地远方控制模式值失败，使用数据库配置",
			zap.Error(err))
		// 从数据库配置读取控制模式
		if config.ControlAuth != nil && *config.ControlAuth == 1 {
			isRemoteControl = true
		} else {
			isRemoteControl = false
		}
	} else {
		isRemoteControl = signalVal == 1 //0:就地控制，1:远方控制
	}

	global.GVA_LOG.Info("AVC控制模式判断",
		zap.Int("bwdNo", bwdNo),
		zap.Bool("isRemoteControl", isRemoteControl))

	// 如果是远程模式，无论调控过程是否成功，都必须发送结果到1189端口
	if isRemoteControl {
		// 初始化默认值，用于发送结果
		var actualVoltage float64 = 0.0
		var actualReactive float64 = 0.0

		// 使用defer确保无论函数如何返回，都会发送结果
		defer func() {
			if sendErr := s.sendAVCResultToDispatch(bwdNo, config, actualVoltage, actualReactive); sendErr != nil {
				global.GVA_LOG.Error("发送AVC结果到调度失败",
					zap.Int("bwdNo", bwdNo),
					zap.Error(sendErr))
			}
		}()

		// 执行远程调控逻辑
		return s.executeRemoteAVCControl(bwdNo, config)
	}

	// 执行就地调控逻辑
	return s.executeLocalAVCControl(bwdNo, config)
}

// executeRemoteAVCControl 执行远程AVC调控逻辑
func (s *avc) executeRemoteAVCControl(bwdNo int, config agvc.AgvcBwdSetting) error {
	// 步骤1：检查AVC投退信号（从内存读取）
	var avcEnabled bool
	pointID, err := PointMapper.GetPointID(cons.TYPE_AVC, cons.AVC_YX_SIGNAL)
	if err != nil {
		pointID = "401" // 使用默认点标识
	}
	signalVal, err := DataStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_AVC, cons.YX, pointID)
	if err != nil {
		global.GVA_LOG.Debug("从调度存储读取AVC投退信号失败，使用数据库配置", zap.Error(err))
		avcEnabled = config.AvcIsEnabled != nil && *config.AvcIsEnabled == 1
	} else {
		avcEnabled = signalVal == 1
		global.GVA_LOG.Debug("从调度存储读取AVC投退信号",
			zap.Int("bwdNo", bwdNo),
			zap.Float64("signalVal", signalVal),
			zap.Bool("avcEnabled", avcEnabled))
	}

	if !avcEnabled {
		global.GVA_LOG.Debug("AVC系统未投入（远程模式）", zap.Int("bwdNo", bwdNo))
		return nil
	}

	// 步骤2：检查AVC就地远方控制模式
	pointID, err = PointMapper.GetPointID(cons.TYPE_AVC, cons.AVC_YX_CONTROL_MODE)
	if err == nil {
		controlMode, err := DataStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_AVC, cons.YX, pointID)
		if err == nil && controlMode != 1 {
			global.GVA_LOG.Debug("AVC控制模式不是远程模式，跳过调控",
				zap.Int("bwdNo", bwdNo),
				zap.Float64("controlMode", controlMode))
			return nil
		}
	}

	actualVoltage := 0.0
	actualReactive := 0.0

	// 步骤3：采集并网点数据
	// collectData, err := s.collectAVCData(bwdNo)
	// if err != nil {
	// 	global.GVA_LOG.Warn("采集AVC数据失败，使用默认值",
	// 		zap.Int("bwdNo", bwdNo),
	// 		zap.Error(err))

	// }

	//从内存获取电压执行值和无功执行值
	pointID, err = PointMapper.GetPointID(cons.TYPE_AVC, cons.AVC_YC_VOLTAGE_EXEC_VALUE)
	if err != nil {
		pointID = "403" // 使用默认点标识
	}
	voltageExec, err := DataStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_AVC, cons.YX, pointID)
	if err == nil {
		actualVoltage = voltageExec
	}
	pointID, err = PointMapper.GetPointID(cons.TYPE_AVC, cons.AVC_YC_REACTIVE_EXEC_VALUE)
	if err != nil {
		pointID = "404" // 使用默认点标识
	}
	reactiveExec, err := DataStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_AVC, cons.YX, pointID)
	if err == nil {
		actualReactive = reactiveExec
	}

	// 步骤4：检查AVC开/闭环状态
	var isOpenLoop bool
	pointID, err = PointMapper.GetPointID(cons.TYPE_AVC, cons.AVC_YX_LOOP_STATUS)
	if err == nil {
		loopStatus, err := DataStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_AVC, cons.YX, pointID)
		if err == nil {
			isOpenLoop = loopStatus == 1
		} else {
			isOpenLoop = config.RunMode != nil && *config.RunMode == 1
		}
	} else {
		isOpenLoop = config.RunMode != nil && *config.RunMode == 1
	}

	global.GVA_LOG.Debug("AVC数据采集（远程模式）",
		zap.Int("bwdNo", bwdNo),
		zap.Bool("isOpenLoop", isOpenLoop),
		zap.Float64("并网点电压", actualVoltage),
		zap.Float64("总无功", actualReactive))

	// 步骤5：获取目标电压范围
	targetLow := config.AvcAdjustmentRangeMin
	targetHigh := config.AvcAdjustmentRangeMax

	// 步骤6：判断电压是否在合格范围
	if actualVoltage >= *targetLow && actualVoltage <= *targetHigh {
		global.GVA_LOG.Debug("电压在合格范围，无需调节（远程模式）",
			zap.Int("bwdNo", bwdNo),
			zap.Float64("电压", actualVoltage))
		return nil
	}

	// 步骤7：计算电压偏差
	var deltaV float64
	if actualVoltage < *targetLow {
		deltaV = *targetLow - actualVoltage
	} else {
		deltaV = *targetHigh - actualVoltage
	}

	// 判断是否在死区内
	if math.Abs(deltaV) <= *config.AvcVibrationRange {
		global.GVA_LOG.Debug("电压偏差在死区内，无需调节（远程模式）",
			zap.Int("bwdNo", bwdNo),
			zap.Float64("偏差", deltaV))
		return nil
	}

	// 步骤8：检查闭锁信号
	//闭锁pointId
	pointID = "405"
	upRegLock, err := DispatchStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_AVC, cons.YX, pointID)
	if err != nil {
		upRegLock = 0
		global.GVA_LOG.Error("获取AVC闭锁信号失败，",
			zap.Error(err))
	}
	if deltaV > 0 && upRegLock == 1 {
		global.GVA_LOG.Warn("上调节闭锁，无法增加电压（远程模式）", zap.Int("bwdNo", bwdNo))
		return nil
	}

	pointID = "406"
	downRegLock, err := DispatchStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_AVC, cons.YX, pointID)
	if err != nil {
		downRegLock = 0
		global.GVA_LOG.Error("获取AVC闭锁信号失败，",
			zap.Error(err))
	}
	if deltaV < 0 && downRegLock == 1 {
		global.GVA_LOG.Warn("下调节闭锁，无法降低电压（远程模式）", zap.Int("bwdNo", bwdNo))
		return nil
	}

	// 步骤9：计算所需无功调节量
	requiredDeltaQ := s.calcRequiredReactive(deltaV, actualVoltage, *config.AvcSystemImpedance)

	// 步骤10：执行无功调节
	return s.executeRemoteReactiveControl(bwdNo, requiredDeltaQ)
}

// executeLocalAVCControl 执行就地AVC调控逻辑
func (s *avc) executeLocalAVCControl(bwdNo int, config agvc.AgvcBwdSetting) error {
	// 步骤1：检查AVC投退信号（从数据库读取）
	avcEnabled := config.AvcIsEnabled != nil && *config.AvcIsEnabled == 1

	if !avcEnabled {
		global.GVA_LOG.Debug("AVC系统未投入（就地模式）", zap.Int("bwdNo", bwdNo))
		return nil
	}

	// 步骤2：采集并网点数据
	collectData, err := s.collectAVCData(bwdNo)
	if err != nil {
		return fmt.Errorf("采集数据失败: %v", err)
	}

	pointVoltage := collectData["pointVoltage"].(float64)
	totalReactive := collectData["totalReactive"].(float64)

	global.GVA_LOG.Debug("AVC数据采集（就地模式）",
		zap.Int("bwdNo", bwdNo),
		zap.Float64("并网点电压", pointVoltage),
		zap.Float64("总无功", totalReactive))

	// 步骤3：获取目标电压范围
	targetLow := config.AvcAdjustmentRangeMin
	targetHigh := config.AvcAdjustmentRangeMax

	// 步骤4：判断电压是否在合格范围
	if pointVoltage >= *targetLow && pointVoltage <= *targetHigh {
		global.GVA_LOG.Debug("电压在合格范围，无需调节（就地模式）",
			zap.Int("bwdNo", bwdNo),
			zap.Float64("电压", pointVoltage))
		return nil
	}

	// 步骤5：计算电压偏差
	var deltaV float64
	if pointVoltage < *targetLow {
		deltaV = *targetLow - pointVoltage
	} else {
		deltaV = *targetHigh - pointVoltage
	}

	// 判断是否在死区内
	if math.Abs(deltaV) <= *config.AvcVibrationRange {
		global.GVA_LOG.Debug("电压偏差在死区内，无需调节（就地模式）",
			zap.Int("bwdNo", bwdNo),
			zap.Float64("偏差", deltaV))
		return nil
	}

	// 步骤7：计算所需无功调节量
	requiredDeltaQ := s.calcRequiredReactive(deltaV, pointVoltage, *config.AvcSystemImpedance)

	// 步骤8：执行无功调节
	return s.executeReactiveControl(bwdNo, requiredDeltaQ)
}

// executeRemoteReactiveControl 执行远程无功调节
func (s *avc) executeRemoteReactiveControl(bwdNo int, requiredDeltaQ float64) error {
	// 筛选可用设备（逆变器）
	availableDevices, err := s.filterAvailableDevices(bwdNo, requiredDeltaQ)
	if err != nil || len(availableDevices) == 0 {
		return fmt.Errorf("没有可用的调节设备: %v", err)
	}

	// 平均分配无功调节量
	perDeviceQ := requiredDeltaQ / float64(len(availableDevices))

	for _, dev := range availableDevices {
		// 限制调节量不超过设备最大无功容量
		actualQ := perDeviceQ
		if dev.RatedReactivePower != nil && math.Abs(actualQ) > *dev.RatedReactivePower {
			if actualQ > 0 {
				actualQ = *dev.RatedReactivePower
			} else {
				actualQ = -*dev.RatedReactivePower
			}
		}

		// 获取设备当前无功
		currentReactive, err := DataStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_BWG, cons.YC, "27")
		if err != nil {
			currentReactive = 0
		}

		// 计算目标无功
		targetReactive := currentReactive + actualQ

		// 根据逆变器品牌选择不同的控制方式
		brand := cons.INVERTER_BRAND_HUAWEI // 默认华为
		if dev.InverterBrand != nil {
			brand = *dev.InverterBrand
		}

		// 华为逆变器使用专用控制器，通过Q/S模式调节
		if brand == cons.INVERTER_BRAND_HUAWEI {
			// 计算目标Q/S比值
			// 获取设备当前视在功率
			apparentPower, err := DataStorage.GetDataAsFloat64(1, *dev.InverterNo, cons.TYPE_NBQ, cons.YC, "202")
			if err != nil || apparentPower == 0 {
				global.GVA_LOG.Warn("无法获取视在功率，跳过该逆变器",
					zap.Int("inverterNo", *dev.InverterNo),
					zap.Error(err))
				continue
			}

			targetQS := targetReactive / apparentPower

			// 使用华为控制器进行AVC调控
			if err := HuaweiController.ControlAVCByReactivePower(&dev, targetQS); err != nil {
				global.GVA_LOG.Error("华为逆变器AVC控制失败",
					zap.Int("inverterNo", *dev.InverterNo),
					zap.Float64("targetQS", targetQS),
					zap.Error(err))
				continue
			}
		} else {
			// 其他品牌逆变器使用原有方式
			reactivePowerPoint, err := InverterBrandMapper.GetReactivePowerPoint(brand)
			if err != nil {
				global.GVA_LOG.Warn("获取逆变器品牌点位失败，使用默认点位",
					zap.Int("inverterNo", *dev.InverterNo),
					zap.Int("brand", brand),
					zap.Error(err))
				reactivePowerPoint = "402"
			}

			var pointID int
			fmt.Sscanf(reactivePowerPoint, "%d", &pointID)

			// 构建指令
			commands := map[int]interface{}{
				pointID: targetReactive,
			}

			// 发送CoAP指令
			host := CoapSender.GetDefaultCoapHost()
			port := CoapSender.GetDefaultCoapPort()
			psid := 1
			CoapSender.SendInverterCommand(host, port, psid, *dev.InverterNo, commands)
		}
	}

	return nil
}

// executeReactiveControl 执行就地无功调节
func (s *avc) executeReactiveControl(bwdNo int, requiredDeltaQ float64) error {
	// 筛选可用设备（逆变器）
	availableDevices, err := s.filterAvailableDevices(bwdNo, requiredDeltaQ)
	if err != nil || len(availableDevices) == 0 {
		return fmt.Errorf("没有可用的调节设备: %v", err)
	}

	// 平均分配无功调节量
	regulationDetails := s.assignReactiveToDevices(bwdNo, requiredDeltaQ, availableDevices)

	// 下发无功调节指令
	successCount := s.sendReactiveCommands(bwdNo, 0, regulationDetails)

	global.GVA_LOG.Info("AVC调节周期完成（就地模式）",
		zap.Int("bwdNo", bwdNo),
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
	pointVoltage, err := DataStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_BWG, cons.YC, voltagePointID)
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
	systemFreq, err := DataStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_BWG, cons.YC, freqPointID)
	if err != nil {
		systemFreq = 50.0
	}
	data["systemFreq"] = systemFreq

	return data, nil
}

// calcRequiredReactive 计算所需无功调节量
func (s *avc) calcRequiredReactive(deltaV, voltage, systemFreq float64) float64 {
	// 使用无功灵敏度计算
	// deltaQ = deltaV * V_real / 系统电抗
	// 正值表示需要发出无功（增加电压），负值表示需要吸收无功（降低电压）
	return deltaV * voltage / systemFreq
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
		currentReactive, err := DataStorage.GetDataAsFloat64(1, nbqNo, cons.TYPE_BWG, cons.YC, "27")
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

		// 获取逆变器品牌并确定控制点位
		var brand int = cons.INVERTER_BRAND_HUAWEI // 默认华为
		var invSetting agvc.AgvcNbqSetting
		if err := global.GVA_DB.Where("inverter_no = ?", detail.EQID).First(&invSetting).Error; err == nil {
			if invSetting.InverterBrand != nil {
				brand = *invSetting.InverterBrand
			}
		}

		reactivePowerPoint, err := InverterBrandMapper.GetReactivePowerPoint(brand)
		if err != nil {
			global.GVA_LOG.Warn("获取逆变器品牌点位失败，使用默认点位",
				zap.Int("inverterNo", detail.EQID),
				zap.Int("brand", brand),
				zap.Error(err))
			reactivePowerPoint = "402"
		}

		var pointID int
		fmt.Sscanf(reactivePowerPoint, "%d", &pointID)

		// 构建指令
		commands := map[int]interface{}{
			pointID: targetReactive, // 无功功率补偿执行值（遥调）
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
func (s *avc) GetAVCConfig(bwdNo int) (agvc.AgvcBwdSetting, error) {
	var config agvc.AgvcBwdSetting
	err := global.GVA_DB.Where("number = ?", bwdNo).First(&config).Error
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

	return global.GVA_DB.Model(&agvc.AgvcBwdSetting{}).
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
func (s *avc) CreateOrUpdateConfig(config *agvc.AgvcBwdSetting) error {
	var existing agvc.AgvcBwdSetting
	err := global.GVA_DB.Where("number = ?", config.Number).First(&existing).Error

	if err != nil {
		// 不存在，创建新配置
		return global.GVA_DB.Create(config).Error
	}

	// 存在，更新配置
	return global.GVA_DB.Model(&existing).Updates(config).Error
}

// sendAVCResultToDispatch 发送AVC计算结果到调度
func (s *avc) sendAVCResultToDispatch(bwdNo int, config agvc.AgvcBwdSetting, actualVoltage, actualReactive float64) error {
	results := make(map[string]interface{})

	// AVC遥信标准点
	// 401: AVC功能投退信号
	pointID, _ := PointMapper.GetPointID(cons.TYPE_AVC, cons.AVC_YX_SIGNAL)
	if pointID == "" {
		pointID = "401"
	}
	avcSignal, err := DispatchStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_AVC, cons.YK, pointID)
	if err != nil {
		avcSignal = 0
	}
	results["avcSignal"] = avcSignal

	// 402: AVC功能就地远方控制模式 (0=本地, 1=远程)
	pointID, _ = PointMapper.GetPointID(cons.TYPE_AVC, cons.AVC_YX_CONTROL_MODE)
	if pointID == "" {
		pointID = "402"
	}
	avcControlMode, err := DispatchStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_AVC, cons.YK, pointID)
	if err != nil {
		avcControlMode = 0
	}
	results["avcControlMode"] = avcControlMode

	// 403: AVC功能当前指令状态（0=无指令, 1=有指令）
	pointID, _ = PointMapper.GetPointID(cons.TYPE_AVC, cons.AVC_YX_COMMAND_STATUS)
	if pointID == "" {
		pointID = "403"
	}
	avcCmdStatus, err := DispatchStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_AVC, cons.YK, pointID)
	if err != nil {
		avcCmdStatus = 0
	}
	results["avcCmdStatus"] = avcCmdStatus

	// 404: AVC功能开闭环状态 (0=闭环, 1=开环)
	pointID, _ = PointMapper.GetPointID(cons.TYPE_AVC, cons.AVC_YX_LOOP_STATUS)
	if pointID == "" {
		pointID = "404"
	}
	avcLoopStatus, err := DispatchStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_AVC, cons.YK, pointID)
	if err != nil {
		avcLoopStatus = 0
	}
	results["avcLoopStatus"] = avcLoopStatus

	// 405: AVC功能上调节闭锁
	// if config.UpRegLock != nil {
	//     results["avcUpRegLock"] = float64(*config.UpRegLock)
	// } else {
	//     results["avcUpRegLock"] = float64(0)
	// }

	// // 406: AVC功能下调节闭锁
	// if config.DownRegLock != nil {
	//     results["avcDownRegLock"] = float64(*config.DownRegLock)
	// } else {
	//     results["avcDownRegLock"] = float64(0)
	// }

	// AVC遥测标准点
	// 401: 无功可增容量
	totalCapacity := float64(0)
	inverters, err := Device.GetOnlineInvertersByBwdNo(bwdNo)
	if err == nil {
		for _, inv := range inverters {
			if inv.RatedReactivePower != nil {
				totalCapacity += *inv.RatedReactivePower
			}
		}
	}
	results["reactiveIncreaseCap"] = totalCapacity

	// 402: 无功可减容量
	results["reactiveDecreaseCap"] = totalCapacity

	// 403: 电压执行值（当前并网点电压）
	pointID, _ = PointMapper.GetPointID(cons.TYPE_AVC, cons.AVC_YC_VOLTAGE_EXEC_VALUE)
	if pointID == "" {
		pointID = "403"
	}
	voltageExecValue, err := DispatchStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_AVC, cons.YT, pointID)
	if err != nil {
		voltageExecValue = 0
	}
	results["voltageExecValue"] = voltageExecValue

	// 404: 无功执行值（当前总无功）
	pointID, _ = PointMapper.GetPointID(cons.TYPE_AVC, cons.AVC_YC_REACTIVE_EXEC_VALUE)
	if pointID == "" {
		pointID = "404"
	}
	reactiveExecValue, err := DispatchStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_AVC, cons.YT, pointID)
	if err != nil {
		reactiveExecValue = 0
	}
	results["reactiveExecValue"] = reactiveExecValue

	fmt.Println("AVC结果数据:", results)

	// 发送到调度
	return CoapSender.SendAVCResultToDispatch(bwdNo, results)
}
