package agvc

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
	agvcReq "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/request"
	agcvMain "github.com/flipped-aurora/gin-vue-admin/server/service/agvc/agcv_main"
	"github.com/flipped-aurora/gin-vue-admin/server/service/agvc/cons"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"go.uber.org/zap"
)

type AgvcBwdHisService struct{}

// CreateAgvcBwdHis 创建agvcBwdHis表记录
// Author [yourname](https://github.com/yourname)
func (agvcBwdHisService *AgvcBwdHisService) CreateAgvcBwdHis(ctx context.Context, agvcBwdHis *agvc.AgvcBwdHis) (err error) {
	err = global.GVA_DB.Create(agvcBwdHis).Error
	return err
}

// DeleteAgvcBwdHis 删除agvcBwdHis表记录
// Author [yourname](https://github.com/yourname)
func (agvcBwdHisService *AgvcBwdHisService) DeleteAgvcBwdHis(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&agvc.AgvcBwdHis{}, "id = ?", ID).Error
	return err
}

// DeleteAgvcBwdHisByIds 批量删除agvcBwdHis表记录
// Author [yourname](https://github.com/yourname)
func (agvcBwdHisService *AgvcBwdHisService) DeleteAgvcBwdHisByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]agvc.AgvcBwdHis{}, "id in ?", IDs).Error
	return err
}

// UpdateAgvcBwdHis 更新agvcBwdHis表记录
// Author [yourname](https://github.com/yourname)
// func (agvcBwdHisService *AgvcBwdHisService) UpdateAgvcBwdHis(ctx context.Context, agvcBwdHis agvc.AgvcBwdHis) (err error) {
//     err = global.GVA_DB.Model(&agvc.AgvcBwdHis{}).Where("id = ?", agvcBwdHis.ID).Updates(&agvcBwdHis).Error
//     return err
// }

// GetAgvcBwdHis 根据ID获取agvcBwdHis表记录
// Author [yourname](https://github.com/yourname)
func (agvcBwdHisService *AgvcBwdHisService) GetAgvcBwdHis(ctx context.Context, ID string) (agvcBwdHis agvc.AgvcBwdHis, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&agvcBwdHis).Error
	return
}

// GetAgvcBwdHisInfoList 分页获取并网点历史数据列表
// 从InfluxDB中查询设备列表及其最新数据
func (agvcBwdHisService *AgvcBwdHisService) GetAgvcBwdHisInfoList(ctx context.Context, info agvcReq.AgvcBwdHisSearch) (list []agvc.AgvcBwdHis, total int64, err error) {
	// 检查InfluxDB客户端是否已初始化
	if global.GVA_INFLUXDB == nil {
		return nil, 0, fmt.Errorf("InfluxDB客户端未初始化")
	}

	// 获取查询API
	queryAPI := global.GVA_INFLUXDB.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)

	// 构建Flux查询语句 - 获取所有并网点设备的最新数据
	// 使用group by eqid来获取每个设备的最新记录
	flux := fmt.Sprintf(`
    from(bucket: "%s")
        |> range(start: -30d)
        |> filter(fn: (r) => r["_measurement"] == "agvc")
        |> filter(fn: (r) => r["eqType"] == "5")
        |> group(columns: ["eqid"])
        |> last()
    `, global.GVA_CONFIG.InfluxDB.Bucket)

	// 执行查询
	result, err := queryAPI.Query(ctx, flux)
	if err != nil {
		return nil, 0, fmt.Errorf("查询InfluxDB失败: %v", err)
	}

	// 收集所有唯一设备ID
	deviceSet := make(map[string]bool)

	for result.Next() {
		record := result.Record()
		eqid := fmt.Sprintf("%v", record.ValueByKey("eqid"))
		deviceSet[eqid] = true
	}

	// 检查查询错误
	if result.Err() != nil {
		return nil, 0, fmt.Errorf("解析InfluxDB结果失败: %v", result.Err())
	}

	// 查询配置表获取设备名称
	var bwdSettings []agvc.AgvcBwdSetting
	global.GVA_DB.Find(&bwdSettings)

	// 创建设备编号到名称的映射
	numberToNameMap := make(map[string]string)
	for _, setting := range bwdSettings {
		if setting.Number != nil && setting.Name != nil {
			numberToNameMap[*setting.Number] = *setting.Name
		}
	}

	// 转换为AgvcBwdHis结构体列表
	var agvcBwdHisList []agvc.AgvcBwdHis
	for eqid := range deviceSet {
		number := eqid
		name := eqid // 默认使用设备编号

		// 如果在配置表中找到名称,则使用配置的名称
		if deviceName, ok := numberToNameMap[eqid]; ok {
			name = deviceName
		}

		bwdHis := agvc.AgvcBwdHis{
			Number: &number,
			Name:   &name,
		}
		agvcBwdHisList = append(agvcBwdHisList, bwdHis)
	}

	// 计算总数
	total = int64(len(agvcBwdHisList))

	// 应用分页
	offset := info.PageSize * (info.Page - 1)
	limit := info.PageSize

	if offset >= int(total) {
		return []agvc.AgvcBwdHis{}, total, nil
	}

	end := offset + limit
	if end > int(total) {
		end = int(total)
	}

	if limit != 0 {
		return agvcBwdHisList[offset:end], total, nil
	}

	return agvcBwdHisList, total, nil
}

