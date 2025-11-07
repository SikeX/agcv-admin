package agcv_main

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc/agvc_main"
	"go.uber.org/zap"
)

// HuaweiPointInitializer 华为逆变器点位初始化器
type HuaweiPointInitializer struct{}

var HuaweiPointInit = &HuaweiPointInitializer{}

// InitializeHuaweiPoints 初始化华为逆变器5xx系列点位到数据库
func (h *HuaweiPointInitializer) InitializeHuaweiPoints() error {
	global.GVA_LOG.Info("开始初始化华为逆变器5xx点位映射...")

	pointMappings := []agvc_main.PointMapping{
		{
			EQType:      "02",
			DataType:    "05",
			Point:       "501",
			PointName:   "无功功率变化梯度",
			Unit:        "%/s",
			Description: "无功功率变化梯度(%/s) - 寄存器地址:42015",
			Category:    "HUAWEI_AVC",
		},
		{
			EQType:      "02",
			DataType:    "05",
			Point:       "502",
			PointName:   "有功功率变化梯度",
			Unit:        "%/s",
			Description: "有功功率变化梯度(%/s) - 寄存器地址:42017",
			Category:    "HUAWEI_AGC",
		},
		{
			EQType:      "02",
			DataType:    "05",
			Point:       "503",
			PointName:   "调度指令维持时间",
			Unit:        "s",
			Description: "调度指令维持时间(s) - 寄存器地址:42019",
			Category:    "HUAWEI_CONTROL",
		},
		{
			EQType:      "02",
			DataType:    "05",
			Point:       "504",
			PointName:   "有功功率固定值降额(kW)",
			Unit:        "kW",
			Description: "有功功率固定值降额(kW) - 寄存器地址:40120",
			Category:    "HUAWEI_AGC",
		},
		{
			EQType:      "02",
			DataType:    "05",
			Point:       "505",
			PointName:   "无功功率补偿(PF)",
			Unit:        "PF",
			Description: "无功功率补偿(功率因数) - 寄存器地址:40122",
			Category:    "HUAWEI_AVC",
		},
		{
			EQType:      "02",
			DataType:    "05",
			Point:       "506",
			PointName:   "无功功率补偿(Q/S)",
			Unit:        "Q/S",
			Description: "无功功率补偿(Q/S) - 寄存器地址:40123",
			Category:    "HUAWEI_AVC",
		},
		{
			EQType:      "02",
			DataType:    "05",
			Point:       "507",
			PointName:   "有功功率百分比降额",
			Unit:        "0.1%",
			Description: "有功功率百分比降额(0.1%) - 寄存器地址:40125",
			Category:    "HUAWEI_AGC",
		},
		{
			EQType:      "02",
			DataType:    "05",
			Point:       "508",
			PointName:   "有功功率固定值降额(W)",
			Unit:        "W",
			Description: "有功功率固定值降额(W) - 寄存器地址:40126",
			Category:    "HUAWEI_AGC",
		},
		{
			EQType:      "02",
			DataType:    "05",
			Point:       "510",
			PointName:   "有功调节模式",
			Unit:        "",
			Description: "有功调节模式(0:百分比 1:绝对值) - 寄存器地址:35300",
			Category:    "HUAWEI_CONTROL",
		},
		{
			EQType:      "02",
			DataType:    "05",
			Point:       "511",
			PointName:   "有功调节值",
			Unit:        "",
			Description: "有功调节值 - 寄存器地址:35301",
			Category:    "HUAWEI_CONTROL",
		},
		{
			EQType:      "02",
			DataType:    "05",
			Point:       "512",
			PointName:   "有功调节指令",
			Unit:        "",
			Description: "有功调节指令(执行) - 寄存器地址:35303",
			Category:    "HUAWEI_CONTROL",
		},
		{
			EQType:      "02",
			DataType:    "05",
			Point:       "513",
			PointName:   "无功调节模式",
			Unit:        "",
			Description: "无功调节模式(1:PF 2:Q/S) - 寄存器地址:35304",
			Category:    "HUAWEI_CONTROL",
		},
		{
			EQType:      "02",
			DataType:    "05",
			Point:       "514",
			PointName:   "无功调节值",
			Unit:        "",
			Description: "无功调节值 - 寄存器地址:35305",
			Category:    "HUAWEI_CONTROL",
		},
		{
			EQType:      "02",
			DataType:    "05",
			Point:       "515",
			PointName:   "无功调节指令",
			Unit:        "",
			Description: "无功调节指令(执行) - 寄存器地址:35307",
			Category:    "HUAWEI_CONTROL",
		},
	}

	// 批量插入或更新点位映射
	successCount := 0
	for _, mapping := range pointMappings {
		// 先查询是否存在
		var existingMapping agvc_main.PointMapping
		err := global.GVA_DB.Where("eq_type = ? AND data_type = ? AND point = ?",
			mapping.EQType, mapping.DataType, mapping.Point).First(&existingMapping).Error

		if err != nil {
			// 不存在，创建新记录
			if err := global.GVA_DB.Create(&mapping).Error; err != nil {
				global.GVA_LOG.Warn("创建华为点位映射失败",
					zap.String("point", mapping.Point),
					zap.String("pointName", mapping.PointName),
					zap.Error(err))
				continue
			}
			global.GVA_LOG.Info("创建华为点位映射成功",
				zap.String("point", mapping.Point),
				zap.String("pointName", mapping.PointName))
		} else {
			// 已存在，更新记录
			existingMapping.PointName = mapping.PointName
			existingMapping.Unit = mapping.Unit
			existingMapping.Description = mapping.Description
			existingMapping.Category = mapping.Category

			if err := global.GVA_DB.Save(&existingMapping).Error; err != nil {
				global.GVA_LOG.Warn("更新华为点位映射失败",
					zap.String("point", mapping.Point),
					zap.String("pointName", mapping.PointName),
					zap.Error(err))
				continue
			}
			global.GVA_LOG.Info("更新华为点位映射成功",
				zap.String("point", mapping.Point),
				zap.String("pointName", mapping.PointName))
		}
		successCount++
	}

	global.GVA_LOG.Info("华为逆变器5xx点位映射初始化完成",
		zap.Int("成功数量", successCount),
		zap.Int("总数量", len(pointMappings)))

	return nil
}

// GetHuaweiPointMapping 获取华为逆变器点位映射
func (h *HuaweiPointInitializer) GetHuaweiPointMapping(pointID string) (*agvc_main.PointMapping, error) {
	var mapping agvc_main.PointMapping
	err := global.GVA_DB.Where("eq_type = ? AND data_type = ? AND point = ?",
		"02", "05", pointID).First(&mapping).Error
	if err != nil {
		return nil, err
	}
	return &mapping, nil
}

// GetAllHuaweiPointMappings 获取所有华为逆变器点位映射
func (h *HuaweiPointInitializer) GetAllHuaweiPointMappings() ([]agvc_main.PointMapping, error) {
	var mappings []agvc_main.PointMapping
	err := global.GVA_DB.Where("point LIKE ? AND eq_type = ?", "5%", "02").Find(&mappings).Error
	if err != nil {
		return nil, err
	}
	return mappings, nil
}
