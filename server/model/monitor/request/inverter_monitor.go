package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type InverterMonitorSearch struct {
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	Inverter_no    *string     `json:"inverter_no" form:"inverter_no"`
	Name           *string     `json:"name" form:"name"`
	Status         *string     `json:"status" form:"status"`
	request.PageInfo
}
