package agcv_main

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc/agvc_main/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type avc struct{}

var AVC = new(avc)

// GetAVCConfig 获取AVC配置
// @Tags     AGVC_AVC
// @Summary  获取AVC配置
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    psid query string true "电站ID"
// @Success  200  {object} response.Response{data=model.AVCConfig,msg=string} "获取成功"
// @Router   /agvc/avc/config [get]
func (a *avc) GetAVCConfig(c *gin.Context) {
	bwdNo := c.Query("bwdNo")
	if bwdNo == "" {
		response.FailWithMessage("并网点编号不能为空", c)
		return
	}
	bwdNoInt, err := strconv.Atoi(bwdNo)
	if err != nil {
		response.FailWithMessage("并网点编号必须是整数", c)
		return
	}

	config, err := serviceAVC.GetAVCConfig(bwdNoInt)
	if err != nil {
		global.GVA_LOG.Error("获取AVC配置失败", zap.Error(err))
		response.FailWithMessage("获取AVC配置失败", c)
		return
	}

	response.OkWithData(config, c)
}

// // UpdateAVCConfig 更新AVC配置
// // @Tags     AGVC_AVC
// // @Summary  更新AVC配置
// // @Security ApiKeyAuth
// // @accept   application/json
// // @Produce  application/json
// // @Param    data body request.AVCConfigUpdate true "配置信息"
// // @Success  200  {object} response.Response{msg=string} "更新成功"
// // @Router   /agvc/avc/config [put]
func (a *avc) UpdateAVCConfig(c *gin.Context) {
	var req request.AVCConfigUpdate
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = serviceAVC.UpdateAVCConfig(req)
	if err != nil {
		global.GVA_LOG.Error("更新AVC配置失败", zap.Error(err))
		response.FailWithMessage("更新AVC配置失败", c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

// // CreateAVCConfig 创建AVC配置
// // @Tags     AGVC_AVC
// // @Summary  创建AVC配置
// // @Security ApiKeyAuth
// // @accept   application/json
// // @Produce  application/json
// // @Param    data body model.AVCConfig true "配置信息"
// // @Success  200  {object} response.Response{msg=string} "创建成功"
// // @Router   /agvc/avc/config [post]
func (a *avc) CreateAVCConfig(c *gin.Context) {
	var config agvc.AgvcBwdSetting
	err := c.ShouldBindJSON(&config)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = serviceAVC.CreateOrUpdateConfig(&config)
	if err != nil {
		global.GVA_LOG.Error("创建AVC配置失败", zap.Error(err))
		response.FailWithMessage("创建AVC配置失败", c)
		return
	}

	response.OkWithMessage("创建成功", c)
}

// StartAVC 启动AVC控制
// @Tags     AGVC_AVC
// @Summary  启动AVC控制
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    psid query string true "电站ID"
// @Success  200  {object} response.Response{msg=string} "启动成功"
// @Router   /agvc/avc/start [post]
func (a *avc) StartAVC(c *gin.Context) {
	bwdNo := c.Query("bwdNo")
	if bwdNo == "" {
		response.FailWithMessage("并网点编号不能为空", c)
		return
	}
	bwdNoInt, err := strconv.Atoi(bwdNo)
	if err != nil {
		response.FailWithMessage("并网点编号必须是整数", c)
		return
	}

	err = serviceAVC.StartAVC(bwdNoInt)
	if err != nil {
		global.GVA_LOG.Error("启动AVC失败", zap.Error(err))
		response.FailWithMessage("启动AVC失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("AVC已启动", c)
}

// StopAVC 停止AVC控制
// @Tags     AGVC_AVC
// @Summary  停止AVC控制
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    psid query string true "电站ID"
// @Success  200  {object} response.Response{msg=string} "停止成功"
// @Router   /agvc/avc/stop [post]
func (a *avc) StopAVC(c *gin.Context) {
	bwdNo := c.Query("bwdNo")
	if bwdNo == "" {
		response.FailWithMessage("并网点编号不能为空", c)
		return
	}
	bwdNoInt, err := strconv.Atoi(bwdNo)
	if err != nil {
		response.FailWithMessage("并网点编号必须是整数", c)
		return
	}

	err = serviceAVC.StopAVC(bwdNoInt)
	if err != nil {
		global.GVA_LOG.Error("停止AVC失败", zap.Error(err))
		response.FailWithMessage("停止AVC失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("AVC已停止", c)
}

// GetAVCRecords 获取AVC调节记录
// @Tags     AGVC_AVC
// @Summary  获取AVC调节记录
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query request.AVCRegulationRecordSearch true "查询条件"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /agvc/avc/records [get]
func (a *avc) GetAVCRecords(c *gin.Context) {
	var req request.AVCRegulationRecordSearch
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	records, total, err := serviceAVC.GetAVCRecords(req)
	if err != nil {
		global.GVA_LOG.Error("获取AVC记录失败", zap.Error(err))
		response.FailWithMessage("获取AVC记录失败", c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     records,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}
