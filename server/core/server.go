package core

import (
    "fmt"
    "time"

    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/initialize"
    agvcMain "github.com/flipped-aurora/gin-vue-admin/server/service/agvc/agcv_main"
    "github.com/flipped-aurora/gin-vue-admin/server/service/system"
    "go.uber.org/zap"
)

func RunServer() {
    if global.GVA_CONFIG.System.UseRedis {
        // 初始化redis服务
        initialize.Redis()
        if global.GVA_CONFIG.System.UseMultipoint {
            initialize.RedisList()
        }
    }

    if global.GVA_CONFIG.System.UseMongo {
        err := initialize.Mongo.Initialization()
        if err != nil {
            zap.L().Error(fmt.Sprintf("%+v", err))
        }
    }
    // 从db加载jwt数据
    if global.GVA_DB != nil {
        system.LoadAll()
    }

    initialize.InfluxDB() // 初始化InfluxDB

    //initialize.MQTT() // 初始化MQTT
    //initialize.SaveRealData("iot/real-data/iot-ViCgfkLkdPk8Ih2T9AT")

    initialize.CoapServer()         // 初始化并启动CoAP服务（5683端口-采集数据）
    initialize.CoapDispatchServer() // 初始化并启动调度CoAP服务器（1190端口-接收调度数据，1189端口-返回计算结果）

    agvcMain.DataStorage.Initialize()     // 初始化数据存储
    agvcMain.DispatchStorage.Initialize() // 初始化调度数据存储

    // 初始化设备点位映射
    if err := agvcMain.PointMapper.Initialize(); err != nil {
        zap.L().Error("初始化设备点位映射失败", zap.Error(err))
    }

    // 初始化逆变器品牌点位映射
    agvcMain.InverterBrandMapper.Initialize()

    // 初始化华为逆变器5xx点位映射
    if err := agvcMain.HuaweiPointInit.InitializeHuaweiPoints(); err != nil {
        zap.L().Error("初始化华为逆变器点位映射失败", zap.Error(err))
    }

    // 初始化并网点配置缓存
    if err := agvcMain.SettingCache.Initialize(); err != nil {
        zap.L().Error("初始化并网点配置缓存失败", zap.Error(err))
    }

    // 初始化AGC服务
    agvcMain.AGC.Initialize()

    // 初始化AVC服务
    agvcMain.AVC.Initialize()

    // 初始化计划曲线服务
    agvcMain.ScheduleService.Initialize()

    // 自动启动所有并网点的AGC和AVC功能
    go func() {
        // 延迟3秒启动，确保所有依赖服务已完全初始化
        time.Sleep(3 * time.Second)

        global.GVA_LOG.Info("========== 开始自动启动所有并网点的AGC和AVC功能 ==========")

        // 自动启动所有并网点的AGC功能
        agvcMain.AGC.AutoStartAllGridPoints()

        // 自动启动所有并网点的AVC功能
        agvcMain.AVC.AutoStartAllGridPoints()

        global.GVA_LOG.Info("========== 自动启动流程完成 ==========")
    }()

    Router := initialize.Routers()

    address := fmt.Sprintf(":%d", global.GVA_CONFIG.System.Addr)

    initServer(address, Router, 10*time.Minute, 10*time.Minute)
}
