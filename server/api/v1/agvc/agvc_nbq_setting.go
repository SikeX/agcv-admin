package agvc

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
	agvcReq "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AgvcNbqSettingApi struct{}

// CreateAgvcNbqSetting 创建逆变器配置
// @Tags AgvcNbqSetting
// @Summary 创建逆变器配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body agvc.AgvcNbqSetting true "创建逆变器配置"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /agvcNbqSetting/createAgvcNbqSetting [post]
func (agvcNbqSettingApi *AgvcNbqSettingApi) CreateAgvcNbqSetting(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var agvcNbqSetting agvc.AgvcNbqSetting
	err := c.ShouldBindJSON(&agvcNbqSetting)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = agvcNbqSettingService.CreateAgvcNbqSetting(ctx, &agvcNbqSetting)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteAgvcNbqSetting 删除逆变器配置
// @Tags AgvcNbqSetting
// @Summary 删除逆变器配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body agvc.AgvcNbqSetting true "删除逆变器配置"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /agvcNbqSetting/deleteAgvcNbqSetting [delete]
func (agvcNbqSettingApi *AgvcNbqSettingApi) DeleteAgvcNbqSetting(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	id := c.Query("id")
	err := agvcNbqSettingService.DeleteAgvcNbqSetting(ctx, id)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteAgvcNbqSettingByIds 批量删除逆变器配置
// @Tags AgvcNbqSetting
// @Summary 批量删除逆变器配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /agvcNbqSetting/deleteAgvcNbqSettingByIds [delete]
func (agvcNbqSettingApi *AgvcNbqSettingApi) DeleteAgvcNbqSettingByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ids := c.QueryArray("ids[]")
	err := agvcNbqSettingService.DeleteAgvcNbqSettingByIds(ctx, ids)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateAgvcNbqSetting 更新逆变器配置
// @Tags AgvcNbqSetting
// @Summary 更新逆变器配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body agvc.AgvcNbqSetting true "更新逆变器配置"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /agvcNbqSetting/updateAgvcNbqSetting [put]
func (agvcNbqSettingApi *AgvcNbqSettingApi) UpdateAgvcNbqSetting(c *gin.Context) {
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var agvcNbqSetting agvc.AgvcNbqSetting
	err := c.ShouldBindJSON(&agvcNbqSetting)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = agvcNbqSettingService.UpdateAgvcNbqSetting(ctx, agvcNbqSetting)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindAgvcNbqSetting 用id查询逆变器配置
// @Tags AgvcNbqSetting
// @Summary 用id查询逆变器配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id query int true "用id查询逆变器配置"
// @Success 200 {object} response.Response{data=agvc.AgvcNbqSetting,msg=string} "查询成功"
// @Router /agvcNbqSetting/findAgvcNbqSetting [get]
func (agvcNbqSettingApi *AgvcNbqSettingApi) FindAgvcNbqSetting(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	id := c.Query("id")
	reagvcNbqSetting, err := agvcNbqSettingService.GetAgvcNbqSetting(ctx, id)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(reagvcNbqSetting, c)
}

// GetAgvcNbqSettingList 分页获取逆变器配置列表
// @Tags AgvcNbqSetting
// @Summary 分页获取逆变器配置列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query agvcReq.AgvcNbqSettingSearch true "分页获取逆变器配置列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /agvcNbqSetting/getAgvcNbqSettingList [get]
func (agvcNbqSettingApi *AgvcNbqSettingApi) GetAgvcNbqSettingList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo agvcReq.AgvcNbqSettingSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := agvcNbqSettingService.GetAgvcNbqSettingInfoList(ctx, pageInfo)
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

// GetAgvcNbqSettingPublic 不需要鉴权的逆变器配置接口
// @Tags AgvcNbqSetting
// @Summary 不需要鉴权的逆变器配置接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /agvcNbqSetting/getAgvcNbqSettingPublic [get]
func (agvcNbqSettingApi *AgvcNbqSettingApi) GetAgvcNbqSettingPublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口,一般本接口用于C端服务,需要自己实现业务逻辑
	agvcNbqSettingService.GetAgvcNbqSettingPublic(ctx)
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的逆变器配置接口信息",
	}, "获取成功", c)
}

// FindAgvcNbqSettingByInverterNo 根据逆变器编号查询逆变器配置
// @Tags AgvcNbqSetting
// @Summary 根据逆变器编号查询逆变器配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param inverterNo query string true "逆变器编号"
// @Success 200 {object} response.Response{data=agvc.AgvcNbqSetting,msg=string} "查询成功"
// @Router /agvcNbqSetting/findAgvcNbqSettingByInverterNo [get]
func (agvcNbqSettingApi *AgvcNbqSettingApi) FindAgvcNbqSettingByInverterNo(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	inverterNo := c.Query("inverterNo")
	agvcNbqSetting, err := agvcNbqSettingService.GetAgvcNbqSettingByInverterNo(ctx, inverterNo)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(agvcNbqSetting, c)
}
