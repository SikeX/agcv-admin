package inverter_brands

import (
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
	"go.uber.org/zap"
)

// 华为逆变器点位标识常量(501-515系列)
const (
	// AGC/AVC调节参数
	HW_POINT_REACTIVE_GRADIENT = "501" // 无功功率变化梯度(%/s)
	HW_POINT_ACTIVE_GRADIENT   = "502" // 有功功率变化梯度(%/s)
	HW_POINT_MAINTAIN_TIME     = "503" // 调度指令维持时间(s)

	// 有功调节
	HW_POINT_ACTIVE_KW      = "501" // 有功功率固定值降额(kW)
	HW_POINT_REACTIVE_PF    = "502" // 无功功率补偿(功率因数)
	HW_POINT_REACTIVE_QS    = "503" // 无功功率补偿(Q/S)
	HW_POINT_ACTIVE_PERCENT = "504" // 有功功率百分比降额(0.1%)
	HW_POINT_ACTIVE_W       = "505" // 有功功率固定值降额(W)

	// 调节模式和指令
	HW_POINT_ACTIVE_MODE    = "510" // [有功]调节模式
	HW_POINT_ACTIVE_VALUE   = "511" // [有功]调节值
	HW_POINT_ACTIVE_CMD     = "512" // [有功]调节指令
	HW_POINT_REACTIVE_MODE  = "513" // [无功]调节模式
	HW_POINT_REACTIVE_VALUE = "514" // [无功]调节值
	HW_POINT_REACTIVE_CMD   = "515" // [无功]调节指令
)

// HuaweiControlMode 华为逆变器有功调节模式
type HuaweiControlMode int

const (
	HuaweiModePercentage HuaweiControlMode = 0 // 百分比调节模式
	HuaweiModeAbsolute   HuaweiControlMode = 1 // 绝对值调节模式
)

// HuaweiReactiveMode 华为逆变器无功调节模式
type HuaweiReactiveMode int

const (
	HuaweiReactiveModePF HuaweiReactiveMode = 1 // 功率因数模式
	HuaweiReactiveModeQS HuaweiReactiveMode = 2 // Q/S模式
)

// HuaweiInverterController 华为逆变器控制器
type HuaweiInverterController struct{}

var HuaweiController = &HuaweiInverterController{}

// ControlAGC AGC调控统一入口 - 模式驱动控制
// @param inverter 逆变器配置
// @param mode 调节模式 (HuaweiModePercentage/HuaweiModeAbsolute)
// @param targetValue 目标值 (百分比模式:0-100, 绝对值模式:kW)
func (h *HuaweiInverterController) ControlAGC(inverter *agvc.AgvcNbqSetting, mode HuaweiControlMode, targetValue float64) error {
	if inverter.InverterNo == nil {
		return fmt.Errorf("逆变器编号为空")
	}

	invNo := *inverter.InverterNo
	host := GetCoapAdapter().GetDefaultCoapHost()
	port := GetCoapAdapter().GetDefaultCoapPort()
	psid := 1

	global.GVA_LOG.Info("华为逆变器AGC模式驱动调控开始",
		zap.Int("逆变器编号", invNo),
		zap.Int("调节模式", int(mode)),
		zap.Float64("目标值", targetValue))

	// 步骤1: 设置有功调节模式
	modeCommands := map[int]interface{}{
		401: int(mode),
	}

	if err := GetCoapAdapter().SendInverterCommand(host, port, psid, invNo, modeCommands); err != nil {
		return fmt.Errorf("设置有功调节模式失败: %v", err)
	}

	// 步骤2: 根据调节模式选择对应的调节值寄存器并设置值
	var valueCommands map[int]interface{}
	switch mode {
	case HuaweiModePercentage:
		// 百分比模式：限制范围0-100%，转换为0.1%单位
		if targetValue < 0 {
			targetValue = 0
		} else if targetValue > 100 {
			targetValue = 100
		}
		valueInt := int(targetValue * 10) // 转换为0.1%单位
		valueCommands = map[int]interface{}{
			404: valueInt,
		}
		global.GVA_LOG.Info("AGC百分比模式", zap.Float64("目标百分比", targetValue), zap.Int("寄存器值", valueInt))

	case HuaweiModeAbsolute:
		// 绝对值模式：限制不超过额定功率，转换为0.1kW单位
		if targetValue < 0 {
			targetValue = 0
		}
		if inverter.RatedActivePower != nil && targetValue > *inverter.RatedActivePower {
			targetValue = *inverter.RatedActivePower
		}
		valueInt := int(targetValue * 10) // 转换为0.1kW单位
		valueCommands = map[int]interface{}{
			401: valueInt, // 同时设置固定值降额寄存器
		}
		global.GVA_LOG.Info("AGC绝对值模式", zap.Float64("目标功率(kW)", targetValue), zap.Int("寄存器值", valueInt))

	default:
		return fmt.Errorf("不支持的有功调节模式: %d", mode)
	}

	if err := GetCoapAdapter().SendInverterCommand(host, port, psid, invNo, valueCommands); err != nil {
		return fmt.Errorf("设置有功调节值失败: %v", err)
	}

	// 步骤3: 发送有功调节指令
	execCommands := map[int]interface{}{
		402: 1, // 执行调节指令
	}

	if err := GetCoapAdapter().SendInverterCommand(host, port, psid, invNo, execCommands); err != nil {
		return fmt.Errorf("发送有功调节指令失败: %v", err)
	}

	global.GVA_LOG.Info("华为逆变器AGC模式驱动调控成功",
		zap.Int("逆变器编号", invNo),
		zap.Int("调节模式", int(mode)),
		zap.Float64("目标值", targetValue))

	return nil
}

