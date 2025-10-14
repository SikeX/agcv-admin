import service from '@/utils/request'
// @Tags SysInverterSetting
// @Summary 创建逆变器设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysInverterSetting true "创建逆变器设置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /sysInverterSetting/createSysInverterSetting [post]
export const createSysInverterSetting = (data) => {
  return service({
    url: '/sysInverterSetting/createSysInverterSetting',
    method: 'post',
    data
  })
}

// @Tags SysInverterSetting
// @Summary 删除逆变器设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysInverterSetting true "删除逆变器设置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sysInverterSetting/deleteSysInverterSetting [delete]
export const deleteSysInverterSetting = (params) => {
  return service({
    url: '/sysInverterSetting/deleteSysInverterSetting',
    method: 'delete',
    params
  })
}

// @Tags SysInverterSetting
// @Summary 批量删除逆变器设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除逆变器设置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sysInverterSetting/deleteSysInverterSetting [delete]
export const deleteSysInverterSettingByIds = (params) => {
  return service({
    url: '/sysInverterSetting/deleteSysInverterSettingByIds',
    method: 'delete',
    params
  })
}

// @Tags SysInverterSetting
// @Summary 更新逆变器设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysInverterSetting true "更新逆变器设置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /sysInverterSetting/updateSysInverterSetting [put]
export const updateSysInverterSetting = (data) => {
  return service({
    url: '/sysInverterSetting/updateSysInverterSetting',
    method: 'put',
    data
  })
}

// @Tags SysInverterSetting
// @Summary 用id查询逆变器设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.SysInverterSetting true "用id查询逆变器设置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /sysInverterSetting/findSysInverterSetting [get]
export const findSysInverterSetting = (params) => {
  return service({
    url: '/sysInverterSetting/findSysInverterSetting',
    method: 'get',
    params
  })
}

// @Tags SysInverterSetting
// @Summary 根据逆变器编号查询逆变器设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param inverterNo query string true "逆变器编号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /sysInverterSetting/findSysInverterSettingByInverterNo [get]
export const findSysInverterSettingByInverterNo = (params) => {
  return service({
    url: '/sysInverterSetting/findSysInverterSettingByInverterNo',
    method: 'get',
    params
  })
}

// @Tags SysInverterSetting
// @Summary 分页获取逆变器设置列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取逆变器设置列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sysInverterSetting/getSysInverterSettingList [get]
export const getSysInverterSettingList = (params) => {
  return service({
    url: '/sysInverterSetting/getSysInverterSettingList',
    method: 'get',
    params
  })
}

// @Tags SysInverterSetting
// @Summary 不需要鉴权的逆变器设置接口
// @Accept application/json
// @Produce application/json
// @Param data query settingReq.SysInverterSettingSearch true "分页获取逆变器设置列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /sysInverterSetting/getSysInverterSettingPublic [get]
export const getSysInverterSettingPublic = () => {
  return service({
    url: '/sysInverterSetting/getSysInverterSettingPublic',
    method: 'get',
  })
}
