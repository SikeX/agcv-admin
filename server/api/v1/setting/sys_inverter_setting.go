package setting

import (
	
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/model/setting"
    settingReq "github.com/flipped-aurora/gin-vue-admin/server/model/setting/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type SysInverterSettingApi struct {}



// CreateSysInverterSetting 创建逆变器设置
// @Tags SysInverterSetting
// @Summary 创建逆变器设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body setting.SysInverterSetting true "创建逆变器设置"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /sysInverterSetting/createSysInverterSetting [post]
func (sysInverterSettingApi *SysInverterSettingApi) CreateSysInverterSetting(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var sysInverterSetting setting.SysInverterSetting
	err := c.ShouldBindJSON(&sysInverterSetting)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = sysInverterSettingService.CreateSysInverterSetting(ctx,&sysInverterSetting)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteSysInverterSetting 删除逆变器设置
// @Tags SysInverterSetting
// @Summary 删除逆变器设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body setting.SysInverterSetting true "删除逆变器设置"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /sysInverterSetting/deleteSysInverterSetting [delete]
func (sysInverterSettingApi *SysInverterSettingApi) DeleteSysInverterSetting(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	id := c.Query("id")
	err := sysInverterSettingService.DeleteSysInverterSetting(ctx,id)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteSysInverterSettingByIds 批量删除逆变器设置
// @Tags SysInverterSetting
// @Summary 批量删除逆变器设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /sysInverterSetting/deleteSysInverterSettingByIds [delete]
func (sysInverterSettingApi *SysInverterSettingApi) DeleteSysInverterSettingByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ids := c.QueryArray("ids[]")
	err := sysInverterSettingService.DeleteSysInverterSettingByIds(ctx,ids)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateSysInverterSetting 更新逆变器设置
// @Tags SysInverterSetting
// @Summary 更新逆变器设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body setting.SysInverterSetting true "更新逆变器设置"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /sysInverterSetting/updateSysInverterSetting [put]
func (sysInverterSettingApi *SysInverterSettingApi) UpdateSysInverterSetting(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var sysInverterSetting setting.SysInverterSetting
	err := c.ShouldBindJSON(&sysInverterSetting)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = sysInverterSettingService.UpdateSysInverterSetting(ctx,sysInverterSetting)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindSysInverterSetting 用id查询逆变器设置
// @Tags SysInverterSetting
// @Summary 用id查询逆变器设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id query int true "用id查询逆变器设置"
// @Success 200 {object} response.Response{data=setting.SysInverterSetting,msg=string} "查询成功"
// @Router /sysInverterSetting/findSysInverterSetting [get]
func (sysInverterSettingApi *SysInverterSettingApi) FindSysInverterSetting(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	id := c.Query("id")
	resysInverterSetting, err := sysInverterSettingService.GetSysInverterSetting(ctx,id)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(resysInverterSetting, c)
}
// GetSysInverterSettingList 分页获取逆变器设置列表
// @Tags SysInverterSetting
// @Summary 分页获取逆变器设置列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query settingReq.SysInverterSettingSearch true "分页获取逆变器设置列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /sysInverterSetting/getSysInverterSettingList [get]
func (sysInverterSettingApi *SysInverterSettingApi) GetSysInverterSettingList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo settingReq.SysInverterSettingSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := sysInverterSettingService.GetSysInverterSettingInfoList(ctx,pageInfo)
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

// GetSysInverterSettingPublic 不需要鉴权的逆变器设置接口
// @Tags SysInverterSetting
// @Summary 不需要鉴权的逆变器设置接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /sysInverterSetting/getSysInverterSettingPublic [get]
func (sysInverterSettingApi *SysInverterSettingApi) GetSysInverterSettingPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    sysInverterSettingService.GetSysInverterSettingPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的逆变器设置接口信息",
    }, "获取成功", c)
}
