package inverter_brands

import (
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
)

// InverterBrandInterface 逆变器品牌接口
// 定义了所有品牌逆变器必须实现的控制方法
type InverterBrandInterface interface {
    // GetBrandCode 获取品牌代码
    GetBrandCode() int

    // GetBrandName 获取品牌名称
    GetBrandName() string

    // ControlAGCByPercentage 按百分比进行AGC调控
    // @param inverter 逆变器配置
    // @param targetPercentage 目标百分比 (0-100)
    ControlAGCByPercentage(inverter *agvc.AgvcNbqSetting, targetPercentage float64) error

    // ControlAGCByAbsolute 按绝对值进行AGC调控
    // @param inverter 逆变器配置
    // @param targetPowerKW 目标功率(kW)
    ControlAGCByAbsolute(inverter *agvc.AgvcNbqSetting, targetPowerKW float64) error

    // ControlAVCByPowerFactor 按功率因数进行AVC调控
    // @param inverter 逆变器配置
    // @param targetPF 目标功率因数 (-1到1)
    ControlAVCByPowerFactor(inverter *agvc.AgvcNbqSetting, targetPF float64) error

    // ControlAVCByReactivePower 按无功功率(Q/S)进行AVC调控
    // @param inverter 逆变器配置
    // @param targetQS 目标Q/S (-1到1)
    ControlAVCByReactivePower(inverter *agvc.AgvcNbqSetting, targetQS float64) error

    // InitializePoints 初始化品牌特定的点位映射到数据库
    InitializePoints() error
}
