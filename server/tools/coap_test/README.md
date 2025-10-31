# CoAP测试客户端

## 功能说明

这是一个用于测试 `/agvc/data` CoAP接口的测试客户端工具。

## 构建

```bash
cd /home/engine/project/server/tools/coap_test
go build -o coap_test_client
```

## 使用方法

### 1. 启动GVA服务器

首先确保InfluxDB已启动并正确配置，然后启动GVA服务器：

```bash
cd /home/engine/project/server
./server  # 或者 go run main.go
```

### 2. 运行测试客户端

在另一个终端窗口中运行：

```bash
cd /home/engine/project/server/tools/coap_test
./coap_test_client
```

## 测试场景

测试客户端会执行以下测试场景：

1. **单个数据点测试**：发送一条数据
2. **多个数据点测试**：发送多条来自同一电站的数据
3. **不同电站测试**：发送来自不同电站的数据

## 预期输出

成功的输出示例：

```
CoAP客户端测试
================

测试1: 单个数据点
------------------
发送数据: [{"psid":1,"eqid":1,"eqType":2,"dataType":2,"point":"2","value":32.32}]
响应代码: 65
响应内容: {"success":true}
✓ 数据发送成功!

测试2: 多个数据点
------------------
发送数据: [{"psid":1,"eqid":101,"eqType":1,"dataType":1,"point":"temperature","value":25.5}...]
响应代码: 65
响应内容: {"success":true}
✓ 数据发送成功!

测试3: 不同电站
------------------
发送数据: [{"psid":1,"eqid":1,"eqType":2,"dataType":2,"point":"2","value":32.32}...]
响应代码: 65
响应内容: {"success":true}
✓ 数据发送成功!

所有测试完成!
```

## 响应代码说明

- **65 (Created)**: 数据创建成功
- **128 (Bad Request)**: 请求格式错误或数据为空
- **132 (Not Found)**: 路径未找到
- **133 (Method Not Allowed)**: HTTP方法不允许
- **160 (Internal Server Error)**: 服务器内部错误

## 验证数据

### 查看服务器日志

服务器日志会显示接收到的数据：

```
INFO    AGVC data written to InfluxDB   {"psid": 1, "eqid": 1, "eqType": 2, "dataType": 2, "point": "2", "value": 32.32}
INFO    AGVC data saved successfully    {"count": 1}
```

### 查询InfluxDB

使用InfluxDB CLI查询数据：

```bash
# InfluxDB 2.x
influx query 'from(bucket:"test") |> range(start: -1h) |> filter(fn: (r) => r._measurement == "agvc_data")'

# 或使用HTTP API
curl -XPOST "http://localhost:8086/api/v2/query?org=test" \
  -H "Authorization: Token YOUR_TOKEN" \
  -H "Content-Type: application/vnd.flux" \
  -d 'from(bucket:"test") |> range(start: -1h) |> filter(fn: (r) => r._measurement == "agvc_data")'
```

## 自定义测试

可以修改 `main.go` 中的测试数据来进行自定义测试：

```go
customData := []AgvcDataItem{
    {
        Psid:     100,
        Eqid:     200,
        EqType:   3,
        DataType: 4,
        Point:    "custom_point",
        Value:    99.99,
    },
}
if err := sendCoapRequest(host, port, path, customData); err != nil {
    fmt.Printf("错误: %v\n", err)
}
```

## 故障排查

### 连接超时

- 确保GVA服务器已启动
- 检查防火墙设置
- 确认CoAP配置中 `enable: true`

### 数据未写入InfluxDB

- 检查InfluxDB服务是否运行
- 验证config.yaml中的InfluxDB配置
- 查看GVA服务器日志获取详细错误信息

### 端口被占用

如果5683端口被占用，可以修改config.yaml中的端口号，并相应修改测试客户端中的端口。
