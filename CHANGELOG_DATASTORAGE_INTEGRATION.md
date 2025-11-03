# DataStorage 集成改造说明

## 改造概述

将 `coap_handler` 中的 `SaveAgvcData` 功能与 `agvc_main` 里面的 `dataStorage` 结合起来，实现接收到数据后每5分钟定时保存到 InfluxDB。

## 改造内容

### 1. 数据模型扩展 (`server/model/agvc/agvc_main/device.go`)

新增 `AgvcDataItem` 结构体，用于接收 CoAP 数据：

```go
// AgvcDataItem AGVC数据项（从CoAP接收的数据格式）
type AgvcDataItem struct {
    Psid     int     `json:"psid"`     // 电站ID
    Eqid     int     `json:"eqid"`     // 设备ID
    EqType   int     `json:"eqType"`   // 设备类型
    DataType int     `json:"dataType"` // 数据类型
    Point    string  `json:"point"`    // 数据点
    Value    float64 `json:"value"`    // 数值
}
```

### 2. DataStorage 新增批量存储方法 (`server/service/agvc/agcv_main/data_storage.go`)

新增 `StoreAgvcDataBatch` 方法，支持从 CoAP 接收的数据格式批量存储：

```go
// StoreAgvcDataBatch 批量存储AGVC数据（从CoAP接收的数据格式）
func (s *dataStorage) StoreAgvcDataBatch(dataBatch []agvc_main.AgvcDataItem) {
    s.mu.Lock()
    defer s.mu.Unlock()

    now := time.Now().Unix()
    for _, item := range dataBatch {
        // 转换为内部RealtimeData格式
        data := &agvc_main.RealtimeData{
            PSID:      fmt.Sprintf("%03d", item.Psid),
            EQID:      fmt.Sprintf("%04d", item.Eqid),
            EQType:    fmt.Sprintf("%02d", item.EqType),
            DataType:  fmt.Sprintf("%02d", item.DataType),
            Point:     item.Point,
            Value:     item.Value,
            Timestamp: now,
        }
        key := s.makeKey(data.PSID, data.EQID, data.EQType, data.DataType, data.Point)
        s.realtimeData[key] = data
    }
}
```

**关键点：**
- 自动将 int 类型的 ID 转换为固定长度的字符串格式
- Psid: 3位数字 (如 001, 002)
- Eqid: 4位数字 (如 0001, 0002)
- EqType: 2位数字 (如 01, 02)
- DataType: 2位数字 (如 01, 02)

### 3. 改进定时保存逻辑

优化 `periodicSave` 方法：
- 增加数据量检查，空数据时跳过保存
- 优化日志输出
- 避免重复加锁

### 4. 统一 InfluxDB Measurement 名称

将 `data_storage.go` 中的 measurement 名称从 `"device_data"` 统一为 `"agvc_data"`，与现有的 `agvc_data.go` 保持一致。

### 5. CoAP Handler 改造 (`server/initialize/coap_handler.go`)

**改造前：**
```go
// 调用服务层保存数据
agvcDataService := service.ServiceGroupApp.AgvcServiceGroup.AgvcDataService
if err := agvcDataService.SaveAgvcData(ctx, dataBatch); err != nil {
    global.GVA_LOG.Error("Failed to save AGVC data to InfluxDB", zap.Error(err))
    return coapCodeInternalServerError, []byte(`{"error":"save failed"}`)
}
```

**改造后：**
```go
// 转换数据格式并存储到DataStorage
convertedData := make([]agvc_main.AgvcDataItem, len(dataBatch))
for i, item := range dataBatch {
    convertedData[i] = agvc_main.AgvcDataItem{
        Psid:     item.Psid,
        Eqid:     item.Eqid,
        EqType:   item.EqType,
        DataType: item.DataType,
        Point:    item.Point,
        Value:    item.Value,
    }
}

// 存储到内存，由DataStorage每5分钟定时保存到InfluxDB
agvcMainService.DataStorage.StoreAgvcDataBatch(convertedData)

global.GVA_LOG.Debug("AGVC data stored to memory", zap.Int("count", len(dataBatch)))
```

