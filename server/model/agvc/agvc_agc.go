// 自动生成模板AgvcAgcHis
package agvc

import (
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
	TargetActivePower   *float64 `json:"targetActivePower" form:"targetActivePower" point:"name:目标有功(kW),value:20"`
	DispatchActivePower *float64 `json:"dispatchActivePower" form:"dispatchActivePower" point:"name:调度下发(kW),value:21"`
	AgcFunctionState    *int64   `json:"agcFunctionState" form:"agcFunctionState" point:"name:AGC功能投退状态,value:22"`
	AgcControlMode      *int64   `json:"agcControlMode" form:"agcControlMode" point:"name:AGC调节方式,value:23"`
	AgcControlAuthority *int64   `json:"agcControlAuthority" form:"agcControlAuthority" point:"name:AGC控制权限,value:24"`
	AgcUpperLimit       *float64 `json:"agcUpperLimit" form:"agcUpperLimit" point:"name:可调上限(kW),value:25"`
	AgcLowerLimit       *float64 `json:"agcLowerLimit" form:"agcLowerLimit" point:"name:可调下限(kW),value:26"`
}

// TableName agvcAgcHis表 AgvcAgcHis自定义表名 agvc_agc_his
func (AgvcAgc) TableName() string {
	return "agvc_agc_his"
}

// PointToFieldMappingAgc AGC的point值到字段名的映射缓存
var pointToFieldMappingAgc map[string]string
var pointMappingOnceAgc sync.Once

// InitPointMappingAgc 初始化AGC的point标签映射，在系统启动时调用
func InitPointMappingAgc() {
	pointMappingOnceAgc.Do(func() {
		pointToFieldMappingAgc = make(map[string]string)

		// 使用反射解析AgvcAgc结构体的point标签
		t := reflect.TypeOf(AgvcAgc{})
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
						pointToFieldMappingAgc[pointValue] = field.Name
						break
					}
				}
			}
		}
	})
}

// GetFieldNameByPointAgc 根据point值获取AGC字段名
func GetFieldNameByPointAgc(point string) (string, bool) {
	fieldName, ok := pointToFieldMappingAgc[point]
	return fieldName, ok
}

// GetAllPointValuesAgc 获取AGC所有point值列表
func GetAllPointValuesAgc() []string {
	points := make([]string, 0, len(pointToFieldMappingAgc))
	for point := range pointToFieldMappingAgc {
		points = append(points, point)
	}
	return points
}
