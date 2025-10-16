import service from '@/utils/request'
// @Tags SysTestPoint
// @Summary 创建测试管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysTestPoint true "创建测试管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /sysTestPoint/createSysTestPoint [post]
export const createSysTestPoint = (data) => {
  return service({
    url: '/sysTestPoint/createSysTestPoint',
    method: 'post',
    data
  })
}

// @Tags SysTestPoint
// @Summary 删除测试管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysTestPoint true "删除测试管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sysTestPoint/deleteSysTestPoint [delete]
export const deleteSysTestPoint = (params) => {
  return service({
    url: '/sysTestPoint/deleteSysTestPoint',
    method: 'delete',
    params
  })
}

// @Tags SysTestPoint
// @Summary 批量删除测试管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除测试管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sysTestPoint/deleteSysTestPoint [delete]
export const deleteSysTestPointByIds = (params) => {
  return service({
    url: '/sysTestPoint/deleteSysTestPointByIds',
    method: 'delete',
    params
  })
}

// @Tags SysTestPoint
// @Summary 更新测试管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysTestPoint true "更新测试管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /sysTestPoint/updateSysTestPoint [put]
export const updateSysTestPoint = (data) => {
  return service({
    url: '/sysTestPoint/updateSysTestPoint',
    method: 'put',
    data
  })
}

// @Tags SysTestPoint
// @Summary 用id查询测试管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.SysTestPoint true "用id查询测试管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /sysTestPoint/findSysTestPoint [get]
export const findSysTestPoint = (params) => {
  return service({
    url: '/sysTestPoint/findSysTestPoint',
    method: 'get',
    params
  })
}

// @Tags SysTestPoint
// @Summary 分页获取测试管理列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取测试管理列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sysTestPoint/getSysTestPointList [get]
export const getSysTestPointList = (params) => {
  return service({
    url: '/sysTestPoint/getSysTestPointList',
    method: 'get',
    params
  })
}

// @Tags SysTestPoint
// @Summary 不需要鉴权的测试管理接口
// @Accept application/json
// @Produce application/json
// @Param data query systemReq.SysTestPointSearch true "分页获取测试管理列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /sysTestPoint/getSysTestPointPublic [get]
export const getSysTestPointPublic = () => {
  return service({
    url: '/sysTestPoint/getSysTestPointPublic',
    method: 'get',
  })
}
