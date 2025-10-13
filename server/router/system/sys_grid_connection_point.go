package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type SysGridConnectionPointRouter struct {}

// InitSysGridConnectionPointRouter 初始化 并网点配置 路由信息
func (s *SysGridConnectionPointRouter) InitSysGridConnectionPointRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	sysGridConnectionPointRouter := Router.Group("sysGridConnectionPoint").Use(middleware.OperationRecord())
	sysGridConnectionPointRouterWithoutRecord := Router.Group("sysGridConnectionPoint")
	sysGridConnectionPointRouterWithoutAuth := PublicRouter.Group("sysGridConnectionPoint")
	{
		sysGridConnectionPointRouter.POST("createSysGridConnectionPoint", sysGridConnectionPointApi.CreateSysGridConnectionPoint)   // 新建并网点配置
		sysGridConnectionPointRouter.DELETE("deleteSysGridConnectionPoint", sysGridConnectionPointApi.DeleteSysGridConnectionPoint) // 删除并网点配置
		sysGridConnectionPointRouter.DELETE("deleteSysGridConnectionPointByIds", sysGridConnectionPointApi.DeleteSysGridConnectionPointByIds) // 批量删除并网点配置
		sysGridConnectionPointRouter.PUT("updateSysGridConnectionPoint", sysGridConnectionPointApi.UpdateSysGridConnectionPoint)    // 更新并网点配置
	}
	{
		sysGridConnectionPointRouterWithoutRecord.GET("findSysGridConnectionPoint", sysGridConnectionPointApi.FindSysGridConnectionPoint)        // 根据ID获取并网点配置
		sysGridConnectionPointRouterWithoutRecord.GET("getSysGridConnectionPointList", sysGridConnectionPointApi.GetSysGridConnectionPointList)  // 获取并网点配置列表
	}
	{
	    sysGridConnectionPointRouterWithoutAuth.GET("getSysGridConnectionPointPublic", sysGridConnectionPointApi.GetSysGridConnectionPointPublic)  // 并网点配置开放接口
	}
}
