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

	initialize.CoapServer() // 初始化并启动CoAP服务

	agvcMain.DataStorage.Initialize()

	// 初始化AGC服务
	agvcMain.AGC.Initialize()

	// 初始化AVC服务
	agvcMain.AVC.Initialize()

	Router := initialize.Routers()

	address := fmt.Sprintf(":%d", global.GVA_CONFIG.System.Addr)

	initServer(address, Router, 10*time.Minute, 10*time.Minute)
}
