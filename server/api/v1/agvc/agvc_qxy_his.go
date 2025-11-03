package agvc

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
	agvcReq "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AgvcQxyHisApi struct{}

// CreateAgvcQxyHis 创建气象仪监控
// @Tags AgvcQxyHis
// @Summary 创建气象仪监控
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body agvc.AgvcQxyHis true "创建气象仪监控"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /agvcQxyHis/createAgvcQxyHis [post]
func (agvcQxyHisApi *AgvcQxyHisApi) CreateAgvcQxyHis(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var agvcQxyHis agvc.AgvcQxyHis
	err := c.ShouldBindJSON(&agvcQxyHis)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = agvcQxyHisService.CreateAgvcQxyHis(ctx, &agvcQxyHis)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteAgvcQxyHis 删除气象仪监控
// @Tags AgvcQxyHis
// @Summary 删除气象仪监控
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body agvc.AgvcQxyHis true "删除气象仪监控"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /agvcQxyHis/deleteAgvcQxyHis [delete]
func (agvcQxyHisApi *AgvcQxyHisApi) DeleteAgvcQxyHis(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	err := agvcQxyHisService.DeleteAgvcQxyHis(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteAgvcQxyHisByIds 批量删除气象仪监控
// @Tags AgvcQxyHis
// @Summary 批量删除气象仪监控
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /agvcQxyHis/deleteAgvcQxyHisByIds [delete]
func (agvcQxyHisApi *AgvcQxyHisApi) DeleteAgvcQxyHisByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := agvcQxyHisService.DeleteAgvcQxyHisByIds(ctx, IDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateAgvcQxyHis 更新气象仪监控
// @Tags AgvcQxyHis
// @Summary 更新气象仪监控
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body agvc.AgvcQxyHis true "更新气象仪监控"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /agvcQxyHis/updateAgvcQxyHis [put]
func (agvcQxyHisApi *AgvcQxyHisApi) UpdateAgvcQxyHis(c *gin.Context) {
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var agvcQxyHis agvc.AgvcQxyHis
	err := c.ShouldBindJSON(&agvcQxyHis)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = agvcQxyHisService.UpdateAgvcQxyHis(ctx, agvcQxyHis)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindAgvcQxyHis 用id查询气象仪监控
// @Tags AgvcQxyHis
// @Summary 用id查询气象仪监控
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询气象仪监控"
// @Success 200 {object} response.Response{data=agvc.AgvcQxyHis,msg=string} "查询成功"
// @Router /agvcQxyHis/findAgvcQxyHis [get]
func (agvcQxyHisApi *AgvcQxyHisApi) FindAgvcQxyHis(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	reagvcQxyHis, err := agvcQxyHisService.GetAgvcQxyHis(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(reagvcQxyHis, c)
}

// GetAgvcQxyHisList 分页获取气象仪监控列表
// @Tags AgvcQxyHis
// @Summary 分页获取气象仪监控列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query agvcReq.AgvcQxyHisSearch true "分页获取气象仪监控列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /agvcQxyHis/getAgvcQxyHisList [get]
func (agvcQxyHisApi *AgvcQxyHisApi) GetAgvcQxyHisList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo agvcReq.AgvcQxyHisSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := agvcQxyHisService.GetAgvcQxyHisInfoList(ctx, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetAgvcQxyHisPublic 不需要鉴权的气象仪监控接口
// @Tags AgvcQxyHis
// @Summary 不需要鉴权的气象仪监控接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /agvcQxyHis/getAgvcQxyHisPublic [get]
func (agvcQxyHisApi *AgvcQxyHisApi) GetAgvcQxyHisPublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	agvcQxyHisService.GetAgvcQxyHisPublic(ctx)
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的气象仪监控接口信息",
	}, "获取成功", c)
}

// GetAgvcQxyHistory 获取气象仪历史数据
// @Tags AgvcQxyHis
// @Summary 获取气象仪历史数据
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param eqid query string true "设备编号"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Success 200 {object} response.Response{data=[]map[string]interface{},msg=string} "获取成功"
// @Router /agvcQxyHis/getAgvcQxyHistory [get]
func (agvcQxyHisApi *AgvcQxyHisApi) GetAgvcQxyHistory(c *gin.Context) {
	ctx := c.Request.Context()

	eqid := c.Query("eqid")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	historyData, err := agvcQxyHisService.GetAgvcQxyHistory(ctx, eqid, startTime, endTime)
	if err != nil {
		global.GVA_LOG.Error("获取历史数据失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithData(historyData, c)
}

// GenerateAgvcQxyTestData 生成气象仪测试数据
// @Tags AgvcQxyHis
// @Summary 生成气象仪测试数据
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body map[string]interface{} true "设备编号和数据条数"
// @Success 200 {object} response.Response{msg=string} "生成成功"
// @Router /agvcQxyHis/generateTestData [post]
func (agvcQxyHisApi *AgvcQxyHisApi) GenerateAgvcQxyTestData(c *gin.Context) {
	ctx := c.Request.Context()

	var req struct {
		Eqid  string `json:"eqid" binding:"required"`
		Count int    `json:"count"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 如果没有指定count,默认100条
	if req.Count <= 0 {
		req.Count = 100
	}

	err := agvcQxyHisService.GenerateAgvcQxyTestData(ctx, req.Eqid, req.Count)
	if err != nil {
		global.GVA_LOG.Error("生成测试数据失败!", zap.Error(err))
		response.FailWithMessage("生成失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("测试数据生成成功", c)
}
