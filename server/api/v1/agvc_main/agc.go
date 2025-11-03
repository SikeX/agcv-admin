package agcv_main

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc/agvc_main"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc/agvc_main/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type agc struct{}

var AGC = new(agc)

// GetAGCConfig 获取AGC配置
// @Tags     AGVC_AGC
// @Summary  获取AGC配置
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    psid query string true "电站ID"
// @Success  200  {object} response.Response{data=model.AGCConfig,msg=string} "获取成功"
// @Router   /agvc/agc/config [get]
func (a *agc) GetAGCConfig(c *gin.Context) {
	psid := c.Query("psid")
	if psid == "" {
		response.FailWithMessage("电站ID不能为空", c)
		return
	}

	config, err := serviceAGC.GetAGCConfig(psid)
	if err != nil {
		global.GVA_LOG.Error("获取AGC配置失败", zap.Error(err))
		response.FailWithMessage("获取AGC配置失败", c)
		return
	}

	response.OkWithData(config, c)
}

// UpdateAGCConfig 更新AGC配置
// @Tags     AGVC_AGC
// @Summary  更新AGC配置
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.AGCConfigUpdate true "配置信息"
// @Success  200  {object} response.Response{msg=string} "更新成功"
// @Router   /agvc/agc/config [put]
func (a *agc) UpdateAGCConfig(c *gin.Context) {
	var req request.AGCConfigUpdate
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = serviceAGC.UpdateAGCConfig(req)
	if err != nil {
		global.GVA_LOG.Error("更新AGC配置失败", zap.Error(err))
		response.FailWithMessage("更新AGC配置失败", c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

// CreateAGCConfig 创建AGC配置
// @Tags     AGVC_AGC
// @Summary  创建AGC配置
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body model.AGCConfig true "配置信息"
// @Success  200  {object} response.Response{msg=string} "创建成功"
// @Router   /agvc/agc/config [post]
func (a *agc) CreateAGCConfig(c *gin.Context) {
	var config agvc_main.AGCConfig
	err := c.ShouldBindJSON(&config)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = serviceAGC.CreateOrUpdateConfig(&config)
	if err != nil {
		global.GVA_LOG.Error("创建AGC配置失败", zap.Error(err))
		response.FailWithMessage("创建AGC配置失败", c)
		return
	}

	response.OkWithMessage("创建成功", c)
}

// StartAGC 启动AGC控制
// @Tags     AGVC_AGC
// @Summary  启动AGC控制
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    psid query string true "电站ID"
// @Success  200  {object} response.Response{msg=string} "启动成功"
// @Router   /agvc/agc/start [post]
func (a *agc) StartAGC(c *gin.Context) {
	psid := c.Query("psid")
	if psid == "" {
		response.FailWithMessage("电站ID不能为空", c)
		return
	}

	err := serviceAGC.StartAGC(psid)
	if err != nil {
		global.GVA_LOG.Error("启动AGC失败", zap.Error(err))
		response.FailWithMessage("启动AGC失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("AGC已启动", c)
}

// StopAGC 停止AGC控制
// @Tags     AGVC_AGC
// @Summary  停止AGC控制
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    psid query string true "电站ID"
// @Success  200  {object} response.Response{msg=string} "停止成功"
// @Router   /agvc/agc/stop [post]
func (a *agc) StopAGC(c *gin.Context) {
	psid := c.Query("psid")
	if psid == "" {
		response.FailWithMessage("电站ID不能为空", c)
		return
	}

	err := serviceAGC.StopAGC(psid)
	if err != nil {
		global.GVA_LOG.Error("停止AGC失败", zap.Error(err))
		response.FailWithMessage("停止AGC失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("AGC已停止", c)
}

// GetAGCRecords 获取AGC调节记录
// @Tags     AGVC_AGC
// @Summary  获取AGC调节记录
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query request.AGCRegulationRecordSearch true "查询条件"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /agvc/agc/records [get]
func (a *agc) GetAGCRecords(c *gin.Context) {
	var req request.AGCRegulationRecordSearch
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	records, total, err := serviceAGC.GetAGCRecords(req)
	if err != nil {
		global.GVA_LOG.Error("获取AGC记录失败", zap.Error(err))
		response.FailWithMessage("获取AGC记录失败", c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     records,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}
