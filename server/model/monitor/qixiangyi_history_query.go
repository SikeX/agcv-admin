package monitor

// QixiangyiHistoryQuery 气象仪历史数据查询参数
type QixiangyiHistoryQuery struct {
	QixiangyiID int    `json:"qixiangyiID" form:"qixiangyiID" binding:"required"` // 气象仪ID
	StartTime   string `json:"startTime" form:"startTime" binding:"required"`    // 开始时间
	EndTime     string `json:"endTime" form:"endTime" binding:"required"`        // 结束时间
}

// QixiangyiHistoryData 气象仪历史数据点
type QixiangyiHistoryData struct {
	Time  string  `json:"time"`  // 时间戳
	Value float64 `json:"value"` // 数值
}

// QixiangyiHistoryResponse 气象仪历史数据响应
type QixiangyiHistoryResponse struct {
	PointID   int                        `json:"pointId"`   // 点标识
	PointName string                     `json:"pointName"` // 数据点名称
	Data      []QixiangyiHistoryData     `json:"data"`      // 历史数据
}