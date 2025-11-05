# AGC/AVC 远程同步与逆变器品牌控制实现文档

## 概述

本次修改实现了以下两个主要功能：

### 1. AGC/AVC 远程同步逻辑改造

**原逻辑**：先判断是远程还是就地，如果是远程就使用调度数据。

**新逻辑**：
- 远程推送的值会修改 AgvcBwdSetting 数据库的值
- 用户也可以调用接口修改 AgvcBwdSetting 数据库的值
- 修改后需要：
  - 同步到调度（coap1189端口）
  - 更新内存中对应的数据（如果有就替换，没有就新增）

### 2. 增加逆变器品牌控制

需要针对三个品牌的逆变器进行 AGC/AVC 调控：
- 华为（HUAWEI）
- 阳光（SUNGROW）
- 固德威（GOODWE）

## 修改内容

### 一、数据模型修改

#### 1. AgvcNbqSetting（逆变器配置）增加品牌字段

**文件**: `server/model/agvc/agvc_nbq_setting.go`

```go
InverterBrand *int `json:"inverterBrand" form:"inverterBrand" gorm:"column:inverter_brand;"` //逆变器品牌(1-华为,2-阳光,3-固德威)
```

### 二、逆变器品牌点位定义

#### 1. 品牌常量定义

**文件**: `server/service/agvc/cons/def.go`

```go
// 逆变器品牌常量
const (
    INVERTER_BRAND_HUAWEI  = 1 // 华为
    INVERTER_BRAND_SUNGROW = 2 // 阳光
    INVERTER_BRAND_GOODWE  = 3 // 固德威
)

// 逆变器控制点位常量
const (
    // 遥控点位
    INV_YK_POWER_SWITCH = "开关机" // 点标识: 401
    
    // 遥调点位
    INV_YT_ACTIVE_POWER_LIMIT   = "有功功率降额执行值" // 点标识: 401
    INV_YT_REACTIVE_POWER_COMP  = "无功功率补偿执行值" // 点标识: 402
)
```

#### 2. 三大品牌逆变器控制点位映射表

| 品牌 | 开关机（遥控） | 有功功率控制（遥调） | 无功功率控制（遥调） |
|------|---------------|-------------------|-------------------|
| **华为** | 401 | 401 | 402 |
| **阳光** | 401 | 401 | 402 |
| **固德威** | 401 | 401 | 402 |

> 注：虽然点位编号相同，但不同品牌可能在具体实现上有差异，框架已预留扩展能力。

### 三、新增核心组件

#### 1. 逆变器品牌点位映射器

**文件**: `server/service/agvc/agcv_main/inverter_brand_mapper.go`

**功能**：
- 管理三大品牌逆变器的控制点位映射
- 提供品牌点位查询接口
- 支持动态扩展新品牌

**主要方法**：
```go
// Initialize 初始化品牌映射
func (m *inverterBrandMapper) Initialize()

// GetActivePowerPoint 获取有功功率控制点位
func (m *inverterBrandMapper) GetActivePowerPoint(brand int) (string, error)

// GetReactivePowerPoint 获取无功功率控制点位
func (m *inverterBrandMapper) GetReactivePowerPoint(brand int) (string, error)

// GetBrandName 获取品牌名称
func (m *inverterBrandMapper) GetBrandName(brand int) string
```

#### 2. 并网点配置缓存管理器

**文件**: `server/service/agvc/agcv_main/bwd_setting_cache.go`

**功能**：
- 管理 AgvcBwdSetting 的内存缓存
- 同步配置到调度（1189端口）
- 处理远程推送的配置更新

**主要方法**：
```go
// Initialize 初始化缓存，从数据库加载所有配置
func (c *BwdSettingCache) Initialize() error

// Get 获取并网点配置（优先从缓存获取）
func (c *BwdSettingCache) Get(bwdNo int) (agvc.AgvcBwdSetting, bool)

// Update 更新并网点配置（同时更新数据库和缓存，并同步到调度）
func (c *BwdSettingCache) Update(bwdNo int, setting agvc.AgvcBwdSetting) error

// UpdateFromRemote 从远程调度更新配置（调度推送过来的数据）
func (c *BwdSettingCache) UpdateFromRemote(bwdNo int, updates map[string]interface{}) error

// syncToDispatch 同步配置到调度（1189端口）
func (c *BwdSettingCache) syncToDispatch(bwdNo int, setting agvc.AgvcBwdSetting) error
```

