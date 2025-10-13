package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type SysHisEventRouter struct {}

// InitSysHisEventRouter 初始化 历史事件 路由信息
func (s *SysHisEventRouter) InitSysHisEventRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	sysHisEventRouter := Router.Group("sysHisEvent").Use(middleware.OperationRecord())
	sysHisEventRouterWithoutRecord := Router.Group("sysHisEvent")
	sysHisEventRouterWithoutAuth := PublicRouter.Group("sysHisEvent")
	{
		sysHisEventRouter.POST("createSysHisEvent", sysHisEventApi.CreateSysHisEvent)   // 新建历史事件
		sysHisEventRouter.DELETE("deleteSysHisEvent", sysHisEventApi.DeleteSysHisEvent) // 删除历史事件
		sysHisEventRouter.DELETE("deleteSysHisEventByIds", sysHisEventApi.DeleteSysHisEventByIds) // 批量删除历史事件
		sysHisEventRouter.PUT("updateSysHisEvent", sysHisEventApi.UpdateSysHisEvent)    // 更新历史事件
	}
	{
		sysHisEventRouterWithoutRecord.GET("findSysHisEvent", sysHisEventApi.FindSysHisEvent)        // 根据ID获取历史事件
		sysHisEventRouterWithoutRecord.GET("getSysHisEventList", sysHisEventApi.GetSysHisEventList)  // 获取历史事件列表
	}
	{
	    sysHisEventRouterWithoutAuth.GET("getSysHisEventPublic", sysHisEventApi.GetSysHisEventPublic)  // 历史事件开放接口
	}
}
