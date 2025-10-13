import service from '@/utils/request'
// @Tags SysGridConnectionPoint
// @Summary 创建并网点配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysGridConnectionPoint true "创建并网点配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /sysGridConnectionPoint/createSysGridConnectionPoint [post]
export const createSysGridConnectionPoint = (data) => {
  return service({
    url: '/sysGridConnectionPoint/createSysGridConnectionPoint',
    method: 'post',
    data
  })
}

// @Tags SysGridConnectionPoint
// @Summary 删除并网点配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysGridConnectionPoint true "删除并网点配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sysGridConnectionPoint/deleteSysGridConnectionPoint [delete]
export const deleteSysGridConnectionPoint = (params) => {
  return service({
    url: '/sysGridConnectionPoint/deleteSysGridConnectionPoint',
    method: 'delete',
    params
  })
}

// @Tags SysGridConnectionPoint
// @Summary 批量删除并网点配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除并网点配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sysGridConnectionPoint/deleteSysGridConnectionPoint [delete]
export const deleteSysGridConnectionPointByIds = (params) => {
  return service({
    url: '/sysGridConnectionPoint/deleteSysGridConnectionPointByIds',
    method: 'delete',
    params
  })
}

// @Tags SysGridConnectionPoint
// @Summary 更新并网点配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysGridConnectionPoint true "更新并网点配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /sysGridConnectionPoint/updateSysGridConnectionPoint [put]
export const updateSysGridConnectionPoint = (data) => {
  return service({
    url: '/sysGridConnectionPoint/updateSysGridConnectionPoint',
    method: 'put',
    data
  })
}

// @Tags SysGridConnectionPoint
// @Summary 用id查询并网点配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.SysGridConnectionPoint true "用id查询并网点配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /sysGridConnectionPoint/findSysGridConnectionPoint [get]
export const findSysGridConnectionPoint = (params) => {
  return service({
    url: '/sysGridConnectionPoint/findSysGridConnectionPoint',
    method: 'get',
    params
  })
}

// @Tags SysGridConnectionPoint
// @Summary 分页获取并网点配置列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取并网点配置列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sysGridConnectionPoint/getSysGridConnectionPointList [get]
export const getSysGridConnectionPointList = (params) => {
  return service({
    url: '/sysGridConnectionPoint/getSysGridConnectionPointList',
    method: 'get',
    params
  })
}

// @Tags SysGridConnectionPoint
// @Summary 不需要鉴权的并网点配置接口
// @Accept application/json
// @Produce application/json
// @Param data query systemReq.SysGridConnectionPointSearch true "分页获取并网点配置列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /sysGridConnectionPoint/getSysGridConnectionPointPublic [get]
export const getSysGridConnectionPointPublic = () => {
  return service({
    url: '/sysGridConnectionPoint/getSysGridConnectionPointPublic',
    method: 'get',
  })
}
