# AGC/AVC 启动检查和常量化点名称实现总结

## 实现的功能

### 1. PointMapper 点名称常量化

**问题**: PointMapper.GetPointID 中使用硬编码的中文点名称，容易输入错误。

**解决方案**: 在 `server/service/agvc/cons/def.go` 中定义了所有标准点名称常量。

#### 新增常量

##### AGC调度标准点名称常量
- **遥信(YX)**:
  - `AGC_YX_SIGNAL` - AGC投退信号 (点标识: 401)
  - `AGC_YX_CONTROL_MODE` - AGC就地远方控制模式 (点标识: 402)
  - `AGC_YX_LOOP_STATUS` - AGC开/闭环状态 (点标识: 404)
  - `AGC_YX_UP_REG_LOCK` - AGC有功上调节闭锁 (点标识: 405)
  - `AGC_YX_DOWN_REG_LOCK` - AGC有功下调节闭锁 (点标识: 406)

- **遥测(YC)**:
  - `AGC_YC_POWER_UPPER_LIMIT` - 有功调节上限 (点标识: 401)
  - `AGC_YC_POWER_LOWER_LIMIT` - 有功调节下限 (点标识: 402)
  - `AGC_YC_POWER_EXEC_VALUE` - 有功执行值 (点标识: 403)

- **遥控(YK)**:
  - `AGC_YK_SIGNAL` - AGC投退信号 (点标识: 401)
  - `AGC_YK_CONTROL_MODE` - AGC就地远方控制模式 (点标识: 402)
  - `AGC_YK_LOOP_STATUS` - AGC开/闭环状态 (点标识: 404)

- **遥调(YT)**:
  - `AGC_YT_POWER_EXEC` - 有功执行 (点标识: 401)

##### AVC调度标准点名称常量
- **遥信(YX)**:
  - `AVC_YX_SIGNAL` - AVC功能投退信号 (点标识: 401)
  - `AVC_YX_CONTROL_MODE` - AVC功能就地远方控制模式 (点标识: 402)
  - `AVC_YX_COMMAND_STATUS` - AVC功能当前指令状态 (点标识: 403)
  - `AVC_YX_LOOP_STATUS` - AVC功能开闭环状态 (点标识: 404)
  - `AVC_YX_UP_REG_LOCK` - AVC功能上调节闭锁 (点标识: 405)
  - `AVC_YX_DOWN_REG_LOCK` - AVC功能下调节闭锁 (点标识: 406)

- **遥测(YC)**:
  - `AVC_YC_REACTIVE_INC_CAP` - 无功可增容量 (点标识: 401)
  - `AVC_YC_REACTIVE_DEC_CAP` - 无功可减容量 (点标识: 402)
  - `AVC_YC_VOLTAGE_EXEC_VALUE` - 电压执行值 (点标识: 403)
  - `AVC_YC_REACTIVE_EXEC_VALUE` - 无功执行值 (点标识: 404)

- **遥控(YK)**:
  - `AVC_YK_SIGNAL` - AVC功能投退信号 (点标识: 401)
  - `AVC_YK_CONTROL_MODE` - AVC功能就地远方控制模式 (点标识: 402)
  - `AVC_YK_COMMAND_STATUS` - AVC功能当前指令状态 (点标识: 403)
  - `AVC_YK_LOOP_STATUS` - AVC功能开闭环状态 (点标识: 404)

- **遥调(YT)**:
  - `AVC_YT_VOLTAGE_EXEC` - 电压执行 (点标识: 401)
  - `AVC_YT_REACTIVE_EXEC` - 无功执行 (点标识: 402)

##### 并网柜常用点名称常量
- `BWG_YC_ACTIVE_POWER` - 有功功率P(kW)
- `BWG_YC_REACTIVE_POWER` - 无功功率Q(kvar)
- `BWG_YC_VOLTAGE_AB` - AB线电压Uab(kV)
- `BWG_YC_FREQUENCY` - F(频率Hz)
- `BWG_YC_APPARENT_POWER` - S(视在功率kVA)

##### 逆变器常用点名称常量
- `NBQ_YC_ACTIVE_POWER` - 交流功率(kW)
- `NBQ_YC_REACTIVE_POWER` - 无功功率(kVar)
- `NBQ_YC_APPARENT_POWER` - 视在功率(kVa)

### 2. AGC/AVC 调控启动前的启用检查

**问题**: AGC和AVC调控启动前需要检查是否启用，并根据控制模式（本地/远程）从不同数据源读取值。

**解决方案**: 

#### AGC检查逻辑 (agc.go)

在 `executeAGCCycle` 函数中添加了完整的启动前检查：

1. **步骤2: 判断控制权限模式**
   - `ControlAuth = 0 或 nil`: 本地控制模式
   - `ControlAuth = 1`: 远程调度控制模式

2. **步骤3: 检查AGC投退信号**
   - **本地模式**: 从数据库 `AgvcBwdSetting.AgcIsEnabled` 读取
   - **远程模式**: 从调度内存 `DispatchStorage` 读取 AGC投退信号(YX/401)
   - 如果未投入，直接返回，不执行调控

