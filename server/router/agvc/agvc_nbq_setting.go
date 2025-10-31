package agvc

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type AgvcNbqSettingRouter struct{}

// InitAgvcNbqSettingRouter 初始化 逆变器配置 路由信息
func (s *AgvcNbqSettingRouter) InitAgvcNbqSettingRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	agvcNbqSettingRouter := Router.Group("agvcNbqSetting").Use(middleware.OperationRecord())
	agvcNbqSettingRouterWithoutRecord := Router.Group("agvcNbqSetting")
	agvcNbqSettingRouterWithoutAuth := PublicRouter.Group("agvcNbqSetting")
	{
		agvcNbqSettingRouter.POST("createAgvcNbqSetting", agvcNbqSettingApi.CreateAgvcNbqSetting)             // 新建逆变器配置
		agvcNbqSettingRouter.DELETE("deleteAgvcNbqSetting", agvcNbqSettingApi.DeleteAgvcNbqSetting)           // 删除逆变器配置
		agvcNbqSettingRouter.DELETE("deleteAgvcNbqSettingByIds", agvcNbqSettingApi.DeleteAgvcNbqSettingByIds) // 批量删除逆变器配置
		agvcNbqSettingRouter.PUT("updateAgvcNbqSetting", agvcNbqSettingApi.UpdateAgvcNbqSetting)              // 更新逆变器配置
	}
	{
		agvcNbqSettingRouterWithoutRecord.GET("findAgvcNbqSetting", agvcNbqSettingApi.FindAgvcNbqSetting)                         // 根据ID获取逆变器配置
		agvcNbqSettingRouterWithoutRecord.GET("findAgvcNbqSettingByInverterNo", agvcNbqSettingApi.FindAgvcNbqSettingByInverterNo) // 根据逆变器编号获取逆变器配置
		agvcNbqSettingRouterWithoutRecord.GET("getAgvcNbqSettingList", agvcNbqSettingApi.GetAgvcNbqSettingList)                   // 获取逆变器配置列表
	}
	{
		agvcNbqSettingRouterWithoutAuth.GET("getAgvcNbqSettingPublic", agvcNbqSettingApi.GetAgvcNbqSettingPublic) // 逆变器配置开放接口
	}
}
