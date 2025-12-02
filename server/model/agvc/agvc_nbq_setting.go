// 自动生成模板AgvcNbqSetting
package agvc

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 逆变器配置 结构体  AgvcNbqSetting
type AgvcNbqSetting struct {
	global.GVA_MODEL
	Psid                *int     `json:"psid" form:"psid" gorm:"column:psid;"`                                                //电站编号
	InverterNo          *int     `json:"inverterNo" form:"inverterNo" gorm:"column:inverter_no;"`                             //逆变器编号
	Name                *string  `json:"name" form:"name" gorm:"column:name;"`                                                //逆变器名称
	InverterBrand       *int     `json:"inverterBrand" form:"inverterBrand" gorm:"column:inverter_brand;"`                    //逆变器品牌(1-华为,2-阳光,3-固德威)
	RatedActivePower    *float64 `json:"ratedActivePower" form:"ratedActivePower" gorm:"column:rated_active_power;"`          //额定有功功率
	RatedReactivePower  *float64 `json:"ratedReactivePower" form:"ratedReactivePower" gorm:"column:rated_reactive_power;"`    //额定无功功率
	JitterRange         *float64 `json:"jitterRange" form:"jitterRange" gorm:"column:jitter_range;"`                          //功率抖动区间
	DeadbandRange       *float64 `json:"deadbandRange" form:"deadbandRange" gorm:"column:deadband_range;"`                    //功率死区区间
	UpgradePriority     *float64 `json:"upgradePriority" form:"upgradePriority" gorm:"column:upgrade_priority;"`              //升额优先级
	DowngradePriority   *float64 `json:"downgradePriority" form:"downgradePriority" gorm:"column:downgrade_priority;"`        //降级优先级
	IsParticipateAdjust *bool    `json:"isParticipateAdjust" form:"isParticipateAdjust" gorm:"column:is_participate_adjust;"` //参与调节
	IsBenchmarkInverter *bool    `json:"isBenchmarkInverter" form:"isBenchmarkInverter" gorm:"column:is_benchmark_inverter;"` //标杆逆变器
	BwdNo               *int     `json:"bwdNo" form:"bwdNo" gorm:"column:bwd_no;"`                                            //并网点编号
}

// TableName 逆变器配置 AgvcNbqSetting自定义表名 agvc_nbq_setting
func (AgvcNbqSetting) TableName() string {
	return "agvc_nbq_setting"
}
