package initialize

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/agvc/service"
	"go.uber.org/zap"
)

// Service 初始化服务
func Service(ctx context.Context) {
	// 初始化数据存储服务
	service.DataStorage.Initialize()
	
	// 初始化AGC服务
	service.AGC.Initialize()
	
	// 初始化AVC服务
	service.AVC.Initialize()
	
	// 如果有CoAP服务器，注册处理器
	// if global.GVA_COAP_SERVER != nil {
	// 	service.CoapReceiver.RegisterHandlers(global.GVA_COAP_SERVER)
	// }
	
	global.GVA_LOG.Info("AGVC插件服务初始化完成", zap.String("plugin", "agvc"))
}
