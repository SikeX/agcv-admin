package agcv_main

import (
	"fmt"
	"math"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
	"go.uber.org/zap"
)

// HuaweiInverterPoint 华为逆变器控制点位定义(5xx系列)
type HuaweiInverterPoint struct {
	PointID       string
	PointName     string
	Unit          string
	Multiplier    float64
	Offset        float64
	RegisterAddr  int
	RegisterCount int
	Description   string
}

// 华为逆变器控制点位映射表(5xx系列点标识)
var HuaweiInverterPoints = map[string]HuaweiInverterPoint{
	// AGC/AVC调节参数
	"501": {PointID: "501", PointName: "无功功率变化梯度", Unit: "%/s", Multiplier: 0.001, Offset: 0, RegisterAddr: 42015, RegisterCount: 5, Description: "无功功率变化梯度(%/s)"},
	"502": {PointID: "502", PointName: "有功功率变化梯度", Unit: "%/s", Multiplier: 0.001, Offset: 0, RegisterAddr: 42017, RegisterCount: 5, Description: "有功功率变化梯度(%/s)"},
	"503": {PointID: "503", PointName: "调度指令维持时间", Unit: "s", Multiplier: 1, Offset: 0, RegisterAddr: 42019, RegisterCount: 5, Description: "调度指令维持时间(s)"},

	// 有功调节
	"504": {PointID: "504", PointName: "有功功率固定值降额(kW)", Unit: "kW", Multiplier: 0.1, Offset: 0, RegisterAddr: 40120, RegisterCount: 1, Description: "有功功率固定值降额(kW)"},
	"507": {PointID: "507", PointName: "有功功率百分比降额", Unit: "0.1%", Multiplier: 0.1, Offset: 0, RegisterAddr: 40125, RegisterCount: 1, Description: "有功功率百分比降额(0.1%)"},
	"508": {PointID: "508", PointName: "有功功率固定值降额(W)", Unit: "W", Multiplier: 1, Offset: 0, RegisterAddr: 40126, RegisterCount: 5, Description: "有功功率固定值降额(W)"},

	// 无功调节
	"505": {PointID: "505", PointName: "无功功率补偿(PF)", Unit: "PF", Multiplier: 0.001, Offset: 0, RegisterAddr: 40122, RegisterCount: 2, Description: "无功功率补偿(功率因数)"},
	"506": {PointID: "506", PointName: "无功功率补偿(Q/S)", Unit: "Q/S", Multiplier: 0.001, Offset: 0, RegisterAddr: 40123, RegisterCount: 2, Description: "无功功率补偿(Q/S)"},

	// 调节模式和指令
	"510": {PointID: "510", PointName: "有功调节模式", Unit: "", Multiplier: 1, Offset: 0, RegisterAddr: 35300, RegisterCount: 1, Description: "[有功]调节模式"},
	"511": {PointID: "511", PointName: "有功调节值", Unit: "", Multiplier: 1, Offset: 0, RegisterAddr: 35301, RegisterCount: 5, Description: "[有功]调节值"},
	"512": {PointID: "512", PointName: "有功调节指令", Unit: "", Multiplier: 1, Offset: 0, RegisterAddr: 35303, RegisterCount: 1, Description: "[有功]调节指令"},
	"513": {PointID: "513", PointName: "无功调节模式", Unit: "", Multiplier: 1, Offset: 0, RegisterAddr: 35304, RegisterCount: 1, Description: "[无功]调节模式"},
	"514": {PointID: "514", PointName: "无功调节值", Unit: "", Multiplier: 1, Offset: 0, RegisterAddr: 35305, RegisterCount: 5, Description: "[无功]调节值"},
	"515": {PointID: "515", PointName: "无功调节指令", Unit: "", Multiplier: 1, Offset: 0, RegisterAddr: 35307, RegisterCount: 1, Description: "[无功]调节指令"},
}

// yx
const (
	HW_NBQ_YG_MODE  = "501" //[有功]调节模式
	HW_NBQ_YG_VALUE = "502" //[有功]调节值
)

// HuaweiControlMode 华为逆变器调节模式
type HuaweiControlMode int

const (
	HuaweiModePercentage HuaweiControlMode = 0 // 百分比调节模式
	HuaweiModeAbsolute   HuaweiControlMode = 1 // 绝对值调节模式
)

// HuaweiInverterController 华为逆变器控制器
type HuaweiInverterController struct{}

