package agvc

import (
    "context"
    "fmt"
    "reflect"
    "strconv"
    "strings"
    "time"

    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
    agvcReq "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/request"
    "github.com/flipped-aurora/gin-vue-admin/server/service/agvc/cons"
    influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type AgvcNbqHisService struct{}

func init() {
    // 系统启动时初始化point标签映射
    agvc.InitPointMapping()
}

// CreateAgvcNbqHis 创建agvcNbqHis表记录
// Author [yourname](https://github.com/yourname)
func (agvcNbqHisService *AgvcNbqHisService) CreateAgvcNbqHis(ctx context.Context, agvcNbqHis *agvc.AgvcNbqHis) (err error) {
    err = global.GVA_DB.Create(agvcNbqHis).Error
    return err
}

// DeleteAgvcNbqHis 删除agvcNbqHis表记录
// Author [yourname](https://github.com/yourname)
func (agvcNbqHisService *AgvcNbqHisService) DeleteAgvcNbqHis(ctx context.Context, ID string) (err error) {
    err = global.GVA_DB.Delete(&agvc.AgvcNbqHis{}, "id = ?", ID).Error
    return err
}

// DeleteAgvcNbqHisByIds 批量删除agvcNbqHis表记录
// Author [yourname](https://github.com/yourname)
func (agvcNbqHisService *AgvcNbqHisService) DeleteAgvcNbqHisByIds(ctx context.Context, IDs []string) (err error) {
    err = global.GVA_DB.Delete(&[]agvc.AgvcNbqHis{}, "id in ?", IDs).Error
    return err
}

// UpdateAgvcNbqHis 更新agvcNbqHis表记录
// Author [yourname](https://github.com/yourname)
func (agvcNbqHisService *AgvcNbqHisService) UpdateAgvcNbqHis(ctx context.Context, agvcNbqHis agvc.AgvcNbqHis) (err error) {
    err = global.GVA_DB.Model(&agvc.AgvcNbqHis{}).Where("id = ?", agvcNbqHis.InverterNo).Updates(&agvcNbqHis).Error
    return err
}

// GetAgvcNbqHis 根据ID获取agvcNbqHis表记录
// Author [yourname](https://github.com/yourname)
func (agvcNbqHisService *AgvcNbqHisService) GetAgvcNbqHis(ctx context.Context, ID string) (agvcNbqHis agvc.AgvcNbqHis, err error) {
    err = global.GVA_DB.Where("id = ?", ID).First(&agvcNbqHis).Error
    return
}

// GetAgvcNbqHisInfoList 分页获取逆变器历史数据列表
// 从InfluxDB中查询设备列表及其最新数据
func (agvcNbqHisService *AgvcNbqHisService) GetAgvcNbqHisInfoList(ctx context.Context, info agvcReq.AgvcNbqHisSearch) (list []agvc.AgvcNbqHis, total int64, err error) {
    limit := info.PageSize
    offset := info.PageSize * (info.Page - 1)

    // 创建db，从AgvcNbqSetting获取逆变器配置列表
    db := global.GVA_DB.Model(&agvc.AgvcNbqSetting{})
    var inverterSettings []agvc.AgvcNbqSetting

    // 如果有条件搜索 下方会自动创建搜索语句
    if len(info.CreatedAtRange) == 2 {
        db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
    }

    if info.Psid != nil {
        db = db.Where("psid = ?", *info.Psid)
    }
    if info.Inverter_no != nil {
        db = db.Where("inverter_no LIKE ?", "%"+*info.Inverter_no+"%")
    }
    if info.Name != nil && *info.Name != "" {
        db = db.Where("name LIKE ?", "%"+*info.Name+"%")
    }
    err = db.Count(&total).Error
    if err != nil {
        return
    }
    err = db.Find(&inverterSettings).Error
    if err != nil {
        return
    }

    // 检查InfluxDB客户端是否已初始化
    if global.GVA_INFLUXDB == nil {
        return nil, 0, fmt.Errorf("InfluxDB客户端未初始化")
    }

    // 获取查询API
    queryAPI := global.GVA_INFLUXDB.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)

    agvcNbqHises := make([]agvc.AgvcNbqHis, 0)

    // 对每个逆变器查询其最新数据
    for _, setting := range inverterSettings {
        // 构建Flux查询语句 - 查询该逆变器的所有点位的最新值
        flux := fmt.Sprintf(`
        from(bucket: "%s")
            |> range(start: -1y)
            |> filter(fn: (r) => r["_measurement"] == "%s")
                        |> filter(fn: (r) => r["psid"] == "%d")
            |> filter(fn: (r) => r["eqid"] == "%d")
            |> filter(fn: (r) => r["eqType"] == "%s")
            |> filter(fn: (r) => r["dataType"] == "%d")`,
            global.GVA_CONFIG.InfluxDB.Bucket, global.GVA_CONFIG.InfluxDB.GetMeasurement(), *setting.Psid, *setting.InverterNo, agvc.EqTypeNBQ, cons.YC)

        flux += `
            |> last()`

        // 执行查询
        result, err := queryAPI.Query(ctx, flux)
        if err != nil {
            global.GVA_LOG.Error(fmt.Sprintf("查询逆变器%d的InfluxDB数据失败: %v", *setting.InverterNo, err))
            continue
        }

        // 创建一个AgvcNbqHis对象
        nbqHis := agvc.AgvcNbqHis{
            Psid:       setting.Psid,
            InverterNo: setting.InverterNo,
            Name:       setting.Name,
        }

        // 解析结果，将各个点位的值填充到结构体中
        for result.Next() {
            record := result.Record()
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

            // 根据point值设置对应的字段
            setFieldByPoint(&nbqHis, point, value)
        }

        if result.Err() != nil {
            global.GVA_LOG.Error(fmt.Sprintf("解析逆变器%d的InfluxDB结果失败: %v", *setting.InverterNo, result.Err()))
        }

        agvcNbqHises = append(agvcNbqHises, nbqHis)
    }

    // 对inverterMonitors手动分页
    // 计算总数
    total = int64(len(agvcNbqHises))

    // 应用分页
    offset = info.PageSize * (info.Page - 1)
    limit = info.PageSize

    if offset >= int(total) {
        return []agvc.AgvcNbqHis{}, total, nil
    }

    end := offset + limit
    if end > int(total) {
        end = int(total)
    }

    if limit != 0 {
        return agvcNbqHises[offset:end], total, nil
    }

    return agvcNbqHises, total, err
}

