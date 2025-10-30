
package monitor

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	agvcReq "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/monitor"
	serviceAgvc "github.com/flipped-aurora/gin-vue-admin/server/service/agvc"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type SysQixiangyiHistoryService struct{}

var agvcQxySettingService serviceAgvc.AgvcQxySettingService

// GetSysQixiangyiHistoryList 获取气象仪配置列表（只返回ID和名称）
// 参考并网点监控的实现，从AGVC气象仪设置中获取配置列表
func (sysQixiangyiHistoryService *SysQixiangyiHistoryService) GetSysQixiangyiHistoryList(ctx context.Context) ([]monitor.SysQixiangyiHistory, error) {
	// 获取AGVC气象仪配置列表
	list, _, err := agvcQxySettingService.GetAgvcQxySettingInfoList(ctx, agvcReq.AgvcQxySettingSearch{})
	if err != nil {
		return nil, err
	}

	var result []monitor.SysQixiangyiHistory
	for _, setting := range list {
		// 处理指针类型的DeviceName字段
		deviceName := ""
		if setting.DeviceName != nil {
			deviceName = *setting.DeviceName
		}
		
		result = append(result, monitor.SysQixiangyiHistory{
			Name:        &deviceName,
			QixiangyiID: int(setting.ID), // 使用设置ID作为气象仪ID
		})
	}

	return result, nil
}

// GetQixiangyiHistoryData 获取气象仪历史数据
// 参考并网点监控的实现，通过InfluxDB获取历史数据
func (sysQixiangyiHistoryService *SysQixiangyiHistoryService) GetQixiangyiHistoryData(ctx context.Context, query monitor.QixiangyiHistoryQuery) ([]monitor.QixiangyiHistoryResponse, error) {
	// 验证气象仪ID
	if query.QixiangyiID <= 0 {
		return nil, fmt.Errorf("无效的气象仪ID: %d", query.QixiangyiID)
	}

	// 获取气象仪对应的点标识列表
	pointIdentifiers := sysQixiangyiHistoryService.GetPointIdentifiersByQixiangyiId(query.QixiangyiID)

	var result []monitor.QixiangyiHistoryResponse
	
	// 为每个点标识查询历史数据
	for _, pointID := range pointIdentifiers {
		data, err := sysQixiangyiHistoryService.queryInfluxDBData(pointID, query.StartTime, query.EndTime)
		if err != nil {
			global.GVA_LOG.Error(fmt.Sprintf("查询点标识 %d 的历史数据失败: %v", pointID, err))
			continue
		}

		result = append(result, monitor.QixiangyiHistoryResponse{
			PointID:   pointID,
			PointName: monitor.GetQixiangyiPointName(pointID),
			Data:      data,
		})
	}

	return result, nil
}

// queryInfluxDBData 查询InfluxDB中的历史数据
// 参考并网点监控的实现
func (sysQixiangyiHistoryService *SysQixiangyiHistoryService) queryInfluxDBData(pointID int, startTime, endTime string) ([]monitor.QixiangyiHistoryData, error) {
	// 检查InfluxDB配置
	if global.GVA_CONFIG.InfluxDB.Host == "" {
		return nil, fmt.Errorf("InfluxDB未配置")
	}

	// 创建InfluxDB客户端
	client := influxdb2.NewClient(global.GVA_CONFIG.InfluxDB.Dsn(), global.GVA_CONFIG.InfluxDB.Token)
	defer client.Close()

	// 获取查询API
	queryAPI := client.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)

	// 构建查询语句
	query := fmt.Sprintf(`
		from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r["_measurement"] == "qixiangyi_data")
		|> filter(fn: (r) => r["point_id"] == "%d")
		|> filter(fn: (r) => r["_field"] == "value")
		|> sort(columns: ["_time"])
	`, global.GVA_CONFIG.InfluxDB.Bucket, startTime, endTime, pointID)

	// 执行查询
	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("查询InfluxDB失败: %v", err)
	}
	defer result.Close()

	var data []monitor.QixiangyiHistoryData

	// 处理查询结果
	for result.Next() {
		record := result.Record()
		
		// 获取时间和值
		timestamp := record.Time().Format(time.RFC3339)
		value, ok := record.Value().(float64)
		if !ok {
			// 尝试转换其他数值类型
			if intValue, ok := record.Value().(int64); ok {
				value = float64(intValue)
			} else if strValue, ok := record.Value().(string); ok {
				if parsedValue, err := strconv.ParseFloat(strValue, 64); err == nil {
					value = parsedValue
				} else {
					continue // 跳过无法解析的值
				}
			} else {
				continue // 跳过非数值类型
			}
		}

		data = append(data, monitor.QixiangyiHistoryData{
			Time:  timestamp,
			Value: value,
		})
	}

	// 检查查询错误
	if result.Err() != nil {
		return nil, fmt.Errorf("处理查询结果失败: %v", result.Err())
	}

	return data, nil
}

