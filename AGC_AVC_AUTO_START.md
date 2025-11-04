# AGC和AVC自动启动功能说明

## 功能概述

本次更新实现了系统启动时自动启动所有并网点的AGC和AVC功能。之前需要通过API接口手动启动每个并网点的AGC和AVC，现在系统会在启动时自动检查所有并网点配置并启动相应的功能。

## 主要变更

### 1. AGC服务 (`server/service/agvc/agcv_main/agc.go`)

新增方法：`AutoStartAllGridPoints()`

**功能说明：**
- 从数据库查询所有并网点配置（`agvc_bwd_setting`表）
- 遍历每个并网点配置
- 检查并网点的AGC是否启用（`agc_is_enabled`字段）
- 如果启用，自动调用`StartAGC(bwdNo)`启动AGC控制循环
- 记录详细的启动日志，包括成功/失败的并网点信息

**智能过滤：**
- 跳过`number`字段为空的配置
- 跳过`agc_is_enabled`为0或NULL的配置
- 自动处理并网点编号格式错误的情况

### 2. AVC服务 (`server/service/agvc/agcv_main/avc.go`)

新增方法：`AutoStartAllGridPoints()`

**功能说明：**
- 从数据库查询所有并网点配置（`agvc_bwd_setting`表）
- 遍历每个并网点配置
- 检查并网点的AVC是否启用（`avc_is_enabled`字段）
- 如果启用，自动调用`StartAVC(bwdNo)`启动AVC控制循环
- 记录详细的启动日志，包括成功/失败的并网点信息

**智能过滤：**
- 跳过`number`字段为空的配置
- 跳过`avc_is_enabled`为0或NULL的配置
- 自动处理并网点编号格式错误的情况

### 3. 系统启动 (`server/core/server.go`)

在`RunServer()`函数中添加自动启动逻辑：

```go
// 自动启动所有并网点的AGC和AVC功能
go func() {
    // 延迟3秒启动，确保所有依赖服务已完全初始化
    time.Sleep(3 * time.Second)
    
    global.GVA_LOG.Info("========== 开始自动启动所有并网点的AGC和AVC功能 ==========")
    
    // 自动启动所有并网点的AGC功能
    agvcMain.AGC.AutoStartAllGridPoints()
    
    // 自动启动所有并网点的AVC功能
    agvcMain.AVC.AutoStartAllGridPoints()
    
    global.GVA_LOG.Info("========== 自动启动流程完成 ==========")
}()
```

**设计说明：**
- 使用goroutine异步执行，不阻塞系统启动
- 延迟3秒启动，确保数据库连接、CoAP服务等依赖已完全初始化
- 在日志中使用明显的分隔符，便于查看启动状态

## 使用说明

### 1. 并网点配置要求

在`agvc_bwd_setting`表中，需要确保以下字段配置正确：

| 字段 | 说明 | 是否必需 |
|------|------|----------|
| `number` | 并网点编号（整数字符串） | 是 |
| `name` | 并网点名称 | 否 |
| `agc_is_enabled` | AGC启用标志（1=启用, 0=禁用） | 是 |
| `avc_is_enabled` | AVC启用标志（1=启用, 0=禁用） | 是 |
| `agc_control_period` | AGC控制周期（秒） | 是 |
| `avc_control_period` | AVC控制周期（秒） | 是 |

### 2. 启动日志示例

系统启动时，会在日志中看到以下信息：

```
[INFO] AGC控制服务初始化成功
[INFO] AVC控制服务初始化成功
[INFO] ========== 开始自动启动所有并网点的AGC和AVC功能 ==========
[INFO] 开始自动启动所有并网点的AGC功能  {"并网点数量": 5}
[INFO] 自动启动AGC成功  {"bwdNo": 1, "name": "并网点1"}
[DEBUG] 并网点AGC未启用，跳过  {"bwdNo": 2, "name": "并网点2"}
[INFO] 自动启动AGC成功  {"bwdNo": 3, "name": "并网点3"}
[WARN] 自动启动AGC失败  {"bwdNo": 4, "name": "并网点4", "error": "获取并网点配置失败"}
[INFO] AGC自动启动完成  {"成功数量": 2, "总数量": 5}
[INFO] 开始自动启动所有并网点的AVC功能  {"并网点数量": 5}
[INFO] 自动启动AVC成功  {"bwdNo": 1, "name": "并网点1"}
[INFO] 自动启动AVC成功  {"bwdNo": 3, "name": "并网点3"}
[INFO] AVC自动启动完成  {"成功数量": 2, "总数量": 5}
[INFO] ========== 自动启动流程完成 ==========
```

### 3. 启动条件

