package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// AVCConfig AVC配置 结构体
type AVCConfig struct {
	global.GVA_MODEL
	PSID               string  `json:"psid" form:"psid" gorm:"column:psid;comment:电站ID;uniqueIndex"`
	IsActive           *int    `json:"isActive" form:"isActive" gorm:"column:is_active;comment:是否投入;default:0"`              // 0:未投入 1:投入
	ControlAuth        *int    `json:"controlAuth" form:"controlAuth" gorm:"column:control_auth;comment:控制权限;default:1"`      // 1:调度 2:站内
	CommandStatus      *int    `json:"commandStatus" form:"commandStatus" gorm:"column:command_status;comment:指令状态;default:0"` // 0:无指令 1:有指令
	RunMode            *int    `json:"runMode" form:"runMode" gorm:"column:run_mode;comment:运行模式;default:2"`                  // 1:开环 2:闭环
	TargetVoltageLow   float64 `json:"targetVoltageLow" form:"targetVoltageLow" gorm:"column:target_voltage_low;comment:目标电压下限(kV)"`
	TargetVoltageHigh  float64 `json:"targetVoltageHigh" form:"targetVoltageHigh" gorm:"column:target_voltage_high;comment:目标电压上限(kV)"`
	VoltageExecValue   float64 `json:"voltageExecValue" form:"voltageExecValue" gorm:"column:voltage_exec_value;comment:电压执行值(kV)"`
	ReactiveExecValue  float64 `json:"reactiveExecValue" form:"reactiveExecValue" gorm:"column:reactive_exec_value;comment:无功执行值(MVar)"`
	ReactiveIncCap     float64 `json:"reactiveIncCap" form:"reactiveIncCap" gorm:"column:reactive_inc_cap;comment:无功可增容量(MVar)"`
	ReactiveDecCap     float64 `json:"reactiveDecCap" form:"reactiveDecCap" gorm:"column:reactive_dec_cap;comment:无功可减容量(MVar)"`
	UpRegLock          *int    `json:"upRegLock" form:"upRegLock" gorm:"column:up_reg_lock;comment:上调节闭锁;default:0"`           // 0:未闭锁 1:闭锁
	DownRegLock        *int    `json:"downRegLock" form:"downRegLock" gorm:"column:down_reg_lock;comment:下调节闭锁;default:0"`     // 0:未闭锁 1:闭锁
	RegPeriod          int     `json:"regPeriod" form:"regPeriod" gorm:"column:reg_period;comment:调节周期(秒);default:60"`
	VoltageDeadZone    float64 `json:"voltageDeadZone" form:"voltageDeadZone" gorm:"column:voltage_dead_zone;comment:电压死区(kV);default:0.5"`
	ReactiveSensitivity float64 `json:"reactiveSensitivity" form:"reactiveSensitivity" gorm:"column:reactive_sensitivity;comment:无功灵敏度(MVar/kV);default:1.0"`
}

// TableName AVCConfig表名
func (AVCConfig) TableName() string {
	return "agvc_avc_config"
}

// AVCRegulationRecord AVC调节记录
type AVCRegulationRecord struct {
	global.GVA_MODEL
	PSID              string  `json:"psid" form:"psid" gorm:"column:psid;comment:电站ID;index"`
	TargetVoltage     float64 `json:"targetVoltage" form:"targetVoltage" gorm:"column:target_voltage;comment:目标电压(kV)"`
	ActualVoltage     float64 `json:"actualVoltage" form:"actualVoltage" gorm:"column:actual_voltage;comment:实际电压(kV)"`
	VoltageDeviation  float64 `json:"voltageDeviation" form:"voltageDeviation" gorm:"column:voltage_deviation;comment:电压偏差(kV)"`
	RequiredReactive  float64 `json:"requiredReactive" form:"requiredReactive" gorm:"column:required_reactive;comment:需求无功(MVar)"`
	RegulationReactive float64 `json:"regulationReactive" form:"regulationReactive" gorm:"column:regulation_reactive;comment:调节无功(MVar)"`
	TotalReactive     float64 `json:"totalReactive" form:"totalReactive" gorm:"column:total_reactive;comment:总无功(MVar)"`
	SystemFreq        float64 `json:"systemFreq" form:"systemFreq" gorm:"column:system_freq;comment:系统频率(Hz)"`
	Status            string  `json:"status" form:"status" gorm:"column:status;comment:状态"` // success, failed, timeout
	Message           string  `json:"message" form:"message" gorm:"column:message;comment:信息;type:text"`
}

// TableName AVCRegulationRecord表名
func (AVCRegulationRecord) TableName() string {
	return "agvc_avc_regulation_record"
}

// DeviceReactiveRegulation 设备无功调节详情
type DeviceReactiveRegulation struct {
	global.GVA_MODEL
	RecordID          uint    `json:"recordId" form:"recordId" gorm:"column:record_id;comment:调节记录ID;index"`
	PSID              string  `json:"psid" form:"psid" gorm:"column:psid;comment:电站ID"`
	EQID              string  `json:"eqid" form:"eqid" gorm:"column:eqid;comment:设备ID"`
	EQType            string  `json:"eqType" form:"eqType" gorm:"column:eq_type;comment:设备类型"`
	RegulationReactive float64 `json:"regulationReactive" form:"regulationReactive" gorm:"column:regulation_reactive;comment:分配无功(MVar)"`
	BeforeReactive    float64 `json:"beforeReactive" form:"beforeReactive" gorm:"column:before_reactive;comment:调节前无功(MVar)"`
	AfterReactive     float64 `json:"afterReactive" form:"afterReactive" gorm:"column:after_reactive;comment:调节后无功(MVar)"`
	Status            string  `json:"status" form:"status" gorm:"column:status;comment:状态"` // success, failed, timeout
}

// TableName DeviceReactiveRegulation表名
func (DeviceReactiveRegulation) TableName() string {
	return "agvc_device_reactive_regulation"
}
