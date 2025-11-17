# CoAP客户端轮询功能实现总结

## 需求概述

将AGVC模块的数据获取方式从**被动接收**改为**主动轮询**：
- 从启动CoAP服务端接收数据 → 改为启动CoAP客户端主动请求数据
- CoAP服务器地址：`127.0.0.1:5589`
- 请求接口：`/agvc/data`
- 请求三种设备类型数据：AGC(60)、AVC(61)、并网柜(5)
- 数据存储到内存和InfluxDB

## 架构变更

### 原架构
```
外部设备 --推送数据--> CoAP服务端(1188) --接收--> 数据存储
```

### 新架构
```
CoAP客户端 --轮询请求--> CoAP服务器(5589) --返回数据--> 批量存储
     ↑                                              ↓
     |                                          内存存储
     |                                              ↓
     └----定时轮询(5秒)                       InfluxDB持久化
```

## 实现的文件

### 1. 新增文件

#### `server/service/agvc/agcv_main/coap_client_poller.go`
CoAP客户端轮询服务的核心实现：

**主要功能**：
- 初始化轮询服务，从配置读取参数
- 定时轮询所有并网点的数据
- 对每个并网点请求三种设备类型的数据
- 批量存储数据到内存和InfluxDB
- 支持动态修改服务器地址和轮询间隔

**关键方法**：
- `Initialize()`: 初始化服务并启动轮询goroutine
- `pollAllGridPoints()`: 轮询所有并网点
- `pollGridPointData(eqid)`: 轮询单个并网点的数据
- `requestData(eqid, eqType)`: 发送CoAP请求获取数据
- `storeData(messages)`: 存储数据到内存
- `Stop()`: 停止轮询服务

**包级变量**：
```go
var CoapClientPoller = new(coapClientPoller)
```

#### `server/service/agvc/agcv_main/COAP_CLIENT_POLLER_README.md`
详细的功能文档，包含：
- 功能描述和架构说明
- 数据请求规范
- 配置说明
- 轮询逻辑
- 使用方法

### 2. 修改文件

#### `server/config/coap.go`
添加CoAP客户端配置结构：
```go
type Coap struct {
    // 原有的服务端配置
    Enable bool
    Host   string
    Port   int
    
    // 新增的客户端配置
    ClientEnable   bool   // 是否启用客户端轮询
    ClientHost     string // CoAP服务器地址
    ClientPort     int    // CoAP服务器端口
    PollInterval   int    // 轮询间隔（秒）
}
```

#### `server/config.yaml`
添加配置项：
```yaml
coap:
  enable: true
  host: 0.0.0.0
  port: 1188
  # CoAP客户端轮询配置
  client-enable: true
  client-host: 127.0.0.1
  client-port: 5589
  poll-interval: 5  # 轮询间隔（秒）
```

#### `server/core/agcv.go`
在 `agvcInitialize()` 函数中添加轮询服务初始化：
```go
func agvcInitialize() {
    agvcMain.DataStorage.Initialize()         // 初始化数据存储
    agvcMain.DispatchStorage.Initialize()     // 初始化调度数据存储
    agvcMain.CoapClientPoller.Initialize()    // 初始化CoAP客户端轮询服务（新增）
}
```

## 数据流程

### 1. 启动流程
```
系统启动
  ↓
core/agcv.go: agvcInitialize()
  ↓
CoapClientPoller.Initialize()
  ↓
读取配置文件（client-enable, client-host, client-port, poll-interval）
  ↓
启动轮询goroutine
  ↓
立即执行第一次轮询
  ↓
启动定时器（每5秒触发一次）
```

### 2. 轮询流程
```
定时器触发
  ↓
从数据库读取所有并网点配置（agvc_bwd_setting表）
  ↓
遍历每个并网点
  ↓
对每个并网点请求三种类型数据：
  ├─ EQType=60 (AGC)
  ├─ EQType=61 (AVC)
  └─ EQType=5  (并网柜)
  ↓
发送CoAP POST请求到 /agvc/data
  ↓
接收响应数据（CoAPDataMessage数组）
  ↓
转换为RealtimeData格式
  ↓
批量存储到内存（DataStorage.StoreBatch）
  ↓
DataStorage自动定期保存到InfluxDB（5秒一次）
```

### 3. 请求格式
**请求**：
```http
POST /agvc/data
Content-Type: application/json

{
  "psid": 1,
  "eqid": <并网点编号>,
  "eqType": <设备类型>
}
```

**响应**：
```json
[
  {
    "psid": 1,
    "eqid": 1001,
    "eqType": 60,
    "dataType": 1,
    "point": "401",
    "value": 1
  },
  ...
]
```

## 数据存储

### 内存存储
使用 `DataStorage` 服务的 `StoreBatch` 方法：
- Key格式：`{psid}_{eqid}_{eqType}_{dataType}_{point}`
- 存储结构：`map[string]*agvc_main.RealtimeData`
- 线程安全：使用 `sync.RWMutex`

### InfluxDB持久化
- 自动触发：`DataStorage` 的 `periodicSave` 每5秒自动保存
- Measurement: `agvc_data`
- Tags: `psid`, `eqid`, `eqType`, `dataType`, `point`
- Field: `value`
- 时间戳：数据接收时间

## 配置说明

### 必填配置
- `client-enable`: 必须设置为 `true` 才能启用轮询
- `client-host`: CoAP服务器地址
- `client-port`: CoAP服务器端口（默认5589）

### 可选配置
- `poll-interval`: 轮询间隔（秒），默认5秒

