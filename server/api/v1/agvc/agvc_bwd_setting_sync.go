package agvc

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	agvcMainService "github.com/flipped-aurora/gin-vue-admin/server/service/agvc/agcv_main"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AgvcBwdSettingSyncApi struct{}

// UpdateFromRemote 从远程调度更新并网点配置
// @Tags AgvcBwdSetting
// @Summary 从远程调度更新并网点配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body map[string]interface{} true "并网点编号和更新数据"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /agvcBwdSetting/updateFromRemote [post]
func (agvcBwdSettingSyncApi *AgvcBwdSettingSyncApi) UpdateFromRemote(c *gin.Context) {
	var req struct {
		BwdNo   int                    `json:"bwdNo" binding:"required"`   // 并网点编号
		Updates map[string]interface{} `json:"updates" binding:"required"` // 更新数据
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 通过缓存管理器更新配置
	if err := agvcMainService.SettingCache.UpdateFromRemote(req.BwdNo, req.Updates); err != nil {
		global.GVA_LOG.Error("从远程更新并网点配置失败",
			zap.Int("bwdNo", req.BwdNo),
			zap.Error(err))
		response.FailWithMessage("更新失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

// SyncToDispatch 同步并网点配置到调度
// @Tags AgvcBwdSetting
// @Summary 同步并网点配置到调度
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body map[string]interface{} true "并网点编号"
// @Success 200 {object} response.Response{msg=string} "同步成功"
// @Router /agvcBwdSetting/syncToDispatch [post]
func (agvcBwdSettingSyncApi *AgvcBwdSettingSyncApi) SyncToDispatch(c *gin.Context) {
	var req struct {
		BwdNo int `json:"bwdNo" binding:"required"` // 并网点编号
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 从缓存获取配置
	setting, exists := agvcMainService.SettingCache.Get(req.BwdNo)
	if !exists {
		response.FailWithMessage("并网点配置不存在", c)
		return
	}

	// 通过缓存管理器同步到调度
	if err := agvcMainService.SettingCache.Update(req.BwdNo, setting); err != nil {
		global.GVA_LOG.Error("同步并网点配置到调度失败",
			zap.Int("bwdNo", req.BwdNo),
			zap.Error(err))
		response.FailWithMessage("同步失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("同步成功", c)
}
