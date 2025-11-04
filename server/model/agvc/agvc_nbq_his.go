// 自动生成模板AgvcNbqHis
package agvc

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// agvcNbqHis表 结构体  AgvcNbqHis
type AgvcNbqHis struct {
	global.GVA_MODEL
	InverterNo          *int     `json:"inverterNo" form:"inverterNo" gorm:"column:inverter_no;comment:逆变器编号"`
	Name                *string  `json:"name" form:"name" gorm:"column:name;comment:逆变器名称"`
	Status              *string  `json:"status" form:"status" gorm:"column:status;comment:逆变器运行状态"`
	IsParticipateAdjust *bool    `json:"isParticipateAdjust" form:"isParticipateAdjust" gorm:"column:is_participate_adjust;comment:是否调节"`
	RatedPower          *float64 `json:"ratedPower" form:"ratedPower" gorm:"column:rated_power;comment:额定功率(kW)"`
	ActivePower         *float64 `json:"activePower" form:"activePower" gorm:"column:active_power;comment:有功功率"`
	ApTargetValue       *float64 `json:"apTargetValue" form:"apTargetValue" gorm:"column:ap_target_value;comment:有功目标(kW)"`
	ReactivePower       *float64 `json:"reactivePower" form:"reactivePower" gorm:"column:reactive_power;comment:无功功率(kVar)"`
	RpTargetValue       *float64 `json:"rpTargetValue" form:"rpTargetValue" gorm:"column:rp_target_value;comment:无功目标(kVar)"`
	PowerFactor         *float64 `json:"powerFactor" form:"powerFactor" gorm:"column:power_factor;comment:功率因数"`
}

// TableName agvcNbqHis表 AgvcNbqHis自定义表名 agvc_nbq_his
func (AgvcNbqHis) TableName() string {
	return "agvc_nbq_his"
}