func (agvcBwdHisService *AgvcBwdHisService) GetAgvcBwdHisPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}

// GetAgvcBwdHistoryNew 获取并网点历史数据(新版本，按AgvcNbqHis格式)
func (agvcBwdHisService *AgvcBwdHisService) GetAgvcBwdHistoryNew(ctx context.Context, bwdHis agvc.AgvcBwdHis, startTime, endTime string) ([]map[string]interface{}, error) {
	// 检查InfluxDB客户端是否已初始化
	if global.GVA_INFLUXDB == nil {
		return nil, fmt.Errorf("InfluxDB客户端未初始化")
	}

	// 获取查询API
	queryAPI := global.GVA_INFLUXDB.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)

	// 提取结构体中所有带point标签的字段的value值
	pointValues := extractPointValuesFromBwd(bwdHis)

	// 如果没有提取到任何point值，返回错误
	if len(pointValues) == 0 {
		return nil, fmt.Errorf("未找到需要查询的点位信息")
	}

	// 检查必要的参数
	if bwdHis.Number == nil {
		return nil, fmt.Errorf("设备编号不能为空")
	}

	var historyData []map[string]interface{}

	// 对每个点位进行查询
	for _, pointValue := range pointValues {
		// 构建Flux查询语句
		flux := fmt.Sprintf(`
        from(bucket: "%s")
            |> range(start: %s, stop: %s)
            |> filter(fn: (r) => r["_measurement"] == "%s")
            |> filter(fn: (r) => r["eqid"] == "%s")
            |> filter(fn: (r) => r["eqType"] == "%s")
            |> filter(fn: (r) => r["_field"] == "point_%s")`,
			global.GVA_CONFIG.InfluxDB.Bucket,
			startTime,
			endTime,
			global.GVA_CONFIG.InfluxDB.GetMeasurement(),
			*bwdHis.Number,
			agvc.EqTypeBWD,
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
				"time":   record.Time().Format(time.RFC3339),
				"eqid":   record.ValueByKey("eqid"),
				"eqType": record.ValueByKey("eqType"),
				"point":  pointValue,
				"value":  record.Value(),
			}
			historyData = append(historyData, data)
		}

		// 检查是否有错误
		if result.Err() != nil {
			global.GVA_LOG.Error(fmt.Sprintf("解析点位%s的InfluxDB结果失败: %v", pointValue, result.Err()))
		}
	}

	// 如果没有查询到数据，返回空数组而不是错误
	if len(historyData) == 0 {
		global.GVA_LOG.Info(fmt.Sprintf("并网点%s在时间范围%s到%s内没有历史数据", *bwdHis.Number, startTime, endTime))
		return []map[string]interface{}{}, nil
	}

	return historyData, nil
}

