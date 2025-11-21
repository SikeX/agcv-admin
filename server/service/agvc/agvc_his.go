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

// GetHistoryGeneric 通用的历史数据查询方法（返回interface{}）
// psid: 电站编号, eqid: 设备编号, eqType: 设备类型, pointValues: 点位列表, startTime/endTime: 时间范围
// 支持任意设备类型的历史数据查询，包括：NBQ（逆变器）、BWD（并网点）、AGC、AVC、QXY（气象仪）等
// 根据设备类型返回对应的数据结构：NBQ->[]AgvcNbq, BWD->[]AgvcBwd, AGC->[]AgvcAgc, AVC->[]AgvcAvc
func (agvcHisService *AgvcHisService) GetHistoryGeneric(ctx context.Context, psid, eqid int, eqType string, pointValues []string, eqName, startTime, endTime string) (interface{}, error) {
    switch eqType {
    case agvc.EqTypeNBQ:
        return agvcHisService.getHistoryNbq(ctx, psid, eqid, pointValues, eqName, startTime, endTime)
    case agvc.EqTypeBWD:
        return agvcHisService.getHistoryBwd(ctx, psid, eqid, pointValues, startTime, endTime)
    case agvc.EqTypeAGC:
        return agvcHisService.getHistoryAgc(ctx, psid, eqid, pointValues, startTime, endTime)
    case agvc.EqTypeAVC:
        return agvcHisService.getHistoryAvc(ctx, psid, eqid, pointValues, startTime, endTime)
    default:
        return nil, fmt.Errorf("不支持的设备类型: %s", eqType)
    }
}

// getHistoryNbq 查询NBQ历史数据
func (agvcHisService *AgvcHisService) getHistoryNbq(ctx context.Context, psid, eqid int, pointValues []string, eqName, startTime, endTime string) ([]agvc.AgvcNbq, error) {
    // 检查InfluxDB客户端是否已初始化
    if global.GVA_INFLUXDB == nil {
        return nil, fmt.Errorf("InfluxDB客户端未初始化")
    }

    measurement := global.GVA_CONFIG.InfluxDB.GetNBQMeasurement()

    // 获取查询API
    queryAPI := global.GVA_INFLUXDB.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)

    // 如果没有指定点位，返回错误
    // if len(pointValues) == 0 {
    //     return nil, fmt.Errorf("未找到需要查询的点位信息")
    // }

    // 构建Flux查询语句
    // 查询所有点位的历史数据
    flux := fmt.Sprintf(`
        from(bucket: "%s")
            |> range(start: %s, stop: %s)
            |> filter(fn: (r) => r["_measurement"] == "%s")
            |> filter(fn: (r) => r["psid"] == "%d")
            |> filter(fn: (r) => r["eqid"] == "%d")
            |> filter(fn: (r) => r["eqType"] == "%d")
            |> filter(fn: (r) => r["dataType"] == "%d")
            |> aggregateWindow(every: 5m, fn: last, createEmpty: false)
            |> pivot(rowKey: ["_time"], columnKey: ["_field"], valueColumn: "_value")
            |> sort(columns: ["_time"], desc: true)`,
        global.GVA_CONFIG.InfluxDB.Bucket,
        startTime,
        endTime,
        measurement,
        psid,
        eqid,
        NBQ_DEVICE_TYPE,
        cons.YC,
    )

    fmt.Println(flux)

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

    nbqhises := make([]agvc.AgvcNbq, 0)

    // 解析结果
    for result.Next() {
        record := result.Record()
        //truncate到5分钟
        ctime := record.Time().Format(time.RFC3339)

        allPointValues := agvc.GetAllPointValues()

        nbqHis := &agvc.AgvcNbq{
            Ctime:      &ctime,
            Psid:       &psid,
            InverterNo: &eqid,
            Name:       &eqName,
        }

        for _, point := range allPointValues {

            value := record.ValueByKey("point_" + point)
            if value == nil {
                continue
            }
            valueFloat64, _ := value.(float64)
            // 使用缓存的映射关系设置字段值，NBQ历史查询使用YC类型(type=2)
            if fieldName, found := agvc.GetFieldNameByPoint(cons.YC, point); found {
                setFieldValue(nbqHis, fieldName, valueFloat64)
            }
        }
        nbqhises = append(nbqhises, *nbqHis)

        // 获取或创建对应时间的AgvcNbqHis

        // 根据设备类型设置设备编号字段
        // NBQ 使用 InverterNo，其他类型可以扩展
        // if eqType == agvc.EqTypeNBQ {
        //     eqidInt := 0
        //     fmt.Sscanf(eqid, "%d", &eqidInt)
        //     nbqHis.InverterNo = &eqidInt
        // }

    }

    // 检查是否有错误
    if result.Err() != nil {
        global.GVA_LOG.Error(fmt.Sprintf("解析InfluxDB结果失败: %v", result.Err()))
        return nil, result.Err()
    }

    // 如果没有查询到数据，返回空数组
    if len(nbqhises) == 0 {
        global.GVA_LOG.Info(fmt.Sprintf("设备%d在时间范围%s到%s内没有历史数据", eqid, startTime, endTime))
        return []agvc.AgvcNbq{}, nil
    }

    return nbqhises, nil
}

