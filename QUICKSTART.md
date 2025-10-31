# CoAP AGVC Data 接口 - 快速开始

## 🚀 5分钟快速开始

### 前置条件

1. ✅ InfluxDB 已安装并运行 (默认: localhost:8086)
2. ✅ Go 1.23+ 已安装
3. ✅ 配置文件 `server/config.yaml` 已正确设置

### 步骤1: 检查配置

确保 `server/config.yaml` 中的配置正确：

```yaml
# CoAP配置
coap:
  enable: true        # 必须为true
  host: "0.0.0.0"
  port: 5683

# InfluxDB配置
influxdb:
  host: "127.0.0.1"
  port: "8086"
  token: "your-token"    # 替换为你的token
  org: "test"            # 替换为你的组织
  bucket: "test"         # 替换为你的bucket
```

### 步骤2: 启动服务器

```bash
cd server
go run main.go
```

或编译后运行：

```bash
cd server
go build
./server
```

你应该看到日志输出：

```
[GIN-debug] Listening and serving HTTP on :8888
INFO    CoAP server started    {"address": "0.0.0.0:5683"}
INFO    InfluxDB连接成功
```

### 步骤3: 测试接口

#### 选项A: 使用Go测试客户端 (推荐)

```bash
cd server/tools/coap_test
go run main.go
```

#### 选项B: 使用Python测试脚本

```bash
cd server
python3 test_coap_client.py
```

#### 选项C: 使用coap-client命令行工具

```bash
echo '[{"psid":1,"eqid":1,"eqType":2,"dataType":2,"point":"2","value":32.32}]' | \
  coap-client -m post -t application/json coap://localhost:5683/agvc/data
```

### 步骤4: 验证数据

查看服务器日志，应该看到：

```
INFO    AGVC data written to InfluxDB   {"psid": 1, "eqid": 1, "eqType": 2, "dataType": 2, "point": "2", "value": 32.32}
INFO    AGVC data saved successfully    {"count": 1}
```

查询InfluxDB验证数据已保存：

```bash
# 使用InfluxDB CLI
influx query 'from(bucket:"test") |> range(start: -1h) |> filter(fn: (r) => r._measurement == "agvc_data")'
```

## 📊 数据格式

### 请求 (POST /agvc/data)

```json
[
  {
    "psid": 1,          // 电站ID
    "eqid": 1,          // 设备ID
    "eqType": 2,        // 设备类型
    "dataType": 2,      // 数据类型
    "point": "2",       // 数据点
    "value": 32.32      // 数值
  }
]
```

### 响应

成功:
```json
{"success": true}
```

失败:
```json
{"error": "error message"}
```

## 🛠️ 故障排查

### 问题: 连接超时

**解决方案:**
- 检查服务器是否启动: `ps aux | grep server`
- 检查端口是否监听: `netstat -tunlp | grep 5683`
- 检查防火墙: `sudo ufw status`

### 问题: 数据未写入InfluxDB

**解决方案:**
- 检查InfluxDB服务: `systemctl status influxdb`
- 验证配置正确性
- 查看服务器日志获取详细错误

### 问题: JSON解析错误

**解决方案:**
- 确保发送的是有效的JSON数组
- 检查所有必需字段都存在
- 验证数据类型正确

## 📚 更多文档

- **完整API文档**: `server/COAP_AGVC_DATA_README.md`
- **实现总结**: `IMPLEMENTATION_SUMMARY.md`
- **测试工具说明**: `server/tools/coap_test/README.md`

## 🎯 下一步

1. 修改测试脚本发送真实的设备数据
2. 集成到你的设备端
3. 设置InfluxDB数据可视化
4. 配置告警规则

## 💡 提示

- CoAP使用UDP协议，网络不稳定时可能丢包
- 建议批量发送数据以提高效率
- 定期检查InfluxDB存储空间
- 为生产环境配置认证和加密

## 🔗 相关链接

- [InfluxDB文档](https://docs.influxdata.com/)
- [CoAP协议规范](https://datatracker.ietf.org/doc/html/rfc7252)
- [gin-vue-admin文档](https://www.gin-vue-admin.com/)

---

遇到问题？请查看详细文档或提交Issue。
