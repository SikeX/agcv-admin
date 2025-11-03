# AGVC 插件快速开始指南

## 5分钟快速体验

### 前置条件

- ✅ gin-vue-admin主系统已安装并运行
- ✅ MySQL数据库已配置
- ✅ InfluxDB已安装并初始化
- ✅ 已获取JWT Token（通过登录系统）

### 步骤1：验证插件安装

启动GVA服务器后，检查日志：

```bash
# 查看启动日志
tail -f server/log/server.log | grep AGVC

# 应该看到类似输出：
# [INFO] AGVC插件数据库表初始化成功
# [INFO] 测点映射初始化成功
# [INFO] 数据存储服务初始化成功
# [INFO] AGC控制服务初始化成功
# [INFO] AVC控制服务初始化成功
```

### 步骤2：创建测试设备

```bash
# 设置变量
TOKEN="YOUR_JWT_TOKEN"
BASE_URL="http://localhost:8888"

# 创建并网点
curl -X POST "$BASE_URL/api/agvc/device/create" \
  -H "x-token: $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "psid": "001",
    "eqid": "0000",
    "eqType": "02",
    "dataType": "02",
    "name": "测试电站并网点",
    "status": 1,
    "maxPower": 10.0,
    "maxReg": 2.0,
    "maxReact": 5.0
  }'

# 创建逆变器1
curl -X POST "$BASE_URL/api/agvc/device/create" \
  -H "x-token: $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "psid": "001",
    "eqid": "0001",
    "eqType": "01",
    "dataType": "02",
    "name": "测试逆变器1",
    "status": 1,
    "maxPower": 2.5,
    "maxReg": 0.5,
    "maxReact": 1.0
  }'

# 创建逆变器2
curl -X POST "$BASE_URL/api/agvc/device/create" \
  -H "x-token: $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "psid": "001",
    "eqid": "0002",
    "eqType": "01",
    "dataType": "02",
    "name": "测试逆变器2",
    "status": 1,
    "maxPower": 2.5,
    "maxReg": 0.5,
    "maxReact": 1.0
  }'
```

### 步骤3：查看设备列表

```bash
curl -X GET "$BASE_URL/api/agvc/device/list?psid=001&page=1&pageSize=10" \
  -H "x-token: $TOKEN"
```

### 步骤4：配置AGC

```bash
curl -X POST "$BASE_URL/api/agvc/agc/config" \
  -H "x-token: $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "psid": "001",
    "isActive": 0,
    "controlAuth": 1,
    "runMode": 2,
    "jitterRange": 0.3,
    "regPeriod": 30,
    "regStep": 0.1,
    "powerUpperLimit": 10.0,
    "powerLowerLimit": 0.0,
    "dispatchExecValue": 5.0,
    "stationExecValue": 5.0
  }'
```

### 步骤5：模拟设备数据上报（可选）

如果没有真实设备，可以手动插入测试数据到内存：

```bash
# 使用CoAP客户端模拟（需要安装coap-client工具）
# Ubuntu: sudo apt-get install libcoap2-bin
# 或者使用自定义脚本通过HTTP API插入
```

### 步骤6：启动AGC控制

```bash
# 先投入AGC
curl -X PUT "$BASE_URL/api/agvc/agc/config" \
  -H "x-token: $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "psid": "001",
    "isActive": 1,
    "dispatchExecValue": 6.0
  }'

# 启动AGC控制循环
curl -X POST "$BASE_URL/api/agvc/agc/start?psid=001" \
  -H "x-token: $TOKEN"

# 响应: {"code":0,"msg":"AGC已启动"}
```

### 步骤7：查看AGC运行状态

```bash
# 查看AGC配置
curl -X GET "$BASE_URL/api/agvc/agc/config?psid=001" \
  -H "x-token: $TOKEN"

# 查看调节记录
curl -X GET "$BASE_URL/api/agvc/agc/records?psid=001&page=1&pageSize=10" \
  -H "x-token: $TOKEN"
```

### 步骤8：停止AGC（测试完成后）

```bash
curl -X POST "$BASE_URL/api/agvc/agc/stop?psid=001" \
  -H "x-token: $TOKEN"
```

## 完整测试脚本

创建一个测试脚本 `test_agvc.sh`：

