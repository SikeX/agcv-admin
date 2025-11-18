package agvc

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
	"github.com/flipped-aurora/gin-vue-admin/server/service/agvc/cons"
)

// AgvcHisService 通用历史数据服务
type AgvcHisService struct{}

// GetHistory 通用的历史数据查询方法
// psid: 电站编号, eqid: 设备编号, eqType: 设备类型, pointValues: 点位列表, startTime/endTime: 时间范围
// 支持任意设备类型的历史数据查询，包括：NBQ（逆变器）、BWD（并网点）、AGC、AVC、QXY（气象仪）等
func (agvcHisService *AgvcHisService) GetHistory(ctx context.Context, psid int, eqid string, eqType string, pointValues []string, startTime, endTime string) ([]agvc.AgvcNbqHis, error) {
	// 检查InfluxDB客户端是否已初始化
	if global.GVA_INFLUXDB == nil {
		return nil, fmt.Errorf("InfluxDB客户端未初始化")
	}

	// 获取查询API
	queryAPI := global.GVA_INFLUXDB.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)

	// 如果没有指定点位，返回错误
	if len(pointValues) == 0 {
		return nil, fmt.Errorf("未找到需要查询的点位信息")
	}

	// 构建Flux查询语句
	// 查询所有点位的历史数据
	flux := fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: %s, stop: %s)
			|> filter(fn: (r) => r["_measurement"] == "%s")
			|> filter(fn: (r) => r["psid"] == "%d")
			|> filter(fn: (r) => r["eqid"] == "%s")
			|> filter(fn: (r) => r["eqType"] == "%s")
			|> filter(fn: (r) => r["dataType"] == "%d")`,
		global.GVA_CONFIG.InfluxDB.Bucket,
		startTime,
		endTime,
		global.GVA_CONFIG.InfluxDB.GetMeasurement(),
		psid,
		eqid,
		eqType,
		cons.YC,
	)

	// 如果只查询特定点位，添加点位过滤
	if len(pointValues) > 0 {
		pointFilters := make([]string, 0, len(pointValues))
		for _, point := range pointValues {
			pointFilters = append(pointFilters, fmt.Sprintf(`r["_field"] == "point_%s"`, point))
		}
		flux += fmt.Sprintf(`
			|> filter(fn: (r) => %s)`, strings.Join(pointFilters, " or "))
	}

	// 执行查询
	result, err := queryAPI.Query(ctx, flux)
	if err != nil {
		global.GVA_LOG.Error(fmt.Sprintf("查询InfluxDB数据失败: %v", err))
		return nil, err
	}

	// 按时间分组存储数据
	timeDataMap := make(map[string]*agvc.AgvcNbqHis)

	// 解析结果
	for result.Next() {
		record := result.Record()
		timeKey := record.Time().Format(time.RFC3339)
		
		// 提取point值
		field := record.Field()
		pointArr := strings.Split(field, "_")
		if len(pointArr) < 2 {
			continue
		}
		point := pointArr[1]

		// 获取或创建对应时间的AgvcNbqHis
		nbqHis, exists := timeDataMap[timeKey]
		if !exists {
			ctimeStr := record.Time().Format(time.RFC3339)
			nbqHis = &agvc.AgvcNbqHis{
				Ctime: &ctimeStr,
				Psid:  &psid,
			}
			// 根据设备类型设置设备编号字段
			// NBQ 使用 InverterNo，其他类型可以扩展
			if eqType == agvc.EqTypeNBQ {
				eqidInt := 0
				fmt.Sscanf(eqid, "%d", &eqidInt)
				nbqHis.InverterNo = &eqidInt
			}
			timeDataMap[timeKey] = nbqHis
		}

		// 设置字段值
		value, ok := record.Value().(float64)
		if !ok {
			continue
		}

		// 使用缓存的映射关系设置字段值
		if fieldName, found := agvc.GetFieldNameByPoint(point); found {
			setFieldValue(nbqHis, fieldName, value)
		}
	}

	// 检查是否有错误
	if result.Err() != nil {
		global.GVA_LOG.Error(fmt.Sprintf("解析InfluxDB结果失败: %v", result.Err()))
		return nil, result.Err()
	}

	// 将map转换为slice
	historyData := make([]agvc.AgvcNbqHis, 0, len(timeDataMap))
	for _, data := range timeDataMap {
		historyData = append(historyData, *data)
	}

	// 如果没有查询到数据，返回空数组
	if len(historyData) == 0 {
		global.GVA_LOG.Info(fmt.Sprintf("设备%s在时间范围%s到%s内没有历史数据", eqid, startTime, endTime))
		return []agvc.AgvcNbqHis{}, nil
	}

	return historyData, nil
}

// GetRealData 通用的实时数据查询方法
// psid: 电站编号, eqid: 设备编号, eqType: 设备类型
// 从InfluxDB获取指定设备的最新数据
func (agvcHisService *AgvcHisService) GetRealData(ctx context.Context, psid int, eqid string, eqType string) (*agvc.AgvcNbqHis, error) {
	// 检查InfluxDB客户端是否已初始化
	if global.GVA_INFLUXDB == nil {
		return nil, fmt.Errorf("InfluxDB客户端未初始化")
	}

	// 获取查询API
	queryAPI := global.GVA_INFLUXDB.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)

	// 构建Flux查询语句 - 查询该设备所有点位的最新值
	flux := fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: -1y)
			|> filter(fn: (r) => r["_measurement"] == "%s")
			|> filter(fn: (r) => r["psid"] == "%d")
			|> filter(fn: (r) => r["eqid"] == "%s")
			|> filter(fn: (r) => r["eqType"] == "%s")
			|> filter(fn: (r) => r["dataType"] == "%d")
			|> last()`,
		global.GVA_CONFIG.InfluxDB.Bucket,
		global.GVA_CONFIG.InfluxDB.GetMeasurement(),
		psid,
		eqid,
		eqType,
		cons.YC,
	)

	// 执行查询
	result, err := queryAPI.Query(ctx, flux)
	if err != nil {
		global.GVA_LOG.Error(fmt.Sprintf("查询设备%s的实时数据失败: %v", eqid, err))
		return nil, err
	}

	// 创建数据对象
	realData := &agvc.AgvcNbqHis{
		Psid: &psid,
	}
	
	// 根据设备类型设置设备编号字段
	if eqType == agvc.EqTypeNBQ {
		eqidInt := 0
		fmt.Sscanf(eqid, "%d", &eqidInt)
		realData.InverterNo = &eqidInt
	}

	// 解析结果
	for result.Next() {
		record := result.Record()
		
		// 设置采集时间
		if realData.Ctime == nil {
			ctimeStr := record.Time().Format(time.RFC3339)
			realData.Ctime = &ctimeStr
		}

		// 提取point值
		field := record.Field()
		pointArr := strings.Split(field, "_")
		if len(pointArr) < 2 {
			continue
		}
		point := pointArr[1]

		value, ok := record.Value().(float64)
		if !ok {
			continue
		}

		// 使用缓存的映射关系设置字段值
		if fieldName, found := agvc.GetFieldNameByPoint(point); found {
			setFieldValue(realData, fieldName, value)
		}
	}

	// 检查是否有错误
	if result.Err() != nil {
		global.GVA_LOG.Error(fmt.Sprintf("解析设备%s的实时数据失败: %v", eqid, result.Err()))
		return nil, result.Err()
	}

	// 如果没有采集时间，说明没有数据
	if realData.Ctime == nil {
		return nil, fmt.Errorf("未找到设备%s的实时数据", eqid)
	}

	return realData, nil
}

// setFieldValue 根据字段名设置AgvcNbqHis结构体对应字段的值
func setFieldValue(nbqHis *agvc.AgvcNbqHis, fieldName string, value float64) {
	v := reflect.ValueOf(nbqHis).Elem()
	field := v.FieldByName(fieldName)
	
	if !field.IsValid() || !field.CanSet() {
		return
	}

	// 设置指针类型的float64值
	if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.Float64 {
		field.Set(reflect.ValueOf(&value))
	}
}
