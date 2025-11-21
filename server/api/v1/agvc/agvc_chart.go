package agvc

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	agvcReq "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AgvcChartApi struct{}

// GetPowerChartData 获取电站出力图表数据
// @Tags AgvcChart
// @Summary 获取电站出力图表数据（BWD当前有功 + AGC目标有功）
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query agvcReq.AgvcPowerChartRequest true "查询参数"
// @Success 200 {object} response.Response{data=[]agvc.PowerChartData,msg=string} "获取成功"
// @Router /agvcChart/getPowerChartData [get]
func (chartApi *AgvcChartApi) GetPowerChartData(c *gin.Context) {
	ctx := c.Request.Context()

	var req agvcReq.AgvcPowerChartRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage("参数绑定失败:"+err.Error(), c)
		return
	}

	// 验证时间参数
	if req.StartTime == "" || req.EndTime == "" {
		response.FailWithMessage("开始时间和结束时间不能为空", c)
		return
	}

	// 调用服务层获取数据
	chartData, err := agvcChartService.GetPowerChartData(ctx, req.Psid, req.Eqid, req.StartTime, req.EndTime)
	if err != nil {
		global.GVA_LOG.Error("获取电站出力图表数据失败!", zap.Error(err))
		response.FailWithMessage("获取电站出力图表数据失败:"+err.Error(), c)
		return
	}

	response.OkWithData(chartData, c)
}

// GetVoltageReactiveChartData 获取电压和无功图表数据
// @Tags AgvcChart
// @Summary 获取电压和无功图表数据（BWD当前电压/无功 + AVC目标电压/无功）
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query agvcReq.AgvcVoltageReactiveChartRequest true "查询参数"
// @Success 200 {object} response.Response{data=[]agvc.VoltageReactiveChartData,msg=string} "获取成功"
// @Router /agvcChart/getVoltageReactiveChartData [get]
func (chartApi *AgvcChartApi) GetVoltageReactiveChartData(c *gin.Context) {
	ctx := c.Request.Context()

	var req agvcReq.AgvcVoltageReactiveChartRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage("参数绑定失败:"+err.Error(), c)
		return
	}

	// 验证时间参数
	if req.StartTime == "" || req.EndTime == "" {
		response.FailWithMessage("开始时间和结束时间不能为空", c)
		return
	}

	// 调用服务层获取数据
	chartData, err := agvcChartService.GetVoltageReactiveChartData(ctx, req.Psid, req.Eqid, req.StartTime, req.EndTime)
	if err != nil {
		global.GVA_LOG.Error("获取电压和无功图表数据失败!", zap.Error(err))
		response.FailWithMessage("获取电压和无功图表数据失败:"+err.Error(), c)
		return
	}

	response.OkWithData(chartData, c)
}
