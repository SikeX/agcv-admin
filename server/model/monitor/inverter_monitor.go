// 自动生成模板InverterMonitor
package monitor

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 逆变器监控 结构体  InverterMonitor
type InverterMonitor struct {
	global.GVA_MODEL
	Inverter_no         *string  `json:"inverter_no" form:"inverter_no" gorm:"column:inverter_no;"`                           //逆变器编号
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
	Time         time.Time `json:"_time"`                                                          // 时间戳
	YgPreal      float64   `json:"ygPreal" form:"ygPreal" gorm:"column:yg_preal;"`                 //有功功率点号
	WgPreal      float64   `json:"wgPreal" form:"wgPreal" gorm:"column:wg_preal;"`                 //无功功率点号
	YgPmax       float64   `json:"ygPmax" form:"ygPmax" gorm:"column:yg_pmax;"`                    //有功调节上限点号
	YgPMin       float64   `json:"ygPMin" form:"ygPMin" gorm:"column:yg_p_min;"`                   //有功调节下限点号
	WgPmax       float64   `json:"wgPmax" form:"wgPmax" gorm:"column:wg_pmax;"`                    //无功调节上限点号
	WgPMin       float64   `json:"wgPMin" form:"wgPMin" gorm:"column:wg_p_min;"`                   //无功调节下限点号
	Pf           float64   `json:"pf" form:"pf" gorm:"column:pf;"`                                 //功率因数点号
	YgCmd        float64   `json:"ygCmd" form:"ygCmd" gorm:"column:yg_cmd;"`                       //有功调节模式点号
	YgMode       float64   `json:"ygMode" form:"ygMode" gorm:"column:yg_mode;"`                    //有功调节模式点号
	YgRate       float64   `json:"ygRate" form:"ygRate" gorm:"column:yg_rate;"`                    //有功功率变化梯度点号
	YgFixW       float64   `json:"ygFixW" form:"ygFixW" gorm:"column:yg_fix_w;"`                   //有功固定值降额点号
	YgFixPercent float64   `json:"ygFixPercent" form:"ygFixPercent" gorm:"column:yg_fix_percent;"` //有功百分比降额点号
	WgU          float64   `json:"wgU" form:"wgU" gorm:"column:wg_u;"`                             //电网电压点号
}
