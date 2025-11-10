# 调度值与逆变器配置同步功能实现总结

## 任务说明

实现调度值与 `agvc.AgvcNbqSetting`（逆变器配置）的双向同步功能：
1. 当逆变器配置改变时，自动同步更新到调度系统
2. 当调度系统通过接口改变逆变器配置时，通过 `CoapSender.SendAGCResultToDispatch` 发送到调度

## 实现概述

本实现通过以下方式确保调度值与逆变器配置的同步：

### 1. 本地到调度的同步（Web界面/API → 调度系统）

**修改文件**：`server/service/agvc/agvc_nbq_setting.go`

#### 1.1 添加必要的导入

```go
import (
    "context"
    "fmt"

    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/agvc"
    agvcReq "github.com/flipped-aurora/gin-vue-admin/server/model/agvc/request"
    "github.com/flipped-aurora/gin-vue-admin/server/service/agvc/agcv_main"
    "go.uber.org/zap"
)
```

#### 1.2 修改创建方法

在 `CreateAgvcNbqSetting` 方法中添加同步逻辑：
- 创建成功后自动调用 `syncNbqSettingToDispatch` 同步到调度
- 同步失败只记录警告日志，不影响创建操作

#### 1.3 修改更新方法

在 `UpdateAgvcNbqSetting` 方法中添加同步逻辑：
- 更新成功后自动调用 `syncNbqSettingToDispatch` 同步到调度
- 同步失败只记录警告日志，不影响更新操作

#### 1.4 新增同步方法

实现 `syncNbqSettingToDispatch` 方法：
```go
func (agvcNbqSettingService *AgvcNbqSettingService) syncNbqSettingToDispatch(
    ctx context.Context, 
    setting agvc.AgvcNbqSetting
) error
```

**功能**：
- 验证必要字段（逆变器编号、并网点编号）
- 构建同步数据包（包含所有配置项）
- 通过 CoAP 发送到调度系统的 1189 端口
- 记录详细的操作日志

**同步的配置项**：
- `ratedActivePower` - 额定有功功率
- `ratedReactivePower` - 额定无功功率
- `isParticipateAdjust` - 参与调节标志
- `isBenchmarkInverter` - 标杆逆变器标志
- `upgradePriority` - 升额优先级
- `downgradePriority` - 降级优先级
- `jitterRange` - 功率抖动区间
- `deadbandRange` - 功率死区区间
- `inverterNo` - 逆变器编号

### 2. 调度到本地的同步（调度系统 → Web界面/API）

**修改文件**：`server/initialize/coap_handler.go`

#### 2.1 注册新的CoAP路由

在 `coapRoutes` 中添加逆变器配置更新端点：
```go
"/agvc/nbq_setting": {
    coapCodePost: handleNbqSettingUpdate,
}
```

#### 2.2 实现处理器函数

实现 `handleNbqSettingUpdate` 函数：
```go
func handleNbqSettingUpdate(
    ctx context.Context, 
    msg coapMessage
) (code byte, payload []byte)
```

**功能**：
- 接收来自调度系统的配置更新请求
- 解析 JSON 格式的请求数据
- 验证逆变器编号并查询现有配置
- 更新指定的配置字段
- 保存到数据库
- 返回操作结果

**接口信息**：
- **端口**：1188（CoAP服务器监听端口）
- **路径**：`/agvc/nbq_setting`
- **方法**：POST
- **数据格式**：JSON

**请求示例**：
```json
{
  "inverterNo": 1,
  "ratedActivePower": 1000.0,
  "ratedReactivePower": 800.0,
  "jitterRange": 10.0,
  "deadbandRange": 5.0,
  "upgradePriority": 1.0,
  "downgradePriority": 2.0,
  "isParticipateAdjust": true,
  "isBenchmarkInverter": false
}
```

**响应示例**：
```json
{
  "success": true
}
```

**错误码**：
- `400` - 请求数据格式错误或缺少必需字段
- `404` - 指定的逆变器不存在
- `500` - 数据库更新失败

## 数据流向

### 流向1：Web界面 → 数据库 → 调度系统

