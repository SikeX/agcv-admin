// 自动生成模板AgvcAgcHis
package agvc

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// agvcAgcHis表 结构体  AgvcAgc
// 用于存储和展示AGC实时数据和历史数据
type AgvcAgc struct {
	Ctime               *string  `json:"ctime" form:"ctime" point:"name:采集时间"`
	Psid                *int     `json:"psid" form:"psid" point:"name:电站编号"`
	Eqid                *int     `json:"eqid" form:"eqid" point:"name:设备编号"`
	Name                *string  `json:"name" form:"name" point:"name:并网点名称"`
	TargetActivePower   *float64 `json:"targetActivePower" form:"targetActivePower" point:"name:目标有功(kW),value:403,type:2"`
	DispatchActivePower *float64 `json:"dispatchActivePower" form:"dispatchActivePower" point:"name:调度下发(kW),value:21,type:2"`
	AgcFunctionState    *float64 `json:"agcFunctionState" form:"agcFunctionState" point:"name:AGC功能投退状态,value:401,type:1"`
	AgcControlMode      *float64 `json:"agcControlMode" form:"agcControlMode" point:"name:AGC调节方式,value:404,type:1"`
	AgcControlAuthority *float64 `json:"agcControlAuthority" form:"agcControlAuthority" point:"name:AGC控制权限,value:402,type:1"`
	AgcUpperLimit       *float64 `json:"agcUpperLimit" form:"agcUpperLimit" point:"name:可调上限(kW),value:401,type:2"`
	AgcLowerLimit       *float64 `json:"agcLowerLimit" form:"agcLowerLimit" point:"name:可调下限(kW),value:402,type:2"`
}

// TableName agvcAgcHis表 AgvcAgcHis自定义表名 agvc_agc_his
func (AgvcAgc) TableName() string {
	return "agvc_agc_his"
}

// PointInfoAgc AGC的point信息结构体
type PointInfoAgc struct {
	FieldName string
	DataType  int // 1=遥信(YX), 2=遥测(YC)
}

// PointToFieldMappingAgc AGC的point值到字段信息的映射缓存
// key格式为 "type-value"，例如 "1-22" 或 "2-20"
var pointToFieldMappingAgc map[string]PointInfoAgc
var pointMappingOnceAgc sync.Once

// InitPointMappingAgc 初始化AGC的point标签映射，在系统启动时调用
func InitPointMappingAgc() {
	pointMappingOnceAgc.Do(func() {
		pointToFieldMappingAgc = make(map[string]PointInfoAgc)

		// 使用反射解析AgvcAgc结构体的point标签
		t := reflect.TypeOf(AgvcAgc{})
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tag := field.Tag.Get("point")

			if tag == "" {
				continue
			}

			// 解析point标签，提取value和type值
			if strings.Contains(tag, "value:") {
				parts := strings.Split(tag, ",")
				var pointValue string
				var dataType int = 2 // 默认为遥测

				for _, part := range parts {
					part = strings.TrimSpace(part)
					if strings.HasPrefix(part, "value:") {
						pointValue = strings.TrimPrefix(part, "value:")
					} else if strings.HasPrefix(part, "type:") {
						typeStr := strings.TrimPrefix(part, "type:")
						if typeStr == "1" {
							dataType = 1
						} else if typeStr == "2" {
							dataType = 2
						}
					}
				}

				if pointValue != "" {
					// 使用 "type-value" 作为复合key
					key := fmt.Sprintf("%d-%s", dataType, pointValue)
					pointToFieldMappingAgc[key] = PointInfoAgc{
						FieldName: field.Name,
						DataType:  dataType,
					}
				}
			}
		}
	})
}

// GetFieldNameByPointAgc 根据dataType和point值获取AGC字段名
func GetFieldNameByPointAgc(dataType int, point string) (string, bool) {
	key := fmt.Sprintf("%d-%s", dataType, point)
	fmt.Println()
	pointInfo, ok := pointToFieldMappingAgc[key]
	return pointInfo.FieldName, ok
}

// GetPointInfoByPointAgc 根据dataType和point值获取AGC点位信息
func GetPointInfoByPointAgc(dataType int, point string) (PointInfoAgc, bool) {
	key := fmt.Sprintf("%d-%s", dataType, point)
	fmt.Println("aaaaaaaaaa", pointToFieldMappingAgc)
	pointInfo, ok := pointToFieldMappingAgc[key]
	return pointInfo, ok
}

// GetAllPointValuesAgc 获取AGC所有point值列表（去重）
func GetAllPointValuesAgc() []string {
	pointSet := make(map[string]bool)
	for key := range pointToFieldMappingAgc {
		// key格式为 "type-value"，提取value部分
		parts := strings.Split(key, "-")
		if len(parts) == 2 {
			pointSet[parts[1]] = true
		}
	}

	points := make([]string, 0, len(pointSet))
	for point := range pointSet {
		points = append(points, point)
	}
	return points
}

// GetPointValuesByTypeAgc 根据数据类型获取AGC的point值列表
func GetPointValuesByTypeAgc(dataType int) []string {
	points := make([]string, 0)
	for key, info := range pointToFieldMappingAgc {
		if info.DataType == dataType {
			// key格式为 "type-value"，提取value部分
			parts := strings.Split(key, "-")
			if len(parts) == 2 {
				points = append(points, parts[1])
			}
		}
	}
	return points
}
