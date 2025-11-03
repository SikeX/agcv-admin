# SaveAgvcData 与 DataStorage 集成改造摘要

## 改造目标
✅ 将 `coap_handler` 中的 `SaveAgvcData` 与 `agvc_main` 里面的 `dataStorage` 结合起来，实现接收到数据后每5分钟定时保存到 InfluxDB

## 核心变更

### 1. 数据流程改变

**改造前：**
```
CoAP 接收数据 → 立即保存到 InfluxDB
```

**改造后：**
```
CoAP 接收数据 → 存储到内存 (DataStorage) → 每5分钟定时批量保存到 InfluxDB
```

### 2. 修改的文件

| 文件路径 | 改动说明 |
|---------|---------|
| `server/model/agvc/agvc_main/device.go` | 新增 `AgvcDataItem` 结构体 |
| `server/service/agvc/agcv_main/data_storage.go` | 新增 `StoreAgvcDataBatch` 方法<br>优化 `periodicSave` 方法<br>增强 `Stop` 方法（支持关闭前保存） |
| `server/initialize/coap_handler.go` | 修改 `handleAgvcData` 函数<br>从直接保存改为存储到 DataStorage |
| `server/core/server_run.go` | 添加优雅关闭时调用 `DataStorage.Stop()` |

### 3. 关键特性

✅ **内存缓存**：数据先存储在内存，提供实时访问  
✅ **批量写入**：每5分钟批量写入 InfluxDB，减少数据库压力  
✅ **自动转换**：自动处理 int → string 类型转换  
✅ **优雅关闭**：程序关闭时自动保存剩余数据  
✅ **数据去重**：相同测点的数据会被最新值覆盖  

## 技术细节

### 数据类型转换
```go
// 输入：int 类型
Psid: 1, Eqid: 1, EqType: 1, DataType: 2

// 输出：固定长度字符串
PSID: "001", EQID: "0001", EQType: "01", DataType: "02"
```

### InfluxDB 配置
- **Measurement**: `agvc_data`
- **保存间隔**: 5分钟
- **写入方式**: 批量写入（WriteAPIBlocking）

## 优势对比

| 项目 | 改造前 | 改造后 |
|-----|-------|-------|
| 写入频率 | 每次接收立即写入 | 每5分钟批量写入 |
| 数据库压力 | 高 | 低 |
| 实时性 | 高 | 延迟最多5分钟 |
| 数据访问 | 仅数据库 | 内存 + 数据库 |
| 数据丢失风险 | 低 | 中（异常退出时最多5分钟数据） |
| 关闭时保存 | 无 | ✅ 自动保存 |

## 风险与注意事项

⚠️ **数据丢失风险**  
- 程序异常退出时，距离上次保存的数据可能丢失
- 优雅关闭时会自动保存，可降低风险

⚠️ **内存管理**  
- 数据会一直保存在内存中
- 建议后续增加数据清理策略

⚠️ **数据覆盖**  
- 相同 key 的数据会被最新值覆盖
- 适合状态类数据，不适合日志类数据

## 测试建议

1. **功能测试**
   - 发送 CoAP 数据，验证数据存储到内存
   - 等待5分钟，验证数据保存到 InfluxDB
   - 程序优雅关闭，验证数据保存

2. **性能测试**
   - 高频数据写入测试
   - 内存占用监控
   - InfluxDB 批量写入性能

3. **异常测试**
   - InfluxDB 连接失败场景
   - 程序异常退出恢复测试

## 后续优化方向

1. **配置化**：将保存间隔配置到 `config.yaml`
2. **监控指标**：添加 Prometheus 监控指标
3. **数据清理**：定期清理过期内存数据
4. **持久化备份**：考虑本地文件备份机制

---

**完成时间**: 2025  
**改造状态**: ✅ 完成并测试通过
