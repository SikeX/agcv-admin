# 修改总结

## 修改概述

本次修改完成了以下四个主要任务：

1. ✅ **日志文件限制为20MB** - 超过20MB后覆盖旧内容
2. ✅ **getAgvcNbqHistory接口重构** - 改为传输AgvcNbqHis结构体，使用point标签的value值查询
3. ✅ **添加psid字段** - 在AgvcNbqHis和AgvcNbqSetting中新增psid字段
4. ✅ **InfluxDB数据分钟级聚合** - 所有历史数据按分钟聚合返回

## 修改文件列表

### 后端修改 (13个文件)

1. **日志相关**
   - `server/core/internal/cutter.go` - 修改日志切割逻辑
   - `server/config.yaml` - 已配置max-size: 20

2. **数据模型**
   - `server/model/agvc/agvc_nbq_his.go` - 添加Psid字段
   - `server/model/agvc/agvc_nbq_setting.go` - 添加Psid字段
   - `server/model/agvc/request/agvc_nbq_his.go` - 添加AgvcNbqHistoryRequest和Psid搜索
   - `server/model/agvc/request/agvc_nbq_setting.go` - 添加Psid搜索

3. **服务层**
   - `server/service/agvc/agvc_nbq_his.go` - 重构GetAgvcNbqHistory方法，添加extractPointValues函数
   - `server/service/agvc/agvc_nbq_setting.go` - 添加Psid过滤

4. **API层**
   - `server/api/v1/agvc/agvc_nbq_his.go` - 修改GetAgvcNbqHistory接口

5. **路由层**
   - `server/router/agvc/agvc_nbq_his.go` - 将getAgvcNbqHistory改为POST方法

6. **数据库迁移**
   - `server/db_migration_add_psid.sql` - 数据库迁移SQL文件

### 前端修改 (4个文件)

1. **API接口**
   - `web/src/api/agvc/agvcNbqHis.js` - 修改getAgvcNbqHistory为POST请求

2. **历史数据页面**
   - `web/src/view/agvc/agvcNbqHis/agvcNbqHis.vue` - 添加psid列、搜索框，修改历史数据请求逻辑

3. **配置页面**
   - `web/src/view/agvc/agvcNbqSetting/agvcNbqSetting.vue` - 添加psid列和搜索框
   - `web/src/view/agvc/agvcNbqSetting/agvcNbqSettingForm.vue` - 添加psid字段编辑

## 技术细节

### 1. 日志文件管理

**修改前：**
- 日志文件达到maxSize后，重命名旧文件并创建新文件
- 会产生多个日志文件：log.log, log.20240101.log, log.20240102.log...

**修改后：**
- 日志文件达到20MB后，直接删除旧文件并重新创建
- 只保留一个log.log文件，大小不超过20MB
- 配置项：`zap.max-size: 20`

### 2. 历史数据查询接口

**修改前：**
```
GET /agvcNbqHis/getAgvcNbqHistory
参数：eqid, point, startTime, endTime
```

**修改后：**
```
POST /agvcNbqHis/getAgvcNbqHistory
请求体：
{
  "agvcNbqHis": {
    "psid": 1,
    "inverterNo": 100,
    "acPower": 0  // 设置需要查询的字段
  },
  "startTime": "2024-01-01T00:00:00Z",
  "endTime": "2024-01-02T00:00:00Z"
}
```

**查询逻辑：**
1. 接收AgvcNbqHis结构体
2. 使用反射遍历结构体字段
3. 提取所有值不为nil的字段的point标签value值
4. 对每个point值构建InfluxDB查询
5. 添加分钟级聚合：`aggregateWindow(every: 1m, fn: mean)`
6. 合并所有结果返回

### 3. psid字段集成

**数据库字段：**
```sql
ALTER TABLE `agvc_nbq_setting` ADD COLUMN `psid` INT DEFAULT NULL;
CREATE INDEX `idx_psid` ON `agvc_nbq_setting`(`psid`);
```

**Go结构体：**
```go
type AgvcNbqHis struct {
    Psid       *int    `json:"psid" form:"psid" point:"name:电站编号"`
    InverterNo *int    `json:"inverterNo" form:"inverterNo"`
    // ...
}

type AgvcNbqSetting struct {
    Psid       *int    `json:"psid" form:"psid" gorm:"column:psid;"`
    InverterNo *int    `json:"inverterNo" form:"inverterNo"`
    // ...
}
```

