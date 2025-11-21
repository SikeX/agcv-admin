// 自动生成模板AgvcBwdHis
package agvc

import (
    "reflect"
    "strings"
    "sync"
)

// agvcBwdHis表 结构体  AgvcBwd
// 用于存储和展示并网点实时数据和历史数据
type AgvcBwd struct {
    Ctime                *string  `json:"ctime" form:"ctime" point:"name:采集时间"`
    Psid                 *int     `json:"psid" form:"psid" point:"name:电站编号"`
    Eqid                 *int     `json:"eqid" form:"eqid" point:"name:设备编号"`
    Name                 *string  `json:"name" form:"name" point:"name:并网点名称"`
    CurrentActivePower   *float64 `json:"currentActivePower" form:"currentActivePower" point:"name:当前有功(kW),value:10,type:2"`
    CurrentVoltage       *float64 `json:"currentVoltage" form:"currentVoltage" point:"name:当前电压(kV),value:11,type:2"`
    CurrentReactivePower *float64 `json:"currentReactivePower" form:"currentReactivePower" point:"name:当前无功(kVar),value:12,type:2"`
    SystemFrequency      *float64 `json:"systemFrequency" form:"systemFrequency" point:"name:系统频率(Hz),value:13,type:2"`
    SystemImpedance      *float64 `json:"systemImpedance" form:"systemImpedance" point:"name:系统阻抗,value:14,type:2"`
}

// TableName agvcBwdHis表 AgvcBwdHis自定义表名 agvc_bwd_his
func (AgvcBwd) TableName() string {
    return "agvc_bwd_his"
}

// PointInfoBwd BWD的point信息结构体
type PointInfoBwd struct {
    FieldName string
    DataType  int // 1=遥信(YX), 2=遥测(YC)
}

// PointToFieldMappingBwd BWD的point值到字段信息的映射缓存
var pointToFieldMappingBwd map[string]PointInfoBwd
var pointMappingOnceBwd sync.Once

// InitPointMappingBwd 初始化BWD的point标签映射，在系统启动时调用
func InitPointMappingBwd() {
    pointMappingOnceBwd.Do(func() {
        pointToFieldMappingBwd = make(map[string]PointInfoBwd)

        // 使用反射解析AgvcBwd结构体的point标签
        t := reflect.TypeOf(AgvcBwd{})
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
                    pointToFieldMappingBwd[pointValue] = PointInfoBwd{
                        FieldName: field.Name,
                        DataType:  dataType,
                    }
                }
            }
        }
    })
}

// GetFieldNameByPointBwd 根据point值获取BWD字段名
func GetFieldNameByPointBwd(point string) (string, bool) {
    pointInfo, ok := pointToFieldMappingBwd[point]
    return pointInfo.FieldName, ok
}

// GetPointInfoByPointBwd 根据point值获取BWD点位信息
func GetPointInfoByPointBwd(point string) (PointInfoBwd, bool) {
    pointInfo, ok := pointToFieldMappingBwd[point]
    return pointInfo, ok
}

// GetAllPointValuesBwd 获取BWD所有point值列表
func GetAllPointValuesBwd() []string {
    points := make([]string, 0, len(pointToFieldMappingBwd))
    for point := range pointToFieldMappingBwd {
        points = append(points, point)
    }
    return points
}

// GetPointValuesByTypeBwd 根据数据类型获取BWD的point值列表
func GetPointValuesByTypeBwd(dataType int) []string {
    points := make([]string, 0)
    for point, info := range pointToFieldMappingBwd {
        if info.DataType == dataType {
            points = append(points, point)
        }
    }
    return points
}
