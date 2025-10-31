package agvc

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
	agvcReq "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/request"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type AgvcNbqHisService struct{}

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
	err = global.GVA_DB.Model(&agvc.AgvcNbqHis{}).Where("id = ?", agvcNbqHis.ID).Updates(&agvcNbqHis).Error
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
	// 检查InfluxDB客户端是否已初始化
	if global.GVA_INFLUXDB == nil {
		return nil, 0, fmt.Errorf("InfluxDB客户端未初始化")
	}

	// 获取查询API
	queryAPI := global.GVA_INFLUXDB.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)

	// 构建Flux查询语句 - 获取所有逆变器设备的最新数据
	// 使用group by eqid来获取每个设备的最新记录
	flux := fmt.Sprintf(`
	from(bucket: "%s")
		|> range(start: -30d)
		|> filter(fn: (r) => r["_measurement"] == "agvc")
		|> filter(fn: (r) => r["eqType"] == "2")
		|> group(columns: ["eqid"])
		|> last()
	`, global.GVA_CONFIG.InfluxDB.Bucket)

	// 执行查询
	result, err := queryAPI.Query(ctx, flux)
	if err != nil {
		return nil, 0, fmt.Errorf("查询InfluxDB失败: %v", err)
	}

	// 按设备ID分组收集数据
	deviceDataMap := make(map[string]map[string]interface{})

	for result.Next() {
		record := result.Record()
		eqid := fmt.Sprintf("%v", record.ValueByKey("eqid"))
		point := fmt.Sprintf("%v", record.ValueByKey("point"))
		value := record.Value()

		if _, exists := deviceDataMap[eqid]; !exists {
			deviceDataMap[eqid] = make(map[string]interface{})
			deviceDataMap[eqid]["eqid"] = eqid
		}

		// 根据点位映射存储数据
		switch point {
		case "10": // 交流功率(有功功率)
			if v, ok := value.(float64); ok {
				deviceDataMap[eqid]["activePower"] = v
			}
		case "27": // 无功功率
			if v, ok := value.(float64); ok {
				deviceDataMap[eqid]["reactivePower"] = v
			}
		case "21": // 功率因数
			if v, ok := value.(float64); ok {
				deviceDataMap[eqid]["powerFactor"] = v
			}
		case "202": // 直流功率(额定功率)
			if v, ok := value.(float64); ok {
				deviceDataMap[eqid]["ratedPower"] = v
			}
		}
	}

	// 检查查询错误
	if result.Err() != nil {
		return nil, 0, fmt.Errorf("解析InfluxDB结果失败: %v", result.Err())
	}

	// 转换为AgvcNbqHis结构体列表
	var agvcNbqHisList []agvc.AgvcNbqHis
	for eqid, data := range deviceDataMap {
		inverterNo, _ := data["eqid"].(string)
		var inverterNoInt int64
		fmt.Sscanf(inverterNo, "%d", &inverterNoInt)

		activePower := 0.0
		if v, ok := data["activePower"].(float64); ok {
			activePower = v
		}

		reactivePower := 0.0
		if v, ok := data["reactivePower"].(float64); ok {
			reactivePower = v
		}

		powerFactor := 0.0
		if v, ok := data["powerFactor"].(float64); ok {
			powerFactor = v
		}

		ratedPower := 0.0
		if v, ok := data["ratedPower"].(float64); ok {
			ratedPower = v
		}

		nbqHis := agvc.AgvcNbqHis{
			InverterNo:    &inverterNoInt,
			Name:          &eqid,
			ActivePower:   &activePower,
			ReactivePower: &reactivePower,
			PowerFactor:   &powerFactor,
			RatedPower:    &ratedPower,
		}
		agvcNbqHisList = append(agvcNbqHisList, nbqHis)
	}

	// 计算总数
	total = int64(len(agvcNbqHisList))

	// 应用分页
	offset := info.PageSize * (info.Page - 1)
	limit := info.PageSize

	if offset >= int(total) {
		return []agvc.AgvcNbqHis{}, total, nil
	}

	end := offset + limit
	if end > int(total) {
		end = int(total)
	}

	if limit != 0 {
		return agvcNbqHisList[offset:end], total, nil
	}

	return agvcNbqHisList, total, nil
}
func (agvcNbqHisService *AgvcNbqHisService) GetAgvcNbqHisPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}

