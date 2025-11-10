# AGC流程修改总结

## 修改内容

本次修改主要针对AGC（Automatic Generation Control）流程的三个问题进行了修复和改进：

### 1. 修复华为逆变器百分比降额计算问题

**问题描述**：
华为逆变器在使用百分比模式（mode=0）进行AGC调控时，代码错误地将功率值（kW）直接传递给了期望百分比（0-100）的函数。

**修改位置**：
- `server/service/agvc/agcv_main/agc.go`
  - `executeRemoteOpenLoopControl` 函数（第464-492行）
  - `executeRemoteClosedLoopControl` 函数（第613-643行）

**修改内容**：
- 当华为逆变器使用百分比模式时，现在会将目标功率（kW）转换为百分比
- 转换公式：`targetPercentage = (targetPower / RatedActivePower) * 100.0`
- 添加了额定功率校验，确保不会除以0
- 增强了日志记录，便于问题追踪

**代码示例**：
```go
if mode == 0 {
    // 百分比模式：需要将目标功率转换为百分比（基于额定功率）
    var targetPercentage float64
    if inv.RatedActivePower != nil && *inv.RatedActivePower > 0 {
        targetPercentage = (targetPower / *inv.RatedActivePower) * 100.0
    } else {
        global.GVA_LOG.Warn("逆变器额定功率未配置，无法计算百分比")
        continue
    }
    if err := HuaweiController.ControlAGCByPercentage(&inv, targetPercentage); err != nil {
        // 错误处理
    }
}
```

---

### 2. 添加基于辐射的有功调节上限计算

**需求描述**：
有功调节上限值需要根据倾斜辐射和水平辐射动态计算，而不是使用固定配置值。

**修改位置**：
- `server/service/agvc/agcv_main/agc.go`
  - 新增函数 `calculateActivePowerLimitFromIrradiance`（第1111-1179行）
  - `executeRemoteClosedLoopControl` 函数中调用（第552-566行）
  - `executeClosedLoopControl` 函数中调用（第717-731行）

**实现逻辑**：
1. 从气象仪数据存储中获取水平辐射强度（点位8）和倾斜辐射强度（点位9）
2. 选择两者中的最大值作为实际辐射强度
3. 计算公式：`有功调节上限 = Σ(逆变器额定功率 × 实际辐射强度 / 1000)`
4. 标准测试条件（STC）下的辐射强度为 1000 W/㎡
5. 如果无法获取辐射数据，则回退使用配置中的执行值

**关键代码**：
```go
func (s *agc) calculateActivePowerLimitFromIrradiance(bwdNo int, inverters []agvc.AgvcNbqSetting) (float64, error) {
    // 获取水平辐射强度和倾斜辐射强度
    horizontalIrradiance, _ := DataStorage.GetDataAsFloat64(1, qxyNo, 8, cons.YC, "8")
    tiltedIrradiance, _ := DataStorage.GetDataAsFloat64(1, qxyNo, 8, cons.YC, "9")
    
    // 使用最大值
    actualIrradiance := max(horizontalIrradiance, tiltedIrradiance)
    
    // 计算总上限
    var totalActivePowerLimit float64 = 0
    for _, inv := range inverters {
        if inv.RatedActivePower != nil {
            invPowerLimit := *inv.RatedActivePower * (actualIrradiance / 1000.0)
            totalActivePowerLimit += invPowerLimit
        }
    }
    
    return totalActivePowerLimit, nil
}
```

**上限应用**：
- 在闭环控制中，计算每个逆变器的功率上限按额定功率比例分配
- 如果目标功率超过上限，则限制为上限值
- 添加了详细的调试日志，记录限制过程

---

### 3. 加入AGC有功上/下调节闭锁流程

**需求描述**：
AGC流程需要支持有功上调节闭锁和下调节闭锁功能，当闭锁时应跳过相应的调节。

**修改位置**：
- `server/service/agvc/agcv_main/agc.go`
  - `executeRemoteClosedLoopControl` 函数（第517-544行）
  - `executeClosedLoopControl` 函数（第682-709行）

