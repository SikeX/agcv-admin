package system

import (
	
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/model/system"
    systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type SysHisEventApi struct {}



// CreateSysHisEvent 创建历史事件
// @Tags SysHisEvent
// @Summary 创建历史事件
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body system.SysHisEvent true "创建历史事件"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /sysHisEvent/createSysHisEvent [post]
func (sysHisEventApi *SysHisEventApi) CreateSysHisEvent(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var sysHisEvent system.SysHisEvent
	err := c.ShouldBindJSON(&sysHisEvent)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = sysHisEventService.CreateSysHisEvent(ctx,&sysHisEvent)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteSysHisEvent 删除历史事件
// @Tags SysHisEvent
// @Summary 删除历史事件
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body system.SysHisEvent true "删除历史事件"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /sysHisEvent/deleteSysHisEvent [delete]
func (sysHisEventApi *SysHisEventApi) DeleteSysHisEvent(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	err := sysHisEventService.DeleteSysHisEvent(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteSysHisEventByIds 批量删除历史事件
// @Tags SysHisEvent
// @Summary 批量删除历史事件
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /sysHisEvent/deleteSysHisEventByIds [delete]
func (sysHisEventApi *SysHisEventApi) DeleteSysHisEventByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := sysHisEventService.DeleteSysHisEventByIds(ctx,IDs)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateSysHisEvent 更新历史事件
// @Tags SysHisEvent
// @Summary 更新历史事件
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body system.SysHisEvent true "更新历史事件"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /sysHisEvent/updateSysHisEvent [put]
func (sysHisEventApi *SysHisEventApi) UpdateSysHisEvent(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var sysHisEvent system.SysHisEvent
	err := c.ShouldBindJSON(&sysHisEvent)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = sysHisEventService.UpdateSysHisEvent(ctx,sysHisEvent)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindSysHisEvent 用id查询历史事件
// @Tags SysHisEvent
// @Summary 用id查询历史事件
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询历史事件"
// @Success 200 {object} response.Response{data=system.SysHisEvent,msg=string} "查询成功"
// @Router /sysHisEvent/findSysHisEvent [get]
func (sysHisEventApi *SysHisEventApi) FindSysHisEvent(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	resysHisEvent, err := sysHisEventService.GetSysHisEvent(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(resysHisEvent, c)
}
// GetSysHisEventList 分页获取历史事件列表
// @Tags SysHisEvent
// @Summary 分页获取历史事件列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query systemReq.SysHisEventSearch true "分页获取历史事件列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /sysHisEvent/getSysHisEventList [get]
func (sysHisEventApi *SysHisEventApi) GetSysHisEventList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo systemReq.SysHisEventSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := sysHisEventService.GetSysHisEventInfoList(ctx,pageInfo)
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

// GetSysHisEventPublic 不需要鉴权的历史事件接口
// @Tags SysHisEvent
// @Summary 不需要鉴权的历史事件接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /sysHisEvent/getSysHisEventPublic [get]
func (sysHisEventApi *SysHisEventApi) GetSysHisEventPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    sysHisEventService.GetSysHisEventPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的历史事件接口信息",
    }, "获取成功", c)
}
