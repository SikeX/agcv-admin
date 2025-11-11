# 逆变器品牌管理模块

## 概述

本目录包含了所有与逆变器品牌相关的代码，采用工厂模式和策略模式实现了品牌的统一管理和扩展。

## 目录结构

```
inverter_brands/
├── README.md                        # 本文档
├── inverter_brand_interface.go      # 品牌接口定义
├── inverter_brand_factory.go        # 品牌工厂（注册和获取品牌实现）
├── inverter_brand_mapper.go         # 品牌点位映射器
├── coap_adapter.go                  # CoAP适配器接口
├── huawei_brand_impl.go             # 华为品牌实现
├── huawei_inverter_controller.go    # 华为逆变器控制器
└── huawei_point_initializer.go      # 华为点位初始化器
```

## 核心组件

### 1. InverterBrandInterface (品牌接口)

定义了所有逆变器品牌必须实现的方法：

- `GetBrandCode()` - 获取品牌代码
- `GetBrandName()` - 获取品牌名称
- `ControlAGCByPercentage()` - AGC百分比调控
- `ControlAGCByAbsolute()` - AGC绝对值调控
- `ControlAVCByPowerFactor()` - AVC功率因数调控
- `ControlAVCByReactivePower()` - AVC无功功率调控
- `InitializePoints()` - 初始化品牌点位

### 2. InverterBrandFactory (品牌工厂)

负责品牌的注册、获取和管理：

```go
// 注册品牌
BrandFactory.RegisterBrand(&HuaweiInverterBrand{})

// 获取品牌实现
brand, err := BrandFactory.GetBrand(cons.INVERTER_BRAND_HUAWEI)

// 初始化所有品牌点位
BrandFactory.InitializeAllBrands()
```

### 3. InverterBrandMapper (品牌映射器)

管理不同品牌的点位映射关系：

```go
// 获取品牌的无功功率控制点位
pointID, err := InverterBrandMapper.GetReactivePowerPoint(brandCode)

// 获取品牌的有功功率控制点位
pointID, err := InverterBrandMapper.GetActivePowerPoint(brandCode)
```

### 4. CoapAdapter (CoAP适配器)

提供CoAP通信接口，解耦品牌模块与主模块的依赖：

```go
// 发送逆变器控制指令
GetCoapAdapter().SendInverterCommand(host, port, psid, inverterNo, commands)
```

## 华为品牌实现

### 支持的调节模式

华为逆变器支持以下无功调节模式：

1. **功率因数模式 (40122)** - `HUAWEI_WG_MODE_PF`
   - 通过设置目标功率因数进行调控
   - 功率因数范围：-1 到 1

2. **Q/S模式 (40123)** - `HUAWEI_WG_MODE_QS`
   - 通过设置无功与视在功率的比值进行调控
   - Q/S比值范围：-1 到 1

3. **夜间无功功率模式 (40129)** - `HUAWEI_WG_MODE_NIGHT`
   - 在夜间（无有功输出）时直接设置无功功率
   - 无功功率范围：受逆变器额定容量限制

4. **夜间无功Q/S模式 (42809)** - `HUAWEI_WG_MODE_NIGHT_QS`
   - 在夜间使用Q/S比值进行调控
   - 使用额定视在功率作为基准

### 控制流程

```
调用方 → BrandFactory.GetBrand()
       → HuaweiInverterBrand
       → HuaweiController
       → GetCoapAdapter()
       → CoAP指令发送
       → 逆变器执行
```

## 添加新品牌

### 步骤1：定义品牌常量

在 `server/service/agvc/cons/def.go` 中添加：

```go
const (
    INVERTER_BRAND_HUAWEI  = 1 // 华为
    INVERTER_BRAND_SUNGROW = 2 // 阳光
    INVERTER_BRAND_NEW     = 4 // 新品牌
)
```

### 步骤2：创建品牌实现

创建 `new_brand_impl.go`：

