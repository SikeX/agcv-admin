// 自动生成模板AgvcAgcHis
package agvc

// agvcAgcHis表 结构体  AgvcAgcHis
// 用于存储和展示AGC实时数据和历史数据
type AgvcAgcHis struct {
	Psid                  *int     `json:"psid" form:"psid" point:"name:电站编号"`
	Number                *string  `json:"number" form:"number" point:"name:设备编号"`
	Name                  *string  `json:"name" form:"name" point:"name:并网点名称"`
	TargetActivePower     *float64 `json:"targetActivePower" form:"targetActivePower" point:"name:目标有功(kW),value:20"`
	DispatchActivePower   *float64 `json:"dispatchActivePower" form:"dispatchActivePower" point:"name:调度下发(kW),value:21"`
	AgcFunctionState      *int64   `json:"agcFunctionState" form:"agcFunctionState" point:"name:AGC功能投退状态,value:22"`
	AgcControlMode        *int64   `json:"agcControlMode" form:"agcControlMode" point:"name:AGC调节方式,value:23"`
	AgcControlAuthority   *int64   `json:"agcControlAuthority" form:"agcControlAuthority" point:"name:AGC控制权限,value:24"`
	AgcUpperLimit         *float64 `json:"agcUpperLimit" form:"agcUpperLimit" point:"name:可调上限(kW),value:25"`
	AgcLowerLimit         *float64 `json:"agcLowerLimit" form:"agcLowerLimit" point:"name:可调下限(kW),value:26"`
}

// TableName agvcAgcHis表 AgvcAgcHis自定义表名 agvc_agc_his
func (AgvcAgcHis) TableName() string {
	return "agvc_agc_his"
}
