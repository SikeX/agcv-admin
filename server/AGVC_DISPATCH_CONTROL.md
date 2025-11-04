# AGVC调度控制功能说明

## 功能概述

本次更新实现了以下核心功能：

1. **1187端口CoAP服务器** - 接收调度命令
2. **调度数据内存存储** - 存储调度下发的控制值
3. **设备点位映射** - 从Excel自动加载设备类型到测点的映射关系
4. **AGC/AVC调度控制** - 支持从内存读取调度值进行控制
5. **并网柜功率聚合** - 当并网柜数据缺失时自动从逆变器聚合
6. **AGC/AVC计划曲线** - 支持设置本地和调度的计划曲线

## 架构设计

### 1. CoAP服务器端口

- **5683端口**（原有）：用于接收采集数据，数据存储到 `DataStorage`
- **1187端口**（新增）：用于接收调度命令，数据存储到 `DispatchStorage`

### 2. 数据流向

```
调度系统 --[CoAP 1187]--> DispatchStorage(内存) --> AGC/AVC控制逻辑
                                                           |
采集系统 --[CoAP 5683]--> DataStorage(内存) ------------> |
                                                           |
                                                           v
                                                   逆变器控制指令
                                                           |
                                                           v
                                              --[CoAP 1188]--> 采集系统
```

### 3. 核心服务

#### DispatchStorage (调度数据存储)
- 存储从1187端口接收的调度命令数据
- 提供读取接口供AGC/AVC使用

#### PointMapper (设备点位映射)
- 从Excel文件加载设备映射关系
- 提供设备类型→测点名称→点标识的查询
- Excel位置：`server/设备及测点标准.xlsx`

#### PowerAggregator (功率聚合)
- 当并网柜功率数据缺失时，自动从下属逆变器聚合
- 支持有功功率、无功功率、视在功率的聚合

#### ScheduleService (计划曲线服务)
- 定时检查并执行到期的计划曲线
- 支持本地曲线和调度曲线
- 每分钟检查一次，自动更新执行值

## 数据格式

### 接收数据格式（1187端口）

```json
[{
  "psid": 1,
  "eqid": 1,
  "eqType": 5,
  "dataType": 2,
  "point": "7",
  "value": 1000.5
}]
```

### 发送数据格式（1188端口）

```json
[{
  "psid": 1,
  "eqid": 1,
  "eqType": 2,
  "dataType": 2,
  "point": "401",
  "value": 500.0
}]
```

## 控制逻辑

### AGC控制流程

1. 判断控制权限（`control_auth`字段）
   - `0`：站内控制，使用 `station_exec_value`
   - `1`：调度控制，优先从 `DispatchStorage` 获取值，无值则使用 `dispatch_exec_value`

2. 判断运行模式（`run_mode`字段）
   - `0`：闭环运行，根据实际出力与目标值的偏差进行调节
   - `1`：开环运行，直接下发目标值

3. 采集并网柜功率
   - 优先从 `DataStorage` 获取并网柜的有功功率
   - 如无数据，则从下属逆变器聚合

4. 计算调节量并分配到逆变器

### AVC控制流程

1. 判断控制权限（同AGC）
2. 采集并网点电压和总无功
   - 电压：从并网柜采集
   - 无功：优先从并网柜获取，无数据则从逆变器聚合
3. 判断电压是否超出范围
4. 计算所需无功调节量
5. 分配到各逆变器并下发指令

## 计划曲线

### 数据模型

```go
type AgvcScheduleCurve struct {
    BwdNo       int     // 并网点编号
    Type        int     // 类型 1:AGC 2:AVC
    Source      int     // 来源 1:本地 2:调度
    StartTime   string  // 开始时间 HH:MM格式
    TargetValue float64 // 目标值
    Enabled     int     // 是否启用 0:禁用 1:启用
    Executed    int     // 今日是否已执行
    LastExecAt  *int64  // 最后执行时间戳
}
```

### 执行逻辑

