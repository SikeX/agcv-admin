
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

type AgvcQxyHisService struct {}
// CreateAgvcQxyHis 创建气象仪监控记录
// Author [yourname](https://github.com/yourname)
func (agvcQxyHisService *AgvcQxyHisService) CreateAgvcQxyHis(ctx context.Context, agvcQxyHis *agvc.AgvcQxyHis) (err error) {
	err = global.GVA_DB.Create(agvcQxyHis).Error
	return err
}

// DeleteAgvcQxyHis 删除气象仪监控记录
// Author [yourname](https://github.com/yourname)
func (agvcQxyHisService *AgvcQxyHisService)DeleteAgvcQxyHis(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&agvc.AgvcQxyHis{},"id = ?",ID).Error
	return err
}

// DeleteAgvcQxyHisByIds 批量删除气象仪监控记录
// Author [yourname](https://github.com/yourname)
func (agvcQxyHisService *AgvcQxyHisService)DeleteAgvcQxyHisByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]agvc.AgvcQxyHis{},"id in ?",IDs).Error
	return err
}

// UpdateAgvcQxyHis 更新气象仪监控记录
// Author [yourname](https://github.com/yourname)
func (agvcQxyHisService *AgvcQxyHisService)UpdateAgvcQxyHis(ctx context.Context, agvcQxyHis agvc.AgvcQxyHis) (err error) {
	err = global.GVA_DB.Model(&agvc.AgvcQxyHis{}).Where("id = ?",agvcQxyHis.ID).Updates(&agvcQxyHis).Error
	return err
}

// GetAgvcQxyHis 根据ID获取气象仪监控记录
// Author [yourname](https://github.com/yourname)
func (agvcQxyHisService *AgvcQxyHisService)GetAgvcQxyHis(ctx context.Context, ID string) (agvcQxyHis agvc.AgvcQxyHis, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&agvcQxyHis).Error
	return
}
// GetAgvcQxyHisInfoList 分页获取气象仪监控列表
// 从 InfluxDB中查询设备列表及其最新数据
func (agvcQxyHisService *AgvcQxyHisService)GetAgvcQxyHisInfoList(ctx context.Context, info agvcReq.AgvcQxyHisSearch) (list []agvc.AgvcQxyHis, total int64, err error) {
	// 检查InfluxDB客户端是否已初始化
	if global.GVA_INFLUXDB == nil {
		return nil, 0, fmt.Errorf("InfluxDB客户端未初始化")
	}

	// 获取查询API
	queryAPI := global.GVA_INFLUXDB.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)

	// 构建Flux查询语句 - 获取所有气象仪设备的最新数据
	// 使用group by eqid来获取每个设备的最新记录
	flux := fmt.Sprintf(`
	from(bucket: "%s")
		|> range(start: -30d)
		|> filter(fn: (r) => r["_measurement"] == "agvc")
		|> filter(fn: (r) => r["eqType"] == "8")
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
	var qxySettings []agvc.AgvcQxySetting
	global.GVA_DB.Find(&qxySettings)

	// 创建设备编号到名称的映射
	numberToNameMap := make(map[string]string)
	for _, setting := range qxySettings {
		if setting.Number != nil && setting.DeviceName != nil {
			numberToNameMap[*setting.Number] = *setting.DeviceName
		}
	}

	// 转换为AgvcQxyHis结构体列表
	var agvcQxyHisList []agvc.AgvcQxyHis
	for eqid := range deviceSet {
		number := eqid
		name := eqid // 默认使用设备编号
		
		// 如果在配置表中找到名称,则使用配置的名称
		if deviceName, ok := numberToNameMap[eqid]; ok {
			name = deviceName
		}

		qxyHis := agvc.AgvcQxyHis{
			Number: &number,
			Name:   &name,
		}
		agvcQxyHisList = append(agvcQxyHisList, qxyHis)
	}

	// 计算总数
	total = int64(len(agvcQxyHisList))

	// 应用分页
	offset := info.PageSize * (info.Page - 1)
	limit := info.PageSize

	if offset >= int(total) {
		return []agvc.AgvcQxyHis{}, total, nil
	}

	end := offset + limit
	if end > int(total) {
		end = int(total)
	}

	if limit != 0 {
		return agvcQxyHisList[offset:end], total, nil
	}

	return agvcQxyHisList, total, nil
}
func (agvcQxyHisService *AgvcQxyHisService)GetAgvcQxyHisPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}

// GetAgvcQxyHistory 获取气象仪历史数据
func (agvcQxyHisService *AgvcQxyHisService) GetAgvcQxyHistory(ctx context.Context, eqid, startTime, endTime string) ([]map[string]interface{}, error) {
	// 检查InfluxDB客户端是否已初始化
	if global.GVA_INFLUXDB == nil {
		return nil, fmt.Errorf("InfluxDB客户端未初始化")
	}

	// 获取查询API
	queryAPI := global.GVA_INFLUXDB.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)

	// 构建Flux查询语句 - 使用统一的agvc表
	// eqType=8 代表气象仪
	flux := fmt.Sprintf(`
	from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r["_measurement"] == "agvc")
		|> filter(fn: (r) => r["eqid"] == "%s")
		|> filter(fn: (r) => r["eqType"] == "8")
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
			"pointName": agvc.GetQxyPointName(point), // 添加点位名称
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

// GenerateAgvcQxyTestData 生成气象仪测试数据
// 随机选择10-15个点位,每个点位生成10-20条时间序列数据
func (agvcQxyHisService *AgvcQxyHisService) GenerateAgvcQxyTestData(ctx context.Context, eqid string, count int) error {
	// 检查InfluxDB客户端是否已初始化
	if global.GVA_INFLUXDB == nil {
		return fmt.Errorf("InfluxDB客户端未初始化")
	}

	// 获取写API
	writeAPI := global.GVA_INFLUXDB.WriteAPIBlocking(global.GVA_CONFIG.InfluxDB.Org, global.GVA_CONFIG.InfluxDB.Bucket)

	// 获取所有气象仪点位
	allPoints := agvc.GetAllQxyPoints()

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
			case "2": // 环境温度(-20~50℃)
				value = 25 + 15*math.Sin(float64(i)/10) + rand.Float64()*5 - 2.5
			case "46": // 环境湿度(20~100%RH)
				value = 60 + 20*math.Sin(float64(i)/8) + rand.Float64()*10
			case "4", "5", "6": // 风速(0~15m/s)
				value = 5 + 3*math.Sin(float64(i)/12) + rand.Float64()*2
			case "7": // 风向(0~360°)
				value = 180 + 90*math.Sin(float64(i)/15) + rand.Float64()*30
			case "11": // 组件温度(-10~70℃)
				value = 30 + 20*math.Sin(float64(i)/10) + rand.Float64()*5
			case "8", "9": // 辐射强度(0~1300W/m²)
				value = 600 + 400*math.Sin(float64(i)/10) + rand.Float64()*100
			case "39", "40": // 总辐射量(0~30MJ/m²)
				value = 15 + 10*math.Sin(float64(i)/12) + rand.Float64()*3
			case "52": // 气压(900~1100Pa)
				value = 1000 + 50*math.Sin(float64(i)/20) + rand.Float64()*20
			case "12": // 雨量
				value = rand.Float64() * 10
			case "3": // 露点温度
				value = 15 + 10*math.Sin(float64(i)/10) + rand.Float64()*5
			case "51": // 日照小时数(1~10h)
				value = 5 + 3*math.Sin(float64(i)/8) + rand.Float64()*2
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
					"eqType":   agvc.EqTypeQXY, // 使用常量
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