### 四、AGC/AVC 控制流程优化

#### 1. AGC 控制流程改进

**文件**: `server/service/agvc/agcv_main/agc.go`

**修改点**：
- `GetAGCConfig`: 优先从缓存获取配置
- `executeRemoteOpenLoopControl`: 根据逆变器品牌选择控制点位
- `executeRemoteClosedLoopControl`: 根据逆变器品牌选择控制点位
- `executeClosedLoopControl`: 根据逆变器品牌选择控制点位

**示例代码**：
```go
// 根据逆变器品牌获取控制点位
brand := cons.INVERTER_BRAND_HUAWEI // 默认华为
if inv.InverterBrand != nil {
    brand = *inv.InverterBrand
}

activePowerPoint, err := InverterBrandMapper.GetActivePowerPoint(brand)
if err != nil {
    activePowerPoint = "401" // 使用默认点位
}

var pointID int
fmt.Sscanf(activePowerPoint, "%d", &pointID)

commands := map[int]interface{}{
    pointID: targetPower,
}
```

#### 2. AVC 控制流程改进

**文件**: `server/service/agvc/agcv_main/avc.go`

**修改点**：
- `executeRemoteReactiveControl`: 根据逆变器品牌选择控制点位
- `sendReactiveCommands`: 根据逆变器品牌选择控制点位

### 五、API 接口新增

#### 1. 从远程调度更新配置接口

**文件**: `server/api/v1/agvc/agvc_bwd_setting_sync.go`

**接口**: `POST /agvcBwdSetting/updateFromRemote`

**请求参数**：
```json
{
  "bwdNo": 1,
  "updates": {
    "dispatch_exec_value": 100.5,
    "control_auth": 1,
    "agc_is_enabled": 1,
    "run_mode": 2
  }
}
```

**功能**：
- 接收调度推送的配置更新
- 更新数据库
- 更新内存缓存
- 不再需要同步回调度（因为是调度推送的）

#### 2. 同步配置到调度接口

**接口**: `POST /agvcBwdSetting/syncToDispatch`

**请求参数**：
```json
{
  "bwdNo": 1
}
```

**功能**：
- 用户修改配置后，主动同步到调度
- 从缓存获取配置
- 发送到调度的1189端口

### 六、系统初始化流程

**文件**: `server/core/server.go`

**初始化顺序**：
1. 初始化数据存储（DataStorage、DispatchStorage）
2. 初始化设备点位映射（PointMapper）
3. **初始化逆变器品牌点位映射（InverterBrandMapper）**
4. **初始化并网点配置缓存（SettingCache）**
5. 初始化AGC服务
6. 初始化AVC服务
7. 初始化计划曲线服务
8. 自动启动所有并网点的AGC和AVC功能

## 使用场景

### 场景1：调度推送配置更新

1. 调度系统通过 CoAP 1190端口推送数据
2. 系统接收到数据后调用内部逻辑，执行 `SettingCache.UpdateFromRemote()`
3. 配置更新到数据库和内存缓存
4. AGC/AVC 控制循环会自动使用最新配置

### 场景2：用户通过界面修改配置

1. 用户在前端修改 AgvcBwdSetting 配置
2. 前端调用 `PUT /agvcBwdSetting/updateAgvcBwdSetting` 接口
3. 后端更新数据库
4. 调用 `syncToCache()` 同步到缓存和调度
5. AGC/AVC 控制循环会自动使用最新配置

### 场景3：不同品牌逆变器的AGC控制

1. AGC 控制周期执行
2. 遍历所有在线逆变器
3. 根据每个逆变器的品牌字段（InverterBrand）
4. 通过 `InverterBrandMapper.GetActivePowerPoint(brand)` 获取对应的控制点位
5. 使用正确的点位发送 CoAP 控制指令
6. 华为、阳光、固德威逆变器分别使用各自品牌的标准点位

