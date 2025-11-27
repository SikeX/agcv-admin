# AGC/AVC状态发送到TCP 1187端口功能实现

## 功能描述

在 `agvcBwdHis` 页面中，当用户操作以下AGC/AVC状态字段时，系统会自动将数据发送到TCP的1187端口：

### AGC状态字段
- **agcFunctionState** - AGC功能投退状态
  - 点位: 401
  - 数据类型: 1 (遥信)
  - 值: 1=投入, 0=退出

- **agcControlMode** - AGC调节方式
  - 点位: 404
  - 数据类型: 1 (遥信)
  - 值: 1=闭环指导, 0=开环指导

- **agcControlAuthority** - AGC控制权限
  - 点位: 402
  - 数据类型: 1 (遥信)
  - 值: 1=调度控制, 0=站内控制

### AVC状态字段
- **avcFunctionState** - AVC功能投退状态
  - 点位: 401
  - 数据类型: 1 (遥信)
  - 值: 1=投入, 0=退出

- **avcControlMode** - AVC调节方式
  - 点位: 404
  - 数据类型: 1 (遥信)
  - 值: 1=开环指导, 0=闭环调节

- **avcControlAuthority** - AVC控制权限
  - 点位: 402
  - 数据类型: 1 (遥信)
  - 值: 1=站内控制, 0=调度控制

## 数据格式

发送到TCP 1187端口的数据格式为JSON数组：

```json
[
  {
    "psid": 1,
    "eqid": 1,
    "eqType": 60,
    "dataType": 1,
    "point": "401",
    "value": 1
  }
]
```

### 字段说明
- **psid**: 电站编号
- **eqid**: 设备编号（并网点编号）
- **eqType**: 设备类型
  - 60 = AGC
  - 61 = AVC
- **dataType**: 数据类型
  - 1 = 遥信 (YX)
  - 2 = 遥测 (YC)
- **point**: 点位标识
- **value**: 点位值

## 实现文件

### 后端

#### 1. API层
**文件**: `server/api/v1/agvc/agvc_bwd_his.go`

新增接口：
```go
// SendAgcAvcStatesToTcp 发送AGC/AVC状态到TCP 1187端口
func (agvcBwdHisApi *AgvcBwdHisApi) SendAgcAvcStatesToTcp(c *gin.Context)
```

#### 2. 路由层
**文件**: `server/router/agvc/agvc_bwd_his.go`

新增路由：
```go
agvcBwdHisRouter.POST("sendAgcAvcStatesToTcp", agvcBwdHisApi.SendAgcAvcStatesToTcp)
```

#### 3. 服务层
**文件**: `server/service/agvc/agvc_bwd_his.go`

新增服务方法：
```go
// SendAgcAvcStatesToTcp 发送AGC/AVC状态到TCP 1187端口
func (agvcBwdHisService *AgvcBwdHisService) SendAgcAvcStatesToTcp(ctx context.Context, messages []agvcMainReq.CoAPDataMessage) error
```

使用 `agcvMain.CoapSender` 的以下方法：
- `GetDispatchCoapHost()` - 获取主机地址
- `GetDispatchBackCoapPort()` - 获取1187端口
- `SendData(host, port, messages)` - 发送数据

### 前端

#### 1. API接口
**文件**: `web/src/api/agvc/agvcBwdHis.js`

新增API函数：
```javascript
// 发送AGC/AVC状态到TCP 1187端口
export const sendAgcAvcStatesToTcp = (data) => {
  return service({
    url: '/agvcBwdHis/sendAgcAvcStatesToTcp',
    method: 'post',
    donNotShowLoading: true,
    data
  })
}
```

#### 2. 页面组件
**文件**: `web/src/view/agvc/agvcBwdHis/agvcBwdHis.vue`

**变更内容**：
1. 导入 `sendAgcAvcStatesToTcp` API函数
2. 为AGC状态字段添加watch监听器
3. 为AVC状态字段添加watch监听器
4. 当字段值变化时，自动构建数据并调用API发送

