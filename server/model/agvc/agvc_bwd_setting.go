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
	AgcIsEnabled                  *int64   `json:"agcIsEnabled" form:"agcIsEnabled" gorm:"column:agc_is_enabled;"`                                                     //AGC是否启用
	AvcIsEnabled                  *int64   `json:"avcIsEnabled" form:"avcIsEnabled" gorm:"column:avc_is_enabled;"`                                                     //AVC是否启用
	AgcRemoteMode                 *int64   `json:"agcRemoteMode" form:"agcRemoteMode" gorm:"column:agc_remote_mode;"`                                                  //AGC远程模式
	AvcRemoteMode                 *int64   `json:"avcRemoteMode" form:"avcRemoteMode" gorm:"column:avc_remote_mode;"`                                                  //AVC远程模式
	AgcDispatchExecValue          *float64 `json:"agcDispatchExecValue" form:"agcDispatchExecValue" gorm:"column:agc_dispatch_exec_value;"`                            //agc度执行值(kW)
	StationExecValue              *float64 `json:"stationExecValue" form:"stationExecValue" gorm:"column:station_exec_value;"`                                         //站点执行值(kW)
	AgcLoopStatus                 *int64   `json:"agcLoopStatus" form:"agcLoopStatus" gorm:"column:agc_loop_status;"`
	AvcLoopStatus                 *int64   `json:"avcLoopStatus" form:"avcLoopStatus" gorm:"column:avc_loop_status;"`              // avc运行模式,1-开环，2-闭环
	AvcVolteExecValue             *float64 `json:"avcVolteExecValue" form:"avcVolteExecValue" gorm:"column:avc_volte_exec_value;"` //avc电压执行值(kV)
	AvcWGExecValue                *float64 `json:"avcWGExecValue" form:"avcWGExecValue" gorm:"column:avc_wg_exec_value;"`          //avc无功执行值
}

// TableName 并网点配置 AgvcBwdSetting自定义表名 agvc_bwd_setting
func (AgvcBwdSetting) TableName() string {
	return "agvc_bwd_setting"
}