## 数据流向图

```
┌─────────────────┐
│  调度系统(1190)  │
└────────┬────────┘
         │ 推送配置
         ▼
┌─────────────────────────┐
│  UpdateFromRemote API   │
└────────┬────────────────┘
         │
         ▼
┌─────────────────────────┐      ┌──────────────┐
│   SettingCache          │─────▶│   数据库      │
│  (内存缓存管理器)         │      └──────────────┘
└────────┬────────────────┘
         │
         │ 同步到调度
         ▼
┌─────────────────┐
│ 调度系统(1189)  │
└─────────────────┘

         │ 获取配置
         ▼
┌─────────────────────────┐
│    AGC/AVC 控制循环      │
└────────┬────────────────┘
         │
         ▼
┌─────────────────────────┐
│ InverterBrandMapper     │
│  (品牌点位映射器)         │
└────────┬────────────────┘
         │ 获取控制点位
         ▼
┌─────────────────────────┐
│   逆变器 (CoAP 1188)     │
│  华为/阳光/固德威         │
└─────────────────────────┘
```

## 技术要点

### 1. 线程安全

- `BwdSettingCache` 使用 `sync.RWMutex` 保证并发安全
- 读操作使用读锁，写操作使用写锁
- 避免竞态条件

### 2. 配置优先级

- 内存缓存 > 数据库 > 默认值
- 优先从缓存读取，提高性能
- 缓存未命中时从数据库加载并更新缓存

### 3. 品牌扩展性

- 通过映射表管理品牌点位
- 新增品牌只需添加映射即可
- 不需要修改业务逻辑代码

### 4. 错误处理

- 品牌点位获取失败时使用默认点位
- 记录警告日志但不阻塞流程
- 保证系统容错性

## 数据库变更

### 新增字段

**表**: `agvc_nbq_setting`

```sql
ALTER TABLE agvc_nbq_setting 
ADD COLUMN inverter_brand INT DEFAULT 1 COMMENT '逆变器品牌(1-华为,2-阳光,3-固德威)';
```

## 配置示例

### 逆变器配置示例

```json
{
  "id": 1,
  "inverterNo": 1001,
  "name": "1#逆变器",
  "inverterBrand": 1,
  "ratedActivePower": 500.0,
  "ratedReactivePower": 300.0,
  "bwdNo": "1"
}
```

### 远程推送配置示例

```json
{
  "bwdNo": 1,
  "updates": {
    "dispatch_exec_value": 450.5,
    "control_auth": 1,
    "agc_is_enabled": 1,
    "avc_is_enabled": 1,
    "run_mode": 2
  }
}
```

## 注意事项

1. **品牌字段默认值**: 新增逆变器时，如果不指定品牌，默认为华为（1）
2. **配置同步时机**: 调度推送的数据不会再次同步回调度，避免循环
3. **缓存初始化**: 系统启动时会自动加载所有并网点配置到缓存
4. **点位兼容性**: 当前三大品牌的点位编号相同，但框架已预留差异化扩展能力
5. **并发安全**: 所有缓存操作都是线程安全的，可以在多协程环境下使用

## 测试建议

1. **单元测试**：
   - 测试品牌映射器的各个方法
   - 测试缓存管理器的 CRUD 操作
   
2. **集成测试**：
   - 测试远程推送配置接口
   - 测试同步到调度接口
   - 测试不同品牌逆变器的控制流程

3. **压力测试**：
   - 并发更新配置
   - 高频读取缓存
   - 大量逆变器同时控制

## 未来扩展

1. **品牌特性支持**：可以为不同品牌定义不同的控制策略
2. **点位动态映射**：支持从配置文件或数据库加载点位映射
3. **配置版本控制**：记录配置变更历史
4. **配置回滚功能**：支持配置回滚到历史版本

## 总结

本次修改实现了：
- ✅ 远程调度数据推送更新配置
- ✅ 用户修改配置同步到调度
- ✅ 内存缓存机制提高性能
- ✅ 支持三大品牌逆变器的差异化控制
- ✅ 完整的 AGC/AVC 控制流程优化
- ✅ 良好的扩展性和容错性
