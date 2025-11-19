package agvc

import (
    "fmt"
    
    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
    agvcReq "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/request"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type AgvcBwdHisApi struct{}

// CreateAgvcBwdHis 创建agvcBwdHis表
// @Tags AgvcBwdHis
// @Summary 创建agvcBwdHis表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body agvc.AgvcBwdHis true "创建agvcBwdHis表"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /agvcBwdHis/createAgvcBwdHis [post]
func (agvcBwdHisApi *AgvcBwdHisApi) CreateAgvcBwdHis(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    var agvcBwdHis agvc.AgvcBwdHis
    err := c.ShouldBindJSON(&agvcBwdHis)
    if err != nil {
        response.FailWithMessage(err.Error(), c)
        return
    }
    err = agvcBwdHisService.CreateAgvcBwdHis(ctx, &agvcBwdHis)
    if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
        response.FailWithMessage("创建失败:"+err.Error(), c)
        return
    }
    response.OkWithMessage("创建成功", c)
}

// DeleteAgvcBwdHis 删除agvcBwdHis表
// @Tags AgvcBwdHis
// @Summary 删除agvcBwdHis表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body agvc.AgvcBwdHis true "删除agvcBwdHis表"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /agvcBwdHis/deleteAgvcBwdHis [delete]
func (agvcBwdHisApi *AgvcBwdHisApi) DeleteAgvcBwdHis(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    ID := c.Query("ID")
    err := agvcBwdHisService.DeleteAgvcBwdHis(ctx, ID)
    if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
        response.FailWithMessage("删除失败:"+err.Error(), c)
        return
    }
    response.OkWithMessage("删除成功", c)
}

// DeleteAgvcBwdHisByIds 批量删除agvcBwdHis表
// @Tags AgvcBwdHis
// @Summary 批量删除agvcBwdHis表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /agvcBwdHis/deleteAgvcBwdHisByIds [delete]
func (agvcBwdHisApi *AgvcBwdHisApi) DeleteAgvcBwdHisByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    IDs := c.QueryArray("IDs[]")
    err := agvcBwdHisService.DeleteAgvcBwdHisByIds(ctx, IDs)
    if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
        response.FailWithMessage("批量删除失败:"+err.Error(), c)
        return
    }
    response.OkWithMessage("批量删除成功", c)
}

// UpdateAgvcBwdHis 更新agvcBwdHis表
// @Tags AgvcBwdHis
// @Summary 更新agvcBwdHis表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body agvc.AgvcBwdHis true "更新agvcBwdHis表"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /agvcBwdHis/updateAgvcBwdHis [put]
// func (agvcBwdHisApi *AgvcBwdHisApi) UpdateAgvcBwdHis(c *gin.Context) {
//     // 从ctx获取标准context进行业务行为
//     ctx := c.Request.Context()

//     var agvcBwdHis agvc.AgvcBwdHis
//     err := c.ShouldBindJSON(&agvcBwdHis)
//     if err != nil {
//         response.FailWithMessage(err.Error(), c)
//         return
//     }
//     err = agvcBwdHisService.UpdateAgvcBwdHis(ctx, agvcBwdHis)
//     if err != nil {
//         global.GVA_LOG.Error("更新失败!", zap.Error(err))
//         response.FailWithMessage("更新失败:"+err.Error(), c)
//         return
//     }
//     response.OkWithMessage("更新成功", c)
// }

// FindAgvcBwdHis 用id查询agvcBwdHis表
// @Tags AgvcBwdHis
// @Summary 用id查询agvcBwdHis表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询agvcBwdHis表"
// @Success 200 {object} response.Response{data=agvc.AgvcBwdHis,msg=string} "查询成功"
// @Router /agvcBwdHis/findAgvcBwdHis [get]
func (agvcBwdHisApi *AgvcBwdHisApi) FindAgvcBwdHis(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    ID := c.Query("ID")
    reagvcBwdHis, err := agvcBwdHisService.GetAgvcBwdHis(ctx, ID)
    if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
        response.FailWithMessage("查询失败:"+err.Error(), c)
        return
    }
    response.OkWithData(reagvcBwdHis, c)
}

// GetAgvcBwdHisList 分页获取agvcBwdHis表列表
// @Tags AgvcBwdHis
// @Summary 分页获取agvcBwdHis表列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query agvcReq.AgvcBwdHisSearch true "分页获取agvcBwdHis表列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /agvcBwdHis/getAgvcBwdHisList [get]
func (agvcBwdHisApi *AgvcBwdHisApi) GetAgvcBwdHisList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    var pageInfo agvcReq.AgvcBwdHisSearch
    err := c.ShouldBindQuery(&pageInfo)
    if err != nil {
        response.FailWithMessage(err.Error(), c)
        return
    }
    list, total, err := agvcBwdHisService.GetAgvcBwdHisInfoList(ctx, pageInfo)
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

