package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type SysSvgSvcSettingRouter struct {}

// InitSysSvgSvcSettingRouter 初始化 SVG/SVC设置 路由信息
func (s *SysSvgSvcSettingRouter) InitSysSvgSvcSettingRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	sysSvgSvcSettingRouter := Router.Group("sysSvgSvcSetting").Use(middleware.OperationRecord())
	sysSvgSvcSettingRouterWithoutRecord := Router.Group("sysSvgSvcSetting")
	sysSvgSvcSettingRouterWithoutAuth := PublicRouter.Group("sysSvgSvcSetting")
	{
		sysSvgSvcSettingRouter.POST("createSysSvgSvcSetting", sysSvgSvcSettingApi.CreateSysSvgSvcSetting)   // 新建SVG/SVC设置
		sysSvgSvcSettingRouter.DELETE("deleteSysSvgSvcSetting", sysSvgSvcSettingApi.DeleteSysSvgSvcSetting) // 删除SVG/SVC设置
		sysSvgSvcSettingRouter.DELETE("deleteSysSvgSvcSettingByIds", sysSvgSvcSettingApi.DeleteSysSvgSvcSettingByIds) // 批量删除SVG/SVC设置
		sysSvgSvcSettingRouter.PUT("updateSysSvgSvcSetting", sysSvgSvcSettingApi.UpdateSysSvgSvcSetting)    // 更新SVG/SVC设置
	}
	{
		sysSvgSvcSettingRouterWithoutRecord.GET("findSysSvgSvcSetting", sysSvgSvcSettingApi.FindSysSvgSvcSetting)        // 根据ID获取SVG/SVC设置
		sysSvgSvcSettingRouterWithoutRecord.GET("getSysSvgSvcSettingList", sysSvgSvcSettingApi.GetSysSvgSvcSettingList)  // 获取SVG/SVC设置列表
	}
	{
	    sysSvgSvcSettingRouterWithoutAuth.GET("getSysSvgSvcSettingPublic", sysSvgSvcSettingApi.GetSysSvgSvcSettingPublic)  // SVG/SVC设置开放接口
	}
}