// extractPointValuesFromBwd 从AgvcBwdHis结构体中提取所有带point标签的字段的value值
func extractPointValuesFromBwd(bwdHis agvc.AgvcBwdHis) []string {
	var pointValues []string
	v := reflect.ValueOf(bwdHis)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("point")

		if tag == "" {
			continue
		}

		// 检查字段值是否为nil
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

// GetAgvcBwdHistory 获取并网点历史数据(旧版本，保持兼容)
func (agvcBwdHisService *AgvcBwdHisService) GetAgvcBwdHistory(ctx context.Context, eqid, startTime, endTime string) ([]map[string]interface{}, error) {
	// 检查InfluxDB客户端是否已初始化
	if global.GVA_INFLUXDB == nil {
		return nil, fmt.Errorf("InfluxDB客户端未初始化")
	}

	// 获取查询API
	queryAPI := global.GVA_INFLUXDB.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)

	// 构建Flux查询语句 - 使用统一的agvc表
	// eqType=5 代表并网点
	flux := fmt.Sprintf(`
    from(bucket: "%s")
        |> range(start: %s, stop: %s)
        |> filter(fn: (r) => r["_measurement"] == "agvc_data")
        |> filter(fn: (r) => r["eqid"] == "%s")
        |> filter(fn: (r) => r["eqType"] == "5")
        |> filter(fn: (r) => r["psid"] == "1")
    `, global.GVA_CONFIG.InfluxDB.Bucket, startTime, endTime, eqid)

	// 执行查询
	result, err := queryAPI.Query(ctx, flux)
	if err != nil {
		return nil, fmt.Errorf("查询InfluxDB失败: %v", err)
	}

	// 解析结果
	var historyData []map[string]interface{}
	for result.Next() {
		record := result.Record()
		point := fmt.Sprintf("%v", record.ValueByKey("point"))

		timestamp := record.Time().Truncate(time.Minute) // 向下取整到秒，确保同一秒的数据能正确聚合

		data := map[string]interface{}{
			"time":      timestamp,
			"psid":      record.ValueByKey("psid"),
			"eqid":      record.ValueByKey("eqid"),
			"eqType":    record.ValueByKey("eqType"),
			"dataType":  record.ValueByKey("dataType"),
			"point":     point,
			"pointName": agvc.GetBwdPointName(point), // 添加点位名称
			"value":     record.Value(),
		}
		historyData = append(historyData, data)
	}

	// 检查是否有错误
	if result.Err() != nil {
		return nil, fmt.Errorf("解析InfluxDB结果失败: %v", result.Err())
	}

	return historyData, nil
}

// GenerateAgvcBwdTestData 生成并网点测试数据
// 随机选择10-15个点位,每个点位生成10-20条时间序列数据
func (agvcBwdHisService *AgvcBwdHisService) GenerateAgvcBwdTestData(ctx context.Context, eqid string, count int) error {
	// 检查InfluxDB客户端是否已初始化
	if global.GVA_INFLUXDB == nil {
		return fmt.Errorf("InfluxDB客户端未初始化")
	}

	// 获取写API
	writeAPI := global.GVA_INFLUXDB.WriteAPIBlocking(global.GVA_CONFIG.InfluxDB.Org, global.GVA_CONFIG.InfluxDB.Bucket)

	// 获取所有并网点点位
	allPoints := agvc.GetAllBwdPoints()

	// 随机选择10-15个点位
	rand.Seed(time.Now().UnixNano())
	pointCount := len(allPoints) //10 + rand.Intn(6) // 10-15个点位
	selectedPoints := make([]string, 0, pointCount)

	//// 随机打乱点位顺序
	//rand.Shuffle(len(allPoints), func(i, j int) {
	//    allPoints[i], allPoints[j] = allPoints[j], allPoints[i]
	//})

	// 选择前pointCount个点位
	for i := 0; i < pointCount && i < len(allPoints); i++ {
		selectedPoints = append(selectedPoints, allPoints[i])
	}

	// 为每个选中的点位生成10-20条数据
	for _, point := range selectedPoints {
		dataCount := 10 + rand.Intn(11) // 10-20条数据
		now := time.Now()

		for i := 0; i < dataCount; i++ {
			timestamp := now.Add(-time.Duration(dataCount-i) * time.Minute)

			// 根据不同点位生成不同范围的模拟数值
			var value float64
			switch point {
			case "1", "2", "3": // 线电压(10kV)
				value = 10 + 0.5*math.Sin(float64(i)/10) + rand.Float64()*0.2 - 0.1
			case "12", "13", "14": // 相电压
				value = 230 + 10*math.Sin(float64(i)/8) + rand.Float64()*5
			case "4", "5", "6": // 电流
				value = 50 + 20*math.Sin(float64(i)/10) + rand.Float64()*5
			case "7": // 有功功率
				value = 100 + 50*math.Sin(float64(i)/12) + rand.Float64()*10
			case "8": // 无功功率
				value = 50 + 25*math.Sin(float64(i)/8) + rand.Float64()*5
			case "9": // 功率因数
				value = 0.95 + 0.04*math.Sin(float64(i)/15) + rand.Float64()*0.02 - 0.01
			case "10": // 频率
				value = 50 + 0.3*math.Sin(float64(i)/20) + rand.Float64()*0.1
			default:
				// 其他点位使用通用范围
				value = 100 + 50*math.Sin(float64(i)/8) + rand.Float64()*10 - 5
			}

			// 创建数据点
			p := influxdb2.NewPoint(
				"agvc",
				map[string]string{
					"psid":     "1",
					"eqid":     eqid,
					"eqType":   agvc.EqTypeBWD, // 使用常量
					"dataType": "1",            // 遥测
					"point":    point,
				},
				map[string]interface{}{
					"value": value,
				},
				timestamp,
			)

			// 写入数据点
			if err := writeAPI.WritePoint(ctx, p); err != nil {
				return fmt.Errorf("写入InfluxDB失败: %v", err)
			}
		}
	}

	return nil
}

