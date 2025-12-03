package system

import (
	"context"
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	"go.uber.org/zap"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetSystemConfig
//@description: 读取配置文件
//@return: conf config.Server, err error

type SystemConfigService struct{}

var SystemConfigServiceApp = new(SystemConfigService)

func (systemConfigService *SystemConfigService) GetSystemConfig() (conf config.Server, err error) {
	return global.GVA_CONFIG, nil
}

// @description   set system config,
//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetSystemConfig
//@description: 设置配置文件
//@param: system model.System
//@return: err error

func (systemConfigService *SystemConfigService) SetSystemConfig(system system.System) (err error) {
	cs := utils.StructToMap(system.Config)
	for k, v := range cs {
		global.GVA_VP.Set(k, v)
	}
	err = global.GVA_VP.WriteConfig()
	if err != nil {
		return err
	}

	// 重新加载配置到全局变量
	if err = global.GVA_VP.Unmarshal(&global.GVA_CONFIG); err != nil {
		global.GVA_LOG.Error("重新加载配置失败", zap.Error(err))
		return err
	}

	// 如果InfluxDB配置有变更，更新bucket保留策略
	if global.GVA_INFLUXDB != nil {
		if updateErr := UpdateBucketRetentionPolicies(); updateErr != nil {
			global.GVA_LOG.Warn("更新InfluxDB bucket保留策略失败", zap.Error(updateErr))
			// 不返回错误，因为配置已经保存成功，只是保留策略更新失败
		}
	}

	return nil
}

func UpdateBucketRetentionPolicies() error {
	client := global.GVA_INFLUXDB
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

	shardGroupDurationSeconds := int64(retentionSeconds)

	// 更新保留策略
	if needsUpdate {
		bucket.RetentionRules = domain.RetentionRules{
			domain.RetentionRule{
				EverySeconds:              int64(retentionSeconds),
				ShardGroupDurationSeconds: &shardGroupDurationSeconds,
			},
		}

		_, err := bucketsAPI.UpdateBucket(ctx, bucket)
		if err != nil {
			return fmt.Errorf("更新bucket失败: %w", err)
		}
	}

	return nil
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: GetServerInfo
//@description: 获取服务器信息
//@return: server *utils.Server, err error

func (systemConfigService *SystemConfigService) GetServerInfo() (server *utils.Server, err error) {
	var s utils.Server
	s.Os = utils.InitOS()
	if s.Cpu, err = utils.InitCPU(); err != nil {
		global.GVA_LOG.Error("func utils.InitCPU() Failed", zap.String("err", err.Error()))
		return &s, err
	}
	if s.Ram, err = utils.InitRAM(); err != nil {
		global.GVA_LOG.Error("func utils.InitRAM() Failed", zap.String("err", err.Error()))
		return &s, err
	}
	if s.Disk, err = utils.InitDisk(); err != nil {
		global.GVA_LOG.Error("func utils.InitDisk() Failed", zap.String("err", err.Error()))
		return &s, err
	}

	return &s, nil
}
