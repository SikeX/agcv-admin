
// 自动生成模板SysSvgSvcSetting
package system
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// SVG/SVC设置 结构体  SysSvgSvcSetting
type SysSvgSvcSetting struct {
    global.GVA_MODEL
  WugongName  *string `json:"wugongName" form:"wugongName" gorm:"column:wugong_name;"`  //无功补偿装置名称
  WugongCapacity  *float64 `json:"wugongCapacity" form:"wugongCapacity" gorm:"column:wugong_capacity;"`  //无功容量
  CapacityCoefficient  *float64 `json:"capacityCoefficient" form:"capacityCoefficient" gorm:"column:capacity_coefficient;"`  //容量系数
  IsAdjustment  *string `json:"isAdjustment" form:"isAdjustment" gorm:"column:is_adjustment;"`  //参与调节
  Note  *string `json:"note" form:"note" gorm:"column:note;"`  //备注
}


// TableName SVG/SVC设置 SysSvgSvcSetting自定义表名 sys_svg_svc_setting
func (SysSvgSvcSetting) TableName() string {
    return "sys_svg_svc_setting"
}





