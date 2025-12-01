package initialize

import (
    "context"
    "fmt"

    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/utils"
    influxdb2 "github.com/influxdata/influxdb-client-go/v2"
    "github.com/influxdata/influxdb-client-go/v2/api"
    "github.com/influxdata/influxdb-client-go/v2/domain"
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
    configureRetention(client)
}

func configureRetention(client influxdb2.Client) {
    config := global.GVA_CONFIG.InfluxDB
    if config.Retention == "" || config.Bucket == "" {
        return
    }

    duration, err := utils.ParseDuration(config.Retention)
    if err != nil {
        global.GVA_LOG.Error("解析InfluxDB保留时间失败", zap.Error(err))
        return
    }

    ctx := context.Background()
    bucketsAPI := client.BucketsAPI()

    bucket, err := bucketsAPI.FindBucketByName(ctx, config.Bucket)
    if err != nil {
        global.GVA_LOG.Error("获取InfluxDB Bucket失败", zap.String("bucket", config.Bucket), zap.Error(err))
        return
    }

    // 检查是否需要更新
    retentionSeconds := int(duration.Seconds())
    needsUpdate := false

    if len(bucket.RetentionRules) == 0 {
        if retentionSeconds > 0 {
            needsUpdate = true
        }
    } else {
        // 假设第一条规则是保留规则
        currentSeconds := bucket.RetentionRules[0].EverySeconds
        if currentSeconds != retentionSeconds {
            needsUpdate = true
        }
    }

    if needsUpdate {
        bucket.RetentionRules = []domain.RetentionRule{
            {
                Type:         "expire",
                EverySeconds: retentionSeconds,
            },
        }

        _, err := bucketsAPI.UpdateBucket(ctx, bucket)
        if err != nil {
            global.GVA_LOG.Error("更新InfluxDB保留策略失败", zap.Error(err))
        } else {
            global.GVA_LOG.Info("更新InfluxDB保留策略成功", zap.String("retention", config.Retention))
        }
    }
}

// GetInfluxDBClient 获取InfluxDB客户端实例
func GetInfluxDBClient() influxdb2.Client {
    if client, ok := global.GVA_INFLUXDB.(influxdb2.Client); ok {
        return client
    }
    return nil
}

// GetInfluxDBWriteAPI 获取InfluxDB写API
func GetInfluxDBWriteAPI() api.WriteAPI {
    client := GetInfluxDBClient()
    if client != nil {
        return client.WriteAPI(global.GVA_CONFIG.InfluxDB.Org, global.GVA_CONFIG.InfluxDB.Bucket)
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