// setFieldByPoint 根据point值设置AgvcNbqHis结构体对应的字段
func setFieldByPoint(nbqHis *agvc.AgvcNbqHis, point string, value float64) {
    // 使用反射获取结构体
    v := reflect.ValueOf(nbqHis).Elem()
    t := v.Type()

    // 遍历结构体字段，找到匹配的point标签
    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i)
        tag := t.Field(i).Tag.Get("point")

        if tag == "" {
            continue
        }

        // 解析point标签，提取value值
        if strings.Contains(tag, "value:") {
            parts := strings.Split(tag, ",")
            for _, part := range parts {
                part = strings.TrimSpace(part)
                if strings.HasPrefix(part, "value:") {
                    pointValue := strings.TrimPrefix(part, "value:")
                    if pointValue == point {
                        // 找到匹配的字段，设置值
                        if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.Float64 {
                            field.Set(reflect.ValueOf(&value))
                            return
                        }
                    }
                }
            }
        }
    }
}

func (agvcNbqHisService *AgvcNbqHisService) GetAgvcNbqHisPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}

// GetHistory 通用的历史数据查询方法
// psid: 电站编号, eqid: 设备编号, eqType: 设备类型, pointValues: 点位列表, startTime/endTime: 时间范围
func (agvcNbqHisService *AgvcNbqHisService) GetHistory(ctx context.Context, psid, eqid int, eqType string, pointValues []string, startTime, endTime string) ([]agvc.AgvcNbqHis, error) {
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
            |> filter(fn: (r) => r["eqid"] == "%d")
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
                Ctime:      &ctimeStr,
                Psid:       &psid,
                InverterNo: &eqid,
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
        global.GVA_LOG.Info(fmt.Sprintf("设备%d在时间范围%s到%s内没有历史数据", eqid, startTime, endTime))
        return []agvc.AgvcNbqHis{}, nil
    }

    return historyData, nil
}

// GetAgvcNbqHistory 获取逆变器历史数据
// 调用通用GetHistory方法，传入设备类型
func (agvcNbqHisService *AgvcNbqHisService) GetAgvcNbqHistory(ctx context.Context, nbqHis agvc.AgvcNbqHis, startTime, endTime string) ([]agvc.AgvcNbqHis, error) {
    // 检查必要的参数
    if nbqHis.InverterNo == nil {
        return nil, fmt.Errorf("逆变器编号不能为空")
    }
    if nbqHis.Psid == nil {
        return nil, fmt.Errorf("电站编号不能为空")
    }

    // 提取结构体中所有带point标签的字段的value值
    pointValues := extractPointValues(nbqHis)

    // 如果没有提取到任何point值，返回错误
    if len(pointValues) == 0 {
        return nil, fmt.Errorf("未找到需要查询的点位信息")
    }

    // 调用通用查询方法
    return agvcNbqHisService.GetHistory(ctx, *nbqHis.Psid, *nbqHis.InverterNo, agvc.EqTypeNBQ, pointValues, startTime, endTime)
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

// extractPointValues 从AgvcNbqHis结构体中提取所有带point标签的字段的value值
func extractPointValues(nbqHis agvc.AgvcNbqHis) []string {
    var pointValues []string

    // 使用反射获取结构体
    v := reflect.ValueOf(nbqHis)
    t := v.Type()

    // 遍历结构体字段，找到所有带point标签的字段
    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i)
        tag := t.Field(i).Tag.Get("point")

        if tag == "" {
            continue
        }

        // 检查字段值是否为nil（因为所有字段都是指针类型）
        if field.Kind() == reflect.Ptr && field.IsNil() {
            continue
        }

        // 解析point标签，提取value值
        if strings.Contains(tag, "value:") {
            parts := strings.Split(tag, ",")
            for _, part := range parts {
                part = strings.TrimSpace(part)
                if strings.HasPrefix(part, "value:") {
                    pointValue := strings.TrimPrefix(part, "value:")
                    pointValues = append(pointValues, pointValue)
                    break
                }
            }
        }
    }

    return pointValues
}

