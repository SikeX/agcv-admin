package monitor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/monitor"
	monitorReq "github.com/flipped-aurora/gin-vue-admin/server/model/monitor/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/setting"
	serviceSetting "github.com/flipped-aurora/gin-vue-admin/server/service/setting"
)

type InverterMonitorService struct{}

var sysInverterSettingService serviceSetting.SysInverterSettingService

// CreateInverterMonitor 创建逆变器监控记录
// Author [yourname](https://github.com/yourname)
func (inverterMonitorService *InverterMonitorService) CreateInverterMonitor(ctx context.Context, inverterMonitor *monitor.InverterMonitor) (err error) {
	err = global.GVA_DB.Create(inverterMonitor).Error
	return err
}

// DeleteInverterMonitor 删除逆变器监控记录
// Author [yourname](https://github.com/yourname)
func (inverterMonitorService *InverterMonitorService) DeleteInverterMonitor(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&monitor.InverterMonitor{}, "id = ?", ID).Error
	return err
}

// DeleteInverterMonitorByIds 批量删除逆变器监控记录
// Author [yourname](https://github.com/yourname)
func (inverterMonitorService *InverterMonitorService) DeleteInverterMonitorByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]monitor.InverterMonitor{}, "id in ?", IDs).Error
	return err
}

// UpdateInverterMonitor 更新逆变器监控记录
// Author [yourname](https://github.com/yourname)
func (inverterMonitorService *InverterMonitorService) UpdateInverterMonitor(ctx context.Context, inverterMonitor monitor.InverterMonitor) (err error) {
	err = global.GVA_DB.Model(&monitor.InverterMonitor{}).Where("id = ?", inverterMonitor.ID).Updates(&inverterMonitor).Error
	return err
}

// GetInverterMonitor 根据ID获取逆变器监控记录
// Author [yourname](https://github.com/yourname)
func (inverterMonitorService *InverterMonitorService) GetInverterMonitor(ctx context.Context, ID string) (inverterMonitor monitor.InverterMonitor, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&inverterMonitor).Error
	return
}

// GetInverterMonitorInfoList 分页获取逆变器监控记录
// Author [yourname](https://github.com/yourname)
func (inverterMonitorService *InverterMonitorService) GetInverterMonitorInfoList(ctx context.Context, info monitorReq.InverterMonitorSearch) (list []monitor.InverterMonitor, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&monitor.InverterMonitor{})
	var inverterMonitors []monitor.InverterMonitor
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
	if info.Status != nil && *info.Status != "" {
		db = db.Where("status = ?", *info.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&inverterMonitors).Error
	return inverterMonitors, total, err
}
func (inverterMonitorService *InverterMonitorService) GetInverterMonitorPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}

