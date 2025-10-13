package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type SysHisEventSearch struct{
    IsRead  *string `json:"isRead" form:"isRead"`
    DataType  *string `json:"dataType" form:"dataType"`
    HappenTimeRange []time.Time `json:"happenTimeRange" form:"happenTimeRange[]"`
    request.PageInfo
}