```
用户在Web界面修改逆变器配置
    ↓
前端调用API（POST /agvcNbqSetting/updateAgvcNbqSetting）
    ↓
API层（agvc_nbq_setting.go）
    ↓
Service层更新数据库
    ↓
调用syncNbqSettingToDispatch
    ↓
通过CoapSender.SendAGCResultToDispatch发送
    ↓
CoAP客户端发送到调度系统（1189端口）
    ↓
调度系统接收并应用配置
```

### 流向2：调度系统 → 数据库 → Web界面

```
调度系统需要修改逆变器配置
    ↓
调度系统发送CoAP请求（POST /agvc/nbq_setting）
    ↓
CoAP服务器接收（1188端口）
    ↓
handleNbqSettingUpdate处理器
    ↓
解析JSON数据并验证
    ↓
查询逆变器配置
    ↓
更新配置字段
    ↓
保存到数据库
    ↓
返回成功响应
    ↓
Web界面刷新后显示新配置
```

## 关键特性

### 1. 双向同步

- **正向同步**：本地修改自动同步到调度
- **反向同步**：调度修改自动更新到本地数据库

### 2. 容错机制

- 同步失败不影响主要操作（创建/更新）
- 详细的错误日志记录
- 优雅的错误处理

### 3. 灵活性

- 支持部分字段更新
- 只同步非空字段
- 自动跳过无效数据

### 4. 可追踪性

- 完整的操作日志
- 记录同步时间和数据量
- 记录失败原因

## 测试场景

### 场景1：通过Web界面修改配置

1. 登录Web管理界面
2. 进入逆变器配置页面
3. 修改某个逆变器的额定功率
4. 点击保存
5. 验证：
   - 数据库中的配置已更新
   - 调度系统收到同步数据
   - 日志中有同步成功的记录

### 场景2：调度系统远程更新配置

1. 调度系统发送CoAP请求到1188端口
2. 携带逆变器配置更新数据
3. 验证：
   - 数据库中的配置已更新
   - Web界面刷新后显示新配置
   - 日志中有更新成功的记录

### 场景3：网络异常处理

1. 断开与调度系统的网络连接
2. 通过Web界面修改逆变器配置
3. 验证：
   - 数据库更新成功
   - 日志中有同步失败的警告
   - 操作返回成功（不受同步失败影响）

## 配置要求

### CoAP服务器配置

在 `config.yaml` 中确保CoAP服务器已启用：

```yaml
coap:
  enable: true
  host: "0.0.0.0"
  port: 1188
```

### 调度系统配置

确保 `CoapSender` 中的调度系统地址配置正确：
- 主机地址：通过 `GetDispatchCoapHost()` 获取
- 端口：1189（通过 `GetDispatchCoapPort()` 获取）

## 日志监控

系统在以下情况记录日志：

### Info级别
- 逆变器配置已同步到调度
- 从调度系统接收到配置更新

### Warn级别
- 同步到调度失败
- 缺少必需字段（并网点编号）

### Debug级别
- 配置无需同步（无数据变更）

### Error级别
- JSON解析失败
- 数据库操作失败
- 逆变器不存在

## 性能考虑

1. **异步处理**：同步操作不阻塞主要业务流程
2. **批量优化**：可扩展为批量同步多个配置
3. **缓存机制**：可考虑缓存逆变器配置减少数据库查询
4. **重试机制**：可实现失败重试机制

## 安全考虑

1. **数据验证**：严格验证输入数据的有效性
2. **权限控制**：确保只有授权的调度系统可以更新配置
3. **审计日志**：记录所有配置变更操作
4. **加密传输**：考虑使用CoAPS（CoAP over DTLS）

## 文档

详细的实现文档请参考：
- `NBQ_SETTING_SYNC_IMPLEMENTATION.md` - 详细的功能说明和使用指南

## 总结

本实现完成了调度值与逆变器配置的双向同步功能，确保了数据的一致性和实时性。主要特点包括：

1. ✅ 自动同步：配置更新后自动同步到调度系统
2. ✅ 双向通信：支持调度系统远程更新配置
3. ✅ 容错性强：同步失败不影响主要操作
4. ✅ 日志完整：详细记录所有操作和错误
5. ✅ 易于维护：代码结构清晰，符合GVA规范

通过这个实现，系统能够实时保持本地配置与调度系统的一致性，提高了系统的可靠性和可维护性。