// GetAgvcNbqHistory 获取逆变器历史数据
func (agvcNbqHisService *AgvcNbqHisService) GetAgvcNbqHistory(ctx context.Context, eqid, startTime, endTime string) ([]map[string]interface{}, error) {
	// 检查InfluxDB客户端是否已初始化
	if global.GVA_INFLUXDB == nil {
		return nil, fmt.Errorf("InfluxDB客户端未初始化")
	}

	// 获取查询API
	queryAPI := global.GVA_INFLUXDB.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)

	// 构建Flux查询语句 - 使用统一的agvc表
	// eqType=2 代表逆变器
	flux := fmt.Sprintf(`
	from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r["_measurement"] == "agvc")
		|> filter(fn: (r) => r["eqid"] == "%s")
		|> filter(fn: (r) => r["eqType"] == "2")
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

		data := map[string]interface{}{
			"time":      record.Time().Format(time.RFC3339),
			"psid":      record.ValueByKey("psid"),
			"eqid":      record.ValueByKey("eqid"),
			"eqType":    record.ValueByKey("eqType"),
			"dataType":  record.ValueByKey("dataType"),
			"point":     point,
			"pointName": agvc.GetNbqPointName(point), // 添加点位名称
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

// GenerateAgvcNbqTestData 生成逆变器测试数据
// 随机选择10-15个点位,每个点位生成10-20条时间序列数据
func (agvcNbqHisService *AgvcNbqHisService) GenerateAgvcNbqTestData(ctx context.Context, eqid string, count int) error {
	// 检查InfluxDB客户端是否已初始化
	if global.GVA_INFLUXDB == nil {
		return fmt.Errorf("InfluxDB客户端未初始化")
	}

	// 获取写API
	writeAPI := global.GVA_INFLUXDB.WriteAPIBlocking(global.GVA_CONFIG.InfluxDB.Org, global.GVA_CONFIG.InfluxDB.Bucket)

	// 获取所有逆变器点位
	allPoints := agvc.GetAllNbqPoints()

	// 随机选择10-15个点位
	rand.Seed(time.Now().UnixNano())
	pointCount := 10 + rand.Intn(6) // 10-15个点位
	selectedPoints := make([]string, 0, pointCount)

	// 随机打乱点位顺序
	rand.Shuffle(len(allPoints), func(i, j int) {
		allPoints[i], allPoints[j] = allPoints[j], allPoints[i]
	})

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
			case "1", "2", "3", "4": // 发电量
				value = 1000 + 500*math.Sin(float64(i)/5) + rand.Float64()*100
			case "10", "202": // 功率
				value = 100 + 50*math.Sin(float64(i)/10) + rand.Float64()*10 - 5
			case "94", "95", "96": // 相电压
				value = 220 + 10*math.Sin(float64(i)/8) + rand.Float64()*5
			case "12", "13", "14": // 相电流
				value = 50 + 20*math.Sin(float64(i)/6) + rand.Float64()*5
			case "22": // 温度
				value = 45 + 5*math.Sin(float64(i)/15) + rand.Float64()*2
			case "21", "410", "411", "412": // 功率因数
				value = 0.95 + 0.04*math.Sin(float64(i)/12) + rand.Float64()*0.02 - 0.01
			default:
				// 其他点位使用通用范围
				value = 50 + 30*math.Sin(float64(i)/8) + rand.Float64()*10 - 5
			}

			// 创建数据点
			p := influxdb2.NewPoint(
				"agvc",
				map[string]string{
					"psid":     "1",
					"eqid":     eqid,
					"eqType":   agvc.EqTypeNBQ, // 使用常量
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
