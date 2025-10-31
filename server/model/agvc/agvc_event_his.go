// 自动生成模板AgvcEventHis
package agvc

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 历史事件 结构体  AgvcEventHis
type AgvcEventHis struct {
	global.GVA_MODEL
	Psid       *string    `json:"psid" form:"psid" gorm:"column:psid;"`                    //电站名称
	Eqid       *string    `json:"eqid" form:"eqid" gorm:"column:eqid;"`                    //设备编号
	EqType     *string    `json:"eqType" form:"eqType" gorm:"column:eq_type;"`             //设备类型
	DataType   *string    `json:"dataType" form:"dataType" gorm:"column:data_type;"`       //数据类型
	Datapoint  *int64     `json:"datapoint" form:"datapoint" gorm:"column:datapoint;"`     //点号
	RecordTime *time.Time `json:"RecordTime" form:"RecordTime" gorm:"column:record_time;"` //发生时间/恢复时间
	DataValue  *string    `json:"dataValue" form:"dataValue" gorm:"column:data_value;"`    //数值 (0或1)
	IsRead     *string    `json:"isRead" form:"isRead" gorm:"column:is_read;"`             //已读 (0或1)
}

// TableName 历史事件 AgvcEventHis自定义表名 agvc_event_his
func (AgvcEventHis) TableName() string {
	return "agvc_event_his"
}
