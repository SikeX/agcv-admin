# AGVC 系统架构文档

## 系统概述

AGVC（AGC/AVC Control）是一个完整的光伏电站能量管理系统，实现了：
- **AGC（Automatic Generation Control）**：自动发电控制/有功控制
- **AVC（Automatic Voltage Control）**：自动电压控制/无功电压控制

## 架构图

```
┌─────────────────────────────────────────────────────────────┐
│                       GVA Main System                        │
│  ┌────────────┐  ┌────────────┐  ┌──────────────┐          │
│  │   JWT      │  │  Casbin    │  │   InfluxDB   │          │
│  │   Auth     │  │   RBAC     │  │   Client     │          │
│  └────────────┘  └────────────┘  └──────────────┘          │
└─────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    AGVC Plugin                               │
│                                                               │
│  ┌──────────────────────────────────────────────────────┐  │
│  │                   API Layer                           │  │
│  │  ┌──────┐  ┌──────┐  ┌──────┐  ┌──────────┐        │  │
│  │  │Device│  │ AGC  │  │ AVC  │  │ History  │        │  │
│  │  └──────┘  └──────┘  └──────┘  └──────────┘        │  │
│  └──────────────────────────────────────────────────────┘  │
│                           │                                  │
│  ┌──────────────────────────────────────────────────────┐  │
│  │                 Service Layer                         │  │
│  │  ┌─────────────────────────────────────────────┐    │  │
│  │  │         Data Storage Service                 │    │  │
│  │  │  ┌────────────────────────────────────┐     │    │  │
│  │  │  │   Thread-Safe Memory Map           │     │    │  │
│  │  │  │   Key: psid_eqid_eqType_dataType   │     │    │  │
│  │  │  │   Value: RealtimeData              │     │    │  │
│  │  │  └────────────────────────────────────┘     │    │  │
│  │  │  ┌────────────────────────────────────┐     │    │  │
│  │  │  │   5-Minute Timer                   │     │    │  │
│  │  │  │   → Save to InfluxDB               │     │    │  │
│  │  │  └────────────────────────────────────┘     │    │  │
│  │  └─────────────────────────────────────────────┘    │  │
│  │  ┌─────────────────────────────────────────────┐    │  │
│  │  │           AGC Control Service                │    │  │
│  │  │  ┌────────────────────────────────────┐     │    │  │
│  │  │  │  Control Loop (Goroutine)          │     │    │  │
│  │  │  │  1. Check if active                │     │    │  │
│  │  │  │  2. Get exec value                 │     │    │  │
│  │  │  │  3. Collect data                   │     │    │  │
│  │  │  │  4. Calculate deviation            │     │    │  │
│  │  │  │  5. Check jitter range             │     │    │  │
│  │  │  │  6. Check lock signals             │     │    │  │
│  │  │  │  7. Allocate to inverters          │     │    │  │
│  │  │  │  8. Send CoAP commands             │     │    │  │
│  │  │  └────────────────────────────────────┘     │    │  │
│  │  └─────────────────────────────────────────────┘    │  │
│  │  ┌─────────────────────────────────────────────┐    │  │
│  │  │           AVC Control Service                │    │  │
│  │  │  ┌────────────────────────────────────┐     │    │  │
│  │  │  │  Control Loop (Goroutine)          │     │    │  │
│  │  │  │  1. Check if active                │     │    │  │
│  │  │  │  2. Get target voltage range       │     │    │  │
│  │  │  │  3. Collect voltage data           │     │    │  │
│  │  │  │  4. Check voltage in range         │     │    │  │
│  │  │  │  5. Calculate reactive needed      │     │    │  │
│  │  │  │  6. Filter available devices       │     │    │  │
│  │  │  │  7. Allocate reactive power        │     │    │  │
│  │  │  │  8. Send CoAP commands             │     │    │  │
│  │  │  └────────────────────────────────────┘     │    │  │
│  │  └─────────────────────────────────────────────┘    │  │
│  │  ┌─────────────────────────────────────────────┐    │  │
│  │  │       CoAP Sender/Receiver Service           │    │  │
│  │  └─────────────────────────────────────────────┘    │  │
│  │  ┌─────────────────────────────────────────────┐    │  │
│  │  │         Device Management Service            │    │  │
│  │  └─────────────────────────────────────────────┘    │  │
│  │  ┌─────────────────────────────────────────────┐    │  │
│  │  │         History Query Service                │    │  │
│  │  └─────────────────────────────────────────────┘    │  │
│  └──────────────────────────────────────────────────────┘  │
│                           │                                  │
│  ┌──────────────────────────────────────────────────────┐  │
│  │                   Model Layer                         │  │
│  │  • Device                                             │  │
│  │  • PointMapping                                       │  │
│  │  • AGCConfig / AVCConfig                             │  │
│  │  • AGCRegulationRecord / AVCRegulationRecord         │  │
│  │  • InverterRegulation / DeviceReactiveRegulation     │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                   External Systems                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   InfluxDB   │  │  CoAP Devices│  │  GORM/MySQL  │     │
│  │  (History)   │  │  (IoT Edge)  │  │  (Config)    │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────────────────────────────────────────────┘
```

