package agvc

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type AgvcQxySettingRouter struct {}

// InitAgvcQxySettingRouter 初始化 气象仪配置 路由信息
func (s *AgvcQxySettingRouter) InitAgvcQxySettingRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	agvcQxySettingRouter := Router.Group("agvcQxySetting").Use(middleware.OperationRecord())
	agvcQxySettingRouterWithoutRecord := Router.Group("agvcQxySetting")
	agvcQxySettingRouterWithoutAuth := PublicRouter.Group("agvcQxySetting")
	{
		agvcQxySettingRouter.POST("createAgvcQxySetting", agvcQxySettingApi.CreateAgvcQxySetting)   // 新建气象仪配置
		agvcQxySettingRouter.DELETE("deleteAgvcQxySetting", agvcQxySettingApi.DeleteAgvcQxySetting) // 删除气象仪配置
		agvcQxySettingRouter.DELETE("deleteAgvcQxySettingByIds", agvcQxySettingApi.DeleteAgvcQxySettingByIds) // 批量删除气象仪配置
		agvcQxySettingRouter.PUT("updateAgvcQxySetting", agvcQxySettingApi.UpdateAgvcQxySetting)    // 更新气象仪配置
	}
	{
		agvcQxySettingRouterWithoutRecord.GET("findAgvcQxySetting", agvcQxySettingApi.FindAgvcQxySetting)        // 根据ID获取气象仪配置
		agvcQxySettingRouterWithoutRecord.GET("getAgvcQxySettingList", agvcQxySettingApi.GetAgvcQxySettingList)  // 获取气象仪配置列表
	}
	{
	    agvcQxySettingRouterWithoutAuth.GET("getAgvcQxySettingPublic", agvcQxySettingApi.GetAgvcQxySettingPublic)  // 气象仪配置开放接口
	}
}