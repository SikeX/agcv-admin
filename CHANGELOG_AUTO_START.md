# 变更日志 - AGC/AVC自动启动功能

## 版本：v1.0.0
## 日期：2024-11-04

### 新增功能

#### 1. 系统启动时自动启动所有并网点的AGC和AVC功能

**背景：**
- 原先需要通过API接口手动调用启动每个并网点的AGC和AVC功能
- 系统重启后需要重新手动启动，操作繁琐且容易遗漏

**改进：**
- 系统启动时自动检查所有并网点配置
- 根据`agc_is_enabled`和`avc_is_enabled`标志自动启动相应功能
- 减少人工干预，提高系统可用性

### 修改文件清单

#### 后端代码

1. **server/service/agvc/agcv_main/agc.go**
   - 新增方法：`AutoStartAllGridPoints()`
   - 功能：自动启动所有启用AGC的并网点
   - 特性：
     - 智能过滤（跳过未启用或配置错误的并网点）
     - 详细的启动日志
     - 错误处理（单个失败不影响其他）

2. **server/service/agvc/agcv_main/avc.go**
   - 新增方法：`AutoStartAllGridPoints()`
   - 功能：自动启动所有启用AVC的并网点
   - 特性：
     - 智能过滤（跳过未启用或配置错误的并网点）
     - 详细的启动日志
     - 错误处理（单个失败不影响其他）

3. **server/core/server.go**
   - 在`RunServer()`函数中添加自动启动逻辑
   - 使用goroutine异步执行，不阻塞主流程
   - 延迟3秒启动，确保依赖服务完全初始化

#### 文档

4. **AGC_AVC_AUTO_START.md**
   - 完整的功能说明文档
   - 使用指南和配置示例
   - 技术细节和注意事项

5. **CHANGELOG_AUTO_START.md**
   - 变更日志（本文件）

### 代码变更统计

```
server/core/server.go                          | +16 行
server/service/agvc/agcv_main/agc.go          | +73 行
server/service/agvc/agcv_main/avc.go          | +73 行
AGC_AVC_AUTO_START.md                         | +308 行
CHANGELOG_AUTO_START.md                       | +本文件
---------------------------------------------------
总计：5 个文件修改，470+ 行新增
```

### 关键实现

#### AGC自动启动逻辑

```go
// AutoStartAllGridPoints 自动启动所有并网点的AGC功能
func (s *agc) AutoStartAllGridPoints() {
    var settings []agvc.AgvcBwdSetting
    if err := global.GVA_DB.Find(&settings).Error; err != nil {
        global.GVA_LOG.Error("查询并网点配置失败", zap.Error(err))
        return
    }

    for _, setting := range settings {
        // 检查AGC是否启用
        if setting.AgcIsEnabled != nil && *setting.AgcIsEnabled == 1 {
            // 启动AGC
            s.StartAGC(bwdNo)
        }
    }
}
```

#### 系统启动集成

```go
// 在RunServer()中
go func() {
    time.Sleep(3 * time.Second)
    agvcMain.AGC.AutoStartAllGridPoints()
    agvcMain.AVC.AutoStartAllGridPoints()
}()
```

### 兼容性

- ✅ 保持原有API接口不变
- ✅ 向后兼容现有配置
- ✅ 不影响手动启动/停止功能
- ✅ 支持运行时动态调整

### 测试建议

#### 功能测试

1. **正常启动测试**
   - 配置多个并网点，部分启用AGC/AVC
   - 启动系统
   - 检查日志确认自动启动成功

2. **错误处理测试**
   - 配置错误的并网点编号
   - 未创建对应的配置记录
   - 验证系统能正常处理错误

3. **并发测试**
   - 多个并网点同时启动
   - 验证不会出现重复启动

4. **重启测试**
   - 系统重启后验证AGC/AVC自动恢复运行

#### 性能测试

1. **启动延迟测试**
   - 测量自动启动对系统启动时间的影响
   - 建议：< 5秒额外延迟

2. **资源占用测试**
   - 监控多个并网点启动后的CPU和内存占用
   - 确保不会导致资源耗尽

### 配置要求

#### 数据库表

**agvc_bwd_setting**（并网点配置表）
```sql
-- 必需字段
number              VARCHAR     -- 并网点编号（整数字符串）
agc_is_enabled      INT         -- AGC启用标志（0/1）
avc_is_enabled      INT         -- AVC启用标志（0/1）
agc_control_period  INT         -- AGC控制周期（秒）
avc_control_period  INT         -- AVC控制周期（秒）

-- 可选字段
name                VARCHAR     -- 并网点名称
agc_vibration_range FLOAT       -- AGC抖动区间
avc_adjustment_range_min FLOAT  -- AVC调节范围最小值
avc_adjustment_range_max FLOAT  -- AVC调节范围最大值
```

### 部署指南

#### 1. 更新代码
```bash
git pull origin main
cd server
go mod tidy
go build
```

#### 2. 配置数据库
```sql
-- 确保并网点配置正确
UPDATE agvc_bwd_setting SET agc_is_enabled = 1 WHERE number = '1';
UPDATE agvc_bwd_setting SET avc_is_enabled = 1 WHERE number = '1';
```

#### 3. 重启服务
```bash
systemctl restart gin-vue-admin
```

#### 4. 验证日志
```bash
tail -f /var/log/gin-vue-admin/app.log | grep -E "AGC|AVC|自动启动"
```

### 监控和告警

建议添加以下监控指标：

1. **启动成功率**
   - 监控自动启动的成功数量和失败数量
   - 失败率 > 10% 时告警

2. **启动时间**
   - 监控从系统启动到所有并网点启动完成的时间
   - 超过预期时间告警

3. **运行状态**
   - 定期检查AGC/AVC是否仍在运行
   - 异常停止时告警

### 已知限制

1. **配置变更**：修改`agc_is_enabled`或`avc_is_enabled`后需要重启服务才能生效
2. **启动延迟**：固定3秒延迟，可能不适合所有场景
3. **错误恢复**：单个并网点启动失败不会自动重试

### 后续优化计划

1. **配置热加载**：支持运行时重新加载配置，无需重启
2. **自动重试**：启动失败时自动重试机制
3. **健康检查**：定期检查AGC/AVC运行状态，自动重启失败的实例
4. **优先级控制**：支持设置并网点启动优先级
5. **状态查询API**：提供API查询各并网点的启动状态

### 反馈和支持

如遇到问题或有改进建议，请通过以下方式反馈：
- 创建Issue描述问题
- 提交Pull Request改进代码
- 联系技术支持团队

---

**版本历史：**
- v1.0.0 (2024-11-04): 初始版本，实现基本的自动启动功能
