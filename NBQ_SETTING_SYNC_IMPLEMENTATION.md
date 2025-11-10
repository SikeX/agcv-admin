# 逆变器配置与调度系统同步实现说明

## 概述

本实现确保逆变器配置（`AgvcNbqSetting`）与调度系统之间的双向同步，当配置发生变化时能够实时同步到调度系统，反之亦然。

## 功能特性

### 1. 本地到调度的同步

当通过Web界面或API更新/创建逆变器配置时，系统会自动将配置同步到调度系统。

**实现位置**：
- `server/service/agvc/agvc_nbq_setting.go`

**触发时机**：
- 创建逆变器配置时（`CreateAgvcNbqSetting`）
- 更新逆变器配置时（`UpdateAgvcNbqSetting`）

**同步的配置项**：
- 额定有功功率（`ratedActivePower`）
- 额定无功功率（`ratedReactivePower`）
- 参与调节标志（`isParticipateAdjust`）
- 标杆逆变器标志（`isBenchmarkInverter`）
- 升额优先级（`upgradePriority`）
- 降级优先级（`downgradePriority`）
- 功率抖动区间（`jitterRange`）
- 功率死区区间（`deadbandRange`）
- 逆变器编号（`inverterNo`）

**同步方式**：
通过CoAP协议发送到调度系统的1189端口，使用`CoapSender.SendAGCResultToDispatch`方法。

### 2. 调度到本地的同步

当调度系统需要更新逆变器配置时，可通过CoAP接口发送更新请求。

**实现位置**：
- `server/initialize/coap_handler.go`

**接口信息**：
- **路径**：`/agvc/nbq_setting`
- **方法**：POST
- **端口**：1188（CoAP服务器监听端口）

**请求数据格式**：
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

**响应格式**：
```json
{
  "success": true
}
```

**错误响应**：
- `400 Bad Request`: JSON格式错误或缺少必需字段
- `404 Not Found`: 指定的逆变器编号不存在
- `500 Internal Server Error`: 数据库更新失败

## 技术细节

### 同步流程

#### 本地到调度的同步流程：

```
用户操作（创建/更新）
    ↓
API层接收请求
    ↓
Service层处理业务逻辑
    ↓
更新数据库
    ↓
调用syncNbqSettingToDispatch
    ↓
构建同步数据
    ↓
通过CoapSender.SendAGCResultToDispatch发送
    ↓
调度系统1189端口接收
```

#### 调度到本地的同步流程：

```
调度系统发送配置更新
    ↓
CoAP服务器1188端口接收
    ↓
handleNbqSettingUpdate处理器
    ↓
解析JSON数据
    ↓
查询逆变器配置
    ↓
更新配置字段
    ↓
保存到数据库
    ↓
返回成功响应
```

### 关键方法

#### syncNbqSettingToDispatch

```go
func (agvcNbqSettingService *AgvcNbqSettingService) syncNbqSettingToDispatch(
    ctx context.Context, 
    setting agvc.AgvcNbqSetting
) error
```

**功能**：将逆变器配置同步到调度系统

**参数**：
- `ctx`：上下文对象
- `setting`：逆变器配置对象

**返回值**：
- `error`：同步过程中的错误，如果成功则为nil

**特点**：
- 自动检查必需字段（逆变器编号、并网点编号）
- 只同步非空字段
- 失败时记录警告日志，但不影响主流程

#### handleNbqSettingUpdate

```go
func handleNbqSettingUpdate(
    ctx context.Context, 
    msg coapMessage
) (code byte, payload []byte)
```

**功能**：处理来自调度系统的逆变器配置更新请求

**参数**：
- `ctx`：上下文对象
- `msg`：CoAP消息对象

**返回值**：
- `code`：CoAP响应码
- `payload`：响应数据

**特点**：
- 支持部分字段更新
- 自动验证数据有效性
- 记录详细的操作日志

## 使用场景

### 场景1：通过Web界面修改逆变器配置

1. 管理员在Web界面修改逆变器的额定功率
2. 前端调用更新API
3. 后端更新数据库
4. 系统自动将新配置同步到调度系统
5. 调度系统接收到更新并调整控制策略

### 场景2：调度系统远程调整逆变器参数

1. 调度系统检测到需要调整逆变器参与调节状态
2. 调度系统通过CoAP发送配置更新请求到1188端口
3. 系统接收请求并更新数据库
4. 新配置立即生效
5. AGC/AVC控制逻辑使用新配置进行调节

### 场景3：批量同步配置

1. 系统启动时或定时任务触发
2. 遍历所有逆变器配置
3. 逐个同步到调度系统
4. 确保调度系统与本地配置一致

## 错误处理

### 同步失败处理

当配置同步到调度系统失败时：
- 记录警告日志，包含错误详情
- **不影响**数据库更新操作（已提交）
- 后续可通过重试机制或手动触发重新同步

### 字段验证

- 逆变器编号：必需字段，不能为空
- 并网点编号：必需字段，用于确定目标调度点
- 其他字段：可选，只同步非空值

## 日志记录

系统在关键操作点记录日志：

1. **Info级别**：
   - 成功同步配置到调度
   - 从调度接收到配置更新

2. **Warn级别**：
   - 同步到调度失败
   - 缺少必需字段

3. **Debug级别**：
   - 配置无需同步（无变更）
   - 数据转换过程

4. **Error级别**：
   - 数据格式错误
   - 数据库操作失败

## 配置要求

确保以下配置正确：

1. **CoAP服务器配置**（`config.yaml`）：
```yaml
coap:
  enable: true
  host: "0.0.0.0"
  port: 1188
```

2. **调度系统地址**：
   - 在`CoapSender`中配置调度系统的主机和端口
   - 默认端口：1189（发送到调度）

## 监控建议

1. 监控同步失败的日志
2. 定期检查配置一致性
3. 监控CoAP服务器的可用性
4. 跟踪同步延迟和性能

## 未来扩展

1. **批量同步接口**：支持一次同步多个逆变器配置
2. **同步状态跟踪**：记录每次同步的时间和状态
3. **冲突解决机制**：当本地和调度配置不一致时的处理策略
4. **重试机制**：同步失败时的自动重试
5. **配置版本管理**：跟踪配置变更历史

## 注意事项

1. **网络依赖**：同步功能依赖于与调度系统的网络连接
2. **性能考虑**：频繁更新可能产生大量网络流量
3. **数据一致性**：确保同步数据的原子性和一致性
4. **安全性**：CoAP通信建议使用加密和认证机制