// getHistoryBwd 查询BWD历史数据
func (agvcHisService *AgvcHisService) getHistoryBwd(ctx context.Context, psid, eqid int, pointValues []string, startTime, endTime string) ([]agvc.AgvcBwd, error) {
    // 检查InfluxDB客户端是否已初始化
    if global.GVA_INFLUXDB == nil {
        return nil, fmt.Errorf("InfluxDB客户端未初始化")
    }

    measurement := global.GVA_CONFIG.InfluxDB.GetMeasurement()

    // 获取查询API
    queryAPI := global.GVA_INFLUXDB.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)

    // 如果没有指定点位，使用所有BWD点位
    if len(pointValues) == 0 {
        pointValues = agvc.GetAllPointValuesBwd()
        if len(pointValues) == 0 {
            return nil, fmt.Errorf("未找到需要查询的点位信息")
        }
    }

    // 构建Flux查询语句
    flux := fmt.Sprintf(`
        from(bucket: "%s")
            |> range(start: %s, stop: %s)
            |> filter(fn: (r) => r["_measurement"] == "%s")
            |> filter(fn: (r) => r["psid"] == "%d")
            |> filter(fn: (r) => r["eqid"] == "%d")
            |> filter(fn: (r) => r["eqType"] == "%s")
            |> filter(fn: (r) => r["dataType"] == "%d")`,
        global.GVA_CONFIG.InfluxDB.Bucket,
        startTime,
        endTime,
        measurement,
        psid,
        eqid,
        agvc.EqTypeBWD,
        cons.YC,
    )

    // 添加点位过滤
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
    timeDataMap := make(map[string]*agvc.AgvcBwd)

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

        // 获取或创建对应时间的数据对象
        bwdHis, exists := timeDataMap[timeKey]
        if !exists {
            ctimeStr := record.Time().Format(time.RFC3339)
            bwdHis = &agvc.AgvcBwd{
                Ctime: &ctimeStr,
                Psid:  &psid,
                Eqid:  &eqid,
            }
            timeDataMap[timeKey] = bwdHis
        }

        // 设置字段值
        value, ok := record.Value().(float64)
        if !ok {
            continue
        }

        // 使用缓存的映射关系设置字段值，BWD历史查询使用YC类型(type=2)
        if fieldName, found := agvc.GetFieldNameByPointBwd(cons.YC, point); found {
            setFieldValueBwd(bwdHis, fieldName, value)
        }
    }

    // 检查是否有错误
    if result.Err() != nil {
        global.GVA_LOG.Error(fmt.Sprintf("解析InfluxDB结果失败: %v", result.Err()))
        return nil, result.Err()
    }

    // 将map转换为slice
    historyData := make([]agvc.AgvcBwd, 0, len(timeDataMap))
    for _, data := range timeDataMap {
        historyData = append(historyData, *data)
    }

    // 如果没有查询到数据，返回空数组
    if len(historyData) == 0 {
        global.GVA_LOG.Info(fmt.Sprintf("设备%s在时间范围%s到%s内没有历史数据", eqid, startTime, endTime))
        return []agvc.AgvcBwd{}, nil
    }

    return historyData, nil
}

