import service from '@/utils/request'
// @Tags SysSvgSvcSetting
// @Summary 创建SVG/SVC设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysSvgSvcSetting true "创建SVG/SVC设置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /sysSvgSvcSetting/createSysSvgSvcSetting [post]
export const createSysSvgSvcSetting = (data) => {
  return service({
    url: '/sysSvgSvcSetting/createSysSvgSvcSetting',
    method: 'post',
    data
  })
}

// @Tags SysSvgSvcSetting
// @Summary 删除SVG/SVC设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysSvgSvcSetting true "删除SVG/SVC设置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sysSvgSvcSetting/deleteSysSvgSvcSetting [delete]
export const deleteSysSvgSvcSetting = (params) => {
  return service({
    url: '/sysSvgSvcSetting/deleteSysSvgSvcSetting',
    method: 'delete',
    params
  })
}

// @Tags SysSvgSvcSetting
// @Summary 批量删除SVG/SVC设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除SVG/SVC设置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sysSvgSvcSetting/deleteSysSvgSvcSetting [delete]
export const deleteSysSvgSvcSettingByIds = (params) => {
  return service({
    url: '/sysSvgSvcSetting/deleteSysSvgSvcSettingByIds',
    method: 'delete',
    params
  })
}

// @Tags SysSvgSvcSetting
// @Summary 更新SVG/SVC设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysSvgSvcSetting true "更新SVG/SVC设置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /sysSvgSvcSetting/updateSysSvgSvcSetting [put]
export const updateSysSvgSvcSetting = (data) => {
  return service({
    url: '/sysSvgSvcSetting/updateSysSvgSvcSetting',
    method: 'put',
    data
  })
}

// @Tags SysSvgSvcSetting
// @Summary 用id查询SVG/SVC设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.SysSvgSvcSetting true "用id查询SVG/SVC设置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /sysSvgSvcSetting/findSysSvgSvcSetting [get]
export const findSysSvgSvcSetting = (params) => {
  return service({
    url: '/sysSvgSvcSetting/findSysSvgSvcSetting',
    method: 'get',
    params
  })
}

// @Tags SysSvgSvcSetting
// @Summary 分页获取SVG/SVC设置列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取SVG/SVC设置列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sysSvgSvcSetting/getSysSvgSvcSettingList [get]
export const getSysSvgSvcSettingList = (params) => {
  return service({
    url: '/sysSvgSvcSetting/getSysSvgSvcSettingList',
    method: 'get',
    params
  })
}

// @Tags SysSvgSvcSetting
// @Summary 不需要鉴权的SVG/SVC设置接口
// @Accept application/json
// @Produce application/json
// @Param data query systemReq.SysSvgSvcSettingSearch true "分页获取SVG/SVC设置列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /sysSvgSvcSetting/getSysSvgSvcSettingPublic [get]
export const getSysSvgSvcSettingPublic = () => {
  return service({
    url: '/sysSvgSvcSetting/getSysSvgSvcSettingPublic',
    method: 'get',
  })
}
