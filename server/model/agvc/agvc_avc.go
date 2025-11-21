// 自动生成模板AgvcAvcHis
package agvc

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// agvcAvcHis表 结构体  AgvcAvc
// 用于存储和展示AVC实时数据和历史数据
type AgvcAvc struct {
	Ctime               *string  `json:"ctime" form:"ctime" point:"name:采集时间"`
	Psid                *int     `json:"psid" form:"psid" point:"name:电站编号"`
	Eqid                *int     `json:"eqid" form:"eqid" point:"name:设备编号"`
	Name                *string  `json:"name" form:"name" point:"name:并网点名称"`
	TargetVoltage       *float64 `json:"targetVoltage" form:"targetVoltage" point:"name:目标电压(kV),value:403,type:2"`
	TargetReactivePower *float64 `json:"targetReactivePower" form:"targetReactivePower" point:"name:目标无功(kVar),value:404,type:2"`
	AvcFunctionState    *float64 `json:"avcFunctionState" form:"avcFunctionState" point:"name:AVC功能投退状态,value:401,type:1"`
	AvcControlMode      *float64 `json:"avcControlMode" form:"avcControlMode" point:"name:AVC调节方式,value:404,type:1"`
	AvcControlAuthority *float64 `json:"avcControlAuthority" form:"avcControlAuthority" point:"name:AVC控制权限,value:402,type:1"`
}

// TableName agvcAvcHis表 AgvcAvcHis自定义表名 agvc_avc_his
func (AgvcAvc) TableName() string {
	return "agvc_avc_his"
}

// PointInfoAvc AVC的point信息结构体
type PointInfoAvc struct {
	FieldName string
	DataType  int // 1=遥信(YX), 2=遥测(YC)
}

// PointToFieldMappingAvc AVC的point值到字段信息的映射缓存
// key格式为 "type-value"，例如 "1-32" 或 "2-30"
var pointToFieldMappingAvc map[string]PointInfoAvc
var pointMappingOnceAvc sync.Once

// InitPointMappingAvc 初始化AVC的point标签映射，在系统启动时调用
func InitPointMappingAvc() {
	pointMappingOnceAvc.Do(func() {
		pointToFieldMappingAvc = make(map[string]PointInfoAvc)

		// 使用反射解析AgvcAvc结构体的point标签
		t := reflect.TypeOf(AgvcAvc{})
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
					pointToFieldMappingAvc[key] = PointInfoAvc{
						FieldName: field.Name,
						DataType:  dataType,
					}
				}
			}
		}
	})
}

// GetFieldNameByPointAvc 根据dataType和point值获取AVC字段名
func GetFieldNameByPointAvc(dataType int, point string) (string, bool) {
	key := fmt.Sprintf("%d-%s", dataType, point)
	pointInfo, ok := pointToFieldMappingAvc[key]
	return pointInfo.FieldName, ok
}

// GetPointInfoByPointAvc 根据dataType和point值获取AVC点位信息
func GetPointInfoByPointAvc(dataType int, point string) (PointInfoAvc, bool) {
	key := fmt.Sprintf("%d-%s", dataType, point)
	pointInfo, ok := pointToFieldMappingAvc[key]
	return pointInfo, ok
}

// GetAllPointValuesAvc 获取AVC所有point值列表（去重）
func GetAllPointValuesAvc() []string {
	pointSet := make(map[string]bool)
	for key := range pointToFieldMappingAvc {
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

// GetPointValuesByTypeAvc 根据数据类型获取AVC的point值列表
func GetPointValuesByTypeAvc(dataType int) []string {
	points := make([]string, 0)
	for key, info := range pointToFieldMappingAvc {
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
