package system

import (
	
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/model/system"
    systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type SysTestPointApi struct {}



// CreateSysTestPoint 创建测试管理
// @Tags SysTestPoint
// @Summary 创建测试管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body system.SysTestPoint true "创建测试管理"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /sysTestPoint/createSysTestPoint [post]
func (sysTestPointApi *SysTestPointApi) CreateSysTestPoint(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var sysTestPoint system.SysTestPoint
	err := c.ShouldBindJSON(&sysTestPoint)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = sysTestPointService.CreateSysTestPoint(ctx,&sysTestPoint)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteSysTestPoint 删除测试管理
// @Tags SysTestPoint
// @Summary 删除测试管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body system.SysTestPoint true "删除测试管理"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /sysTestPoint/deleteSysTestPoint [delete]
func (sysTestPointApi *SysTestPointApi) DeleteSysTestPoint(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	err := sysTestPointService.DeleteSysTestPoint(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteSysTestPointByIds 批量删除测试管理
// @Tags SysTestPoint
// @Summary 批量删除测试管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /sysTestPoint/deleteSysTestPointByIds [delete]
func (sysTestPointApi *SysTestPointApi) DeleteSysTestPointByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := sysTestPointService.DeleteSysTestPointByIds(ctx,IDs)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateSysTestPoint 更新测试管理
// @Tags SysTestPoint
// @Summary 更新测试管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body system.SysTestPoint true "更新测试管理"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /sysTestPoint/updateSysTestPoint [put]
func (sysTestPointApi *SysTestPointApi) UpdateSysTestPoint(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var sysTestPoint system.SysTestPoint
	err := c.ShouldBindJSON(&sysTestPoint)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = sysTestPointService.UpdateSysTestPoint(ctx,sysTestPoint)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindSysTestPoint 用id查询测试管理
// @Tags SysTestPoint
// @Summary 用id查询测试管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询测试管理"
// @Success 200 {object} response.Response{data=system.SysTestPoint,msg=string} "查询成功"
// @Router /sysTestPoint/findSysTestPoint [get]
func (sysTestPointApi *SysTestPointApi) FindSysTestPoint(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	resysTestPoint, err := sysTestPointService.GetSysTestPoint(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(resysTestPoint, c)
}
// GetSysTestPointList 分页获取测试管理列表
// @Tags SysTestPoint
// @Summary 分页获取测试管理列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query systemReq.SysTestPointSearch true "分页获取测试管理列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /sysTestPoint/getSysTestPointList [get]
func (sysTestPointApi *SysTestPointApi) GetSysTestPointList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo systemReq.SysTestPointSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := sysTestPointService.GetSysTestPointInfoList(ctx,pageInfo)
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

// GetSysTestPointPublic 不需要鉴权的测试管理接口
// @Tags SysTestPoint
// @Summary 不需要鉴权的测试管理接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /sysTestPoint/getSysTestPointPublic [get]
func (sysTestPointApi *SysTestPointApi) GetSysTestPointPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    sysTestPointService.GetSysTestPointPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的测试管理接口信息",
    }, "获取成功", c)
}
