package agcv_main

import (
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/service/agvc/cons"
)

// InverterBrandPointMapping 逆变器品牌点位映射
type InverterBrandPointMapping struct {
	Brand                    int    // 品牌代码
	BrandName                string // 品牌名称
	PowerSwitchYK            string // 开关机遥控点位
	ActivePowerLimitYT       string // 有功功率降额遥调点位
	ReactivePowerYT          string // 无功功率补偿遥调点位
	ActivePowerLimitPercetYT string // 有功功率降额百分比遥调点位
}

// inverterBrandMapper 逆变器品牌映射器
type inverterBrandMapper struct {
	mappings map[int]InverterBrandPointMapping
}

// Initialize 初始化逆变器品牌点位映射
func (m *inverterBrandMapper) Initialize() {
	// 华为逆变器点位映射
	m.mappings[cons.INVERTER_BRAND_HUAWEI] = InverterBrandPointMapping{
		Brand:                    cons.INVERTER_BRAND_HUAWEI,
		BrandName:                "华为",
		PowerSwitchYK:            "401", // 开关机
		ActivePowerLimitYT:       "501", // 有功功率降额执行值
		ReactivePowerYT:          "503", // 无功功率补偿执行值
		ActivePowerLimitPercetYT: "503", // 有功功率降额百分比遥调点位
	}

	// 阳光逆变器点位映射
	m.mappings[cons.INVERTER_BRAND_SUNGROW] = InverterBrandPointMapping{
		Brand:              cons.INVERTER_BRAND_SUNGROW,
		BrandName:          "阳光",
		PowerSwitchYK:      "401", // 启停控制
		ActivePowerLimitYT: "401", // 有功设定值
		ReactivePowerYT:    "402", // 无功设定值
	}

	// 固德威逆变器点位映射
	m.mappings[cons.INVERTER_BRAND_GOODWE] = InverterBrandPointMapping{
		Brand:              cons.INVERTER_BRAND_GOODWE,
		BrandName:          "固德威",
		PowerSwitchYK:      "401", // 运行控制
		ActivePowerLimitYT: "401", // 有功功率控制
		ReactivePowerYT:    "402", // 无功功率控制
	}
}

// GetBrandMapping 获取品牌点位映射
func (m *inverterBrandMapper) GetBrandMapping(brand int) (InverterBrandPointMapping, error) {
	mapping, exists := m.mappings[brand]
	if !exists {
		return InverterBrandPointMapping{}, fmt.Errorf("未找到品牌代码 %d 的点位映射", brand)
	}
	return mapping, nil
}

// GetActivePowerPoint 获取有功功率控制点位
func (m *inverterBrandMapper) GetActivePowerPoint(brand int) (string, error) {
	mapping, err := m.GetBrandMapping(brand)
	if err != nil {
		return "", err
	}
	return mapping.ActivePowerLimitYT, nil
}

// GetReactivePowerPoint 获取无功功率控制点位
func (m *inverterBrandMapper) GetReactivePowerPoint(brand int) (string, error) {
	mapping, err := m.GetBrandMapping(brand)
	if err != nil {
		return "", err
	}
	return mapping.ReactivePowerYT, nil
}

// GetPowerSwitchPoint 获取开关机控制点位
func (m *inverterBrandMapper) GetPowerSwitchPoint(brand int) (string, error) {
	mapping, err := m.GetBrandMapping(brand)
	if err != nil {
		return "", err
	}
	return mapping.PowerSwitchYK, nil
}

// GetBrandName 获取品牌名称
func (m *inverterBrandMapper) GetBrandName(brand int) string {
	mapping, err := m.GetBrandMapping(brand)
	if err != nil {
		return "未知品牌"
	}
	return mapping.BrandName
}

// GetAllBrands 获取所有品牌映射
func (m *inverterBrandMapper) GetAllBrands() []InverterBrandPointMapping {
	brands := make([]InverterBrandPointMapping, 0, len(m.mappings))
	for _, mapping := range m.mappings {
		brands = append(brands, mapping)
	}
	return brands
}
