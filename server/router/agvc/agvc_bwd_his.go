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
		agvcBwdHisRouter.POST("createAgvcBwdHis", agvcBwdHisApi.CreateAgvcBwdHis)             // 新建agvcBwdHis表
		agvcBwdHisRouter.DELETE("deleteAgvcBwdHis", agvcBwdHisApi.DeleteAgvcBwdHis)           // 删除agvcBwdHis表
		agvcBwdHisRouter.DELETE("deleteAgvcBwdHisByIds", agvcBwdHisApi.DeleteAgvcBwdHisByIds) // 批量删除agvcBwdHis表
		// agvcBwdHisRouter.PUT("updateAgvcBwdHis", agvcBwdHisApi.UpdateAgvcBwdHis)              // 更新agvcBwdHis表
		agvcBwdHisRouter.POST("sendAgcAvcStatesToTcp", agvcBwdHisApi.SendAgcAvcStatesToTcp) // 发送AGC/AVC状态到TCP
		// agvcBwdHisRouter.PUT("updatePlanCurves", agvcBwdHisApi.UpdatePlanCurves)         // 更新计划曲线
		// agvcBwdHisRouter.PUT("updateAgcParameters", agvcBwdHisApi.UpdateAgcParameters) // 更新AGC参数
		// agvcBwdHisRouter.PUT("updateAvcParameters", agvcBwdHisApi.UpdateAvcParameters) // 更新AVC参数
		// agvcBwdHisRouter.PUT("updateAgcStatus", agvcBwdHisApi.UpdateAgcStatus)         // 更新AGC状态
		// agvcBwdHisRouter.PUT("updateAvcStatus", agvcBwdHisApi.UpdateAvcStatus)         // 更新AVC状态
	}
	{
		agvcBwdHisRouterWithoutRecord.GET("findAgvcBwdHis", agvcBwdHisApi.FindAgvcBwdHis)       // 根据ID获取agvcBwdHis表
		agvcBwdHisRouterWithoutRecord.GET("getAgvcBwdHisList", agvcBwdHisApi.GetAgvcBwdHisList) // 获取agvcBwdHis表列表
		agvcBwdHisRouterWithoutRecord.GET("getAgvcBwdHistory", agvcBwdHisApi.GetAgvcBwdHistory) // 获取并网点历史数据
		// agvcBwdHisRouterWithoutRecord.GET("getAgcParameters", agvcBwdHisApi.GetAgcParameters)     // 获取AGC参数
		// agvcBwdHisRouterWithoutRecord.GET("getAvcParameters", agvcBwdHisApi.GetAvcParameters)     // 获取AVC参数
		// agvcBwdHisRouterWithoutRecord.GET("getPlanCurves", agvcBwdHisApi.GetPlanCurves)           // 获取计划曲线
		agvcBwdHisRouterWithoutRecord.GET("getAgcStatus", agvcBwdHisApi.GetAgcStatus)             // 获取AGC状态
		agvcBwdHisRouterWithoutRecord.GET("getAvcStatus", agvcBwdHisApi.GetAvcStatus)             // 获取AVC状态
		agvcBwdHisRouterWithoutRecord.GET("getBwdRealtimeData", agvcBwdHisApi.GetBwdRealtimeData) // 获取并网点实时数据
	}
	{
		agvcBwdHisRouterWithoutAuth.GET("getAgvcBwdHisPublic", agvcBwdHisApi.GetAgvcBwdHisPublic) // agvcBwdHis表开放接口
	}
}
