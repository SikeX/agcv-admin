
// 自动生成模板SysHisEvent
package system
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

// 历史事件 结构体  SysHisEvent
type SysHisEvent struct {
    global.GVA_MODEL
  Psid  *int64 `json:"psid" form:"psid" gorm:"column:psid;"`  //电站名称
  Eqid  *int64 `json:"eqid" form:"eqid" gorm:"column:eqid;"`  //设备编号
  EqType  *int64 `json:"eqType" form:"eqType" gorm:"column:eq_type;"`  //设备类型
  DataType  *string `json:"dataType" form:"dataType" gorm:"column:data_type;"`  //数据类型
  Datapoint  *string `json:"datapoint" form:"datapoint" gorm:"column:datapoint;"`  //点号
  HappenTime  *time.Time `json:"happenTime" form:"happenTime" gorm:"column:happen_time;"`  //发生时间
  DataValue  *string `json:"dataValue" form:"dataValue" gorm:"column:data_value;"`  //数值
  IsRead  *string `json:"isRead" form:"isRead" gorm:"column:is_read;"`  //已读
  Name  *string `json:"name" form:"name" gorm:"column:name;"`  //名称
  RecoverTime  *time.Time `json:"recoverTime" form:"recoverTime" gorm:"column:recover_time;"`  //恢复时间
}


// TableName 历史事件 SysHisEvent自定义表名 sys_his_event
func (SysHisEvent) TableName() string {
    return "sys_his_event"
}





