package agvc

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
)

type AgvcAvcHisService struct{}

func init() {
	// 系统启动时初始化AVC的point标签映射
	agvc.InitPointMappingAvc()
}

// setFieldByPointForAvc 根据point值设置AgvcAvcHis结构体对应的字段
func setFieldByPointForAvc(avcHis *agvc.AgvcAvc, point string, value float64) {
	v := reflect.ValueOf(avcHis).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("point")

		if tag == "" {
			continue
		}

		if strings.Contains(tag, "value:") {
			parts := strings.Split(tag, ",")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if strings.HasPrefix(part, "value:") {
					pointValue := strings.TrimPrefix(part, "value:")
					if pointValue == point {
						if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.Float64 {
							field.Set(reflect.ValueOf(&value))
							return
						} else if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.Int64 {
							intValue := int64(value)
							field.Set(reflect.ValueOf(&intValue))
							return
						}
					}
				}
			}
		}
	}
}

// extractPointValuesFromAvc 从AgvcAvcHis结构体中提取所有带point标签的字段的value值
func extractPointValuesFromAvc(avcHis agvc.AgvcAvc) []string {
	var pointValues []string
	v := reflect.ValueOf(avcHis)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("point")

		if tag == "" {
			continue
		}

		if field.Kind() == reflect.Ptr && field.IsNil() {
			continue
		}

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

// GetAvcRealData 获取AVC实时数据（从InfluxDB中获取最后一条数据）
// psid: 电站编号, eqid: 并网点编号
func (agvcAvcHisService *AgvcAvcHisService) GetAvcRealData(ctx context.Context, psid, eqid int) (*agvc.AgvcAvc, error) {
	agvcRealData, err := agvcRealService.GetRealData(ctx, psid, eqid, AVC_DEVICE_TYPE)
	if err != nil {
		return nil, err
	}
	if agvcRealData == nil {
		return nil, fmt.Errorf("未找到AVC实时数据")
	}
	realData := agvcRealData.(*agvc.AgvcAvc)

	return realData, nil
}

// setFieldByPointForAvcRealData 根据point值设置AgvcAvc结构体对应的字段
func setFieldByPointForAvcRealData(avcData *agvc.AgvcAvc, point string, value float64) {
	v := reflect.ValueOf(avcData).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("point")

		if tag == "" {
			continue
		}

		if strings.Contains(tag, "value:") {
			parts := strings.Split(tag, ",")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if strings.HasPrefix(part, "value:") {
					pointValue := strings.TrimPrefix(part, "value:")
					if pointValue == point {
						if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.Float64 {
							field.Set(reflect.ValueOf(&value))
							return
						} else if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.Int64 {
							intValue := int64(value)
							field.Set(reflect.ValueOf(&intValue))
							return
						}
					}
				}
			}
		}
	}
}
