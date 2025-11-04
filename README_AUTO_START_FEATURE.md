# AGC/AVC自动启动功能 - 快速上手

## 📋 功能说明

本次更新实现了**系统启动时自动启动所有并网点的AGC和AVC功能**，无需手动调用API接口。

### 之前
```bash
# 需要手动为每个并网点调用启动接口
curl -X POST "http://localhost:8888/agvc/agc/start?bwdNo=1"
curl -X POST "http://localhost:8888/agvc/avc/start?bwdNo=1"
curl -X POST "http://localhost:8888/agvc/agc/start?bwdNo=2"
curl -X POST "http://localhost:8888/agvc/avc/start?bwdNo=2"
# ...
```

### 现在
```bash
# 系统启动后自动启动所有启用的并网点
# 无需手动操作！
```

## 🚀 快速开始

### 1. 配置并网点

在数据库中确保并网点配置正确：

```sql
-- 查看当前配置
SELECT number, name, agc_is_enabled, avc_is_enabled 
FROM agvc_bwd_setting;

-- 启用AGC和AVC
UPDATE agvc_bwd_setting 
SET agc_is_enabled = 1, avc_is_enabled = 1 
WHERE number = '1';
```

### 2. 启动系统

```bash
# 直接启动服务
./server

# 或使用systemctl
systemctl restart gin-vue-admin
```

### 3. 查看日志

```bash
# 查看自动启动日志
tail -f logs/server.log | grep "自动启动"

# 期望看到：
# [INFO] ========== 开始自动启动所有并网点的AGC和AVC功能 ==========
# [INFO] 自动启动AGC成功 {"bwdNo": 1, "name": "并网点1"}
# [INFO] 自动启动AVC成功 {"bwdNo": 1, "name": "并网点1"}
# [INFO] ========== 自动启动流程完成 ==========
```

## 📁 修改的文件

- `server/core/server.go` - 添加自动启动调用
- `server/service/agvc/agcv_main/agc.go` - 添加`AutoStartAllGridPoints()`方法
- `server/service/agvc/agcv_main/avc.go` - 添加`AutoStartAllGridPoints()`方法

## 📖 详细文档

- [AGC_AVC_AUTO_START.md](./AGC_AVC_AUTO_START.md) - 完整功能说明和使用指南
- [CHANGELOG_AUTO_START.md](./CHANGELOG_AUTO_START.md) - 变更日志和部署指南
- [IMPLEMENTATION_SUMMARY.md](./IMPLEMENTATION_SUMMARY.md) - 技术实现总结

## ✅ 启动条件

系统会自动启动满足以下条件的并网点：

**AGC启动条件：**
- ✓ 配置存在于`agvc_bwd_setting`表
- ✓ `number`字段不为空
- ✓ `agc_is_enabled = 1`

**AVC启动条件：**
- ✓ 配置存在于`agvc_bwd_setting`表
- ✓ `number`字段不为空
- ✓ `avc_is_enabled = 1`

## 🔍 常见问题

### Q: 某个并网点没有自动启动？

**A:** 检查以下几点：
1. 数据库中`agc_is_enabled`或`avc_is_enabled`是否为1
2. `number`字段是否正确填写
3. 查看日志中的错误信息

### Q: 如何禁用某个并网点的自动启动？

**A:** 设置相应的启用标志为0：
```sql
UPDATE agvc_bwd_setting 
SET agc_is_enabled = 0, avc_is_enabled = 0 
WHERE number = '2';
```
然后重启系统。

### Q: 可以手动启动/停止吗？

**A:** 可以！原有的API接口保持不变：
```bash
# 手动启动
curl -X POST "http://localhost:8888/agvc/agc/start?bwdNo=1"

# 手动停止
curl -X POST "http://localhost:8888/agvc/agc/stop?bwdNo=1"
```

### Q: 系统启动后多久开始自动启动？

**A:** 延迟3秒后开始，确保依赖服务完全初始化。

## 🎯 主要特性

- ✅ **智能过滤**：自动跳过未启用或配置错误的并网点
- ✅ **错误隔离**：单个失败不影响其他并网点
- ✅ **详细日志**：记录每个并网点的启动状态
- ✅ **完全兼容**：保持原有API接口不变
- ✅ **异步执行**：不阻塞系统启动流程

## 📊 日志示例

```
2024-11-04 08:30:15 [INFO] AGC控制服务初始化成功
2024-11-04 08:30:15 [INFO] AVC控制服务初始化成功
2024-11-04 08:30:18 [INFO] ========== 开始自动启动所有并网点的AGC和AVC功能 ==========
2024-11-04 08:30:18 [INFO] 开始自动启动所有并网点的AGC功能 {"并网点数量": 3}
2024-11-04 08:30:18 [INFO] 自动启动AGC成功 {"bwdNo": 1, "name": "并网点1"}
2024-11-04 08:30:18 [DEBUG] 并网点AGC未启用，跳过 {"bwdNo": 2, "name": "并网点2"}
2024-11-04 08:30:18 [INFO] 自动启动AGC成功 {"bwdNo": 3, "name": "并网点3"}
2024-11-04 08:30:18 [INFO] AGC自动启动完成 {"成功数量": 2, "总数量": 3}
2024-11-04 08:30:18 [INFO] 开始自动启动所有并网点的AVC功能 {"并网点数量": 3}
2024-11-04 08:30:18 [INFO] 自动启动AVC成功 {"bwdNo": 1, "name": "并网点1"}
2024-11-04 08:30:18 [INFO] 自动启动AVC成功 {"bwdNo": 3, "name": "并网点3"}
2024-11-04 08:30:18 [INFO] AVC自动启动完成 {"成功数量": 2, "总数量": 3}
2024-11-04 08:30:18 [INFO] ========== 自动启动流程完成 ==========
```

## 🔧 配置示例

### 启用所有功能的并网点

```sql
INSERT INTO agvc_bwd_setting (
    number, 
    name, 
    agc_is_enabled, 
    avc_is_enabled,
    agc_control_period,
    avc_control_period,
    agc_vibration_range
) VALUES (
    '1',
    '并网点1',
    1,  -- 启用AGC
    1,  -- 启用AVC
    30, -- AGC控制周期30秒
    60, -- AVC控制周期60秒
    50  -- AGC抖动区间50kW
);
```

### 仅启用AGC的并网点

```sql
INSERT INTO agvc_bwd_setting (
    number, 
    name, 
    agc_is_enabled, 
    avc_is_enabled,
    agc_control_period
) VALUES (
    '2',
    '并网点2',
    1,  -- 启用AGC
    0,  -- 禁用AVC
    30  -- AGC控制周期30秒
);
```

## 💡 小贴士

1. **首次使用**：建议先在测试环境验证功能正常
2. **日志监控**：建议设置日志监控，及时发现启动失败
3. **配置验证**：启动前确保并网点配置完整
4. **性能考虑**：并网点数量多时，启动可能需要稍长时间

## 🆘 获取帮助

- 查看完整文档：[AGC_AVC_AUTO_START.md](./AGC_AVC_AUTO_START.md)
- 查看变更日志：[CHANGELOG_AUTO_START.md](./CHANGELOG_AUTO_START.md)
- 查看实现细节：[IMPLEMENTATION_SUMMARY.md](./IMPLEMENTATION_SUMMARY.md)

---

**版本**：v1.0.0  
**更新日期**：2024-11-04  
**状态**：✅ Ready for Production
