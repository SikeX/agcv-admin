package agvc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
	"github.com/flipped-aurora/gin-vue-admin/server/service/agvc/cons"
)

type AgvcChartService struct{}

// PowerChartData 电站出力图表数据
type PowerChartData struct {
	Time               string   `json:"time"`
	CurrentActivePower *float64 `json:"currentActivePower"` // BWD当前有功
	TargetActivePower  *float64 `json:"targetActivePower"`  // AGC目标有功
}

// VoltageReactiveChartData 电压和无功图表数据
type VoltageReactiveChartData struct {
	Time                 string   `json:"time"`
	CurrentReactivePower *float64 `json:"currentReactivePower"` // BWD当前无功
	CurrentVoltage       *float64 `json:"currentVoltage"`       // BWD当前电压
	CurrentVoltageA      *float64 `json:"currentVoltageA"`      //a相电压
	CurrentVoltageB      *float64 `json:"currentVoltageB"`      //b相电压
	CurrentVoltageC      *float64 `json:"currentVoltageC"`      //c相电压
	TargetVoltage        *float64 `json:"targetVoltage"`        // AVC目标电压
	TargetReactivePower  *float64 `json:"targetReactivePower"`  // AVC目标无功
}

// GetPowerChartData 获取电站出力图表数据
// 查询BWD的CurrentActivePower和AGC的TargetActivePower
func (chartService *AgvcChartService) GetPowerChartData(ctx context.Context, psid, eqid int, startTime, endTime string) ([]PowerChartData, error) {
	// 检查InfluxDB客户端是否已初始化
	if global.GVA_INFLUXDB == nil {
		return nil, fmt.Errorf("InfluxDB客户端未初始化")
	}

	// BWD和AGC数据使用agvc储存桶，Measurement统一使用"agvc"
	bucket := global.GVA_CONFIG.InfluxDB.GetAgvcBucket()
	measurement := global.GVA_CONFIG.InfluxDB.GetMeasurement()
	queryAPI := global.GVA_INFLUXDB.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)

	// 构建Flux查询语句 - 查询BWD的CurrentActivePower
	fluxBwd := fmt.Sprintf(`
        from(bucket: "%s")
            |> range(start: %s, stop: %s)
            |> filter(fn: (r) => r["_measurement"] == "%s")
            |> filter(fn: (r) => r["psid"] == "%d")
            |> filter(fn: (r) => r["eqid"] == "%d")
            |> filter(fn: (r) => r["eqType"] == "%s")
            |> filter(fn: (r) => r["dataType"] == "%d")
            |> filter(fn: (r) => r["_field"] == "point_7")
            |> aggregateWindow(every: 5s, fn: last, createEmpty: false)
            |> yield(name: "bwd")`,
		bucket,
		startTime,
		endTime,
		measurement,
		psid,
		eqid,
		agvc.EqTypeBWD,
		cons.YC,
	)

	// 构建Flux查询语句 - 查询AGC的TargetActivePower
	fluxAgc := fmt.Sprintf(`
        from(bucket: "%s")
            |> range(start: %s, stop: %s)
            |> filter(fn: (r) => r["_measurement"] == "%s")
            |> filter(fn: (r) => r["psid"] == "%d")
            |> filter(fn: (r) => r["eqid"] == "%d")
            |> filter(fn: (r) => r["eqType"] == "%s")
            |> filter(fn: (r) => r["dataType"] == "%d")
            |> filter(fn: (r) => r["_field"] == "point_403")
            |> aggregateWindow(every: 5s, fn: last, createEmpty: false)
            |> yield(name: "agc")`,
		bucket,
		startTime,
		endTime,
		measurement,
		psid,
		eqid,
		agvc.EqTypeAGC,
		cons.YC,
	)

	// 合并两个查询
	combinedFlux := fluxBwd + "\n\n" + fluxAgc

	//global.GVA_LOG.Info(fmt.Sprintf("执行Flux查询: %s", combinedFlux))

	// 执行查询
	result, err := queryAPI.Query(ctx, combinedFlux)
	if err != nil {
		global.GVA_LOG.Error(fmt.Sprintf("查询InfluxDB数据失败: %v", err))
		return nil, err
	}

	// 使用map按时间组织数据
	timeDataMap := make(map[string]*PowerChartData)

	// 解析结果
	for result.Next() {
		record := result.Record()
		timeKey := record.Time().Format(time.RFC3339)

		// 获取或创建对应时间的数据对象
		data, exists := timeDataMap[timeKey]
		if !exists {
			data = &PowerChartData{
				Time: timeKey,
			}
			timeDataMap[timeKey] = data
		}

		// 获取值
		value, ok := record.Value().(float64)
		if !ok {
			continue
		}

		// 根据yield的name区分是BWD还是AGC的数据
		resultName := record.Result()
		field := record.Field()

		if resultName == "bwd" && field == "point_7" {
			// BWD的CurrentActivePower
			data.CurrentActivePower = &value
		} else if resultName == "agc" && field == "point_403" {
			// AGC的TargetActivePower
			data.TargetActivePower = &value
		}
	}

	// 检查是否有错误
	if result.Err() != nil {
		global.GVA_LOG.Error(fmt.Sprintf("解析InfluxDB结果失败: %v", result.Err()))
		return nil, result.Err()
	}

	// 将map转换为slice并按时间排序
	chartData := make([]PowerChartData, 0, len(timeDataMap))
	for _, data := range timeDataMap {
		chartData = append(chartData, *data)
	}

	// 如果没有查询到数据，返回空数组
	if len(chartData) == 0 {
		global.GVA_LOG.Info(fmt.Sprintf("设备%d在时间范围%s到%s内没有数据", eqid, startTime, endTime))
		return []PowerChartData{}, nil
	}

	return chartData, nil
}

