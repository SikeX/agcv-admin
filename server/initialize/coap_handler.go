package initialize

import (
	"context"
	"encoding/json"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
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
	if err := json.Unmarshal(msg.payload, &dataBatch); err != nil {
		global.GVA_LOG.Error("Failed to parse AGVC data", zap.Error(err))
		return coapCodeBadRequest, []byte(`{"error":"invalid json"}`)
	}

	// 验证数据
	if len(dataBatch) == 0 {
		global.GVA_LOG.Warn("Received empty AGVC data batch")
		return coapCodeBadRequest, []byte(`{"error":"empty data"}`)
	}

	// 调用服务层保存数据
	agvcDataService := service.ServiceGroupApp.AgvcServiceGroup.AgvcDataService
	if err := agvcDataService.SaveAgvcData(ctx, dataBatch); err != nil {
		global.GVA_LOG.Error("Failed to save AGVC data to InfluxDB", zap.Error(err))
		return coapCodeInternalServerError, []byte(`{"error":"save failed"}`)
	}

	global.GVA_LOG.Info("AGVC data saved successfully", zap.Int("count", len(dataBatch)))
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
