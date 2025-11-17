import service from '@/utils/request'
// @Tags AgvcBwdHis
// @Summary 创建agvcBwdHis表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.AgvcBwdHis true "创建agvcBwdHis表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /agvcBwdHis/createAgvcBwdHis [post]
export const createAgvcBwdHis = (data) => {
  return service({
    url: '/agvcBwdHis/createAgvcBwdHis',
    method: 'post',
    data
  })
}

// @Tags AgvcBwdHis
// @Summary 删除agvcBwdHis表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.AgvcBwdHis true "删除agvcBwdHis表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /agvcBwdHis/deleteAgvcBwdHis [delete]
export const deleteAgvcBwdHis = (params) => {
  return service({
    url: '/agvcBwdHis/deleteAgvcBwdHis',
    method: 'delete',
    params
  })
}

// @Tags AgvcBwdHis
// @Summary 批量删除agvcBwdHis表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除agvcBwdHis表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /agvcBwdHis/deleteAgvcBwdHis [delete]
export const deleteAgvcBwdHisByIds = (params) => {
  return service({
    url: '/agvcBwdHis/deleteAgvcBwdHisByIds',
    method: 'delete',
    params
  })
}

// @Tags AgvcBwdHis
// @Summary 更新agvcBwdHis表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.AgvcBwdHis true "更新agvcBwdHis表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /agvcBwdHis/updateAgvcBwdHis [put]
export const updateAgvcBwdHis = (data) => {
  return service({
    url: '/agvcBwdHis/updateAgvcBwdHis',
    method: 'put',
    data
  })
}

// @Tags AgvcBwdHis
// @Summary 用id查询agvcBwdHis表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.AgvcBwdHis true "用id查询agvcBwdHis表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /agvcBwdHis/findAgvcBwdHis [get]
export const findAgvcBwdHis = (params) => {
  return service({
    url: '/agvcBwdHis/findAgvcBwdHis',
    method: 'get',
    params
  })
}

// @Tags AgvcBwdHis
// @Summary 分页获取agvcBwdHis表列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取agvcBwdHis表列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /agvcBwdHis/getAgvcBwdHisList [get]
export const getAgvcBwdHisList = (params) => {
  return service({
    url: '/agvcBwdHis/getAgvcBwdHisList',
    method: 'get',
    params
  })
}

// @Tags AgvcBwdHis
// @Summary 不需要鉴权的agvcBwdHis表接口
// @Accept application/json
// @Produce application/json
// @Param data query agvcReq.AgvcBwdHisSearch true "分页获取agvcBwdHis表列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /agvcBwdHis/getAgvcBwdHisPublic [get]
export const getAgvcBwdHisPublic = () => {
  return service({
    url: '/agvcBwdHis/getAgvcBwdHisPublic',
    method: 'get'
  })
}

// @Tags AgvcBwdHis
// @Summary 获取并网点历史数据
// @Accept application/json
// @Produce application/json
// @Param eqid query string true "设备编号"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Success 200 {object} response.Response{data=[]map[string]interface{},msg=string} "获取成功"
// @Router /agvcBwdHis/getAgvcBwdHistory [get]
export const getAgvcBwdHistory = (params) => {
  return service({
    url: '/agvcBwdHis/getAgvcBwdHistory',
    method: 'get',
    donNotShowLoading: true,
    params
  })
}

// @Tags AgvcBwdHis
// @Summary 生成并网点测试数据
// @Accept application/json
// @Produce application/json
// @Param eqid query string false "设备编号"
// @Param count query int false "生成数据点数量"
// @Success 200 {object} response.Response{msg=string} "生成成功"
// @Router /agvcBwdHis/generateTestData [post]
export const generateAgvcBwdTestData = (params) => {
  return service({
    url: '/agvcBwdHis/generateTestData',
    method: 'post',
    params
  })
}

// @Tags AgvcBwdHis
// @Summary 更新AGC参数设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body object true "AGC参数"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /agvcBwdHis/updateAgcParameters [put]
export const updateAgcParameters = (data) => {
  return service({
    url: '/agvcBwdHis/updateAgcParameters',
    method: 'put',
    data
  })
}

// @Tags AgvcBwdHis
// @Summary 更新AVC参数设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body object true "AVC参数"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /agvcBwdHis/updateAvcParameters [put]
export const updateAvcParameters = (data) => {
  return service({
    url: '/agvcBwdHis/updateAvcParameters',
    method: 'put',
    data
  })
}

// @Tags AgvcBwdHis
// @Summary 更新计划曲线
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body object true "计划曲线数据"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /agvcBwdHis/updatePlanCurves [put]
export const updatePlanCurves = (data) => {
  return service({
    url: '/agvcBwdHis/updatePlanCurves',
    method: 'put',
    data
  })
}

// @Tags AgvcBwdHis
// @Summary 获取AGC参数设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param number query string true "并网点编号"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /agvcBwdHis/getAgcParameters [get]
export const getAgcParameters = (params) => {
  return service({
    url: '/agvcBwdHis/getAgcParameters',
    method: 'get',
    donNotShowLoading: true,
    params
  })
}

// @Tags AgvcBwdHis
// @Summary 获取AVC参数设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param number query string true "并网点编号"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /agvcBwdHis/getAvcParameters [get]
export const getAvcParameters = (params) => {
  return service({
    url: '/agvcBwdHis/getAvcParameters',
    method: 'get',
    params
  })
}

// @Tags AgvcBwdHis
// @Summary 获取计划曲线
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param number query string true "并网点编号"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /agvcBwdHis/getPlanCurves [get]
export const getPlanCurves = (params) => {
  return service({
    url: '/agvcBwdHis/getPlanCurves',
    method: 'get',
    params
  })
}

// @Tags AgvcBwdHis
// @Summary 更新AGC状态
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body object true "AGC状态数据"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /agvcBwdHis/updateAgcStatus [put]
export const updateAgcStatus = (data) => {
  return service({
    url: '/agvcBwdHis/updateAgcStatus',
    method: 'put',
    donNotShowLoading: true,
    data
  })
}

// @Tags AgvcBwdHis
// @Summary 获取AGC状态
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param number query string true "并网点编号"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /agvcBwdHis/getAgcStatus [get]
export const getAgcStatus = (params) => {
  return service({
    url: '/agvcBwdHis/getAgcStatus',
    method: 'get',
    donNotShowLoading: true,
    params
  })
}

// @Tags AgvcBwdHis
// @Summary 更新AVC状态
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body object true "AVC状态数据"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /agvcBwdHis/updateAvcStatus [put]
export const updateAvcStatus = (data) => {
  return service({
    url: '/agvcBwdHis/updateAvcStatus',
    method: 'put',
    donNotShowLoading: true,
    data
  })
}

// @Tags AgvcBwdHis
// @Summary 获取AVC状态
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param number query string true "并网点编号"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /agvcBwdHis/getAvcStatus [get]
export const getAvcStatus = (params) => {
  return service({
    url: '/agvcBwdHis/getAvcStatus',
    method: 'get',
    donNotShowLoading: true,
    params
  })
}

// @Tags AgvcBwdHis
// @Summary 获取并网点实时数据
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param number query string true "并网点编号"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /agvcBwdHis/getBwdRealtimeData [get]
export const getBwdRealtimeData = (params) => {
  return service({
    url: '/agvcBwdHis/getBwdRealtimeData',
    method: 'get',
    donNotShowLoading: true,
    params
  })
}
