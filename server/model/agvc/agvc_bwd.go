// 自动生成模板AgvcBwdHis
package agvc

// agvcBwdHis表 结构体  AgvcBwd
// 用于存储和展示并网点实时数据和历史数据
type AgvcBwd struct {
	Psid                 *int     `json:"psid" form:"psid" point:"name:电站编号"`
	Number               *string  `json:"number" form:"number" point:"name:设备编号"`
	Name                 *string  `json:"name" form:"name" point:"name:并网点名称"`
	CurrentActivePower   *float64 `json:"currentActivePower" form:"currentActivePower" point:"name:当前有功(kW),value:10"`
	CurrentVoltage       *float64 `json:"currentVoltage" form:"currentVoltage" point:"name:当前电压(kV),value:11"`
	CurrentReactivePower *float64 `json:"currentReactivePower" form:"currentReactivePower" point:"name:当前无功(kVar),value:12"`
	SystemFrequency      *float64 `json:"systemFrequency" form:"systemFrequency" point:"name:系统频率(Hz),value:13"`
	SystemImpedance      *float64 `json:"systemImpedance" form:"systemImpedance" point:"name:系统阻抗,value:14"`
}

// TableName agvcBwdHis表 AgvcBwdHis自定义表名 agvc_bwd_his
func (AgvcBwd) TableName() string {
	return "agvc_bwd_his"
}
