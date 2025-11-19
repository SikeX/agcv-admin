package agvc

import (
    "fmt"
    "net/http"
    "sync"
    "time"

    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
    agvcReq "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/request"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/utils"
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

    var agvcNbqHis agvc.AgvcNbq
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

    var agvcNbqHis agvc.AgvcNbq
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
// func (agvcNbqHisApi *AgvcNbqHisApi) GetAgvcNbqHistory(c *gin.Context) {
//     // 创建业务用Context
//     ctx := c.Request.Context()

//     var req agvcReq.AgvcNbqHistoryRequest
//     err := c.ShouldBindJSON(&req)
//     if err != nil {
//         response.FailWithMessage("参数绑定失败:"+err.Error(), c)
//         return
//     }

//     // 检查必要的参数
//     if req.StartTime == "" || req.EndTime == "" {
//         response.FailWithMessage("开始时间和结束时间不能为空", c)
//         return
//     }

//     // 如果传入的时间是日期格式，转换为RFC3339格式
//     if len(req.StartTime) == 10 {
//         req.StartTime += "T00:00:00Z"
//         req.EndTime += "T23:59:59Z"
//     }

//     historyData, err := agvcNbqHisService.GetAgvcNbqHistory(ctx, req.AgvcNbqHis, req.StartTime, req.EndTime)
//     if err != nil {
//         global.GVA_LOG.Error("获取历史数据失败!", zap.Error(err))
//         response.FailWithMessage("获取历史数据失败:"+err.Error(), c)
//         return
//     }
//     response.OkWithData(historyData, c)
// }

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

// 用于导出token一次性存储
var (
    nbqExportTokenCache      = make(map[string]interface{})
    nbqExportTokenExpiration = make(map[string]time.Time)
    nbqTokenMutex            sync.RWMutex
)

// 五分钟检测窗口过期
func cleanupExpiredNbqExportTokens() {
    for {
        time.Sleep(5 * time.Minute)
        nbqTokenMutex.Lock()
        now := time.Now()
        for token, expiry := range nbqExportTokenExpiration {
            if now.After(expiry) {
                delete(nbqExportTokenCache, token)
                delete(nbqExportTokenExpiration, token)
            }
        }
        nbqTokenMutex.Unlock()
    }
}

func init() {
    go cleanupExpiredNbqExportTokens()
}

// ExportAgvcNbqHistory 导出逆变器历史数据token
// @Tags AgvcNbqHis
// @Summary 导出逆变器历史数据
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query agvcReq.AgvcNbqHisSearch true "查询参数"
// @Success 200 {object} response.Response{data=string,msg=string} "获取导出链接成功"
// @Router /agvcNbqHis/exportAgvcNbqHistory [get]
func (agvcNbqHisApi *AgvcNbqHisApi) ExportAgvcNbqHistory(c *gin.Context) {
    var pageInfo agvcReq.AgvcNbqHisSearch
    err := c.ShouldBindQuery(&pageInfo)
    if err != nil {
        response.FailWithMessage(err.Error(), c)
        return
    }

    // 创造一次性token
    token := utils.RandomString(32)

    // 记录本次请求参数
    exportParams := map[string]interface{}{
        "searchInfo": pageInfo,
    }

    // 参数保留记录完成鉴权
    nbqTokenMutex.Lock()
    nbqExportTokenCache[token] = exportParams
    nbqExportTokenExpiration[token] = time.Now().Add(30 * time.Minute)
    nbqTokenMutex.Unlock()

    // 生成一次性链接
    exportUrl := fmt.Sprintf("/agvcNbqHis/exportAgvcNbqHistoryByToken?token=%s", token)
    response.OkWithData(exportUrl, c)
}

// ExportAgvcNbqHistoryByToken 通过token导出逆变器历史数据
// @Tags AgvcNbqHis
// @Summary 通过token导出逆变器历史数据
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/octet-stream
// @Param token query string true "导出token"
// @Success 200 {file} file "Excel文件"
// @Router /agvcNbqHis/exportAgvcNbqHistoryByToken [get]
func (agvcNbqHisApi *AgvcNbqHisApi) ExportAgvcNbqHistoryByToken(c *gin.Context) {
    ctx := c.Request.Context()
    
    token := c.Query("token")
    if token == "" {
        response.FailWithMessage("导出token不能为空", c)
        return
    }

    // 获取token并且从缓存中剔除
    nbqTokenMutex.RLock()
    exportParamsRaw, exists := nbqExportTokenCache[token]
    expiry, _ := nbqExportTokenExpiration[token]
    nbqTokenMutex.RUnlock()

    if !exists || time.Now().After(expiry) {
        global.GVA_LOG.Error("导出token无效或已过期!")
        response.FailWithMessage("导出token无效或已过期", c)
        return
    }

    // 从token获取参数
    exportParams, ok := exportParamsRaw.(map[string]interface{})
    if !ok {
        global.GVA_LOG.Error("解析导出参数失败!")
        response.FailWithMessage("解析导出参数失败", c)
        return
    }

    // 获取导出参数
    searchInfo, ok := exportParams["searchInfo"].(agvcReq.AgvcNbqHisSearch)
    if !ok {
        global.GVA_LOG.Error("解析查询参数失败!")
        response.FailWithMessage("解析查询参数失败", c)
        return
    }

    // 清理一次性token
    nbqTokenMutex.Lock()
    delete(nbqExportTokenCache, token)
    delete(nbqExportTokenExpiration, token)
    nbqTokenMutex.Unlock()

    // 导出
    fileBuffer, fileName, err := agvcNbqHisService.ExportAgvcNbqHistory(ctx, searchInfo)
    if err != nil {
        global.GVA_LOG.Error("导出失败!", zap.Error(err))
        response.FailWithMessage("导出失败: "+err.Error(), c)
        return
    }

    c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
    c.Header("success", "true")
    c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileBuffer)
}
