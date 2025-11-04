# AGC/AVC 主流程端口和数据流修改说明

## 修改概述

本次修改实现了AGC和AVC主流程的端口变更和数据回传功能，使系统能够：
1. 从1190端口接收调度数据
2. 将接收的数据存储到内存（DispatchStorage）
3. AGC/AVC主流程进行计算
4. 通过1189端口返回计算结果给调度

## 端口变更

### 修改前
- **接收端口**: 1187（调度CoAP服务器）
- **发送端口**: 1188（默认CoAP端口）

### 修改后
- **接收端口**: 1190（调度CoAP服务器）
- **发送端口**: 1189（调度结果返回端口）
- **数据路径**: `/agvc/data`

## 修改的文件

### 1. `/server/initialize/coap_dispatch.go`
- **修改**: 将调度CoAP服务器端口从1187改为1190
- **影响**: 所有从调度接收的数据现在通过1190端口

### 2. `/server/model/agvc/agvc_main/request/common.go`
- **修改**: CoAPDataMessage结构中的Point字段从`int`改为`string`
- **原因**: 符合CoAP数据格式标准（point应为字符串）

### 3. `/server/service/agvc/agcv_main/coap_sender.go`
- **新增方法**:
  - `GetDispatchCoapHost()`: 获取调度CoAP主机地址
  - `GetDispatchCoapPort()`: 获取调度CoAP返回端口（1189）
  - `SendAGCResultToDispatch()`: 发送AGC计算结果到调度
  - `SendAVCResultToDispatch()`: 发送AVC计算结果到调度
- **修改**: 所有发送方法中的Point字段改为string类型

### 4. `/server/service/agvc/agcv_main/agc.go`
- **修改**: `executeAGCCycle()` 和 `executeOpenLoopControl()` 方法
- **新增**: `sendAGCResultToDispatch()` 方法
- **功能**: 在AGC计算完成后，自动将结果发送到1189端口

### 5. `/server/service/agvc/agcv_main/avc.go`
- **修改**: `executeAVCCycle()` 方法
- **新增**: `sendAVCResultToDispatch()` 方法
- **功能**: 在AVC计算完成后，自动将结果发送到1189端口

### 6. `/server/core/server.go`
- **修改**: 更新注释，说明新的端口配置

## 标准点位映射

### AGC标准点

#### 遥信（YX）
| 点标识 | 点名称 | 数据键 |
|--------|--------|--------|
| 401 | AGC投退信号 | agcSignal |
| 402 | AGC就地远方控制模式 | agcControlMode |
| 404 | AGC开/闭环状态 | agcLoopStatus |
| 405 | AGC有功上调节闭锁 | agcUpRegLock |
| 406 | AGC有功下调节闭锁 | agcDownRegLock |

#### 遥测（YC）
| 点标识 | 点名称 | 数据键 |
|--------|--------|--------|
| 401 | 有功调节上限 | powerUpperLimit |
| 402 | 有功调节下限 | powerLowerLimit |
| 403 | 有功执行值 | powerExecValue |

### AVC标准点

#### 遥信（YX）
| 点标识 | 点名称 | 数据键 |
|--------|--------|--------|
| 401 | AVC功能投退信号 | avcSignal |
| 402 | AVC功能就地远方控制模式 | avcControlMode |
| 403 | AVC功能当前指令状态 | avcCmdStatus |
| 404 | AVC功能开闭环状态 | avcLoopStatus |
| 405 | AVC功能上调节闭锁 | avcUpRegLock |
| 406 | AVC功能下调节闭锁 | avcDownRegLock |

#### 遥测（YC）
| 点标识 | 点名称 | 数据键 |
|--------|--------|--------|
| 401 | 无功可增容量 | reactiveIncreaseCap |
| 402 | 无功可减容量 | reactiveDecreaseCap |
| 403 | 电压执行值 | voltageExecValue |
| 404 | 无功执行值 | reactiveExecValue |

## 数据格式

### 接收数据格式（1190端口）
```json
[{
  "psid": 1,
  "eqid": 1,
  "eqType": 2,
  "dataType": 2,
  "point": "401",
  "value": 32.32
}]
```

### 发送数据格式（1189端口）
```json
[{
  "psid": 1,
  "eqid": 1,
  "eqType": 2,
  "dataType": 2,
  "point": "401",
  "value": 1.0
}]
```

## 数据流程

```
调度系统
    ↓ (1190端口, /agvc/data)
CoAP Dispatch Server (接收)
    ↓
DispatchStorage (内存存储)
    ↓
AGC/AVC 主流程计算
    ↓
生成标准点结果
    ↓ (1189端口, /agvc/data)
CoAP Sender (发送)
    ↓
调度系统
```

## 主流程计算逻辑

### AGC主流程
1. 从1190端口接收调度数据
2. 读取AGC投退信号、控制模式、开/闭环状态
3. 根据有功执行值进行计算和调节
4. 采集实际输出功率
5. 组装AGC标准点数据
6. 通过1189端口返回计算结果

### AVC主流程
1. 从1190端口接收调度数据
2. 读取AVC投退信号、控制模式、开/闭环状态
3. 根据电压执行值进行计算和调节
4. 采集实际电压和无功
5. 计算无功容量
6. 组装AVC标准点数据
7. 通过1189端口返回计算结果

## 配置说明

### 调度CoAP配置
- **主机**: 127.0.0.1 (可通过配置文件修改)
- **接收端口**: 1190
- **发送端口**: 1189
- **路径**: /agvc/data

## 兼容性说明

1. **向后兼容**: 原有的CoAP服务（5683端口）和逆变器控制（1188端口）保持不变
2. **数据格式**: Point字段统一为string类型，保持与调度系统一致
3. **内存存储**: DispatchStorage独立管理调度数据，不影响原有DataStorage

## 测试建议

1. **端口测试**: 验证1190端口能正常接收调度数据
2. **存储测试**: 验证数据能正确存储到DispatchStorage
3. **计算测试**: 验证AGC/AVC主流程能正确读取和计算
4. **发送测试**: 验证1189端口能正常发送结果数据
5. **格式测试**: 验证发送的数据格式符合标准点要求

## 注意事项

1. 确保调度系统配置正确的端口（1190用于发送，1189用于接收）
2. Point字段现在是string类型，所有点标识需使用字符串格式
3. 结果数据在每个AGC/AVC控制周期结束后自动发送
4. 如果发送失败，会记录错误日志但不影响主流程继续运行
5. 所有标准点数据均按照规范的点标识进行映射

## 日志标识

可以通过以下关键字在日志中追踪数据流：
- `调度CoAP服务器已启动`: 1190端口启动成功
- `调度数据已存储到内存`: 数据接收并存储成功
- `发送AGC计算结果到调度`: AGC结果发送
- `发送AVC计算结果到调度`: AVC结果发送
- `发送AGC结果到调度失败`: AGC结果发送失败
- `发送AVC结果到调度失败`: AVC结果发送失败
