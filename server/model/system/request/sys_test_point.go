
package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type SysTestPointSearch struct{
    InStorageName  string      `json:"inStorageName" form:"inStorageName"`     // 入库点名
    PointName      string      `json:"pointName" form:"pointName"`             // 测点名称
    PointType      *string     `json:"pointType" form:"pointType"`             // 测点类型
    DeviceType     *string     `json:"deviceType" form:"deviceType"`           // 设备类型
    CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
    request.PageInfo
}
