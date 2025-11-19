package agvc

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
	agvcReq "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service/agvc/cons"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type AgvcNbqHisService struct{}

var agvcHisService = AgvcHisService{}

func init() {
	// 系统启动时初始化point标签映射
	agvc.InitPointMapping()
}

// CreateAgvcNbqHis 创建agvcNbqHis表记录
// Author [yourname](https://github.com/yourname)
func (agvcNbqHisService *AgvcNbqHisService) CreateAgvcNbqHis(ctx context.Context, agvcNbqHis *agvc.AgvcNbq) (err error) {
	err = global.GVA_DB.Create(agvcNbqHis).Error
	return err
}

// DeleteAgvcNbqHis 删除agvcNbqHis表记录
// Author [yourname](https://github.com/yourname)
func (agvcNbqHisService *AgvcNbqHisService) DeleteAgvcNbqHis(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&agvc.AgvcNbq{}, "id = ?", ID).Error
	return err
}

// DeleteAgvcNbqHisByIds 批量删除agvcNbqHis表记录
// Author [yourname](https://github.com/yourname)
func (agvcNbqHisService *AgvcNbqHisService) DeleteAgvcNbqHisByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]agvc.AgvcNbq{}, "id in ?", IDs).Error
	return err
}

// UpdateAgvcNbqHis 更新agvcNbqHis表记录
// Author [yourname](https://github.com/yourname)
func (agvcNbqHisService *AgvcNbqHisService) UpdateAgvcNbqHis(ctx context.Context, agvcNbqHis agvc.AgvcNbq) (err error) {
	err = global.GVA_DB.Model(&agvc.AgvcNbq{}).Where("id = ?", agvcNbqHis.InverterNo).Updates(&agvcNbqHis).Error
	return err
}

// GetAgvcNbqHis 根据ID获取agvcNbqHis表记录
// Author [yourname](https://github.com/yourname)
func (agvcNbqHisService *AgvcNbqHisService) GetAgvcNbqHis(ctx context.Context, ID string) (agvcNbqHis agvc.AgvcNbq, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&agvcNbqHis).Error
	return
}

// GetAgvcNbqHisInfoList 分页获取逆变器历史数据列表
// 从InfluxDB中查询设备列表及其最新数据
func (agvcNbqHisService *AgvcNbqHisService) GetAgvcNbqHisInfoList(ctx context.Context, info agvcReq.AgvcNbqHisSearch) (list []agvc.AgvcNbq, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	// 创建db，从AgvcNbqSetting获取逆变器配置列表
	db := global.GVA_DB.Model(&agvc.AgvcNbqSetting{})
	var inverterSettings []agvc.AgvcNbqSetting

	if info.Psid != nil {
		db = db.Where("psid = ?", *info.Psid)
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
	if err != nil {
		return
	}

	agvcNbqHises := make([]agvc.AgvcNbq, 0)

	// 对每个逆变器查询其最新数据
	for _, setting := range inverterSettings {
		agvcNbqHisInterface, err := agvcHisService.GetHistoryGeneric(ctx, *setting.Psid, *setting.InverterNo, agvc.EqTypeNBQ, nil, *setting.Name, info.StartTime, info.EndTime)
		if err != nil {
			global.GVA_LOG.Error(fmt.Sprintf("查询逆变器%d的历史数据失败: %v", *setting.InverterNo, err))
			continue
		}
		agvcNbqHis, ok := agvcNbqHisInterface.([]agvc.AgvcNbq)
		if !ok {
			global.GVA_LOG.Error(fmt.Sprintf("类型断言失败: %v", agvcNbqHisInterface))
			continue
		}
		agvcNbqHises = append(agvcNbqHises, agvcNbqHis...)
	}

	// 对inverterMonitors手动分页
	// 计算总数
	total = int64(len(agvcNbqHises))

	// 应用分页
	offset = info.PageSize * (info.Page - 1)
	limit = info.PageSize

	if offset >= int(total) {
		return []agvc.AgvcNbq{}, total, nil
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

// setFieldByPoint 根据point值设置AgvcNbqHis结构体对应的字段
func setFieldByPoint(nbqHis *agvc.AgvcNbq, point string, value float64) {
	// 使用反射获取结构体
	v := reflect.ValueOf(nbqHis).Elem()
	t := v.Type()

	// 遍历结构体字段，找到匹配的point标签
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("point")

		if tag == "" {
			continue
		}

		// 解析point标签，提取value值
		if strings.Contains(tag, "value:") {
			parts := strings.Split(tag, ",")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if strings.HasPrefix(part, "value:") {
					pointValue := strings.TrimPrefix(part, "value:")
					if pointValue == point {
						// 找到匹配的字段，设置值
						if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.Float64 {
							field.Set(reflect.ValueOf(&value))
							return
						}
					}
				}
			}
		}
	}
}

func (agvcNbqHisService *AgvcNbqHisService) GetAgvcNbqHisPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}

