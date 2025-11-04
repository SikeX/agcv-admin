package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type AgvcNbqHisSearch struct {
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	Inverter_no    *string     `json:"inverter_no" form:"inverter_no"`
	Name           *string     `json:"name" form:"name"`
	request.PageInfo
}
