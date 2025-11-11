package inverter_brands

import (
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
    "github.com/flipped-aurora/gin-vue-admin/server/service/agvc/cons"
)

// HuaweiInverterBrand 华为逆变器品牌实现
type HuaweiInverterBrand struct{}

// GetBrandCode 获取品牌代码
func (h *HuaweiInverterBrand) GetBrandCode() int {
    return cons.INVERTER_BRAND_HUAWEI
}

// GetBrandName 获取品牌名称
func (h *HuaweiInverterBrand) GetBrandName() string {
    return "华为"
}

// ControlAGCByPercentage 按百分比进行AGC调控
func (h *HuaweiInverterBrand) ControlAGCByPercentage(inverter *agvc.AgvcNbqSetting, targetPercentage float64) error {
    return HuaweiController.ControlAGCByPercentage(inverter, targetPercentage)
}

// ControlAGCByAbsolute 按绝对值进行AGC调控
func (h *HuaweiInverterBrand) ControlAGCByAbsolute(inverter *agvc.AgvcNbqSetting, targetPowerKW float64) error {
    return HuaweiController.ControlAGCByAbsolute(inverter, targetPowerKW)
}

// ControlAVCByPowerFactor 按功率因数进行AVC调控
func (h *HuaweiInverterBrand) ControlAVCByPowerFactor(inverter *agvc.AgvcNbqSetting, targetPF float64) error {
    return HuaweiController.ControlAVCByPowerFactor(inverter, targetPF)
}

// ControlAVCByReactivePower 按无功功率(Q/S)进行AVC调控
func (h *HuaweiInverterBrand) ControlAVCByReactivePower(inverter *agvc.AgvcNbqSetting, targetQS float64) error {
    return HuaweiController.ControlAVCByReactivePower(inverter, targetQS)
}

// InitializePoints 初始化华为逆变器点位映射
func (h *HuaweiInverterBrand) InitializePoints() error {
    return HuaweiPointInit.InitializeHuaweiPoints()
}