// GetAgvcBwdHisPublic 不需要鉴权的agvcBwdHis表接口
// @Tags AgvcBwdHis
// @Summary 不需要鉴权的agvcBwdHis表接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /agvcBwdHis/getAgvcBwdHisPublic [get]
func (agvcBwdHisApi *AgvcBwdHisApi) GetAgvcBwdHisPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    agvcBwdHisService.GetAgvcBwdHisPublic(ctx)
    response.OkWithDetailed(gin.H{
        "info": "不需要鉴权的agvcBwdHis表接口信息",
    }, "获取成功", c)
}

// GetAgvcBwdHistory 获取并网点历史数据
// @Tags AgvcBwdHis
// @Summary 获取并网点历史数据
// @Accept application/json
// @Produce application/json
// @Param eqid query string true "设备编号"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Success 200 {object} response.Response{data=[]map[string]interface{},msg=string} "获取成功"
// @Router /agvcBwdHis/getAgvcBwdHistory [get]
func (agvcBwdHisApi *AgvcBwdHisApi) GetAgvcBwdHistory(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    eqid := c.Query("eqid")
    startTime := c.Query("startTime")
    endTime := c.Query("endTime")

    if eqid == "" || startTime == "" || endTime == "" {
        response.FailWithMessage("参数不完整", c)
        return
    }

    //如果传入的时间是日期转换为2025-10-06T07:06:10.245Z格式
    if len(startTime) == 10 {
        startTime += "T00:00:00Z"
        endTime += "T23:59:59Z"
    }

    historyData, err := agvcBwdHisService.GetAgvcBwdHistory(ctx, eqid, startTime, endTime)
    if err != nil {
        global.GVA_LOG.Error("获取历史数据失败!", zap.Error(err))
        response.FailWithMessage("获取历史数据失败:"+err.Error(), c)
        return
    }
    response.OkWithData(historyData, c)
}

// GenerateAgvcBwdTestData 生成并网点测试数据
// @Tags AgvcBwdHis
// @Summary 生成并网点测试数据
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param eqid query string false "设备编号"
// @Param count query int false "生成数据点数量，默认100"
// @Success 200 {object} response.Response{msg=string} "生成成功"
// @Router /agvcBwdHis/generateTestData [post]
func (agvcBwdHisApi *AgvcBwdHisApi) GenerateAgvcBwdTestData(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    eqid := c.Query("eqid")
    count := 100
    if c.Query("count") != "" {
        c.ShouldBindQuery(&count)
    }

    err := agvcBwdHisService.GenerateAgvcBwdTestData(ctx, eqid, count)
    if err != nil {
        global.GVA_LOG.Error("生成测试数据失败!", zap.Error(err))
        response.FailWithMessage("生成测试数据失败:"+err.Error(), c)
        return
    }
    response.OkWithMessage("生成测试数据成功", c)
}

// UpdatePlanCurves 更新计划曲线
// @Tags AgvcBwdHis
// @Summary 更新计划曲线
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body agvcReq.PlanCurvesRequest true "计划曲线数据"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /agvcBwdHis/updatePlanCurves [put]
// func (agvcBwdHisApi *AgvcBwdHisApi) UpdatePlanCurves(c *gin.Context) {
//     ctx := c.Request.Context()

//     var req agvcReq.PlanCurvesRequest
//     err := c.ShouldBindJSON(&req)
//     if err != nil {
//         response.FailWithMessage(err.Error(), c)
//         return
//     }

//     err = agvcBwdHisService.UpdatePlanCurves(ctx, req)
//     if err != nil {
//         global.GVA_LOG.Error("更新计划曲线失败!", zap.Error(err))
//         response.FailWithMessage("更新计划曲线失败:"+err.Error(), c)
//         return
//     }
//     response.OkWithMessage("更新计划曲线成功", c)
// }

// UpdateAgcParameters 更新AGC参数设置
// @Tags AgvcBwdHis
// @Summary 更新AGC参数设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body agvcReq.AgcParametersRequest true "AGC参数数据"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /agvcBwdHis/updateAgcParameters [put]
// func (agvcBwdHisApi *AgvcBwdHisApi) UpdateAgcParameters(c *gin.Context) {
//     ctx := c.Request.Context()

