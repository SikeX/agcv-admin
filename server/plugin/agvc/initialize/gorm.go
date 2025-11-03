package initialize

import (
    "context"

    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/plugin/agvc/model"
    "go.uber.org/zap"
)

// Gorm 初始化数据库表
func Gorm(ctx context.Context) {
    db := global.GVA_DB
    if db == nil {
        global.GVA_LOG.Error("AGVC插件初始化失败：数据库未初始化")
        return
    }

    // 自动迁移表结构
    err := db.WithContext(ctx).AutoMigrate(
        &model.Device{},
        &model.PointMapping{},
        &model.AGCConfig{},
        &model.AGCRegulationRecord{},
        &model.InverterRegulation{},
        &model.AVCConfig{},
        &model.AVCRegulationRecord{},
        &model.DeviceReactiveRegulation{},
    )

    if err != nil {
        global.GVA_LOG.Error("AGVC插件数据库迁移失败", zap.Error(err))
    } else {
        global.GVA_LOG.Info("AGVC插件数据库表初始化成功")
        
        // 初始化测点映射
        InitPointMappings(ctx)
    }
}