**实现逻辑**：
1. 在执行调节前，根据调节方向（outputDeviation的正负）判断是上调还是下调
2. 如果是上调（outputDeviation > 0），检查上调节闭锁状态（点位405）
3. 如果是下调（outputDeviation < 0），检查下调节闭锁状态（点位406）
4. 闭锁值为1时表示已闭锁，跳过调节并记录日志
5. 闭锁值为0时或读取失败时，正常执行调节

**关键代码**：
```go
// 检查上/下调节闭锁
if outputDeviation > 0 {
    // 需要上调，检查上调节闭锁
    pointID, _ := PointMapper.GetPointID(cons.TYPE_AGC, cons.AGC_YX_UP_REG_LOCK)
    if pointID == "" {
        pointID = "405"
    }
    upRegLock, err := DataStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_AGC, cons.YX, pointID)
    if err == nil && upRegLock == 1 {
        global.GVA_LOG.Info("AGC有功上调节已闭锁，跳过上调",
            zap.Int("bwdNo", bwdNo),
            zap.Float64("outputDeviation", outputDeviation))
        return nil
    }
} else if outputDeviation < 0 {
    // 需要下调，检查下调节闭锁
    pointID, _ := PointMapper.GetPointID(cons.TYPE_AGC, cons.AGC_YX_DOWN_REG_LOCK)
    if pointID == "" {
        pointID = "406"
    }
    downRegLock, err := DataStorage.GetDataAsFloat64(1, bwdNo, cons.TYPE_AGC, cons.YX, pointID)
    if err == nil && downRegLock == 1 {
        global.GVA_LOG.Info("AGC有功下调节已闭锁，跳过下调",
            zap.Int("bwdNo", bwdNo),
            zap.Float64("outputDeviation", outputDeviation))
        return nil
    }
}
```

---

## 相关点位定义

### AGC点位（已在 `server/service/agvc/cons/def.go` 中定义）

**遥信(YX)**：
- `405`: AGC有功上调节闭锁 (0=未闭锁, 1=闭锁)
- `406`: AGC有功下调节闭锁 (0=未闭锁, 1=闭锁)

**遥测(YC)**：
- `401`: 有功调节上限 (kW)
- `402`: 有功调节下限 (kW)

### 气象仪点位（已在 `server/model/agvc/point_mapping.go` 中定义）

- `8`: 水平辐射强度 (W/㎡)
- `9`: 倾斜辐射强度 (W/㎡)

---

## 测试建议

### 1. 华为逆变器百分比模式测试
- 配置华为逆变器为百分比模式（mode=0）
- 观察日志，确认目标百分比计算正确
- 验证寄存器504接收到的值是否在0-1000范围内（对应0-100%）

### 2. 辐射上限计算测试
- 确保气象仪数据正常上报
- 观察日志中的辐射数据和计算出的功率上限
- 在不同辐射条件下验证上限是否合理变化
- 测试辐射数据缺失时的降级处理

### 3. 闭锁功能测试
- 设置上调节闭锁（点位405 = 1），测试是否能阻止上调
- 设置下调节闭锁（点位406 = 1），测试是否能阻止下调
- 取消闭锁（设置为0），验证调节恢复正常

### 4. 集成测试
- 在实际运行环境中观察AGC控制效果
- 验证功率调节是否符合辐射条件限制
- 检查闭锁状态下系统的稳定性

---

## 注意事项

1. **气象站编号配置**：当前代码中气象站编号硬编码为1，实际项目中可能需要根据并网点配置动态获取对应的气象站编号

2. **辐射数据类型**：当前代码假设气象仪的设备类型(eqType)为8，如实际不同需要调整

3. **额定功率配置**：所有逆变器必须配置正确的额定功率，否则无法进行百分比转换和辐射上限计算

4. **闭锁状态同步**：确保闭锁状态能够及时从调度端同步到系统，避免控制延迟

5. **日志级别**：部分调试日志使用Debug级别，生产环境可能需要调整日志级别以避免过多输出

---

## 文件清单

修改的文件：
- `server/service/agvc/agcv_main/agc.go` - 主要修改文件

相关参考文件（未修改）：
- `server/service/agvc/agcv_main/huawei_inverter_controller.go`
- `server/service/agvc/cons/def.go`
- `server/model/agvc/point_mapping.go`
- `server/model/agvc/agvc_main/agc_config.go`
- `server/model/agvc/agvc_bwd_setting.go`
- `server/model/agvc/agvc_nbq_setting.go`