// ControlAVC AVC调控统一入口 - 模式驱动控制
// @param inverter 逆变器配置
// @param mode 调节模式 (HuaweiReactiveModePF/HuaweiReactiveModeQS)
// @param targetValue 目标值 (功率因数模式:-1到1, Q/S模式:-1到1)
func (h *HuaweiInverterController) ControlAVC(inverter *agvc.AgvcNbqSetting, mode HuaweiReactiveMode, targetValue float64) error {
	if inverter.InverterNo == nil {
		return fmt.Errorf("逆变器编号为空")
	}

	invNo := *inverter.InverterNo
	host := GetCoapAdapter().GetDefaultCoapHost()
	port := GetCoapAdapter().GetDefaultCoapPort()
	psid := 1

	global.GVA_LOG.Info("华为逆变器AVC模式驱动调控开始",
		zap.Int("逆变器编号", invNo),
		zap.Int("调节模式", int(mode)),
		zap.Float64("目标值", targetValue))

	// 步骤1: 设置无功调节模式
	modeCommands := map[int]interface{}{
		401: int(mode),
	}

	if err := GetCoapAdapter().SendInverterCommand(host, port, psid, invNo, modeCommands); err != nil {
		return fmt.Errorf("设置无功调节模式失败: %v", err)
	}

	// 步骤2: 根据调节模式选择对应的调节值寄存器并设置值
	var valueCommands map[int]interface{}
	switch mode {
	case HuaweiReactiveModePF:
		// 功率因数模式：限制范围-1到1，转换为0.001单位
		if targetValue < -1 {
			targetValue = -1
		} else if targetValue > 1 {
			targetValue = 1
		}
		valueInt := int(targetValue * 1000) // 转换为0.001单位
		valueCommands = map[int]interface{}{
			405: valueInt, // 同时设置功率因数寄存器
		}
		global.GVA_LOG.Info("AVC功率因数模式", zap.Float64("目标功率因数", targetValue), zap.Int("寄存器值", valueInt))

	case HuaweiReactiveModeQS:
		// Q/S模式：限制范围-1到1，转换为0.001单位
		if targetValue < -1 {
			targetValue = -1
		} else if targetValue > 1 {
			targetValue = 1
		}
		valueInt := int(targetValue * 1000) // 转换为0.001单位
		valueCommands = map[int]interface{}{
			402: valueInt, // 同时设置Q/S寄存器
		}
		global.GVA_LOG.Info("AVC Q/S模式", zap.Float64("目标Q/S", targetValue), zap.Int("寄存器值", valueInt))

	default:
		return fmt.Errorf("不支持的无功调节模式: %d", mode)
	}

	if err := GetCoapAdapter().SendInverterCommand(host, port, psid, invNo, valueCommands); err != nil {
		return fmt.Errorf("设置无功调节值失败: %v", err)
	}

	// 步骤3: 发送无功调节指令
	execCommands := map[int]interface{}{
		402: 1, // 执行调节指令
	}

	if err := GetCoapAdapter().SendInverterCommand(host, port, psid, invNo, execCommands); err != nil {
		return fmt.Errorf("发送无功调节指令失败: %v", err)
	}

	global.GVA_LOG.Info("华为逆变器AVC模式驱动调控成功",
		zap.Int("逆变器编号", invNo),
		zap.Int("调节模式", int(mode)),
		zap.Float64("目标值", targetValue))

	return nil
}