## 工作流程

```
CoAP 数据接收
    ↓
handleAgvcData (coap_handler.go)
    ↓
数据格式转换 (int → string)
    ↓
DataStorage.StoreAgvcDataBatch
    ↓
存储到内存 Map (realtimeData)
    ↓
定时器触发 (每5分钟)
    ↓
saveToInfluxDB
    ↓
批量写入 InfluxDB
```

## 优势

1. **内存缓存**：数据先存储在内存中，提供实时访问能力
2. **批量写入**：每5分钟批量写入 InfluxDB，减少数据库压力
3. **数据聚合**：相同测点的数据会被最新值覆盖，避免重复存储
4. **统一管理**：所有实时数据通过 DataStorage 统一管理
5. **类型安全**：自动处理数据类型转换（int → string）

## 数据流

### 输入数据格式 (CoAP)
```json
[
  {
    "psid": 1,
    "eqid": 1,
    "eqType": 1,
    "dataType": 2,
    "point": "401",
    "value": 123.45
  }
]
```

### 内存存储格式 (RealtimeData)
```go
{
    PSID:      "001",
    EQID:      "0001",
    EQType:    "01",
    DataType:  "02",
    Point:     "401",
    Value:     123.45,
    Timestamp: 1234567890
}
```

### InfluxDB 存储格式
- **Measurement**: `agvc_data`
- **Tags**: 
  - psid: "001"
  - eqid: "0001"
  - eqType: "01"
  - dataType: "02"
  - point: "401"
- **Fields**: 
  - value: 123.45
- **Timestamp**: 当前时间

## 配置项

定时保存间隔在 `data_storage.go` 的 `Initialize` 方法中设置：

```go
// 启动5分钟定时保存到InfluxDB
s.saveTimer = time.NewTicker(5 * time.Minute)
```

如需修改保存间隔，可调整此处的时间参数。

## 注意事项

1. **数据覆盖**：相同 key (psid_eqid_eqType_dataType_point) 的数据会被最新值覆盖
2. **内存管理**：数据会一直保存在内存中，直到程序重启
3. **数据丢失风险**：程序异常退出时，距离上次保存后的数据可能会丢失
4. **类型转换**：确保输入的 int 类型 ID 值在合理范围内（psid: 0-999, eqid: 0-9999, etc.）

### 6. 优雅关闭支持 (`server/core/server_run.go`)

在服务器优雅关闭时，自动保存内存中的数据到 InfluxDB：

```go
// 停止数据存储服务（保存数据到InfluxDB）
agvcMain.DataStorage.Stop()
```

**Stop 方法实现：**
```go
func (s *dataStorage) Stop() {
    // 停止定时器
    if s.saveTimer != nil {
        s.saveTimer.Stop()
    }
    
    // 关闭前强制保存一次数据
    s.mu.RLock()
    dataCount := len(s.realtimeData)
    s.mu.RUnlock()
    
    if dataCount > 0 {
        global.GVA_LOG.Info("服务关闭前保存数据到InfluxDB", zap.Int("数据量", dataCount))
        if err := s.saveToInfluxDB(); err != nil {
            global.GVA_LOG.Error("服务关闭前保存数据失败", zap.Error(err))
        } else {
            global.GVA_LOG.Info("服务关闭前数据保存成功")
        }
    }
    
    close(s.stopChan)
    global.GVA_LOG.Info("数据存储服务已停止")
}
```

## 后续优化建议

1. ✅ ~~**持久化策略**：考虑程序关闭时强制保存一次数据~~ (已实现)
2. **数据清理**：定期清理过期的内存数据
3. **可配置化**：将保存间隔、数据保留策略等配置化
4. **监控指标**：添加数据保存成功率、内存使用量等监控指标
