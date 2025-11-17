# 变更摘要 - CoAP客户端轮询功能

## 概述
将AGVC模块的数据获取方式从**被动接收**改为**主动轮询**，实现定期从CoAP服务器（端口5589）请求AGC、AVC和并网柜数据。

## 架构变更
```
旧架构：外部设备 → CoAP服务端(1188) → 被动接收数据
新架构：CoAP客户端 → 定期轮询(5589) → 主动请求数据 → 存储到内存+InfluxDB
```

## 修改的文件

### 1. `server/config/coap.go`
**修改内容**：添加CoAP客户端配置字段
```go
type Coap struct {
    // 原有字段...
    
    // 新增字段
    ClientEnable   bool   // 是否启用客户端轮询
    ClientHost     string // CoAP服务器地址
    ClientPort     int    // CoAP服务器端口
    PollInterval   int    // 轮询间隔（秒）
}
```

### 2. `server/config.yaml`
**修改内容**：添加CoAP客户端配置项
```yaml
coap:
  # 原有配置...
  
  # 新增配置
  client-enable: true
  client-host: 127.0.0.1
  client-port: 5589
  poll-interval: 5
```

### 3. `server/core/agcv.go`
**修改内容**：在 `agvcInitialize()` 函数中添加轮询服务初始化
```go
func agvcInitialize() {
    agvcMain.DataStorage.Initialize()
    agvcMain.DispatchStorage.Initialize()
    agvcMain.CoapClientPoller.Initialize() // 新增
}
```

## 新增的文件

### 1. `server/service/agvc/agcv_main/coap_client_poller.go`
**核心实现文件**（249行）

**主要功能**：
- CoAP客户端轮询服务实现
- 定期从所有并网点请求三种设备类型数据
- 批量存储数据到内存和InfluxDB

**关键结构**：
```go
type coapClientPoller struct {
    stopChan     chan struct{}
    pollTicker   *time.Ticker
    serverHost   string
    serverPort   int
    pollInterval time.Duration
}
```

**核心方法**：
- `Initialize()`: 初始化服务
- `pollAllGridPoints()`: 轮询所有并网点
- `requestData()`: 发送CoAP请求
- `storeData()`: 存储数据

### 2. 文档文件
- `server/service/agvc/agcv_main/COAP_CLIENT_POLLER_README.md` - 功能详细文档
- `COAP_CLIENT_MIGRATION.md` - 实现总结和架构说明
- `IMPLEMENTATION_CHECKLIST.md` - 实现检查清单
- `CHANGES_SUMMARY.md` - 本文件

## 功能特性

### 数据请求
- **接口**: POST `/agvc/data`
- **参数**: 
  - `psid`: 1（固定）
  - `eqid`: 并网点编号
  - `eqType`: 60/61/5（AGC/AVC/并网柜）
- **响应**: `CoAPDataMessage` 数组

### 轮询策略
- 从数据库读取所有并网点（`agvc_bwd_setting`表）
- 对每个并网点请求三种设备类型数据
- 默认5秒轮询一次（可配置）

### 数据存储
- **内存**: 使用 `DataStorage` 存储到Map
- **持久化**: 自动定期保存到InfluxDB（5秒一次）

### 错误处理
- 连接失败：记录日志，继续处理
- 解析失败：记录日志，跳过该条数据
- 超时：5秒超时自动取消

## 配置说明

### 启用功能
```yaml
coap:
  client-enable: true  # 必须设置为true
```

### 自定义配置
```yaml
coap:
  client-enable: true
  client-host: 192.168.1.100  # 自定义服务器地址
  client-port: 5589           # 自定义端口
  poll-interval: 3            # 自定义轮询间隔
```

### 禁用功能
```yaml
coap:
  client-enable: false  # 禁用轮询
```

## 兼容性

### ✅ 向后兼容
- 原有CoAP服务端功能保留
- 通过配置开关控制
- 不影响其他AGVC功能
- 数据格式完全兼容

### ✅ 数据结构复用
- 复用 `CoAPDataMessage`
- 复用 `RealtimeData`
- 复用 `DataStorage`

## 测试要点

### 功能测试
- [ ] 服务正常启动
- [ ] 定时轮询执行
- [ ] 数据正确请求
- [ ] 数据正确存储

### 配置测试
- [ ] 启用/禁用开关
- [ ] 自定义地址和端口
- [ ] 自定义轮询间隔

### 异常测试
- [ ] CoAP服务器不可达
- [ ] 数据格式错误
- [ ] 网络超时
- [ ] 并网点配置为空

## 日志示例

### 启动日志
```
[INFO] CoAP客户端轮询服务初始化 server=127.0.0.1:5589 interval=5s
```

### 轮询日志
```
[DEBUG] 开始轮询并网点数据 并网点数量=3
[DEBUG] 成功获取并存储数据 EQID=1001 EQType=60 数据量=15
```

### 错误日志
```
[ERROR] 请求数据失败 EQID=1002 EQType=61 error=connection timeout
```

## 部署清单

### 部署前检查
- [ ] 更新代码
- [ ] 修改配置文件
- [ ] 确认并网点配置
- [ ] 确认CoAP服务器运行

### 部署步骤
1. 拉取最新代码
2. 更新 `config.yaml`
3. 启动服务
4. 检查日志
5. 验证数据

### 验证步骤
1. 查看启动日志
2. 查看轮询日志
3. 检查内存数据
4. 检查InfluxDB数据

## 影响范围

### ✅ 无影响
- 原有CoAP服务端
- 其他AGVC功能
- 数据库结构
- API接口

### ✨ 新增功能
- CoAP客户端轮询
- 配置项
- 日志记录
- 文档

## 性能影响

### 资源占用
- **CPU**: 轮询和数据处理（低）
- **内存**: 数据存储（MB级别）
- **网络**: 定期请求（每个并网点3个请求/轮询周期）

### 优化建议
- 合理设置轮询间隔
- 根据并网点数量调整
- 监控资源使用情况

## 文档清单

1. ✅ **COAP_CLIENT_POLLER_README.md** - 功能详细文档
2. ✅ **COAP_CLIENT_MIGRATION.md** - 实现总结和架构说明
3. ✅ **IMPLEMENTATION_CHECKLIST.md** - 实现检查清单
4. ✅ **CHANGES_SUMMARY.md** - 本变更摘要

## Git提交信息建议

```
feat: 实现CoAP客户端轮询功能，支持主动从CoAP服务器获取AGVC数据

- 新增CoAP客户端轮询服务（coap_client_poller.go）
- 添加客户端配置项（client-enable, client-host, client-port, poll-interval）
- 支持定期轮询所有并网点的AGC、AVC和并网柜数据
- 数据自动存储到内存和InfluxDB
- 完善的错误处理和日志记录
- 向后兼容，不影响原有功能

涉及文件：
- 新增：server/service/agvc/agcv_main/coap_client_poller.go
- 修改：server/config/coap.go
- 修改：server/config.yaml
- 修改：server/core/agcv.go
- 文档：3个README文档

架构变更：从被动接收改为主动轮询
```

## 联系信息

如有问题，请查看详细文档：
- 功能文档：`server/service/agvc/agcv_main/COAP_CLIENT_POLLER_README.md`
- 实现总结：`COAP_CLIENT_MIGRATION.md`
- 检查清单：`IMPLEMENTATION_CHECKLIST.md`

---

**版本**: 1.0.0  
**日期**: 2024-11-17  
**状态**: ✅ 已完成
