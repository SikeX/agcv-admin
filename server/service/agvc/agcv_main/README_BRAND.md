# 逆变器品牌接口实现指南

## 概述

本系统采用面向接口编程的设计模式，将不同品牌的逆变器控制逻辑抽象为统一的接口。这样做的好处是：

1. **扩展性强**：添加新品牌只需实现接口，无需修改现有代码
2. **维护性好**：每个品牌的实现相互独立，互不影响
3. **代码清晰**：通过品牌工厂统一管理，调用逻辑简洁明了

## 架构说明

### 核心组件

1. **InverterBrandInterface** (`inverter_brand_interface.go`)
   - 定义了所有品牌逆变器必须实现的方法接口
   - 包括AGC控制、AVC控制和点位初始化等方法

2. **InverterBrandFactory** (`inverter_brand_factory.go`)
   - 品牌工厂，负责注册和管理所有品牌实现
   - 提供统一的品牌获取接口

3. **品牌实现** (如 `huawei_brand_impl.go`)
   - 各品牌的具体实现文件
   - 实现 InverterBrandInterface 接口的所有方法

## 接口定义

```go
type InverterBrandInterface interface {
    // GetBrandCode 获取品牌代码
    GetBrandCode() int

    // GetBrandName 获取品牌名称
    GetBrandName() string

    // ControlAGCByPercentage 按百分比进行AGC调控
    ControlAGCByPercentage(inverter *agvc.AgvcNbqSetting, targetPercentage float64) error

    // ControlAGCByAbsolute 按绝对值进行AGC调控
    ControlAGCByAbsolute(inverter *agvc.AgvcNbqSetting, targetPowerKW float64) error

    // ControlAVCByPowerFactor 按功率因数进行AVC调控
    ControlAVCByPowerFactor(inverter *agvc.AgvcNbqSetting, targetPF float64) error

    // ControlAVCByReactivePower 按无功功率(Q/S)进行AVC调控
    ControlAVCByReactivePower(inverter *agvc.AgvcNbqSetting, targetQS float64) error

    // InitializePoints 初始化品牌特定的点位映射到数据库
    InitializePoints() error
}
```

## 添加新品牌步骤

### 步骤1：定义品牌常量

在 `server/service/agvc/cons/def.go` 文件中添加品牌常量：

```go
const (
    INVERTER_BRAND_HUAWEI  = 1 // 华为
    INVERTER_BRAND_SUNGROW = 2 // 阳光
    INVERTER_BRAND_GOODWE  = 3 // 固德威
    INVERTER_BRAND_YOUR_BRAND = 4 // 你的品牌名称
)
```

### 步骤2：创建品牌实现文件

创建新文件 `your_brand_impl.go`，实现品牌接口：

```go
package agcv_main

import (
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
    "github.com/flipped-aurora/gin-vue-admin/server/service/agvc/cons"
)

// YourBrandInverter 你的品牌逆变器实现
type YourBrandInverter struct{}

// GetBrandCode 获取品牌代码
func (y *YourBrandInverter) GetBrandCode() int {
    return cons.INVERTER_BRAND_YOUR_BRAND
}

// GetBrandName 获取品牌名称
func (y *YourBrandInverter) GetBrandName() string {
    return "你的品牌名称"
}

// ControlAGCByPercentage 按百分比进行AGC调控
func (y *YourBrandInverter) ControlAGCByPercentage(inverter *agvc.AgvcNbqSetting, targetPercentage float64) error {
    // 实现你的品牌的AGC百分比控制逻辑
    // 可以创建专门的控制器，如 YourBrandController
    return YourBrandController.ControlAGCByPercentage(inverter, targetPercentage)
}

// ControlAGCByAbsolute 按绝对值进行AGC调控
func (y *YourBrandInverter) ControlAGCByAbsolute(inverter *agvc.AgvcNbqSetting, targetPowerKW float64) error {
    // 实现你的品牌的AGC绝对值控制逻辑
    return YourBrandController.ControlAGCByAbsolute(inverter, targetPowerKW)
}

// ControlAVCByPowerFactor 按功率因数进行AVC调控
func (y *YourBrandInverter) ControlAVCByPowerFactor(inverter *agvc.AgvcNbqSetting, targetPF float64) error {
    // 实现你的品牌的AVC功率因数控制逻辑
    return YourBrandController.ControlAVCByPowerFactor(inverter, targetPF)
}

// ControlAVCByReactivePower 按无功功率(Q/S)进行AVC调控
func (y *YourBrandInverter) ControlAVCByReactivePower(inverter *agvc.AgvcNbqSetting, targetQS float64) error {
    // 实现你的品牌的AVC无功功率控制逻辑
    return YourBrandController.ControlAVCByReactivePower(inverter, targetQS)
}

// InitializePoints 初始化品牌特定的点位映射
func (y *YourBrandInverter) InitializePoints() error {
    // 实现你的品牌的点位初始化逻辑
    return YourBrandPointInit.InitializePoints()
}
```

### 步骤3：创建品牌控制器（可选）

如果需要复杂的控制逻辑，可以创建专门的控制器文件 `your_brand_controller.go`：