// UpdatePlanCurves 更新计划曲线（AGC和AVC分开存储，AVC又分电压和无功）
// func (agvcBwdHisService *AgvcBwdHisService) UpdatePlanCurves(ctx context.Context, req agvcReq.PlanCurvesRequest) error {
//     // 将计划曲线序列化为JSON
//     curvesJSON, err := json.Marshal(req.PlanCurves)
//     if err != nil {
//         return fmt.Errorf("序列化计划曲线失败: %v", err)
//     }

//     curvesStr := string(curvesJSON)

//     // 更新或创建记录
//     var bwdHis agvc.AgvcBwdHis
//     err = global.GVA_DB.Where("number = ?", req.Number).First(&bwdHis).Error

//     if err != nil {
//         // 记录不存在，创建新记录
//         bwdHis = agvc.AgvcBwdHis{
//             Number: &req.Number,
//         }
//         // 根据curveType决定存储到哪个字段
//         if req.CurveType == "agc" {
//             bwdHis.AgcPlanCurves = &curvesStr
//         } else if req.CurveType == "avc-voltage" {
//             bwdHis.AvcVoltagePlanCurves = &curvesStr
//         } else if req.CurveType == "avc-reactive" {
//             bwdHis.AvcReactivePlanCurves = &curvesStr
//         }
//         return global.GVA_DB.Create(&bwdHis).Error
//     }

//     // 更新现有记录
//     if req.CurveType == "agc" {
//         bwdHis.AgcPlanCurves = &curvesStr
//     } else if req.CurveType == "avc-voltage" {
//         bwdHis.AvcVoltagePlanCurves = &curvesStr
//     } else if req.CurveType == "avc-reactive" {
//         bwdHis.AvcReactivePlanCurves = &curvesStr
//     }
//     return global.GVA_DB.Save(&bwdHis).Error
// }

// UpdateAgcStatus 更新AGC状态
// func (agvcBwdHisService *AgvcBwdHisService) UpdateAgcStatus(ctx context.Context, req agvcReq.AgcStatusRequest) error {
// 	// 更新或创建记录
// 	var bwdHis agvc.AgvcBwdHis
// 	err := global.GVA_DB.Where("number = ?", req.Number).First(&bwdHis).Error

// 	if err != nil {
// 		// 记录不存在，创建新记录
// 		bwdHis = agvc.AgvcBwdHis{
// 			Number:              &req.Number,
// 			AgcFunctionState:    &req.AgcFunctionState,
// 			AgcControlMode:      &req.AgcControlMode,
// 			AgcControlAuthority: &req.AgcControlAuthority,
// 		}
// 		return global.GVA_DB.Create(&bwdHis).Error
// 	}

// 	// 更新现有记录
// 	bwdHis.AgcFunctionState = &req.AgcFunctionState
// 	bwdHis.AgcControlMode = &req.AgcControlMode
// 	bwdHis.AgcControlAuthority = &req.AgcControlAuthority
// 	return global.GVA_DB.Save(&bwdHis).Error
// }

