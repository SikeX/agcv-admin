// 自动生成模板AgvcBwdHis
package agvc

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// agvcBwdHis表 结构体  AgvcBwdHis
type AgvcBwdHis struct {
	global.GVA_MODEL
	Number                *string  `json:"number" form:"number" gorm:"column:number;comment:设备编号"`
	Name                  *string  `json:"name" form:"name" gorm:"column:name;comment:并网点名称"`
	TargetActivePower     *float64 `json:"targetActivePower" form:"targetActivePower" gorm:"column:target_active_power;comment:目标有功(kW)"`
	TargetVoltage         *float64 `json:"targetVoltage" form:"targetVoltage" gorm:"column:target_voltage;comment:目标电压(kV)"`
	AgcParameters         *string  `json:"agcParameters" form:"agcParameters" gorm:"column:agc_parameters;type:json;comment:AGC参数配置"`
	AvcParameters         *string  `json:"avcParameters" form:"avcParameters" gorm:"column:avc_parameters;type:json;comment:AVC参数配置"`
	AgcPlanCurves         *string  `json:"agcPlanCurves" form:"agcPlanCurves" gorm:"column:agc_plan_curves;type:json;comment:AGC计划曲线(本地曲线)"`
	AvcVoltagePlanCurves  *string  `json:"avcVoltagePlanCurves" form:"avcVoltagePlanCurves" gorm:"column:avc_voltage_plan_curves;type:json;comment:AVC电压计划曲线"`
	AvcReactivePlanCurves *string  `json:"avcReactivePlanCurves" form:"avcReactivePlanCurves" gorm:"column:avc_reactive_plan_curves;type:json;comment:AVC无功计划曲线"`
	CurrentActivePower    *float64 `json:"currentActivePower" form:"currentActivePower" gorm:"column:current_active_power;comment:当前有功(kW)"`
	CurrentVoltage        *float64 `json:"currentVoltage" form:"currentVoltage" gorm:"column:current_voltage;comment:当前电压(kV)"`
	SystemFrequency       *float64 `json:"systemFrequency" form:"systemFrequency" gorm:"column:system_frequency;comment:系统频率(Hz)"`
	SystemImpedance       *float64 `json:"systemImpedance" form:"systemImpedance" gorm:"column:system_impedance;comment:系统阻抗"`
	// AGC状态字段
	AgcFunctionState    *int64 `json:"agcFunctionState" form:"agcFunctionState" gorm:"column:agc_function_state;comment:AGC功能投退状态(1:投入,0:退出)"`
	AgcControlMode      *int64 `json:"agcControlMode" form:"agcControlMode" gorm:"column:agc_control_mode;comment:AGC调节方式(1:闭环指导,0:开环指导)"`
	AgcControlAuthority *int64 `json:"agcControlAuthority" form:"agcControlAuthority" gorm:"column:agc_control_authority;comment:AGC控制权限(1:调度控制,0:站内控制)"`
	// AVC状态字段
	AvcFunctionState    *int64 `json:"avcFunctionState" form:"avcFunctionState" gorm:"column:avc_function_state;comment:AVC功能投退状态(1:投入,0:退出)"`
	AvcControlMode      *int64 `json:"avcControlMode" form:"avcControlMode" gorm:"column:avc_control_mode;comment:AVC调节方式(1:开环指导,0:闭环调节)"`
	AvcControlAuthority *int64 `json:"avcControlAuthority" form:"avcControlAuthority" gorm:"column:avc_control_authority;comment:AVC控制权限(1:站内控制,0:调度控制)"`
}

// TableName agvcBwdHis表 AgvcBwdHis自定义表名 agvc_bwd_his
func (AgvcBwdHis) TableName() string {
	return "agvc_bwd_his"
}
