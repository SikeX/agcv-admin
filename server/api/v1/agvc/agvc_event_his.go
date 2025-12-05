package agvc

import (
    
    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
    agvcReq "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type AgvcEventHisApi struct {}



// CreateAgvcEventHis 创建历史事件
// @Tags AgvcEventHis
// @Summary 创建历史事件
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body agvc.AgvcEventHis true "创建历史事件"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /agvcEventHis/createAgvcEventHis [post]
func (agvcEventHisApi *AgvcEventHisApi) CreateAgvcEventHis(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    var agvcEventHis agvc.AgvcEventHis
    err := c.ShouldBindJSON(&agvcEventHis)
    if err != nil {
        response.FailWithMessage(err.Error(), c)
        return
    }
    err = agvcEventHisService.CreateAgvcEventHis(ctx,&agvcEventHis)
    if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
        response.FailWithMessage("创建失败:" + err.Error(), c)
        return
    }
    response.OkWithMessage("创建成功", c)
}

// DeleteAgvcEventHis 删除历史事件
// @Tags AgvcEventHis
// @Summary 删除历史事件
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body agvc.AgvcEventHis true "删除历史事件"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /agvcEventHis/deleteAgvcEventHis [delete]
func (agvcEventHisApi *AgvcEventHisApi) DeleteAgvcEventHis(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    ID := c.Query("ID")
    err := agvcEventHisService.DeleteAgvcEventHis(ctx,ID)
    if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
        response.FailWithMessage("删除失败:" + err.Error(), c)
        return
    }
    response.OkWithMessage("删除成功", c)
}

// DeleteAgvcEventHisByIds 批量删除历史事件
// @Tags AgvcEventHis
// @Summary 批量删除历史事件
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /agvcEventHis/deleteAgvcEventHisByIds [delete]
func (agvcEventHisApi *AgvcEventHisApi) DeleteAgvcEventHisByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    IDs := c.QueryArray("IDs[]")
    err := agvcEventHisService.DeleteAgvcEventHisByIds(ctx,IDs)
    if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
        response.FailWithMessage("批量删除失败:" + err.Error(), c)
        return
    }
    response.OkWithMessage("批量删除成功", c)
}

// UpdateAgvcEventHis 更新历史事件
// @Tags AgvcEventHis
// @Summary 更新历史事件
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body agvc.AgvcEventHis true "更新历史事件"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /agvcEventHis/updateAgvcEventHis [put]
func (agvcEventHisApi *AgvcEventHisApi) UpdateAgvcEventHis(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

    var agvcEventHis agvc.AgvcEventHis
    err := c.ShouldBindJSON(&agvcEventHis)
    if err != nil {
        response.FailWithMessage(err.Error(), c)
        return
    }
    err = agvcEventHisService.UpdateAgvcEventHis(ctx,agvcEventHis)
    if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
        response.FailWithMessage("更新失败:" + err.Error(), c)
        return
    }
    response.OkWithMessage("更新成功", c)
}

// FindAgvcEventHis 用id查询历史事件
// @Tags AgvcEventHis
// @Summary 用id查询历史事件
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询历史事件"
// @Success 200 {object} response.Response{data=agvc.AgvcEventHis,msg=string} "查询成功"
// @Router /agvcEventHis/findAgvcEventHis [get]
func (agvcEventHisApi *AgvcEventHisApi) FindAgvcEventHis(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    ID := c.Query("ID")
    reagvcEventHis, err := agvcEventHisService.GetAgvcEventHis(ctx,ID)
    if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
        response.FailWithMessage("查询失败:" + err.Error(), c)
        return
    }
    response.OkWithData(reagvcEventHis, c)
}
// GetAgvcEventHisList 分页获取历史事件列表
// @Tags AgvcEventHis
// @Summary 分页获取历史事件列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query agvcReq.AgvcEventHisSearch true "分页获取历史事件列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /agvcEventHis/getAgvcEventHisList [get]
func (agvcEventHisApi *AgvcEventHisApi) GetAgvcEventHisList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    var pageInfo agvcReq.AgvcEventHisSearch
    err := c.ShouldBindQuery(&pageInfo)
    if err != nil {
        response.FailWithMessage(err.Error(), c)
        return
    }
    list, total, err := agvcEventHisService.GetAgvcEventHisInfoList(ctx,pageInfo)
    if err != nil {
        global.GVA_LOG.Error("获取失败!", zap.Error(err))
        response.FailWithMessage("获取失败:" + err.Error(), c)
        return
    }
    response.OkWithDetailed(response.PageResult{
        List:     list,
        Total:    total,
        Page:     pageInfo.Page,
        PageSize: pageInfo.PageSize,
    }, "获取成功", c)
}

// GetAgvcEventHisPublic 不需要鉴权的历史事件接口
// @Tags AgvcEventHis
// @Summary 不需要鉴权的历史事件接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /agvcEventHis/getAgvcEventHisPublic [get]
func (agvcEventHisApi *AgvcEventHisApi) GetAgvcEventHisPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    agvcEventHisService.GetAgvcEventHisPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的历史事件接口信息",
    }, "获取成功", c)
}

// ReceiveEventPush 接收事件推送（无需认证）
// @Tags AgvcEventHis
// @Summary 接收事件推送
// @Accept application/json
// @Produce application/json
// @Param data body agvc.Event true "事件数据"
// @Success 200 {object} response.Response{msg=string} "接收成功"
// @Router /agvcEventHis/receiveEventPush [post]
func (agvcEventHisApi *AgvcEventHisApi) ReceiveEventPush(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    var event agvc.Event
    err := c.ShouldBindJSON(&event)
    if err != nil {
        global.GVA_LOG.Error("参数绑定失败!", zap.Error(err))
        response.FailWithMessage("参数错误:" + err.Error(), c)
        return
    }

    err = agvcEventHisService.ReceiveEvent(ctx, &event)
    if err != nil {
        global.GVA_LOG.Error("事件入库失败!", zap.Error(err))
        response.FailWithMessage("事件入库失败:" + err.Error(), c)
        return
    }

    global.GVA_LOG.Info("接收到事件推送", 
        zap.Int("psid", event.Psid),
        zap.Int("eqid", event.Eqid),
        zap.Int("eqType", event.EqType),
        zap.String("datapoint", event.Datapoint),
    )

    response.OkWithMessage("接收成功", c)
}
