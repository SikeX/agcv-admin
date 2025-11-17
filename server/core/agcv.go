package core

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	agvcMain "github.com/flipped-aurora/gin-vue-admin/server/service/agvc/agcv_main"
	"go.uber.org/zap"
	"time"
)

func RunAGCV() {
	agvcInitialize()
	//run()
}

func agvcInitialize() {
	agvcMain.DataStorage.Initialize()     // 初始化数据存储
	agvcMain.DispatchStorage.Initialize() // 初始化调度数据存储
}

func run() {
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
}
