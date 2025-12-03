package request

// AgvcPowerChartRequest 电站出力图表查询请求
type AgvcPowerChartRequest struct {
	Psid      int    `json:"psid" form:"psid" `
	Eqid      int    `json:"eqid" form:"eqid" binding:"required"`
	StartTime string `json:"startTime" form:"startTime" binding:"required"`
	EndTime   string `json:"endTime" form:"endTime" binding:"required"`
}

// AgvcVoltageReactiveChartRequest 电压和无功出力图表查询请求
type AgvcVoltageReactiveChartRequest struct {
	Psid      int    `json:"psid" form:"psid" `
	Eqid      int    `json:"eqid" form:"eqid" binding:"required"`
	StartTime string `json:"startTime" form:"startTime" binding:"required"`
	EndTime   string `json:"endTime" form:"endTime" binding:"required"`
}
