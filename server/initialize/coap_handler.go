package initialize

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc/agvc_main"
	agvcMainService "github.com/flipped-aurora/gin-vue-admin/server/service/agvc/agcv_main"
	"go.uber.org/zap"
)

// CoapHandler CoAP处理器函数类型
type CoapHandler func(ctx context.Context, msg coapMessage) (code byte, payload []byte)

// coapRoutes CoAP路由表
var coapRoutes = map[string]map[byte]CoapHandler{
	// path -> method -> handler
	"/health": {
		coapCodeGet: handleHealthCheck,
	},
	"/agvc/data": {
		coapCodePost: handleAgvcData,
	},
}

// handleHealthCheck 健康检查处理器
func handleHealthCheck(ctx context.Context, msg coapMessage) (code byte, payload []byte) {
	return coapCodeContent, []byte("ok")
}

// handleAgvcData AGVC数据接收处理器
func handleAgvcData(ctx context.Context, msg coapMessage) (code byte, payload []byte) {
	// 解析JSON数据
	var dataBatch agvc.AgvcDataBatch
	// 查看原始JSON数据
	fmt.Println("cccccccc", string(msg.payload))
	if err := json.Unmarshal(msg.payload, &dataBatch); err != nil {
		global.GVA_LOG.Error("Failed to parse AGVC data", zap.Error(err))
		return coapCodeBadRequest, []byte(`{"error":"invalid json"}`)
	}

	// 验证数据
	if len(dataBatch) == 0 {
		global.GVA_LOG.Warn("Received empty AGVC data batch")
		return coapCodeBadRequest, []byte(`{"error":"empty data"}`)
	}

	// 转换数据格式并存储到DataStorage
	convertedData := make([]agvc_main.AgvcDataItem, len(dataBatch))
	for i, item := range dataBatch {
		convertedData[i] = agvc_main.AgvcDataItem{
			Psid:     item.Psid,
			Eqid:     item.Eqid,
			EqType:   item.EqType,
			DataType: item.DataType,
			Point:    item.Point,
			Value:    item.Value,
		}
	}
	//查看转换后的数据
	global.GVA_LOG.Debug("Converted AGVC data", zap.Any("data", convertedData))

	// 存储到内存，由DataStorage每5分钟定时保存到InfluxDB
	agvcMainService.DataStorage.StoreAgvcDataBatch(convertedData)

	global.GVA_LOG.Debug("AGVC data stored to memory", zap.Int("count", len(dataBatch)))
	return coapCodeCreated, []byte(`{"success":true}`)
}

// getCoapHandler 获取CoAP处理器
func getCoapHandler(path string, method byte) CoapHandler {
	if methods, exists := coapRoutes[path]; exists {
		if handler, exists := methods[method]; exists {
			return handler
		}
	}
	return nil
}