## 数据流程

### 1. CoAP数据采集流程

```
IoT Device → CoAP Client → [POST /agvc/data] → CoAP Receiver
                                                      ↓
                                              Validate Data
                                                      ↓
                                              Memory Map (Thread-Safe)
                                                      ↓
                                              [5-Minute Timer]
                                                      ↓
                                              InfluxDB (Long-term Storage)
```

### 2. 实时数据查询流程

```
Web Client → [GET /agvc/device/realtimeData] → Device Service
                                                      ↓
                                              Memory Map Query
                                                      ↓
                                              Point Mapping Join
                                                      ↓
                                              Formatted Response
                                                      ↓
                                              Web Client
```

### 3. 历史数据查询流程

```
Web Client → [POST /agvc/history/query] → History Service
                                                ↓
                                         Build Flux Query
                                                ↓
                                         InfluxDB Query
                                                ↓
                                         Parse Results
                                                ↓
                                         Web Client
```

### 4. AGC控制流程

```
[Timer Trigger] → AGC Control Loop
                        ↓
                 1. Get AGC Config from DB
                        ↓
                 2. Check IsActive
                        ↓
                 3. Get Exec Value (Dispatch/Station)
                        ↓
                 4. Collect Data from Memory Map
                    • Target Power
                    • Actual Power
                    • System Frequency
                        ↓
                 5. Calculate Deviation
                    Deviation = Target - Actual
                        ↓
                 6. Check Jitter Range
                    |Deviation| > JitterRange ?
                        ↓
                 7. Check Lock Signals
                    • Up Regulation Lock
                    • Down Regulation Lock
                        ↓
                 8. Get Available Inverters
                    Status = Online
                        ↓
                 9. Allocate Regulation Power
                    Per Inverter = Deviation / Count
                    Limit to MaxReg
                        ↓
                10. Create Regulation Record in DB
                        ↓
                11. Send CoAP Commands to Each Inverter
                    Host: CoapHost (default: 127.0.0.1)
                    Port: CoapPort (default: 1188)
                    URL: /agvc/data
                    Data: [{psid, eqid, eqType, dataType, point, value}]
                        ↓
                12. Update Record Status
                    • success
                    • partial
                    • failed
                        ↓
                [Next Timer Trigger after RegPeriod seconds]
```

### 5. AVC控制流程

```
[Timer Trigger] → AVC Control Loop
                        ↓
                 1. Get AVC Config from DB
                        ↓
                 2. Check IsActive
                        ↓
                 3. Get Target Voltage Range
                    • TargetVoltageLow
                    • TargetVoltageHigh
                        ↓
                 4. Collect Data from Memory Map
                    • Point Voltage
                    • Total Reactive
                    • System Frequency
                        ↓
                 5. Check Voltage in Range
                    Low ≤ Voltage ≤ High ?
                        ↓
                 6. Calculate Voltage Deviation
                    If Voltage < Low: ΔV = Low - Voltage
                    If Voltage > High: ΔV = High - Voltage
                        ↓
                 7. Check Voltage Dead Zone
                    |ΔV| > VoltageDeadZone ?
                        ↓
                 8. Check Lock Signals
                        ↓
                 9. Calculate Required Reactive
                    ΔQ = ΔV × ReactiveSensitivity
                        ↓
                10. Filter Available Devices
                    • Status = Online
                    • MaxReact > 0
                        ↓
                11. Allocate Reactive Power
                    Per Device = ΔQ / Count
                    Limit to MaxReact
                        ↓
                12. Create Regulation Record in DB
                        ↓
                13. Send CoAP Commands to Each Device
                        ↓
                14. Update Record Status
                        ↓
                [Next Timer Trigger after RegPeriod seconds]
```

