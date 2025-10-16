package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type SysQixiangyiSettingRouter struct {}

// InitSysQixiangyiSettingRouter 初始化 气象仪配置 路由信息
func (s *SysQixiangyiSettingRouter) InitSysQixiangyiSettingRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	sysQixiangyiSettingRouter := Router.Group("sysQixiangyiSetting").Use(middleware.OperationRecord())
	sysQixiangyiSettingRouterWithoutRecord := Router.Group("sysQixiangyiSetting")
	sysQixiangyiSettingRouterWithoutAuth := PublicRouter.Group("sysQixiangyiSetting")
	{
		sysQixiangyiSettingRouter.POST("createSysQixiangyiSetting", sysQixiangyiSettingApi.CreateSysQixiangyiSetting)   // 新建气象仪配置
		sysQixiangyiSettingRouter.DELETE("deleteSysQixiangyiSetting", sysQixiangyiSettingApi.DeleteSysQixiangyiSetting) // 删除气象仪配置
		sysQixiangyiSettingRouter.DELETE("deleteSysQixiangyiSettingByIds", sysQixiangyiSettingApi.DeleteSysQixiangyiSettingByIds) // 批量删除气象仪配置
		sysQixiangyiSettingRouter.PUT("updateSysQixiangyiSetting", sysQixiangyiSettingApi.UpdateSysQixiangyiSetting)    // 更新气象仪配置
	}
	{
		sysQixiangyiSettingRouterWithoutRecord.GET("findSysQixiangyiSetting", sysQixiangyiSettingApi.FindSysQixiangyiSetting)        // 根据ID获取气象仪配置
		sysQixiangyiSettingRouterWithoutRecord.GET("getSysQixiangyiSettingList", sysQixiangyiSettingApi.GetSysQixiangyiSettingList)  // 获取气象仪配置列表
	}
	{
	    sysQixiangyiSettingRouterWithoutAuth.GET("getSysQixiangyiSettingPublic", sysQixiangyiSettingApi.GetSysQixiangyiSettingPublic)  // 气象仪配置开放接口
	}
}
