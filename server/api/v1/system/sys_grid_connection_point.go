package system

import (
	
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/model/system"
    systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type SysGridConnectionPointApi struct {}



// CreateSysGridConnectionPoint 创建并网点配置
// @Tags SysGridConnectionPoint
// @Summary 创建并网点配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body system.SysGridConnectionPoint true "创建并网点配置"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /sysGridConnectionPoint/createSysGridConnectionPoint [post]
func (sysGridConnectionPointApi *SysGridConnectionPointApi) CreateSysGridConnectionPoint(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var sysGridConnectionPoint system.SysGridConnectionPoint
	err := c.ShouldBindJSON(&sysGridConnectionPoint)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = sysGridConnectionPointService.CreateSysGridConnectionPoint(ctx,&sysGridConnectionPoint)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteSysGridConnectionPoint 删除并网点配置
// @Tags SysGridConnectionPoint
// @Summary 删除并网点配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body system.SysGridConnectionPoint true "删除并网点配置"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /sysGridConnectionPoint/deleteSysGridConnectionPoint [delete]
func (sysGridConnectionPointApi *SysGridConnectionPointApi) DeleteSysGridConnectionPoint(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	err := sysGridConnectionPointService.DeleteSysGridConnectionPoint(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteSysGridConnectionPointByIds 批量删除并网点配置
// @Tags SysGridConnectionPoint
// @Summary 批量删除并网点配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /sysGridConnectionPoint/deleteSysGridConnectionPointByIds [delete]
func (sysGridConnectionPointApi *SysGridConnectionPointApi) DeleteSysGridConnectionPointByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := sysGridConnectionPointService.DeleteSysGridConnectionPointByIds(ctx,IDs)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateSysGridConnectionPoint 更新并网点配置
// @Tags SysGridConnectionPoint
// @Summary 更新并网点配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body system.SysGridConnectionPoint true "更新并网点配置"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /sysGridConnectionPoint/updateSysGridConnectionPoint [put]
func (sysGridConnectionPointApi *SysGridConnectionPointApi) UpdateSysGridConnectionPoint(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var sysGridConnectionPoint system.SysGridConnectionPoint
	err := c.ShouldBindJSON(&sysGridConnectionPoint)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = sysGridConnectionPointService.UpdateSysGridConnectionPoint(ctx,sysGridConnectionPoint)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindSysGridConnectionPoint 用id查询并网点配置
// @Tags SysGridConnectionPoint
// @Summary 用id查询并网点配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询并网点配置"
// @Success 200 {object} response.Response{data=system.SysGridConnectionPoint,msg=string} "查询成功"
// @Router /sysGridConnectionPoint/findSysGridConnectionPoint [get]
func (sysGridConnectionPointApi *SysGridConnectionPointApi) FindSysGridConnectionPoint(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	resysGridConnectionPoint, err := sysGridConnectionPointService.GetSysGridConnectionPoint(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(resysGridConnectionPoint, c)
}
// GetSysGridConnectionPointList 分页获取并网点配置列表
// @Tags SysGridConnectionPoint
// @Summary 分页获取并网点配置列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query systemReq.SysGridConnectionPointSearch true "分页获取并网点配置列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /sysGridConnectionPoint/getSysGridConnectionPointList [get]
func (sysGridConnectionPointApi *SysGridConnectionPointApi) GetSysGridConnectionPointList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo systemReq.SysGridConnectionPointSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := sysGridConnectionPointService.GetSysGridConnectionPointInfoList(ctx,pageInfo)
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

// GetSysGridConnectionPointPublic 不需要鉴权的并网点配置接口
// @Tags SysGridConnectionPoint
// @Summary 不需要鉴权的并网点配置接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /sysGridConnectionPoint/getSysGridConnectionPointPublic [get]
func (sysGridConnectionPointApi *SysGridConnectionPointApi) GetSysGridConnectionPointPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    sysGridConnectionPointService.GetSysGridConnectionPointPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的并网点配置接口信息",
    }, "获取成功", c)
}
