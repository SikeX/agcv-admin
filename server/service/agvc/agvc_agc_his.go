package agvc

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
	agvcReq "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service/agvc/cons"
)

type AgvcAgcHisService struct{}

var agvcRealService = AgvcRealService{}

func init() {
	// 系统启动时初始化AGC的point标签映射
	agvc.InitPointMappingAgc()
}

// GetAgvcAgcHisInfoList 分页获取AGC历史数据列表
// 从InfluxDB中查询设备列表及其最新数据
func (agvcAgcHisService *AgvcAgcHisService) GetAgvcAgcHisInfoList(ctx context.Context, info agvcReq.AgvcAgcHisSearch) (list []agvc.AgvcAgc, total int64, err error) {
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

	agvcAgcHisList := make([]agvc.AgvcAgc, 0)

	// 对每个并网点查询其最新AGC数据
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
			global.GVA_CONFIG.InfluxDB.GetAgvcBucket(), global.GVA_CONFIG.InfluxDB.GetMeasurement(),
			*setting.Psid, *setting.Eqid, agvc.EqTypeBWD, cons.YC)

		flux += `
            |> last()`

		// 执行查询
		result, err := queryAPI.Query(ctx, flux)
		if err != nil {
			global.GVA_LOG.Error(fmt.Sprintf("查询AGC %s的InfluxDB数据失败: %v", *setting.Eqid, err))
			continue
		}

		// 创建一个AgvcAgcHis对象
		agcHis := agvc.AgvcAgc{
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
			setFieldByPointForAgc(&agcHis, point, value)
		}

		if result.Err() != nil {
			global.GVA_LOG.Error(fmt.Sprintf("解析AGC %s的InfluxDB结果失败: %v", *setting.Eqid, result.Err()))
		}

		agvcAgcHisList = append(agvcAgcHisList, agcHis)
	}

	// 手动分页
	total = int64(len(agvcAgcHisList))

	offset = info.PageSize * (info.Page - 1)
	limit = info.PageSize

	if offset >= int(total) {
		return []agvc.AgvcAgc{}, total, nil
	}

	end := offset + limit
	if end > int(total) {
		end = int(total)
	}

	if limit != 0 {
		return agvcAgcHisList[offset:end], total, nil
	}

	return agvcAgcHisList, total, err
}

// setFieldByPointForAgc 根据point值设置AgvcAgcHis结构体对应的字段
func setFieldByPointForAgc(agcHis *agvc.AgvcAgc, point string, value float64) {
	v := reflect.ValueOf(agcHis).Elem()
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

// GetAgvcAgcHistory 获取AGC历史数据
// func (agvcAgcHisService *AgvcAgcHisService) GetAgvcAgcHistory(ctx context.Context, agcHis agvc.AgvcAgc, startTime, endTime string) ([]map[string]interface{}, error) {
// 	// 检查InfluxDB客户端是否已初始化
// 	if global.GVA_INFLUXDB == nil {
// 		return nil, fmt.Errorf("InfluxDB客户端未初始化")
// 	}

// 	// 获取查询API
// 	queryAPI := global.GVA_INFLUXDB.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)

// 	// 提取结构体中所有带point标签的字段的value值
// 	pointValues := extractPointValuesFromAgc(agcHis)

// 	// 如果没有提取到任何point值，返回错误
// 	if len(pointValues) == 0 {
// 		return nil, fmt.Errorf("未找到需要查询的点位信息")
// 	}

// 	// 检查必要的参数
// 	if agcHis.Number == nil {
// 		return nil, fmt.Errorf("设备编号不能为空")
// 	}

// 	var historyData []map[string]interface{}

// 	// 对每个点位进行查询
// 	for _, pointValue := range pointValues {
// 		flux := fmt.Sprintf(`
//         from(bucket: "%s")
//             |> range(start: %s, stop: %s)
//             |> filter(fn: (r) => r["_measurement"] == "%s")
//             |> filter(fn: (r) => r["psid"] == "%d")
//             |> filter(fn: (r) => r["eqid"] == "%s")
//             |> filter(fn: (r) => r["eqType"] == "%s")
//             |> filter(fn: (r) => r["dataType"] == "%d")
//             |> filter(fn: (r) => r["_field"] == "point_%s")`,
// 			global.GVA_CONFIG.InfluxDB.Bucket,
// 			startTime,
// 			endTime,
// 			global.GVA_CONFIG.InfluxDB.GetMeasurement(),
// 			*agcHis.Psid,
// 			*agcHis.Number,
// 			agvc.EqTypeAGC,
// 			cons.YC,
// 			pointValue,
// 		)

// 		// 执行查询
// 		result, err := queryAPI.Query(ctx, flux)
// 		if err != nil {
// 			global.GVA_LOG.Error(fmt.Sprintf("查询点位%s的InfluxDB数据失败: %v", pointValue, err))
// 			continue
// 		}

// 		// 解析结果
// 		for result.Next() {
// 			record := result.Record()

// 			data := map[string]interface{}{
// 				"time":  record.Time().Format(time.RFC3339),
// 				"psid":  record.ValueByKey("psid"),
// 				"eqid":  record.ValueByKey("eqid"),
// 				"point": pointValue,
// 				"value": record.Value(),
// 			}
// 			historyData = append(historyData, data)
// 		}

// 		if result.Err() != nil {
// 			global.GVA_LOG.Error(fmt.Sprintf("解析点位%s的InfluxDB结果失败: %v", pointValue, result.Err()))
// 		}
// 	}

// 	// 如果没有查询到数据，返回空数组
// 	if len(historyData) == 0 {
// 		global.GVA_LOG.Info(fmt.Sprintf("AGC %s在时间范围%s到%s内没有历史数据", *agcHis.Number, startTime, endTime))
// 		return []map[string]interface{}{}, nil
// 	}

// 	return historyData, nil
// }

// extractPointValuesFromAgc 从AgvcAgcHis结构体中提取所有带point标签的字段的value值
func extractPointValuesFromAgc(agcHis agvc.AgvcAgc) []string {
	var pointValues []string
	v := reflect.ValueOf(agcHis)
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

// GetAgcRealData 获取AGC实时数据（从InfluxDB中获取最后一条数据）
// psid: 电站编号, eqid: 设备编号
func (agvcAgcHisService *AgvcAgcHisService) GetAgcRealData(ctx context.Context, psid, eqid int) (*agvc.AgvcAgc, error) {

	agvcRealData, err := agvcRealService.GetRealData(ctx, psid, eqid, AGC_DEVICE_TYPE)
	if err != nil {
		return nil, err
	}
	if agvcRealData == nil {
		return nil, fmt.Errorf("未找到AGC实时数据")
	}
	realData := agvcRealData.(*agvc.AgvcAgc)

	return realData, nil
}

// setFieldByPointForAgcRealData 根据point值设置AgvcAgc结构体对应的字段
func setFieldByPointForAgcRealData(agcData *agvc.AgvcAgc, point string, value float64) {
	v := reflect.ValueOf(agcData).Elem()
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
