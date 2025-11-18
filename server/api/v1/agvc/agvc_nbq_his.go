package agvc

import (
    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
    agvcReq "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/request"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type AgvcNbqHisApi struct{}

// CreateAgvcNbqHis 创建agvcNbqHis表
// @Tags AgvcNbqHis
// @Summary 创建agvcNbqHis表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body agvc.AgvcNbqHis true "创建agvcNbqHis表"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /agvcNbqHis/createAgvcNbqHis [post]
func (agvcNbqHisApi *AgvcNbqHisApi) CreateAgvcNbqHis(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    var agvcNbqHis agvc.AgvcNbqHis
    err := c.ShouldBindJSON(&agvcNbqHis)
    if err != nil {
        response.FailWithMessage(err.Error(), c)
        return
    }
    err = agvcNbqHisService.CreateAgvcNbqHis(ctx, &agvcNbqHis)
    if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
        response.FailWithMessage("创建失败:"+err.Error(), c)
        return
    }
    response.OkWithMessage("创建成功", c)
}

// DeleteAgvcNbqHis 删除agvcNbqHis表
// @Tags AgvcNbqHis
// @Summary 删除agvcNbqHis表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body agvc.AgvcNbqHis true "删除agvcNbqHis表"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /agvcNbqHis/deleteAgvcNbqHis [delete]
func (agvcNbqHisApi *AgvcNbqHisApi) DeleteAgvcNbqHis(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    ID := c.Query("ID")
    err := agvcNbqHisService.DeleteAgvcNbqHis(ctx, ID)
    if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
        response.FailWithMessage("删除失败:"+err.Error(), c)
        return
    }
    response.OkWithMessage("删除成功", c)
}

// DeleteAgvcNbqHisByIds 批量删除agvcNbqHis表
// @Tags AgvcNbqHis
// @Summary 批量删除agvcNbqHis表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /agvcNbqHis/deleteAgvcNbqHisByIds [delete]
func (agvcNbqHisApi *AgvcNbqHisApi) DeleteAgvcNbqHisByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    IDs := c.QueryArray("IDs[]")
    err := agvcNbqHisService.DeleteAgvcNbqHisByIds(ctx, IDs)
    if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
        response.FailWithMessage("批量删除失败:"+err.Error(), c)
        return
    }
    response.OkWithMessage("批量删除成功", c)
}

// UpdateAgvcNbqHis 更新agvcNbqHis表
// @Tags AgvcNbqHis
// @Summary 更新agvcNbqHis表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body agvc.AgvcNbqHis true "更新agvcNbqHis表"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /agvcNbqHis/updateAgvcNbqHis [put]
func (agvcNbqHisApi *AgvcNbqHisApi) UpdateAgvcNbqHis(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

    var agvcNbqHis agvc.AgvcNbqHis
    err := c.ShouldBindJSON(&agvcNbqHis)
    if err != nil {
        response.FailWithMessage(err.Error(), c)
        return
    }
    err = agvcNbqHisService.UpdateAgvcNbqHis(ctx, agvcNbqHis)
    if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
        response.FailWithMessage("更新失败:"+err.Error(), c)
        return
    }
    response.OkWithMessage("更新成功", c)
}

// FindAgvcNbqHis 用id查询agvcNbqHis表
// @Tags AgvcNbqHis
// @Summary 用id查询agvcNbqHis表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询agvcNbqHis表"
// @Success 200 {object} response.Response{data=agvc.AgvcNbqHis,msg=string} "查询成功"
// @Router /agvcNbqHis/findAgvcNbqHis [get]
func (agvcNbqHisApi *AgvcNbqHisApi) FindAgvcNbqHis(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    ID := c.Query("ID")
    reagvcNbqHis, err := agvcNbqHisService.GetAgvcNbqHis(ctx, ID)
    if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
        response.FailWithMessage("查询失败:"+err.Error(), c)
        return
    }
    response.OkWithData(reagvcNbqHis, c)
}

// GetAgvcNbqHisList 分页获取agvcNbqHis表列表
// @Tags AgvcNbqHis
// @Summary 分页获取agvcNbqHis表列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query agvcReq.AgvcNbqHisSearch true "分页获取agvcNbqHis表列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /agvcNbqHis/getAgvcNbqHisList [get]
func (agvcNbqHisApi *AgvcNbqHisApi) GetAgvcNbqHisList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    var pageInfo agvcReq.AgvcNbqHisSearch
    err := c.ShouldBindQuery(&pageInfo)
    if err != nil {
        response.FailWithMessage(err.Error(), c)
        return
    }
    list, total, err := agvcNbqHisService.GetAgvcNbqHisInfoList(ctx, pageInfo)
    if err != nil {
        global.GVA_LOG.Error("获取失败!", zap.Error(err))
        response.FailWithMessage("获取失败:"+err.Error(), c)
        return
    }
    response.OkWithDetailed(response.PageResult{
        List:     list,
        Total:    total,
        Page:     pageInfo.Page,
        PageSize: pageInfo.PageSize,
    }, "获取成功", c)
}

