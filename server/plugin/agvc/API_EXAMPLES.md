# AGVC API 使用示例

## 基础URL

所有API的基础URL为：`/api/agvc`

## 认证

所有私有API需要在请求头中携带JWT Token：

```
x-token: YOUR_JWT_TOKEN
```

## 完整工作流程示例

### 1. 创建电站并网点设备

```bash
curl -X POST "http://localhost:8888/api/agvc/device/create" \
  -H "x-token: YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "psid": "001",
    "eqid": "0000",
    "eqType": "02",
    "dataType": "02",
    "name": "1号电站并网点",
    "status": 1,
    "maxPower": 20.0,
    "maxReg": 5.0,
    "maxReact": 10.0,
    "remark": "主并网点"
  }'
```

### 2. 创建逆变器设备

```bash
# 创建逆变器1
curl -X POST "http://localhost:8888/api/agvc/device/create" \
  -H "x-token: YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "psid": "001",
    "eqid": "0001",
    "eqType": "01",
    "dataType": "02",
    "name": "1号逆变器",
    "status": 1,
    "maxPower": 2.5,
    "maxReg": 0.5,
    "maxReact": 1.0
  }'

# 创建逆变器2
curl -X POST "http://localhost:8888/api/agvc/device/create" \
  -H "x-token: YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "psid": "001",
    "eqid": "0002",
    "eqType": "01",
    "dataType": "02",
    "name": "2号逆变器",
    "status": 1,
    "maxPower": 2.5,
    "maxReg": 0.5,
    "maxReact": 1.0
  }'
```

### 3. 获取设备列表

```bash
curl -X GET "http://localhost:8888/api/agvc/device/list?psid=001&page=1&pageSize=10" \
  -H "x-token: YOUR_TOKEN"
```

### 4. 获取逆变器列表

```bash
curl -X GET "http://localhost:8888/api/agvc/device/inverters?psid=001" \
  -H "x-token: YOUR_TOKEN"
```

### 5. 创建AGC配置

```bash
curl -X POST "http://localhost:8888/api/agvc/agc/config" \
  -H "x-token: YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "psid": "001",
    "isActive": 0,
    "controlAuth": 1,
    "runMode": 2,
    "jitterRange": 0.5,
    "regPeriod": 30,
    "regStep": 0.1,
    "powerUpperLimit": 20.0,
    "powerLowerLimit": 0.0,
    "powerExecValue": 0.0,
    "upRegLock": 0,
    "downRegLock": 0,
    "dispatchExecValue": 10.0,
    "stationExecValue": 8.0
  }'
```

### 6. 更新AGC配置（投入AGC）

```bash
curl -X PUT "http://localhost:8888/api/agvc/agc/config" \
  -H "x-token: YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "psid": "001",
    "isActive": 1,
    "controlAuth": 1,
    "runMode": 2,
    "dispatchExecValue": 15.0
  }'
```

### 7. 启动AGC控制

```bash
curl -X POST "http://localhost:8888/api/agvc/agc/start?psid=001" \
  -H "x-token: YOUR_TOKEN"
```

### 8. 获取AGC配置

```bash
curl -X GET "http://localhost:8888/api/agvc/agc/config?psid=001" \
  -H "x-token: YOUR_TOKEN"
```

### 9. 获取AGC调节记录

```bash
curl -X GET "http://localhost:8888/api/agvc/agc/records?psid=001&page=1&pageSize=10" \
  -H "x-token: YOUR_TOKEN"
```

### 10. 停止AGC控制

```bash
curl -X POST "http://localhost:8888/api/agvc/agc/stop?psid=001" \
  -H "x-token: YOUR_TOKEN"
```

### 11. 创建AVC配置

```bash
curl -X POST "http://localhost:8888/api/agvc/avc/config" \
  -H "x-token: YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "psid": "001",
    "isActive": 0,
    "controlAuth": 1,
    "commandStatus": 0,
    "runMode": 2,
    "targetVoltageLow": 10.0,
    "targetVoltageHigh": 10.5,
    "voltageExecValue": 10.2,
    "reactiveExecValue": 0.0,
    "reactiveIncCap": 10.0,
    "reactiveDecCap": 10.0,
    "upRegLock": 0,
    "downRegLock": 0,
    "regPeriod": 60,
    "voltageDeadZone": 0.2,
    "reactiveSensitivity": 1.0
  }'
```

