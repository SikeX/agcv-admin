# AGVC调度控制功能实现总结

## 实现概述

本次开发实现了完整的AGVC调度控制系统，包括：
1. ✅ CoAP 1187端口接收调度命令
2. ✅ 调度数据内存存储（DispatchStorage）
3. ✅ 设备点位映射（从Excel自动加载）
4. ✅ AGC/AVC控制逻辑修改（支持调度控制）
5. ✅ 并网柜功率聚合（从逆变器求和）
6. ✅ AGC/AVC计划曲线功能
7. ✅ CoAP发送控制值回采集端

## 文件清单

### 新增文件

#### 模型层
- `server/model/agvc/agvc_schedule.go` - 计划曲线数据模型

#### 服务层
- `server/service/agvc/agcv_main/point_mapper.go` - 设备点位映射服务
- `server/service/agvc/agcv_main/dispatch_storage.go` - 调度数据存储服务
- `server/service/agvc/agcv_main/schedule_service.go` - 计划曲线服务
- `server/service/agvc/agcv_main/power_aggregator.go` - 功率聚合服务

#### 初始化层
- `server/initialize/coap_dispatch.go` - 调度CoAP服务器（1187端口）

#### API层
- `server/api/v1/agvc_main/schedule.go` - 计划曲线API

#### 路由层
- `server/router/agvc_main/enter.go` - 路由组入口
- `server/router/agvc_main/schedule.go` - 计划曲线路由

#### 文档
- `server/AGVC_DISPATCH_CONTROL.md` - 功能详细说明文档
- `server/test_dispatch.sh` - 测试脚本

### 修改文件

#### 核心初始化
- `server/core/server.go` 
  - 添加调度CoAP服务器初始化
  - 添加DispatchStorage初始化
  - 添加PointMapper初始化
  - 添加ScheduleService初始化

#### 全局变量
- `server/global/global.go`
  - 添加GVA_COAP_DISPATCH_SERVER变量

#### AGC服务
- `server/service/agvc/agcv_main/agc.go`
  - 修改executeAGCCycle：支持从DispatchStorage获取调度值
  - 修改collectAGCData：使用PowerAggregator获取并网柜功率

#### AVC服务
- `server/service/agvc/agcv_main/avc.go`
  - 修改collectAVCData：使用PowerAggregator获取并网柜功率

#### 路由注册
- `server/initialize/router_biz.go`
  - 注册计划曲线路由

#### 数据库迁移
- `server/initialize/gorm_biz.go`
  - 添加AgvcScheduleCurve表迁移

#### API入口
- `server/api/v1/enter.go`
  - 注册AgvcMainApiGroup
- `server/api/v1/agvc_main/enter.go`
  - 添加Schedule API

#### 路由入口
- `server/router/enter.go`
  - 注册AgvcMain路由组

## 核心功能说明

### 1. 调度数据接收（1187端口）

**端口说明：**
- 5683端口（原有）：接收采集数据 → DataStorage
- 1187端口（新增）：接收调度命令 → DispatchStorage

**数据流：**
```
调度系统 --[POST /agvc/data]--> CoAP Server(1187) --> DispatchStorage(内存)
```

**代码位置：**
- `server/initialize/coap_dispatch.go::handleDispatchData`

### 2. 设备点位映射

**功能：**
从Excel文件读取设备类型→测点名称→点标识的映射关系

**Excel文件：**
- 位置：`server/设备及测点标准.xlsx`
- 格式：每个sheet对应一种设备类型（逆变器、并网柜、气象仪等）

**使用示例：**
```go
// 获取并网柜有功功率的点标识
pointID, _ := PointMapper.GetPointID(5, "有功功率P(kW)")
// 返回: "7"
```

**代码位置：**
- `server/service/agvc/agcv_main/point_mapper.go`

### 3. 调度控制逻辑

**AGC控制：**
```go
if config.ControlAuth == 1 {  // 调度控制
    // 优先从DispatchStorage获取
    dispatchVal, err := DispatchStorage.GetDataAsFloat64(...)
    if err == nil {
        execVal = dispatchVal  // 使用调度值
    } else {
        execVal = config.DispatchExecValue  // 使用配置值
    }
}
```

**代码位置：**
- `server/service/agvc/agcv_main/agc.go::executeAGCCycle`

### 4. 功率聚合

**逻辑：**
1. 先尝试从DataStorage获取并网柜功率
2. 如果没有数据，则从该并网柜下的所有逆变器聚合

**支持的功率类型：**
- 有功功率（P）
- 无功功率（Q）
- 视在功率（S）

**使用示例：**
```go
// 获取并网柜1的有功功率（自动聚合）
power, err := PowerAggregator.GetBwgActivePower(1)
```

**代码位置：**
- `server/service/agvc/agcv_main/power_aggregator.go`

### 5. 计划曲线

**功能：**
- 支持本地曲线和调度曲线
- 支持AGC和AVC
- 定时检查（每分钟）
- 自动执行到期计划