// SetMaintainTime 设置调度指令维持时间
// @param inverter 逆变器配置
// @param maintainTimeSeconds 维持时间(秒)
// func (h *HuaweiInverterController) SetMaintainTime(inverter *agvc.AgvcNbqSetting, maintainTimeSeconds int) error {
//     if inverter.InverterNo == nil {
//         return fmt.Errorf("逆变器编号为空")
//     }

//     invNo := *inverter.InverterNo
//     host := GetCoapAdapter().GetDefaultCoapHost()
//     port := GetCoapAdapter().GetDefaultCoapPort()
//     psid := 1

//     if maintainTimeSeconds < 0 {
//         maintainTimeSeconds = 0
//     }

//     global.GVA_LOG.Info("设置华为逆变器调度指令维持时间",
//         zap.Int("逆变器编号", invNo),
//         zap.Int("维持时间(秒)", maintainTimeSeconds))

//     commands := map[int]interface{}{
//         HW_REG_MAINTAIN_TIME: maintainTimeSeconds,
//     }

//     if err := GetCoapAdapter().SendInverterCommand(host, port, psid, invNo, commands); err != nil {
//         return fmt.Errorf("设置调度指令维持时间失败: %v", err)
//     }

//     global.GVA_LOG.Info("华为逆变器调度指令维持时间设置成功",
//         zap.Int("逆变器编号", invNo),
//         zap.Int("维持时间(秒)", maintainTimeSeconds))

//     return nil
// }

// ControlAGCByPercentage 按百分比进行AGC调控
func (h *HuaweiInverterController) ControlAGCByPercentage(inverter *agvc.AgvcNbqSetting, targetPercentage float64) error {
	if inverter.InverterNo == nil {
		return fmt.Errorf("逆变器编号为空")
	}

	invNo := *inverter.InverterNo

	// 限制百分比范围 0-100%
	// if targetPercentage < 0 {
	// 	targetPercentage = 0
	// } else if targetPercentage > 100 {
	// 	targetPercentage = 100
	// }

	global.GVA_LOG.Info("华为逆变器AGC百分比调控",
		zap.Int("逆变器编号", invNo),
		zap.Float64("目标百分比", targetPercentage))

	// 1. 设置有功调节模式为百分比模式(0)
	// modeCommands := map[int]interface{}{
	//     HW_REG_ACTIVE_MODE: int(HuaweiModePercentage),
	// }

	host := GetCoapAdapter().GetDefaultCoapHost()
	port := GetCoapAdapter().GetDispatchBackCoapPort()
	psid := 1

	// if err := GetCoapAdapter().SendInverterCommand(host, port, psid, invNo, modeCommands); err != nil {
	//     return fmt.Errorf("设置调节模式失败: %v", err)
	// }

	// 2. 设置有功调节值
	valueInt := int(targetPercentage)

	valueCommands := map[int]interface{}{
		504: valueInt, // 有功功率百分比降额寄存器
	}

	if err := GetCoapAdapter().SendInverterCommand(host, port, psid, invNo, valueCommands); err != nil {
		return fmt.Errorf("设置调节值失败: %v", err)
	}

	global.GVA_LOG.Info("华为逆变器AGC百分比调控成功",
		zap.Int("逆变器编号", invNo),
		zap.Float64("目标百分比", targetPercentage),
		zap.Int("寄存器值", valueInt))

	return nil
}

