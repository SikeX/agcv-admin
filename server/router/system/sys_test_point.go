package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type SysTestPointRouter struct {}

// InitSysTestPointRouter 初始化 测试管理 路由信息
func (s *SysTestPointRouter) InitSysTestPointRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	sysTestPointRouter := Router.Group("sysTestPoint").Use(middleware.OperationRecord())
	sysTestPointRouterWithoutRecord := Router.Group("sysTestPoint")
	sysTestPointRouterWithoutAuth := PublicRouter.Group("sysTestPoint")
	{
		sysTestPointRouter.POST("createSysTestPoint", sysTestPointApi.CreateSysTestPoint)   // 新建测试管理
		sysTestPointRouter.DELETE("deleteSysTestPoint", sysTestPointApi.DeleteSysTestPoint) // 删除测试管理
		sysTestPointRouter.DELETE("deleteSysTestPointByIds", sysTestPointApi.DeleteSysTestPointByIds) // 批量删除测试管理
		sysTestPointRouter.PUT("updateSysTestPoint", sysTestPointApi.UpdateSysTestPoint)    // 更新测试管理
	}
	{
		sysTestPointRouterWithoutRecord.GET("findSysTestPoint", sysTestPointApi.FindSysTestPoint)        // 根据ID获取测试管理
		sysTestPointRouterWithoutRecord.GET("getSysTestPointList", sysTestPointApi.GetSysTestPointList)  // 获取测试管理列表
	}
	{
	    sysTestPointRouterWithoutAuth.GET("getSysTestPointPublic", sysTestPointApi.GetSysTestPointPublic)  // 测试管理开放接口
	}
}
