package initialize

import (
    "context"
    "fmt"
    "time"

    "github.com/flipped-aurora/gin-vue-admin/server/global"
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
    UpdateBucketRetentionPolicies()
}

// UpdateBucketRetentionPolicies 更新InfluxDB bucket的保留策略
func UpdateBucketRetentionPolicies() error {
    client := GetInfluxDBClient()
    if client == nil {
        return fmt.Errorf("InfluxDB客户端未初始化")
    }

    config := global.GVA_CONFIG.InfluxDB
    ctx := context.Background()
    bucketsAPI := client.BucketsAPI()

    // 更新agvc储存桶保留策略
    agvcBucketName := config.GetAgvcBucket()
    agvcRetention := config.GetAgvcRetention()
    if err := updateBucketRetention(ctx, bucketsAPI, agvcBucketName, agvcRetention); err != nil {
        global.GVA_LOG.Error("更新agvc储存桶保留策略失败", 
            zap.String("bucket", agvcBucketName), 
            zap.String("retention", agvcRetention),
            zap.Error(err))
    } else {
        global.GVA_LOG.Info("更新agvc储存桶保留策略成功", 
            zap.String("bucket", agvcBucketName), 
            zap.String("retention", agvcRetention))
    }

    // 更新逆变器储存桶保留策略
    nbqBucketName := config.GetNbqBucket()
    nbqRetention := config.GetNbqRetention()
    if err := updateBucketRetention(ctx, bucketsAPI, nbqBucketName, nbqRetention); err != nil {
        global.GVA_LOG.Error("更新逆变器储存桶保留策略失败", 
            zap.String("bucket", nbqBucketName), 
            zap.String("retention", nbqRetention),
            zap.Error(err))
    } else {
        global.GVA_LOG.Info("更新逆变器储存桶保留策略成功", 
            zap.String("bucket", nbqBucketName), 
            zap.String("retention", nbqRetention))
    }

    return nil
}

// updateBucketRetention 更新单个bucket的保留策略
func updateBucketRetention(ctx context.Context, bucketsAPI api.BucketsAPI, bucketName, retention string) error {
    // 解析保留时间
    duration, err := time.ParseDuration(retention)
    if err != nil {
        return fmt.Errorf("解析保留时间失败: %w", err)
    }
    retentionSeconds := int(duration.Seconds())

    // 查找bucket
    bucket, err := bucketsAPI.FindBucketByName(ctx, bucketName)
    if err != nil {
        return fmt.Errorf("查找bucket失败: %w", err)
    }
    if bucket == nil {
        return fmt.Errorf("bucket不存在: %s", bucketName)
    }

    // 检查是否需要更新
    needsUpdate := false
    if len(bucket.RetentionRules) == 0 {
        if retentionSeconds > 0 {
            needsUpdate = true
        }
    } else {
        // 检查第一条保留规则
        currentSeconds := bucket.RetentionRules[0].EverySeconds
        if currentSeconds != int64(retentionSeconds) {
            needsUpdate = true
        }
    }

    // 更新保留策略
    if needsUpdate {
        bucket.RetentionRules = domain.RetentionRules{
            domain.RetentionRule{
                EverySeconds: int64(retentionSeconds),
            },
        }

        _, err := bucketsAPI.UpdateBucket(ctx, bucket)
        if err != nil {
            return fmt.Errorf("更新bucket失败: %w", err)
        }
    }

    return nil
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