// GetVoltageReactiveChartData 获取电压和无功图表数据
// 查询BWD的CurrentVoltage、CurrentReactivePower和AVC的TargetVoltage、TargetReactivePower
func (chartService *AgvcChartService) GetVoltageReactiveChartData(ctx context.Context, psid, eqid int, startTime, endTime string) ([]VoltageReactiveChartData, error) {
	// 检查InfluxDB客户端是否已初始化
	if global.GVA_INFLUXDB == nil {
		return nil, fmt.Errorf("InfluxDB客户端未初始化")
	}

	// BWD和AVC数据使用agvc储存桶，Measurement统一使用"agvc"
	bucket := global.GVA_CONFIG.InfluxDB.GetAgvcBucket()
	measurement := global.GVA_CONFIG.InfluxDB.GetMeasurement()
	queryAPI := global.GVA_INFLUXDB.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)

	// 构建Flux查询语句 - 查询BWD的CurrentVoltage和CurrentReactivePower
	fluxBwd := fmt.Sprintf(`
        from(bucket: "%s")
            |> range(start: %s, stop: %s)
            |> filter(fn: (r) => r["_measurement"] == "%s")
            |> filter(fn: (r) => r["psid"] == "%d")
            |> filter(fn: (r) => r["eqid"] == "%d")
            |> filter(fn: (r) => r["eqType"] == "%s")
            |> filter(fn: (r) => r["dataType"] == "%d")
            |> filter(fn: (r) => r["_field"] == "point_1" or r["_field"] == "point_2" or r["_field"] == "point_3" or r["_field"] == "point_8")
            |> aggregateWindow(every: 5s, fn: last, createEmpty: false)
            |> yield(name: "bwd")`,
		bucket,
		startTime,
		endTime,
		measurement,
		psid,
		eqid,
		agvc.EqTypeBWD,
		cons.YC,
	)

	// 构建Flux查询语句 - 查询AVC的TargetVoltage和TargetReactivePower
	fluxAvc := fmt.Sprintf(`
        from(bucket: "%s")
            |> range(start: %s, stop: %s)
            |> filter(fn: (r) => r["_measurement"] == "%s")
            |> filter(fn: (r) => r["psid"] == "%d")
            |> filter(fn: (r) => r["eqid"] == "%d")
            |> filter(fn: (r) => r["eqType"] == "%s")
            |> filter(fn: (r) => r["dataType"] == "%d")
            |> filter(fn: (r) => r["_field"] == "point_403" or r["_field"] == "point_404")
            |> aggregateWindow(every: 5s, fn: last, createEmpty: false)
            |> yield(name: "avc")`,
		bucket,
		startTime,
		endTime,
		measurement,
		psid,
		eqid,
		agvc.EqTypeAVC,
		cons.YC,
	)

	// 合并两个查询
	combinedFlux := fluxBwd + "\n\n" + fluxAvc

	//global.GVA_LOG.Info(fmt.Sprintf("执行Flux查询: %s", combinedFlux))

	// 执行查询
	result, err := queryAPI.Query(ctx, combinedFlux)
	if err != nil {
		global.GVA_LOG.Error(fmt.Sprintf("查询InfluxDB数据失败: %v", err))
		return nil, err
	}

	// 使用map按时间组织数据
	timeDataMap := make(map[string]*VoltageReactiveChartData)

	// 解析结果
	for result.Next() {
		record := result.Record()
		timeKey := record.Time().Format(time.RFC3339)

		// 获取或创建对应时间的数据对象
		data, exists := timeDataMap[timeKey]
		if !exists {
			data = &VoltageReactiveChartData{
				Time: timeKey,
			}
			timeDataMap[timeKey] = data
		}

		// 获取值
		value, ok := record.Value().(float64)
		if !ok {
			continue
		}

		// 根据yield的name和field区分不同的数据
		resultName := record.Result()
		field := record.Field()

		// 提取point值
		pointArr := strings.Split(field, "_")
		if len(pointArr) < 2 {
			continue
		}
		point := pointArr[1]

		if resultName == "bwd" {
			// BWD的数据
			switch point {
			case "1":
				// CurrentVoltage
				data.CurrentVoltageA = &value
			case "2":
				// CurrentVoltageB
				data.CurrentVoltageB = &value
			case "3":
				// CurrentVoltageC
				data.CurrentVoltageC = &value
			case "8":
				// CurrentReactivePower
				data.CurrentReactivePower = &value
			}
		} else if resultName == "avc" {
			// AVC的数据
			switch point {
			case "403":
				// TargetVoltage
				data.TargetVoltage = &value
			case "404":
				// TargetReactivePower
				data.TargetReactivePower = &value
			}
		}
		if data.CurrentVoltageA != nil && data.CurrentVoltageB != nil && data.CurrentVoltageC != nil {
			currentVoltage := (*data.CurrentVoltageA + *data.CurrentVoltageB + *data.CurrentVoltageC) / 3
			data.CurrentVoltage = &currentVoltage
		}
	}

	// 检查是否有错误
	if result.Err() != nil {
		global.GVA_LOG.Error(fmt.Sprintf("解析InfluxDB结果失败: %v", result.Err()))
		return nil, result.Err()
	}

	// 将map转换为slice
	chartData := make([]VoltageReactiveChartData, 0, len(timeDataMap))
	for _, data := range timeDataMap {
		chartData = append(chartData, *data)
	}

	// 如果没有查询到数据，返回空数组
	if len(chartData) == 0 {
		global.GVA_LOG.Info(fmt.Sprintf("设备%d在时间范围%s到%s内没有数据", eqid, startTime, endTime))
		return []VoltageReactiveChartData{}, nil
	}

	return chartData, nil
}
