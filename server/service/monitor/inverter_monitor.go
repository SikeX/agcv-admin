package monitor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/monitor"
	monitorReq "github.com/flipped-aurora/gin-vue-admin/server/model/monitor/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service/setting"
)

type InverterMonitorService struct{}

var sysInverterSettingService setting.SysInverterSettingService

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
		db = db.Where("inverter_no = ?", *info.Inverter_no)
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

type InverterPointSet struct {
	Time        time.Time
	PowerFactor float64
}

// GetInverterHistory 从InfluxDB获取逆变器历史数据
func (inverterMonitorService *InverterMonitorService) GetInverterHistory(ctx context.Context, inverterNo string, startTime, endTime string) ([]monitor.InverterHistory, error) {
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

	//有功功率点号
	ygPrealPoints := strings.Split(*sysInverterSetting.YgPreal, ",")
	wgPrealPonts := strings.Split(*sysInverterSetting.WgPreal, ",")

	inverterNos := []string{"00100010202326", "00100010202329"}

	// 获取查询API
	queryAPI := global.GVA_INFLUXDB.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)

	queryStr := ""
	for i, no := range inverterNos {
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
	fmt.Println(result.Record())

	// 解析结果
	// var records []map[string]interface{}
	var historys []monitor.InverterHistory
	for result.Next() {
		history := monitor.InverterHistory{}
		record := result.Record()
		//获取code
		code := record.ValueByKey("code").(string)
		if code == "00100010202326" {
		}
		// records = append(records, data)
		//时间转换为时间戳
		history.Time = record.Time()
		//判断record.Value()的类型
		switch v := record.Value().(type) {
		case int64:
			history.PowerFactor = float64(v)
		case float64:
			history.PowerFactor = v
		}
		historys = append(historys, history)
	}

	fmt.Println(historys)

	// 检查是否有错误
	if result.Err() != nil {
		return nil, fmt.Errorf("解析InfluxDB结果失败: %v", result.Err())
	}

	return historys, nil
}
