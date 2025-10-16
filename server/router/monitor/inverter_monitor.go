package monitor

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type InverterMonitorRouter struct {}

// InitInverterMonitorRouter 初始化 逆变器监控 路由信息
func (s *InverterMonitorRouter) InitInverterMonitorRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	inverterMonitorRouter := Router.Group("inverterMonitor").Use(middleware.OperationRecord())
	inverterMonitorRouterWithoutRecord := Router.Group("inverterMonitor")
	inverterMonitorRouterWithoutAuth := PublicRouter.Group("inverterMonitor")
	{
		inverterMonitorRouter.POST("createInverterMonitor", inverterMonitorApi.CreateInverterMonitor)   // 新建逆变器监控
		inverterMonitorRouter.DELETE("deleteInverterMonitor", inverterMonitorApi.DeleteInverterMonitor) // 删除逆变器监控
		inverterMonitorRouter.DELETE("deleteInverterMonitorByIds", inverterMonitorApi.DeleteInverterMonitorByIds) // 批量删除逆变器监控
		inverterMonitorRouter.PUT("updateInverterMonitor", inverterMonitorApi.UpdateInverterMonitor)    // 更新逆变器监控
	}
	{
		inverterMonitorRouterWithoutRecord.GET("findInverterMonitor", inverterMonitorApi.FindInverterMonitor)        // 根据ID获取逆变器监控
		inverterMonitorRouterWithoutRecord.GET("getInverterMonitorList", inverterMonitorApi.GetInverterMonitorList)  // 获取逆变器监控列表
		inverterMonitorRouterWithoutRecord.GET("getInverterHistory", inverterMonitorApi.GetInverterHistory)  // 获取逆变器历史数据
	}
	{
	    inverterMonitorRouterWithoutAuth.GET("getInverterMonitorPublic", inverterMonitorApi.GetInverterMonitorPublic)  // 逆变器监控开放接口
	}
}