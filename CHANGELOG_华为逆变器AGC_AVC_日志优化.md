# 更新日志 - 华为逆变器AGC/AVC调控及日志优化

## 版本信息
- 日期: 2024
- 分支: feat-huawei-inverter-agc-avc-5xx-ids-log-20m
- 类型: 功能增强 + 配置优化

## 主要更新

### 1. 日志大小限制优化（20MB）

#### 修改文件：
- `server/config/zap.go` - 添加MaxSize配置字段
- `server/config.yaml` - 添加max-size: 20配置项
- `server/core/internal/cutter.go` - 实现文件大小切割功能
- `server/core/internal/zap_core.go` - 集成maxSize配置

#### 功能说明：
- 单个日志文件最大限制为20MB
- 达到大小限制后自动切割，生成带时间戳的新文件
- 旧文件格式：`原文件名.20240101120000.log`
- 保持原有的按日期切割功能不变

#### 配置示例：
```yaml
zap:
  level: info
  director: log
  max-size: 20  # MB
  retention-day: -1
```

### 2. 华为逆变器5xx点位标识定义

#### 新增文件：
- `server/service/agvc/agcv_main/huawei_inverter_controller.go` - 华为逆变器控制器
- `server/service/agvc/agcv_main/huawei_point_initializer.go` - 点位初始化器
- `server/华为逆变器AGC_AVC调控说明.md` - 详细技术文档

#### 点位定义（5xx系列）：

| 点标识 | 名称 | 用途 | 寄存器地址 |
|--------|------|------|------------|
| 501 | 无功功率变化梯度 | AVC调节速率控制 | 42015 |
| 502 | 有功功率变化梯度 | AGC调节速率控制 | 42017 |
| 503 | 调度指令维持时间 | 控制指令持续时间 | 42019 |
| 504 | 有功功率固定值降额(kW) | AGC绝对值控制 | 40120 |
| 505 | 无功功率补偿(PF) | AVC功率因数控制 | 40122 |
| 506 | 无功功率补偿(Q/S) | AVC无功比值控制 | 40123 |
| 507 | 有功功率百分比降额 | AGC百分比控制 | 40125 |
| 508 | 有功功率固定值降额(W) | AGC精细控制 | 40126 |
| 510 | 有功调节模式 | 模式选择 | 35300 |
| 511 | 有功调节值 | 调节设定值 | 35301 |
| 512 | 有功调节指令 | 执行指令 | 35303 |
| 513 | 无功调节模式 | 模式选择 | 35304 |
| 514 | 无功调节值 | 调节设定值 | 35305 |
| 515 | 无功调节指令 | 执行指令 | 35307 |

### 3. AGC调控流程实现

#### 百分比调节模式
```go
// 控制逆变器按80%额定功率输出
err := HuaweiController.ControlAGCByPercentage(inverter, 80.0)
```

特点：
- 所有逆变器统一按比例调节
- 适合负载均衡分配
- 快速响应调度指令

#### 绝对值调节模式
```go
// 控制逆变器输出500kW
err := HuaweiController.ControlAGCByAbsolute(inverter, 500.0)
```

特点：
- 精确功率控制
- 支持不同容量逆变器差异化控制
- 适合基于实时负荷的闭环调节

#### 集成到现有系统
- 修改 `server/service/agvc/agcv_main/agc.go`
- 在 `executeRemoteOpenLoopControl` 和 `executeRemoteClosedLoopControl` 中集成
- 根据逆变器品牌自动选择控制方式

### 4. AVC调控流程实现

#### 功率因数模式
```go
// 设置功率因数为0.95
err := HuaweiController.ControlAVCByPowerFactor(inverter, 0.95)
```

特点：
- 功率因数调节
- 电网功率因数考核
- 标准电能质量控制

#### Q/S模式
```go
// 设置Q/S比值为0.3
err := HuaweiController.ControlAVCByReactivePower(inverter, 0.3)
```

特点：
- 精确无功调节
- 电压调节
- 无功补偿

#### 集成到现有系统
- 修改 `server/service/agvc/agcv_main/avc.go`
- 在 `executeRemoteReactiveControl` 中集成
- 自动计算Q/S比值并控制

### 5. 系统初始化增强

#### 修改文件：
- `server/core/server.go`

