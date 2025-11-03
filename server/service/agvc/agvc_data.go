package agvc

import (
	"context"
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"go.uber.org/zap"
)

type AgvcDataService struct{}

// SaveAgvcData 保存AGVC数据到InfluxDB
func (agvcDataService *AgvcDataService) SaveAgvcData(ctx context.Context, dataBatch agvc.AgvcDataBatch) error {
	if global.GVA_INFLUXDB == nil {
		return fmt.Errorf("InfluxDB client not initialized")
	}

	client, ok := global.GVA_INFLUXDB.(influxdb2.Client)
	if !ok {
		return fmt.Errorf("invalid InfluxDB client type")
	}

	// 获取写API
	writeAPI := client.WriteAPI(global.GVA_CONFIG.InfluxDB.Org, global.GVA_CONFIG.InfluxDB.Bucket)

	// 处理异步写入错误
	errorsCh := writeAPI.Errors()
	go func() {
		for err := range errorsCh {
			global.GVA_LOG.Error("InfluxDB write error", zap.Error(err))
		}
	}()

	// 批量写入数据点
	for _, dataItem := range dataBatch {
		// 创建数据点
		point := write.NewPoint(
			"agvc_data", // measurement名称
			map[string]string{ // tags
				"psid":     fmt.Sprintf("%d", dataItem.Psid),
				"eqid":     fmt.Sprintf("%d", dataItem.Eqid),
				"eqType":   fmt.Sprintf("%d", dataItem.EqType),
				"dataType": fmt.Sprintf("%d", dataItem.DataType),
				"point":    dataItem.Point,
			},
			map[string]interface{}{ // fields
				"value": dataItem.Value,
			},
			time.Now(),
		)

		// 写入数据点
		writeAPI.WritePoint(point)

		//global.GVA_LOG.Info("AGVC data written to InfluxDB",
		//	zap.Int("psid", dataItem.Psid),
		//	zap.Int("eqid", dataItem.Eqid),
		//	zap.Int("eqType", dataItem.EqType),
		//	zap.Int("dataType", dataItem.DataType),
		//	zap.String("point", dataItem.Point),
		//	zap.Float64("value", dataItem.Value),
		//)
	}

	// 强制刷新缓冲区
	writeAPI.Flush()

	return nil
}