//     var req agvcReq.AgcParametersRequest
//     err := c.ShouldBindJSON(&req)
//     if err != nil {
//         response.FailWithMessage(err.Error(), c)
//         return
//     }

//     err = agvcBwdHisService.UpdateAgcParameters(ctx, req)
//     if err != nil {
//         global.GVA_LOG.Error("更新AGC参数失败!", zap.Error(err))
//         response.FailWithMessage("更新AGC参数失败:"+err.Error(), c)
//         return
//     }
//     response.OkWithMessage("更新AGC参数成功", c)
// }

// UpdateAvcParameters 更新AVC参数设置
// @Tags AgvcBwdHis
// @Summary 更新AVC参数设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body agvcReq.AvcParametersRequest true "AVC参数数据"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /agvcBwdHis/updateAvcParameters [put]
// func (agvcBwdHisApi *AgvcBwdHisApi) UpdateAvcParameters(c *gin.Context) {
//     ctx := c.Request.Context()

//     var req agvcReq.AvcParametersRequest
//     err := c.ShouldBindJSON(&req)
//     if err != nil {
//         response.FailWithMessage(err.Error(), c)
//         return
//     }

//     err = agvcBwdHisService.UpdateAvcParameters(ctx, req)
//     if err != nil {
//         global.GVA_LOG.Error("更新AVC参数失败!", zap.Error(err))
//         response.FailWithMessage("更新AVC参数失败:"+err.Error(), c)
//         return
//     }
//     response.OkWithMessage("更新AVC参数成功", c)
// }

// GetAgcParameters 获取AGC参数设置
// @Tags AgvcBwdHis
// @Summary 获取AGC参数设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param number query string true "并网点编号"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /agvcBwdHis/getAgcParameters [get]
// func (agvcBwdHisApi *AgvcBwdHisApi) GetAgcParameters(c *gin.Context) {
//     ctx := c.Request.Context()

//     number := c.Query("number")
//     if number == "" {
//         response.FailWithMessage("参数不完整", c)
//         return
//     }

//     data, err := agvcBwdHisService.GetAgcParameters(ctx, number)
//     if err != nil {
//         global.GVA_LOG.Error("获取AGC参数失败!", zap.Error(err))
//         response.FailWithMessage("获取AGC参数失败:"+err.Error(), c)
//         return
//     }
//     response.OkWithData(data, c)
// }

// GetAvcParameters 获取AVC参数设置
// @Tags AgvcBwdHis
// @Summary 获取AVC参数设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param number query string true "并网点编号"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /agvcBwdHis/getAvcParameters [get]
// func (agvcBwdHisApi *AgvcBwdHisApi) GetAvcParameters(c *gin.Context) {
//     ctx := c.Request.Context()

//     number := c.Query("number")
//     if number == "" {
//         response.FailWithMessage("参数不完整", c)
//         return
//     }

//     data, err := agvcBwdHisService.GetAvcParameters(ctx, number)
//     if err != nil {
//         global.GVA_LOG.Error("获取AVC参数失败!", zap.Error(err))
//         response.FailWithMessage("获取AVC参数失败:"+err.Error(), c)
//         return
//     }
//     response.OkWithData(data, c)
// }

// GetPlanCurves 获取计划曲线
// @Tags AgvcBwdHis
// @Summary 获取计划曲线
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param number query string true "并网点编号"
// @Param curveType query string false "曲线类型: agc/avc"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /agvcBwdHis/getPlanCurves [get]
// func (agvcBwdHisApi *AgvcBwdHisApi) GetPlanCurves(c *gin.Context) {
//     ctx := c.Request.Context()

//     number := c.Query("number")
//     if number == "" {
//         response.FailWithMessage("参数不完整", c)
//         return
//     }

//     curveType := c.Query("curveType") // agc 或 avc

//     data, err := agvcBwdHisService.GetPlanCurves(ctx, number, curveType)
//     if err != nil {
//         global.GVA_LOG.Error("获取计划曲线失败!", zap.Error(err))
//         response.FailWithMessage("获取计划曲线失败:"+err.Error(), c)
//         return
//     }
//     response.OkWithData(data, c)
// }

// UpdateAgcStatus 更新AGC状态
// @Tags AgvcBwdHis
// @Summary 更新AGC状态
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body agvcReq.AgcStatusRequest true "AGC状态数据"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /agvcBwdHis/updateAgcStatus [put]
// func (agvcBwdHisApi *AgvcBwdHisApi) UpdateAgcStatus(c *gin.Context) {
//     ctx := c.Request.Context()

