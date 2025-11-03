# SaveAgvcData 与 DataStorage 集成使用指南

## 概述

本次改造将 CoAP 数据接收的即时保存机制改为基于内存缓存的定时批量保存机制，提高了系统性能并支持实时数据访问。

## 功能特性

### ✅ 主要特性

1. **内存缓存**
   - 接收到的数据先存储在内存中
   - 支持高速读写操作
   - 自动管理数据更新和覆盖

2. **定时批量保存**
   - 每 5 分钟自动保存到 InfluxDB
   - 减少数据库写入压力
   - 提高整体系统性能

3. **优雅关闭**
   - 程序退出前自动保存所有数据
   - 避免数据丢失
   - 支持 SIGINT 和 SIGTERM 信号

4. **自动类型转换**
   - 自动将 int 类型 ID 转换为固定长度字符串
   - 确保数据格式一致性
   - 简化数据访问逻辑

## 数据流程

```
┌─────────────┐
│ CoAP 客户端  │
└──────┬──────┘
       │ 发送数据
       ▼
┌─────────────────────┐
│   CoAP Handler      │
│ handleAgvcData      │
└──────┬──────────────┘
       │ 数据转换
       ▼
┌─────────────────────┐
│   DataStorage       │
│ StoreAgvcDataBatch  │
│   (内存 Map)        │
└──────┬──────────────┘
       │
       │ ┌─────────────┐
       ├─►│ 实时查询API  │
       │ └─────────────┘
       │
       │ 每5分钟
       ▼
┌─────────────────────┐
│   InfluxDB          │
│ agvc_data measurement│
└─────────────────────┘
```

## API 使用

### 1. 存储数据

CoAP 客户端发送数据到 `/agvc/data` 端点：

```json
POST /agvc/data
Content-Type: application/json

[
  {
    "psid": 1,
    "eqid": 1001,
    "eqType": 1,
    "dataType": 2,
    "point": "401",
    "value": 123.45
  },
  {
    "psid": 1,
    "eqid": 1001,
    "eqType": 1,
    "dataType": 2,
    "point": "402",
    "value": 67.89
  }
]
```

响应：
```json
{
  "success": true
}
```

### 2. 访问实时数据

在代码中访问实时数据：

```go
import agvcMain "github.com/flipped-aurora/gin-vue-admin/server/service/agvc/agcv_main"

// 获取单个数据点
data, exists := agvcMain.DataStorage.GetData("001", "1001", "01", "02", "401")
if exists {
    fmt.Printf("Value: %v\n", data.Value)
}

// 获取设备所有数据
deviceData := agvcMain.DataStorage.GetDeviceData("001", "1001", "01", "02")
for point, data := range deviceData {
    fmt.Printf("Point %s: %v\n", point, data.Value)
}

// 获取数据为 float64
value, err := agvcMain.DataStorage.GetDataAsFloat64("001", "1001", "01", "02", "401")
if err == nil {
    fmt.Printf("Float Value: %.2f\n", value)
}

// 获取数据为 int
intValue, err := agvcMain.DataStorage.GetDataAsInt("001", "1001", "01", "02", "401")
if err == nil {
    fmt.Printf("Int Value: %d\n", intValue)
}
```

## 配置说明

### 定时保存间隔

修改 `server/service/agvc/agcv_main/data_storage.go` 中的 `Initialize` 方法：

```go
func (s *dataStorage) Initialize() {
    s.realtimeData = make(map[string]*agvc_main.RealtimeData)
    s.stopChan = make(chan struct{})

    // 修改这里的时间间隔
    s.saveTimer = time.NewTicker(5 * time.Minute)  // 默认 5 分钟
    go s.periodicSave()

    global.GVA_LOG.Info("数据存储服务初始化成功")
}
```

可选的时间间隔：
- `1 * time.Minute` - 1 分钟
- `5 * time.Minute` - 5 分钟 (默认)
- `10 * time.Minute` - 10 分钟
- `30 * time.Minute` - 30 分钟

### InfluxDB 配置

确保 `config.yaml` 中配置了正确的 InfluxDB 连接信息：

```yaml
influxdb:
  url: "http://localhost:8086"
  token: "your-token"
  org: "your-org"
  bucket: "your-bucket"
```

## 数据存储格式

### 内存存储格式

```go
Key: "001_1001_01_02_401"
Value: {
    PSID:      "001",
    EQID:      "1001",
    EQType:    "01",
    DataType:  "02",
    Point:     "401",
    Value:     123.45,
    Timestamp: 1699012345
}
```

### InfluxDB 存储格式

```
Measurement: agvc_data
Tags:
  - psid: "001"
  - eqid: "1001"
  - eqType: "01"
  - dataType: "02"
  - point: "401"
Fields:
  - value: 123.45
Timestamp: 2024-11-03T08:00:00Z
```

