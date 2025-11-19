package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type AgvcAgcHisSearch struct {
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	Psid           *int        `json:"psid" form:"psid"`
	Number         *string     `json:"number" form:"number"`
	Name           *string     `json:"name" form:"name"`
	request.PageInfo
}

// AgvcAgcHistoryRequest 历史数据查询请求
type AgvcAgcHistoryRequest struct {
	AgvcAgcHis agvc.AgvcAgc `json:"agvcAgcHis" form:"agvcAgcHis"`
	StartTime  string       `json:"startTime" form:"startTime" binding:"required"`
	EndTime    string       `json:"endTime" form:"endTime" binding:"required"`
}
