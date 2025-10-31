package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type AgvcEventHisSearch struct {
	Psid            string      `json:"psid" form:"psid"`                         // 电站名称
	EqType          string      `json:"eqType" form:"eqType"`                     // 设备类型
	Eqid            string      `json:"eqid" form:"eqid"`                         // 设备编号
	RecordTimeRange []time.Time `json:"recordTimeRange" form:"recordTimeRange[]"` // 发生时间范围
	request.PageInfo
}
