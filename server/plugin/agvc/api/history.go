package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/agvc/model/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type history struct{}

var History = new(history)

// QueryHistoryData 查询历史数据
// @Tags     AGVC_History
// @Summary  查询历史数据
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.HistoryDataRequest true "查询条件"
// @Success  200  {object} response.Response{data=[]map[string]interface{},msg=string} "查询成功"
// @Router   /agvc/history/query [post]
func (a *history) QueryHistoryData(c *gin.Context) {
	var req request.HistoryDataRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	data, err := serviceHistory.QueryHistoryData(req)
	if err != nil {
		global.GVA_LOG.Error("查询历史数据失败", zap.Error(err))
		response.FailWithMessage("查询历史数据失败: "+err.Error(), c)
		return
	}

	response.OkWithData(data, c)
}
