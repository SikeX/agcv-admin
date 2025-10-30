import service from '@/utils/request'
// @Tags AgvcQxySetting
// @Summary 创建气象仪配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.AgvcQxySetting true "创建气象仪配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /agvcQxySetting/createAgvcQxySetting [post]
export const createAgvcQxySetting = (data) => {
  return service({
    url: '/agvcQxySetting/createAgvcQxySetting',
    method: 'post',
    data
  })
}

// @Tags AgvcQxySetting
// @Summary 删除气象仪配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.AgvcQxySetting true "删除气象仪配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /agvcQxySetting/deleteAgvcQxySetting [delete]
export const deleteAgvcQxySetting = (params) => {
  return service({
    url: '/agvcQxySetting/deleteAgvcQxySetting',
    method: 'delete',
    params
  })
}

// @Tags AgvcQxySetting
// @Summary 批量删除气象仪配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除气象仪配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /agvcQxySetting/deleteAgvcQxySetting [delete]
export const deleteAgvcQxySettingByIds = (params) => {
  return service({
    url: '/agvcQxySetting/deleteAgvcQxySettingByIds',
    method: 'delete',
    params
  })
}

// @Tags AgvcQxySetting
// @Summary 更新气象仪配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.AgvcQxySetting true "更新气象仪配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /agvcQxySetting/updateAgvcQxySetting [put]
export const updateAgvcQxySetting = (data) => {
  return service({
    url: '/agvcQxySetting/updateAgvcQxySetting',
    method: 'put',
    data
  })
}

// @Tags AgvcQxySetting
// @Summary 用id查询气象仪配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.AgvcQxySetting true "用id查询气象仪配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /agvcQxySetting/findAgvcQxySetting [get]
export const findAgvcQxySetting = (params) => {
  return service({
    url: '/agvcQxySetting/findAgvcQxySetting',
    method: 'get',
    params
  })
}

// @Tags AgvcQxySetting
// @Summary 分页获取气象仪配置列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取气象仪配置列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /agvcQxySetting/getAgvcQxySettingList [get]
export const getAgvcQxySettingList = (params) => {
  return service({
    url: '/agvcQxySetting/getAgvcQxySettingList',
    method: 'get',
    params
  })
}

// @Tags AgvcQxySetting
// @Summary 不需要鉴权的气象仪配置接口
// @Accept application/json
// @Produce application/json
// @Param data query agvcReq.AgvcQxySettingSearch true "分页获取气象仪配置列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /agvcQxySetting/getAgvcQxySettingPublic [get]
export const getAgvcQxySettingPublic = () => {
  return service({
    url: '/agvcQxySetting/getAgvcQxySettingPublic',
    method: 'get',
  })
}