package agvc

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type AgvcBwdSettingRouter struct {}

// InitAgvcBwdSettingRouter 初始化 并网点配置 路由信息
func (s *AgvcBwdSettingRouter) InitAgvcBwdSettingRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	agvcBwdSettingRouter := Router.Group("agvcBwdSetting").Use(middleware.OperationRecord())
	agvcBwdSettingRouterWithoutRecord := Router.Group("agvcBwdSetting")
	agvcBwdSettingRouterWithoutAuth := PublicRouter.Group("agvcBwdSetting")
	{
		agvcBwdSettingRouter.POST("createAgvcBwdSetting", agvcBwdSettingApi.CreateAgvcBwdSetting)   // 新建并网点配置
		agvcBwdSettingRouter.DELETE("deleteAgvcBwdSetting", agvcBwdSettingApi.DeleteAgvcBwdSetting) // 删除并网点配置
		agvcBwdSettingRouter.DELETE("deleteAgvcBwdSettingByIds", agvcBwdSettingApi.DeleteAgvcBwdSettingByIds) // 批量删除并网点配置
		agvcBwdSettingRouter.PUT("updateAgvcBwdSetting", agvcBwdSettingApi.UpdateAgvcBwdSetting)    // 更新并网点配置
	}
	{
		agvcBwdSettingRouterWithoutRecord.GET("findAgvcBwdSetting", agvcBwdSettingApi.FindAgvcBwdSetting)        // 根据ID获取并网点配置
		agvcBwdSettingRouterWithoutRecord.GET("getAgvcBwdSettingList", agvcBwdSettingApi.GetAgvcBwdSettingList)  // 获取并网点配置列表
	}
	{
	    agvcBwdSettingRouterWithoutAuth.GET("getAgvcBwdSettingPublic", agvcBwdSettingApi.GetAgvcBwdSettingPublic)  // 并网点配置开放接口
	}
}
