package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// DeviceSearch 设备搜索条件
type DeviceSearch struct {
	request.PageInfo
	PSID     string `json:"psid" form:"psid"`
	EQID     string `json:"eqid" form:"eqid"`
	EQType   string `json:"eqType" form:"eqType"`
	DataType string `json:"dataType" form:"dataType"`
	Status   *int   `json:"status" form:"status"`
}

// RealtimeDataRequest 实时数据查询请求
type RealtimeDataRequest struct {
	PSID     string   `json:"psid" form:"psid" binding:"required"`
	EQID     string   `json:"eqid" form:"eqid" binding:"required"`
	EQType   string   `json:"eqType" form:"eqType" binding:"required"`
	DataType string   `json:"dataType" form:"dataType"`
	Points   []string `json:"points" form:"points"` // 可选，指定查询的点标识
}

// HistoryDataRequest 历史数据查询请求
type HistoryDataRequest struct {
	PSID      string   `json:"psid" form:"psid" binding:"required"`
	EQID      string   `json:"eqid" form:"eqid" binding:"required"`
	EQType    string   `json:"eqType" form:"eqType" binding:"required"`
	DataType  string   `json:"dataType" form:"dataType"`
	Points    []string `json:"points" form:"points" binding:"required"`
	StartTime int64    `json:"startTime" form:"startTime" binding:"required"` // Unix时间戳（秒）
	EndTime   int64    `json:"endTime" form:"endTime" binding:"required"`     // Unix时间戳（秒）
	Interval  string   `json:"interval" form:"interval"`                      // 聚合间隔，如"1m", "5m", "1h"
}

// AGCConfigUpdate AGC配置更新请求
type AGCConfigUpdate struct {
	PSID              string  `json:"psid" binding:"required"`
	IsActive          *int    `json:"isActive"`
	ControlAuth       *int    `json:"controlAuth"`
	RunMode           *int    `json:"runMode"`
	JitterRange       float64 `json:"jitterRange"`
	RegPeriod         int     `json:"regPeriod"`
	RegStep           float64 `json:"regStep"`
	PowerUpperLimit   float64 `json:"powerUpperLimit"`
	PowerLowerLimit   float64 `json:"powerLowerLimit"`
	PowerExecValue    float64 `json:"powerExecValue"`
	UpRegLock         *int    `json:"upRegLock"`
	DownRegLock       *int    `json:"downRegLock"`
	DispatchExecValue float64 `json:"dispatchExecValue"`
	StationExecValue  float64 `json:"stationExecValue"`
}

// AVCConfigUpdate AVC配置更新请求
type AVCConfigUpdate struct {
	PSID                string  `json:"psid" binding:"required"`
	IsActive            *int    `json:"isActive"`
	ControlAuth         *int    `json:"controlAuth"`
	CommandStatus       *int    `json:"commandStatus"`
	RunMode             *int    `json:"runMode"`
	TargetVoltageLow    float64 `json:"targetVoltageLow"`
	TargetVoltageHigh   float64 `json:"targetVoltageHigh"`
	VoltageExecValue    float64 `json:"voltageExecValue"`
	ReactiveExecValue   float64 `json:"reactiveExecValue"`
	ReactiveIncCap      float64 `json:"reactiveIncCap"`
	ReactiveDecCap      float64 `json:"reactiveDecCap"`
	UpRegLock           *int    `json:"upRegLock"`
	DownRegLock         *int    `json:"downRegLock"`
	RegPeriod           int     `json:"regPeriod"`
	VoltageDeadZone     float64 `json:"voltageDeadZone"`
	ReactiveSensitivity float64 `json:"reactiveSensitivity"`
}

// AGCRegulationRecordSearch AGC调节记录搜索
type AGCRegulationRecordSearch struct {
	request.PageInfo
	PSID      string `json:"psid" form:"psid"`
	Status    string `json:"status" form:"status"`
	StartTime *int64 `json:"startTime" form:"startTime"`
	EndTime   *int64 `json:"endTime" form:"endTime"`
}

// AVCRegulationRecordSearch AVC调节记录搜索
type AVCRegulationRecordSearch struct {
	request.PageInfo
	PSID      string `json:"psid" form:"psid"`
	Status    string `json:"status" form:"status"`
	StartTime *int64 `json:"startTime" form:"startTime"`
	EndTime   *int64 `json:"endTime" form:"endTime"`
}

// CoAPDataMessage CoAP数据消息（用于接收和发送）
type CoAPDataMessage struct {
	PSID     int         `json:"psid"`
	EQID     int         `json:"eqid" binding:"required"`
	EQType   int         `json:"eqType" binding:"required"`
	DataType int         `json:"dataType" binding:"required"`
	Point    string      `json:"point" binding:"required"`
	Value    interface{} `json:"value" binding:"required"`
}