// GetNbqRealData 获取逆变器实时数据（从InfluxDB中获取最后一条数据）
func (agvcNbqHisService *AgvcNbqHisService) GetNbqRealData(ctx context.Context, psid, inverterNo int) (*agvc.AgvcNbqHis, error) {
    // 检查InfluxDB客户端是否已初始化
    if global.GVA_INFLUXDB == nil {
        return nil, fmt.Errorf("InfluxDB客户端未初始化")
    }

    // 获取查询API
    queryAPI := global.GVA_INFLUXDB.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)

    // 构建Flux查询语句 - 查询该逆变器所有点位的最新值
    flux := fmt.Sprintf(`
        from(bucket: "%s")
            |> range(start: -1y)
            |> filter(fn: (r) => r["_measurement"] == "%s")
            |> filter(fn: (r) => r["psid"] == "%d")
            |> filter(fn: (r) => r["eqid"] == "%d")
            |> filter(fn: (r) => r["eqType"] == "%s")
            |> filter(fn: (r) => r["dataType"] == "%d")
            |> last()`,
        global.GVA_CONFIG.InfluxDB.Bucket,
        global.GVA_CONFIG.InfluxDB.GetMeasurement(),
        psid,
        inverterNo,
        agvc.EqTypeNBQ,
        cons.YC,
    )

    // 执行查询
    result, err := queryAPI.Query(ctx, flux)
    if err != nil {
        global.GVA_LOG.Error(fmt.Sprintf("查询逆变器%d的实时数据失败: %v", inverterNo, err))
        return nil, err
    }

    // 创建AgvcNbqHis对象
    nbqHis := &agvc.AgvcNbqHis{
        Psid:       &psid,
        InverterNo: &inverterNo,
    }

    // 解析结果
    for result.Next() {
        record := result.Record()
        
        // 设置采集时间
        if nbqHis.Ctime == nil {
            ctimeStr := record.Time().Format(time.RFC3339)
            nbqHis.Ctime = &ctimeStr
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
            setFieldValue(nbqHis, fieldName, value)
        }
    }

    // 检查是否有错误
    if result.Err() != nil {
        global.GVA_LOG.Error(fmt.Sprintf("解析逆变器%d的实时数据失败: %v", inverterNo, result.Err()))
        return nil, result.Err()
    }

    // 如果没有采集时间，说明没有数据
    if nbqHis.Ctime == nil {
        return nil, fmt.Errorf("未找到逆变器%d的实时数据", inverterNo)
    }

    return nbqHis, nil
}

// WriteAgvcNbqDataToInfluxDB 将逆变器数据写入InfluxDB
// 该方法用于将CoAP或其他来源的数据写入InfluxDB
func (agvcNbqHisService *AgvcNbqHisService) WriteAgvcNbqDataToInfluxDB(ctx context.Context, psid int, eqid int, point string, value float64) error {
    // 检查InfluxDB客户端是否已初始化
    if global.GVA_INFLUXDB == nil {
        return fmt.Errorf("InfluxDB客户端未初始化")
    }

    // 获取写API
    writeAPI := global.GVA_INFLUXDB.WriteAPIBlocking(global.GVA_CONFIG.InfluxDB.Org, global.GVA_CONFIG.InfluxDB.Bucket)

    // 创建数据点
    p := influxdb2.NewPoint(
        global.GVA_CONFIG.InfluxDB.GetMeasurement(),
        map[string]string{
            "psid":     strconv.Itoa(psid),
            "eqType":   agvc.EqTypeNBQ,
            "eqid":     strconv.Itoa(eqid),
            "dataType": strconv.Itoa(cons.YC),
            "point":    point,
        },
        map[string]interface{}{
            "value": value,
        },
        time.Now(),
    )

    // 写入数据点
    if err := writeAPI.WritePoint(ctx, p); err != nil {
        return fmt.Errorf("写入InfluxDB失败: %v", err)
    }

    return nil
}
