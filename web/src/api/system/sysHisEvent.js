import service from '@/utils/request'
// @Tags SysHisEvent
// @Summary 创建历史事件
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysHisEvent true "创建历史事件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /sysHisEvent/createSysHisEvent [post]
export const createSysHisEvent = (data) => {
  return service({
    url: '/sysHisEvent/createSysHisEvent',
    method: 'post',
    data
  })
}

// @Tags SysHisEvent
// @Summary 删除历史事件
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysHisEvent true "删除历史事件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sysHisEvent/deleteSysHisEvent [delete]
export const deleteSysHisEvent = (params) => {
  return service({
    url: '/sysHisEvent/deleteSysHisEvent',
    method: 'delete',
    params
  })
}

// @Tags SysHisEvent
// @Summary 批量删除历史事件
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除历史事件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sysHisEvent/deleteSysHisEvent [delete]
export const deleteSysHisEventByIds = (params) => {
  return service({
    url: '/sysHisEvent/deleteSysHisEventByIds',
    method: 'delete',
    params
  })
}

// @Tags SysHisEvent
// @Summary 更新历史事件
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SysHisEvent true "更新历史事件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /sysHisEvent/updateSysHisEvent [put]
export const updateSysHisEvent = (data) => {
  return service({
    url: '/sysHisEvent/updateSysHisEvent',
    method: 'put',
    data
  })
}

// @Tags SysHisEvent
// @Summary 用id查询历史事件
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.SysHisEvent true "用id查询历史事件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /sysHisEvent/findSysHisEvent [get]
export const findSysHisEvent = (params) => {
  return service({
    url: '/sysHisEvent/findSysHisEvent',
    method: 'get',
    params
  })
}

// @Tags SysHisEvent
// @Summary 分页获取历史事件列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取历史事件列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sysHisEvent/getSysHisEventList [get]
export const getSysHisEventList = (params) => {
  return service({
    url: '/sysHisEvent/getSysHisEventList',
    method: 'get',
    params
  })
}

// @Tags SysHisEvent
// @Summary 不需要鉴权的历史事件接口
// @Accept application/json
// @Produce application/json
// @Param data query systemReq.SysHisEventSearch true "分页获取历史事件列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /sysHisEvent/getSysHisEventPublic [get]
export const getSysHisEventPublic = () => {
  return service({
    url: '/sysHisEvent/getSysHisEventPublic',
    method: 'get',
  })
}