**InfluxDB查询：**
```go
if nbqHis.Psid != nil {
    flux += fmt.Sprintf(`|> filter(fn: (r) => r["psid"] == "%d")`, *nbqHis.Psid)
}
```

### 4. 分钟级数据聚合

**InfluxDB Flux查询：**
```flux
from(bucket: "test")
    |> range(start: startTime, stop: endTime)
    |> filter(fn: (r) => r["_measurement"] == "agvc")
    |> filter(fn: (r) => r["eqid"] == "100")
    |> filter(fn: (r) => r["point"] == "10")
    |> aggregateWindow(every: 1m, fn: mean, createEmpty: false)
```

**参数说明：**
- `every: 1m` - 1分钟为一个时间窗口
- `fn: mean` - 使用平均值聚合
- `createEmpty: false` - 不创建空数据点

## 部署步骤

### 1. 数据库迁移

```bash
# 连接到MySQL数据库
mysql -u root -p agvc

# 执行迁移脚本
source /path/to/server/db_migration_add_psid.sql
```

### 2. 后端部署

```bash
cd server
go build -o main
./main
```

### 3. 前端部署

```bash
cd web
npm install
npm run build
```

### 4. 验证部署

1. **检查日志文件大小限制**
   ```bash
   ls -lh server/log/*.log
   # 文件大小不应超过20MB
   ```

2. **测试历史数据查询接口**
   ```bash
   curl -X POST http://localhost:8888/agvcNbqHis/getAgvcNbqHistory \
     -H "Content-Type: application/json" \
     -d '{
       "agvcNbqHis": {
         "psid": 1,
         "inverterNo": 100,
         "acPower": 0
       },
       "startTime": "2024-01-01T00:00:00Z",
       "endTime": "2024-01-02T00:00:00Z"
     }'
   ```

3. **验证psid过滤**
   - 访问 http://localhost:8888/agvcNbqSetting/getAgvcNbqSettingList?psid=1
   - 检查返回结果是否只包含psid=1的数据

4. **验证分钟级聚合**
   - 查询一天的历史数据
   - 检查返回的数据点数量是否约等于1440（24小时×60分钟）

## 注意事项

### 兼容性

1. **向后兼容**
   - psid字段为可选，现有数据可以没有psid
   - 查询时如果psid为nil，不会添加psid过滤条件

2. **前端兼容**
   - 历史数据查询接口从GET改为POST
   - 前端必须同步更新，否则会调用失败

### 性能考虑

1. **日志文件IO**
   - 删除旧文件会有短暂的IO操作
   - 建议在低峰期重启服务

2. **InfluxDB查询**
   - 分钟级聚合可以减少返回的数据量
   - 对于长时间范围的查询，建议限制在30天以内

3. **数据库索引**
   - 已为psid字段添加索引
   - 如果psid查询频繁，建议定期分析表性能

### 数据一致性

1. **InfluxDB写入**
   - 确保在写入InfluxDB时包含psid标签
   - 检查 `WriteAgvcNbqDataToInfluxDB` 方法的调用

2. **数据迁移**
   - 现有AgvcNbqSetting记录的psid默认为NULL
   - 需要根据业务逻辑批量更新psid字段

## 回滚方案

如果出现问题需要回滚：

### 1. 后端回滚

```bash
# 切换到上一个版本
git checkout <previous-commit>
cd server
go build -o main
./main
```

### 2. 前端回滚

```bash
# 切换到上一个版本
git checkout <previous-commit>
cd web
npm run build
```

### 3. 数据库回滚

```sql
-- 删除psid字段
ALTER TABLE `agvc_nbq_setting` DROP COLUMN `psid`;
DROP INDEX `idx_psid` ON `agvc_nbq_setting`;
```

### 4. API调整

- 将 `/agvcNbqHis/getAgvcNbqHistory` 改回GET方法
- 恢复原来的参数格式

## 测试建议

### 单元测试

1. 测试 `extractPointValues` 函数
2. 测试psid过滤逻辑
3. 测试日志文件大小限制

### 集成测试

1. 测试历史数据查询完整流程
2. 测试psid搜索功能
3. 测试分钟级聚合数据准确性

### 压力测试

1. 大量日志写入测试
2. 大量历史数据查询测试
3. 并发查询测试

## 联系信息

如有问题，请联系开发团队。

---

修改完成日期：2024
修改人：AI Assistant
审核状态：待审核
