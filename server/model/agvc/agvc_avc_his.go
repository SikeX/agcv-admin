// 自动生成模板AgvcAvcHis
package agvc

// agvcAvcHis表 结构体  AgvcAvcHis
// 用于存储和展示AVC实时数据和历史数据
type AgvcAvcHis struct {
	Psid                *int     `json:"psid" form:"psid" point:"name:电站编号"`
	Number              *string  `json:"number" form:"number" point:"name:设备编号"`
	Name                *string  `json:"name" form:"name" point:"name:并网点名称"`
	TargetVoltage       *float64 `json:"targetVoltage" form:"targetVoltage" point:"name:目标电压(kV),value:30"`
	TargetReactivePower *float64 `json:"targetReactivePower" form:"targetReactivePower" point:"name:目标无功(kVar),value:31"`
	AvcFunctionState    *int64   `json:"avcFunctionState" form:"avcFunctionState" point:"name:AVC功能投退状态,value:32"`
	AvcControlMode      *int64   `json:"avcControlMode" form:"avcControlMode" point:"name:AVC调节方式,value:33"`
	AvcControlAuthority *int64   `json:"avcControlAuthority" form:"avcControlAuthority" point:"name:AVC控制权限,value:34"`
}

// TableName agvcAvcHis表 AgvcAvcHis自定义表名 agvc_avc_his
func (AgvcAvcHis) TableName() string {
	return "agvc_avc_his"
}
