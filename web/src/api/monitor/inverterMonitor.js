import service from '@/utils/request'

// @Tags InverterMonitor
// @Summary 创建逆变器监控
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.InverterMonitor true "创建逆变器监控"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /inverterMonitor/createInverterMonitor [post]
export const createInverterMonitor = (data) => {
  return service({
    url: '/inverterMonitor/createInverterMonitor',
    method: 'post',
    data
  })
}

// @Tags InverterMonitor
// @Summary 删除逆变器监控
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.InverterMonitor true "删除逆变器监控"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /inverterMonitor/deleteInverterMonitor [delete]
export const deleteInverterMonitor = (params) => {
  return service({
    url: '/inverterMonitor/deleteInverterMonitor',
    method: 'delete',
    params
  })
}

// @Tags InverterMonitor
// @Summary 批量删除逆变器监控
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除逆变器监控"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /inverterMonitor/deleteInverterMonitorByIds [delete]
export const deleteInverterMonitorByIds = (params) => {
  return service({
    url: '/inverterMonitor/deleteInverterMonitorByIds',
    method: 'delete',
    params
  })
}

// @Tags InverterMonitor
// @Summary 更新逆变器监控
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.InverterMonitor true "更新逆变器监控"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /inverterMonitor/updateInverterMonitor [put]
export const updateInverterMonitor = (data) => {
  return service({
    url: '/inverterMonitor/updateInverterMonitor',
    method: 'put',
    data
  })
}

// @Tags InverterMonitor
// @Summary 用id查询逆变器监控
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.InverterMonitor true "用id查询逆变器监控"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /inverterMonitor/findInverterMonitor [get]
export const findInverterMonitor = (params) => {
  return service({
    url: '/inverterMonitor/findInverterMonitor',
    method: 'get',
    params
  })
}

// @Tags InverterMonitor
// @Summary 分页获取逆变器监控列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取逆变器监控列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /inverterMonitor/getInverterMonitorList [get]
export const getInverterMonitorList = (params) => {
  return service({
    url: '/inverterMonitor/getInverterMonitorList',
    method: 'get',
    params
  })
}

// @Tags InverterMonitor
// @Summary 获取逆变器历史数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param inverterNo query string true "逆变器编号"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /inverterMonitor/getInverterHistory [get]
export const getInverterHistory = (params) => {
  return service({
    url: '/inverterMonitor/getInverterHistory',
    method: 'get',
    params
  })
}