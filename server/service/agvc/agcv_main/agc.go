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

// AutoStartAllGridPoints 自动启动所有并网点的AGC功能
func (s *agc) AutoStartAllGridPoints() {
	// 查询所有并网点配置
	var settings []agvc.AgvcBwdSetting
	if err := global.GVA_DB.Find(&settings).Error; err != nil {
		global.GVA_LOG.Error("查询并网点配置失败", zap.Error(err))
		return
	}

	if len(settings) == 0 {
		global.GVA_LOG.Info("未找到并网点配置，跳过AGC自动启动")
		return
	}

	global.GVA_LOG.Info("开始自动启动所有并网点的AGC功能", zap.Int("并网点数量", len(settings)))

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
		////判断并网点调控模式
		//pointID, err := PointMapper.GetPointID(cons.TYPE_AGC, cons.AGC_YX_CONTROL_MODE)
		//if err != nil {
		//    global.GVA_LOG.Warn("获取AGC就地远方控制模式点位失败，使用数据库配置",
		//        zap.Error(err))
		//    continue
		//}
		//signalVal, err := DispatchStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_AGC, cons.YX, pointID)
		//if err != nil {
		//    global.GVA_LOG.Warn("获取AGC就地远方控制模式值失败，使用数据库配置",
		//        zap.Error(err))
		//    continue
		//}
		//
		////0:就地控制，1:远方控制
		//if signalVal == 0 {
		//    // 检查AGC是否启用
		//    if setting.AgcIsEnabled == nil || *setting.AgcIsEnabled == 0 {
		//        global.GVA_LOG.Debug("并网点AGC未启用，跳过",
		//            zap.Int("bwdNo", bwdNo),
		//            zap.String("name", func() string {
		//                if setting.Name != nil {
		//                    return *setting.Name
		//                }
		//                return ""
		//            }()))
		//        continue
		//    }
		//} else {
		//    AgcOnPointID, err := PointMapper.GetPointID(cons.TYPE_AGC, cons.AGC_YX_SIGNAL)
		//    if err != nil {
		//        global.GVA_LOG.Warn("获取AGC就地远方控制模式点位失败，使用数据库配置",
		//            zap.Error(err))
		//        continue
		//    }
		//    AgcOnPoint, err := DispatchStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_AGC, cons.YX, AgcOnPointID)
		//    if err != nil {
		//        global.GVA_LOG.Warn("获取AGC就地远方控制模式值失败，使用数据库配置",
		//            zap.Error(err))
		//        continue
		//    }
		//    // 检查AGC是否启用
		//    if AgcOnPoint == 0 {
		//        global.GVA_LOG.Debug("并网点AGC未启用，跳过",
		//            zap.Int("bwdNo", bwdNo),
		//            zap.String("name", func() string {
		//                if setting.Name != nil {
		//                    return *setting.Name
		//                }
		//                return ""
		//            }()))
		//        continue
		//    }
		//}

		// 启动AGC
		if err := s.StartAGC(bwdNo); err != nil {
			global.GVA_LOG.Warn("自动启动AGC失败",
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
			global.GVA_LOG.Info("自动启动AGC成功",
				zap.Int("bwdNo", bwdNo),
				zap.String("name", func() string {
					if setting.Name != nil {
						return *setting.Name
					}
					return ""
				}()))
		}
	}

	global.GVA_LOG.Info("AGC自动启动完成",
		zap.Int("成功数量", successCount),
		zap.Int("总数量", len(settings)))
}

