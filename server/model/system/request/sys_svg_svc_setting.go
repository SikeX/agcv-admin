
package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type SysSvgSvcSettingSearch struct{
    ID  *uint   `json:"ID" form:"ID"`
    WugongName  *string `json:"wugongName" form:"wugongName"`
    IsAdjustment  *string `json:"isAdjustment" form:"isAdjustment"`
    request.PageInfo
}