3. **步骤4: 检查AGC就地远方控制模式**
   - **远程模式**: 从调度内存读取 AGC就地远方控制模式(YX/402)
   - 如果不是远程模式(值!=1)，跳过调控

4. **步骤5: 检查AGC开/闭环状态**
   - **本地模式**: 从数据库 `AgvcBwdSetting.RunMode` 读取
   - **远程模式**: 从调度内存读取 AGC开/闭环状态(YX/404)

5. **步骤6: 获取执行值**
   - **本地模式**: 从数据库 `AgvcBwdSetting.StationExecValue` 读取
   - **远程模式**: 从调度内存读取 有功执行值(YC/403)
     - 如果调度内存无数据，则使用数据库的 `DispatchExecValue`

6. **步骤7: 判断是否开环运行**
   - 如果是开环，直接执行目标值
   - 如果是闭环，进入PID调节逻辑

#### AVC检查逻辑 (avc.go)

在 `executeAVCCycle` 函数中添加了类似的启动前检查：

1. **步骤2: 判断控制权限模式**
   - `ControlAuth = 0 或 nil`: 本地控制模式
   - `ControlAuth = 1`: 远程调度控制模式

2. **步骤3: 检查AVC投退信号**
   - **本地模式**: 从数据库 `AVCConfig.IsActive` 读取
   - **远程模式**: 从调度内存 `DispatchStorage` 读取 AVC功能投退信号(YX/401)
   - 如果未投入，直接返回，不执行调控

3. **步骤4: 检查AVC就地远方控制模式**
   - **远程模式**: 从调度内存读取 AVC功能就地远方控制模式(YX/402)
   - 如果不是远程模式(值!=1)，跳过调控

4. **步骤5: 检查AVC开/闭环状态**
   - **本地模式**: 从数据库 `AVCConfig.RunMode` 读取
   - **远程模式**: 从调度内存读取 AVC功能开闭环状态(YX/404)
   - (目前AVC暂不支持开环模式，变量仅用于日志记录)

### 3. 更新所有使用硬编码点名称的地方

#### 修改的文件:
1. `server/service/agvc/agcv_main/agc.go`
   - 使用 `cons.AGC_YX_SIGNAL` 替代 "AGC投退信号"
   - 使用 `cons.AGC_YX_CONTROL_MODE` 替代 "AGC就地远方控制模式"
   - 使用 `cons.AGC_YX_LOOP_STATUS` 替代 "AGC开/闭环状态"
   - 使用 `cons.AGC_YC_POWER_EXEC_VALUE` 替代 "有功执行值"
   - 使用 `cons.BWG_YC_FREQUENCY` 替代 "F(频率Hz)"

2. `server/service/agvc/agcv_main/avc.go`
   - 使用 `cons.AVC_YX_SIGNAL` 替代 "AVC功能投退信号"
   - 使用 `cons.AVC_YX_CONTROL_MODE` 替代 "AVC功能就地远方控制模式"
   - 使用 `cons.AVC_YX_LOOP_STATUS` 替代 "AVC功能开闭环状态"
   - 使用 `cons.BWG_YC_VOLTAGE_AB` 替代 "AB线电压Uab(kV)"
   - 使用 `cons.BWG_YC_FREQUENCY` 替代 "F(频率Hz)"

3. `server/service/agvc/agcv_main/power_aggregator.go`
   - 使用 `cons.BWG_YC_ACTIVE_POWER` 替代 "有功功率P(kW)"
   - 使用 `cons.BWG_YC_REACTIVE_POWER` 替代 "无功功率Q(kvar)"
   - 使用 `cons.BWG_YC_APPARENT_POWER` 替代 "S(视在功率kVA)"
   - 使用 `cons.NBQ_YC_ACTIVE_POWER` 替代 "交流功率(kW)"
   - 使用 `cons.NBQ_YC_REACTIVE_POWER` 替代 "无功功率(kVar)"
   - 使用 `cons.NBQ_YC_APPARENT_POWER` 替代 "视在功率(kVa)"

## 优势

1. **减少错误**: 使用常量避免拼写错误和输入错误
2. **代码可维护性**: 点名称统一管理，便于修改和维护
3. **智能提示**: IDE可以提供自动补全和类型检查
4. **符合规范**: 符合调度AGC/AVC数据标准点规范
5. **灵活的控制模式**: 支持本地和远程两种控制模式
6. **数据源分离**: 本地模式读数据库，远程模式读内存，职责清晰

## 测试验证

- ✅ 代码编译通过
- ✅ 所有硬编码点名称已替换为常量
- ✅ AGC/AVC启动检查逻辑已实现
- ✅ 支持本地/远程两种控制模式

## 注意事项

1. 调度数据存储在 `DispatchStorage` 内存中，由1187端口接收的CoAP数据填充
2. 数据库配置作为fallback，当调度数据不可用时使用
3. 所有点标识需要在Excel文件 `/home/engine/project/server/设备及测点标准.xlsx` 中正确配置
4. 如果点位映射失败，会使用默认的点标识作为后备方案

## 前端说明

根据用户要求，本次修改**仅涉及后端代码**，前端无需修改，也无需重新编译前端。