// StartAGC 启动AGC控制循环
func (s *agc) StartAGC(bwdNo int) error {
	// 检查是否已经在运行
	if _, exists := s.bwdNoChans[bwdNo]; exists {
		return fmt.Errorf("电站%d的AGC控制已在运行", bwdNo)
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
		return fmt.Errorf("电站%d的AGC控制未运行", bwdNo)
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

	// 步骤2：判断控制权限模式
	var isRemoteControl bool

	pointID, err := PointMapper.GetPointID(cons.TYPE_AGC, cons.AGC_YX_CONTROL_MODE)
	if err != nil {
		global.GVA_LOG.Warn("获取AGC就地远方控制模式点位失败，使用数据库配置",
			zap.Error(err))
		// 从数据库配置读取控制模式
		if config.ControlAuth != nil && *config.ControlAuth == 1 {
			isRemoteControl = true
		} else {
			isRemoteControl = false
		}
	} else {
		signalVal, err := DispatchStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_AGC, cons.YX, pointID)
		if err != nil {
			global.GVA_LOG.Warn("获取AGC就地远方控制模式值失败，使用数据库配置",
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
	}

	global.GVA_LOG.Info("AGC控制模式判断",
		zap.Int("bwdNo", bwdNo),
		zap.Bool("isRemoteControl", isRemoteControl))

	// 如果是远程模式，无论调控过程是否成功，都必须发送结果到1189端口
	if isRemoteControl {
		// 初始化默认值，用于发送结果
		var actualOutput float64 = 0.0
		var targetOutput float64 = 0.0

		// 使用defer确保无论函数如何返回，都会发送结果
		defer func() {
			if sendErr := s.sendAGCResultToDispatch(bwdNo, config, actualOutput, targetOutput); sendErr != nil {
				global.GVA_LOG.Error("发送AGC结果到调度失败",
					zap.Int("bwdNo", bwdNo),
					zap.Error(sendErr))
			}
		}()

		// 执行远程调控逻辑
		return s.executeRemoteAGCControl(bwdNo, config, &actualOutput, &targetOutput)
	}

	// 执行就地调控逻辑
	return s.executeLocalAGCControl(bwdNo, config)
}

// executeRemoteAGCControl 执行远程AGC调控逻辑
func (s *agc) executeRemoteAGCControl(bwdNo int, config agvc.AgvcBwdSetting, actualOutput, targetOutput *float64) error {
	// 步骤1：检查AGC投退信号（从内存读取）
	var agcEnabled bool
	pointID, err := PointMapper.GetPointID(cons.TYPE_AGC, cons.AGC_YX_SIGNAL)
	if err != nil {
		global.GVA_LOG.Warn("获取AGC投退信号点位失败，使用数据库配置", zap.Error(err))
		agcEnabled = config.AgcIsEnabled != nil && *config.AgcIsEnabled == 1
	} else {
		signalVal, err := DispatchStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_AGC, cons.YX, pointID)
		if err != nil {
			global.GVA_LOG.Debug("从调度存储读取AGC投退信号失败，使用数据库配置", zap.Error(err))
			agcEnabled = config.AgcIsEnabled != nil && *config.AgcIsEnabled == 1
		} else {
			agcEnabled = signalVal == 1
			global.GVA_LOG.Debug("从调度存储读取AGC投退信号",
				zap.Int("bwdNo", bwdNo),
				zap.Float64("signalVal", signalVal),
				zap.Bool("agcEnabled", agcEnabled))
		}
	}

	if !agcEnabled {
		global.GVA_LOG.Debug("AGC系统未投入（远程模式）", zap.Int("bwdNo", bwdNo))
		return nil
	}

	// 步骤2：检查AGC就地远方控制模式
	pointID, err = PointMapper.GetPointID(cons.TYPE_AGC, cons.AGC_YX_CONTROL_MODE)
	if err == nil {
		controlMode, err := DispatchStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_AGC, cons.YX, pointID)
		if err == nil && controlMode != 1 {
			global.GVA_LOG.Debug("AGC控制模式不是远程模式，跳过调控",
				zap.Int("bwdNo", bwdNo),
				zap.Float64("controlMode", controlMode))
			return nil
		}
	}

	// 步骤3：检查AGC开/闭环状态
	var isOpenLoop bool
	isOpenLoop = true
	pointID, err = PointMapper.GetPointID(cons.TYPE_BWG, cons.AGC_YX_LOOP_STATUS)
	if err == nil {
		loopStatus, err := DispatchStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_AGC, cons.YX, pointID)
		if err == nil {
			isOpenLoop = loopStatus == 1
		}
	}

	// 步骤4：获取执行值（从内存获取有功执行值）
	var execVal float64
	pointID, err = PointMapper.GetPointID(cons.TYPE_AGC, cons.AGC_YC_POWER_EXEC_VALUE)
	if err != nil {
		pointID = "403" // 使用默认点标识
	}

	dispatchVal, err := DispatchStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_BWG, cons.YC, pointID)
	if err == nil {
		execVal = dispatchVal
		global.GVA_LOG.Debug("从调度存储获取AGC目标值",
			zap.Int("bwdNo", bwdNo),
			zap.Float64("execVal", execVal))
	} else {
		// 调度存储没有值，使用配置中的值
		if config.DispatchExecValue != nil {
			execVal = *config.DispatchExecValue
		}
		global.GVA_LOG.Debug("调度存储无数据，使用配置中的调度执行值",
			zap.Int("bwdNo", bwdNo),
			zap.Float64("execVal", execVal))
	}

	*targetOutput = execVal

	// 步骤5：采集实际出力
	actual, err := PowerAggregator.GetBwgActivePower(bwdNo)
	if err != nil {
		global.GVA_LOG.Warn("采集实际出力失败，使用默认值0",
			zap.Int("bwdNo", bwdNo),
			zap.Error(err))
		actual = 0.0
	}
	*actualOutput = actual

	// 步骤6：判断是否开环运行
	if isOpenLoop {
		// 开环运行，直接执行目标值
		global.GVA_LOG.Debug("AGC开环运行（远程模式），直接执行目标值",
			zap.Int("bwdNo", bwdNo),
			zap.Float64("execVal", execVal))
		return s.executeRemoteOpenLoopControl(bwdNo, execVal)
	}

	// 步骤7：闭环运行，计算偏差
	outputDeviation := *targetOutput - *actualOutput

	global.GVA_LOG.Debug("AGC数据采集（远程模式）",
		zap.Int("bwdNo", bwdNo),
		zap.Float64("目标出力", *targetOutput),
		zap.Float64("实际出力", *actualOutput),
		zap.Float64("出力偏差", outputDeviation))

	// 步骤8：判断偏差是否在抖动区间内
	if config.AgcVibrationRange != nil && math.Abs(outputDeviation) <= *config.AgcVibrationRange {
		global.GVA_LOG.Debug("偏差在抖动区间内，无需调节（远程模式）",
			zap.Int("bwdNo", bwdNo),
			zap.Float64("偏差", outputDeviation),
			zap.Float64("抖动区间", *config.AgcVibrationRange))
		return nil
	}

	// 步骤9：执行闭环调节
	return s.executeRemoteClosedLoopControl(bwdNo, outputDeviation)
}

