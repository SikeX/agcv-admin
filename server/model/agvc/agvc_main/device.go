package agvc_main

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Device 设备信息 结构体
type Device struct {
	global.GVA_MODEL
	PSID       string  `json:"psid" form:"psid" gorm:"column:psid;comment:电站ID(3位);size:3;index"`
	EQID       string  `json:"eqid" form:"eqid" gorm:"column:eqid;comment:设备ID(4位);size:4;index"`
	EQType     string  `json:"eqType" form:"eqType" gorm:"column:eq_type;comment:设备类型(2位);size:2;index"` // 01:逆变器 02:并网点
	DataType   string  `json:"dataType" form:"dataType" gorm:"column:data_type;comment:数据类型(2位);size:2"` // 01:遥信 02:遥测 03:遥控 04:遥调
	Name       string  `json:"name" form:"name" gorm:"column:name;comment:设备名称"`
	DeviceCode string  `json:"deviceCode" form:"deviceCode" gorm:"column:device_code;comment:设备编号;uniqueIndex"` // psid+eqid+dataType+eqType
	Status     *int    `json:"status" form:"status" gorm:"column:status;comment:设备状态;default:1"`                // 1:在线 0:离线
	MaxPower   float64 `json:"maxPower" form:"maxPower" gorm:"column:max_power;comment:最大功率(MW)"`
	MaxReg     float64 `json:"maxReg" form:"maxReg" gorm:"column:max_reg;comment:最大调节能力(MW)"`
	MaxReact   float64 `json:"maxReact" form:"maxReact" gorm:"column:max_react;comment:最大无功容量(MVar)"`
	Remark     string  `json:"remark" form:"remark" gorm:"column:remark;comment:备注"`
}

// TableName Device表名
func (Device) TableName() string {
	return "agvc_device"
}

// RealtimeData 实时数据存储结构（用于内存Map）
type RealtimeData struct {
	PSID      string      `json:"psid"`
	EQID      string      `json:"eqid"`
	EQType    string      `json:"eqType"`
	DataType  string      `json:"dataType"`
	Point     string      `json:"point"`
	Value     interface{} `json:"value"`
	Timestamp int64       `json:"timestamp"` // 数据时间戳
}

// AgvcDataItem AGVC数据项（从CoAP接收的数据格式）
type AgvcDataItem struct {
	Psid     int     `json:"psid"`     // 电站ID
	Eqid     int     `json:"eqid"`     // 设备ID
	EqType   int     `json:"eqType"`   // 设备类型
	DataType int     `json:"dataType"` // 数据类型
	Point    string  `json:"point"`    // 数据点
	Value    float64 `json:"value"`    // 数值
}

// DeviceDataKey 设备数据键（用于Map索引）
type DeviceDataKey struct {
	PSID     string
	EQID     string
	EQType   string
	DataType string
	Point    string
}

// PointMapping 测点映射表
type PointMapping struct {
	global.GVA_MODEL
	EQType      string `json:"eqType" form:"eqType" gorm:"column:eq_type;comment:设备类型;index"`       // 01:逆变器 02:并网点
	DataType    string `json:"dataType" form:"dataType" gorm:"column:data_type;comment:数据类型;index"` // 01:遥信 02:遥测 03:遥控 04:遥调
	Point       string `json:"point" form:"point" gorm:"column:point;comment:点标识;index"`            // 401, 402等
	PointName   string `json:"pointName" form:"pointName" gorm:"column:point_name;comment:点名称"`     // 如：AGC投退信号
	Unit        string `json:"unit" form:"unit" gorm:"column:unit;comment:单位"`                      // MW, kV等
	Description string `json:"description" form:"description" gorm:"column:description;comment:描述"` // 详细描述
	Category    string `json:"category" form:"category" gorm:"column:category;comment:分类"`          // AGC, AVC, INVERTER等
}

// TableName PointMapping表名
func (PointMapping) TableName() string {
	return "agvc_point_mapping"
}
