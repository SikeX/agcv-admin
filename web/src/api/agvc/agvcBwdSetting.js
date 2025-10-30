import service from '@/utils/request'
// @Tags AgvcBwdSetting
// @Summary 创建并网点配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.AgvcBwdSetting true "创建并网点配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /agvcBwdSetting/createAgvcBwdSetting [post]
export const createAgvcBwdSetting = (data) => {
  return service({
    url: '/agvcBwdSetting/createAgvcBwdSetting',
    method: 'post',
    data
  })
}

// @Tags AgvcBwdSetting
// @Summary 删除并网点配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.AgvcBwdSetting true "删除并网点配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /agvcBwdSetting/deleteAgvcBwdSetting [delete]
export const deleteAgvcBwdSetting = (params) => {
  return service({
    url: '/agvcBwdSetting/deleteAgvcBwdSetting',
    method: 'delete',
    params
  })
}

// @Tags AgvcBwdSetting
// @Summary 批量删除并网点配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除并网点配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /agvcBwdSetting/deleteAgvcBwdSetting [delete]
export const deleteAgvcBwdSettingByIds = (params) => {
  return service({
    url: '/agvcBwdSetting/deleteAgvcBwdSettingByIds',
    method: 'delete',
    params
  })
}

// @Tags AgvcBwdSetting
// @Summary 更新并网点配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.AgvcBwdSetting true "更新并网点配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /agvcBwdSetting/updateAgvcBwdSetting [put]
export const updateAgvcBwdSetting = (data) => {
  return service({
    url: '/agvcBwdSetting/updateAgvcBwdSetting',
    method: 'put',
    data
  })
}

// @Tags AgvcBwdSetting
// @Summary 用id查询并网点配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.AgvcBwdSetting true "用id查询并网点配置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /agvcBwdSetting/findAgvcBwdSetting [get]
export const findAgvcBwdSetting = (params) => {
  return service({
    url: '/agvcBwdSetting/findAgvcBwdSetting',
    method: 'get',
    params
  })
}

// @Tags AgvcBwdSetting
// @Summary 分页获取并网点配置列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取并网点配置列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /agvcBwdSetting/getAgvcBwdSettingList [get]
export const getAgvcBwdSettingList = (params) => {
  return service({
    url: '/agvcBwdSetting/getAgvcBwdSettingList',
    method: 'get',
    params
  })
}

// @Tags AgvcBwdSetting
// @Summary 不需要鉴权的并网点配置接口
// @Accept application/json
// @Produce application/json
// @Param data query agvcReq.AgvcBwdSettingSearch true "分页获取并网点配置列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /agvcBwdSetting/getAgvcBwdSettingPublic [get]
export const getAgvcBwdSettingPublic = () => {
  return service({
    url: '/agvcBwdSetting/getAgvcBwdSettingPublic',
    method: 'get',
  })
}
