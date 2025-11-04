package agvc

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// AgvcScheduleCurve AGC/AVC计划曲线
type AgvcScheduleCurve struct {
	global.GVA_MODEL
	BwdNo       int     `json:"bwdNo" gorm:"column:bwd_no;comment:并网点编号;not null"`
	Type        int     `json:"type" gorm:"column:type;comment:类型 1:AGC 2:AVC;not null"`
	Source      int     `json:"source" gorm:"column:source;comment:来源 1:本地 2:调度;not null"`
	StartTime   string  `json:"startTime" gorm:"column:start_time;comment:开始时间 HH:MM格式;not null"`
	TargetValue float64 `json:"targetValue" gorm:"column:target_value;comment:目标值;not null"`
	Enabled     int     `json:"enabled" gorm:"column:enabled;comment:是否启用 0:禁用 1:启用;default:1"`
	Executed    int     `json:"executed" gorm:"column:executed;comment:今日是否已执行 0:未执行 1:已执行;default:0"`
	LastExecAt  *int64  `json:"lastExecAt" gorm:"column:last_exec_at;comment:最后执行时间戳"`
}

func (AgvcScheduleCurve) TableName() string {
	return "agvc_schedule_curve"
}
