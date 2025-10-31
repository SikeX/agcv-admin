import service from '@/utils/request'
// @Tags AgvcNbqSetting
// @Summary 创建逆变器配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.AgvcNbqSetting true "创建逆变器配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /agvcNbqSetting/createAgvcNbqSetting [post]
export const createAgvcNbqSetting = (data) => {
  return service({
    url: '/agvcNbqSetting/createAgvcNbqSetting',
    method: 'post',
    data
  })
}

// @Tags AgvcNbqSetting
// @Summary 删除逆变器配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.AgvcNbqSetting true "删除逆变器配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /agvcNbqSetting/deleteAgvcNbqSetting [delete]
export const deleteAgvcNbqSetting = (params) => {
  return service({
    url: '/agvcNbqSetting/deleteAgvcNbqSetting',
    method: 'delete',
    params
  })
}

// @Tags AgvcNbqSetting
// @Summary 批量删除逆变器配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除逆变器配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /agvcNbqSetting/deleteAgvcNbqSetting [delete]
export const deleteAgvcNbqSettingByIds = (params) => {
  return service({
    url: '/agvcNbqSetting/deleteAgvcNbqSettingByIds',
    method: 'delete',
    params
  })
}

// @Tags AgvcNbqSetting
// @Summary 更新逆变器配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.AgvcNbqSetting true "更新逆变器配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /agvcNbqSetting/updateAgvcNbqSetting [put]
export const updateAgvcNbqSetting = (data) => {
  return service({
    url: '/agvcNbqSetting/updateAgvcNbqSetting',
    method: 'put',
    data
  })
}

// @Tags AgvcNbqSetting
// @Summary 用id查询逆变器配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.AgvcNbqSetting true "用id查询逆变器配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /agvcNbqSetting/findAgvcNbqSetting [get]
export const findAgvcNbqSetting = (params) => {
  return service({
    url: '/agvcNbqSetting/findAgvcNbqSetting',
    method: 'get',
    params
  })
}

// @Tags AgvcNbqSetting
// 根据逆变器编号查询逆变器配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param inverterNo query string true "逆变器编号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /agvcNbqSetting/findAgvcNbqSettingByInverterNo [get]
export const findAgvcNbqSettingByInverterNo = (params) => {
  return service({
    url: '/agvcNbqSetting/findAgvcNbqSettingByInverterNo',
    method: 'get',
    params
  })
}

// @Tags AgvcNbqSetting
// @Summary 分页获取逆变器配置列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取逆变器配置列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /agvcNbqSetting/getAgvcNbqSettingList [get]
export const getAgvcNbqSettingList = (params) => {
  return service({
    url: '/agvcNbqSetting/getAgvcNbqSettingList',
    method: 'get',
    params
  })
}

// @Tags AgvcNbqSetting
// @Summary 不需要鉴权的逆变器配置接口
// @Accept application/json
// @Produce application/json
// @Param data query agvcReq.AgvcNbqSettingSearch true "分页获取逆变器配置列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /agvcNbqSetting/getAgvcNbqSettingPublic [get]
export const getAgvcNbqSettingPublic = () => {
  return service({
    url: '/agvcNbqSetting/getAgvcNbqSettingPublic',
    method: 'get',
  })
}
