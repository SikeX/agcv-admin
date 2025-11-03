package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/agvc/router"
	"github.com/gin-gonic/gin"
)

// Router 初始化路由
func Router(engine *gin.Engine) {
	public := engine.Group(global.GVA_CONFIG.System.RouterPrefix).Group("agvc")
	private := engine.Group(global.GVA_CONFIG.System.RouterPrefix).Group("agvc")
	private.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())

	// 注册各模块路由
	router.Router.Device.InitDeviceRouter(public, private)
	router.Router.History.InitHistoryRouter(public, private)
	router.Router.AGC.InitAGCRouter(public, private)
	router.Router.AVC.InitAVCRouter(public, private)
}
