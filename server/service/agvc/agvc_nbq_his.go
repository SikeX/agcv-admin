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
	"github.com/flipped-aurora/gin-vue-admin/server/service/agvc/agcv_main"
	"github.com/flipped-aurora/gin-vue-admin/server/service/agvc/cons"
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

type codePointKey struct {
	Code int
	Attr string
}

// GetAgvcNbqHisInfoList 分页获取逆变器历史数据列表
// 从InfluxDB中查询设备列表及其最新数据
func (agvcNbqHisService *AgvcNbqHisService) GetAgvcNbqHisInfoList(ctx context.Context, info agvcReq.AgvcNbqHisSearch) (list []agvc.AgvcNbqHis, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	// 创建db
	db := global.GVA_DB.Model(&agvc.AgvcNbqSetting{})
	var inverterSettings []agvc.AgvcNbqSetting
	// 如果有条件搜索 下方会自动创建搜索语句
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
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

	//逆变器监控要查询的属性
	attrs := []string{
		"apparentPower", //有功功率
		"reactivePower", //无功功率
		"powerFactor",   // 功率因数
	}

	//提取出inverterSettings的inverterNo
	//code: 逆变器编号
	//attr: 逆变器属性
	codeAttrValueMap := make(map[int]map[string]float64)
	nbqPoints := make([]codePointKey, 0)
	for _, setting := range inverterSettings {
		for _, attr := range attrs {
			key := codePointKey{
				Code: *setting.InverterNo,
				Attr: attr,
			}
			nbqPoints = append(nbqPoints, key)
			// inverterNos = append(inverterNos, *setting.InverterNo+cons.NBLabelPointMap[attr])
			// codeAttrValueMap[*setting.InverterNo][attr] = 0
			codeAttrValueMap[*setting.InverterNo] = map[string]float64{}
		}
	}

	for _, codePoint := range nbqPoints {
		//从内存中获取数据
		if value, err := agcv_main.DataStorage.GetDataAsFloat64(codePoint.Code, cons.TYPE_NBQ, cons.YC, cons.NBLabelPointMap[codePoint.Attr]); err == nil {
			codeAttrValueMap[codePoint.Code][codePoint.Attr] = value
		}
	}

	agvcNbqHises := make([]agvc.AgvcNbqHis, 0)
	status := "1"
	for _, setting := range inverterSettings {
		apparentPower := 0.0
		reactivePower := 0.0
		powerFactor := 0.0
		if attrValue, ok := codeAttrValueMap[*setting.InverterNo]["apparentPower"]; ok {
			apparentPower = attrValue
		}
		if attrValue, ok := codeAttrValueMap[*setting.InverterNo]["reactivePower"]; ok {
			reactivePower = attrValue
		}
		if attrValue, ok := codeAttrValueMap[*setting.InverterNo]["powerFactor"]; ok {
			powerFactor = attrValue
		}
		agvcNbqHises = append(agvcNbqHises, agvc.AgvcNbqHis{
			InverterNo:          setting.InverterNo,
			Name:                setting.Name,
			Status:              &status,
			IsParticipateAdjust: setting.IsParticipateAdjust,
			RatedPower:          setting.RatedActivePower,
			ActivePower:         &apparentPower,
			ReactivePower:       &reactivePower,
			PowerFactor:         &powerFactor,
		})
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
