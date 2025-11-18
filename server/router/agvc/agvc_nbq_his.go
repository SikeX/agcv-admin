package agvc

import (
    "github.com/flipped-aurora/gin-vue-admin/server/middleware"
    "github.com/gin-gonic/gin"
)

type AgvcNbqHisRouter struct{}

// InitAgvcNbqHisRouter 初始化 agvcNbqHis表 路由信息
func (s *AgvcNbqHisRouter) InitAgvcNbqHisRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
    agvcNbqHisRouter := Router.Group("agvcNbqHis").Use(middleware.OperationRecord())
    agvcNbqHisRouterWithoutRecord := Router.Group("agvcNbqHis")
    agvcNbqHisRouterWithoutAuth := PublicRouter.Group("agvcNbqHis")
    {
        agvcNbqHisRouter.POST("createAgvcNbqHis", agvcNbqHisApi.CreateAgvcNbqHis)             // 新廻agvcNbqHis表
        agvcNbqHisRouter.DELETE("deleteAgvcNbqHis", agvcNbqHisApi.DeleteAgvcNbqHis)           // 删除agvcNbqHis表
        agvcNbqHisRouter.DELETE("deleteAgvcNbqHisByIds", agvcNbqHisApi.DeleteAgvcNbqHisByIds) // 批量删除agvcNbqHis表
        agvcNbqHisRouter.PUT("updateAgvcNbqHis", agvcNbqHisApi.UpdateAgvcNbqHis)              // 更新agvcNbqHis表
        //agvcNbqHisRouter.POST("generateTestData", agvcNbqHisApi.GenerateAgvcNbqTestData)      // 生成测试数据
    }
    {
        agvcNbqHisRouterWithoutRecord.GET("findAgvcNbqHis", agvcNbqHisApi.FindAgvcNbqHis)         // 根据ID获取agvcNbqHis表
        agvcNbqHisRouterWithoutRecord.GET("getAgvcNbqHisList", agvcNbqHisApi.GetAgvcNbqHisList)   // 获取agvcNbqHis表列表
        agvcNbqHisRouterWithoutRecord.POST("getAgvcNbqHistory", agvcNbqHisApi.GetAgvcNbqHistory)  // 获取逆变器历史数据
        agvcNbqHisRouterWithoutRecord.GET("getNbqRealData", agvcNbqHisApi.GetNbqRealData)         // 获取逆变器实时数据
    }
    {
        agvcNbqHisRouterWithoutAuth.GET("getAgvcNbqHisPublic", agvcNbqHisApi.GetAgvcNbqHisPublic) // agvcNbqHis表开放接口
    }
}
