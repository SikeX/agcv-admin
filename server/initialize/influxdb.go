package initialize

import (
	"context"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"go.uber.org/zap"
)

// InfluxDB 初始化InfluxDB客户端
func InfluxDB() {
	influxConfig := global.GVA_CONFIG.InfluxDB
	if influxConfig.Host == "" || influxConfig.Port == "" {
		global.GVA_LOG.Info("InfluxDB配置不完整，跳过初始化")
		return
	}

	// 创建InfluxDB客户端
	client := influxdb2.NewClient(
		fmt.Sprintf("http://%s:%s", influxConfig.Host, influxConfig.Port),
		influxConfig.Token,
	)

	// 检查连接是否正常
	_, err := client.Ping(context.Background())
	if err != nil {
		global.GVA_LOG.Error("InfluxDB连接失败", zap.Error(err))
		return
	}

	global.GVA_INFLUXDB = client
	global.GVA_LOG.Info("InfluxDB连接成功")
}

// GetInfluxDBClient 获取InfluxDB客户端实例
func GetInfluxDBClient() influxdb2.Client {
	if client, ok := global.GVA_INFLUXDB.(influxdb2.Client); ok {
		return client
	}
	return nil
}

// GetInfluxDBQueryAPI 获取InfluxDB查询API
func GetInfluxDBQueryAPI() api.QueryAPI {
	client := GetInfluxDBClient()
	if client != nil {
		return client.QueryAPI(global.GVA_CONFIG.InfluxDB.Org)
	}
	return nil
}
