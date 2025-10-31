package agvc

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type AgvcBwdHisRouter struct{}

// InitAgvcBwdHisRouter 初始化 agvcBwdHis表 路由信息
func (s *AgvcBwdHisRouter) InitAgvcBwdHisRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	agvcBwdHisRouter := Router.Group("agvcBwdHis").Use(middleware.OperationRecord())
	agvcBwdHisRouterWithoutRecord := Router.Group("agvcBwdHis")
	agvcBwdHisRouterWithoutAuth := PublicRouter.Group("agvcBwdHis")
	{
		agvcBwdHisRouter.POST("createAgvcBwdHis", agvcBwdHisApi.CreateAgvcBwdHis)             // 新廻agvcBwdHis表
		agvcBwdHisRouter.DELETE("deleteAgvcBwdHis", agvcBwdHisApi.DeleteAgvcBwdHis)           // 删除agvcBwdHis表
		agvcBwdHisRouter.DELETE("deleteAgvcBwdHisByIds", agvcBwdHisApi.DeleteAgvcBwdHisByIds) // 批量删除agvcBwdHis表
		agvcBwdHisRouter.PUT("updateAgvcBwdHis", agvcBwdHisApi.UpdateAgvcBwdHis)              // 更新agvcBwdHis表
		agvcBwdHisRouter.POST("generateTestData", agvcBwdHisApi.GenerateAgvcBwdTestData)      // 生成测试数据
	}
	{
		agvcBwdHisRouterWithoutRecord.GET("findAgvcBwdHis", agvcBwdHisApi.FindAgvcBwdHis)       // 根据ID获取agvcBwdHis表
		agvcBwdHisRouterWithoutRecord.GET("getAgvcBwdHisList", agvcBwdHisApi.GetAgvcBwdHisList) // 获取agvcBwdHis表列表
		agvcBwdHisRouterWithoutRecord.GET("getAgvcBwdHistory", agvcBwdHisApi.GetAgvcBwdHistory) // 获取并网点历史数据
	}
	{
		agvcBwdHisRouterWithoutAuth.GET("getAgvcBwdHisPublic", agvcBwdHisApi.GetAgvcBwdHisPublic) // agvcBwdHis表开放接口
	}
}