var HuaweiController = &HuaweiInverterController{}

// ControlAGCByPercentage 按百分比进行AGC调控
func (h *HuaweiInverterController) ControlAGCByPercentage(inverter *agvc.AgvcNbqSetting, targetPercentage float64) error {
	if inverter.InverterNo == nil {
		return fmt.Errorf("逆变器编号为空")
	}

	invNo := *inverter.InverterNo

	// 限制百分比范围 0-100%
	if targetPercentage < 0 {
		targetPercentage = 0
	} else if targetPercentage > 100 {
		targetPercentage = 100
	}

	global.GVA_LOG.Info("华为逆变器AGC百分比调控",
		zap.Int("逆变器编号", invNo),
		zap.Float64("目标百分比", targetPercentage))

	// 1. 设置有功调节模式为百分比模式(0)
	modeCommands := map[int]interface{}{
		35300: int(HuaweiModePercentage),
	}

	host := CoapSender.GetDefaultCoapHost()
	port := CoapSender.GetDefaultCoapPort()
	psid := 1

	if err := CoapSender.SendInverterCommand(host, port, psid, invNo, modeCommands); err != nil {
		return fmt.Errorf("设置调节模式失败: %v", err)
	}

	// 2. 设置有功调节值(百分比，0.1%精度，所以需要乘以10)
	valueInt := int(targetPercentage * 10) // 转换为0.1%单位

	valueCommands := map[int]interface{}{
		40125: valueInt, // 有功功率百分比降额寄存器
	}

	if err := CoapSender.SendInverterCommand(host, port, psid, invNo, valueCommands); err != nil {
		return fmt.Errorf("设置调节值失败: %v", err)
	}

	// 3. 发送调节指令(执行)
	execCommands := map[int]interface{}{
		35303: 1, // 有功调节指令执行
	}

	if err := CoapSender.SendInverterCommand(host, port, psid, invNo, execCommands); err != nil {
		return fmt.Errorf("发送调节指令失败: %v", err)
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
	modeCommands := map[int]interface{}{
		35300: int(HuaweiModeAbsolute),
	}

	host := CoapSender.GetDefaultCoapHost()
	port := CoapSender.GetDefaultCoapPort()
	psid := 1

	if err := CoapSender.SendInverterCommand(host, port, psid, invNo, modeCommands); err != nil {
		return fmt.Errorf("设置调节模式失败: %v", err)
	}

	// 2. 设置有功调节值(kW，0.1精度)
	valueInt := int(targetPowerKW * 10) // 转换为0.1kW单位

	valueCommands := map[int]interface{}{
		40120: valueInt, // 有功功率固定值降额(kW)寄存器
	}

	if err := CoapSender.SendInverterCommand(host, port, psid, invNo, valueCommands); err != nil {
		return fmt.Errorf("设置调节值失败: %v", err)
	}

	// 3. 发送调节指令(执行)
	execCommands := map[int]interface{}{
		35303: 1, // 有功调节指令执行
	}

	if err := CoapSender.SendInverterCommand(host, port, psid, invNo, execCommands); err != nil {
		return fmt.Errorf("发送调节指令失败: %v", err)
	}

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
		35304: 1, // 功率因数模式
	}

	host := CoapSender.GetDefaultCoapHost()
	port := CoapSender.GetDefaultCoapPort()
	psid := 1

	if err := CoapSender.SendInverterCommand(host, port, psid, invNo, modeCommands); err != nil {
		return fmt.Errorf("设置调节模式失败: %v", err)
	}

	// 2. 设置无功调节值(功率因数，0.001精度)
	valueInt := int(targetPF * 1000) // 转换为0.001单位

	valueCommands := map[int]interface{}{
		40122: valueInt, // 无功功率补偿(功率因数)寄存器
	}

	if err := CoapSender.SendInverterCommand(host, port, psid, invNo, valueCommands); err != nil {
		return fmt.Errorf("设置调节值失败: %v", err)
	}

	// 3. 发送调节指令(执行)
	execCommands := map[int]interface{}{
		35307: 1, // 无功调节指令执行
	}

	if err := CoapSender.SendInverterCommand(host, port, psid, invNo, execCommands); err != nil {
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
		35304: 2, // Q/S模式
	}

	host := CoapSender.GetDefaultCoapHost()
	port := CoapSender.GetDefaultCoapPort()
	psid := 1

	if err := CoapSender.SendInverterCommand(host, port, psid, invNo, modeCommands); err != nil {
		return fmt.Errorf("设置调节模式失败: %v", err)
	}

	// 2. 设置无功调节值(Q/S，0.001精度)
	valueInt := int(targetQS * 1000) // 转换为0.001单位

	valueCommands := map[int]interface{}{
		40123: valueInt, // 无功功率补偿(Q/S)寄存器
	}

	if err := CoapSender.SendInverterCommand(host, port, psid, invNo, valueCommands); err != nil {
		return fmt.Errorf("设置调节值失败: %v", err)
	}

	// 3. 发送调节指令(执行)
	execCommands := map[int]interface{}{
		35307: 1, // 无功调节指令执行
	}

	if err := CoapSender.SendInverterCommand(host, port, psid, invNo, execCommands); err != nil {
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

	host := CoapSender.GetDefaultCoapHost()
	port := CoapSender.GetDefaultCoapPort()
	psid := 1

	global.GVA_LOG.Info("设置华为逆变器调节梯度",
		zap.Int("逆变器编号", invNo),
		zap.Float64("有功梯度(%/s)", activePowerGradient),
		zap.Float64("无功梯度(%/s)", reactivePowerGradient))

	commands := make(map[int]interface{})

	// 设置有功功率变化梯度(0.001精度)
	if activePowerGradient > 0 {
		commands[42017] = int(activePowerGradient * 1000)
	}

	// 设置无功功率变化梯度(0.001精度)
	if reactivePowerGradient > 0 {
		commands[42015] = int(reactivePowerGradient * 1000)
	}

	if len(commands) > 0 {
		if err := CoapSender.SendInverterCommand(host, port, psid, invNo, commands); err != nil {
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

// GetPointDefinition 获取华为逆变器点位定义
func (h *HuaweiInverterController) GetPointDefinition(pointID string) (HuaweiInverterPoint, error) {
	point, exists := HuaweiInverterPoints[pointID]
	if !exists {
		return HuaweiInverterPoint{}, fmt.Errorf("未找到点位定义: %s", pointID)
	}
	return point, nil
}

// GetAllPointDefinitions 获取所有华为逆变器点位定义
func (h *HuaweiInverterController) GetAllPointDefinitions() map[string]HuaweiInverterPoint {
	return HuaweiInverterPoints
}

// ValidateControlValue 验证控制值是否在合理范围内
func (h *HuaweiInverterController) ValidateControlValue(pointID string, value float64) error {
	point, err := h.GetPointDefinition(pointID)
	if err != nil {
		return err
	}

	switch pointID {
	case "507": // 有功功率百分比降额
		if value < 0 || value > 1000 { // 0-100% (0.1%精度)
			return fmt.Errorf("%s值超出范围[0-1000]: %.2f", point.PointName, value)
		}
	case "504": // 有功功率固定值降额(kW)
		if value < 0 {
			return fmt.Errorf("%s值不能为负: %.2f", point.PointName, value)
		}
	case "505": // 功率因数
		realValue := value * point.Multiplier
		if realValue < -1 || realValue > 1 {
			return fmt.Errorf("%s值超出范围[-1, 1]: %.3f", point.PointName, realValue)
		}
	case "506": // Q/S
		realValue := value * point.Multiplier
		if realValue < -1 || realValue > 1 {
			return fmt.Errorf("%s值超出范围[-1, 1]: %.3f", point.PointName, realValue)
		}
	case "501", "502": // 变化梯度
		if value < 0 {
			return fmt.Errorf("%s值不能为负: %.2f", point.PointName, value)
		}
	}

	return nil
}

// ConvertToRegisterValue 将实际值转换为寄存器值
func (h *HuaweiInverterController) ConvertToRegisterValue(pointID string, realValue float64) (int, error) {
	point, err := h.GetPointDefinition(pointID)
	if err != nil {
		return 0, err
	}

	// 应用倍率和偏移
	registerValue := int(math.Round((realValue - point.Offset) / point.Multiplier))

	return registerValue, nil
}

// ConvertFromRegisterValue 将寄存器值转换为实际值
func (h *HuaweiInverterController) ConvertFromRegisterValue(pointID string, registerValue int) (float64, error) {
	point, err := h.GetPointDefinition(pointID)
	if err != nil {
		return 0, err
	}

	// 应用倍率和偏移
	realValue := float64(registerValue)*point.Multiplier + point.Offset

	return realValue, nil
}
