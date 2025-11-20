package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type AgvcNbqHisSearch struct {
	StartTime   string  `json:"startTime" form:"startTime" binding:"required"`
	EndTime     string  `json:"endTime" form:"endTime" binding:"required"`
	Psid        *int    `json:"psid" form:"psid"`
	Inverter_no *string `json:"inverter_no" form:"inverter_no"`
	Name        *string `json:"name" form:"name"`
	request.PageInfo
}

// AgvcNbqHistoryRequest 历史数据查询请求
type AgvcNbqHistoryRequest struct {
	AgvcNbqHis agvc.AgvcNbq `json:"agvcNbqHis" form:"agvcNbqHis"`
	StartTime  string       `json:"startTime" form:"startTime" binding:"required"`
	EndTime    string       `json:"endTime" form:"endTime" binding:"required"`
}

// AgvcNbqRealDataRequest 实时数据查询请求
type AgvcNbqRealDataRequest struct {
	Psid       *int    `json:"psid" form:"psid"`
	InverterNo *int    `json:"inverter_no" form:"inverter_no"`
	Name       *string `json:"name" form:"name"`
	request.PageInfo
}