```go
package inverter_brands

import (
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
    "github.com/flipped-aurora/gin-vue-admin/server/service/agvc/cons"
)

type NewInverterBrand struct{}

func (n *NewInverterBrand) GetBrandCode() int {
    return cons.INVERTER_BRAND_NEW
}

func (n *NewInverterBrand) GetBrandName() string {
    return "新品牌"
}

// 实现其他接口方法...
```

### 步骤3：注册品牌

在 `inverter_brand_factory.go` 的 `init()` 函数中添加：

```go
func init() {
    BrandFactory.RegisterBrand(&HuaweiInverterBrand{})
    BrandFactory.RegisterBrand(&NewInverterBrand{}) // 新增
}
```

### 步骤4：添加点位映射

在 `inverter_brand_mapper.go` 的 `Initialize()` 函数中添加：

```go
m.mappings[cons.INVERTER_BRAND_NEW] = InverterBrandPointMapping{
    Brand:              cons.INVERTER_BRAND_NEW,
    BrandName:          "新品牌",
    PowerSwitchYK:      "401",
    ActivePowerLimitYT: "401",
    ReactivePowerYT:    "402",
}
```

## 使用示例

### AGC控制示例

```go
// 获取品牌实现
brand, err := BrandFactory.GetBrand(inverter.InverterBrand)
if err != nil {
    return err
}

// 按百分比调控
err = brand.ControlAGCByPercentage(inverter, 80.0)

// 按绝对值调控
err = brand.ControlAGCByAbsolute(inverter, 500.0) // 500kW
```

### AVC控制示例

```go
// 获取品牌实现
brand, err := BrandFactory.GetBrand(inverter.InverterBrand)
if err != nil {
    return err
}

// 功率因数调控
err = brand.ControlAVCByPowerFactor(inverter, 0.95)

// Q/S调控
err = brand.ControlAVCByReactivePower(inverter, 0.3)
```

## AVC无功调控策略抽象

为了避免过多的if判断，我们将不同调节模式的处理逻辑抽象成了独立函数：

### 1. executeReactiveControlByPF
功率因数模式调控，根据目标无功和当前有功计算功率因数。

### 2. executeReactiveControlByQS
Q/S模式调控，根据目标无功和视在功率计算Q/S比值。

### 3. executeReactiveControlByNightReactive
夜间无功功率模式调控，直接设置无功值。

### 4. executeReactiveControlByNightQS
夜间Q/S模式调控，使用额定视在功率计算Q/S比值。

### 使用示例

在 `executeRemoteReactiveControl` 函数中：

```go
switch int(currentAdjustMode) {
case cons.HUAWEI_WG_MODE_PF:
    err = s.executeReactiveControlByPF(brandImpl, &dev, targetReactive, ratedActivePower)
case cons.HUAWEI_WG_MODE_QS:
    err = s.executeReactiveControlByQS(brandImpl, &dev, targetReactive)
case cons.HUAWEI_WG_MODE_NIGHT:
    err = s.executeReactiveControlByNightReactive(brandImpl, &dev, targetReactive)
case cons.HUAWEI_WG_MODE_NIGHT_QS:
    err = s.executeReactiveControlByNightQS(brandImpl, &dev, targetReactive)
}
```

## 注意事项

1. **线程安全**：BrandFactory 使用读写锁保证并发安全
2. **错误处理**：所有方法都返回error，调用方需要妥善处理
3. **向后兼容**：通过 `agcv_main/inverter_brand_exports.go` 保持与原有代码的兼容
4. **依赖注入**：通过 CoapAdapter 解耦模块间依赖，提高可测试性

## 维护指南

- 添加新品牌时，确保实现所有接口方法
- 修改接口定义时，需要同步更新所有品牌实现
- 新增调节模式时，在对应的控制器中添加处理逻辑
- 定期检查点位映射的准确性，确保与硬件协议一致

## 相关文档

- [品牌接口设计文档](./BRAND_INTERFACE_DESIGN.md) (待补充)
- [华为逆变器通信协议](./HUAWEI_PROTOCOL.md) (待补充)
- [点位映射说明](./POINT_MAPPING.md) (待补充)
