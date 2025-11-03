# AGVC 插件 - AGC/AVC能量管理系统

## 功能概述

AGVC插件是一个完整的光伏电站能量管理系统（EMS），实现了AGC（有功控制）和AVC（无功电压控制）功能。

### 核心功能

1. **CoAP数据采集**
   - 接收来自设备的CoAP数据
   - 每5分钟自动保存到InfluxDB
   - 实时数据存储在内存Map中（线程安全）

2. **设备管理**
   - 设备注册和配置
   - 逆变器管理
   - 并网点管理
   - 测点映射管理

3. **实时数据监控**
   - 设备实时数据查询
   - 按PSID、EQID、EQType、DataType、Point组织
   - 实时数据与测点映射关联展示

4. **历史数据查询**
   - InfluxDB历史数据查询
   - 支持时间范围、聚合间隔
   - 支持多测点并行查询

5. **AGC有功控制**
   - 自动有功功率调节
   - 支持调度/站内控制模式
   - 支持开环/闭环运行
   - 抖动区间判断
   - 闭锁信号保护
   - 调节量平均分配到逆变器
   - 完整的调节记录

6. **AVC无功电压控制**
   - 自动电压调节
   - 无功功率优化
   - 电压死区判断
   - 设备筛选和调节量分配
   - 完整的调节记录

## 数据结构

### 设备编号规则

设备编号格式：`PSID(3位) + EQID(4位) + DataType(2位) + EQType(2位)`

- **PSID**: 电站ID（3位数字）
- **EQID**: 设备ID（4位数字）
- **EQType**: 设备类型
  - `01`: 逆变器
  - `02`: 并网点
- **DataType**: 数据类型
  - `01`: 遥信
  - `02`: 遥测
  - `03`: 遥控
  - `04`: 遥调

### CoAP数据格式

```json
[{
  "psid": "001",
  "eqid": "0001",
  "eqType": "01",
  "dataType": "02",
  "point": "401",
  "value": 32.32
}]
```

### 标准测点

#### AGC标准点

| 类型 | 点名称 | 点标识 |
|------|--------|--------|
| 遥信 | AGC投退信号 | 401 |
| 遥信 | AGC就地远方控制模式 | 402 |
| 遥信 | AGC开/闭环状态 | 404 |
| 遥信 | AGC有功上调节闭锁 | 405 |
| 遥信 | AGC有功下调节闭锁 | 406 |
| 遥测 | 有功调节上限 | 401 |
| 遥测 | 有功调节下限 | 402 |
| 遥测 | 有功执行值 | 403 |
| 遥调 | 有功执行 | 401 |

#### AVC标准点

| 类型 | 点名称 | 点标识 |
|------|--------|--------|
| 遥信 | AVC功能投退信号 | 401 |
| 遥信 | AVC功能就地远方控制模式 | 402 |
| 遥信 | AVC功能当前指令状态 | 403 |
| 遥信 | AVC功能开闭环状态 | 404 |
| 遥信 | AVC功能上调节闭锁 | 405 |
| 遥信 | AVC功能下调节闭锁 | 406 |
| 遥测 | 无功可增容量 | 401 |
| 遥测 | 无功可减容量 | 402 |
| 遥测 | 电压执行值 | 403 |
| 遥测 | 无功执行值 | 404 |
| 遥调 | 电压执行 | 401 |
| 遥调 | 无功执行 | 402 |

#### 逆变器标准点

| 类型 | 点名称 | 点标识 |
|------|--------|--------|
| 遥控 | 开关机 | 401 |
| 遥调 | 有功功率降额执行值 | 401 |
| 遥调 | 无功功率补偿执行值 | 402 |

## API接口

### 设备管理

- `POST /agvc/device/create` - 创建设备
- `DELETE /agvc/device/delete` - 删除设备
- `PUT /agvc/device/update` - 更新设备
- `GET /agvc/device/find` - 获取设备详情
- `GET /agvc/device/list` - 获取设备列表
- `GET /agvc/device/realtimeData` - 获取设备实时数据
- `GET /agvc/device/inverters` - 获取逆变器列表

### 历史数据

- `POST /agvc/history/query` - 查询历史数据

### AGC控制

- `GET /agvc/agc/config` - 获取AGC配置
- `POST /agvc/agc/config` - 创建AGC配置
- `PUT /agvc/agc/config` - 更新AGC配置
- `POST /agvc/agc/start` - 启动AGC控制
- `POST /agvc/agc/stop` - 停止AGC控制
- `GET /agvc/agc/records` - 获取AGC调节记录

### AVC控制

- `GET /agvc/avc/config` - 获取AVC配置
- `POST /agvc/avc/config` - 创建AVC配置
- `PUT /agvc/avc/config` - 更新AVC配置
- `POST /agvc/avc/start` - 启动AVC控制
- `POST /agvc/avc/stop` - 停止AVC控制
- `GET /agvc/avc/records` - 获取AVC调节记录

## 使用说明

### 1. 初始化设备

首先创建电站设备和逆变器设备：

```json
POST /agvc/device/create
{
  "psid": "001",
  "eqid": "0001",
  "eqType": "01",
  "dataType": "02",
  "name": "1号逆变器",
  "maxPower": 2.5,
  "maxReg": 0.5,
  "maxReact": 1.0
}
```

### 2. 配置AGC

创建AGC配置：

```json
POST /agvc/agc/config
{
  "psid": "001",
  "isActive": 1,
  "controlAuth": 1,
  "runMode": 2,
  "jitterRange": 0.5,
  "regPeriod": 30,
  "regStep": 0.1,
  "powerUpperLimit": 10.0,
  "powerLowerLimit": 0.0
}
```

### 3. 启动AGC控制

```
POST /agvc/agc/start?psid=001
```

### 4. 查看调节记录

```
GET /agvc/agc/records?psid=001&page=1&pageSize=10
```

## 技术特点

1. **并发安全**: 使用sync.RWMutex保护内存数据
2. **定时保存**: 每5分钟自动保存到InfluxDB
3. **闭环控制**: 实时监测并自动调节
4. **故障保护**: 闭锁信号、设备状态检查
5. **详细记录**: 完整的调节历史记录
6. **灵活配置**: 支持多种运行模式

## 注意事项

1. 确保InfluxDB已正确配置和初始化
2. 确保CoAP服务器地址和端口正确
3. 设备编号必须遵循规定格式
4. AGC/AVC启动前需要先创建配置
5. 调节周期不宜过短，建议≥30秒

## 依赖

- InfluxDB v2.x
- CoAP (github.com/plgd-dev/go-coap/v3)
- GORM (数据持久化)
- Gin (HTTP路由)
