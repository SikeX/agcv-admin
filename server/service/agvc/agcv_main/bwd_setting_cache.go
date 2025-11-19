package agcv_main

import (
	"fmt"
	"sync"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
	"go.uber.org/zap"
)

// BwdSettingCache 并网点配置缓存管理器
type BwdSettingCache struct {
	cache map[int]agvc.AgvcBwdSetting // key: bwdNo, value: AgvcBwdSetting
	mutex sync.RWMutex
}

var SettingCache = &BwdSettingCache{
	cache: make(map[int]agvc.AgvcBwdSetting),
}

// Initialize 初始化缓存，从数据库加载所有配置
func (c *BwdSettingCache) Initialize() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var settings []agvc.AgvcBwdSetting
	if err := global.GVA_DB.Find(&settings).Error; err != nil {
		return fmt.Errorf("加载并网点配置失败: %v", err)
	}

	for _, setting := range settings {
		if setting.Eqid != nil {
			var bwdNo int
			if _, err := fmt.Sscanf(fmt.Sprintf("%d", *setting.Eqid), "%d", &bwdNo); err == nil {
				c.cache[bwdNo] = setting
			}
		}
	}

	global.GVA_LOG.Info("并网点配置缓存初始化完成", zap.Int("配置数量", len(c.cache)))
	return nil
}

// Get 获取并网点配置（优先从缓存获取）
func (c *BwdSettingCache) Get(bwdNo int) (agvc.AgvcBwdSetting, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	setting, exists := c.cache[bwdNo]
	return setting, exists
}

// Set 设置并网点配置到缓存
func (c *BwdSettingCache) Set(bwdNo int, setting agvc.AgvcBwdSetting) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache[bwdNo] = setting
	global.GVA_LOG.Debug("更新并网点配置缓存", zap.Int("bwdNo", bwdNo))
}

// Update 更新并网点配置（同时更新数据库和缓存，并同步到调度）
func (c *BwdSettingCache) Update(bwdNo int, setting agvc.AgvcBwdSetting) error {
	// 更新数据库
	if err := global.GVA_DB.Model(&agvc.AgvcBwdSetting{}).
		Where("number = ?", fmt.Sprintf("%d", bwdNo)).
		Updates(&setting).Error; err != nil {
		return fmt.Errorf("更新数据库失败: %v", err)
	}

	// 更新缓存
	c.Set(bwdNo, setting)

	// 同步到调度
	if err := c.syncToDispatch(bwdNo, setting); err != nil {
		global.GVA_LOG.Warn("同步配置到调度失败",
			zap.Int("bwdNo", bwdNo),
			zap.Error(err))
	}

	global.GVA_LOG.Info("并网点配置更新成功",
		zap.Int("bwdNo", bwdNo))

	return nil
}

// UpdateFromRemote 从远程调度更新配置（调度推送过来的数据）
func (c *BwdSettingCache) UpdateFromRemote(bwdNo int, updates map[string]interface{}) error {
	// 获取当前配置
	currentSetting, exists := c.Get(bwdNo)
	if !exists {
		// 如果缓存中不存在，从数据库加载
		var setting agvc.AgvcBwdSetting
		if err := global.GVA_DB.Where("number = ?", fmt.Sprintf("%d", bwdNo)).
			First(&setting).Error; err != nil {
			return fmt.Errorf("配置不存在: %v", err)
		}
		currentSetting = setting
	}

	// 应用更新
	if val, ok := updates["dispatch_exec_value"].(float64); ok {
		currentSetting.AgcDispatchExecValue = &val
	}
	if val, ok := updates["control_auth"].(int64); ok {
		currentSetting.AgcRemoteMode = &val
	}
	if val, ok := updates["agc_is_enabled"].(int64); ok {
		currentSetting.AgcIsEnabled = &val
	}
	if val, ok := updates["avc_is_enabled"].(int64); ok {
		currentSetting.AvcIsEnabled = &val
	}
	if val, ok := updates["run_mode"].(int64); ok {
		currentSetting.AgcLoopStatus = &val
	}

	// 更新数据库
	if err := global.GVA_DB.Model(&agvc.AgvcBwdSetting{}).
		Where("number = ?", fmt.Sprintf("%d", bwdNo)).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("更新数据库失败: %v", err)
	}

	// 更新缓存
	c.Set(bwdNo, currentSetting)

	global.GVA_LOG.Info("从远程更新并网点配置成功",
		zap.Int("bwdNo", bwdNo))

	return nil
}

// syncToDispatch 同步配置到调度（1189端口）
func (c *BwdSettingCache) syncToDispatch(bwdNo int, setting agvc.AgvcBwdSetting) error {
	results := make(map[string]interface{})

	// AGC相关配置
	if setting.AgcIsEnabled != nil {
		results["agcSignal"] = *setting.AgcIsEnabled
	}
	if setting.AgcRemoteMode != nil {
		results["agcControlMode"] = *setting.AgcRemoteMode
	}
	if setting.AgcLoopStatus != nil {
		results["agcLoopStatus"] = *setting.AgcLoopStatus
	}
	if setting.AgcDispatchExecValue != nil {
		results["powerExecValue"] = *setting.AgcDispatchExecValue
	}

	// AVC相关配置
	if setting.AvcIsEnabled != nil {
		results["avcSignal"] = *setting.AvcIsEnabled
	}

	// 发送AGC配置到调度
	if len(results) > 0 {
		if err := CoapSender.SendAGCResultToDispatch(bwdNo, results); err != nil {
			return fmt.Errorf("发送AGC配置到调度失败: %v", err)
		}
	}

	return nil
}

// Delete 删除并网点配置
func (c *BwdSettingCache) Delete(bwdNo int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.cache, bwdNo)
	global.GVA_LOG.Debug("删除并网点配置缓存", zap.Int("bwdNo", bwdNo))
}

// GetAll 获取所有并网点配置
func (c *BwdSettingCache) GetAll() map[int]agvc.AgvcBwdSetting {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// 返回副本以避免并发问题
	result := make(map[int]agvc.AgvcBwdSetting, len(c.cache))
	for k, v := range c.cache {
		result[k] = v
	}
	return result
}