- 每分钟检查一次是否有到期的计划
- 根据计划类型和来源更新对应的执行值：
  - AGC本地曲线 → 更新 `station_exec_value`
  - AGC调度曲线 → 更新 `dispatch_exec_value`
  - AVC本地曲线 → 更新 `station_exec_value`
  - AVC调度曲线 → 更新 `dispatch_exec_value`
- 每日自动重置执行状态

### API接口

#### 创建计划曲线
```
POST /agvcMain/schedule/create
```

#### 更新计划曲线
```
PUT /agvcMain/schedule/update
```

#### 删除计划曲线
```
DELETE /agvcMain/schedule/delete?id=1
```

#### 获取计划曲线列表
```
GET /agvcMain/schedule/list?bwdNo=1&type=1&source=1
```

## 设备点位映射

### Excel文件结构

Excel文件包含多个sheet，每个sheet对应一种设备类型：
- 逆变器 (eqType=2)
- 箱变 (eqType=4)
- 并网柜 (eqType=5)
- 电表 (eqType=7)
- 气象仪 (eqType=8)

每个sheet格式：
| 测点名称 | 点标识 |
|---------|--------|
| 有功功率P(kW) | 7 |
| 无功功率Q(kVar) | 8 |
| ... | ... |

### 使用示例

```go
// 获取并网柜的有功功率点标识
pointID, err := PointMapper.GetPointID(5, "有功功率P(kW)")
// 返回: "7"

// 获取逆变器的视在功率点标识
pointID, err := PointMapper.GetPointID(2, "视在功率(kVa)")
// 返回: "28"
```

## 配置说明

### 并网点配置（agvc_bwd_setting表）

| 字段 | 说明 |
|-----|------|
| control_auth | 控制权限 0:站内 1:调度 |
| run_mode | 运行模式 0:闭环 1:开环 |
| agc_is_enabled | AGC是否投入 |
| station_exec_value | 站内执行值 |
| dispatch_exec_value | 调度执行值 |
| agc_vibration_range | 抖动区间 |
| agc_control_period | 控制周期（秒）|

## 测试步骤

### 1. 测试调度数据接收

```bash
# 发送调度命令到1187端口
coap-client -m post -t json coap://localhost:1187/agvc/data \
  -e '[{"psid":1,"eqid":1,"eqType":5,"dataType":5,"point":"7","value":1000}]'
```

### 2. 测试计划曲线

```bash
# 创建一个10:00执行的AGC本地曲线
curl -X POST http://localhost:8888/agvcMain/schedule/create \
  -H "Content-Type: application/json" \
  -d '{
    "bwdNo": 1,
    "type": 1,
    "source": 1,
    "startTime": "10:00",
    "targetValue": 1000,
    "enabled": 1
  }'
```

### 3. 查看日志

```bash
# 查看调度数据接收日志
grep "调度数据已存储" server.log

# 查看计划曲线执行日志
grep "执行计划曲线" server.log

# 查看功率聚合日志
grep "从逆变器聚合功率" server.log
```

## 注意事项

1. **psid固定为1**：本项目中psid永远是1
2. **数据类型**：
   - YX=1 (遥信)
   - YC=2 (遥测)
   - YM=3 (遥脉)
   - YK=4 (遥控)
   - YT=5 (遥调)
3. **时间格式**：计划曲线的startTime使用HH:MM格式（如"10:00"）
4. **Excel加载**：系统启动时自动加载Excel，如更新Excel需重启服务

## 扩展说明

### 添加新的设备类型映射

1. 在Excel中添加新的sheet
2. 在 `point_mapper.go` 的 `sheetToEqType` map中添加映射
3. 重启服务自动加载

### 自定义功率聚合策略

在 `power_aggregator.go` 中修改 `aggregateFromInverters` 方法，可以实现：
- 加权平均
- 优先选择特定逆变器
- 排除故障设备
- 等等

### 自定义计划曲线执行逻辑

在 `schedule_service.go` 中修改 `executeSchedule` 方法，可以实现：
- 渐变执行
- 分段执行
- 条件执行
- 等等
