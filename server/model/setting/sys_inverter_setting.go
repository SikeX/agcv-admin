
// 自动生成模板SysInverterSetting
package setting
import (
	"time"
)

// 逆变器设置 结构体  SysInverterSetting
type SysInverterSetting struct {
  Id  *int64 `json:"id" form:"id" gorm:"primarykey;column:id;"`  //id字段
  InverterNo  *int64 `json:"inverterNo" form:"inverterNo" gorm:"column:inverter_no;"`  //逆变器编号
  Name  *string `json:"name" form:"name" gorm:"column:name;"`  //逆变器名称
  RatedActivePower  *float64 `json:"ratedActivePower" form:"ratedActivePower" gorm:"column:rated_active_power;"`  //额定有功功率
  RatedReactivePower  *float64 `json:"ratedReactivePower" form:"ratedReactivePower" gorm:"column:rated_reactive_power;"`  //额定无功功率
  JitterRange  *float64 `json:"jitterRange" form:"jitterRange" gorm:"column:jitter_range;"`  //功率抖动区间
  DeadbandRange  *float64 `json:"deadbandRange" form:"deadbandRange" gorm:"column:deadband_range;"`  //功率死区区间
  UpgradePriority  *float64 `json:"upgradePriority" form:"upgradePriority" gorm:"column:upgrade_priority;"`  //升额优先级
  DowngradePriority  *float64 `json:"downgradePriority" form:"downgradePriority" gorm:"column:downgrade_priority;"`  //降级优先级
  IsParticipateAdjust  *bool `json:"isParticipateAdjust" form:"isParticipateAdjust" gorm:"column:is_participate_adjust;"`  //是否参与调节
  IsBenchmarkInverter  *bool `json:"isBenchmarkInverter" form:"isBenchmarkInverter" gorm:"column:is_benchmark_inverter;"`  //是否为标杆逆变器
  CreatedAt  *time.Time `json:"createdAt" form:"createdAt" gorm:"column:created_at;"`  //创建时间
  UpdatedAt  *time.Time `json:"updatedAt" form:"updatedAt" gorm:"column:updated_at;"`  //更新时间
  DeletedAt  *time.Time `json:"deletedAt" form:"deletedAt" gorm:"column:deleted_at;"`  //删除时间
}


// TableName 逆变器设置 SysInverterSetting自定义表名 sys_inverter_setting
func (SysInverterSetting) TableName() string {
    return "sys_inverter_setting"
}