**数据表：**
```sql
CREATE TABLE agvc_schedule_curve (
    id BIGINT PRIMARY KEY,
    bwd_no INT,          -- 并网点编号
    type INT,            -- 1:AGC 2:AVC
    source INT,          -- 1:本地 2:调度
    start_time VARCHAR,  -- HH:MM格式
    target_value FLOAT,  -- 目标值
    enabled INT,         -- 是否启用
    executed INT,        -- 今日是否已执行
    last_exec_at BIGINT  -- 最后执行时间戳
)
```

**执行逻辑：**
- AGC本地曲线 → 更新station_exec_value
- AGC调度曲线 → 更新dispatch_exec_value
- AVC本地曲线 → 更新station_exec_value
- AVC调度曲线 → 更新dispatch_exec_value

**代码位置：**
- `server/service/agvc/agcv_main/schedule_service.go`
- `server/api/v1/agvc_main/schedule.go`

## API接口

### 创建计划曲线
```http
POST /agvcMain/schedule/create
Content-Type: application/json

{
  "bwdNo": 1,
  "type": 1,
  "source": 1,
  "startTime": "10:00",
  "targetValue": 1000,
  "enabled": 1
}
```

### 更新计划曲线
```http
PUT /agvcMain/schedule/update
Content-Type: application/json

{
  "id": 1,
  "targetValue": 1500,
  "enabled": 1
}
```

### 删除计划曲线
```http
DELETE /agvcMain/schedule/delete?id=1
```

### 查询计划曲线列表
```http
GET /agvcMain/schedule/list?bwdNo=1&type=1&source=1
```

## 测试方法

### 1. 启动服务
```bash
cd server
go run main.go
```

### 2. 测试调度数据接收
```bash
# 使用coap-client（需要先安装）
echo '[{"psid":1,"eqid":1,"eqType":5,"dataType":5,"point":"7","value":1000}]' | \
  coap-client -m post -t application/json coap://localhost:1187/agvc/data
```

### 3. 测试计划曲线
```bash
# 创建计划
curl -X POST http://localhost:8888/agvcMain/schedule/create \
  -H "Content-Type: application/json" \
  -H "x-token: YOUR_TOKEN" \
  -d '{
    "bwdNo": 1,
    "type": 1,
    "source": 1,
    "startTime": "10:00",
    "targetValue": 1000,
    "enabled": 1
  }'
```

### 4. 查看日志
```bash
# 实时查看相关日志
tail -f server.log | grep -E '调度|计划|聚合'
```

## 配置说明

### 并网点配置（agvc_bwd_setting表）

| 字段 | 说明 | 示例值 |
|-----|------|--------|
| control_auth | 控制权限<br>0:站内控制<br>1:调度控制 | 1 |
| run_mode | 运行模式<br>0:闭环<br>1:开环 | 0 |
| station_exec_value | 站内执行值 | 1000.0 |
| dispatch_exec_value | 调度执行值 | 1500.0 |
| agc_is_enabled | AGC是否投入 | 1 |
| agc_vibration_range | 抖动区间 | 50.0 |
| agc_control_period | 控制周期（秒） | 60 |

## 数据类型常量

```go
const (
    YX = 1  // 遥信
    YC = 2  // 遥测
    YM = 3  // 遥脉
    YK = 4  // 遥控
    YT = 5  // 遥调
)

const (
    TYPE_BWG = 5  // 并网柜
    TYPE_NBQ = 2  // 逆变器
)
```

## 注意事项

1. **psid固定为1**：本项目中psid永远是1
2. **Excel文件位置**：`server/设备及测点标准.xlsx`必须存在
3. **端口占用**：确保1187和5683端口未被占用
4. **时间格式**：计划曲线使用HH:MM格式（如"10:00"）
5. **每日重置**：计划曲线的executed字段会在00:00自动重置

## 扩展建议

### 1. 添加新设备类型映射
- 在Excel添加新sheet
- 在point_mapper.go添加映射关系
- 重启服务

### 2. 自定义聚合策略
- 修改power_aggregator.go
- 可实现加权平均、优先级选择等

### 3. 增强计划曲线
- 添加渐变执行
- 添加条件判断
- 添加冲突检测

### 4. 添加监控告警
- 调度数据接收状态
- 计划执行成功率
- 聚合失败告警

## 技术亮点

1. **双CoAP服务器**：分离采集和调度通道
2. **内存存储**：高性能实时数据访问
3. **自动映射**：Excel驱动的配置管理
4. **智能聚合**：自动处理数据缺失
5. **定时调度**：灵活的计划曲线系统
6. **分层架构**：严格遵循GVA框架规范

## 性能考虑

1. **内存使用**：DispatchStorage和DataStorage都使用内存存储，需要监控内存使用
2. **并发安全**：所有存储服务都使用sync.RWMutex保护
3. **定时任务**：计划曲线每分钟检查一次，轻量级操作
4. **Excel加载**：仅在启动时加载一次

## 维护建议

1. 定期检查Excel文件完整性
2. 监控调度数据接收情况
3. 定期清理过期的计划曲线
4. 监控内存使用情况
5. 定期备份配置和计划数据

## 相关文档

- 详细功能说明：`server/AGVC_DISPATCH_CONTROL.md`
- 测试脚本：`server/test_dispatch.sh`
- GVA框架文档：`README.md`

## 联系方式

如有问题，请查阅以上文档或联系开发团队。