// GetAgcStatus 获取AGC状态
func (agvcBwdHisService *AgvcBwdHisService) GetAgcStatus(ctx context.Context, numberStr string) (map[string]interface{}, error) {

	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return nil, fmt.Errorf("无效的数字：%v", err)
	}

	agcDatas := agcvMain.DataStorage.GetDeviceData(1, number, AGC_DEVICE_TYPE, cons.YX)
	if agcDatas == nil {
		return nil, fmt.Errorf("未找到数据")
	}
	global.GVA_LOG.Info("AGC数据：", zap.Any("agcDatas", agcDatas))
	agcStatus := make(map[string]interface{})
	for point, data := range agcDatas {
		switch point {
		case AGC_YX_CONTROL_MODE:
			agcStatus["agcControlMode"] = data.Value
		case AGC_YX_SIGNAL:
			agcStatus["agcFunctionState"] = data.Value
		case AGC_YX_LOOP_STATUS:
			agcStatus["agcControlAuthority"] = data.Value
		}
	}

	return agcStatus, nil
}

// GetAgcParameters 获取AGC参数设置
// func (agvcBwdHisService *AgvcBwdHisService) GetAgcParameters(ctx context.Context, number string) (map[string]interface{}, error) {

// 	// number, err := strconv.Atoi(numberStr)
// 	// if err != nil {
// 	//     return nil, fmt.Errorf("转换设备编号失败: %v", err)
// 	// }
// 	// //从datastore中获取AGC参数
// 	// agcParameters, err := agcv_main.AGC.GetAGCConfigFromMem(number)
// 	var bwdHis agvc.AgvcBwdHis
// 	err := global.GVA_DB.Where("number = ?", number).First(&bwdHis).Error

// 	if err != nil {
// 		// 如果没有记录，创建一条新记录（只设置设备编号）
// 		bwdHis = agvc.AgvcBwdHis{
// 			Number: &number,
// 		}
// 		err = global.GVA_DB.Create(&bwdHis).Error
// 		if err != nil {
// 			return nil, fmt.Errorf("创建记录失败: %v", err)
// 		}
// 	}

// 	if bwdHis.AgcParameters == nil {
// 		// 如果没有配置，从配置表获取默认值
// 		var bwdSetting agvc.AgvcBwdSetting
// 		err = global.GVA_DB.Where("number = ?", number).First(&bwdSetting).Error
// 		if err != nil {
// 			return map[string]interface{}{}, nil
// 		}

// 		// 返回配置表中的AGC参数
// 		return map[string]interface{}{
// 			"controlMethod":        bwdSetting.AgcFunctionExit,
// 			"stepSize":             bwdSetting.AgcStepSize,
// 			"stepPeriod":           bwdSetting.AgcStepPeriod,
// 			"vibrationRange":       bwdSetting.AgcVibrationRange,
// 			"controlPeriod":        bwdSetting.AgcControlPeriod,
// 			"microAdjustmentCoeff": bwdSetting.AgcMicroAdjustmentCoefficient,
// 		}, nil
// 	}

// 	// 解析JSON
// 	var params map[string]interface{}
// 	err = json.Unmarshal([]byte(*bwdHis.AgcParameters), &params)
// 	if err != nil {
// 		return nil, fmt.Errorf("解析AGC参数失败: %v", err)
// 	}

// 	return params, nil
// }

// UpdateAgcParameters 更新AGC参数设置
// func (agvcBwdHisService *AgvcBwdHisService) UpdateAgcParameters(ctx context.Context, req agvcReq.AgcParametersRequest) error {
// 	// 将参数序列化为JSON
// 	paramsJSON, err := json.Marshal(req)
// 	if err != nil {
// 		return fmt.Errorf("序列化AGC参数失败: %v", err)
// 	}

// 	paramsStr := string(paramsJSON)

// 	// 更新或创建记录
// 	var bwdHis agvc.AgvcBwdHis
// 	err = global.GVA_DB.Where("number = ?", req.Number).First(&bwdHis).Error

// 	if err != nil {
// 		// 记录不存在，创建新记录
// 		bwdHis = agvc.AgvcBwdHis{
// 			Number:        &req.Number,
// 			AgcParameters: &paramsStr,
// 		}
// 		return global.GVA_DB.Create(&bwdHis).Error
// 	}

// 	// 更新现有记录
// 	bwdHis.AgcParameters = &paramsStr
// 	return global.GVA_DB.Save(&bwdHis).Error
// }

// GetAvcParameters 获取AVC参数设置
// func (agvcBwdHisService *AgvcBwdHisService) GetAvcParameters(ctx context.Context, number string) (map[string]interface{}, error) {
// 	var bwdHis agvc.AgvcBwdHis
// 	err := global.GVA_DB.Where("number = ?", number).First(&bwdHis).Error

