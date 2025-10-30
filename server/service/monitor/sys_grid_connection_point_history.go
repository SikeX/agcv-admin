
package monitor

import (
	"context"
	"fmt"
	"time"
	
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/monitor"
	agvcReq "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/request"
	serviceAgvc "github.com/flipped-aurora/gin-vue-admin/server/service/agvc"
	"go.uber.org/zap"
)

type SysGridConnectionPointHistoryService struct{}

var agvcBwdSettingService serviceAgvc.AgvcBwdSettingService

// GetSysGridConnectionPointHistoryList 获取并网点配置列表（只返回ID和名称）
// Author [yourname](https://github.com/yourname)
func (sysGridConnectionPointHistoryService *SysGridConnectionPointHistoryService) GetSysGridConnectionPointHistoryList(ctx context.Context) ([]monitor.SysGridConnectionPointHistory, error) {
	// 获取并网点配置列表
	list, _, err := agvcBwdSettingService.GetAgvcBwdSettingInfoList(ctx, agvcReq.AgvcBwdSettingSearch{})
	if err != nil {
		return nil, err
	}
	
	// 只返回ID和名称
	var result []monitor.SysGridConnectionPointHistory
	for _, item := range list {
		result = append(result, monitor.SysGridConnectionPointHistory{
			GVA_MODEL: item.GVA_MODEL,
			Name:      item.Name,
		})
	}
	
	return result, nil
}

// GetGridPointHistoryData 获取并网点历史数据
// Author [yourname](https://github.com/yourname)
func (sysGridConnectionPointHistoryService *SysGridConnectionPointHistoryService) GetGridPointHistoryData(ctx context.Context, query monitor.GridPointHistoryQuery) ([]monitor.GridPointHistoryResponse, error) {
	// 验证并网点ID
	if query.GridPointID <= 0 {
		return nil, fmt.Errorf("无效的并网点ID: %d", query.GridPointID)
	}

	// 获取并网点对应的点标识列表
	pointIdentifiers, err := sysGridConnectionPointHistoryService.GetPointIdentifiersByGridPoint(query.GridPointID)
	if err != nil {
		return nil, fmt.Errorf("获取点标识列表失败: %v", err)
	}

	// 查询所有点标识的数据
	var results []monitor.GridPointHistoryResponse
	for _, identifier := range pointIdentifiers {
		// 查询数据
		data, err := sysGridConnectionPointHistoryService.queryInfluxDBData(identifier.PointID, query.StartTime, query.EndTime)
		if err != nil {
			// 如果某个点标识查询失败，记录日志并继续查询其他点标识
			global.GVA_LOG.Error("查询点标识数据失败", zap.String("pointID", fmt.Sprintf("%d", identifier.PointID)), zap.Error(err))
			continue
		}
		
		results = append(results, monitor.GridPointHistoryResponse{
			PointID:   identifier.PointID,
			PointName: identifier.PointName,
			Data:      data,
		})
	}

	return results, nil
}

// queryInfluxDBData 从InfluxDB查询历史数据
func (sysGridConnectionPointHistoryService *SysGridConnectionPointHistoryService) queryInfluxDBData(pointID int, startTime, endTime string) ([]monitor.GridPointHistoryData, error) {
	// 解析时间参数
	var start, end time.Time
	var err error
	
	if startTime != "" {
		start, err = time.Parse(time.RFC3339, startTime)
		if err != nil {
			return nil, fmt.Errorf("开始时间格式错误: %v", err)
		}
	} else {
		// 默认查询最近24小时的数据
		start = time.Now().Add(-24 * time.Hour)
	}
	
	if endTime != "" {
		end, err = time.Parse(time.RFC3339, endTime)
		if err != nil {
			return nil, fmt.Errorf("结束时间格式错误: %v", err)
		}
	} else {
		end = time.Now()
	}
	
	// 如果有配置InfluxDB连接，则查询真实数据
	if global.GVA_INFLUXDB != nil {
		// 获取查询API
		queryAPI := global.GVA_INFLUXDB.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)
		
		// 构建Flux查询语句
		flux := fmt.Sprintf(`
	from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r["_measurement"] == "mqtt")
		|> filter(fn: (r) => r["point_id"] == "%d")
		|> yield(name: "mean")
	`, global.GVA_CONFIG.InfluxDB.Bucket, start.Format(time.RFC3339), end.Format(time.RFC3339), pointID)
		
		// 执行查询
		result, err := queryAPI.Query(context.Background(), flux)
		if err != nil {
			return nil, fmt.Errorf("InfluxDB查询失败: %v", zap.Error(err))
		}
		
		// 解析查询结果
		var historyData []monitor.GridPointHistoryData
		for result.Next() {
			record := result.Record()
			
			// 获取时间
			timestamp := record.Time()
			
			// 获取值并转换为float64
			var value float64
			switch v := record.Value().(type) {
			case int64:
				value = float64(v)
			case float64:
				value = v
			}
			
			historyData = append(historyData, monitor.GridPointHistoryData{
				Time:  timestamp.Format(time.RFC3339),
				Value: value,
			})
		}
		
		// 检查是否有错误
		if result.Err() != nil {
			return nil, fmt.Errorf("解析InfluxDB结果失败 %v", zap.Error(result.Err()))
		}
		
		return historyData, nil
	}
	
	// 如果没有配置InfluxDB或查询失败，返回模拟数据用于测试
	return nil,nil
}

// GetPointIdentifiersByGridPoint 根据并网点ID获取对应的点标识列表
func (sysGridConnectionPointHistoryService *SysGridConnectionPointHistoryService) GetPointIdentifiersByGridPoint(gridPointID int) ([]monitor.PointIdentifierInfo, error) {
	pointIDs := monitor.GetPointIdentifiersByGridPoint(gridPointID)
	if len(pointIDs) == 0 {
		return nil, fmt.Errorf("未找到并网点ID %d 对应的点标识", gridPointID)
	}

	var result []monitor.PointIdentifierInfo
	for _, pointID := range pointIDs {
		pointName := monitor.GetPointName(pointID)
		if pointName == "未知点位" {
			continue
		}
		result = append(result, monitor.PointIdentifierInfo{
			PointID:   pointID,
			PointName: pointName,
		})
	}

	return result, nil
}