// ControlAGCByAbsolute 按绝对值进行AGC调控
func (h *HuaweiInverterController) ControlAGCByAbsolute(inverter *agvc.AgvcNbqSetting, targetPowerKW float64) error {
	if inverter.InverterNo == nil {
		return fmt.Errorf("逆变器编号为空")
	}

	invNo := *inverter.InverterNo

	// 限制功率范围
	if targetPowerKW < 0 {
		targetPowerKW = 0
	}
	if inverter.RatedActivePower != nil && targetPowerKW > *inverter.RatedActivePower {
		targetPowerKW = *inverter.RatedActivePower
	}

	global.GVA_LOG.Info("华为逆变器AGC绝对值调控",
		zap.Int("逆变器编号", invNo),
		zap.Float64("目标功率(kW)", targetPowerKW))

	// 1. 设置有功调节模式为绝对值模式(1)
	// modeCommands := map[int]interface{}{
	//     HW_REG_ACTIVE_MODE: int(HuaweiModeAbsolute),
	// }

	host := GetCoapAdapter().GetDefaultCoapHost()
	port := GetCoapAdapter().GetDispatchBackCoapPort()
	psid := 1

	// if err := GetCoapAdapter().SendInverterCommand(host, port, psid, invNo, modeCommands); err != nil {
	//     return fmt.Errorf("设置调节模式失败: %v", err)
	// }

	// 2. 设置有功调节值(kW)
	valueInt := int(targetPowerKW)

	valueCommands := map[int]interface{}{
		501: valueInt, // 有功功率固定值降额(kW)寄存器
	}

	if err := GetCoapAdapter().SendInverterCommand(host, port, psid, invNo, valueCommands); err != nil {
		return fmt.Errorf("设置调节值失败: %v", err)
	}

	//// 3. 发送调节指令(执行)
	//execCommands := map[int]interface{}{
	//    501: 1, // 有功调节指令执行
	//}

	//if err := GetCoapAdapter().SendInverterCommand(host, port, psid, invNo, execCommands); err != nil {
	//    return fmt.Errorf("发送调节指令失败: %v", err)
	//}

	global.GVA_LOG.Info("华为逆变器AGC绝对值调控成功",
		zap.Int("逆变器编号", invNo),
		zap.Float64("目标功率(kW)", targetPowerKW),
		zap.Int("寄存器值", valueInt))

	return nil
}

// ControlAVCByPowerFactor 按功率因数进行AVC调控
func (h *HuaweiInverterController) ControlAVCByPowerFactor(inverter *agvc.AgvcNbqSetting, targetPF float64) error {
	if inverter.InverterNo == nil {
		return fmt.Errorf("逆变器编号为空")
	}

	invNo := *inverter.InverterNo

	// 限制功率因数范围 -1 到 1
	if targetPF < -1 {
		targetPF = -1
	} else if targetPF > 1 {
		targetPF = 1
	}

	global.GVA_LOG.Info("华为逆变器AVC功率因数调控",
		zap.Int("逆变器编号", invNo),
		zap.Float64("目标功率因数", targetPF))

	// 1. 设置无功调节模式为功率因数模式
	modeCommands := map[int]interface{}{
		401: int(HuaweiReactiveModePF), // 功率因数模式
	}

	host := GetCoapAdapter().GetDefaultCoapHost()
	port := GetCoapAdapter().GetDefaultCoapPort()
	psid := 1

	if err := GetCoapAdapter().SendInverterCommand(host, port, psid, invNo, modeCommands); err != nil {
		return fmt.Errorf("设置调节模式失败: %v", err)
	}

	// 2. 设置无功调节值(功率因数，0.001精度)
	valueInt := int(targetPF * 1000) // 转换为0.001单位

	valueCommands := map[int]interface{}{
		405: valueInt, // 无功功率补偿(功率因数)寄存器
	}

	if err := GetCoapAdapter().SendInverterCommand(host, port, psid, invNo, valueCommands); err != nil {
		return fmt.Errorf("设置调节值失败: %v", err)
	}

	// 3. 发送调节指令(执行)
	execCommands := map[int]interface{}{
		402: 1, // 无功调节指令执行
	}

	if err := GetCoapAdapter().SendInverterCommand(host, port, psid, invNo, execCommands); err != nil {
		return fmt.Errorf("发送调节指令失败: %v", err)
	}

	global.GVA_LOG.Info("华为逆变器AVC功率因数调控成功",
		zap.Int("逆变器编号", invNo),
		zap.Float64("目标功率因数", targetPF),
		zap.Int("寄存器值", valueInt))

	return nil
}

