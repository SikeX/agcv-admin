package monitor

import (
	"context"
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/monitor"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SysQixiangyiHistoryApi struct{}

// GetSysQixiangyiHistoryList 获取气象仪监控列表
// @Tags     SysQixiangyiHistory
// @Summary  获取气象仪监控列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{data=[]monitor.SysQixiangyiHistory,msg=string} "获取成功"
// @Router   /sysQixiangyiHistory/getSysQixiangyiHistoryList [get]
func (sysQixiangyiHistoryApi *SysQixiangyiHistoryApi) GetSysQixiangyiHistoryList(c *gin.Context) {
	list, err := sysQixiangyiHistoryService.GetSysQixiangyiHistoryList(context.Background())
	if err != nil {
		global.GVA_LOG.Error("获取气象仪监控列表失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"list": list}, "获取成功", c)
}

// GetQixiangyiHistoryData 获取气象仪历史数据
// @Tags     SysQixiangyiHistory
// @Summary  获取气象仪历史数据
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    qixiangyiId query int true "气象仪ID"
// @Param    startTime query string true "开始时间"
// @Param    endTime query string true "结束时间"
// @Success  200  {object} response.Response{data=[]monitor.QixiangyiHistoryResponse,msg=string} "获取成功"
// @Router   /sysQixiangyiHistory/getQixiangyiHistoryData [get]
func (sysQixiangyiHistoryApi *SysQixiangyiHistoryApi) GetQixiangyiHistoryData(c *gin.Context) {
	// 获取查询参数
	qixiangyiIdStr := c.Query("qixiangyiId")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	// 验证参数
	if qixiangyiIdStr == "" || startTime == "" || endTime == "" {
		response.FailWithMessage("参数不完整", c)
		return
	}

	// 转换气象仪ID
	qixiangyiId, err := strconv.Atoi(qixiangyiIdStr)
	if err != nil {
		response.FailWithMessage("气象仪ID格式错误", c)
		return
	}

	// 构建查询对象
	query := monitor.QixiangyiHistoryQuery{
		QixiangyiID: qixiangyiId,
		StartTime:   startTime,
		EndTime:     endTime,
	}

	// 获取历史数据
	data, err := sysQixiangyiHistoryService.GetQixiangyiHistoryData(context.Background(), query)
	if err != nil {
		global.GVA_LOG.Error("获取气象仪历史数据失败!", zap.Error(err))
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(gin.H{"data": data}, "获取成功", c)
}

// GenerateTestData 生成测试数据
// @Tags     SysQixiangyiHistory
// @Summary  生成气象仪测试数据
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    qixiangyiId query int true "气象仪ID"
// @Success  200  {object} response.Response{msg=string} "生成成功"
// @Router   /sysQixiangyiHistory/generateTestData [post]
func (sysQixiangyiHistoryApi *SysQixiangyiHistoryApi) GenerateTestData(c *gin.Context) {
	// 获取气象仪ID
	qixiangyiIdStr := c.Query("qixiangyiId")
	if qixiangyiIdStr == "" {
		response.FailWithMessage("气象仪ID不能为空", c)
		return
	}

	// 转换气象仪ID
	qixiangyiId, err := strconv.Atoi(qixiangyiIdStr)
	if err != nil {
		response.FailWithMessage("气象仪ID格式错误", c)
		return
	}

	// 生成测试数据
	err = sysQixiangyiHistoryService.GenerateTestData(qixiangyiId)
	if err != nil {
		global.GVA_LOG.Error("生成测试数据失败!", zap.Error(err))
		response.FailWithMessage("生成失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("测试数据生成成功", c)
}
