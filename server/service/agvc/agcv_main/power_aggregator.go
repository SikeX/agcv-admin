package agcv_main

import (
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/service/agvc/cons"
	"go.uber.org/zap"
)

// PowerAggregator 功率聚合服务
type powerAggregator struct{}

var PowerAggregator = new(powerAggregator)

// GetBwgPower 获取并网柜功率，如果没有则从逆变器聚合
// powerType: "P"=有功功率, "Q"=无功功率, "S"=视在功率
func (pa *powerAggregator) GetBwgPower(bwgNo int, powerType string) (float64, error) {
	// 获取点标识
	var pointID string
	var err error

	switch powerType {
	case "P": // 有功功率
		pointID, err = PointMapper.GetPointID(cons.TYPE_BWG, "有功功率P(kW)")
	case "Q": // 无功功率
		pointID, err = PointMapper.GetPointID(cons.TYPE_BWG, "无功功率Q(kVar)")
	case "S": // 视在功率
		pointID, err = PointMapper.GetPointID(cons.TYPE_BWG, "S(视在功率kVA)")
	default:
		return 0, fmt.Errorf("不支持的功率类型: %s", powerType)
	}

	if err != nil {
		global.GVA_LOG.Warn("获取并网柜点位映射失败，使用默认值",
			zap.String("powerType", powerType),
			zap.Error(err))
		// 使用默认点标识
		switch powerType {
		case "P":
			pointID = "7"
		case "Q":
			pointID = "8"
		case "S":
			pointID = "11"
		}
	}

	// 尝试从DataStorage获取并网柜数据
	power, err := DataStorage.GetDataAsFloat64(bwgNo, cons.TYPE_BWG, cons.YC, pointID)
	if err == nil {
		global.GVA_LOG.Debug("从DataStorage获取并网柜功率",
			zap.Int("bwgNo", bwgNo),
			zap.String("powerType", powerType),
			zap.Float64("power", power))
		return power, nil
	}

	// 如果没有并网柜数据，从逆变器聚合
	global.GVA_LOG.Debug("并网柜数据不存在，从逆变器聚合",
		zap.Int("bwgNo", bwgNo),
		zap.String("powerType", powerType))

	return pa.aggregateFromInverters(bwgNo, powerType)
}

// aggregateFromInverters 从逆变器聚合功率
func (pa *powerAggregator) aggregateFromInverters(bwgNo int, powerType string) (float64, error) {
	// 获取该并网柜下的所有在线逆变器
	inverters, err := Device.GetOnlineInvertersByBwdNo(bwgNo)
	if err != nil || len(inverters) == 0 {
		return 0, fmt.Errorf("没有可用的逆变器: %v", err)
	}

	// 获取逆变器的点标识
	var invPointID string
	switch powerType {
	case "P": // 有功功率
		invPointID, err = PointMapper.GetPointID(cons.TYPE_NBQ, "交流功率(kW)")
		if err != nil {
			invPointID = "10" // 默认值
		}
	case "Q": // 无功功率
		invPointID, err = PointMapper.GetPointID(cons.TYPE_NBQ, "无功功率(kVar)")
		if err != nil {
			invPointID = "27" // 默认值
		}
	case "S": // 视在功率
		invPointID, err = PointMapper.GetPointID(cons.TYPE_NBQ, "视在功率(kVa)")
		if err != nil {
			invPointID = "28" // 默认值
		}
	}

	// 累加所有逆变器的功率
	var totalPower float64
	successCount := 0

	for _, inv := range inverters {
		invPower, err := DataStorage.GetDataAsFloat64(*inv.InverterNo, cons.TYPE_NBQ, cons.YC, invPointID)
		if err != nil {
			global.GVA_LOG.Warn("获取逆变器功率失败",
				zap.Int("inverterNo", *inv.InverterNo),
				zap.String("powerType", powerType),
				zap.Error(err))
			continue
		}

		totalPower += invPower
		successCount++
	}

	if successCount == 0 {
		return 0, fmt.Errorf("所有逆变器都没有%s数据", powerType)
	}

	global.GVA_LOG.Debug("从逆变器聚合功率完成",
		zap.Int("bwgNo", bwgNo),
		zap.String("powerType", powerType),
		zap.Int("逆变器总数", len(inverters)),
		zap.Int("成功数", successCount),
		zap.Float64("totalPower", totalPower))

	return totalPower, nil
}

// GetBwgActivePower 获取并网柜有功功率
func (pa *powerAggregator) GetBwgActivePower(bwgNo int) (float64, error) {
	return pa.GetBwgPower(bwgNo, "P")
}

// GetBwgReactivePower 获取并网柜无功功率
func (pa *powerAggregator) GetBwgReactivePower(bwgNo int) (float64, error) {
	return pa.GetBwgPower(bwgNo, "Q")
}

// GetBwgApparentPower 获取并网柜视在功率
func (pa *powerAggregator) GetBwgApparentPower(bwgNo int) (float64, error) {
	return pa.GetBwgPower(bwgNo, "S")
}