// ControlAVCByReactivePower 按无功功率(Q/S)进行AVC调控
func (h *HuaweiInverterController) ControlAVCByReactivePower(inverter *agvc.AgvcNbqSetting, targetQS float64) error {
	if inverter.InverterNo == nil {
		return fmt.Errorf("逆变器编号为空")
	}

	invNo := *inverter.InverterNo

	// 限制Q/S范围 -1 到 1
	if targetQS < -1 {
		targetQS = -1
	} else if targetQS > 1 {
		targetQS = 1
	}

	global.GVA_LOG.Info("华为逆变器AVC无功调控",
		zap.Int("逆变器编号", invNo),
		zap.Float64("目标Q/S", targetQS))

	// 1. 设置无功调节模式为Q/S模式
	modeCommands := map[int]interface{}{
		401: int(HuaweiReactiveModeQS), // Q/S模式
	}

	host := GetCoapAdapter().GetDefaultCoapHost()
	port := GetCoapAdapter().GetDefaultCoapPort()
	psid := 1

	if err := GetCoapAdapter().SendInverterCommand(host, port, psid, invNo, modeCommands); err != nil {
		return fmt.Errorf("设置调节模式失败: %v", err)
	}

	// 2. 设置无功调节值(Q/S，0.001精度)
	valueInt := int(targetQS * 1000) // 转换为0.001单位

	valueCommands := map[int]interface{}{
		405: valueInt, // 无功功率补偿(Q/S)寄存器
	}

	if err := GetCoapAdapter().SendInverterCommand(host, port, psid, invNo, valueCommands); err != nil {
		return fmt.Errorf("设置调节值失败: %v", err)
	}

	// 3. 发送调节指令(执行)
	execCommands := map[int]interface{}{
		402: 1, // 无功调节指令执行
	}

	if err := GetCoapAdapter().SendInverterCommand(host, port, psid, invNo, execCommands); err != nil {
		return fmt.Errorf("发送调节指令失败: %v", err)
	}

	global.GVA_LOG.Info("华为逆变器AVC无功调控成功",
		zap.Int("逆变器编号", invNo),
		zap.Float64("目标Q/S", targetQS),
		zap.Int("寄存器值", valueInt))

	return nil
}

// SetControlGradient 设置调节梯度
func (h *HuaweiInverterController) SetControlGradient(inverter *agvc.AgvcNbqSetting, activePowerGradient, reactivePowerGradient float64) error {
	if inverter.InverterNo == nil {
		return fmt.Errorf("逆变器编号为空")
	}

	invNo := *inverter.InverterNo

	host := GetCoapAdapter().GetDefaultCoapHost()
	port := GetCoapAdapter().GetDefaultCoapPort()
	psid := 1

	global.GVA_LOG.Info("设置华为逆变器调节梯度",
		zap.Int("逆变器编号", invNo),
		zap.Float64("有功梯度(%/s)", activePowerGradient),
		zap.Float64("无功梯度(%/s)", reactivePowerGradient))

	commands := make(map[int]interface{})

	// 设置有功功率变化梯度(0.001精度)
	if activePowerGradient > 0 {
		commands[402] = int(activePowerGradient * 1000)
	}

	// 设置无功功率变化梯度(0.001精度)
	if reactivePowerGradient > 0 {
		commands[405] = int(reactivePowerGradient * 1000)
	}

	if len(commands) > 0 {
		if err := GetCoapAdapter().SendInverterCommand(host, port, psid, invNo, commands); err != nil {
			return fmt.Errorf("设置调节梯度失败: %v", err)
		}
	}

	return nil
}

