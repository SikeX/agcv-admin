package request

import (
    "time"

    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type AgvcNbqHisSearch struct {
    CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
    Psid           *int        `json:"psid" form:"psid"`
    Inverter_no    *string     `json:"inverter_no" form:"inverter_no"`
    Name           *string     `json:"name" form:"name"`
    request.PageInfo
}

// AgvcNbqHistoryRequest 历史数据查询请求
type AgvcNbqHistoryRequest struct {
    AgvcNbqHis agvc.AgvcNbqHis `json:"agvcNbqHis" form:"agvcNbqHis"`
    StartTime  string           `json:"startTime" form:"startTime" binding:"required"`
    EndTime    string           `json:"endTime" form:"endTime" binding:"required"`
}

// AgvcNbqRealDataRequest 实时数据查询请求
type AgvcNbqRealDataRequest struct {
    Psid       int `json:"psid" form:"psid" binding:"required"`
    InverterNo int `json:"inverterNo" form:"inverterNo" binding:"required"`
}
