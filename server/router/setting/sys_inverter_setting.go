package setting

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type SysInverterSettingRouter struct {}

// InitSysInverterSettingRouter 初始化 逆变器设置 路由信息
func (s *SysInverterSettingRouter) InitSysInverterSettingRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	sysInverterSettingRouter := Router.Group("sysInverterSetting").Use(middleware.OperationRecord())
	sysInverterSettingRouterWithoutRecord := Router.Group("sysInverterSetting")
	sysInverterSettingRouterWithoutAuth := PublicRouter.Group("sysInverterSetting")
	{
		sysInverterSettingRouter.POST("createSysInverterSetting", sysInverterSettingApi.CreateSysInverterSetting)   // 新建逆变器设置
		sysInverterSettingRouter.DELETE("deleteSysInverterSetting", sysInverterSettingApi.DeleteSysInverterSetting) // 删除逆变器设置
		sysInverterSettingRouter.DELETE("deleteSysInverterSettingByIds", sysInverterSettingApi.DeleteSysInverterSettingByIds) // 批量删除逆变器设置
		sysInverterSettingRouter.PUT("updateSysInverterSetting", sysInverterSettingApi.UpdateSysInverterSetting)    // 更新逆变器设置
	}
	{
		sysInverterSettingRouterWithoutRecord.GET("findSysInverterSetting", sysInverterSettingApi.FindSysInverterSetting)        // 根据ID获取逆变器设置
		sysInverterSettingRouterWithoutRecord.GET("getSysInverterSettingList", sysInverterSettingApi.GetSysInverterSettingList)  // 获取逆变器设置列表
	}
	{
	    sysInverterSettingRouterWithoutAuth.GET("getSysInverterSettingPublic", sysInverterSettingApi.GetSysInverterSettingPublic)  // 逆变器设置开放接口
	}
}
