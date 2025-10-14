// 自动生成模板InverterMonitor
package monitor

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 逆变器监控 结构体  InverterMonitor
type InverterMonitor struct {
	global.GVA_MODEL
	Inverter_no         *int64   `json:"inverter_no" form:"inverter_no" gorm:"column:inverter_no;"`                           //逆变器编号
	Name                *string  `json:"name" form:"name" gorm:"column:name;"`                                                //逆变器名称
	Status              *string  `json:"status" form:"status" gorm:"column:status;"`                                          //运行状态
	IsParticipateAdjust *bool    `json:"isParticipateAdjust" form:"isParticipateAdjust" gorm:"column:is_participate_adjust;"` //是否调节
	RatedPower          *float64 `json:"ratedPower" form:"ratedPower" gorm:"column:rated_power;"`                             //额定功率(kW)
	ActivePower         *float64 `json:"activePower" form:"activePower" gorm:"column:active_power;"`                          //有功功率
	ApTargetValue       *float64 `json:"apTargetValue" form:"apTargetValue" gorm:"column:ap_target_value;"`                   //有功目标(kW)
	ReactivePower       *float64 `json:"reactivePower" form:"reactivePower" gorm:"column:reactive_power;"`                    //无功功率(kVar)
	RpTargetValue       *float64 `json:"rpTargetValue" form:"rpTargetValue" gorm:"column:rp_target_value;"`                   //无功目标(kVar)
	PowerFactor         *float64 `json:"powerFactor" form:"powerFactor" gorm:"column:power_factor;"`                          //功率因数
}

// TableName 逆变器监控 InverterMonitor自定义表名 inverter_monitor
func (InverterMonitor) TableName() string {
	return "inverter_monitor"
}

type InverterHistory struct {
	Time        time.Time `json:"time"`        // 时间戳
	Code        string    `json:"code"`        // 设备编码
	Value       float64   `json:"value"`       // 值
	PowerFactor float64   `json:"powerFactor"` // 功率因数
	ActivePower float64   `json:"activePower"` // 有功功率
}