// 	if err != nil {
// 		// 如果没有记录，创建一条新记录（只设置设备编号）
// 		bwdHis = agvc.AgvcBwdHis{
// 			Number: &number,
// 		}
// 		err = global.GVA_DB.Create(&bwdHis).Error
// 		if err != nil {
// 			return nil, fmt.Errorf("创建记录失败: %v", err)
// 		}
// 	}

// 	if bwdHis.AvcParameters == nil {
// 		// 如果没有配置，从配置表获取默认值
// 		var bwdSetting agvc.AgvcBwdSetting
// 		err = global.GVA_DB.Where("number = ?", number).First(&bwdSetting).Error
// 		if err != nil {
// 			return map[string]interface{}{}, nil
// 		}

// 		// 返回配置表中的AVC参数
// 		return map[string]interface{}{
// 			"stepSize":           bwdSetting.AvcStepSize,
// 			"stepPeriod":         bwdSetting.AvcStepPeriod,
// 			"vibrationRange":     bwdSetting.AvcVibrationRange,
// 			"controlPeriod":      bwdSetting.AvcControlPeriod,
// 			"systemImpedance":    bwdSetting.AvcSystemImpedance,
// 			"adjustmentRangeMin": bwdSetting.AvcAdjustmentRangeMin,
// 			"adjustmentRangeMax": bwdSetting.AvcAdjustmentRangeMax,
// 		}, nil
// 	}

// 	// 解析JSON
// 	var params map[string]interface{}
// 	err = json.Unmarshal([]byte(*bwdHis.AvcParameters), &params)
// 	if err != nil {
// 		return nil, fmt.Errorf("解析AVC参数失败: %v", err)
// 	}

// 	return params, nil
// }

// UpdateAvcParameters 更新AVC参数设置
// func (agvcBwdHisService *AgvcBwdHisService) UpdateAvcParameters(ctx context.Context, req agvcReq.AvcParametersRequest) error {
// 	// 将参数序列化为JSON
// 	paramsJSON, err := json.Marshal(req)
// 	if err != nil {
// 		return fmt.Errorf("序列化AVC参数失败: %v", err)
// 	}

// 	paramsStr := string(paramsJSON)

// 	// 更新或创建记录
// 	var bwdHis agvc.AgvcBwdHis
// 	err = global.GVA_DB.Where("number = ?", req.Number).First(&bwdHis).Error

// 	if err != nil {
// 		// 记录不存在，创建新记录
// 		bwdHis = agvc.AgvcBwdHis{
// 			Number:        &req.Number,
// 			AvcParameters: &paramsStr,
// 		}
// 		return global.GVA_DB.Create(&bwdHis).Error
// 	}

// 	// 更新现有记录
// 	bwdHis.AvcParameters = &paramsStr
// 	return global.GVA_DB.Save(&bwdHis).Error
// }

// GetPlanCurves 获取计划曲线（AGC和AVC分开存储，AVC又分电压和无功）
// func (agvcBwdHisService *AgvcBwdHisService) GetPlanCurves(ctx context.Context, number string, curveType string) (map[string]interface{}, error) {
// 	var bwdHis agvc.AgvcBwdHis
// 	err := global.GVA_DB.Where("number = ?", number).First(&bwdHis).Error

// 	if err != nil {
// 		// 如果没有记录，创建一条新记录（只设置设备编号）
// 		bwdHis = agvc.AgvcBwdHis{
// 			Number: &number,
// 		}
// 		err = global.GVA_DB.Create(&bwdHis).Error
// 		if err != nil {
// 			return nil, fmt.Errorf("创建记录失败: %v", err)
// 		}
// 		// 新创建的记录，返回默认空数据
// 		return map[string]interface{}{"localCurve": []agvcReq.PlanCurvePoint{}}, nil
// 	}