// executeLocalAGCControl 执行就地AGC调控逻辑
func (s *agc) executeLocalAGCControl(bwdNo int, config agvc.AgvcBwdSetting) error {
	// 步骤1：检查AGC投退信号（从数据库读取）
	agcEnabled := config.AgcIsEnabled != nil && *config.AgcIsEnabled == 1

	if !agcEnabled {
		global.GVA_LOG.Debug("AGC系统未投入（就地模式）", zap.Int("bwdNo", bwdNo))
		return nil
	}

	// 步骤2：检查AGC开/闭环状态
	isOpenLoop := config.RunMode != nil && *config.RunMode == 1

	// 步骤3：获取执行值（从数据库配置读取）
	var execVal float64
	if config.StationExecValue != nil {
		execVal = *config.StationExecValue
	}

	// 步骤4：判断是否开环运行
	if isOpenLoop {
		// 开环运行，直接执行目标值
		global.GVA_LOG.Debug("AGC开环运行（就地模式），直接执行目标值",
			zap.Int("bwdNo", bwdNo),
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

	global.GVA_LOG.Debug("AGC数据采集（就地模式）",
		zap.Int("bwdNo", bwdNo),
		zap.Float64("目标出力", targetOutput),
		zap.Float64("实际出力", actualOutput),
		zap.Float64("出力偏差", outputDeviation))

	// 步骤7：判断偏差是否在抖动区间内
	if config.AgcVibrationRange != nil && math.Abs(outputDeviation) <= *config.AgcVibrationRange {
		global.GVA_LOG.Debug("偏差在抖动区间内，无需调节（就地模式）",
			zap.Int("bwdNo", bwdNo),
			zap.Float64("偏差", outputDeviation),
			zap.Float64("抖动区间", *config.AgcVibrationRange))
		return nil
	}

	// 步骤8：执行闭环调节
	return s.executeClosedLoopControl(bwdNo, outputDeviation)
}

// executeRemoteOpenLoopControl 执行远程开环控制
func (s *agc) executeRemoteOpenLoopControl(bwdNo int, execVal float64) error {
	// 获取可用逆变器
	inverters, err := Device.GetOnlineInvertersByBwdNo(bwdNo)
	if err != nil || len(inverters) == 0 {
		return fmt.Errorf("没有可用的逆变器: %v", err)
	}

	// 平均分配目标功率
	perInvPower := execVal / float64(len(inverters))

	for _, inv := range inverters {
		// 根据逆变器品牌选择不同的控制方式
		brand := cons.INVERTER_BRAND_HUAWEI // 默认华为
		if inv.InverterBrand != nil {
			brand = *inv.InverterBrand
		}

		// 华为逆变器使用专用控制器
		if brand == cons.INVERTER_BRAND_HUAWEI {
			// 使用绝对值模式控制
			if err := HuaweiController.ControlAGCByAbsolute(&inv, perInvPower); err != nil {
				global.GVA_LOG.Error("华为逆变器AGC控制失败",
					zap.Int("inverterNo", *inv.InverterNo),
					zap.Error(err))
				continue
			}
		} else {
			// 其他品牌逆变器使用原有方式
			activePowerPoint, err := InverterBrandMapper.GetActivePowerPoint(brand)
			if err != nil {
				global.GVA_LOG.Warn("获取逆变器品牌点位失败，使用默认点位",
					zap.Int("inverterNo", *inv.InverterNo),
					zap.Int("brand", brand),
					zap.Error(err))
				activePowerPoint = "401"
			}

			var pointID int
			fmt.Sscanf(activePowerPoint, "%d", &pointID)

			commands := map[int]interface{}{
				pointID: perInvPower,
			}

			host := CoapSender.GetDefaultCoapHost()
			port := CoapSender.GetDefaultCoapPort()
			psid := 1
			CoapSender.SendInverterCommand(host, port, psid, *inv.InverterNo, commands)
		}
	}

	return nil
}

// executeRemoteClosedLoopControl 执行远程闭环控制
func (s *agc) executeRemoteClosedLoopControl(bwdNo int, outputDeviation float64) error {
	// 获取可用逆变器
	inverters, err := Device.GetOnlineInvertersByBwdNo(bwdNo)
	if err != nil || len(inverters) == 0 {
		return fmt.Errorf("没有可用的逆变器: %v", err)
	}

	// 平均分配调节量
	perInvDeviation := outputDeviation / float64(len(inverters))

	for _, inv := range inverters {
		// 限制调节量不超过逆变器最大调节能力
		actualReg := perInvDeviation
		if inv.RatedActivePower != nil && math.Abs(actualReg) > *inv.RatedActivePower {
			if actualReg > 0 {
				actualReg = *inv.RatedActivePower
			} else {
				actualReg = -*inv.RatedActivePower
			}
		}

		// 获取逆变器当前功率
		currentPower, err := DataStorage.GetDataAsFloat64(*inv.InverterNo, 2, cons.YC, "28")
		if err != nil {
			currentPower = 0
		}

		// 计算目标功率
		targetPower := currentPower + actualReg

		// 根据逆变器品牌选择不同的控制方式
		brand := cons.INVERTER_BRAND_HUAWEI // 默认华为
		if inv.InverterBrand != nil {
			brand = *inv.InverterBrand
		}

		// 华为逆变器使用专用控制器
		if brand == cons.INVERTER_BRAND_HUAWEI {
			// 使用绝对值模式控制
			if err := HuaweiController.ControlAGCByAbsolute(&inv, targetPower); err != nil {
				global.GVA_LOG.Error("华为逆变器AGC闭环控制失败",
					zap.Int("inverterNo", *inv.InverterNo),
					zap.Float64("currentPower", currentPower),
					zap.Float64("targetPower", targetPower),
					zap.Error(err))
				continue
			}
		} else {
			// 其他品牌逆变器使用原有方式
			activePowerPoint, err := InverterBrandMapper.GetActivePowerPoint(brand)
			if err != nil {
				global.GVA_LOG.Warn("获取逆变器品牌点位失败，使用默认点位",
					zap.Int("inverterNo", *inv.InverterNo),
					zap.Int("brand", brand),
					zap.Error(err))
				activePowerPoint = "401"
			}

			var pointID int
			fmt.Sscanf(activePowerPoint, "%d", &pointID)

			// 构建指令
			commands := map[int]interface{}{
				pointID: targetPower,
			}

			// 发送CoAP指令
			host := CoapSender.GetDefaultCoapHost()
			port := CoapSender.GetDefaultCoapPort()
			psid := 1
			CoapSender.SendInverterCommand(host, port, psid, *inv.InverterNo, commands)
		}

	}
	return nil
}

// executeClosedLoopControl 执行就地闭环控制
func (s *agc) executeClosedLoopControl(bwdNo int, outputDeviation float64) error {
	// 获取可用逆变器
	inverters, err := Device.GetOnlineInvertersByBwdNo(bwdNo)
	if err != nil || len(inverters) == 0 {
		return fmt.Errorf("没有可用的逆变器: %v", err)
	}

	// 平均分配调节量
	perInvDeviation := outputDeviation / float64(len(inverters))
	regulationDetails := make([]agvc_main.InverterRegulation, 0, len(inverters))

	for _, inv := range inverters {
		// 限制调节量不超过逆变器最大调节能力
		actualReg := perInvDeviation
		if inv.RatedActivePower != nil && math.Abs(actualReg) > *inv.RatedActivePower {
			if actualReg > 0 {
				actualReg = *inv.RatedActivePower
			} else {
				actualReg = -*inv.RatedActivePower
			}
		}

		// 获取逆变器当前功率
		currentPower, err := DataStorage.GetDataAsFloat64(*inv.InverterNo, 2, cons.YC, "28")
		if err != nil {
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

	// 下发调节指令
	successCount := 0
	for i := range regulationDetails {
		detail := &regulationDetails[i]

		// 计算目标功率
		targetPower := detail.BeforePower + detail.RegulationPower

		// 根据逆变器品牌获取控制点位
		// 需要获取逆变器配置来确定品牌
		var brand int = cons.INVERTER_BRAND_HUAWEI // 默认华为
		for _, inv := range inverters {
			if *inv.InverterNo == detail.EQID && inv.InverterBrand != nil {
				brand = *inv.InverterBrand
				break
			}
		}

		activePowerPoint, err := InverterBrandMapper.GetActivePowerPoint(brand)
		if err != nil {
			global.GVA_LOG.Warn("获取逆变器品牌点位失败，使用默认点位",
				zap.Int("inverterNo", detail.EQID),
				zap.Int("brand", brand),
				zap.Error(err))
			activePowerPoint = "401"
		}

		var pointID int
		fmt.Sscanf(activePowerPoint, "%d", &pointID)

		// 构建指令
		commands := map[int]interface{}{
			pointID: targetPower,
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
	}

	global.GVA_LOG.Info("AGC调节周期完成（就地模式）",
		zap.Int("bwdNo", bwdNo),
		zap.Int("成功数", successCount),
		zap.Int("总数", len(regulationDetails)))

	return nil
}

// executeOpenLoopControl 执行开环控制
func (s *agc) executeOpenLoopControl(bwdNo int, execVal float64) error {
	// 获取并网点配置
	config, err := s.GetAGCConfig(bwdNo)
	if err != nil {
		return fmt.Errorf("获取并网点配置失败: %v", err)
	}

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

	// 采集实际出力
	actualOutput, _ := PowerAggregator.GetBwgActivePower(bwdNo)

	// 发送AGC计算结果到调度
	if err := s.sendAGCResultToDispatch(bwdNo, config, actualOutput, execVal); err != nil {
		global.GVA_LOG.Error("发送AGC结果到调度失败（开环）",
			zap.Int("bwdNo", bwdNo),
			zap.Error(err))
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
	pointID, err := PointMapper.GetPointID(cons.TYPE_BWG, cons.BWG_YC_FREQUENCY)
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

// GetAGCConfig 获取AGC配置（优先从缓存获取）
func (s *agc) GetAGCConfig(bwdNo int) (agvc.AgvcBwdSetting, error) {
	// 优先从缓存获取
	if config, exists := SettingCache.Get(bwdNo); exists {
		return config, nil
	}

	// 缓存不存在，从数据库获取
	var config agvc.AgvcBwdSetting
	err := global.GVA_DB.Where("number = ?", fmt.Sprintf("%d", bwdNo)).First(&config).Error
	if err != nil {
		return config, err
	}

	// 加载到缓存
	SettingCache.Set(bwdNo, config)
	return config, nil
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

// sendAGCResultToDispatch 发送AGC计算结果到调度
func (s *agc) sendAGCResultToDispatch(bwdNo int, config agvc.AgvcBwdSetting, actualOutput, targetOutput float64) error {
	results := make(map[string]interface{})

	// AGC遥信标准点
	// 401: AGC投退信号
	pointID, _ := PointMapper.GetPointID(cons.TYPE_AGC, cons.AGC_YX_SIGNAL)
	if pointID == "" {
		pointID = "401"
	}
	agcSignal, err := DispatchStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_AGC, cons.YK, pointID)
	if err != nil {
		agcSignal = 0
	}
	results["agcSignal"] = agcSignal

	// 402: AGC就地远方控制模式 (0=本地, 1=远程)
	pointID, _ = PointMapper.GetPointID(cons.TYPE_AGC, cons.AGC_YX_CONTROL_MODE)
	if pointID == "" {
		pointID = "402"
	}
	controlMode, err := DispatchStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_AGC, cons.YK, pointID)
	if err != nil {
		controlMode = 0
	}
	results["agcControlMode"] = controlMode

	// 404: AGC开/闭环状态 (0=闭环, 1=开环)
	pointID, _ = PointMapper.GetPointID(cons.TYPE_AGC, cons.AGC_YX_LOOP_STATUS)
	if pointID == "" {
		pointID = "404"
	}
	loopStatus, err := DispatchStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_AGC, cons.YK, pointID)
	if err != nil {
		loopStatus = 0
	}
	results["agcLoopStatus"] = loopStatus

	// 405: AGC有功上调节闭锁（从调度存储读取）
	// pointID, _ = PointMapper.GetPointID(cons.TYPE_BWG, cons.AGC_YX_UP_REG_LOCK)
	// if pointID == "" {
	//     pointID = "405"
	// }
	// upRegLock, err := DispatchStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_BWG, cons.YX, pointID)
	// if err != nil {
	//     upRegLock = 0
	// }
	// results["agcUpRegLock"] = upRegLock

	// 406: AGC有功下调节闭锁（从调度存储读取）
	// pointID, _ = PointMapper.GetPointID(cons.TYPE_BWG, cons.AGC_YX_DOWN_REG_LOCK)
	// if pointID == "" {
	//     pointID = "406"
	// }
	// downRegLock, err := DispatchStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_BWG, cons.YX, pointID)
	// if err != nil {
	//     downRegLock = 0
	// }
	// results["agcDownRegLock"] = downRegLock

	// AGC遥测标准点
	// 401: 有功调节上限（从调度存储读取）
	pointID, _ = PointMapper.GetPointID(cons.TYPE_BWG, cons.AGC_YC_POWER_UPPER_LIMIT)
	if pointID == "" {
		pointID = "401"
	}
	powerUpperLimit, err := DispatchStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_BWG, cons.YC, pointID)
	if err != nil {
		powerUpperLimit = 0
	}
	results["powerUpperLimit"] = powerUpperLimit

	// 402: 有功调节下限（从调度存储读取）
	pointID, _ = PointMapper.GetPointID(cons.TYPE_BWG, cons.AGC_YC_POWER_LOWER_LIMIT)
	if pointID == "" {
		pointID = "402"
	}
	powerLowerLimit, err := DispatchStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_BWG, cons.YK, pointID)
	if err != nil {
		powerLowerLimit = 0
	}
	results["powerLowerLimit"] = powerLowerLimit

	// 403: 有功执行值（当前实际输出功率）
	pointID, _ = PointMapper.GetPointID(cons.TYPE_AGC, cons.AGC_YC_POWER_EXEC_VALUE)
	if pointID == "" {
		pointID = "401"
	}
	powerExecValue, err := DispatchStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_AGC, cons.YT, pointID)
	if err != nil {
		powerExecValue = 0
	}
	results["powerExecValue"] = powerExecValue

	fmt.Println("AGC结果:", results)

	// 发送到调度
	return CoapSender.SendAGCResultToDispatch(bwdNo, results)
}