```go
package agcv_main

import (
    "fmt"
    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
    "go.uber.org/zap"
)

// YourBrandController 你的品牌逆变器控制器
type yourBrandController struct{}

var YourBrandController = &yourBrandController{}

// ControlAGCByPercentage 按百分比进行AGC调控
func (c *yourBrandController) ControlAGCByPercentage(inverter *agvc.AgvcNbqSetting, targetPercentage float64) error {
    if inverter.InverterNo == nil {
        return fmt.Errorf("逆变器编号为空")
    }

    invNo := *inverter.InverterNo

    // 限制百分比范围 0-100%
    if targetPercentage < 0 {
        targetPercentage = 0
    } else if targetPercentage > 100 {
        targetPercentage = 100
    }

    global.GVA_LOG.Info("你的品牌逆变器AGC百分比调控",
        zap.Int("逆变器编号", invNo),
        zap.Float64("目标百分比", targetPercentage))

    // 实现具体的控制逻辑
    // 1. 构建CoAP命令
    // 2. 发送命令到逆变器
    // 3. 记录日志

    return nil
}

// 其他控制方法的实现...
```

### 步骤4：创建点位初始化器（可选）

创建 `your_brand_point_initializer.go` 文件：

```go
package agcv_main

import (
    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/agvc_main"
    "go.uber.org/zap"
)

// YourBrandPointInitializer 你的品牌点位初始化器
type yourBrandPointInitializer struct{}

var YourBrandPointInit = &yourBrandPointInitializer{}

// InitializePoints 初始化你的品牌点位映射
func (i *yourBrandPointInitializer) InitializePoints() error {
    global.GVA_LOG.Info("开始初始化你的品牌点位映射...")

    pointMappings := []agvc_main.PointMapping{
        {
            EQType:      "02",
            DataType:    "05",
            Point:       "601",
            PointName:   "有功功率设定",
            Unit:        "kW",
            Description: "有功功率设定点",
            Category:    "YOUR_BRAND_AGC",
        },
        // 添加更多点位映射...
    }

    // 批量插入或更新点位映射
    for _, mapping := range pointMappings {
        // 实现插入/更新逻辑
    }

    return nil
}
```

### 步骤5：注册品牌到工厂

在 `inverter_brand_factory.go` 的 `init()` 函数中注册新品牌：

```go
func init() {
    // 注册华为品牌
    BrandFactory.RegisterBrand(&HuaweiInverterBrand{})

    // 注册你的品牌
    BrandFactory.RegisterBrand(&YourBrandInverter{})

    global.GVA_LOG.Info("逆变器品牌工厂初始化完成")
}
```

### 步骤6：更新品牌映射器（如需要）

如果需要在 `inverter_brand_mapper.go` 中添加点位映射：

```go
func (m *inverterBrandMapper) Initialize() {
    // 华为逆变器点位映射
    m.mappings[cons.INVERTER_BRAND_HUAWEI] = InverterBrandPointMapping{
        Brand:              cons.INVERTER_BRAND_HUAWEI,
        BrandName:          "华为",
        PowerSwitchYK:      "401",
        ActivePowerLimitYT: "401",
        ReactivePowerYT:    "402",
    }

    // 你的品牌点位映射
    m.mappings[cons.INVERTER_BRAND_YOUR_BRAND] = InverterBrandPointMapping{
        Brand:              cons.INVERTER_BRAND_YOUR_BRAND,
        BrandName:          "你的品牌名称",
        PowerSwitchYK:      "601",
        ActivePowerLimitYT: "601",
        ReactivePowerYT:    "602",
    }
}
```

## 使用示例

添加完成后，系统会自动通过品牌工厂获取对应的品牌实现：

```go
// 在 agc.go 或 avc.go 中的使用
brandCode := cons.INVERTER_BRAND_YOUR_BRAND
brandImpl, err := BrandFactory.GetBrand(brandCode)
if err != nil {
    // 处理错误
    return err
}

// 调用AGC控制
err = brandImpl.ControlAGCByPercentage(&inverter, 80.0)

// 调用AVC控制
err = brandImpl.ControlAVCByPowerFactor(&inverter, 0.95)
```

## 华为品牌实现参考

可以参考现有的华为品牌实现：

- `huawei_brand_impl.go` - 华为品牌接口实现
- `huawei_inverter_controller.go` - 华为逆变器控制器
- `huawei_point_initializer.go` - 华为点位初始化器

## 注意事项

1. **品牌代码唯一性**：确保品牌代码在 `cons/def.go` 中是唯一的
2. **接口完整性**：必须实现接口定义的所有方法
3. **错误处理**：妥善处理各种异常情况，返回有意义的错误信息
4. **日志记录**：在关键步骤添加日志，便于调试和监控
5. **参数校验**：对输入参数进行有效性检查
6. **线程安全**：如果有共享状态，注意并发安全问题

## 测试建议

1. 单元测试：为每个控制方法编写单元测试
2. 集成测试：在测试环境中验证与实际逆变器的通信
3. 压力测试：测试在高并发场景下的稳定性
4. 边界测试：测试参数边界情况的处理

## 总结

通过实现品牌接口，你可以轻松地为系统添加新的逆变器品牌支持。这种设计模式使得代码更加模块化、可维护，并且便于扩展。
