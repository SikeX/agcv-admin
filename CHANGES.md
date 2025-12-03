# InfluxDB Bucket 保留策略配置功能

## 功能概述

为系统配置页面的InfluxDB配置部分添加了保留时间设置功能，支持用户通过前端界面配置bucket的数据保留策略，并在配置更新时自动应用到InfluxDB。

## 实现的功能

1. ✅ **agvc储存桶和逆变器储存桶保留时间配置**
   - 添加了两个输入框用于设置保留时间
   - 默认值显示在输入框中（而不是placeholder）
   - agvc储存桶默认保留时间：24h（1天）
   - 逆变器储存桶默认保留时间：168h（7天）

2. ✅ **配置自动更新**
   - 点击"立即更新"按钮后自动保存配置到config.yaml
   - 重新加载配置到全局变量
   - 自动调用InfluxDB API更新bucket的保留策略

3. ✅ **Bucket清除策略自动设置**
   - 使用InfluxDB BucketsAPI更新RetentionRules
   - 智能检查是否需要更新（避免不必要的API调用）
   - 支持的时间格式：24h, 7d, 168h 等（符合Go time.ParseDuration）

## 修改的文件

### 后端文件

1. **server/config/influxdb.go**
   - 添加了 `GetAgvcRetention()` 方法：返回agvc保留时间，默认24h
   - 添加了 `GetNbqRetention()` 方法：返回nbq保留时间，默认168h

2. **server/initialize/influxdb.go**
   - 添加导入：`time`, `domain` 包
   - 实现 `UpdateBucketRetentionPolicies()` 函数：更新所有bucket的保留策略
   - 实现 `updateBucketRetention()` 函数：更新单个bucket的保留策略
   - 在 `InfluxDB()` 初始化函数中调用策略更新

3. **server/service/system/sys_system.go**
   - 添加导入：`initialize` 包
   - 修改 `SetSystemConfig()` 方法：
     - 保存配置后重新加载到全局变量
     - 自动调用 `UpdateBucketRetentionPolicies()` 更新bucket策略
     - 即使策略更新失败也不影响配置保存（记录警告日志）

4. **server/config.yaml**
   - 更新了influxdb配置格式
   - 添加了 `agvcBucket`, `agvcRetention`, `nbqBucket`, `nbqRetention` 字段
   - 移除了旧的 `bucket`, `measurement`, `nbqMeasurement`, `retention` 字段

### 前端文件

1. **web/src/view/systemTools/system/system.vue**
   - 添加了"agvc储存桶保留时间"输入框
   - 添加了"逆变器储存桶保留时间"输入框
   - 在 `initForm()` 函数中设置默认值：
     - 如果配置中没有保留时间，设置 agvcRetention = '24h'
     - 如果配置中没有保留时间，设置 nbqRetention = '168h'
   - 添加了提示信息（1天=24h，7天=168h）

## 技术实现细节

### 保留策略更新流程

```
用户修改配置 → 点击"立即更新" 
    ↓
SetSystemConfig() 保存配置到config.yaml
    ↓
重新加载配置到 global.GVA_CONFIG
    ↓
调用 UpdateBucketRetentionPolicies()
    ↓
├─ updateBucketRetention(agvc_data, 24h)
│   ├─ 解析时间字符串
│   ├─ 查找bucket
│   ├─ 检查是否需要更新
│   └─ 调用 BucketsAPI.UpdateBucket()
│
└─ updateBucketRetention(nbq_data, 168h)
    ├─ 解析时间字符串
    ├─ 查找bucket
    ├─ 检查是否需要更新
    └─ 调用 BucketsAPI.UpdateBucket()
```

### 时间格式支持

支持的时间格式（符合Go `time.ParseDuration`）：
- 小时：`24h`, `48h`, `168h`
- 天：`1d` (需要换算为小时，如 24h)
- 混合：`24h30m`, `1h30m`

### 错误处理

1. **配置解析错误**：记录日志，返回错误
2. **Bucket查找失败**：记录日志，返回错误
3. **策略更新失败**：记录警告日志，不影响配置保存

### 默认值策略

1. **后端默认值**：在 `GetAgvcRetention()` 和 `GetNbqRetention()` 方法中定义
2. **前端默认值**：在 `initForm()` 函数中设置，确保用户看到实际的默认值
3. **一致性**：前后端默认值保持一致（24h和168h）

## 使用说明

### 配置保留时间

1. 登录系统，进入"系统工具" → "系统配置"
2. 切换到"InfluxDB 配置"标签页
3. 找到保留时间配置：
   - "agvc储存桶保留时间"（默认24h）
   - "逆变器储存桶保留时间"（默认168h）
4. 输入新的保留时间（格式：24h, 7d, 168h等）
5. 点击"立即更新"按钮
6. 系统会自动保存配置并更新InfluxDB的bucket策略

### 时间格式示例

- `24h` = 1天
- `168h` = 7天
- `720h` = 30天
- `8760h` = 365天

### 注意事项

1. 修改保留时间后需要点击"立即更新"才会生效
2. 保留时间更新会立即应用到InfluxDB
3. 旧数据不会立即删除，会在保留期过后自动清理
4. 如果InfluxDB连接失败，配置仍会保存，但策略不会更新

## 兼容性说明

- 向后兼容：保留了原有配置的读取方式
- 如果配置文件中没有保留时间字段，会使用默认值
- 前端初始化时会自动设置默认值，确保用户体验
