package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type SysGridConnectionPointSearch struct {
	Name string `json:"name" form:"name"`
	request.PageInfo
}