**AGC自动启动条件：**
1. 并网点配置存在于`agvc_bwd_setting`表
2. `number`字段不为空且为有效整数
3. `agc_is_enabled = 1`
4. 并网点配置在`agvc_bwd_setting`表中存在对应记录

**AVC自动启动条件：**
1. 并网点配置存在于`agvc_bwd_setting`表
2. `number`字段不为空且为有效整数
3. `avc_is_enabled = 1`
4. AVC配置在`agvc_avc_config`表中存在对应记录（可选，启动时会自动创建）

### 4. 错误处理

系统会智能处理以下错误情况：

- **并网点编号格式错误**：记录警告日志并跳过该并网点
- **配置不存在**：记录警告日志并跳过该并网点
- **AGC/AVC已在运行**：记录警告日志，避免重复启动
- **数据库查询失败**：记录错误日志，终止自动启动流程

## 与原有接口的关系

### 原有API接口保持不变

- `POST /agvc/agc/start?bwdNo=xxx` - 手动启动指定并网点的AGC
- `POST /agvc/agc/stop?bwdNo=xxx` - 停止指定并网点的AGC
- `POST /agvc/avc/start?bwdNo=xxx` - 手动启动指定并网点的AVC
- `POST /agvc/avc/stop?bwdNo=xxx` - 停止指定并网点的AVC

### 使用场景

- **自动启动**：系统启动时自动启动所有启用的并网点
- **手动控制**：运行时可通过API接口手动启动/停止特定并网点
- **动态调整**：修改`agc_is_enabled`或`avc_is_enabled`后，需要重启服务或手动调用API生效

## 配置示例

### 启用AGC和AVC的并网点配置

```sql
INSERT INTO agvc_bwd_setting (
    name, 
    number, 
    agc_is_enabled, 
    avc_is_enabled, 
    agc_control_period, 
    avc_control_period,
    agc_vibration_range,
    avc_adjustment_range_min,
    avc_adjustment_range_max
) VALUES (
    '并网点1',
    '1',
    1,  -- AGC启用
    1,  -- AVC启用
    30, -- AGC控制周期30秒
    60, -- AVC控制周期60秒
    50, -- AGC抖动区间50kW
    9.5,  -- AVC调节范围最小值9.5kV
    10.5  -- AVC调节范围最大值10.5kV
);
```

### 仅启用AGC的并网点配置

```sql
INSERT INTO agvc_bwd_setting (
    name, 
    number, 
    agc_is_enabled, 
    avc_is_enabled, 
    agc_control_period, 
    agc_vibration_range
) VALUES (
    '并网点2',
    '2',
    1,  -- AGC启用
    0,  -- AVC禁用
    30, -- AGC控制周期30秒
    50  -- AGC抖动区间50kW
);
```

## 注意事项

1. **数据库依赖**：自动启动功能依赖数据库连接，确保数据库配置正确
2. **启动延迟**：系统有3秒的启动延迟，确保所有依赖服务初始化完成
3. **并发安全**：AGC和AVC服务已实现并发安全，不会重复启动同一并网点
4. **日志监控**：建议监控启动日志，及时发现配置问题
5. **配置验证**：确保并网点配置完整且有效，避免启动失败

## 技术细节

### 启动流程

```
系统启动
  ↓
初始化数据库
  ↓
初始化AGC服务
  ↓
初始化AVC服务
  ↓
启动goroutine（异步）
  ↓
延迟3秒
  ↓
查询所有并网点配置
  ↓
遍历并网点配置
  ├─ 检查AGC是否启用 → 启动AGC
  └─ 检查AVC是否启用 → 启动AVC
  ↓
记录启动结果
  ↓
启动完成
```

### 并发控制

- AGC服务维护`bwdNoChans map[int]chan struct{}`记录已启动的并网点
- AVC服务维护`psidChans map[int]chan struct{}`记录已启动的并网点
- 启动前检查是否已在运行，避免重复启动
- 停止时清理对应的channel，释放资源

### 错误恢复

- 单个并网点启动失败不影响其他并网点
- 记录详细的错误日志，便于排查问题
- 系统启动失败不影响主服务运行

## 后续优化建议

1. **健康检查**：添加定时健康检查，自动重启失败的并网点
2. **配置热加载**：支持运行时动态加载配置，无需重启服务
3. **启动顺序优化**：根据并网点优先级控制启动顺序
4. **启动状态API**：提供API查询各并网点的启动状态
5. **启动失败告警**：集成告警系统，启动失败时发送通知

## 版本信息

- **版本**：v1.0.0
- **更新时间**：2024-11-04
- **更新内容**：实现AGC和AVC自动启动功能
