package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type AgvcBwdHisSearch struct {
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	Psid           *int        `json:"psid" form:"psid"`
	Number         *string     `json:"number" form:"number"`
	Name           *string     `json:"name" form:"name"`
	request.PageInfo
}

// AgvcBwdHistoryRequest 历史数据查询请求
type AgvcBwdHistoryRequest struct {
	AgvcBwdHis agvc.AgvcBwd `json:"agvcBwdHis" form:"agvcBwdHis"`
	StartTime  string       `json:"startTime" form:"startTime" binding:"required"`
	EndTime    string       `json:"endTime" form:"endTime" binding:"required"`
}

// 计划曲线点
type PlanCurvePoint struct {
	Time        string  `json:"time"`        // 时间 HH:mm:ss
	TargetValue float64 `json:"targetValue"` // 目标值
}

// AGC/AVC计划曲线（只保留本地曲线）
type PlanCurves struct {
	LocalCurve []PlanCurvePoint `json:"localCurve"` // 本地曲线
}

// AGC参数设置请求
type AgcParametersRequest struct {
	Number               string  `json:"number" binding:"required"` // 并网点编号
	ControlMethod        string  `json:"controlMethod"`             // 调节方式
	AdjustmentMode       string  `json:"adjustmentMode"`            // 调节模式
	StepSize             float64 `json:"stepSize"`                  // 调节步长(kW)
	StepPeriod           int     `json:"stepPeriod"`                // 步长周期(秒)
	VibrationRange       float64 `json:"vibrationRange"`            // 抖动区间(kW)
	ControlPeriod        int     `json:"controlPeriod"`             // 调控周期(秒)
	MicroAdjustmentCoeff float64 `json:"microAdjustmentCoeff"`      // 微调系数
}

// AVC参数设置请求
type AvcParametersRequest struct {
	Number             string  `json:"number" binding:"required"` // 并网点编号
	ControlMethod      string  `json:"controlMethod"`             // 调节方式
	StepSize           float64 `json:"stepSize"`                  // 调节步长(kV)
	StepPeriod         int     `json:"stepPeriod"`                // 步长周期(秒)
	VibrationRange     float64 `json:"vibrationRange"`            // 抖动区间(kV)
	ControlPeriod      int     `json:"controlPeriod"`             // 调控周期(秒)
	SystemImpedance    float64 `json:"systemImpedance"`           // 系统阻抗
	AdjustmentRangeMin float64 `json:"adjustmentRangeMin"`        // 调节范围最小值(kV)
	AdjustmentRangeMax float64 `json:"adjustmentRangeMax"`        // 调节范围最大值(kV)
}

// 计划曲线请求
type PlanCurvesRequest struct {
	Number     string     `json:"number" binding:"required"` // 并网点编号
	CurveType  string     `json:"curveType"`                 // 曲线类型: local/dispatch
	PlanCurves PlanCurves `json:"planCurves"`                // 计划曲线数据
}

// AGC状态请求
type AgcStatusRequest struct {
	Number              string `json:"number" binding:"required"` // 并网点编号
	AgcFunctionState    int64  `json:"agcFunctionState"`          // AGC功能投退状态(1:投入,0:退出)
	AgcControlMode      int64  `json:"agcControlMode"`            // AGC调节方式(1:闭环指导,0:开环指导)
	AgcControlAuthority int64  `json:"agcControlAuthority"`       // AGC控制权限(1:调度控制,0:站内控制)
}

// AVC状态请求
type AvcStatusRequest struct {
	Number              string `json:"number" binding:"required"` // 并网点编号
	AvcFunctionState    int64  `json:"avcFunctionState"`          // AVC功能投退状态(1:投入,0:退出)
	AvcControlMode      int64  `json:"avcControlMode"`            // AVC调节方式(1:开环指导,0:闭环调节)
	AvcControlAuthority int64  `json:"avcControlAuthority"`       // AVC控制权限(1:站内控制,0:调度控制)
}
