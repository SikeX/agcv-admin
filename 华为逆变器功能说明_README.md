# 华为逆变器AGC/AVC功能及日志优化说明

## 概览

本次更新实现了两个主要功能：
1. **日志大小限制**：限制单个日志文件最大为20MB
2. **华为逆变器AGC/AVC调控**：实现华为逆变器的5xx系列点位定义及完整的AGC/AVC调控流程

## 快速开始

### 1. 日志配置

在 `server/config.yaml` 中已配置：

```yaml
zap:
  level: info
  max-size: 20  # 单个日志文件最大20MB
  director: log
  retention-day: -1
```

### 2. 华为逆变器点位（5xx系列）

系统启动时会自动初始化以下点位到数据库：

#### AGC控制点位
- **501**: 无功功率变化梯度(%/s)
- **502**: 有功功率变化梯度(%/s)
- **503**: 调度指令维持时间(s)
- **504**: 有功功率固定值降额(kW)
- **507**: 有功功率百分比降额(0.1%)
- **508**: 有功功率固定值降额(W)

#### AVC控制点位
- **505**: 无功功率补偿(功率因数)
- **506**: 无功功率补偿(Q/S)

#### 控制模式点位
- **510-512**: 有功调节模式、调节值、调节指令
- **513-515**: 无功调节模式、调节值、调节指令

## 使用示例

### AGC控制

#### 百分比模式（推荐用于统一调度）
```go
// 控制逆变器按80%额定功率输出
err := HuaweiController.ControlAGCByPercentage(inverter, 80.0)
```

#### 绝对值模式（推荐用于精确控制）
```go
// 控制逆变器输出500kW
err := HuaweiController.ControlAGCByAbsolute(inverter, 500.0)
```

### AVC控制

#### 功率因数模式
```go
// 设置功率因数为0.95
err := HuaweiController.ControlAVCByPowerFactor(inverter, 0.95)
```

#### Q/S模式（推荐）
```go
// 设置无功/视在功率比值为0.3
err := HuaweiController.ControlAVCByReactivePower(inverter, 0.3)
```

### 设置调节梯度
```go
// 设置有功梯度10%/s，无功梯度5%/s
err := HuaweiController.SetControlGradient(inverter, 10.0, 5.0)
```

## 文件清单

### 新增文件
```
server/service/agvc/agcv_main/
├── huawei_inverter_controller.go      # 华为逆变器控制器（核心）
└── huawei_point_initializer.go        # 点位初始化器

server/
└── 华为逆变器AGC_AVC调控说明.md        # 详细技术文档

CHANGELOG_华为逆变器AGC_AVC_日志优化.md  # 完整更新日志
```

### 修改文件
```
server/
├── config.yaml                         # 添加日志大小配置
├── config/zap.go                       # 添加MaxSize字段
├── core/
│   ├── server.go                       # 集成华为点位初始化
│   └── internal/
│       ├── cutter.go                   # 实现文件大小切割
│       └── zap_core.go                 # 集成maxSize配置
└── service/agvc/agcv_main/
    ├── agc.go                          # 集成华为AGC控制器
    └── avc.go                          # 集成华为AVC控制器
```

## 功能特点

### 1. 日志管理
- ✅ 单文件大小限制20MB
- ✅ 自动切割并重命名（带时间戳）
- ✅ 保持按日期切割功能
- ✅ 支持配置保留天数

### 2. 华为逆变器控制
- ✅ 双模式AGC（百分比/绝对值）
- ✅ 双模式AVC（功率因数/Q/S）
- ✅ 调节梯度控制
- ✅ 自动数值范围检查
- ✅ 详细日志记录

### 3. 系统集成
- ✅ 自动初始化点位映射
- ✅ 品牌自动识别
- ✅ 向后兼容其他品牌
- ✅ 无缝集成现有AGC/AVC服务

## 调控流程

### AGC百分比模式流程
```
1. 设置调节模式 = 0 (百分比)
2. 写入百分比值到寄存器40125
3. 发送执行指令到寄存器35303
```

### AGC绝对值模式流程
```
1. 设置调节模式 = 1 (绝对值)
2. 写入功率值(kW)到寄存器40120
3. 发送执行指令到寄存器35303
```

### AVC Q/S模式流程
```
1. 获取视在功率S
2. 计算Q/S比值
3. 设置调节模式 = 2 (Q/S)
4. 写入Q/S值到寄存器40123
5. 发送执行指令到寄存器35307
```

## 数据精度

| 数据类型 | 精度 | 示例 |
|----------|------|------|
| 有功功率(kW) | 0.1kW | 500.5kW → 5005 |
| 有功百分比 | 0.1% | 85.5% → 855 |
| 功率因数 | 0.001 | 0.95 → 950 |
| Q/S比值 | 0.001 | 0.328 → 328 |
| 变化梯度 | 0.001%/s | 10.5%/s → 10500 |

## 注意事项

### 部署前
1. 确保数据库连接正常（用于点位初始化）
2. 检查逆变器品牌字段设置（华为 = 1）
3. 备份现有配置文件

### 运行时
1. 首次启动会初始化5xx点位到数据库
2. 日志切割后旧文件保留（根据retention-day清理）
3. 华为逆变器自动使用新控制器
4. 其他品牌继续使用原控制方式

### 调试
- 查看日志：`server/log/` 目录
- 点位映射：数据库 `agvc_point_mapping` 表
- 控制日志：搜索 "华为逆变器AGC/AVC"

## 技术支持

### 文档
- 📖 详细技术文档：`server/华为逆变器AGC_AVC调控说明.md`
- 📋 更新日志：`CHANGELOG_华为逆变器AGC_AVC_日志优化.md`
- 📊 Excel模板：`server/华为逆变器-AGVC模板数据.xlsx`

### 代码位置
- 控制器：`server/service/agvc/agcv_main/huawei_inverter_controller.go`
- 初始化：`server/service/agvc/agcv_main/huawei_point_initializer.go`
- AGC集成：`server/service/agvc/agcv_main/agc.go`
- AVC集成：`server/service/agvc/agcv_main/avc.go`

## 常见问题

### Q1: 如何验证日志切割功能？
A: 查看 `server/log/` 目录，当日志文件达到20MB时会自动生成带时间戳的备份文件。

### Q2: 如何确认华为点位初始化成功？
A: 查看启动日志，搜索 "华为逆变器5xx点位映射初始化完成"，或查询数据库 `agvc_point_mapping` 表中点位501-515。

### Q3: 如何选择AGC模式（百分比vs绝对值）？
A: 
- **百分比模式**：适合所有逆变器统一调度，负载均衡
- **绝对值模式**：适合不同容量逆变器，精确控制

### Q4: 如何选择AVC模式（PF vs Q/S）？
A:
- **PF模式**：适合功率因数考核场景
- **Q/S模式**：适合精确无功控制和电压调节

### Q5: 其他品牌逆变器是否受影响？
A: 不受影响。系统根据品牌字段自动选择控制方式，其他品牌继续使用原有控制逻辑。

## 版本信息

- **分支**: feat-huawei-inverter-agc-avc-5xx-ids-log-20m
- **日期**: 2024
- **主要功能**: 华为逆变器AGC/AVC + 日志大小限制20MB
- **兼容性**: 向后兼容，不影响现有功能

## 下一步

1. ✅ 基础功能已完成
2. 🔄 建议进行集成测试
3. 📈 监控运行状态
4. 🔧 根据实际情况优化参数

---

**注意**: 本功能已集成到AGC/AVC主流程中，启动后即可自动工作。
