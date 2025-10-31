import service from '@/utils/request'
// @Tags AgvcEventHis
// @Summary 创建历史事件
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.AgvcEventHis true "创建历史事件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /agvcEventHis/createAgvcEventHis [post]
export const createAgvcEventHis = (data) => {
  return service({
    url: '/agvcEventHis/createAgvcEventHis',
    method: 'post',
    data
  })
}

// @Tags AgvcEventHis
// @Summary 删除历史事件
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.AgvcEventHis true "删除历史事件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /agvcEventHis/deleteAgvcEventHis [delete]
export const deleteAgvcEventHis = (params) => {
  return service({
    url: '/agvcEventHis/deleteAgvcEventHis',
    method: 'delete',
    params
  })
}

// @Tags AgvcEventHis
// @Summary 批量删除历史事件
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除历史事件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /agvcEventHis/deleteAgvcEventHis [delete]
export const deleteAgvcEventHisByIds = (params) => {
  return service({
    url: '/agvcEventHis/deleteAgvcEventHisByIds',
    method: 'delete',
    params
  })
}

// @Tags AgvcEventHis
// @Summary 更新历史事件
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.AgvcEventHis true "更新历史事件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /agvcEventHis/updateAgvcEventHis [put]
export const updateAgvcEventHis = (data) => {
  return service({
    url: '/agvcEventHis/updateAgvcEventHis',
    method: 'put',
    data
  })
}

// @Tags AgvcEventHis
// @Summary 用id查询历史事件
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.AgvcEventHis true "用id查询历史事件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /agvcEventHis/findAgvcEventHis [get]
export const findAgvcEventHis = (params) => {
  return service({
    url: '/agvcEventHis/findAgvcEventHis',
    method: 'get',
    params
  })
}

// @Tags AgvcEventHis
// @Summary 分页获取历史事件列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取历史事件列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /agvcEventHis/getAgvcEventHisList [get]
export const getAgvcEventHisList = (params) => {
  return service({
    url: '/agvcEventHis/getAgvcEventHisList',
    method: 'get',
    params
  })
}

// @Tags AgvcEventHis
// @Summary 不需要鉴权的历史事件接口
// @Accept application/json
// @Produce application/json
// @Param data query agvcReq.AgvcEventHisSearch true "分页获取历史事件列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /agvcEventHis/getAgvcEventHisPublic [get]
export const getAgvcEventHisPublic = () => {
  return service({
    url: '/agvcEventHis/getAgvcEventHisPublic',
    method: 'get',
  })
}
