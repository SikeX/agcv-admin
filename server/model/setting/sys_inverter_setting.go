// 自动生成模板SysInverterSetting
package setting

import (
	"time"
)

// 逆变器设置 结构体  SysInverterSetting
type SysInverterSetting struct {
	Id                  *int64     `json:"id" form:"id" gorm:"primarykey;column:id;"`                                             //id字段
	InverterNo          *string    `json:"inverterNo" form:"inverterNo" gorm:"primarykey;column:inverter_no;" binding:"required"` //逆变器编号
	Name                *string    `json:"name" form:"name" gorm:"column:name;"`                                                  //逆变器名称
	RatedActivePower    *float64   `json:"ratedActivePower" form:"ratedActivePower" gorm:"column:rated_active_power;"`            //额定有功功率
	RatedReactivePower  *float64   `json:"ratedReactivePower" form:"ratedReactivePower" gorm:"column:rated_reactive_power;"`      //额定无功功率
	JitterRange         *float64   `json:"jitterRange" form:"jitterRange" gorm:"column:jitter_range;"`                            //功率抖动区间
	DeadbandRange       *float64   `json:"deadbandRange" form:"deadbandRange" gorm:"column:deadband_range;"`                      //功率死区区间
	UpgradePriority     *float64   `json:"upgradePriority" form:"upgradePriority" gorm:"column:upgrade_priority;"`                //升额优先级
	DowngradePriority   *float64   `json:"downgradePriority" form:"downgradePriority" gorm:"column:downgrade_priority;"`          //降级优先级
	IsParticipateAdjust *bool      `json:"isParticipateAdjust" form:"isParticipateAdjust" gorm:"column:is_participate_adjust;"`   //参与调节
	IsBenchmarkInverter *bool      `json:"isBenchmarkInverter" form:"isBenchmarkInverter" gorm:"column:is_benchmark_inverter;"`   //标杆逆变器
	CreatedAt           *time.Time `json:"createdAt" form:"createdAt" gorm:"column:created_at;"`                                  //创建时间
	UpdatedAt           *time.Time `json:"updatedAt" form:"updatedAt" gorm:"column:updated_at;"`                                  //更新时间
	DeletedAt           *time.Time `json:"deletedAt" form:"deletedAt" gorm:"column:deleted_at;"`                                  //删除时间
	YgPreal             *string    `json:"ygPreal" form:"ygPreal" gorm:"column:yg_preal;"`                                        //有功功率点号
	WgPreal             *string    `json:"wgPreal" form:"wgPreal" gorm:"column:wg_preal;"`                                        //无功功率点号
	YgPmax              *string    `json:"ygPmax" form:"ygPmax" gorm:"column:yg_pmax;"`                                           //有功调节上限点号
	YgPMin              *string    `json:"ygPMin" form:"ygPMin" gorm:"column:yg_p_min;"`                                          //有功调节下限点号
	WgPmax              *string    `json:"wgPmax" form:"wgPmax" gorm:"column:wg_pmax;"`                                           //无功调节上限点号
	WgPMin              *string    `json:"wgPMin" form:"wgPMin" gorm:"column:wg_p_min;"`                                          //无功调节下限点号
	Pf                  *string    `json:"pf" form:"pf" gorm:"column:pf;"`                                                        //功率因数点号
	YgCmd               *string    `json:"ygCmd" form:"ygCmd" gorm:"column:yg_cmd;"`                                              //有功调节模式点号
	YgMode              *string    `json:"ygMode" form:"ygMode" gorm:"column:yg_mode;"`                                           //有功调节模式点号
	YgRate              *string    `json:"ygRate" form:"ygRate" gorm:"column:yg_rate;"`                                           //有功功率变化梯度点号
	YgFixW              *string    `json:"ygFixW" form:"ygFixW" gorm:"column:yg_fix_w;"`                                          //有功固定值降额点号
	YgFixPercent        *string    `json:"ygFixPercent" form:"ygFixPercent" gorm:"column:yg_fix_percent;"`                        //有功百分比降额点号
	WgU                 *string    `json:"wgU" form:"wgU" gorm:"column:wg_u;"`                                                    //电网电压点号
	QuMode              *string    `json:"quMode" form:"quMode" gorm:"column:qu_mode;"`                                           //Q-U特征曲线模式点号
}

// TableName 逆变器设置 SysInverterSetting自定义表名 sys_inverter_setting
func (SysInverterSetting) TableName() string {
	return "sys_inverter_setting"
}
