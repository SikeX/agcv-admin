// 自动生成模板AgvcAvcHis
package agvc

import (
    "reflect"
    "strings"
    "sync"
)

// agvcAvcHis表 结构体  AgvcAvc
// 用于存储和展示AVC实时数据和历史数据
type AgvcAvc struct {
    Ctime               *string  `json:"ctime" form:"ctime" point:"name:采集时间"`
    Psid                *int     `json:"psid" form:"psid" point:"name:电站编号"`
    Number              *string  `json:"number" form:"number" point:"name:设备编号"`
    Name                *string  `json:"name" form:"name" point:"name:并网点名称"`
    TargetVoltage       *float64 `json:"targetVoltage" form:"targetVoltage" point:"name:目标电压(kV),value:30"`
    TargetReactivePower *float64 `json:"targetReactivePower" form:"targetReactivePower" point:"name:目标无功(kVar),value:31"`
    AvcFunctionState    *int64   `json:"avcFunctionState" form:"avcFunctionState" point:"name:AVC功能投退状态,value:32"`
    AvcControlMode      *int64   `json:"avcControlMode" form:"avcControlMode" point:"name:AVC调节方式,value:33"`
    AvcControlAuthority *int64   `json:"avcControlAuthority" form:"avcControlAuthority" point:"name:AVC控制权限,value:34"`
}

// TableName agvcAvcHis表 AgvcAvcHis自定义表名 agvc_avc_his
func (AgvcAvc) TableName() string {
    return "agvc_avc_his"
}

// PointToFieldMappingAvc AVC的point值到字段名的映射缓存
var pointToFieldMappingAvc map[string]string
var pointMappingOnceAvc sync.Once

// InitPointMappingAvc 初始化AVC的point标签映射，在系统启动时调用
func InitPointMappingAvc() {
    pointMappingOnceAvc.Do(func() {
        pointToFieldMappingAvc = make(map[string]string)

        // 使用反射解析AgvcAvc结构体的point标签
        t := reflect.TypeOf(AgvcAvc{})
        for i := 0; i < t.NumField(); i++ {
            field := t.Field(i)
            tag := field.Tag.Get("point")

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
                        pointToFieldMappingAvc[pointValue] = field.Name
                        break
                    }
                }
            }
        }
    })
}

// GetFieldNameByPointAvc 根据point值获取AVC字段名
func GetFieldNameByPointAvc(point string) (string, bool) {
    fieldName, ok := pointToFieldMappingAvc[point]
    return fieldName, ok
}

// GetAllPointValuesAvc 获取AVC所有point值列表
func GetAllPointValuesAvc() []string {
    points := make([]string, 0, len(pointToFieldMappingAvc))
    for point := range pointToFieldMappingAvc {
        points = append(points, point)
    }
    return points
}
