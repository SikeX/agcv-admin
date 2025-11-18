package request

import (
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type AgvcNbqSettingSearch struct {
    Psid                *int    `json:"psid" form:"psid"`
    InverterNo          *string `json:"inverterNo" form:"inverterNo"`
    Name                *string `json:"name" form:"name"`
    IsParticipateAdjust *bool   `json:"isParticipateAdjust" form:"isParticipateAdjust"`
    IsBenchmarkInverter *bool   `json:"isBenchmarkInverter" form:"isBenchmarkInverter"`
    request.PageInfo
    Sort  string `json:"sort" form:"sort"`
    Order string `json:"order" form:"order"`
}