**AGC Watch监听器逻辑**：
```javascript
watch([agcFunctionState, agcControlMode, agcControlAuthority], async (newValues, oldValues) => {
  // 跳过初始化赋值
  if (!oldValues || oldValues.every((val, index) => val === newValues[index])) {
    return
  }
  
  // 构建消息数组
  const messages = [
    { psid, eqid, eqType: 60, dataType: 1, point: "401", value: agcFunctionState },
    { psid, eqid, eqType: 60, dataType: 1, point: "402", value: agcControlAuthority },
    { psid, eqid, eqType: 60, dataType: 1, point: "404", value: agcControlMode }
  ]
  
  // 发送到TCP
  await sendAgcAvcStatesToTcp(messages)
})
```

**AVC Watch监听器逻辑**：
```javascript
watch([avcFunctionState, avcControlMode, avcControlAuthority], async (newValues, oldValues) => {
  // 跳过初始化赋值
  if (!oldValues || oldValues.every((val, index) => val === newValues[index])) {
    return
  }
  
  // 构建消息数组
  const messages = [
    { psid, eqid, eqType: 61, dataType: 1, point: "401", value: avcFunctionState },
    { psid, eqid, eqType: 61, dataType: 1, point: "402", value: avcControlAuthority },
    { psid, eqid, eqType: 61, dataType: 1, point: "404", value: avcControlMode }
  ]
  
  // 发送到TCP
  await sendAgcAvcStatesToTcp(messages)
})
```

## 关键特性

1. **自动触发**: 用户在页面上点击单选按钮改变AGC/AVC状态时，自动触发发送
2. **跳过初始化**: Watch监听器会跳过页面初始化时的数据加载，只响应用户操作
3. **错误处理**: 发送失败时显示错误消息给用户
4. **日志记录**: 后端记录发送日志，便于调试和追踪
5. **批量发送**: 一次watch触发会发送该模块的所有3个状态点位

## 通信流程

```
用户操作页面
    ↓
Vue Watch监听器触发
    ↓
构建CoAPDataMessage数组
    ↓
调用前端API: sendAgcAvcStatesToTcp()
    ↓
POST请求到后端: /agvcBwdHis/sendAgcAvcStatesToTcp
    ↓
后端API层: SendAgcAvcStatesToTcp()
    ↓
服务层: SendAgcAvcStatesToTcp()
    ↓
CoAP发送器: CoapSender.SendData()
    ↓
UDP CoAP协议发送到 127.0.0.1:1187
```

## 使用说明

### 前端使用
1. 打开 `agvcBwdHis` 页面
2. 选择并网点设备
3. 在AGC控制面板或AVC控制面板中，点击任意单选按钮改变状态
4. 系统会自动发送数据到TCP 1187端口
5. 如果发送成功，控制台会输出成功日志
6. 如果发送失败，会弹出错误提示消息

### 后端日志
发送成功时的日志示例：
```
[INFO] 发送AGC/AVC状态到TCP host=127.0.0.1 port=1187 dataCount=3
[INFO] AGC/AVC状态发送成功 count=3
```

## 注意事项

1. **网络配置**: 确保1187端口可访问，CoAP服务正常运行
2. **数据类型**: AGC的eqType=60，AVC的eqType=61，必须正确设置
3. **点位映射**: 点位标识与后端点位映射保持一致
4. **初始化跳过**: Watch监听器会跳过初始化加载，避免不必要的发送
5. **错误提示**: 发送失败时会显示用户友好的错误消息

## 扩展建议

1. **重试机制**: 可以添加失败重试逻辑
2. **队列管理**: 对于频繁操作，可以考虑添加发送队列
3. **状态反馈**: 可以添加发送状态的UI反馈（loading状态）
4. **历史记录**: 可以记录发送历史便于审计
