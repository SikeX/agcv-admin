
package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type AgvcBwdSettingSearch struct {
	Name string `json:"name" form:"name"`
	Number string `json:"number" form:"number"`
	request.PageInfo
}