//     var req agvcReq.AgcStatusRequest
//     err := c.ShouldBindJSON(&req)
//     if err != nil {
//         response.FailWithMessage(err.Error(), c)
//         return
//     }

//     err = agvcBwdHisService.UpdateAgcStatus(ctx, req)
//     if err != nil {
//         global.GVA_LOG.Error("更新AGC状态失败!", zap.Error(err))
//         response.FailWithMessage("更新AGC状态失败:"+err.Error(), c)
//         return
//     }
//     response.OkWithMessage("更新AGC状态成功", c)
// }

// GetAgcStatus 获取AGC状态（从InfluxDB获取实时数据）
// @Tags AgvcBwdHis
// @Summary 获取AGC状态
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param number query string true "并网点编号"
// @Param psid query int false "电站编号，默认为1"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /agvcBwdHis/getAgcStatus [get]
func (agvcBwdHisApi *AgvcBwdHisApi) GetAgcStatus(c *gin.Context) {
    ctx := c.Request.Context()

    number := c.Query("number")
    if number == "" {
        response.FailWithMessage("参数不完整", c)
        return
    }

    // 获取电站编号，默认为1
    psid := 1
    if c.Query("psid") != "" {
        fmt.Sscanf(c.Query("psid"), "%d", &psid)
    }

    // 从InfluxDB获取AGC实时数据
    data, err := agvcAgcHisService.GetAgcRealData(ctx, psid, number)
    if err != nil {
        global.GVA_LOG.Error("获取AGC状态失败!", zap.Error(err))
        response.FailWithMessage("获取AGC状态失败:"+err.Error(), c)
        return
    }
    response.OkWithData(data, c)
}

// UpdateAvcStatus 更新AVC状态
// @Tags AgvcBwdHis
// @Summary 更新AVC状态
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body agvcReq.AvcStatusRequest true "AVC状态数据"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /agvcBwdHis/updateAvcStatus [put]
// func (agvcBwdHisApi *AgvcBwdHisApi) UpdateAvcStatus(c *gin.Context) {
//     ctx := c.Request.Context()

//     var req agvcReq.AvcStatusRequest
//     err := c.ShouldBindJSON(&req)
//     if err != nil {
//         response.FailWithMessage(err.Error(), c)
//         return
//     }

//     err = agvcBwdHisService.UpdateAvcStatus(ctx, req)
//     if err != nil {
//         global.GVA_LOG.Error("更新AVC状态失败!", zap.Error(err))
//         response.FailWithMessage("更新AVC状态失败:"+err.Error(), c)
//         return
//     }
//     response.OkWithMessage("更新AVC状态成功", c)
// }

// GetAvcStatus 获取AVC状态（从InfluxDB获取实时数据）
// @Tags AgvcBwdHis
// @Summary 获取AVC状态
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param number query string true "并网点编号"
// @Param psid query int false "电站编号，默认为1"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /agvcBwdHis/getAvcStatus [get]
func (agvcBwdHisApi *AgvcBwdHisApi) GetAvcStatus(c *gin.Context) {
    ctx := c.Request.Context()

    number := c.Query("number")
    if number == "" {
        response.FailWithMessage("参数不完整", c)
        return
    }

    // 获取电站编号，默认为1
    psid := 1
    if c.Query("psid") != "" {
        fmt.Sscanf(c.Query("psid"), "%d", &psid)
    }

    // 从InfluxDB获取AVC实时数据
    data, err := agvcAvcHisService.GetAvcRealData(ctx, psid, number)
    if err != nil {
        global.GVA_LOG.Error("获取AVC状态失败!", zap.Error(err))
        response.FailWithMessage("获取AVC状态失败:"+err.Error(), c)
        return
    }
    response.OkWithData(data, c)
}

// GetBwdRealtimeData 获取并网点实时数据
// @Tags AgvcBwdHis
// @Summary 获取并网点实时数据
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param number query string true "并网点编号"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /agvcBwdHis/getBwdRealtimeData [get]
func (agvcBwdHisApi *AgvcBwdHisApi) GetBwdRealtimeData(c *gin.Context) {
    ctx := c.Request.Context()

    number := c.Query("number")
    if number == "" {
        response.FailWithMessage("参数不完整", c)
        return
    }

    realtimeData, err := agvcBwdHisService.GetBwdRealtimeData(ctx, number)
    if err != nil {
        global.GVA_LOG.Error("获取并网点实时数据失败!", zap.Error(err))
        response.FailWithMessage("获取并网点实时数据失败:"+err.Error(), c)
        return
    }
    response.OkWithData(realtimeData, c)
}