## 日志说明

### 正常运行日志

```
INFO  数据存储服务初始化成功
DEBUG AGVC data stored to memory  count=10
INFO  数据已保存到InfluxDB  数据量=100
```

### 关闭时日志

```
INFO  关闭WEB服务...
INFO  服务关闭前保存数据到InfluxDB  数据量=50
INFO  服务关闭前数据保存成功
INFO  数据存储服务已停止
INFO  WEB服务已关闭
```

### 错误日志

```
ERROR 保存数据到InfluxDB失败  error="connection refused"
ERROR 服务关闭前保存数据失败  error="write timeout"
```

## 监控指标

### 关键指标

1. **内存数据量**
   ```go
   count := agvcMain.DataStorage.GetDataCount()
   fmt.Printf("Current data count: %d\n", count)
   ```

2. **数据快照**
   ```go
   snapshot := agvcMain.DataStorage.ExportSnapshot()
   fmt.Println(snapshot)
   ```

## 性能优化建议

### 1. 调整保存间隔

根据业务需求调整保存间隔：
- 实时性要求高：1-2 分钟
- 平衡性能：5 分钟 (推荐)
- 性能优先：10-30 分钟

### 2. 数据清理策略

未来可以添加数据清理策略：

```go
// 示例：清理超过 1 小时的数据
func (s *dataStorage) CleanOldData(maxAge time.Duration) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    now := time.Now().Unix()
    for key, data := range s.realtimeData {
        if now - data.Timestamp > int64(maxAge.Seconds()) {
            delete(s.realtimeData, key)
        }
    }
}
```

### 3. 批量大小控制

当数据量很大时，可以分批保存：

```go
const batchSize = 1000

// 在 saveToInfluxDB 中分批写入
for i := 0; i < len(points); i += batchSize {
    end := i + batchSize
    if end > len(points) {
        end = len(points)
    }
    batch := points[i:end]
    writeAPI.WritePoint(ctx, batch...)
}
```

## 故障排查

### 问题 1: 数据没有保存到 InfluxDB

**可能原因：**
1. InfluxDB 连接配置错误
2. 网络连接问题
3. 权限不足

**解决方法：**
```bash
# 检查 InfluxDB 连接
curl http://localhost:8086/health

# 查看日志
tail -f logs/server.log | grep InfluxDB

# 测试写入
influx write -b your-bucket -o your-org -p s 'agvc_data,psid=001,eqid=0001 value=100'
```

### 问题 2: 内存占用过高

**可能原因：**
1. 数据累积过多
2. 没有定期清理

**解决方法：**
1. 减少保存间隔
2. 实现数据清理策略
3. 监控内存使用

### 问题 3: 程序退出时数据丢失

**可能原因：**
1. 强制终止（kill -9）
2. Stop() 方法未调用

**解决方法：**
1. 使用优雅关闭（kill 或 Ctrl+C）
2. 确保 server_run.go 中调用了 DataStorage.Stop()

## 测试

### 运行集成测试

```bash
cd /home/engine/project
./test_integration.sh
```

### 手动测试

1. **启动服务**
   ```bash
   cd server
   go run main.go
   ```

2. **发送测试数据**
   ```bash
   # 使用 CoAP 客户端发送数据
   coap-client -m post -t json -e '[{"psid":1,"eqid":1,"eqType":1,"dataType":2,"point":"401","value":123.45}]' coap://localhost:5683/agvc/data
   ```

3. **等待 5 分钟或优雅关闭**
   ```bash
   # Ctrl+C 触发优雅关闭
   ```

4. **检查 InfluxDB**
   ```bash
   influx query 'from(bucket:"your-bucket") |> range(start: -1h) |> filter(fn: (r) => r._measurement == "agvc_data")'
   ```

## 相关文档

- [详细改造说明](CHANGELOG_DATASTORAGE_INTEGRATION.md)
- [改造摘要](INTEGRATION_SUMMARY.md)
- [集成测试脚本](test_integration.sh)

## 注意事项

⚠️ **重要提示**

1. **数据丢失风险**
   - 程序异常退出时，最多丢失 5 分钟数据
   - 建议配合监控系统使用

2. **内存管理**
   - 数据会永久保存在内存中
   - 需要根据业务场景评估内存需求

3. **数据覆盖**
   - 相同测点的数据会被覆盖
   - 仅保留最新值

4. **并发安全**
   - 已使用 RWMutex 保证并发安全
   - 可以从多个 goroutine 安全访问

## 支持

如有问题，请：
1. 查看日志文件
2. 运行集成测试
3. 检查 InfluxDB 连接
4. 查阅相关文档

---

**更新时间**: 2025  
**版本**: v1.0  
**状态**: ✅ 生产就绪
