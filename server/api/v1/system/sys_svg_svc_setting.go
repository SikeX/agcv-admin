package system

import (
	
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/model/system"
    systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type SysSvgSvcSettingApi struct {}



// CreateSysSvgSvcSetting 创建SVG/SVC设置
// @Tags SysSvgSvcSetting
// @Summary 创建SVG/SVC设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body system.SysSvgSvcSetting true "创建SVG/SVC设置"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /sysSvgSvcSetting/createSysSvgSvcSetting [post]
func (sysSvgSvcSettingApi *SysSvgSvcSettingApi) CreateSysSvgSvcSetting(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var sysSvgSvcSetting system.SysSvgSvcSetting
	err := c.ShouldBindJSON(&sysSvgSvcSetting)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = sysSvgSvcSettingService.CreateSysSvgSvcSetting(ctx,&sysSvgSvcSetting)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteSysSvgSvcSetting 删除SVG/SVC设置
// @Tags SysSvgSvcSetting
// @Summary 删除SVG/SVC设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body system.SysSvgSvcSetting true "删除SVG/SVC设置"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /sysSvgSvcSetting/deleteSysSvgSvcSetting [delete]
func (sysSvgSvcSettingApi *SysSvgSvcSettingApi) DeleteSysSvgSvcSetting(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	err := sysSvgSvcSettingService.DeleteSysSvgSvcSetting(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteSysSvgSvcSettingByIds 批量删除SVG/SVC设置
// @Tags SysSvgSvcSetting
// @Summary 批量删除SVG/SVC设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /sysSvgSvcSetting/deleteSysSvgSvcSettingByIds [delete]
func (sysSvgSvcSettingApi *SysSvgSvcSettingApi) DeleteSysSvgSvcSettingByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := sysSvgSvcSettingService.DeleteSysSvgSvcSettingByIds(ctx,IDs)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateSysSvgSvcSetting 更新SVG/SVC设置
// @Tags SysSvgSvcSetting
// @Summary 更新SVG/SVC设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body system.SysSvgSvcSetting true "更新SVG/SVC设置"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /sysSvgSvcSetting/updateSysSvgSvcSetting [put]
func (sysSvgSvcSettingApi *SysSvgSvcSettingApi) UpdateSysSvgSvcSetting(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var sysSvgSvcSetting system.SysSvgSvcSetting
	err := c.ShouldBindJSON(&sysSvgSvcSetting)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = sysSvgSvcSettingService.UpdateSysSvgSvcSetting(ctx,sysSvgSvcSetting)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindSysSvgSvcSetting 用id查询SVG/SVC设置
// @Tags SysSvgSvcSetting
// @Summary 用id查询SVG/SVC设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询SVG/SVC设置"
// @Success 200 {object} response.Response{data=system.SysSvgSvcSetting,msg=string} "查询成功"
// @Router /sysSvgSvcSetting/findSysSvgSvcSetting [get]
func (sysSvgSvcSettingApi *SysSvgSvcSettingApi) FindSysSvgSvcSetting(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	resysSvgSvcSetting, err := sysSvgSvcSettingService.GetSysSvgSvcSetting(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(resysSvgSvcSetting, c)
}
// GetSysSvgSvcSettingList 分页获取SVG/SVC设置列表
// @Tags SysSvgSvcSetting
// @Summary 分页获取SVG/SVC设置列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query systemReq.SysSvgSvcSettingSearch true "分页获取SVG/SVC设置列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /sysSvgSvcSetting/getSysSvgSvcSettingList [get]
func (sysSvgSvcSettingApi *SysSvgSvcSettingApi) GetSysSvgSvcSettingList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo systemReq.SysSvgSvcSettingSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := sysSvgSvcSettingService.GetSysSvgSvcSettingInfoList(ctx,pageInfo)
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

// GetSysSvgSvcSettingPublic 不需要鉴权的SVG/SVC设置接口
// @Tags SysSvgSvcSetting
// @Summary 不需要鉴权的SVG/SVC设置接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /sysSvgSvcSetting/getSysSvgSvcSettingPublic [get]
func (sysSvgSvcSettingApi *SysSvgSvcSettingApi) GetSysSvgSvcSettingPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    sysSvgSvcSettingService.GetSysSvgSvcSettingPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的SVG/SVC设置接口信息",
    }, "获取成功", c)
}
