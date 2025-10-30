package agvc

import (
	"strconv"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
    agvcReq "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type AgvcQxySettingApi struct {}


// CreateAgvcQxySetting 创建气象仪配置
// @Tags AgvcQxySetting
// @Summary 创建气象仪配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body agvc.AgvcQxySetting true "创建气象仪配置"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /agvcQxySetting/createAgvcQxySetting [post]
func (agvcQxySettingApi *AgvcQxySettingApi) CreateAgvcQxySetting(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var agvcQxySetting agvc.AgvcQxySetting
	err := c.ShouldBindJSON(&agvcQxySetting)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = agvcQxySettingService.CreateAgvcQxySetting(ctx,&agvcQxySetting)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteAgvcQxySetting 删除气象仪配置
// @Tags AgvcQxySetting
// @Summary 删除气象仪配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body agvc.AgvcQxySetting true "删除气象仪配置"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /agvcQxySetting/deleteAgvcQxySetting [delete]
func (agvcQxySettingApi *AgvcQxySettingApi) DeleteAgvcQxySetting(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	err := agvcQxySettingService.DeleteAgvcQxySetting(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteAgvcQxySettingByIds 批量删除气象仪配置
// @Tags AgvcQxySetting
// @Summary 批量删除气象仪配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /agvcQxySetting/deleteAgvcQxySettingByIds [delete]
func (agvcQxySettingApi *AgvcQxySettingApi) DeleteAgvcQxySettingByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var IDS request.IdsReq
    err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 将[]int转换为[]string
	var stringIds []string
	for _, id := range IDS.Ids {
		stringIds = append(stringIds, strconv.Itoa(id))
	}
	err = agvcQxySettingService.DeleteAgvcQxySettingByIds(ctx, stringIds)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateAgvcQxySetting 更新气象仪配置
// @Tags AgvcQxySetting
// @Summary 更新气象仪配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body agvc.AgvcQxySetting true "更新气象仪配置"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /agvcQxySetting/updateAgvcQxySetting [put]
func (agvcQxySettingApi *AgvcQxySettingApi) UpdateAgvcQxySetting(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var agvcQxySetting agvc.AgvcQxySetting
	err := c.ShouldBindJSON(&agvcQxySetting)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = agvcQxySettingService.UpdateAgvcQxySetting(ctx,agvcQxySetting)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindAgvcQxySetting 用id查询气象仪配置
// @Tags AgvcQxySetting
// @Summary 用id查询气象仪配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询气象仪配置"
// @Success 200 {object} response.Response{data=agvc.AgvcQxySetting,msg=string} "查询成功"
// @Router /agvcQxySetting/findAgvcQxySetting [get]
func (agvcQxySettingApi *AgvcQxySettingApi) FindAgvcQxySetting(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	reagvcQxySetting, err := agvcQxySettingService.GetAgvcQxySetting(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(reagvcQxySetting, c)
}

// GetAgvcQxySettingList 分页获取气象仪配置列表
// @Tags AgvcQxySetting
// @Summary 分页获取气象仪配置列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query agvcReq.AgvcQxySettingSearch true "分页获取气象仪配置列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /agvcQxySetting/getAgvcQxySettingList [get]
func (agvcQxySettingApi *AgvcQxySettingApi) GetAgvcQxySettingList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo agvcReq.AgvcQxySettingSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := agvcQxySettingService.GetAgvcQxySettingInfoList(ctx,pageInfo)
	if err != nil {
	    global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:" + err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetAgvcQxySettingPublic 不需要鉴权的气象仪配置接口
// @Tags AgvcQxySetting
// @Summary 不需要鉴权的气象仪配置接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /agvcQxySetting/getAgvcQxySettingPublic [get]
func (agvcQxySettingApi *AgvcQxySettingApi) GetAgvcQxySettingPublic(c *gin.Context) {
    // 此方法不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端项目，需要自己实现业务逻辑
    agvcQxySettingService.GetAgvcQxySettingPublic(c.Request.Context())
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的气象仪配置接口信息",
    }, "获取成功", c)
}