// GetAgvcNbqHistory 获取逆变器历史数据
// 调用通用GetHistory方法，传入设备类型
// func (agvcNbqHisService *AgvcNbqHisService) GetAgvcNbqHistory(ctx context.Context, nbqHis agvc.AgvcNbq, eqName, startTime, endTime string) ([]agvc.AgvcNbq, error) {
// 	// 检查必要的参数
// 	if nbqHis.InverterNo == nil {
// 		return nil, fmt.Errorf("逆变器编号不能为空")
// 	}
// 	if nbqHis.Psid == nil {
// 		return nil, fmt.Errorf("电站编号不能为空")
// 	}

// 	// 提取结构体中所有带point标签的字段的value值
// 	pointValues := extractPointValues(nbqHis)

// 	// 如果没有提取到任何point值，返回错误
// 	if len(pointValues) == 0 {
// 		return nil, fmt.Errorf("未找到需要查询的点位信息")
// 	}

// 	// 调用通用历史数据服务
// 	hisService := &AgvcHisService{}

// 	return hisService.GetHistory(ctx, *nbqHis.Psid, *nbqHis.InverterNo, agvc.EqTypeNBQ, pointValues, eqName, startTime, endTime)
// }

// setFieldValue 根据字段名设置AgvcNbqHis结构体对应字段的值
// func setFieldValue(nbqHis *agvc.AgvcNbqHis, fieldName string, value float64) {
//     v := reflect.ValueOf(nbqHis).Elem()
//     field := v.FieldByName(fieldName)

//     if !field.IsValid() || !field.CanSet() {
//         return
//     }

//     // 设置指针类型的float64值
//     if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.Float64 {
//         field.Set(reflect.ValueOf(&value))
//     }
// }

// extractPointValues 从AgvcNbqHis结构体中提取所有带point标签的字段的value值
func extractPointValues(nbqHis agvc.AgvcNbq) []string {
	var pointValues []string

	// 使用反射获取结构体
	v := reflect.ValueOf(nbqHis)
	t := v.Type()

	// 遍历结构体字段，找到所有带point标签的字段
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("point")

		if tag == "" {
			continue
		}

		// 检查字段值是否为nil（因为所有字段都是指针类型）
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

// GetNbqRealData 获取逆变器实时数据（从InfluxDB中获取最后一条数据）
func (agvcNbqHisService *AgvcNbqHisService) GetNbqRealData(ctx context.Context, psid, inverterNo int) (*agvc.AgvcNbq, error) {
	// 调用通用实时数据服务
	hisService := &AgvcHisService{}
	eqidStr := fmt.Sprintf("%d", inverterNo)
	return hisService.GetRealData(ctx, psid, eqidStr, agvc.EqTypeNBQ)
}

// WriteAgvcNbqDataToInfluxDB 将逆变器数据写入InfluxDB
// 该方法用于将CoAP或其他来源的数据写入InfluxDB
func (agvcNbqHisService *AgvcNbqHisService) WriteAgvcNbqDataToInfluxDB(ctx context.Context, psid int, eqid int, point string, value float64) error {
	// 检查InfluxDB客户端是否已初始化
	if global.GVA_INFLUXDB == nil {
		return fmt.Errorf("InfluxDB客户端未初始化")
	}

	// 获取写API
	writeAPI := global.GVA_INFLUXDB.WriteAPIBlocking(global.GVA_CONFIG.InfluxDB.Org, global.GVA_CONFIG.InfluxDB.Bucket)

	// 创建数据点
	p := influxdb2.NewPoint(
		global.GVA_CONFIG.InfluxDB.GetMeasurement(),
		map[string]string{
			"psid":     strconv.Itoa(psid),
			"eqType":   agvc.EqTypeNBQ,
			"eqid":     strconv.Itoa(eqid),
			"dataType": strconv.Itoa(cons.YC),
			"point":    point,
		},
		map[string]interface{}{
			"value": value,
		},
		time.Now(),
	)

	// 写入数据点
	if err := writeAPI.WritePoint(ctx, p); err != nil {
		return fmt.Errorf("写入InfluxDB失败: %v", err)
	}

	return nil
}
