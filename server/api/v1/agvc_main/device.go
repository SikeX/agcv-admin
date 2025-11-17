package agcv_main

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc/agvc_main"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc/agvc_main/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type device struct{}

var Device = new(device)

// CreateDevice 创建设备
// @Tags     AGVC_Device
// @Summary  创建设备
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body model.Device true "设备信息"
// @Success  200  {object} response.Response{msg=string} "创建成功"
// @Router   /agvc/device/create [post]
func (a *device) CreateDevice(c *gin.Context) {
	var dev agvc_main.Device
	err := c.ShouldBindJSON(&dev)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = serviceDevice.CreateDevice(&dev)
	if err != nil {
		global.GVA_LOG.Error("创建设备失败", zap.Error(err))
		response.FailWithMessage("创建设备失败", c)
		return
	}

	response.OkWithData(dev, c)
}

// DeleteDevice 删除设备
// @Tags     AGVC_Device
// @Summary  删除设备
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    id query uint true "设备ID"
// @Success  200  {object} response.Response{msg=string} "删除成功"
// @Router   /agvc/device/delete [delete]
// func (a *device) DeleteDevice(c *gin.Context) {
// 	var id request.GetById
// 	err := c.ShouldBindQuery(&id)
// 	if err != nil {
// 		response.FailWithMessage(err.Error(), c)
// 		return
// 	}

// 	err = serviceDevice.DeleteDevice(id.Uint())
// 	if err != nil {
// 		global.GVA_LOG.Error("删除设备失败", zap.Error(err))
// 		response.FailWithMessage("删除设备失败", c)
// 		return
// 	}

// 	response.OkWithMessage("删除成功", c)
// }

// UpdateDevice 更新设备
// @Tags     AGVC_Device
// @Summary  更新设备
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body model.Device true "设备信息"
// @Success  200  {object} response.Response{msg=string} "更新成功"
// @Router   /agvc/device/update [put]
func (a *device) UpdateDevice(c *gin.Context) {
	var dev agvc_main.Device
	err := c.ShouldBindJSON(&dev)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = serviceDevice.UpdateDevice(&dev)
	if err != nil {
		global.GVA_LOG.Error("更新设备失败", zap.Error(err))
		response.FailWithMessage("更新设备失败", c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

// GetDevice 获取设备详情
// @Tags     AGVC_Device
// @Summary  获取设备详情
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    id query uint true "设备ID"
// @Success  200  {object} response.Response{data=model.Device,msg=string} "获取成功"
// @Router   /agvc/device/find [get]
// func (a *device) GetDevice(c *gin.Context) {
// 	var id request.GetById
// 	err := c.ShouldBindQuery(&id)
// 	if err != nil {
// 		response.FailWithMessage(err.Error(), c)
// 		return
// 	}

// 	dev, err := serviceDevice.GetDevice(id.Uint())
// 	if err != nil {
// 		global.GVA_LOG.Error("获取设备失败", zap.Error(err))
// 		response.FailWithMessage("获取设备失败", c)
// 		return
// 	}

// 	response.OkWithData(dev, c)
// }

// GetDeviceList 获取设备列表
// @Tags     AGVC_Device
// @Summary  获取设备列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query request.DeviceSearch true "查询条件"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /agvc/device/list [get]
func (a *device) GetDeviceList(c *gin.Context) {
	var req request.DeviceSearch
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := serviceDevice.GetDeviceList(req)
	if err != nil {
		global.GVA_LOG.Error("获取设备列表失败", zap.Error(err))
		response.FailWithMessage("获取设备列表失败", c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}

// GetDeviceRealtimeData 获取设备实时数据
// @Tags     AGVC_Device
// @Summary  获取设备实时数据
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    psid query string true "电站ID"
// @Param    eqid query string true "设备ID"
// @Param    eqType query string true "设备类型"
// @Success  200  {object} response.Response{data=object,msg=string} "获取成功"
// @Router   /agvc/device/realtimeData [get]
func (a *device) GetDeviceRealtimeData(c *gin.Context) {
	psid := c.Query("psid")
	eqid := c.Query("eqid")
	eqType := c.Query("eqType")

	psidInt, err := strconv.Atoi(psid)
	if err != nil {
		response.FailWithMessage("电站ID无效", c)
		return
	}
	eqidInt, err := strconv.Atoi(eqid)
	if err != nil {
		response.FailWithMessage("设备ID无效", c)
		return
	}
	eqTypeInt, err := strconv.Atoi(eqType)
	if err != nil {
		response.FailWithMessage("设备类型无效", c)
		return
	}

	if psidInt == 0 || eqidInt == 0 || eqTypeInt == 0 {
		response.FailWithMessage("参数不完整", c)
		return
	}

	data, err := serviceDevice.GetDeviceRealtimeData(psidInt, eqidInt, eqTypeInt)
	if err != nil {
		global.GVA_LOG.Error("获取设备实时数据失败", zap.Error(err))
		response.FailWithMessage("获取设备实时数据失败", c)
		return
	}

	response.OkWithData(data, c)
}

// GetInverters 获取逆变器列表
// @Tags     AGVC_Device
// @Summary  获取逆变器列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    psid query string true "电站ID"
// @Success  200  {object} response.Response{data=[]model.Device,msg=string} "获取成功"
// @Router   /agvc/device/inverters [get]
func (a *device) GetInverters(c *gin.Context) {
	psid := c.Query("psid")
	if psid == "" {
		response.FailWithMessage("电站ID不能为空", c)
		return
	}

	inverters, err := serviceDevice.GetInvertersByPSID(psid)
	if err != nil {
		global.GVA_LOG.Error("获取逆变器列表失败", zap.Error(err))
		response.FailWithMessage("获取逆变器列表失败", c)
		return
	}

	response.OkWithData(inverters, c)
}
