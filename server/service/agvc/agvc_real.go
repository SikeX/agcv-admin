package agvc

import (
    "context"
    "fmt"
    "time"

    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
    "github.com/flipped-aurora/gin-vue-admin/server/service/agvc/cons"
)

// AgvcRealService 实时数据服务
type AgvcRealService struct{}

// GetRealData 通用的实时数据查询方法
// psid: 电站编号, eqid: 设备编号, eqType: 设备类型
// 从InfluxDB获取指定设备的最新数据
func (agvcRealService *AgvcRealService) GetRealData(ctx context.Context, psid, eqid, eqType int) (interface{}, error) {

    measurement := ""
    // 创建数据对象
    var res interface{}
    allPointValues := []string{}

    switch eqType {
    case NBQ_DEVICE_TYPE:
        measurement = global.GVA_CONFIG.InfluxDB.GetNBQMeasurement()
        allPointValues = agvc.GetAllPointValues()
    case BWD_DEVICE_TYPE:
        measurement = global.GVA_CONFIG.InfluxDB.GetMeasurement()
        allPointValues = agvc.GetAllPointValuesBwd()
    case AGC_DEVICE_TYPE:
        measurement = global.GVA_CONFIG.InfluxDB.GetMeasurement()
        allPointValues = agvc.GetAllPointValuesAgc()
    case AVC_DEVICE_TYPE:
        measurement = global.GVA_CONFIG.InfluxDB.GetMeasurement()
        allPointValues = agvc.GetAllPointValuesAvc()
    default:
        return nil, fmt.Errorf("未知设备类型: %d", eqType)
    }

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
            |> filter(fn: (r) => r["eqid"] == "%d")
            |> filter(fn: (r) => r["eqType"] == "%d")
                        |> aggregateWindow(every: 5s, fn: last, createEmpty: false)
                        |> pivot(rowKey: ["_time"], columnKey: ["_field"], valueColumn: "_value")
                        |> sort(columns: ["_time"], desc: true)
                        |> limit(n: 1)
                        `,
        global.GVA_CONFIG.InfluxDB.Bucket,
        measurement,
        psid,
        eqid,
        eqType,
    )

    // fmt.Println(flux)

    // 执行查询
    result, err := queryAPI.Query(ctx, flux)
    if err != nil {
        global.GVA_LOG.Error(fmt.Sprintf("查询设备%s的实时数据失败: %v", eqid, err))
        return nil, err
    }

    // 根据设备类型设置设备编号字段
    // if eqType == agvc.EqTypeNBQ {
    //     eqidInt := 0
    //     fmt.Sscanf(eqid, "%d", &eqidInt)
    //     realData.InverterNo = &eqidInt
    // }

    // 解析结果
    for result.Next() {
        record := result.Record()
        ctime := record.Time().Format(time.RFC3339)

        // 获取dataType标签
        dataTypeValue := record.ValueByKey("dataType")
        if dataTypeValue == nil {
            continue
        }
        dataTypeStr, ok := dataTypeValue.(string)
        if !ok {
            continue
        }

        // fmt.Println(record.Values())

        for _, point := range allPointValues {

            value := record.ValueByKey("point_" + point)
            if value == nil {
                continue
            }
            valueFloat64, _ := value.(float64)
            
            // 根据设备类型获取点位信息
            var pointDataType int
            switch eqType {
            case NBQ_DEVICE_TYPE:
                if pointInfo, found := agvc.GetPointInfoByPoint(point); found {
                    pointDataType = pointInfo.DataType
                }
            case BWD_DEVICE_TYPE:
                if pointInfo, found := agvc.GetPointInfoByPointBwd(point); found {
                    pointDataType = pointInfo.DataType
                }
            case AGC_DEVICE_TYPE:
                if pointInfo, found := agvc.GetPointInfoByPointAgc(point); found {
                    pointDataType = pointInfo.DataType
                }
            case AVC_DEVICE_TYPE:
                if pointInfo, found := agvc.GetPointInfoByPointAvc(point); found {
                    pointDataType = pointInfo.DataType
                }
            }
            
            // 判断dataType是否匹配
            if dataTypeStr != fmt.Sprintf("%d", pointDataType) {
                continue
            }
            
            // 使用缓存的映射关系设置字段值
            switch eqType {
            case NBQ_DEVICE_TYPE:
                resEntity := &agvc.AgvcNbq{
                    Ctime:      &ctime,
                    Psid:       &psid,
                    InverterNo: &eqid,
                }
                if fieldName, found := agvc.GetFieldNameByPoint(point); found {
                    setFieldValue(resEntity, fieldName, valueFloat64)
                }
                res = resEntity
            case BWD_DEVICE_TYPE:
                resEntity := &agvc.AgvcBwd{
                    Ctime: &ctime,
                    Psid:  &psid,
                    Eqid:  &eqid,
                }
                if fieldName, found := agvc.GetFieldNameByPointBwd(point); found {
                    setFieldValueBwd(resEntity, fieldName, valueFloat64)
                }
                res = resEntity
            case AGC_DEVICE_TYPE:
                resEntity := &agvc.AgvcAgc{
                    Ctime: &ctime,
                    Psid:  &psid,
                    Eqid:  &eqid,
                }
                if fieldName, found := agvc.GetFieldNameByPointAgc(point); found {
                    setFieldValueAgc(resEntity, fieldName, valueFloat64)
                }
                res = resEntity
            case AVC_DEVICE_TYPE:
                resEntity := &agvc.AgvcAvc{
                    Ctime: &ctime,
                    Psid:  &psid,
                    Eqid:  &eqid,
                }
                if fieldName, found := agvc.GetFieldNameByPointAvc(point); found {
                    setFieldValueAvc(resEntity, fieldName, valueFloat64)
                }
                res = resEntity
            }

        }

        // 检查是否有错误
        if result.Err() != nil {
            global.GVA_LOG.Error(fmt.Sprintf("解析设备%s的实时数据失败: %v", eqid, result.Err()))
            return nil, result.Err()
        }
    }

    return res, nil
}
