package agvc

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type AgvcQxyHisRouter struct{}

// InitAgvcQxyHisRouter 初始化 气象仪监控 路由信息
func (s *AgvcQxyHisRouter) InitAgvcQxyHisRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	agvcQxyHisRouter := Router.Group("agvcQxyHis").Use(middleware.OperationRecord())
	agvcQxyHisRouterWithoutRecord := Router.Group("agvcQxyHis")
	agvcQxyHisRouterWithoutAuth := PublicRouter.Group("agvcQxyHis")
	{
		agvcQxyHisRouter.POST("createAgvcQxyHis", agvcQxyHisApi.CreateAgvcQxyHis)             // 新建气象仪监控
		agvcQxyHisRouter.DELETE("deleteAgvcQxyHis", agvcQxyHisApi.DeleteAgvcQxyHis)           // 删除气象仪监控
		agvcQxyHisRouter.DELETE("deleteAgvcQxyHisByIds", agvcQxyHisApi.DeleteAgvcQxyHisByIds) // 批量删除气象仪监控
		agvcQxyHisRouter.PUT("updateAgvcQxyHis", agvcQxyHisApi.UpdateAgvcQxyHis)              // 更新气象仪监控
	}
	{
		agvcQxyHisRouterWithoutRecord.GET("findAgvcQxyHis", agvcQxyHisApi.FindAgvcQxyHis)       // 根据ID获取气象仪监控
		agvcQxyHisRouterWithoutRecord.GET("getAgvcQxyHisList", agvcQxyHisApi.GetAgvcQxyHisList) // 获取气象仪监控列表
		agvcQxyHisRouterWithoutRecord.GET("getAgvcQxyHistory", agvcQxyHisApi.GetAgvcQxyHistory) // 获取气象仪历史数据
	}
	{
		agvcQxyHisRouter.POST("generateTestData", agvcQxyHisApi.GenerateAgvcQxyTestData) // 生成气象仪测试数据
	}
	{
		agvcQxyHisRouterWithoutAuth.GET("getAgvcQxyHisPublic", agvcQxyHisApi.GetAgvcQxyHisPublic) // 气象仪监控开放接口
	}
}
