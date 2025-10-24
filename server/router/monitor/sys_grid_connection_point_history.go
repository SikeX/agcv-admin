package monitor

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type SysGridConnectionPointHistoryRouter struct {}

// InitSysGridConnectionPointHistoryRouter 初始化 并网点监控 路由信息
func (s *SysGridConnectionPointHistoryRouter) InitSysGridConnectionPointHistoryRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	sysGridConnectionPointHistoryRouter := Router.Group("sysGridConnectionPointHistory").Use(middleware.OperationRecord())
	sysGridConnectionPointHistoryRouterWithoutRecord := Router.Group("sysGridConnectionPointHistory")
	{
		sysGridConnectionPointHistoryRouterWithoutRecord.GET("getSysGridConnectionPointHistoryList", sysGridConnectionPointHistoryApi.GetSysGridConnectionPointHistoryList)  // 获取并网点监控列表
		sysGridConnectionPointHistoryRouterWithoutRecord.POST("getGridPointHistoryData", sysGridConnectionPointHistoryApi.GetGridPointHistoryData)  // 获取并网点历史数据
		sysGridConnectionPointHistoryRouter.POST("generateTestData", sysGridConnectionPointHistoryApi.GenerateTestData)  // 生成测试数据
	}
}