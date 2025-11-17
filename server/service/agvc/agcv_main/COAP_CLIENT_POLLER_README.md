# CoAP 客户端轮询服务

## 功能描述

CoAP客户端轮询服务用于主动从CoAP服务器（端口5589）请求AGVC相关数据，替代原有的被动接收模式。

## 架构变更

### 原架构
- **模式**：被动接收
- **实现**：启动CoAP服务端（1188端口）监听数据
- **流程**：等待外部设备推送数据

### 新架构
- **模式**：主动轮询
- **实现**：启动CoAP客户端定期请求数据
- **流程**：定期向CoAP服务器（5589端口）请求数据

## 数据请求规范

### 请求接口
- **URL**: `/agvc/data`
- **方法**: POST
- **Content-Type**: application/json

### 请求参数
```json
{
  "psid": 1,        // 电站ID（固定值）
  "eqid": <并网点编号>,  // 从 agvc_bwd_setting 表的 number 字段获取
  "eqType": <设备类型>   // 60(AGC) / 61(AVC) / 5(并网柜)
}
```

### 设备类型说明
- **60**: AGC（自动有功功率控制）
- **61**: AVC（自动无功电压控制）
- **5**: 并网柜

### 响应数据格式
```json
[
  {
    "psid": 1,
    "eqid": <并网点编号>,
    "eqType": <设备类型>,
    "dataType": <数据类型>,
    "point": "<点位标识>",
    "value": <数值>
  },
  ...
]
```

## 配置说明

### 配置文件位置
`server/config.yaml`

### 配置项
```yaml
coap:
  # CoAP服务端配置（原有配置，保留）
  enable: true
  host: 0.0.0.0
  port: 1188
  
  # CoAP客户端轮询配置（新增）
  client-enable: true          # 是否启用客户端轮询
  client-host: 127.0.0.1       # CoAP服务器地址
  client-port: 5589            # CoAP服务器端口
  poll-interval: 5             # 轮询间隔（秒）
```

### 配置说明
- `client-enable`: 控制是否启用客户端轮询功能
- `client-host`: CoAP服务器的IP地址
- `client-port`: CoAP服务器的端口，默认5589
- `poll-interval`: 轮询间隔，单位为秒，默认5秒

## 数据存储

### 存储位置
1. **内存存储**: 使用 `DataStorage` 服务存储到内存Map
2. **InfluxDB**: 通过 `DataStorage` 的定期保存机制自动写入InfluxDB

### 存储流程
```
CoAP请求 → 获取数据 → 转换格式 → 存入内存 → 定期持久化到InfluxDB
```

## 轮询逻辑

### 启动流程
1. 系统启动时初始化CoAP客户端轮询服务
2. 从配置文件读取服务器地址和轮询间隔
3. 启动轮询goroutine

### 轮询流程
1. 从数据库读取所有并网点配置（`agvc_bwd_setting`表）
2. 遍历每个并网点
3. 对每个并网点分别请求三种设备类型的数据：
   - AGC (EQType=60)
   - AVC (EQType=61)
   - 并网柜 (EQType=5)
4. 将获取的数据批量存储到内存
5. 等待下一个轮询周期

### 错误处理
- 连接失败：记录错误日志，继续处理下一个请求
- 数据解析失败：记录错误日志，继续处理下一个请求
- 响应超时：5秒超时时间，超时后记录错误并继续

## 代码文件

### 新增文件
- `server/service/agvc/agcv_main/coap_client_poller.go`: CoAP客户端轮询服务实现

### 修改文件
- `server/core/agcv.go`: 在启动时初始化轮询服务
- `server/config/coap.go`: 添加客户端配置结构
- `server/config.yaml`: 添加客户端配置项

## 使用的数据结构

### CoAPDataMessage
```go
type CoAPDataMessage struct {
    PSID     int         `json:"psid" binding:"required"`
    EQID     int         `json:"eqid" binding:"required"`
    EQType   int         `json:"eqType" binding:"required"`
    DataType int         `json:"dataType" binding:"required"`
    Point    string      `json:"point" binding:"required"`
    Value    interface{} `json:"value" binding:"required"`
}
```

### RealtimeData
```go
type RealtimeData struct {
    PSID      int         `json:"psid"`
    EQID      int         `json:"eqid"`
    EQType    int         `json:"eqType"`
    DataType  int         `json:"dataType"`
    Point     string      `json:"point"`
    Value     interface{} `json:"value"`
    Timestamp int64       `json:"timestamp"`
}
```

## 调试和监控

### 日志级别
- **Info**: 服务启动、停止、配置信息
- **Debug**: 轮询详情、数据获取成功
- **Error**: 连接失败、数据解析失败、其他错误

### 监控指标
- 轮询周期执行情况
- 数据获取成功/失败次数
- 数据存储量
- 响应时间

## 注意事项

1. **并网点配置**: 确保 `agvc_bwd_setting` 表中有正确的并网点配置
2. **服务器地址**: 确保CoAP服务器（5589端口）正常运行
3. **轮询间隔**: 根据实际需求调整轮询间隔，避免过于频繁
4. **资源占用**: 轮询会产生网络请求和数据存储，注意资源占用
5. **错误处理**: 轮询失败不会影响其他功能，会继续下一次轮询

## 停止服务

服务会在系统关闭时自动停止。如需手动停止：
```go
agvcMain.CoapClientPoller.Stop()
```

## 运行时配置修改

可以在运行时动态修改配置：

### 修改服务器地址
```go
agvcMain.CoapClientPoller.SetServerAddress("192.168.1.100", 5589)
```

### 修改轮询间隔
```go
agvcMain.CoapClientPoller.SetPollInterval(10 * time.Second)
```