// getHistoryAgc 查询AGC历史数据
func (agvcHisService *AgvcHisService) getHistoryAgc(ctx context.Context, psid, eqid int, pointValues []string, startTime, endTime string) ([]agvc.AgvcAgc, error) {
    // 检查InfluxDB客户端是否已初始化
    if global.GVA_INFLUXDB == nil {
        return nil, fmt.Errorf("InfluxDB客户端未初始化")
    }

    measurement := global.GVA_CONFIG.InfluxDB.GetMeasurement()

    // 获取查询API
    queryAPI := global.GVA_INFLUXDB.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)

    // 如果没有指定点位，使用所有AGC点位
    if len(pointValues) == 0 {
        pointValues = agvc.GetAllPointValuesAgc()
        if len(pointValues) == 0 {
            return nil, fmt.Errorf("未找到需要查询的点位信息")
        }
    }

    // 构建Flux查询语句
    flux := fmt.Sprintf(`
        from(bucket: "%s")
            |> range(start: %s, stop: %s)
            |> filter(fn: (r) => r["_measurement"] == "%s")
            |> filter(fn: (r) => r["psid"] == "%d")
            |> filter(fn: (r) => r["eqid"] == "%d")
            |> filter(fn: (r) => r["eqType"] == "%s")
            |> filter(fn: (r) => r["dataType"] == "%d")`,
        global.GVA_CONFIG.InfluxDB.Bucket,
        startTime,
        endTime,
        measurement,
        psid,
        eqid,
        agvc.EqTypeAGC,
        cons.YC,
    )

    // 添加点位过滤
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
    timeDataMap := make(map[string]*agvc.AgvcAgc)

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

        // 获取或创建对应时间的数据对象
        agcHis, exists := timeDataMap[timeKey]
        if !exists {
            ctimeStr := record.Time().Format(time.RFC3339)
            agcHis = &agvc.AgvcAgc{
                Ctime: &ctimeStr,
                Psid:  &psid,
                Eqid:  &eqid,
            }
            timeDataMap[timeKey] = agcHis
        }

        // 设置字段值
        value, ok := record.Value().(float64)
        if !ok {
            // AGC的某些字段是int64类型
            if intValue, ok := record.Value().(int64); ok {
                // AGC的遥信字段使用YX类型(type=1)
                if fieldName, found := agvc.GetFieldNameByPointAgc(cons.YX, point); found {
                    setFieldValueAgcInt(agcHis, fieldName, intValue)
                }
                continue
            }
            continue
        }

        // 使用缓存的映射关系设置字段值，AGC历史查询使用YC类型(type=2)
        if fieldName, found := agvc.GetFieldNameByPointAgc(cons.YC, point); found {
            setFieldValueAgc(agcHis, fieldName, value)
        }
    }

    // 检查是否有错误
    if result.Err() != nil {
        global.GVA_LOG.Error(fmt.Sprintf("解析InfluxDB结果失败: %v", result.Err()))
        return nil, result.Err()
    }

    // 将map转换为slice
    historyData := make([]agvc.AgvcAgc, 0, len(timeDataMap))
    for _, data := range timeDataMap {
        historyData = append(historyData, *data)
    }

    // 如果没有查询到数据，返回空数组
    if len(historyData) == 0 {
        global.GVA_LOG.Info(fmt.Sprintf("设备%s在时间范围%s到%s内没有历史数据", eqid, startTime, endTime))
        return []agvc.AgvcAgc{}, nil
    }

    return historyData, nil
}

