// 自动生成模板AgvcQxySetting
package agvc
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 气象仪配置 结构体  AgvcQxySetting
type AgvcQxySetting struct {
    global.GVA_MODEL
  Number  *string `json:"number" form:"number" gorm:"column:number;"`  //设备编号
  DeviceName  *string `json:"deviceName" form:"deviceName" gorm:"column:device_name;"`  //设备名称
  DeviceType  *string `json:"deviceType" form:"deviceType" gorm:"column:device_type;"`  //设备类型
  DevicePosition  *string `json:"devicePosition" form:"devicePosition" gorm:"column:device_position;"`  //设备位置
  DeviceFactory  *string `json:"deviceFactory" form:"deviceFactory" gorm:"column:device_factory;"`  //设备厂家
  DeviceModel  *int64 `json:"deviceModel" form:"deviceModel" gorm:"column:device_model;"`  //设备型号
  GcpName  *string `json:"gcpName" form:"gcpName" gorm:"column:gcp_name;"`  //并网点
  InstallAngle  *float64 `json:"installAngle" form:"installAngle" gorm:"column:install_angle;"`  //安装角度
  IsMaster  *string `json:"isMaster" form:"isMaster" gorm:"column:is_master;"`  //是否主气象仪
}


// TableName 气象仪配置 AgvcQxySetting自定义表名 agvc_qxy_setting
func (AgvcQxySetting) TableName() string {
    return "agvc_qxy_setting"
}