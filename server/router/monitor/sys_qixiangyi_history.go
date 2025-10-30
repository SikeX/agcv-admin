package monitor

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type SysQixiangyiHistoryRouter struct{}

// InitSysQixiangyiHistoryRouter 初始化气象仪监控路由信息
func (s *SysQixiangyiHistoryRouter) InitSysQixiangyiHistoryRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	sysQixiangyiHistoryRouter := Router.Group("sysQixiangyiHistory").Use(middleware.OperationRecord())
	sysQixiangyiHistoryRouterWithoutRecord := Router.Group("sysQixiangyiHistory")

	var sysQixiangyiHistoryApi = v1.ApiGroupApp.MonitorApiGroup.SysQixiangyiHistoryApi
	{
		sysQixiangyiHistoryRouter.POST("generateTestData", sysQixiangyiHistoryApi.GenerateTestData) // 生成测试数据
	}
	{
		sysQixiangyiHistoryRouterWithoutRecord.GET("getSysQixiangyiHistoryList", sysQixiangyiHistoryApi.GetSysQixiangyiHistoryList) // 获取气象仪监控列表
		sysQixiangyiHistoryRouterWithoutRecord.GET("getQixiangyiHistoryData", sysQixiangyiHistoryApi.GetQixiangyiHistoryData)       // 获取气象仪历史数据
	}
}