// 	// 根据curveType返回对应的曲线
// 	if curveType == "agc" {
// 		// 解析AGC计划曲线
// 		if bwdHis.AgcPlanCurves != nil && *bwdHis.AgcPlanCurves != "" {
// 			var agcCurves map[string]interface{}
// 			err = json.Unmarshal([]byte(*bwdHis.AgcPlanCurves), &agcCurves)
// 			if err == nil {
// 				return agcCurves, nil
// 			}
// 		}
// 		return map[string]interface{}{"localCurve": []agvcReq.PlanCurvePoint{}}, nil
// 	} else if curveType == "avc-voltage" {
// 		// 解析AVC电压计划曲线
// 		if bwdHis.AvcVoltagePlanCurves != nil && *bwdHis.AvcVoltagePlanCurves != "" {
// 			var avcCurves map[string]interface{}
// 			err = json.Unmarshal([]byte(*bwdHis.AvcVoltagePlanCurves), &avcCurves)
// 			if err == nil {
// 				return avcCurves, nil
// 			}
// 		}
// 		return map[string]interface{}{"localCurve": []agvcReq.PlanCurvePoint{}}, nil
// 	} else if curveType == "avc-reactive" {
// 		// 解析AVC无功计划曲线
// 		if bwdHis.AvcReactivePlanCurves != nil && *bwdHis.AvcReactivePlanCurves != "" {
// 			var avcCurves map[string]interface{}
// 			err = json.Unmarshal([]byte(*bwdHis.AvcReactivePlanCurves), &avcCurves)
// 			if err == nil {
// 				return avcCurves, nil
// 			}
// 		}
// 		return map[string]interface{}{"localCurve": []agvcReq.PlanCurvePoint{}}, nil
// 	}

// 	// 如果没有指定curveType，返回AGC和AVC各自的计划曲线
// 	result := make(map[string]interface{})

// 	// 解析AGC计划曲线
// 	if bwdHis.AgcPlanCurves != nil && *bwdHis.AgcPlanCurves != "" {
// 		var agcCurves map[string]interface{}
// 		err = json.Unmarshal([]byte(*bwdHis.AgcPlanCurves), &agcCurves)
// 		if err == nil {
// 			result["agc"] = agcCurves
// 		}
// 	} else {
// 		result["agc"] = map[string]interface{}{"localCurve": []agvcReq.PlanCurvePoint{}}
// 	}

// 	// 解析AVC电压计划曲线
// 	if bwdHis.AvcVoltagePlanCurves != nil && *bwdHis.AvcVoltagePlanCurves != "" {
// 		var avcVoltageCurves map[string]interface{}
// 		err = json.Unmarshal([]byte(*bwdHis.AvcVoltagePlanCurves), &avcVoltageCurves)
// 		if err == nil {
// 			result["avc-voltage"] = avcVoltageCurves
// 		}
// 	} else {
// 		result["avc-voltage"] = map[string]interface{}{"localCurve": []agvcReq.PlanCurvePoint{}}
// 	}

// 	// 解析AVC无功计划曲线
// 	if bwdHis.AvcReactivePlanCurves != nil && *bwdHis.AvcReactivePlanCurves != "" {
// 		var avcReactiveCurves map[string]interface{}
// 		err = json.Unmarshal([]byte(*bwdHis.AvcReactivePlanCurves), &avcReactiveCurves)
// 		if err == nil {
// 			result["avc-reactive"] = avcReactiveCurves
// 		}
// 	} else {
// 		result["avc-reactive"] = map[string]interface{}{"localCurve": []agvcReq.PlanCurvePoint{}}
// 	}

// 	return result, nil
// }

// UpdateAvcStatus 更新AVC状态
// func (agvcBwdHisService *AgvcBwdHisService) UpdateAvcStatus(ctx context.Context, req agvcReq.AvcStatusRequest) error {
// 	// 更新或创建记录
// 	var bwdHis agvc.AgvcBwdHis
// 	err := global.GVA_DB.Where("number = ?", req.Number).First(&bwdHis).Error

// 	if err != nil {
// 		// 记录不存在，创建新记录
// 		bwdHis = agvc.AgvcBwdHis{
// 			Number:              &req.Number,
// 			AvcFunctionState:    &req.AvcFunctionState,
// 			AvcControlMode:      &req.AvcControlMode,
// 			AvcControlAuthority: &req.AvcControlAuthority,
// 		}
// 		return global.GVA_DB.Create(&bwdHis).Error
// 	}

// 	// 更新现有记录
// 	bwdHis.AvcFunctionState = &req.AvcFunctionState
// 	bwdHis.AvcControlMode = &req.AvcControlMode
// 	bwdHis.AvcControlAuthority = &req.AvcControlAuthority
// 	return global.GVA_DB.Save(&bwdHis).Error
// }

