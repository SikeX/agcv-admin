package agvc

import (
	"github.com/gin-gonic/gin"
)

type AgvcChartRouter struct{}

// InitAgvcChartRouter 初始化 图表数据 路由信息
func (s *AgvcChartRouter) InitAgvcChartRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	agvcChartRouterWithoutRecord := Router.Group("agvcChart")
	{
		agvcChartRouterWithoutRecord.GET("getPowerChartData", agvcChartApi.GetPowerChartData)                         // 获取电站出力图表数据
		agvcChartRouterWithoutRecord.GET("getVoltageReactiveChartData", agvcChartApi.GetVoltageReactiveChartData)     // 获取电压和无功图表数据
	}
}
