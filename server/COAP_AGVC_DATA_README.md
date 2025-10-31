# CoAP AGVC Data 接口文档

## 功能概述

本功能实现了一个CoAP接口 `/agvc/data`，用于接收AGVC设备的数据，并将数据存储到InfluxDB中。

## 接口信息

- **路径**: `/agvc/data`
- **方法**: POST (CoAP Code: 2)
- **协议**: CoAP (UDP)
- **默认端口**: 5683
- **数据格式**: JSON

## 数据结构

### 请求体 (JSON数组)

```json
[
  {
    "psid": 1,          // 电站ID (整数)
    "eqid": 1,          // 设备ID (整数)
    "eqType": 2,        // 设备类型 (整数)
    "dataType": 2,      // 数据类型 (整数)
    "point": "2",       // 数据点 (字符串)
    "value": 32.32      // 数值 (浮点数)
  }
]
```

### 响应

#### 成功响应 (Code: 65 - Created)
```json
{"success": true}
```

#### 错误响应

- **Code: 128 (Bad Request)**: 数据格式错误或数据为空
  ```json
  {"error": "invalid json"}
  ```
  或
  ```json
  {"error": "empty data"}
  ```

- **Code: 160 (Internal Server Error)**: 保存到InfluxDB失败
  ```json
  {"error": "save failed"}
  ```

## InfluxDB 存储格式

### Measurement
- **名称**: `agvc_data`

### Tags
- `psid`: 电站ID
- `eqid`: 设备ID
- `eqType`: 设备类型
- `dataType`: 数据类型
- `point`: 数据点

### Fields
- `value`: 数值

### 时间戳
- 使用接收数据时的服务器时间

## 配置

### CoAP 服务配置 (config.yaml)

```yaml
coap:
  enable: true        # 启用CoAP服务
  host: "0.0.0.0"    # 监听地址
  port: 5683         # 监听端口
```

### InfluxDB 配置 (config.yaml)

```yaml
influxdb:
  host: "localhost"
  port: "8086"
  token: "your-influxdb-token"
  org: "your-org"
  bucket: "your-bucket"
  username: "your-username"  # 可选
  password: "your-password"  # 可选
```

## 代码结构

### 模型层 (Model)
- **文件**: `server/model/agvc/agvc_data.go`
- **结构体**: 
  - `AgvcDataItem`: 单条数据项
  - `AgvcDataBatch`: 数据批量数组

### 服务层 (Service)
- **文件**: `server/service/agvc/agvc_data.go`
- **服务**: `AgvcDataService`
- **方法**: `SaveAgvcData(ctx, dataBatch)` - 保存数据到InfluxDB

### 初始化层 (Initialize)
- **文件**: `server/initialize/coap_handler.go`
- **处理器**: `handleAgvcData` - 处理CoAP请求
- **路由表**: `coapRoutes` - CoAP路由映射

### CoAP服务器
- **文件**: `server/initialize/coap.go`
- **功能**: CoAP协议实现和请求分发

## 测试

### 使用Python测试脚本

```bash
# 确保服务器已启动
cd /home/engine/project/server
python3 test_coap_client.py
```

### 使用CoAP客户端工具

```bash
# 使用 coap-client (需要安装 libcoap)
echo '[{"psid":1,"eqid":1,"eqType":2,"dataType":2,"point":"2","value":32.32}]' | \
  coap-client -m post -t application/json coap://localhost:5683/agvc/data
```

### 使用curl (通过CoAP-HTTP代理)

如果有CoAP-HTTP代理，可以使用curl测试：

```bash
curl -X POST http://localhost:5683/agvc/data \
  -H "Content-Type: application/json" \
  -d '[{"psid":1,"eqid":1,"eqType":2,"dataType":2,"point":"2","value":32.32}]'
```

## 查询InfluxDB数据

使用InfluxDB CLI或UI查询数据：

```flux
from(bucket: "your-bucket")
  |> range(start: -1h)
  |> filter(fn: (r) => r._measurement == "agvc_data")
  |> filter(fn: (r) => r.psid == "1")
```

或使用SQL风格查询（InfluxDB 1.x）：

```sql
SELECT * FROM agvc_data WHERE psid='1' AND time > now() - 1h
```

## 日志

服务器会记录以下日志：

- **接收数据**: `AGVC data written to InfluxDB`
- **解析错误**: `Failed to parse AGVC data`
- **保存失败**: `Failed to save AGVC data to InfluxDB`
- **保存成功**: `AGVC data saved successfully`

## 错误处理

1. **InfluxDB未初始化**: 返回错误并记录日志
2. **JSON解析失败**: 返回Bad Request响应
3. **数据为空**: 返回Bad Request响应
4. **InfluxDB写入失败**: 返回Internal Server Error响应

## 性能考虑

- 使用InfluxDB的异步写入API，提高写入性能
- 支持批量数据接收，减少网络往返
- CoAP协议基于UDP，网络开销小

## 扩展性

可以通过修改 `coapRoutes` 映射表添加更多CoAP路由：

```go
var coapRoutes = map[string]map[byte]CoapHandler{
    "/agvc/data": {
        coapCodePost: handleAgvcData,
    },
    "/your/new/path": {
        coapCodeGet: handleYourNewPath,
    },
}
```

## 安全建议

1. 使用防火墙限制CoAP端口的访问
2. 考虑实现CoAP DTLS加密
3. 实施数据验证和速率限制
4. 定期审计InfluxDB访问日志

## 故障排查

### 服务器未启动
检查配置文件中 `coap.enable` 是否为 `true`

### 连接超时
检查防火墙设置和端口是否被占用

### 数据未写入InfluxDB
- 检查InfluxDB配置是否正确
- 查看服务器日志获取详细错误信息
- 验证InfluxDB服务是否运行

### CoAP请求失败
- 使用 `tcpdump` 或 `wireshark` 抓包分析
- 检查CoAP消息格式是否正确
- 查看服务器日志

## 相关文件

- `server/model/agvc/agvc_data.go` - 数据模型
- `server/service/agvc/agvc_data.go` - 业务逻辑
- `server/service/agvc/enter.go` - 服务注册
- `server/initialize/coap.go` - CoAP服务器
- `server/initialize/coap_handler.go` - CoAP路由处理
- `server/test_coap_client.py` - Python测试脚本

## 版本历史

- v1.0.0: 初始版本，实现基本的CoAP数据接收和InfluxDB存储功能
