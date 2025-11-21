package agvc

import (
    "context"
    "fmt"
    "reflect"
    "strings"
    "time"

    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
    agvcReq "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/request"
    "github.com/flipped-aurora/gin-vue-admin/server/service/agvc/cons"
)

type AgvcAvcHisService struct{}

func init() {
    // 系统启动时初始化AVC的point标签映射
    agvc.InitPointMappingAvc()
}

// GetAgvcAvcHisInfoList 分页获取AVC历史数据列表
// 从InfluxDB中查询设备列表及其最新数据
func (agvcAvcHisService *AgvcAvcHisService) GetAgvcAvcHisInfoList(ctx context.Context, info agvcReq.AgvcAvcHisSearch) (list []agvc.AgvcAvc, total int64, err error) {
    limit := info.PageSize
    offset := info.PageSize * (info.Page - 1)

    // 从AgvcBwdSetting获取并网点配置列表
    db := global.GVA_DB.Model(&agvc.AgvcBwdSetting{})
    var bwdSettings []agvc.AgvcBwdSetting

    // 条件搜索
    if len(info.CreatedAtRange) == 2 {
        db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
    }

    if info.Psid != nil {
        db = db.Where("psid = ?", *info.Psid)
    }
    if info.Number != nil {
        db = db.Where("number LIKE ?", "%"+*info.Number+"%")
    }
    if info.Name != nil && *info.Name != "" {
        db = db.Where("name LIKE ?", "%"+*info.Name+"%")
    }

    err = db.Count(&total).Error
    if err != nil {
        return
    }
    err = db.Find(&bwdSettings).Error
    if err != nil {
        return
    }

    // 检查InfluxDB客户端是否已初始化
    if global.GVA_INFLUXDB == nil {
        return nil, 0, fmt.Errorf("InfluxDB客户端未初始化")
    }

    // 获取查询API
    queryAPI := global.GVA_INFLUXDB.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)

    agvcAvcHisList := make([]agvc.AgvcAvc, 0)

    // 对每个并网点查询其最新AVC数据
    for _, setting := range bwdSettings {
        // 构建Flux查询语句
        flux := fmt.Sprintf(`
        from(bucket: "%s")
            |> range(start: -1y)
            |> filter(fn: (r) => r["_measurement"] == "%s")
            |> filter(fn: (r) => r["psid"] == "%d")
            |> filter(fn: (r) => r["eqid"] == "%s")
            |> filter(fn: (r) => r["eqType"] == "%s")
            |> filter(fn: (r) => r["dataType"] == "%d")`,
            global.GVA_CONFIG.InfluxDB.Bucket, global.GVA_CONFIG.InfluxDB.GetMeasurement(),
            *setting.Psid, *setting.Eqid, agvc.EqTypeAVC, cons.YC)

        flux += `
            |> last()`

        // 执行查询
        result, err := queryAPI.Query(ctx, flux)
        if err != nil {
            global.GVA_LOG.Error(fmt.Sprintf("查询AVC %s的InfluxDB数据失败: %v", *setting.Eqid, err))
            continue
        }

        // 创建一个AgvcAvcHis对象
        avcHis := agvc.AgvcAvc{
            Psid: setting.Psid,
            Eqid: setting.Eqid,
            Name: setting.Name,
        }

        // 解析结果
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
            setFieldByPointForAvc(&avcHis, point, value)
        }

        if result.Err() != nil {
            global.GVA_LOG.Error(fmt.Sprintf("解析AVC %s的InfluxDB结果失败: %v", *setting.Eqid, result.Err()))
        }

        agvcAvcHisList = append(agvcAvcHisList, avcHis)
    }

    // 手动分页
    total = int64(len(agvcAvcHisList))

    offset = info.PageSize * (info.Page - 1)
    limit = info.PageSize

    if offset >= int(total) {
        return []agvc.AgvcAvc{}, total, nil
    }

    end := offset + limit
    if end > int(total) {
        end = int(total)
    }

    if limit != 0 {
        return agvcAvcHisList[offset:end], total, nil
    }

    return agvcAvcHisList, total, err
}

// setFieldByPointForAvc 根据point值设置AgvcAvcHis结构体对应的字段
func setFieldByPointForAvc(avcHis *agvc.AgvcAvc, point string, value float64) {
    v := reflect.ValueOf(avcHis).Elem()
    t := v.Type()

    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i)
        tag := t.Field(i).Tag.Get("point")

        if tag == "" {
            continue
        }

        if strings.Contains(tag, "value:") {
            parts := strings.Split(tag, ",")
            for _, part := range parts {
                part = strings.TrimSpace(part)
                if strings.HasPrefix(part, "value:") {
                    pointValue := strings.TrimPrefix(part, "value:")
                    if pointValue == point {
                        if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.Float64 {
                            field.Set(reflect.ValueOf(&value))
                            return
                        } else if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.Int64 {
                            intValue := int64(value)
                            field.Set(reflect.ValueOf(&intValue))
                            return
                        }
                    }
                }
            }
        }
    }
}