// CalculateInverterPowerDistribution 计算逆变器功率分配
func (h *HuaweiInverterController) CalculateInverterPowerDistribution(inverters []*agvc.AgvcNbqSetting, totalPowerKW float64, usePercentage bool) map[int]float64 {
	distribution := make(map[int]float64)

	if len(inverters) == 0 {
		return distribution
	}

	if usePercentage {
		// 百分比模式：所有逆变器设置相同百分比
		var totalRatedPower float64
		for _, inv := range inverters {
			if inv.RatedActivePower != nil {
				totalRatedPower += *inv.RatedActivePower
			}
		}

		if totalRatedPower > 0 {
			percentage := (totalPowerKW / totalRatedPower) * 100
			for _, inv := range inverters {
				if inv.InverterNo != nil {
					distribution[*inv.InverterNo] = percentage
				}
			}
		}
	} else {
		// 绝对值模式：按额定功率比例分配
		var totalRatedPower float64
		for _, inv := range inverters {
			if inv.RatedActivePower != nil {
				totalRatedPower += *inv.RatedActivePower
			}
		}

		if totalRatedPower > 0 {
			for _, inv := range inverters {
				if inv.InverterNo != nil && inv.RatedActivePower != nil {
					ratio := *inv.RatedActivePower / totalRatedPower
					distribution[*inv.InverterNo] = totalPowerKW * ratio
				}
			}
		} else {
			// 平均分配
			avgPower := totalPowerKW / float64(len(inverters))
			for _, inv := range inverters {
				if inv.InverterNo != nil {
					distribution[*inv.InverterNo] = avgPower
				}
			}
		}
	}

	return distribution
}

// ValidateControlValue 验证控制值是否在合理范围内
// func (h *HuaweiInverterController) ValidateControlValue(pointID string, value float64) error {
//     switch pointID {
//     case HW_POINT_ACTIVE_PERCENT: // 有功功率百分比降额
//         if value < 0 || value > 1000 { // 0-100% (0.1%精度)
//             return fmt.Errorf("有功功率百分比降额值超出范围[0-1000]: %.2f", value)
//         }
//     case HW_POINT_ACTIVE_KW: // 有功功率固定值降额(kW)
//         if value < 0 {
//             return fmt.Errorf("有功功率固定值降额值不能为负: %.2f", value)
//         }
//     case HW_POINT_REACTIVE_PF: // 功率因数
//         if value < -1000 || value > 1000 { // -1到1 (0.001精度)
//             return fmt.Errorf("功率因数值超出范围[-1000, 1000]: %.3f", value)
//         }
//     case HW_POINT_REACTIVE_QS: // Q/S
//         if value < -1000 || value > 1000 { // -1到1 (0.001精度)
//             return fmt.Errorf("Q/S值超出范围[-1000, 1000]: %.3f", value)
//         }
//     case HW_POINT_REACTIVE_GRADIENT, HW_POINT_ACTIVE_GRADIENT: // 变化梯度
//         if value < 0 {
//             return fmt.Errorf("变化梯度值不能为负: %.2f", value)
//         }
//     case HW_POINT_MAINTAIN_TIME: // 维持时间
//         if value < 0 {
//             return fmt.Errorf("维持时间不能为负: %.2f", value)
//         }
//     default:
//         return fmt.Errorf("不支持的点位ID: %s", pointID)
//     }
//     return nil
// }