### 配置示例
```yaml
# 开发环境
coap:
  client-enable: true
  client-host: 127.0.0.1
  client-port: 5589
  poll-interval: 5

# 生产环境
coap:
  client-enable: true
  client-host: 192.168.1.100
  client-port: 5589
  poll-interval: 3

# 禁用轮询
coap:
  client-enable: false
```

## 错误处理

### 连接失败
- 记录错误日志
- 继续处理下一个请求
- 不影响其他功能

### 数据解析失败
- 记录错误日志和原始数据
- 跳过该条数据
- 继续处理其他数据

### 超时处理
- 请求超时：5秒
- 超时后自动取消请求
- 记录超时日志

### 配置错误
- `client-enable=false`: 不启动轮询服务
- 端口为0: 使用默认端口5589
- 轮询间隔为0: 使用默认间隔5秒

## 日志监控

### 日志级别

#### Info级别
- 服务启动/停止
- 配置信息
- 服务器地址变更
- 轮询间隔变更

#### Debug级别
- 轮询开始/结束
- 数据获取成功
- 数据量统计

#### Error级别
- 数据库查询失败
- CoAP连接失败
- 请求发送失败
- 数据解析失败

### 监控指标
```
CoAP客户端轮询服务初始化
  server=127.0.0.1:5589
  interval=5s

开始轮询并网点数据
  并网点数量=3

成功获取并存储数据
  EQID=1001
  EQType=60
  数据量=15

请求数据失败
  EQID=1002
  EQType=61
  error=connection timeout
```

## 运行时管理

### 启动服务
自动启动，无需手动调用

### 停止服务
```go
agvcMain.CoapClientPoller.Stop()
```

### 动态修改配置

#### 修改服务器地址
```go
agvcMain.CoapClientPoller.SetServerAddress("192.168.1.200", 5589)
```

#### 修改轮询间隔
```go
agvcMain.CoapClientPoller.SetPollInterval(10 * time.Second)
```

## 性能考虑

### 资源占用
- **网络请求**: 每个并网点每次轮询3个请求（AGC、AVC、并网柜）
- **内存占用**: 取决于数据点数量，通常在MB级别
- **CPU占用**: 轮询和数据处理占用低

### 优化建议
1. **合理设置轮询间隔**: 根据业务需求调整，避免过于频繁
2. **并发控制**: 当前是串行处理，可考虑并发请求提高性能
3. **失败重试**: 可添加重试机制提高数据可靠性
4. **连接复用**: 可考虑复用UDP连接减少开销

## 测试方法

### 1. 模拟CoAP服务器
可以使用以下工具模拟CoAP服务器：
- go-coap库创建测试服务器
- Python CoAPthon库
- libcoap工具

### 2. 测试步骤
1. 启动CoAP测试服务器（5589端口）
2. 配置数据库中的并网点信息
3. 启动gin-vue-admin服务
4. 观察日志输出
5. 检查内存数据（DataStorage）
6. 检查InfluxDB数据

### 3. 验证要点
- [ ] 服务正常启动
- [ ] 定时轮询执行
- [ ] 所有并网点都被轮询
- [ ] 三种设备类型数据都被请求
- [ ] 数据正确存储到内存
- [ ] 数据正确持久化到InfluxDB
- [ ] 错误处理正确

## 兼容性说明

### 向后兼容
- 原有的CoAP服务端（1188端口）保留，不受影响
- 通过 `client-enable` 配置控制是否启用新功能
- 数据存储接口保持不变
- 不影响其他AGVC功能

### 数据格式
- 使用已有的 `CoAPDataMessage` 结构体
- 使用已有的 `RealtimeData` 结构体
- 存储格式与原有系统完全兼容

## 注意事项

1. **并网点配置**: 确保 `agvc_bwd_setting` 表中有正确的并网点配置
2. **CoAP服务器**: 确保目标CoAP服务器（5589端口）正常运行
3. **网络连接**: 确保服务器可以访问CoAP服务器
4. **数据格式**: CoAP服务器返回的数据格式必须符合 `CoAPDataMessage` 数组
5. **配置启用**: 必须设置 `client-enable: true` 才能启用功能

## 故障排查

### 服务未启动
- 检查配置：`client-enable` 是否为 `true`
- 检查日志：是否有初始化错误

### 没有数据
- 检查CoAP服务器是否运行
- 检查网络连接
- 检查并网点配置是否存在
- 检查日志中的错误信息

### 数据不完整
- 检查CoAP服务器响应
- 检查数据格式是否正确
- 检查日志中的解析错误

### 性能问题
- 检查轮询间隔设置
- 检查并网点数量
- 检查网络延迟
- 考虑优化轮询策略

## 未来优化方向

1. **并发请求**: 使用goroutine池并发处理多个并网点
2. **失败重试**: 添加重试机制提高数据可靠性
3. **连接池**: 复用UDP连接减少连接开销
4. **数据缓存**: 添加本地缓存减少请求频率
5. **动态调整**: 根据负载自动调整轮询间隔
6. **健康检查**: 添加CoAP服务器健康检查机制
7. **统计信息**: 添加详细的统计和监控指标

## 总结

本次实现完成了从被动接收到主动轮询的架构转变，主要优势：

✅ **主动控制**: 系统主动控制数据获取频率和时机
✅ **配置灵活**: 支持动态配置服务器地址和轮询间隔
✅ **向后兼容**: 不影响现有功能，平滑迁移
✅ **错误处理**: 完善的错误处理和日志记录
✅ **易于维护**: 代码结构清晰，文档完整

实现的核心文件：
- ✅ `coap_client_poller.go` - 核心实现
- ✅ `config/coap.go` - 配置结构
- ✅ `config.yaml` - 配置文件
- ✅ `core/agcv.go` - 启动集成
- ✅ 完整文档