// GetAgvcAvcHistory 获取AVC历史数据
func (agvcAvcHisService *AgvcAvcHisService) GetAgvcAvcHistory(ctx context.Context, avcHis agvc.AgvcAvc, startTime, endTime string) ([]map[string]interface{}, error) {
    // 检查InfluxDB客户端是否已初始化
    if global.GVA_INFLUXDB == nil {
        return nil, fmt.Errorf("InfluxDB客户端未初始化")
    }

    // 获取查询API
    queryAPI := global.GVA_INFLUXDB.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)

    // 提取结构体中所有带point标签的字段的value值
    pointValues := extractPointValuesFromAvc(avcHis)

    // 如果没有提取到任何point值，返回错误
    if len(pointValues) == 0 {
        return nil, fmt.Errorf("未找到需要查询的点位信息")
    }

    // 检查必要的参数
    if avcHis.Eqid == nil {
        return nil, fmt.Errorf("设备编号不能为空")
    }

    var historyData []map[string]interface{}

    // 对每个点位进行查询
    for _, pointValue := range pointValues {
        flux := fmt.Sprintf(`
        from(bucket: "%s")
            |> range(start: %s, stop: %s)
            |> filter(fn: (r) => r["_measurement"] == "%s")
            |> filter(fn: (r) => r["psid"] == "%d")
            |> filter(fn: (r) => r["eqid"] == "%s")
            |> filter(fn: (r) => r["eqType"] == "%s")
            |> filter(fn: (r) => r["dataType"] == "%d")
            |> filter(fn: (r) => r["_field"] == "point_%s")`,
            global.GVA_CONFIG.InfluxDB.Bucket,
            startTime,
            endTime,
            global.GVA_CONFIG.InfluxDB.GetMeasurement(),
            *avcHis.Psid,
            *avcHis.Eqid,
            agvc.EqTypeAVC,
            cons.YC,
            pointValue,
        )

        // 执行查询
        result, err := queryAPI.Query(ctx, flux)
        if err != nil {
            global.GVA_LOG.Error(fmt.Sprintf("查询点位%s的InfluxDB数据失败: %v", pointValue, err))
            continue
        }

        // 解析结果
        for result.Next() {
            record := result.Record()

            data := map[string]interface{}{
                "time":  record.Time().Format(time.RFC3339),
                "psid":  record.ValueByKey("psid"),
                "eqid":  record.ValueByKey("eqid"),
                "point": pointValue,
                "value": record.Value(),
            }
            historyData = append(historyData, data)
        }

        if result.Err() != nil {
            global.GVA_LOG.Error(fmt.Sprintf("解析点位%s的InfluxDB结果失败: %v", pointValue, result.Err()))
        }
    }

    // 如果没有查询到数据，返回空数组
    if len(historyData) == 0 {
        global.GVA_LOG.Info(fmt.Sprintf("AVC %s在时间范围%s到%s内没有历史数据", *avcHis.Eqid, startTime, endTime))
        return []map[string]interface{}{}, nil
    }

    return historyData, nil
}

// extractPointValuesFromAvc 从AgvcAvcHis结构体中提取所有带point标签的字段的value值
func extractPointValuesFromAvc(avcHis agvc.AgvcAvc) []string {
    var pointValues []string
    v := reflect.ValueOf(avcHis)
    t := v.Type()

    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i)
        tag := t.Field(i).Tag.Get("point")

        if tag == "" {
            continue
        }

        if field.Kind() == reflect.Ptr && field.IsNil() {
            continue
        }

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

// GetAvcRealData 获取AVC实时数据（从InfluxDB中获取最后一条数据）
// psid: 电站编号, eqid: 并网点编号
func (agvcAvcHisService *AgvcAvcHisService) GetAvcRealData(ctx context.Context, psid, eqid int) (*agvc.AgvcAvc, error) {
    agvcRealData, err := agvcRealService.GetRealData(ctx, psid, eqid, AVC_DEVICE_TYPE)
    if err != nil {
        return nil, err
    }
    if agvcRealData == nil {
        return nil, fmt.Errorf("未找到AVC实时数据")
    }
    realData := agvcRealData.(*agvc.AgvcAvc)

    return realData, nil
}

// setFieldByPointForAvcRealData 根据point值设置AgvcAvc结构体对应的字段
func setFieldByPointForAvcRealData(avcData *agvc.AgvcAvc, point string, value float64) {
    v := reflect.ValueOf(avcData).Elem()
    t := v.Type()

    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i)
        tag := t.Field(i).Tag.Get("point")

        if tag == "" {
            continue
        }

        if strings.Contains(tag, "value:") {
            parts := strings.Split(tag, ",")
            for _, part := range parts {
                part = strings.TrimSpace(part)
                if strings.HasPrefix(part, "value:") {
                    pointValue := strings.TrimPrefix(part, "value:")
                    if pointValue == point {
                        if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.Float64 {
                            field.Set(reflect.ValueOf(&value))
                            return
                        } else if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.Int64 {
                            intValue := int64(value)
                            field.Set(reflect.ValueOf(&intValue))
                            return
                        }
                    }
                }
            }
        }
    }
}
