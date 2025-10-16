package system

import (
	
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/model/system"
    systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type SysQixiangyiSettingApi struct {}



// CreateSysQixiangyiSetting 创建气象仪配置
// @Tags SysQixiangyiSetting
// @Summary 创建气象仪配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body system.SysQixiangyiSetting true "创建气象仪配置"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /sysQixiangyiSetting/createSysQixiangyiSetting [post]
func (sysQixiangyiSettingApi *SysQixiangyiSettingApi) CreateSysQixiangyiSetting(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var sysQixiangyiSetting system.SysQixiangyiSetting
	err := c.ShouldBindJSON(&sysQixiangyiSetting)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = sysQixiangyiSettingService.CreateSysQixiangyiSetting(ctx,&sysQixiangyiSetting)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteSysQixiangyiSetting 删除气象仪配置
// @Tags SysQixiangyiSetting
// @Summary 删除气象仪配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body system.SysQixiangyiSetting true "删除气象仪配置"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /sysQixiangyiSetting/deleteSysQixiangyiSetting [delete]
func (sysQixiangyiSettingApi *SysQixiangyiSettingApi) DeleteSysQixiangyiSetting(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	err := sysQixiangyiSettingService.DeleteSysQixiangyiSetting(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteSysQixiangyiSettingByIds 批量删除气象仪配置
// @Tags SysQixiangyiSetting
// @Summary 批量删除气象仪配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /sysQixiangyiSetting/deleteSysQixiangyiSettingByIds [delete]
func (sysQixiangyiSettingApi *SysQixiangyiSettingApi) DeleteSysQixiangyiSettingByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := sysQixiangyiSettingService.DeleteSysQixiangyiSettingByIds(ctx,IDs)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateSysQixiangyiSetting 更新气象仪配置
// @Tags SysQixiangyiSetting
// @Summary 更新气象仪配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body system.SysQixiangyiSetting true "更新气象仪配置"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /sysQixiangyiSetting/updateSysQixiangyiSetting [put]
func (sysQixiangyiSettingApi *SysQixiangyiSettingApi) UpdateSysQixiangyiSetting(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var sysQixiangyiSetting system.SysQixiangyiSetting
	err := c.ShouldBindJSON(&sysQixiangyiSetting)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = sysQixiangyiSettingService.UpdateSysQixiangyiSetting(ctx,sysQixiangyiSetting)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindSysQixiangyiSetting 用id查询气象仪配置
// @Tags SysQixiangyiSetting
// @Summary 用id查询气象仪配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询气象仪配置"
// @Success 200 {object} response.Response{data=system.SysQixiangyiSetting,msg=string} "查询成功"
// @Router /sysQixiangyiSetting/findSysQixiangyiSetting [get]
func (sysQixiangyiSettingApi *SysQixiangyiSettingApi) FindSysQixiangyiSetting(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	resysQixiangyiSetting, err := sysQixiangyiSettingService.GetSysQixiangyiSetting(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(resysQixiangyiSetting, c)
}
// GetSysQixiangyiSettingList 分页获取气象仪配置列表
// @Tags SysQixiangyiSetting
// @Summary 分页获取气象仪配置列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint false "ID"
// @Param deviceName query string false "设备名称"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /sysQixiangyiSetting/getSysQixiangyiSettingList [get]
func (sysQixiangyiSettingApi *SysQixiangyiSettingApi) GetSysQixiangyiSettingList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo systemReq.SysQixiangyiSettingSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := sysQixiangyiSettingService.GetSysQixiangyiSettingInfoList(ctx,pageInfo)
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

// GetSysQixiangyiSettingPublic 不需要鉴权的气象仪配置接口
// @Tags SysQixiangyiSetting
// @Summary 不需要鉴权的气象仪配置接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /sysQixiangyiSetting/getSysQixiangyiSettingPublic [get]
func (sysQixiangyiSettingApi *SysQixiangyiSettingApi) GetSysQixiangyiSettingPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    sysQixiangyiSettingService.GetSysQixiangyiSettingPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的气象仪配置接口信息",
    }, "获取成功", c)
}