## 关键设计

### 1. 线程安全的实时数据存储

```go
type dataStorage struct {
    mu             sync.RWMutex
    realtimeData   map[string]*model.RealtimeData
    saveTimer      *time.Ticker
    stopChan       chan struct{}
}
```

- 使用`sync.RWMutex`保护并发访问
- 读多写少场景优化
- Key格式：`psid_eqid_eqType_dataType_point`

### 2. 独立的控制循环

每个电站的AGC和AVC都有独立的Goroutine运行控制循环：

```go
// 启动AGC
go s.agcControlLoop(psid, config, stopChan)

// 启动AVC
go s.avcControlLoop(psid, config, stopChan)
```

通过channel实现优雅停止：

```go
select {
case <-ticker.C:
    // 执行控制周期
case <-stopChan:
    // 退出循环
    return
}
```

### 3. 分层架构

严格遵循GVA的分层架构：

- **Router层**：路由定义，中间件挂载
- **API层**：HTTP请求处理，参数验证
- **Service层**：业务逻辑，数据操作
- **Model层**：数据模型定义

### 4. 插件化设计

实现GVA的Plugin接口：

```go
func (p *plugin) Register(group *gin.Engine) {
    initialize.Gorm(ctx)      // 数据库表初始化
    initialize.Service(ctx)   // 服务初始化
    initialize.Router(group)  // 路由初始化
}
```

### 5. 完整的记录体系

每次调节都会记录：

- **主记录**：AGCRegulationRecord / AVCRegulationRecord
  - 目标值、实际值、偏差
  - 调节量、系统状态
  - 执行状态、消息
  
- **详细记录**：InverterRegulation / DeviceReactiveRegulation
  - 每个设备的分配量
  - 调节前后的值
  - 执行状态

## 配置说明

### AGC配置参数

| 参数 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| IsActive | int | 是否投入 | 0 |
| ControlAuth | int | 控制权限（1:调度 2:站内） | 1 |
| RunMode | int | 运行模式（1:开环 2:闭环） | 2 |
| JitterRange | float64 | 抖动区间(MW) | 0.5 |
| RegPeriod | int | 调节周期(秒) | 30 |
| RegStep | float64 | 调节步长(MW) | 0.1 |
| PowerUpperLimit | float64 | 有功调节上限(MW) | - |
| PowerLowerLimit | float64 | 有功调节下限(MW) | - |
| UpRegLock | int | 上调节闭锁 | 0 |
| DownRegLock | int | 下调节闭锁 | 0 |

### AVC配置参数

| 参数 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| IsActive | int | 是否投入 | 0 |
| TargetVoltageLow | float64 | 目标电压下限(kV) | - |
| TargetVoltageHigh | float64 | 目标电压上限(kV) | - |
| VoltageDeadZone | float64 | 电压死区(kV) | 0.5 |
| ReactiveSensitivity | float64 | 无功灵敏度(MVar/kV) | 1.0 |
| RegPeriod | int | 调节周期(秒) | 60 |
| UpRegLock | int | 上调节闭锁 | 0 |
| DownRegLock | int | 下调节闭锁 | 0 |

## 性能考虑

1. **内存数据存储**：实时数据保存在内存中，查询速度快
2. **定时批量保存**：每5分钟批量写入InfluxDB，减少IO
3. **并发控制**：使用RWMutex，读操作不互斥
4. **异步控制循环**：使用Goroutine，不阻塞主线程
5. **索引优化**：数据库表添加必要索引

## 扩展性

1. **多电站支持**：每个电站独立的配置和控制循环
2. **可配置策略**：支持不同的调节策略（当前是平均分配）
3. **测点可扩展**：通过PointMapping表动态管理测点
4. **插件化架构**：可以独立部署和升级

## 安全性

1. **JWT认证**：所有API需要认证
2. **Casbin权限控制**：基于角色的访问控制
3. **数据验证**：严格的参数验证
4. **并发保护**：线程安全的数据访问
5. **闭锁保护**：防止误操作
