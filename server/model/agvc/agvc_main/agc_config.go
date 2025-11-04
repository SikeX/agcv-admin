package agvc_main

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// AGCConfig AGC配置 结构体
type AGCConfig struct {
	global.GVA_MODEL
	PSID              string  `json:"psid" form:"psid" gorm:"column:psid;comment:电站ID;uniqueIndex"`
	IsActive          *int    `json:"isActive" form:"isActive" gorm:"column:is_active;comment:是否投入;default:0"`          // 0:未投入 1:投入
	ControlAuth       *int    `json:"controlAuth" form:"controlAuth" gorm:"column:control_auth;comment:控制权限;default:1"` // 1:调度 2:站内
	RunMode           *int    `json:"runMode" form:"runMode" gorm:"column:run_mode;comment:运行模式;default:2"`             // 1:开环 2:闭环
	JitterRange       float64 `json:"jitterRange" form:"jitterRange" gorm:"column:jitter_range;comment:抖动区间(MW);default:0.5"`
	RegPeriod         int     `json:"regPeriod" form:"regPeriod" gorm:"column:reg_period;comment:调节周期(秒);default:30"`
	RegStep           float64 `json:"regStep" form:"regStep" gorm:"column:reg_step;comment:调节步长(MW);default:0.1"`
	PowerUpperLimit   float64 `json:"powerUpperLimit" form:"powerUpperLimit" gorm:"column:power_upper_limit;comment:有功调节上限(MW)"`
	PowerLowerLimit   float64 `json:"powerLowerLimit" form:"powerLowerLimit" gorm:"column:power_lower_limit;comment:有功调节下限(MW)"`
	PowerExecValue    float64 `json:"powerExecValue" form:"powerExecValue" gorm:"column:power_exec_value;comment:有功执行值(MW)"`
	UpRegLock         *int    `json:"upRegLock" form:"upRegLock" gorm:"column:up_reg_lock;comment:上调节闭锁;default:0"`       // 0:未闭锁 1:闭锁
	DownRegLock       *int    `json:"downRegLock" form:"downRegLock" gorm:"column:down_reg_lock;comment:下调节闭锁;default:0"` // 0:未闭锁 1:闭锁
	DispatchExecValue float64 `json:"dispatchExecValue" form:"dispatchExecValue" gorm:"column:dispatch_exec_value;comment:调度执行值(MW)"`
	StationExecValue  float64 `json:"stationExecValue" form:"stationExecValue" gorm:"column:station_exec_value;comment:站内执行值(MW)"`
}

// TableName AGCConfig表名
func (AGCConfig) TableName() string {
	return "agvc_agc_config"
}

// AGCRegulationRecord AGC调节记录
type AGCRegulationRecord struct {
	global.GVA_MODEL
	BwdNo           int     `json:"bwdNo" form:"bwdNo" gorm:"column:bwdNo;comment:并网点编号;index"`
	TargetPower     float64 `json:"targetPower" form:"targetPower" gorm:"column:target_power;comment:目标出力(MW)"`
	ActualPower     float64 `json:"actualPower" form:"actualPower" gorm:"column:actual_power;comment:实际出力(MW)"`
	PowerDeviation  float64 `json:"powerDeviation" form:"powerDeviation" gorm:"column:power_deviation;comment:出力偏差(MW)"`
	RegulationPower float64 `json:"regulationPower" form:"regulationPower" gorm:"column:regulation_power;comment:调节量(MW)"`
	SystemFreq      float64 `json:"systemFreq" form:"systemFreq" gorm:"column:system_freq;comment:系统频率(Hz)"`
	Status          string  `json:"status" form:"status" gorm:"column:status;comment:状态"` // success, failed, timeout
	Message         string  `json:"message" form:"message" gorm:"column:message;comment:信息;type:text"`
}

// TableName AGCRegulationRecord表名
func (AGCRegulationRecord) TableName() string {
	return "agvc_agc_regulation_record"
}

// InverterRegulation 逆变器调节详情
type InverterRegulation struct {
	global.GVA_MODEL
	RecordID        uint    `json:"recordId" form:"recordId" gorm:"column:record_id;comment:调节记录ID;index"`
	BwdNo           int     `json:"bwdNo" form:"bwdNo" gorm:"column:bwdNo;comment:并网点编号"`
	EQID            int     `json:"eqid" form:"eqid" gorm:"column:eqid;comment:逆变器ID"`
	RegulationPower float64 `json:"regulationPower" form:"regulationPower" gorm:"column:regulation_power;comment:分配调节量(MW)"`
	BeforePower     float64 `json:"beforePower" form:"beforePower" gorm:"column:before_power;comment:调节前功率(MW)"`
	AfterPower      float64 `json:"afterPower" form:"afterPower" gorm:"column:after_power;comment:调节后功率(MW)"`
	Status          string  `json:"status" form:"status" gorm:"column:status;comment:状态"` // success, failed, timeout
}

// TableName InverterRegulation表名
func (InverterRegulation) TableName() string {
	return "agvc_inverter_regulation"
}