// 获取逆变器历史数据
func (inverterMonitorService *InverterMonitorService) GetInverterHistory(ctx context.Context, inverterNo, startTime, endTime string) ([]monitor.InverterHistory, error) {
	// 检查InfluxDB客户端是否已初始化
	if global.GVA_INFLUXDB == nil {
		return nil, fmt.Errorf("InfluxDB客户端未初始化")
	}
	//如果没有传入结束时间，则使用当前时间
	if endTime == "" {
		endTime = "now()"
	}
	if startTime == "" {
		startTime = "now() - 1d"
	}

	sysInverterSetting, err := sysInverterSettingService.GetSysInverterSettingByInverterNo(ctx, inverterNo)
	if err != nil {
		return nil, fmt.Errorf("获取逆变器设置失败: %v", err)
	}

	pointMap, codes := inverterMonitorService.getPointMapAndCodes(sysInverterSetting, inverterNo)

	// 获取查询API
	queryAPI := global.GVA_INFLUXDB.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)

	queryStr := ""
	for i, no := range codes {
		if i > 0 {
			queryStr += " or "
		}
		queryStr += fmt.Sprintf(`r["code"] == "%s"`, no)
	}
	// 构建Flux查询语句
	// 注意：这里假设数据存储在sensor_data measurement中，tag为code，field包含value等字段
	flux := fmt.Sprintf(`
	from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r["_measurement"] == "mqtt")
  	|> filter(fn: (r) => %s)
	`, global.GVA_CONFIG.InfluxDB.Bucket, startTime, endTime, queryStr)

	fmt.Println(flux)

	// 执行查询
	result, err := queryAPI.Query(ctx, flux)
	if err != nil {
		return nil, fmt.Errorf("查询InfluxDB失败: %v", err)
	}
	if !result.Next() {
		return nil, nil
	}

	// 解析结果
	// 使用map来存储同一时间点的不同code值
	// historyMap := make(map[time.Time]monitor.InverterHistory)
	timeCodeValueMap := make(map[time.Time]map[string]float64)

	for result.Next() {
		record := result.Record()
		//获取code
		code := record.ValueByKey("code").(string)

		// 获取时间并作为map的key
		timestamp := record.Time().Truncate(time.Minute) // 向下取整到秒，确保同一秒的数据能正确聚合

		// 从map中获取已有的history，如果不存在则创建新的
		codeValueMap, exists := timeCodeValueMap[timestamp]
		if !exists {
			codeValueMap = make(map[string]float64)
			timeCodeValueMap[timestamp] = codeValueMap
		}

		// 根据pointMap中code对应的值设置history中的相应字段
		// fieldName := pointMap[code]
		var value float64
		//判断record.Value()的类型
		switch v := record.Value().(type) {
		case int64:
			value = float64(v)
		case float64:
			value = v
		}

		codeValueMap[code] = value

		// // 根据字段名设置对应的值
		// switch fieldName {
		// case "ygPreal":
		// 	history.ActivePower = value
		// case "wgPreal":
		// 	history.PowerFactor = value
		// }

		// 更新map中的值
		// historyMap[timestamp] = history
		timeCodeValueMap[timestamp] = codeValueMap
	}

	var historys []monitor.InverterHistory
	//处理timeCodeValueMap
	for time, codeValueMap := range timeCodeValueMap {
		history := monitor.InverterHistory{
			Time: time,
		}
		for code, value := range codeValueMap {
			switch pointMap[code] {
			case "ygPreal":
				history.YgPreal += value
			case "wgPreal":
				history.WgPreal += value
			case "ygPmax":
				history.YgPmax += value
			case "ygPMin":
				history.YgPMin += value
			case "wgPmax":
				history.WgPmax += value
			case "wgPMin":
				history.WgPMin += value
			case "pf":
				history.Pf += value
			}
		}
		historys = append(historys, history)
	}

	// 按时间排序
	// sort.Slice(historys, func(i, j int) bool {
	// 	return historys[i].Time.Before(historys[j].Time)
	// })

	// 检查是否有错误
	if result.Err() != nil {
		return nil, fmt.Errorf("解析InfluxDB结果失败: %v", result.Err())
	}

	return historys, nil
}

func (inverterMonitorService *InverterMonitorService) getPointMapAndCodes(sysInverterSetting setting.SysInverterSetting, inverterNo string) (map[string]string, []string) {
	var codes []string
	pointMap := make(map[string]string)
	//有功功率点号
	if sysInverterSetting.YgPreal != nil {
		ygPrealPoints := strings.Split(*sysInverterSetting.YgPreal, ",")
		for _, point := range ygPrealPoints {
			code := inverterNo + point
			codes = append(codes, code)
			pointMap[code] = "ygPreal"
		}
	}
	//无功功率点号
	if sysInverterSetting.WgPreal != nil {
		wgPrealPoints := strings.Split(*sysInverterSetting.WgPreal, ",")
		for _, point := range wgPrealPoints {
			code := inverterNo + point
			codes = append(codes, code)
			pointMap[code] = "wgPreal"
		}
	}
	//有功功率最大值点号
	if sysInverterSetting.YgPmax != nil {
		ygPmaxPoints := strings.Split(*sysInverterSetting.YgPmax, ",")
		for _, point := range ygPmaxPoints {
			code := inverterNo + point
			codes = append(codes, code)
			pointMap[code] = "ygPmax"
		}
	}
	//有功功率最小值点号
	if sysInverterSetting.YgPMin != nil {
		ygPMinPoints := strings.Split(*sysInverterSetting.YgPMin, ",")
		for _, point := range ygPMinPoints {
			code := inverterNo + point
			codes = append(codes, code)
			pointMap[code] = "ygPMin"
		}
	}
	//无功功率最大值点号
	if sysInverterSetting.WgPmax != nil {
		wgPmaxPoints := strings.Split(*sysInverterSetting.WgPmax, ",")
		for _, point := range wgPmaxPoints {
			code := inverterNo + point
			codes = append(codes, code)
			pointMap[code] = "wgPmax"
		}
	}
	return pointMap, codes
}