// GetAgvcNbqHisPublic 不需要鉴权的agvcNbqHis表接口
// @Tags AgvcNbqHis
// @Summary 不需要鉴权的agvcNbqHis表接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /agvcNbqHis/getAgvcNbqHisPublic [get]
func (agvcNbqHisApi *AgvcNbqHisApi) GetAgvcNbqHisPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    agvcNbqHisService.GetAgvcNbqHisPublic(ctx)
    response.OkWithDetailed(gin.H{
        "info": "不需要鉴权的agvcNbqHis表接口信息",
    }, "获取成功", c)
}

// GetAgvcNbqHistory 获取逆变器历史数据
// @Tags AgvcNbqHis
// @Summary 获取逆变器历史数据
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body agvcReq.AgvcNbqHistoryRequest true "包含AgvcNbqHis结构体和时间范围"
// @Success 200 {object} response.Response{data=[]agvc.AgvcNbqHis,msg=string} "获取成功"
// @Router /agvcNbqHis/getAgvcNbqHistory [post]
func (agvcNbqHisApi *AgvcNbqHisApi) GetAgvcNbqHistory(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    var req agvcReq.AgvcNbqHistoryRequest
    err := c.ShouldBindJSON(&req)
    if err != nil {
        response.FailWithMessage("参数绑定失败:"+err.Error(), c)
        return
    }

    // 检查必要的参数
    if req.StartTime == "" || req.EndTime == "" {
        response.FailWithMessage("开始时间和结束时间不能为空", c)
        return
    }

    // 如果传入的时间是日期格式，转换为RFC3339格式
    if len(req.StartTime) == 10 {
        req.StartTime += "T00:00:00Z"
        req.EndTime += "T23:59:59Z"
    }

    historyData, err := agvcNbqHisService.GetAgvcNbqHistory(ctx, req.AgvcNbqHis, req.StartTime, req.EndTime)
    if err != nil {
        global.GVA_LOG.Error("获取历史数据失败!", zap.Error(err))
        response.FailWithMessage("获取历史数据失败:"+err.Error(), c)
        return
    }
    response.OkWithData(historyData, c)
}

// GetNbqRealData 获取逆变器实时数据
// @Tags AgvcNbqHis
// @Summary 获取逆变器实时数据（最新一条数据）
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param psid query int true "电站编号"
// @Param inverterNo query int true "逆变器编号"
// @Success 200 {object} response.Response{data=agvc.AgvcNbqHis,msg=string} "获取成功"
// @Router /agvcNbqHis/getNbqRealData [get]
func (agvcNbqHisApi *AgvcNbqHisApi) GetNbqRealData(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    var req agvcReq.AgvcNbqRealDataRequest
    err := c.ShouldBindQuery(&req)
    if err != nil {
        response.FailWithMessage("参数绑定失败:"+err.Error(), c)
        return
    }

    realData, err := agvcNbqHisService.GetNbqRealData(ctx, req.Psid, req.InverterNo)
    if err != nil {
        global.GVA_LOG.Error("获取实时数据失败!", zap.Error(err))
        response.FailWithMessage("获取实时数据失败:"+err.Error(), c)
        return
    }
    response.OkWithData(realData, c)
}

// GenerateAgvcNbqTestData 生成逆变器测试数据
// @Tags AgvcNbqHis
// @Summary 生成逆变器测试数据
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param eqid query string false "设备编号"
// @Param count query int false "生成数据点数量，默认100"
// @Success 200 {object} response.Response{msg=string} "生成成功"
// @Router /agvcNbqHis/generateTestData [post]
//func (agvcNbqHisApi *AgvcNbqHisApi) GenerateAgvcNbqTestData(c *gin.Context) {
//    // 创建业务用Context
//    ctx := c.Request.Context()
//
//    eqid := c.Query("eqid")
//    count := 100
//    if c.Query("count") != "" {
//        c.ShouldBindQuery(&count)
//    }
//
//    err := agvcNbqHisService.GenerateAgvcNbqTestData(ctx, eqid, count)
//    if err != nil {
//        global.GVA_LOG.Error("生成测试数据失败!", zap.Error(err))
//        response.FailWithMessage("生成测试数据失败:"+err.Error(), c)
//        return
//    }
//    response.OkWithMessage("生成测试数据成功", c)
//}
