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
func (agvcBwdHisService *AgvcBwdHisService) UpdateAgvcBwdHis(ctx context.Context, agvcBwdHis agvc.AgvcBwdHis) (err error) {
	err = global.GVA_DB.Model(&agvc.AgvcBwdHis{}).Where("id = ?", agvcBwdHis.ID).Updates(&agvcBwdHis).Error
	return err
}

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

// GetAgvcBwdHistory 获取并网点历史数据
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
		|> filter(fn: (r) => r["_measurement"] == "agvc")
		|> filter(fn: (r) => r["eqid"] == "%s")
		|> filter(fn: (r) => r["eqType"] == "5")
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
