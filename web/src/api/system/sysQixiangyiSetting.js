import service from '@/utils/request'
// @Tags SysQixiangyiSetting
// @Summary 创建气象仪配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysQixiangyiSetting true "创建气象仪配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /sysQixiangyiSetting/createSysQixiangyiSetting [post]
export const createSysQixiangyiSetting = (data) => {
  return service({
    url: '/sysQixiangyiSetting/createSysQixiangyiSetting',
    method: 'post',
    data
  })
}

// @Tags SysQixiangyiSetting
// @Summary 删除气象仪配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysQixiangyiSetting true "删除气象仪配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sysQixiangyiSetting/deleteSysQixiangyiSetting [delete]
export const deleteSysQixiangyiSetting = (params) => {
  return service({
    url: '/sysQixiangyiSetting/deleteSysQixiangyiSetting',
    method: 'delete',
    params
  })
}

// @Tags SysQixiangyiSetting
// @Summary 批量删除气象仪配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除气象仪配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sysQixiangyiSetting/deleteSysQixiangyiSetting [delete]
export const deleteSysQixiangyiSettingByIds = (params) => {
  return service({
    url: '/sysQixiangyiSetting/deleteSysQixiangyiSettingByIds',
    method: 'delete',
    params
  })
}

// @Tags SysQixiangyiSetting
// @Summary 更新气象仪配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysQixiangyiSetting true "更新气象仪配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /sysQixiangyiSetting/updateSysQixiangyiSetting [put]
export const updateSysQixiangyiSetting = (data) => {
  return service({
    url: '/sysQixiangyiSetting/updateSysQixiangyiSetting',
    method: 'put',
    data
  })
}

// @Tags SysQixiangyiSetting
// @Summary 用id查询气象仪配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.SysQixiangyiSetting true "用id查询气象仪配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /sysQixiangyiSetting/findSysQixiangyiSetting [get]
export const findSysQixiangyiSetting = (params) => {
  return service({
    url: '/sysQixiangyiSetting/findSysQixiangyiSetting',
    method: 'get',
    params
  })
}

// @Tags SysQixiangyiSetting
// @Summary 分页获取气象仪配置列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取气象仪配置列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sysQixiangyiSetting/getSysQixiangyiSettingList [get]
export const getSysQixiangyiSettingList = (params) => {
  return service({
    url: '/sysQixiangyiSetting/getSysQixiangyiSettingList',
    method: 'get',
    params
  })
}

// @Tags SysQixiangyiSetting
// @Summary 不需要鉴权的气象仪配置接口
// @Accept application/json
// @Produce application/json
// @Param data query systemReq.SysQixiangyiSettingSearch true "分页获取气象仪配置列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /sysQixiangyiSetting/getSysQixiangyiSettingPublic [get]
export const getSysQixiangyiSettingPublic = () => {
  return service({
    url: '/sysQixiangyiSetting/getSysQixiangyiSettingPublic',
    method: 'get',
  })
}
