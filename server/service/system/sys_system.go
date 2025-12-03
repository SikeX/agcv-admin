package system

import (
    "github.com/flipped-aurora/gin-vue-admin/server/config"
    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/initialize"
    "github.com/flipped-aurora/gin-vue-admin/server/model/system"
    "github.com/flipped-aurora/gin-vue-admin/server/utils"
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
        if updateErr := initialize.UpdateBucketRetentionPolicies(); updateErr != nil {
            global.GVA_LOG.Warn("更新InfluxDB bucket保留策略失败", zap.Error(updateErr))
            // 不返回错误，因为配置已经保存成功，只是保留策略更新失败
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
