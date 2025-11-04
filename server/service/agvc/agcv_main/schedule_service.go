package agcv_main

import (
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
	"go.uber.org/zap"
)

// ScheduleService 计划曲线服务
type scheduleService struct {
	stopChan chan struct{}
	running  bool
}

var ScheduleService = new(scheduleService)

// Initialize 初始化计划曲线服务
func (ss *scheduleService) Initialize() {
	ss.stopChan = make(chan struct{})
	ss.running = true

	// 启动定时检查任务（每分钟检查一次）
	go ss.checkAndExecuteSchedules()

	global.GVA_LOG.Info("计划曲线服务初始化成功")
}

// Stop 停止计划曲线服务
func (ss *scheduleService) Stop() {
	if ss.running {
		close(ss.stopChan)
		ss.running = false
		global.GVA_LOG.Info("计划曲线服务已停止")
	}
}

// checkAndExecuteSchedules 定期检查并执行计划曲线
func (ss *scheduleService) checkAndExecuteSchedules() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ss.executeSchedules()
		case <-ss.stopChan:
			return
		}
	}
}

// executeSchedules 执行到期的计划曲线
func (ss *scheduleService) executeSchedules() {
	now := time.Now()
	currentTime := now.Format("15:04") // HH:MM格式
	currentDate := now.Format("2006-01-02")

	// 查询所有启用且到期的计划
	var schedules []agvc.AgvcScheduleCurve
	err := global.GVA_DB.Where("enabled = ? AND start_time = ?", 1, currentTime).Find(&schedules).Error
	if err != nil {
		global.GVA_LOG.Error("查询计划曲线失败", zap.Error(err))
		return
	}

	if len(schedules) == 0 {
		return
	}

	global.GVA_LOG.Info("发现待执行的计划曲线",
		zap.Int("数量", len(schedules)),
		zap.String("时间", currentTime))

	for _, schedule := range schedules {
		// 检查今日是否已执行
		if schedule.Executed == 1 && schedule.LastExecAt != nil {
			lastExecTime := time.Unix(*schedule.LastExecAt, 0)
			if lastExecTime.Format("2006-01-02") == currentDate {
				global.GVA_LOG.Debug("计划今日已执行，跳过",
					zap.Uint("ID", schedule.ID),
					zap.Int("bwdNo", schedule.BwdNo))
				continue
			}
		}

		// 执行计划
		if err := ss.executeSchedule(&schedule); err != nil {
			global.GVA_LOG.Error("执行计划曲线失败",
				zap.Uint("ID", schedule.ID),
				zap.Int("bwdNo", schedule.BwdNo),
				zap.Error(err))
			continue
		}

		// 更新执行状态
		nowTimestamp := now.Unix()
		global.GVA_DB.Model(&schedule).Updates(map[string]interface{}{
			"executed":     1,
			"last_exec_at": nowTimestamp,
		})
	}
}

// executeSchedule 执行单个计划曲线
func (ss *scheduleService) executeSchedule(schedule *agvc.AgvcScheduleCurve) error {
	global.GVA_LOG.Info("执行计划曲线",
		zap.Uint("ID", schedule.ID),
		zap.Int("bwdNo", schedule.BwdNo),
		zap.Int("type", schedule.Type),
		zap.Int("source", schedule.Source),
		zap.Float64("targetValue", schedule.TargetValue))

	// 获取并网点配置
	var config agvc.AgvcBwdSetting
	err := global.GVA_DB.Where("number = ?", schedule.BwdNo).First(&config).Error
	if err != nil {
		return fmt.Errorf("获取并网点配置失败: %v", err)
	}

	// 根据类型和来源更新对应的执行值
	updates := make(map[string]interface{})

	switch schedule.Type {
	case 1: // AGC
		if schedule.Source == 1 {
			// 本地曲线：更新站内执行值
			updates["station_exec_value"] = schedule.TargetValue
			global.GVA_LOG.Info("更新AGC站内执行值",
				zap.Int("bwdNo", schedule.BwdNo),
				zap.Float64("value", schedule.TargetValue))
		} else if schedule.Source == 2 {
			// 调度曲线：更新调度执行值
			updates["dispatch_exec_value"] = schedule.TargetValue
			global.GVA_LOG.Info("更新AGC调度执行值",
				zap.Int("bwdNo", schedule.BwdNo),
				zap.Float64("value", schedule.TargetValue))
		}

	case 2: // AVC
		if schedule.Source == 1 {
			// 本地曲线：更新站内执行值
			updates["station_exec_value"] = schedule.TargetValue
			global.GVA_LOG.Info("更新AVC站内执行值",
				zap.Int("bwdNo", schedule.BwdNo),
				zap.Float64("value", schedule.TargetValue))
		} else if schedule.Source == 2 {
			// 调度曲线：更新调度执行值
			updates["dispatch_exec_value"] = schedule.TargetValue
			global.GVA_LOG.Info("更新AVC调度执行值",
				zap.Int("bwdNo", schedule.BwdNo),
				zap.Float64("value", schedule.TargetValue))
		}
	}

	// 更新数据库
	if len(updates) > 0 {
		err = global.GVA_DB.Model(&agvc.AgvcBwdSetting{}).
			Where("number = ?", schedule.BwdNo).
			Updates(updates).Error
		if err != nil {
			return fmt.Errorf("更新配置失败: %v", err)
		}
	}

	return nil
}

// CreateSchedule 创建计划曲线
func (ss *scheduleService) CreateSchedule(schedule *agvc.AgvcScheduleCurve) error {
	return global.GVA_DB.Create(schedule).Error
}

// UpdateSchedule 更新计划曲线
func (ss *scheduleService) UpdateSchedule(schedule *agvc.AgvcScheduleCurve) error {
	return global.GVA_DB.Save(schedule).Error
}

// DeleteSchedule 删除计划曲线
func (ss *scheduleService) DeleteSchedule(id uint) error {
	return global.GVA_DB.Delete(&agvc.AgvcScheduleCurve{}, id).Error
}

// GetSchedulesByBwdNo 获取某并网点的所有计划曲线
func (ss *scheduleService) GetSchedulesByBwdNo(bwdNo int, scheduleType int, source int) ([]agvc.AgvcScheduleCurve, error) {
	var schedules []agvc.AgvcScheduleCurve
	db := global.GVA_DB.Where("bwd_no = ?", bwdNo)

	if scheduleType > 0 {
		db = db.Where("type = ?", scheduleType)
	}

	if source > 0 {
		db = db.Where("source = ?", source)
	}

	err := db.Order("start_time ASC").Find(&schedules).Error
	return schedules, err
}

// ResetDailyExecutionStatus 重置每日执行状态（在每天00:00调用）
func (ss *scheduleService) ResetDailyExecutionStatus() error {
	return global.GVA_DB.Model(&agvc.AgvcScheduleCurve{}).
		Where("enabled = ?", 1).
		Update("executed", 0).Error
}
