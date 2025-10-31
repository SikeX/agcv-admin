package agvc

// AgvcDataItem AGVC数据项
type AgvcDataItem struct {
	Psid     int     `json:"psid"`     // 电站ID
	Eqid     int     `json:"eqid"`     // 设备ID
	EqType   int     `json:"eqType"`   // 设备类型
	DataType int     `json:"dataType"` // 数据类型
	Point    string  `json:"point"`    // 数据点
	Value    float64 `json:"value"`    // 数值
}

// AgvcDataBatch AGVC数据批量接收结构
type AgvcDataBatch []AgvcDataItem
