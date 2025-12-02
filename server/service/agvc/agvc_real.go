package agvc

import (
    "context"
    "fmt"
    "strconv"
    "time"

    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
)

// AgvcRealService 实时数据服务
type AgvcRealService struct{}

// GetRealData 通用的实时数据查询方法
// psid: 电站编号, eqid: 设备编号, eqType: 设备类型
// 从InfluxDB获取指定设备的最新数据
func (agvcRealService *AgvcRealService) GetRealData(ctx context.Context, psid, eqid, eqType int) (interface{}, error) {

    // 根据设备类型选择bucket，Measurement统一使用"agvc"
    var bucket string
    var res interface{}
    var allPointValues []string

    switch eqType {
    case NBQ_DEVICE_TYPE:
        bucket = global.GVA_CONFIG.InfluxDB.GetNbqBucket()
        allPointValues = agvc.GetAllPointValues()
        res = &agvc.AgvcNbq{
            Psid:       &psid,
            InverterNo: &eqid,
        }
    case BWD_DEVICE_TYPE:
        bucket = global.GVA_CONFIG.InfluxDB.GetAgvcBucket()
        allPointValues = agvc.GetAllPointValuesBwd()
        res = &agvc.AgvcBwd{
            Psid: &psid,
            Eqid: &eqid,
        }
    case AGC_DEVICE_TYPE:
        bucket = global.GVA_CONFIG.InfluxDB.GetAgvcBucket()
        allPointValues = agvc.GetAllPointValuesAgc()
        res = &agvc.AgvcAgc{
            Psid: &psid,
            Eqid: &eqid,
        }
    case AVC_DEVICE_TYPE:
        bucket = global.GVA_CONFIG.InfluxDB.GetAgvcBucket()
        allPointValues = agvc.GetAllPointValuesAvc()
        res = &agvc.AgvcAvc{
            Psid: &psid,
            Eqid: &eqid,
        }
    default:
        return nil, fmt.Errorf("未知设备类型: %d", eqType)
    }

    // 检查InfluxDB客户端是否已初始化
    if global.GVA_INFLUXDB == nil {
        return nil, fmt.Errorf("InfluxDB客户端未初始化")
    }

    // 获取查询API
    queryAPI := global.GVA_INFLUXDB.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)

    // 构建Flux查询语句 - 查询该设备所有点位的最新值，Measurement统一使用"agvc"
    flux := fmt.Sprintf(`
        from(bucket: "%s")
            |> range(start: -1h)
            |> filter(fn: (r) => r["_measurement"] == "%s")
            |> filter(fn: (r) => r["psid"] == "%d")
            |> filter(fn: (r) => r["eqid"] == "%d")
            |> filter(fn: (r) => r["eqType"] == "%d")
                        |> aggregateWindow(every: 5s, fn: last, createEmpty: false)
                        |> pivot(rowKey: ["_time"], columnKey: ["_field"], valueColumn: "_value")
                        |> sort(columns: ["_time"], desc: true)
                        |> limit(n: 1)
                        `,
        bucket,
        global.GVA_CONFIG.InfluxDB.GetMeasurement(),
        psid,
        eqid,
        eqType,
    )

    //fmt.Println(flux)

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

            dataTypeInt, _ := strconv.Atoi(dataTypeStr)

            // 根据设备类型获取点位信息，需要同时匹配dataType和point
            switch eqType {
            case NBQ_DEVICE_TYPE:
                if _, found := agvc.GetPointInfoByPoint(dataTypeInt, point); found {
                    if fieldName, found := agvc.GetFieldNameByPoint(dataTypeInt, point); found {
                        setFieldValue(res.(*agvc.AgvcNbq), fieldName, valueFloat64)
                    }
                    res.(*agvc.AgvcNbq).Ctime = &ctime
                }
            case BWD_DEVICE_TYPE:
                if _, found := agvc.GetPointInfoByPointBwd(dataTypeInt, point); found {
                    if fieldName, found := agvc.GetFieldNameByPointBwd(dataTypeInt, point); found {
                        setFieldValueBwd(res.(*agvc.AgvcBwd), fieldName, valueFloat64)
                    }
                    res.(*agvc.AgvcBwd).Ctime = &ctime
                }
            case AGC_DEVICE_TYPE:
                if _, found := agvc.GetPointInfoByPointAgc(dataTypeInt, point); found {
                    if fieldName, found := agvc.GetFieldNameByPointAgc(dataTypeInt, point); found {
                        setFieldValueAgc(res.(*agvc.AgvcAgc), fieldName, valueFloat64)
                    }
                    res.(*agvc.AgvcAgc).Ctime = &ctime
                }
            case AVC_DEVICE_TYPE:
                if _, found := agvc.GetPointInfoByPointAvc(dataTypeInt, point); found {
                    if fieldName, found := agvc.GetFieldNameByPointAvc(dataTypeInt, point); found {
                        setFieldValueAvc(res.(*agvc.AgvcAvc), fieldName, valueFloat64)
                    }
                    res.(*agvc.AgvcAvc).Ctime = &ctime
                }
            }
        }

        // 检查是否有错误
        if result.Err() != nil {
            global.GVA_LOG.Error(fmt.Sprintf("解析设备%d的实时数据失败: %v", eqid, result.Err()))
            return nil, result.Err()
        }
    }

    return res, nil
}
