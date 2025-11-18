package agcv_main

import (
    "context"
    "fmt"
    "time"

    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/agvc_main/request"
    "go.uber.org/zap"
)

type history struct{}

var History = new(history)

// QueryHistoryData 查询历史数据
func (s *history) QueryHistoryData(req request.HistoryDataRequest) ([]map[string]interface{}, error) {
    if global.GVA_INFLUXDB == nil {
        return nil, fmt.Errorf("InfluxDB客户端未初始化")
    }

    // 构建查询语句
    query := s.buildFluxQuery(req)

    global.GVA_LOG.Debug("执行InfluxDB查询", zap.String("query", query))

    // 执行查询
    queryAPI := global.GVA_INFLUXDB.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)
    result, err := queryAPI.Query(context.Background(), query)
    if err != nil {
        global.GVA_LOG.Error("InfluxDB查询失败", zap.Error(err))
        return nil, fmt.Errorf("查询失败: %v", err)
    }
    defer result.Close()

    // 解析结果
    var data []map[string]interface{}
    for result.Next() {
        record := result.Record()
        item := map[string]interface{}{
            "time":     record.Time().Unix(),
            "psid":     record.ValueByKey("psid"),
            "eqid":     record.ValueByKey("eqid"),
            "eqType":   record.ValueByKey("eqType"),
            "dataType": record.ValueByKey("dataType"),
            "point":    record.ValueByKey("point"),
            "value":    record.Value(),
        }
        data = append(data, item)
    }

    if result.Err() != nil {
        global.GVA_LOG.Error("读取查询结果失败", zap.Error(result.Err()))
        return nil, fmt.Errorf("读取结果失败: %v", result.Err())
    }

    // 如果没有查询到数据，返回空数组而不是错误
    if len(data) == 0 {
        global.GVA_LOG.Info("查询的时间范围内没有历史数据")
        return []map[string]interface{}{}, nil
    }

    return data, nil
}

// buildFluxQuery 构建Flux查询语句
func (s *history) buildFluxQuery(req request.HistoryDataRequest) string {
    bucket := global.GVA_CONFIG.InfluxDB.Bucket
    startTime := time.Unix(req.StartTime, 0).Format(time.RFC3339)
    endTime := time.Unix(req.EndTime, 0).Format(time.RFC3339)

    // 基础查询
    query := fmt.Sprintf(`from(bucket: "%s")
  |> range(start: %s, stop: %s)
  |> filter(fn: (r) => r["_measurement"] == "%s")
  |> filter(fn: (r) => r["psid"] == "%s")
  |> filter(fn: (r) => r["eqid"] == "%s")
  |> filter(fn: (r) => r["eqType"] == "%s")`,
        bucket, startTime, endTime, global.GVA_CONFIG.InfluxDB.GetMeasurement(), req.PSID, req.EQID, req.EQType)

    // 添加数据类型过滤
    if req.DataType != "" {
        query += fmt.Sprintf(`
  |> filter(fn: (r) => r["dataType"] == "%s")`, req.DataType)
    }

    // 添加点标识过滤
    if len(req.Points) > 0 {
        pointFilter := ""
        for i, point := range req.Points {
            if i == 0 {
                pointFilter = fmt.Sprintf(`r["point"] == "%s"`, point)
            } else {
                pointFilter += fmt.Sprintf(` or r["point"] == "%s"`, point)
            }
        }
        query += fmt.Sprintf(`
  |> filter(fn: (r) => %s)`, pointFilter)
    }

    // 添加聚合间隔
    if req.Interval != "" {
        query += fmt.Sprintf(`
  |> aggregateWindow(every: %s, fn: mean, createEmpty: false)`, req.Interval)
    }

    // 排序
    query += `
  |> sort(columns: ["_time"])`

    return query
}

// QueryLatestData 查询最新数据（用于验证调节效果）
func (s *history) QueryLatestData(psid, eqid, eqType, dataType, point string, duration time.Duration) (float64, error) {
    if global.GVA_INFLUXDB == nil {
        return 0, fmt.Errorf("InfluxDB客户端未初始化")
    }

    bucket := global.GVA_CONFIG.InfluxDB.Bucket
    startTime := time.Now().Add(-duration).Format(time.RFC3339)

    query := fmt.Sprintf(`from(bucket: "%s")
  |> range(start: %s)
  |> filter(fn: (r) => r["_measurement"] == "%s")
  |> filter(fn: (r) => r["psid"] == "%s")
  |> filter(fn: (r) => r["eqid"] == "%s")
  |> filter(fn: (r) => r["eqType"] == "%s")
  |> filter(fn: (r) => r["dataType"] == "%s")
  |> filter(fn: (r) => r["point"] == "%s")
  |> last()`,
        bucket, startTime, global.GVA_CONFIG.InfluxDB.GetMeasurement(), psid, eqid, eqType, dataType, point)

    queryAPI := global.GVA_INFLUXDB.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)
    result, err := queryAPI.Query(context.Background(), query)
    if err != nil {
        return 0, fmt.Errorf("查询失败: %v", err)
    }
    defer result.Close()

    if result.Next() {
        record := result.Record()
        if value, ok := record.Value().(float64); ok {
            return value, nil
        }
        return 0, fmt.Errorf("数据类型错误")
    }

    if result.Err() != nil {
        return 0, fmt.Errorf("读取结果失败: %v", result.Err())
    }

    return 0, fmt.Errorf("没有找到数据")
}

// QueryAggregateData 查询聚合数据（平均值、最大值、最小值等）
func (s *history) QueryAggregateData(psid, eqid, eqType, dataType, point string, startTime, endTime int64, aggregateFunc string) (float64, error) {
    if global.GVA_INFLUXDB == nil {
        return 0, fmt.Errorf("InfluxDB客户端未初始化")
    }

    // 验证聚合函数
    validFuncs := map[string]bool{
        "mean": true, "max": true, "min": true, "sum": true, "count": true,
    }
    if !validFuncs[aggregateFunc] {
        return 0, fmt.Errorf("不支持的聚合函数: %s", aggregateFunc)
    }

    bucket := global.GVA_CONFIG.InfluxDB.Bucket
    start := time.Unix(startTime, 0).Format(time.RFC3339)
    end := time.Unix(endTime, 0).Format(time.RFC3339)

    query := fmt.Sprintf(`from(bucket: "%s")
  |> range(start: %s, stop: %s)
  |> filter(fn: (r) => r["_measurement"] == "%s")
  |> filter(fn: (r) => r["psid"] == "%s")
  |> filter(fn: (r) => r["eqid"] == "%s")
  |> filter(fn: (r) => r["eqType"] == "%s")
  |> filter(fn: (r) => r["dataType"] == "%s")
  |> filter(fn: (r) => r["point"] == "%s")
  |> %s()`,
        bucket, start, end, global.GVA_CONFIG.InfluxDB.GetMeasurement(), psid, eqid, eqType, dataType, point, aggregateFunc)

    queryAPI := global.GVA_INFLUXDB.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)
    result, err := queryAPI.Query(context.Background(), query)
    if err != nil {
        return 0, fmt.Errorf("查询失败: %v", err)
    }
    defer result.Close()

    if result.Next() {
        record := result.Record()
        if value, ok := record.Value().(float64); ok {
            return value, nil
        }
        return 0, fmt.Errorf("数据类型错误")
    }

    if result.Err() != nil {
        return 0, fmt.Errorf("读取结果失败: %v", result.Err())
    }

    return 0, fmt.Errorf("没有找到数据")
}