### 12. 启动AVC控制

```bash
curl -X POST "http://localhost:8888/api/agvc/avc/start?psid=001" \
  -H "x-token: YOUR_TOKEN"
```

### 13. 获取设备实时数据

```bash
curl -X GET "http://localhost:8888/api/agvc/device/realtimeData?psid=001&eqid=0001&eqType=01" \
  -H "x-token: YOUR_TOKEN"
```

### 14. 查询历史数据

```bash
curl -X POST "http://localhost:8888/api/agvc/history/query" \
  -H "x-token: YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "psid": "001",
    "eqid": "0001",
    "eqType": "01",
    "dataType": "02",
    "points": ["401", "402"],
    "startTime": 1704067200,
    "endTime": 1704153600,
    "interval": "5m"
  }'
```

## CoAP数据上报示例

设备通过CoAP协议向服务器上报数据：

```bash
# 使用coap-client工具（需要安装libcoap）
coap-client -m post -t json \
  -e '[{
    "psid":"001",
    "eqid":"0001",
    "eqType":"01",
    "dataType":"02",
    "point":"401",
    "value":2.35
  },{
    "psid":"001",
    "eqid":"0001",
    "eqType":"01",
    "dataType":"02",
    "point":"402",
    "value":0.15
  }]' \
  coap://localhost:1188/agvc/data
```

## 响应格式

### 成功响应

```json
{
  "code": 0,
  "data": {...},
  "msg": "操作成功"
}
```

### 错误响应

```json
{
  "code": 7,
  "data": {},
  "msg": "错误信息"
}
```

### 分页响应

```json
{
  "code": 0,
  "data": {
    "list": [...],
    "total": 100,
    "page": 1,
    "pageSize": 10
  },
  "msg": "获取成功"
}
```

## 数据模型示例

### 设备信息

```json
{
  "id": 1,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z",
  "psid": "001",
  "eqid": "0001",
  "eqType": "01",
  "dataType": "02",
  "name": "1号逆变器",
  "deviceCode": "00100010201",
  "status": 1,
  "maxPower": 2.5,
  "maxReg": 0.5,
  "maxReact": 1.0,
  "remark": ""
}
```

### AGC配置

```json
{
  "id": 1,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z",
  "psid": "001",
  "isActive": 1,
  "controlAuth": 1,
  "runMode": 2,
  "jitterRange": 0.5,
  "regPeriod": 30,
  "regStep": 0.1,
  "powerUpperLimit": 20.0,
  "powerLowerLimit": 0.0,
  "powerExecValue": 15.0,
  "upRegLock": 0,
  "downRegLock": 0,
  "dispatchExecValue": 15.0,
  "stationExecValue": 15.0
}
```

### AGC调节记录

```json
{
  "id": 1,
  "createdAt": "2024-01-01T12:00:00Z",
  "updatedAt": "2024-01-01T12:00:30Z",
  "psid": "001",
  "targetPower": 15.0,
  "actualPower": 14.2,
  "powerDeviation": 0.8,
  "regulationPower": 0.8,
  "systemFreq": 50.02,
  "status": "success",
  "message": "调节成功，2个逆变器全部响应"
}
```

### 设备实时数据

```json
{
  "device": {
    "id": 1,
    "psid": "001",
    "eqid": "0001",
    "eqType": "01",
    "name": "1号逆变器",
    "status": 1
  },
  "realtimeData": {
    "02": {
      "401": {
        "value": 2.35,
        "pointName": "有功功率",
        "unit": "MW",
        "description": "逆变器有功功率",
        "timestamp": 1704153600
      },
      "402": {
        "value": 0.15,
        "pointName": "无功功率",
        "unit": "MVar",
        "description": "逆变器无功功率",
        "timestamp": 1704153600
      }
    }
  }
}
```

## 错误码

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 7 | 参数错误或业务逻辑错误 |
| 401 | 未认证 |
| 403 | 无权限 |
| 500 | 服务器内部错误 |

## 注意事项

1. 所有时间戳使用Unix时间戳（秒）
2. 设备编号必须严格遵循格式要求
3. AGC/AVC配置中的指针类型字段（如isActive）在更新时传null表示不修改
4. 启动AGC/AVC前必须确保配置已创建
5. CoAP数据上报的JSON数组中可以包含多个数据点
6. 历史数据查询的时间范围不宜过大，建议不超过7天