// // GetAvcStatus 获取AVC状态
// func (agvcBwdHisService *AgvcBwdHisService) GetAvcStatus(ctx context.Context, numberStr string) (map[string]interface{}, error) {
// 	number, err := strconv.Atoi(numberStr)
// 	if err != nil {
// 		return nil, fmt.Errorf("无效的数字：%v", err)
// 	}

// 	avcDatas := agcvMain.DataStorage.GetDeviceData(1, number, AVC_DEVICE_TYPE, cons.YX)
// 	if avcDatas == nil {
// 		return nil, fmt.Errorf("未找到数据")
// 	}
// 	global.GVA_LOG.Info("AVC数据：", zap.Any("avcDatas", avcDatas))
// 	avcStatus := make(map[string]interface{})
// 	for point, data := range avcDatas {
// 		switch point {
// 		case AVC_YX_CONTROL_MODE:
// 			avcStatus["avcControlMode"] = data.Value
// 		case AVC_YX_SIGNAL:
// 			avcStatus["avcFunctionState"] = data.Value
// 		case AVC_YX_LOOP_STATUS:
// 			avcStatus["avcControlAuthority"] = data.Value
// 		}
// 	}

// 	return avcStatus, nil
// }

// GetBwdRealtimeData 获取并网点实时数据
func (agvcBwdHisService *AgvcBwdHisService) GetBwdRealtimeData(ctx context.Context, number string) (map[string][]map[string]interface{}, error) {
	// 检查InfluxDB客户端是否已初始化
	if global.GVA_INFLUXDB == nil {
		return nil, fmt.Errorf("InfluxDB客户端未初始化")
	}

	//开始时间今天00:00:00
	startTime := time.Now().Truncate(24 * time.Hour)
	// 结束时间今天23:59:59
	endTime := startTime.Add(24 * time.Hour).Add(-1 * time.Second)

	// 获取查询API
	queryAPI := global.GVA_INFLUXDB.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)

	// 构建Flux查询语句 - 获取最近1分钟的数据
	flux := fmt.Sprintf(`
    from(bucket: "%s")
        |> range(start: %s, stop: %s)
        |> filter(fn: (r) => r["_measurement"] == "agvc_data")
        |> filter(fn: (r) => r["eqid"] == "%s")
        |> filter(fn: (r) => r["psid"] == "1")
        |> filter(fn: (r) => r["dataType"] == "2")
        |> filter(fn: (r) => r["_field"] == "value")
    `, global.GVA_CONFIG.InfluxDB.Bucket, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339), number)

	// 执行查询
	result, err := queryAPI.Query(ctx, flux)
	if err != nil {
		return nil, fmt.Errorf("查询InfluxDB失败: %v", err)
	}

	// 解析结果
	totalData := make(map[string][]map[string]interface{}, 0)
	realtimeData := make([]map[string]interface{}, 0)
	dispatchData := make([]map[string]interface{}, 0)
	for result.Next() {
		record := result.Record()
		eqType := fmt.Sprintf("%v", record.ValueByKey("eqType"))
		point := fmt.Sprintf("%v", record.ValueByKey("point"))

		timestamp := record.Time().Truncate(time.Minute) // 向下取整到秒，确保同一秒的数据能正确聚合

		//获取实际出力
		if eqType != "5" {
			if point == BWD_YC_POWER_REAL {
				data := map[string]interface{}{
					"time":      timestamp,
					"psid":      record.ValueByKey("psid"),
					"eqid":      record.ValueByKey("eqid"),
					"eqType":    record.ValueByKey("eqType"),
					"dataType":  record.ValueByKey("dataType"),
					"point":     point,
					"pointName": agvc.GetBwdPointName(point), // 添加点位名称
					"value":     record.Value(),
				}
				realtimeData = append(realtimeData, data)
			}
		}
		//采集调度出力
		if eqType == "60" {
			if point == AGC_YC_POWER_EXEC_VALUE {
				data := map[string]interface{}{
					"time":      timestamp,
					"psid":      record.ValueByKey("psid"),
					"eqid":      record.ValueByKey("eqid"),
					"eqType":    record.ValueByKey("eqType"),
					"dataType":  record.ValueByKey("dataType"),
					"point":     point,
					"pointName": agvc.GetBwdPointName(point), // 添加点位名称
					"value":     record.Value(),
				}
				dispatchData = append(dispatchData, data)
			}
		}

	}
	totalData["realtimeData"] = realtimeData
	totalData["dispatchData"] = dispatchData
	return totalData, nil
}
