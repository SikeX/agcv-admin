# 逆变器调节闭锁功能实现 - 变更总结

## 功能描述

实现了AGC系统中的逆变器上下调节闭锁功能：
- **上调节闭锁**：当所有逆变器都不能再往上调节时，自动设置上调节闭锁
- **下调节闭锁**：当所有逆变器都不能再往下调节时，自动设置下调节闭锁

## 修改的文件

### 1. server/service/agvc/agcv_main/agc.go

#### 新增方法
- `checkAndUpdateRegulationLocks(bwdNo int) error` - 检查并更新上下调节闭锁状态

#### 修改内容
- `executeAGCCycle()` - 在控制周期开始时调用闭锁检查
- `sendAGCResultToDispatch()` - 启用闭锁状态发送给调度系统

### 2. server/service/agvc/agcv_main/data_storage.go

#### 新增方法
- `SetData(psid, eqid, eqType, dataType int, point string, value interface{}) error` - 设置数据点值

### 3. server/service/agvc/agcv_main/dispatch_storage.go

#### 新增方法
- `SetData(psid, eqid, eqType, dataType int, point string, value interface{}) error` - 设置调度数据点值

## 核心逻辑

1. **检查时机**：每个AGC控制周期开始时自动检查
2. **判断标准**：
   - 上调能力：`(功率上限 - 当前功率) > 0.5kW`
   - 下调能力：`(当前功率 - 0) > 0.5kW`
3. **闭锁条件**：
   - 所有逆变器都不能上调 → 设置上调节闭锁
   - 所有逆变器都不能下调 → 设置下调节闭锁
4. **数据存储**：同时更新到DataStorage和DispatchStorage

## 使用的点位

- **点位405**：AGC有功上调节闭锁（遥信YX）
- **点位406**：AGC有功下调节闭锁（遥信YX）

值定义：
- `0` = 未闭锁，允许调节
- `1` = 已闭锁，禁止调节

## 依赖功能

- 逆变器在线状态管理
- 气象站辐射数据采集
- 逆变器当前功率数据采集
- DataStorage和DispatchStorage数据存储服务

## 测试验证

建议测试以下场景：
1. 高辐射场景（逆变器功率接近上限）- 验证上调节闭锁
2. 低辐射场景（逆变器功率接近0）- 验证下调节闭锁
3. 正常运行场景（功率在中间范围）- 验证闭锁解除

## 文档

详细实现文档请参考：`INVERTER_REGULATION_LOCKS_IMPLEMENTATION.md`