```bash
#!/bin/bash

# 配置
TOKEN="YOUR_JWT_TOKEN"
BASE_URL="http://localhost:8888"
PSID="001"

echo "=== AGVC 快速测试脚本 ==="
echo ""

# 1. 创建设备
echo "1. 创建测试设备..."
curl -s -X POST "$BASE_URL/api/agvc/device/create" \
  -H "x-token: $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "psid": "'$PSID'",
    "eqid": "0000",
    "eqType": "02",
    "dataType": "02",
    "name": "测试并网点",
    "status": 1,
    "maxPower": 10.0
  }' | jq .
echo ""

curl -s -X POST "$BASE_URL/api/agvc/device/create" \
  -H "x-token: $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "psid": "'$PSID'",
    "eqid": "0001",
    "eqType": "01",
    "dataType": "02",
    "name": "测试逆变器",
    "status": 1,
    "maxPower": 2.5,
    "maxReg": 0.5
  }' | jq .
echo ""

# 2. 查看设备列表
echo "2. 查看设备列表..."
curl -s -X GET "$BASE_URL/api/agvc/device/list?psid=$PSID" \
  -H "x-token: $TOKEN" | jq .
echo ""

# 3. 创建AGC配置
echo "3. 创建AGC配置..."
curl -s -X POST "$BASE_URL/api/agvc/agc/config" \
  -H "x-token: $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "psid": "'$PSID'",
    "isActive": 0,
    "controlAuth": 1,
    "runMode": 2,
    "jitterRange": 0.3,
    "regPeriod": 30,
    "dispatchExecValue": 5.0
  }' | jq .
echo ""

# 4. 投入AGC
echo "4. 投入AGC..."
curl -s -X PUT "$BASE_URL/api/agvc/agc/config" \
  -H "x-token: $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "psid": "'$PSID'",
    "isActive": 1
  }' | jq .
echo ""

# 5. 启动AGC
echo "5. 启动AGC控制..."
curl -s -X POST "$BASE_URL/api/agvc/agc/start?psid=$PSID" \
  -H "x-token: $TOKEN" | jq .
echo ""

# 6. 等待一段时间
echo "6. 等待30秒观察运行..."
sleep 30

# 7. 查看记录
echo "7. 查看调节记录..."
curl -s -X GET "$BASE_URL/api/agvc/agc/records?psid=$PSID" \
  -H "x-token: $TOKEN" | jq .
echo ""

# 8. 停止AGC
echo "8. 停止AGC控制..."
curl -s -X POST "$BASE_URL/api/agvc/agc/stop?psid=$PSID" \
  -H "x-token: $TOKEN" | jq .
echo ""

echo "=== 测试完成 ==="
```

使用脚本：

```bash
chmod +x test_agvc.sh
./test_agvc.sh
```

## 常见问题快速解决

### Q1: 提示"数据库未初始化"

```bash
# 检查配置文件
cat server/config.yaml | grep -A 10 "mysql"

# 确认数据库连接
mysql -h localhost -u root -p -e "SHOW DATABASES;"
```

### Q2: 提示"InfluxDB客户端未初始化"

```bash
# 检查InfluxDB服务
curl http://localhost:8086/health

# 检查配置
cat server/config.yaml | grep -A 5 "influxdb"
```

### Q3: 找不到设备

```bash
# 确认设备已创建
curl -X GET "$BASE_URL/api/agvc/device/list" -H "x-token: $TOKEN"

# 检查设备状态
mysql -u root -p -e "SELECT * FROM gva.agvc_device WHERE psid='001'"
```

### Q4: AGC无法启动

```bash
# 检查AGC配置是否存在
curl -X GET "$BASE_URL/api/agvc/agc/config?psid=001" -H "x-token: $TOKEN"

# 检查日志
tail -f server/log/server.log | grep AGC

# 确认有可用的逆变器
curl -X GET "$BASE_URL/api/agvc/device/inverters?psid=001" -H "x-token: $TOKEN"
```

## 下一步

完成快速开始后，建议阅读以下文档：

1. **README.md** - 了解完整功能
2. **API_EXAMPLES.md** - 学习所有API的使用
3. **ARCHITECTURE.md** - 理解系统架构
4. **DEPLOYMENT.md** - 生产环境部署

## 获取帮助

如遇到问题：

1. 查看 `server/log/server.log` 日志
2. 检查数据库表是否正确创建
3. 确认InfluxDB连接正常
4. 查阅项目文档
5. 提交Issue到项目仓库

## 视频教程

TODO: 添加视频教程链接

## 在线演示

TODO: 添加在线演示地址

---

**祝您使用愉快！**
