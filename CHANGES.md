# 修改说明文档

## 概述
本次修改实现了以下四个主要功能：
1. 日志文件限制为20MB，超过后覆盖旧内容
2. getAgvcNbqHistory接口改为传输AgvcNbqHis结构体
3. AgvcNbqHis和AgvcNbqSetting中新增psid字段
4. InfluxDB数据聚合为分钟数据返回

## 详细修改内容

### 1. 日志文件限制为20MB

#### 修改文件：
- `server/core/internal/cutter.go`
- `server/config.yaml` (已配置 max-size: 20)

#### 修改说明：
- 修改了 `getRotatedFilename` 方法
- 当日志文件达到20MB时，删除旧文件并重新创建，而不是重命名旧文件
- 这样确保只有一个日志文件存在，超过大小后覆盖旧内容

#### 关键代码：
```go
// 文件达到最大大小，删除旧文件并重新创建
// 这样可以确保只有一个日志文件，超过大小后覆盖旧内容
os.Remove(baseFilename)
return baseFilename
```

### 2. getAgvcNbqHistory接口改为传输AgvcNbqHis结构体

#### 修改文件：
- `server/model/agvc/request/agvc_nbq_his.go` - 新增 AgvcNbqHistoryRequest 结构体
- `server/service/agvc/agvc_nbq_his.go` - 修改 GetAgvcNbqHistory 方法签名和逻辑
- `server/api/v1/agvc/agvc_nbq_his.go` - 修改 GetAgvcNbqHistory API处理逻辑
- `server/router/agvc/agvc_nbq_his.go` - 将路由从GET改为POST
- `web/src/api/agvc/agvcNbqHis.js` - 修改前端API调用
- `web/src/view/agvc/agvcNbqHis/agvcNbqHis.vue` - 修改前端调用逻辑

#### 修改说明：
- **后端API** 从接收独立参数改为接收 `AgvcNbqHistoryRequest` 结构体
- **Service层** 新增 `extractPointValues` 函数，从AgvcNbqHis结构体中提取所有带point标签的字段的value值
- **查询逻辑** 只查询结构体中值不为nil的字段对应的点位
- **HTTP方法** 从GET改为POST，因为需要传输复杂的结构体数据
- **前端调用** 构建AgvcNbqHis对象，只设置需要查询的字段

#### 关键逻辑：
```go
// 提取结构体中所有带point标签的字段的value值
pointValues := extractPointValues(nbqHis)

// 对每个点位进行查询
for _, pointValue := range pointValues {
    // 构建Flux查询语句
    // ...
}
```

### 3. AgvcNbqHis和AgvcNbqSetting中新增psid字段

#### 修改文件：
- `server/model/agvc/agvc_nbq_his.go` - 添加 Psid 字段
- `server/model/agvc/agvc_nbq_setting.go` - 添加 Psid 字段
- `server/model/agvc/request/agvc_nbq_his.go` - 添加 Psid 搜索条件
- `server/model/agvc/request/agvc_nbq_setting.go` - 添加 Psid 搜索条件
- `server/service/agvc/agvc_nbq_his.go` - 在所有InfluxDB查询中添加psid过滤
- `server/service/agvc/agvc_nbq_setting.go` - 在列表查询中添加psid过滤
- `web/src/view/agvc/agvcNbqHis/agvcNbqHis.vue` - 添加psid列和搜索框
- `web/src/view/agvc/agvcNbqSetting/agvcNbqSetting.vue` - 添加psid列和搜索框
- `web/src/view/agvc/agvcNbqSetting/agvcNbqSettingForm.vue` - 添加psid字段编辑

#### 修改说明：
- **数据模型** 在AgvcNbqHis和AgvcNbqSetting结构体中添加 `Psid *int` 字段
- **搜索条件** AgvcNbqHisSearch中添加Psid搜索参数
- **InfluxDB查询** 在所有查询语句中添加psid过滤条件（当psid不为nil时）
- **前端显示** 在表格中添加电站编号列，在搜索表单中添加电站编号搜索框

#### 关键代码：
```go
// 如果有psid过滤条件，添加psid过滤
if nbqHis.Psid != nil {
    flux += fmt.Sprintf(`
    |> filter(fn: (r) => r["psid"] == "%d")`, *nbqHis.Psid)
}
```

### 4. InfluxDB数据聚合为分钟数据

#### 修改文件：
- `server/service/agvc/agvc_nbq_his.go` - 在GetAgvcNbqHistory方法中添加聚合

#### 修改说明：
- 在Flux查询语句中添加 `aggregateWindow` 函数
- 使用 `every: 1m` 进行分钟级聚合
- 使用 `fn: mean` 计算平均值
- `createEmpty: false` 不创建空时间窗口

#### 关键代码：
```go
// 添加分钟级聚合
flux += `
    |> aggregateWindow(every: 1m, fn: mean, createEmpty: false)`
```

## 数据库迁移

需要在数据库中为以下表添加 `psid` 字段：
- `agvc_nbq_setting` 表
- 如果 `agvc_nbq_his` 是数据库表，也需要添加

SQL语句示例：
```sql
ALTER TABLE agvc_nbq_setting ADD COLUMN psid INT;
```

## API接口变更

### 变更前：
```
GET /agvcNbqHis/getAgvcNbqHistory
Query参数: eqid, point, startTime, endTime
```

### 变更后：
```
POST /agvcNbqHis/getAgvcNbqHistory
Body参数: 
{
  "agvcNbqHis": {
    "psid": 1,
    "inverterNo": 100,
    "name": "逆变器1",
    "acPower": 0,  // 设置需要查询的字段为非空值
    ...
  },
  "startTime": "2024-01-01T00:00:00Z",
  "endTime": "2024-01-02T00:00:00Z"
}
```

## 前端修改说明

### 主要变更：
1. API调用方法从GET改为POST
2. 请求参数从简单参数改为结构化对象
3. 添加psid字段的显示和搜索
4. 历史数据查询时构建AgvcNbqHis对象

### 使用方式：
用户在查看历史数据时，只需点击对应字段的"历史"按钮，系统会自动：
1. 提取该行的设备信息（psid, inverterNo, name）
2. 设置该字段的值为0（触发查询）
3. 发送POST请求获取历史数据
4. 显示分钟级聚合后的数据

## 测试建议

1. **日志测试**：
   - 启动服务，持续产生日志
   - 观察日志文件大小是否在达到20MB后被重置

2. **历史数据查询测试**：
   - 测试单个字段的历史数据查询
   - 测试多个字段的历史数据查询
   - 验证数据是否按分钟聚合

3. **psid过滤测试**：
   - 测试带psid的数据查询
   - 测试不带psid的数据查询
   - 验证过滤逻辑是否正确

4. **前端测试**：
   - 测试列表页的psid搜索
   - 测试历史数据图表显示
   - 验证时间范围选择功能

## 注意事项

1. **数据类型一致性**：前后端的psid字段都使用int类型
2. **向后兼容**：如果现有数据没有psid，查询时会自动忽略psid过滤
3. **性能考虑**：分钟级聚合可以减少返回的数据量，提高前端渲染性能
4. **日志覆盖**：旧日志会被覆盖，如需长期保存请配置日志备份策略