// GetPointIdentifiersByQixiangyiId 根据气象仪ID获取对应的点标识列表
func (sysQixiangyiHistoryService *SysQixiangyiHistoryService) GetPointIdentifiersByQixiangyiId(qixiangyiId int) []int {
	return monitor.GetPointIdentifiersByQixiangyiId(qixiangyiId)
}

// GenerateTestData 生成测试数据
// 参考并网点监控的实现
func (sysQixiangyiHistoryService *SysQixiangyiHistoryService) GenerateTestData(qixiangyiId int) error {
	// 检查InfluxDB配置
	if global.GVA_CONFIG.InfluxDB.Host == "" {
		return fmt.Errorf("InfluxDB未配置")
	}


	// 创建InfluxDB客户端
	client := influxdb2.NewClient(global.GVA_CONFIG.InfluxDB.Dsn(), global.GVA_CONFIG.InfluxDB.Token)
	defer client.Close()

	// 获取写入API
	writeAPI := client.WriteAPIBlocking(global.GVA_CONFIG.InfluxDB.Org, global.GVA_CONFIG.InfluxDB.Bucket)

	// 获取气象仪对应的点标识列表
	pointIdentifiers := sysQixiangyiHistoryService.GetPointIdentifiersByQixiangyiId(qixiangyiId)

	// 生成过去24小时的测试数据
	now := time.Now()
	startTime := now.Add(-24 * time.Hour)

	// 为每个点标识生成数据
	for _, pointID := range pointIdentifiers {
		// 每5分钟生成一个数据点
		for t := startTime; t.Before(now); t = t.Add(5 * time.Minute) {
			// 根据点标识生成合理的测试值
			value := sysQixiangyiHistoryService.generateTestValue(pointID)
			
			// 创建数据点
			point := influxdb2.NewPointWithMeasurement("qixiangyi_data").
				AddTag("point_id", fmt.Sprintf("%d", pointID)).
				AddTag("point_name", monitor.GetQixiangyiPointName(pointID)).
				AddField("value", value).
				SetTime(t)

			// 写入数据
			if err := writeAPI.WritePoint(context.Background(), point); err != nil {
				global.GVA_LOG.Error(fmt.Sprintf("写入测试数据失败: %v", err))
				continue
			}
		}
	}

	return nil
}

// generateTestValue 根据点标识生成合理的测试值
func (sysQixiangyiHistoryService *SysQixiangyiHistoryService) generateTestValue(pointID int) float64 {
	// 根据不同的点标识生成不同范围的随机值
	switch pointID {
	case 2: // 环境温度
		return 20 + float64(time.Now().Unix()%20) // 20-40℃
	case 3: // 露点温度
		return 10 + float64(time.Now().Unix()%15) // 10-25℃
	case 4: // 风速
		return float64(time.Now().Unix()%10) // 0-10 m/s
	case 7: // 风向
		return float64(time.Now().Unix() % 360) // 0-360°
	case 8, 9: // 辐射强度
		return float64(time.Now().Unix()%800) + 200 // 200-1000 W/㎡
	case 11: // 组件温度
		return 25 + float64(time.Now().Unix()%30) // 25-55℃
	case 46: // 环境湿度
		return 40 + float64(time.Now().Unix()%40) // 40-80%
	case 52: // 气压
		return 1000 + float64(time.Now().Unix()%50) // 1000-1050 Pa
	default:
		return float64(time.Now().Unix() % 100) // 默认0-100
	}
}
