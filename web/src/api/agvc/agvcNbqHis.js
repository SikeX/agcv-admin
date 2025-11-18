import service from '@/utils/request'
// @Tags AgvcNbqHis
// @Summary 创建agvcNbqHis表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.AgvcNbqHis true "创建agvcNbqHis表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /agvcNbqHis/createAgvcNbqHis [post]
export const createAgvcNbqHis = (data) => {
  return service({
    url: '/agvcNbqHis/createAgvcNbqHis',
    method: 'post',
    data
  })
}

// @Tags AgvcNbqHis
// @Summary 删除agvcNbqHis表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.AgvcNbqHis true "删除agvcNbqHis表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /agvcNbqHis/deleteAgvcNbqHis [delete]
export const deleteAgvcNbqHis = (params) => {
  return service({
    url: '/agvcNbqHis/deleteAgvcNbqHis',
    method: 'delete',
    params
  })
}

// @Tags AgvcNbqHis
// @Summary 批量删除agvcNbqHis表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除agvcNbqHis表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /agvcNbqHis/deleteAgvcNbqHis [delete]
export const deleteAgvcNbqHisByIds = (params) => {
  return service({
    url: '/agvcNbqHis/deleteAgvcNbqHisByIds',
    method: 'delete',
    params
  })
}

// @Tags AgvcNbqHis
// @Summary 更新agvcNbqHis表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.AgvcNbqHis true "更新agvcNbqHis表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /agvcNbqHis/updateAgvcNbqHis [put]
export const updateAgvcNbqHis = (data) => {
  return service({
    url: '/agvcNbqHis/updateAgvcNbqHis',
    method: 'put',
    data
  })
}

// @Tags AgvcNbqHis
// @Summary 用id查询agvcNbqHis表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.AgvcNbqHis true "用id查询agvcNbqHis表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /agvcNbqHis/findAgvcNbqHis [get]
export const findAgvcNbqHis = (params) => {
  return service({
    url: '/agvcNbqHis/findAgvcNbqHis',
    method: 'get',
    params
  })
}

// @Tags AgvcNbqHis
// @Summary 分页获取agvcNbqHis表列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取agvcNbqHis表列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /agvcNbqHis/getAgvcNbqHisList [get]
export const getAgvcNbqHisList = (params) => {
  return service({
    url: '/agvcNbqHis/getAgvcNbqHisList',
    method: 'get',
    params
  })
}

// @Tags AgvcNbqHis
// @Summary 不需要鉴权的agvcNbqHis表接口
// @Accept application/json
// @Produce application/json
// @Param data query agvcReq.AgvcNbqHisSearch true "分页获取agvcNbqHis表列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /agvcNbqHis/getAgvcNbqHisPublic [get]
export const getAgvcNbqHisPublic = () => {
  return service({
    url: '/agvcNbqHis/getAgvcNbqHisPublic',
    method: 'get',
  })
}

// @Tags AgvcNbqHis
// @Summary 获取逆变器历史数据
// @Accept application/json
// @Produce application/json
// @Param data body agvcReq.AgvcNbqHistoryRequest true "包含AgvcNbqHis结构体和时间范围"
// @Success 200 {object} response.Response{data=[]map[string]interface{},msg=string} "获取成功"
// @Router /agvcNbqHis/getAgvcNbqHistory [post]
export const getAgvcNbqHistory = (data) => {
  return service({
    url: '/agvcNbqHis/getAgvcNbqHistory',
    method: 'post',
    data
  })
}

// @Tags AgvcNbqHis
// @Summary 生成逆变器测试数据
// @Accept application/json
// @Produce application/json
// @Param eqid query string false "设备编号"
// @Param count query int false "生成数据点数量"
// @Success 200 {object} response.Response{msg=string} "生成成功"
// @Router /agvcNbqHis/generateTestData [post]
export const generateAgvcNbqTestData = (params) => {
  return service({
    url: '/agvcNbqHis/generateTestData',
    method: 'post',
    params
  })
}