// getHistoryAvc 查询AVC历史数据
func (agvcHisService *AgvcHisService) getHistoryAvc(ctx context.Context, psid, eqid int, pointValues []string, startTime, endTime string) ([]agvc.AgvcAvc, error) {
    // 检查InfluxDB客户端是否已初始化
    if global.GVA_INFLUXDB == nil {
        return nil, fmt.Errorf("InfluxDB客户端未初始化")
    }

    measurement := global.GVA_CONFIG.InfluxDB.GetMeasurement()

    // 获取查询API
    queryAPI := global.GVA_INFLUXDB.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)

    // 如果没有指定点位，使用所有AVC点位
    if len(pointValues) == 0 {
        pointValues = agvc.GetAllPointValuesAvc()
        if len(pointValues) == 0 {
            return nil, fmt.Errorf("未找到需要查询的点位信息")
        }
    }

    // 构建Flux查询语句
    flux := fmt.Sprintf(`
        from(bucket: "%s")
            |> range(start: %s, stop: %s)
            |> filter(fn: (r) => r["_measurement"] == "%s")
            |> filter(fn: (r) => r["psid"] == "%d")
            |> filter(fn: (r) => r["eqid"] == "%d")
            |> filter(fn: (r) => r["eqType"] == "%s")
            |> filter(fn: (r) => r["dataType"] == "%d")`,
        global.GVA_CONFIG.InfluxDB.Bucket,
        startTime,
        endTime,
        measurement,
        psid,
        eqid,
        agvc.EqTypeAVC,
        cons.YC,
    )

    // 添加点位过滤
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
    timeDataMap := make(map[string]*agvc.AgvcAvc)

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

        // 获取或创建对应时间的数据对象
        avcHis, exists := timeDataMap[timeKey]
        if !exists {
            ctimeStr := record.Time().Format(time.RFC3339)
            avcHis = &agvc.AgvcAvc{
                Ctime: &ctimeStr,
                Psid:  &psid,
                Eqid:  &eqid,
            }
            timeDataMap[timeKey] = avcHis
        }

        // 设置字段值
        value, ok := record.Value().(float64)
        if !ok {
            // AVC的某些字段是int64类型
            if intValue, ok := record.Value().(int64); ok {
                // AVC的遥信字段使用YX类型(type=1)
                if fieldName, found := agvc.GetFieldNameByPointAvc(cons.YX, point); found {
                    setFieldValueAvcInt(avcHis, fieldName, intValue)
                }
                continue
            }
            continue
        }

        // 使用缓存的映射关系设置字段值，AVC历史查询使用YC类型(type=2)
        if fieldName, found := agvc.GetFieldNameByPointAvc(cons.YC, point); found {
            setFieldValueAvc(avcHis, fieldName, value)
        }
    }

    // 检查是否有错误
    if result.Err() != nil {
        global.GVA_LOG.Error(fmt.Sprintf("解析InfluxDB结果失败: %v", result.Err()))
        return nil, result.Err()
    }

    // 将map转换为slice
    historyData := make([]agvc.AgvcAvc, 0, len(timeDataMap))
    for _, data := range timeDataMap {
        historyData = append(historyData, *data)
    }

    // 如果没有查询到数据，返回空数组
    if len(historyData) == 0 {
        global.GVA_LOG.Info(fmt.Sprintf("设备%s在时间范围%s到%s内没有历史数据", eqid, startTime, endTime))
        return []agvc.AgvcAvc{}, nil
    }

    return historyData, nil
}

// setFieldValue 根据字段名设置AgvcNbq结构体对应字段的值
func setFieldValue(nbqHis *agvc.AgvcNbq, fieldName string, value float64) {
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

// setFieldValueBwd 根据字段名设置AgvcBwd结构体对应字段的值
func setFieldValueBwd(bwdHis *agvc.AgvcBwd, fieldName string, value float64) {
    v := reflect.ValueOf(bwdHis).Elem()
    field := v.FieldByName(fieldName)

    if !field.IsValid() || !field.CanSet() {
        return
    }

    // 设置指针类型的float64值
    if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.Float64 {
        field.Set(reflect.ValueOf(&value))
    }
}

// setFieldValueAgc 根据字段名设置AgvcAgc结构体对应字段的值（float64类型）
func setFieldValueAgc(agcHis *agvc.AgvcAgc, fieldName string, value float64) {
    v := reflect.ValueOf(agcHis).Elem()
    field := v.FieldByName(fieldName)

    if !field.IsValid() || !field.CanSet() {
        return
    }

    // 设置指针类型的float64值
    if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.Float64 {
        field.Set(reflect.ValueOf(&value))
    }
}

// setFieldValueAgcInt 根据字段名设置AgvcAgc结构体对应字段的值（int64类型）
func setFieldValueAgcInt(agcHis *agvc.AgvcAgc, fieldName string, value int64) {
    v := reflect.ValueOf(agcHis).Elem()
    field := v.FieldByName(fieldName)

    if !field.IsValid() || !field.CanSet() {
        return
    }

    // 设置指针类型的int64值
    if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.Int64 {
        field.Set(reflect.ValueOf(&value))
    }
}

// setFieldValueAvc 根据字段名设置AgvcAvc结构体对应字段的值（float64类型）
func setFieldValueAvc(avcHis *agvc.AgvcAvc, fieldName string, value float64) {
    v := reflect.ValueOf(avcHis).Elem()
    field := v.FieldByName(fieldName)

    if !field.IsValid() || !field.CanSet() {
        return
    }

    // 设置指针类型的float64值
    if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.Float64 {
        field.Set(reflect.ValueOf(&value))
    }
}

// setFieldValueAvcInt 根据字段名设置AgvcAvc结构体对应字段的值（int64类型）
func setFieldValueAvcInt(avcHis *agvc.AgvcAvc, fieldName string, value int64) {
    v := reflect.ValueOf(avcHis).Elem()
    field := v.FieldByName(fieldName)

    if !field.IsValid() || !field.CanSet() {
        return
    }

    // 设置指针类型的int64值
    if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.Int64 {
        field.Set(reflect.ValueOf(&value))
    }
}
