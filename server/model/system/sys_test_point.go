
// 自动生成模板SysTestPoint
package system
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 测试管理 结构体  SysTestPoint
type SysTestPoint struct {
    global.GVA_MODEL
  InStorageName  *string `json:"inStorageName" form:"inStorageName" gorm:"column:in_storage_name;"`  //入库点名
  PointName  *string `json:"pointName" form:"pointName" gorm:"column:point_name;"`  //测点名称
  PointType  *string `json:"pointType" form:"pointType" gorm:"column:point_type;"`  //测点类型
  DeviceType  *string `json:"deviceType" form:"deviceType" gorm:"column:device_type;"`  //设备类型
  IsAllClose  *string `json:"isAllClose" form:"isAllClose" gorm:"column:is_all_close;"`  //是否全栈闭锁
  IsDeviceClose  *string `json:"isDeviceClose" form:"isDeviceClose" gorm:"column:is_device_close;"`  //是否设备闭锁
  ClosePosition  *int64 `json:"closePosition" form:"closePosition" gorm:"column:close_position;"`  //闭锁位置
  CloseRangeMax  *float64 `json:"closeRangeMax" form:"closeRangeMax" gorm:"column:close_range_max;"`  //闭锁上限
  CloseRangeMin  *float64 `json:"closeRangeMin" form:"closeRangeMin" gorm:"column:close_range_min;"`  //闭锁下限
}


// TableName 测试管理 SysTestPoint自定义表名 sys_test_point
func (SysTestPoint) TableName() string {
    return "sys_test_point"
}





