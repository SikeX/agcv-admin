import service from '@/utils/request'
// @Tags AgvcQxyHis
// @Summary 创建气象仪监控
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.AgvcQxyHis true "创建气象仪监控"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /agvcQxyHis/createAgvcQxyHis [post]
export const createAgvcQxyHis = (data) => {
  return service({
    url: '/agvcQxyHis/createAgvcQxyHis',
    method: 'post',
    data
  })
}

// @Tags AgvcQxyHis
// @Summary 删除气象仪监控
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.AgvcQxyHis true "删除气象仪监控"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /agvcQxyHis/deleteAgvcQxyHis [delete]
export const deleteAgvcQxyHis = (params) => {
  return service({
    url: '/agvcQxyHis/deleteAgvcQxyHis',
    method: 'delete',
    params
  })
}

// @Tags AgvcQxyHis
// @Summary 批量删除气象仪监控
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除气象仪监控"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /agvcQxyHis/deleteAgvcQxyHis [delete]
export const deleteAgvcQxyHisByIds = (params) => {
  return service({
    url: '/agvcQxyHis/deleteAgvcQxyHisByIds',
    method: 'delete',
    params
  })
}

// @Tags AgvcQxyHis
// @Summary 更新气象仪监控
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.AgvcQxyHis true "更新气象仪监控"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /agvcQxyHis/updateAgvcQxyHis [put]
export const updateAgvcQxyHis = (data) => {
  return service({
    url: '/agvcQxyHis/updateAgvcQxyHis',
    method: 'put',
    data
  })
}

// @Tags AgvcQxyHis
// @Summary 用id查询气象仪监控
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.AgvcQxyHis true "用id查询气象仪监控"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /agvcQxyHis/findAgvcQxyHis [get]
export const findAgvcQxyHis = (params) => {
  return service({
    url: '/agvcQxyHis/findAgvcQxyHis',
    method: 'get',
    params
  })
}

// @Tags AgvcQxyHis
// @Summary 分页获取气象仪监控列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取气象仪监控列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /agvcQxyHis/getAgvcQxyHisList [get]
export const getAgvcQxyHisList = (params) => {
  return service({
    url: '/agvcQxyHis/getAgvcQxyHisList',
    method: 'get',
    params
  })
}

// @Tags AgvcQxyHis
// @Summary 不需要鉴权的气象仪监控接口
// @Accept application/json
// @Produce application/json
// @Param data query agvcReq.AgvcQxyHisSearch true "分页获取气象仪监控列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /agvcQxyHis/getAgvcQxyHisPublic [get]
export const getAgvcQxyHisPublic = () => {
  return service({
    url: '/agvcQxyHis/getAgvcQxyHisPublic',
    method: 'get',
  })
}

// @Tags AgvcQxyHis
// @Summary 获取气象仪历史数据
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param eqid query string true "设备编号"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Success 200 {object} response.Response{data=[]map[string]interface{},msg=string} "获取成功"
// @Router /agvcQxyHis/getAgvcQxyHistory [get]
export const getAgvcQxyHistory = (params) => {
  return service({
    url: '/agvcQxyHis/getAgvcQxyHistory',
    method: 'get',
    params
  })
}

// @Tags AgvcQxyHis
// @Summary 生成气象仪测试数据
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body {eqid: string, count: number} true "设备编号和数据条数"
// @Success 200 {object} response.Response{msg=string} "生成成功"
// @Router /agvcQxyHis/generateTestData [post]
export const generateAgvcQxyTestData = (data) => {
  return service({
    url: '/agvcQxyHis/generateTestData',
    method: 'post',
    data
  })
}
