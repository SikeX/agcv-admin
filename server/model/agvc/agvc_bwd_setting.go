// 自动生成模板AgvcBwdSetting
package agvc

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 并网点配置 结构体  AgvcBwdSetting
type AgvcBwdSetting struct {
	global.GVA_MODEL
	Name                          *string  `json:"name" form:"name" gorm:"column:name;"`                                                                               //并网点名称
	Number                        *string  `json:"number" form:"number" gorm:"column:number;"`                                                                         //设备编号
	VoltageLevel                  *float64 `json:"voltageLevel" form:"voltageLevel" gorm:"column:voltage_level;"`                                                      //电压等级(kV)
	AgcFunctionExit               *int64   `json:"agcFunctionExit" form:"agcFunctionExit" gorm:"column:agc_function_exit;"`                                            //AGC功能退出模式
	AgcStepSize                   *float64 `json:"agcStepSize" form:"agcStepSize" gorm:"column:agc_step_size;"`                                                        //AGC调节步长(kW)
	AgcStepPeriod                 *int64   `json:"agcStepPeriod" form:"agcStepPeriod" gorm:"column:agc_step_period;"`                                                  //AGC步长周期(秒)
	AgcVibrationRange             *float64 `json:"agcVibrationRange" form:"agcVibrationRange" gorm:"column:agc_vibration_range;"`                                      //AGC抖动区间(kW)
	AgcControlPeriod              *int64   `json:"agcControlPeriod" form:"agcControlPeriod" gorm:"column:agc_control_period;"`                                         //AGC调控周期(秒)
	AgcMicroAdjustmentCoefficient *float64 `json:"agcMicroAdjustmentCoefficient" form:"agcMicroAdjustmentCoefficient" gorm:"column:agc_micro_adjustment_coefficient;"` //AGC微调系数
	AvcStepSize                   *float64 `json:"avcStepSize" form:"avcStepSize" gorm:"column:avc_step_size;"`                                                        //AVC调节步长(kV)
	AvcStepPeriod                 *int64   `json:"avcStepPeriod" form:"avcStepPeriod" gorm:"column:avc_step_period;"`                                                  //AVC步长周期(秒)
	AvcVibrationRange             *float64 `json:"avcVibrationRange" form:"avcVibrationRange" gorm:"column:avc_vibration_range;"`                                      //AVC抖动区间(kV)
	AvcControlPeriod              *int64   `json:"avcControlPeriod" form:"avcControlPeriod" gorm:"column:avc_control_period;"`                                         //AVC调控周期(秒)
	AvcSystemImpedance            *float64 `json:"avcSystemImpedance" form:"avcSystemImpedance" gorm:"column:avc_system_impedance;"`                                   //AVC系统阻抗
	AvcAdjustmentRangeMin         *float64 `json:"avcAdjustmentRangeMin" form:"avcAdjustmentRangeMin" gorm:"column:avc_adjustment_range_min;"`                         //AVC调节范围最小值(kV)
	AvcAdjustmentRangeMax         *float64 `json:"avcAdjustmentRangeMax" form:"avcAdjustmentRangeMax" gorm:"column:avc_adjustment_range_max;"`                         //AVC调节范围最大值(kV)
	AgcIsEnabled                  *int64   `json:"agcIsEnabled" form:"agcIsEnabled" gorm:"column:agc_is_enabled;"`
	AvcIsEnabled                  *int64   `json:"avcIsEnabled" form:"avcIsEnabled" gorm:"column:avc_is_enabled;"`
	ControlAuth                   *int64   `json:"controlAuth" form:"controlAuth" gorm:"column:control_auth;"`
	DispatchExecValue             *float64 `json:"dispatchExecValue" form:"dispatchExecValue" gorm:"column:dispatch_exec_value;"`
	StationExecValue              *float64 `json:"stationExecValue" form:"stationExecValue" gorm:"column:station_exec_value;"`
	RunMode                       *int64   `json:"runMode" form:"runMode" gorm:"column:run_mode;"` // 运行模式,1-开环，2-闭环
}

// TableName 并网点配置 AgvcBwdSetting自定义表名 agvc_bwd_setting
func (AgvcBwdSetting) TableName() string {
	return "agvc_bwd_setting"
}