#### 新增初始化步骤：
```go
// 初始化华为逆变器5xx点位映射
if err := agvcMain.HuaweiPointInit.InitializeHuaweiPoints(); err != nil {
    zap.L().Error("初始化华为逆变器点位映射失败", zap.Error(err))
}
```

功能：
- 系统启动时自动初始化点位映射到数据库
- 支持增量更新，不会重复创建
- 自动更新点位描述信息

## 技术细节

### 日志切割机制

```go
// 检查文件大小
info, err := os.Stat(baseFilename)
if err != nil || info.Size() < c.maxSize {
    return baseFilename  // 未达到限制，继续使用
}

// 达到限制，切割文件
newFilename := nameWithoutExt + "." + time.Now().Format("20060102150405") + ext
os.Rename(baseFilename, newFilename)
```

### 华为逆变器控制器核心方法

| 方法 | 功能 | 参数 |
|------|------|------|
| ControlAGCByPercentage | AGC百分比控制 | inverter, targetPercentage |
| ControlAGCByAbsolute | AGC绝对值控制 | inverter, targetPowerKW |
| ControlAVCByPowerFactor | AVC功率因数控制 | inverter, targetPF |
| ControlAVCByReactivePower | AVC无功比值控制 | inverter, targetQS |
| SetControlGradient | 设置调节梯度 | inverter, activePowerGradient, reactivePowerGradient |
| GetPointDefinition | 获取点位定义 | pointID |
| ValidateControlValue | 验证控制值 | pointID, value |

### 数据精度处理

| 数据类型 | 精度 | 转换公式 | 示例 |
|----------|------|----------|------|
| 有功功率(kW) | 0.1kW | value * 10 | 500.5kW → 5005 |
| 有功百分比 | 0.1% | value * 10 | 85.5% → 855 |
| 功率因数 | 0.001 | value * 1000 | 0.95 → 950 |
| Q/S比值 | 0.001 | value * 1000 | 0.328 → 328 |
| 变化梯度 | 0.001%/s | value * 1000 | 10.5%/s → 10500 |

## 兼容性说明

### 向后兼容
- ✅ 现有AGC/AVC功能不受影响
- ✅ 其他品牌逆变器使用原有控制方式
- ✅ 日志系统保持向下兼容

### 品牌识别
系统根据逆变器品牌字段自动选择控制方式：
- 华为逆变器 (brand = 1)：使用新的5xx点位控制器
- 其他品牌：使用原有的点位映射控制

## 测试要点

### 日志功能测试
1. 验证日志文件大小限制为20MB
2. 验证文件切割功能正常
3. 验证切割后文件命名格式
4. 验证日志保留天数功能

### 华为逆变器控制测试
1. 测试AGC百分比模式控制
2. 测试AGC绝对值模式控制
3. 测试AVC功率因数模式控制
4. 测试AVC Q/S模式控制
5. 验证调节梯度设置
6. 验证点位初始化功能

### 集成测试
1. 验证系统启动时点位自动初始化
2. 验证AGC服务集成华为控制器
3. 验证AVC服务集成华为控制器
4. 验证多品牌逆变器混合场景

## 部署说明

### 配置更新
1. 更新 `config.yaml`，添加 `max-size: 20`
2. 重启服务使配置生效

### 数据库更新
系统启动时会自动创建或更新华为逆变器点位映射，无需手动操作。

### 注意事项
1. 首次启动会在数据库中创建5xx系列点位映射
2. 日志文件切割后，旧文件不会自动删除（受retention-day控制）
3. 华为逆变器需要正确设置品牌字段（brand = 1）

## 文档
- 详细技术文档：`server/华为逆变器AGC_AVC调控说明.md`
- Excel数据源：`server/华为逆变器-AGVC模板数据.xlsx`

## 后续优化建议

1. **性能优化**
   - 考虑批量发送控制指令，减少网络通信次数
   - 优化日志切割性能，避免高频写入时的锁竞争

2. **功能增强**
   - 添加控制指令执行状态监控
   - 实现控制指令重试机制
   - 增加控制历史记录功能

3. **监控告警**
   - 添加调节梯度超限告警
   - 实现控制失败告警通知
   - 监控日志文件大小和数量

4. **测试覆盖**
   - 增加单元测试覆盖率
   - 编写集成测试用例
   - 进行压力测试和边界测试

## 贡献者
- 系统架构：基于GVA框架
- 华为逆变器控制：新增功能
- 日志优化：配置增强
