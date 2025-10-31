package agvc

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type AgvcEventHisRouter struct {}

// InitAgvcEventHisRouter 初始化 历史事件 路由信息
func (s *AgvcEventHisRouter) InitAgvcEventHisRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	agvcEventHisRouter := Router.Group("agvcEventHis").Use(middleware.OperationRecord())
	agvcEventHisRouterWithoutRecord := Router.Group("agvcEventHis")
	agvcEventHisRouterWithoutAuth := PublicRouter.Group("agvcEventHis")
	{
		agvcEventHisRouter.POST("createAgvcEventHis", agvcEventHisApi.CreateAgvcEventHis)   // 新建历史事件
		agvcEventHisRouter.DELETE("deleteAgvcEventHis", agvcEventHisApi.DeleteAgvcEventHis) // 删除历史事件
		agvcEventHisRouter.DELETE("deleteAgvcEventHisByIds", agvcEventHisApi.DeleteAgvcEventHisByIds) // 批量删除历史事件
		agvcEventHisRouter.PUT("updateAgvcEventHis", agvcEventHisApi.UpdateAgvcEventHis)    // 更新历史事件
	}
	{
		agvcEventHisRouterWithoutRecord.GET("findAgvcEventHis", agvcEventHisApi.FindAgvcEventHis)        // 根据ID获取历史事件
		agvcEventHisRouterWithoutRecord.GET("getAgvcEventHisList", agvcEventHisApi.GetAgvcEventHisList)  // 获取历史事件列表
	}
	{
	    agvcEventHisRouterWithoutAuth.GET("getAgvcEventHisPublic", agvcEventHisApi.GetAgvcEventHisPublic)  // 历史事件开放接口
	}
}
