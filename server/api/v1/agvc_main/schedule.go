package agcv_main

import (
    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    agvcMainService "github.com/flipped-aurora/gin-vue-admin/server/service/agvc/agcv_main"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type ScheduleApi struct{}

// CreateSchedule 创建计划曲线
// @Tags     AGVCSchedule
// @Summary  创建AGC/AVC计划曲线
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body agvc.AgvcScheduleCurve true "计划曲线信息"
// @Success  200  {object} response.Response{msg=string} "创建成功"
// @Router   /agvcMain/schedule/create [post]
func (s *ScheduleApi) CreateSchedule(c *gin.Context) {
    var schedule agvc.AgvcScheduleCurve
    if err := c.ShouldBindJSON(&schedule); err != nil {
        response.FailWithMessage(err.Error(), c)
        return
    }

    if err := agvcMainService.ScheduleService.CreateSchedule(&schedule); err != nil {
        global.GVA_LOG.Error("创建计划曲线失败", zap.Error(err))
        response.FailWithMessage("创建失败: "+err.Error(), c)
        return
    }

    response.OkWithMessage("创建成功", c)
}

// UpdateSchedule 更新计划曲线
// @Tags     AGVCSchedule
// @Summary  更新AGC/AVC计划曲线
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body agvc.AgvcScheduleCurve true "计划曲线信息"
// @Success  200  {object} response.Response{msg=string} "更新成功"
// @Router   /agvcMain/schedule/update [put]
func (s *ScheduleApi) UpdateSchedule(c *gin.Context) {
    var schedule agvc.AgvcScheduleCurve
    if err := c.ShouldBindJSON(&schedule); err != nil {
        response.FailWithMessage(err.Error(), c)
        return
    }

    if err := agvcMainService.ScheduleService.UpdateSchedule(&schedule); err != nil {
        global.GVA_LOG.Error("更新计划曲线失败", zap.Error(err))
        response.FailWithMessage("更新失败: "+err.Error(), c)
        return
    }

    response.OkWithMessage("更新成功", c)
}

// DeleteSchedule 删除计划曲线
// @Tags     AGVCSchedule
// @Summary  删除AGC/AVC计划曲线
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    id query int true "计划曲线ID"
// @Success  200  {object} response.Response{msg=string} "删除成功"
// @Router   /agvcMain/schedule/delete [delete]
func (s *ScheduleApi) DeleteSchedule(c *gin.Context) {
    var id struct {
        ID uint `json:"id" form:"id"`
    }
    if err := c.ShouldBindQuery(&id); err != nil {
        response.FailWithMessage(err.Error(), c)
        return
    }

    if err := agvcMainService.ScheduleService.DeleteSchedule(id.ID); err != nil {
        global.GVA_LOG.Error("删除计划曲线失败", zap.Error(err))
        response.FailWithMessage("删除失败: "+err.Error(), c)
        return
    }

    response.OkWithMessage("删除成功", c)
}

// GetScheduleList 获取计划曲线列表
// @Tags     AGVCSchedule
// @Summary  获取AGC/AVC计划曲线列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    bwdNo query int true "并网点编号"
// @Param    type query int false "类型 1:AGC 2:AVC"
// @Param    source query int false "来源 1:本地 2:调度"
// @Success  200  {object} response.Response{data=[]agvc.AgvcScheduleCurve,msg=string} "获取成功"
// @Router   /agvcMain/schedule/list [get]
func (s *ScheduleApi) GetScheduleList(c *gin.Context) {
    var req struct {
        BwdNo  int `json:"bwdNo" form:"bwdNo"`
        Type   int `json:"type" form:"type"`
        Source int `json:"source" form:"source"`
    }

    if err := c.ShouldBindQuery(&req); err != nil {
        response.FailWithMessage(err.Error(), c)
        return
    }

    schedules, err := agvcMainService.ScheduleService.GetSchedulesByBwdNo(req.BwdNo, req.Type, req.Source)
    if err != nil {
        global.GVA_LOG.Error("获取计划曲线列表失败", zap.Error(err))
        response.FailWithMessage("获取失败: "+err.Error(), c)
        return
    }

    response.OkWithData(schedules, c)
}
