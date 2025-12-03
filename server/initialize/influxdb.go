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

	// 配置数据保留策略
	// configureRetention(client)
}

// func configureRetention(client influxdb2.Client) {
// 	config := global.GVA_CONFIG.InfluxDB
// 	if config.AgvcBucket == "" || config.NbqBucket == "" {
// 		return
// 	}

// 	// 配置agvc储存桶保留策略
// 	agvcDuration, err := utils.ParseDuration(config.AgvcRetention)
// 	if err != nil {
// 		global.GVA_LOG.Error("解析InfluxDB agvc储存桶保留时间失败", zap.Error(err))
// 		return
// 	}
// 	agvcRetentionSeconds := int64(agvcDuration.Seconds())

// 	ctx := context.Background()
// 	bucketsAPI := client.BucketsAPI()

// 	bucket, err := bucketsAPI.FindBucketByName(ctx, config.AgvcBucket)
// 	if err != nil {
// 		global.GVA_LOG.Error("获取InfluxDB agvc储存桶失败", zap.String("bucket", config.AgvcBucket), zap.Error(err))
// 		return
// 	}

// 	// 检查是否需要更新
// 	retentionSeconds := int64(duration.Seconds())
// 	needsUpdate := false

// 	if len(bucket.RetentionRules) == 0 {
// 		if retentionSeconds > 0 {
// 			needsUpdate = true
// 		}
// 	} else {
// 		// 假设第一条规则是保留规则
// 		currentSeconds := bucket.RetentionRules[0].EverySeconds
// 		if currentSeconds != retentionSeconds {
// 			needsUpdate = true
// 		}
// 	}

// 	if needsUpdate {
// 		bucket.RetentionRules = []domain.RetentionRule{
// 			{
// 				EverySeconds: retentionSeconds,
// 			},
// 		}

// 		_, err := bucketsAPI.UpdateBucket(ctx, bucket)
// 		if err != nil {
// 			global.GVA_LOG.Error("更新InfluxDB保留策略失败", zap.Error(err))
// 		} else {
// 			global.GVA_LOG.Info("更新InfluxDB保留策略成功", zap.String("retention", config.Retention))
// 		}
// 	}
// }

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
