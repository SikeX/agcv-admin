
package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type SysQixiangyiSettingSearch struct{
    ID  uint   `json:"ID" form:"ID"`
    DeviceName  string  `json:"deviceName" form:"deviceName"`
    request.PageInfo
}
