package initialize

import (
    "github.com/flipped-aurora/gin-vue-admin/server/router"
    "github.com/gin-gonic/gin"
)

func holder(routers ...*gin.RouterGroup) {
    _ = routers
    _ = router.RouterGroupApp
}
func initBizRouter(routers ...*gin.RouterGroup) {
    privateGroup := routers[0]
    publicGroup := routers[1]
    holder(publicGroup, privateGroup)
    {
        agvcRouter := router.RouterGroupApp.Agvc
        agvcRouter.InitAgvcBwdSettingRouter(privateGroup, publicGroup)
        agvcRouter.InitAgvcQxySettingRouter(privateGroup, publicGroup)
        agvcRouter.InitAgvcNbqSettingRouter(privateGroup, publicGroup)
        agvcRouter.InitAgvcEventHisRouter(privateGroup, publicGroup)
        agvcRouter.InitAgvcNbqHisRouter(privateGroup, publicGroup)
        agvcRouter.InitAgvcBwdHisRouter(privateGroup, publicGroup)
        agvcRouter.InitAgvcQxyHisRouter(privateGroup, publicGroup)

        router.RouterGroupApp.AgvcMainDevice.InitDeviceRouter(privateGroup, privateGroup)
        
        // AGVC主控路由
        agvcMainRouter := router.RouterGroupApp.AgvcMain
        agvcMainRouter.ScheduleRouter.InitScheduleRouter(privateGroup)
    }
}
